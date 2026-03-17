package chess_test

import (
	"strings"
	"testing"

	"github.com/RchrdHndrcks/gochess/chess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToPGN(t *testing.T) {
	t.Run("Empty game with default tags", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		pgn := c.ToPGN(chess.PGNTags{})

		assert.Contains(t, pgn, `[Event "?"]`)
		assert.Contains(t, pgn, `[Site "?"]`)
		assert.Contains(t, pgn, `[Date "?"]`)
		assert.Contains(t, pgn, `[Round "?"]`)
		assert.Contains(t, pgn, `[White "?"]`)
		assert.Contains(t, pgn, `[Black "?"]`)
		assert.Contains(t, pgn, `[Result "*"]`)
		assert.True(t, strings.HasSuffix(strings.TrimSpace(pgn), "*"))
	})

	t.Run("Empty game with custom tags", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		tags := chess.PGNTags{
			Event: "Test Tournament",
			Site:  "Internet",
			Date:  "2026.03.17",
			Round: "1",
			White: "Player1",
			Black: "Player2",
		}
		pgn := c.ToPGN(tags)

		assert.Contains(t, pgn, `[Event "Test Tournament"]`)
		assert.Contains(t, pgn, `[Site "Internet"]`)
		assert.Contains(t, pgn, `[Date "2026.03.17"]`)
		assert.Contains(t, pgn, `[Round "1"]`)
		assert.Contains(t, pgn, `[White "Player1"]`)
		assert.Contains(t, pgn, `[Black "Player2"]`)
		assert.Contains(t, pgn, `[Result "*"]`)
	})

	t.Run("Scholar's mate", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		moves := []string{"e2e4", "e7e5", "f1c4", "b8c6", "d1h5", "g8f6", "h5f7"}
		for _, m := range moves {
			require.NoError(t, c.MakeMove(m))
		}

		require.True(t, c.IsCheckmate())

		pgn := c.ToPGN(chess.PGNTags{})

		assert.Contains(t, pgn, `[Result "1-0"]`)
		assert.Contains(t, pgn, "1. e2e4 e7e5")
		assert.Contains(t, pgn, "2. f1c4 b8c6")
		assert.Contains(t, pgn, "3. d1h5 g8f6")
		assert.Contains(t, pgn, "4. h5f7")
		assert.True(t, strings.HasSuffix(strings.TrimSpace(pgn), "1-0"))
	})

	t.Run("Line wrapping", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		pgn := c.ToPGN(chess.PGNTags{})

		for _, line := range strings.Split(pgn, "\n") {
			assert.LessOrEqual(t, len(line), 80, "line exceeds 80 characters: %s", line)
		}
	})
}

func TestParsePGN(t *testing.T) {
	t.Run("Simple PGN", func(t *testing.T) {
		pgn := `[Event "Test"]
[Site "Internet"]
[Date "2026.03.17"]
[Round "1"]
[White "Alice"]
[Black "Bob"]
[Result "1-0"]

1. e2e4 e7e5 2. f1c4 b8c6 3. d1h5 g8f6 4. h5f7 1-0
`
		tags, moves, err := chess.ParsePGN(pgn)
		require.NoError(t, err)

		assert.Equal(t, "Test", tags.Event)
		assert.Equal(t, "Internet", tags.Site)
		assert.Equal(t, "2026.03.17", tags.Date)
		assert.Equal(t, "1", tags.Round)
		assert.Equal(t, "Alice", tags.White)
		assert.Equal(t, "Bob", tags.Black)
		assert.Equal(t, "1-0", tags.Result)

		expectedMoves := []string{"e2e4", "e7e5", "f1c4", "b8c6", "d1h5", "g8f6", "h5f7"}
		assert.Equal(t, expectedMoves, moves)
	})

	t.Run("PGN with comments and variations", func(t *testing.T) {
		pgn := `[Event "?"]
[Result "*"]

1. e2e4 {Best move} e7e5 (1... d7d5 2. e4d5) 2. g1f3 $1 b8c6 *
`
		tags, moves, err := chess.ParsePGN(pgn)
		require.NoError(t, err)

		assert.Equal(t, "?", tags.Event)
		assert.Equal(t, "*", tags.Result)

		expectedMoves := []string{"e2e4", "e7e5", "g1f3", "b8c6"}
		assert.Equal(t, expectedMoves, moves)
	})

	t.Run("Empty PGN", func(t *testing.T) {
		pgn := `[Event "?"]
[Site "?"]
[Date "?"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]

*
`
		tags, moves, err := chess.ParsePGN(pgn)
		require.NoError(t, err)

		assert.Equal(t, "?", tags.Event)
		assert.Equal(t, "*", tags.Result)
		assert.Empty(t, moves)
	})
}

func TestPGNRoundtrip(t *testing.T) {
	t.Run("Play moves export parse verify", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		playedMoves := []string{"e2e4", "e7e5", "g1f3", "b8c6", "f1c4", "f8c5"}
		for _, m := range playedMoves {
			require.NoError(t, c.MakeMove(m))
		}

		tags := chess.PGNTags{
			Event: "Roundtrip Test",
			White: "W",
			Black: "B",
		}
		pgn := c.ToPGN(tags)

		parsedTags, parsedMoves, err := chess.ParsePGN(pgn)
		require.NoError(t, err)

		assert.Equal(t, "Roundtrip Test", parsedTags.Event)
		assert.Equal(t, "W", parsedTags.White)
		assert.Equal(t, "B", parsedTags.Black)
		assert.Equal(t, "*", parsedTags.Result)
		assert.Equal(t, playedMoves, parsedMoves)
	})
}
