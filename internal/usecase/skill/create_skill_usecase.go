package skill

import (
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

type ClaudeSkillCreator interface {
	CreateSkill(skillContent string) error
	SkillContent() string
}

type OpenCodeSkillCreator interface {
	CreateSkill(skillContent string) error
	SkillContent() string
}

type CopilotInstructionAppender interface {
	AppendDifLogInstructions(projectRoot string) error
}

type CreateSkillUseCase struct {
	claudeAdapter   ClaudeSkillCreator
	openCodeAdapter OpenCodeSkillCreator
	copilotAdapter  CopilotInstructionAppender
	projectRoot     string
}

func NewCreateSkillUseCase(
	claudeAdapter ClaudeSkillCreator,
	openCodeAdapter OpenCodeSkillCreator,
	copilotAdapter CopilotInstructionAppender,
	projectRoot string,
) *CreateSkillUseCase {
	return &CreateSkillUseCase{
		claudeAdapter:   claudeAdapter,
		openCodeAdapter: openCodeAdapter,
		copilotAdapter:  copilotAdapter,
		projectRoot:     projectRoot,
	}
}

type SkillCreationResult struct {
	AITool  entity.AITool
	Message string
}

func (uc *CreateSkillUseCase) Execute(aiTool entity.AITool) (SkillCreationResult, error) {
	switch aiTool {
	case entity.AIToolClaude:
		if err := uc.claudeAdapter.CreateSkill(uc.claudeAdapter.SkillContent()); err != nil {
			return SkillCreationResult{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create Claude skill", err)
		}
		return SkillCreationResult{
			AITool:  aiTool,
			Message: "Skill created at ~/.config/opencode/skills/diflog/SKILL.md",
		}, nil

	case entity.AIToolOpenCode:
		if err := uc.openCodeAdapter.CreateSkill(uc.openCodeAdapter.SkillContent()); err != nil {
			return SkillCreationResult{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create OpenCode skill", err)
		}
		return SkillCreationResult{
			AITool:  aiTool,
			Message: "Skill created at ~/.config/opencode/skills/diflog/SKILL.md",
		}, nil

	case entity.AIToolCopilot:
		if err := uc.copilotAdapter.AppendDifLogInstructions(uc.projectRoot); err != nil {
			return SkillCreationResult{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create Copilot instructions", err)
		}
		return SkillCreationResult{
			AITool:  aiTool,
			Message: "DifLog instructions appended to .github/copilot-instructions.md",
		}, nil

	default:
		return SkillCreationResult{}, dlerrors.NewDifLogError(dlerrors.ErrAIToolNotDetected, "unknown AI tool. Use 'claude', 'opencode', or 'copilot'", nil)
	}
}
