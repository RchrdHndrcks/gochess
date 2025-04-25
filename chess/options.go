package chess

import (
	"fmt"
)

// Option is a function that configures a chess.
type Option func(*Chess) error

const (
	defaultFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

// WithBoard sets the board of the chess.
// If the board is nil, it returns an error.
// If you want to use this option, it must be the first one.
func WithBoard(b Board) Option {
	return func(c *Chess) error {
		c.board = b
		return nil
	}
}

// WithFEN sets the FEN of the chess.
// If the FEN is invalid, it returns an error.
// If you try to set the FEN before the board, it will set the default board.
func WithFEN(FEN string) Option {
	return func(c *Chess) error {
		if err := c.LoadPosition(FEN); err != nil {
			return fmt.Errorf("failed to load position: %w", err)
		}

		return nil
	}
}
