package chess_test

import (
	"testing"

	"github.com/RchrdHndrcks/gochess"
	"github.com/RchrdHndrcks/gochess/chess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUCI(t *testing.T) {
	t.Run("a1 to a2", func(t *testing.T) {
		oCor := gochess.Coor(0, 7)
		tCor := gochess.Coor(0, 6)
		var promotionPiece int8
		expected := "a1a2"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("a1 to b2", func(t *testing.T) {
		oCor := gochess.Coor(0, 7)
		tCor := gochess.Coor(1, 6)
		var promotionPiece int8
		expected := "a1b2"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("a1 to b1", func(t *testing.T) {
		oCor := gochess.Coor(0, 7)
		tCor := gochess.Coor(1, 7)
		var promotionPiece int8
		expected := "a1b1"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("a1 to b3", func(t *testing.T) {
		oCor := gochess.Coor(0, 7)
		tCor := gochess.Coor(1, 5)
		var promotionPiece int8
		expected := "a1b3"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("a1 to h8", func(t *testing.T) {
		oCor := gochess.Coor(0, 7)
		tCor := gochess.Coor(7, 0)
		var promotionPiece int8
		expected := "a1h8"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("h8 promotion to queen", func(t *testing.T) {
		oCor := gochess.Coor(7, 1)
		tCor := gochess.Coor(7, 0)
		promotionPiece := gochess.White | gochess.Queen
		expected := "h7h8q"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("h8 promotion to knight", func(t *testing.T) {
		oCor := gochess.Coor(7, 1)
		tCor := gochess.Coor(7, 0)
		promotionPiece := gochess.White | gochess.Knight
		expected := "h7h8n"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("h8 promotion to bishop", func(t *testing.T) {
		oCor := gochess.Coor(7, 1)
		tCor := gochess.Coor(7, 0)
		promotionPiece := gochess.White | gochess.Bishop
		expected := "h7h8b"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("h8 promotion to rook", func(t *testing.T) {
		oCor := gochess.Coor(7, 1)
		tCor := gochess.Coor(7, 0)
		promotionPiece := gochess.White | gochess.Rook
		expected := "h7h8r"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})

	t.Run("h1 promotion to queen", func(t *testing.T) {
		oCor := gochess.Coor(7, 6)
		tCor := gochess.Coor(7, 7)
		promotionPiece := gochess.Black | gochess.Queen
		expected := "h2h1q"

		got := chess.UCI(oCor, tCor, promotionPiece)

		assert.Equal(t, expected, got)
	})
}

func TestAlgebraicToCoordinate(t *testing.T) {
	t.Run("a8", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("a8")

		// Assert
		expected := gochess.Coor(0, 0)
		assert.Equal(t, expected, got)
		assert.Nil(t, err)
	})

	t.Run("h1", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("h1")

		// Assert
		expected := gochess.Coor(7, 7)
		assert.Equal(t, expected, got)
		assert.Nil(t, err)
	})

	t.Run("a1", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("a1")

		// Assert
		expected := gochess.Coor(0, 7)
		assert.Equal(t, expected, got)
		assert.Nil(t, err)
	})

	t.Run("h8", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("h8")

		// Assert
		expected := gochess.Coor(7, 0)
		assert.Equal(t, expected, got)
		assert.Nil(t, err)
	})

	t.Run("e4", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("e4")

		// Assert
		expected := gochess.Coor(4, 4)
		assert.Equal(t, expected, got)
		assert.Nil(t, err)
	})

	t.Run("e3", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("e3")

		// Assert
		expected := gochess.Coor(4, 5)
		assert.Equal(t, expected, got)
		assert.Nil(t, err)
	})

	t.Run("d5", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("d5")

		// Assert
		expected := gochess.Coor(3, 3)
		assert.Equal(t, expected, got)
		assert.Nil(t, err)
	})

	t.Run("i9 - invalid coordinate", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("i9")

		// Assert
		expected := gochess.Coor(0, 0)
		assert.Equal(t, expected, got)
		require.NotNil(t, err)
		assert.Equal(t, "coordinate out of bounds", err.Error())
	})

	t.Run("ab3 - invalid coordinate", func(t *testing.T) {
		// Act
		got, err := chess.AlgebraicToCoordinate("ab3")

		// Assert
		expected := gochess.Coor(0, 0)
		assert.Equal(t, expected, got)
		require.NotNil(t, err)
		assert.Equal(t, "invalid text notation", err.Error())
	})
}

func TestCoordinateToAlgebraic(t *testing.T) {
	t.Run("0,0 to a8", func(t *testing.T) {
		// Act
		coor := gochess.Coor(0, 0)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := "a8"
		assert.Equal(t, expected, got)
	})

	t.Run("7,7 to h1", func(t *testing.T) {
		// Act
		coor := gochess.Coor(7, 7)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := "h1"
		assert.Equal(t, expected, got)
	})

	t.Run("0,7 to a1", func(t *testing.T) {
		// Act
		coor := gochess.Coor(0, 7)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := "a1"
		assert.Equal(t, expected, got)
	})

	t.Run("7,0 to h8", func(t *testing.T) {
		// Act
		coor := gochess.Coor(7, 0)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := "h8"
		assert.Equal(t, expected, got)
	})

	t.Run("4,4 to e4", func(t *testing.T) {
		// Act
		coor := gochess.Coor(4, 4)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := "e4"
		assert.Equal(t, expected, got)
	})

	t.Run("4,5 to e3", func(t *testing.T) {
		// Act
		coor := gochess.Coor(4, 5)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := "e3"
		assert.Equal(t, expected, got)
	})

	t.Run("3,3 to d5", func(t *testing.T) {
		// Act
		coor := gochess.Coor(3, 3)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := "d5"
		assert.Equal(t, expected, got)
	})

	t.Run("8,8 - invalid coordinate", func(t *testing.T) {
		// Act
		coor := gochess.Coor(8, 8)
		got := chess.CoordinateToAlgebraic(coor)

		// Assert
		expected := ""
		assert.Equal(t, expected, got)
	})
}
