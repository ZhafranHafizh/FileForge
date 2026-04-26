package interactive

import (
	"strings"
	"testing"
)

func TestRenderComingSoonContainsFeatureName(t *testing.T) {
	got := RenderComingSoon("OCR")
	for _, part := range []string{"OCR", "This feature is not available yet.", "Planned for a future release."} {
		if !strings.Contains(got, part) {
			t.Fatalf("RenderComingSoon() missing %q in %q", part, got)
		}
	}
}

func TestRenderSuccessContainsOutputPathAndText(t *testing.T) {
	got := RenderSuccess("./output/result.pdf")
	for _, part := range []string{"The process is complete. Please access the following path to view the results.", "./output/result.pdf"} {
		if !strings.Contains(got, part) {
			t.Fatalf("RenderSuccess() missing %q in %q", part, got)
		}
	}
}

func TestRenderSummaryContainsTitleAndRows(t *testing.T) {
	got := RenderSummary("Summary", []SummaryRow{
		{Label: "Action", Value: "PDF to Image"},
		{Label: "Input", Value: "invoice.pdf"},
	})
	for _, part := range []string{"Summary", "Action", "PDF to Image", "Input", "invoice.pdf"} {
		if !strings.Contains(got, part) {
			t.Fatalf("RenderSummary() missing %q in %q", part, got)
		}
	}
}

func TestRenderSummaryContainsMultilineValues(t *testing.T) {
	got := RenderSummary("Summary", []SummaryRow{
		{Label: "Inputs", Value: "a.pdf\nb.pdf"},
	})
	for _, part := range []string{"Inputs", "a.pdf", "b.pdf"} {
		if !strings.Contains(got, part) {
			t.Fatalf("RenderSummary() missing %q in %q", part, got)
		}
	}
}
