package cmd

import (
	"context"
	"fmt"

	convertpkg "fileforge/internal/convert"
	"fileforge/internal/runner"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVar(&convertOpts.To, "to", "", "target format: jpg, jpeg, png, webp, pdf")
	convertCmd.Flags().StringVar(&convertOpts.Out, "out", "", "output file or directory path")
	convertCmd.Flags().IntVar(&convertOpts.DPI, "dpi", 150, "output DPI for PDF to image conversion")
	convertCmd.Flags().IntVar(&convertOpts.FirstPage, "first-page", 0, "first page to render for PDF to image conversion")
	convertCmd.Flags().IntVar(&convertOpts.LastPage, "last-page", 0, "last page to render for PDF to image conversion")
	_ = convertCmd.MarkFlagRequired("to")
	_ = convertCmd.MarkFlagRequired("out")
}

type convertOptions struct {
	To        string
	Out       string
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
		return convertpkg.NewImageConverter(run).Convert(context.Background(), convertpkg.ImageConvertRequest{
			InputPath:  input,
			OutputPath: convertOpts.Out,
			ToFormat:   convertOpts.To,
			Force:      rootOpts.Force,
		})
	case convertpkg.RoutePDFToImage:
		return convertpkg.NewPDFToImageConverter(run).Convert(context.Background(), convertpkg.PDFToImageRequest{
			InputPath: input,
			OutputDir: convertOpts.Out,
			ToFormat:  convertOpts.To,
			DPI:       convertOpts.DPI,
			FirstPage: convertOpts.FirstPage,
			LastPage:  convertOpts.LastPage,
			Force:     rootOpts.Force,
		})
	case convertpkg.RouteImageToPDF:
		return convertpkg.NewImageToPDFConverter(run).Convert(context.Background(), convertpkg.ImageToPDFRequest{
			InputPath:  input,
			OutputPath: convertOpts.Out,
			Force:      rootOpts.Force,
		})
	default:
		return convertpkg.ValidationError{Err: fmt.Errorf("unsupported conversion route")}
	}
}
