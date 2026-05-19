package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newBranchCommand(c *Container) *cobra.Command {
	branchCmd := &cobra.Command{
		Use:   "branch",
		Short: "Manage branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}
			m := model.NewBranchModel(c.ListBranchesUC)
			p := tea.NewProgram(m, tea.WithAltScreen())
			_, err := p.Run()
			return err
		},
	}

	branchCmd.AddCommand(
		&cobra.Command{
			Use:   "create <name>",
			Short: "Create a new branch",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := c.RequireInitialized(); err != nil {
					return err
				}
				branch, err := c.CreateBranchUC.Execute(args[0])
				if err != nil {
					return err
				}
				fmt.Printf("Created branch: %s\n", branch.Name)
				return nil
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all branches",
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := c.RequireInitialized(); err != nil {
					return err
				}
				m := model.NewBranchModel(c.ListBranchesUC)
				p := tea.NewProgram(m, tea.WithAltScreen())
				_, err := p.Run()
				return err
			},
		},
		&cobra.Command{
			Use:   "switch <name>",
			Short: "Switch to a branch",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := c.RequireInitialized(); err != nil {
					return err
				}
				if err := c.SwitchBranchUC.Execute(args[0]); err != nil {
					return err
				}
				fmt.Printf("Switched to branch: %s\n", args[0])
				return nil
			},
		},
	)

	return branchCmd
}
