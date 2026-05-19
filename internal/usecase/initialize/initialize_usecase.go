package initialize

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
)

type Initializer interface {
	Initialize(aiTool entity.AITool) error
}

type InitializeUseCase struct {
	initializer Initializer
}

func NewInitializeUseCase(initializer Initializer) *InitializeUseCase {
	return &InitializeUseCase{initializer: initializer}
}

func (uc *InitializeUseCase) Execute(aiTool entity.AITool) error {
	return uc.initializer.Initialize(aiTool)
}
