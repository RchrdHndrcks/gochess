package chess

import (
	"fmt"

	"github.com/RchrdHndrcks/gochess"
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
		if c.board == nil {
			b, _ := gochess.NewBoard(8)
			_ = WithBoard(b)(c)
		}

		if err := c.LoadPosition(FEN); err != nil {
			return fmt.Errorf("failed to load position: %w", err)
		}

		return nil
	}
}

// defaultOptions check if the setted options are valid and if not, set the default options.
func defaultOptions(chess *Chess) {
	if chess.board == nil {
		b, _ := gochess.NewBoard(8)
		_ = WithBoard(b)(chess)
	}

	if chess.FEN() == "" {
		_ = WithFEN(defaultFEN)(chess)
	}
}
