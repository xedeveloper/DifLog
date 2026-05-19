package skill

import (
	"os"
	"time"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/domain/service"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type SaveSkillContextUseCase struct {
	skillContextRepo repository.SkillContextRepository
	hasherService    service.HasherService
}

func NewSaveSkillContextUseCase(
	skillContextRepo repository.SkillContextRepository,
	hasherService service.HasherService,
) *SaveSkillContextUseCase {
	return &SaveSkillContextUseCase{
		skillContextRepo: skillContextRepo,
		hasherService:    hasherService,
	}
}

type SaveSkillContextResult struct {
	Hash     string
	FilePath string
}

func (uc *SaveSkillContextUseCase) ExecuteFromFile(contextFilePath string) (SaveSkillContextResult, error) {
	content, err := os.ReadFile(contextFilePath)
	if err != nil {
		return SaveSkillContextResult{}, dlerrors.NewDifLogError(
			dlerrors.ErrContextNotFound,
			"failed to read context file: "+contextFilePath,
			err,
		)
	}
	return uc.ExecuteFromContent(string(content), contextFilePath)
}

func (uc *SaveSkillContextUseCase) ExecuteFromContent(content, source string) (SaveSkillContextResult, error) {
	if content == "" {
		return SaveSkillContextResult{}, dlerrors.NewDifLogError(
			dlerrors.ErrInvalidArgument,
			"context content cannot be empty",
			nil,
		)
	}

	hash := uc.hasherService.HashContent([]byte(content))

	skillCtx := entity.SkillContext{
		Hash:      hash,
		Timestamp: time.Now(),
		Content:   content,
		Source:    source,
	}

	if err := uc.skillContextRepo.SaveSkillContext(skillCtx); err != nil {
		return SaveSkillContextResult{}, err
	}

	return SaveSkillContextResult{
		Hash:     hash,
		FilePath: ".difLog/skill-contexts/" + hash + ".json",
	}, nil
}
