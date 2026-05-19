package detect

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/domain/service"
)

type DetectAIToolUseCase struct {
	aiDetector  service.AIDetectorService
	projectRoot string
}

func NewDetectAIToolUseCase(aiDetector service.AIDetectorService, projectRoot string) *DetectAIToolUseCase {
	return &DetectAIToolUseCase{aiDetector: aiDetector, projectRoot: projectRoot}
}

type DetectionResult struct {
	AITool       entity.AITool
	ContextFiles []string
}

func (uc *DetectAIToolUseCase) Execute() (DetectionResult, error) {
	tool, err := uc.aiDetector.DetectAITool(uc.projectRoot)
	if err != nil {
		return DetectionResult{}, err
	}

	contextFiles := uc.aiDetector.GetContextFilePaths(uc.projectRoot, tool)

	return DetectionResult{
		AITool:       tool,
		ContextFiles: contextFiles,
	}, nil
}
