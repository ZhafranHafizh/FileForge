package interactive

import (
	"path/filepath"
	"strings"
)

const DefaultInteractiveOutputDir = "FileForge-Output"

func NormalizePath(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.Trim(trimmed, `"`)
	trimmed = strings.Trim(trimmed, `'`)
	if trimmed == "" {
		return ""
	}
	return filepath.Clean(trimmed)
}

func ResolveInteractiveOutputDir(raw string) string {
	normalized := NormalizePath(raw)
	if normalized == "" {
		return filepath.Clean(DefaultInteractiveOutputDir)
	}
	return normalized
}
