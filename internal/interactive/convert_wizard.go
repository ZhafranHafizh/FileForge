package interactive

import (
	"context"
	"fmt"

	"fileforge/internal/convert"
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
}

type PDFToImageWizard struct {
	service pdfToImageService
}

type ImageToPDFWizard struct {
	service imageToPDFService
}

type ConvertInput struct {
	InputPath  string
	ToFormat   string
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
}

type ImageToPDFInput struct {
	InputPath  string
	OutputPath string
	Force      bool
}

func NewConvertWizard(service imageConvertService) *ConvertWizard {
	return &ConvertWizard{service: service}
}

func NewPDFToImageWizard(service pdfToImageService) *PDFToImageWizard {
	return &PDFToImageWizard{service: service}
}

func NewImageToPDFWizard(service imageToPDFService) *ImageToPDFWizard {
	return &ImageToPDFWizard{service: service}
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
				Title("Output path").
				Value(&state.OutputPath).
				Validate(func(v string) error {
					path := NormalizePath(v)
					if path == "" {
						return fmt.Errorf("output path is required")
					}
					return nil
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)
	state.OutputPath = NormalizePath(state.OutputPath)

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

	return w.Execute(ctx, state)
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
