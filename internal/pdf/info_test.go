package pdf

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildInfoPlanArgs(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "document.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildInfoPlan(InfoRequest{InputPath: input})
	if err != nil {
		t.Fatalf("BuildInfoPlan() error = %v", err)
	}

	want := []string{input}
	if !reflect.DeepEqual(plan.Args(), want) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), want)
	}
}

func TestExtractPages(t *testing.T) {
	raw := "Title: sample\nPages: 12\nEncrypted: no"
	pages := extractPages(raw)
	if pages == nil {
		t.Fatal("expected page count")
	}
	if *pages != 12 {
		t.Fatalf("pages = %d, want 12", *pages)
	}
}
