# DifLog

A Git-like version control tool for AI contexts — built for developers who use **ClaudeCode (opencode)** and **GitHub Copilot CLI**.

DifLog tracks, versions, commits, diffs, and restores your AI context files (like `CLAUDE.md` and `.github/copilot-instructions.md`) with a full Bubble Tea TUI experience.

---

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Build from Source](#build-from-source)
  - [Install Globally](#install-globally)
- [Quick Start](#quick-start)
- [Commands Reference](#commands-reference)
- [Project Architecture](#project-architecture)
  - [CLEAN Architecture Layers](#clean-architecture-layers)
  - [Directory Structure](#directory-structure)
  - [Domain Entities](#domain-entities)
  - [Data Flow](#data-flow)
- [.difLog Storage Format](#diflog-storage-format)
- [AI Tool Support](#ai-tool-support)
  - [ClaudeCode (opencode)](#claudecode-opencode)
  - [GitHub Copilot CLI](#github-copilot-cli)
- [Remote Repository (GitHub)](#remote-repository-github)
- [Skill System](#skill-system)
- [Environment Variables](#environment-variables)
- [Running Tests](#running-tests)
- [Dependencies](#dependencies)

---

## Features

| Feature | Description |
|---|---|
| `init` | Initialize `.difLog` in any project directory with interactive AI tool selection |
| `add` | Stage AI context files (auto-detects based on configured AI tool) |
| `commit` | Commit staged contexts with a message, SHA-256 hash, and timestamp |
| `log` | Browse full commit history in a scrollable TUI list |
| `diff` | View line-by-line diff between any two commits |
| `checkout` | Restore your project context files to any previous commit |
| `branch` | Create, list, and switch branches for isolated context versions |
| `push` | Push committed contexts to a remote GitHub repository |
| `pull` | Pull contexts from a remote GitHub repository |
| `detect` | Auto-detect which AI tool is used in the current project |
| `skill create` | Install a DifLog skill into opencode or Copilot CLI globally |
| `skill save` | Snapshot the current AI chat context with a SHA-256 hash |
| `status` | View currently staged contexts and active branch |

---

## Requirements

- **Go 1.24.3+**
- **macOS / Linux** (Windows support is untested)
- An interactive terminal (TTY) — required for the Bubble Tea TUI
- A GitHub Personal Access Token (for `push` / `pull` only) — see [Remote Repository](#remote-repository-github)

---

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/xedeveloper/DifLog.git
cd DifLog

# Install dependencies
go mod download

# Build the binary
go build -o diflog .
```

### Install Globally

To install `diflog` as a system-wide command:

```bash
go install github.com/xedeveloper/DifLog@latest
```

Or after building from source, move the binary to your PATH:

```bash
# Build
go build -o diflog .

# Move to a directory in your PATH (e.g. /usr/local/bin)
mv diflog /usr/local/bin/diflog

# Verify
diflog --help
```

---

## Quick Start

```bash
# 1. Navigate to your project
cd your-project/

# 2. Initialize DifLog (interactive AI tool selection TUI)
diflog init

# 3. Stage your AI context files
diflog add

# 4. Commit the staged contexts
diflog commit -m "initial context snapshot"

# 5. View the commit history
diflog log

# 6. Install the DifLog skill into your AI tool
diflog skill create

# 7. (Optional) Push to a remote GitHub repo for backup
export DIFLOG_GITHUB_TOKEN=your_token_here
diflog push https://github.com/your-org/your-ai-contexts-repo
```

---

## Commands Reference

### `diflog init`

Initialize DifLog in the current directory. Creates a `.difLog/` directory and launches an interactive AI tool selection TUI that auto-detects ClaudeCode or Copilot CLI.

```bash
diflog init

# Skip the interactive TUI and specify the AI tool directly
diflog init --ai claude
diflog init --ai copilot
```

**Flags:**

| Flag | Short | Description |
|---|---|---|
| `--ai` | `-a` | AI tool to use: `claude` or `copilot`. Skips interactive selection. |

---

### `diflog add`

Stage AI context files into the DifLog staging area. When called with no arguments, auto-discovers context files for the configured AI tool.

```bash
# Auto-detect and stage all context files for the configured AI tool
diflog add

# Stage a specific file
diflog add CLAUDE.md

# Stage multiple files
diflog add CLAUDE.md .github/copilot-instructions.md
```

**ClaudeCode files auto-staged:** `CLAUDE.md`, `.claude/settings.json`, `.claude/settings.local.json`, any `.md` files in `.claude/`

**Copilot CLI files auto-staged:** `.github/copilot-instructions.md`, any `.md` files in `.github/`

---

### `diflog commit`

Commit all staged contexts with a message. Each commit records a SHA-256 hash, timestamp, branch, parent commit hash, and author.

```bash
diflog commit -m "add project-specific coding rules"
```

**Flags:**

| Flag | Short | Description |
|---|---|---|
| `--message` | `-m` | Commit message (required) |

---

### `diflog log`

View commit history in a scrollable Bubble Tea TUI list. Shows hash, message, author, branch, and timestamp for each commit.

```bash
diflog log
```

**Keyboard shortcuts in TUI:**

| Key | Action |
|---|---|
| `↑` / `↓` | Navigate commits |
| `q` / `Esc` | Quit |

---

### `diflog diff`

View a side-by-side line diff between two commits in a scrollable TUI viewport.

```bash
diflog diff <hash1> <hash2>

# Example
diflog diff abc12345 def67890
```

Hashes must be at least 4 characters long.

**Keyboard shortcuts in TUI:**

| Key | Action |
|---|---|
| `↑` / `↓` | Scroll diff |
| `q` / `Esc` | Quit |

---

### `diflog checkout`

Restore your project's AI context files to the state captured in a specific commit.

```bash
diflog checkout <hash>

# Example
diflog checkout abc12345678
```

> **Note:** This overwrites the current context files on disk. The full commit hash is required.

---

### `diflog branch`

Manage branches for isolated context version lines.

```bash
# Open interactive branch list TUI
diflog branch

# List all branches (same as above)
diflog branch list

# Create a new branch (branching from current HEAD)
diflog branch create feature/new-rules

# Switch to a branch
diflog branch switch feature/new-rules
```

---

### `diflog push`

Push all local committed contexts to a remote GitHub repository for backup and collaboration.

```bash
# Push to saved remote URL
diflog push

# Push to a specific remote URL (saves it for future use)
diflog push https://github.com/your-org/your-ai-contexts-repo
```

Requires the `DIFLOG_GITHUB_TOKEN` environment variable — see [Environment Variables](#environment-variables).

---

### `diflog pull`

Pull context commits from a remote GitHub repository into your local `.difLog/`.

```bash
# Pull from saved remote URL
diflog pull

# Pull from a specific remote URL
diflog pull https://github.com/your-org/your-ai-contexts-repo
```

---

### `diflog detect`

Detect which AI tool is configured in the current project directory.

```bash
diflog detect
```

Looks for indicators such as `CLAUDE.md`, `.claude/`, `opencode.jsonc` (ClaudeCode) or `.github/copilot-instructions.md` (Copilot CLI).

---

### `diflog status`

View currently staged contexts and the active branch.

```bash
diflog status
```

---

### `diflog skill create`

Install a DifLog skill/instructions file into your AI tool globally.

```bash
# Auto-detect AI tool from .difLog/config.json
diflog skill create

# Specify AI tool explicitly
diflog skill create claude
diflog skill create copilot
```

**ClaudeCode:** Creates `~/.config/opencode/skills/diflog/SKILL.md` — a trigger-based skill that lets the AI read your context history and save chat snapshots.

**GitHub Copilot CLI:** Appends DifLog instructions to `.github/copilot-instructions.md` in the current project.

---

### `diflog skill save`

Snapshot the current AI chat context as a versioned, hashed file inside `.difLog/skill-contexts/`.

```bash
# Snapshot from a file
diflog skill save context.md

# Snapshot from inline content
diflog skill save --content "Summary of today's session..." --source "session-2024-01-15"
```

**Flags:**

| Flag | Description |
|---|---|
| `--content` | Inline content to snapshot |
| `--source` | Label for the snapshot source (default: `manual`) |

Each snapshot is saved as `.difLog/skill-contexts/{sha256-hash}.json` containing the hash, timestamp, source label, and full content.

---

## Project Architecture

DifLog follows **CLEAN Architecture** with strict separation between domain logic, application use cases, infrastructure, and the presentation layer.

### CLEAN Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                  Presentation Layer                  │
│         CLI (Cobra)  +  TUI (Bubble Tea)             │
├─────────────────────────────────────────────────────┤
│                  Use Case Layer                      │
│   init │ add │ commit │ log │ diff │ checkout        │
│   branch │ push │ pull │ detect │ skill              │
├─────────────────────────────────────────────────────┤
│                  Domain Layer                        │
│    Entities  │  Repository Interfaces  │  Services   │
├─────────────────────────────────────────────────────┤
│               Infrastructure Layer                   │
│  Local Storage │ GitHub Remote │ Hashing │ Diffing   │
│  AI Detector │ ClaudeCode Adapter │ Copilot Adapter  │
└─────────────────────────────────────────────────────┘
```

**Dependency Rule:** Each layer only depends on the layer directly below it. The domain layer has zero external dependencies.

---

### Directory Structure

```
DifLog/
├── main.go                                  # Entry point — wires CLI root command
├── go.mod
├── go.sum
│
├── internal/
│   ├── domain/                              # Domain Layer (no external deps)
│   │   ├── entity/
│   │   │   ├── context.go                  # AIContext entity + AITool constants
│   │   │   ├── commit.go                   # Commit entity
│   │   │   ├── branch.go                   # Branch entity
│   │   │   ├── staging.go                  # StagingIndex + StagedContext
│   │   │   ├── skill_context.go            # SkillContext snapshot entity
│   │   │   └── remote.go                   # RemoteConfig entity
│   │   ├── repository/                     # Repository interfaces (ports)
│   │   │   ├── context_repository.go
│   │   │   ├── commit_repository.go
│   │   │   ├── branch_repository.go
│   │   │   ├── staging_repository.go
│   │   │   ├── skill_context_repository.go
│   │   │   └── remote_repository.go
│   │   └── service/                        # Domain service interfaces
│   │       ├── hasher_service.go           # Content hashing interface
│   │       ├── differ_service.go           # Diff computation interface
│   │       └── ai_detector_service.go      # AI tool detection interface
│   │
│   ├── usecase/                            # Use Case Layer (business logic)
│   │   ├── initialize/initialize_usecase.go
│   │   ├── add/add_context_usecase.go
│   │   ├── commit/commit_context_usecase.go
│   │   ├── log/view_commit_log_usecase.go
│   │   ├── diff/diff_commits_usecase.go
│   │   ├── checkout/checkout_commit_usecase.go
│   │   ├── branch/
│   │   │   ├── create_branch_usecase.go
│   │   │   ├── list_branches_usecase.go
│   │   │   └── switch_branch_usecase.go
│   │   ├── push/push_to_remote_usecase.go
│   │   ├── pull/pull_from_remote_usecase.go
│   │   ├── detect/detect_ai_tool_usecase.go
│   │   ├── skill/
│   │   │   ├── create_skill_usecase.go
│   │   │   └── save_skill_context_usecase.go
│   │   ├── status/view_status_usecase.go
│   │   ├── pathutil/pathutil.go            # Safe path utilities (prevent traversal)
│   │   └── integration_test.go             # init → add → commit integration test
│   │
│   ├── infrastructure/                     # Infrastructure Layer (concrete implementations)
│   │   ├── storage/
│   │   │   ├── local_storage_initializer.go     # Creates/reads .difLog directory
│   │   │   ├── context_local_repository.go      # Stores raw context content (objects/)
│   │   │   ├── commit_local_repository.go       # Reads/writes logs/commits.json
│   │   │   ├── branch_local_repository.go       # Reads/writes refs/heads/*
│   │   │   ├── staging_local_repository.go      # Reads/writes staging/index.json
│   │   │   └── skill_context_local_repository.go # Reads/writes skill-contexts/*.json
│   │   ├── remote/
│   │   │   └── github_remote_repository.go      # GitHub API push/pull via go-github
│   │   ├── hashing/
│   │   │   └── sha256_hasher_service.go         # SHA-256 content hashing
│   │   ├── diffing/
│   │   │   └── text_differ_service.go           # Line diff via go-diff
│   │   └── ai/
│   │       ├── ai_tool_detector.go              # Detects Claude/Copilot by file markers
│   │       ├── claude_code_adapter.go           # ClaudeCode skill content + context paths
│   │       └── copilot_cli_adapter.go           # Copilot CLI skill content + context paths
│   │
│   └── presentation/                       # Presentation Layer
│       ├── tui/
│       │   ├── app.go                           # TUI application bootstrap
│       │   ├── model/                           # Bubble Tea models (one per command)
│       │   │   ├── initialize_model.go
│       │   │   ├── ai_selection_model.go        # Interactive AI tool picker
│       │   │   ├── add_model.go
│       │   │   ├── commit_model.go
│       │   │   ├── log_model.go
│       │   │   ├── diff_model.go
│       │   │   ├── checkout_model.go
│       │   │   ├── branch_model.go
│       │   │   ├── detect_model.go
│       │   │   └── remote_model.go              # Shared push/pull model
│       │   └── component/                       # Reusable TUI components
│       │       ├── spinner_component.go
│       │       ├── list_component.go
│       │       ├── input_component.go
│       │       ├── diff_viewer_component.go
│       │       └── progress_component.go
│       └── cli/                                 # Cobra CLI commands + DI container
│           ├── root_command.go                  # Root command + dependency injection
│           ├── init_command.go
│           ├── add_command.go
│           ├── commit_command.go
│           ├── log_command.go
│           ├── diff_command.go
│           ├── checkout_command.go
│           ├── branch_command.go
│           ├── push_command.go                  # Contains both push + pull commands
│           ├── skill_command.go
│           ├── detect_command.go
│           └── status_command.go
│
└── pkg/
    └── errors/
        └── diflog_errors.go                     # Custom DifLogError type + error codes
```

---

### Domain Entities

#### `AIContext`
Represents a single AI context file captured at a point in time.

```go
type AIContext struct {
    Hash      string    // SHA-256 of file content
    FilePath  string    // Relative path of context file
    Content   string    // Raw file content
    AITool    AITool    // "claude" | "copilot"
    Timestamp time.Time
}
```

#### `Commit`
A snapshot of one or more staged `AIContext` entries.

```go
type Commit struct {
    Hash       string      // SHA-256 of commit contents
    Message    string
    ParentHash string      // Hash of previous commit (empty for first commit)
    Branch     string
    Contexts   []AIContext
    Timestamp  time.Time
    Author     string
}
```

#### `Branch`
A named pointer to a commit hash.

```go
type Branch struct {
    Name       string
    CommitHash string // Latest commit hash on this branch
    IsActive   bool
}
```

#### `StagingIndex`
The staging area — a list of contexts ready to be committed.

```go
type StagingIndex struct {
    Entries []StagedContext
}
type StagedContext struct {
    Hash      string
    FilePath  string
    AITool    AITool
    Timestamp time.Time
}
```

#### `SkillContext`
A versioned snapshot of a chat session saved by the AI skill.

```go
type SkillContext struct {
    Hash      string    // SHA-256 of content
    Timestamp time.Time
    Content   string    // Full snapshot content
    Source    string    // Label (e.g. "chat-2024-01-15")
    FilePath  string    // Path to saved .json file
}
```

---

### Data Flow

#### Staging and Committing

```
User runs: diflog add CLAUDE.md
    │
    ▼
AddContextUseCase
    │  reads file → computes SHA-256
    │  stores raw content in .difLog/objects/{hash}
    │  writes entry to .difLog/staging/index.json
    ▼
User runs: diflog commit -m "message"
    │
    ▼
CommitContextUseCase
    │  reads staging/index.json
    │  builds Commit entity with parent hash from current branch
    │  computes commit hash
    │  appends to .difLog/logs/commits.json
    │  updates .difLog/refs/heads/{branch} to new commit hash
    │  clears staging/index.json
    ▼
Done
```

#### Dependency Injection

All dependencies are wired in `internal/presentation/cli/root_command.go` through a `Container` struct. No global state is used. The container is constructed once per CLI invocation.

```
NewRootCommand()
    └── NewContainer(projectRoot)
            ├── storage.NewLocalStorageInitializer(...)
            ├── storage.NewCommitLocalRepository(...)
            ├── hashing.NewSHA256HasherService()
            ├── ai.NewAIToolDetector()
            ├── remote.NewGitHubRemoteRepository(token)
            └── usecase constructors (inject repos + services)
```

---

## .difLog Storage Format

All version control data is stored locally in the `.difLog/` directory at your project root.

```
.difLog/
├── config.json              # { "branch": "main", "remote": "", "aiTool": "claude" }
├── HEAD                     # Text: "ref: refs/heads/main"
│
├── staging/
│   └── index.json           # Array of StagedContext entries
│
├── objects/
│   └── {sha256-hash}        # Raw file content (content-addressed)
│
├── refs/
│   ├── heads/
│   │   └── main             # Text: "{latest-commit-hash}"
│   └── remote/
│       └── origin/
│           └── main         # Text: "{remote-commit-hash}"
│
├── logs/
│   └── commits.json         # Array of Commit records (full history)
│
└── skill-contexts/
    └── {sha256-hash}.json   # SkillContext snapshot files
```

**`config.json` schema:**
```json
{
  "branch": "main",
  "remote": "https://github.com/your-org/ai-contexts",
  "aiTool": "claude"
}
```

**`staging/index.json` schema:**
```json
[
  {
    "hash": "abc123...",
    "filePath": "CLAUDE.md",
    "aiTool": "claude",
    "timestamp": "2024-01-15T10:00:00Z"
  }
]
```

---

## AI Tool Support

### ClaudeCode (opencode)

**Detection markers** (checked in project root):
- `CLAUDE.md`
- `.claude/` directory
- `opencode.jsonc`
- `opencode.json`

**Context files tracked:**
- `CLAUDE.md`
- `.claude/settings.json`
- `.claude/settings.local.json`
- Any `*.md` files inside `.claude/`

**Skill installation path:** `~/.config/opencode/skills/diflog/SKILL.md`

---

### GitHub Copilot CLI

**Detection markers** (checked in project root):
- `.github/copilot-instructions.md`
- `.github/` directory

**Context files tracked:**
- `.github/copilot-instructions.md`
- Any `*.md` files inside `.github/`

**Skill installation:** Appends a DifLog instructions section to `.github/copilot-instructions.md`

---

## Remote Repository (GitHub)

DifLog uses the GitHub API to push and pull context data. It stores all commit and context data as a single JSON file (`diflog/data.json`) inside a GitHub repository.

### Setup

1. Create a GitHub Personal Access Token with `repo` scope:
   - Go to **GitHub → Settings → Developer Settings → Personal Access Tokens**
   - Generate a token with `repo` (full repo access) scope

2. Export the token:
   ```bash
   export DIFLOG_GITHUB_TOKEN=ghp_your_token_here
   ```

3. Push to your remote repository:
   ```bash
   diflog push https://github.com/your-org/your-ai-contexts-repo
   ```

The remote URL is saved to `.difLog/config.json` on first push, so subsequent `diflog push` / `diflog pull` commands don't need the URL argument.

---

## Skill System

Running `diflog skill create` installs DifLog awareness directly into your AI tool.

### For ClaudeCode (opencode)

Creates `~/.config/opencode/skills/diflog/SKILL.md` — a skill that the AI can trigger to:
- Browse `.difLog/logs/commits.json` and present your context history
- Guide you through restoring a previous context via `diflog checkout`
- Save the current chat context snapshot via `diflog skill save`

### For GitHub Copilot CLI

Appends a DifLog instructions block to `.github/copilot-instructions.md` so Copilot is aware of the version control workflow and can assist with context management commands.

### Saving Chat Context Snapshots

```bash
# Save an existing context file as a versioned snapshot
diflog skill save CLAUDE.md

# Save inline content (e.g., pasted from a chat summary)
diflog skill save --content "Today we refactored auth and added OAuth..." \
                  --source "session-2024-01-15"
```

Snapshots are stored in `.difLog/skill-contexts/{sha256}.json` and can be referenced by hash in future sessions.

---

## Environment Variables

| Variable | Description | Required |
|---|---|---|
| `DIFLOG_GITHUB_TOKEN` | GitHub Personal Access Token for push/pull | Only for `push` / `pull` |

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/usecase/...
go test ./internal/infrastructure/...
go test ./internal/presentation/tui/model/...

# Run with race detector
go test -race ./...
```

### Test Coverage by Package

| Package | Tests |
|---|---|
| `internal/infrastructure/ai` | AI tool detection, file path resolution |
| `internal/infrastructure/hashing` | SHA-256 determinism, consistency |
| `internal/infrastructure/diffing` | Line-level diff computation |
| `internal/infrastructure/storage` | Save/load commits, branches, skill-contexts |
| `internal/usecase/initialize` | Init creates correct .difLog structure |
| `internal/usecase/skill` | Save skill context from content and from file |
| `internal/presentation/tui/model` | AI selection TUI keyboard nav, defaults |
| `internal/usecase` (integration) | Full init → add → commit flow |

---

## Dependencies

| Package | Version | Purpose |
|---|---|---|
| `github.com/charmbracelet/bubbletea` | v1.3.10 | TUI framework |
| `github.com/charmbracelet/bubbles` | v1.0.0 | TUI components (spinner, list, textinput, viewport) |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | TUI styling and layout |
| `github.com/spf13/cobra` | v1.10.2 | CLI command framework |
| `github.com/google/go-github/v66` | v66.0.0 | GitHub API client |
| `golang.org/x/oauth2` | v0.36.0 | OAuth2 token transport for GitHub |
| `github.com/sergi/go-diff` | v1.4.0 | Text diff computation |
| `github.com/stretchr/testify` | v1.11.1 | Test assertions |
