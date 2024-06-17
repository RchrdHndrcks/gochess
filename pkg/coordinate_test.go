package pkg_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordinateToAlgebraic(t *testing.T) {
	tests := []struct {
		coor pkg.Coordinate
		want string
	}{
		{
			coor: pkg.Coor(0, 0),
			want: "a8",
		},
		{
			coor: pkg.Coor(7, 7),
			want: "h1",
		},
		{
			coor: pkg.Coor(0, 7),
			want: "a1",
		},
		{
			coor: pkg.Coor(7, 0),
			want: "h8",
		},
		{
			coor: pkg.Coor(4, 4),
			want: "e4",
		},
		{
			coor: pkg.Coor(4, 5),
			want: "e3",
		},
		{
			coor: pkg.Coor(3, 3),
			want: "d5",
		},
		{
			coor: pkg.Coor(8, 8),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := pkg.CoordinateToAlgebraic(tt.coor)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestAlgebraicToCoordinate(t *testing.T) {
	tests := []struct {
		algebraic string
		want      pkg.Coordinate
		errMsg    string
	}{
		{
			algebraic: "a8",
			want:      pkg.Coor(0, 0),
		},
		{
			algebraic: "h1",
			want:      pkg.Coor(7, 7),
		},
		{
			algebraic: "a1",
			want:      pkg.Coor(0, 7),
		},
		{
			algebraic: "h8",
			want:      pkg.Coor(7, 0),
		},
		{
			algebraic: "e4",
			want:      pkg.Coor(4, 4),
		},
		{
			algebraic: "e3",
			want:      pkg.Coor(4, 5),
		},
		{
			algebraic: "d5",
			want:      pkg.Coor(3, 3),
		},
		{
			algebraic: "i9",
			want:      pkg.Coor(0, 0),
			errMsg:    "coordinate out of bounds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.algebraic, func(t *testing.T) {
			// Act
			got, err := pkg.AlgebraicToCoordinate(tt.algebraic)

			// Assert
			assert.Equal(t, tt.want, got, "expected %v for %s, got %v", tt.want, tt.algebraic, got)
			if tt.errMsg != "" {
				require.NotNil(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}
