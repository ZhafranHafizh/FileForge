package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print FileForge version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		if rootOpts.Quiet {
			return nil
		}

		_, err := fmt.Fprintf(cmd.OutOrStdout(), "fileforge %s\ncommit: %s\nbuilt: %s\n", version, commit, date)
		return err
	},
}
