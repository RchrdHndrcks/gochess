package chess

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/RchrdHndrcks/gochess/pkg"
)

type (
	// Board represents a chess board.
	Board interface {
		// SetSquare sets a piece in a square.
		SetSquare(c pkg.Coordinate, p int8) error
		// Square returns the piece in a square.
		Square(c pkg.Coordinate) (int8, error)
		// MakeMove makes a move without checking if it is legal.
		MakeMove(origin, target pkg.Coordinate) error
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
		// enPassantSquare is the square where a pawn can capture in passant.
		enPassantSquare string
	}

	// Chess represents a Chess game.
	Chess struct {
		board Board
		// turn is the current turn.
		turn int8
		// movesCount is the number of moves played in algebaric notation.
		// It will increase by 1 after each black move.
		movesCount uint64
		// halfMoves is the number of half moves since the last capture or pawn move.
		halfMoves int
		// enPassantSquare is the square where a pawn can capture in passant.
		enPassantSquare string
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

	castleRook = map[string]pkg.Coordinate{
		"e1g1": pkg.Coor(7, 7),
		"e1c1": pkg.Coor(0, 7),
		"e8g8": pkg.Coor(7, 0),
		"e8c8": pkg.Coor(0, 0),
	}
)

// New creates a new chess game.
func New(opts ...Option) (*Chess, error) {
	c := &Chess{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	defaultOptions(c)
	return c, nil
}

// LoadPosition loads a board from a FEN string.
//
// The function will read the entire FEN string and will return an error if
// the FEN string is invalid.
//
// The board and properties will not be modified if the FEN string is invalid.
func (c *Chess) LoadPosition(FEN string) error {
	fenRows := strings.Split(FEN, "/")
	if len(fenRows) != 8 {
		return fmt.Errorf("invalid FEN: %s", FEN)
	}

	props := strings.Split(fenRows[7], " ")
	if len(props) != 6 {
		return fmt.Errorf("invalid FEN: %s", FEN)
	}

	fenRows[7] = props[0]

	var whiteKing, blackKing int

	brd := [8][8]int8{}
	for y := range 8 {
		row := [8]int8{}

		if len(fenRows[y]) == 0 || len(fenRows[y]) > 8 {
			return fmt.Errorf("invalid FEN: %s", FEN)
		}

		for x := range 8 {
			char := string(fenRows[y][0])
			fenRows[y] = fenRows[y][1:]

			n, err := strconv.Atoi(char)
			// If c is not a number, it's a piece.
			if err != nil {
				p := pkg.Pieces[char]
				row[x] = p
				if p == pkg.King|pkg.White {
					whiteKing++
				}
				if p == pkg.King|pkg.Black {
					blackKing++
				}
				continue
			}

			// If c is a number, reduce the number of empty squares from the FEN string.
			n--
			if n > 0 {
				fenRows[y] = fmt.Sprintf("%d%s", n, fenRows[y])
			}
		}

		if fenRows[y] != "" && fenRows[y] != "0" {
			return fmt.Errorf("invalid FEN: %s", FEN)
		}

		brd[y] = row
	}

	// If any of the kings is not in the board, the position is invalid.
	if whiteKing != 1 || blackKing != 1 {
		return errors.New("invalid FEN: both kings must be in the board once")
	}

	// Make a copy of the actual chess struct because if
	// the properties are invalid or the position is invalid
	// the struct will not be modified.
	copy := *c
	c.setBoard(brd)

	// If the FEN is invalid, setProperties will
	// return an error without modifying the board or the properties.
	if err := c.setProperties(FEN); err != nil {
		*c = copy
		return fmt.Errorf("invalid FEN: %w", err)
	}

	if !c.isPositionLegal() {
		*c = copy
		return errors.New("invalid FEN: the current turn can capture the opponent king")
	}

	return nil
}

// FEN returns the FEN string of the current position.
//
// If any of the kings is not in the board, the function returns an empty string.
func (c *Chess) FEN() string {
	if !c.isPositionValid() {
		return ""
	}

	fen := ""
	for y := range c.board.Width() {
		empty := 0
		for x := range c.board.Width() {
			// Ignore errors since the coordinates are valid.
			piece, _ := c.board.Square(pkg.Coor(x, y))

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

	ac := cmp.Or(c.availableCastles, "-")
	ips := cmp.Or(c.enPassantSquare, "-")

	fen += fmt.Sprintf(" %s %s %s %d %d", pkg.ColorNames[c.turn], ac, ips, c.halfMoves, c.movesCount)
	return fen
}

// AvailableMoves returns the available legal moves for the current turn.
// It returns an empty slice if position is stalemate.
// It returns a nil slice if position is checkmate.
func (c *Chess) AvailableMoves() []string {
	moves := c.availableMoves()

	legalMoves := []string{}
	for _, move := range moves {
		if c.isLegalMove(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	// If there are no legal moves, we need to check if position is checkmate or
	// stalemate.
	if len(legalMoves) == 0 {
		if c.isCheck() {
			legalMoves = nil
		}
	}

	return legalMoves
}

func (c Chess) isLegalMove(move string) bool {
	kingsColor := c.turn
	c.makeMove(move)

	availableMoves := c.availableMoves()
	kingPosition := c.kingsPosition(kingsColor)

	kingUnderAttack := destinationMatch(availableMoves, kingPosition)

	c.unmakeMove()
	kingWayUnderAttack := false
	if c.isCastleMove(move) {
		kingWayUnderAttack = destinationMatch(availableMoves, castleKingWay[move])
	}

	return !kingUnderAttack && !kingWayUnderAttack
}

// MakeMove checks if the move is legal and makes it.
// It returns an error if the move is not legal.
func (c *Chess) MakeMove(move string) error {
	moves := c.AvailableMoves()

	if !slices.Contains(moves, move) {
		return fmt.Errorf("move is not legal: %s", move)
	}

	c.makeMove(move)
	return nil
}

// makeMove makes a move without checking if it is legal.
func (c *Chess) makeMove(move string) {
	lastFEN := c.FEN()

	// Ignore the error because the move should be already validated.
	o, _ := AlgebraicToCoordinate(move[:2])
	t, _ := AlgebraicToCoordinate(move[2:4])

	if c.isCastleMove(move) {
		// If the move is a castle move, we need to move the rook too.
		rookOrigin := castleRook[move]
		rookTarget := pkg.Coor((o.X+t.X)/2, o.Y)

		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.MakeMove(rookOrigin, rookTarget)
	}

	if c.isEnPassantMove(move) {
		// If the move is an en passant capture, we need to remove the captured pawn.
		// The captured pawn is behind the target square.
		behindTarget := pkg.Coor(t.X, o.Y)
		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.SetSquare(behindTarget, pkg.Empty)
	}

	var madeMove bool
	// UCI moves only permit 5 characters if the move is a pawn coronation.
	isPromotion := len(move) == 5
	if isPromotion {
		p := pkg.PiecesWithoutColor[move[4:5]]
		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.SetSquare(t, p|c.turn)
		_ = c.board.SetSquare(o, pkg.Empty)
		madeMove = true
	}

	if !madeMove {
		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.MakeMove(o, t)
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
}

// UnmakeMove unmake the last move.
//
// It searches for the last move in the history and unmake it.
// If there are no moves in the history, the function does nothing.
func (c *Chess) UnmakeMove() {
	c.unmakeMove()
}

// unmakeMove is a helper function to unmake the last move.
func (c *Chess) unmakeMove() {
	if len(c.history) == 0 {
		return
	}

	lastMove := c.history[len(c.history)-1]
	c.history = c.history[:len(c.history)-1]

	lastFEN := lastMove.fen

	// Ignore the error because the FEN is valid since it was on the board.
	_ = c.LoadPosition(lastFEN)

	c.halfMoves = lastMove.halfMove
	c.availableCastles = lastMove.availableCastles
	c.enPassantSquare = lastMove.enPassantSquare

	// If turn color is white, last move was black.
	// So we decrease the moves count.
	if c.turn == pkg.White && c.movesCount > 1 {
		c.movesCount--
	}
}

// IsCheck returns if the current turn is in check.
func (c *Chess) IsCheck() bool {
	return c.isCheck()
}

// Square returns the piece in a square.
// The square is represented by an algebraic notation.
//
// If the square is not valid, the function returns an error.
func (c Chess) Square(square string) (string, error) {
	coor, err := AlgebraicToCoordinate(square)
	if err != nil {
		return "", fmt.Errorf("failed to convert algebraic notation to coordinate: %w", err)
	}

	// Ignore the error because the coordinates is
	// already validated.
	p, _ := c.board.Square(coor)
	return pkg.PieceNames[p], nil
}

// setBoard is a helper function that sets the board with a 2D array of pieces.
// The array should have 8 rows and 8 columns.
// The pieces are not validated.
func (c *Chess) setBoard(rows [8][8]int8) {
	for y := range 8 {
		for x := range 8 {
			c.board.SetSquare(pkg.Coor(x, y), rows[y][x])
		}
	}
}

// setProperties is a helper function that sets the properties of the Chess struct.
// It verifies the properties of the FEN string.
// It returns an error if the FEN string is invalid.
func (c *Chess) setProperties(FEN string) error {
	props := strings.Split(FEN, " ")[1:]

	color, ok := pkg.Colors[props[0]]
	if !ok {
		return fmt.Errorf("invalid color: %s", props[0])
	}

	availableCastles := props[1]
	if err := c.validateCastles(availableCastles); err != nil {
		return fmt.Errorf("invalid castles: %s", availableCastles)
	}

	enPassantSquare := props[2]
	if err := c.validateInPassant(enPassantSquare); err != nil {
		return fmt.Errorf("invalid en passant square: %s", enPassantSquare)
	}

	halfMoves, err := strconv.Atoi(props[3])
	if err != nil {
		return fmt.Errorf("invalid half moves: %s", props[3])
	}

	movesCount, err := strconv.ParseUint(props[4], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid moves count: %s", props[4])
	}

	c.turn = color
	c.availableCastles = availableCastles
	c.enPassantSquare = props[2]
	c.halfMoves = halfMoves
	c.movesCount = movesCount
	return nil
}

// kingsPosition returns the position of the king of the given color.
func (c Chess) kingsPosition(color int8) pkg.Coordinate {
	coor := pkg.Coordinate{}
	for y := range c.board.Width() {
		for x := range c.board.Width() {
			piece, _ := c.board.Square(pkg.Coor(x, y))
			if piece == pkg.King|color {
				coor = pkg.Coor(x, y)
				break
			}
		}
	}

	return coor
}

func (c Chess) isCheck() bool {
	kingPosition := c.kingsPosition(c.turn)

	c.toggleColor()
	defer c.toggleColor()
	moves := c.availableMoves()

	return destinationMatch(moves, kingPosition)
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

	k, _ := c.board.Square(pkg.Coor(4, 0))
	rr, _ := c.board.Square(pkg.Coor(7, 0))
	lr, _ := c.board.Square(pkg.Coor(0, 0))
	toBeRemoved["k"] = rr != pkg.Rook|pkg.Black || k != pkg.King|pkg.Black
	toBeRemoved["q"] = lr != pkg.Rook|pkg.Black || k != pkg.King|pkg.Black

	K, _ := c.board.Square(pkg.Coor(4, 7))
	rR, _ := c.board.Square(pkg.Coor(7, 7))
	lR, _ := c.board.Square(pkg.Coor(0, 7))
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
//
// It must be called after a move is made. If no move was made,
// the function will panic.
func (c *Chess) updateHalfMoves() {
	c.halfMoves++

	// First we look for a change in the board.
	// If we have less pieces than before, a capture was made so we reset the counter.
	h := c.history[len(c.history)-1]
	aux, _ := New(WithFEN(h.fen))
	piecesCount := 0
	auxPiecesCount := 0
	for y := range c.board.Width() {
		for x := range c.board.Width() {
			piece, _ := c.board.Square(pkg.Coor(x, y))
			auxPiece, _ := aux.board.Square(pkg.Coor(x, y))
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
	coor, _ := AlgebraicToCoordinate(origin)
	p, _ := aux.board.Square(coor)

	piece := p &^ (pkg.White | pkg.Black)
	if piece == pkg.Pawn {
		c.halfMoves = 0
	}
}

// updateInPassantSquare updates the en passant square.
//
// It must be called after a move is made. If no move was made,
// the function will panic.
func (c *Chess) updateInPassantSquare() {
	c.enPassantSquare = ""

	lastMove := c.history[len(c.history)-1].move
	if len(lastMove) != 4 {
		return
	}

	dest, _ := AlgebraicToCoordinate(lastMove[2:])
	p, _ := c.board.Square(dest)

	if p&^(pkg.White|pkg.Black) != pkg.Pawn {
		return
	}

	destRow, _ := strconv.Atoi(lastMove[3:4])
	orgRow, _ := strconv.Atoi(lastMove[1:2])
	if destRow == orgRow+2 || destRow == orgRow-2 {
		c.enPassantSquare = fmt.Sprintf("%s%d", lastMove[2:3], (destRow+orgRow)/2)
		return
	}
}

// isCastleMove returns if the move is a castle move.
//
// The passed move must be valid.
func (c Chess) isCastleMove(move string) bool {
	if castlesMoves[move] != c.turn {
		return false
	}

	origin, _ := AlgebraicToCoordinate(move[:2])
	p, _ := c.board.Square(origin)

	return p == pkg.King|c.turn
}

// isEnPassantMove returns if the move is an en passant capture move.
//
// The passed move must be valid.
func (c Chess) isEnPassantMove(move string) bool {
	if c.enPassantSquare == "" {
		return false
	}

	origin, _ := AlgebraicToCoordinate(move[:2])

	if move[2:4] != c.enPassantSquare {
		return false
	}

	p, _ := c.board.Square(origin)
	return p&^(pkg.White|pkg.Black) == pkg.Pawn
}

// isPositionValid checks if both kings are in the board once.
func (c Chess) isPositionValid() bool {
	var whiteKings, blackKings int
	for y := range c.board.Width() {
		for x := range c.board.Width() {
			piece, _ := c.board.Square(pkg.Coor(x, y))

			if piece == pkg.King|pkg.White {
				whiteKings++
			}

			if piece == pkg.King|pkg.Black {
				blackKings++
			}
		}
	}

	return whiteKings == 1 && blackKings == 1
}

// isPositionLegal verifies if the current turn can capture the opponent king.
//
// If the current turn can capture the opponent king, the position is not legal
// and the function returns false.
func (c Chess) isPositionLegal() bool {
	c.toggleColor()
	return !c.isCheck()
}

// validateCastles validates the castles string.
func (Chess) validateCastles(castles string) error {
	if castles == "-" {
		return nil
	}

	castlePieces := map[rune]bool{'K': true, 'Q': true, 'k': true, 'q': true}
	for _, castle := range castles {
		if !castlePieces[castle] {
			return errors.New("invalid castles")
		}

		delete(castlePieces, castle)
	}

	return nil
}

// validateInPassant validates the in passant square.
func (c Chess) validateInPassant(square string) error {
	if square == "-" {
		return nil
	}

	coor, err := AlgebraicToCoordinate(square)
	if err != nil {
		return errors.New("invalid in passant square")
	}

	if coor.Y != 2 && coor.Y != 5 {
		return errors.New("invalid in passant square")
	}

	yCoor := 4
	color := pkg.White
	if coor.Y == 2 {
		yCoor = 3
		color = pkg.Black
	}

	auxCoor := pkg.Coor(coor.X, yCoor)
	p, _ := c.board.Square(auxCoor)
	if p&^color != pkg.Pawn {
		return errors.New("invalid in passant square")
	}

	return nil
}

// destinationMatch looks for a destination in a list of moves.
// It returns true if any of the moves has the destination.
// The function expects the moves in UCI format.
func destinationMatch(moves []string, destination pkg.Coordinate) bool {
	algCoor := CoordinateToAlgebraic(destination)
	for _, move := range moves {
		dest := move[2:4]
		if dest == algCoor {
			return true
		}
	}

	return false
}

// availableMoves returns the available moves for the current turn without checking if they are legal.
func (c *Chess) availableMoves() []string {
	moves := []string{}
	for x := range 8 {
		for y := range 8 {
			origin := pkg.Coor(x, y)
			piece, _ := c.board.Square(origin)
			if piece&c.turn == pkg.Empty {
				continue
			}

			moves = append(moves, c.movesForPiece(piece, origin)...)
		}
	}

	return moves
}

// movesForPiece returns the available moves for a piece.
//
// The function returns a slice of UCI moves.
// (e.g. "e2e4" for moving the piece at e2 to e4.)
// Disclaimer: This function does not check if the move is legal for a Chess game.
func (c Chess) movesForPiece(piece int8, origin pkg.Coordinate) []string {
	switch piece &^ (pkg.White | pkg.Black) {
	case pkg.Pawn:
		return c.pawnMoves(origin)
	case pkg.Rook:
		return c.rookMoves(origin)
	case pkg.Queen:
		return c.queenMoves(origin)
	case pkg.King:
		return append(c.kingMoves(origin), c.kingCastleMoves(origin)...)
	case pkg.Bishop:
		return c.bishopMoves(origin)
	case pkg.Knight:
		return c.knightMoves(origin)
	}

	return nil
}

// pawnMoves returns all the valid pawn moves.
func (c Chess) pawnMoves(origin pkg.Coordinate) []string {
	p, _ := c.board.Square(origin)
	dir := -1
	if p&pkg.White == pkg.Empty {
		dir = 1
	}

	isPromotion := false
	tCor := pkg.Coor(origin.X, origin.Y+1*dir)
	if tCor.Y == 7 || tCor.Y == 0 {
		isPromotion = true
	}

	moves := make([]string, 0, 2)
	s, _ := c.board.Square(tCor)
	if s == pkg.Empty {
		moves = append(moves, UCI(origin, tCor))
	}

	if isPromotion {
		return append(c.pawnCaptureMoves(origin, true), c.promotionPosibilities(origin, tCor)...)
	}

	if !(dir == 1 && origin.Y == 1) && !(dir == -1 && origin.Y == 6) {
		return append(c.pawnCaptureMoves(origin, false), moves...)
	}

	tCor = pkg.Coor(origin.X, origin.Y+2*dir)
	s, _ = c.board.Square(tCor)
	if s == pkg.Empty {
		moves = append(moves, UCI(origin, tCor))
	}

	return append(c.pawnCaptureMoves(origin, false), moves...)
}

// pawnCaptureMoves returns all the valid pawn capture moves.
func (c Chess) pawnCaptureMoves(origin pkg.Coordinate, isPromotion bool) []string {
	p, _ := c.board.Square(origin)
	pColor := p & (pkg.White | pkg.Black)
	dir := -1
	if pColor == pkg.Black {
		dir = 1
	}

	moves := []string{}
	offsets := []int{-1, 1}
	for _, o := range offsets {
		tCor := pkg.Coor(origin.X+o, origin.Y+1*dir)
		if tCor.X < 0 || tCor.X > 7 || tCor.Y < 0 || tCor.Y > 7 {
			continue
		}

		if CoordinateToAlgebraic(tCor) == c.enPassantSquare {
			moves = append(moves, UCI(origin, tCor))
			continue
		}

		ts, _ := c.board.Square(tCor)
		if ts == pkg.Empty || ts&pColor != pkg.Empty {
			continue
		}

		if !isPromotion {
			moves = append(moves, UCI(origin, tCor))
			continue
		}

		moves = append(moves, c.promotionPosibilities(origin, tCor)...)
	}

	return moves
}

// promotionPosibilities is a helper function that returns the UCI moves with
// the value of the piece to be promoted.
func (c Chess) promotionPosibilities(origin, target pkg.Coordinate) []string {
	moves := make([]string, 4)
	for i, p := range []int8{pkg.Queen, pkg.Rook, pkg.Bishop, pkg.Knight} {
		moves[i] = UCI(origin, target, p)
	}

	return moves
}

// knightMoves returns valid knight moves.
func (c Chess) knightMoves(origin pkg.Coordinate) []string {
	offsets := []pkg.Coordinate{
		{X: 1, Y: 2}, {X: 2, Y: 1},
		{X: 1, Y: -2}, {X: 2, Y: -1},
		{X: -1, Y: 2}, {X: -2, Y: 1},
		{X: -1, Y: -2}, {X: -2, Y: -1},
	}

	return c.oneStepPieces(origin, offsets)
}

// kingMoves returns valid king moves.
func (c Chess) kingMoves(origin pkg.Coordinate) []string {
	offsets := []pkg.Coordinate{
		{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: -1},
		{X: 0, Y: 1}, {X: 0, Y: -1},
		{X: -1, Y: 1}, {X: -1, Y: 0}, {X: -1, Y: -1},
	}

	return c.oneStepPieces(origin, offsets)
}

// kingCastleMoves returns valid castle moves.
func (c Chess) kingCastleMoves(origin pkg.Coordinate) []string {
	if c.availableCastles == "-" {
		return nil
	}

	p, _ := c.board.Square(origin)
	kingColor := p & (pkg.White | pkg.Black)

	castleDirections := map[string]int{
		"k": 1, "K": 1,
		"q": -1, "Q": -1,
	}

	moves := []string{}
	for castle, dir := range castleDirections {
		if !strings.Contains(c.availableCastles, castle) {
			continue
		}

		if pkg.Pieces[castle]&kingColor == pkg.Empty {
			continue
		}

		ts, err := c.board.Square(pkg.Coor(origin.X+dir, origin.Y))
		if err != nil || ts != pkg.Empty {
			continue
		}

		ts, err = c.board.Square(pkg.Coor(origin.X+2*dir, origin.Y))
		if err != nil || ts != pkg.Empty {
			continue
		}

		moves = append(moves, UCI(origin, pkg.Coor(origin.X+2*dir, origin.Y)))

		if len(moves) == 2 {
			break
		}
	}

	return moves
}

// rookMoves returns valid rook moves.
func (c Chess) rookMoves(origin pkg.Coordinate) []string {
	offsets := []pkg.Coordinate{{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: -1}}
	return c.slidingPieces(origin, offsets)
}

// bishopMoves returns valid bishop moves.
func (c Chess) bishopMoves(origin pkg.Coordinate) []string {
	offsets := []pkg.Coordinate{{X: 1, Y: 1}, {X: -1, Y: 1}, {X: 1, Y: -1}, {X: -1, Y: -1}}
	return c.slidingPieces(origin, offsets)
}

// queenMoves returns valid queen moves.
func (c Chess) queenMoves(origin pkg.Coordinate) []string {
	return append(c.rookMoves(origin), c.bishopMoves(origin)...)
}

// slidingPieces returns valid moves for sliding pieces.
func (c Chess) slidingPieces(origin pkg.Coordinate, offsets []pkg.Coordinate) []string {
	p, _ := c.board.Square(origin)

	color := p & (pkg.White | pkg.Black)
	moves := []string{}
	for _, d := range offsets {
		for i := 1; ; i++ {
			tCor := pkg.Coor(origin.X+i*d.X, origin.Y+i*d.Y)
			ts, err := c.board.Square(tCor)
			if err != nil {
				break
			}

			if ts == pkg.Empty {
				moves = append(moves, UCI(origin, tCor))
				continue
			}

			if ts&color == pkg.Empty {
				moves = append(moves, UCI(origin, tCor))
				break
			}

			// If the piece is the same color, stop looking in that direction.
			break
		}
	}

	return moves
}

func (c Chess) oneStepPieces(origin pkg.Coordinate, offsets []pkg.Coordinate) []string {
	p, _ := c.board.Square(origin)

	color := p & (pkg.White | pkg.Black)
	moves := []string{}
	for _, d := range offsets {
		tCor := pkg.Coor(origin.X+d.X, origin.Y+d.Y)
		ts, err := c.board.Square(tCor)
		if err != nil {
			continue
		}

		if ts == pkg.Empty {
			moves = append(moves, UCI(origin, tCor))
			continue
		}

		if ts&color == pkg.Empty {
			moves = append(moves, UCI(origin, tCor))
		}
	}

	return moves
}
