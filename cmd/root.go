package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	interactivepkg "fileforge/internal/interactive"
	"fileforge/internal/runner"

	"github.com/spf13/cobra"
)

const (
	ExitSuccess           = 0
	ExitGeneralError      = 1
	ExitInvalidInput      = 2
	ExitMissingDependency = 3
	ExitConversionFailed  = 4
	ExitCompressionFailed = 5
)

type GlobalOptions struct {
	Verbose bool
	Quiet   bool
	Force   bool
	Config  string
}

type exitCoder interface {
	error
	ExitCode() int
}

type commandError struct {
	code int
	err  error
}

func (e commandError) Error() string {
	return e.err.Error()
}

func (e commandError) Unwrap() error {
	return e.err
}

func (e commandError) ExitCode() int {
	return e.code
}

var rootOpts GlobalOptions

var rootCmd = &cobra.Command{
	Use:           "fileforge",
	Short:         "Local-first file conversion, compression, and PDF utility CLI",
	Long:          "FileForge is an offline-capable CLI for file conversion, compression, and PDF utilities using local binaries.",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractiveApp(cmd)
	},
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		printError(rootCmd.ErrOrStderr(), err)
		return exitCodeFor(err)
	}

	return ExitSuccess
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&rootOpts.Verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&rootOpts.Quiet, "quiet", false, "suppress non-error output")
	rootCmd.PersistentFlags().BoolVar(&rootOpts.Force, "force", false, "overwrite existing output files")
	rootCmd.PersistentFlags().StringVar(&rootOpts.Config, "config", "", "path to config file")
}

func printError(w io.Writer, err error) {
	if err == nil {
		return
	}

	if w == nil {
		w = os.Stderr
	}

	_, _ = fmt.Fprintf(w, "Error: %v\n", err)
}

func exitCodeFor(err error) int {
	if err == nil {
		return ExitSuccess
	}

	var coded exitCoder
	if errors.As(err, &coded) {
		return coded.ExitCode()
	}

	return ExitGeneralError
}

func newCommandError(code int, err error) error {
	if err == nil {
		return nil
	}

	return commandError{code: code, err: err}
}

func runInteractiveApp(cmd *cobra.Command) error {
	run := runner.New(runner.Options{
		Verbose: rootOpts.Verbose,
		Stdout:  cmd.ErrOrStderr(),
		Stderr:  cmd.ErrOrStderr(),
	})

	app := interactivepkg.NewApp(interactivepkg.Options{
		Runner:       run,
		Doctor:       func(ctx context.Context, out io.Writer) error { return runDoctor(ctx, run, out) },
		Stdout:       cmd.OutOrStdout(),
		ImageForce:   rootOpts.Force,
		AccessibleUI: false,
	})

	if err := app.Run(context.Background()); err != nil {
		if interactivepkg.IsCancelled(err) {
			return nil
		}
		if interactivepkg.IsValidationError(err) {
			return newCommandError(ExitInvalidInput, err)
		}
		if interactivepkg.IsDependencyError(err) {
			return newCommandError(ExitMissingDependency, err)
		}
		return newCommandError(ExitGeneralError, err)
	}

	return nil
}
