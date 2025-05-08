package chess

import "github.com/RchrdHndrcks/gochess"

// boardAdapter is an adapter for gochess.Board that implements the Board interface.
//
// It is used to take advantage of the parallelism implementing the Cloner interface.
// The Cloner interface is implemented by the adapter to avoid dependencies to the
// gochess/chess package from the gochess package.
type boardAdapter struct {
	*gochess.Board
}

func newBoardAdapter(board *gochess.Board) *boardAdapter {
	return &boardAdapter{
		Board: board,
	}
}

// Clone implements the Cloner interface.
func (b *boardAdapter) Clone() Board {
	return &boardAdapter{
		Board: b.Board.Clone(),
	}
}
