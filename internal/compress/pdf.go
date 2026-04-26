package compress

import (
	"context"
	"fmt"
	"strings"

	"fileforge/internal/engine"
	"fileforge/internal/validation"
)

var supportedPDFPresets = map[string]string{
	"screen":   "/screen",
	"ebook":    "/ebook",
	"printer":  "/printer",
	"prepress": "/prepress",
	"default":  "/default",
}

type PDFCompressRequest struct {
	InputPath  string
	OutputPath string
	Preset     string
	Force      bool
}

type PDFCompressResult struct {
	InputPath     string
	OutputPath    string
	BeforeBytes   int64
	AfterBytes    int64
	PercentChange *float64
}

type PDFCompressPlan struct {
	Binary     string
	InputPath  string
	OutputPath string
	Preset     string
}

type PDFCompressor struct {
	runner commandRunner
}

func NewPDFCompressor(run commandRunner) *PDFCompressor {
	return &PDFCompressor{runner: run}
}

func (c *PDFCompressor) Compress(ctx context.Context, req PDFCompressRequest) (PDFCompressResult, error) {
	binary, err := engine.DetectGhostscriptBinary(ctx, c.runner)
	if err != nil {
		return PDFCompressResult{}, DependencyError{Err: fmt.Errorf("ghostscript is required: %w", err)}
	}

	plan, err := BuildPDFCompressPlan(req, binary)
	if err != nil {
		return PDFCompressResult{}, err
	}

	beforeBytes, err := validation.FileSize(plan.InputPath)
	if err != nil {
		return PDFCompressResult{}, ValidationError{Err: err}
	}

	if _, err := c.runner.Run(ctx, plan.Binary, plan.Args()); err != nil {
		return PDFCompressResult{}, fmt.Errorf("pdf compression failed: %w", err)
	}

	afterBytes, err := validation.FileSize(plan.OutputPath)
	if err != nil {
		return PDFCompressResult{}, fmt.Errorf("read compressed output size: %w", err)
	}

	return PDFCompressResult{
		InputPath:     plan.InputPath,
		OutputPath:    plan.OutputPath,
		BeforeBytes:   beforeBytes,
		AfterBytes:    afterBytes,
		PercentChange: percentChange(beforeBytes, afterBytes),
	}, nil
}

func BuildPDFCompressPlan(req PDFCompressRequest, binary string) (PDFCompressPlan, error) {
	if err := validatePDFInput(req.InputPath); err != nil {
		return PDFCompressPlan{}, err
	}
	if err := validation.EnsureOutputPath(req.OutputPath, req.Force); err != nil {
		return PDFCompressPlan{}, ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(req.OutputPath, []string{"pdf"}); err != nil {
		return PDFCompressPlan{}, ValidationError{Err: err}
	}

	preset, err := normalizePDFPreset(req.Preset)
	if err != nil {
		return PDFCompressPlan{}, err
	}

	return PDFCompressPlan{
		Binary:     binary,
		InputPath:  req.InputPath,
		OutputPath: req.OutputPath,
		Preset:     preset,
	}, nil
}

func (p PDFCompressPlan) Args() []string {
	return []string{
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=" + supportedPDFPresets[p.Preset],
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile=" + p.OutputPath,
		p.InputPath,
	}
}

func normalizePDFPreset(preset string) (string, error) {
	normalized := validation.NormalizeExtension(strings.TrimSpace(preset))
	if normalized == "" {
		normalized = "ebook"
	}

	if _, ok := supportedPDFPresets[normalized]; !ok {
		return "", ValidationError{Err: fmt.Errorf("unsupported PDF preset %q; supported presets: screen, ebook, printer, prepress, default", preset)}
	}

	return normalized, nil
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
