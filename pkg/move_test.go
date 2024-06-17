package pkg_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess/pkg"
)

func TestUCI(t *testing.T) {
	tests := []struct {
		name string
		oCor pkg.Coordinate
		tCor pkg.Coordinate
		want string
	}{
		{
			name: "a1 to a2",
			oCor: pkg.Coor(0, 7),
			tCor: pkg.Coor(0, 6),
			want: "a1a2",
		},
		{
			name: "a1 to b2",
			oCor: pkg.Coor(0, 7),
			tCor: pkg.Coor(1, 6),
			want: "a1b2",
		},
		{
			name: "a1 to b1",
			oCor: pkg.Coor(0, 7),
			tCor: pkg.Coor(1, 7),
			want: "a1b1",
		},
		{
			name: "a1 to b3",
			oCor: pkg.Coor(0, 7),
			tCor: pkg.Coor(1, 5),
			want: "a1b3",
		},
		{
			name: "a1 to h8",
			oCor: pkg.Coor(0, 7),
			tCor: pkg.Coor(7, 0),
			want: "a1h8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pkg.UCI(tt.oCor, tt.tCor); got != tt.want {
				t.Errorf("UCI() = %v, want %v", got, tt.want)
			}
		})
	}
}
