package engine

import (
	"context"
	"fmt"
)

var ghostscriptCandidates = []string{"gs", "gswin64c", "gswin32c"}

func GhostscriptCandidates() []string {
	candidates := make([]string, len(ghostscriptCandidates))
	copy(candidates, ghostscriptCandidates)
	return candidates
}

func DetectGhostscriptBinary(ctx context.Context, runner CommandRunner) (string, error) {
	for _, candidate := range ghostscriptCandidates {
		if _, err := VersionWithRunner(ctx, runner, candidate, "--version"); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("ghostscript not available; tried: %v", ghostscriptCandidates)
}
