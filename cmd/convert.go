package cmd

import (
	"context"
	"fmt"

	convertpkg "fileforge/internal/convert"
	outputpkg "fileforge/internal/output"
	"fileforge/internal/runner"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVar(&convertOpts.To, "to", "", "target format: jpg, jpeg, png, webp, pdf")
	convertCmd.Flags().StringVar(&convertOpts.Out, "out", "", "output file or directory path")
	convertCmd.Flags().StringVar(&convertOpts.OutputDir, "output-dir", "", "output directory for generated files")
	convertCmd.Flags().StringVar(&convertOpts.OutputDir, "output", "", "alias for --output-dir")
	convertCmd.Flags().IntVar(&convertOpts.DPI, "dpi", 150, "output DPI for PDF to image conversion")
	convertCmd.Flags().IntVar(&convertOpts.FirstPage, "first-page", 0, "first page to render for PDF to image conversion")
	convertCmd.Flags().IntVar(&convertOpts.LastPage, "last-page", 0, "last page to render for PDF to image conversion")
	_ = convertCmd.MarkFlagRequired("to")
}

type convertOptions struct {
	To        string
	Out       string
	OutputDir string
	DPI       int
	FirstPage int
	LastPage  int
}

var convertOpts convertOptions

var convertCmd = &cobra.Command{
	Use:   "convert <input>",
	Short: "Convert supported image and PDF formats",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		run := runner.New(runner.Options{
			Verbose: rootOpts.Verbose,
			Stdout:  cmd.ErrOrStderr(),
			Stderr:  cmd.ErrOrStderr(),
		})

		err := runConvert(cmd, run, args[0])
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

func runConvert(cmd *cobra.Command, run *runner.Runner, input string) error {
	route, err := convertpkg.DetectRoute(input, convertOpts.To)
	if err != nil {
		return err
	}

	switch route {
	case convertpkg.RouteImageToImage:
		outputPath, err := outputpkg.ResolveOutputPath(input, convertOpts.Out, convertOpts.OutputDir, "", convertOpts.To, rootOpts.Force)
		if err != nil {
			return convertpkg.ValidationError{Err: err}
		}
		err = convertpkg.NewImageConverter(run).Convert(context.Background(), convertpkg.ImageConvertRequest{
			InputPath:  input,
			OutputPath: outputPath,
			ToFormat:   convertOpts.To,
			Force:      rootOpts.Force,
		})
		if err == nil {
			return outputpkg.PrintSuccess(cmd.OutOrStdout(), outputPath)
		}
		return err
	case convertpkg.RoutePDFToImage:
		outputDir, err := outputpkg.ResolveOutputDir(convertOpts.Out, convertOpts.OutputDir, outputpkg.GeneratedPDFToImageDirName(input))
		if err != nil {
			return convertpkg.ValidationError{Err: err}
		}
		err = convertpkg.NewPDFToImageConverter(run).Convert(context.Background(), convertpkg.PDFToImageRequest{
			InputPath: input,
			OutputDir: outputDir,
			ToFormat:  convertOpts.To,
			DPI:       convertOpts.DPI,
			FirstPage: convertOpts.FirstPage,
			LastPage:  convertOpts.LastPage,
			Force:     rootOpts.Force,
		})
		if err == nil {
			return outputpkg.PrintSuccess(cmd.OutOrStdout(), outputDir)
		}
		return err
	case convertpkg.RouteImageToPDF:
		outputPath, err := outputpkg.ResolveOutputPath(input, convertOpts.Out, convertOpts.OutputDir, "", "pdf", rootOpts.Force)
		if err != nil {
			return convertpkg.ValidationError{Err: err}
		}
		err = convertpkg.NewImageToPDFConverter(run).Convert(context.Background(), convertpkg.ImageToPDFRequest{
			InputPath:  input,
			OutputPath: outputPath,
			Force:      rootOpts.Force,
		})
		if err == nil {
			return outputpkg.PrintSuccess(cmd.OutOrStdout(), outputPath)
		}
		return err
	default:
		return convertpkg.ValidationError{Err: fmt.Errorf("unsupported conversion route")}
	}
}
