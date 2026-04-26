package interactive

import (
	"context"
	"errors"
	"fmt"
	"io"

	"fileforge/internal/compress"
	"fileforge/internal/convert"
	"fileforge/internal/pdf"
	"fileforge/internal/runner"

	"charm.land/huh/v2"
)

var ErrCancelled = errors.New("interactive flow cancelled")

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

type DependencyError struct {
	Err error
}

func (e DependencyError) Error() string {
	return e.Err.Error()
}

func (e DependencyError) Unwrap() error {
	return e.Err
}

type Options struct {
	Runner       *runner.Runner
	Doctor       func(context.Context, io.Writer) error
	Stdout       io.Writer
	ImageForce   bool
	AccessibleUI bool
}

type App struct {
	convertWizard  *ConvertWizard
	compressWizard *CompressWizard
	pdfCompress    *PDFCompressWizard
	pdfMerge       *PDFMergeWizard
	pdfSplit       *PDFSplitWizard
	pdfInfo        *PDFInfoWizard
	doctor         func(context.Context, io.Writer) error
	stdout         io.Writer
	accessibleUI   bool
}

func NewApp(opts Options) *App {
	return &App{
		convertWizard:  NewConvertWizard(convert.NewImageConverter(opts.Runner)),
		compressWizard: NewCompressWizard(compress.NewImageCompressor(opts.Runner)),
		pdfCompress:    NewPDFCompressWizard(compress.NewPDFCompressor(opts.Runner)),
		pdfMerge:       NewPDFMergeWizard(pdf.NewMerger(opts.Runner)),
		pdfSplit:       NewPDFSplitWizard(pdf.NewSplitter(opts.Runner)),
		pdfInfo:        NewPDFInfoWizard(pdf.NewInspector(opts.Runner)),
		doctor:         opts.Doctor,
		stdout:         opts.Stdout,
		accessibleUI:   opts.AccessibleUI,
	}
}

func (a *App) Run(ctx context.Context) error {
	for {
		choice, err := a.promptMainMenu()
		if err != nil {
			return cancelIfInterrupted(err)
		}

		switch choice {
		case "convert_image":
			if err := a.convertWizard.Run(ctx); err != nil && !IsCancelled(err) {
				return err
			}
		case "compress_image":
			if err := a.compressWizard.Run(ctx); err != nil && !IsCancelled(err) {
				return err
			}
		case "doctor":
			if a.doctor != nil {
				if err := a.doctor(ctx, a.stdout); err != nil {
					return err
				}
			}
		case "pdf_tools":
			pdfChoice, err := a.promptPDFMenu()
			if err != nil {
				return cancelIfInterrupted(err)
			}
			if err := a.runPDFChoice(ctx, pdfChoice); err != nil && !IsCancelled(err) {
				return err
			}
		case "pdf_compress":
			if err := a.pdfCompress.Run(ctx); err != nil && !IsCancelled(err) {
				return err
			}
		case "pdf_merge":
			if err := a.pdfMerge.Run(ctx); err != nil && !IsCancelled(err) {
				return err
			}
		case "pdf_split":
			if err := a.pdfSplit.Run(ctx); err != nil && !IsCancelled(err) {
				return err
			}
		case "pdf_info":
			if err := a.pdfInfo.Run(ctx); err != nil && !IsCancelled(err) {
				return err
			}
		case "pdf_to_image", "image_to_pdf", "ocr", "office":
			if err := showNotImplemented("This feature is not implemented yet."); err != nil {
				return err
			}
		case "exit":
			return nil
		}
	}
}

func (a *App) promptMainMenu() (string, error) {
	var choice string
	err := huh.NewSelect[string]().
		Title("Choose an action").
		Options(
			huh.NewOption("Convert file", "convert_image"),
			huh.NewOption("Compress file", "compress_image"),
			huh.NewOption("PDF tools", "pdf_tools"),
			huh.NewOption("PDF to Image", "pdf_to_image"),
			huh.NewOption("Image to PDF", "image_to_pdf"),
			huh.NewOption("PDF Compress", "pdf_compress"),
			huh.NewOption("PDF Merge", "pdf_merge"),
			huh.NewOption("PDF Split", "pdf_split"),
			huh.NewOption("PDF Info", "pdf_info"),
			huh.NewOption("Doctor / Check dependencies", "doctor"),
			huh.NewOption("OCR", "ocr"),
			huh.NewOption("Office conversion", "office"),
			huh.NewOption("Exit", "exit"),
		).
		Value(&choice).
		Run()
	return choice, err
}

func (a *App) promptPDFMenu() (string, error) {
	var choice string
	err := huh.NewSelect[string]().
		Title("Choose a PDF action").
		Options(
			huh.NewOption("PDF Compress", "pdf_compress"),
			huh.NewOption("PDF Merge", "pdf_merge"),
			huh.NewOption("PDF Split", "pdf_split"),
			huh.NewOption("PDF Info", "pdf_info"),
			huh.NewOption("PDF to Image", "pdf_to_image"),
			huh.NewOption("Image to PDF", "image_to_pdf"),
			huh.NewOption("Back", "back"),
		).
		Value(&choice).
		Run()
	return choice, err
}

func (a *App) runPDFChoice(ctx context.Context, choice string) error {
	switch choice {
	case "pdf_compress":
		return a.pdfCompress.Run(ctx)
	case "pdf_merge":
		return a.pdfMerge.Run(ctx)
	case "pdf_split":
		return a.pdfSplit.Run(ctx)
	case "pdf_info":
		return a.pdfInfo.Run(ctx)
	case "pdf_to_image", "image_to_pdf":
		return showNotImplemented("This feature is not implemented yet.")
	case "back":
		return nil
	default:
		return nil
	}
}

func showNotImplemented(message string) error {
	var confirm bool
	return cancelIfInterrupted(
		huh.NewConfirm().
			Title(message).
			Affirmative("OK").
			Negative("OK").
			Value(&confirm).
			Run(),
	)
}

func IsCancelled(err error) bool {
	return errors.Is(err, ErrCancelled)
}

func IsValidationError(err error) bool {
	var target ValidationError
	return errors.As(err, &target)
}

func IsDependencyError(err error) bool {
	var target DependencyError
	return errors.As(err, &target)
}

func cancelIfInterrupted(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %v", ErrCancelled, err)
}
