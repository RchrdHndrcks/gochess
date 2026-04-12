package chess

import (
	"sync"

	"github.com/RchrdHndrcks/gochess/v2"
)

// capacityByPiece is a helper map where the key is the piece and the value is
// the max moves count for that piece. It is useful to know which capacity asign
// to the slice where the moves of these pieces are stored.
var capacityByPiece = map[gochess.Piece]int{
	gochess.White | gochess.Queen:  27,
	gochess.Black | gochess.Queen:  27,
	gochess.White | gochess.Rook:   14,
	gochess.Black | gochess.Rook:   14,
	gochess.White | gochess.Bishop: 13,
	gochess.Black | gochess.Bishop: 13,
}

// moveData holds pre-parsed move information for the internal applyMove
// method. Both makeMove (UCI string entry point) and MakeMoveCompact build
// this struct so that the actual board mutation is string-free.
type moveData struct {
	from          gochess.Coordinate
	to            gochess.Coordinate
	capturedPiece gochess.Piece // piece captured (with color); Empty if none
	promotionType gochess.Piece // promotion piece type without color; Empty if none
	isCastle      bool
	isEnPassant   bool
	uci           string // UCI string stored in history for unmakeMove
}

// makeMove makes a move without checking if it is legal.
func (c *Chess) makeMove(move string) {
	o, _ := AlgebraicToCoordinate(move[:2])
	t, _ := AlgebraicToCoordinate(move[2:4])

	md := moveData{
		from:        o,
		to:          t,
		isCastle:    c.isCastleMove(move),
		isEnPassant: c.isEnPassantMove(move),
		uci:         move,
	}

	if md.isEnPassant {
		md.capturedPiece = gochess.Pawn | opponentColor(c.turn)
	} else {
		dst, _ := c.board.Square(t)
		md.capturedPiece = dst
	}

	if len(move) == 5 {
		md.promotionType = gochess.PiecesWithoutColor[move[4:5]]
	}

	c.applyMove(md)
}

// applyMove mutates the board for a pre-parsed move, updates the history,
// king positions, side-to-move and derived counters. It is shared between
// makeMove (string entry point) and MakeMoveCompact.
func (c *Chess) applyMove(md moveData) {
	o, t := md.from, md.to

	// Capture the moving piece's type before mutating the board so that
	// the incremental piece-list updates below can reference it.
	movingPiece, _ := c.board.Square(o)
	movingType := gochess.PieceType(movingPiece)
	moverColor := c.turn
	oppColor := opponentColor(moverColor)

	if md.isCastle {
		// Move the rook for the castle.
		c.makeMoveOnBoard(castleRook[md.uci], gochess.Coor((o.X+t.X)/2, o.Y))
	}

	if md.isEnPassant {
		// Remove the captured pawn, which lives behind the target square.
		_ = c.board.SetSquare(gochess.Coor(t.X, o.Y), gochess.Empty)
	}

	if md.promotionType != gochess.Empty {
		_ = c.board.SetSquare(t, md.promotionType|c.turn)
		_ = c.board.SetSquare(o, gochess.Empty)
	} else {
		c.makeMoveOnBoard(o, t)
	}

	// Incremental piece-list updates. These mirror the board mutations
	// above and avoid the cost of rescanning the board on every move.
	switch {
	case md.isCastle:
		c.pieceLists[colorIndex(moverColor)][gochess.King].move(o, t)
		rookFrom := castleRook[md.uci]
		rookTo := gochess.Coor((o.X+t.X)/2, o.Y)
		c.pieceLists[colorIndex(moverColor)][gochess.Rook].move(rookFrom, rookTo)
	case md.isEnPassant:
		c.pieceLists[colorIndex(moverColor)][gochess.Pawn].move(o, t)
		c.pieceLists[colorIndex(oppColor)][gochess.Pawn].remove(gochess.Coor(t.X, o.Y))
	case md.promotionType != gochess.Empty:
		c.pieceLists[colorIndex(moverColor)][gochess.Pawn].remove(o)
		if md.capturedPiece != gochess.Empty {
			c.pieceLists[colorIndex(oppColor)][gochess.PieceType(md.capturedPiece)].remove(t)
		}
		c.pieceLists[colorIndex(moverColor)][md.promotionType].add(t)
	default:
		if md.capturedPiece != gochess.Empty {
			c.pieceLists[colorIndex(oppColor)][gochess.PieceType(md.capturedPiece)].remove(t)
		}
		c.pieceLists[colorIndex(moverColor)][movingType].move(o, t)
	}

	c.history = append(
		c.history,
		chessContext{
			move:              md.uci,
			compactMove:       NullMove,
			capturedPiece:     md.capturedPiece,
			positionKey:       positionKey(c.actualFEN),
			halfMove:          c.halfMoves,
			availableCastles:  c.availableCastles,
			enPassantFile:     c.enPassantFile,
			whiteKingPosition: c.whiteKingPosition,
			blackKingPosition: c.blackKingPosition,
			check:             c.check,
			checkmate:         c.checkmate,
			stalemate:         c.stalemate,
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

// opponentColor returns the opposing color piece bit.
func opponentColor(color gochess.Piece) gochess.Piece {
	if color == gochess.White {
		return gochess.Black
	}
	return gochess.White
}

// makeMoveOnBoard is a helper function to make a move on the board.
//
// It must be used only when the move is already validated.
func (c *Chess) makeMoveOnBoard(origin, target gochess.Coordinate) {
	p, _ := c.board.Square(origin)
	_ = c.board.SetSquare(origin, gochess.Empty)
	_ = c.board.SetSquare(target, p)
}

// unmakeMove is a helper function to unmake the last move.
//
// Until v1.1.0 unmakeMove was just reloading the last FEN from the history.
// Now it makes the move back to the board. This change was made because
// makeMove modifies the current board and unmakeMove was just modifying
// the board reference.
func (c *Chess) unmakeMove() {
	if len(c.history) == 0 {
		return
	}

	lastContext := c.history[len(c.history)-1]
	c.history = c.history[:len(c.history)-1]

	c.halfMoves = lastContext.halfMove
	c.availableCastles = lastContext.availableCastles
	c.enPassantFile = lastContext.enPassantFile
	c.whiteKingPosition = lastContext.whiteKingPosition
	c.blackKingPosition = lastContext.blackKingPosition
	c.check = lastContext.check
	c.checkmate = lastContext.checkmate
	c.stalemate = lastContext.stalemate

	c.toggleColor()

	if c.turn == gochess.Black {
		c.movesCount--
	}

	move := lastContext.move
	// Note: o = original target, t = original origin (variables are swapped
	// here so that "moving back" reads naturally as o -> t).
	o, _ := AlgebraicToCoordinate(move[2:4])
	t, _ := AlgebraicToCoordinate(move[:2])

	if len(move) == 5 {
		// Promotion: restore the pawn on the origin square.
		_ = c.board.SetSquare(t, gochess.Pawn|c.turn)
		// Restore whatever lived on the target square: a captured piece for
		// a promotion-capture, or Empty for a quiet promotion. Without
		// clearing here, a quiet promotion would leave the promoted piece
		// on the destination square.
		if lastContext.capturedPiece != gochess.Empty {
			_ = c.board.SetSquare(o, lastContext.capturedPiece)
		} else {
			_ = c.board.SetSquare(o, gochess.Empty)
		}
	} else {
		// Move the piece back to its original square.
		c.makeMoveOnBoard(o, t)
		// Restore the captured piece on the destination, if any. For en
		// passant the captured pawn lives on a different square (handled
		// below), so leave the diagonal destination empty in that case.
		isEP := isEnPassantMoveByContext(move, lastContext)
		if lastContext.capturedPiece != gochess.Empty && !isEP {
			_ = c.board.SetSquare(o, lastContext.capturedPiece)
		}
	}

	if c.isCastleMove(move) {
		c.makeMoveOnBoard(gochess.Coor((o.X+t.X)/2, o.Y), castleRook[move])
	}

	if isEnPassantMoveByContext(move, lastContext) {
		// The captured pawn belongs to the opponent of the side that made
		// the EP capture. After toggleColor, c.turn IS the side that made
		// the capture, so the captured pawn color is the opposite.
		_ = c.board.SetSquare(gochess.Coor(o.X, t.Y), gochess.Pawn|opponentColor(c.turn))
	}

	// Incremental piece-list updates: reverse the operations performed by
	// applyMove. c.turn is now the side that made the move (post toggle).
	moverColor := c.turn
	oppColor := opponentColor(moverColor)
	isEP := isEnPassantMoveByContext(move, lastContext)
	isCastle := c.isCastleMove(move)
	switch {
	case isCastle:
		rookFrom := castleRook[move]
		rookTo := gochess.Coor((o.X+t.X)/2, o.Y)
		c.pieceLists[colorIndex(moverColor)][gochess.Rook].move(rookTo, rookFrom)
		c.pieceLists[colorIndex(moverColor)][gochess.King].move(o, t)
	case isEP:
		c.pieceLists[colorIndex(oppColor)][gochess.Pawn].add(gochess.Coor(o.X, t.Y))
		c.pieceLists[colorIndex(moverColor)][gochess.Pawn].move(o, t)
	case len(move) == 5:
		// Promotion (with or without capture).
		promoType := gochess.PiecesWithoutColor[move[4:5]]
		c.pieceLists[colorIndex(moverColor)][promoType].remove(o)
		if lastContext.capturedPiece != gochess.Empty {
			c.pieceLists[colorIndex(oppColor)][gochess.PieceType(lastContext.capturedPiece)].add(o)
		}
		c.pieceLists[colorIndex(moverColor)][gochess.Pawn].add(t)
	default:
		// Regular move (possibly capture). Determine moving type from board:
		// after the board restore, the piece is back on its origin square (t).
		movedPiece, _ := c.board.Square(t)
		movedType := gochess.PieceType(movedPiece)
		c.pieceLists[colorIndex(moverColor)][movedType].move(o, t)
		if lastContext.capturedPiece != gochess.Empty {
			c.pieceLists[colorIndex(oppColor)][gochess.PieceType(lastContext.capturedPiece)].add(o)
		}
	}

	c.actualFEN = c.calculateFEN()
}

// isEnPassantMoveByContext reports whether the just-undone move was an
// en passant capture. We can't use Chess.isEnPassantMove because the en
// passant file has already been restored to the pre-move state and the
// destination square no longer holds the moving pawn (we just moved it
// back). Instead, we look at the captured piece slot in the history: an
// EP capture is the only pawn-target move where the captured pawn lives
// on a different square than the move's destination.
func isEnPassantMoveByContext(move string, ctx chessContext) bool {
	if len(move) != 4 {
		return false
	}
	if gochess.PieceType(ctx.capturedPiece) != gochess.Pawn {
		return false
	}
	// The destination file differs from the origin file (diagonal pawn move)
	// AND the captured piece slot is set, but the target square in the
	// pre-move position was empty. The cheapest detector: the move is a
	// pawn diagonal whose target square's pre-move occupant was empty.
	// We approximate this via the EP file stored in the restored context:
	// if the destination square equals the EP target, this is EP.
	if ctx.enPassantFile < 0 {
		return false
	}
	dest, _ := AlgebraicToCoordinate(move[2:4])
	if int8(dest.X) != ctx.enPassantFile {
		return false
	}
	// EP target rank: y=2 if white captured (target rank 6), y=5 if black.
	return dest.Y == 2 || dest.Y == 5
}

// movesForPiece returns the available moves for a piece.
//
// The function returns a slice of UCI moves.
// (e.g. "e2e4" for moving the piece at e2 to e4.)
// Disclaimer: This function does not check if the move is legal for a Chess game.
func (c Chess) movesForPiece(piece gochess.Piece, origin gochess.Coordinate) []string {
	switch gochess.PieceType(piece) {
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
		captureMoves := c.pawnCaptureMoves(origin, true)
		// Only add push-promotion moves if the target square is empty.
		// Previously this called promotionPossibilities unconditionally, which
		// generated pseudo-legal captures to the promotion square even when it
		// was occupied (e.g. by the opponent's king), incorrectly marking that
		// square as attacked.
		if s == gochess.Empty {
			captureMoves = append(captureMoves, c.promotionPossibilities(origin, tCor)...)
		}
		return captureMoves
	}

	if !(dir == 1 && origin.Y == 1) && !(dir == -1 && origin.Y == 6) {
		return append(c.pawnCaptureMoves(origin, false), moves...)
	}

	// Double push is only allowed when the single-push square is empty, i.e.
	// the previous step appended a move. Otherwise the pawn cannot leap over
	// an occupied intermediate square.
	if len(moves) == 0 {
		return c.pawnCaptureMoves(origin, false)
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
	pColor := gochess.PieceColor(p)
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

		if c.enPassantFile >= 0 && int8(tCor.X) == c.enPassantFile && tCor.Y == expectedEPRank(pColor) {
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

		moves = append(moves, c.promotionPossibilities(origin, tCor)...)
	}

	return moves
}

// promotionPossibilities is a helper function that returns the UCI moves with
// the value of the piece to be promoted.
func (c Chess) promotionPossibilities(origin, target gochess.Coordinate) []string {
	moves := make([]string, 4)
	for i, p := range []gochess.Piece{gochess.Queen, gochess.Rook, gochess.Bishop, gochess.Knight} {
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
	if c.availableCastles == NoCastling {
		return nil
	}

	p, _ := c.board.Square(origin)
	kingColor := gochess.PieceColor(p)

	type castleOption struct {
		right CastleRights
		color gochess.Piece
		dir   int
	}

	options := []castleOption{
		{WhiteKingside, gochess.White, 1},
		{WhiteQueenside, gochess.White, -1},
		{BlackKingside, gochess.Black, 1},
		{BlackQueenside, gochess.Black, -1},
	}

	moves := make([]string, 0, 2)
	for _, opt := range options {
		if !c.availableCastles.Has(opt.right) {
			continue
		}
		if opt.color != kingColor {
			continue
		}

		ts, err := c.board.Square(gochess.Coor(origin.X+opt.dir, origin.Y))
		if err != nil || ts != gochess.Empty {
			continue
		}

		ts, err = c.board.Square(gochess.Coor(origin.X+2*opt.dir, origin.Y))
		if err != nil || ts != gochess.Empty {
			continue
		}

		// For queenside castling, the b-file square (3 squares from king) must
		// also be empty even though the king does not pass through it.
		if opt.dir == -1 {
			ts, err = c.board.Square(gochess.Coor(origin.X+3*opt.dir, origin.Y))
			if err != nil || ts != gochess.Empty {
				continue
			}
		}

		moves = append(moves, UCI(origin, gochess.Coor(origin.X+2*opt.dir, origin.Y)))
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

	color := gochess.PieceColor(p)
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

	color := gochess.PieceColor(p)
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
	if c.enPassantFile < 0 {
		return false
	}

	origin, _ := AlgebraicToCoordinate(move[:2])
	target, _ := AlgebraicToCoordinate(move[2:4])

	p, _ := c.board.Square(origin)
	if gochess.PieceType(p) != gochess.Pawn {
		return false
	}

	if int8(target.X) != c.enPassantFile {
		return false
	}

	return target.Y == expectedEPRank(gochess.PieceColor(p))
}

// expectedEPRank returns the y-coordinate of the en passant target square for
// a pawn of the given color performing the capture.
//
// White pawns capture en passant onto rank 6 (y=2); black pawns onto rank 3 (y=5).
func expectedEPRank(capturingColor gochess.Piece) int {
	if capturingColor == gochess.White {
		return 2
	}
	return 5
}

// legalMoves returns the legal moves for the current turn.
func (c Chess) legalMoves() []string {
	moves := c.availableMoves()
	legalMoves := make([]string, 0, len(moves))

	goroutinesCount := c.config.Parallelism
	_, ok := c.board.(Cloner)
	if !ok || goroutinesCount <= 1 {
		return c.calculateLegalMovesSecuentially(moves)
	}

	wg := &sync.WaitGroup{}
	availableMovesChan := make(chan string, goroutinesCount)
	legalMovesChan := make(chan string, len(moves))
	wg.Add(goroutinesCount)
	for range goroutinesCount {
		go func() {
			defer wg.Done()
			copy := c.clone()

			for move := range availableMovesChan {
				if copy.isLegalMove(move) {
					legalMovesChan <- move
				}
			}
		}()
	}

	for _, move := range moves {
		availableMovesChan <- move
	}

	close(availableMovesChan)
	wg.Wait()
	close(legalMovesChan)

	for move := range legalMovesChan {
		legalMoves = append(legalMoves, move)
	}

	return legalMoves
}

func (c Chess) calculateLegalMovesSecuentially(moves []string) []string {
	legalMoves := make([]string, 0, len(moves))
	for _, move := range moves {
		if c.isLegalMove(move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

// availableMoves returns the available moves for the current turn without checking if they are legal.
func (c Chess) availableMoves() []string {
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
func (c Chess) isLegalMove(move string) bool {
	kingsColor := c.turn

	c.makeMove(move)
	kingPosition := c.kingsPosition(kingsColor)
	opponent := gochess.Black
	if kingsColor == gochess.Black {
		opponent = gochess.White
	}
	kingUnderAttack := c.IsAttacked(kingPosition, opponent)
	c.unmakeMove()

	// If the king is under attack, the move is not legal.
	if kingUnderAttack {
		return false
	}

	// FIDE rule 3.8.2: castling has three restrictions on attacked squares.
	if c.isCastleMove(move) {
		// (1) Cannot castle while in check. isCheck() is called on the restored
		//     pre-castle position, so the king is still on its starting square and
		//     pawn attacks to that square are generated correctly.
		if c.isCheck() {
			return false
		}
		// (2) Cannot castle through check (king passage square under attack).
		if c.IsAttacked(castleKingWay[move], opponent) {
			return false
		}
	}

	return true
}
