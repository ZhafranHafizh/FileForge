package convert

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"fileforge/internal/engine"
	"fileforge/internal/validation"
)

type PDFToImageRequest struct {
	InputPath string
	OutputDir string
	ToFormat  string
	DPI       int
	FirstPage int
	LastPage  int
	Force     bool
}

type PDFToImagePlan struct {
	InputPath  string
	OutputDir  string
	OutputBase string
	ToFormat   string
	DPI        int
	FirstPage  int
	LastPage   int
}

type PDFToImageConverter struct {
	runner commandRunner
}

func NewPDFToImageConverter(run commandRunner) *PDFToImageConverter {
	return &PDFToImageConverter{runner: run}
}

func (c *PDFToImageConverter) Convert(ctx context.Context, req PDFToImageRequest) error {
	plan, err := BuildPDFToImagePlan(req)
	if err != nil {
		return err
	}

	if err := engine.RequireWithRunner(ctx, c.runner, "pdftoppm", "-v"); err != nil {
		return DependencyError{Err: fmt.Errorf("pdftoppm is required: %w", err)}
	}

	if _, err := c.runner.Run(ctx, "pdftoppm", plan.Args()); err != nil {
		return fmt.Errorf("pdf to image conversion failed: %w", err)
	}

	return nil
}

func validatePDFInput(path string) error {
	if err := validation.EnsureInputFile(path); err != nil {
		return ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(path, []string{"pdf"}); err != nil {
		return ValidationError{Err: err}
	}
	return nil
}

func BuildPDFToImagePlan(req PDFToImageRequest) (PDFToImagePlan, error) {
	if err := validatePDFInput(req.InputPath); err != nil {
		return PDFToImagePlan{}, err
	}

	toFormat := validation.NormalizeExtension(req.ToFormat)
	if !validation.IsSupportedExtension(toFormat, supportedPDFImageFormats) {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("unsupported PDF to image target format %q; supported formats: jpg, jpeg, png", req.ToFormat)}
	}

	if req.DPI <= 0 {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("dpi must be greater than 0")}
	}
	if req.FirstPage < 0 || req.LastPage < 0 {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("first-page and last-page must be greater than 0 when provided")}
	}
	if req.FirstPage == 0 && req.LastPage > 0 {
		// allowed
	}
	if req.FirstPage > 0 && req.LastPage > 0 && req.LastPage < req.FirstPage {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("last-page must be greater than or equal to first-page")}
	}
	if req.FirstPage < 0 || req.LastPage < 0 {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("page values must be greater than 0")}
	}
	if req.FirstPage == 0 && req.LastPage == 0 {
		// okay
	}
	if req.FirstPage != 0 && req.FirstPage < 1 {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("first-page must be greater than 0")}
	}
	if req.LastPage != 0 && req.LastPage < 1 {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("last-page must be greater than 0")}
	}

	outputDir := filepath.Clean(req.OutputDir)
	if outputDir == "" || outputDir == "." {
		return PDFToImagePlan{}, ValidationError{Err: fmt.Errorf("output directory is required")}
	}
	if err := ensureOutputDirectory(outputDir, toFormat, req.Force); err != nil {
		return PDFToImagePlan{}, ValidationError{Err: err}
	}

	return PDFToImagePlan{
		InputPath:  req.InputPath,
		OutputDir:  outputDir,
		OutputBase: filepath.Join(outputDir, "page"),
		ToFormat:   toFormat,
		DPI:        req.DPI,
		FirstPage:  req.FirstPage,
		LastPage:   req.LastPage,
	}, nil
}

func (p PDFToImagePlan) Args() []string {
	args := []string{}
	if p.ToFormat == "png" {
		args = append(args, "-png")
	} else {
		args = append(args, "-jpeg")
	}
	args = append(args, "-r", strconv.Itoa(p.DPI))
	if p.FirstPage > 0 {
		args = append(args, "-f", strconv.Itoa(p.FirstPage))
	}
	if p.LastPage > 0 {
		args = append(args, "-l", strconv.Itoa(p.LastPage))
	}
	args = append(args, p.InputPath, p.OutputBase)
	return args
}

func ensureOutputDirectory(path string, format string, force bool) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("output path must be a directory: %s", path)
		}
		if !force {
			pattern := filepath.Join(path, "page-*."+format)
			matches, globErr := filepath.Glob(pattern)
			if globErr != nil {
				return fmt.Errorf("scan output directory %q: %w", path, globErr)
			}
			if len(matches) > 0 {
				return fmt.Errorf("output directory already contains generated %s files (use --force to overwrite)", format)
			}
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("stat output directory %q: %w", path, err)
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", path, err)
	}
	return nil
}
