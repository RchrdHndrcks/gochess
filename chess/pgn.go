package chess

import (
	"fmt"
	"strings"

	"github.com/RchrdHndrcks/gochess"
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

// ToPGN generates a PGN string from the current game state.
//
// It writes the seven required tag pairs and the move text using UCI notation.
// Empty tag values default to "?". The Result tag is determined automatically
// if not provided: "1-0" or "0-1" for checkmate, "1/2-1/2" for stalemate,
// and "*" for an ongoing game.
func (c *Chess) ToPGN(tags PGNTags) string {
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

	moveText := strings.Join(lines[moveTextStart:], " ")
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
			return "0-1"
		}
		return "1-0"
	}

	if c.stalemate {
		return "1/2-1/2"
	}

	return "*"
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
func writeTag(sb *strings.Builder, name, value string) {
	sb.WriteString(fmt.Sprintf("[%s \"%s\"]\n", name, value))
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
	value := strings.Trim(parts[1], "\"")

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
		"1-0":     true,
		"0-1":     true,
		"1/2-1/2": true,
		"*":       true,
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
