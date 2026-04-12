package chess

import "github.com/RchrdHndrcks/gochess/v2"

// pieceList tracks up to 10 squares of a given piece type and color.
type pieceList struct {
	squares [10]gochess.Coordinate
	count   int
}

func (pl *pieceList) add(sq gochess.Coordinate) {
	if pl.count >= len(pl.squares) {
		return
	}
	pl.squares[pl.count] = sq
	pl.count++
}

func (pl *pieceList) remove(sq gochess.Coordinate) {
	for i := 0; i < pl.count; i++ {
		if pl.squares[i] == sq {
			pl.count--
			pl.squares[i] = pl.squares[pl.count]
			return
		}
	}
}

func (pl *pieceList) move(from, to gochess.Coordinate) {
	for i := 0; i < pl.count; i++ {
		if pl.squares[i] == from {
			pl.squares[i] = to
			return
		}
	}
}

func (pl *pieceList) slice() []gochess.Coordinate {
	result := make([]gochess.Coordinate, pl.count)
	copy(result, pl.squares[:pl.count])
	return result
}

// colorIndex maps gochess.White → 0, gochess.Black → 1.
func colorIndex(c gochess.Piece) int {
	if c == gochess.White {
		return 0
	}
	return 1
}

// initPieceLists resets and repopulates all piece lists from the current board.
// Called once during construction (New/LoadPosition), NOT on every move.
func (c *Chess) initPieceLists() {
	for i := range c.pieceLists {
		for j := range c.pieceLists[i] {
			c.pieceLists[i][j].count = 0
		}
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			sq := gochess.Coor(x, y)
			piece, _ := c.board.Square(sq)
			if piece == gochess.Empty {
				continue
			}
			color := gochess.PieceColor(piece)
			ptype := gochess.PieceType(piece)
			c.pieceLists[colorIndex(color)][ptype].add(sq)
		}
	}
}

// PieceSquares returns all squares occupied by pieces of the given color and type.
func (c *Chess) PieceSquares(color, pieceType gochess.Piece) []gochess.Coordinate {
	return c.pieceLists[colorIndex(color)][pieceType].slice()
}
