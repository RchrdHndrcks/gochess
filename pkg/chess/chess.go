package chess

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/RchrdHndrcks/gochess/pkg"
)

type (
	// Board represents a chess board.
	Board interface {
		// LoadPosition loads a position into the board from a string.
		LoadPosition(string) error
		// Square returns the piece in a square.
		Square(c pkg.Coordinate) (int8, error)
		// AvailableMoves returns all the available moves for the current turn, without checking
		// if they are legal.
		AvailableMoves(turn int8, inPassantSquare, castlePossibilities string) ([]string, error)
		// MakeMove makes a move without checking if it is legal.
		MakeMove(string) error
		// Width returns the width of the board.
		Width() int
	}

	// gameHistory represents the history of a game.
	gameHistory struct {
		// move is a played move.
		move string
		// fen is a FEN strings that represents the position after the move.
		fen string
		// halfMove is the number of half moves since the last capture or pawn move.
		halfMove int
		// availableCastles is the castles that are available.
		availableCastles string
		// inPassantSquare is the square where a pawn can capture in passant.
		inPassantSquare string
	}

	// Chess represents a Chess game.
	Chess struct {
		board Board
		// turn is the current turn.
		turn int8
		// movesCount is the number of moves played in algebaric notation.
		// It will increase by 1 after each black move.
		movesCount uint32
		// halfMoves is the number of half moves since the last capture or pawn move.
		halfMoves int
		// inPassantSquare is the square where a pawn can capture in passant.
		inPassantSquare string
		// availableCastles is the castles that are available.
		// It will has the same format as the FEN castles.
		availableCastles string

		history []gameHistory
	}
)

var (
	castlesMoves = map[string]int8{
		"e1g1": pkg.White,
		"e1c1": pkg.White,
		"e8g8": pkg.Black,
		"e8c8": pkg.Black,
	}

	castleKingWay = map[string]pkg.Coordinate{
		"e1g1": pkg.Coor(5, 7),
		"e1c1": pkg.Coor(3, 7),
		"e8g8": pkg.Coor(5, 0),
		"e8c8": pkg.Coor(3, 0),
	}
)

// NewChess creates a new chess game.
func NewChess(opts ...Option) (*Chess, error) {
	c := &Chess{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	defaultOptions(c)
	return c, nil
}

// FEN returns the FEN string of the current position.
func (c Chess) FEN() string {
	if c.board == nil || reflect.ValueOf(c.board).IsNil() {
		return ""
	}

	fen := ""
	for y := 0; y < c.board.Width(); y++ {
		empty := 0
		for x := 0; x < c.board.Width(); x++ {
			piece, err := c.board.Square(pkg.Coor(x, y))
			if err != nil {
				return ""
			}

			if piece == pkg.Empty {
				empty++
				continue
			}

			if empty > 0 {
				fen += fmt.Sprintf("%d", empty)
				empty = 0
			}

			fen += pkg.PieceNames[piece]
		}

		if empty > 0 {
			fen += fmt.Sprintf("%d", empty)
		}

		if y < 7 {
			fen += "/"
		}
	}

	turn := c.turn
	if turn == 0 {
		turn = pkg.White
	}

	ac := c.availableCastles
	if ac == "" {
		ac = "-"
	}

	ips := c.inPassantSquare
	if ips == "" {
		ips = "-"
	}

	fen += fmt.Sprintf(" %s %s %s %d %d",
		pkg.ColorNames[turn], ac, ips, c.halfMoves, c.movesCount)
	return fen
}

// AvailableLegalMoves returns the available legal moves for the current turn.
// It returns an empty slice if position is stalemate.
// It returns a nil slice if position is checkmate.
func (c Chess) AvailableLegalMoves() ([]string, error) {
	if c.board == nil || reflect.ValueOf(c.board).IsNil() {
		return nil, fmt.Errorf("board cannot be nil")
	}

	moves, err := c.board.AvailableMoves(c.turn, c.inPassantSquare, c.availableCastles)
	if err != nil {
		return nil, fmt.Errorf("failed to get available moves: %w", err)
	}

	// If any of the kings is not in the board, there are not legal moves.
	if !c.isPositionValid() {
		return nil, fmt.Errorf("poition is not valid")
	}

	legalMoves := []string{}
	for _, move := range moves {
		isLegalMove, err := c.isLegalMove(move)
		if err != nil {
			return nil, fmt.Errorf("failed to check if move is legal: %w", err)
		}

		if isLegalMove {
			legalMoves = append(legalMoves, move)
		}
	}

	// If there are no legal moves, we need to check if position is checkmate or
	// stalemate.
	if len(legalMoves) == 0 {
		isCheck, err := c.isCheck()
		if err != nil {
			return nil, fmt.Errorf("failed to check if position is check: %w", err)
		}

		if isCheck {
			legalMoves = nil
		}
	}

	return legalMoves, nil
}

func (c Chess) isLegalMove(move string) (bool, error) {
	turn := c.turn
	err := c.makeMove(move)
	if err != nil {
		return false, fmt.Errorf("failed to make move: %w", err)
	}

	availableMoves, err := c.board.AvailableMoves(c.turn,
		c.inPassantSquare, c.availableCastles)
	if err != nil {
		return false, fmt.Errorf("failed to get legal moves after check move %s: %w", move, err)
	}

	kingPosition, err := c.kingsPosition(turn)
	if err != nil {
		return false, fmt.Errorf("failed to get king position: %w", err)
	}

	kingUnderAttack := destinationMatch(availableMoves, kingPosition)

	c.unmakeMove() // nolint:errcheck

	kingWayUnderAttack := false
	if c.isCastleMove(move) {
		kingWayUnderAttack = destinationMatch(availableMoves, castleKingWay[move])
	}

	return !kingUnderAttack && !kingWayUnderAttack, nil
}

// MakeMove checks if the move is legal and makes it.
// It returns an error if the move is not legal.
func (c *Chess) MakeMove(move string) error {
	if c.board == nil || reflect.ValueOf(c.board).IsNil() {
		return fmt.Errorf("board cannot be nil")
	}

	moves, err := c.AvailableLegalMoves()
	if err != nil {
		return fmt.Errorf("failed to get available moves: %w", err)
	}

	if !slices.Contains(moves, move) {
		return fmt.Errorf("move is not legal: %s", move)
	}

	return c.makeMove(move)
}

// makeMove makes a move without checking if it is legal.
func (c *Chess) makeMove(move string) error {
	lastFEN := c.FEN()

	err := c.board.MakeMove(move)
	if err != nil {
		return fmt.Errorf("failed to make move: %w", err)
	}

	c.history = append(
		c.history,
		gameHistory{
			move:             move,
			fen:              lastFEN,
			halfMove:         c.halfMoves,
			availableCastles: c.availableCastles,
		},
	)

	c.toggleColor()
	c.updateMovesCount()
	c.updateCastlePossibilities()
	c.updateHalfMoves()
	c.updateInPassantSquare()
	return nil
}

// UnmakeMove unmake the last move.
func (c *Chess) UnmakeMove() error {
	return c.unmakeMove()
}

func (c *Chess) unmakeMove() error {
	if len(c.history) == 0 {
		return nil
	}

	lastMove := c.history[len(c.history)-1]
	c.history = c.history[:len(c.history)-1]

	lastFEN := lastMove.fen
	err := c.board.LoadPosition(lastFEN)
	if err != nil {
		return fmt.Errorf("failed to load position: %w", err)
	}

	c.halfMoves = lastMove.halfMove
	c.availableCastles = lastMove.availableCastles
	c.inPassantSquare = lastMove.inPassantSquare

	// If turn color is white, last move was black.
	// So we decrease the moves count.
	if c.turn == pkg.White {
		c.movesCount--
	}

	c.toggleColor()
	return nil
}

// IsCheck returns if the current turn is in check.
func (c *Chess) IsCheck() (bool, error) {
	if c.board == nil || reflect.ValueOf(c.board).IsNil() {
		return false, fmt.Errorf("board cannot be nil")
	}

	return c.isCheck()
}

func (c *Chess) setProperties(FEN string) error {
	props := strings.Split(FEN, " ")
	if len(props) != 6 {
		return fmt.Errorf("invalid FEN: %s", FEN)
	}

	props = props[1:]

	color, ok := pkg.Colors[props[0]]
	if !ok {
		return fmt.Errorf("invalid FEN color: %s", props[0])
	}

	availableCastles := props[1]
	if err := validateCastles(availableCastles); err != nil {
		return fmt.Errorf("invalid FEN castles: %w", err)
	}

	inPassantSquare := props[2]
	if err := validateInPassant(inPassantSquare); err != nil {
		return fmt.Errorf("invalid FEN in passant square: %w", err)
	}

	halfMoves, err := strconv.Atoi(props[3])
	if err != nil {
		return fmt.Errorf("invalid FEN half moves: %s", props[3])
	}

	movesCount, err := strconv.Atoi(props[4])
	if err != nil {
		return fmt.Errorf("invalid FEN moves count: %s", props[4])
	}

	c.turn = color
	c.availableCastles = availableCastles
	c.inPassantSquare = props[2]
	c.halfMoves = halfMoves
	c.movesCount = uint32(movesCount)
	return nil
}

// kingsPosition returns de position of the king of the current turn.
func (c Chess) kingsPosition(color int8) (pkg.Coordinate, error) {
	if color != pkg.White && color != pkg.Black {
		return pkg.Coordinate{}, fmt.Errorf("invalid color: %d", color)
	}

	for y := 0; y < c.board.Width(); y++ {
		for x := 0; x < c.board.Width(); x++ {
			piece, err := c.board.Square(pkg.Coor(x, y))
			if err != nil {
				return pkg.Coordinate{}, fmt.Errorf("failed to get square: %w", err)
			}

			if piece == pkg.King|color {
				return pkg.Coor(x, y), nil
			}
		}
	}

	return pkg.Coordinate{}, fmt.Errorf("king not found")
}

func (c Chess) isCheck() (bool, error) {
	kingPosition, err := c.kingsPosition(c.turn)
	if err != nil {
		return false, fmt.Errorf("failed to get king position: %w", err)
	}

	c.toggleColor()
	moves, err := c.board.AvailableMoves(c.turn, c.inPassantSquare, c.availableCastles)
	if err != nil {
		return false, fmt.Errorf("failed to get available moves: %w", err)
	}

	return destinationMatch(moves, kingPosition), nil
}

func (c *Chess) toggleColor() {
	if c.turn == pkg.White {
		c.turn = pkg.Black
		return
	}

	c.turn = pkg.White
}

func (c *Chess) updateMovesCount() {
	if c.turn == pkg.White {
		c.movesCount++
	}
}

// updateCastlePossibilities checks if the castles are still available.
func (c *Chess) updateCastlePossibilities() {
	toBeRemoved := map[string]bool{}

	k, _ := c.board.Square(pkg.Coor(4, 0))  // nolint:errcheck
	rr, _ := c.board.Square(pkg.Coor(7, 0)) // nolint:errcheck
	lr, _ := c.board.Square(pkg.Coor(0, 0)) // nolint:errcheck
	toBeRemoved["k"] = rr != pkg.Rook|pkg.Black || k != pkg.King|pkg.Black
	toBeRemoved["q"] = lr != pkg.Rook|pkg.Black || k != pkg.King|pkg.Black

	K, _ := c.board.Square(pkg.Coor(4, 7))  // nolint:errcheck
	rR, _ := c.board.Square(pkg.Coor(7, 7)) // nolint:errcheck
	lR, _ := c.board.Square(pkg.Coor(0, 7)) // nolint:errcheck
	toBeRemoved["K"] = rR != pkg.Rook|pkg.White || K != pkg.King|pkg.White
	toBeRemoved["Q"] = lR != pkg.Rook|pkg.White || K != pkg.King|pkg.White

	for castle, mustDelete := range toBeRemoved {
		if !mustDelete {
			continue
		}

		c.availableCastles = strings.ReplaceAll(c.availableCastles, castle, "")
	}
}

// updateHalfMoves updates the half moves counter.
func (c *Chess) updateHalfMoves() {
	c.halfMoves++
	if len(c.history) == 0 {
		return
	}

	// First we look for a change in the board.
	// If we have less pieces than before, a capture was made so we reset the counter.
	h := c.history[len(c.history)-1]
	aux, _ := NewChess(WithFEN(h.fen)) // nolint:errcheck
	piecesCount := 0
	auxPiecesCount := 0
	for y := 0; y < c.board.Width(); y++ {
		for x := 0; x < c.board.Width(); x++ {
			piece, _ := c.board.Square(pkg.Coor(x, y))      // nolint:errcheck
			auxPiece, _ := aux.board.Square(pkg.Coor(x, y)) // nolint:errcheck
			if piece != pkg.Empty {
				piecesCount++
			}

			if auxPiece != pkg.Empty {
				auxPiecesCount++
			}
		}
	}

	if piecesCount != auxPiecesCount {
		c.halfMoves = 0
		return
	}

	// If no capture was made, we check if last move was a pawn move.
	origin := h.move[:2]
	coor, _ := pkg.AlgebraicToCoordinate(origin) // nolint:errcheck

	p, err := aux.board.Square(coor)
	if err != nil {
		return
	}

	piece := p &^ (pkg.White | pkg.Black)
	if piece == pkg.Pawn {
		c.halfMoves = 0
	}
}

func (c *Chess) updateInPassantSquare() {
	c.inPassantSquare = ""

	if c.history == nil || len(c.history) == 0 {
		return
	}

	lastMove := c.history[len(c.history)-1].move
	if len(lastMove) != 4 {
		return
	}

	dest, err := pkg.AlgebraicToCoordinate(lastMove[2:])
	if err != nil {
		return
	}

	p, err := c.board.Square(dest)
	if err != nil {
		return
	}

	if p&^(pkg.White|pkg.Black) != pkg.Pawn {
		return
	}

	destRow, _ := strconv.Atoi(lastMove[3:4]) // nolint:errcheck
	orgRow, _ := strconv.Atoi(lastMove[1:2])  // nolint:errcheck
	if destRow == orgRow+2 || destRow == orgRow-2 {
		c.inPassantSquare = fmt.Sprintf("%s%d", lastMove[2:3], (destRow+orgRow)/2)
		return
	}
}

// isCastleMove returns if the move is a castle move.
func (c Chess) isCastleMove(move string) bool {
	if castlesMoves[move] != c.turn {
		return false
	}

	origin, err := pkg.AlgebraicToCoordinate(move[:2])
	if err != nil {
		return false
	}

	p, err := c.board.Square(origin)
	if err != nil {
		return false
	}

	return p == pkg.King|c.turn
}

// isPositionValid checks if both kings are in the board.
func (c Chess) isPositionValid() bool {
	if _, err := c.kingsPosition(pkg.White); err != nil {
		return false
	}

	if _, err := c.kingsPosition(pkg.Black); err != nil {
		return false
	}

	return true
}

// validateCastles validates the castles string.
func validateCastles(castles string) error {
	if castles == "-" {
		return nil
	}

	castlePieces := map[rune]bool{'K': true, 'Q': true, 'k': true, 'q': true}
	for _, castle := range castles {
		if !castlePieces[castle] {
			return fmt.Errorf("invalid castle: %s", castles)
		}

		delete(castlePieces, castle)
	}

	return nil
}

// validateInPassant validates the in passant square.
func validateInPassant(square string) error {
	if square == "-" {
		return nil
	}

	if len(square) != 2 {
		return fmt.Errorf("invalid in passant square: %s", square)
	}

	if square[0] < 'a' || square[0] > 'h' {
		return fmt.Errorf("invalid in passant square: %s", square)
	}

	if square[1] < '1' || square[1] > '8' {
		return fmt.Errorf("invalid in passant square: %s", square)
	}

	return nil
}

// destinationMatch looks for a destination in a list of moves.
// It returns true if any of the moves has the destination.
// The function expects the moves in UCI format.
func destinationMatch(moves []string, destination pkg.Coordinate) bool {
	algCoor := pkg.CoordinateToAlgebraic(destination)
	for _, move := range moves {
		dest := move[2:4]
		if dest == algCoor {
			return true
		}
	}

	return false
}
