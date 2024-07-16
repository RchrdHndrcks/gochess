package pkg_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailableMoves(t *testing.T) {
	tests := []struct {
		name                string
		turn                int8
		inPassantSquare     string
		castlePossibilities string
		FEN                 string
		availableMoves      []string
	}{
		{
			name:            "Default",
			turn:            pkg.White,
			FEN:             "8/8/8/8/8/8/8/8 w - - 0 0",
			inPassantSquare: "",
			availableMoves:  []string{},
		},
		{
			name:            "Custom FEN",
			turn:            pkg.White,
			FEN:             "8/8/8/k7/8/K2P4/8/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"a3a4", "a3b4", "a3b3", "a3b2",
				"a3a2", "d3d4"},
		},
		{
			name:            "Initial position - available moves 1",
			turn:            pkg.White,
			FEN:             "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"b1a3", "b1c3", "g1f3", "g1h3",
				"b2b3", "b2b4", "c2c3", "c2c4", "d2d3", "d2d4", "e2e3",
				"e2e4", "f2f3", "f2f4", "h2h3", "h2h4", "a2a3", "a2a4",
				"g2g3", "g2g4"},
		},
		{
			name:            "Initial position - available moves 2",
			turn:            pkg.Black,
			FEN:             "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			inPassantSquare: "e3",
			availableMoves: []string{"b8a6", "b8c6", "g8f6", "g8h6",
				"a7a6", "a7a5", "b7b6", "b7b5", "c7c6", "c7c5", "d7d6",
				"d7d5", "e7e6", "e7e5", "f7f6", "f7f5", "g7g6", "g7g5",
				"h7h6", "h7h5"},
		},
		{
			name:            "Initial position - available moves 3",
			turn:            pkg.White,
			FEN:             "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			inPassantSquare: "d6",
			availableMoves: []string{"b1a3", "b1c3", "g1f3", "g1h3", "g1e2",
				"a2a3", "a2a4", "b2b3", "b2b4", "c2c3", "c2c4", "d2d3",
				"d2d4", "f2f3", "f2f4", "g2g3", "g2g4", "h2h3", "h2h4",
				"e4d5", "e4e5", "e1e2", "f1e2", "f1d3", "f1c4", "f1b5",
				"f1a6", "d1e2", "d1f3", "d1g4", "d1h5"},
		},
		{
			name:            "Initial position - available moves 4",
			turn:            pkg.Black,
			FEN:             "rnbqkbnr/ppp1pppp/8/3p4/4P3/2N5/PPPP1PPP/R1BQKBNR b KQkq - 0 2",
			inPassantSquare: "",
			availableMoves: []string{"b8a6", "b8c6", "b8d7", "g8f6", "g8h6",
				"a7a6", "a7a5", "b7b6", "b7b5", "c7c6", "c7c5", "d5d4",
				"d5e4", "e7e6", "e7e5", "f7f6", "f7f5", "g7g6", "g7g5",
				"h7h6", "h7h5", "d8d7", "d8d6", "c8d7", "c8e6", "c8f5",
				"c8g4", "c8h3", "e8d7"},
		},
		{
			name:            "Rook - available moves",
			turn:            pkg.White,
			FEN:             "8/8/8/8/4R3/8/8/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"e4e1", "e4e2", "e4e3", "e4e5",
				"e4e6", "e4e7", "e4e8", "e4d4", "e4c4", "e4b4", "e4a4",
				"e4f4", "e4g4", "e4h4"},
		},
		{
			name:            "Rook with pieces - available moves",
			turn:            pkg.White,
			FEN:             "4n3/8/8/5p2/1n2RP2/4r3/8/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"e4e3", "e4e5", "e4e6", "e4e7",
				"e4e8", "e4d4", "e4c4", "e4b4"},
		},
		{
			name:            "Knight - available moves",
			turn:            pkg.White,
			FEN:             "8/8/8/8/4N3/8/8/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"e4c3", "e4d2", "e4f2", "e4g3",
				"e4g5", "e4f6", "e4d6", "e4c5"},
		},
		{
			name:            "Knight with pieces - available moves",
			turn:            pkg.White,
			FEN:             "8/8/8/8/2p1N3/2P5/3n4/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"e4d2", "e4f2", "e4g3", "e4g5",
				"e4f6", "e4d6", "e4c5"},
		},
		{
			name:            "Bishop - available moves",
			turn:            pkg.White,
			FEN:             "8/8/8/8/4B3/8/8/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"e4d3", "e4c2", "e4b1", "e4f3",
				"e4g2", "e4h1", "e4d5", "e4c6", "e4b7", "e4a8", "e4f5",
				"e4g6", "e4h7"},
		},
		{
			name:            "Queen - available moves",
			turn:            pkg.White,
			FEN:             "8/8/8/8/4Q3/8/8/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"e4d3", "e4c2", "e4b1", "e4f3",
				"e4g2", "e4h1", "e4d5", "e4c6", "e4b7", "e4a8", "e4f5",
				"e4g6", "e4h7", "e4e3", "e4e2", "e4e1", "e4e5", "e4e6",
				"e4e7", "e4e8", "e4d4", "e4c4", "e4b4", "e4a4", "e4f4",
				"e4g4", "e4h4"},
		},
		{
			name:            "King - available moves",
			turn:            pkg.White,
			FEN:             "8/8/8/8/4K3/8/8/8 w - - 0 1",
			inPassantSquare: "",
			availableMoves: []string{"e4d3", "e4d4", "e4d5", "e4e3",
				"e4e5", "e4f3", "e4f4", "e4f5"},
		},
		{
			name:                "King castle - available moves 1",
			turn:                pkg.White,
			FEN:                 "8/8/8/8/8/8/r6r/R3K2R w KQ - 0 1",
			castlePossibilities: "KQ",
			availableMoves: []string{"e1d1", "e1f1", "e1d2", "e1e2", "e1f2",
				"e1c1", "e1g1", "a1b1", "a1c1", "a1d1", "a1a2", "h1h2", "h1g1",
				"h1f1"},
		},
		{
			name:                "King castle - available moves 2",
			turn:                pkg.White,
			FEN:                 "8/8/8/8/8/8/7r/4K2R w K - 0 1",
			castlePossibilities: "K",
			availableMoves: []string{"e1d1", "e1f1", "e1d2", "e1e2", "e1f2",
				"e1g1", "h1h2", "h1g1", "h1f1"},
		},
		{
			name:                "King castle - available moves 3",
			turn:                pkg.White,
			FEN:                 "8/8/8/8/8/8/7r/4K2R w - - 0 1",
			castlePossibilities: "",
			availableMoves: []string{"e1d1", "e1f1", "e1d2", "e1e2", "e1f2",
				"h1h2", "h1g1", "h1f1"},
		},
		{
			name:                "King castle - available moves 4",
			turn:                pkg.White,
			FEN:                 "8/8/8/8/8/8/7r/4K2R w - - 0 1",
			castlePossibilities: "x",
			availableMoves: []string{"e1d1", "e1f1", "e1d2", "e1e2", "e1f2",
				"h1h2", "h1g1", "h1f1"},
		},
		{
			name:           "Pawn - available moves",
			turn:           pkg.White,
			FEN:            "8/8/8/8/4P3/8/P7/8 w - - 0 1",
			availableMoves: []string{"e4e5", "a2a3", "a2a4"},
		},
		{
			name:           "Pawn - available moves 2",
			turn:           pkg.White,
			FEN:            "8/P7/8/8/8/8/8/8 w - - 0 1",
			availableMoves: []string{"a7a8q", "a7a8r", "a7a8b", "a7a8n"},
		},
		{
			name:           "Pawn - available moves 3",
			turn:           pkg.Black,
			FEN:            "8/8/8/8/8/8/p7/8 b - - 0 1",
			availableMoves: []string{"a2a1q", "a2a1r", "a2a1b", "a2a1n"},
		},
		{
			name:            "Pawn - available moves - in passant",
			turn:            pkg.White,
			FEN:             "8/8/8/3pP3/8/8/8/8 w - d6 0 1",
			inPassantSquare: "d6",
			availableMoves:  []string{"e5d6", "e5e6"},
		},
		{
			name:           "Pawn - promotion 1",
			turn:           pkg.White,
			FEN:            "8/3P4/8/8/8/8/8/8 w - - 0 1",
			availableMoves: []string{"d7d8q", "d7d8r", "d7d8b", "d7d8n"},
		},
		{
			name: "Pawn - promotion 2",
			turn: pkg.White,
			FEN:  "2r5/3P4/8/8/8/8/8/8 w - - 0 1",
			availableMoves: []string{"d7d8q", "d7d8r", "d7d8b", "d7d8n", "d7c8q", "d7c8r",
				"d7c8b", "d7c8n"},
		},
		{
			name: "Pawn - promotion 3",
			turn: pkg.White,
			FEN:  "2rbr3/3P4/8/8/8/8/8/8 w - - 0 1",
			availableMoves: []string{"d7c8q", "d7c8r", "d7c8b", "d7c8n", "d7e8q", "d7e8r",
				"d7e8b", "d7e8n"},
		},
		{
			name:           "Pawn - promotion 4",
			turn:           pkg.Black,
			FEN:            "8/8/8/8/8/8/p7/8 b - - 0 1",
			availableMoves: []string{"a2a1q", "a2a1r", "a2a1b", "a2a1n"},
		},
		{
			name: "Pawn - promotion 5",
			turn: pkg.Black,
			FEN:  "8/8/8/8/8/8/p7/1R6 b - - 0 1",
			availableMoves: []string{"a2a1q", "a2a1r", "a2a1b", "a2a1n", "a2b1q", "a2b1r",
				"a2b1b", "a2b1n"},
		},
		{
			name:           "Pawn - promotion 6",
			turn:           pkg.Black,
			FEN:            "8/8/8/8/8/8/p7/RR6 b - - 0 1",
			availableMoves: []string{"a2b1q", "a2b1r", "a2b1b", "a2b1n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			b := pkg.NewBoard()
			errLoad := b.LoadPosition(tt.FEN)

			// Act
			moves, err := b.AvailableMoves(tt.turn, tt.inPassantSquare, tt.castlePossibilities)

			// Assert
			require.Nil(t, errLoad)
			require.Nil(t, err)
			assert.ElementsMatch(t, tt.availableMoves, moves)
		})
	}
}

func TestLoadPosition_Errors(t *testing.T) {
	tests := []struct {
		name   string
		turn   int8
		FEN    string
		errMsg string
	}{
		{
			name:   "Invalid FEN",
			FEN:    "invalid",
			errMsg: "invalid FEN: invalid",
		},
		{
			name:   "Invalid FEN - invalid number of properties 1",
			FEN:    "8/8/8/8/8/8/8/8",
			errMsg: "invalid FEN: 8/8/8/8/8/8/8/8",
		},
		{
			name:   "Invalid FEN - invalid number of properties 2",
			FEN:    "8/8/8/8/8/8/8/8 w",
			errMsg: "invalid FEN: 8/8/8/8/8/8/8/8 w",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			b := pkg.NewBoard()

			// Act
			err := b.LoadPosition(tt.FEN)

			// Assert
			require.NotNil(t, err)
			assert.Equal(t, tt.errMsg, err.Error())
		})
	}
}

func TestMakeMove(t *testing.T) {
	tests := []struct {
		name   string
		FEN    string
		move   string
		want   string
		errMsg string
	}{
		{
			name: "Valid move 1",
			FEN:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			move: "e2e4",
			want: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			name: "Valid move 2 - Coronation",
			FEN:  "8/7P/8/8/8/8/8/8 w - - 0 1",
			move: "h7h8q",
			want: "7Q/8/8/8/8/8/8/8 b - - 0 1",
		},
		{
			name: "Valid move 3 - Castle white kingside",
			FEN:  "8/8/8/8/8/8/8/4K2R w K - 0 1",
			move: "e1g1",
			want: "8/8/8/8/8/8/8/5RK1 b - - 0 1",
		},
		{
			name: "Valid move 4 - Castle white queenside",
			FEN:  "8/8/8/8/8/8/8/R3K3 w Q - 0 1",
			move: "e1c1",
			want: "8/8/8/8/8/8/8/2KR4 b - - 0 1",
		},
		{
			name: "Valid move 5 - Castle black kingside",
			FEN:  "4k2r/8/8/8/8/8/8/8 b k - 0 1",
			move: "e8g8",
			want: "5rk1/8/8/8/8/8/8/8 w - - 0 1",
		},
		{
			name: "Valid move 6 - Castle black queenside",
			FEN:  "r3k3/8/8/8/8/8/8/8 b q - 0 1",
			move: "e8c8",
			want: "2kr4/8/8/8/8/8/8/8 w - - 0 1",
		},
		{
			name: "Valid move 7 - Queen moves e1g1",
			FEN:  "8/8/8/8/8/8/8/4Q2R w - - 0 1",
			move: "e1g1",
			want: "8/8/8/8/8/8/8/6QR b - - 0 1",
		},
		{
			name: "Valid move 8 - Capture",
			FEN:  "8/8/8/3p4/4P3/8/8/8 w - - 0 1",
			move: "e4d5",
			want: "8/8/8/3P4/8/8/8/8 b - - 0 1",
		},
		{
			name:   "Invalid move - invalid origin square",
			FEN:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			move:   "i2e4",
			errMsg: "error in parsing origin square to algebraic: coordinate out of bounds",
		},
		{
			name:   "Invalid move - invalid destination square",
			FEN:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			move:   "e4i4",
			errMsg: "error in parsing target square to algebraic: coordinate out of bounds",
		},
		{
			name:   "Invalid move - invalid move length",
			FEN:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			move:   "e4",
			errMsg: "move length should be 4 or 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			b := pkg.NewBoard()
			errLoad := b.LoadPosition(tt.FEN)
			wantB := pkg.NewBoard()
			wantErrLoad := wantB.LoadPosition(tt.want)

			// Act
			err := b.MakeMove(tt.move)

			// Assert
			require.Nil(t, errLoad)

			if tt.errMsg != "" {
				require.NotNil(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}

			require.Nil(t, err)

			if tt.want != "" {
				require.Nil(t, wantErrLoad)
				assert.Equal(t, wantB, b)
			}
		})
	}
}
