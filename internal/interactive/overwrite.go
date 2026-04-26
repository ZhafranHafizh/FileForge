package interactive

import (
	"os"

	"charm.land/huh/v2"
)

func confirmOverwrite(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	var overwrite bool
	if err := huh.NewConfirm().
		Title("Output already exists. Overwrite?").
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
