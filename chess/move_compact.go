package chess

import (
	"github.com/RchrdHndrcks/gochess/v2"
)

// Move is a packed representation of a chess move.
//
// Layout (32 bits):
//
//	bits  0-5:  from square (0-63)
//	bits  6-11: to square (0-63)
//	bits 12-14: promotion piece index (None=0, Knight=1, Bishop=2, Rook=3, Queen=4)
//	bits 15-17: flags (Quiet, Capture, EnPassant, Castle, DoublePush, Promotion, PromotionCapture)
//	bits 18-20: captured piece index (Pawn=1, Knight=2, Bishop=3, Rook=4, Queen=5, King=6)
//	bit  21:    gives check (1 = move delivers check to opponent)
//	bits 22-31: reserved for future use (e.g., move ordering score)
type Move uint32

// MoveFlag categorises the kind of move encoded in a Move.
type MoveFlag uint8

// Move flag values stored in bits 15-17 of Move.
const (
	FlagQuiet            MoveFlag = 0
	FlagCapture          MoveFlag = 1
	FlagEnPassant        MoveFlag = 2
	FlagCastle           MoveFlag = 3
	FlagDoublePush       MoveFlag = 4
	FlagPromotion        MoveFlag = 5
	FlagPromotionCapture MoveFlag = 6
)

// NullMove is the zero Move, used as a sentinel for "no move".
const NullMove Move = 0

// Bit layout helpers.
const (
	moveFromShift      = 0
	moveToShift        = 6
	movePromoShift     = 12
	moveFlagsShift     = 15
	moveCapturedShift  = 18
	moveGivesCheckBit  = 21
	moveSquareMask     = 0x3F
	movePromoMask      = 0x7
	moveFlagsMask      = 0x7
	moveCapturedMask   = 0x7
	moveGivesCheckMask = uint32(1) << moveGivesCheckBit
)

// promotionIndex maps a promotion piece type to its 3-bit index in Move.
var promotionIndex = map[gochess.Piece]uint32{
	gochess.Empty:  0,
	gochess.Knight: 1,
	gochess.Bishop: 2,
	gochess.Rook:   3,
	gochess.Queen:  4,
}

// promotionPiece is the inverse lookup for promotionIndex.
var promotionPiece = [...]gochess.Piece{
	0: gochess.Empty,
	1: gochess.Knight,
	2: gochess.Bishop,
	3: gochess.Rook,
	4: gochess.Queen,
}

// capturedIndex maps a piece type (without color) to its 3-bit captured index.
var capturedIndex = map[gochess.Piece]uint32{
	gochess.Empty:  0,
	gochess.Pawn:   1,
	gochess.Knight: 2,
	gochess.Bishop: 3,
	gochess.Rook:   4,
	gochess.Queen:  5,
	gochess.King:   6,
}

// capturedPieceFromIndex is the inverse lookup for capturedIndex.
var capturedPieceFromIndex = [...]gochess.Piece{
	0: gochess.Empty,
	1: gochess.Pawn,
	2: gochess.Knight,
	3: gochess.Bishop,
	4: gochess.Rook,
	5: gochess.Queen,
	6: gochess.King,
}

// squareFromCoordinate converts a (X,Y) board coordinate (Y=0 is rank 8) into a
// 0-63 square index where a1=0 and h8=63.
func squareFromCoordinate(c gochess.Coordinate) uint32 {
	return uint32((7-c.Y)*8 + c.X)
}

// coordinateFromSquare is the inverse of squareFromCoordinate.
func coordinateFromSquare(sq uint32) gochess.Coordinate {
	x := int(sq) % 8
	rank := int(sq) / 8 // 0..7 where 0 is rank 1
	y := 7 - rank
	return gochess.Coor(x, y)
}

// NewMove builds a Move with the given from/to squares and flag, no
// promotion, no captured piece, no gives-check bit set. Use the With*
// builder methods to layer additional information.
func NewMove(from, to gochess.Coordinate, flags MoveFlag) Move {
	return Move(squareFromCoordinate(from)<<moveFromShift |
		squareFromCoordinate(to)<<moveToShift |
		uint32(flags)<<moveFlagsShift)
}

// From returns the origin coordinate of the move.
func (m Move) From() gochess.Coordinate {
	return coordinateFromSquare((uint32(m) >> moveFromShift) & moveSquareMask)
}

// To returns the destination coordinate of the move.
func (m Move) To() gochess.Coordinate {
	return coordinateFromSquare((uint32(m) >> moveToShift) & moveSquareMask)
}

// Flags returns the move flag.
func (m Move) Flags() MoveFlag {
	return MoveFlag((uint32(m) >> moveFlagsShift) & moveFlagsMask)
}

// Promotion returns the promotion piece type (without color), or
// gochess.Empty if the move is not a promotion.
func (m Move) Promotion() gochess.Piece {
	idx := (uint32(m) >> movePromoShift) & movePromoMask
	if int(idx) >= len(promotionPiece) {
		return gochess.Empty
	}
	return promotionPiece[idx]
}

// CapturedPiece returns the captured piece type (without color), or
// gochess.Empty if the move did not capture.
func (m Move) CapturedPiece() gochess.Piece {
	idx := (uint32(m) >> moveCapturedShift) & moveCapturedMask
	if int(idx) >= len(capturedPieceFromIndex) {
		return gochess.Empty
	}
	return capturedPieceFromIndex[idx]
}

// IsCapture reports whether the move captures a piece (regular capture,
// en passant, or promotion-capture).
func (m Move) IsCapture() bool {
	f := m.Flags()
	return f == FlagCapture || f == FlagEnPassant || f == FlagPromotionCapture
}

// GivesCheck reports whether bit 21 is set indicating the move delivers
// check to the opponent. Only guaranteed for moves freshly returned by
// Moves()/Captures()/QuietMoves(); reconstructed moves do not have it set.
func (m Move) GivesCheck() bool {
	return uint32(m)&moveGivesCheckMask != 0
}

// WithPromotion returns a copy of m with the given promotion piece type set.
// piece must be one of Knight, Bishop, Rook, Queen (color is ignored).
func (m Move) WithPromotion(piece gochess.Piece) Move {
	pt := gochess.PieceType(piece)
	idx, ok := promotionIndex[pt]
	if !ok {
		idx = 0
	}
	cleared := uint32(m) &^ (movePromoMask << movePromoShift)
	return Move(cleared | idx<<movePromoShift)
}

// WithCapturedPiece returns a copy of m with the captured piece type stored.
// piece is the piece that lived on the destination (or behind it for EP)
// before the move; color is ignored.
func (m Move) WithCapturedPiece(piece gochess.Piece) Move {
	pt := gochess.PieceType(piece)
	idx, ok := capturedIndex[pt]
	if !ok {
		idx = 0
	}
	cleared := uint32(m) &^ (moveCapturedMask << moveCapturedShift)
	return Move(cleared | idx<<moveCapturedShift)
}

// WithGivesCheck returns a copy of m with bit 21 set or cleared.
func (m Move) WithGivesCheck(givesCheck bool) Move {
	if givesCheck {
		return Move(uint32(m) | moveGivesCheckMask)
	}
	return Move(uint32(m) &^ moveGivesCheckMask)
}

// UCI returns the UCI string representation of the move.
func (m Move) UCI() string {
	from := m.From()
	to := m.To()
	if promo := m.Promotion(); promo != gochess.Empty {
		return UCI(from, to, promo)
	}
	return UCI(from, to)
}
