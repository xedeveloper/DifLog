package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type StagingLocalRepository struct {
	difLogDir string
}

func NewStagingLocalRepository(difLogDir string) *StagingLocalRepository {
	return &StagingLocalRepository{difLogDir: difLogDir}
}

func (r *StagingLocalRepository) LoadStagingIndex() (entity.StagingIndex, error) {
	indexPath := filepath.Join(r.difLogDir, "staging", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return entity.StagingIndex{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read staging index", err)
	}

	var index entity.StagingIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return entity.StagingIndex{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to parse staging index", err)
	}

	if index.Entries == nil {
		index.Entries = []entity.StagedContext{}
	}

	return index, nil
}

func (r *StagingLocalRepository) SaveStagingIndex(index entity.StagingIndex) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to marshal staging index", err)
	}

	indexPath := filepath.Join(r.difLogDir, "staging", "index.json")
	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to write staging index", err)
	}

	return nil
}

func (r *StagingLocalRepository) ClearStagingIndex() error {
	empty := entity.StagingIndex{Entries: []entity.StagedContext{}}
	return r.SaveStagingIndex(empty)
}
