package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/infrastructure/storage"
)

func setupInitializer(t *testing.T) (*storage.LocalStorageInitializer, string) {
	dir := t.TempDir()
	initializer := storage.NewLocalStorageInitializer(dir)
	return initializer, dir
}

func TestLocalStorageInitializer_Initialize(t *testing.T) {
	initializer, _ := setupInitializer(t)

	err := initializer.Initialize(entity.AIToolClaude)
	require.NoError(t, err)

	assert.True(t, initializer.IsInitialized())
}

func TestLocalStorageInitializer_AlreadyInitialized(t *testing.T) {
	initializer, _ := setupInitializer(t)

	require.NoError(t, initializer.Initialize(entity.AIToolClaude))
	err := initializer.Initialize(entity.AIToolClaude)
	assert.Error(t, err)
}

func TestLocalStorageInitializer_LoadConfig(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	config, err := initializer.LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "main", config.Branch)
	assert.Equal(t, entity.AIToolClaude, config.AITool)
}

func TestContextLocalRepository_SaveAndLoad(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewContextLocalRepository(initializer.DifLogDir())

	content := []byte("# Test Context")
	hash := "abc123def456"

	err := repo.SaveObject(hash, content)
	require.NoError(t, err)

	loaded, err := repo.LoadObject(hash)
	require.NoError(t, err)
	assert.Equal(t, content, loaded)
}

func TestContextLocalRepository_ObjectExists(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewContextLocalRepository(initializer.DifLogDir())

	assert.False(t, repo.ObjectExists("nonexistent"))

	require.NoError(t, repo.SaveObject("exists123", []byte("content")))
	assert.True(t, repo.ObjectExists("exists123"))
}

func TestStagingLocalRepository_SaveAndLoad(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewStagingLocalRepository(initializer.DifLogDir())

	index, err := repo.LoadStagingIndex()
	require.NoError(t, err)
	assert.Empty(t, index.Entries)

	index.Entries = append(index.Entries, entity.StagedContext{
		Hash:     "hash1",
		FilePath: "CLAUDE.md",
		AITool:   entity.AIToolClaude,
	})

	require.NoError(t, repo.SaveStagingIndex(index))

	loaded, err := repo.LoadStagingIndex()
	require.NoError(t, err)
	assert.Len(t, loaded.Entries, 1)
	assert.Equal(t, "CLAUDE.md", loaded.Entries[0].FilePath)
}

func TestBranchLocalRepository_GetActiveBranch(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewBranchLocalRepository(initializer.DifLogDir(), initializer)

	branch, err := repo.GetActiveBranch()
	require.NoError(t, err)
	assert.Equal(t, "main", branch.Name)
	assert.True(t, branch.IsActive)
}

func TestBranchLocalRepository_CreateAndSwitch(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewBranchLocalRepository(initializer.DifLogDir(), initializer)

	newBranch := entity.Branch{Name: "feature", CommitHash: ""}
	require.NoError(t, repo.SaveBranch(newBranch))
	assert.True(t, repo.BranchExists("feature"))

	require.NoError(t, repo.SetActiveBranch("feature"))

	active, err := repo.GetActiveBranch()
	require.NoError(t, err)
	assert.Equal(t, "feature", active.Name)
}

func TestCommitLocalRepository_SaveAndLoad(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	// Create the main ref file
	mainRefPath := filepath.Join(initializer.DifLogDir(), "refs", "heads", "main")
	os.WriteFile(mainRefPath, []byte(""), 0644)

	repo := storage.NewCommitLocalRepository(initializer.DifLogDir())

	commit := entity.Commit{
		Hash:    "abc123",
		Message: "test commit",
		Branch:  "main",
		Author:  "tester",
	}

	require.NoError(t, repo.SaveCommit(commit))

	loaded, err := repo.LoadCommit("abc123")
	require.NoError(t, err)
	assert.Equal(t, "test commit", loaded.Message)

	commits, err := repo.ListCommits("main")
	require.NoError(t, err)
	assert.Len(t, commits, 1)
}
