# chess/pgn

## Overview

The `pgn` package provides types and parsing utilities for the [Portable Game
Notation (PGN)](https://www.chessclub.com/help/PGN-spec) standard.

It is a sub-package of `chess/` and has no dependency on the chess engine.
This means it can be imported independently to inspect or process PGN data
without instantiating a game.

The companion generation function (`Chess.PGN`) lives in the parent `chess/`
package because it reads the engine's internal move history.

## Key Types

### `PGNTags`

```go
type PGNTags struct {
    Event  string
    Site   string
    Date   string
    Round  string
    White  string
    Black  string
    Result string
}
```

Holds the seven required tag pairs defined by the PGN standard (the "Seven
Tag Roster"). Empty fields in a generated PGN default to `"?"`.

### Result constants

```go
const (
    ResultWhiteWins = "1-0"
    ResultBlackWins = "0-1"
    ResultDraw      = "1/2-1/2"
    ResultOngoing   = "*"
)
```

## API

```go
func Parse(pgn string) (PGNTags, []string, error)
```

### `Parse`

Parses a PGN string and returns:

1. The seven tag pairs extracted from the bracket-enclosed headers.
2. A slice of moves in the same notation used in the move text (UCI for games
   produced by this library).
3. An error if a tag line is malformed.

The parser ignores:
- Brace comments `{ ... }`
- Semicolon rest-of-line comments `; ...`
- Variations `( ... )` (including nested)
- NAGs (e.g. `$1`, `$18`)
- Move numbers and result tokens

## Usage examples

### Parse a PGN file

```go
import (
    chesspgn "github.com/RchrdHndrcks/gochess/chess/pgn"
)

pgnText := `[Event "World Championship"]
[Site "Reykjavik"]
[Date "1972.07.11"]
[Round "1"]
[White "Spassky"]
[Black "Fischer"]
[Result "1-0"]

1. d4 Nf6 2. c4 e6 3. Nf3 1-0
`

tags, moves, err := chesspgn.Parse(pgnText)
if err != nil {
    log.Fatal(err)
}

fmt.Println(tags.White)  // Spassky
fmt.Println(tags.Black)  // Fischer
fmt.Println(moves)       // [d4 Nf6 c4 e6 Nf3]
```

### Generate and round-trip a game

```go
import (
    "github.com/RchrdHndrcks/gochess/chess"
    chesspgn "github.com/RchrdHndrcks/gochess/chess/pgn"
)

game, _ := chess.New()
game.MakeMove("e2e4")
game.MakeMove("e7e5")
game.MakeMove("g1f3")

tags := chesspgn.PGNTags{
    Event: "My Game",
    White: "Alice",
    Black: "Bob",
}
pgnText := game.PGN(tags)

// Parse it back
parsedTags, parsedMoves, _ := chesspgn.Parse(pgnText)
fmt.Println(parsedTags.Event) // My Game
fmt.Println(parsedMoves)      // [e2e4 e7e5 g1f3]
```

### Check draw/win result constants

```go
if tags.Result == chesspgn.ResultDraw {
    fmt.Println("The game was drawn.")
}
```

## PGN escaping rules

Tag values are escaped per the PGN specification:

| Character | Escaped as |
|-----------|------------|
| `\`       | `\\`       |
| `"`       | `\"`       |
| `\n`, `\r` | stripped  |

`Parse` reverses the escaping when reading tag values back.

## Dependencies

- Standard library only (`fmt`, `strings`).

## Interactions with other packages

| Package | Relationship |
|---------|-------------|
| `chess/` | Imports `chess/pgn` for `PGNTags` and result constants. Provides `Chess.PGN()` for generation. |
| `gochess` (root) | No direct dependency. |
