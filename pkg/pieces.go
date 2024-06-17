package pkg

const (
	White int8 = 8
	Black int8 = 16

	Empty  int8 = 0
	Pawn   int8 = 1
	Knight int8 = 2
	Bishop int8 = 3
	Rook   int8 = 4
	Queen  int8 = 5
	King   int8 = 6
)

var (
	// Colors is a map of color names to their integer values.
	Colors = map[string]int8{
		"w": White,
		"b": Black,
	}

	// ColorNames is a map of color integer values to their names.
	ColorNames = map[int8]string{
		White: "w",
		Black: "b",
	}

	// Pieces is a map of piece names to their integer values.
	Pieces = map[string]int8{
		"p": Black | Pawn, "P": White | Pawn,
		"n": Black | Knight, "N": White | Knight,
		"b": Black | Bishop, "B": White | Bishop,
		"r": Black | Rook, "R": White | Rook,
		"q": Black | Queen, "Q": White | Queen,
		"k": Black | King, "K": White | King,
	}

	// PieceNames is a map of piece integer values to their names.
	PieceNames = map[int8]string{
		Black | Pawn: "p", White | Pawn: "P",
		Black | Knight: "n", White | Knight: "N",
		Black | Bishop: "b", White | Bishop: "B",
		Black | Rook: "r", White | Rook: "R",
		Black | Queen: "q", White | Queen: "Q",
		Black | King: "k", White | King: "K",
	}
)
