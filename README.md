# Gobuild

Gobuild is a terminal UI application for discovering, orchestrating, and monitoring build and test tasks across a developer workspace. It is designed to bring project execution, live log output, and plugin-based support into a centralized, keyboard-driven interface.

## Current Status
This repository is in the initial project setup phase. The foundational project structure, documentation, and module initialization are being established before moving into the configuration and TUI engine work.

## Goals
- Discover Go, Rust, Node, and Make-based projects automatically.
- Run build and test tasks with live, streamed output.
- Provide a small, responsive Bubble Tea-based terminal UI.
- Support extensible language automation through Git-backed plugins.

## Project Structure
- `cmd/gobuild/` — Application entrypoint.
- `internal/config/` — YAML configuration and project parsing.
- `internal/runner/` — Task execution and process management.
- `internal/plugin/` — Plugin loading and Git-based install/update support.
- `internal/ui/` — Bubble Tea models, views, and events.
- `pkg/` — Shared helpers and reusable packages.

## Quick Start

### Prerequisites
- Go 1.22 or newer
- Git

### Initialize the module
```bash
go mod init github.com/Ngbcoder/gobuild
```

### Build the app
```bash
mkdir -p bin
go build -o ./bin/gobuild ./cmd/gobuild
```

### Run the app
```bash
go run ./cmd/gobuild
```

## Documentation
- `design.md` — High-level architecture and system design.
- `todo.md` — Development roadmap and task tracking.
- `AGENTS.md` — Operating guidance for contributors and future work.

## License
This project is licensed under the MIT License.
