package storage

import (
	"os"
	"path/filepath"

	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type ContextLocalRepository struct {
	difLogDir string
}

func NewContextLocalRepository(difLogDir string) *ContextLocalRepository {
	return &ContextLocalRepository{difLogDir: difLogDir}
}

func (r *ContextLocalRepository) SaveObject(hash string, content []byte) error {
	objectPath := filepath.Join(r.difLogDir, "objects", hash)
	if err := os.WriteFile(objectPath, content, 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to save object: "+hash, err)
	}
	return nil
}

func (r *ContextLocalRepository) LoadObject(hash string) ([]byte, error) {
	objectPath := filepath.Join(r.difLogDir, "objects", hash)
	data, err := os.ReadFile(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, dlerrors.NewDifLogError(dlerrors.ErrContextNotFound, "object not found: "+hash, err)
		}
		return nil, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to load object: "+hash, err)
	}
	return data, nil
}

func (r *ContextLocalRepository) ObjectExists(hash string) bool {
	objectPath := filepath.Join(r.difLogDir, "objects", hash)
	_, err := os.Stat(objectPath)
	return err == nil
}
