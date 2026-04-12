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

func TestCapturesInto_MatchesCaptures(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}
	slice := c.Captures()
	var ml MoveList
	c.CapturesInto(&ml)
	if ml.Len() != len(slice) {
		t.Fatalf("CapturesInto Len=%d, Captures len=%d", ml.Len(), len(slice))
	}
	for i := range slice {
		if ml.At(i) != slice[i] {
			t.Errorf("mismatch at %d: ml=%s slice=%s", i, ml.At(i).UCI(), slice[i].UCI())
		}
	}
}

func TestQuietMovesInto_MatchesQuietMoves(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.LoadPosition(kiwipeteFEN); err != nil {
		t.Fatalf("LoadPosition: %v", err)
	}
	slice := c.QuietMoves()
	var ml MoveList
	c.QuietMovesInto(&ml)
	if ml.Len() != len(slice) {
		t.Fatalf("QuietMovesInto Len=%d, QuietMoves len=%d", ml.Len(), len(slice))
	}
	for i := range slice {
		if ml.At(i) != slice[i] {
			t.Errorf("mismatch at %d: ml=%s slice=%s", i, ml.At(i).UCI(), slice[i].UCI())
		}
	}
}
