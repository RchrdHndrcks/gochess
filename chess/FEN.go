package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/RchrdHndrcks/gochess/v2"
)

// loadPosition is a helper function that loads a board from a FEN string.
//
// The function will read the entire FEN string and will return an error if
// the FEN string is invalid.
//
// The board and properties will not be modified if the FEN string is invalid.
func (c *Chess) loadPosition(FEN string) error {
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
	var whiteKingPosition, blackKingPosition *gochess.Coordinate

	brd := make([][]gochess.Piece, 8)
	for y := range 8 {
		row := make([]gochess.Piece, 8)

		if len(fenRows[y]) == 0 || len(fenRows[y]) > 8 {
			return fmt.Errorf("invalid FEN: %s", FEN)
		}

		for x := range 8 {
			char := string(fenRows[y][0])
			fenRows[y] = fenRows[y][1:]

			n, err := strconv.Atoi(char)
			// If c is not a number, it's a piece.
			if err != nil {
				p, ok := gochess.Pieces[char]
				if !ok {
					return fmt.Errorf("invalid FEN: unknown piece character: %s", char)
				}
				row[x] = p
				coor := gochess.Coor(x, y)
				if p == gochess.King|gochess.White {
					whiteKing++
					whiteKingPosition = &coor
				}
				if p == gochess.King|gochess.Black {
					blackKing++
					blackKingPosition = &coor
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
	b, _ := gochess.NewBoard(8, brd...)
	c.board = b

	// If the FEN is invalid, setProperties will
	// return an error without modifying the board or the properties.
	if err := c.setProperties(FEN); err != nil {
		*c = copy
		return fmt.Errorf("invalid FEN: %w", err)
	}

	legacyWhiteKingPosition := c.whiteKingPosition
	legacyBlackKingPosition := c.blackKingPosition

	c.whiteKingPosition = whiteKingPosition
	c.blackKingPosition = blackKingPosition

	if !c.isPositionLegal() {
		c.whiteKingPosition = legacyWhiteKingPosition
		c.blackKingPosition = legacyBlackKingPosition
		*c = copy
		return errors.New("invalid FEN: the current turn can capture the opponent king")
	}

	return nil
}

// calculateFEN returns the FEN string of the current position.
//
// If move is not empty, it will only update the specified move.
// If more than one move is passed, it will update only the first move.
func (c *Chess) calculateFEN(move ...string) string {
	if c.blackKingPosition == nil || c.whiteKingPosition == nil {
		return ""
	}

	ac := c.availableCastles.String()
	ips := "-"
	if c.enPassantFile >= 0 {
		// EP rank: rank 6 (index 2) when black to move, rank 3 (index 5) when white to move.
		// After a white double push, side to move is black -> EP square is on rank 3 (y=5).
		// After a black double push, side to move is white -> EP square is on rank 6 (y=2).
		epY := 5
		if c.turn == gochess.White {
			epY = 2
		}
		ips = CoordinateToAlgebraic(gochess.Coor(int(c.enPassantFile), epY))
	}

	var boardFEN string
	if len(move) == 0 {
		boardFEN = c.calculateEntireBoardFEN()
	} else {
		m := move[0]
		origin, _ := AlgebraicToCoordinate(m[:2])
		target, _ := AlgebraicToCoordinate(m[2:4])
		boardFEN = c.calculateBoardFEN(origin.Y, target.Y)
	}

	return boardFEN + fmt.Sprintf(" %s %s %s %d %d", gochess.ColorNames[c.turn], ac, ips, c.halfMoves, c.movesCount)
}

// calculateBoardFEN returns the FEN string of the board without the properties.
// If rows is not empty, it will only update the specified rows.
func (c *Chess) calculateBoardFEN(rows ...int) string {
	if len(rows) == 0 {
		return c.calculateEntireBoardFEN()
	}

	r := strings.Split(strings.Split(c.actualFEN, " ")[0], "/")
	for _, row := range rows {
		r[row] = c.calculateRowFEN(row)
	}

	return strings.Join(r, "/")
}

func (c *Chess) calculateEntireBoardFEN() string {
	fen := ""

	for y := range c.board.Width() {
		fen += c.calculateRowFEN(y)
		if y < 7 {
			fen += "/"
		}
	}

	return fen
}

func (c *Chess) calculateRowFEN(y int) string {
	fen := ""
	empty := 0
	for x := range c.board.Width() {
		// Ignore errors since the coordinates are valid.
		piece, _ := c.board.Square(gochess.Coor(x, y))

		if piece == gochess.Empty {
			empty++
			continue
		}

		if empty > 0 {
			fen += fmt.Sprintf("%d", empty)
			empty = 0
		}

		fen += gochess.PieceNames[piece]
	}

	if empty > 0 {
		fen += fmt.Sprintf("%d", empty)
	}

	return fen
}

// setProperties is a helper function that sets the properties of the Chess struct.
// It verifies the properties of the FEN string.
// It returns an error if the FEN string is invalid.
func (c *Chess) setProperties(FEN string) error {
	props := strings.Split(FEN, " ")[1:]

	color, ok := gochess.Colors[props[0]]
	if !ok {
		return fmt.Errorf("invalid color: %s", props[0])
	}

	castleRights, err := parseCastleRights(props[1])
	if err != nil {
		return fmt.Errorf("invalid castles: %s", props[1])
	}

	enPassantSquare := props[2]
	if err := c.validateEnPassant(enPassantSquare); err != nil {
		return fmt.Errorf("invalid en passant square: %s", enPassantSquare)
	}

	enPassantFile := int8(-1)
	if enPassantSquare != "-" {
		coor, _ := AlgebraicToCoordinate(enPassantSquare)
		enPassantFile = int8(coor.X)
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
	c.availableCastles = castleRights
	c.enPassantFile = enPassantFile
	c.halfMoves = halfMoves
	c.movesCount = movesCount
	return nil
}

// updateMovesCount updates the moves count.
func (c *Chess) updateMovesCount() {
	if c.turn == gochess.White {
		c.movesCount++
	}
}

// updateCastlePossibilities checks if the castles are still available.
func (c *Chess) updateCastlePossibilities() {
	k, _ := c.board.Square(gochess.Coor(4, 0))
	rr, _ := c.board.Square(gochess.Coor(7, 0))
	lr, _ := c.board.Square(gochess.Coor(0, 0))
	if rr != gochess.Rook|gochess.Black || k != gochess.King|gochess.Black {
		c.availableCastles &^= BlackKingside
	}
	if lr != gochess.Rook|gochess.Black || k != gochess.King|gochess.Black {
		c.availableCastles &^= BlackQueenside
	}

	K, _ := c.board.Square(gochess.Coor(4, 7))
	rR, _ := c.board.Square(gochess.Coor(7, 7))
	lR, _ := c.board.Square(gochess.Coor(0, 7))
	if rR != gochess.Rook|gochess.White || K != gochess.King|gochess.White {
		c.availableCastles &^= WhiteKingside
	}
	if lR != gochess.Rook|gochess.White || K != gochess.King|gochess.White {
		c.availableCastles &^= WhiteQueenside
	}
}

// updateHalfMoves updates the half moves counter.
//
// It must be called after a move is made. If no move was made
// (i.e. the history is empty), the function returns early without
// modifying the counter.
func (c *Chess) updateHalfMoves() {
	if len(c.history) == 0 {
		return
	}

	c.halfMoves++
	h := c.history[len(c.history)-1]

	// If the move was promotion, reset the counter.
	if len(h.move) == 5 {
		c.halfMoves = 0
		return
	}

	// If a capture was made, reset.
	if h.capturedPiece != gochess.Empty {
		c.halfMoves = 0
		return
	}

	// If no capture was made, check whether the last move was a pawn move.
	target := h.move[2:4]
	coor, _ := AlgebraicToCoordinate(target)
	p, _ := c.board.Square(coor)

	if gochess.PieceType(p) == gochess.Pawn {
		c.halfMoves = 0
	}
}

// updateEnPassantSquare updates the en passant square.
//
// It must be called after a move is made. If no move was made,
// the function will panic.
func (c *Chess) updateEnPassantSquare() {
	c.enPassantFile = -1

	lastMove := c.history[len(c.history)-1].move
	if len(lastMove) != 4 {
		return
	}

	dest, _ := AlgebraicToCoordinate(lastMove[2:])
	p, _ := c.board.Square(dest)

	if gochess.PieceType(p) != gochess.Pawn {
		return
	}

	destRow, _ := strconv.Atoi(lastMove[3:4])
	orgRow, _ := strconv.Atoi(lastMove[1:2])
	if destRow == orgRow+2 || destRow == orgRow-2 {
		c.enPassantFile = int8(dest.X)
		return
	}
}

// validateEnPassant validates the in passant square.
func (c Chess) validateEnPassant(square string) error {
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
	if coor.Y == 2 {
		yCoor = 3
	}

	auxCoor := gochess.Coor(coor.X, yCoor)
	p, _ := c.board.Square(auxCoor)
	if gochess.PieceType(p) != gochess.Pawn {
		return errors.New("invalid in passant square")
	}

	return nil
}

// parseCastleRights parses the FEN castles field into a CastleRights bitmask.
func parseCastleRights(castles string) (CastleRights, error) {
	if castles == "-" {
		return NoCastling, nil
	}

	bits := map[rune]CastleRights{
		'K': WhiteKingside,
		'Q': WhiteQueenside,
		'k': BlackKingside,
		'q': BlackQueenside,
	}

	var cr CastleRights
	for _, ch := range castles {
		bit, ok := bits[ch]
		if !ok {
			return NoCastling, errors.New("invalid castles")
		}
		if cr.Has(bit) {
			return NoCastling, errors.New("duplicate castle character")
		}
		cr |= bit
	}

	return cr, nil
}

// isPositionLegal verifies if the current turn can capture the opponent king.
//
// If the current turn can capture the opponent king, the position is not legal
// and the function returns false.
func (c Chess) isPositionLegal() bool {
	c.toggleColor()
	defer c.toggleColor()
	return !c.isCheck()
}

func (c *Chess) toggleColor() {
	if c.turn == gochess.White {
		c.turn = gochess.Black
		return
	}

	c.turn = gochess.White
}

// isCheck is the helper function that checks if the current turn is in check.
func (c Chess) isCheck() bool {
	if c.blackKingPosition == nil || c.whiteKingPosition == nil {
		return false
	}

	kingPosition := c.kingsPosition(c.turn)
	return c.IsAttacked(kingPosition, opponentColor(c.turn))
}

// kingsPosition returns the position of the king of the given color.
func (c Chess) kingsPosition(color gochess.Piece) gochess.Coordinate {
	if color == gochess.White {
		return *c.whiteKingPosition
	}

	return *c.blackKingPosition
}

// pieceFromFEN is a helper function that returns the piece at the given coordinate
// in the FEN string.
//
// It must be called with a valid coordinate and a valid FEN string.
// If there is no piece at the given coordinate, it returns gochess.Empty.
func pieceFromFEN(fen string, coord gochess.Coordinate) gochess.Piece {
	fenRow := strings.Split(strings.Split(fen, " ")[0], "/")[coord.Y]
	var count int
	for _, c := range fenRow {
		n, err := strconv.Atoi(string(c))
		if err == nil {
			count += n
			if count > coord.X {
				return gochess.Empty
			}

			continue
		}

		if count == coord.X {
			if p, ok := gochess.Pieces[string(c)]; ok {
				return p
			}
			return gochess.Empty
		}

		count++
	}

	return gochess.Empty
}
