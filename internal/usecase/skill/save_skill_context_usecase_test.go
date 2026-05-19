package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/infrastructure/storage"
	"github.com/xedeveloper/DifLog/internal/infrastructure/hashing"
	skillUC "github.com/xedeveloper/DifLog/internal/usecase/skill"
)

func setupSkillContextRepo(t *testing.T) *storage.SkillContextLocalRepository {
	dir := t.TempDir()
	initializer := storage.NewLocalStorageInitializer(dir)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))
	return storage.NewSkillContextLocalRepository(initializer.DifLogDir())
}

func TestSaveSkillContextUseCase_ExecuteFromContent(t *testing.T) {
	repo := setupSkillContextRepo(t)
	hasher := hashing.NewSHA256HasherService()
	uc := skillUC.NewSaveSkillContextUseCase(repo, hasher)

	result, err := uc.ExecuteFromContent("# My Chat Context\n\nSome important discussion.", "test-session")
	require.NoError(t, err)

	assert.NotEmpty(t, result.Hash)
	assert.Equal(t, ".difLog/skill-contexts/"+result.Hash+".json", result.FilePath)

	saved, err := repo.LoadSkillContext(result.Hash)
	require.NoError(t, err)
	assert.Equal(t, "# My Chat Context\n\nSome important discussion.", saved.Content)
	assert.Equal(t, "test-session", saved.Source)
}

func TestSaveSkillContextUseCase_ExecuteFromFile(t *testing.T) {
	repo := setupSkillContextRepo(t)
	hasher := hashing.NewSHA256HasherService()
	uc := skillUC.NewSaveSkillContextUseCase(repo, hasher)

	tmpFile := filepath.Join(t.TempDir(), "context.md")
	require.NoError(t, os.WriteFile(tmpFile, []byte("# File Context\n\nContent from file."), 0644))

	result, err := uc.ExecuteFromFile(tmpFile)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Hash)

	saved, err := repo.LoadSkillContext(result.Hash)
	require.NoError(t, err)
	assert.Equal(t, "# File Context\n\nContent from file.", saved.Content)
	assert.Equal(t, tmpFile, saved.Source)
}

func TestSaveSkillContextUseCase_EmptyContentReturnsError(t *testing.T) {
	repo := setupSkillContextRepo(t)
	hasher := hashing.NewSHA256HasherService()
	uc := skillUC.NewSaveSkillContextUseCase(repo, hasher)

	_, err := uc.ExecuteFromContent("", "test")
	assert.Error(t, err)
}

func TestSaveSkillContextUseCase_SameContentProducesSameHash(t *testing.T) {
	repo := setupSkillContextRepo(t)
	hasher := hashing.NewSHA256HasherService()
	uc := skillUC.NewSaveSkillContextUseCase(repo, hasher)

	content := "# Deterministic content"
	r1, err := uc.ExecuteFromContent(content, "session-1")
	require.NoError(t, err)

	r2, err := uc.ExecuteFromContent(content, "session-2")
	require.NoError(t, err)

	assert.Equal(t, r1.Hash, r2.Hash, "same content must produce same hash")
}
