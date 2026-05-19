package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type SkillContextLocalRepository struct {
	difLogDir string
}

func NewSkillContextLocalRepository(difLogDir string) *SkillContextLocalRepository {
	return &SkillContextLocalRepository{difLogDir: difLogDir}
}

func (r *SkillContextLocalRepository) skillContextsDir() string {
	return filepath.Join(r.difLogDir, "skill-contexts")
}

func (r *SkillContextLocalRepository) SaveSkillContext(ctx entity.SkillContext) error {
	if err := os.MkdirAll(r.skillContextsDir(), 0755); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create skill-contexts directory", err)
	}

	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to marshal skill context", err)
	}

	filePath := filepath.Join(r.skillContextsDir(), ctx.Hash+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to write skill context file", err)
	}

	return nil
}

func (r *SkillContextLocalRepository) ListSkillContexts() ([]entity.SkillContext, error) {
	entries, err := os.ReadDir(r.skillContextsDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []entity.SkillContext{}, nil
		}
		return nil, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to list skill contexts", err)
	}

	var contexts []entity.SkillContext
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(r.skillContextsDir(), entry.Name()))
		if err != nil {
			continue
		}
		var ctx entity.SkillContext
		if err := json.Unmarshal(data, &ctx); err != nil {
			continue
		}
		contexts = append(contexts, ctx)
	}

	return contexts, nil
}

func (r *SkillContextLocalRepository) LoadSkillContext(hash string) (entity.SkillContext, error) {
	filePath := filepath.Join(r.skillContextsDir(), hash+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return entity.SkillContext{}, dlerrors.NewDifLogError(dlerrors.ErrContextNotFound, "skill context not found: "+hash, err)
		}
		return entity.SkillContext{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read skill context", err)
	}

	var ctx entity.SkillContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		return entity.SkillContext{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to parse skill context", err)
	}

	return ctx, nil
}
