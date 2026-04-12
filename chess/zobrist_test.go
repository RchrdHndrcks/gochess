package chess

import "testing"

func TestZobristDeterminism(t *testing.T) {
	c1, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c2, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c1.Hash() != c2.Hash() {
		t.Fatalf("starting Hash mismatch: %d vs %d", c1.Hash(), c2.Hash())
	}
	if c1.PawnHash() != c2.PawnHash() {
		t.Fatalf("starting PawnHash mismatch: %d vs %d", c1.PawnHash(), c2.PawnHash())
	}
	if c1.Hash() == 0 {
		t.Fatalf("Hash should be non-zero for the starting position")
	}
}

func TestZobristChangesOnMove(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	before := c.Hash()
	if err := c.MakeMove("e2e4"); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
	if c.Hash() == before {
		t.Fatalf("Hash did not change after a move")
	}
}

func TestZobristRestoredAfterUnmake(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	beforeHash := c.Hash()
	beforePawn := c.PawnHash()
	if err := c.MakeMove("e2e4"); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
	c.UnmakeMove()
	if c.Hash() != beforeHash {
		t.Fatalf("Hash not restored: got %d want %d", c.Hash(), beforeHash)
	}
	if c.PawnHash() != beforePawn {
		t.Fatalf("PawnHash not restored: got %d want %d", c.PawnHash(), beforePawn)
	}
}

// TestZobristIncrementalMatchesScratch is the critical regression test:
// after every move in a multi-move game, the incrementally maintained
// Hash() must equal a fresh full recomputation. If this ever fails, the
// XOR delta sequence in applyMove is wrong.
func TestZobristIncrementalMatchesScratch(t *testing.T) {
	// Sequence covers a quiet pawn double-push, knight/bishop development,
	// castling, an exchange capture (Bxc6 / bxc6), and finishing development.
	// The capture is essential: an earlier version of the test never
	// exercised the "remove captured piece" XOR delta in applyMove.
	moves := []string{
		"e2e4", "e7e5",
		"g1f3", "b8c6",
		"f1b5", "a7a6",
		"b5c6", "d7c6", // Bxc6 then bxc6 — captures by both sides
		"e1g1", "f8e7", // castling
	}
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	startHash := c.Hash()
	startPawn := c.PawnHash()
	hashes := make([]uint64, 0, len(moves))
	pawnHashes := make([]uint64, 0, len(moves))
	hashes = append(hashes, startHash)
	pawnHashes = append(pawnHashes, startPawn)

	for i, m := range moves {
		if err := c.MakeMove(m); err != nil {
			t.Fatalf("move %d (%s): %v", i, m, err)
		}
		want := computeHashFromScratch(c)
		if c.Hash() != want {
			t.Fatalf("move %d (%s): incremental Hash %d != scratch %d", i, m, c.Hash(), want)
		}
		wantPawn := computePawnHashFromScratch(c)
		if c.PawnHash() != wantPawn {
			t.Fatalf("move %d (%s): incremental PawnHash %d != scratch %d", i, m, c.PawnHash(), wantPawn)
		}
		hashes = append(hashes, c.Hash())
		pawnHashes = append(pawnHashes, c.PawnHash())
	}

	// Unmake all moves and assert hashes match the snapshot from before
	// each move. This catches asymmetric XOR deltas in unmakeMove that the
	// forward-only assertions above would miss.
	for i := len(moves) - 1; i >= 0; i-- {
		c.UnmakeMove()
		if c.Hash() != hashes[i] {
			t.Fatalf("after unmake of move %d (%s): Hash %d != %d", i, moves[i], c.Hash(), hashes[i])
		}
		if c.PawnHash() != pawnHashes[i] {
			t.Fatalf("after unmake of move %d (%s): PawnHash %d != %d", i, moves[i], c.PawnHash(), pawnHashes[i])
		}
	}
	if c.Hash() != startHash {
		t.Fatalf("after full unmake: Hash %d != start %d", c.Hash(), startHash)
	}
}

// TestZobristIncrementalMatchesScratchPromotionWithCapture exercises the
// combined capture+promotion XOR delta. Regression for the case where
// applyMove forgot to subtract either the captured piece's contribution
// or the moving pawn's contribution before adding the promoted piece.
func TestZobristIncrementalMatchesScratchPromotionWithCapture(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Black knight on b8, white pawn on a7 ready to capture+promote on b8.
	if err := c.LoadPosition("1n6/P7/8/8/8/8/8/k6K w - - 0 1"); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}
	startHash := c.Hash()
	startPawn := c.PawnHash()
	if err := c.MakeMove("a7b8q"); err != nil {
		t.Fatalf("promotion-with-capture: %v", err)
	}
	if c.Hash() != computeHashFromScratch(c) {
		t.Fatalf("after a7xb8=Q: incremental Hash != scratch")
	}
	if c.PawnHash() != computePawnHashFromScratch(c) {
		t.Fatalf("after a7xb8=Q: incremental PawnHash != scratch")
	}
	c.UnmakeMove()
	if c.Hash() != startHash {
		t.Fatalf("after unmake of a7xb8=Q: Hash %d != start %d", c.Hash(), startHash)
	}
	if c.PawnHash() != startPawn {
		t.Fatalf("after unmake of a7xb8=Q: PawnHash %d != start %d", c.PawnHash(), startPawn)
	}
}

func TestZobristIncrementalMatchesScratchEnPassant(t *testing.T) {
	// Force an en passant capture and verify hash is correct after.
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, m := range []string{"e2e4", "a7a6", "e4e5", "d7d5", "e5d6"} {
		if err := c.MakeMove(m); err != nil {
			t.Fatalf("move %s: %v", m, err)
		}
		if c.Hash() != computeHashFromScratch(c) {
			t.Fatalf("after %s: hash mismatch", m)
		}
		if c.PawnHash() != computePawnHashFromScratch(c) {
			t.Fatalf("after %s: pawn hash mismatch", m)
		}
	}
}

func TestZobristIncrementalMatchesScratchPromotion(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.LoadPosition("8/P7/8/8/8/8/8/k6K w - - 0 1"); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}
	if err := c.MakeMove("a7a8q"); err != nil {
		t.Fatalf("promotion: %v", err)
	}
	if c.Hash() != computeHashFromScratch(c) {
		t.Fatalf("promotion: hash mismatch")
	}
	if c.PawnHash() != computePawnHashFromScratch(c) {
		t.Fatalf("promotion: pawn hash mismatch")
	}
}

func TestPawnHashUnchangedOnKnightMove(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	before := c.PawnHash()
	if err := c.MakeMove("g1f3"); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
	if c.PawnHash() != before {
		t.Fatalf("PawnHash changed on knight move: %d -> %d", before, c.PawnHash())
	}
}

func TestPawnHashChangesOnPawnMove(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	before := c.PawnHash()
	if err := c.MakeMove("e2e4"); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
	if c.PawnHash() == before {
		t.Fatalf("PawnHash did not change on pawn move")
	}
}

func TestZobristLoadPositionMatchesScratch(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1",
		"8/8/8/3pP3/8/8/8/k6K w - d6 0 1",
	}
	for _, f := range fens {
		if err := c.LoadPosition(f); err != nil {
			t.Fatalf("LoadPosition(%q): %v", f, err)
		}
		if c.Hash() != computeHashFromScratch(c) {
			t.Fatalf("LoadPosition(%q): incremental Hash != scratch", f)
		}
	}
}
