package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newLogCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "log",
		Short: "View commit history",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			m := model.NewLogModel(c.LogUC)
			p := tea.NewProgram(m, tea.WithAltScreen())
			_, err := p.Run()
			return err
		},
	}
}
