package chess

import (
	"github.com/RchrdHndrcks/gochess/v2"
)

// zobristKeys holds the random keys used to compute the Zobrist hash of a
// chess position. Indices use a small fixed shape suitable for incremental
// updates:
//
//   - pieces[colorIndex][pieceType][square]
//   - castling[castleRightsBitmask] (16 entries, one per combination)
//   - enPassant[file] (8 entries, one per file 0-7)
//   - sideToMove is XOR'd into the hash when it is Black's turn to move.
//
// The pieceType axis has 7 slots so that gochess.PieceType values (1..6)
// can be used directly as the index.
type zobristKeys struct {
	pieces     [2][7][64]uint64
	castling   [16]uint64
	enPassant  [8]uint64
	sideToMove uint64
}

// zobrist holds the global table of Zobrist keys. It is initialized once
// at package load with a deterministic SplitMix64 PRNG so that the same
// position always produces the same hash across runs.
var zobrist zobristKeys

func init() {
	rng := splitMix64{state: 0x1234567890ABCDEF}
	for c := 0; c < 2; c++ {
		for pt := 0; pt < 7; pt++ {
			for sq := 0; sq < 64; sq++ {
				zobrist.pieces[c][pt][sq] = rng.next()
			}
		}
	}
	for i := 0; i < 16; i++ {
		zobrist.castling[i] = rng.next()
	}
	for i := 0; i < 8; i++ {
		zobrist.enPassant[i] = rng.next()
	}
	zobrist.sideToMove = rng.next()
}

// splitMix64 is a deterministic 64-bit PRNG used to fill the Zobrist tables.
// It is taken from Sebastiano Vigna's reference implementation.
type splitMix64 struct {
	state uint64
}

func (s *splitMix64) next() uint64 {
	s.state += 0x9E3779B97F4A7C15
	z := s.state
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

// computeHashFromScratch builds the full Zobrist hash for the current
// position by XORing the key for every occupied square, the castling
// rights mask, the en passant file (if any) and the side-to-move key
// when it is Black to move. It is invoked once on construction (or
// FEN load) to seed the incremental hash.
func computeHashFromScratch(c *Chess) uint64 {
	h := uint64(0)
	for sq := 0; sq < 64; sq++ {
		pt, color, ok := c.PieceAt(sq)
		if !ok {
			continue
		}
		h ^= zobrist.pieces[colorIndex(color)][pt][sq]
	}
	h ^= zobrist.castling[uint8(c.availableCastles)]
	if c.enPassantFile >= 0 {
		h ^= zobrist.enPassant[c.enPassantFile]
	}
	if c.turn == gochess.Black {
		h ^= zobrist.sideToMove
	}
	return h
}

// computePawnHashFromScratch builds the pawn-only Zobrist hash for the
// current position by XORing the piece key for every pawn on the board.
// Castling rights, en passant and side-to-move are intentionally excluded
// — the pawn hash is meant to identify pawn structures for evaluation
// caching, not full positions.
func computePawnHashFromScratch(c *Chess) uint64 {
	h := uint64(0)
	for sq := 0; sq < 64; sq++ {
		pt, color, ok := c.PieceAt(sq)
		if !ok || pt != gochess.Pawn {
			continue
		}
		h ^= zobrist.pieces[colorIndex(color)][gochess.Pawn][sq]
	}
	return h
}
