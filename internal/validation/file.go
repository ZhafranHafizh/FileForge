package validation

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EnsureInputFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("input path is required")
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("input file does not exist: %s", path)
		}
		return fmt.Errorf("stat input file %q: %w", path, err)
	}

	if err := EnsureRegularFile(path); err != nil {
		return fmt.Errorf("invalid input file: %w", err)
	}

	return nil
}

func EnsureRegularFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("file path is required")
	}

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("stat file %q: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, expected a regular file: %s", path)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("path is not a regular file: %s", path)
	}

	return nil
}

func EnsureOutputPath(path string, force bool) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("output path is required")
	}

	if err := EnsureOutputDir(path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("output path must be a regular file: %s", path)
		}
		if !force {
			return fmt.Errorf("output file already exists: %s (use --force to overwrite)", path)
		}
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat output path %q: %w", path, err)
	}

	return nil
}

func EnsureOutputDir(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("output path is required")
	}

	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", dir, err)
	}

	return nil
}

func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, fmt.Errorf("file does not exist: %s", path)
		}
		return 0, fmt.Errorf("stat file %q: %w", path, err)
	}

	if !info.Mode().IsRegular() {
		return 0, fmt.Errorf("path is not a regular file: %s", path)
	}

	return info.Size(), nil
}

func Extension(path string) string {
	return NormalizeExtension(filepath.Ext(path))
}
