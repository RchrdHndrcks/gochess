package chess

import (
	"strings"

	"github.com/RchrdHndrcks/gochess"
)

// capacityByPiece is a helper map where the key is the piece and the value is
// the max moves count for that piece. It is useful to know which capacity asign
// to the slice where the moves of these pieces are stored.
var capacityByPiece = map[int8]int{
	gochess.White | gochess.Queen:  27,
	gochess.Black | gochess.Queen:  27,
	gochess.White | gochess.Rook:   14,
	gochess.Black | gochess.Rook:   14,
	gochess.White | gochess.Bishop: 13,
	gochess.Black | gochess.Bishop: 13,
}

// makeMove makes a move without checking if it is legal.
func (c *Chess) makeMove(move string) {
	lastFEN := c.actualFEN

	// Ignore the error because the move should be already validated.
	o, _ := AlgebraicToCoordinate(move[:2])
	t, _ := AlgebraicToCoordinate(move[2:4])

	if c.isCastleMove(move) {
		// If the move is a castle move, we need to move the rook too.
		rookOrigin := castleRook[move]
		rookTarget := gochess.Coor((o.X+t.X)/2, o.Y)

		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.MakeMove(rookOrigin, rookTarget)
	}

	if c.isEnPassantMove(move) {
		// If the move is an en passant capture, we need to remove the captured pawn.
		// The captured pawn is behind the target square.
		behindTarget := gochess.Coor(t.X, o.Y)
		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.SetSquare(behindTarget, gochess.Empty)
	}

	var madeMove bool
	// UCI moves only permit 5 characters if the move is a pawn coronation.
	isPromotion := len(move) == 5
	if isPromotion {
		p := gochess.PiecesWithoutColor[move[4:5]]
		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.SetSquare(t, p|c.turn)
		_ = c.board.SetSquare(o, gochess.Empty)
		madeMove = true
	}

	if !madeMove {
		// Ignore the error because the coordinates is valid because
		// the move is already validated.
		_ = c.board.MakeMove(o, t)
	}

	c.history = append(
		c.history,
		chessContext{
			move:              move,
			fen:               lastFEN,
			halfMove:          c.halfMoves,
			availableCastles:  c.availableCastles,
			enPassantSquare:   c.enPassantSquare,
			whiteKingPosition: c.whiteKingPosition,
			blackKingPosition: c.blackKingPosition,
		},
	)

	// If the origin is the king, update the king position.
	if o == *c.whiteKingPosition {
		c.whiteKingPosition = &t
	}

	if o == *c.blackKingPosition {
		c.blackKingPosition = &t
	}

	c.toggleColor()
	c.updateMovesCount()
	c.updateCastlePossibilities()
	c.updateHalfMoves()
	c.updateEnPassantSquare()
}

// unmakeMove is a helper function to unmake the last move.
func (c *Chess) unmakeMove() {
	if len(c.history) == 0 {
		return
	}

	lastContext := c.history[len(c.history)-1]
	c.history = c.history[:len(c.history)-1]

	lastFEN := lastContext.fen

	// Ignore the error because the FEN is valid since it was on the board.
	_ = c.loadPosition(lastFEN)
	c.whiteKingPosition = lastContext.whiteKingPosition
	c.blackKingPosition = lastContext.blackKingPosition
}

// movesForPiece returns the available moves for a piece.
//
// The function returns a slice of UCI moves.
// (e.g. "e2e4" for moving the piece at e2 to e4.)
// Disclaimer: This function does not check if the move is legal for a Chess game.
func (c Chess) movesForPiece(piece int8, origin gochess.Coordinate) []string {
	switch piece &^ (gochess.White | gochess.Black) {
	case gochess.Pawn:
		return c.pawnMoves(origin)
	case gochess.Rook:
		return c.rookMoves(origin)
	case gochess.Queen:
		return c.queenMoves(origin)
	case gochess.King:
		return append(c.kingMoves(origin), c.kingCastleMoves(origin)...)
	case gochess.Bishop:
		return c.bishopMoves(origin)
	case gochess.Knight:
		return c.knightMoves(origin)
	}

	return nil
}

// pawnMoves returns all the valid pawn moves.
func (c Chess) pawnMoves(origin gochess.Coordinate) []string {
	p, _ := c.board.Square(origin)
	dir := -1
	if p&gochess.White == gochess.Empty {
		dir = 1
	}

	isPromotion := false
	tCor := gochess.Coor(origin.X, origin.Y+1*dir)
	if tCor.Y == 7 || tCor.Y == 0 {
		isPromotion = true
	}

	moves := make([]string, 0, 2)
	s, _ := c.board.Square(tCor)
	if s == gochess.Empty {
		moves = append(moves, UCI(origin, tCor))
	}

	if isPromotion {
		return append(c.pawnCaptureMoves(origin, true), c.promotionPosibilities(origin, tCor)...)
	}

	if !(dir == 1 && origin.Y == 1) && !(dir == -1 && origin.Y == 6) {
		return append(c.pawnCaptureMoves(origin, false), moves...)
	}

	tCor = gochess.Coor(origin.X, origin.Y+2*dir)
	s, _ = c.board.Square(tCor)
	if s == gochess.Empty {
		moves = append(moves, UCI(origin, tCor))
	}

	return append(c.pawnCaptureMoves(origin, false), moves...)
}

// pawnCaptureMoves returns all the valid pawn capture moves.
func (c Chess) pawnCaptureMoves(origin gochess.Coordinate, isPromotion bool) []string {
	p, _ := c.board.Square(origin)
	pColor := p & (gochess.White | gochess.Black)
	dir := -1
	if pColor == gochess.Black {
		dir = 1
	}

	moves := make([]string, 0, 2)
	offsets := []int{-1, 1}
	for _, o := range offsets {
		tCor := gochess.Coor(origin.X+o, origin.Y+1*dir)
		if tCor.X < 0 || tCor.X > 7 || tCor.Y < 0 || tCor.Y > 7 {
			continue
		}

		if CoordinateToAlgebraic(tCor) == c.enPassantSquare {
			moves = append(moves, UCI(origin, tCor))
			continue
		}

		ts, _ := c.board.Square(tCor)
		if ts == gochess.Empty || ts&pColor != gochess.Empty {
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
func (c Chess) promotionPosibilities(origin, target gochess.Coordinate) []string {
	moves := make([]string, 4)
	for i, p := range []int8{gochess.Queen, gochess.Rook, gochess.Bishop, gochess.Knight} {
		moves[i] = UCI(origin, target, p)
	}

	return moves
}

// knightMoves returns valid knight moves.
func (c Chess) knightMoves(origin gochess.Coordinate) []string {
	offsets := []gochess.Coordinate{
		{X: 1, Y: 2}, {X: 2, Y: 1},
		{X: 1, Y: -2}, {X: 2, Y: -1},
		{X: -1, Y: 2}, {X: -2, Y: 1},
		{X: -1, Y: -2}, {X: -2, Y: -1},
	}

	return c.oneStepPieces(origin, offsets)
}

// kingMoves returns valid king moves.
func (c Chess) kingMoves(origin gochess.Coordinate) []string {
	offsets := []gochess.Coordinate{
		{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: -1},
		{X: 0, Y: 1}, {X: 0, Y: -1},
		{X: -1, Y: 1}, {X: -1, Y: 0}, {X: -1, Y: -1},
	}

	return c.oneStepPieces(origin, offsets)
}

// kingCastleMoves returns valid castle moves.
func (c Chess) kingCastleMoves(origin gochess.Coordinate) []string {
	if c.availableCastles == "-" {
		return nil
	}

	p, _ := c.board.Square(origin)
	kingColor := p & (gochess.White | gochess.Black)

	castleDirections := map[string]int{
		"k": 1, "K": 1,
		"q": -1, "Q": -1,
	}

	moves := make([]string, 0, 2)
	for castle, dir := range castleDirections {
		if !strings.Contains(c.availableCastles, castle) {
			continue
		}

		if gochess.Pieces[castle]&kingColor == gochess.Empty {
			continue
		}

		ts, err := c.board.Square(gochess.Coor(origin.X+dir, origin.Y))
		if err != nil || ts != gochess.Empty {
			continue
		}

		ts, err = c.board.Square(gochess.Coor(origin.X+2*dir, origin.Y))
		if err != nil || ts != gochess.Empty {
			continue
		}

		moves = append(moves, UCI(origin, gochess.Coor(origin.X+2*dir, origin.Y)))

		if len(moves) == 2 {
			break
		}
	}

	return moves
}

// rookMoves returns valid rook moves.
func (c Chess) rookMoves(origin gochess.Coordinate) []string {
	offsets := []gochess.Coordinate{{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: -1}}
	return c.slidingPieces(origin, offsets)
}

// bishopMoves returns valid bishop moves.
func (c Chess) bishopMoves(origin gochess.Coordinate) []string {
	offsets := []gochess.Coordinate{{X: 1, Y: 1}, {X: -1, Y: 1}, {X: 1, Y: -1}, {X: -1, Y: -1}}
	return c.slidingPieces(origin, offsets)
}

// queenMoves returns valid queen moves.
func (c Chess) queenMoves(origin gochess.Coordinate) []string {
	return append(c.rookMoves(origin), c.bishopMoves(origin)...)
}

// slidingPieces returns valid moves for sliding pieces.
func (c Chess) slidingPieces(origin gochess.Coordinate, offsets []gochess.Coordinate) []string {
	p, _ := c.board.Square(origin)

	color := p & (gochess.White | gochess.Black)
	moves := make([]string, 0, capacityByPiece[p])
	for _, d := range offsets {
		for i := 1; ; i++ {
			tCor := gochess.Coor(origin.X+i*d.X, origin.Y+i*d.Y)
			ts, err := c.board.Square(tCor)
			if err != nil {
				break
			}

			if ts == gochess.Empty {
				moves = append(moves, UCI(origin, tCor))
				continue
			}

			if ts&color == gochess.Empty {
				moves = append(moves, UCI(origin, tCor))
				break
			}

			// If the piece is the same color, stop looking in that direction.
			break
		}
	}

	return moves
}

func (c Chess) oneStepPieces(origin gochess.Coordinate, offsets []gochess.Coordinate) []string {
	p, _ := c.board.Square(origin)

	color := p & (gochess.White | gochess.Black)
	moves := make([]string, 0, 8)
	for _, d := range offsets {
		tCor := gochess.Coor(origin.X+d.X, origin.Y+d.Y)
		ts, err := c.board.Square(tCor)
		if err != nil {
			continue
		}

		if ts == gochess.Empty {
			moves = append(moves, UCI(origin, tCor))
			continue
		}

		if ts&color == gochess.Empty {
			moves = append(moves, UCI(origin, tCor))
		}
	}

	return moves
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

	return p == gochess.King|c.turn
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
	return p&^(gochess.White|gochess.Black) == gochess.Pawn
}

// destinationMatch looks for a destination in a list of moves.
// It returns true if any of the moves has the destination.
// The function expects the moves in UCI format.
func destinationMatch(moves []string, destination gochess.Coordinate) bool {
	algCoor := CoordinateToAlgebraic(destination)
	for _, move := range moves {
		dest := move[2:4]
		if dest == algCoor {
			return true
		}
	}

	return false
}

// legalMoves returns the legal moves for the current turn.
func (c *Chess) legalMoves() []string {
	moves := c.availableMoves()

	legalMoves := make([]string, 0, len(moves))
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

// availableMoves returns the available moves for the current turn without checking if they are legal.
func (c *Chess) availableMoves() []string {
	moves := make([]string, 0, 40)
	for x := range 8 {
		for y := range 8 {
			origin := gochess.Coor(x, y)
			piece, _ := c.board.Square(origin)
			if piece&c.turn == gochess.Empty {
				continue
			}

			moves = append(moves, c.movesForPiece(piece, origin)...)
		}
	}

	return moves
}

// isLegalMove is a helper function that verifies if the move is legal.
//
// It verifies it making the move in a temporary board and checking if the
// king is in check or the king way is under attack in castling moves.
func (c *Chess) isLegalMove(move string) bool {
	kingsColor := c.turn
	c.makeMove(move)

	availableMoves := c.availableMoves()
	kingPosition := c.kingsPosition(kingsColor)

	kingUnderAttack := destinationMatch(availableMoves, kingPosition)
	c.unmakeMove()

	// If the king is under attack, the move is not legal.
	if kingUnderAttack {
		return false
	}

	// If the move is a castle and the king way is under attack, the move is not legal.
	if c.isCastleMove(move) && destinationMatch(availableMoves, castleKingWay[move]) {
		return false
	}

	return true
}
