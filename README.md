## GOCHESS

Gochess implements all the logic related to a Chess game. With this library, you will be able to 
play a Chess game or create a variant of it!

### Implementation

This library separates the logic related to the Chess board from the game logic. The Board struct 
handles all logic related to piece movements, while the Chess struct wraps the logic related to 
chess rules, such as determining which color must play each turn, checking castling availability, 
and identifying squares available for en passant capture.