package output

import (
	"bytes"
	"testing"
)

func TestPrintSuccess(t *testing.T) {
	var buf bytes.Buffer
	if err := PrintSuccess(&buf, "./output/result.pdf"); err != nil {
		t.Fatalf("PrintSuccess() error = %v", err)
	}
	want := "The process is complete. Please access the following path to view the results.\n./output/result.pdf\n"
	if buf.String() != want {
		t.Fatalf("PrintSuccess() = %q", buf.String())
	}
}
