package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/usecase/pathutil"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type BranchLocalRepository struct {
	difLogDir   string
	initializer *LocalStorageInitializer
}

func NewBranchLocalRepository(difLogDir string, initializer *LocalStorageInitializer) *BranchLocalRepository {
	return &BranchLocalRepository{difLogDir: difLogDir, initializer: initializer}
}

func (r *BranchLocalRepository) SaveBranch(branch entity.Branch) error {
	if err := pathutil.ValidateBranchName(branch.Name); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "invalid branch name", err)
	}
	refPath := filepath.Join(r.difLogDir, "refs", "heads", branch.Name)
	if err := os.WriteFile(refPath, []byte(branch.CommitHash), 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to save branch ref", err)
	}
	return nil
}

func (r *BranchLocalRepository) LoadBranch(name string) (entity.Branch, error) {
	if err := pathutil.ValidateBranchName(name); err != nil {
		return entity.Branch{}, dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "invalid branch name", err)
	}
	refPath := filepath.Join(r.difLogDir, "refs", "heads", name)
	data, err := os.ReadFile(refPath)
	if err != nil {
		if os.IsNotExist(err) {
			return entity.Branch{}, dlerrors.NewDifLogError(dlerrors.ErrBranchNotFound, "branch not found: "+name, err)
		}
		return entity.Branch{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read branch ref", err)
	}

	activeBranch, _ := r.GetActiveBranch()
	return entity.Branch{
		Name:       name,
		CommitHash: string(data),
		IsActive:   activeBranch.Name == name,
	}, nil
}

func (r *BranchLocalRepository) ListBranches() ([]entity.Branch, error) {
	headsDir := filepath.Join(r.difLogDir, "refs", "heads")
	entries, err := os.ReadDir(headsDir)
	if err != nil {
		return nil, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to list branches", err)
	}

	activeBranch, _ := r.GetActiveBranch()
	var branches []entity.Branch
	for _, entry := range entries {
		if !entry.IsDir() {
			data, _ := os.ReadFile(filepath.Join(headsDir, entry.Name()))
			branches = append(branches, entity.Branch{
				Name:       entry.Name(),
				CommitHash: string(data),
				IsActive:   activeBranch.Name == entry.Name(),
			})
		}
	}

	return branches, nil
}

func (r *BranchLocalRepository) GetActiveBranch() (entity.Branch, error) {
	headPath := filepath.Join(r.difLogDir, "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return entity.Branch{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read HEAD", err)
	}

	headContent := strings.TrimSpace(string(data))
	branchName := strings.TrimPrefix(headContent, "ref: refs/heads/")
	if err := pathutil.ValidateBranchName(branchName); err != nil {
		return entity.Branch{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "invalid active branch in HEAD", err)
	}

	refPath := filepath.Join(r.difLogDir, "refs", "heads", branchName)
	refData, _ := os.ReadFile(refPath)

	return entity.Branch{
		Name:       branchName,
		CommitHash: strings.TrimSpace(string(refData)),
		IsActive:   true,
	}, nil
}

func (r *BranchLocalRepository) SetActiveBranch(name string) error {
	return r.initializer.UpdateHEAD(name)
}

func (r *BranchLocalRepository) BranchExists(name string) bool {
	refPath := filepath.Join(r.difLogDir, "refs", "heads", name)
	_, err := os.Stat(refPath)
	return err == nil
}

type branchesFile struct {
	Branches []entity.Branch `json:"branches"`
}

func marshalBranches(branches []entity.Branch) ([]byte, error) {
	return json.MarshalIndent(branchesFile{Branches: branches}, "", "  ")
}
