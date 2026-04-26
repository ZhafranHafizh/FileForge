package cmd

import (
	"context"
	"fmt"

	outputpkg "fileforge/internal/output"
	pdfpkg "fileforge/internal/pdf"
	"fileforge/internal/runner"

	"github.com/spf13/cobra"
)

type pdfOptions struct {
	Out       string
	OutputDir string
	Pages     string
}

var pdfOpts pdfOptions

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "PDF utilities",
}

var pdfMergeCmd = &cobra.Command{
	Use:   "merge <input1.pdf> <input2.pdf> [input3.pdf...]",
	Short: "Merge multiple PDFs into one file",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		run := runner.New(runner.Options{
			Verbose: rootOpts.Verbose,
			Stdout:  cmd.ErrOrStderr(),
			Stderr:  cmd.ErrOrStderr(),
		})
		outputPath, err := outputpkg.ResolveOutputPath("merged.pdf", pdfOpts.Out, pdfOpts.OutputDir, "", "pdf", rootOpts.Force)
		if err != nil {
			return newCommandError(ExitInvalidInput, err)
		}

		merger := pdfpkg.NewMerger(run)
		err = merger.Merge(context.Background(), pdfpkg.MergeRequest{
			InputPaths: args,
			OutputPath: outputPath,
			Force:      rootOpts.Force,
		})
		if err != nil {
			if pdfpkg.IsValidationError(err) {
				return newCommandError(ExitInvalidInput, err)
			}
			if pdfpkg.IsDependencyError(err) {
				return newCommandError(ExitMissingDependency, err)
			}
			return newCommandError(ExitConversionFailed, err)
		}

		out := cmd.OutOrStdout()
		if !rootOpts.Quiet {
			_, _ = fmt.Fprintf(out, "Merged %d PDFs into %s\n", len(args), outputPath)
		}
		return outputpkg.PrintSuccess(out, outputPath)
	},
}

var pdfSplitCmd = &cobra.Command{
	Use:   "split <input.pdf>",
	Short: "Extract a page range from a PDF",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		run := runner.New(runner.Options{
			Verbose: rootOpts.Verbose,
			Stdout:  cmd.ErrOrStderr(),
			Stderr:  cmd.ErrOrStderr(),
		})
		outputPath, err := outputpkg.ResolveOutputPath(args[0], pdfOpts.Out, pdfOpts.OutputDir, "", "pdf", rootOpts.Force)
		if pdfOpts.OutputDir != "" && pdfOpts.Out == "" {
			outputPath, err = outputpkg.ResolveOutputPath(args[0], "", pdfOpts.OutputDir, "-pages-"+outputpkg.SanitizeFilenamePart(pdfOpts.Pages), "pdf", rootOpts.Force)
		}
		if err != nil {
			return newCommandError(ExitInvalidInput, err)
		}

		splitter := pdfpkg.NewSplitter(run)
		err = splitter.Split(context.Background(), pdfpkg.SplitRequest{
			InputPath:  args[0],
			PageRange:  pdfOpts.Pages,
			OutputPath: outputPath,
			Force:      rootOpts.Force,
		})
		if err != nil {
			if pdfpkg.IsValidationError(err) {
				return newCommandError(ExitInvalidInput, err)
			}
			if pdfpkg.IsDependencyError(err) {
				return newCommandError(ExitMissingDependency, err)
			}
			return newCommandError(ExitConversionFailed, err)
		}

		out := cmd.OutOrStdout()
		if !rootOpts.Quiet {
			_, _ = fmt.Fprintf(out, "Extracted pages %s to %s\n", pdfOpts.Pages, outputPath)
		}
		return outputpkg.PrintSuccess(out, outputPath)
	},
}

var pdfInfoCmd = &cobra.Command{
	Use:   "info <input.pdf>",
	Short: "Show PDF metadata using pdfinfo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		run := runner.New(runner.Options{
			Verbose: rootOpts.Verbose,
			Stdout:  cmd.ErrOrStderr(),
			Stderr:  cmd.ErrOrStderr(),
		})

		inspector := pdfpkg.NewInspector(run)
		info, err := inspector.Info(context.Background(), pdfpkg.InfoRequest{
			InputPath: args[0],
		})
		if err != nil {
			if pdfpkg.IsValidationError(err) {
				return newCommandError(ExitInvalidInput, err)
			}
			if pdfpkg.IsDependencyError(err) {
				return newCommandError(ExitMissingDependency, err)
			}
			return newCommandError(ExitGeneralError, err)
		}

		if rootOpts.Quiet {
			return nil
		}

		out := cmd.OutOrStdout()
		if info.Pages != nil {
			_, _ = fmt.Fprintf(out, "Pages: %d\n", *info.Pages)
		}
		_, _ = fmt.Fprint(out, info.RawOutput)
		if len(info.RawOutput) > 0 && info.RawOutput[len(info.RawOutput)-1] != '\n' {
			_, _ = fmt.Fprintln(out)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pdfCmd)

	pdfCmd.AddCommand(pdfMergeCmd)
	pdfMergeCmd.Flags().StringVar(&pdfOpts.Out, "out", "", "output PDF path")
	pdfMergeCmd.Flags().StringVar(&pdfOpts.OutputDir, "output-dir", "", "output directory for generated files")
	pdfMergeCmd.Flags().StringVar(&pdfOpts.OutputDir, "output", "", "alias for --output-dir")

	pdfCmd.AddCommand(pdfSplitCmd)
	pdfSplitCmd.Flags().StringVar(&pdfOpts.Out, "out", "", "output PDF path")
	pdfSplitCmd.Flags().StringVar(&pdfOpts.OutputDir, "output-dir", "", "output directory for generated files")
	pdfSplitCmd.Flags().StringVar(&pdfOpts.OutputDir, "output", "", "alias for --output-dir")
	pdfSplitCmd.Flags().StringVar(&pdfOpts.Pages, "pages", "", "page range to extract")
	_ = pdfSplitCmd.MarkFlagRequired("pages")

	pdfCmd.AddCommand(pdfInfoCmd)
}
