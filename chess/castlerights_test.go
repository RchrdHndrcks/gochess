package chess_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/v2/chess"
	"github.com/stretchr/testify/assert"
)

func TestCastleRightsHas(t *testing.T) {
	tests := []struct {
		name  string
		cr    chess.CastleRights
		right chess.CastleRights
		want  bool
	}{
		{"all has WK", chess.AllCastling, chess.WhiteKingside, true},
		{"all has WQ", chess.AllCastling, chess.WhiteQueenside, true},
		{"all has BK", chess.AllCastling, chess.BlackKingside, true},
		{"all has BQ", chess.AllCastling, chess.BlackQueenside, true},
		{"none has WK", chess.NoCastling, chess.WhiteKingside, false},
		{"WK has WK", chess.WhiteKingside, chess.WhiteKingside, true},
		{"WK has WQ", chess.WhiteKingside, chess.WhiteQueenside, false},
		{"WK|BQ has both", chess.WhiteKingside | chess.BlackQueenside, chess.WhiteKingside | chess.BlackQueenside, true},
		{"WK has WK|WQ", chess.WhiteKingside, chess.WhiteKingside | chess.WhiteQueenside, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cr.Has(tt.right))
		})
	}
}

func TestCastleRightsString(t *testing.T) {
	tests := []struct {
		name string
		cr   chess.CastleRights
		want string
	}{
		{"none", chess.NoCastling, "-"},
		{"all", chess.AllCastling, "KQkq"},
		{"WK only", chess.WhiteKingside, "K"},
		{"WQ only", chess.WhiteQueenside, "Q"},
		{"BK only", chess.BlackKingside, "k"},
		{"BQ only", chess.BlackQueenside, "q"},
		{"K and q", chess.WhiteKingside | chess.BlackQueenside, "Kq"},
		{"KQ", chess.WhiteKingside | chess.WhiteQueenside, "KQ"},
		{"kq", chess.BlackKingside | chess.BlackQueenside, "kq"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cr.String())
		})
	}
}

func TestCastleRightsRemoval(t *testing.T) {
	tests := []struct {
		name   string
		start  chess.CastleRights
		remove chess.CastleRights
		want   chess.CastleRights
	}{
		{"remove WK from all", chess.AllCastling, chess.WhiteKingside, chess.WhiteQueenside | chess.BlackKingside | chess.BlackQueenside},
		{"remove all white", chess.AllCastling, chess.WhiteKingside | chess.WhiteQueenside, chess.BlackKingside | chess.BlackQueenside},
		{"remove non-present", chess.WhiteKingside, chess.BlackQueenside, chess.WhiteKingside},
		{"remove all", chess.AllCastling, chess.AllCastling, chess.NoCastling},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.start &^ tt.remove
			assert.Equal(t, tt.want, got)
		})
	}
}
