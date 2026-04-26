package compress

import "fileforge/internal/validation"

type Kind string

const (
	KindImage Kind = "image"
	KindPDF   Kind = "pdf"
)

func DetectKind(path string) Kind {
	if validation.Extension(path) == "pdf" {
		return KindPDF
	}
	return KindImage
}
