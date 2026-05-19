package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newDiffCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "diff <hash1> <hash2>",
		Short: "Show diff between two commits",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			if len(args[0]) < 4 || len(args[1]) < 4 {
				return fmt.Errorf("commit hashes must be at least 4 characters")
			}

			m := model.NewDiffModel(c.DiffUC, args[0], args[1])
			p := tea.NewProgram(m, tea.WithAltScreen())
			_, err := p.Run()
			return err
		},
	}
}
