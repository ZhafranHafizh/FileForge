package cmd

import "testing"

func TestToolPurpose(t *testing.T) {
	tests := map[string]string{
		"magick":   "required for image conversion, image compression, image to PDF",
		"qpdf":     "required for PDF merge and PDF split",
		"gs":       "required for PDF compression",
		"pdftoppm": "required for PDF to image",
		"pdfinfo":  "required for PDF info",
	}

	for name, want := range tests {
		if got := toolPurpose(name); got != want {
			t.Fatalf("toolPurpose(%q) = %q, want %q", name, got, want)
		}
	}
}
