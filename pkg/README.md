# GoChess Package Architecture

This package contains the core components of the GoChess library. The architecture is deliberately modular to facilitate the creation of chess variants and custom implementations.

## Modular Design for Chess Variants

The separation between the board representation (`Board`) and the game rules (`Chess`) is a key design decision that enables easy creation of chess variants. By implementing the appropriate interfaces, you can:

- Create custom board sizes and shapes
- Implement alternative piece movement rules
- Design new chess variants with minimal code changes
- Reuse the core game logic while replacing specific components

This approach allows you to focus only on the aspects that differ in your variant, while leveraging the existing implementation for everything else.

## Pieces

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

## Board

The `Board` struct is responsible for the physical representation of the chess board and basic piece movements, without enforcing chess-specific rules. This separation allows for implementing chess variants by simply modifying or replacing the rules layer.

The `Board` exports the following functions:

```go
Width() int
Square(c Coordinate) (int8, error)
MakeMove(origin, target Coordinate) error
SetSquare(c Coordinate, p int8) error
isValidCoordinate(c Coordinate) bool
```

- `Width()`: Returns the number of squares on one side of the board (8). This is implemented to follow the Chess board interface.

- `Square(c Coordinate)`: Returns the piece at the given coordinate. It returns an error if the coordinate is out of bounds.

- `MakeMove(origin, target Coordinate)`: Executes a move from origin to target. It handles the basic piece movement without validating chess-specific rules. Returns an error if coordinates are invalid.

- `SetSquare(c Coordinate, p int8)`: Places a piece at the specified coordinate. Returns an error if the coordinate is invalid.

## Coordinates

The `Coordinate` struct represents a position on the board with X and Y values. The package provides utility functions to work with coordinates, including conversion between different notation systems.

## Creating Chess Variants

To create a chess variant:

1. Implement a custom board by either extending the existing `Board` or creating a new one that satisfies the board interface
2. Modify or extend the game rules in a custom implementation that uses your board
3. Reuse as much of the existing code as makes sense for your variant

See the `chess` package for more details on how the standard chess rules are implemented on top of this foundation.
