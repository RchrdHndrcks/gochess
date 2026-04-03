package chess

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/RchrdHndrcks/gochess"
)

// SAN converts a UCI move (like "e2e4") to Standard Algebraic Notation (like "e4").
//
// The move must be present in AvailableMoves(). Disambiguation, captures, check
// and checkmate suffixes are determined automatically from the current position.
func (c *Chess) SAN(uciMove string) (string, error) {
	if len(uciMove) < 4 || len(uciMove) > 5 {
		return "", fmt.Errorf("invalid UCI move: %s", uciMove)
	}

	// Verify the move is legal.
	found := false
	for _, m := range c.AvailableMoves() {
		if m == uciMove {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("move is not legal: %s", uciMove)
	}

	origin, _ := AlgebraicToCoordinate(uciMove[:2])
	target, _ := AlgebraicToCoordinate(uciMove[2:4])

	// Use FEN to get piece info, as the board may be temporarily modified
	// by legal move calculations in sequential mode.
	piece := pieceFromFEN(c.FEN(), origin)
	pieceType := piece &^ (gochess.White | gochess.Black)

	var san string

	// Handle castling.
	isCastle := pieceType == gochess.King && castlesMoves[uciMove] == c.turn
	if isCastle {
		if target.X > origin.X {
			san = "O-O"
		} else {
			san = "O-O-O"
		}
		return san + checkSuffix(c, uciMove), nil
	}

	// Determine capture.
	targetPiece := pieceFromFEN(c.FEN(), target)
	isCapture := targetPiece != gochess.Empty

	// En passant is also a capture.
	if pieceType == gochess.Pawn && c.enPassantSquare != "" && uciMove[2:4] == c.enPassantSquare {
		isCapture = true
	}

	if pieceType == gochess.Pawn {
		san = pawnSAN(origin, target, isCapture, uciMove)
	} else {
		// Piece letter — reuse gochess.PieceNames (uppercase = white-colored pieces).
		san = gochess.PieceNames[pieceType|gochess.White]

		// Disambiguation.
		san += disambiguation(c, piece, origin, target)

		if isCapture {
			san += "x"
		}

		san += CoordinateToAlgebraic(target)
	}

	return san + checkSuffix(c, uciMove), nil
}

// FromSAN converts a SAN string (like "Nf3") to a UCI move (like "g1f3").
//
// The SAN must correspond to a legal move in the current position.
func (c *Chess) FromSAN(san string) (string, error) {
	san = strings.TrimRight(san, "+#")

	// Handle castling.
	if san == "O-O" || san == "0-0" {
		return findCastleMove(c, true)
	}
	if san == "O-O-O" || san == "0-0-0" {
		return findCastleMove(c, false)
	}

	if len(san) == 0 {
		return "", fmt.Errorf("invalid SAN: empty string")
	}

	if unicode.IsUpper(rune(san[0])) && san[0] != 'O' {
		return parsePieceMoveSAN(c, san)
	}

	return parsePawnMoveSAN(c, san)
}

// pawnSAN builds the SAN string for a pawn move.
func pawnSAN(origin, target gochess.Coordinate, isCapture bool, uciMove string) string {
	var san string

	if isCapture {
		san += string(rune('a' + origin.X))
		san += "x"
	}

	san += CoordinateToAlgebraic(target)

	// Promotion.
	if len(uciMove) == 5 {
		promoChar := strings.ToUpper(uciMove[4:5])
		san += "=" + promoChar
	}

	return san
}

// disambiguation returns the disambiguation string needed for a piece move.
func disambiguation(c *Chess, piece int8, origin, target gochess.Coordinate) string {
	pieceType := piece &^ (gochess.White | gochess.Black)
	targetAlg := CoordinateToAlgebraic(target)
	fen := c.FEN()

	sameFile := false
	sameRank := false
	ambiguous := false

	for _, m := range c.AvailableMoves() {
		if len(m) < 4 {
			continue
		}

		mOrigin, _ := AlgebraicToCoordinate(m[:2])
		if mOrigin == origin {
			continue
		}

		if m[2:4] != targetAlg {
			continue
		}

		mPiece := pieceFromFEN(fen, mOrigin)
		mPieceType := mPiece &^ (gochess.White | gochess.Black)
		if mPieceType != pieceType {
			continue
		}

		ambiguous = true
		if mOrigin.X == origin.X {
			sameFile = true
		}
		if mOrigin.Y == origin.Y {
			sameRank = true
		}
	}

	if !ambiguous {
		return ""
	}

	if sameFile && sameRank {
		return CoordinateToAlgebraic(origin)
	}

	if sameFile {
		return fmt.Sprintf("%d", 8-origin.Y)
	}

	return string(rune('a' + origin.X))
}

// checkSuffix determines if a move results in check or checkmate.
func checkSuffix(c *Chess, uciMove string) string {
	cloned := c.clone()
	_ = cloned.MakeMove(uciMove)

	if cloned.IsCheckmate() {
		return "#"
	}

	if cloned.IsCheck() {
		return "+"
	}

	return ""
}

// findCastleMove finds the castling UCI move from available moves.
func findCastleMove(c *Chess, kingside bool) (string, error) {
	for _, m := range c.AvailableMoves() {
		if !c.isCastleMove(m) {
			continue
		}

		origin, _ := AlgebraicToCoordinate(m[:2])
		target, _ := AlgebraicToCoordinate(m[2:4])

		if kingside && target.X > origin.X {
			return m, nil
		}
		if !kingside && target.X < origin.X {
			return m, nil
		}
	}

	return "", fmt.Errorf("castling move not available")
}

// parsePieceMoveSAN parses SAN for non-pawn pieces (e.g., "Nf3", "Raxe1", "R1e1").
func parsePieceMoveSAN(c *Chess, san string) (string, error) {
	pieceChar := san[0]
	// Reuse gochess.Pieces; strip color to get bare piece type.
	p, ok := gochess.Pieces[string(pieceChar)]
	if !ok || p == gochess.Empty {
		return "", fmt.Errorf("invalid piece in SAN: %c", pieceChar)
	}
	pieceType := p &^ (gochess.White | gochess.Black)

	rest := san[1:]
	rest = strings.ReplaceAll(rest, "x", "")

	if len(rest) < 2 {
		return "", fmt.Errorf("invalid SAN: %s", san)
	}

	targetAlg := rest[len(rest)-2:]
	disambig := rest[:len(rest)-2]

	var fileDisambig int = -1
	var rankDisambig int = -1

	for _, ch := range disambig {
		if ch >= 'a' && ch <= 'h' {
			fileDisambig = int(ch - 'a')
		} else if ch >= '1' && ch <= '8' {
			rankDisambig = 8 - int(ch-'0')
		}
	}

	fen := c.FEN()
	for _, m := range c.AvailableMoves() {
		if m[2:4] != targetAlg {
			continue
		}

		mOrigin, _ := AlgebraicToCoordinate(m[:2])
		mPiece := pieceFromFEN(fen, mOrigin)
		mPieceType := mPiece &^ (gochess.White | gochess.Black)

		if mPieceType != pieceType {
			continue
		}

		if fileDisambig >= 0 && mOrigin.X != fileDisambig {
			continue
		}

		if rankDisambig >= 0 && mOrigin.Y != rankDisambig {
			continue
		}

		return m, nil
	}

	return "", fmt.Errorf("no matching move found for SAN: %s", san)
}

// parsePawnMoveSAN parses SAN for pawn moves (e.g., "e4", "exd5", "e8=Q").
func parsePawnMoveSAN(c *Chess, san string) (string, error) {
	var fileDisambig int = -1
	var promotion string
	isCaptureSAN := false

	// Check for promotion.
	if idx := strings.Index(san, "="); idx >= 0 {
		if idx+1 >= len(san) {
			return "", fmt.Errorf("invalid SAN: promotion piece missing after '=': %s", san)
		}
		promotion = strings.ToLower(san[idx+1 : idx+2])
		san = san[:idx]
	}

	// Remove capture indicator.
	parts := strings.Split(san, "x")
	if len(parts) == 2 {
		if len(parts[0]) == 0 {
			return "", fmt.Errorf("invalid capture SAN: missing file before 'x': %s", san)
		}
		fileDisambig = int(parts[0][0] - 'a')
		isCaptureSAN = true
		san = parts[1]
	}

	targetAlg := san
	if len(targetAlg) != 2 {
		return "", fmt.Errorf("invalid pawn move SAN: %s", san)
	}

	fen := c.FEN()
	for _, m := range c.AvailableMoves() {
		if m[2:4] != targetAlg {
			continue
		}

		mOrigin, _ := AlgebraicToCoordinate(m[:2])
		mPiece := pieceFromFEN(fen, mOrigin)
		mPieceType := mPiece &^ (gochess.White | gochess.Black)

		if mPieceType != gochess.Pawn {
			continue
		}

		if fileDisambig >= 0 && mOrigin.X != fileDisambig {
			continue
		}

		// A non-capture SAN (no 'x') must not match diagonal moves (captures/en passant).
		mTarget, _ := AlgebraicToCoordinate(m[2:4])
		if !isCaptureSAN && mOrigin.X != mTarget.X {
			continue
		}

		// Check promotion match.
		if promotion != "" {
			if len(m) != 5 || m[4:5] != promotion {
				continue
			}
		} else if len(m) == 5 {
			continue
		}

		return m, nil
	}

	return "", fmt.Errorf("no matching move found for pawn SAN: %s", targetAlg)
}
