package branch

import (
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/usecase/pathutil"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type SwitchBranchUseCase struct {
	branchRepo repository.BranchRepository
}

func NewSwitchBranchUseCase(branchRepo repository.BranchRepository) *SwitchBranchUseCase {
	return &SwitchBranchUseCase{branchRepo: branchRepo}
}

func (uc *SwitchBranchUseCase) Execute(name string) error {
	if err := pathutil.ValidateBranchName(name); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, err.Error(), err)
	}

	if !uc.branchRepo.BranchExists(name) {
		return dlerrors.NewDifLogError(dlerrors.ErrBranchNotFound, "branch not found: "+name, nil)
	}

	return uc.branchRepo.SetActiveBranch(name)
}
