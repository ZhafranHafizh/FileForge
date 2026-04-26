package validation

import (
	"fmt"
	"strings"
)

func NormalizeExtension(ext string) string {
	normalized := strings.TrimSpace(strings.ToLower(ext))
	return strings.TrimPrefix(normalized, ".")
}

func IsSupportedExtension(ext string, supported []string) bool {
	normalized := NormalizeExtension(ext)
	if normalized == "" {
		return false
	}

	for _, candidate := range supported {
		if normalized == NormalizeExtension(candidate) {
			return true
		}
	}

	return false
}

func EnsureSupportedExtension(path string, supported []string) error {
	ext := Extension(path)
	if IsSupportedExtension(ext, supported) {
		return nil
	}

	return fmt.Errorf("unsupported file extension %q; supported extensions: %s", ext, strings.Join(normalizeExtensions(supported), ", "))
}

func normalizeExtensions(values []string) []string {
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, value := range values {
		ext := NormalizeExtension(value)
		if ext == "" {
			continue
		}
		if _, ok := seen[ext]; ok {
			continue
		}
		seen[ext] = struct{}{}
		normalized = append(normalized, ext)
	}

	return normalized
}
