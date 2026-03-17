package chess

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/RchrdHndrcks/gochess"
)

// pieceToSAN maps piece types (without color) to their SAN letter.
var pieceToSAN = map[int8]string{
	gochess.King:   "K",
	gochess.Queen:  "Q",
	gochess.Rook:   "R",
	gochess.Bishop: "B",
	gochess.Knight: "N",
}

// sanToPiece maps SAN piece letters to piece types (without color).
var sanToPiece = map[byte]int8{
	'K': gochess.King,
	'Q': gochess.Queen,
	'R': gochess.Rook,
	'B': gochess.Bishop,
	'N': gochess.Knight,
}

// ToSAN converts a UCI move (like "e2e4") to Standard Algebraic Notation (like "e4").
//
// The function requires the current Chess state to determine disambiguation,
// captures, check, and checkmate. The UCI move must be present in AvailableMoves().
func ToSAN(c *Chess, uciMove string) (string, error) {
	if len(uciMove) < 4 || len(uciMove) > 5 {
		return "", fmt.Errorf("invalid UCI move: %s", uciMove)
	}

	// Check that the move is legal.
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

	// Handle castling. Check using piece type and known castle moves.
	isCastle := pieceType == gochess.King && castlesMoves[uciMove] == c.turn
	if isCastle {
		if target.X > origin.X {
			san = "O-O"
		} else {
			san = "O-O-O"
		}
		return san + checkSuffix(c, uciMove), nil
	}

	// Determine capture by checking the FEN rather than the board directly,
	// as the board may be temporarily modified by legal move calculations.
	targetPiece := pieceFromFEN(c.FEN(), target)
	isCapture := targetPiece != gochess.Empty

	// En passant is also a capture: pawn moves diagonally to en passant square.
	if pieceType == gochess.Pawn && c.enPassantSquare != "" && uciMove[2:4] == c.enPassantSquare {
		isCapture = true
	}

	if pieceType == gochess.Pawn {
		san = pawnSAN(origin, target, isCapture, uciMove)
	} else {
		// Piece letter.
		san = pieceToSAN[pieceType]

		// Disambiguation.
		san += disambiguation(c, piece, origin, target)

		if isCapture {
			san += "x"
		}

		san += CoordinateToAlgebraic(target)
	}

	return san + checkSuffix(c, uciMove), nil
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
//
// It checks all legal moves to see if another piece of the same type can
// reach the same target square. If so, it adds file, rank, or both.
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

		// Skip this move itself.
		mOrigin, _ := AlgebraicToCoordinate(m[:2])
		if mOrigin == origin {
			continue
		}

		// Check if this move targets the same square.
		if m[2:4] != targetAlg {
			continue
		}

		// Check if the piece at origin is the same type.
		// Use FEN to avoid board corruption issues.
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
		// Use rank for disambiguation.
		return fmt.Sprintf("%d", 8-origin.Y)
	}

	// Use file for disambiguation (default).
	return string(rune('a' + origin.X))
}

// checkSuffix determines if a move results in check or checkmate by temporarily
// making the move and checking the resulting position state.
func checkSuffix(c *Chess, uciMove string) string {
	// Clone the chess to avoid modifying the original.
	cloned := c.clone()

	// Make the move on the clone.
	_ = cloned.MakeMove(uciMove)

	if cloned.IsCheckmate() {
		return "#"
	}

	if cloned.IsCheck() {
		return "+"
	}

	return ""
}

// FromSAN converts a SAN string (like "Nf3") to a UCI move (like "g1f3").
//
// The function requires the current Chess state to find the matching move
// among AvailableMoves().
func FromSAN(c *Chess, san string) (string, error) {
	san = strings.TrimRight(san, "+#")

	// Handle castling.
	if san == "O-O" || san == "0-0" {
		return findCastleMove(c, true)
	}
	if san == "O-O-O" || san == "0-0-0" {
		return findCastleMove(c, false)
	}

	// Determine if it's a piece move or pawn move.
	if len(san) == 0 {
		return "", fmt.Errorf("invalid SAN: empty string")
	}

	if unicode.IsUpper(rune(san[0])) && san[0] != 'O' {
		return parsePieceMoveSAN(c, san)
	}

	return parsePawnMoveSAN(c, san)
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
	pieceType, ok := sanToPiece[pieceChar]
	if !ok {
		return "", fmt.Errorf("invalid piece in SAN: %c", pieceChar)
	}

	rest := san[1:]

	// Remove capture indicator.
	rest = strings.ReplaceAll(rest, "x", "")

	if len(rest) < 2 {
		return "", fmt.Errorf("invalid SAN: %s", san)
	}

	// The last two characters are always the target square.
	targetAlg := rest[len(rest)-2:]
	disambig := rest[:len(rest)-2]

	// Parse disambiguation.
	var fileDisambig int = -1
	var rankDisambig int = -1

	for _, ch := range disambig {
		if ch >= 'a' && ch <= 'h' {
			fileDisambig = int(ch - 'a')
		} else if ch >= '1' && ch <= '8' {
			rankDisambig = 8 - int(ch-'0')
		}
	}

	// Find matching move.
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

	// Check for promotion.
	if idx := strings.Index(san, "="); idx >= 0 {
		promotion = strings.ToLower(san[idx+1 : idx+2])
		san = san[:idx]
	}

	// Remove capture indicator.
	parts := strings.Split(san, "x")
	if len(parts) == 2 {
		fileDisambig = int(parts[0][0] - 'a')
		san = parts[1]
	} else if len(parts[0]) == 1 && parts[0][0] >= 'a' && parts[0][0] <= 'h' {
		// Could be just a destination like "e4", but we handle below.
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

// MoveToSAN converts a UCI move string to Standard Algebraic Notation.
//
// The move must be present in AvailableMoves().
func (c *Chess) MoveToSAN(move string) (string, error) {
	return ToSAN(c, move)
}

// MoveFromSAN converts a SAN string to a UCI move string.
//
// The SAN must correspond to a legal move in the current position.
func (c *Chess) MoveFromSAN(san string) (string, error) {
	return FromSAN(c, san)
}
