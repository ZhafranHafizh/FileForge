package convert

import (
	"fmt"
	"strings"

	"fileforge/internal/validation"
)

type Route string

const (
	RouteImageToImage Route = "image_to_image"
	RoutePDFToImage   Route = "pdf_to_image"
	RouteImageToPDF   Route = "image_to_pdf"
)

var (
	supportedImageFormats    = []string{"jpg", "jpeg", "png", "webp"}
	supportedPDFImageFormats = []string{"jpg", "jpeg", "png"}
)

func DetectRoute(inputPath string, targetFormat string) (Route, error) {
	inputExt := validation.Extension(inputPath)
	target := validation.NormalizeExtension(targetFormat)

	switch {
	case inputExt == "pdf" && validation.IsSupportedExtension(target, supportedPDFImageFormats):
		return RoutePDFToImage, nil
	case validation.IsSupportedExtension(inputExt, supportedImageFormats) && target == "pdf":
		return RouteImageToPDF, nil
	case validation.IsSupportedExtension(inputExt, supportedImageFormats) && validation.IsSupportedExtension(target, supportedImageFormats):
		return RouteImageToImage, nil
	default:
		return "", ValidationError{Err: fmt.Errorf("unsupported conversion route: %s -> %s", inputExt, strings.TrimSpace(targetFormat))}
	}
}
