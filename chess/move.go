package chess

import (
	"fmt"

	"github.com/RchrdHndrcks/gochess"
)

// AlgebraicToCoordinate returns a new Coordinate from text notation.
// For example, "a1" would return (0, 0).
// If the text notation is invalid, an empty Coordinate is returned.
func AlgebraicToCoordinate(s string) (gochess.Coordinate, error) {
	if len(s) != 2 {
		return gochess.Coordinate{}, fmt.Errorf("invalid text notation")
	}

	x := int(s[0] - 'a')
	y := 8 - int(s[1]-'0')
	if x < 0 || x > 7 || y < 0 || y > 7 {
		return gochess.Coordinate{}, fmt.Errorf("coordinate out of bounds")
	}

	return gochess.Coor(x, y), nil
}

// CoordinateToAlgebraic returns a new algebraic notation from a Coordinate.
// For example, (0, 0) would return "a1".
// If the Coordinate is out of bounds, an empty string is returned.
func CoordinateToAlgebraic(c gochess.Coordinate) string {
	if c.X > 7 || c.Y > 7 || c.X < 0 || c.Y < 0 {
		return ""
	}

	return fmt.Sprintf("%c%d", 'a'+c.X, 8-c.Y)
}

// UCI returns the UCI notation of a move.
//
// It receives the origin and target coordinates of the move.
// For example, if the origin is (0, 0) and the target is (0, 1), it would return "a1a2".
//
// If the move is a promotion, it receives the piece to promote to. If it receives more
// than one piece, it returns the first one.
func UCI(origin, target gochess.Coordinate, piece ...int8) string {
	p := ""
	if len(piece) > 0 {
		pi := piece[0]
		// First, we need to uncolor the piece to get the piece name.
		// We do this by doing a bitwise AND with ^White.
		pi &= ^gochess.White

		// The UCI notation for promotion only uses lowercase letters, so we need to
		// convert the piece to lowercase doing a bitwise OR with Black.
		p = gochess.PieceNames[pi|gochess.Black]
	}

	return CoordinateToAlgebraic(origin) + CoordinateToAlgebraic(target) + p
}
