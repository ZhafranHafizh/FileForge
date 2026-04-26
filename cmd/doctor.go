package cmd

import (
	"context"
	"fmt"
	"io"

	"fileforge/internal/engine"
	"fileforge/internal/runner"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check availability of required local engine binaries",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		run := runner.New(runner.Options{
			Verbose: rootOpts.Verbose,
			Stdout:  cmd.ErrOrStderr(),
			Stderr:  cmd.ErrOrStderr(),
		})

		if err := runDoctor(ctx, run, cmd.OutOrStdout()); err != nil {
			if len(diagnoseRequiredTools(ctx, run).Missing) > 0 {
				return newCommandError(ExitMissingDependency, err)
			}
			return err
		}

		return nil
	},
}

func runDoctor(ctx context.Context, run *runner.Runner, out io.Writer) error {
	report := diagnoseRequiredTools(ctx, run)

	if !rootOpts.Quiet {
		_, _ = fmt.Fprintln(out, "FileForge Doctor")
		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintln(out, "Required tools:")

		for _, result := range report.Results {
			if result.Available {
				_, _ = fmt.Fprintf(out, "[OK] %s %s\n", result.Name, result.Version)
				continue
			}

			_, _ = fmt.Fprintf(out, "[Missing] %s\n", result.Name)
		}

		if len(report.Missing) > 0 {
			_, _ = fmt.Fprintln(out)
			for _, missing := range report.Missing {
				hint := engine.InstallHints(missing.Name)
				_, _ = fmt.Fprintf(out, "Missing dependency: %s\n", missing.Name)
				_, _ = fmt.Fprintf(out, "macOS:   %s\n", hint.MacOS)
				_, _ = fmt.Fprintf(out, "Ubuntu:  %s\n", hint.Ubuntu)
				_, _ = fmt.Fprintf(out, "Windows: %s\n", hint.Windows)
				_, _ = fmt.Fprintln(out)
			}
		}

		_, _ = fmt.Fprintln(out, "Status:")
		if len(report.Missing) == 0 {
			_, _ = fmt.Fprintln(out, "Ready for Milestone 2 features.")
		} else {
			_, _ = fmt.Fprintln(out, "Missing required dependencies.")
		}
	}

	if len(report.Missing) > 0 {
		return fmt.Errorf("required dependencies are missing")
	}

	return nil
}

func diagnoseRequiredTools(ctx context.Context, run *runner.Runner) engine.Report {
	tools := []engine.Tool{
		{Name: "qpdf", VersionArgs: []string{"--version"}},
		{Name: "pdftoppm", VersionArgs: []string{"-v"}},
		{Name: "pdfinfo", VersionArgs: []string{"-v"}},
		{Name: "magick", VersionArgs: []string{"-version"}},
	}

	report := engine.Diagnose(ctx, run, tools)

	gsBinary, err := engine.DetectGhostscriptBinary(ctx, run)
	gsResult := engine.ToolResult{Name: "gs", Available: false, Err: err}
	if err == nil {
		version, versionErr := engine.VersionWithRunner(ctx, run, gsBinary, "--version")
		if versionErr == nil {
			gsResult.Available = true
			gsResult.Version = version
		} else {
			gsResult.Err = versionErr
		}
	}

	results := make([]engine.ToolResult, 0, len(report.Results)+1)
	results = append(results, report.Results[0:1]...)
	results = append(results, gsResult)
	results = append(results, report.Results[1:]...)

	missing := make([]engine.ToolResult, 0, len(report.Missing)+1)
	for _, result := range results {
		if !result.Available {
			missing = append(missing, result)
		}
	}

	return engine.Report{
		Results: results,
		Missing: missing,
	}
}
