package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newCheckoutCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "checkout <hash>",
		Short: "Restore context to a specific commit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			m := model.NewCheckoutModel(c.CheckoutUC, args[0])
			p := tea.NewProgram(m)
			_, err := p.Run()
			return err
		},
	}
}
