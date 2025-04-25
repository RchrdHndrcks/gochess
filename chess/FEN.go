package chess

import (
	"cmp"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/RchrdHndrcks/gochess"
)

var fenAnalysisRegex = regexp.MustCompile("[/0-9]")

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

	brd := make([][]int8, 8, 8)
	for y := range 8 {
		row := make([]int8, 8, 8)

		if len(fenRows[y]) == 0 || len(fenRows[y]) > 8 {
			return fmt.Errorf("invalid FEN: %s", FEN)
		}

		for x := range 8 {
			char := string(fenRows[y][0])
			fenRows[y] = fenRows[y][1:]

			n, err := strconv.Atoi(char)
			// If c is not a number, it's a piece.
			if err != nil {
				p := gochess.Pieces[char]
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

	ac := cmp.Or(c.availableCastles, "-")
	ips := cmp.Or(c.enPassantSquare, "-")

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

	availableCastles := props[1]
	if err := c.validateCastles(availableCastles); err != nil {
		return fmt.Errorf("invalid castles: %s", availableCastles)
	}

	enPassantSquare := props[2]
	if err := c.validateEnPassant(enPassantSquare); err != nil {
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

// updateMovesCount updates the moves count.
func (c *Chess) updateMovesCount() {
	if c.turn == gochess.White {
		c.movesCount++
	}
}

// updateCastlePossibilities checks if the castles are still available.
func (c *Chess) updateCastlePossibilities() {
	toBeRemoved := map[string]bool{}

	k, _ := c.board.Square(gochess.Coor(4, 0))
	rr, _ := c.board.Square(gochess.Coor(7, 0))
	lr, _ := c.board.Square(gochess.Coor(0, 0))
	toBeRemoved["k"] = rr != gochess.Rook|gochess.Black || k != gochess.King|gochess.Black
	toBeRemoved["q"] = lr != gochess.Rook|gochess.Black || k != gochess.King|gochess.Black

	K, _ := c.board.Square(gochess.Coor(4, 7))
	rR, _ := c.board.Square(gochess.Coor(7, 7))
	lR, _ := c.board.Square(gochess.Coor(0, 7))
	toBeRemoved["K"] = rR != gochess.Rook|gochess.White || K != gochess.King|gochess.White
	toBeRemoved["Q"] = lR != gochess.Rook|gochess.White || K != gochess.King|gochess.White

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
	h := c.history[len(c.history)-1]

	// If the move was promotion, reset the counter.
	if len(h.move) == 5 {
		c.halfMoves = 0
		return
	}

	// Look for a change in the board.
	// If we have less pieces than before, a capture was made so we reset the counter.
	lastFENPiecePart := strings.Split(h.fen, " ")[0]

	lastFENPiecePart = fenAnalysisRegex.ReplaceAllString(lastFENPiecePart, "")
	fenPiecePart := fenAnalysisRegex.ReplaceAllString(c.calculateBoardFEN(), "")

	if len(lastFENPiecePart) > len(fenPiecePart) {
		c.halfMoves = 0
		return
	}

	// If no capture was made, we check if last move was a pawn move.
	target := h.move[2:4]
	coor, _ := AlgebraicToCoordinate(target)
	p, _ := c.board.Square(coor)

	piece := p &^ (gochess.White | gochess.Black)
	if piece == gochess.Pawn {
		c.halfMoves = 0
	}
}

// updateEnPassantSquare updates the en passant square.
//
// It must be called after a move is made. If no move was made,
// the function will panic.
func (c *Chess) updateEnPassantSquare() {
	c.enPassantSquare = ""

	lastMove := c.history[len(c.history)-1].move
	if len(lastMove) != 4 {
		return
	}

	dest, _ := AlgebraicToCoordinate(lastMove[2:])
	p, _ := c.board.Square(dest)

	if p&^(gochess.White|gochess.Black) != gochess.Pawn {
		return
	}

	destRow, _ := strconv.Atoi(lastMove[3:4])
	orgRow, _ := strconv.Atoi(lastMove[1:2])
	if destRow == orgRow+2 || destRow == orgRow-2 {
		c.enPassantSquare = fmt.Sprintf("%s%d", lastMove[2:3], (destRow+orgRow)/2)
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
	color := gochess.White
	if coor.Y == 2 {
		yCoor = 3
		color = gochess.Black
	}

	auxCoor := gochess.Coor(coor.X, yCoor)
	p, _ := c.board.Square(auxCoor)
	if p&^color != gochess.Pawn {
		return errors.New("invalid in passant square")
	}

	return nil
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

func (c Chess) isCheck() bool {
	kingPosition := c.kingsPosition(c.turn)

	c.toggleColor()
	defer c.toggleColor()
	moves := c.availableMoves()

	return destinationMatch(moves, kingPosition)
}

// kingsPosition returns the position of the king of the given color.
func (c Chess) kingsPosition(color int8) gochess.Coordinate {
	if color == gochess.White {
		return *c.whiteKingPosition
	}

	return *c.blackKingPosition
}
