package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/infrastructure/storage"
	"time"
)

func TestSkillContextLocalRepository_SaveAndLoad(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewSkillContextLocalRepository(initializer.DifLogDir())

	ctx := entity.SkillContext{
		Hash:      "abc123def456",
		Timestamp: time.Now(),
		Content:   "# Chat Context\n\nThis is a test context.",
		Source:    "test-session",
	}

	require.NoError(t, repo.SaveSkillContext(ctx))

	loaded, err := repo.LoadSkillContext("abc123def456")
	require.NoError(t, err)
	assert.Equal(t, ctx.Hash, loaded.Hash)
	assert.Equal(t, ctx.Content, loaded.Content)
	assert.Equal(t, ctx.Source, loaded.Source)
}

func TestSkillContextLocalRepository_ListSkillContexts(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewSkillContextLocalRepository(initializer.DifLogDir())

	empty, err := repo.ListSkillContexts()
	require.NoError(t, err)
	assert.Empty(t, empty)

	for _, hash := range []string{"hash001", "hash002", "hash003"} {
		require.NoError(t, repo.SaveSkillContext(entity.SkillContext{
			Hash:      hash,
			Timestamp: time.Now(),
			Content:   "content for " + hash,
			Source:    "test",
		}))
	}

	list, err := repo.ListSkillContexts()
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestSkillContextLocalRepository_LoadNonExistent(t *testing.T) {
	initializer, _ := setupInitializer(t)
	require.NoError(t, initializer.Initialize(entity.AIToolClaude))

	repo := storage.NewSkillContextLocalRepository(initializer.DifLogDir())

	_, err := repo.LoadSkillContext("doesnotexist")
	assert.Error(t, err)
}
