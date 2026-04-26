package interactive

import (
	"path/filepath"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "trim spaces", in: "  ./folder/file.png  ", want: filepath.Clean("./folder/file.png")},
		{name: "double quotes", in: `"C:\temp\image.png"`, want: filepath.Clean(`C:\temp\image.png`)},
		{name: "single quotes", in: `'./images/photo.jpg'`, want: filepath.Clean("./images/photo.jpg")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePath(tt.in); got != tt.want {
				t.Fatalf("NormalizePath() = %q, want %q", got, tt.want)
			}
		})
	}
}
