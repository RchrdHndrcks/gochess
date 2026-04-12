package chess

// Perft counts leaf nodes at the given depth. Uses bulk counting at depth 1.
func (c *Chess) Perft(depth int) uint64 {
	if depth == 0 {
		return 1
	}
	moves := c.Moves()
	if depth == 1 {
		return uint64(len(moves))
	}
	var nodes uint64
	for _, m := range moves {
		if err := c.MakeMoveCompact(m); err != nil {
			continue
		}
		nodes += c.Perft(depth - 1)
		c.UnmakeMoveCompact()
	}
	return nodes
}

// PerftDivide returns per-move node counts at the given depth.
func (c *Chess) PerftDivide(depth int) map[string]uint64 {
	result := make(map[string]uint64)
	if depth == 0 {
		return result
	}
	for _, m := range c.Moves() {
		if err := c.MakeMoveCompact(m); err != nil {
			continue
		}
		nodes := c.Perft(depth - 1)
		c.UnmakeMoveCompact()
		result[m.UCI()] = nodes
	}
	return result
}
