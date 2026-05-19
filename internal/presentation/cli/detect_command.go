package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/presentation/tui/model"
)

func newDetectCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect AI tool in the current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := model.NewDetectModel(c.DetectUC)
			p := tea.NewProgram(m)
			_, err := p.Run()
			return err
		},
	}
}
