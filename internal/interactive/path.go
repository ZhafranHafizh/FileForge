package interactive

import (
	"path/filepath"
	"strings"
)

func NormalizePath(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.Trim(trimmed, `"`)
	trimmed = strings.Trim(trimmed, `'`)
	if trimmed == "" {
		return ""
	}
	return filepath.Clean(trimmed)
}
