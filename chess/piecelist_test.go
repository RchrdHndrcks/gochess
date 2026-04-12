package chess

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/v2"
)

// scanBoardPieces returns a [2][7][]gochess.Coordinate of all piece squares
// found on the board by direct scan, suitable for comparison against the
// incrementally-maintained piece lists.
func scanBoardPieces(c *Chess) [2][7]map[gochess.Coordinate]struct{} {
	var out [2][7]map[gochess.Coordinate]struct{}
	for i := range out {
		for j := range out[i] {
			out[i][j] = make(map[gochess.Coordinate]struct{})
		}
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			sq := gochess.Coor(x, y)
			p, _ := c.board.Square(sq)
			if p == gochess.Empty {
				continue
			}
			out[colorIndex(gochess.PieceColor(p))][gochess.PieceType(p)][sq] = struct{}{}
		}
	}
	return out
}

func pieceListsMatchBoard(t *testing.T, c *Chess) {
	t.Helper()
	want := scanBoardPieces(c)
	for ci := 0; ci < 2; ci++ {
		for pt := 1; pt <= 6; pt++ {
			got := c.pieceLists[ci][pt].slice()
			if len(got) != len(want[ci][pt]) {
				t.Errorf("color=%d type=%d count mismatch: got %d, want %d", ci, pt, len(got), len(want[ci][pt]))
				continue
			}
			for _, sq := range got {
				if _, ok := want[ci][pt][sq]; !ok {
					t.Errorf("color=%d type=%d unexpected square %v", ci, pt, sq)
				}
			}
		}
	}
}

func TestPieceLists_StartPositionCounts(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	cases := []struct {
		color gochess.Piece
		ptype gochess.Piece
		want  int
	}{
		{gochess.White, gochess.Pawn, 8},
		{gochess.White, gochess.Knight, 2},
		{gochess.White, gochess.Bishop, 2},
		{gochess.White, gochess.Rook, 2},
		{gochess.White, gochess.Queen, 1},
		{gochess.White, gochess.King, 1},
		{gochess.Black, gochess.Pawn, 8},
		{gochess.Black, gochess.Knight, 2},
		{gochess.Black, gochess.Bishop, 2},
		{gochess.Black, gochess.Rook, 2},
		{gochess.Black, gochess.Queen, 1},
		{gochess.Black, gochess.King, 1},
	}
	for _, tc := range cases {
		got := c.PieceSquares(tc.color, tc.ptype)
		if len(got) != tc.want {
			t.Errorf("PieceSquares(%v,%v) = %d, want %d", tc.color, tc.ptype, len(got), tc.want)
		}
	}
}

func TestPieceLists_MatchBoardForVariousPositions(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
	}
	for _, fen := range fens {
		c, err := New()
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		if err := c.LoadPosition(fen); err != nil {
			t.Fatalf("LoadPosition(%q): %v", fen, err)
		}
		pieceListsMatchBoard(t, c)
	}
}

func TestPieceLists_MakeUnmakeStability_Kiwipete(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}

	moves := c.AvailableMoves()
	for _, mv := range moves {
		before := scanBoardPieces(c)
		if err := c.MakeMove(mv); err != nil {
			t.Fatalf("MakeMove(%s): %v", mv, err)
		}
		// After make, piece lists must still match the board.
		pieceListsMatchBoard(t, c)
		c.UnmakeMove()
		// After unmake, piece lists must match the original board state.
		after := scanBoardPieces(c)
		for ci := 0; ci < 2; ci++ {
			for pt := 1; pt <= 6; pt++ {
				if len(before[ci][pt]) != len(after[ci][pt]) {
					t.Fatalf("move=%s color=%d type=%d: count changed across make/unmake", mv, ci, pt)
				}
			}
		}
		pieceListsMatchBoard(t, c)
	}
}

func TestPieceLists_PieceSquaresIsIndependent(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	got := c.PieceSquares(gochess.White, gochess.Pawn)
	if len(got) != 8 {
		t.Fatalf("expected 8 white pawns, got %d", len(got))
	}
	// Mutate the returned slice; internal state must not change.
	got[0] = gochess.Coor(0, 0)
	again := c.PieceSquares(gochess.White, gochess.Pawn)
	if again[0] == got[0] && got[0] == (gochess.Coordinate{X: 0, Y: 0}) {
		// only fails if internal state was actually shared
		t.Fatalf("PieceSquares returned slice that shares internal state")
	}
}
