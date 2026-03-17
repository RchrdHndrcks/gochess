package chess

import (
	"testing"
)

func TestToSAN_PawnMoves(t *testing.T) {
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	tests := []struct {
		uci string
		san string
	}{
		{"e2e4", "e4"},
		{"d2d4", "d4"},
		{"a2a3", "a3"},
		{"h2h4", "h4"},
	}

	for _, tt := range tests {
		san, err := ToSAN(c, tt.uci)
		if err != nil {
			t.Errorf("ToSAN(%s) error: %v", tt.uci, err)
			continue
		}
		if san != tt.san {
			t.Errorf("ToSAN(%s) = %s, want %s", tt.uci, san, tt.san)
		}
	}
}

func TestToSAN_PieceMoves(t *testing.T) {
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	// Nf3
	san, err := ToSAN(c, "g1f3")
	if err != nil {
		t.Fatalf("ToSAN(g1f3) error: %v", err)
	}
	if san != "Nf3" {
		t.Errorf("ToSAN(g1f3) = %s, want Nf3", san)
	}

	// Nc3
	san, err = ToSAN(c, "b1c3")
	if err != nil {
		t.Fatalf("ToSAN(b1c3) error: %v", err)
	}
	if san != "Nc3" {
		t.Errorf("ToSAN(b1c3) = %s, want Nc3", san)
	}
}

func TestToSAN_BishopMove(t *testing.T) {
	// Play e4, e5 to open diagonals, then Bc4.
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	moves := []string{"e2e4", "e7e5", "f1c4"}
	for i, m := range moves {
		if i < 2 {
			if err := c.MakeMove(m); err != nil {
				t.Fatalf("MakeMove(%s) error: %v", m, err)
			}
		} else {
			san, err := ToSAN(c, m)
			if err != nil {
				t.Fatalf("ToSAN(%s) error: %v", m, err)
			}
			if san != "Bc4" {
				t.Errorf("ToSAN(%s) = %s, want Bc4", m, san)
			}
		}
	}
}

func TestToSAN_Captures(t *testing.T) {
	// Set up a position where exd5 is possible.
	c, err := New(
		WithFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san, err := ToSAN(c, "e4d5")
	if err != nil {
		t.Fatalf("ToSAN(e4d5) error: %v", err)
	}
	if san != "exd5" {
		t.Errorf("ToSAN(e4d5) = %s, want exd5", san)
	}
}

func TestToSAN_KnightCapture(t *testing.T) {
	// Set up a position where Nxf3 can happen (after a knight is on f3 and can be captured).
	c, err := New(
		WithFEN("rnbqkb1r/pppppppp/5n2/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 1 2"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	// White plays e5 attacking the knight, then we'll set up a capture scenario.
	// Instead, let's use a position where a knight captures.
	c2, err := New(
		WithFEN("rnbqkb1r/pppp1ppp/4pn2/8/3PP3/8/PPP2PPP/RNBQKBNR w KQkq - 0 3"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	// Play e4-e5, then Nxe4 would be ...Nxe4 but let's just check the notation is correct.
	// Actually let's use a simpler approach: verify a knight capturing a piece.
	_ = c
	// White to move: Play Nc3 first, then setup a capture.
	// Let's use a direct FEN with a capture available.
	c3, err := New(
		WithFEN("r1bqkbnr/pppppppp/2n5/4P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}
	_ = c2

	san, err := ToSAN(c3, "c6e5")
	if err != nil {
		t.Fatalf("ToSAN(c6e5) error: %v", err)
	}
	if san != "Nxe5" {
		t.Errorf("ToSAN(c6e5) = %s, want Nxe5", san)
	}
}

func TestToSAN_Castling(t *testing.T) {
	// Kingside castling.
	c, err := New(
		WithFEN("rnbqk2r/ppppbppp/4pn2/8/4P3/5N2/PPPPBPPP/RNBQK2R w KQkq - 4 4"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san, err := ToSAN(c, "e1g1")
	if err != nil {
		t.Fatalf("ToSAN(e1g1) error: %v", err)
	}
	if san != "O-O" {
		t.Errorf("ToSAN(e1g1) = %s, want O-O", san)
	}

	// Queenside castling.
	c2, err := New(
		WithFEN("r3kbnr/pppqpppp/2n1b3/3p4/3P4/2N1B3/PPPQPPPP/R3KBNR w KQkq - 6 5"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san, err = ToSAN(c2, "e1c1")
	if err != nil {
		t.Fatalf("ToSAN(e1c1) error: %v", err)
	}
	if san != "O-O-O" {
		t.Errorf("ToSAN(e1c1) = %s, want O-O-O", san)
	}
}

func TestToSAN_Check(t *testing.T) {
	// Position where Qf7+ gives check.
	c, err := New(
		WithFEN("rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	// White queen can go to h5 to attack f7. Let's use Qh5.
	san, err := ToSAN(c, "d1h5")
	if err != nil {
		t.Fatalf("ToSAN(d1h5) error: %v", err)
	}
	if san != "Qh5" {
		t.Errorf("ToSAN(d1h5) = %s, want Qh5", san)
	}

	// Set up a position where Qf7+ is check.
	c2, err := New(
		WithFEN("rnbqkbnr/pppp1ppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQkq - 0 2"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san2, err := ToSAN(c2, "h5f7")
	if err != nil {
		t.Fatalf("ToSAN(h5f7) error: %v", err)
	}
	if san2 != "Qxf7+" {
		t.Errorf("ToSAN(h5f7) = %s, want Qxf7+", san2)
	}
}

func TestToSAN_ScholarsMate(t *testing.T) {
	// Scholar's mate: 1.e4 e5 2.Bc4 Nc6 3.Qh5 Nf6?? 4.Qxf7#
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	movesAndSAN := []struct {
		uci string
		san string
	}{
		{"e2e4", "e4"},
		{"e7e5", "e5"},
		{"f1c4", "Bc4"},
		{"b8c6", "Nc6"},
		{"d1h5", "Qh5"},
		{"g8f6", "Nf6"},
		{"h5f7", "Qxf7#"},
	}

	for _, ms := range movesAndSAN {
		san, err := ToSAN(c, ms.uci)
		if err != nil {
			t.Fatalf("ToSAN(%s) error: %v", ms.uci, err)
		}
		if san != ms.san {
			t.Errorf("ToSAN(%s) = %s, want %s", ms.uci, san, ms.san)
		}

		if err := c.MakeMove(ms.uci); err != nil {
			t.Fatalf("MakeMove(%s) error: %v", ms.uci, err)
		}
	}

	if !c.IsCheckmate() {
		t.Error("expected checkmate after Scholar's mate")
	}
}

func TestToSAN_Promotion(t *testing.T) {
	// Position where a pawn can promote. Kings far apart, pawn not adjacent to black king.
	c, err := New(
		WithFEN("8/4P1k1/8/8/8/8/8/K7 w - - 0 1"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san, err := ToSAN(c, "e7e8q")
	if err != nil {
		t.Fatalf("ToSAN(e7e8q) error: %v", err)
	}
	if san != "e8=Q" {
		t.Errorf("ToSAN(e7e8q) = %s, want e8=Q", san)
	}
}

func TestToSAN_Disambiguation(t *testing.T) {
	// Two rooks on the same rank, need file disambiguation.
	// Rooks on a1 and f1, king on h1, black king on h8.
	c, err := New(
		WithFEN("7k/8/8/8/8/8/8/R4RK1 w - - 0 1"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san, err := ToSAN(c, "a1d1")
	if err != nil {
		t.Fatalf("ToSAN(a1d1) error: %v", err)
	}
	if san != "Rad1" {
		t.Errorf("ToSAN(a1d1) = %s, want Rad1", san)
	}

	// Two rooks on the same file, need rank disambiguation.
	c2, err := New(
		WithFEN("7k/8/8/8/8/R7/8/R3K3 w - - 0 1"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san, err = ToSAN(c2, "a1a2")
	if err != nil {
		t.Fatalf("ToSAN(a1a2) error: %v", err)
	}
	if san != "R1a2" {
		t.Errorf("ToSAN(a1a2) = %s, want R1a2", san)
	}
}

func TestFromSAN_PawnMoves(t *testing.T) {
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	tests := []struct {
		san string
		uci string
	}{
		{"e4", "e2e4"},
		{"d4", "d2d4"},
		{"a3", "a2a3"},
	}

	for _, tt := range tests {
		uci, err := FromSAN(c, tt.san)
		if err != nil {
			t.Errorf("FromSAN(%s) error: %v", tt.san, err)
			continue
		}
		if uci != tt.uci {
			t.Errorf("FromSAN(%s) = %s, want %s", tt.san, uci, tt.uci)
		}
	}
}

func TestFromSAN_PieceMoves(t *testing.T) {
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	uci, err := FromSAN(c, "Nf3")
	if err != nil {
		t.Fatalf("FromSAN(Nf3) error: %v", err)
	}
	if uci != "g1f3" {
		t.Errorf("FromSAN(Nf3) = %s, want g1f3", uci)
	}
}

func TestFromSAN_Castling(t *testing.T) {
	c, err := New(
		WithFEN("rnbqk2r/ppppbppp/4pn2/8/4P3/5N2/PPPPBPPP/RNBQK2R w KQkq - 4 4"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	uci, err := FromSAN(c, "O-O")
	if err != nil {
		t.Fatalf("FromSAN(O-O) error: %v", err)
	}
	if uci != "e1g1" {
		t.Errorf("FromSAN(O-O) = %s, want e1g1", uci)
	}
}

func TestFromSAN_Captures(t *testing.T) {
	c, err := New(
		WithFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	uci, err := FromSAN(c, "exd5")
	if err != nil {
		t.Fatalf("FromSAN(exd5) error: %v", err)
	}
	if uci != "e4d5" {
		t.Errorf("FromSAN(exd5) = %s, want e4d5", uci)
	}
}

func TestFromSAN_Promotion(t *testing.T) {
	c, err := New(
		WithFEN("8/4P1k1/8/8/8/8/8/K7 w - - 0 1"),
		WithParallelism(1),
	)
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	uci, err := FromSAN(c, "e8=Q")
	if err != nil {
		t.Fatalf("FromSAN(e8=Q) error: %v", err)
	}
	if uci != "e7e8q" {
		t.Errorf("FromSAN(e8=Q) = %s, want e7e8q", uci)
	}
}

func TestFromSAN_Roundtrip(t *testing.T) {
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	// Test roundtrip for all available moves in the starting position.
	for _, uci := range c.AvailableMoves() {
		san, err := ToSAN(c, uci)
		if err != nil {
			t.Errorf("ToSAN(%s) error: %v", uci, err)
			continue
		}

		roundtrip, err := FromSAN(c, san)
		if err != nil {
			t.Errorf("FromSAN(%s) error: %v", san, err)
			continue
		}

		if roundtrip != uci {
			t.Errorf("roundtrip failed: %s -> %s -> %s", uci, san, roundtrip)
		}
	}
}

func TestMoveToSAN_Method(t *testing.T) {
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	san, err := c.MoveToSAN("e2e4")
	if err != nil {
		t.Fatalf("MoveToSAN(e2e4) error: %v", err)
	}
	if san != "e4" {
		t.Errorf("MoveToSAN(e2e4) = %s, want e4", san)
	}
}

func TestMoveFromSAN_Method(t *testing.T) {
	c, err := New(WithParallelism(1))
	if err != nil {
		t.Fatalf("failed to create chess: %v", err)
	}

	uci, err := c.MoveFromSAN("e4")
	if err != nil {
		t.Fatalf("MoveFromSAN(e4) error: %v", err)
	}
	if uci != "e2e4" {
		t.Errorf("MoveFromSAN(e4) = %s, want e2e4", uci)
	}
}
