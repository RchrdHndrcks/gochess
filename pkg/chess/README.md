## Chess

Chess implements all logic related to game rules. It use as default the
pkg.Board but it could receive anotherone if it follows Board interface.

Chess mission is wrap all chess logic but also be the most generic as possible.
This could be usefull if you want to implement a chess variant, like 960Chess.

Chess exports the following functions:

``` go
func (c Chess) FEN() string
func (c Chess) AvailableLegalMoves() ([]string, error)
func (c *Chess) MakeMove(move string) error 
func (c *Chess) UnmakeMove() error
```

`FEN` returns a string with the FEN of the current position.

`AvailableLegalMoves` returns all legal moves of the current player. If position
is stalemate, it will return an empty slice and no error. If position is checkmate,
it will return a nil slice and no error.

`MakeMove` checks if move passed by parameter is legal and if it is makes it.

`Unmakemove` unmakes the last move.

### Options

Chess is exported, so it could be created manually. If you follow this approach
you are going to have a Chess with an empty board. 

Also, you can create a Chess using `NewChess` function. It receives Options
as parameter. If you don't use any, Chess board will be created in initial position.

Those are the posssible options:

`WithBoard` receives a Board and put it in Chess. If you want to use this, it 
should be the first option.

`WithFEN` receives a string and calls Board.LoadPosition. 

### Board

This is Board interface shape:

```go
Board interface {
		LoadPosition(string) error
		Square(c pkg.Coordinate) (int8, error)
		AvailableMoves(turn int8, inPassantSquare, castlePossibilities string) ([]string, error)
		MakeMove(string) error
		Width() int
	}
```