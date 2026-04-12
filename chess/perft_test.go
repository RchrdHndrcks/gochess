package chess

import (
	"strconv"
	"testing"
)

const (
	startFEN    = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	perftKiwipeteFEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	position3   = "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"
)

type perftCase struct {
	name  string
	fen   string
	depth int
	nodes uint64
}

func TestPerft(t *testing.T) {
	cases := []perftCase{
		{"start", startFEN, 1, 20},
		{"start", startFEN, 2, 400},
		{"start", startFEN, 3, 8902},
		{"start", startFEN, 4, 197281},
		{"start", startFEN, 5, 4865609},
		{"kiwipete", perftKiwipeteFEN, 1, 48},
		{"kiwipete", perftKiwipeteFEN, 2, 2039},
		{"kiwipete", perftKiwipeteFEN, 3, 97862},
		{"kiwipete", perftKiwipeteFEN, 4, 4085603},
		{"position3", position3, 1, 14},
		{"position3", position3, 2, 191},
		{"position3", position3, 3, 2812},
		{"position3", position3, 4, 43238},
		{"position3", position3, 5, 674624},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(name(tt), func(t *testing.T) {
			// Skip expensive depths in -short mode (CI).
			if testing.Short() {
				switch tt.name {
				case "start":
					if tt.depth >= 5 {
						t.Skipf("skipping %s depth %d in short mode", tt.name, tt.depth)
					}
				case "kiwipete":
					if tt.depth >= 3 {
						t.Skipf("skipping %s depth %d in short mode", tt.name, tt.depth)
					}
				case "position3":
					if tt.depth >= 5 {
						t.Skipf("skipping %s depth %d in short mode", tt.name, tt.depth)
					}
				}
			}
			c, err := New(WithFEN(tt.fen))
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			got := c.Perft(tt.depth)
			if got != tt.nodes {
				t.Errorf("Perft(%s, %d) = %d, want %d", tt.name, tt.depth, got, tt.nodes)
			}
		})
	}
}

func TestPerftDivideConsistency(t *testing.T) {
	c, err := New(WithFEN(startFEN))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	const depth = 3
	divide := c.PerftDivide(depth)
	var sum uint64
	for _, n := range divide {
		sum += n
	}
	want := c.Perft(depth)
	if sum != want {
		t.Errorf("sum of PerftDivide(%d) = %d, want Perft = %d", depth, sum, want)
	}
}

func name(tt perftCase) string {
	return tt.name + "_d" + strconv.Itoa(tt.depth)
}
