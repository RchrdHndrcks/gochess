package chess

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/v2"
)

func mustNew(t *testing.T, fen string) *Chess {
	t.Helper()
	var c *Chess
	var err error
	if fen == "" {
		c, err = New()
	} else {
		c, err = New(WithFEN(fen))
	}
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func sq(t *testing.T, alg string) gochess.Coordinate {
	t.Helper()
	co, err := AlgebraicToCoordinate(alg)
	if err != nil {
		t.Fatalf("AlgebraicToCoordinate(%q): %v", alg, err)
	}
	return co
}

func TestIsAttacked(t *testing.T) {
	t.Run("start position pawn attacks", func(t *testing.T) {
		c := mustNew(t, "")
		// White pawns on rank 2 attack rank 3.
		for _, s := range []string{"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3"} {
			if !c.IsAttacked(sq(t, s), gochess.White) {
				t.Errorf("expected %s attacked by white", s)
			}
		}
		// Black pawns on rank 7 attack rank 6.
		for _, s := range []string{"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6"} {
			if !c.IsAttacked(sq(t, s), gochess.Black) {
				t.Errorf("expected %s attacked by black", s)
			}
		}
		// e4 not attacked by anyone in start position.
		if c.IsAttacked(sq(t, "e4"), gochess.White) {
			t.Errorf("e4 should not be attacked by white in start position")
		}
		if c.IsAttacked(sq(t, "e4"), gochess.Black) {
			t.Errorf("e4 should not be attacked by black in start position")
		}
	})

	t.Run("knight attacks", func(t *testing.T) {
		c := mustNew(t, "8/8/8/3N4/8/8/8/k6K w - - 0 1")
		for _, s := range []string{"b6", "c7", "e7", "f6", "f4", "e3", "c3", "b4"} {
			if !c.IsAttacked(sq(t, s), gochess.White) {
				t.Errorf("expected %s attacked by white knight", s)
			}
		}
		if c.IsAttacked(sq(t, "d6"), gochess.White) {
			t.Errorf("d6 should not be attacked by knight on d5")
		}
	})

	t.Run("rook attacks", func(t *testing.T) {
		c := mustNew(t, "8/8/8/3R4/8/8/8/k6K w - - 0 1")
		if !c.IsAttacked(sq(t, "d1"), gochess.White) {
			t.Error("d1 should be attacked by rook on d5")
		}
		if !c.IsAttacked(sq(t, "a5"), gochess.White) {
			t.Error("a5 should be attacked by rook on d5")
		}
		if c.IsAttacked(sq(t, "e6"), gochess.White) {
			t.Error("e6 should not be attacked by rook on d5")
		}
	})

	t.Run("bishop attacks", func(t *testing.T) {
		c := mustNew(t, "8/8/8/3B4/8/8/8/k6K w - - 0 1")
		if !c.IsAttacked(sq(t, "a8"), gochess.White) {
			t.Error("a8 should be attacked by bishop on d5")
		}
		if !c.IsAttacked(sq(t, "h1"), gochess.White) {
			t.Error("h1 should be attacked by bishop on d5")
		}
		if c.IsAttacked(sq(t, "d4"), gochess.White) {
			t.Error("d4 should not be attacked by bishop on d5")
		}
	})

	t.Run("queen attacks file and diagonal", func(t *testing.T) {
		c := mustNew(t, "8/8/8/3Q4/8/8/8/k6K w - - 0 1")
		if !c.IsAttacked(sq(t, "d1"), gochess.White) {
			t.Error("d1 should be attacked by queen on d5 (file)")
		}
		if !c.IsAttacked(sq(t, "h1"), gochess.White) {
			t.Error("h1 should be attacked by queen on d5 (diagonal)")
		}
	})

	t.Run("blocked ray", func(t *testing.T) {
		c := mustNew(t, "8/8/8/3R4/3P4/8/8/k6K w - - 0 1")
		if c.IsAttacked(sq(t, "d1"), gochess.White) {
			t.Error("d1 should not be attacked: pawn on d4 blocks rook")
		}
	})

	t.Run("replaces check detection", func(t *testing.T) {
		// Black king on e8, white rook on e1. King is in check.
		c := mustNew(t, "4k3/8/8/8/8/8/8/4R1K1 b - - 0 1")
		if !c.IsAttacked(sq(t, "e8"), gochess.White) {
			t.Error("e8 should be attacked by white rook on e1")
		}
		if !c.IsCheck() {
			t.Error("position should report check")
		}
	})
}

func TestAttackedBy(t *testing.T) {
	t.Run("pawn and knight together", func(t *testing.T) {
		// Square e4: attacked by white pawn d3 and white knight f2? Place pieces.
		c := mustNew(t, "k6K/8/8/8/8/3P4/5N2/8 w - - 0 1")
		ps := c.AttackedBy(sq(t, "e4"), gochess.White)
		if !ps.Has(PawnSet) {
			t.Errorf("expected PawnSet, got %s", ps)
		}
		if !ps.Has(KnightSet) {
			t.Errorf("expected KnightSet, got %s", ps)
		}
		if ps.Has(BishopSet) || ps.Has(RookSet) || ps.Has(QueenSet) || ps.Has(KingSet) {
			t.Errorf("unexpected extra bits: %s", ps)
		}
	})

	t.Run("no attackers", func(t *testing.T) {
		c := mustNew(t, "")
		if got := c.AttackedBy(sq(t, "e4"), gochess.White); got != NoPieces {
			t.Errorf("expected NoPieces, got %s", got)
		}
	})

	t.Run("rook and queen on same file", func(t *testing.T) {
		// White queen d8, white rook d4, target d1 — only rook reaches first.
		// AttackedBy stops at first same-side piece on each ray.
		c := mustNew(t, "3Q4/8/8/8/3R4/8/8/k2K4 w - - 0 1")
		ps := c.AttackedBy(sq(t, "d1"), gochess.White)
		if !ps.Has(RookSet) {
			t.Errorf("expected RookSet, got %s", ps)
		}
		if ps.Has(QueenSet) {
			t.Errorf("queen should be x-ray-blocked for AttackedBy: %s", ps)
		}
	})
}

func TestPawnAttackMap(t *testing.T) {
	c := mustNew(t, "")
	white := c.PawnAttackMap(gochess.White)
	black := c.PawnAttackMap(gochess.Black)

	// White pawns on rank 2 (y=6) attack rank 3 (y=5): all 8 files.
	var expectedWhite uint64
	for x := 0; x < 8; x++ {
		expectedWhite |= uint64(1) << uint(5*8+x)
	}
	if white != expectedWhite {
		t.Errorf("white pawn map mismatch: got %064b want %064b", white, expectedWhite)
	}

	var expectedBlack uint64
	for x := 0; x < 8; x++ {
		expectedBlack |= uint64(1) << uint(2*8+x)
	}
	if black != expectedBlack {
		t.Errorf("black pawn map mismatch: got %064b want %064b", black, expectedBlack)
	}
}

func TestAttackers(t *testing.T) {
	t.Run("two rooks on same file (x-ray)", func(t *testing.T) {
		// White rooks on d4 and d8 attacking d1.
		c := mustNew(t, "3R4/8/8/8/3R4/8/8/k2K4 w - - 0 1")
		got := c.Attackers(sq(t, "d1"), gochess.White, nil)
		if len(got) != 2 {
			t.Fatalf("expected 2 attackers, got %d: %v", len(got), got)
		}
		// Both should be rooks (cheapest tier here).
		want1, want2 := sq(t, "d4"), sq(t, "d8")
		ok := (got[0] == want1 && got[1] == want2) || (got[0] == want2 && got[1] == want1)
		if !ok {
			t.Errorf("unexpected attackers: %v", got)
		}
	})

	t.Run("queen behind rook on file (x-ray cross-type)", func(t *testing.T) {
		// White rook d4, white queen d8. Both should be returned for d1.
		c := mustNew(t, "3Q4/8/8/8/3R4/8/8/k2K4 w - - 0 1")
		got := c.Attackers(sq(t, "d1"), gochess.White, nil)
		if len(got) != 2 {
			t.Fatalf("expected 2 attackers, got %d: %v", len(got), got)
		}
		// Rook (cheaper) must come before queen.
		if got[0] != sq(t, "d4") {
			t.Errorf("rook should be first, got %v", got)
		}
		if got[1] != sq(t, "d8") {
			t.Errorf("queen should be second, got %v", got)
		}
	})

	t.Run("bishop behind queen on diagonal", func(t *testing.T) {
		// White queen on b2, white bishop on a1. Target h8.
		// Queen reaches h8 diagonally; bishop on a1 is x-ray behind queen.
		c := mustNew(t, "8/8/8/8/8/8/1Q6/B3k2K w - - 0 1")
		got := c.Attackers(sq(t, "h8"), gochess.White, nil)
		if len(got) != 2 {
			t.Fatalf("expected 2 attackers, got %d: %v", len(got), got)
		}
		// Bishop is cheaper than queen, so it must come first.
		if got[0] != sq(t, "a1") {
			t.Errorf("bishop should be first (cheapest), got %v", got)
		}
		if got[1] != sq(t, "b2") {
			t.Errorf("queen should be second, got %v", got)
		}
	})

	t.Run("mixed cheapest-first ordering", func(t *testing.T) {
		// White pawn c2 (no, c2 attacks b3/d3), let's place:
		// Target e4: attacker pawn d3, knight f2, bishop b1, queen e1 (file).
		c := mustNew(t, "7k/8/7K/8/8/3P4/5N2/4Q2B w - - 0 1")
		got := c.Attackers(sq(t, "e4"), gochess.White, nil)
		// Expect pawn first, then knight, then bishop, then queen.
		if len(got) != 4 {
			t.Fatalf("expected 4 attackers, got %d: %v", len(got), got)
		}
		if got[0] != sq(t, "d3") {
			t.Errorf("pawn first, got %v", got)
		}
		if got[1] != sq(t, "f2") {
			t.Errorf("knight second, got %v", got)
		}
		if got[2] != sq(t, "h1") {
			t.Errorf("bishop third, got %v", got)
		}
		if got[3] != sq(t, "e1") {
			t.Errorf("queen fourth, got %v", got)
		}
	})

	t.Run("appends to dst", func(t *testing.T) {
		c := mustNew(t, "")
		dst := make([]gochess.Coordinate, 0, 8)
		dst = append(dst, sq(t, "a1"))
		got := c.Attackers(sq(t, "e3"), gochess.White, dst)
		if len(got) < 1 || got[0] != sq(t, "a1") {
			t.Errorf("dst prefix not preserved: %v", got)
		}
	})
}

func TestIsAttackedQueen(t *testing.T) {
	t.Run("queen on file", func(t *testing.T) {
		c := mustNew(t, "4k3/8/8/8/3q4/8/7K/8 b - - 0 1")
		if !c.IsAttacked(sq(t, "d1"), gochess.Black) {
			t.Error("d1 should be attacked by black queen on d4 (file)")
		}
	})
	t.Run("queen on diagonal", func(t *testing.T) {
		c := mustNew(t, "4k3/8/8/8/3q4/8/7K/8 b - - 0 1")
		if !c.IsAttacked(sq(t, "g1"), gochess.Black) {
			t.Error("g1 should be attacked by black queen on d4 (diagonal)")
		}
	})
}
