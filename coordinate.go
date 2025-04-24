package gochess

import "fmt"

// coordinate represents a 2D coordinate.
type Coordinate struct {
	X, Y int
}

// String returns the string representation of a Coordinate.
func (c Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

// Coordinate returns a new coordinate.
func Coor(x, y int) Coordinate {
	return Coordinate{x, y}
}
