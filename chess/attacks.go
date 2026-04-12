package chess

import (
	"strings"

	"github.com/RchrdHndrcks/gochess/v2"
)

// PieceSet is a bitmask of piece types. Each bit represents one piece type
// and is used to describe which kinds of pieces attack a given square.
type PieceSet uint8

// Piece type bits used by PieceSet.
const (
	// NoPieces is the empty PieceSet (no piece types present).
	NoPieces PieceSet = 0
	// PawnSet is the bit representing pawns.
	PawnSet PieceSet = 1 << iota
	// KnightSet is the bit representing knights.
	KnightSet
	// BishopSet is the bit representing bishops.
	BishopSet
	// RookSet is the bit representing rooks.
	RookSet
	// QueenSet is the bit representing queens.
	QueenSet
	// KingSet is the bit representing kings.
	KingSet
)

// Has reports whether ps contains all bits of p.
func (ps PieceSet) Has(p PieceSet) bool { return ps&p == p && p != 0 }

// String returns a human-readable representation of the PieceSet for debugging.
func (ps PieceSet) String() string {
	if ps == NoPieces {
		return "NoPieces"
	}
	parts := make([]string, 0, 6)
	if ps&PawnSet != 0 {
		parts = append(parts, "Pawn")
	}
	if ps&KnightSet != 0 {
		parts = append(parts, "Knight")
	}
	if ps&BishopSet != 0 {
		parts = append(parts, "Bishop")
	}
	if ps&RookSet != 0 {
		parts = append(parts, "Rook")
	}
	if ps&QueenSet != 0 {
		parts = append(parts, "Queen")
	}
	if ps&KingSet != 0 {
		parts = append(parts, "King")
	}
	return strings.Join(parts, "|")
}

// straightDirs are the four orthogonal ray directions (rook/queen).
var straightDirs = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// diagonalDirs are the four diagonal ray directions (bishop/queen).
var diagonalDirs = [4][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

// knightOffsets are the eight knight jump offsets.
var knightOffsets = [8][2]int{
	{1, 2}, {2, 1}, {-1, 2}, {-2, 1},
	{1, -2}, {2, -1}, {-1, -2}, {-2, -1},
}

// kingOffsets are the eight king step offsets.
var kingOffsets = [8][2]int{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
}

// pawnAttackerOffsets returns the two relative offsets that a pawn of `side`
// must occupy in order to attack a square. White pawns move toward y=0, so a
// white pawn attacking square (x,y) sits at (x±1, y+1). Black pawns are mirrored.
func pawnAttackerOffsets(side gochess.Piece) [2][2]int {
	if side == gochess.White {
		return [2][2]int{{-1, 1}, {1, 1}}
	}
	return [2][2]int{{-1, -1}, {1, -1}}
}

// inBounds reports whether (x,y) lies on the 8x8 board.
func inBounds(x, y int) bool {
	return x >= 0 && x < 8 && y >= 0 && y < 8
}

// IsAttacked reports whether the given square is attacked by any piece of
// the given side. It performs targeted ray, knight and pawn lookups from
// sq outward and is recomputed on every call (no cache).
func (c *Chess) IsAttacked(sq gochess.Coordinate, side gochess.Piece) bool {
	// Pawn attacks.
	for _, off := range pawnAttackerOffsets(side) {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.Pawn|side {
			return true
		}
	}

	// Knight attacks.
	for _, off := range knightOffsets {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.Knight|side {
			return true
		}
	}

	// King attacks.
	for _, off := range kingOffsets {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.King|side {
			return true
		}
	}

	// Straight rays (rook/queen).
	for _, dir := range straightDirs {
		for i := 1; i < 8; i++ {
			x, y := sq.X+dir[0]*i, sq.Y+dir[1]*i
			if !inBounds(x, y) {
				break
			}
			p, _ := c.board.Square(gochess.Coor(x, y))
			if p == gochess.Empty {
				continue
			}
			if gochess.PieceColor(p) == side {
				pt := gochess.PieceType(p)
				if pt == gochess.Rook || pt == gochess.Queen {
					return true
				}
			}
			break
		}
	}

	// Diagonal rays (bishop/queen).
	for _, dir := range diagonalDirs {
		for i := 1; i < 8; i++ {
			x, y := sq.X+dir[0]*i, sq.Y+dir[1]*i
			if !inBounds(x, y) {
				break
			}
			p, _ := c.board.Square(gochess.Coor(x, y))
			if p == gochess.Empty {
				continue
			}
			if gochess.PieceColor(p) == side {
				pt := gochess.PieceType(p)
				if pt == gochess.Bishop || pt == gochess.Queen {
					return true
				}
			}
			break
		}
	}

	return false
}

// AttackedBy returns the set of piece types of the given side that attack sq.
// It performs targeted lookups outward from sq and is recomputed on every
// call (no cache).
func (c *Chess) AttackedBy(sq gochess.Coordinate, side gochess.Piece) PieceSet {
	var ps PieceSet

	for _, off := range pawnAttackerOffsets(side) {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.Pawn|side {
			ps |= PawnSet
			break
		}
	}

	for _, off := range knightOffsets {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.Knight|side {
			ps |= KnightSet
			break
		}
	}

	for _, off := range kingOffsets {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.King|side {
			ps |= KingSet
			break
		}
	}

	for _, dir := range straightDirs {
		for i := 1; i < 8; i++ {
			x, y := sq.X+dir[0]*i, sq.Y+dir[1]*i
			if !inBounds(x, y) {
				break
			}
			p, _ := c.board.Square(gochess.Coor(x, y))
			if p == gochess.Empty {
				continue
			}
			if gochess.PieceColor(p) == side {
				pt := gochess.PieceType(p)
				if pt == gochess.Rook {
					ps |= RookSet
				} else if pt == gochess.Queen {
					ps |= QueenSet
				}
			}
			break
		}
	}

	for _, dir := range diagonalDirs {
		for i := 1; i < 8; i++ {
			x, y := sq.X+dir[0]*i, sq.Y+dir[1]*i
			if !inBounds(x, y) {
				break
			}
			p, _ := c.board.Square(gochess.Coor(x, y))
			if p == gochess.Empty {
				continue
			}
			if gochess.PieceColor(p) == side {
				pt := gochess.PieceType(p)
				if pt == gochess.Bishop {
					ps |= BishopSet
				} else if pt == gochess.Queen {
					ps |= QueenSet
				}
			}
			break
		}
	}

	return ps
}

// PawnAttackMap returns a 64-bit bitmap of every square attacked by at least
// one pawn of the given side. Bit (y*8 + x) is set when square (x,y) is
// attacked. Returned as uint64 to avoid heap allocations in hot paths.
func (c *Chess) PawnAttackMap(side gochess.Piece) uint64 {
	var bb uint64
	pawns := c.pieceLists[colorIndex(side)][gochess.Pawn]
	dy := -1
	if side == gochess.Black {
		dy = 1
	}
	for i := 0; i < pawns.count; i++ {
		sq := pawns.squares[i]
		ty := sq.Y + dy
		if ty < 0 || ty > 7 {
			continue
		}
		if sq.X-1 >= 0 {
			bb |= uint64(1) << uint(ty*8+(sq.X-1))
		}
		if sq.X+1 <= 7 {
			bb |= uint64(1) << uint(ty*8+(sq.X+1))
		}
	}
	return bb
}

// Attackers returns the squares of every piece of the given side that attacks
// sq, including x-ray attackers along sliding rays (a same-side rook/queen
// behind another rook/queen on a straight ray, or a same-side bishop/queen
// behind another bishop/queen on a diagonal ray, is included).
//
// Results are appended to dst in cheapest-first order: pawns, knights,
// bishops, rooks, queens, king. Within the same piece type the order is
// undefined.
func (c *Chess) Attackers(sq gochess.Coordinate, side gochess.Piece, dst []gochess.Coordinate) []gochess.Coordinate {
	var pawns, knights, bishops, rooks, queens, kings []gochess.Coordinate

	// Pawns attacking sq.
	for _, off := range pawnAttackerOffsets(side) {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.Pawn|side {
			pawns = append(pawns, gochess.Coor(x, y))
		}
	}

	// Knights.
	for _, off := range knightOffsets {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.Knight|side {
			knights = append(knights, gochess.Coor(x, y))
		}
	}

	// King.
	for _, off := range kingOffsets {
		x, y := sq.X+off[0], sq.Y+off[1]
		if !inBounds(x, y) {
			continue
		}
		p, _ := c.board.Square(gochess.Coor(x, y))
		if p == gochess.King|side {
			kings = append(kings, gochess.Coor(x, y))
		}
	}

	// Straight rays: x-ray through any same-side rook or queen.
	for _, dir := range straightDirs {
		for i := 1; i < 8; i++ {
			x, y := sq.X+dir[0]*i, sq.Y+dir[1]*i
			if !inBounds(x, y) {
				break
			}
			p, _ := c.board.Square(gochess.Coor(x, y))
			if p == gochess.Empty {
				continue
			}
			if gochess.PieceColor(p) != side {
				break
			}
			pt := gochess.PieceType(p)
			if pt == gochess.Rook {
				rooks = append(rooks, gochess.Coor(x, y))
				continue
			}
			if pt == gochess.Queen {
				queens = append(queens, gochess.Coor(x, y))
				continue
			}
			break
		}
	}

	// Diagonal rays: x-ray through any same-side bishop or queen.
	for _, dir := range diagonalDirs {
		for i := 1; i < 8; i++ {
			x, y := sq.X+dir[0]*i, sq.Y+dir[1]*i
			if !inBounds(x, y) {
				break
			}
			p, _ := c.board.Square(gochess.Coor(x, y))
			if p == gochess.Empty {
				continue
			}
			if gochess.PieceColor(p) != side {
				break
			}
			pt := gochess.PieceType(p)
			if pt == gochess.Bishop {
				bishops = append(bishops, gochess.Coor(x, y))
				continue
			}
			if pt == gochess.Queen {
				queens = append(queens, gochess.Coor(x, y))
				continue
			}
			break
		}
	}

	dst = append(dst, pawns...)
	dst = append(dst, knights...)
	dst = append(dst, bishops...)
	dst = append(dst, rooks...)
	dst = append(dst, queens...)
	dst = append(dst, kings...)
	return dst
}

