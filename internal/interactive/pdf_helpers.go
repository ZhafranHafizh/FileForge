package interactive

import (
	"fmt"
	"strings"

	"fileforge/internal/validation"
)

func ParsePathList(raw string) ([]string, error) {
	parts := strings.Split(raw, ",")
	paths := make([]string, 0, len(parts))

	for _, part := range parts {
		path := NormalizePath(part)
		if path == "" {
			return nil, fmt.Errorf("path list contains an empty entry")
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func validatePDFPath(path string) error {
	if path == "" {
		return fmt.Errorf("input path is required")
	}
	if err := validation.EnsureInputFile(path); err != nil {
		return err
	}
	return validation.EnsureSupportedExtension(path, []string{"pdf"})
}
