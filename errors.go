package gochess

import "errors"

// ErrInvalidWidth is returned when the width is less than 1.
var ErrInvalidWidth = errors.New("invalid width")

// ErrInvalidSquare is returned when a square is invalid.
var ErrInvalidSquare = errors.New("invalid square")

// ErrInvalidCoordinate is returned when a coordinate is out of bounds.
var ErrInvalidCoordinate = errors.New("invalid coordinate")
