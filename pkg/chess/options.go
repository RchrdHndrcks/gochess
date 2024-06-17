package chess

import (
	"fmt"
	"reflect"

	"github.com/RchrdHndrcks/gochess/pkg"
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
		if b == nil || reflect.ValueOf(b).IsNil() {
			return fmt.Errorf("board cannot be nil")
		}

		c.board = b
		return nil
	}
}

// WithFEN sets the FEN of the chess.
// If the FEN is invalid, it returns an error.
// If you try to set the FEN before the board, it will set the default board.
func WithFEN(FEN string) Option {
	return func(c *Chess) error {
		if c.board == nil || reflect.ValueOf(c.board).IsNil() {
			_ = WithBoard(pkg.NewBoard())(c) // nolint:errcheck
		}

		if err := c.board.LoadPosition(FEN); err != nil {
			return fmt.Errorf("failed to load FEN: %w", err)
		}

		if err := c.setProperties(FEN); err != nil {
			return fmt.Errorf("failed to set properties: %w", err)
		}

		return nil
	}
}

// defaultOptions check if the setted options are valid and if not, set the default options.
func defaultOptions(chess *Chess) {
	if chess.board == nil || reflect.ValueOf(chess.board).IsNil() {
		_ = WithBoard(pkg.NewBoard())(chess) // nolint:errcheck
	}
	if chess.FEN() == "8/8/8/8/8/8/8/8 w - - 0 0" {
		_ = WithFEN(defaultFEN)(chess) // nolint:errcheck
	}
}
