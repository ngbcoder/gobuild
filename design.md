# System Design Document: Gobuild

**Project Name:** Gobuild  
**Author:** Nolan Bragan (@Ngbcoder)

**Status:** Draft / In Review  
**Date:** October 30, 2026

---
## 0.5 Main Insturctions
1. Document as much as you can.
2. Commit and push frequently.
3. Make a AGENTS.MD file as you go.
4. If you dont know ASK.

## 1. Overview & Problem Statement

### 1.1 Context
Developers working in compiled or heavily tested languages (Go, Rust, C++, TypeScript) frequently switch between their code editor/terminal and build/test commands. Existing workflows often rely on standard terminal outputs that clutter scrollback buffers or require manually re-running `go test` / `cargo check` commands.

### 1.2 Objective
**Gobuild** is a terminal user interface (TUI) application written in Go that provides a centralized, interactive workspace to monitor, trigger, and inspect project builds, unit tests, and custom task pipelines in real time.

---

## 2. Key Requirements & Features

### 2.1 Functional Requirements
* **Project & Task Discovery:** Automatically detect project configuration files (`go.mod`, `Cargo.toml`, `Makefile`, `package.json`) or load a user-defined `gobuild.yaml`.
* **Task Execution & Control:** Run, restart, and cancel build/test commands asynchronously using keyboard shortcuts.
* **Live Log Streaming:** Stream real-time standard output (`stdout`) and standard error (`stderr`) logs in a scrollable view pane.
* **Status Dashboard:** Display clear pass/fail status, duration, and execution history across configured projects and sub-tasks.
* **Configurability:** Support custom keybindings, custom task triggers, and environment variable overrides per project.
* **Plugin System:** Allow community-contributed language support via Git repository links. Users can add, update, or remove plugins dynamically.

### 2.2 Non-Functional Requirements
* **Low Latency:** Responsive UI updates (60fps rendering without blocking the background build/test processes).
* **Minimal Resource Footprint:** Low CPU and memory overhead during long-running background tasks.
* **Cross-Platform Support:** Native functionality across macOS, Linux, and Windows terminals.

---

## 3. Architecture & Tech Stack

### 3.1 Tech Stack
* **Language:** Go (1.22+)
* **TUI Framework:** [Charm Ecosystem](https://charm.sh/)
    * `bubbletea` (Elm-architecture UI state management)
    * `lipgloss` (Styling engine, borders, colors)
    * `bubbles` (UI components: list, viewport, spinner, text input)
* **Configuration:** `yaml` parsing via `gopkg.in/yaml.v3`
* **Process Execution:** Go `os/exec` with pipe readers (`io.Pipe`) for concurrent log streaming and `git` command execution for plugin installation.

### 3.2 High-Level Architecture

```
+-------------------------------------------------------------------+
|                            Gobuild TUI                            |
|                                                                   |
|  +---------------------------+   +-----------------------------+  |
|  |     Project List View     |   |      Task Output / Logs    |  |
|  | (Sidebar / Main Menu)     |   |     (Viewport Component)    |  |
|  +---------------------------+   +-----------------------------+  |
|  |              Status Bar / Keybinding Hints                  |  |
+-------------------------------------------------------------------+
                                  |
                                  v
                        [ Bubble Tea Event Loop ]
                                  |
            +---------------------+---------------------+
            |                     |                     |
            v                     v                     v
  [ Task Controller ]     [ Log Streamer ]     [ Plugin Manager ]
            |                     |                     |
            +---------------------+---------------------+
                                  |
                                  v
                       [ OS Process Execution ]
             (go test, make, git clone, cargo, etc.)
```

---

## 4. Component Design

### 4.1 State Management (Bubble Tea Model)
The application state follows the Model-Update-View pattern:

* **Model (`model.go`):** Stores current view mode (Dashboard, Log View, Config View, Plugin View), list of projects, active running process context (`exec.Cmd`), loaded plugin definitions, and log buffers.
* **Update (`update.go`):** Receives keyboard events (`tea.KeyMsg`), system signals, and custom message types (`TaskFinishedMsg`, `LogChunkMsg`, `PluginInstalledMsg`).
* **View (`view.go`):** Renders the UI layout using `lipgloss` styled blocks based on the current dimensions (`tea.WindowSizeMsg`).

### 4.2 Process Execution & Streaming Engine
1. **Command Spawn:** When a task is triggered, the engine initializes an `exec.CommandContext` with cancelation support.
2. **Asynchronous Piping:** `StdoutPipe` and `StderrPipe` read raw stdout/stderr streams into thread-safe buffers.
3. **Message Dispatching:** Stream chunks produce `tea.Cmd` updates that append lines to the viewport log without locking the UI thread.

### 4.3 Plugin Engine Architecture
The Plugin Engine allows extending language support without recompiling Gobuild:

1. **Storage Location:** Installed plugins are stored in `~/.gobuild/plugins/<repo-name>/`.
2. **Installation via Git:** When a user passes a repository URL, Gobuild executes a Shallow Git Clone (`git clone --depth 1 <url>`) into the plugins directory.
3. **Manifest Format (`plugin.yaml`):**
   ```yaml
   name: "rust-support"
   language: "rust"
   file_extensions: [".rs"]
   detect_files: ["Cargo.toml"]
   commands:
     build: "cargo build --release"
     test: "cargo test"
     clean: "cargo clean"
   ```
4. **Plugin Lifecycle:**
    * **Discovery & Load:** On application startup, Gobuild parses all manifest files in `~/.gobuild/plugins/`.
    * **Match & Resolve:** During project scanning, if project files match a plugin's `detect_files` or `file_extensions`, the plugin's build and test commands are bound to that project.
    * **Management Commands:**
        * `gobuild plugin add <git-repo-url>`
        * `gobuild plugin list`
        * `gobuild plugin update` (runs `git pull` across installed repos)
        * `gobuild plugin remove <plugin-name>`

---

## 5. User Interface & Workflow

```
+--------------------------------------------------------------------+
| Gobuild v0.1.0                [Projects: 3] [Plugins: 2] [Run: 1]  |
+-------------------------------+------------------------------------+
| PROJECTS                      | LOG OUTPUT: cargo test             |
| > rust-service    [RUNNING...] |   Compiling rust-service v0.1.0    |
|   worker-daemon   [PASSED]     |    Finished test [unoptimized]    |
|   web-client      [FAILED]     |     Running unittests src/lib.rs  |
|                               | running 4 tests                    |
| PLUGINS LOADED                | test tests::test_auth ... ok       |
| * rust-support (Cargo.toml)   | test tests::test_db ... ok         |
| * node-support (package.json) |                                    |
+-------------------------------+------------------------------------+
| [b] Build  [t] Test  [p] Plugins  [r] Rerun  [c] Config  [q] Quit |
+--------------------------------------------------------------------+
```

### Keybindings
* `Tab` / `Shift+Tab`: Switch focus between panels.
* `b`: Run build command (native or plugin-provided) for selected project.
* `t`: Run test suite.
* `p`: Open Plugin Manager overlay.
* `k`: Kill active process.
* `q` / `Ctrl+C`: Exit application safely.

---

## 6. Configuration Schema (`gobuild.yaml`)

```yaml
version: "1"
plugins_dir: "~/.gobuild/plugins"
projects:
  - name: rust-service
    path: ./backend
    language: rust # Matches plugin.yaml language key
    env:
      RUST_LOG: debug
  - name: web-client
    path: ./frontend
    build_cmd: "npm run build"
    test_cmd: "npm test"
```

---

## 7. Development Milestones

```
Phase 1: Foundation (Weeks 1-2)
├── Core CLI runner & process execution engine
├── Basic Bubble Tea layout (Sidebar + Viewport)
└── YAML config parser (gobuild.yaml)

Phase 2: Live Execution & Control (Weeks 3-4)
├── Asynchronous stdout/stderr log streaming
├── Interactive process cancelation & re-runs
└── Status bar state handling (Passed/Failed/Running)

Phase 3: Plugin System (Weeks 5-6)
├── Plugin engine (`~/.gobuild/plugins` directory scanner)
├── Manifest parser (`plugin.yaml`)
├── Git repository manager (`add`, `update`, `list`, `remove`)
└── Dynamic build/test command resolution

Phase 4: Polishing & Ecosystem (Weeks 7-8)
├── File watcher auto-trigger (`fsnotify`)
├── Discord Webhook notifications for build statuses
└── Cross-platform verification & release packaging
```

---

## 8. Future Enhancements

* **Discord Webhook Notifications:** Send pings (`@everyone` or target roles) to Discord channels when long-running build/test pipelines succeed or fail.
* **File Watcher Mode:** Re-trigger build/test tasks automatically upon file edits using `fsnotify`.
* **Central Plugin Registry:** Maintain an official curated list of language plugin repositories searchable directly within the TUI.
* **Test Filter View:** Filter viewport log outputs by failing test cases or specific error tags.