package pdf

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"fileforge/internal/engine"
	"fileforge/internal/validation"
)

type InfoRequest struct {
	InputPath string
}

type InfoPlan struct {
	InputPath string
}

type InfoResult struct {
	RawOutput string
	Pages     *int
}

type Inspector struct {
	runner commandRunner
}

func NewInspector(run commandRunner) *Inspector {
	return &Inspector{runner: run}
}

func (i *Inspector) Info(ctx context.Context, req InfoRequest) (InfoResult, error) {
	plan, err := BuildInfoPlan(req)
	if err != nil {
		return InfoResult{}, err
	}

	if err := engine.RequireWithRunner(ctx, i.runner, "pdfinfo", "-v"); err != nil {
		return InfoResult{}, DependencyError{Err: fmt.Errorf("pdfinfo is required: %w", err)}
	}

	result, err := i.runner.Run(ctx, "pdfinfo", plan.Args())
	if err != nil {
		return InfoResult{}, fmt.Errorf("pdf info failed: %w", err)
	}

	raw := strings.TrimSpace(result.Stdout)
	if raw == "" {
		raw = strings.TrimSpace(result.Stderr)
	}

	return InfoResult{
		RawOutput: raw,
		Pages:     extractPages(raw),
	}, nil
}

func BuildInfoPlan(req InfoRequest) (InfoPlan, error) {
	if err := validatePDFInput(req.InputPath); err != nil {
		return InfoPlan{}, err
	}

	return InfoPlan{InputPath: req.InputPath}, nil
}

func (p InfoPlan) Args() []string {
	return []string{p.InputPath}
}

func extractPages(raw string) *int {
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "Pages:") {
			continue
		}

		value := strings.TrimSpace(strings.TrimPrefix(line, "Pages:"))
		pages, err := strconv.Atoi(value)
		if err != nil {
			return nil
		}
		return &pages
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
