package entity

import "time"

type StagingIndex struct {
	Entries []StagedContext
}

type StagedContext struct {
	Hash      string
	FilePath  string
	AITool    AITool
	Timestamp time.Time
}
