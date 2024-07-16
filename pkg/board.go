package pkg

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// board is a 2D array of pieces.
type board [8][8]int8

// NewBoard creates a new board.
func NewBoard() *board {
	return &board{}
}

// Width returns the width of the board.
func (b *board) Width() int {
	return 8
}

// Square returns the piece at the given Coordinate.
func (b *board) Square(c Coordinate) (int8, error) {
	if c.x > 7 || c.y > 7 || c.x < 0 || c.y < 0 {
		return 0, fmt.Errorf("Coordinate out of bounds")
	}

	return b[c.y][c.x], nil
}

// AvailableMoves returns the available moves for the current turn.
func (b *board) AvailableMoves(turn int8, inPassantSquare, castlePossibilities string,
) ([]string, error) {
	if turn != White && turn != Black {
		return nil, fmt.Errorf("invalid turn color: %d", turn)
	}

	moves := []string{}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			piece, _ := b.Square(Coordinate{x, y}) // nolint:errcheck
			if piece&turn == Empty {
				continue
			}

			origin := Coordinate{x, y}
			moves = append(moves, b.movesForPiece(origin, inPassantSquare, castlePossibilities)...)
		}
	}

	return moves, nil
}

// MakeMove makes a move on the board.
// The move should be in UCI format.
func (b *board) MakeMove(move string) error {
	if len(move) < 4 || len(move) > 5 {
		return fmt.Errorf("move length should be 4 or 5")
	}

	origin, err := AlgebraicToCoordinate(move[:2])
	if err != nil {
		return fmt.Errorf("error in parsing origin square to algebraic: %w", err)
	}

	target, err := AlgebraicToCoordinate(move[2:4])
	if err != nil {
		return fmt.Errorf("error in parsing target square to algebraic: %w", err)
	}

	p, _ := b.Square(origin) // nolint:errcheck
	if p == Empty {
		return fmt.Errorf("origin square is empty")
	}

	// If the move is a coronation move, the origin piece should be replaced
	// by the coronated piece.
	if len(move) == 5 {
		var err error
		p, err = coronationPieceFromMove(move)
		if err != nil {
			return fmt.Errorf("error getting new piece: %w", err)
		}
	}

	// If the move is a castle move, the rook should be moved as well.
	isCastleMove, rookCoor, rookDest := b.isCastleMove(move)
	if isCastleMove {
		b[rookDest.y][rookDest.x] = b[rookCoor.y][rookCoor.x]
		b[rookCoor.y][rookCoor.x] = Empty
	}

	// If the move is a in passant capture, the captured pawn should be removed.
	if p&^(White|Black) == Pawn && origin.x != target.x && b[target.y][target.x] == Empty {
		dir := -1
		if p&Black == Empty {
			dir = 1
		}

		b[target.y+dir][target.x] = Empty
	}

	b[origin.y][origin.x] = Empty
	b[target.y][target.x] = p

	return nil
}

// LoadPosition loads a board from a FEN string.
func (b *board) LoadPosition(FEN string) error {
	fenRows := strings.Split(FEN, "/")
	if len(fenRows) != 8 {
		return fmt.Errorf("invalid FEN: %s", FEN)
	}

	props := strings.Split(fenRows[7], " ")
	if len(props) != 6 {
		return fmt.Errorf("invalid FEN: %s", FEN)
	}

	fenRows[7] = props[0]

	brd := [8][8]int8{}
	for y := 0; y < 8; y++ {
		row := [8]int8{}

		if len(fenRows[y]) == 0 || len(fenRows[y]) > 8 {
			return fmt.Errorf("invalid FEN: %s", FEN)
		}

		for x := 0; x < 8; x++ {
			char := string(fenRows[y][0])
			fenRows[y] = fenRows[y][1:]

			n, err := strconv.Atoi(char)
			// If c is not a number, it's a piece.
			if err != nil {
				row[x] = Pieces[char]
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

	*b = board(brd)
	return nil
}

// movesForPiece returns the available moves for a piece.
// It receives the origin of the piece, the in passant square and castle possibilities
// The in passant square could be empty if there is no in passant square.
// The castlePossibilities expected format is:
// "KQkq" where:
// - K means White can castle kingside.
// - Q means White can castle queenside.
// - k means Black can castle kingside.
// - q means Black can castle queenside.
// If any of the letters is missing, it means that the side cannot castle.
// If there is no castle possibility, the string should be empty.
// The function doesnt check if inPassantSquare is a valid square.
// The function doesnt check if castlePossibilities is a valid string.
// The function a slice of UCI moves.
// (e.g. "e2e4" for moving the piece at e2 to e4.)
// Disclaimer: This function does not check if the move is legal for a Chess game.
func (b *board) movesForPiece(origin Coordinate, inPassantSquare string, castlePossibilities string,
) []string {
	p, err := b.Square(origin)
	if err != nil || p == Empty {
		return nil
	}

	switch p &^ (White | Black) {
	case Pawn:
		coronationMoves := b.pawnCoronationMoves(origin)
		if len(coronationMoves) > 0 {
			return coronationMoves
		}

		return append(b.pawnMoves(origin), b.pawnCaptureMoves(origin, inPassantSquare)...)
	case Rook:
		return b.rookMoves(origin)
	case Queen:
		return b.queenMoves(origin)
	case King:
		return append(b.kingMoves(origin),
			b.kingCastleMoves(origin, castlePossibilities)...)
	case Bishop:
		return b.bishopMoves(origin)
	case Knight:
		return b.knightMoves(origin)
	}

	return nil
}

// pawnMoves returns valid pawn moves.
func (b board) pawnMoves(origin Coordinate) []string {
	p, _ := b.Square(origin) // nolint:errcheck
	dir := -1
	if p&White == Empty {
		dir = 1
	}

	tCor := Coor(origin.x, origin.y+1*dir)
	// If the target square is the last or first row, it is a coronation move,
	// so it will be handled by the coronation function.
	if tCor.y == 7 || tCor.y == 0 {
		return nil
	}

	ts, err := b.Square(tCor)
	if err != nil {
		return nil
	}

	if ts != Empty {
		return nil
	}

	moves := []string{UCI(origin, tCor)}

	if origin.y != 1 && origin.y != 6 {
		return moves
	}

	tCor = Coordinate{origin.x, origin.y + 2*dir}
	ts, err = b.Square(tCor)
	if err != nil {
		return moves
	}

	if ts != Empty {
		return moves
	}

	return append(moves, UCI(origin, tCor))
}

// pawnCaptureMoves returns valid pawn capture moves.
// It doesn't check if the move is a coronation move, so pawnCoronationMoves
// should be called first.
func (b board) pawnCaptureMoves(origin Coordinate, inPassantSquare string) []string {
	p, _ := b.Square(origin) // nolint:errcheck
	pColor := p & (White | Black)
	dir := -1
	if p&White == Empty {
		dir = 1
	}

	moves := []string{}
	offsets := []int{-1, 1}
	for _, o := range offsets {
		tx := origin.x + o
		ty := origin.y + 1*dir
		if tx < 0 || tx > 7 || ty < 0 || ty > 7 {
			continue
		}

		tCor := Coor(origin.x+o, origin.y+1*dir)

		if CoordinateToAlgebraic(tCor) == inPassantSquare {
			moves = append(moves, UCI(origin, tCor))
			continue
		}

		ts, err := b.Square(tCor)
		if err != nil {
			continue
		}

		if ts == Empty || ts&pColor != Empty {
			continue
		}

		moves = append(moves, UCI(origin, tCor))
	}

	return moves
}

// pawnCoronationMoves returns valid pawn coronation moves.
// It also call the pawnCaptureMoves function to check if there is a capture move.
// If any move is returned there is not necessary to check pawnCaptureMoves.
// It also handle normal pawn moves.
func (b board) pawnCoronationMoves(origin Coordinate) []string {
	p, _ := b.Square(origin) // nolint:errcheck
	dir := -1
	if p&White == Empty {
		dir = 1
	}

	ts := Coor(origin.x, origin.y+1*dir)
	if ts.y != 7 && ts.y != 0 {
		return nil
	}

	possibleCoor := []Coordinate{}
	s, err := b.Square(ts)
	if err != nil {
		return nil
	}

	if s == Empty {
		possibleCoor = append(possibleCoor, ts)
	}

	// pawnCaptureMoves doesn't check if the move is a coronation move,
	// so we take possible captures at this point and convert them to
	// coronation moves.
	// TODO: Design a better way to handle this.
	captureMoves := b.pawnCaptureMoves(origin, "")
	for _, move := range captureMoves {
		tCor, err := AlgebraicToCoordinate(move[2:4])
		if err != nil {
			continue
		}

		possibleCoor = append(possibleCoor, tCor)
	}

	moves := []string{}
	for _, t := range possibleCoor {
		for _, p := range []int8{Black | Queen, Black | Rook, Black | Bishop, Black | Knight} {
			move := UCI(origin, t) + PieceNames[p]
			moves = append(moves, move)
		}
	}

	return moves
}

// knightMoves returns valid knight moves.
func (b board) knightMoves(origin Coordinate) []string {
	offsets := []Coordinate{
		{1, 2}, {2, 1},
		{1, -2}, {2, -1},
		{-1, 2}, {-2, 1},
		{-1, -2}, {-2, -1},
	}

	return b.oneStepPieces(origin, offsets)
}

// kingMoves returns valid king moves.
func (b board) kingMoves(origin Coordinate) []string {
	offsets := []Coordinate{
		{1, 1}, {1, 0}, {1, -1},
		{0, 1}, {0, -1},
		{-1, 1}, {-1, 0}, {-1, -1},
	}

	return b.oneStepPieces(origin, offsets)
}

// kingCastleMoves returns valid castle moves.
func (b board) kingCastleMoves(origin Coordinate, castlePossibilities string) []string {
	if castlePossibilities == "" {
		return nil
	}

	p, err := b.Square(origin)
	if err != nil {
		return nil
	}

	kingColor := p & (White | Black)

	castleDirections := map[string]int{
		"k": 1, "K": 1,
		"q": -1, "Q": -1,
	}

	moves := []string{}
	for castle, dir := range castleDirections {
		if !strings.Contains(castlePossibilities, castle) {
			continue
		}

		if Pieces[castle]&kingColor == Empty {
			continue
		}

		ts, err := b.Square(Coor(origin.x+dir, origin.y))
		if err != nil || ts != Empty {
			continue
		}

		ts, err = b.Square(Coor(origin.x+2*dir, origin.y))
		if err != nil || ts != Empty {
			continue
		}

		moves = append(moves, UCI(origin, Coor(origin.x+2*dir, origin.y)))

		if len(moves) == 2 {
			break
		}
	}

	return moves
}

// rookMoves returns valid rook moves.
func (b board) rookMoves(origin Coordinate) []string {
	offsets := []Coordinate{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	return b.slidingPieces(origin, offsets)
}

// bishopMoves returns valid bishop moves.
func (b board) bishopMoves(origin Coordinate) []string {
	offsets := []Coordinate{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	return b.slidingPieces(origin, offsets)
}

// bishopMoves returns valid bishop moves.
func (b board) queenMoves(origin Coordinate) []string {
	return append(b.rookMoves(origin), b.bishopMoves(origin)...)
}

// slidingPieces returns valid moves for sliding pieces.
func (b board) slidingPieces(origin Coordinate, offsets []Coordinate) []string {
	p, err := b.Square(origin)
	if err != nil {
		return nil
	}

	color := p & (White | Black)
	moves := []string{}
	for _, d := range offsets {
		for i := 1; ; i++ {
			tCor := Coordinate{origin.x + i*d.x, origin.y + i*d.y}
			ts, err := b.Square(tCor)
			if err != nil {
				break
			}

			if ts == Empty {
				moves = append(moves, UCI(origin, tCor))
				continue
			}

			if ts&color == Empty {
				moves = append(moves, UCI(origin, tCor))
				break
			}

			// If the piece is the same color, stop looking in that direction.
			break
		}
	}

	return moves
}

func (b board) oneStepPieces(origin Coordinate, offsets []Coordinate) []string {
	p, err := b.Square(origin)
	if err != nil {
		return nil
	}

	color := p & (White | Black)
	moves := []string{}
	for _, d := range offsets {
		tCor := Coordinate{origin.x + d.x, origin.y + d.y}
		ts, err := b.Square(tCor)
		if err != nil {
			continue
		}

		if ts == Empty {
			moves = append(moves, UCI(origin, tCor))
			continue
		}

		if ts&color == Empty {
			moves = append(moves, UCI(origin, tCor))
		}
	}

	return moves
}

// isCastleMove returns if the move is a castle move and the rook origin and target.
// If the move is not a castle move, it returns false and empty Coordinates.
func (b board) isCastleMove(move string) (bool, Coordinate, Coordinate) {
	castleMoves := []string{"e1g1", "e1c1", "e8g8", "e8c8"}
	if !slices.Contains(castleMoves, move) {
		return false, Coordinate{}, Coordinate{}
	}

	origin, err := AlgebraicToCoordinate(move[:2])
	if err != nil {
		return false, Coordinate{}, Coordinate{}
	}

	p, _ := b.Square(origin) // nolint:errcheck
	if p&King != King {
		return false, Coordinate{}, Coordinate{}
	}

	rookCoors := map[string]Coordinate{
		"e1g1": {7, 7}, "e1c1": {0, 7},
		"e8g8": {7, 0}, "e8c8": {0, 0},
	}

	rookDest := map[string]Coordinate{
		"e1g1": {5, 7}, "e1c1": {3, 7},
		"e8g8": {5, 0}, "e8c8": {3, 0},
	}

	return true, rookCoors[move], rookDest[move]
}

// coronationPieceFromMove returns the piece to coronate to.
// The move should be in UCI format.
// If the move is not a coronation move, it returns an error.
func coronationPieceFromMove(move string) (int8, error) {
	if len(move) != 5 {
		return Empty, fmt.Errorf("invalid move length: %s", move)
	}

	target, err := AlgebraicToCoordinate(move[2:4])
	if err != nil {
		return Empty, fmt.Errorf("error in parsing target square to algebraic: %w", err)
	}

	coronationPiece := Pieces[string(move[4])]
	if target.y != 0 && target.y != 7 {
		return Empty, fmt.Errorf("invalid coronation target row: %d", target.y)
	}

	if target.y == 0 {
		coronationPiece = (coronationPiece &^ Black) | White
	}

	return coronationPiece, nil
}
