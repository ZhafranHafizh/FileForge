package interactive

import (
	"fmt"
	"io"
	"path/filepath"

	"charm.land/huh/v2"
)

func printSuccess(w io.Writer, path string) error {
	if w == nil {
		return nil
	}
	_, err := fmt.Fprintln(w, RenderSuccess(path))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w)
	return err
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

func confirmSummary(w io.Writer, rows []SummaryRow) (bool, error) {
	if w != nil {
		_, _ = fmt.Fprintln(w, RenderSummary("Summary", rows))
		_, _ = fmt.Fprintln(w)
	}

	var confirm bool
	if err := huh.NewConfirm().
		Title("Proceed?").
		Affirmative("Run").
		Negative("Cancel").
		Value(&confirm).
		Run(); err != nil {
		return false, cancelIfInterrupted(err)
	}
	if !confirm {
		return false, ErrCancelled
	}
	return true, nil
}

func boolText(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
