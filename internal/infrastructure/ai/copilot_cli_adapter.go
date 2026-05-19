package ai

import (
	"fmt"
	"os"
	"path/filepath"
)

type CopilotCLIAdapter struct{}

func NewCopilotCLIAdapter() *CopilotCLIAdapter {
	return &CopilotCLIAdapter{}
}

func (a *CopilotCLIAdapter) AppendDifLogInstructions(projectRoot string) error {
	githubDir := filepath.Join(projectRoot, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		return err
	}

	instructionsPath := filepath.Join(githubDir, "copilot-instructions.md")
	diflogSection := a.diflogInstructionsSection()

	existing, err := os.ReadFile(instructionsPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	content := string(existing)
	if len(content) > 0 {
		content += "\n\n"
	}
	content += diflogSection

	return os.WriteFile(instructionsPath, []byte(content), 0644)
}

func (a *CopilotCLIAdapter) diflogInstructionsSection() string {
	return fmt.Sprintf(`## DifLog Context Version Control

This project uses DifLog for AI context version control.
Context history is stored in the ` + "`.difLog`" + ` directory.

### Viewing Context History

Read ` + "`.difLog/logs/commits.json`" + ` to view committed context history.
Each entry has: hash, message, timestamp, branch, and tracked context files.

### Saving Current Chat Context

When asked to save the current chat context, summarise the conversation and run:
` + "```" + `
diflog skill save --content "<summary>" --source "chat-session"
` + "```" + `
The snapshot is saved to ` + "`.difLog/skill-contexts/{hash}.json`" + ` with:
` + "```json" + `
{ "hash": "<sha256>", "timestamp": "<ISO-8601>", "content": "...", "source": "..." }
` + "```" + `

### Available Commands

| Command | Description |
|---|---|
| ` + "`diflog log`" + ` | View commit history |
| ` + "`diflog status`" + ` | Show staged contexts |
| ` + "`diflog add`" + ` | Stage context files |
| ` + "`diflog commit -m \"msg\"`" + ` | Commit staged contexts |
| ` + "`diflog checkout <hash>`" + ` | Restore context to commit |
| ` + "`diflog diff <hash1> <hash2>`" + ` | Compare two context versions |
| ` + "`diflog skill save <file>`" + ` | Save chat context snapshot |
`)
}
