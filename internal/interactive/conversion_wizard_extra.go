package interactive

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"fileforge/internal/convert"
	outputpkg "fileforge/internal/output"
	"fileforge/internal/validation"

	"charm.land/huh/v2"
)

func (w *PDFToImageWizard) Run(ctx context.Context) error {
	state := PDFToImageInput{DPI: 150}
	dpiText := "150"
	firstPageText := ""
	lastPageText := ""

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input PDF path").
				Description("You can type or drag-and-drop a file path.").
				Value(&state.InputPath).
				Validate(func(v string) error { return validatePDFPath(NormalizePath(v)) }),
			huh.NewInput().
				Title("Output folder (optional, default: ./FileForge-Output)").
				Value(&state.OutputDir).
				Validate(func(v string) error { return nil }),
			huh.NewSelect[string]().
				Title("Output image format").
				Options(
					huh.NewOption("jpg", "jpg"),
					huh.NewOption("png", "png"),
				).
				Value(&state.ToFormat),
			huh.NewInput().
				Title("DPI").
				Value(&dpiText).
				Validate(func(v string) error {
					value, err := parsePositiveInteger(v, "dpi")
					if err != nil {
						return err
					}
					state.DPI = value
					return nil
				}),
			huh.NewInput().
				Title("First page (optional)").
				Value(&firstPageText).
				Validate(func(v string) error {
					value, err := parseOptionalPositiveInteger(v, "first page")
					if err != nil {
						return err
					}
					state.FirstPage = value
					return validatePageBounds(state.FirstPage, state.LastPage)
				}),
			huh.NewInput().
				Title("Last page (optional)").
				Value(&lastPageText).
				Validate(func(v string) error {
					value, err := parseOptionalPositiveInteger(v, "last page")
					if err != nil {
						return err
					}
					state.LastPage = value
					return validatePageBounds(state.FirstPage, state.LastPage)
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)
	state.OutputDir = ResolveInteractiveOutputDir(state.OutputDir)
	generatedDir, err := outputpkg.ResolveOutputDir("", state.OutputDir, outputpkg.GeneratedPDFToImageDirName(state.InputPath))
	if err != nil {
		return ValidationError{Err: err}
	}
	state.OutputDir = generatedDir
	force, err := confirmGeneratedDirOverwrite(state.OutputDir, state.ToFormat)
	if err != nil {
		return err
	}
	state.Force = force

	var confirm bool
	summary := fmt.Sprintf("PDF to Image\nInput: %s\nOutput directory: %s\nFormat: %s\nDPI: %d\nFirst page: %s\nLast page: %s\nOverwrite: %t",
		state.InputPath, state.OutputDir, state.ToFormat, state.DPI, emptyAsAuto(state.FirstPage), emptyAsAuto(state.LastPage), state.Force)
	if err := huh.NewConfirm().Title(summary).Affirmative("Run").Negative("Cancel").Value(&confirm).Run(); err != nil {
		return cancelIfInterrupted(err)
	}
	if !confirm {
		return ErrCancelled
	}
	if err := w.Execute(ctx, state); err != nil {
		return err
	}
	return printSuccess(w.stdout, state.OutputDir)
}

func (w *PDFToImageWizard) Execute(ctx context.Context, in PDFToImageInput) error {
	err := w.service.Convert(ctx, convert.PDFToImageRequest{
		InputPath: NormalizePath(in.InputPath),
		OutputDir: NormalizePath(in.OutputDir),
		ToFormat:  in.ToFormat,
		DPI:       in.DPI,
		FirstPage: in.FirstPage,
		LastPage:  in.LastPage,
		Force:     in.Force,
	})
	if err != nil {
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

func (w *ImageToPDFWizard) Run(ctx context.Context) error {
	state := ImageToPDFInput{}
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
			huh.NewInput().
				Title("Output folder (optional, default: ./FileForge-Output)").
				Value(&state.OutputDir).
				Validate(func(v string) error { return nil }),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)
	state.OutputDir = ResolveInteractiveOutputDir(state.OutputDir)
	outputPath, err := outputpkg.ResolveOutputPath(state.InputPath, "", state.OutputDir, "", "pdf", true)
	if err != nil {
		return ValidationError{Err: err}
	}
	state.OutputPath = outputPath

	force, err := confirmOverwrite(state.OutputPath)
	if err != nil {
		return err
	}
	state.Force = force

	var confirm bool
	summary := fmt.Sprintf("Image to PDF\nInput: %s\nOutput: %s\nOverwrite: %t", state.InputPath, state.OutputPath, state.Force)
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

func (w *ImageToPDFWizard) Execute(ctx context.Context, in ImageToPDFInput) error {
	err := w.service.Convert(ctx, convert.ImageToPDFRequest{
		InputPath:  NormalizePath(in.InputPath),
		OutputPath: NormalizePath(in.OutputPath),
		Force:      in.Force,
	})
	if err != nil {
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

func parsePositiveInteger(raw string, field string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", field)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be greater than 0", field)
	}
	return value, nil
}

func parseOptionalPositiveInteger(raw string, field string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	return parsePositiveInteger(raw, field)
}

func validatePageBounds(first int, last int) error {
	if first > 0 && last > 0 && last < first {
		return fmt.Errorf("last page must be greater than or equal to first page")
	}
	return nil
}

func emptyAsAuto(page int) string {
	if page == 0 {
		return "auto"
	}
	return strconv.Itoa(page)
}
