package main

import (
	"fmt"
	"os"

	"github.com/xedeveloper/DifLog/internal/presentation/cli"
)

func main() {
	rootCmd := cli.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
