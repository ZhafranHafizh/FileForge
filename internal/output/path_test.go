package output

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureOnlyOneOutputMode(t *testing.T) {
	if err := EnsureOnlyOneOutputMode("a", "b"); err == nil {
		t.Fatal("expected mutual exclusion error")
	}
}

func TestResolveOutputPathPreservesOut(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputPath("image.jpg", filepath.Join(dir, "custom.jpg"), "", "-compressed", "jpg", false)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	if got != filepath.Join(dir, "custom.jpg") {
		t.Fatalf("ResolveOutputPath() = %q", got)
	}
}

func TestResolveOutputPathFromOutputDir(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputPath("image.jpg", "", dir, "-compressed", "jpg", false)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	if got != filepath.Join(dir, "image-compressed.jpg") {
		t.Fatalf("ResolveOutputPath() = %q", got)
	}
}

func TestResolvePDFCompressionOutputPath(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputPath("report.pdf", "", dir, "-compressed", "pdf", false)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	if got != filepath.Join(dir, "report-compressed.pdf") {
		t.Fatalf("ResolveOutputPath() = %q", got)
	}
}

func TestResolveImageToImageOutputPath(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputPath("photo.png", "", dir, "", "webp", false)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	if got != filepath.Join(dir, "photo.webp") {
		t.Fatalf("ResolveOutputPath() = %q", got)
	}
}

func TestResolveImageToPDFOutputPath(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputPath("photo.jpg", "", dir, "", "pdf", false)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	if got != filepath.Join(dir, "photo.pdf") {
		t.Fatalf("ResolveOutputPath() = %q", got)
	}
}

func TestResolvePDFToImageOutputDir(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputDir("", dir, GeneratedPDFToImageDirName("document.pdf"))
	if err != nil {
		t.Fatalf("ResolveOutputDir() error = %v", err)
	}
	if got != filepath.Join(dir, "document-pages") {
		t.Fatalf("ResolveOutputDir() = %q", got)
	}
}

func TestResolvePDFMergeOutputPath(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputPath("merged.pdf", "", dir, "", "pdf", false)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	if got != filepath.Join(dir, "merged.pdf") {
		t.Fatalf("ResolveOutputPath() = %q", got)
	}
}

func TestResolvePDFSplitOutputPath(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveOutputPath("book.pdf", "", dir, "-pages-"+SanitizeFilenamePart("1-3,7,10-12"), "pdf", false)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	if got != filepath.Join(dir, "book-pages-1-3_7_10-12.pdf") {
		t.Fatalf("ResolveOutputPath() = %q", got)
	}
}

func TestGeneratedPDFSplitName(t *testing.T) {
	got := GeneratedPDFSplitName("book.pdf", "1-3,7,10-12")
	if got != "book-pages-1-3_7_10-12" {
		t.Fatalf("GeneratedPDFSplitName() = %q", got)
	}
}

func TestResolveOutputPathOverwriteProtection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.pdf")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if _, err := ResolveOutputPath("a.pdf", "", dir, "", "pdf", false); err == nil {
		t.Fatal("expected overwrite error")
	}
}
