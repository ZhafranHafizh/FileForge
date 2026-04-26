package convert

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildPDFToImagePlanArgs(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	outDir := filepath.Join(dir, "pages")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildPDFToImagePlan(PDFToImageRequest{
		InputPath: input,
		OutputDir: outDir,
		ToFormat:  "jpg",
		DPI:       150,
		FirstPage: 1,
		LastPage:  3,
	})
	if err != nil {
		t.Fatalf("BuildPDFToImagePlan() error = %v", err)
	}

	want := []string{"-jpeg", "-r", "150", "-f", "1", "-l", "3", input, filepath.Join(outDir, "page")}
	if !reflect.DeepEqual(plan.Args(), want) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), want)
	}
}

func TestBuildPDFToImagePlanRejectsFormat(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildPDFToImagePlan(PDFToImageRequest{InputPath: input, OutputDir: filepath.Join(dir, "pages"), ToFormat: "webp", DPI: 150})
	if err == nil {
		t.Fatal("expected format validation error")
	}
}

func TestBuildPDFToImagePlanRejectsDPI(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildPDFToImagePlan(PDFToImageRequest{InputPath: input, OutputDir: filepath.Join(dir, "pages"), ToFormat: "png", DPI: 0})
	if err == nil {
		t.Fatal("expected dpi validation error")
	}
}

func TestBuildPDFToImagePlanRejectsPageRange(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildPDFToImagePlan(PDFToImageRequest{InputPath: input, OutputDir: filepath.Join(dir, "pages"), ToFormat: "png", DPI: 150, FirstPage: 3, LastPage: 1})
	if err == nil {
		t.Fatal("expected page range validation error")
	}
}
