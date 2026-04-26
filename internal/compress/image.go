package compress

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fileforge/internal/engine"
	"fileforge/internal/runner"
	"fileforge/internal/validation"
)

var supportedImageFormats = []string{"jpg", "jpeg", "png", "webp"}

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

type ImageCompressRequest struct {
	InputPath  string
	OutputPath string
	Quality    int
	Force      bool
}

type ImageCompressResult struct {
	InputPath     string
	OutputPath    string
	BeforeBytes   int64
	AfterBytes    int64
	PercentChange *float64
}

type ImageCompressor struct {
	runner commandRunner
}

type commandRunner interface {
	Run(ctx context.Context, name string, args []string) (runner.Result, error)
}

type CompressPlan struct {
	InputPath  string
	OutputPath string
	Quality    int
}

func NewImageCompressor(run commandRunner) *ImageCompressor {
	return &ImageCompressor{runner: run}
}

func IsValidationError(err error) bool {
	var target ValidationError
	return errors.As(err, &target)
}

func IsDependencyError(err error) bool {
	var target DependencyError
	return errors.As(err, &target)
}

func (c *ImageCompressor) Compress(ctx context.Context, req ImageCompressRequest) (ImageCompressResult, error) {
	plan, err := BuildCompressPlan(req)
	if err != nil {
		return ImageCompressResult{}, err
	}

	beforeBytes, err := validation.FileSize(plan.InputPath)
	if err != nil {
		return ImageCompressResult{}, ValidationError{Err: err}
	}

	if err := engine.RequireWithRunner(ctx, c.runner, "magick", "-version"); err != nil {
		return ImageCompressResult{}, DependencyError{Err: fmt.Errorf("imagemagick is required: %w", err)}
	}

	if _, err := c.runner.Run(ctx, "magick", plan.Args()); err != nil {
		return ImageCompressResult{}, fmt.Errorf("image compression failed: %w", err)
	}

	afterBytes, err := validation.FileSize(plan.OutputPath)
	if err != nil {
		return ImageCompressResult{}, fmt.Errorf("read compressed output size: %w", err)
	}

	return ImageCompressResult{
		InputPath:     plan.InputPath,
		OutputPath:    plan.OutputPath,
		BeforeBytes:   beforeBytes,
		AfterBytes:    afterBytes,
		PercentChange: percentChange(beforeBytes, afterBytes),
	}, nil
}

func BuildCompressPlan(req ImageCompressRequest) (CompressPlan, error) {
	if err := validateInput(req.InputPath); err != nil {
		return CompressPlan{}, err
	}
	if err := validation.EnsureOutputPath(req.OutputPath, req.Force); err != nil {
		return CompressPlan{}, ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(req.OutputPath, supportedImageFormats); err != nil {
		return CompressPlan{}, ValidationError{Err: err}
	}
	if err := validateQuality(req.Quality); err != nil {
		return CompressPlan{}, err
	}

	return CompressPlan{
		InputPath:  req.InputPath,
		OutputPath: req.OutputPath,
		Quality:    req.Quality,
	}, nil
}

func (p CompressPlan) Args() []string {
	return []string{p.InputPath, "-quality", strconv.Itoa(p.Quality), p.OutputPath}
}

func validateInput(path string) error {
	if err := validation.EnsureInputFile(path); err != nil {
		return ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(path, supportedImageFormats); err != nil {
		return ValidationError{Err: err}
	}
	return nil
}

func validateQuality(quality int) error {
	if quality < 1 || quality > 100 {
		return ValidationError{Err: fmt.Errorf("quality must be between 1 and 100")}
	}
	return nil
}

func percentChange(before int64, after int64) *float64 {
	if before <= 0 {
		return nil
	}

	change := (float64(before-after) / float64(before)) * 100
	return &change
}

func FormatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	value := float64(size)
	suffixes := []string{"KB", "MB", "GB", "TB"}
	index := -1
	for value >= unit && index < len(suffixes)-1 {
		value /= unit
		index++
	}

	if index < 0 {
		return fmt.Sprintf("%d B", size)
	}

	return fmt.Sprintf("%.2f %s", value, suffixes[index])
}

func SupportedFormats() string {
	return strings.Join(supportedImageFormats, ", ")
}
