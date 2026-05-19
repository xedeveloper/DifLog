package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/infrastructure/ai"
	"github.com/xedeveloper/DifLog/internal/infrastructure/diffing"
	"github.com/xedeveloper/DifLog/internal/infrastructure/hashing"
	"github.com/xedeveloper/DifLog/internal/infrastructure/storage"
	addUC "github.com/xedeveloper/DifLog/internal/usecase/add"
	commitUC "github.com/xedeveloper/DifLog/internal/usecase/commit"
	initUC "github.com/xedeveloper/DifLog/internal/usecase/initialize"
	logUC "github.com/xedeveloper/DifLog/internal/usecase/log"
	diffUC "github.com/xedeveloper/DifLog/internal/usecase/diff"
)

func setupProject(t *testing.T) string {
	dir := t.TempDir()
	claudeMD := filepath.Join(dir, "CLAUDE.md")
	require.NoError(t, os.WriteFile(claudeMD, []byte("# Claude Context\n\nThis is a test context."), 0644))
	return dir
}

func TestInitAddCommitFlow(t *testing.T) {
	projectRoot := setupProject(t)

	initializer := storage.NewLocalStorageInitializer(projectRoot)
	contextRepo := storage.NewContextLocalRepository(initializer.DifLogDir())
	commitRepo := storage.NewCommitLocalRepository(initializer.DifLogDir())
	branchRepo := storage.NewBranchLocalRepository(initializer.DifLogDir(), initializer)
	stagingRepo := storage.NewStagingLocalRepository(initializer.DifLogDir())
	hasherSvc := hashing.NewSHA256HasherService()
	aiDetector := ai.NewAIToolDetector()

	initUseCase := initUC.NewInitializeUseCase(initializer)
	require.NoError(t, initUseCase.Execute(entity.AIToolClaude))

	addUseCase := addUC.NewAddContextUseCase(contextRepo, stagingRepo, hasherSvc, aiDetector, projectRoot)
	result, err := addUseCase.Execute([]string{}, entity.AIToolClaude)
	require.NoError(t, err)
	assert.NotEmpty(t, result.StagedFiles)
	assert.Contains(t, result.StagedFiles, "CLAUDE.md")

	commitUseCase := commitUC.NewCommitContextUseCase(commitRepo, stagingRepo, contextRepo, branchRepo, hasherSvc)
	commit, err := commitUseCase.Execute("initial context commit")
	require.NoError(t, err)
	assert.Equal(t, "initial context commit", commit.Message)
	assert.NotEmpty(t, commit.Hash)
	assert.Equal(t, "main", commit.Branch)

	logUseCase := logUC.NewViewCommitLogUseCase(commitRepo, branchRepo)
	commits, err := logUseCase.Execute()
	require.NoError(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "initial context commit", commits[0].Message)
}

func TestDiffFlow(t *testing.T) {
	projectRoot := setupProject(t)

	initializer := storage.NewLocalStorageInitializer(projectRoot)
	contextRepo := storage.NewContextLocalRepository(initializer.DifLogDir())
	commitRepo := storage.NewCommitLocalRepository(initializer.DifLogDir())
	branchRepo := storage.NewBranchLocalRepository(initializer.DifLogDir(), initializer)
	stagingRepo := storage.NewStagingLocalRepository(initializer.DifLogDir())
	hasherSvc := hashing.NewSHA256HasherService()
	aiDetector := ai.NewAIToolDetector()
	differSvc := diffing.NewTextDifferService()

	initUseCase := initUC.NewInitializeUseCase(initializer)
	require.NoError(t, initUseCase.Execute(entity.AIToolClaude))

	addUseCase := addUC.NewAddContextUseCase(contextRepo, stagingRepo, hasherSvc, aiDetector, projectRoot)
	commitUseCase := commitUC.NewCommitContextUseCase(commitRepo, stagingRepo, contextRepo, branchRepo, hasherSvc)

	_, err := addUseCase.Execute([]string{}, entity.AIToolClaude)
	require.NoError(t, err)
	commit1, err := commitUseCase.Execute("first commit")
	require.NoError(t, err)

	claudeMD := filepath.Join(projectRoot, "CLAUDE.md")
	require.NoError(t, os.WriteFile(claudeMD, []byte("# Claude Context\n\nUpdated context with more info."), 0644))

	_, err = addUseCase.Execute([]string{}, entity.AIToolClaude)
	require.NoError(t, err)
	commit2, err := commitUseCase.Execute("second commit")
	require.NoError(t, err)

	diffUseCase := diffUC.NewDiffCommitsUseCase(commitRepo, differSvc)
	diffResult, err := diffUseCase.Execute(commit1.Hash, commit2.Hash)
	require.NoError(t, err)
	assert.NotEmpty(t, diffResult.FileDiffs)
}
