package chess

import (
	"testing"
)

const kiwipeteFEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

func TestMakeMoveCompact_KiwipeteFENMatchesMakeMove(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := a.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}

	b, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := b.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}

	for _, uci := range a.AvailableMoves() {
		uci := uci
		t.Run(uci, func(t *testing.T) {
			if err := a.MakeMove(uci); err != nil {
				t.Fatalf("MakeMove(%s): %v", uci, err)
			}
			m, err := b.ParseUCIMove(uci)
			if err != nil {
				t.Fatalf("ParseUCIMove(%s): %v", uci, err)
			}
			if err := b.MakeMoveCompact(m); err != nil {
				t.Fatalf("MakeMoveCompact(%s): %v", uci, err)
			}
			if a.FEN() != b.FEN() {
				t.Errorf("FEN mismatch after %s:\n  MakeMove:        %s\n  MakeMoveCompact: %s",
					uci, a.FEN(), b.FEN())
			}
			a.UnmakeMove()
			b.UnmakeMoveCompact()
			if a.FEN() != b.FEN() {
				t.Errorf("FEN mismatch after unmake %s:\n  UnmakeMove:        %s\n  UnmakeMoveCompact: %s",
					uci, a.FEN(), b.FEN())
			}
		})
	}
}

func TestParseUCIMove_RoundTrip(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}

	for _, uci := range c.AvailableMoves() {
		uci := uci
		t.Run(uci, func(t *testing.T) {
			m, err := c.ParseUCIMove(uci)
			if err != nil {
				t.Fatalf("ParseUCIMove(%s): %v", uci, err)
			}
			if got := m.UCI(); got != uci {
				t.Errorf("Move.UCI() = %q, want %q", got, uci)
			}
		})
	}
}

func TestParseUCIMove_Invalid(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	cases := []string{
		"",
		"e2",
		"e2e4q5",
		"z9z9",
		"e7e8q", // not legal from start position
	}
	for _, uci := range cases {
		uci := uci
		t.Run(uci, func(t *testing.T) {
			if _, err := c.ParseUCIMove(uci); err == nil {
				t.Errorf("ParseUCIMove(%q) want error, got nil", uci)
			}
		})
	}
}

func TestUnmakeMoveCompact_PromotionWithoutCapture(t *testing.T) {
	// Position with white pawn on e7, ready to push-promote without capturing.
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	const fen = "k7/4P3/8/8/8/8/8/4K3 w - - 0 1"
	if err := c.LoadPosition(fen); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}

	beforeFEN := c.FEN()
	m, err := c.ParseUCIMove("e7e8q")
	if err != nil {
		t.Fatalf("ParseUCIMove: %v", err)
	}
	if err := c.MakeMoveCompact(m); err != nil {
		t.Fatalf("MakeMoveCompact: %v", err)
	}
	c.UnmakeMoveCompact()
	if c.FEN() != beforeFEN {
		t.Errorf("FEN after unmake = %q, want %q", c.FEN(), beforeFEN)
	}
	// e8 must be empty.
	if name, _ := c.Square("e8"); name != "" {
		t.Errorf("e8 = %q, want empty", name)
	}
}

func TestMoves_GivesCheckBit(t *testing.T) {
	// For each tested FEN, generate Moves() and assert that the GivesCheck
	// bit on every returned move agrees with ground truth (play the move,
	// query the opponent's king square, and call IsAttacked equivalent via
	// isCheck on the side now to move, which is the opponent we just moved
	// against).
	fens := []string{
		// Starting position (no checks possible).
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		// White queen on h5 and black king on e8 — Qh5xf7 etc. gives check.
		"4k3/5p2/8/7Q/8/8/8/4K3 w - - 0 1",
		// Kiwipete-like middlegame position with multiple checking moves.
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
	}
	for _, fen := range fens {
		fen := fen
		t.Run(fen, func(t *testing.T) {
			c, err := New()
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			if err := c.LoadPosition(fen); err != nil {
				t.Fatalf("LoadPosition: %v", err)
			}
			moves := c.Moves()
			anyChecked := false
			for _, m := range moves {
				// Ground truth: play the move via MakeMoveCompact, then ask
				// isCheck (which checks whether the side to move — the
				// opponent of the player who just moved — is in check).
				clone, err := New()
				if err != nil {
					t.Fatalf("New: %v", err)
				}
				if err := clone.LoadPosition(fen); err != nil {
					t.Fatalf("LoadPosition: %v", err)
				}
				if err := clone.MakeMoveCompact(m.WithGivesCheck(false)); err != nil {
					t.Fatalf("MakeMoveCompact(%s): %v", m.UCI(), err)
				}
				wantCheck := clone.isCheck()
				if got := m.GivesCheck(); got != wantCheck {
					t.Errorf("move %s GivesCheck=%v, want %v", m.UCI(), got, wantCheck)
				}
				if wantCheck {
					anyChecked = true
				}
			}
			_ = anyChecked
		})
	}
}
