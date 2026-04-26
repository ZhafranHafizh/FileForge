package pdf

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildMergePlanArgs(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.pdf")
	b := filepath.Join(dir, "b.pdf")
	out := filepath.Join(dir, "merged.pdf")
	for _, path := range []string{a, b} {
		if err := os.WriteFile(path, []byte("pdf"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
	}

	plan, err := BuildMergePlan(MergeRequest{
		InputPaths: []string{a, b},
		OutputPath: out,
	})
	if err != nil {
		t.Fatalf("BuildMergePlan() error = %v", err)
	}

	want := []string{"--empty", "--pages", a, b, "--", out}
	if !reflect.DeepEqual(plan.Args(), want) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), want)
	}
}

func TestBuildMergePlanRequiresTwoFiles(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.pdf")
	if err := os.WriteFile(a, []byte("pdf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := BuildMergePlan(MergeRequest{
		InputPaths: []string{a},
		OutputPath: filepath.Join(dir, "merged.pdf"),
	})
	if err == nil {
		t.Fatal("expected error for single input")
	}
}

func TestBuildMergePlanPreservesInputOrder(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.pdf")
	b := filepath.Join(dir, "b.pdf")
	c := filepath.Join(dir, "c.pdf")
	out := filepath.Join(dir, "merged.pdf")
	for _, path := range []string{a, b, c} {
		if err := os.WriteFile(path, []byte("pdf"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
	}

	plan, err := BuildMergePlan(MergeRequest{
		InputPaths: []string{c, a, b},
		OutputPath: out,
	})
	if err != nil {
		t.Fatalf("BuildMergePlan() error = %v", err)
	}

	want := []string{"--empty", "--pages", c, a, b, "--", out}
	if !reflect.DeepEqual(plan.Args(), want) {
		t.Fatalf("Args() = %v, want %v", plan.Args(), want)
	}
}
