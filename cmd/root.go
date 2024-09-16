package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xr0-org/progstack-ssg/pkg/ssg"
)

var rootCmd = &cobra.Command{
	Use:   "progstack-ssg [source] [target] [theme]",
	Short: "Generate a site from Markdown files and directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf(
				"must provide source, target and theme directories",
			)
		}
		return ssg.Generate(args[0], args[1], args[2])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
