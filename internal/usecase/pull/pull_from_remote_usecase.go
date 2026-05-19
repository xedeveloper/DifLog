package pull

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type RemoteConfigStore interface {
	LoadRemote() (string, error)
	SaveRemote(remote string) error
}

type PullFromRemoteUseCase struct {
	commitRepo  repository.CommitRepository
	contextRepo repository.ContextRepository
	branchRepo  repository.BranchRepository
	remoteRepo  repository.RemoteRepository
	configStore RemoteConfigStore
}

func NewPullFromRemoteUseCase(
	commitRepo repository.CommitRepository,
	contextRepo repository.ContextRepository,
	branchRepo repository.BranchRepository,
	remoteRepo repository.RemoteRepository,
	configStore RemoteConfigStore,
) *PullFromRemoteUseCase {
	return &PullFromRemoteUseCase{
		commitRepo:  commitRepo,
		contextRepo: contextRepo,
		branchRepo:  branchRepo,
		remoteRepo:  remoteRepo,
		configStore: configStore,
	}
}

func (uc *PullFromRemoteUseCase) Execute(remoteURL string) error {
	configRemote, err := uc.configStore.LoadRemote()
	if err != nil {
		return err
	}

	if remoteURL == "" {
		remoteURL = configRemote
	}

	if remoteURL == "" {
		return dlerrors.NewDifLogError(dlerrors.ErrRemoteNotSet, "no remote URL configured. Use: diflog pull <url>", nil)
	}

	activeBranch, err := uc.branchRepo.GetActiveBranch()
	if err != nil {
		return err
	}

	remote := entity.Remote{Name: "origin", URL: remoteURL}
	commits, objects, err := uc.remoteRepo.PullCommits(remote, activeBranch.Name)
	if err != nil {
		return err
	}

	for hash, content := range objects {
		if !uc.contextRepo.ObjectExists(hash) {
			if err := uc.contextRepo.SaveObject(hash, content); err != nil {
				return err
			}
		}
	}

	for _, commit := range commits {
		if err := uc.commitRepo.SaveCommit(commit); err != nil {
			return err
		}
	}

	return nil
}
