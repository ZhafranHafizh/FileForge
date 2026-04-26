package interactive

import "testing"

func TestIsCancelled(t *testing.T) {
	err := cancelIfInterrupted(assertErr("interrupted"))
	if !IsCancelled(err) {
		t.Fatal("expected cancelled error")
	}
}

func TestTopLevelMenuChoices(t *testing.T) {
	choices := topLevelMenuChoices()
	want := []string{"convert", "compress", "pdf_tools", "doctor", "coming_soon", "exit"}
	if len(choices) != len(want) {
		t.Fatalf("unexpected top-level choice count: %d", len(choices))
	}
	for i, value := range want {
		if choices[i].Value != value {
			t.Fatalf("choice %d = %q, want %q", i, choices[i].Value, value)
		}
	}
}

func TestSubmenuChoices(t *testing.T) {
	tests := []struct {
		name string
		got  []menuChoice
		want []string
	}{
		{
			name: "convert",
			got:  convertMenuChoices(),
			want: []string{"image_to_image", "pdf_to_image", "image_to_pdf", "back"},
		},
		{
			name: "compress",
			got:  compressMenuChoices(),
			want: []string{"image_compression", "pdf_compression", "back"},
		},
		{
			name: "pdf",
			got:  pdfMenuChoices(),
			want: []string{"pdf_merge", "pdf_split", "pdf_info", "back"},
		},
		{
			name: "coming soon",
			got:  comingSoonMenuChoices(),
			want: []string{"ocr", "office_conversion", "metadata_cleanup", "back"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.got) != len(tt.want) {
				t.Fatalf("len() = %d, want %d", len(tt.got), len(tt.want))
			}
			for i := range tt.want {
				if tt.got[i].Value != tt.want[i] {
					t.Fatalf("choice[%d] = %q, want %q", i, tt.got[i].Value, tt.want[i])
				}
			}
		})
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
