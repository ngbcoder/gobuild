# Gobuild Development TODO \& Roadmap

> **Author:** Nolan Bragan (@Ngbcoder)  
> **Status:** Active Tracking  
> **Repository Rules:** 
> - Document as much as you can.  
> - Commit and push frequently.  
> - Maintain `AGENTS.md` continuously.  
> - If you don't know, ASK.


## Phase 0: Project Setup \& Repository Initialization

### 0.1 Repository Setup

* \[x] Initialize Git repository (`git init`).
* \[x] Create `.gitignore` for Go (include binary outputs, test binaries, `.gobuild/` test directories, vendor files).
* \[x] Initialize Go module (`go mod init github.com/Ngbcoder/gobuild`).
* \[x] Create root `AGENTS.md` file with initial project context and contributor instructions.
* \[x] Create initial `README.md` containing project summary and build instructions.
* \[x] Create initial license file (e.g., MIT).

### 0.2 Initial Documentation \& Structural Setup

* \[x] Create project folder structure:

  * `cmd/gobuild/` (Application entrypoint)
  * `internal/config/` (Configuration models \& parsers)
  * `internal/runner/` (Process execution \& stream handlers)
  * `internal/plugin/` (Plugin loading \& Git manager)
  * `internal/ui/` (Bubble Tea components, views, models)
  * `pkg/` (Shared utilities if any)
* \[x] Document project architecture and folder guidelines in `AGENTS.md`.

\---

## Phase 1: Core Foundation \& Configuration Parsing (Weeks 1-2)

### 1.1 YAML Configuration Model (`gobuild.yaml`)

* \[x] Install `gopkg.in/yaml.v3` dependency (`go get gopkg.in/yaml.v3`).
* \[x] Define `ProjectConfig` and `GobuildConfig` Go structs in `internal/config/config.go`.
* \[x] Implement `LoadConfig(path string) (\*GobuildConfig, error)`.
* \[x] Handle missing configuration file gracefully with default fallbacks.
* \[x] Expand path tilde (`\~`) logic for cross-platform support (e.g., `\~/.gobuild/plugins`).
* \[x] Write unit tests for configuration parsing (`internal/config/config\_test.go`).

### 1.2 Native Auto-Discovery Engine

* \[x] Implement detection logic for `go.mod` (Go projects).
* \[x] Implement detection logic for `Cargo.toml` (Rust projects).
* \[x] Implement detection logic for `package.json` (Node/TypeScript projects).
* \[x] Implement detection logic for `Makefile`.
* \[x] Create project scanner utility to traverse current working directory and find sub-projects.
* \[x] Unit test discovery logic against mock directory structures.

### 1.3 Basic TUI Framework Setup (Charm Ecosystem)

* \[ ] Install Charm dependencies (`bubbletea`, `lipgloss`, `bubbles`).
* \[ ] Define root application state model in `internal/ui/model.go`.
* \[ ] Create basic `Init()`, `Update()`, and `View()` functions in `internal/ui/app.go`.
* \[ ] Implement main window resize handler (`tea.WindowSizeMsg`).
* \[ ] Setup initial layout structure using `lipgloss`:

  * \[ ] Top header bar
  * \[ ] Left sidebar pane (Projects \& Plugins)
  * \[ ] Right main view pane (Log Output / Viewport)
  * \[ ] Bottom status bar / help panel
* \[ ] Verify core CLI entrypoint (`cmd/gobuild/main.go`) initializes the Bubble Tea program cleanly.

\---

## Phase 2: Execution Engine \& Real-Time Log Streaming (Weeks 3-4)

### 2.1 Low-Level Process Execution Unit

* \[ ] Create `TaskRunner` struct in `internal/runner/runner.go`.
* \[ ] Implement `exec.CommandContext` wrapping to support explicit cancelation via context cancellation.
* \[ ] Create thread-safe log buffer for task stdout/stderr stream ingestion.
* \[ ] Handle process exit status codes and execution duration metrics logging.
* \[ ] Write unit tests for running shell commands and capturing outputs asynchronously.

### 2.2 Asynchronous Log Streaming to TUI

* \[ ] Implement pipe readers (`io.Pipe`, `bufio.Scanner`) for non-blocking line reading.
* \[ ] Define Bubble Tea custom message types (`LogChunkMsg`, `TaskFinishedMsg`, `TaskStartedMsg`).
* \[ ] Implement dynamic line-by-line append handler in Bubble Tea `Update()` method.
* \[ ] Connect `bubbles/viewport` component to display active command output.
* \[ ] Implement auto-scroll to bottom functionality in viewport when new log lines arrive.
* \[ ] Add manual scrolling controls for the viewport (Up/Down arrow keys, PageUp/PageDown).

### 2.3 Interactive Process Control \& Keybindings

* \[ ] Implement `b` keybinding: Trigger project build command.
* \[ ] Implement `t` keybinding: Trigger project test command.
* \[ ] Implement `k` keybinding: Terminate active task process context (`SIGKILL` / `SIGINT`).
* \[ ] Implement `r` keybinding: Cancel current process (if running) and rerun instantly.
* \[ ] Implement `Tab` / `Shift+Tab` keybindings: Toggle focus between Sidebar and Output Viewport.
* \[ ] Implement `q` / `Ctrl+C` keybindings: Gracefully shut down active child processes and exit app.

### 2.4 Status Bar \& UI Polish

* \[ ] Implement task execution status indicators: `\[IDLE]`, `\[RUNNING...]`, `\[PASSED]`, `\[FAILED]`.
* \[ ] Add execution duration display (e.g., `1.2s`, `450ms`) upon task completion.
* \[ ] Format pass/fail output indicators with distinct colors using `lipgloss` (Green for Pass, Red for Fail).
* \[ ] Create help message panel displaying available keybindings dynamically based on current focus.

\---

## Phase 3: Dynamic Plugin System (Weeks 5-6)

### 3.1 Plugin Manifest \& Directory Management

* \[ ] Define `PluginManifest` Go struct for parsing `plugin.yaml`.
* \[ ] Implement filesystem scanner to load manifests from `\~/.gobuild/plugins/`.
* \[ ] Implement plugin auto-creation of standard plugins directory (`\~/.gobuild/plugins`) on first run.
* \[ ] Add unit tests for manifest parsing and plugin validation.

### 3.2 Git-Based Plugin Manager

* \[ ] Create `internal/plugin/git.go` module for shell execution of Git commands.
* \[ ] Implement `AddPlugin(repoURL string)` using `git clone --depth 1 <url>`.
* \[ ] Implement `ListPlugins()` returning installed plugin instances and manifests.
* \[ ] Implement `UpdatePlugins()` running `git pull` across all directories in `\~/.gobuild/plugins/`.
* \[ ] Implement `RemovePlugin(pluginName string)` removing the plugin folder from disk.

### 3.3 Dynamic Command Resolution Engine

* \[ ] Implement file and project matching logic based on plugin `detect\_files` and `file\_extensions`.
* \[ ] Resolve task execution fallback order:

  1. Explicit project override in `gobuild.yaml`.
  2. Discovered language plugin command (`plugin.yaml`).
  3. Default system discovery rules.
* \[ ] Test command resolution logic against various simulated project layouts.

### 3.4 Plugin Manager TUI Overlay

* \[ ] Implement `p` keybinding to trigger Plugin Manager view overlay.
* \[ ] Display list of installed plugins with repository links and supported file extensions.
* \[ ] Create input dialog modal for adding a new plugin via Git URL (`bubbles/textinput`).
* \[ ] Add keybindings to trigger update or removal of selected plugins from UI.

\---

## Phase 4: Advanced Features, Polishing \& Release (Weeks 7-8)

### 4.1 File Watcher Mode (`fsnotify`)

* \[ ] Install `github.com/fsnotify/fsnotify` dependency.
* \[ ] Create background file watcher service in `internal/watcher/watcher.go`.
* \[ ] Ignore standard excluded patterns (`.git`, `node\_modules`, `vendor`, build outputs).
* \[ ] Trigger task execution commands automatically upon source code modification events.
* \[ ] Add toggle shortcut key in TUI to enable/disable file watcher on the fly.

### 4.2 Discord Webhook Notifications

* \[ ] Create webhook payload builder in `internal/notify/discord.go`.
* \[ ] Add `discord\_webhook\_url` setting to `gobuild.yaml`.
* \[ ] Trigger async HTTP POST requests to Discord on build/test completion (Pass/Fail alerts).
* \[ ] Embed execution duration, project name, and status summary in webhook payload.

### 4.3 Documentation \& AGENTS.MD Updates

* \[ ] Document complete plugin creation guide with sample `plugin.yaml`.
* \[ ] Create detailed troubleshooting guide in `README.md`.
* \[ ] Review and update `AGENTS.md` with final project structure, guidelines, and API specs.

### 4.4 Cross-Platform Testing \& Release Packaging

* \[ ] Test execution, terminal rendering, and keybindings on macOS.
* \[ ] Test execution, terminal rendering, and keybindings on Linux (Ubuntu/Debian).
* \[ ] Test execution, terminal rendering, and keybindings on Windows (CMD, PowerShell, Windows Terminal).
* \[ ] Set up GoReleaser or GitHub Actions pipeline for automated binary generation (`darwin`, `linux`, `windows`).
* \[ ] Tag initial official release (v0.1.0).

```