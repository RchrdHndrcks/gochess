package gochess

// Piece represents a chess piece with its color encoded as a bitfield.
type Piece int8

const (
	// White is the integer value of the white color.
	White Piece = 0b01000
	// Black is the integer value of the black color.
	Black Piece = 0b10000

	// Empty is the integer value of an empty square.
	Empty Piece = 0b00000
	// Pawn is the integer value of a pawn piece.
	Pawn Piece = 0b00001
	// Knight is the integer value of a knight piece.
	Knight Piece = 0b00010
	// Bishop is the integer value of a bishop piece.
	Bishop Piece = 0b00011
	// Rook is the integer value of a rook piece.
	Rook Piece = 0b00100
	// Queen is the integer value of a queen piece.
	Queen Piece = 0b00101
	// King is the integer value of a king piece.
	King Piece = 0b00110
)

// PieceColor returns the color portion of a piece (White or Black).
func PieceColor(piece Piece) Piece {
	return piece & (White | Black)
}

// PieceType returns the type portion of a piece (Pawn, Knight, Bishop, etc.)
// by stripping the color bits.
func PieceType(piece Piece) Piece {
	return piece &^ (White | Black)
}

var (
	// Colors is a map of color names to their integer values.
	Colors = map[string]Piece{
		"w": White,
		"b": Black,
	}

	// ColorNames is a map of color integer values to their names.
	ColorNames = map[Piece]string{
		White: "w",
		Black: "b",
	}

	// PiecesWithoutColor is a map of piece names to their integer values without color.
	PiecesWithoutColor = map[string]Piece{
		"p": Pawn, "P": Pawn,
		"n": Knight, "N": Knight,
		"b": Bishop, "B": Bishop,
		"r": Rook, "R": Rook,
		"q": Queen, "Q": Queen,
		"k": King, "K": King,
	}

	// Pieces is a map of piece names to their integer values.
	Pieces = map[string]Piece{
		"p": Black | Pawn, "P": White | Pawn,
		"n": Black | Knight, "N": White | Knight,
		"b": Black | Bishop, "B": White | Bishop,
		"r": Black | Rook, "R": White | Rook,
		"q": Black | Queen, "Q": White | Queen,
		"k": Black | King, "K": White | King,
	}

	// PieceNames is a map of piece integer values to their names.
	PieceNames = map[Piece]string{
		Black | Pawn: "p", White | Pawn: "P",
		Black | Knight: "n", White | Knight: "N",
		Black | Bishop: "b", White | Bishop: "B",
		Black | Rook: "r", White | Rook: "R",
		Black | Queen: "q", White | Queen: "Q",
		Black | King: "k", White | King: "K",
	}
)
