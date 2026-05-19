package repository

import "github.com/xedeveloper/DifLog/internal/domain/entity"

type BranchRepository interface {
	SaveBranch(branch entity.Branch) error
	LoadBranch(name string) (entity.Branch, error)
	ListBranches() ([]entity.Branch, error)
	GetActiveBranch() (entity.Branch, error)
	SetActiveBranch(name string) error
	BranchExists(name string) bool
}
