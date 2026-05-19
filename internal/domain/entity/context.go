package entity

import "time"

type AITool string

const (
	AIToolClaude   AITool = "claude"
	AIToolOpenCode AITool = "opencode"
	AIToolCopilot  AITool = "copilot"
	AIToolUnknown  AITool = "unknown"
)

type AIContext struct {
	Hash      string
	FilePath  string
	Content   string
	AITool    AITool
	Timestamp time.Time
}
