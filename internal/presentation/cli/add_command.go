package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newAddCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "add [file...]",
		Short: "Stage AI context files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			aiTool := c.LoadAITool()
			m := model.NewAddModel(c.AddUC, args, aiTool)
			p := tea.NewProgram(m)
			_, err := p.Run()
			return err
		},
	}
}
