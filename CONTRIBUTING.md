# Contributing

Thanks for contributing to enumerators.

## Setup

1. Fork and clone the repository.
2. `go mod download`
3. `go test ./...`

## Pull requests

- Run `go fmt ./...` and `go vet ./...`.
- Ensure new pipeline helpers dispose upstream enumerators correctly.
- Add tests for edge cases (empty sources, errors, cancellation).
- Update `docs/operations.md` when adding combinators.

## Changelog

Note changes under `## [Unreleased]` in [CHANGELOG.md](CHANGELOG.md).
