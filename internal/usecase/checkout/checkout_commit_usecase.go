package checkout

import (
	"os"
	"path/filepath"

	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/usecase/pathutil"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type CheckoutCommitUseCase struct {
	commitRepo  repository.CommitRepository
	projectRoot string
}

func NewCheckoutCommitUseCase(commitRepo repository.CommitRepository, projectRoot string) *CheckoutCommitUseCase {
	return &CheckoutCommitUseCase{commitRepo: commitRepo, projectRoot: projectRoot}
}

func (uc *CheckoutCommitUseCase) Execute(commitHash string) error {
	if commitHash == "" {
		return dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "commit hash is required", nil)
	}

	commit, err := uc.commitRepo.LoadCommit(commitHash)
	if err != nil {
		return err
	}

	for _, ctx := range commit.Contexts {
		absPath, err := pathutil.SafeJoin(uc.projectRoot, ctx.FilePath)
		if err != nil {
			return dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "commit contains unsafe file path: "+ctx.FilePath, err)
		}
		dir := filepath.Dir(absPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create directory for checkout", err)
		}
		if err := os.WriteFile(absPath, []byte(ctx.Content), 0644); err != nil {
			return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to restore context file: "+ctx.FilePath, err)
		}
	}

	return nil
}
