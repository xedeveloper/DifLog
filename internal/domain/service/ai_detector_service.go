package service

import "github.com/xedeveloper/DifLog/internal/domain/entity"

type AIDetectorService interface {
	DetectAITool(projectRoot string) (entity.AITool, error)
	DetectAllAITools(projectRoot string) ([]entity.AITool, error)
	GetContextFilePaths(projectRoot string, tool entity.AITool) []string
}
