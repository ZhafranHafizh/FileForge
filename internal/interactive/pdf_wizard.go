package interactive

import (
	"context"
	"fmt"
	"strings"

	"fileforge/internal/compress"
	"fileforge/internal/pdf"
	"fileforge/internal/validation"

	"charm.land/huh/v2"
)

type pdfCompressService interface {
	Compress(ctx context.Context, req compress.PDFCompressRequest) (compress.PDFCompressResult, error)
}

type pdfMergeService interface {
	Merge(ctx context.Context, req pdf.MergeRequest) error
}

type pdfSplitService interface {
	Split(ctx context.Context, req pdf.SplitRequest) error
}

type pdfInfoService interface {
	Info(ctx context.Context, req pdf.InfoRequest) (pdf.InfoResult, error)
}

type PDFCompressWizard struct {
	service pdfCompressService
}

type PDFMergeWizard struct {
	service pdfMergeService
}

type PDFSplitWizard struct {
	service pdfSplitService
}

type PDFInfoWizard struct {
	service pdfInfoService
}

type PDFCompressInput struct {
	InputPath  string
	Preset     string
	OutputPath string
	Force      bool
}

type PDFMergeInput struct {
	InputPaths []string
	OutputPath string
	Force      bool
}

type PDFSplitInput struct {
	InputPath  string
	PageRange  string
	OutputPath string
	Force      bool
}

type PDFInfoInput struct {
	InputPath string
}

func NewPDFCompressWizard(service pdfCompressService) *PDFCompressWizard {
	return &PDFCompressWizard{service: service}
}

func NewPDFMergeWizard(service pdfMergeService) *PDFMergeWizard {
	return &PDFMergeWizard{service: service}
}

func NewPDFSplitWizard(service pdfSplitService) *PDFSplitWizard {
	return &PDFSplitWizard{service: service}
}

func NewPDFInfoWizard(service pdfInfoService) *PDFInfoWizard {
	return &PDFInfoWizard{service: service}
}

func (w *PDFCompressWizard) Run(ctx context.Context) error {
	state := PDFCompressInput{Preset: "ebook"}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input PDF path").
				Description("You can type or drag-and-drop a file path.").
				Value(&state.InputPath).
				Validate(func(v string) error {
					return validatePDFPath(NormalizePath(v))
				}),
			huh.NewSelect[string]().
				Title("Preset").
				Options(
					huh.NewOption("screen", "screen"),
					huh.NewOption("ebook", "ebook"),
					huh.NewOption("printer", "printer"),
					huh.NewOption("prepress", "prepress"),
					huh.NewOption("default", "default"),
				).
				Value(&state.Preset),
			huh.NewInput().
				Title("Output PDF path").
				Value(&state.OutputPath).
				Validate(func(v string) error {
					path := NormalizePath(v)
					if path == "" {
						return fmt.Errorf("output path is required")
					}
					return validation.EnsureSupportedExtension(path, []string{"pdf"})
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)
	state.OutputPath = NormalizePath(state.OutputPath)

	force, err := confirmOverwrite(state.OutputPath)
	if err != nil {
		return err
	}
	state.Force = force

	var confirm bool
	summary := fmt.Sprintf("PDF Compress\nInput: %s\nPreset: %s\nOutput: %s\nOverwrite: %t", state.InputPath, state.Preset, state.OutputPath, state.Force)
	if err := huh.NewConfirm().Title(summary).Affirmative("Run").Negative("Cancel").Value(&confirm).Run(); err != nil {
		return cancelIfInterrupted(err)
	}
	if !confirm {
		return ErrCancelled
	}

	return w.Execute(ctx, state)
}

func (w *PDFCompressWizard) Execute(ctx context.Context, in PDFCompressInput) error {
	_, err := w.service.Compress(ctx, compress.PDFCompressRequest{
		InputPath:  NormalizePath(in.InputPath),
		Preset:     in.Preset,
		OutputPath: NormalizePath(in.OutputPath),
		Force:      in.Force,
	})
	if err != nil {
		if compress.IsValidationError(err) {
			return ValidationError{Err: err}
		}
		if compress.IsDependencyError(err) {
			return DependencyError{Err: err}
		}
		return err
	}
	return nil
}

func (w *PDFMergeWizard) Run(ctx context.Context) error {
	var rawPaths string
	state := PDFMergeInput{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input PDF paths").
				Description("Enter multiple PDF paths separated by commas. Drag-and-drop paths are supported.").
				Value(&rawPaths).
				Validate(func(v string) error {
					paths, err := ParsePathList(v)
					if err != nil {
						return err
					}
					if len(paths) < 2 {
						return fmt.Errorf("at least two PDF paths are required")
					}
					for _, path := range paths {
						if err := validatePDFPath(path); err != nil {
							return err
						}
					}
					return nil
				}),
			huh.NewInput().
				Title("Output PDF path").
				Value(&state.OutputPath).
				Validate(func(v string) error {
					path := NormalizePath(v)
					if path == "" {
						return fmt.Errorf("output path is required")
					}
					return validation.EnsureSupportedExtension(path, []string{"pdf"})
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	paths, err := ParsePathList(rawPaths)
	if err != nil {
		return ValidationError{Err: err}
	}
	state.InputPaths = paths
	state.OutputPath = NormalizePath(state.OutputPath)

	force, err := confirmOverwrite(state.OutputPath)
	if err != nil {
		return err
	}
	state.Force = force

	var confirm bool
	summary := fmt.Sprintf("PDF Merge\nInputs:\n- %s\nOutput: %s\nOverwrite: %t", strings.Join(state.InputPaths, "\n- "), state.OutputPath, state.Force)
	if err := huh.NewConfirm().Title(summary).Affirmative("Run").Negative("Cancel").Value(&confirm).Run(); err != nil {
		return cancelIfInterrupted(err)
	}
	if !confirm {
		return ErrCancelled
	}

	return w.Execute(ctx, state)
}

func (w *PDFMergeWizard) Execute(ctx context.Context, in PDFMergeInput) error {
	paths := make([]string, 0, len(in.InputPaths))
	for _, path := range in.InputPaths {
		paths = append(paths, NormalizePath(path))
	}

	err := w.service.Merge(ctx, pdf.MergeRequest{
		InputPaths: paths,
		OutputPath: NormalizePath(in.OutputPath),
		Force:      in.Force,
	})
	if err != nil {
		if pdf.IsValidationError(err) {
			return ValidationError{Err: err}
		}
		if pdf.IsDependencyError(err) {
			return DependencyError{Err: err}
		}
		return err
	}
	return nil
}

func (w *PDFSplitWizard) Run(ctx context.Context) error {
	state := PDFSplitInput{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input PDF path").
				Description("You can type or drag-and-drop a file path.").
				Value(&state.InputPath).
				Validate(func(v string) error {
					return validatePDFPath(NormalizePath(v))
				}),
			huh.NewInput().
				Title("Page range").
				Description("Examples: 1, 1-5, 1,3,5, 1-3,7,10-12").
				Value(&state.PageRange).
				Validate(func(v string) error {
					_, err := validation.ParsePageRange(v)
					return err
				}),
			huh.NewInput().
				Title("Output PDF path").
				Value(&state.OutputPath).
				Validate(func(v string) error {
					path := NormalizePath(v)
					if path == "" {
						return fmt.Errorf("output path is required")
					}
					return validation.EnsureSupportedExtension(path, []string{"pdf"})
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)
	state.OutputPath = NormalizePath(state.OutputPath)
	state.PageRange = strings.TrimSpace(state.PageRange)

	force, err := confirmOverwrite(state.OutputPath)
	if err != nil {
		return err
	}
	state.Force = force

	var confirm bool
	summary := fmt.Sprintf("PDF Split\nInput: %s\nPages: %s\nOutput: %s\nOverwrite: %t", state.InputPath, state.PageRange, state.OutputPath, state.Force)
	if err := huh.NewConfirm().Title(summary).Affirmative("Run").Negative("Cancel").Value(&confirm).Run(); err != nil {
		return cancelIfInterrupted(err)
	}
	if !confirm {
		return ErrCancelled
	}

	return w.Execute(ctx, state)
}

func (w *PDFSplitWizard) Execute(ctx context.Context, in PDFSplitInput) error {
	err := w.service.Split(ctx, pdf.SplitRequest{
		InputPath:  NormalizePath(in.InputPath),
		PageRange:  strings.TrimSpace(in.PageRange),
		OutputPath: NormalizePath(in.OutputPath),
		Force:      in.Force,
	})
	if err != nil {
		if pdf.IsValidationError(err) {
			return ValidationError{Err: err}
		}
		if pdf.IsDependencyError(err) {
			return DependencyError{Err: err}
		}
		return err
	}
	return nil
}

func (w *PDFInfoWizard) Run(ctx context.Context) error {
	state := PDFInfoInput{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input PDF path").
				Description("You can type or drag-and-drop a file path.").
				Value(&state.InputPath).
				Validate(func(v string) error {
					return validatePDFPath(NormalizePath(v))
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)

	var confirm bool
	summary := fmt.Sprintf("PDF Info\nInput: %s", state.InputPath)
	if err := huh.NewConfirm().Title(summary).Affirmative("Run").Negative("Cancel").Value(&confirm).Run(); err != nil {
		return cancelIfInterrupted(err)
	}
	if !confirm {
		return ErrCancelled
	}

	return w.Execute(ctx, state)
}

func (w *PDFInfoWizard) Execute(ctx context.Context, in PDFInfoInput) error {
	_, err := w.service.Info(ctx, pdf.InfoRequest{
		InputPath: NormalizePath(in.InputPath),
	})
	if err != nil {
		if pdf.IsValidationError(err) {
			return ValidationError{Err: err}
		}
		if pdf.IsDependencyError(err) {
			return DependencyError{Err: err}
		}
		return err
	}
	return nil
}
