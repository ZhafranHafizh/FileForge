package interactive

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParsePathList(t *testing.T) {
	input := `"C:\Users\zhafr\Docs\a.pdf", "C:\Users\zhafr\Docs\b file.pdf"`
	got, err := ParsePathList(input)
	if err != nil {
		t.Fatalf("ParsePathList() error = %v", err)
	}

	want := []string{
		filepath.Clean(`C:\Users\zhafr\Docs\a.pdf`),
		filepath.Clean(`C:\Users\zhafr\Docs\b file.pdf`),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParsePathList() = %v, want %v", got, want)
	}
}

func TestParsePathListRejectsEmptyEntry(t *testing.T) {
	if _, err := ParsePathList(`"a.pdf", , "b.pdf"`); err == nil {
		t.Fatal("expected empty entry error")
	}
}
