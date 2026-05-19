package add

import (
	"os"
	"path/filepath"
	"time"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/domain/service"
	"github.com/xedeveloper/DifLog/internal/usecase/pathutil"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type AddContextUseCase struct {
	contextRepo    repository.ContextRepository
	stagingRepo    repository.StagingRepository
	hasherService  service.HasherService
	aiDetector     service.AIDetectorService
	projectRoot    string
}

func NewAddContextUseCase(
	contextRepo repository.ContextRepository,
	stagingRepo repository.StagingRepository,
	hasherService service.HasherService,
	aiDetector service.AIDetectorService,
	projectRoot string,
) *AddContextUseCase {
	return &AddContextUseCase{
		contextRepo:   contextRepo,
		stagingRepo:   stagingRepo,
		hasherService: hasherService,
		aiDetector:    aiDetector,
		projectRoot:   projectRoot,
	}
}

type AddResult struct {
	StagedFiles []string
	AITool      entity.AITool
}

func (uc *AddContextUseCase) Execute(filePaths []string, aiTool entity.AITool) (AddResult, error) {
	if len(filePaths) == 0 {
		detected := uc.aiDetector.GetContextFilePaths(uc.projectRoot, aiTool)
		filePaths = detected
	}

	if len(filePaths) == 0 {
		return AddResult{}, dlerrors.NewDifLogError(dlerrors.ErrContextNotFound, "no context files found for the selected AI tool", nil)
	}

	index, err := uc.stagingRepo.LoadStagingIndex()
	if err != nil {
		return AddResult{}, err
	}

	var stagedFiles []string
	for _, filePath := range filePaths {
		absPath := filePath
		if !filepath.IsAbs(filePath) {
			absPath = filepath.Join(uc.projectRoot, filePath)
		}
		relPath, err := pathutil.SafeRelativePath(uc.projectRoot, absPath)
		if err != nil {
			return AddResult{}, dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "context file must be inside project root", err)
		}

		content, err := os.ReadFile(absPath)
		if err != nil {
			return AddResult{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read context file: "+filePath, err)
		}

		hash := uc.hasherService.HashContent(content)

		if err := uc.contextRepo.SaveObject(hash, content); err != nil {
			return AddResult{}, err
		}

		alreadyStaged := false
		for i, entry := range index.Entries {
			if entry.FilePath == relPath {
				index.Entries[i].Hash = hash
				index.Entries[i].Timestamp = time.Now()
				alreadyStaged = true
				break
			}
		}

		if !alreadyStaged {
			index.Entries = append(index.Entries, entity.StagedContext{
				Hash:      hash,
				FilePath:  relPath,
				AITool:    aiTool,
				Timestamp: time.Now(),
			})
		}

		stagedFiles = append(stagedFiles, relPath)
	}

	if err := uc.stagingRepo.SaveStagingIndex(index); err != nil {
		return AddResult{}, err
	}

	return AddResult{StagedFiles: stagedFiles, AITool: aiTool}, nil
}
