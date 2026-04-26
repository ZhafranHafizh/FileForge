package convert

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildConvertPlanValid(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.png")
	output := filepath.Join(dir, "output.jpg")
	if err := os.WriteFile(input, []byte("image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildConvertPlan(ImageConvertRequest{
		InputPath:  input,
		OutputPath: output,
		ToFormat:   "jpg",
	})
	if err != nil {
		t.Fatalf("BuildConvertPlan() error = %v", err)
	}

	wantArgs := []string{input, output}
	if !reflect.DeepEqual(plan.Args(), wantArgs) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), wantArgs)
	}
}

func TestBuildConvertPlanRejectsUnsupportedTargetFormat(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.png")
	output := filepath.Join(dir, "output.gif")
	if err := os.WriteFile(input, []byte("image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildConvertPlan(ImageConvertRequest{
		InputPath:  input,
		OutputPath: output,
		ToFormat:   "gif",
	})
	if err == nil {
		t.Fatal("expected unsupported format error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestBuildConvertPlanRejectsMismatchedOutputExtension(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.png")
	output := filepath.Join(dir, "output.webp")
	if err := os.WriteFile(input, []byte("image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildConvertPlan(ImageConvertRequest{
		InputPath:  input,
		OutputPath: output,
		ToFormat:   "jpg",
	})
	if err == nil {
		t.Fatal("expected mismatched output extension error")
	}
}

func TestBuildConvertPlanAcceptsJPEGAliases(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.jpeg")
	output := filepath.Join(dir, "output.jpg")
	if err := os.WriteFile(input, []byte("image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := BuildConvertPlan(ImageConvertRequest{
		InputPath:  input,
		OutputPath: output,
		ToFormat:   "jpeg",
	}); err != nil {
		t.Fatalf("BuildConvertPlan() error = %v", err)
	}
}
