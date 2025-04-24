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

### Modular Design for Chess Variants

The separation between the board representation (`Board`) and the game rules (`Chess`) is a key design decision that enables easy creation of chess variants. By implementing the appropriate interfaces, you can:

- Create custom board sizes and shapes
- Implement alternative piece movement rules
- Design new chess variants with minimal code changes
- Reuse the core game logic while replacing specific components

This approach allows you to focus only on the aspects that differ in your variant, while leveraging the existing implementation for everything else.

### Pieces Implementation

Pieces are implemented using bitwise logic for efficient representation and operations. They follow this bit pattern:

```
Empty:  00000
Pawn:   00001
Knight: 00010
Bishop: 00011
Rook:   00100
Queen:  00101
King:   00110

White:  01000
Black:  10000
```

With this implementation, you can easily create a colored piece using a bitwise OR operation, e.g., `White | Pawn -> 01001`.

Using only five bits, you can represent all Chess pieces, making the `int8` type a powerful and memory-efficient tool for this purpose.

This bit-based representation offers several advantages:
- Compact storage (only 5 bits needed per piece)
- Fast bitwise operations for piece manipulation
- Easy extraction of piece type and color
- Efficient board state representation

This implementation was inspired by Sebastian Lague's chess engine design in C#.

### Board Implementation

The `Board` struct is responsible for the physical representation of the chess board and basic piece movements, without enforcing chess-specific rules. This separation allows for implementing chess variants by simply modifying or replacing the rules layer.

The `Board` exports the following functions:

```go
Width() int
Square(c Coordinate) (int8, error)
MakeMove(origin, target Coordinate) error
SetSquare(c Coordinate, p int8) error
```

### Coordinates

The `Coordinate` struct represents a position on the board with X and Y values. The package provides utility functions to work with coordinates, including conversion between different notation systems.

### Creating Chess Variants

To create a chess variant:

1. Implement a custom board by either extending the existing `Board` or creating a new one that satisfies the board interface
2. Modify or extend the game rules in a custom implementation that uses your board
3. Reuse as much of the existing code as makes sense for your variant

## Usage Examples

### Creating a New Chess Game

```go
package main

import (
	"fmt"
	"github.com/RchrdHndrcks/gochess/chess"
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
	"github.com/RchrdHndrcks/gochess/chess"
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