package validation

import "testing"

func TestExtension(t *testing.T) {
	if got := Extension("sample.PDF"); got != "pdf" {
		t.Fatalf("Extension() = %q, want pdf", got)
	}
}

func TestNormalizeExtension(t *testing.T) {
	if got := NormalizeExtension(" .JpG "); got != "jpg" {
		t.Fatalf("NormalizeExtension() = %q, want jpg", got)
	}
}

func TestIsSupportedExtension(t *testing.T) {
	supported := []string{"jpg", ".png", "webp"}

	if !IsSupportedExtension(".PNG", supported) {
		t.Fatal("expected extension to be supported")
	}
	if IsSupportedExtension("pdf", supported) {
		t.Fatal("expected extension to be unsupported")
	}
}

func TestEnsureSupportedExtension(t *testing.T) {
	if err := EnsureSupportedExtension("photo.jpeg", []string{"jpg", "jpeg", "png"}); err != nil {
		t.Fatalf("EnsureSupportedExtension() error = %v", err)
	}

	if err := EnsureSupportedExtension("document.pdf", []string{"jpg", "png"}); err == nil {
		t.Fatal("expected unsupported extension error")
	}
}
