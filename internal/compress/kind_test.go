package compress

import "testing"

func TestDetectKind(t *testing.T) {
	tests := []struct {
		path string
		want Kind
	}{
		{path: "input.pdf", want: KindPDF},
		{path: "input.PNG", want: KindImage},
	}

	for _, tt := range tests {
		if got := DetectKind(tt.path); got != tt.want {
			t.Fatalf("DetectKind(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}
