package entity

import "time"

type Commit struct {
	Hash       string
	Message    string
	ParentHash string
	Branch     string
	Contexts   []AIContext
	Timestamp  time.Time
	Author     string
}
