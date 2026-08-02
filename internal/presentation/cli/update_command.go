package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update diflog to the latest version",
		Long: `Fetch the latest DifLog release from Girhub and
		Replace the current running library. Currently supports macOS (amd64/arm64) and Linux (amd64)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Checking the latest DifLog release...")

			result, err := c.UpdateUC.Execute()
			if err != nil {
				return err
			}

			if !result.Updated {
				fmt.Printf("Already up to date (%s).\n")
				return nil
			}

			fmt.Printf("Updated diflog %s -> %s\n", result.CurrentVersion, result.LatestVersion)
			fmt.Printf("Installed to %s\n", result.InstalledPath)
			return nil
		},
	}
}
