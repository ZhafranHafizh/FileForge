package convert

import "testing"

func TestDetectRoute(t *testing.T) {
	tests := []struct {
		input string
		to    string
		want  Route
	}{
		{input: "input.png", to: "jpg", want: RouteImageToImage},
		{input: "input.pdf", to: "jpg", want: RoutePDFToImage},
		{input: "input.pdf", to: "png", want: RoutePDFToImage},
		{input: "input.jpg", to: "pdf", want: RouteImageToPDF},
	}

	for _, tt := range tests {
		got, err := DetectRoute(tt.input, tt.to)
		if err != nil {
			t.Fatalf("DetectRoute(%q, %q) error = %v", tt.input, tt.to, err)
		}
		if got != tt.want {
			t.Fatalf("DetectRoute(%q, %q) = %q, want %q", tt.input, tt.to, got, tt.want)
		}
	}
}

func TestDetectRouteRejectsUnsupported(t *testing.T) {
	if _, err := DetectRoute("input.pdf", "webp"); err == nil {
		t.Fatal("expected unsupported route error")
	}
}
