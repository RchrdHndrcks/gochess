package chess

import (
	"fmt"
	"runtime"
	"slices"

	"github.com/RchrdHndrcks/gochess"
)

type (
	// Cloner is a generic interface for objects that can be cloned.
	Cloner interface {
		// Clone returns a clone of the object.
		Clone() Board
	}

	// Board represents a chess board.
	Board interface {
		// SetSquare sets a piece in a square.
		SetSquare(c gochess.Coordinate, p int8) error
		// Square returns the piece in a square.
		Square(c gochess.Coordinate) (int8, error)
		// Width returns the width of the board.
		Width() int
	}

	// config represents configurations of how the methods will work.
	config struct {
		// Parallelism is the number of workers to use for the moves calculation.
		Parallelism int
	}

	// chessContext represents the history of a game.
	chessContext struct {
		// move is a played move.
		move string
		// fen is a FEN strings that represents the position.
		fen string
		// halfMove is the number of half moves since the last capture or pawn move.
		halfMove int
		// availableCastles is the castles that are available.
		availableCastles string
		// enPassantSquare is the square where a pawn can capture in passant.
		enPassantSquare string
		// whiteKingPosition is the position of the white king.
		whiteKingPosition *gochess.Coordinate
		// blackKingPosition is the position of the black king.
		blackKingPosition *gochess.Coordinate
		// check is true if the current turn is in check.
		check bool
		// checkmate is true if the current turn is in checkmate.
		checkmate bool
		// stalemate is true if the current turn is in stalemate.
		stalemate bool
	}

	// Chess represents a Chess game.
	Chess struct {
		board Board
		// turn is the current turn.
		turn int8
		// movesCount is the number of moves played in algebaric notation.
		// It will increase by 1 after each black move.
		movesCount uint64
		// halfMoves is the number of half moves since the last capture or pawn move.
		halfMoves int
		// enPassantSquare is the square where a pawn can capture in passant.
		enPassantSquare string
		// availableCastles is the castles that are available.
		// It will has the same format as the FEN castles.
		availableCastles string
		// moves are the available moves in the current position.
		moves []string
		// actualFEN is the FEN string of the current position.
		actualFEN string
		// blackKingPosition is the position of the black king.
		blackKingPosition *gochess.Coordinate
		// whiteKingPosition is the position of the white king.
		whiteKingPosition *gochess.Coordinate
		// check is true if the current turn is in check.
		check bool
		// checkmate is true if the current turn is in checkmate.
		checkmate bool
		// stalemate is true if the current turn is in stalemate.
		stalemate bool

		// config represents configurations of how the methods will work.
		config config

		// history is the history of the game.
		history []chessContext
	}
)

var (
	castlesMoves = map[string]int8{
		"e1g1": gochess.White,
		"e1c1": gochess.White,
		"e8g8": gochess.Black,
		"e8c8": gochess.Black,
	}

	castleKingWay = map[string]gochess.Coordinate{
		"e1g1": gochess.Coor(5, 7),
		"e1c1": gochess.Coor(3, 7),
		"e8g8": gochess.Coor(5, 0),
		"e8c8": gochess.Coor(3, 0),
	}

	castleRook = map[string]gochess.Coordinate{
		"e1g1": gochess.Coor(7, 7),
		"e1c1": gochess.Coor(0, 7),
		"e8g8": gochess.Coor(7, 0),
		"e8c8": gochess.Coor(0, 0),
	}
)

// New creates a new chess game.
//
// The function will up a pool of workers to maximize performance on the
// moves calculation.
// By default, the number of workers is twice the number of available CPUs.
// If you are running on a container environment or you want to use the
// sequential version, you should set this value manually using the
// WithParallelism option.
// The pool of workers will be used only if the board implements the Cloner
// interface. If you are using a custom Board implementation, you should
// implement the Cloner interface to take advantage of the parallelism.
func New(opts ...Option) (*Chess, error) {
	c := &Chess{
		board:            newBoardAdapter(gochess.DefaultChessBoard()),
		turn:             gochess.White,
		movesCount:       1,
		halfMoves:        0,
		enPassantSquare:  "",
		availableCastles: "KQkq",
		actualFEN:        "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		moves: []string{
			"a2a3", "a2a4", "b2b3", "b2b4", "c2c3", "c2c4", "d2d3", "d2d4",
			"e2e3", "e2e4", "f2f3", "f2f4", "g2g3", "g2g4", "h2h3", "h2h4",
			"b1a3", "b1c3", "g1f3", "g1h3",
		},
		blackKingPosition: &gochess.Coordinate{X: 4, Y: 0},
		whiteKingPosition: &gochess.Coordinate{X: 4, Y: 7},
		check:             false,
		checkmate:         false,
		stalemate:         false,
		config: config{
			// To maximize performance chess uses twice the number of available
			// CPUs. If you are running on a container environment or you want to
			// use the sequential version, you should set this value manually.
			Parallelism: runtime.NumCPU() * 2,
		},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return c, nil
}

// LoadPosition loads a board from a FEN string.
//
// The function will read the entire FEN string and will return an error if
// the FEN string is invalid.
//
// The board and properties will not be modified if the FEN string is invalid.
func (c *Chess) LoadPosition(FEN string) error {
	if err := c.loadPosition(FEN); err != nil {
		return err
	}

	c.actualFEN = FEN
	c.moves = c.legalMoves()
	check := c.isCheck()
	c.check = check && len(c.moves) > 0
	c.checkmate = check && len(c.moves) == 0
	c.stalemate = !check && len(c.moves) == 0
	return nil
}

// FEN returns the FEN string of the current position.
//
// If any of the kings is not in the board, the function returns an empty string.
func (c *Chess) FEN() string {
	return c.actualFEN
}

// AvailableMoves returns the available legal moves for the current turn.
//
// It always returns a non nil slice. It could be empty if the position is
// checkmate or stalemate.
func (c *Chess) AvailableMoves() []string {
	return c.moves
}

// MakeMove checks if the move is legal and makes it.
// It returns an error if the move is not legal.
func (c *Chess) MakeMove(move string) error {
	if !slices.Contains(c.moves, move) {
		return fmt.Errorf("move is not legal: %s", move)
	}

	c.makeMove(move)
	c.actualFEN = c.calculateFEN(move)
	c.moves = c.legalMoves()
	check := c.isCheck()
	c.check = check && len(c.moves) > 0
	c.checkmate = check && len(c.moves) == 0
	c.stalemate = !check && len(c.moves) == 0
	return nil
}

// UnmakeMove unmake the last move.
//
// It searches for the last move in the history and unmake it.
// If there are no moves in the history, the function does nothing.
func (c *Chess) UnmakeMove() {
	c.unmakeMove()
	c.moves = c.legalMoves()
}

// IsCheck returns if the current turn is in check.
func (c *Chess) IsCheck() bool {
	return c.check
}

// IsCheckmate returns if the current turn is in checkmate.
func (c *Chess) IsCheckmate() bool {
	return c.checkmate
}

// IsStalemate returns if the current turn is in stalemate.
func (c *Chess) IsStalemate() bool {
	return c.stalemate
}

// Square returns the piece in a square.
// The square is represented by an algebraic notation.
//
// If the square is not valid, the function returns an error.
func (c *Chess) Square(square string) (string, error) {
	coor, err := AlgebraicToCoordinate(square)
	if err != nil {
		return "", fmt.Errorf("failed to convert algebraic notation to coordinate: %w", err)
	}

	// Ignore the error because the coordinates is
	// already validated.
	p, _ := c.board.Square(coor)
	return gochess.PieceNames[p], nil
}

// clone creates a deep copy of the Chess structure.
func (c Chess) clone() Chess {
	cloned := c

	cloner, ok := c.board.(Cloner)
	if ok {
		cloned.board = cloner.Clone()
	}

	if c.whiteKingPosition != nil {
		whitePos := *c.whiteKingPosition
		cloned.whiteKingPosition = &whitePos
	}
	if c.blackKingPosition != nil {
		blackPos := *c.blackKingPosition
		cloned.blackKingPosition = &blackPos
	}

	if len(c.history) > 0 {
		cloned.history = make([]chessContext, len(c.history))
		for i, ctx := range c.history {
			cloned.history[i] = ctx
			if ctx.whiteKingPosition != nil {
				whitePos := *ctx.whiteKingPosition
				cloned.history[i].whiteKingPosition = &whitePos
			}
			if ctx.blackKingPosition != nil {
				blackPos := *ctx.blackKingPosition
				cloned.history[i].blackKingPosition = &blackPos
			}
		}
	} else {
		cloned.history = nil
	}

	if len(c.moves) > 0 {
		cloned.moves = make([]string, len(c.moves))
		copy(cloned.moves, c.moves)
	}

	return cloned
}
