package pkg

import "fmt"

// coordinate represents a 2D coordinate.
type Coordinate struct {
	x, y int
}

// Coordinate returns a new coordinate.
func Coor(x, y int) Coordinate {
	return Coordinate{x, y}
}

// IndexCoor returns the index of a Coordinate.
// Examples:
// (0, 1) would return 8.
// (1, 0) would return 1.
// If the Coordinate is out of bounds, error is returned.
func IndexCoor(c Coordinate) (int, error) {
	if c.x > 7 || c.y > 7 || c.x < 0 || c.y < 0 {
		return 0, fmt.Errorf("Coordinate out of bounds")
	}

	return c.y*8 + c.x, nil
}

// CoordinateToAlgebraic returns a new algebraic notation from a Coordinate.
// For example, (0, 0) would return "a1".
// If the Coordinate is out of bounds, an empty string is returned.
func CoordinateToAlgebraic(c Coordinate) string {
	if c.x > 7 || c.y > 7 || c.x < 0 || c.y < 0 {
		return ""
	}

	return fmt.Sprintf("%c%d", 'a'+c.x, 8-c.y)
}

// AlgebraicToCoordinate returns a new Coordinate from an algebraic notation.
// For example, "a1" would return (0, 0).
// If the algebraic notation is invalid, an empty Coordinate is returned.
func AlgebraicToCoordinate(s string) (Coordinate, error) {
	if len(s) != 2 {
		return Coordinate{}, fmt.Errorf("invalid algebraic notation")
	}

	x := int(s[0] - 'a')
	y := 8 - int(s[1]-'0')
	if x < 0 || x > 7 || y < 0 || y > 7 {
		return Coordinate{}, fmt.Errorf("coordinate out of bounds")
	}

	return Coordinate{x, y}, nil
}
