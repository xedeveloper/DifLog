package repository

import "github.com/xedeveloper/DifLog/internal/domain/entity"

type CommitRepository interface {
	SaveCommit(commit entity.Commit) error
	LoadCommit(hash string) (entity.Commit, error)
	ListCommits(branch string) ([]entity.Commit, error)
	GetLatestCommitHash(branch string) (string, error)
}
