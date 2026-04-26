package interactive

import (
	"context"
	"fmt"
	"io"

	"fileforge/internal/compress"
	outputpkg "fileforge/internal/output"
	"fileforge/internal/validation"

	"charm.land/huh/v2"
)

type imageCompressService interface {
	Compress(ctx context.Context, req compress.ImageCompressRequest) (compress.ImageCompressResult, error)
}

type CompressWizard struct {
	service imageCompressService
	stdout  io.Writer
}

type CompressInput struct {
	InputPath  string
	OutputDir  string
	OutputPath string
	Quality    int
	Force      bool
}

func NewCompressWizard(service imageCompressService, stdout io.Writer) *CompressWizard {
	return &CompressWizard{service: service, stdout: stdout}
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
				Title("Output folder (optional, default: ./FileForge-Output)").
				Value(&state.OutputDir).
				Validate(func(v string) error { return nil }),
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
	state.OutputDir = ResolveInteractiveOutputDir(state.OutputDir)

	outputPath, err := outputpkg.ResolveOutputPath(state.InputPath, "", state.OutputDir, "-compressed", validation.Extension(state.InputPath), true)
	if err != nil {
		return ValidationError{Err: err}
	}
	state.OutputPath = outputPath

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

	if err := w.Execute(ctx, state); err != nil {
		return err
	}
	return printSuccess(w.stdout, state.OutputPath)
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
