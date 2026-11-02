# AGENTS.md

## Repository Overview
Gobuild is a Go-based terminal UI for discovering, running, and monitoring project build/test tasks in a single workspace. The project is intended to support native project detection for Go, Rust, Node, and Make-based workflows while allowing extension through Git-backed plugins.

## Core Principles
- Keep documentation in sync with implementation changes.
- Prefer small, testable units and clear package boundaries.
- Use Go idioms and standard library patterns before introducing heavier dependencies.
- Document decisions and tradeoffs in the design docs and README as the code evolves.
- If a task is unclear, ask before making a significant design choice.

## Repository Layout
- `cmd/gobuild/`: CLI entrypoint for the application.
- `internal/config/`: Load and validate project configuration.
- `internal/runner/`: Process execution, task lifecycle, and log capture.
- `internal/plugin/`: Plugin discovery, manifest parsing, and Git installation helpers.
- `internal/ui/`: Bubble Tea components and view models.
- `pkg/`: Shared utilities that may be reused across packages.

## Development Guidelines
1. Maintain a clear separation between app logic, configuration, and UI concerns.
2. Add unit tests for parsing and execution logic where possible.
3. Favor async, cancellable process execution patterns for long-running tasks.
4. Keep the README and design notes aligned with the implementation.
5. Avoid introducing breaking changes without updating documentation and tests.

## Useful Commands
- `go mod init github.com/Ngbcoder/gobuild`
- `go test ./...`
- `go build ./cmd/gobuild`
- `go run ./cmd/gobuild`

## Notes
This repository is in the initial setup phase. The project structure and foundational configuration are intentionally minimal and will expand as Phase 1 begins.
