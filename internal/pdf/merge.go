package pdf

import (
	"context"
	"fmt"

	"fileforge/internal/engine"
	"fileforge/internal/runner"
	"fileforge/internal/validation"
)

type commandRunner interface {
	Run(ctx context.Context, name string, args []string) (runner.Result, error)
}

type MergeRequest struct {
	InputPaths []string
	OutputPath string
	Force      bool
}

type MergePlan struct {
	InputPaths []string
	OutputPath string
}

type Merger struct {
	runner commandRunner
}

func NewMerger(run commandRunner) *Merger {
	return &Merger{runner: run}
}

func (m *Merger) Merge(ctx context.Context, req MergeRequest) error {
	plan, err := BuildMergePlan(req)
	if err != nil {
		return err
	}

	if err := engine.RequireWithRunner(ctx, m.runner, "qpdf", "--version"); err != nil {
		return DependencyError{Err: fmt.Errorf("qpdf is required: %w", err)}
	}

	if _, err := m.runner.Run(ctx, "qpdf", plan.Args()); err != nil {
		return fmt.Errorf("pdf merge failed: %w", err)
	}

	return nil
}

func BuildMergePlan(req MergeRequest) (MergePlan, error) {
	if len(req.InputPaths) < 2 {
		return MergePlan{}, ValidationError{Err: fmt.Errorf("pdf merge requires at least two input files")}
	}

	for _, path := range req.InputPaths {
		if err := validatePDFInput(path); err != nil {
			return MergePlan{}, err
		}
	}

	if err := validation.EnsureOutputPath(req.OutputPath, req.Force); err != nil {
		return MergePlan{}, ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(req.OutputPath, []string{"pdf"}); err != nil {
		return MergePlan{}, ValidationError{Err: err}
	}

	return MergePlan{
		InputPaths: req.InputPaths,
		OutputPath: req.OutputPath,
	}, nil
}

func (p MergePlan) Args() []string {
	args := []string{"--empty", "--pages"}
	args = append(args, p.InputPaths...)
	args = append(args, "--", p.OutputPath)
	return args
}
