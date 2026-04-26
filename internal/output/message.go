package output

import (
	"fmt"
	"io"
)

func PrintSuccess(w io.Writer, outputPath string) error {
	_, err := fmt.Fprintf(w, "The process is complete. Please access the following path to view the results.\n%s\n", outputPath)
	return err
}
