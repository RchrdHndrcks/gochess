package chess

// MaxMoves is the fixed capacity of MoveList, chosen as a safe upper bound
// with headroom for move generation. The known maximum number of legal
// moves in any reachable chess position is below this value; 256 is the
// container's hard limit, not a tight theoretical bound.
const MaxMoves = 256

// MoveList is a fixed-array move container that avoids heap allocations
// during move generation in performance-sensitive search paths.
type MoveList struct {
	Moves [MaxMoves]Move
	Count int
}

// Add appends m to the list. It panics if the list is already at MaxMoves
// capacity (or if Count has been externally corrupted) so an overflow does
// not silently corrupt adjacent memory or skip moves.
func (ml *MoveList) Add(m Move) {
	if ml.Count < 0 || ml.Count >= MaxMoves {
		panic("chess: MoveList.Add overflow or invalid Count")
	}
	ml.Moves[ml.Count] = m
	ml.Count++
}

// At returns the move at index i.
func (ml *MoveList) At(i int) Move { return ml.Moves[i] }

// Len returns the number of moves in the list.
func (ml *MoveList) Len() int { return ml.Count }

// Reset clears the list so it can be reused.
func (ml *MoveList) Reset() { ml.Count = 0 }

// Swap exchanges the moves at positions i and j.
func (ml *MoveList) Swap(i, j int) { ml.Moves[i], ml.Moves[j] = ml.Moves[j], ml.Moves[i] }
