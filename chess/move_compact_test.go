package chess

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/v2"
)

func TestMove_NullMove(t *testing.T) {
	t.Run("NullMove is zero", func(t *testing.T) {
		if NullMove != 0 {
			t.Fatalf("NullMove = %d, want 0", NullMove)
		}
	})
}

func TestMove_RoundTrip(t *testing.T) {
	cases := []struct {
		name      string
		from, to  gochess.Coordinate
		flag      MoveFlag
		promo     gochess.Piece
		captured  gochess.Piece
		uci       string
		isCapture bool
	}{
		{
			name: "quiet pawn push e2e4",
			from: mustAlg(t, "e2"), to: mustAlg(t, "e4"),
			flag: FlagDoublePush, uci: "e2e4",
		},
		{
			name: "quiet knight g1f3",
			from: mustAlg(t, "g1"), to: mustAlg(t, "f3"),
			flag: FlagQuiet, uci: "g1f3",
		},
		{
			name: "capture exd5",
			from: mustAlg(t, "e4"), to: mustAlg(t, "d5"),
			flag: FlagCapture, captured: gochess.Pawn,
			uci: "e4d5", isCapture: true,
		},
		{
			name: "en passant",
			from: mustAlg(t, "e5"), to: mustAlg(t, "d6"),
			flag: FlagEnPassant, captured: gochess.Pawn,
			uci: "e5d6", isCapture: true,
		},
		{
			name: "kingside castle",
			from: mustAlg(t, "e1"), to: mustAlg(t, "g1"),
			flag: FlagCastle, uci: "e1g1",
		},
		{
			name: "promotion to queen",
			from: mustAlg(t, "e7"), to: mustAlg(t, "e8"),
			flag: FlagPromotion, promo: gochess.Queen,
			uci: "e7e8q",
		},
		{
			name: "promotion-capture to knight",
			from: mustAlg(t, "e7"), to: mustAlg(t, "d8"),
			flag: FlagPromotionCapture, promo: gochess.Knight, captured: gochess.Rook,
			uci: "e7d8n", isCapture: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m := NewMove(tc.from, tc.to, tc.flag)
			if tc.promo != gochess.Empty {
				m = m.WithPromotion(tc.promo)
			}
			if tc.captured != gochess.Empty {
				m = m.WithCapturedPiece(tc.captured)
			}

			if got := m.From(); got != tc.from {
				t.Errorf("From = %v, want %v", got, tc.from)
			}
			if got := m.To(); got != tc.to {
				t.Errorf("To = %v, want %v", got, tc.to)
			}
			if got := m.Flags(); got != tc.flag {
				t.Errorf("Flags = %d, want %d", got, tc.flag)
			}
			if got := m.Promotion(); got != tc.promo {
				t.Errorf("Promotion = %v, want %v", got, tc.promo)
			}
			if got := m.CapturedPiece(); got != tc.captured {
				t.Errorf("CapturedPiece = %v, want %v", got, tc.captured)
			}
			if got := m.IsCapture(); got != tc.isCapture {
				t.Errorf("IsCapture = %v, want %v", got, tc.isCapture)
			}
			if got := m.UCI(); got != tc.uci {
				t.Errorf("UCI = %q, want %q", got, tc.uci)
			}
			if m.GivesCheck() {
				t.Errorf("GivesCheck = true, want false (not set)")
			}

			// GivesCheck round-trip.
			withCheck := m.WithGivesCheck(true)
			if !withCheck.GivesCheck() {
				t.Errorf("WithGivesCheck(true).GivesCheck() = false, want true")
			}
			// Verify other fields preserved.
			if withCheck.From() != tc.from || withCheck.To() != tc.to ||
				withCheck.Flags() != tc.flag || withCheck.Promotion() != tc.promo ||
				withCheck.CapturedPiece() != tc.captured {
				t.Errorf("WithGivesCheck mutated other fields")
			}
		})
	}
}

func TestMove_IsCapture(t *testing.T) {
	cases := []struct {
		flag MoveFlag
		want bool
	}{
		{FlagQuiet, false},
		{FlagDoublePush, false},
		{FlagCastle, false},
		{FlagPromotion, false},
		{FlagCapture, true},
		{FlagEnPassant, true},
		{FlagPromotionCapture, true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			m := NewMove(mustAlg(t, "a1"), mustAlg(t, "a2"), tc.flag)
			if got := m.IsCapture(); got != tc.want {
				t.Errorf("flag %d: IsCapture = %v, want %v", tc.flag, got, tc.want)
			}
		})
	}
}

func TestMove_SquareCoordinateRoundTrip(t *testing.T) {
	t.Run("a1 == 0", func(t *testing.T) {
		if got := squareFromCoordinate(mustAlg(t, "a1")); got != 0 {
			t.Errorf("a1 = %d, want 0", got)
		}
	})
	t.Run("h8 == 63", func(t *testing.T) {
		if got := squareFromCoordinate(mustAlg(t, "h8")); got != 63 {
			t.Errorf("h8 = %d, want 63", got)
		}
	})
	t.Run("all 64 squares round-trip", func(t *testing.T) {
		for sq := uint32(0); sq < 64; sq++ {
			c := coordinateFromSquare(sq)
			if got := squareFromCoordinate(c); got != sq {
				t.Errorf("sq %d -> %v -> %d", sq, c, got)
			}
		}
	})
}

func mustAlg(t *testing.T, s string) gochess.Coordinate {
	t.Helper()
	c, err := AlgebraicToCoordinate(s)
	if err != nil {
		t.Fatalf("AlgebraicToCoordinate(%q): %v", s, err)
	}
	return c
}
