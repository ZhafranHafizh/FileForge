package interactive

import (
	"context"
	"testing"

	"fileforge/internal/convert"
)

type fakePDFToImageConverter struct {
	lastReq convert.PDFToImageRequest
	called  bool
	err     error
}

func (f *fakePDFToImageConverter) Convert(_ context.Context, req convert.PDFToImageRequest) error {
	f.called = true
	f.lastReq = req
	return f.err
}

type fakeImageToPDFConverter struct {
	lastReq convert.ImageToPDFRequest
	called  bool
	err     error
}

func (f *fakeImageToPDFConverter) Convert(_ context.Context, req convert.ImageToPDFRequest) error {
	f.called = true
	f.lastReq = req
	return f.err
}

func TestPDFToImageWizardExecuteUsesService(t *testing.T) {
	svc := &fakePDFToImageConverter{}
	wizard := NewPDFToImageWizard(svc)
	err := wizard.Execute(context.Background(), PDFToImageInput{
		InputPath: `"./input.pdf"`,
		OutputDir: `'./pages'`,
		ToFormat:  "png",
		DPI:       150,
		FirstPage: 1,
		LastPage:  2,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected converter to be called")
	}
}

func TestImageToPDFWizardExecuteUsesService(t *testing.T) {
	svc := &fakeImageToPDFConverter{}
	wizard := NewImageToPDFWizard(svc)
	err := wizard.Execute(context.Background(), ImageToPDFInput{
		InputPath:  `"./input.jpg"`,
		OutputPath: `'./output.pdf'`,
		Force:      true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected converter to be called")
	}
}

func TestParsePositiveInteger(t *testing.T) {
	if _, err := parsePositiveInteger("0", "dpi"); err == nil {
		t.Fatal("expected positive integer error")
	}
	if _, err := parseOptionalPositiveInteger("", "first page"); err != nil {
		t.Fatalf("parseOptionalPositiveInteger() error = %v", err)
	}
	if err := validatePageBounds(3, 1); err == nil {
		t.Fatal("expected page bounds error")
	}
}
