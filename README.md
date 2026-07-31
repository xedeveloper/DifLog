<div align="center">
  <a href="https://github.com/xedeveloper/DifLog">
    <img
      src="assets/diflog-logo.png"
    />
  </a>
  <br/>
  <br/>
<h1>A Git-like version control tool for AI context & memory</h1>
  <p>Built for developers who use ClaudeCode, OpenCode & Github Copilot CLI</p>
</div>

[![Go Version](https://img.shields.io/github/go-mod/go-version/regent-vcs/regent?style=for-the-badge&logo=go&logoColor=white&color=00ADD8)](go.mod)
[![Copilot CLI Compatible](https://img.shields.io/badge/Copilot%20CLI-Compatible-5dbf33?style=for-the-badge&logo=github&logoColor=white)](https://github.com/xedeveloper/DifLog)[![Claude Code Compatible](https://img.shields.io/badge/Claude%20Code-Compatible-6366f1?style=for-the-badge&logo=anthropic&logoColor=white)](https://github.com/xedeveloper/DifLog) [![Codex Compatible](https://img.shields.io/badge/Codex-Compatible-10b981?style=for-the-badge&logo=openai&logoColor=white)](https://github.com/xedeveloper/DifLog) [![OpenCode Compatible](https://img.shields.io/badge/OpenCode-Compatible-ff6b35?style=for-the-badge)](https://github.com/xedeveloper/DifLog)

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

**ClaudeCode:** Creates `~/.claude/skills/diflog/SKILL.md` — a trigger-based skill that lets the AI read your context history and save chat snapshots.

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

## License

 [MIT License](LICENSE)

---

<div align="center">
  <p>
    <sub>Built by <a href="https://github.com/xedeveloper">xedeveloper</a></sub>
  </p>
</div>
