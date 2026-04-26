package cmd

import (
	"context"
	"fmt"

	compresspkg "fileforge/internal/compress"
	"fileforge/internal/runner"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(compressCmd)
	compressCmd.Flags().IntVar(&compressOpts.Quality, "quality", 80, "image quality from 1 to 100")
	compressCmd.Flags().StringVar(&compressOpts.Preset, "preset", "", "PDF compression preset: screen, ebook, printer, prepress, default")
	compressCmd.Flags().StringVar(&compressOpts.Out, "out", "", "output file path")
	_ = compressCmd.MarkFlagRequired("out")
}

type compressOptions struct {
	Quality int
	Preset  string
	Out     string
}

var compressOpts compressOptions

var compressCmd = &cobra.Command{
	Use:   "compress <input>",
	Short: "Compress supported image and PDF formats",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		run := runner.New(runner.Options{
			Verbose: rootOpts.Verbose,
			Stdout:  cmd.ErrOrStderr(),
			Stderr:  cmd.ErrOrStderr(),
		})

		if compresspkg.DetectKind(args[0]) == compresspkg.KindPDF {
			return runPDFCompression(cmd, run, args[0])
		}

		return runImageCompression(cmd, run, args[0])
	},
}

func runImageCompression(cmd *cobra.Command, run *runner.Runner, input string) error {
	compressor := compresspkg.NewImageCompressor(run)
	result, err := compressor.Compress(context.Background(), compresspkg.ImageCompressRequest{
		InputPath:  input,
		OutputPath: compressOpts.Out,
		Quality:    compressOpts.Quality,
		Force:      rootOpts.Force,
	})
	if err != nil {
		if compresspkg.IsValidationError(err) {
			return newCommandError(ExitInvalidInput, err)
		}
		if compresspkg.IsDependencyError(err) {
			return newCommandError(ExitMissingDependency, err)
		}
		return newCommandError(ExitCompressionFailed, err)
	}

	if rootOpts.Quiet {
		return nil
	}

	out := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(out, "Compressed %s -> %s\n", result.InputPath, result.OutputPath)
	_, _ = fmt.Fprintf(out, "Before: %s\n", compresspkg.FormatBytes(result.BeforeBytes))
	_, _ = fmt.Fprintf(out, "After:  %s\n", compresspkg.FormatBytes(result.AfterBytes))
	if result.PercentChange != nil {
		_, _ = fmt.Fprintf(out, "Change: %.2f%%\n", *result.PercentChange)
	}
	if result.AfterBytes > result.BeforeBytes {
		_, _ = fmt.Fprintln(out, "Warning: output PDF is larger than input.")
	}

	return nil
}

func runPDFCompression(cmd *cobra.Command, run *runner.Runner, input string) error {
	compressor := compresspkg.NewPDFCompressor(run)
	result, err := compressor.Compress(context.Background(), compresspkg.PDFCompressRequest{
		InputPath:  input,
		OutputPath: compressOpts.Out,
		Preset:     compressOpts.Preset,
		Force:      rootOpts.Force,
	})
	if err != nil {
		if compresspkg.IsValidationError(err) {
			return newCommandError(ExitInvalidInput, err)
		}
		if compresspkg.IsDependencyError(err) {
			return newCommandError(ExitMissingDependency, err)
		}
		return newCommandError(ExitCompressionFailed, err)
	}

	if rootOpts.Quiet {
		return nil
	}

	out := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(out, "Compressed %s -> %s\n", result.InputPath, result.OutputPath)
	_, _ = fmt.Fprintf(out, "Before: %s\n", compresspkg.FormatBytes(result.BeforeBytes))
	_, _ = fmt.Fprintf(out, "After:  %s\n", compresspkg.FormatBytes(result.AfterBytes))
	if result.PercentChange != nil {
		_, _ = fmt.Fprintf(out, "Change: %.2f%%\n", *result.PercentChange)
	}

	return nil
}
