package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/usecase/pathutil"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type CommitLocalRepository struct {
	difLogDir string
}

func NewCommitLocalRepository(difLogDir string) *CommitLocalRepository {
	return &CommitLocalRepository{difLogDir: difLogDir}
}

func (r *CommitLocalRepository) SaveCommit(commit entity.Commit) error {
	if err := pathutil.ValidateBranchName(commit.Branch); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrInvalidArgument, "invalid branch name for commit", err)
	}

	commits, err := r.loadAllCommits()
	if err != nil {
		return err
	}
	for i, existing := range commits {
		if existing.Hash == commit.Hash {
			commits[i] = commit
			if err := r.saveAllCommits(commits); err != nil {
				return err
			}
			refPath := filepath.Join(r.difLogDir, "refs", "heads", commit.Branch)
			if err := os.WriteFile(refPath, []byte(commit.Hash), 0644); err != nil {
				return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to update branch ref", err)
			}
			return nil
		}
	}

	commits = append(commits, commit)

	if err := r.saveAllCommits(commits); err != nil {
		return err
	}

	refPath := filepath.Join(r.difLogDir, "refs", "heads", commit.Branch)
	if err := os.WriteFile(refPath, []byte(commit.Hash), 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to update branch ref", err)
	}

	return nil
}

func (r *CommitLocalRepository) LoadCommit(hash string) (entity.Commit, error) {
	commits, err := r.loadAllCommits()
	if err != nil {
		return entity.Commit{}, err
	}

	for _, c := range commits {
		if c.Hash == hash {
			return c, nil
		}
	}

	return entity.Commit{}, dlerrors.NewDifLogError(dlerrors.ErrCommitNotFound, "commit not found: "+hash, nil)
}

func (r *CommitLocalRepository) ListCommits(branch string) ([]entity.Commit, error) {
	commits, err := r.loadAllCommits()
	if err != nil {
		return nil, err
	}

	var branchCommits []entity.Commit
	for _, c := range commits {
		if c.Branch == branch {
			branchCommits = append(branchCommits, c)
		}
	}

	return branchCommits, nil
}

func (r *CommitLocalRepository) GetLatestCommitHash(branch string) (string, error) {
	refPath := filepath.Join(r.difLogDir, "refs", "heads", branch)
	data, err := os.ReadFile(refPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read branch ref", err)
	}
	return string(data), nil
}

func (r *CommitLocalRepository) loadAllCommits() ([]entity.Commit, error) {
	commitsPath := filepath.Join(r.difLogDir, "logs", "commits.json")
	data, err := os.ReadFile(commitsPath)
	if err != nil {
		return nil, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read commits log", err)
	}

	var commits []entity.Commit
	if err := json.Unmarshal(data, &commits); err != nil {
		return nil, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to parse commits log", err)
	}

	return commits, nil
}

func (r *CommitLocalRepository) saveAllCommits(commits []entity.Commit) error {
	data, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to marshal commits", err)
	}

	commitsPath := filepath.Join(r.difLogDir, "logs", "commits.json")
	if err := os.WriteFile(commitsPath, data, 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to write commits log", err)
	}

	return nil
}
