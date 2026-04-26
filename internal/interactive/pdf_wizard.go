package interactive

import (
	"context"
	"fmt"
	"io"
	"strings"

	"fileforge/internal/compress"
	outputpkg "fileforge/internal/output"
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
	stdout  io.Writer
}

type PDFMergeWizard struct {
	service pdfMergeService
	stdout  io.Writer
}

type PDFSplitWizard struct {
	service pdfSplitService
	stdout  io.Writer
}

type PDFInfoWizard struct {
	service pdfInfoService
	stdout  io.Writer
}

type PDFCompressInput struct {
	InputPath  string
	Preset     string
	OutputDir  string
	OutputPath string
	Force      bool
}

type PDFMergeInput struct {
	InputPaths []string
	OutputDir  string
	OutputPath string
	Force      bool
}

type PDFSplitInput struct {
	InputPath  string
	PageRange  string
	OutputDir  string
	OutputPath string
	Force      bool
}

type PDFInfoInput struct {
	InputPath string
}

func NewPDFCompressWizard(service pdfCompressService, stdout io.Writer) *PDFCompressWizard {
	return &PDFCompressWizard{service: service, stdout: stdout}
}

func NewPDFMergeWizard(service pdfMergeService, stdout io.Writer) *PDFMergeWizard {
	return &PDFMergeWizard{service: service, stdout: stdout}
}

func NewPDFSplitWizard(service pdfSplitService, stdout io.Writer) *PDFSplitWizard {
	return &PDFSplitWizard{service: service, stdout: stdout}
}

func NewPDFInfoWizard(service pdfInfoService, stdout io.Writer) *PDFInfoWizard {
	return &PDFInfoWizard{service: service, stdout: stdout}
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

	outputPath, err := outputpkg.ResolveOutputPath(state.InputPath, "", state.OutputDir, "-compressed", "pdf", true)
	if err != nil {
		return ValidationError{Err: err}
	}
	state.OutputPath = outputPath

	force, err := confirmOverwrite(state.OutputPath)
	if err != nil {
		return err
	}
	state.Force = force

	if _, err := confirmSummary(w.stdout, []SummaryRow{
		{Label: "Action", Value: "PDF Compression"},
		{Label: "Input", Value: state.InputPath},
		{Label: "Output", Value: state.OutputPath},
		{Label: "Preset", Value: state.Preset},
		{Label: "Overwrite", Value: boolText(state.Force)},
	}); err != nil {
		return err
	}

	if err := w.Execute(ctx, state); err != nil {
		return err
	}
	return printSuccess(w.stdout, state.OutputPath)
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
				Title("Output folder (optional, default: ./FileForge-Output)").
				Value(&state.OutputDir).
				Validate(func(v string) error { return nil }),
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
	state.OutputDir = ResolveInteractiveOutputDir(state.OutputDir)

	outputPath, err := outputpkg.ResolveOutputPath("merged.pdf", "", state.OutputDir, "", "pdf", true)
	if err != nil {
		return ValidationError{Err: err}
	}
	state.OutputPath = outputPath

	force, err := confirmOverwrite(state.OutputPath)
	if err != nil {
		return err
	}
	state.Force = force

	if _, err := confirmSummary(w.stdout, []SummaryRow{
		{Label: "Action", Value: "Merge PDFs"},
		{Label: "Inputs", Value: strings.Join(state.InputPaths, "\n")},
		{Label: "Output", Value: state.OutputPath},
		{Label: "Overwrite", Value: boolText(state.Force)},
	}); err != nil {
		return err
	}

	if err := w.Execute(ctx, state); err != nil {
		return err
	}
	return printSuccess(w.stdout, state.OutputPath)
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
	state.PageRange = strings.TrimSpace(state.PageRange)

	outputPath, err := outputpkg.ResolveOutputPath(state.InputPath, "", state.OutputDir, "-pages-"+outputpkg.SanitizeFilenamePart(state.PageRange), "pdf", true)
	if err != nil {
		return ValidationError{Err: err}
	}
	state.OutputPath = outputPath

	force, err := confirmOverwrite(state.OutputPath)
	if err != nil {
		return err
	}
	state.Force = force

	if _, err := confirmSummary(w.stdout, []SummaryRow{
		{Label: "Action", Value: "Split PDF"},
		{Label: "Input", Value: state.InputPath},
		{Label: "Pages", Value: state.PageRange},
		{Label: "Output", Value: state.OutputPath},
		{Label: "Overwrite", Value: boolText(state.Force)},
	}); err != nil {
		return err
	}

	if err := w.Execute(ctx, state); err != nil {
		return err
	}
	return printSuccess(w.stdout, state.OutputPath)
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

	if _, err := confirmSummary(w.stdout, []SummaryRow{
		{Label: "Action", Value: "PDF Info"},
		{Label: "Input", Value: state.InputPath},
	}); err != nil {
		return err
	}

	return w.Execute(ctx, state)
}

func (w *PDFInfoWizard) Execute(ctx context.Context, in PDFInfoInput) error {
	info, err := w.service.Info(ctx, pdf.InfoRequest{
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
	if w.stdout != nil {
		if info.Pages != nil {
			_, _ = fmt.Fprintf(w.stdout, "Pages: %d\n", *info.Pages)
		}
		_, _ = fmt.Fprintln(w.stdout, info.RawOutput)
	}
	return nil
}
