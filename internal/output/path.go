package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"fileforge/internal/validation"
)

func EnsureOnlyOneOutputMode(out string, outputDir string) error {
	if strings.TrimSpace(out) != "" && strings.TrimSpace(outputDir) != "" {
		return fmt.Errorf("use either --out or --output-dir, not both")
	}
	return nil
}

func ResolveOutputPath(inputPath string, explicitOut string, outputDir string, suffix string, ext string, force bool) (string, error) {
	if err := EnsureOnlyOneOutputMode(explicitOut, outputDir); err != nil {
		return "", err
	}

	if strings.TrimSpace(explicitOut) != "" {
		path := filepath.Clean(explicitOut)
		if err := validation.EnsureOutputPath(path, force); err != nil {
			return "", err
		}
		return path, nil
	}

	if strings.TrimSpace(outputDir) == "" {
		return "", fmt.Errorf("either --out or --output-dir is required")
	}

	dir := filepath.Clean(outputDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory %q: %w", dir, err)
	}

	fileName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) + suffix + "." + validation.NormalizeExtension(ext)
	path := filepath.Join(dir, fileName)
	if err := validation.EnsureOutputPath(path, force); err != nil {
		return "", err
	}
	return path, nil
}

func ResolveOutputDir(explicitOut string, outputDir string, generatedDirName string) (string, error) {
	if err := EnsureOnlyOneOutputMode(explicitOut, outputDir); err != nil {
		return "", err
	}

	if strings.TrimSpace(explicitOut) != "" {
		path := filepath.Clean(explicitOut)
		if err := ensureDirectory(path); err != nil {
			return "", err
		}
		return path, nil
	}

	if strings.TrimSpace(outputDir) == "" {
		return "", fmt.Errorf("either --out or --output-dir is required")
	}

	root := filepath.Clean(outputDir)
	if err := ensureDirectory(root); err != nil {
		return "", err
	}

	path := filepath.Join(root, generatedDirName)
	if err := ensureDirectory(path); err != nil {
		return "", err
	}
	return path, nil
}

func GeneratedPDFToImageDirName(inputPath string) string {
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	return base + "-pages"
}

func GeneratedPDFSplitName(inputPath string, pageRange string) string {
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	return base + "-pages-" + SanitizeFilenamePart(pageRange)
}

func SanitizeFilenamePart(value string) string {
	var b strings.Builder
	lastUnderscore := false

	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '-', r == '_':
			b.WriteRune(r)
			lastUnderscore = false
		default:
			if !lastUnderscore {
				b.WriteRune('_')
				lastUnderscore = true
			}
		}
	}

	result := strings.Trim(b.String(), "_")
	if result == "" {
		return "output"
	}
	return result
}

func ensureDirectory(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("output path must be a directory: %s", path)
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("stat output directory %q: %w", path, err)
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", path, err)
	}
	return nil
}
