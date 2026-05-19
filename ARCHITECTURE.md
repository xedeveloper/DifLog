## Table of Content

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
