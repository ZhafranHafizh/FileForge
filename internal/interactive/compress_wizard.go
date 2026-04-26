package interactive

import (
	"context"
	"fmt"
	"os"

	"fileforge/internal/compress"
	"fileforge/internal/validation"

	"charm.land/huh/v2"
)

type imageCompressService interface {
	Compress(ctx context.Context, req compress.ImageCompressRequest) (compress.ImageCompressResult, error)
}

type CompressWizard struct {
	service imageCompressService
}

type CompressInput struct {
	InputPath  string
	OutputPath string
	Quality    int
	Force      bool
}

func NewCompressWizard(service imageCompressService) *CompressWizard {
	return &CompressWizard{service: service}
}

func (w *CompressWizard) Run(ctx context.Context) error {
	var state = CompressInput{Quality: 80}
	var quality string = "80"

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
				Title("Output path").
				Value(&state.OutputPath).
				Validate(func(v string) error {
					if NormalizePath(v) == "" {
						return fmt.Errorf("output path is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Quality").
				Description("Enter a value from 1 to 100.").
				Value(&quality).
				Validate(func(v string) error {
					value, err := parseQuality(v)
					if err != nil {
						return err
					}
					state.Quality = value
					return nil
				}),
		),
	)
	if err := form.Run(); err != nil {
		return cancelIfInterrupted(err)
	}

	state.InputPath = NormalizePath(state.InputPath)
	state.OutputPath = NormalizePath(state.OutputPath)

	if force, err := confirmOverwrite(state.OutputPath); err != nil {
		return err
	} else {
		state.Force = force
	}

	var confirm bool
	summary := fmt.Sprintf("Compress\nInput: %s\nOutput: %s\nQuality: %d\nOverwrite: %t", state.InputPath, state.OutputPath, state.Quality, state.Force)
	if err := huh.NewConfirm().Title(summary).Affirmative("Run").Negative("Cancel").Value(&confirm).Run(); err != nil {
		return cancelIfInterrupted(err)
	}
	if !confirm {
		return ErrCancelled
	}

	return w.Execute(ctx, state)
}

func (w *CompressWizard) Execute(ctx context.Context, in CompressInput) error {
	req := compress.ImageCompressRequest{
		InputPath:  NormalizePath(in.InputPath),
		OutputPath: NormalizePath(in.OutputPath),
		Quality:    in.Quality,
		Force:      in.Force,
	}
	_, err := w.service.Compress(ctx, req)
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

func parseQuality(raw string) (int, error) {
	var value int
	if _, err := fmt.Sscanf(raw, "%d", &value); err != nil {
		return 0, fmt.Errorf("quality must be an integer")
	}
	if value < 1 || value > 100 {
		return 0, fmt.Errorf("quality must be between 1 and 100")
	}
	return value, nil
}

func outputExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
