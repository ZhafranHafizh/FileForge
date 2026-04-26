package convert

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildImageToPDFPlanArgs(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.jpg")
	output := filepath.Join(dir, "output.pdf")
	if err := os.WriteFile(input, []byte("img"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildImageToPDFPlan(ImageToPDFRequest{InputPath: input, OutputPath: output})
	if err != nil {
		t.Fatalf("BuildImageToPDFPlan() error = %v", err)
	}
	want := []string{input, output}
	if !reflect.DeepEqual(plan.Args(), want) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), want)
	}
}

func TestBuildImageToPDFPlanRejectsOutputExtension(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.jpg")
	output := filepath.Join(dir, "output.png")
	if err := os.WriteFile(input, []byte("img"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildImageToPDFPlan(ImageToPDFRequest{InputPath: input, OutputPath: output})
	if err == nil {
		t.Fatal("expected output extension validation error")
	}
}
