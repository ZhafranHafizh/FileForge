package interactive

import (
	"context"
	"path/filepath"
	"testing"

	"fileforge/internal/compress"
	outputpkg "fileforge/internal/output"
)

type fakeCompressor struct {
	lastReq compress.ImageCompressRequest
	called  bool
	err     error
}

func (f *fakeCompressor) Compress(_ context.Context, req compress.ImageCompressRequest) (compress.ImageCompressResult, error) {
	f.called = true
	f.lastReq = req
	return compress.ImageCompressResult{}, f.err
}

func TestCompressWizardExecuteUsesService(t *testing.T) {
	svc := &fakeCompressor{}
	wizard := NewCompressWizard(svc, nil)

	err := wizard.Execute(context.Background(), CompressInput{
		InputPath:  `"./input.jpg"`,
		OutputPath: `'./output.jpg'`,
		Quality:    80,
		Force:      true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected compressor to be called")
	}
	if svc.lastReq.Quality != 80 {
		t.Fatalf("unexpected quality: %d", svc.lastReq.Quality)
	}
	if svc.lastReq.InputPath != NormalizePath(`"./input.jpg"`) {
		t.Fatalf("unexpected input path: %q", svc.lastReq.InputPath)
	}
}

func TestParseQuality(t *testing.T) {
	if _, err := parseQuality("abc"); err == nil {
		t.Fatal("expected integer validation error")
	}
	if _, err := parseQuality("101"); err == nil {
		t.Fatal("expected range validation error")
	}
	value, err := parseQuality("80")
	if err != nil {
		t.Fatalf("parseQuality() error = %v", err)
	}
	if value != 80 {
		t.Fatalf("parseQuality() = %d, want 80", value)
	}
}

func TestInteractiveImageCompressionOutputUsesDefaultFolder(t *testing.T) {
	got, err := outputpkg.ResolveOutputPath("image.jpg", "", ResolveInteractiveOutputDir(""), "-compressed", "jpg", true)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	want := filepath.Clean(filepath.Join(DefaultInteractiveOutputDir, "image-compressed.jpg"))
	if got != want {
		t.Fatalf("ResolveOutputPath() = %q, want %q", got, want)
	}
}
