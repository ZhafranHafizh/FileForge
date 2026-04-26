package engine

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"fileforge/internal/runner"
)

type CommandRunner interface {
	Run(ctx context.Context, name string, args []string) (runner.Result, error)
}

type Tool struct {
	Name        string
	VersionArgs []string
}

type ToolResult struct {
	Name      string
	Available bool
	Version   string
	Err       error
}

type Report struct {
	Results []ToolResult
	Missing []ToolResult
}

type InstallHint struct {
	MacOS   string
	Ubuntu  string
	Windows string
}

var versionPattern = regexp.MustCompile(`\d+(?:\.\d+)+(?:-\d+)?`)

func RequiredTools() []Tool {
	return []Tool{
		{Name: "qpdf", VersionArgs: []string{"--version"}},
		{Name: "gs", VersionArgs: []string{"--version"}},
		{Name: "pdftoppm", VersionArgs: []string{"-v"}},
		{Name: "pdfinfo", VersionArgs: []string{"-v"}},
		{Name: "magick", VersionArgs: []string{"-version"}},
	}
}

func Diagnose(ctx context.Context, runner CommandRunner, tools []Tool) Report {
	results := make([]ToolResult, 0, len(tools))
	missing := make([]ToolResult, 0)

	for _, tool := range tools {
		version, err := VersionWithRunner(ctx, runner, tool.Name, tool.VersionArgs...)
		if err != nil {
			result := ToolResult{Name: tool.Name, Available: false, Err: err}
			results = append(results, result)
			missing = append(missing, result)
			continue
		}

		results = append(results, ToolResult{
			Name:      tool.Name,
			Available: true,
			Version:   version,
		})
	}

	return Report{
		Results: results,
		Missing: missing,
	}
}

func IsAvailable(binary string) bool {
	return IsAvailableWithRunner(context.Background(), runner.New(runner.Options{}), binary)
}

func IsAvailableWithRunner(ctx context.Context, runner CommandRunner, binary string) bool {
	_, err := VersionWithRunner(ctx, runner, binary, "--version")
	return err == nil
}

func Version(binary string, args ...string) (string, error) {
	return VersionWithRunner(context.Background(), runner.New(runner.Options{}), binary, args...)
}

func VersionWithRunner(ctx context.Context, runner CommandRunner, binary string, args ...string) (string, error) {
	result, err := runner.Run(ctx, binary, args)
	if err != nil {
		return "", fmt.Errorf("%s not available: %w", binary, err)
	}

	raw := strings.TrimSpace(strings.Join([]string{result.Stdout, result.Stderr}, "\n"))
	if raw == "" {
		return "", fmt.Errorf("no version output from %s", binary)
	}

	version := extractVersion(raw)
	if version == "" {
		return "", fmt.Errorf("unable to parse %s version from %q", binary, firstLine(raw))
	}

	return version, nil
}

func Require(binary string) error {
	return RequireWithRunner(context.Background(), runner.New(runner.Options{}), binary, "--version")
}

func RequireWithRunner(ctx context.Context, runner CommandRunner, binary string, args ...string) error {
	_, err := VersionWithRunner(ctx, runner, binary, args...)
	if err != nil {
		return err
	}

	return nil
}

func InstallHints(binary string) InstallHint {
	switch filepath.Base(binary) {
	case "qpdf":
		return InstallHint{"brew install qpdf", "sudo apt install qpdf", "winget install QPDF.QPDF"}
	case "gs":
		return InstallHint{"brew install ghostscript", "sudo apt install ghostscript", "winget install ArtifexSoftware.GhostScript"}
	case "pdftoppm", "pdfinfo":
		return InstallHint{"brew install poppler", "sudo apt install poppler-utils", "winget install oschwartz10612.Poppler"}
	case "magick":
		return InstallHint{"brew install imagemagick", "sudo apt install imagemagick", "winget install ImageMagick.ImageMagick"}
	default:
		return InstallHint{"See project documentation", "See project documentation", "See project documentation"}
	}
}

func extractVersion(raw string) string {
	return versionPattern.FindString(raw)
}

func firstLine(raw string) string {
	lines := strings.Split(raw, "\n")
	if len(lines) == 0 {
		return ""
	}
	return strings.TrimSpace(lines[0])
}
