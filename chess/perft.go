package chess

import "fmt"

// Perft counts leaf nodes at the given depth. Uses bulk counting at depth 1.
//
// Perft uses the legal-move list directly (c.moves) instead of c.Moves(), so
// the GivesCheck bit is not computed per move — that work would dominate the
// runtime at depth >= 4 and is unnecessary for perft.
func (c *Chess) Perft(depth int) uint64 {
	if depth < 0 {
		panic(fmt.Sprintf("chess: Perft called with negative depth %d", depth))
	}
	if depth == 0 {
		return 1
	}
	if depth == 1 {
		return uint64(len(c.moves))
	}
	// Snapshot the UCI move list because MakeMove recomputes c.moves.
	uciMoves := make([]string, len(c.moves))
	copy(uciMoves, c.moves)
	var nodes uint64
	for _, uci := range uciMoves {
		if err := c.MakeMove(uci); err != nil {
			// MakeMove rejects only moves that aren't in c.moves, but we just
			// took the list from c.moves: any error here is a bug, not a
			// recoverable condition. Fail loudly so it can't undercount.
			panic(fmt.Sprintf("chess: Perft: MakeMove(%q) failed: %v", uci, err))
		}
		nodes += c.Perft(depth - 1)
		c.UnmakeMove()
	}
	return nodes
}

// PerftDivide returns per-move node counts at the given depth.
func (c *Chess) PerftDivide(depth int) map[string]uint64 {
	result := make(map[string]uint64)
	if depth <= 0 {
		return result
	}
	uciMoves := make([]string, len(c.moves))
	copy(uciMoves, c.moves)
	for _, uci := range uciMoves {
		if err := c.MakeMove(uci); err != nil {
			panic(fmt.Sprintf("chess: PerftDivide: MakeMove(%q) failed: %v", uci, err))
		}
		nodes := c.Perft(depth - 1)
		c.UnmakeMove()
		result[uci] = nodes
	}
	return result
}
