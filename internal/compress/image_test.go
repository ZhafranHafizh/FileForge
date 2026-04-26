package compress

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildCompressPlanValid(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.jpg")
	output := filepath.Join(dir, "compressed.jpg")
	if err := os.WriteFile(input, []byte("image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildCompressPlan(ImageCompressRequest{
		InputPath:  input,
		OutputPath: output,
		Quality:    80,
	})
	if err != nil {
		t.Fatalf("BuildCompressPlan() error = %v", err)
	}

	wantArgs := []string{input, "-quality", "80", output}
	if !reflect.DeepEqual(plan.Args(), wantArgs) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), wantArgs)
	}
}

func TestBuildCompressPlanRejectsInvalidQuality(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.jpg")
	output := filepath.Join(dir, "compressed.jpg")
	if err := os.WriteFile(input, []byte("image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildCompressPlan(ImageCompressRequest{
		InputPath:  input,
		OutputPath: output,
		Quality:    101,
	})
	if err == nil {
		t.Fatal("expected invalid quality error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestBuildCompressPlanRejectsUnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.gif")
	output := filepath.Join(dir, "compressed.gif")
	if err := os.WriteFile(input, []byte("image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildCompressPlan(ImageCompressRequest{
		InputPath:  input,
		OutputPath: output,
		Quality:    80,
	})
	if err == nil {
		t.Fatal("expected unsupported extension error")
	}
}

func TestPercentChange(t *testing.T) {
	got := percentChange(200, 150)
	if got == nil {
		t.Fatal("expected percent change")
	}
	if *got != 25 {
		t.Fatalf("percentChange() = %v, want 25", *got)
	}
}
