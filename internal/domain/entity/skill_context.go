package entity

import "time"

type SkillContext struct {
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Source    string    `json:"source"`
}
