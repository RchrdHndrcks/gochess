# Chess Package

## Overview

The `chess` package implements the complete set of standard chess rules on top of the core board representation. It handles all the chess-specific logic while delegating the basic board operations to the underlying board implementation.

This package is designed with flexibility in mind, allowing you to:

1. Use it as a complete chess rules engine
2. Extend it to create chess variants
3. Replace components to customize behavior

## Architecture

The `Chess` struct wraps a board implementation and adds all the chess-specific rules and logic:

- Turn management (white/black)
- Move legality validation
- Special move handling (castling, en passant, promotion)
- Check and checkmate detection
- Game state tracking (halfmove clock, fullmove counter)
- FEN notation support

By default, it uses the standard `pkg.Board` implementation, but it can work with any board that satisfies the `Board` interface, making it ideal for creating chess variants.

## API

The `Chess` package exports the following key functions:

```go
func New(options ...Option) (*Chess, error)
func (c *Chess) FEN() string
func (c *Chess) AvailableMoves() []string
func (c *Chess) MakeMove(move string) error
func (c *Chess) UnmakeMove()
func (c *Chess) IsCheck() bool
func (c *Chess) IsCheckmate() bool
func (c *Chess) IsStalemate() bool
func (c *Chess) LoadPosition(fen string) error
func (c *Chess) Clone() *Chess
```

### Core Functions

- `New(options ...Option)`: Creates a new chess game with the specified options. Without options, it creates a standard starting position.

- `FEN() string`: Returns the current position in Forsyth-Edwards Notation (FEN).

- `AvailableMoves() []string`: Returns all possible legal moves in UCI format.

- `MakeMove(move string) error`: Validates and executes a move in UCI format (e.g., "e2e4"). Returns an error if the move is illegal.

- `UnmakeMove()`: Reverts the last move made, restoring the previous position.

- `IsCheck() bool`: Returns whether the current player's king is in check. If the position is checkmate or stalemate, it returns false.

- `IsCheckmate() bool`: Returns whether the current player's king is in checkmate.

- `IsStalemate() bool`: Returns whether the game is in stalemate.

- `LoadPosition(fen string) error`: Sets up the board according to the provided FEN string.

- `Clone() *Chess`: Returns a copy of the chess game.

## Creating a Chess Game

### Basic Usage

```go
// Create a new chess game with the standard starting position
game, err := chess.New()
if err != nil {
    // Handle error
}

// Get all legal moves for the current player
moves := game.AvailableMoves()

// Make a move
err = game.MakeMove("e2e4")
if err != nil {
    // Handle error
}
```

### Custom Starting Position

```go
// Create a game from a specific FEN position
game, err := chess.New(chess.WithFEN("r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3"))
```

## Options

The `New` function accepts options to customize the chess game:

- `WithBoard(board Board)`: Uses a custom board implementation. This should be the first option if used.

- `WithFEN(fen string)`: Sets up the board using the provided FEN string.

- `WithParallelism(parallelism int)`: Sets the number of parallel workers to use for move generation. The default is twice the number of CPU cores.

## Board Interface

Any board implementation used with the Chess package must satisfy this interface:

```go
type Board interface {
    LoadPosition(string) error
    Square(c pkg.Coordinate) (int8, error)
    AvailableMoves(turn int8, enPassantSquare, castlePossibilities string) ([]string, error)
    Width() int
}
```

## Creating Chess Variants

To create a chess variant:

1. Implement a custom board that satisfies the Board interface
2. Pass it to the Chess constructor using the WithBoard option
3. Extend the Chess struct if needed to add variant-specific rules

Example for a custom variant:

```go
// Create a custom board implementation
customBoard := myvariant.NewCustomBoard()

// Create a chess game using the custom board
game, err := chess.New(chess.WithBoard(customBoard))
```

This modular approach allows you to focus only on the aspects that differ in your variant while leveraging the existing chess logic for everything else.

## Performance Optimization

### Parallel Move Calculation

GoChess uses a parallel processing approach to calculate legal moves, which significantly improves performance on multi-core systems:

- **Automatic Parallelism**: By default, GoChess uses twice as many workers as available CPU cores to maximize performance.
- **Efficient Cloning**: The system implements a `Cloner` interface that allows creating independent copies of the board for each worker, avoiding race conditions.
- **Customizable Configuration**: The level of parallelism can be adjusted using the `WithParallelism` option, allowing optimization based on the execution environment.

```go
// Create a game with a specific level of parallelism
game, err := chess.New(chess.WithParallelism(4))
```

### Memory Optimizations

- **Capacity Preallocation**: Slices for storing moves are preallocated with specific capacities based on the piece type.
- **Efficient Representation**: Pieces are represented using bits, allowing fast operations and minimal memory usage.

### Environment Adaptation

- **Containerized Environments**: In environments with limited resources, it is recommended to manually adjust the level of parallelism. If you are using a containerized environment, you should set the parallelism manually to avoid issues with the number of available CPU cores.
- **Sequential Mode**: For applications that do not require high performance, it can be configured to perform calculations sequentially.