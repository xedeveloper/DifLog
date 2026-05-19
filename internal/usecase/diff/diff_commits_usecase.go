package diff

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/domain/service"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type FileDiff struct {
	FilePath   string
	OldContent string
	NewContent string
	DiffText   string
}

type CommitDiffResult struct {
	OldCommit entity.Commit
	NewCommit entity.Commit
	FileDiffs []FileDiff
}

type DiffCommitsUseCase struct {
	commitRepo    repository.CommitRepository
	differService service.DifferService
}

func NewDiffCommitsUseCase(commitRepo repository.CommitRepository, differService service.DifferService) *DiffCommitsUseCase {
	return &DiffCommitsUseCase{commitRepo: commitRepo, differService: differService}
}

func (uc *DiffCommitsUseCase) Execute(oldHash, newHash string) (CommitDiffResult, error) {
	if oldHash == "" || newHash == "" {
		return CommitDiffResult{}, dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "both commit hashes are required", nil)
	}

	oldCommit, err := uc.commitRepo.LoadCommit(oldHash)
	if err != nil {
		return CommitDiffResult{}, err
	}

	newCommit, err := uc.commitRepo.LoadCommit(newHash)
	if err != nil {
		return CommitDiffResult{}, err
	}

	oldContextMap := buildContextMap(oldCommit.Contexts)
	newContextMap := buildContextMap(newCommit.Contexts)

	var fileDiffs []FileDiff

	for filePath, newCtx := range newContextMap {
		oldCtx, exists := oldContextMap[filePath]
		oldContent := ""
		if exists {
			oldContent = oldCtx.Content
		}
		diffText := uc.differService.UnifiedDiff(
			oldHash[:8]+"/"+filePath,
			newHash[:8]+"/"+filePath,
			oldContent,
			newCtx.Content,
		)
		fileDiffs = append(fileDiffs, FileDiff{
			FilePath:   filePath,
			OldContent: oldContent,
			NewContent: newCtx.Content,
			DiffText:   diffText,
		})
	}

	for filePath, oldCtx := range oldContextMap {
		if _, exists := newContextMap[filePath]; !exists {
			diffText := uc.differService.UnifiedDiff(
				oldHash[:8]+"/"+filePath,
				newHash[:8]+"/"+filePath,
				oldCtx.Content,
				"",
			)
			fileDiffs = append(fileDiffs, FileDiff{
				FilePath:   filePath,
				OldContent: oldCtx.Content,
				NewContent: "",
				DiffText:   diffText,
			})
		}
	}

	return CommitDiffResult{
		OldCommit: oldCommit,
		NewCommit: newCommit,
		FileDiffs: fileDiffs,
	}, nil
}

func buildContextMap(contexts []entity.AIContext) map[string]entity.AIContext {
	m := make(map[string]entity.AIContext)
	for _, ctx := range contexts {
		m[ctx.FilePath] = ctx
	}
	return m
}
