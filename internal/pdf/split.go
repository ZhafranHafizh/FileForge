package pdf

import (
	"context"
	"fmt"
	"strings"

	"fileforge/internal/engine"
	"fileforge/internal/validation"
)

type SplitRequest struct {
	InputPath  string
	PageRange  string
	OutputPath string
	Force      bool
}

type SplitPlan struct {
	InputPath  string
	PageRange  string
	OutputPath string
}

type Splitter struct {
	runner commandRunner
}

func NewSplitter(run commandRunner) *Splitter {
	return &Splitter{runner: run}
}

func (s *Splitter) Split(ctx context.Context, req SplitRequest) error {
	plan, err := BuildSplitPlan(req)
	if err != nil {
		return err
	}

	if err := engine.RequireWithRunner(ctx, s.runner, "qpdf", "--version"); err != nil {
		return DependencyError{Err: fmt.Errorf("qpdf is required: %w", err)}
	}

	if _, err := s.runner.Run(ctx, "qpdf", plan.Args()); err != nil {
		return fmt.Errorf("pdf split failed: %w", err)
	}

	return nil
}

func BuildSplitPlan(req SplitRequest) (SplitPlan, error) {
	if err := validatePDFInput(req.InputPath); err != nil {
		return SplitPlan{}, err
	}
	if _, err := validation.ParsePageRange(req.PageRange); err != nil {
		return SplitPlan{}, ValidationError{Err: err}
	}
	if err := validation.EnsureOutputPath(req.OutputPath, req.Force); err != nil {
		return SplitPlan{}, ValidationError{Err: err}
	}
	if err := validation.EnsureSupportedExtension(req.OutputPath, []string{"pdf"}); err != nil {
		return SplitPlan{}, ValidationError{Err: err}
	}

	return SplitPlan{
		InputPath:  req.InputPath,
		PageRange:  strings.TrimSpace(req.PageRange),
		OutputPath: req.OutputPath,
	}, nil
}

func (p SplitPlan) Args() []string {
	return []string{p.InputPath, "--pages", p.InputPath, p.PageRange, "--", p.OutputPath}
}
