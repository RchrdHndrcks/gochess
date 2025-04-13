package pkg

import (
	"fmt"
)

// Board is a 2D array of pieces.
type Board struct {
	squares [][]int8
	width   int
}

// NewBoard creates a new board.
//
// It receives the width of the board and an optional 2D array of pieces.
// If no 2D array is provided, the board will be initialized with empty squares.
//
// It returns an error if the width is less than 1 or if the squares are not
// valid. Squares could be invalid if:
// - The width of the squares is different from the width of the board.
// - The length of the squares array is different from the width of the board.
// - The length of the inner arrays is different from the width of the board.
func NewBoard(width int, squares ...[]int8) (*Board, error) {
	if width < 1 {
		return nil, fmt.Errorf("board: invalid width: %d", width)
	}

	if len(squares) != 0 {
		if len(squares) != width {
			return nil, fmt.Errorf("board: invalid squares length: %d", len(squares))
		}

		for _, row := range squares {
			if len(row) != width {
				return nil, fmt.Errorf("board: invalid row length: %d", len(row))
			}
		}

		return &Board{
			squares: squares,
			width:   width,
		}, nil
	}

	s := make([][]int8, width)
	for i := range width {
		s[i] = make([]int8, width)
	}

	return &Board{
		squares: s,
		width:   width,
	}, nil
}

// Width returns the width of the board.
func (b *Board) Width() int {
	return b.width
}

// Square returns the piece at the given Coordinate.
//
// It returns an error if the Coordinate is out of bounds.
func (b *Board) Square(c Coordinate) (int8, error) {
	if !b.isValidCoordinate(c) {
		return Empty, fmt.Errorf("board: invalid coordinate: %v", c)
	}

	return b.squares[c.Y][c.X], nil
}

// MakeMove makes a move on the board.
//
// It doesn't make any validation on the move.
// It is the caller's responsibility to make sure the move is valid, including
// if the origin square is not empty, if the target square is empty or if the
// move is valid for the piece.
//
// If the coordinate is out of bounds, it returns an error.
func (b *Board) MakeMove(origin, target Coordinate) error {
	if !b.isValidCoordinate(origin) || !b.isValidCoordinate(target) {
		return fmt.Errorf("invalid coordinates")
	}

	b.squares[target.Y][target.X] = b.squares[origin.Y][origin.X]
	b.squares[origin.Y][origin.X] = Empty
	return nil
}

func (b *Board) SetSquare(c Coordinate, p int8) error {
	if !b.isValidCoordinate(c) {
		return fmt.Errorf("invalid coordinate")
	}

	b.squares[c.Y][c.X] = p
	return nil
}

// isValidCoordinate returns true if the Coordinate is within the board bounds.
func (b *Board) isValidCoordinate(c Coordinate) bool {
	return c.X >= 0 && c.X < b.width && c.Y >= 0 && c.Y < b.width
}
