# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2026-04-04

### Breaking Changes

- [#28](https://github.com/RchrdHndrcks/gochess/pull/28) The `Board` interface now uses `gochess.Piece` instead of `int8` in `Square` and `SetSquare`. Custom board implementations must be updated accordingly.

### Added

- [#31](https://github.com/RchrdHndrcks/gochess/pull/31) `IsFiftyMoveRule() bool` — detects the fifty-move rule draw condition.
- [#32](https://github.com/RchrdHndrcks/gochess/pull/32) `IsInsufficientMaterial() bool` — detects insufficient mating material per FIDE 5.2.2 (K vs K, K+N vs K, K+B vs K, K+B vs K+B same-color squares).
- [#33](https://github.com/RchrdHndrcks/gochess/pull/33) `IsThreefoldRepetition() bool` — detects threefold repetition.
- [#35](https://github.com/RchrdHndrcks/gochess/pull/35) `PieceType(p Piece) Piece` and `PieceColor(p Piece) Piece` helpers in the root package.
- [#37](https://github.com/RchrdHndrcks/gochess/pull/37) PGN support: `(c *Chess) PGN(tags pgn.PGNTags) string` generates a PGN string; `chess/pgn` sub-package exposes `pgn.Parse()`, `pgn.PGNTags`, and result constants (`ResultWhiteWins`, `ResultBlackWins`, `ResultDraw`, `ResultOngoing`).
- [#38](https://github.com/RchrdHndrcks/gochess/pull/38) SAN support: `(c *Chess) SAN(uciMove string) (string, error)` converts UCI to Standard Algebraic Notation; `(c *Chess) FromSAN(san string) (string, error)` converts SAN back to UCI.

### Fixed

- [#24](https://github.com/RchrdHndrcks/gochess/pull/24) `Chess.Square()` now propagates board errors instead of silently returning empty.
- [#25](https://github.com/RchrdHndrcks/gochess/pull/25) Guard against panic in `updateHalfMoves` when history is empty.
- [#27](https://github.com/RchrdHndrcks/gochess/pull/27) Rename `promotionPosibilities` → `promotionPossibilities` (typo fix).
- [#39](https://github.com/RchrdHndrcks/gochess/pull/39) Castling is now correctly disallowed when the king is in check, passes through check, or lands in check (covers pawn attacks).
- [#42](https://github.com/RchrdHndrcks/gochess/pull/42) `AvailableMoves()` returns a copy of the slice; callers can no longer mutate internal state. Fixed stale doc comment that incorrectly described a panic.

### Changed

- [#29](https://github.com/RchrdHndrcks/gochess/pull/29) Map lookups in the board adapter are validated; invalid keys return an error instead of a silent zero value.
- [#36](https://github.com/RchrdHndrcks/gochess/pull/36) Renamed `move.go` to `notation.go` for clarity.

### Documentation

- [#30](https://github.com/RchrdHndrcks/gochess/pull/30) Documented that `Chess` is not safe for concurrent use.
- [#34](https://github.com/RchrdHndrcks/gochess/pull/34) Added negative and edge-case tests for `MakeMove` and `UnmakeMove`.

## [1.3.0] - 2026-01-09

### Added

- [#23](https://github.com/RchrdHndrcks/gochess/pull/23) Add Turn method to Chess.

## [1.2.0] - 2025-05-09

### Changed

- [#19](https://github.com/RchrdHndrcks/gochess/pull/19) Improve API for Chess and Board.
- [#16](https://github.com/RchrdHndrcks/gochess/pull/16) Adjust API and efficiency.

## [1.1.0] - 2025-04-23

### Changed

- [#12](https://github.com/RchrdHndrcks/gochess/pull/12) Change entire API.

### Added

- [#13](https://github.com/RchrdHndrcks/gochess/pull/13) Add MIT license.

## [1.0.4] - 2024-07-15

### Fixed

- [#10](https://github.com/RchrdHndrcks/gochess/pull/10) Fix promotion moves in capture.

## [1.0.3] - 2024-07-15

## Added

- [#8](https://github.com/RchrdHndrcks/gochess/pull/8) Add go.yml file for Github PRs.

## Fixed

- [#7](https://github.com/RchrdHndrcks/gochess/pull/7) Fix castling calculation.

## [1.0.2] - 2024-07-10

### Added

- [#5](https://github.com/RchrdHndrcks/gochess/pull/5) Add Square function to Chess.

## [1.0.1] - 2024-07-07

### Added

- [#3](https://github.com/RchrdHndrcks/gochess/pull/3) Add IsCheck function to Board.