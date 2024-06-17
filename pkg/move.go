package pkg

import "fmt"

// UCI returns the UCI notation of a move.
func UCI(oCor, tCor Coordinate) string {
	return fmt.Sprintf("%s%s", CoordinateToAlgebraic(oCor), CoordinateToAlgebraic(tCor))
}
