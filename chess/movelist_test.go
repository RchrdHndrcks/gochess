package chess

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/v2"
)

func TestMoveList_AddLenAtReset(t *testing.T) {
	var ml MoveList
	if ml.Len() != 0 {
		t.Fatalf("expected empty MoveList, got Len=%d", ml.Len())
	}
	m1 := NewMove(gochess.Coor(4, 6), gochess.Coor(4, 4), FlagDoublePush)
	m2 := NewMove(gochess.Coor(6, 7), gochess.Coor(5, 5), FlagQuiet)
	ml.Add(m1)
	ml.Add(m2)
	if ml.Len() != 2 {
		t.Fatalf("expected Len=2, got %d", ml.Len())
	}
	if ml.At(0) != m1 || ml.At(1) != m2 {
		t.Fatalf("At returned wrong moves")
	}
	ml.Reset()
	if ml.Len() != 0 {
		t.Fatalf("expected empty after Reset, got Len=%d", ml.Len())
	}
}

func TestMoveList_Swap(t *testing.T) {
	var ml MoveList
	m1 := NewMove(gochess.Coor(0, 6), gochess.Coor(0, 5), FlagQuiet)
	m2 := NewMove(gochess.Coor(1, 6), gochess.Coor(1, 5), FlagQuiet)
	ml.Add(m1)
	ml.Add(m2)
	ml.Swap(0, 1)
	if ml.At(0) != m2 || ml.At(1) != m1 {
		t.Fatalf("Swap did not swap entries")
	}
}

func TestMaxMoves(t *testing.T) {
	if MaxMoves != 256 {
		t.Fatalf("expected MaxMoves=256, got %d", MaxMoves)
	}
}

// expectedCaptures derives the expected set of captures by filtering the full
// legal-move list, independently of the staged-gen implementation. Used so
// the *Into tests do not just compare two wrappers around the same code path.
func expectedCaptures(c *Chess) []Move {
	all := c.Moves()
	out := make([]Move, 0, len(all))
	for _, m := range all {
		if m.IsCapture() {
			out = append(out, m)
		}
	}
	return out
}

func expectedQuiets(c *Chess) []Move {
	all := c.Moves()
	out := make([]Move, 0, len(all))
	for _, m := range all {
		if !m.IsCapture() {
			out = append(out, m)
		}
	}
	return out
}

func sameMoveSet(a, b []Move) bool {
	if len(a) != len(b) {
		return false
	}
	seen := map[Move]int{}
	for _, m := range a {
		seen[m]++
	}
	for _, m := range b {
		seen[m]--
		if seen[m] < 0 {
			return false
		}
	}
	return true
}

func TestCapturesInto_MatchesIndependentBaseline(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}
	want := expectedCaptures(c)
	var ml MoveList
	c.CapturesInto(&ml)
	got := make([]Move, ml.Len())
	for i := range got {
		got[i] = ml.At(i)
	}
	for _, m := range got {
		if !m.IsCapture() {
			t.Errorf("CapturesInto returned non-capture: %s", m.UCI())
		}
	}
	if !sameMoveSet(got, want) {
		t.Errorf("CapturesInto set mismatch: got %d, want %d", len(got), len(want))
	}
}

func TestQuietMovesInto_MatchesIndependentBaseline(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}
	want := expectedQuiets(c)
	var ml MoveList
	c.QuietMovesInto(&ml)
	got := make([]Move, ml.Len())
	for i := range got {
		got[i] = ml.At(i)
	}
	for _, m := range got {
		if m.IsCapture() {
			t.Errorf("QuietMovesInto returned capture: %s", m.UCI())
		}
	}
	if !sameMoveSet(got, want) {
		t.Errorf("QuietMovesInto set mismatch: got %d, want %d", len(got), len(want))
	}
}
