package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSAN(t *testing.T) {
	t.Run("PawnMoves", func(t *testing.T) {
		c, err := New(WithParallelism(1))
		require.NoError(t, err)

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
			san, err := c.SAN(tt.uci)
			require.NoError(t, err, "SAN(%s)", tt.uci)
			assert.Equal(t, tt.san, san, "SAN(%s)", tt.uci)
		}
	})

	t.Run("PieceMoves", func(t *testing.T) {
		c, err := New(WithParallelism(1))
		require.NoError(t, err)

		tests := []struct {
			uci string
			san string
		}{
			{"g1f3", "Nf3"},
			{"b1c3", "Nc3"},
		}

		for _, tt := range tests {
			san, err := c.SAN(tt.uci)
			require.NoError(t, err, "SAN(%s)", tt.uci)
			assert.Equal(t, tt.san, san, "SAN(%s)", tt.uci)
		}
	})

	t.Run("BishopMove", func(t *testing.T) {
		c, err := New(WithParallelism(1))
		require.NoError(t, err)

		require.NoError(t, c.MakeMove("e2e4"))
		require.NoError(t, c.MakeMove("e7e5"))

		san, err := c.SAN("f1c4")
		require.NoError(t, err)
		assert.Equal(t, "Bc4", san)
	})

	t.Run("PawnCapture", func(t *testing.T) {
		c, err := New(
			WithFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("e4d5")
		require.NoError(t, err)
		assert.Equal(t, "exd5", san)
	})

	t.Run("KnightCapture", func(t *testing.T) {
		c, err := New(
			WithFEN("r1bqkbnr/pppppppp/2n5/4P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("c6e5")
		require.NoError(t, err)
		assert.Equal(t, "Nxe5", san)
	})

	t.Run("CastlingKingside", func(t *testing.T) {
		c, err := New(
			WithFEN("rnbqk2r/ppppbppp/4pn2/8/4P3/5N2/PPPPBPPP/RNBQK2R w KQkq - 4 4"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("e1g1")
		require.NoError(t, err)
		assert.Equal(t, "O-O", san)
	})

	t.Run("CastlingQueenside", func(t *testing.T) {
		c, err := New(
			WithFEN("r3kbnr/pppqpppp/2n1b3/3p4/3P4/2N1B3/PPPQPPPP/R3KBNR w KQkq - 6 5"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("e1c1")
		require.NoError(t, err)
		assert.Equal(t, "O-O-O", san)
	})

	t.Run("Check", func(t *testing.T) {
		c, err := New(
			WithFEN("rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("d1h5")
		require.NoError(t, err)
		assert.Equal(t, "Qh5", san)
	})

	t.Run("CheckWithCapture", func(t *testing.T) {
		c, err := New(
			WithFEN("rnbqkbnr/pppp1ppp/8/4p2Q/4P3/8/PPPP1PPP/RNB1KBNR w KQkq - 0 2"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("h5f7")
		require.NoError(t, err)
		assert.Equal(t, "Qxf7+", san)
	})

	t.Run("ScholarsMate", func(t *testing.T) {
		c, err := New(WithParallelism(1))
		require.NoError(t, err)

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
			san, err := c.SAN(ms.uci)
			require.NoError(t, err, "SAN(%s)", ms.uci)
			assert.Equal(t, ms.san, san, "SAN(%s)", ms.uci)
			require.NoError(t, c.MakeMove(ms.uci), "MakeMove(%s)", ms.uci)
		}

		assert.True(t, c.IsCheckmate())
	})

	t.Run("Promotion", func(t *testing.T) {
		c, err := New(
			WithFEN("8/4P1k1/8/8/8/8/8/K7 w - - 0 1"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("e7e8q")
		require.NoError(t, err)
		assert.Equal(t, "e8=Q", san)
	})

	t.Run("DisambiguationByFile", func(t *testing.T) {
		// Two rooks on same rank — file disambiguation needed.
		c, err := New(
			WithFEN("7k/8/8/8/8/8/8/R4RK1 w - - 0 1"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("a1d1")
		require.NoError(t, err)
		assert.Equal(t, "Rad1", san)
	})

	t.Run("DisambiguationByRank", func(t *testing.T) {
		// Two rooks on same file — rank disambiguation needed.
		c, err := New(
			WithFEN("7k/8/8/8/8/R7/8/R3K3 w - - 0 1"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		san, err := c.SAN("a1a2")
		require.NoError(t, err)
		assert.Equal(t, "R1a2", san)
	})
}

func TestFromSAN(t *testing.T) {
	t.Run("PawnMoves", func(t *testing.T) {
		c, err := New(WithParallelism(1))
		require.NoError(t, err)

		tests := []struct {
			san string
			uci string
		}{
			{"e4", "e2e4"},
			{"d4", "d2d4"},
			{"a3", "a2a3"},
		}

		for _, tt := range tests {
			uci, err := c.FromSAN(tt.san)
			require.NoError(t, err, "FromSAN(%s)", tt.san)
			assert.Equal(t, tt.uci, uci, "FromSAN(%s)", tt.san)
		}
	})

	t.Run("PieceMoves", func(t *testing.T) {
		c, err := New(WithParallelism(1))
		require.NoError(t, err)

		uci, err := c.FromSAN("Nf3")
		require.NoError(t, err)
		assert.Equal(t, "g1f3", uci)
	})

	t.Run("Castling", func(t *testing.T) {
		c, err := New(
			WithFEN("rnbqk2r/ppppbppp/4pn2/8/4P3/5N2/PPPPBPPP/RNBQK2R w KQkq - 4 4"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		uci, err := c.FromSAN("O-O")
		require.NoError(t, err)
		assert.Equal(t, "e1g1", uci)
	})

	t.Run("Captures", func(t *testing.T) {
		c, err := New(
			WithFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		uci, err := c.FromSAN("exd5")
		require.NoError(t, err)
		assert.Equal(t, "e4d5", uci)
	})

	t.Run("Promotion", func(t *testing.T) {
		c, err := New(
			WithFEN("8/4P1k1/8/8/8/8/8/K7 w - - 0 1"),
			WithParallelism(1),
		)
		require.NoError(t, err)

		uci, err := c.FromSAN("e8=Q")
		require.NoError(t, err)
		assert.Equal(t, "e7e8q", uci)
	})

	t.Run("Roundtrip", func(t *testing.T) {
		c, err := New(WithParallelism(1))
		require.NoError(t, err)

		for _, uci := range c.AvailableMoves() {
			san, err := c.SAN(uci)
			require.NoError(t, err, "SAN(%s)", uci)

			roundtrip, err := c.FromSAN(san)
			require.NoError(t, err, "FromSAN(%s)", san)

			assert.Equal(t, uci, roundtrip, "roundtrip: %s -> %s -> %s", uci, san, roundtrip)
		}
	})
}
