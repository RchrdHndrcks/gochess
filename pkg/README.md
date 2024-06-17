## GOCHESS

Gochess implements all the logic related to a Chess game. With this library, you will be able to 
play a Chess game or create a variant of it!

### Implementation

This library separates the logic related to the Chess board from the game logic. The Board struct 
handles all logic related to piece movements, while the Chess struct wraps the logic related to 
chess rules, such as determining which color must play each turn, checking castling availability, 
and identifying squares available for en passant capture.

### Pieces

Pieces are implemented using bitwise logic. They are represented following this table:

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

With this implementation, you can easily create a colored piece using a bitwise OR operation, e.g., White | Pawn -> 01001.

Using only five bits, you can represent all Chess pieces, making the int8 type a powerful tool for this purpose.

This implementation was inspired by Sebastian Lague, who did something similar in C#.

### Board

The `Board` exports the following functions:

``` go
LoadPosition(string) error
Square(c Coordinate) (int8, error)
AvailableMoves(turn int8, inPassantSquare, castlePossibilities string) ([]string, error)
MakeMove(string) error
Width() int
```

`LoadPosition`: Loads the pieces onto the board using a FEN string as a parameter. It returns an error 
if the FEN string is invalid.

`Square`: Returns the piece at the given coordinate. It returns an error if the coordinate is out of 
bounds.

`AvailableMoves`: Returns all possible moves in UCI format. This function does not check if a move is 
legal; that functionality is provided by the Chess struct.

`MakeMove`: Executes a move in UCI format. It can determine if the move was a capture, castling, or 
promotion and acts accordingly. If the move is invalid, it returns an error.

`Width`: Returns the number of squares on one side of the board (8). This is implemented to follow the
Chess board interface (see more in the Chess package).

