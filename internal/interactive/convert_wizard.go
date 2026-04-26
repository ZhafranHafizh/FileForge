package interactive

import (
	"context"
	"fmt"
	"io"

	"fileforge/internal/convert"
	outputpkg "fileforge/internal/output"
	"fileforge/internal/validation"

	"charm.land/huh/v2"
)

type imageConvertService interface {
	Convert(ctx context.Context, req convert.ImageConvertRequest) error
}

type pdfToImageService interface {
	Convert(ctx context.Context, req convert.PDFToImageRequest) error
}

type imageToPDFService interface {
	Convert(ctx context.Context, req convert.ImageToPDFRequest) error
}

type ConvertWizard struct {
	service imageConvertService
	stdout  io.Writer
}

type PDFToImageWizard struct {
	service pdfToImageService
	stdout  io.Writer
}

type ImageToPDFWizard struct {
	service imageToPDFService
	stdout  io.Writer
}

type ConvertInput struct {
	InputPath  string
	ToFormat   string
	OutputDir  string
	OutputPath string
	Force      bool
}

type PDFToImageInput struct {
	InputPath string
	OutputDir string
	ToFormat  string
	DPI       int
	FirstPage int
	LastPage  int
	Force     bool
}

type ImageToPDFInput struct {
	InputPath  string
	OutputDir  string
	OutputPath string
	Force      bool
}

func NewConvertWizard(service imageConvertService, stdout io.Writer) *ConvertWizard {
	return &ConvertWizard{service: service, stdout: stdout}
}

func NewPDFToImageWizard(service pdfToImageService, stdout io.Writer) *PDFToImageWizard {
	return &PDFToImageWizard{service: service, stdout: stdout}
}

func NewImageToPDFWizard(service imageToPDFService, stdout io.Writer) *ImageToPDFWizard {
	return &ImageToPDFWizard{service: service, stdout: stdout}
}

func (w *ConvertWizard) Run(ctx context.Context) error {
	var state ConvertInput

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input image path").
				Description("You can type or drag-and-drop a file path.").
				Value(&state.InputPath).
				Validate(func(v string) error {
					path := NormalizePath(v)
					if path == "" {
						return fmt.Errorf("input path is required")
					}
					if err := validation.EnsureInputFile(path); err != nil {
						return err
					}
					return validation.EnsureSupportedExtension(path, []string{"jpg", "jpeg", "png", "webp"})
				}),
			huh.NewSelect[string]().
				Title("Target format").
				Options(
					huh.NewOption("jpg", "jpg"),
					huh.NewOption("jpeg", "jpeg"),
					huh.NewOption("png", "png"),
					huh.NewOption("webp", "webp"),
				).
				Value(&state.ToFormat),
			huh.NewInput().
				Title("Output folder").
				Value(&state.OutputDir).
				Validate(func(v string) error {
					path := NormalizePath(v)
					if path == "" {
						return fmt.Errorf("output folder is required")
					}
					return nil
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)
	state.OutputDir = NormalizePath(state.OutputDir)

	outputPath, err := generatedImageConvertOutput(state.InputPath, state.OutputDir, state.ToFormat)
	if err != nil {
		return ValidationError{Err: err}
	}
	state.OutputPath = outputPath

	if force, err := confirmOverwrite(state.OutputPath); err != nil {
		return err
	} else {
		state.Force = force
	}

	var confirm bool
	summary := fmt.Sprintf("Convert\nInput: %s\nFormat: %s\nOutput: %s\nOverwrite: %t", state.InputPath, state.ToFormat, state.OutputPath, state.Force)
	if err := huh.NewConfirm().Title(summary).Affirmative("Run").Negative("Cancel").Value(&confirm).Run(); err != nil {
		return cancelIfInterrupted(err)
	}
	if !confirm {
		return ErrCancelled
	}

	if err := w.Execute(ctx, state); err != nil {
		return err
	}
	return printSuccess(w.stdout, state.OutputPath)
}

func (w *ConvertWizard) Execute(ctx context.Context, in ConvertInput) error {
	req := convert.ImageConvertRequest{
		InputPath:  NormalizePath(in.InputPath),
		OutputPath: NormalizePath(in.OutputPath),
		ToFormat:   in.ToFormat,
		Force:      in.Force,
	}
	if err := w.service.Convert(ctx, req); err != nil {
		if convert.IsValidationError(err) {
			return ValidationError{Err: err}
		}
		if convert.IsDependencyError(err) {
			return DependencyError{Err: err}
		}
		return err
	}
	return nil
}

func generatedImageConvertOutput(inputPath string, outputDir string, toFormat string) (string, error) {
	return outputpkg.ResolveOutputPath(inputPath, "", outputDir, "", toFormat, true)
}
