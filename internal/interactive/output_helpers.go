package interactive

import (
	"io"
	"path/filepath"

	"charm.land/huh/v2"
	outputpkg "fileforge/internal/output"
)

func printSuccess(w io.Writer, path string) error {
	if w == nil {
		return nil
	}
	return outputpkg.PrintSuccess(w, path)
}

func confirmGeneratedDirOverwrite(dir string, format string) (bool, error) {
	pattern := filepath.Join(dir, "page-*."+format)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false, err
	}
	if len(matches) == 0 {
		return false, nil
	}
	var overwrite bool
	if err := huh.NewConfirm().
		Title("Output directory already contains generated files. Overwrite?").
		Affirmative("Overwrite").
		Negative("Cancel").
		Value(&overwrite).
		Run(); err != nil {
		return false, cancelIfInterrupted(err)
	}
	if !overwrite {
		return false, ErrCancelled
	}
	return true, nil
}
