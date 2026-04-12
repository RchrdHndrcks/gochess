package chess

// MaxMoves is the maximum number of legal moves in any chess position.
const MaxMoves = 256

// MoveList is a fixed-array move container that avoids heap allocations
// during move generation in performance-sensitive search paths.
type MoveList struct {
	Moves [MaxMoves]Move
	Count int
}

// Add appends m to the list.
func (ml *MoveList) Add(m Move) { ml.Moves[ml.Count] = m; ml.Count++ }

// At returns the move at index i.
func (ml *MoveList) At(i int) Move { return ml.Moves[i] }

// Len returns the number of moves in the list.
func (ml *MoveList) Len() int { return ml.Count }

// Reset clears the list so it can be reused.
func (ml *MoveList) Reset() { ml.Count = 0 }

// Swap exchanges the moves at positions i and j.
func (ml *MoveList) Swap(i, j int) { ml.Moves[i], ml.Moves[j] = ml.Moves[j], ml.Moves[i] }
