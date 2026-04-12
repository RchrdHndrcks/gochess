package chess

import "github.com/RchrdHndrcks/gochess/v2"

// pieceList tracks up to 10 squares of a given piece type and color.
type pieceList struct {
	squares [10]gochess.Coordinate
	count   int
}

// add appends sq to the piece list. It panics if the list is already full,
// since exceeding the per-piece capacity indicates a malformed position
// (e.g. a FEN with more than 10 pieces of one type) and silently dropping
// the entry would desynchronize the piece lists from the board.
func (pl *pieceList) add(sq gochess.Coordinate) {
	if pl.count >= len(pl.squares) {
		panic("chess: pieceList overflow (more pieces of one type than the board supports)")
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

// colorIndex maps gochess.White → 0, gochess.Black → 1. It panics for
// any other input (gochess.Empty or invalid bit patterns) so that callers
// cannot silently treat an invalid color as Black.
func colorIndex(c gochess.Piece) int {
	switch c {
	case gochess.White:
		return 0
	case gochess.Black:
		return 1
	default:
		panic("chess: invalid piece color passed to colorIndex")
	}
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

// PieceSquares returns all squares occupied by pieces of the given color and
// type. It normalizes pieceType so callers may pass either the bare piece
// type (e.g. gochess.Pawn) or a colored piece (e.g. gochess.White|gochess.Pawn);
// it panics if pieceType does not name a real piece type (1..6) so an out-of-
// range index cannot silently corrupt the result.
func (c *Chess) PieceSquares(color, pieceType gochess.Piece) []gochess.Coordinate {
	pt := gochess.PieceType(pieceType)
	if pt < gochess.Pawn || pt > gochess.King {
		panic("chess: PieceSquares called with invalid piece type")
	}
	return c.pieceLists[colorIndex(color)][pt].slice()
}
