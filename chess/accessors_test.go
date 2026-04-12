package chess

import "testing"

func TestLastMoveStartingPosition(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.LastMove() != NullMove {
		t.Fatalf("LastMove on fresh game = %v, want NullMove", c.LastMove())
	}
}

func TestLastMoveAfterMakeMove(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.MakeMove("e2e4"); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
	last := c.LastMove()
	if last == NullMove {
		t.Fatalf("LastMove = NullMove after a move was made")
	}
	if got := last.UCI(); got != "e2e4" {
		t.Fatalf("LastMove UCI = %q, want %q", got, "e2e4")
	}
}

func TestLastMoveAfterMakeMoveCompact(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	m, err := c.ParseUCIMove("e2e4")
	if err != nil {
		t.Fatalf("ParseUCIMove: %v", err)
	}
	if err := c.MakeMoveCompact(m); err != nil {
		t.Fatalf("MakeMoveCompact: %v", err)
	}
	if got := c.LastMove().UCI(); got != "e2e4" {
		t.Fatalf("LastMove UCI = %q, want %q", got, "e2e4")
	}
}

func TestHalfmoveClock(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.HalfmoveClock() != 0 {
		t.Fatalf("HalfmoveClock fresh game = %d, want 0", c.HalfmoveClock())
	}
	// Pawn move resets to zero (still zero); knight move increments.
	if err := c.MakeMove("g1f3"); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
	if c.HalfmoveClock() != 1 {
		t.Fatalf("HalfmoveClock after knight move = %d, want 1", c.HalfmoveClock())
	}
	if err := c.MakeMove("e7e5"); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
	if c.HalfmoveClock() != 0 {
		t.Fatalf("HalfmoveClock after pawn move = %d, want 0", c.HalfmoveClock())
	}
}

func TestPieceAtAccessor(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// e2 in algebraic = (4,6). squareFromCoordinate gives some sq idx; use it.
	cor, err := AlgebraicToCoordinate("e2")
	if err != nil {
		t.Fatalf("AlgebraicToCoordinate: %v", err)
	}
	sq := int(squareFromCoordinate(cor))
	pt, color, ok := c.PieceAt(sq)
	if !ok {
		t.Fatalf("PieceAt(e2) ok = false, want true")
	}
	if pt != 1 { // Pawn
		t.Fatalf("PieceAt(e2) type = %d, want Pawn(1)", pt)
	}
	if color != 0b01000 { // White
		t.Fatalf("PieceAt(e2) color = %d, want White", color)
	}
	// Empty square.
	cor, _ = AlgebraicToCoordinate("e4")
	if _, _, ok := c.PieceAt(int(squareFromCoordinate(cor))); ok {
		t.Fatalf("PieceAt(e4) ok = true, want false (empty)")
	}
	// Out-of-range.
	if _, _, ok := c.PieceAt(-1); ok {
		t.Fatalf("PieceAt(-1) ok = true, want false")
	}
	if _, _, ok := c.PieceAt(64); ok {
		t.Fatalf("PieceAt(64) ok = true, want false")
	}
}
