package gochess_test

import (
	"fmt"
	"testing"

	"github.com/RchrdHndrcks/gochess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	t.Run("Valid Board Creation", func(t *testing.T) {
		// Arrange & Act
		board, err := gochess.NewBoard(8)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, board)
		assert.Equal(t, 8, board.Width())
	})

	t.Run("Invalid Width", func(t *testing.T) {
		// Arrange & Act
		board, err := gochess.NewBoard(0)

		// Assert
		require.Error(t, err)
		require.Nil(t, board)
		assert.Contains(t, err.Error(), "invalid width")
	})

	t.Run("With Valid Squares", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{gochess.Empty, gochess.Empty, gochess.Empty},
			{gochess.Empty, gochess.White | gochess.King, gochess.Empty},
			{gochess.Empty, gochess.Empty, gochess.Black | gochess.Queen},
		}

		// Act
		board, err := gochess.NewBoard(3, squares...)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, board)
		assert.Equal(t, 3, board.Width())

		// Verify squares were set correctly
		piece, err := board.Square(gochess.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, gochess.White|gochess.King, piece)

		piece, err = board.Square(gochess.Coor(2, 2))
		require.NoError(t, err)
		assert.Equal(t, gochess.Black|gochess.Queen, piece)
	})

	t.Run("With Invalid Squares Length", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{gochess.Empty, gochess.Empty},
			{gochess.Empty, gochess.White | gochess.King},
		}

		// Act
		board, err := gochess.NewBoard(3, squares...)

		// Assert
		require.Error(t, err)
		require.Nil(t, board)
		assert.EqualError(t, err, "board: invalid square: rows count 2 is not equal to width 3")
	})

	t.Run("With Invalid Row Length", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{gochess.Empty, gochess.Empty, gochess.Empty},
			{gochess.Empty, gochess.White | gochess.King},
			{gochess.Empty, gochess.Empty, gochess.Black | gochess.Queen},
		}

		// Act
		board, err := gochess.NewBoard(3, squares...)

		// Assert
		require.Error(t, err)
		require.Nil(t, board)
		assert.EqualError(t, err, "board: invalid square: row 1 has 2 columns, expected 3")
	})
}

func TestBoardWidth(t *testing.T) {
	t.Run("Returns Correct Width", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(8)
		require.NoError(t, err)

		// Act
		width := board.Width()

		// Assert
		assert.Equal(t, 8, width)
	})

	t.Run("Different Widths", func(t *testing.T) {
		// Test different board sizes
		widths := []int{3, 4, 5, 8, 10}

		for _, expectedWidth := range widths {
			// Arrange
			board, err := gochess.NewBoard(expectedWidth)
			require.NoError(t, err)

			// Act
			actualWidth := board.Width()

			// Assert
			assert.Equal(t, expectedWidth, actualWidth)
		}
	})
}

func TestBoardSquare(t *testing.T) {
	t.Run("Valid Coordinates", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{gochess.White | gochess.Rook, gochess.White | gochess.Knight, gochess.White | gochess.Bishop},
			{gochess.White | gochess.Pawn, gochess.Empty, gochess.Empty},
			{gochess.Empty, gochess.Black | gochess.Pawn, gochess.Empty},
		}

		board, err := gochess.NewBoard(3, squares...)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			coord    gochess.Coordinate
			expected int8
		}{
			{gochess.Coor(0, 0), gochess.White | gochess.Rook},
			{gochess.Coor(1, 0), gochess.White | gochess.Knight},
			{gochess.Coor(2, 0), gochess.White | gochess.Bishop},
			{gochess.Coor(0, 1), gochess.White | gochess.Pawn},
			{gochess.Coor(1, 1), gochess.Empty},
			{gochess.Coor(2, 1), gochess.Empty},
			{gochess.Coor(0, 2), gochess.Empty},
			{gochess.Coor(1, 2), gochess.Black | gochess.Pawn},
			{gochess.Coor(2, 2), gochess.Empty},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Coordinate(%d,%d)", tc.coord.X, tc.coord.Y), func(t *testing.T) {
				// Act
				piece, err := board.Square(tc.coord)

				// Assert
				require.NoError(t, err)
				assert.Equal(t, tc.expected, piece)
			})
		}
	})

	t.Run("Invalid Coordinates", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		invalidCoords := []gochess.Coordinate{
			gochess.Coor(-1, 0),
			gochess.Coor(0, -1),
			gochess.Coor(3, 0),
			gochess.Coor(0, 3),
			gochess.Coor(3, 3),
			gochess.Coor(-1, -1),
		}

		for _, coord := range invalidCoords {
			t.Run(fmt.Sprintf("Coordinate(%d,%d)", coord.X, coord.Y), func(t *testing.T) {
				// Act
				piece, err := board.Square(coord)

				// Assert
				require.Error(t, err)
				assert.Equal(t, gochess.Empty, piece)
				assert.Contains(t, err.Error(), "invalid coordinate")
			})
		}
	})
}

func TestBoardMakeMove(t *testing.T) {
	t.Run("Valid Move", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{gochess.Empty, gochess.Empty, gochess.Empty},
			{gochess.Empty, gochess.White | gochess.King, gochess.Empty},
			{gochess.Empty, gochess.Empty, gochess.Empty},
		}

		board, err := gochess.NewBoard(3, squares...)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(gochess.Coor(1, 1), gochess.Coor(2, 2))

		// Assert
		require.NoError(t, err)

		// Verify the move was made correctly
		originPiece, err := board.Square(gochess.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, gochess.Empty, originPiece)

		targetPiece, err := board.Square(gochess.Coor(2, 2))
		require.NoError(t, err)
		assert.Equal(t, gochess.White|gochess.King, targetPiece)
	})

	t.Run("Capture Move", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{gochess.Empty, gochess.Empty, gochess.Empty},
			{gochess.Empty, gochess.White | gochess.King, gochess.Empty},
			{gochess.Empty, gochess.Empty, gochess.Black | gochess.Queen},
		}

		board, err := gochess.NewBoard(3, squares...)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(gochess.Coor(1, 1), gochess.Coor(2, 2))

		// Assert
		require.NoError(t, err)

		// Verify the move was made correctly
		originPiece, err := board.Square(gochess.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, gochess.Empty, originPiece)

		targetPiece, err := board.Square(gochess.Coor(2, 2))
		require.NoError(t, err)
		assert.Equal(t, gochess.White|gochess.King, targetPiece)
	})

	t.Run("Invalid Origin Coordinate", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(3)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(gochess.Coor(-1, 0), gochess.Coor(1, 1))

		// Assert
		require.Error(t, err)
		assert.EqualError(t, err, "board: invalid coordinate: (-1,0)")
	})

	t.Run("Invalid Target Coordinate", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(3)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(gochess.Coor(1, 1), gochess.Coor(3, 3))

		// Assert
		require.Error(t, err)
		assert.EqualError(t, err, "board: invalid coordinate: (3,3)")
	})
}

func TestBoardSetSquare(t *testing.T) {
	t.Run("Valid Coordinate", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(3)
		require.NoError(t, err)

		// Act
		err = board.SetSquare(gochess.Coor(1, 1), gochess.White|gochess.King)

		// Assert
		require.NoError(t, err)

		// Verify the piece was set correctly
		piece, err := board.Square(gochess.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, gochess.White|gochess.King, piece)
	})

	t.Run("Multiple Pieces", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			coord gochess.Coordinate
			piece int8
		}{
			{gochess.Coor(0, 0), gochess.White | gochess.Rook},
			{gochess.Coor(1, 0), gochess.White | gochess.Knight},
			{gochess.Coor(2, 0), gochess.White | gochess.Bishop},
			{gochess.Coor(0, 1), gochess.White | gochess.Pawn},
			{gochess.Coor(1, 1), gochess.Black | gochess.King},
			{gochess.Coor(2, 2), gochess.Black | gochess.Queen},
		}

		for _, tc := range testCases {
			// Act
			err := board.SetSquare(tc.coord, tc.piece)

			// Assert
			require.NoError(t, err)

			// Verify the piece was set correctly
			piece, err := board.Square(tc.coord)
			require.NoError(t, err)
			assert.Equal(t, tc.piece, piece)
		}
	})

	t.Run("Invalid Coordinate", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		invalidCoords := []gochess.Coordinate{
			gochess.Coor(-1, 0),
			gochess.Coor(0, -1),
			gochess.Coor(3, 0),
			gochess.Coor(0, 3),
			gochess.Coor(3, 3),
			gochess.Coor(-1, -1),
		}

		for _, coord := range invalidCoords {
			t.Run(fmt.Sprintf("Coordinate(%d,%d)", coord.X, coord.Y), func(t *testing.T) {
				// Act
				err := board.SetSquare(coord, gochess.White|gochess.King)

				// Assert
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid coordinate")
			})
		}
	})
}

func TestBoardIsValidCoordinate(t *testing.T) {
	// Since isValidCoordinate is a private method, we'll test it indirectly through Square
	t.Run("Valid and Invalid Coordinates", func(t *testing.T) {
		// Arrange
		board, err := gochess.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			coord    gochess.Coordinate
			expected bool
		}{
			// Valid coordinates
			{gochess.Coor(0, 0), true},
			{gochess.Coor(1, 1), true},
			{gochess.Coor(2, 2), true},
			{gochess.Coor(0, 2), true},
			{gochess.Coor(2, 0), true},

			// Invalid coordinates
			{gochess.Coor(-1, 0), false},
			{gochess.Coor(0, -1), false},
			{gochess.Coor(3, 0), false},
			{gochess.Coor(0, 3), false},
			{gochess.Coor(3, 3), false},
			{gochess.Coor(-1, -1), false},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Coordinate(%d,%d)", tc.coord.X, tc.coord.Y), func(t *testing.T) {
				// Act
				_, err := board.Square(tc.coord)

				// Assert
				if tc.expected {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "invalid coordinate")
				}
			})
		}
	})
}
