package push

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type RemoteConfigStore interface {
	LoadRemote() (string, error)
	SaveRemote(remote string) error
}

type PushToRemoteUseCase struct {
	commitRepo    repository.CommitRepository
	contextRepo   repository.ContextRepository
	branchRepo    repository.BranchRepository
	remoteRepo    repository.RemoteRepository
	configStore   RemoteConfigStore
}

func NewPushToRemoteUseCase(
	commitRepo repository.CommitRepository,
	contextRepo repository.ContextRepository,
	branchRepo repository.BranchRepository,
	remoteRepo repository.RemoteRepository,
	configStore RemoteConfigStore,
) *PushToRemoteUseCase {
	return &PushToRemoteUseCase{
		commitRepo:  commitRepo,
		contextRepo: contextRepo,
		branchRepo:  branchRepo,
		remoteRepo:  remoteRepo,
		configStore: configStore,
	}
}

func (uc *PushToRemoteUseCase) Execute(remoteURL string) error {
	configRemote, err := uc.configStore.LoadRemote()
	if err != nil {
		return err
	}

	if remoteURL == "" {
		remoteURL = configRemote
	}

	if remoteURL == "" {
		return dlerrors.NewDifLogError(dlerrors.ErrRemoteNotSet, "no remote URL configured. Use: diflog push <url>", nil)
	}

	activeBranch, err := uc.branchRepo.GetActiveBranch()
	if err != nil {
		return err
	}

	commits, err := uc.commitRepo.ListCommits(activeBranch.Name)
	if err != nil {
		return err
	}

	objects := make(map[string][]byte)
	for _, commit := range commits {
		for _, ctx := range commit.Contexts {
			if _, loaded := objects[ctx.Hash]; !loaded {
				content, err := uc.contextRepo.LoadObject(ctx.Hash)
				if err != nil {
					return err
				}
				objects[ctx.Hash] = content
			}
		}
	}

	remote := entity.Remote{Name: "origin", URL: remoteURL}
	if err := uc.remoteRepo.PushCommits(remote, commits, objects); err != nil {
		return err
	}

	if configRemote != remoteURL {
		if err := uc.configStore.SaveRemote(remoteURL); err != nil {
			return err
		}
	}

	return nil
}
