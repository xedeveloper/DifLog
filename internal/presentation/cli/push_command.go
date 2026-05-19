package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newPushCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "push [remote-url]",
		Short: "Push committed contexts to remote GitHub repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			remoteURL := ""
			if len(args) > 0 {
				remoteURL = args[0]
			}

			m := model.NewPushModel(c.PushUC, remoteURL)
			p := tea.NewProgram(m)
			_, err := p.Run()
			return err
		},
	}
}

func newPullCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "pull [remote-url]",
		Short: "Pull contexts from remote GitHub repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			remoteURL := ""
			if len(args) > 0 {
				remoteURL = args[0]
			}

			m := model.NewPullModel(c.PullUC, remoteURL)
			p := tea.NewProgram(m)
			_, err := p.Run()
			return err
		},
	}
}
