package chess_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/chess"
)

func BenchmarkCapablancaSteiner(b *testing.B) {
	b.ResetTimer()

	for b.Loop() {
		// Arrange
		c, err := chess.New()
		if err != nil {
			b.Fatalf("Error creating new chess game: %v", err)
		}

		// Act
		moves := []string{
			"e2e4", "e7e5", "g1f3", "b8c6", "b1c3", "g8f6", "f1b5", "f8b4",
			"e1g1", "e8g8", "d2d3", "d7d6", "c1g5", "b4c3", "b2c3", "c6e7",
			"f3h4", "c7c6", "b5c4", "c8e6", "g5f6", "g7f6", "c4e6", "f7e6",
			"d1g4", "g8f7", "f2f4", "f8g8", "g4h5", "f7g7", "f4e5", "d6e5",
			"f1f6", "g7f6", "a1f1", "e7f5", "h4f5", "e6f5", "f1f5", "f6e7",
			"h5f7", "e7d6", "f5f6", "d6c5", "f7b7", "d8b6", "f6c6", "b6c6",
			"b7b4",
		}

		for _, move := range moves {
			err = c.MakeMove(move)
			if err != nil {
				b.Fatalf("Error making move %s: %v", move, err)
			}
		}

		if c.FEN() != "r5r1/p6p/2q5/2k1p3/1Q2P3/2PP4/P1P3PP/6K1 b - - 1 25" {
			b.Fatalf("Unexpected final position: %s", c.FEN())
		}

		if c.AvailableMoves() != nil {
			b.Fatalf("Expected no legal moves, but got some")
		}
	}
}

func BenchmarkVariousMoves(b *testing.B) {
	moves := []string{
		"e2e4", "e7e5", "g1f3", "b8c6", "b1c3", "g8f6", "f1b5", "f8b4",
		"e1g1", "e8g8", "d2d3", "d7d6", "c1g5", "b4c3", "b2c3", "c6e7",
		"f3h4", "c7c6", "b5c4", "c8e6", "g5f6", "g7f6", "c4e6", "f7e6",
		"d1g4", "g8f7", "f2f4", "f8g8", "g4h5", "f7g7", "f4e5", "d6e5",
		"f1f6", "g7f6", "a1f1", "e7f5", "h4f5", "e6f5", "f1f5", "f6e7",
		"h5f7", "e7d6", "f5f6", "d6c5", "f7b7", "d8b6", "f6c6", "b6c6",
		"b7b4",
	}

	b.ResetTimer()
	for b.Loop() {
		c, err := chess.New()
		if err != nil {
			b.Fatalf("Error creating new chess game: %v", err)
		}

		for _, move := range moves {
			err = c.MakeMove(move)
			if err != nil {
				b.Fatalf("Error making move %s: %v", move, err)
			}
		}
	}
}
