package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newInitCommand(c *Container) *cobra.Command {
	var aiToolFlag string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize DifLog in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			aiTool := entity.AITool(aiToolFlag)

			if aiTool == "" {
				detectedTools, _ := c.AIDetector.DetectAllAITools(c.ProjectRoot)

				selectionModel := model.NewAISelectionModel(detectedTools)
				p := tea.NewProgram(selectionModel)
				finalModel, err := p.Run()
				if err != nil {
					return fmt.Errorf("AI selection failed: %w", err)
				}

				result, ok := finalModel.(model.AISelectionModel)
				if !ok || result.WasCancelled() {
					return fmt.Errorf("initialization cancelled")
				}

				aiTool = result.SelectedTool()
				if aiTool == entity.AIToolUnknown {
					return fmt.Errorf("no AI tool selected")
				}
			}

			initModel := model.NewInitializeModel(c.InitUC, aiTool)
			p := tea.NewProgram(initModel)
			_, err := p.Run()
			return err
		},
	}

	cmd.Flags().StringVarP(&aiToolFlag, "ai", "a", "", "AI tool to use (claude|opencode|copilot) — skips interactive selection")
	return cmd
}
