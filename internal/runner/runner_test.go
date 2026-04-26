package runner

import (
	"context"
	"testing"
)

func TestRunCapturesStdout(t *testing.T) {
	r := New(Options{})

	result, err := r.Run(context.Background(), "go", []string{"env", "GOOS"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Stdout == "" {
		t.Fatal("expected stdout to be captured")
	}
	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", result.ExitCode)
	}
}

func TestRunMissingBinary(t *testing.T) {
	r := New(Options{})

	result, err := r.Run(context.Background(), "fileforge-does-not-exist", nil)
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
	if result.ExitCode != -1 {
		t.Fatalf("expected exit code -1 for missing binary, got %d", result.ExitCode)
	}
}
