package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStatusCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current staged contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			result, err := c.StatusUC.Execute()
			if err != nil {
				return err
			}

			fmt.Printf("On branch: %s\n\n", result.ActiveBranch.Name)

			if len(result.StagedFiles) == 0 {
				fmt.Println("Nothing staged for commit.")
				return nil
			}

			fmt.Println("Staged contexts:")
			for _, f := range result.StagedFiles {
				fmt.Printf("  + %s (%s) [%s]\n", f.FilePath, f.AITool, f.Hash[:8])
			}

			return nil
		},
	}
}
