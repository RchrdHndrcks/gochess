package chess

import (
	"testing"

	gochess "github.com/RchrdHndrcks/gochess/v2"
)

// TestIntegration_RuyLopez plays a 16-ply Ruy Lopez and after each move
// verifies four invariants:
//  1. Hash() equals computeHashFromScratch (incremental hash correctness).
//  2. PieceSquares matches a fresh board scan (piece-list correctness).
//  3. Captures() + QuietMoves() == Moves() (staged-gen consistency).
//  4. After unmaking every move, FEN and hash return to start position.
func TestIntegration_RuyLopez(t *testing.T) {
	moves := []string{
		"e2e4", "e7e5",
		"g1f3", "b8c6",
		"f1b5", "a7a6",
		"b5a4", "g8f6",
		"e1g1", "f8e7",
		"f1e1", "b7b5",
		"a4b3", "d7d6",
		"c2c3", "e8g8",
	}

	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	startFEN := c.FEN()
	startHash := c.Hash()
	fenHistory := make([]string, 0, len(moves))

	for i, uci := range moves {
		fenHistory = append(fenHistory, c.FEN())
		if err := c.MakeMove(uci); err != nil {
			t.Fatalf("ply %d (%s): MakeMove: %v", i, uci, err)
		}

		// 1. Incremental hash must match recomputed hash.
		if got, want := c.Hash(), computeHashFromScratch(c); got != want {
			t.Fatalf("ply %d (%s): Hash %d != scratch %d", i, uci, got, want)
		}

		// 2. PieceSquares must match a board scan for every (color,type).
		verifyPieceLists(t, c, i, uci)

		// 3. Captures + QuietMoves must equal Moves (set equality on UCI).
		verifyStagedGen(t, c, i, uci)
	}

	// 4. Unmake all moves and verify start FEN / hash are restored.
	for i := len(moves) - 1; i >= 0; i-- {
		c.UnmakeMoveCompact()
		if got := c.FEN(); got != fenHistory[i] {
			t.Fatalf("after unmake ply %d: FEN %q != %q", i, got, fenHistory[i])
		}
	}
	if c.FEN() != startFEN {
		t.Fatalf("after full unmake: FEN %q != start %q", c.FEN(), startFEN)
	}
	if c.Hash() != startHash {
		t.Fatalf("after full unmake: Hash %d != start %d", c.Hash(), startHash)
	}
	if c.Hash() != computeHashFromScratch(c) {
		t.Fatalf("after full unmake: incremental hash diverged from scratch")
	}
}

func verifyPieceLists(t *testing.T, c *Chess, ply int, uci string) {
	t.Helper()
	colors := []gochess.Piece{gochess.White, gochess.Black}
	types := []gochess.Piece{gochess.Pawn, gochess.Knight, gochess.Bishop, gochess.Rook, gochess.Queen, gochess.King}

	// Build map of expected squares per (color,type) by scanning the board.
	expected := map[gochess.Piece]map[gochess.Coordinate]struct{}{}
	for _, col := range colors {
		for _, pt := range types {
			expected[col|pt] = map[gochess.Coordinate]struct{}{}
		}
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			sq := gochess.Coor(x, y)
			p, _ := c.board.Square(sq)
			if p == gochess.Empty {
				continue
			}
			expected[p][sq] = struct{}{}
		}
	}

	for _, col := range colors {
		for _, pt := range types {
			got := c.PieceSquares(col, pt)
			want := expected[col|pt]
			// Build a set from got so a duplicate entry with a missing
			// entry cannot pass the length check silently.
			gotSet := map[gochess.Coordinate]struct{}{}
			for _, sq := range got {
				if _, dup := gotSet[sq]; dup {
					t.Fatalf("ply %d (%s): PieceSquares(%v,%v) returned duplicate sq %v", ply, uci, col, pt, sq)
				}
				gotSet[sq] = struct{}{}
			}
			if len(gotSet) != len(want) {
				t.Fatalf("ply %d (%s): PieceSquares(%v,%v) len=%d, want %d", ply, uci, col, pt, len(gotSet), len(want))
			}
			for sq := range want {
				if _, ok := gotSet[sq]; !ok {
					t.Fatalf("ply %d (%s): PieceSquares(%v,%v) missing expected sq %v", ply, uci, col, pt, sq)
				}
			}
			for sq := range gotSet {
				if _, ok := want[sq]; !ok {
					t.Fatalf("ply %d (%s): PieceSquares(%v,%v) had unexpected sq %v", ply, uci, col, pt, sq)
				}
			}
		}
	}
}

func verifyStagedGen(t *testing.T, c *Chess, ply int, uci string) {
	t.Helper()
	all := c.Moves()
	caps := c.Captures()
	quiets := c.QuietMoves()

	if len(caps)+len(quiets) != len(all) {
		t.Fatalf("ply %d (%s): Captures(%d)+QuietMoves(%d) != Moves(%d)",
			ply, uci, len(caps), len(quiets), len(all))
	}

	seen := map[string]int{}
	for _, m := range all {
		seen[m.UCI()]++
	}
	for _, m := range caps {
		seen[m.UCI()]--
	}
	for _, m := range quiets {
		seen[m.UCI()]--
	}
	for u, n := range seen {
		if n != 0 {
			t.Fatalf("ply %d (%s): staged gen mismatch for %s (delta %d)", ply, uci, u, n)
		}
	}
}
