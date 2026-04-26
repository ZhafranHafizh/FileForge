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

func (e ValidationError) Error() string { return e.Err.Error() }
func (e ValidationError) Unwrap() error { return e.Err }

type DependencyError struct {
	Err error
}

func (e DependencyError) Error() string { return e.Err.Error() }
func (e DependencyError) Unwrap() error { return e.Err }

type Options struct {
	Runner       *runner.Runner
	Doctor       func(context.Context, io.Writer) error
	Stdout       io.Writer
	ImageForce   bool
	AccessibleUI bool
}

type App struct {
	convertWizard  *ConvertWizard
	pdfToImage     *PDFToImageWizard
	imageToPDF     *ImageToPDFWizard
	compressWizard *CompressWizard
	pdfCompress    *PDFCompressWizard
	pdfMerge       *PDFMergeWizard
	pdfSplit       *PDFSplitWizard
	pdfInfo        *PDFInfoWizard
	doctor         func(context.Context, io.Writer) error
	stdout         io.Writer
	accessibleUI   bool
}

type menuChoice struct {
	Label string
	Value string
}

func NewApp(opts Options) *App {
	return &App{
		convertWizard:  NewConvertWizard(convert.NewImageConverter(opts.Runner), opts.Stdout),
		pdfToImage:     NewPDFToImageWizard(convert.NewPDFToImageConverter(opts.Runner), opts.Stdout),
		imageToPDF:     NewImageToPDFWizard(convert.NewImageToPDFConverter(opts.Runner), opts.Stdout),
		compressWizard: NewCompressWizard(compress.NewImageCompressor(opts.Runner), opts.Stdout),
		pdfCompress:    NewPDFCompressWizard(compress.NewPDFCompressor(opts.Runner), opts.Stdout),
		pdfMerge:       NewPDFMergeWizard(pdf.NewMerger(opts.Runner), opts.Stdout),
		pdfSplit:       NewPDFSplitWizard(pdf.NewSplitter(opts.Runner), opts.Stdout),
		pdfInfo:        NewPDFInfoWizard(pdf.NewInspector(opts.Runner), opts.Stdout),
		doctor:         opts.Doctor,
		stdout:         opts.Stdout,
		accessibleUI:   opts.AccessibleUI,
	}
}

func (a *App) Run(ctx context.Context) error {
	if a.stdout != nil {
		_, _ = fmt.Fprintln(a.stdout, RenderHeader())
		_, _ = fmt.Fprintln(a.stdout)
	}

	for {
		choice, err := a.promptMainMenu()
		if err != nil {
			return cancelIfInterrupted(err)
		}

		switch choice {
		case "convert":
			subChoice, err := a.promptMenu("Convert", "Convert images, PDFs, and supported file formats.", convertMenuChoices())
			if err != nil {
				return cancelIfInterrupted(err)
			}
			a.handleInteractiveResult(a.runConvertChoice(ctx, subChoice))
		case "compress":
			subChoice, err := a.promptMenu("Compress", "Reduce image or PDF file size.", compressMenuChoices())
			if err != nil {
				return cancelIfInterrupted(err)
			}
			a.handleInteractiveResult(a.runCompressChoice(ctx, subChoice))
		case "pdf_tools":
			subChoice, err := a.promptMenu("PDF Tools", "Merge, split, inspect, and optimize PDFs.", pdfMenuChoices())
			if err != nil {
				return cancelIfInterrupted(err)
			}
			a.handleInteractiveResult(a.runPDFChoice(ctx, subChoice))
		case "doctor":
			if a.doctor != nil {
				a.handleInteractiveResult(a.doctor(ctx, a.stdout))
			}
		case "coming_soon":
			subChoice, err := a.promptMenu("Coming Soon", "Preview features planned for future releases.", comingSoonMenuChoices())
			if err != nil {
				return cancelIfInterrupted(err)
			}
			a.handleInteractiveResult(a.runComingSoonChoice(subChoice))
		case "exit":
			return nil
		}
	}
}

func (a *App) promptMainMenu() (string, error) {
	return a.promptMenu("What would you like to do?", "Choose a workflow to continue.", topLevelMenuChoices())
}

func (a *App) promptMenu(title string, help string, choices []menuChoice) (string, error) {
	if a.stdout != nil {
		_, _ = fmt.Fprintln(a.stdout, RenderSection(title))
		_, _ = fmt.Fprintln(a.stdout, RenderHelp(help))
	}

	var choice string
	err := huh.NewSelect[string]().
		Title("").
		Options(buildHuhOptions(choices)...).
		Value(&choice).
		Run()
	return choice, err
}

func buildHuhOptions(choices []menuChoice) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(choices))
	for _, choice := range choices {
		options = append(options, huh.NewOption(choice.Label, choice.Value))
	}
	return options
}

func topLevelMenuChoices() []menuChoice {
	return []menuChoice{
		{Label: "Convert", Value: "convert"},
		{Label: "Compress", Value: "compress"},
		{Label: "PDF Tools", Value: "pdf_tools"},
		{Label: "System Check", Value: "doctor"},
		{Label: "Coming Soon", Value: "coming_soon"},
		{Label: "Exit", Value: "exit"},
	}
}

func convertMenuChoices() []menuChoice {
	return []menuChoice{
		{Label: "Image to Image", Value: "image_to_image"},
		{Label: "PDF to Image", Value: "pdf_to_image"},
		{Label: "Image to PDF", Value: "image_to_pdf"},
		{Label: "Back", Value: "back"},
	}
}

func compressMenuChoices() []menuChoice {
	return []menuChoice{
		{Label: "Image Compression", Value: "image_compression"},
		{Label: "PDF Compression", Value: "pdf_compression"},
		{Label: "Back", Value: "back"},
	}
}

func pdfMenuChoices() []menuChoice {
	return []menuChoice{
		{Label: "Merge PDFs", Value: "pdf_merge"},
		{Label: "Split PDF", Value: "pdf_split"},
		{Label: "PDF Info", Value: "pdf_info"},
		{Label: "Back", Value: "back"},
	}
}

func comingSoonMenuChoices() []menuChoice {
	return []menuChoice{
		{Label: "OCR", Value: "ocr"},
		{Label: "Office Conversion", Value: "office_conversion"},
		{Label: "Metadata Cleanup", Value: "metadata_cleanup"},
		{Label: "Back", Value: "back"},
	}
}

func (a *App) runConvertChoice(ctx context.Context, choice string) error {
	switch choice {
	case "image_to_image":
		return a.convertWizard.Run(ctx)
	case "pdf_to_image":
		return a.pdfToImage.Run(ctx)
	case "image_to_pdf":
		return a.imageToPDF.Run(ctx)
	case "back":
		return nil
	default:
		return nil
	}
}

func (a *App) runCompressChoice(ctx context.Context, choice string) error {
	switch choice {
	case "image_compression":
		return a.compressWizard.Run(ctx)
	case "pdf_compression":
		return a.pdfCompress.Run(ctx)
	case "back":
		return nil
	default:
		return nil
	}
}

func (a *App) runPDFChoice(ctx context.Context, choice string) error {
	switch choice {
	case "pdf_merge":
		return a.pdfMerge.Run(ctx)
	case "pdf_split":
		return a.pdfSplit.Run(ctx)
	case "pdf_info":
		return a.pdfInfo.Run(ctx)
	case "pdf_compression":
		return a.pdfCompress.Run(ctx)
	case "back":
		return nil
	default:
		return nil
	}
}

func (a *App) runComingSoonChoice(choice string) error {
	switch choice {
	case "ocr":
		return a.printBlock(RenderComingSoon("OCR"))
	case "office_conversion":
		return a.printBlock(RenderComingSoon("Office Conversion"))
	case "metadata_cleanup":
		return a.printBlock(RenderComingSoon("Metadata Cleanup"))
	case "back":
		return nil
	default:
		return nil
	}
}

func (a *App) printBlock(value string) error {
	if a.stdout != nil {
		_, _ = fmt.Fprintln(a.stdout, value)
		_, _ = fmt.Fprintln(a.stdout)
	}
	return nil
}

func (a *App) handleInteractiveResult(err error) {
	if err == nil || IsCancelled(err) {
		return
	}
	_ = a.printBlock(RenderError(err))
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
