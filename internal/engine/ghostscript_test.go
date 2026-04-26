package engine

import (
	"context"
	"errors"
	"testing"

	"fileforge/internal/runner"
)

type ghostscriptRunner struct {
	results map[string]runner.Result
	errs    map[string]error
}

func (f ghostscriptRunner) Run(_ context.Context, name string, _ []string) (runner.Result, error) {
	if err, ok := f.errs[name]; ok {
		return runner.Result{}, err
	}
	if result, ok := f.results[name]; ok {
		return result, nil
	}
	return runner.Result{}, errors.New("missing stub")
}

func TestDetectGhostscriptBinaryPrefersFirstAvailable(t *testing.T) {
	r := ghostscriptRunner{
		results: map[string]runner.Result{
			"gswin64c": {Stdout: "10.03.1"},
			"gswin32c": {Stdout: "10.03.1"},
		},
		errs: map[string]error{
			"gs": errors.New("not found"),
		},
	}

	got, err := DetectGhostscriptBinary(context.Background(), r)
	if err != nil {
		t.Fatalf("DetectGhostscriptBinary() error = %v", err)
	}
	if got != "gswin64c" {
		t.Fatalf("DetectGhostscriptBinary() = %q, want gswin64c", got)
	}
}

func TestDetectGhostscriptBinaryFailsWhenMissing(t *testing.T) {
	r := ghostscriptRunner{
		errs: map[string]error{
			"gs":       errors.New("not found"),
			"gswin64c": errors.New("not found"),
			"gswin32c": errors.New("not found"),
		},
	}

	if _, err := DetectGhostscriptBinary(context.Background(), r); err == nil {
		t.Fatal("expected error when ghostscript is missing")
	}
}
