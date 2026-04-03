package chess

import (
	"fmt"
	"strings"

	"github.com/RchrdHndrcks/gochess"
)

// PGN result constants per the PGN specification.
const (
	ResultWhiteWins = "1-0"
	ResultBlackWins = "0-1"
	ResultDraw      = "1/2-1/2"
	ResultOngoing   = "*"
)

// PGNTags represents the seven required tag pairs in a PGN file.
type PGNTags struct {
	Event  string
	Site   string
	Date   string
	Round  string
	White  string
	Black  string
	Result string
}

// PGN generates a PGN string from the current game state.
//
// It writes the seven required tag pairs and the move text using UCI notation.
// Empty tag values default to "?". The Result tag is determined automatically
// if not provided: "1-0" or "0-1" for checkmate, "1/2-1/2" for stalemate,
// and "*" for an ongoing game.
func (c *Chess) PGN(tags PGNTags) string {
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

// ParsePGN parses a PGN string and returns the tags and a list of move strings.
//
// It extracts the seven standard tag pairs from bracket-enclosed headers and
// parses the move text section, ignoring comments, variations, and NAGs.
func ParsePGN(pgn string) (PGNTags, []string, error) {
	tags := PGNTags{}
	lines := strings.Split(pgn, "\n")

	moveTextStart := 0
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			moveTextStart = i + 1
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if err := parseTag(&tags, line); err != nil {
				return PGNTags{}, nil, fmt.Errorf("failed to parse tag: %w", err)
			}
			moveTextStart = i + 1
		} else {
			// First non-tag, non-empty line starts the move text.
			moveTextStart = i
			break
		}
	}

	// Strip semicolon comments (rest-of-line) before joining.
	cleanedLines := make([]string, 0, len(lines)-moveTextStart)
	for _, l := range lines[moveTextStart:] {
		if idx := strings.Index(l, ";"); idx >= 0 {
			l = l[:idx]
		}
		cleanedLines = append(cleanedLines, l)
	}
	moveText := strings.Join(cleanedLines, " ")
	moves := parseMoveText(moveText)

	return tags, moves, nil
}

// determineResult returns the game result string.
func (c *Chess) determineResult(provided string) string {
	if provided != "" {
		return provided
	}

	if c.checkmate {
		if c.turn == gochess.White {
			return ResultBlackWins
		}
		return ResultWhiteWins
	}

	if c.stalemate {
		return ResultDraw
	}

	return ResultOngoing
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

// parseTag parses a single PGN tag line and sets it on the tags struct.
func parseTag(tags *PGNTags, line string) error {
	// Format: [Name "Value"]
	line = strings.TrimPrefix(line, "[")
	line = strings.TrimSuffix(line, "]")

	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid tag format: %s", line)
	}

	name := parts[0]
	raw := parts[1]
	// Strip surrounding quotes.
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	// Unescape per PGN specification.
	value := strings.ReplaceAll(raw, `\"`, `"`)
	value = strings.ReplaceAll(value, `\\`, `\`)

	switch name {
	case "Event":
		tags.Event = value
	case "Site":
		tags.Site = value
	case "Date":
		tags.Date = value
	case "Round":
		tags.Round = value
	case "White":
		tags.White = value
	case "Black":
		tags.Black = value
	case "Result":
		tags.Result = value
	}

	return nil
}

// parseMoveText parses the move text portion of a PGN string.
// It ignores comments (enclosed in {} or starting with ;),
// variations (enclosed in ()), and NAGs (starting with $).
func parseMoveText(text string) []string {
	var moves []string

	// Remove brace comments.
	for {
		start := strings.Index(text, "{")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	// Remove variations (parenthesized).
	for {
		start := strings.Index(text, "(")
		if start == -1 {
			break
		}
		depth := 1
		end := start + 1
		for end < len(text) && depth > 0 {
			if text[end] == '(' {
				depth++
			} else if text[end] == ')' {
				depth--
			}
			end++
		}
		text = text[:start] + text[end:]
	}

	tokens := strings.Fields(text)
	results := map[string]bool{
		ResultWhiteWins: true,
		ResultBlackWins: true,
		ResultDraw:      true,
		ResultOngoing:   true,
	}

	for _, token := range tokens {
		// Skip move numbers (e.g., "1.", "12.").
		if strings.HasSuffix(token, ".") {
			continue
		}
		// Also skip move numbers like "1..." for black moves.
		if strings.Contains(token, "...") {
			continue
		}
		// Skip NAGs.
		if strings.HasPrefix(token, "$") {
			continue
		}
		// Skip semicolon comments (rest of line, but we joined lines).
		if strings.HasPrefix(token, ";") {
			continue
		}
		// Skip result tokens.
		if results[token] {
			continue
		}
		// Skip empty tokens.
		if token == "" {
			continue
		}
		moves = append(moves, token)
	}

	return moves
}
