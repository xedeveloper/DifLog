package log

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
)

type ViewCommitLogUseCase struct {
	commitRepo repository.CommitRepository
	branchRepo repository.BranchRepository
}

func NewViewCommitLogUseCase(commitRepo repository.CommitRepository, branchRepo repository.BranchRepository) *ViewCommitLogUseCase {
	return &ViewCommitLogUseCase{commitRepo: commitRepo, branchRepo: branchRepo}
}

func (uc *ViewCommitLogUseCase) Execute() ([]entity.Commit, error) {
	activeBranch, err := uc.branchRepo.GetActiveBranch()
	if err != nil {
		return nil, err
	}
	return uc.commitRepo.ListCommits(activeBranch.Name)
}

func (uc *ViewCommitLogUseCase) ExecuteForBranch(branch string) ([]entity.Commit, error) {
	return uc.commitRepo.ListCommits(branch)
}
