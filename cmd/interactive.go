package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start the interactive terminal wizard",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractiveApp(cmd)
	},
}
