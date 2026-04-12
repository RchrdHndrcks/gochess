package chess

// CastleRights is a 4-bit bitmask representing castling availability.
type CastleRights uint8

const (
	WhiteKingside  CastleRights = 1 << iota // K
	WhiteQueenside                          // Q
	BlackKingside                           // k
	BlackQueenside                          // q

	NoCastling  CastleRights = 0
	AllCastling CastleRights = WhiteKingside | WhiteQueenside | BlackKingside | BlackQueenside
)

// Has reports whether cr contains all the bits in right.
func (cr CastleRights) Has(right CastleRights) bool {
	return cr&right == right
}

// String returns the FEN castling field (e.g. "KQkq", "Kq", "-").
func (cr CastleRights) String() string {
	if cr == NoCastling {
		return "-"
	}
	s := ""
	if cr.Has(WhiteKingside) {
		s += "K"
	}
	if cr.Has(WhiteQueenside) {
		s += "Q"
	}
	if cr.Has(BlackKingside) {
		s += "k"
	}
	if cr.Has(BlackQueenside) {
		s += "q"
	}
	return s
}
