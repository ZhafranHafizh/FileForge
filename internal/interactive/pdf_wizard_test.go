package interactive

import (
	"context"
	"path/filepath"
	"testing"

	"fileforge/internal/compress"
	outputpkg "fileforge/internal/output"
	"fileforge/internal/pdf"
)

type fakePDFCompressor struct {
	lastReq compress.PDFCompressRequest
	called  bool
	err     error
}

func (f *fakePDFCompressor) Compress(_ context.Context, req compress.PDFCompressRequest) (compress.PDFCompressResult, error) {
	f.called = true
	f.lastReq = req
	return compress.PDFCompressResult{}, f.err
}

type fakePDFMerger struct {
	lastReq pdf.MergeRequest
	called  bool
	err     error
}

func (f *fakePDFMerger) Merge(_ context.Context, req pdf.MergeRequest) error {
	f.called = true
	f.lastReq = req
	return f.err
}

type fakePDFSplitter struct {
	lastReq pdf.SplitRequest
	called  bool
	err     error
}

func (f *fakePDFSplitter) Split(_ context.Context, req pdf.SplitRequest) error {
	f.called = true
	f.lastReq = req
	return f.err
}

type fakePDFInspector struct {
	lastReq pdf.InfoRequest
	called  bool
	err     error
}

func (f *fakePDFInspector) Info(_ context.Context, req pdf.InfoRequest) (pdf.InfoResult, error) {
	f.called = true
	f.lastReq = req
	return pdf.InfoResult{}, f.err
}

func TestPDFCompressWizardExecuteUsesService(t *testing.T) {
	svc := &fakePDFCompressor{}
	wizard := NewPDFCompressWizard(svc, nil)

	err := wizard.Execute(context.Background(), PDFCompressInput{
		InputPath:  `"./input.pdf"`,
		Preset:     "ebook",
		OutputPath: `'./output.pdf'`,
		Force:      true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected compressor to be called")
	}
	if svc.lastReq.InputPath != NormalizePath(`"./input.pdf"`) {
		t.Fatalf("unexpected input path: %q", svc.lastReq.InputPath)
	}
}

func TestPDFMergeWizardExecuteUsesService(t *testing.T) {
	svc := &fakePDFMerger{}
	wizard := NewPDFMergeWizard(svc, nil)

	err := wizard.Execute(context.Background(), PDFMergeInput{
		InputPaths: []string{`"./a.pdf"`, `'./b.pdf'`},
		OutputPath: `'./merged.pdf'`,
		Force:      true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected merger to be called")
	}
	if len(svc.lastReq.InputPaths) != 2 {
		t.Fatalf("expected 2 input paths, got %d", len(svc.lastReq.InputPaths))
	}
}

func TestPDFSplitWizardExecuteUsesService(t *testing.T) {
	svc := &fakePDFSplitter{}
	wizard := NewPDFSplitWizard(svc, nil)

	err := wizard.Execute(context.Background(), PDFSplitInput{
		InputPath:  `"./input.pdf"`,
		PageRange:  "1-3",
		OutputPath: `'./chapter.pdf'`,
		Force:      true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected splitter to be called")
	}
	if svc.lastReq.PageRange != "1-3" {
		t.Fatalf("unexpected page range: %q", svc.lastReq.PageRange)
	}
}

func TestPDFInfoWizardExecuteUsesService(t *testing.T) {
	svc := &fakePDFInspector{}
	wizard := NewPDFInfoWizard(svc, nil)

	err := wizard.Execute(context.Background(), PDFInfoInput{
		InputPath: `"./input.pdf"`,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected inspector to be called")
	}
}

func TestInteractivePDFSplitOutputUsesDefaultFolder(t *testing.T) {
	got, err := outputpkg.ResolveOutputPath("book.pdf", "", ResolveInteractiveOutputDir(""), "-pages-"+outputpkg.SanitizeFilenamePart("1-3,7,10-12"), "pdf", true)
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	want := filepath.Clean(filepath.Join(DefaultInteractiveOutputDir, "book-pages-1-3_7_10-12.pdf"))
	if got != want {
		t.Fatalf("ResolveOutputPath() = %q, want %q", got, want)
	}
}
