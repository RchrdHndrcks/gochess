package chess

import (
	"fmt"
	"strings"

	"github.com/RchrdHndrcks/gochess"
	chesspgn "github.com/RchrdHndrcks/gochess/chess/pgn"
)

// PGN generates a PGN string from the current game state.
//
// It writes the seven required tag pairs and the move text using UCI notation.
// Empty tag values default to "?". The Result tag is determined automatically
// if not provided: "1-0" or "0-1" for checkmate, "1/2-1/2" for stalemate,
// and "*" for an ongoing game.
func (c *Chess) PGN(tags chesspgn.PGNTags) string {
	var sb strings.Builder

	result := c.determineResult(tags.Result)

	// Write tag pairs.
	writeTag(&sb, "Event", tagValue(tags.Event))
	writeTag(&sb, "Site", tagValue(tags.Site))
	writeTag(&sb, "Date", tagValue(tags.Date))
	writeTag(&sb, "Round", tagValue(tags.Round))
	writeTag(&sb, "White", tagValue(tags.White))
	writeTag(&sb, "Black", tagValue(tags.Black))
	writeTag(&sb, "Result", result)

	sb.WriteString("\n")

	// Write moves.
	moveText := c.buildMoveText(result)
	sb.WriteString(wrapLines(moveText, 80))
	sb.WriteString("\n")

	return sb.String()
}

// determineResult returns the game result string.
func (c *Chess) determineResult(provided string) string {
	if provided != "" {
		return provided
	}

	if c.checkmate {
		if c.turn == gochess.White {
			return chesspgn.ResultBlackWins
		}
		return chesspgn.ResultWhiteWins
	}

	if c.stalemate {
		return chesspgn.ResultDraw
	}

	return chesspgn.ResultOngoing
}

// buildMoveText builds the move text from the game history.
func (c *Chess) buildMoveText(result string) string {
	var parts []string
	for i, ctx := range c.history {
		moveNum := i/2 + 1
		if i%2 == 0 {
			parts = append(parts, fmt.Sprintf("%d.", moveNum))
		}
		parts = append(parts, ctx.move)
	}
	parts = append(parts, result)
	return strings.Join(parts, " ")
}

// wrapLines wraps text at the given column width, breaking at spaces.
func wrapLines(text string, maxWidth int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var sb strings.Builder
	lineLen := 0

	for i, word := range words {
		if i == 0 {
			sb.WriteString(word)
			lineLen = len(word)
			continue
		}

		if lineLen+1+len(word) > maxWidth {
			sb.WriteString("\n")
			sb.WriteString(word)
			lineLen = len(word)
		} else {
			sb.WriteString(" ")
			sb.WriteString(word)
			lineLen += 1 + len(word)
		}
	}

	return sb.String()
}

// writeTag writes a PGN tag pair to the builder.
// It escapes backslashes and double quotes per PGN specification.
func writeTag(sb *strings.Builder, name, value string) {
	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	sb.WriteString(fmt.Sprintf("[%s \"%s\"]\n", name, escaped))
}

// tagValue returns the value or "?" if empty.
func tagValue(v string) string {
	if v == "" {
		return "?"
	}
	return v
}
