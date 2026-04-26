package engine

import (
	"context"
	"errors"
	"testing"

	"fileforge/internal/runner"
)

type fakeRunner struct {
	results map[string]runner.Result
	errs    map[string]error
}

func (f fakeRunner) Run(_ context.Context, name string, _ []string) (runner.Result, error) {
	if err, ok := f.errs[name]; ok {
		return runner.Result{}, err
	}
	if result, ok := f.results[name]; ok {
		return result, nil
	}
	return runner.Result{}, errors.New("missing stub")
}

func TestVersionParsesOutput(t *testing.T) {
	r := fakeRunner{
		results: map[string]runner.Result{
			"magick": {Stdout: "Version: ImageMagick 7.1.1-33 Q16-HDRI x64"},
		},
	}

	got, err := VersionWithRunner(context.Background(), r, "magick", "-version")
	if err != nil {
		t.Fatalf("Version() error = %v", err)
	}
	if got != "7.1.1-33" {
		t.Fatalf("Version() = %q, want %q", got, "7.1.1-33")
	}
}

func TestDiagnoseTracksMissingTools(t *testing.T) {
	r := fakeRunner{
		results: map[string]runner.Result{
			"qpdf": {Stdout: "qpdf version 11.9.1"},
		},
		errs: map[string]error{
			"gs": errors.New("not found"),
		},
	}

	report := Diagnose(context.Background(), r, []Tool{
		{Name: "qpdf", VersionArgs: []string{"--version"}},
		{Name: "gs", VersionArgs: []string{"--version"}},
	})

	if len(report.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(report.Results))
	}
	if len(report.Missing) != 1 {
		t.Fatalf("expected 1 missing tool, got %d", len(report.Missing))
	}
	if report.Missing[0].Name != "gs" {
		t.Fatalf("missing tool = %q, want gs", report.Missing[0].Name)
	}
}
