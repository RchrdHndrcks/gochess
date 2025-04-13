# GoChess

[![Go Reference](https://pkg.go.dev/badge/github.com/RchrdHndrcks/gochess.svg)](https://pkg.go.dev/github.com/RchrdHndrcks/gochess)
[![Go Report Card](https://goreportcard.com/badge/github.com/RchrdHndrcks/gochess)](https://goreportcard.com/report/github.com/RchrdHndrcks/gochess)

GoChess is a robust Go library that implements the complete logic for chess games. With this library, you can easily build chess applications, create game variants, or integrate chess functionality into your projects.

## Features

- Complete chess rules implementation
- Move validation and generation
- FEN notation parsing and generation
- Simple and intuitive API
- Well-tested codebase
- Lightweight with no external dependencies

## Installation

```bash
go get github.com/RchrdHndrcks/gochess@latest
```

## Architecture

GoChess separates concerns by dividing the implementation into two main components:

1. **Board**: Handles the board representation and basic piece movements
   - Manages the 2D grid of pieces
   - Provides methods for moving pieces and querying board state
   - Coordinates validation and conversion

2. **Chess**: Implements chess-specific rules and game logic
   - Turn management (determining which color plays next)
   - Special move validation (castling, en passant, etc.)
   - Check and checkmate detection
   - Game state tracking
   - FEN notation support

This separation allows for greater flexibility and makes it easier to create chess variants by modifying only the necessary components.

## Usage Examples

### Creating a New Chess Game

```go
package main

import (
	"fmt"
	"github.com/RchrdHndrcks/gochess/pkg/chess"
)

func main() {
	// Create a new chess game with default starting position
	game, err := chess.New()
	if err != nil {
		panic(err)
	}
	
	// Print the current board state in FEN notation
	fmt.Println(game.FEN())
	// Output: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	
	// Get all available moves for the current player
	moves := game.AvailableMoves()
	fmt.Println("Available moves:", moves)
	
	// Make a move
	err = game.MakeMove("e2e4")
	if err != nil {
		panic(err)
	}
	
	// Print the updated board state
	fmt.Println(game.FEN())
}
```

### Creating a Game from a Custom Position

```go
package main

import (
	"fmt"
	"github.com/RchrdHndrcks/gochess/pkg/chess"
)

func main() {
	// Create a game from a specific FEN position
	fenPosition := "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3"
	game, err := chess.New(chess.WithFEN(fenPosition))
	if err != nil {
		panic(err)
	}
	
	// Check if the king is in check
	isCheck := game.IsCheck()
	fmt.Println("Is king in check?", isCheck)
}
```

## Documentation

For detailed documentation, please visit the [GoDoc page](https://pkg.go.dev/github.com/RchrdHndrcks/gochess).

## Testing

The library includes comprehensive tests for all components. To run the tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.