package compress

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildPDFCompressPlanArgs(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	output := filepath.Join(dir, "compressed.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildPDFCompressPlan(PDFCompressRequest{
		InputPath:  input,
		OutputPath: output,
		Preset:     "ebook",
	}, "gs")
	if err != nil {
		t.Fatalf("BuildPDFCompressPlan() error = %v", err)
	}

	want := []string{
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=/ebook",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile=" + output,
		input,
	}
	if !reflect.DeepEqual(plan.Args(), want) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), want)
	}
}

func TestBuildPDFCompressPlanRejectsInvalidPreset(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	output := filepath.Join(dir, "compressed.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildPDFCompressPlan(PDFCompressRequest{
		InputPath:  input,
		OutputPath: output,
		Preset:     "invalid",
	}, "gs")
	if err == nil {
		t.Fatal("expected invalid preset error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestBuildPDFCompressPlanRejectsNonPDFOutput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	output := filepath.Join(dir, "compressed.png")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildPDFCompressPlan(PDFCompressRequest{
		InputPath:  input,
		OutputPath: output,
		Preset:     "ebook",
	}, "gs")
	if err == nil {
		t.Fatal("expected non-pdf output error")
	}
}

func TestBuildPDFCompressPlanDefaultsPreset(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	output := filepath.Join(dir, "compressed.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildPDFCompressPlan(PDFCompressRequest{
		InputPath:  input,
		OutputPath: output,
	}, "gs")
	if err != nil {
		t.Fatalf("BuildPDFCompressPlan() error = %v", err)
	}
	if plan.Preset != "ebook" {
		t.Fatalf("Preset = %q, want ebook", plan.Preset)
	}
}
