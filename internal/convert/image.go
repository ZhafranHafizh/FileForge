package convert

import (
	"context"
	"errors"
	"fmt"

	"fileforge/internal/engine"
	"fileforge/internal/runner"
	"fileforge/internal/validation"
)

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

type ImageConvertRequest struct {
	InputPath  string
	OutputPath string
	ToFormat   string
	Force      bool
}

type ImageConverter struct {
	runner commandRunner
}

type commandRunner interface {
	Run(ctx context.Context, name string, args []string) (runner.Result, error)
}

func NewImageConverter(run commandRunner) *ImageConverter {
	return &ImageConverter{runner: run}
}

func IsValidationError(err error) bool {
	var target ValidationError
	return errors.As(err, &target)
}

func IsDependencyError(err error) bool {
	var target DependencyError
	return errors.As(err, &target)
}

func (c *ImageConverter) Convert(ctx context.Context, req ImageConvertRequest) error {
	plan, err := BuildConvertPlan(req)
	if err != nil {
		return err
	}

	if err := engine.RequireWithRunner(ctx, c.runner, "magick", "-version"); err != nil {
		return DependencyError{Err: fmt.Errorf("imagemagick is required: %w", err)}
	}

	if _, err := c.runner.Run(ctx, "magick", plan.Args()); err != nil {
		return fmt.Errorf("image conversion failed: %w", err)
	}

	return nil
}

type ConvertPlan struct {
	InputPath  string
	OutputPath string
	ToFormat   string
}

func BuildConvertPlan(req ImageConvertRequest) (ConvertPlan, error) {
	if err := validateImageInput(req.InputPath); err != nil {
		return ConvertPlan{}, err
	}
	if err := validation.EnsureOutputPath(req.OutputPath, req.Force); err != nil {
		return ConvertPlan{}, ValidationError{Err: err}
	}

	toFormat := validation.NormalizeExtension(req.ToFormat)
	if !validation.IsSupportedExtension(toFormat, supportedImageFormats) {
		return ConvertPlan{}, ValidationError{Err: fmt.Errorf("unsupported target format %q; supported formats: jpg, jpeg, png, webp", req.ToFormat)}
	}

	if err := validateOutputFormat(req.OutputPath, toFormat); err != nil {
		return ConvertPlan{}, err
	}

	return ConvertPlan{
		InputPath:  req.InputPath,
		OutputPath: req.OutputPath,
		ToFormat:   toFormat,
	}, nil
}

func (p ConvertPlan) Args() []string {
	return []string{p.InputPath, p.OutputPath}
}

func validateImageInput(path string) error {
	if err := validation.EnsureInputFile(path); err != nil {
		return ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(path, supportedImageFormats); err != nil {
		return ValidationError{Err: err}
	}
	return nil
}

func validateOutputFormat(outputPath string, toFormat string) error {
	outputExt := validation.Extension(outputPath)
	if outputExt == "" {
		return ValidationError{Err: fmt.Errorf("output file must include an extension")}
	}
	if !validation.IsSupportedExtension(outputExt, supportedImageFormats) {
		return ValidationError{Err: fmt.Errorf("unsupported output extension %q; supported formats: jpg, jpeg, png, webp", outputExt)}
	}
	if !extensionsMatch(toFormat, outputExt) {
		return ValidationError{Err: fmt.Errorf("output extension %q does not match target format %q", outputExt, toFormat)}
	}
	return nil
}

func extensionsMatch(target string, actual string) bool {
	target = validation.NormalizeExtension(target)
	actual = validation.NormalizeExtension(actual)
	if target == actual {
		return true
	}
	return isJPEGAlias(target) && isJPEGAlias(actual)
}

func isJPEGAlias(ext string) bool {
	normalized := validation.NormalizeExtension(ext)
	return normalized == "jpg" || normalized == "jpeg"
}
