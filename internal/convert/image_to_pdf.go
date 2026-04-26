package convert

import (
	"context"
	"fmt"

	"fileforge/internal/engine"
	"fileforge/internal/validation"
)

type ImageToPDFRequest struct {
	InputPath  string
	OutputPath string
	Force      bool
}

type ImageToPDFPlan struct {
	InputPath  string
	OutputPath string
}

type ImageToPDFConverter struct {
	runner commandRunner
}

func NewImageToPDFConverter(run commandRunner) *ImageToPDFConverter {
	return &ImageToPDFConverter{runner: run}
}

func (c *ImageToPDFConverter) Convert(ctx context.Context, req ImageToPDFRequest) error {
	plan, err := BuildImageToPDFPlan(req)
	if err != nil {
		return err
	}

	if err := engine.RequireWithRunner(ctx, c.runner, "magick", "-version"); err != nil {
		return DependencyError{Err: fmt.Errorf("imagemagick is required: %w", err)}
	}

	if _, err := c.runner.Run(ctx, "magick", plan.Args()); err != nil {
		return fmt.Errorf("image to pdf conversion failed: %w", err)
	}
	return nil
}

func BuildImageToPDFPlan(req ImageToPDFRequest) (ImageToPDFPlan, error) {
	if err := validateImageInput(req.InputPath); err != nil {
		return ImageToPDFPlan{}, err
	}
	if err := validation.EnsureOutputPath(req.OutputPath, req.Force); err != nil {
		return ImageToPDFPlan{}, ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(req.OutputPath, []string{"pdf"}); err != nil {
		return ImageToPDFPlan{}, ValidationError{Err: err}
	}
	return ImageToPDFPlan{InputPath: req.InputPath, OutputPath: req.OutputPath}, nil
}

func (p ImageToPDFPlan) Args() []string {
	return []string{p.InputPath, p.OutputPath}
}
