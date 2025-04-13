package pkg_test

import (
	"fmt"
	"testing"

	"github.com/RchrdHndrcks/gochess/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	t.Run("Valid Board Creation", func(t *testing.T) {
		// Arrange & Act
		board, err := pkg.NewBoard(8)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, board)
		assert.Equal(t, 8, board.Width())
	})

	t.Run("Invalid Width", func(t *testing.T) {
		// Arrange & Act
		board, err := pkg.NewBoard(0)

		// Assert
		require.Error(t, err)
		require.Nil(t, board)
		assert.Contains(t, err.Error(), "invalid width")
	})

	t.Run("With Valid Squares", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{pkg.Empty, pkg.Empty, pkg.Empty},
			{pkg.Empty, pkg.White | pkg.King, pkg.Empty},
			{pkg.Empty, pkg.Empty, pkg.Black | pkg.Queen},
		}

		// Act
		board, err := pkg.NewBoard(3, squares...)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, board)
		assert.Equal(t, 3, board.Width())

		// Verify squares were set correctly
		piece, err := board.Square(pkg.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, pkg.White|pkg.King, piece)

		piece, err = board.Square(pkg.Coor(2, 2))
		require.NoError(t, err)
		assert.Equal(t, pkg.Black|pkg.Queen, piece)
	})

	t.Run("With Invalid Squares Length", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{pkg.Empty, pkg.Empty},
			{pkg.Empty, pkg.White | pkg.King},
		}

		// Act
		board, err := pkg.NewBoard(3, squares...)

		// Assert
		require.Error(t, err)
		require.Nil(t, board)
		assert.Contains(t, err.Error(), "invalid squares length")
	})

	t.Run("With Invalid Row Length", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{pkg.Empty, pkg.Empty, pkg.Empty},
			{pkg.Empty, pkg.White | pkg.King},
			{pkg.Empty, pkg.Empty, pkg.Black | pkg.Queen},
		}

		// Act
		board, err := pkg.NewBoard(3, squares...)

		// Assert
		require.Error(t, err)
		require.Nil(t, board)
		assert.Contains(t, err.Error(), "invalid row length")
	})
}

func TestBoardWidth(t *testing.T) {
	t.Run("Returns Correct Width", func(t *testing.T) {
		// Arrange
		board, err := pkg.NewBoard(8)
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
			board, err := pkg.NewBoard(expectedWidth)
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
			{pkg.White | pkg.Rook, pkg.White | pkg.Knight, pkg.White | pkg.Bishop},
			{pkg.White | pkg.Pawn, pkg.Empty, pkg.Empty},
			{pkg.Empty, pkg.Black | pkg.Pawn, pkg.Empty},
		}

		board, err := pkg.NewBoard(3, squares...)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			coord    pkg.Coordinate
			expected int8
		}{
			{pkg.Coor(0, 0), pkg.White | pkg.Rook},
			{pkg.Coor(1, 0), pkg.White | pkg.Knight},
			{pkg.Coor(2, 0), pkg.White | pkg.Bishop},
			{pkg.Coor(0, 1), pkg.White | pkg.Pawn},
			{pkg.Coor(1, 1), pkg.Empty},
			{pkg.Coor(2, 1), pkg.Empty},
			{pkg.Coor(0, 2), pkg.Empty},
			{pkg.Coor(1, 2), pkg.Black | pkg.Pawn},
			{pkg.Coor(2, 2), pkg.Empty},
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
		board, err := pkg.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		invalidCoords := []pkg.Coordinate{
			pkg.Coor(-1, 0),
			pkg.Coor(0, -1),
			pkg.Coor(3, 0),
			pkg.Coor(0, 3),
			pkg.Coor(3, 3),
			pkg.Coor(-1, -1),
		}

		for _, coord := range invalidCoords {
			t.Run(fmt.Sprintf("Coordinate(%d,%d)", coord.X, coord.Y), func(t *testing.T) {
				// Act
				piece, err := board.Square(coord)

				// Assert
				require.Error(t, err)
				assert.Equal(t, pkg.Empty, piece)
				assert.Contains(t, err.Error(), "invalid coordinate")
			})
		}
	})
}

func TestBoardMakeMove(t *testing.T) {
	t.Run("Valid Move", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{pkg.Empty, pkg.Empty, pkg.Empty},
			{pkg.Empty, pkg.White | pkg.King, pkg.Empty},
			{pkg.Empty, pkg.Empty, pkg.Empty},
		}

		board, err := pkg.NewBoard(3, squares...)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(pkg.Coor(1, 1), pkg.Coor(2, 2))

		// Assert
		require.NoError(t, err)

		// Verify the move was made correctly
		originPiece, err := board.Square(pkg.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, pkg.Empty, originPiece)

		targetPiece, err := board.Square(pkg.Coor(2, 2))
		require.NoError(t, err)
		assert.Equal(t, pkg.White|pkg.King, targetPiece)
	})

	t.Run("Capture Move", func(t *testing.T) {
		// Arrange
		squares := [][]int8{
			{pkg.Empty, pkg.Empty, pkg.Empty},
			{pkg.Empty, pkg.White | pkg.King, pkg.Empty},
			{pkg.Empty, pkg.Empty, pkg.Black | pkg.Queen},
		}

		board, err := pkg.NewBoard(3, squares...)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(pkg.Coor(1, 1), pkg.Coor(2, 2))

		// Assert
		require.NoError(t, err)

		// Verify the move was made correctly
		originPiece, err := board.Square(pkg.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, pkg.Empty, originPiece)

		targetPiece, err := board.Square(pkg.Coor(2, 2))
		require.NoError(t, err)
		assert.Equal(t, pkg.White|pkg.King, targetPiece)
	})

	t.Run("Invalid Origin Coordinate", func(t *testing.T) {
		// Arrange
		board, err := pkg.NewBoard(3)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(pkg.Coor(-1, 0), pkg.Coor(1, 1))

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid coordinates")
	})

	t.Run("Invalid Target Coordinate", func(t *testing.T) {
		// Arrange
		board, err := pkg.NewBoard(3)
		require.NoError(t, err)

		// Act
		err = board.MakeMove(pkg.Coor(1, 1), pkg.Coor(3, 3))

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid coordinates")
	})
}

func TestBoardSetSquare(t *testing.T) {
	t.Run("Valid Coordinate", func(t *testing.T) {
		// Arrange
		board, err := pkg.NewBoard(3)
		require.NoError(t, err)

		// Act
		err = board.SetSquare(pkg.Coor(1, 1), pkg.White|pkg.King)

		// Assert
		require.NoError(t, err)

		// Verify the piece was set correctly
		piece, err := board.Square(pkg.Coor(1, 1))
		require.NoError(t, err)
		assert.Equal(t, pkg.White|pkg.King, piece)
	})

	t.Run("Multiple Pieces", func(t *testing.T) {
		// Arrange
		board, err := pkg.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			coord pkg.Coordinate
			piece int8
		}{
			{pkg.Coor(0, 0), pkg.White | pkg.Rook},
			{pkg.Coor(1, 0), pkg.White | pkg.Knight},
			{pkg.Coor(2, 0), pkg.White | pkg.Bishop},
			{pkg.Coor(0, 1), pkg.White | pkg.Pawn},
			{pkg.Coor(1, 1), pkg.Black | pkg.King},
			{pkg.Coor(2, 2), pkg.Black | pkg.Queen},
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
		board, err := pkg.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		invalidCoords := []pkg.Coordinate{
			pkg.Coor(-1, 0),
			pkg.Coor(0, -1),
			pkg.Coor(3, 0),
			pkg.Coor(0, 3),
			pkg.Coor(3, 3),
			pkg.Coor(-1, -1),
		}

		for _, coord := range invalidCoords {
			t.Run(fmt.Sprintf("Coordinate(%d,%d)", coord.X, coord.Y), func(t *testing.T) {
				// Act
				err := board.SetSquare(coord, pkg.White|pkg.King)

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
		board, err := pkg.NewBoard(3)
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			coord    pkg.Coordinate
			expected bool
		}{
			// Valid coordinates
			{pkg.Coor(0, 0), true},
			{pkg.Coor(1, 1), true},
			{pkg.Coor(2, 2), true},
			{pkg.Coor(0, 2), true},
			{pkg.Coor(2, 0), true},

			// Invalid coordinates
			{pkg.Coor(-1, 0), false},
			{pkg.Coor(0, -1), false},
			{pkg.Coor(3, 0), false},
			{pkg.Coor(0, 3), false},
			{pkg.Coor(3, 3), false},
			{pkg.Coor(-1, -1), false},
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
