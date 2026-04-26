package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type Options struct {
	Verbose bool
	Stdout  io.Writer
	Stderr  io.Writer
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type Runner struct {
	verbose bool
	stdout  io.Writer
	stderr  io.Writer
}

func New(opts Options) *Runner {
	return &Runner{
		verbose: opts.Verbose,
		stdout:  opts.Stdout,
		stderr:  opts.Stderr,
	}
}

func (r *Runner) Run(ctx context.Context, name string, args []string) (Result, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if r.verbose {
		if r.stdout != nil {
			_, _ = fmt.Fprintf(r.stdout, "running: %s %v\n", name, args)
		}
	}

	err := cmd.Run()

	result := Result{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: 0,
	}

	if err == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		return result, fmt.Errorf("command %q failed: %w", name, err)
	}

	result.ExitCode = -1
	return result, fmt.Errorf("command %q failed: %w", name, err)
}
