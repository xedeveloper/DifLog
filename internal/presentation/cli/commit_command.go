package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newCommitCommand(c *Container) *cobra.Command {
	var message string

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit staged AI context files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			m := model.NewCommitModel(c.CommitUC, message)
			p := tea.NewProgram(m)
			_, err := p.Run()
			return err
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message")
	return cmd
}
