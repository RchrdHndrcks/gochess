package pgn_test

import (
	"testing"

	chesspgn "github.com/RchrdHndrcks/gochess/chess/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
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
		tags, moves, err := chesspgn.Parse(pgn)
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

	t.Run("PGN with brace comments and variations", func(t *testing.T) {
		pgn := `[Event "?"]
[Result "*"]

1. e2e4 {Best move} e7e5 (1... d7d5 2. e4d5) 2. g1f3 $1 b8c6 *
`
		tags, moves, err := chesspgn.Parse(pgn)
		require.NoError(t, err)

		assert.Equal(t, "?", tags.Event)
		assert.Equal(t, "*", tags.Result)

		expectedMoves := []string{"e2e4", "e7e5", "g1f3", "b8c6"}
		assert.Equal(t, expectedMoves, moves)
	})

	t.Run("PGN with semicolon comments", func(t *testing.T) {
		pgn := `[Event "?"]
[Result "*"]

1. e2e4 e7e5 ; this is a comment
2. g1f3 b8c6 *
`
		_, moves, err := chesspgn.Parse(pgn)
		require.NoError(t, err)

		expectedMoves := []string{"e2e4", "e7e5", "g1f3", "b8c6"}
		assert.Equal(t, expectedMoves, moves)
	})

	t.Run("Empty move text", func(t *testing.T) {
		pgn := `[Event "?"]
[Site "?"]
[Date "?"]
[Round "?"]
[White "?"]
[Black "?"]
[Result "*"]

*
`
		tags, moves, err := chesspgn.Parse(pgn)
		require.NoError(t, err)

		assert.Equal(t, "?", tags.Event)
		assert.Equal(t, "*", tags.Result)
		assert.Empty(t, moves)
	})

	t.Run("Tag values with special characters roundtrip", func(t *testing.T) {
		// Manually construct escaped PGN and verify unescaping.
		pgn := "[Event \"He said \\\"hello\\\"\"]\n[Result \"*\"]\n\n*\n"
		tags, _, err := chesspgn.Parse(pgn)
		require.NoError(t, err)
		assert.Equal(t, "He said \"hello\"", tags.Event)
	})
}
