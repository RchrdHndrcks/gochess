package chess_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/pkg/chess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChess has a test table where each test case is a sub-test.
func TestChess(t *testing.T) {
	tests := []struct {
		name   string
		opts   []chess.Option
		FEN    string
		errMsg string
	}{
		{
			name: "Default",
			opts: []chess.Option{},
			FEN:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name: "Custom FEN",
			opts: []chess.Option{chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1")},
			FEN:  "8/8/8/k7/8/K2P4/8/8 w - - 0 1",
		},
		{
			name:   "Invalid FEN",
			opts:   []chess.Option{chess.WithFEN("invalid")},
			errMsg: "failed to apply option: failed to load FEN: invalid FEN: invalid",
		},
		{
			name: "Invalid FEN - row too short",
			opts: []chess.Option{chess.WithFEN("8/8/8/8/1P5/8/8 w - - 0 1")},
			errMsg: "failed to apply option: failed to load FEN: invalid FEN: " +
				"8/8/8/8/1P5/8/8 w - - 0 1",
		},
		{
			name:   "Invalid FEN - invalid number of properties",
			opts:   []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8")},
			errMsg: "failed to apply option: failed to load FEN: invalid FEN: 8/8/8/8/8/8/8/8",
		},
		{
			name:   "Invalid FEN - invalid number of properties",
			opts:   []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 w")},
			errMsg: "failed to apply option: failed to load FEN: invalid FEN: 8/8/8/8/8/8/8/8 w",
		},
		{
			name:   "Invalid FEN - invalid color",
			opts:   []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 x KQkq - 0 1")},
			errMsg: "failed to apply option: failed to set properties: invalid FEN color: x",
		},
		{
			name: "Invalid FEN - invalid castling",
			opts: []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 w KGkq - 0 1")},
			errMsg: "failed to apply option: failed to set properties: invalid FEN castles:" +
				" invalid castle: KGkq",
		},
		{
			name: "Invalid FEN - invalid in passant square - invalid len of square",
			opts: []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 w KQkq ab3 0 1")},
			errMsg: "failed to apply option: failed to set properties: invalid FEN in passant " +
				"square: invalid in passant square: ab3",
		},
		{
			name: "Invalid FEN - invalid in passant square - invalid square column",
			opts: []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 w KQkq j2 0 1")},
			errMsg: "failed to apply option: failed to set properties: invalid FEN in passant " +
				"square: invalid in passant square: j2",
		},
		{
			name: "Invalid FEN - invalid in passant square - invalid square row",
			opts: []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 w KQkq a9 0 1")},
			errMsg: "failed to apply option: failed to set properties: invalid FEN in passant " +
				"square: invalid in passant square: a9",
		},
		{
			name:   "Invalid FEN - invalid half moves count - not a number",
			opts:   []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 w KQkq - a 1")},
			errMsg: "failed to apply option: failed to set properties: invalid FEN half moves: a",
		},
		{
			name:   "Invalid FEN - invalid moves count - not a number",
			opts:   []chess.Option{chess.WithFEN("8/8/8/8/8/8/8/8 w KQkq - 0 a")},
			errMsg: "failed to apply option: failed to set properties: invalid FEN moves count: a",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := chess.NewChess(test.opts...)
			if test.errMsg != "" {
				require.NotNil(t, err)
				assert.Equal(t, test.errMsg, err.Error())
				return
			}

			require.NotNil(t, c)
			require.Nil(t, err)
			assert.Equal(t, test.FEN, c.FEN())
		})
	}
}

func TestAvailableLegalMoves(t *testing.T) {
	tests := []struct {
		name   string
		opts   []chess.Option
		moves  []string
		errMsg string
	}{
		{
			name: "Default",
			opts: []chess.Option{},
			moves: []string{
				"a2a3", "a2a4", "b2b3", "b2b4", "c2c3", "c2c4", "d2d3", "d2d4",
				"e2e3", "e2e4", "f2f3", "f2f4", "g2g3", "g2g4", "h2h3", "h2h4",
				"b1a3", "b1c3", "g1f3", "g1h3",
			},
		},
		{
			name:  "Custom FEN 1",
			opts:  []chess.Option{chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1")},
			moves: []string{"a3b3", "a3b2", "a3a2", "d3d4"},
		},
		{
			name:  "Custom FEN 2",
			opts:  []chess.Option{chess.WithFEN("8/8/8/8/8/4k3/7r/5K2 w - - 0 1")},
			moves: []string{"f1e1", "f1g1"},
		},
		{
			name: "Custom FEN - black queenside castle",
			opts: []chess.Option{chess.WithFEN("r3k3/p7/8/8/8/8/8/7K b q - 0 1")},
			moves: []string{"a7a6", "a7a5", "a8b8", "a8c8", "a8d8", "e8d8", "e8f8",
				"e8d7", "e8e7", "e8f7", "e8c8"},
		},
		{
			name: "Custom FEN - black kingside castle",
			opts: []chess.Option{chess.WithFEN("4k2r/7p/8/8/8/8/8/7K b k - 0 1")},
			moves: []string{"e8d8", "e8f8", "e8d7", "e8e7", "e8f7", "e8g8", "h8f8",
				"h8g8", "h7h6", "h7h5"},
		},
		{
			name: "Custom FEN - black both sides castle",
			opts: []chess.Option{chess.WithFEN("r3k2r/p6p/8/8/8/8/8/7K b kq - 0 1")},
			moves: []string{"e8d8", "e8f8", "e8d7", "e8e7", "e8f7", "e8g8", "e8c8",
				"h8f8", "h8g8", "h7h6", "h7h5", "a7a6", "a7a5", "a8b8", "a8c8", "a8d8"},
		},
		{
			name: "Custom FEN - white kingside castle",
			opts: []chess.Option{chess.WithFEN("7k/8/8/8/8/8/7P/4K2R w K - 0 1")},
			moves: []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "e1g1", "h1g1",
				"h1f1", "h2h3", "h2h4"},
		},
		{
			name: "Custom FEN - white queenside castle",
			opts: []chess.Option{chess.WithFEN("7k/8/8/8/8/8/P7/R3K3 w Q - 0 1")},
			moves: []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "e1c1", "a1b1",
				"a1c1", "a1d1", "a2a3", "a2a4"},
		},
		{
			name: "Custom FEN - white both sides castle",
			opts: []chess.Option{chess.WithFEN("7k/8/8/8/8/8/P6P/R3K2R w KQ - 0 1")},
			moves: []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "e1c1", "e1g1",
				"a1b1", "a1c1", "a1d1", "a2a3", "a2a4", "h1g1", "h1f1", "h2h3", "h2h4"},
		},
		{
			name: "Custom FEN - white has not castling rights but black has",
			opts: []chess.Option{chess.WithFEN("7k/8/8/8/8/8/P7/R3K3 w kq - 0 1")},
			moves: []string{"e1d1", "e1d2", "e1e2", "e1f2", "e1f1", "a1b1",
				"a1c1", "a1d1", "a2a3", "a2a4"},
		},
		{
			name: "Custom FEN - black has not castling rights but white has",
			opts: []chess.Option{chess.WithFEN("r3k2r/p6p/8/8/8/8/8/7K b KQ - 0 1")},
			moves: []string{"e8d8", "e8f8", "e8d7", "e8e7", "e8f7", "h8f8", "h8g8",
				"h7h6", "h7h5", "a7a6", "a7a5", "a8b8", "a8c8", "a8d8"},
		},
		{
			name: "Custom FEN - castle way blocked",
			opts: []chess.Option{chess.WithFEN("k7/8/8/8/8/3b4/8/4K2R w K - 0 1")},
			moves: []string{"h1g1", "h1f1", "h1h2", "h1h3", "h1h4", "h1h5", "h1h6",
				"h1h7", "h1h8", "e1d1", "e1d2", "e1f2"},
		},
		{
			name: "Custom FEN - rook under attack in castle",
			opts: []chess.Option{chess.WithFEN("k7/8/8/8/8/7r/8/4K2R w K - 0 1")},
			moves: []string{"h1g1", "h1f1", "h1h2", "h1h3", "e1d1", "e1d2", "e1e2",
				"e1f2", "e1f1", "e1g1"},
		},
		{
			name: "Custom FEN - castle is not possible",
			opts: []chess.Option{chess.WithFEN("k7/8/8/8/7r/8/8/4K2R w - - 0 1")},
			moves: []string{"h1g1", "h1f1", "h1h2", "h1h3", "h1h4", "e1d1", "e1d2",
				"e1e2", "e1f2", "e1f1"},
		},
		{
			name:  "Custom FEN - king is in check",
			opts:  []chess.Option{chess.WithFEN("k7/8/8/8/8/8/3q4/3K4 w - - 0 1")},
			moves: []string{"d1d2"},
		},
		{
			name:  "Custom FEN - king is in checkmate",
			opts:  []chess.Option{chess.WithFEN("k7/8/8/8/8/8/3qr3/3K4 w - - 0 1")},
			moves: nil,
		},
		{
			name:  "Custom FEN - king is in stalemate",
			opts:  []chess.Option{chess.WithFEN("k7/8/8/8/8/8/2q5/K7 w - - 0 1")},
			moves: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			c, errOpts := chess.NewChess(test.opts...)
			require.Nil(t, errOpts)

			// Act
			moves, err := c.AvailableLegalMoves()

			// Assert
			if test.errMsg == "" {
				require.Nil(t, err)
				assert.ElementsMatch(t, test.moves, moves)
				return
			}

			require.NotNil(t, err)
			assert.Equal(t, test.errMsg, err.Error())
		})
	}
}

func TestMakeMove(t *testing.T) {
	tests := []struct {
		name   string
		opts   []chess.Option
		move   string
		FEN    string
		errMsg string
	}{
		{
			name: "Default",
			opts: []chess.Option{},
			move: "e2e4",
			FEN:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			name: "Custom FEN 1",
			opts: []chess.Option{chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1")},
			move: "d3d4",
			FEN:  "8/8/8/k7/3P4/K7/8/8 b - - 0 1",
		},
		{
			name: "Custom FEN 2",
			opts: []chess.Option{chess.WithFEN("8/8/8/8/8/4k3/7r/5K2 w - - 0 1")},
			move: "f1e1",
			FEN:  "8/8/8/8/8/4k3/7r/4K3 b - - 1 1",
		},
		{
			name: "Custom FEN 3 - Castle",
			opts: []chess.Option{chess.WithFEN("k7/8/8/8/8/8/8/4K2R w K - 0 1")},
			move: "e1g1",
			FEN:  "k7/8/8/8/8/8/8/5RK1 b - - 1 1",
		},
		{
			name: "Custom FEN 4 - Capture",
			opts: []chess.Option{chess.WithFEN("k7/8/8/8/8/8/3q4/3K4 w - - 0 1")},
			move: "d1d2",
			FEN:  "k7/8/8/8/8/8/3K4/8 b - - 0 1",
		},
		{
			name: "Custom FEN 5 - Promotion queen",
			opts: []chess.Option{chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1")},
			move: "h7h8q",
			FEN:  "k6Q/8/8/8/8/8/8/7K b - - 0 1",
		},
		{
			name: "Custom FEN 6 - Promotion rook",
			opts: []chess.Option{chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1")},
			move: "h7h8r",
			FEN:  "k6R/8/8/8/8/8/8/7K b - - 0 1",
		},
		{
			name: "Custom FEN 7 - Promotion bishop",
			opts: []chess.Option{chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1")},
			move: "h7h8b",
			FEN:  "k6B/8/8/8/8/8/8/7K b - - 0 1",
		},
		{
			name: "Custom FEN 8 - Promotion knight",
			opts: []chess.Option{chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1")},
			move: "h7h8n",
			FEN:  "k6N/8/8/8/8/8/8/7K b - - 0 1",
		},
		{
			name:   "Custom FEN 9 - Promotion king - error",
			opts:   []chess.Option{chess.WithFEN("k7/7P/8/8/8/8/8/7K w - - 0 1")},
			move:   "h7h8k",
			errMsg: "move is not legal: h7h8k",
		},
		{
			name: "Custom FEN 10 - In passant",
			opts: []chess.Option{chess.WithFEN("7k/8/8/3pP3/8/8/8/7K w - d6 0 1")},
			move: "e5d6",
			FEN:  "7k/8/3P4/8/8/8/8/7K b - - 0 1",
		},
		{
			name:   "Custom FEN 11 - Invalid move - invalid square",
			opts:   []chess.Option{chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1")},
			move:   "d3d9",
			errMsg: "move is not legal: d3d9",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			c, errOpts := chess.NewChess(test.opts...)
			require.Nil(t, errOpts)

			// Act
			err := c.MakeMove(test.move)

			// Assert
			if test.errMsg == "" {
				require.Nil(t, err)
				assert.Equal(t, test.FEN, c.FEN())
				return
			}

			require.NotNil(t, err)
			assert.Equal(t, test.errMsg, err.Error())
		})
	}
}

func TestIsCheck(t *testing.T) {
	tests := []struct {
		name    string
		opts    []chess.Option
		isCheck bool
		errMsg  string
	}{
		{
			name:    "Default",
			opts:    []chess.Option{},
			isCheck: false,
		},
		{
			name:    "Custom FEN - king is in check",
			opts:    []chess.Option{chess.WithFEN("k7/8/8/8/8/8/3q4/3K4 w - - 0 1")},
			isCheck: true,
		},
		{
			name:    "Custom FEN - king is not in check",
			opts:    []chess.Option{chess.WithFEN("k7/8/8/8/8/8/8/3K4 w - - 0 1")},
			isCheck: false,
		},
		{
			name:    "Custom FEN - king is not in the board",
			opts:    []chess.Option{chess.WithFEN("8/8/8/8/8/8/3q4/8 w - - 0 1")},
			isCheck: false,
			errMsg:  "failed to get king position: king not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			c, errOpts := chess.NewChess(test.opts...)
			require.Nil(t, errOpts)

			// Act
			isCheck, err := c.IsCheck()

			// Assert
			assert.Equal(t, test.isCheck, isCheck)
			if test.errMsg != "" {
				require.NotNil(t, err)
				assert.Equal(t, test.errMsg, err.Error())
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestSquare(t *testing.T) {
	tests := []struct {
		name   string
		opts   []chess.Option
		square string
		piece  string
		errMsg string
	}{
		{
			name:   "Default",
			opts:   []chess.Option{},
			square: "e2",
			piece:  "P",
		},
		{
			name:   "Custom FEN",
			opts:   []chess.Option{chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1")},
			square: "d3",
			piece:  "P",
		},
		{
			name:   "Custom FEN - empty square",
			opts:   []chess.Option{chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1")},
			square: "a1",
			piece:  "",
		},
		{
			name:   "Custom FEN - invalid square",
			opts:   []chess.Option{chess.WithFEN("8/8/8/k7/8/K2P4/8/8 w - - 0 1")},
			square: "a9",
			errMsg: "failed to convert algebraic notation to coordinate: coordinate out of bounds",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			c, errOpts := chess.NewChess(test.opts...)
			require.Nil(t, errOpts)

			// Act
			piece, err := c.Square(test.square)

			// Assert
			if test.errMsg != "" {
				require.NotNil(t, err)
				assert.Equal(t, test.errMsg, err.Error())
				return
			}

			require.Nil(t, err)
			assert.Equal(t, test.piece, piece)
		})
	}
}

func TestMakeMove_ScholarMate(t *testing.T) {
	// Arrange
	c, err := chess.NewChess()
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
	assert.Equal(t, "r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4",
		c.FEN())
}

func TestMakeMove_CapablancaSteiner(t *testing.T) {
	// Arrange
	c, err := chess.NewChess()
	require.Nil(t, err)

	// Act
	err = c.MakeMove("e2e4")
	require.Nil(t, err)
	err = c.MakeMove("e7e5")
	require.Nil(t, err)
	err = c.MakeMove("g1f3")
	require.Nil(t, err)
	err = c.MakeMove("b8c6")
	require.Nil(t, err)
	err = c.MakeMove("b1c3")
	require.Nil(t, err)
	err = c.MakeMove("g8f6")
	require.Nil(t, err)
	err = c.MakeMove("f1b5")
	require.Nil(t, err)
	err = c.MakeMove("f8b4")
	require.Nil(t, err)
	err = c.MakeMove("e1g1")
	require.Nil(t, err)
	err = c.MakeMove("e8g8")
	require.Nil(t, err)
	err = c.MakeMove("d2d3")
	require.Nil(t, err)
	err = c.MakeMove("d7d6")
	require.Nil(t, err)
	err = c.MakeMove("c1g5")
	require.Nil(t, err)
	err = c.MakeMove("b4c3")
	require.Nil(t, err)
	err = c.MakeMove("b2c3")
	require.Nil(t, err)
	err = c.MakeMove("c6e7")
	require.Nil(t, err)
	err = c.MakeMove("f3h4")
	require.Nil(t, err)
	err = c.MakeMove("c7c6")
	require.Nil(t, err)
	err = c.MakeMove("b5c4")
	require.Nil(t, err)
	err = c.MakeMove("c8e6")
	require.Nil(t, err)
	err = c.MakeMove("g5f6")
	require.Nil(t, err)
	err = c.MakeMove("g7f6")
	require.Nil(t, err)
	err = c.MakeMove("c4e6")
	require.Nil(t, err)
	err = c.MakeMove("f7e6")
	require.Nil(t, err)
	err = c.MakeMove("d1g4")
	require.Nil(t, err)
	err = c.MakeMove("g8f7")
	require.Nil(t, err)
	err = c.MakeMove("f2f4")
	require.Nil(t, err)
	err = c.MakeMove("f8g8")
	require.Nil(t, err)
	err = c.MakeMove("g4h5")
	require.Nil(t, err)
	err = c.MakeMove("f7g7")
	require.Nil(t, err)
	err = c.MakeMove("f4e5")
	require.Nil(t, err)
	err = c.MakeMove("d6e5")
	require.Nil(t, err)
	err = c.MakeMove("f1f6")
	require.Nil(t, err)
	err = c.MakeMove("g7f6")
	require.Nil(t, err)
	err = c.MakeMove("a1f1")
	require.Nil(t, err)
	err = c.MakeMove("e7f5")
	require.Nil(t, err)
	err = c.MakeMove("h4f5")
	require.Nil(t, err)
	err = c.MakeMove("e6f5")
	require.Nil(t, err)
	err = c.MakeMove("f1f5")
	require.Nil(t, err)
	err = c.MakeMove("f6e7")
	require.Nil(t, err)
	err = c.MakeMove("h5f7")
	require.Nil(t, err)
	err = c.MakeMove("e7d6")
	require.Nil(t, err)
	err = c.MakeMove("f5f6")
	require.Nil(t, err)
	err = c.MakeMove("d6c5")
	require.Nil(t, err)
	err = c.MakeMove("f7b7")
	require.Nil(t, err)
	err = c.MakeMove("d8b6")
	require.Nil(t, err)
	err = c.MakeMove("f6c6")
	require.Nil(t, err)
	err = c.MakeMove("b6c6")
	require.Nil(t, err)
	err = c.MakeMove("b7b4")
	require.Nil(t, err)

	legalMoves, err := c.AvailableLegalMoves()

	// Assert
	assert.Equal(t, "r5r1/p6p/2q5/2k1p3/1Q2P3/2PP4/P1P3PP/6K1 b - - 1 25", c.FEN())
	assert.Nil(t, err)
	assert.Nil(t, legalMoves)
}
