package chess_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/chess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChess(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		c, err := chess.New()
		require.NotNil(t, c)
		require.Nil(t, err)
		assert.Equal(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", c.FEN())
	})

	t.Run("Custom FEN", func(t *testing.T) {
		c, err := chess.New(chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1"))
		require.NotNil(t, c)
		require.Nil(t, err)
		assert.Equal(t, "8/8/8/k7/8/K2P4/8/8 w - - 0 1", c.FEN())
	})

	t.Run("Invalid FEN", func(t *testing.T) {
		_, err := chess.New(chess.WithFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1BNR w KQkq - 0 1"))
		require.Error(t, err)
	})
}

func TestAvailableMoves(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New()
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{
			"a2a3", "a2a4", "b2b3", "b2b4", "c2c3", "c2c4", "d2d3", "d2d4",
			"e2e3", "e2e4", "f2f3", "f2f4", "g2g3", "g2g4", "h2h3", "h2h4",
			"b1a3", "b1c3", "g1f3", "g1h3",
		}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Custom FEN 1", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"a3b3", "a3b2", "a3a2", "d3d4"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Custom FEN 2", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/8/8/4k3/7r/5K2 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"f1e1", "f1g1"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Black Queenside Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("r3k3/p7/8/8/8/8/8/7K b q - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"a7a6", "a7a5", "a8b8", "a8c8", "a8d8", "e8d8", "e8f8",
			"e8d7", "e8e7", "e8f7", "e8c8"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Black Kingside Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("4k2r/7p/8/8/8/8/8/7K b k - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"e8d8", "e8f8", "e8d7", "e8e7", "e8f7", "e8g8", "h8f8",
			"h8g8", "h7h6", "h7h5"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Black Both Sides Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("r3k2r/p6p/8/8/8/8/8/7K b kq - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"e8d8", "e8f8", "e8d7", "e8e7", "e8f7", "e8g8", "e8c8",
			"h8f8", "h8g8", "h7h6", "h7h5", "a7a6", "a7a5", "a8b8", "a8c8", "a8d8"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("White Kingside Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("7k/8/8/8/8/8/7P/4K2R w K - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "e1g1", "h1g1",
			"h1f1", "h2h3", "h2h4"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("White Queenside Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("7k/8/8/8/8/8/P7/R3K3 w Q - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "e1c1", "a1b1",
			"a1c1", "a1d1", "a2a3", "a2a4"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("White Both Sides Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("7k/8/8/8/8/8/P6P/R3K2R w KQ - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "e1c1", "e1g1",
			"a1b1", "a1c1", "a1d1", "a2a3", "a2a4", "h1g1", "h1f1", "h2h3", "h2h4"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("White Has No Castling Rights But Black Has", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("7k/8/8/8/8/8/P7/R3K3 w kq - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "a1b1",
			"a1c1", "a1d1", "a2a3", "a2a4"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Black Has No Castling Rights But White Has", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("r3k2r/p6p/8/8/8/8/8/7K b KQ - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"e8d8", "e8f8", "e8d7", "e8e7", "e8f7", "h8f8", "h8g8",
			"h7h6", "h7h5", "a7a6", "a7a5", "a8b8", "a8c8", "a8d8"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Castle Way Blocked", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/3b4/8/4K2R w K - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"h1g1", "h1f1", "h1h2", "h1h3", "h1h4", "h1h5", "h1h6",
			"h1h7", "h1h8", "e1d1", "e1d2", "e1f2"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Rook Under Attack In Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/7r/8/4K2R w K - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"h1g1", "h1f1", "h1h2", "h1h3", "e1d1", "e1d2", "e1e2",
			"e1f2", "e1f1", "e1g1"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Castle Is Not Possible", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/7r/8/8/4K2R w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"h1g1", "h1f1", "h1h2", "h1h3", "h1h4", "e1d1", "e1d2",
			"e1e2", "e1f2", "e1f1"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("King Is In Check", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/8/3q4/3K4 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"d1d2"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("King Is In Checkmate", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/8/3qr3/3K4 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		assert.Nil(t, moves)
	})

	t.Run("King Is In Stalemate", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/8/2q5/K7 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{}
		assert.ElementsMatch(t, expectedMoves, moves)
	})

	t.Run("Pawn has promotion with capture move", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k6r/6P1/8/8/8/8/8/K7 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		moves := c.AvailableMoves()

		// Assert
		expectedMoves := []string{"g7g8q", "g7g8r", "g7g8b", "g7g8n", "g7h8q", "g7h8r", "g7h8b", "g7h8n",
			"a1a2", "a1b1", "a1b2"}
		assert.ElementsMatch(t, expectedMoves, moves)
	})
}

func TestFEN(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New()
		require.Nil(t, errOpts)

		// Act
		fen := c.FEN()

		// Assert
		require.Equal(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", fen)
	})
}

func TestMakeMove(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New()
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("e2e4")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", c.FEN())
	})

	t.Run("Custom FEN 1", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("d3d4")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "8/8/8/k7/3P4/K7/8/8 b - - 0 1", c.FEN())
	})

	t.Run("Custom FEN 2", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/8/8/4k3/7r/5K2 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("f1e1")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "8/8/8/8/8/4k3/7r/4K3 b - - 1 1", c.FEN())
	})

	t.Run("Castle", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/8/8/4K2R w K - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("e1g1")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "k7/8/8/8/8/8/8/5RK1 b - - 1 1", c.FEN())
	})

	t.Run("Capture", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/8/3q4/3K4 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("d1d2")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "k7/8/8/8/8/8/3K4/8 b - - 0 1", c.FEN())
	})

	t.Run("Promotion Queen", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("h7h8q")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "k6Q/8/8/8/8/8/8/7K b - - 0 1", c.FEN())
	})

	t.Run("Promotion Rook", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("h7h8r")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "k6R/8/8/8/8/8/8/7K b - - 0 1", c.FEN())
	})

	t.Run("Promotion Bishop", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("h7h8b")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "k6B/8/8/8/8/8/8/7K b - - 0 1", c.FEN())
	})

	t.Run("Promotion Knight", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("h7h8n")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "k6N/8/8/8/8/8/8/7K b - - 0 1", c.FEN())
	})

	t.Run("Promotion King Error", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("h7h8k")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "move is not legal: h7h8k", err.Error())
	})

	t.Run("En Passant", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("7k/8/8/3pP3/8/8/8/7K w - d6 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("e5d6")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "7k/8/3P4/8/8/8/8/7K b - - 0 1", c.FEN())
	})

	t.Run("Invalid Move - Invalid Square", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		err := c.MakeMove("d3d9")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "move is not legal: d3d9", err.Error())
	})
}

func TestIsCheck(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New()
		require.Nil(t, errOpts)

		// Assert
		assert.False(t, c.IsCheck())
	})

	t.Run("King Is In Check", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/8/3q4/3K4 w - - 0 1"))
		require.Nil(t, errOpts)

		// Assert
		assert.True(t, c.IsCheck())
	})

	t.Run("King Is Not In Check", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("k7/8/8/8/8/8/8/3K4 w - - 0 1"))
		require.Nil(t, errOpts)

		// Assert
		assert.False(t, c.IsCheck())
	})
}

func TestSquare(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New()
		require.Nil(t, errOpts)

		// Act
		piece, err := c.Square("e2")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "P", piece)
	})

	t.Run("Custom FEN", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		piece, err := c.Square("d3")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "P", piece)
	})

	t.Run("Empty Square", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		piece, err := c.Square("a1")

		// Assert
		require.Nil(t, err)
		assert.Equal(t, "", piece)
	})

	t.Run("Invalid Square", func(t *testing.T) {
		// Arrange
		c, errOpts := chess.New(chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1"))
		require.Nil(t, errOpts)

		// Act
		_, err := c.Square("a9")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "failed to convert algebraic notation to coordinate: coordinate out of bounds", err.Error())
	})
}

func TestMakeMove_ScholarMate(t *testing.T) {
	// Arrange
	c, err := chess.New()
	require.Nil(t, err)

	// Act
	err = c.MakeMove("e2e4")
	require.Nil(t, err)
	err = c.MakeMove("e7e5")
	require.Nil(t, err)
	err = c.MakeMove("f1c4")
	require.Nil(t, err)
	err = c.MakeMove("b8c6")
	require.Nil(t, err)
	err = c.MakeMove("d1h5")
	require.Nil(t, err)
	err = c.MakeMove("g8f6")
	require.Nil(t, err)
	err = c.MakeMove("h5f7")

	// Assert
	require.Nil(t, err)
	assert.Equal(t, "r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4", c.FEN())
}

func TestMakeMove_CapablancaSteiner(t *testing.T) {
	// Arrange
	c, err := chess.New()
	require.Nil(t, err)

	moves := []string{
		"e2e4", "e7e5", "g1f3", "b8c6", "b1c3", "g8f6", "f1b5", "f8b4",
		"e1g1", "e8g8", "d2d3", "d7d6", "c1g5", "b4c3", "b2c3", "c6e7",
		"f3h4", "c7c6", "b5c4", "c8e6", "g5f6", "g7f6", "c4e6", "f7e6",
		"d1g4", "g8f7", "f2f4", "f8g8", "g4h5", "f7g7", "f4e5", "d6e5",
		"f1f6", "g7f6", "a1f1", "e7f5", "h4f5", "e6f5", "f1f5", "f6e7",
		"h5f7", "e7d6", "f5f6", "d6c5", "f7b7", "d8b6", "f6c6", "b6c6",
		"b7b4",
	}

	// Act
	for _, move := range moves {
		err = c.MakeMove(move)
		require.Nil(t, err)
	}

	// Assert
	assert.Equal(t, "r5r1/p6p/2q5/2k1p3/1Q2P3/2PP4/P1P3PP/6K1 b - - 1 25", c.FEN())
}

func TestLoadPosition_Errors(t *testing.T) {
	t.Run("Invalid FEN", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("invalid")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid", err.Error())
	})

	t.Run("Invalid Number of Properties 1", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("8/8/8/8/8/8/8/8")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: 8/8/8/8/8/8/8/8", err.Error())
	})

	t.Run("Invalid Number of Properties 2", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("8/8/8/8/8/8/8/8 w")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: 8/8/8/8/8/8/8/8 w", err.Error())
	})

	t.Run("Invalid FEN - 9 Pawns", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/ppppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: rnbqkbnr/ppppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", err.Error())
	})

	t.Run("Invalid FEN - 1 Pawn and 8 Empty Squares", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/p8/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: rnbqkbnr/p8/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", err.Error())
	})

	t.Run("Invalid FEN - Invalid Color", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR o KQkq - 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid color: o", err.Error())
	})

	t.Run("Invalid FEN - Invalid Castles", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w QQQQ - 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid castles: QQQQ", err.Error())
	})

	t.Run("Invalid FEN - Invalid En Passant square", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e3 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid en passant square: e3", err.Error())
	})

	t.Run("Invalid FEN - Invalid En Passant square", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq h7 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid en passant square: h7", err.Error())
	})

	t.Run("Invalid FEN - Invalid En Passant square", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq h9 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid en passant square: h9", err.Error())
	})

	t.Run("Invalid FEN - Invalid Half Moves", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - J 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid half moves: J", err.Error())
	})

	t.Run("Invalid FEN - Invalid Moves Count", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 Q")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: invalid moves count: Q", err.Error())
	})

	t.Run("No White King", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/4P3/8/PPPPPPPP/RNBQ1BNR b KQkq e3 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: both kings must be in the board once", err.Error())
	})

	t.Run("No Black King", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbq1bnr/pppppppp/8/8/4P3/8/PPPPPPPP/RNBQKBNR b KQkq e3 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: both kings must be in the board once", err.Error())
	})

	t.Run("Two White Kings", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)

		// Act
		err = c.LoadPosition("rnbqkbnr/pppppppp/8/8/4P3/8/PPPPPPPP/RNBKKBNR b KQkq e3 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: both kings must be in the board once", err.Error())
	})

	t.Run("Position is not legal", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		require.NoError(t, err)
		copy := *c

		// Act
		err = c.LoadPosition("k7/8/8/8/8/8/7r/7K b - - 0 1")

		// Assert
		require.NotNil(t, err)
		assert.Equal(t, "invalid FEN: the current turn can capture the opponent king", err.Error())
		assert.Equal(t, copy, *c)
	})
}

func TestUnmakeMove(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		if err != nil {
			t.Errorf("failed to create chess game: %s", err.Error())
		}

		previousFEN := c.FEN()
		err = c.MakeMove("e2e4")
		if err != nil {
			t.Errorf("failed to make move: %s", err.Error())
		}

		// Act
		c.UnmakeMove()

		// Assert
		assert.Equal(t, previousFEN, c.FEN())
	})

	t.Run("No Moves", func(t *testing.T) {
		// Arrange
		c, err := chess.New()
		if err != nil {
			t.Errorf("failed to create chess game: %s", err.Error())
		}

		previousFEN := c.FEN()

		// Act
		c.UnmakeMove()

		// Assert
		assert.Equal(t, previousFEN, c.FEN())
	})
}
