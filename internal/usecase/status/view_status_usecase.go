package status

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
)

type ViewStatusUseCase struct {
	stagingRepo repository.StagingRepository
	branchRepo  repository.BranchRepository
}

func NewViewStatusUseCase(stagingRepo repository.StagingRepository, branchRepo repository.BranchRepository) *ViewStatusUseCase {
	return &ViewStatusUseCase{stagingRepo: stagingRepo, branchRepo: branchRepo}
}

type StatusResult struct {
	ActiveBranch entity.Branch
	StagedFiles  []entity.StagedContext
}

func (uc *ViewStatusUseCase) Execute() (StatusResult, error) {
	activeBranch, err := uc.branchRepo.GetActiveBranch()
	if err != nil {
		return StatusResult{}, err
	}

	index, err := uc.stagingRepo.LoadStagingIndex()
	if err != nil {
		return StatusResult{}, err
	}

	return StatusResult{
		ActiveBranch: activeBranch,
		StagedFiles:  index.Entries,
	}, nil
}
