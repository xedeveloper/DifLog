package initialize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/infrastructure/storage"
	"github.com/xedeveloper/DifLog/internal/usecase/initialize"
)

func TestInitializeUseCase_Execute(t *testing.T) {
	dir := t.TempDir()
	initializer := storage.NewLocalStorageInitializer(dir)
	uc := initialize.NewInitializeUseCase(initializer)

	err := uc.Execute(entity.AIToolClaude)
	require.NoError(t, err)
	assert.True(t, initializer.IsInitialized())
}

func TestInitializeUseCase_AlreadyInitialized(t *testing.T) {
	dir := t.TempDir()
	initializer := storage.NewLocalStorageInitializer(dir)
	uc := initialize.NewInitializeUseCase(initializer)

	require.NoError(t, uc.Execute(entity.AIToolClaude))
	err := uc.Execute(entity.AIToolClaude)
	assert.Error(t, err)
}
