## GOCHESS

Gochess implements all logic related to a Chess game. With this library you would be able
to have a Chess game or create a variant of it!

### Implementation

This library tries to separate the logic related to a Chess Board of the game.
Board has all logic related to pieces movements while Chess wraps logic related to chess rules
like which color must play in each turn, if castle is available or which square is available
to an in passant capture.

### Pieces

Pieces are implemented with bit logic.
These were created following the next table:

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

With this implementation you can easily create a colored piece only making an OR operation,
e.g.: White | Pawn -> `01001`

And with only five bits you can implement all Chess pieces, so, `int8` type is a really powerful
gun!

This kind of implementation was inspired by 
[Sebastian Lague](https://youtu.be/U4ogK0MIzqk?si=YFuh4EK5j4v6dZy5) who made the same in C#.

### Board

Board exports the following functions:

``` go
LoadPosition(string) error
Square(c Coordinate) (int8, error)
AvailableMoves(turn int8, inPassantSquare, castlePossibilities string) ([]string, error)
MakeMove(string) error
Width() int
```

`LoadPosition` load the pieces into the board. This board was created to be used in Chess package, 
so this function expects a FEN as parameter. It returns an error if FEN is not valid.

`Square` returns a piece of the square passed as coordinate. It returns an error if the coordinate
is out of bounds.

`AvailableMoves` returns all possible moves in UCI format. This function does not check if a move
is legal or not. This functionality is provided by Chess.

`MakeMove` makes the movements passed as UCI format. It can determinate if the move was a capture,
a castle or a promotion and it acts in consecuence. If the movement is not valid, it returns
an error.

`Width` returns the square count of a board side (8). This is only implemented to follow Chess
board interface (see more in Chess package).

