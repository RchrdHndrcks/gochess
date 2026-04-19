# [Feature Title]

**Spec ID**: `NNN-short-name`
**Created**: YYYY-MM-DD
**Status**: Draft | Ready | Implemented
**Branch**: `specs/NNN-short-name`

## Context

Why this feature exists. What problem it solves. What use case motivates it
(new variant support, API ergonomics, performance, correctness fix, etc.).

## Requirements

- [ ] Functional requirement 1
- [ ] Functional requirement 2

## Affected Areas

- `path/to/file.go` — what changes and why
- `chess/` — sub-package changes, if any

## Public API

### Added

```go
// Short docstring summary.
func NewThing(opts ...Option) *Thing
```

### Changed

```go
// Before:
func (c *Chess) DoX(arg int) error
// After:
func (c *Chess) DoX(arg int, opts ...Option) error
// Rationale: ...
```

### Removed

(none) — or list deprecated/removed exports with migration hint.

## Backward Compatibility

- Compatible with v2.x: **yes** / **no**
- If **no**: describe the migration for downstream users (`enrok-engine`, `enrok`)
  and which CHANGELOG heading this belongs under (`[2.1.0]`, `[3.0.0]`, etc.).

## Acceptance Criteria

1. **Given** a starting position, **When** `Method(args)` is called, **Then** ...
2. **Given** an edge case (e.g. insufficient material), **When** ..., **Then** ...

## Test Commands

```bash
make test          # go test -v ./... (filters FAILs)
make lint          # golangci-lint
make vet           # go vet
go test -race ./...
```

New tests live alongside the code they test (`board_test.go`, `chess/*_test.go`).

## Benchmarks

Required when the change touches a hot path (move generation, `AvailableMoves`,
FEN/PGN/SAN parsing, repetition detection).

```bash
go test -bench=. -benchmem ./...
```

Baseline expectation: no regression > 5% on existing benchmarks. If the feature
adds a new hot path, record baseline numbers here so reviewers can compare.

## Notes

- CHANGELOG entry draft (Added / Changed / Fixed / Breaking Changes section).
- Edge cases the implementation must handle.
- Design decisions worth preserving for future readers.
