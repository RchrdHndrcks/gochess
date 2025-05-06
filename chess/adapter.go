package chess

import "github.com/RchrdHndrcks/gochess"

type boardAdapter struct {
	*gochess.Board
}

func newBoardAdapter(board *gochess.Board) *boardAdapter {
	return &boardAdapter{
		Board: board,
	}
}

func (b *boardAdapter) Clone() Board {
	return &boardAdapter{
		Board: b.Board.Clone(),
	}
}
