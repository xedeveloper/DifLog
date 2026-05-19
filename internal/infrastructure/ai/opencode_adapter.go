package ai

import (
	"os"
	"path/filepath"
)

type OpenCodeAdapter struct{}

func NewOpenCodeAdapter() *OpenCodeAdapter {
	return &OpenCodeAdapter{}
}

func (a *OpenCodeAdapter) CreateSkill(skillContent string) error {
	skillDir := filepath.Join(os.Getenv("HOME"), ".config", "opencode", "skills", "diflog")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}
	skillPath := filepath.Join(skillDir, "SKILL.md")
	return os.WriteFile(skillPath, []byte(skillContent), 0644)
}

func (a *OpenCodeAdapter) SkillContent() string {
	return `---
name: diflog
description: "Manage DifLog version control for your OpenCode configuration — browse context history, restore opencode.json and .opencode/ configs to a previous commit, or snapshot the current session. Trigger with /diflog."
trigger: /diflog
---

# DifLog — OpenCode Configuration Version Control

DifLog tracks your OpenCode project configuration under version control.
Managed files include ` + "`opencode.json`" + `, ` + "`opencode.jsonc`" + `, and everything inside ` + "`.opencode/`" + ` (agents, skills, rules, plugins).

## Browsing Configuration History

Read ` + "`.difLog/logs/commits.json`" + ` to view the full commit history.

Each commit entry contains:
- ` + "`hash`" + ` — unique commit identifier (SHA-256)
- ` + "`message`" + ` — description of what changed
- ` + "`timestamp`" + ` — when the commit was made
- ` + "`branch`" + ` — which branch this belongs to
- ` + "`contexts`" + ` — array of tracked files with their content at that point

Present history as a numbered list: hash (first 8 chars), message, timestamp.

## Tracked OpenCode Files

DifLog versions the following files for this project:
- ` + "`opencode.json`" + ` / ` + "`opencode.jsonc`" + ` — main OpenCode configuration
- ` + "`.opencode/agents/*.md`" + ` — custom agent definitions
- ` + "`.opencode/skills/*.md`" + ` — project-level skill files
- ` + "`.opencode/rules/*.md`" + ` — coding rules and constraints
- Any other files tracked under ` + "`.opencode/`" + `

## Saving the Current Session Context

When the user asks to save the current chat context, or when triggered with ` + "`/diflog save`" + `:

1. Summarise the current conversation — key decisions, code changes, context discussed
2. Save it as a versioned snapshot:
   ` + "```bash" + `
   diflog skill save --content "<summary>" --source "opencode-session"
   ` + "```" + `
   Or write to a file first:
   ` + "```bash" + `
   diflog skill save <path-to-summary-file>
   ` + "```" + `
3. The snapshot is stored at ` + "`.difLog/skill-contexts/{hash}.json`" + `:
   ` + "```json" + `
   {
     "hash": "<sha256>",
     "timestamp": "<ISO-8601>",
     "content": "<session-summary>",
     "source": "opencode-session"
   }
   ` + "```" + `
4. Confirm: "Session saved with hash ` + "`{hash}`" + `"

## Restoring a Configuration

To restore ` + "`opencode.json`" + ` and ` + "`.opencode/`" + ` files to a previous commit:
` + "```bash" + `
diflog checkout <hash>
` + "```" + `
This overwrites the current config files on disk with the versions from that commit.

## Available CLI Commands

| Command | Description |
|---|---|
| ` + "`diflog log`" + ` | Browse full commit history |
| ` + "`diflog status`" + ` | Show staged files and active branch |
| ` + "`diflog add`" + ` | Stage opencode config files |
| ` + "`diflog commit -m \"msg\"`" + ` | Commit staged configs |
| ` + "`diflog checkout <hash>`" + ` | Restore configs to a previous commit |
| ` + "`diflog diff <h1> <h2>`" + ` | Diff two commits |
| ` + "`diflog skill save <file>`" + ` | Save session context snapshot |
| ` + "`diflog branch list`" + ` | List all branches |
| ` + "`diflog branch switch <name>`" + ` | Switch to a branch |
`
}
