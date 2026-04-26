package validation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureInputFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := EnsureInputFile(path); err != nil {
		t.Fatalf("EnsureInputFile() error = %v", err)
	}
}

func TestEnsureInputFileMissing(t *testing.T) {
	err := EnsureInputFile(filepath.Join(t.TempDir(), "missing.txt"))
	if err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestEnsureRegularFileRejectsDirectory(t *testing.T) {
	err := EnsureRegularFile(t.TempDir())
	if err == nil {
		t.Fatal("expected error for directory input")
	}
}

func TestEnsureOutputPathProtectsOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := EnsureOutputPath(path, false)
	if err == nil {
		t.Fatal("expected overwrite protection error")
	}
}

func TestEnsureOutputPathAllowsOverwriteWithForce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := EnsureOutputPath(path, true); err != nil {
		t.Fatalf("EnsureOutputPath() error = %v", err)
	}
}

func TestEnsureOutputDirCreatesParent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "output.txt")

	if err := EnsureOutputDir(path); err != nil {
		t.Fatalf("EnsureOutputDir() error = %v", err)
	}

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("expected output directory to exist: %v", err)
	}
}

func TestEnsureOutputPathCreatesParent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "output.txt")

	if err := EnsureOutputPath(path, false); err != nil {
		t.Fatalf("EnsureOutputPath() error = %v", err)
	}

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("expected parent directory to exist: %v", err)
	}
}

func TestFileSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("hello")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	size, err := FileSize(path)
	if err != nil {
		t.Fatalf("FileSize() error = %v", err)
	}
	if size != int64(len(content)) {
		t.Fatalf("FileSize() = %d, want %d", size, len(content))
	}
}
