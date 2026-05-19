package branch

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
)

type ListBranchesUseCase struct {
	branchRepo repository.BranchRepository
}

func NewListBranchesUseCase(branchRepo repository.BranchRepository) *ListBranchesUseCase {
	return &ListBranchesUseCase{branchRepo: branchRepo}
}

func (uc *ListBranchesUseCase) Execute() ([]entity.Branch, error) {
	return uc.branchRepo.ListBranches()
}
