package cmd

import (
	"context"

	convertpkg "fileforge/internal/convert"
	"fileforge/internal/runner"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVar(&convertOpts.To, "to", "", "target image format: jpg, jpeg, png, webp")
	convertCmd.Flags().StringVar(&convertOpts.Out, "out", "", "output file path")
	_ = convertCmd.MarkFlagRequired("to")
	_ = convertCmd.MarkFlagRequired("out")
}

type convertOptions struct {
	To  string
	Out string
}

var convertOpts convertOptions

var convertCmd = &cobra.Command{
	Use:   "convert <input>",
	Short: "Convert images between supported formats",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		run := runner.New(runner.Options{
			Verbose: rootOpts.Verbose,
			Stdout:  cmd.ErrOrStderr(),
			Stderr:  cmd.ErrOrStderr(),
		})

		converter := convertpkg.NewImageConverter(run)
		err := converter.Convert(context.Background(), convertpkg.ImageConvertRequest{
			InputPath:  args[0],
			OutputPath: convertOpts.Out,
			ToFormat:   convertOpts.To,
			Force:      rootOpts.Force,
		})
		if err == nil {
			return nil
		}

		if convertpkg.IsValidationError(err) {
			return newCommandError(ExitInvalidInput, err)
		}
		if convertpkg.IsDependencyError(err) {
			return newCommandError(ExitMissingDependency, err)
		}

		return newCommandError(ExitConversionFailed, err)
	},
}
