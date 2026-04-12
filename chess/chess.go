package chess

import (
	"fmt"
	"runtime"
	"slices"
	"strings"

	"github.com/RchrdHndrcks/gochess/v2"
)

type (
	// Cloner is a generic interface for objects that can be cloned.
	Cloner interface {
		// Clone returns a clone of the object.
		Clone() Board
	}

	// Board represents a chess board.
	Board interface {
		// SetSquare sets a piece in a square.
		SetSquare(c gochess.Coordinate, p gochess.Piece) error
		// Square returns the piece in a square.
		Square(c gochess.Coordinate) (gochess.Piece, error)
		// Width returns the width of the board.
		Width() int
	}

	// config represents configurations of how the methods will work.
	config struct {
		// Parallelism is the number of workers to use for the moves calculation.
		Parallelism int
	}

	// chessContext represents the history of a game.
	chessContext struct {
		// move is the played move in UCI notation. Kept for unmakeMove which
		// parses the string to recover origin/target/promotion data.
		move string
		// compactMove is the same move in packed form. Set when MakeMoveCompact
		// is the entry point (NullMove for the legacy makeMove(uci) path, since
		// the engine uses move for unmake either way).
		compactMove Move
		// capturedPiece is the piece captured by this move (Empty if none).
		// For en passant, this is the captured pawn (with color).
		capturedPiece gochess.Piece
		// positionKey is the first three FEN fields (placement / active color /
		// castling rights) of the position before the move was made. Used for
		// threefold-repetition detection without re-parsing FEN strings.
		positionKey string
		// halfMove is the number of half moves since the last capture or pawn move.
		halfMove int
		// availableCastles is the castles that are available.
		availableCastles CastleRights
		// enPassantFile is the file (0-7 = a-h) where en passant capture is
		// available, or -1 if none.
		enPassantFile int8
		// whiteKingPosition is the position of the white king.
		whiteKingPosition *gochess.Coordinate
		// blackKingPosition is the position of the black king.
		blackKingPosition *gochess.Coordinate
		// check is true if the current turn is in check.
		check bool
		// checkmate is true if the current turn is in checkmate.
		checkmate bool
		// stalemate is true if the current turn is in stalemate.
		stalemate bool
	}

	// Chess represents a Chess game.
	//
	// A Chess value is not safe for concurrent use by multiple goroutines.
	Chess struct {
		board Board
		// turn is the current turn.
		turn gochess.Piece
		// movesCount is the number of moves played in algebaric notation.
		// It will increase by 1 after each black move.
		movesCount uint64
		// halfMoves is the number of half moves since the last capture or pawn move.
		halfMoves int
		// enPassantFile is the file (0-7 = a-h) where en passant capture is
		// available, or -1 if none.
		enPassantFile int8
		// availableCastles is the castles that are available.
		availableCastles CastleRights
		// moves are the available moves in the current position.
		moves []string
		// actualFEN is the FEN string of the current position.
		actualFEN string
		// blackKingPosition is the position of the black king.
		blackKingPosition *gochess.Coordinate
		// whiteKingPosition is the position of the white king.
		whiteKingPosition *gochess.Coordinate
		// check is true if the current turn is in check.
		check bool
		// checkmate is true if the current turn is in checkmate.
		checkmate bool
		// stalemate is true if the current turn is in stalemate.
		stalemate bool

		// config represents configurations of how the methods will work.
		config config

		// history is the history of the game.
		history []chessContext

		// pieceLists tracks piece locations per color and type.
		// First index: 0=White, 1=Black. Second index: piece type (1–6).
		pieceLists [2][7]pieceList
	}
)

var (
	castlesMoves = map[string]gochess.Piece{
		"e1g1": gochess.White,
		"e1c1": gochess.White,
		"e8g8": gochess.Black,
		"e8c8": gochess.Black,
	}

	castleKingWay = map[string]gochess.Coordinate{
		"e1g1": gochess.Coor(5, 7),
		"e1c1": gochess.Coor(3, 7),
		"e8g8": gochess.Coor(5, 0),
		"e8c8": gochess.Coor(3, 0),
	}

	castleRook = map[string]gochess.Coordinate{
		"e1g1": gochess.Coor(7, 7),
		"e1c1": gochess.Coor(0, 7),
		"e8g8": gochess.Coor(7, 0),
		"e8c8": gochess.Coor(0, 0),
	}
)

// New creates a new chess game.
//
// The chess.AvailableMoves method will use a pool of workers to maximize
// performance on the moves calculation.
// By default, the number of workers is twice the number of available CPUs.
// If you are running on a container environment or you want to use the
// sequential version, you should set this value manually using the
// WithParallelism option.
// The pool of workers will be used only if the board implements the Cloner
// interface. If you are using a custom Board with the WithBoard option, you
// should implement the Cloner interface to take advantage of the parallelism.
func New(opts ...Option) (*Chess, error) {
	c := &Chess{
		board:            newBoardAdapter(gochess.DefaultChessBoard()),
		turn:             gochess.White,
		movesCount:       1,
		halfMoves:        0,
		enPassantFile:    -1,
		availableCastles: AllCastling,
		actualFEN:        "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		moves: []string{
			"a2a3", "a2a4", "b2b3", "b2b4", "c2c3", "c2c4", "d2d3", "d2d4",
			"e2e3", "e2e4", "f2f3", "f2f4", "g2g3", "g2g4", "h2h3", "h2h4",
			"b1a3", "b1c3", "g1f3", "g1h3",
		},
		blackKingPosition: &gochess.Coordinate{X: 4, Y: 0},
		whiteKingPosition: &gochess.Coordinate{X: 4, Y: 7},
		check:             false,
		checkmate:         false,
		stalemate:         false,
		config: config{
			// To maximize performance chess uses twice the number of available
			// CPUs. If you are running on a container environment or you want to
			// use the sequential version, you should set this value manually.
			Parallelism: runtime.NumCPU() * 2,
		},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	c.initPieceLists()

	return c, nil
}

// LoadPosition loads a board from a FEN string.
//
// The function will read the entire FEN string and will return an error if
// the FEN string is invalid.
//
// The board and properties will not be modified if the FEN string is invalid.
func (c *Chess) LoadPosition(FEN string) error {
	if err := c.loadPosition(FEN); err != nil {
		return err
	}

	c.actualFEN = FEN
	c.initPieceLists()
	c.moves = c.legalMoves()
	check := c.isCheck()
	c.check = check && len(c.moves) > 0
	c.checkmate = check && len(c.moves) == 0
	c.stalemate = !check && len(c.moves) == 0
	return nil
}

// Turn returns the current turn.
//
// It will be gochess.White or gochess.Black.
func (c *Chess) Turn() gochess.Piece {
	return c.turn
}

// FEN returns the FEN string of the current position.
//
// If any of the kings is not in the board, the function returns an empty string.
func (c *Chess) FEN() string {
	return c.actualFEN
}

// AvailableMoves returns the available legal moves for the current turn.
//
// It always returns a non nil slice. It could be empty if the position is
// checkmate or stalemate.
func (c *Chess) AvailableMoves() []string {
	return slices.Clone(c.moves)
}

// MakeMove checks if the move is legal and makes it.
// It returns an error if the move is not legal.
func (c *Chess) MakeMove(move string) error {
	if !slices.Contains(c.moves, move) {
		return fmt.Errorf("move is not legal: %s", move)
	}

	c.makeMove(move)
	c.actualFEN = c.calculateFEN(move)
	c.moves = c.legalMoves()
	check := c.isCheck()
	c.check = check && len(c.moves) > 0
	c.checkmate = check && len(c.moves) == 0
	c.stalemate = !check && len(c.moves) == 0
	return nil
}

// UnmakeMove unmake the last move.
//
// It searches for the last move in the history and unmake it.
// If there are no moves in the history, the function does nothing.
func (c *Chess) UnmakeMove() {
	c.unmakeMove()
	c.moves = c.legalMoves()
}

// IsCheck returns if the current turn is in check.
func (c *Chess) IsCheck() bool {
	return c.check
}

// IsCheckmate returns if the current turn is in checkmate.
func (c *Chess) IsCheckmate() bool {
	return c.checkmate
}

// IsStalemate returns if the current turn is in stalemate.
func (c *Chess) IsStalemate() bool {
	return c.stalemate
}

// IsFiftyMoveRule returns true if the fifty-move rule draw condition is met.
//
// The fifty-move rule states that a player can claim a draw if no capture
// has been made and no pawn has been moved in the last fifty moves
// (100 half-moves).
func (c *Chess) IsFiftyMoveRule() bool {
	return c.halfMoves >= 100
}

// IsInsufficientMaterial returns true if neither side has sufficient material
// to checkmate the opponent, as defined by FIDE Laws of Chess article 5.2.2.
//
// The following positions are considered insufficient material:
//   - King vs King: bare kings cannot deliver checkmate by any legal sequence.
//   - King + Knight vs King: a lone knight cannot force checkmate without the
//     opponent's cooperation.
//   - King + Bishop vs King: a bishop controls only one square color, so the
//     defending king can always evade.
//   - King + Bishop vs King + Bishop (same color squares): the attacking bishop
//     can never reach the squares the defending bishop occupies, so no forced
//     mate exists. Opposite-color bishops are NOT insufficient — they can
//     cooperate to deliver checkmate.
//
// Any other material (pawn, rook, queen, two or more knights, or mixed minor
// pieces not listed above) is considered sufficient.
func (c *Chess) IsInsufficientMaterial() bool {
	width := c.board.Width()

	var knights, bishops int
	var bishopSquareColor int // stores (x+y)%2 of the first bishop found
	bishopSquareColor = -1
	allBishopsSameColor := true

	for y := range width {
		for x := range width {
			piece, _ := c.board.Square(gochess.Coor(x, y))
			if piece == gochess.Empty {
				continue
			}

			pieceType := gochess.PieceType(piece)
			switch pieceType {
			case gochess.King:
				// Kings are always present, skip them.
				continue
			case gochess.Knight:
				knights++
			case gochess.Bishop:
				bishops++
				sc := (x + y) % 2
				if bishopSquareColor == -1 {
					bishopSquareColor = sc
				} else if sc != bishopSquareColor {
					allBishopsSameColor = false
				}
			default:
				// Any other piece (pawn, rook, queen) means sufficient material.
				return false
			}
		}
	}

	totalMinor := knights + bishops

	// King vs King
	if totalMinor == 0 {
		return true
	}

	// King + single minor piece vs King
	if totalMinor == 1 {
		return true
	}

	// King + Bishop vs King + Bishop on same color squares
	if knights == 0 && bishops == 2 && allBishopsSameColor {
		return true
	}

	return false
}

// IsThreefoldRepetition returns true if the current position has occurred
// at least three times during the game.
//
// Two positions are considered the same when piece placement, active color,
// and castling availability are identical. The en passant field and move
// counters are excluded: an en passant target square is set for exactly one
// ply after a double pawn push and the same pawn can never double-push again,
// so the identical non-"-" en passant square cannot appear in two distinct
// positions within a single game.
func (c *Chess) IsThreefoldRepetition() bool {
	currentKey := positionKey(c.actualFEN)
	count := 1 // current position counts as one occurrence
	for _, ctx := range c.history {
		if ctx.positionKey == currentKey {
			count++
			if count >= 3 {
				return true
			}
		}
	}
	return false
}

// positionKey extracts the first three fields (piece placement, active color,
// castling rights) from a FEN string for use in repetition detection.
//
// The en passant field is intentionally excluded: a given en passant target
// square is set for exactly one ply and can never recur within the same game,
// so including it would only produce false negatives without preventing false
// positives.
func positionKey(fen string) string {
	parts := strings.SplitN(fen, " ", 5)
	if len(parts) < 3 {
		return fen
	}
	return parts[0] + " " + parts[1] + " " + parts[2]
}

// Square returns the piece in a square.
// The square is represented by an algebraic notation.
//
// If the square is not valid, the function returns an error.
func (c *Chess) Square(square string) (string, error) {
	coor, err := AlgebraicToCoordinate(square)
	if err != nil {
		return "", fmt.Errorf("failed to convert algebraic notation to coordinate: %w", err)
	}

	p, err := c.board.Square(coor)
	if err != nil {
		return "", fmt.Errorf("failed to get piece from board at %s: %w", square, err)
	}

	return gochess.PieceNames[p], nil
}

// MakeMoveCompact applies the given compact Move to the board.
//
// The move is validated against the legal move list (by UCI string) and
// rejected with an error if it is not legal in the current position. On
// success, the FEN, legal-move list and check/checkmate/stalemate flags
// are recomputed.
//
// Unlike MakeMove(string), the board mutation does not parse strings:
// the from/to/flags/promotion are read directly from the Move bits and
// passed to applyMove. The only string allocated is m.UCI(), which is
// stored in history because unmakeMove parses it.
func (c *Chess) MakeMoveCompact(m Move) error {
	uci := m.UCI()
	if !slices.Contains(c.moves, uci) {
		return fmt.Errorf("move is not legal: %s", uci)
	}

	from := m.From()
	to := m.To()

	md := moveData{
		from:          from,
		to:            to,
		isCastle:      c.isCastleMove(uci),
		isEnPassant:   c.isEnPassantMove(uci),
		promotionType: m.Promotion(),
		uci:           uci,
		compact:       m,
	}

	if md.isEnPassant {
		md.capturedPiece = gochess.Pawn | opponentColor(c.turn)
	} else {
		dst, _ := c.board.Square(to)
		md.capturedPiece = dst
	}

	c.applyMove(md)
	c.actualFEN = c.calculateFEN(uci)
	c.moves = c.legalMoves()
	check := c.isCheck()
	c.check = check && len(c.moves) > 0
	c.checkmate = check && len(c.moves) == 0
	c.stalemate = !check && len(c.moves) == 0
	return nil
}

// UnmakeMoveCompact unmakes the last move and recomputes the legal-move list.
func (c *Chess) UnmakeMoveCompact() {
	c.unmakeMove()
	c.moves = c.legalMoves()
}

// ParseUCIMove validates the given UCI string against the current legal
// move list and, if legal, returns its compact Move representation. The
// returned Move has from/to/flags/promotion/captured filled in but does
// not set the GivesCheck bit.
func (c *Chess) ParseUCIMove(uci string) (Move, error) {
	if len(uci) != 4 && len(uci) != 5 {
		return NullMove, fmt.Errorf("invalid UCI move: %q", uci)
	}
	if !slices.Contains(c.moves, uci) {
		return NullMove, fmt.Errorf("move is not legal: %s", uci)
	}
	return c.uciToCompactMove(uci), nil
}

// Moves returns all legal moves for the current position as compact Move
// values. The GivesCheck bit is computed for each move by playing it on the
// internal board, calling isCheck, and unmaking it. The lightweight
// applyMove/unmakeMove path is used (no FEN recompute, no legal-move-list
// refresh) so this stays O(n) over the legal moves.
func (c *Chess) Moves() []Move {
	uciMoves := c.moves
	out := make([]Move, 0, len(uciMoves))
	for _, uci := range uciMoves {
		m := c.uciToCompactMove(uci)
		md := moveData{
			from:          m.From(),
			to:            m.To(),
			isCastle:      c.isCastleMove(uci),
			isEnPassant:   c.isEnPassantMove(uci),
			promotionType: m.Promotion(),
			uci:           uci,
			compact:       m,
		}
		if md.isEnPassant {
			md.capturedPiece = gochess.Pawn | opponentColor(c.turn)
		} else {
			dst, _ := c.board.Square(md.to)
			md.capturedPiece = dst
		}
		c.applyMove(md)
		if c.isCheck() {
			m = m.WithGivesCheck(true)
		}
		c.unmakeMove()
		out = append(out, m)
	}
	return out
}

// PieceAt returns the piece type, color, and ok=true if a piece exists at
// the given 0-63 square index. Returns zero values and ok=false if the
// square is empty or out of range.
func (c *Chess) PieceAt(sq int) (pieceType, color gochess.Piece, ok bool) {
	if sq < 0 || sq >= 64 {
		return 0, 0, false
	}
	coor := coordinateFromSquare(uint32(sq))
	p, err := c.board.Square(coor)
	if err != nil || p == gochess.Empty {
		return 0, 0, false
	}
	return gochess.PieceType(p), gochess.PieceColor(p), true
}

// uciToCompactMove inspects current board state to determine flags,
// captured piece, and promotion for the given UCI string.
func (c *Chess) uciToCompactMove(uci string) Move {
	from, _ := AlgebraicToCoordinate(uci[:2])
	to, _ := AlgebraicToCoordinate(uci[2:4])

	isCastle := c.isCastleMove(uci)
	isEP := c.isEnPassantMove(uci)

	var captured gochess.Piece
	if isEP {
		captured = gochess.Pawn
	} else {
		dst, _ := c.board.Square(to)
		captured = gochess.PieceType(dst)
	}

	var promo gochess.Piece
	if len(uci) == 5 {
		promo = gochess.PiecesWithoutColor[uci[4:5]]
	}

	// Detect double pawn push.
	isDoublePush := false
	if !isCastle && !isEP && promo == gochess.Empty {
		fp, _ := c.board.Square(from)
		if gochess.PieceType(fp) == gochess.Pawn {
			dy := to.Y - from.Y
			if dy == 2 || dy == -2 {
				isDoublePush = true
			}
		}
	}

	flag := FlagQuiet
	switch {
	case isCastle:
		flag = FlagCastle
	case isEP:
		flag = FlagEnPassant
	case promo != gochess.Empty && captured != gochess.Empty:
		flag = FlagPromotionCapture
	case promo != gochess.Empty:
		flag = FlagPromotion
	case captured != gochess.Empty:
		flag = FlagCapture
	case isDoublePush:
		flag = FlagDoublePush
	}

	m := NewMove(from, to, flag)
	if promo != gochess.Empty {
		m = m.WithPromotion(promo)
	}
	if captured != gochess.Empty {
		m = m.WithCapturedPiece(captured)
	}
	return m
}

// clone creates a deep copy of the Chess structure.
func (c Chess) clone() Chess {
	cloned := c

	cloner, ok := c.board.(Cloner)
	if ok {
		cloned.board = cloner.Clone()
	}

	if c.whiteKingPosition != nil {
		whitePos := *c.whiteKingPosition
		cloned.whiteKingPosition = &whitePos
	}
	if c.blackKingPosition != nil {
		blackPos := *c.blackKingPosition
		cloned.blackKingPosition = &blackPos
	}

	if len(c.history) > 0 {
		cloned.history = make([]chessContext, len(c.history))
		for i, ctx := range c.history {
			cloned.history[i] = ctx
			if ctx.whiteKingPosition != nil {
				whitePos := *ctx.whiteKingPosition
				cloned.history[i].whiteKingPosition = &whitePos
			}
			if ctx.blackKingPosition != nil {
				blackPos := *ctx.blackKingPosition
				cloned.history[i].blackKingPosition = &blackPos
			}
		}
	} else {
		cloned.history = nil
	}

	if len(c.moves) > 0 {
		cloned.moves = make([]string, len(c.moves))
		copy(cloned.moves, c.moves)
	}

	return cloned
}
