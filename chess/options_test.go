package chess

import (
	"testing"

	"github.com/RchrdHndrcks/gochess"
)

func TestWithParallelism(t *testing.T) {
	c, err := New(WithParallelism(2))
	if err != nil {
		t.Fatal(err)
	}

	if c.config.Parallelism != 2 {
		t.Errorf("expected parallelism to be 2, got %d", c.config.Parallelism)
	}
}

func TestWithBoard(t *testing.T) {
	board, err := gochess.NewBoard(8)
	if err != nil {
		t.Fatal(err)
	}

	c, err := New(WithBoard(board))
	if err != nil {
		t.Fatal(err)
	}

	if c.board.Width() != 8 {
		t.Errorf("expected board width to be 8, got %d", c.board.Width())
	}
}
