package chess_test

import (
	"strings"
	"testing"

	"github.com/RchrdHndrcks/gochess/v2/chess"
	chesspgn "github.com/RchrdHndrcks/gochess/v2/chess/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGN(t *testing.T) {
	t.Run("Empty game with default tags", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		pgn := c.PGN(chesspgn.PGNTags{})

		assert.Contains(t, pgn, `[Event "?"]`)
		assert.Contains(t, pgn, `[Site "?"]`)
		assert.Contains(t, pgn, `[Date "?"]`)
		assert.Contains(t, pgn, `[Round "?"]`)
		assert.Contains(t, pgn, `[White "?"]`)
		assert.Contains(t, pgn, `[Black "?"]`)
		assert.Contains(t, pgn, `[Result "*"]`)
		parsedTags, parsedMoves, parseErr := chesspgn.Parse(pgn)
		require.NoError(t, parseErr)
		assert.Equal(t, chesspgn.ResultOngoing, parsedTags.Result)
		assert.Empty(t, parsedMoves)
	})

	t.Run("Empty game with custom tags", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		tags := chesspgn.PGNTags{
			Event: "Test Tournament",
			Site:  "Internet",
			Date:  "2026.03.17",
			Round: "1",
			White: "Player1",
			Black: "Player2",
		}
		pgn := c.PGN(tags)

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

		pgn := c.PGN(chesspgn.PGNTags{})

		parsedTags, parsedMoves, parseErr := chesspgn.Parse(pgn)
		require.NoError(t, parseErr)
		assert.Equal(t, chesspgn.ResultWhiteWins, parsedTags.Result)
		assert.Equal(t, []string{"e2e4", "e7e5", "f1c4", "b8c6", "d1h5", "g8f6", "h5f7"}, parsedMoves)
	})

	t.Run("Line wrapping", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		pgn := c.PGN(chesspgn.PGNTags{})

		for _, line := range strings.Split(pgn, "\n") {
			assert.LessOrEqual(t, len(line), 80, "line exceeds 80 characters: %s", line)
		}
	})

	t.Run("Newlines in tag values are stripped", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		tags := chesspgn.PGNTags{
			Event: "line1\nline2",
			Site:  "cr\r\ntest",
		}
		pgn := c.PGN(tags)

		// Neither the event nor the site tag line should contain a raw newline
		// or carriage return inside the quoted value.
		assert.Contains(t, pgn, `[Event "line1line2"]`)
		assert.Contains(t, pgn, `[Site "crtest"]`)
	})

	t.Run("Tag values with special characters", func(t *testing.T) {
		c, err := chess.New()
		require.NoError(t, err)

		tags := chesspgn.PGNTags{
			Event: "He said \"hello\"",
			Site:  "path\\to\\file",
		}
		pgn := c.PGN(tags)

		assert.Contains(t, pgn, "[Event \"He said \\\"hello\\\"\"]")
		assert.Contains(t, pgn, "[Site \"path\\\\to\\\\file\"]")

		// Roundtrip: parse back and verify unescaped values.
		parsedTags, _, err := chesspgn.Parse(pgn)
		require.NoError(t, err)
		assert.Equal(t, "He said \"hello\"", parsedTags.Event)
		assert.Equal(t, "path\\to\\file", parsedTags.Site)
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

		tags := chesspgn.PGNTags{
			Event: "Roundtrip Test",
			White: "W",
			Black: "B",
		}
		pgn := c.PGN(tags)

		parsedTags, parsedMoves, err := chesspgn.Parse(pgn)
		require.NoError(t, err)

		assert.Equal(t, "Roundtrip Test", parsedTags.Event)
		assert.Equal(t, "W", parsedTags.White)
		assert.Equal(t, "B", parsedTags.Black)
		assert.Equal(t, "*", parsedTags.Result)
		assert.Equal(t, playedMoves, parsedMoves)
	})
}
