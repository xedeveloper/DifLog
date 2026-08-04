package ai

import (
	"os"
	"path/filepath"
)

type ClaudeCodeAdapter struct{}

func NewClaudeCodeAdapter() *ClaudeCodeAdapter {
	return &ClaudeCodeAdapter{}
}

func (a *ClaudeCodeAdapter) CreateSkill(skillContent string) error {
	skillDir := filepath.Join(os.Getenv("HOME"), ".claude", "skills", "diflog")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}
	skillPath := filepath.Join(skillDir, "SKILL.md")
	return os.WriteFile(skillPath, []byte(skillContent), 0644)
}

func (a *ClaudeCodeAdapter) SkillContent() string {
	return `---
name: diflog
description: "Fetch and restore AI context history from DifLog version control. Use when you need to access previous AI context versions, restore a context, view context history, or save the current chat context as a versioned snapshot."
trigger: /diflog
---

# DifLog — AI Context Version Control Skill

This skill interfaces with the DifLog ` + "`.difLog`" + ` directory to manage AI context history.

## Viewing Context History

To view committed context history, read the file at:
` + "`.difLog/logs/commits.json`" + `

Each entry contains:
- ` + "`hash`" + ` — unique commit identifier
- ` + "`message`" + ` — commit message
- ` + "`timestamp`" + ` — when it was committed
- ` + "`branch`" + ` — which branch
- ` + "`contexts`" + ` — array of tracked context files with their content

Present the history as a numbered list with hash (first 8 chars), message, and timestamp.

## Saving Current Chat Context

When the user asks to save the current chat context or when triggered with ` + "`/diflog save`" + `:

1. Summarise the current conversation into a structured markdown document
2. Run the following shell command to snapshot it:
   ` + "```" + `
   diflog skill save --content "<summary>" --source "chat-session"
   ` + "```" + `
   Or write the summary to a temp file and run:
   ` + "```" + `
   diflog skill save <path-to-file>
   ` + "```" + `
3. The snapshot is saved to ` + "`.difLog/skill-contexts/{hash}.json`" + ` with structure:
   ` + "```json" + `
   {
     "hash": "<sha256-of-content>",
     "timestamp": "<ISO-8601>",
     "content": "<context-text>",
     "source": "<source-label>"
   }
   ` + "```" + `
4. Confirm to the user: "Context saved with hash ` + "`{hash}`" + `"

### Creating local Stage

After saving the context, create a local file if not exists with name "context.local.md".

1. If file already exists then override the content of the file.
2. Inside this file create a insert the source generated from above context-text.
3. This should be markdown file that stores the context temporarily for claude code.
4. The structure of markdown should be like this:
	"# Hash: <hash generated from previous step>
	 # Content : <context-text>"

## Restoring a Context

To restore context files to a specific commit:
` + "```" + `
diflog checkout <hash>
` + "```" + `

## Available CLI Commands

| Command | Description |
|---|---|
| ` + "`diflog log`" + ` | View commit history |
| ` + "`diflog status`" + ` | Show staged contexts |
| ` + "`diflog add`" + ` | Stage context files |
| ` + "`diflog commit -m \"msg\"`" + ` | Commit staged contexts |
| ` + "`diflog checkout <hash>`" + ` | Restore context to commit |
| ` + "`diflog diff <h1> <h2>`" + ` | Diff two commits |
| ` + "`diflog skill save <file>`" + ` | Save chat context snapshot |
| ` + "`diflog branch list`" + ` | List branches |
| ` + "`diflog branch switch <name>`" + ` | Switch branch |
`
}
