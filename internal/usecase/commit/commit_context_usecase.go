package commit

import (
	"fmt"
	"os/user"
	"time"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/domain/service"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type CommitContextUseCase struct {
	commitRepo    repository.CommitRepository
	stagingRepo   repository.StagingRepository
	contextRepo   repository.ContextRepository
	branchRepo    repository.BranchRepository
	hasherService service.HasherService
}

func NewCommitContextUseCase(
	commitRepo repository.CommitRepository,
	stagingRepo repository.StagingRepository,
	contextRepo repository.ContextRepository,
	branchRepo repository.BranchRepository,
	hasherService service.HasherService,
) *CommitContextUseCase {
	return &CommitContextUseCase{
		commitRepo:    commitRepo,
		stagingRepo:   stagingRepo,
		contextRepo:   contextRepo,
		branchRepo:    branchRepo,
		hasherService: hasherService,
	}
}

func (uc *CommitContextUseCase) Execute(message string) (entity.Commit, error) {
	if message == "" {
		return entity.Commit{}, dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "commit message cannot be empty", nil)
	}

	index, err := uc.stagingRepo.LoadStagingIndex()
	if err != nil {
		return entity.Commit{}, err
	}

	if len(index.Entries) == 0 {
		return entity.Commit{}, dlerrors.NewDifLogError(dlerrors.ErrNoStagedContexts, "no staged contexts to commit", nil)
	}

	activeBranch, err := uc.branchRepo.GetActiveBranch()
	if err != nil {
		return entity.Commit{}, err
	}

	parentHash, _ := uc.commitRepo.GetLatestCommitHash(activeBranch.Name)

	var contexts []entity.AIContext
	var contentSlices [][]byte
	for _, staged := range index.Entries {
		content, err := uc.contextRepo.LoadObject(staged.Hash)
		if err != nil {
			return entity.Commit{}, err
		}
		contexts = append(contexts, entity.AIContext{
			Hash:      staged.Hash,
			FilePath:  staged.FilePath,
			Content:   string(content),
			AITool:    staged.AITool,
			Timestamp: staged.Timestamp,
		})
		contentSlices = append(contentSlices, content)
	}

	commitHash := uc.hasherService.HashMultiple(contentSlices)
	commitHash = fmt.Sprintf("%s%d", commitHash[:16], time.Now().UnixNano())
	commitHash = uc.hasherService.HashContent([]byte(commitHash))

	author := resolveAuthor()

	commit := entity.Commit{
		Hash:       commitHash,
		Message:    message,
		ParentHash: parentHash,
		Branch:     activeBranch.Name,
		Contexts:   contexts,
		Timestamp:  time.Now(),
		Author:     author,
	}

	if err := uc.commitRepo.SaveCommit(commit); err != nil {
		return entity.Commit{}, err
	}

	if err := uc.stagingRepo.ClearStagingIndex(); err != nil {
		return entity.Commit{}, err
	}

	return commit, nil
}

func resolveAuthor() string {
	u, err := user.Current()
	if err != nil {
		return "unknown"
	}
	if u.Name != "" {
		return u.Name
	}
	return u.Username
}
