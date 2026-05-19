package repository

import "github.com/xedeveloper/DifLog/internal/domain/entity"

type RemoteRepository interface {
	PushCommits(remote entity.Remote, commits []entity.Commit, objects map[string][]byte) error
	PullCommits(remote entity.Remote, branch string) ([]entity.Commit, map[string][]byte, error)
}
