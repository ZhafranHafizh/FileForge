package pdf

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildSplitPlanArgs(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	output := filepath.Join(dir, "chapter.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	plan, err := BuildSplitPlan(SplitRequest{
		InputPath:  input,
		PageRange:  "1-3",
		OutputPath: output,
	})
	if err != nil {
		t.Fatalf("BuildSplitPlan() error = %v", err)
	}

	want := []string{input, "--pages", input, "1-3", "--", output}
	if !reflect.DeepEqual(plan.Args(), want) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), want)
	}
}

func TestBuildSplitPlanRejectsInvalidPageRange(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.pdf")
	output := filepath.Join(dir, "chapter.pdf")
	if err := os.WriteFile(input, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildSplitPlan(SplitRequest{
		InputPath:  input,
		PageRange:  "5-1",
		OutputPath: output,
	})
	if err == nil {
		t.Fatal("expected invalid page range error")
	}
}
