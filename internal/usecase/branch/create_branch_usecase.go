package branch

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/usecase/pathutil"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type CreateBranchUseCase struct {
	branchRepo repository.BranchRepository
	commitRepo repository.CommitRepository
}

func NewCreateBranchUseCase(branchRepo repository.BranchRepository, commitRepo repository.CommitRepository) *CreateBranchUseCase {
	return &CreateBranchUseCase{branchRepo: branchRepo, commitRepo: commitRepo}
}

func (uc *CreateBranchUseCase) Execute(name string) (entity.Branch, error) {
	if err := pathutil.ValidateBranchName(name); err != nil {
		return entity.Branch{}, dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, err.Error(), err)
	}

	if uc.branchRepo.BranchExists(name) {
		return entity.Branch{}, dlerrors.NewDifLogError(dlerrors.ErrBranchExists, "branch already exists: "+name, nil)
	}

	activeBranch, err := uc.branchRepo.GetActiveBranch()
	if err != nil {
		return entity.Branch{}, err
	}

	newBranch := entity.Branch{
		Name:       name,
		CommitHash: activeBranch.CommitHash,
		IsActive:   false,
	}

	if err := uc.branchRepo.SaveBranch(newBranch); err != nil {
		return entity.Branch{}, err
	}

	return newBranch, nil
}
