package gitdiff

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func DefaultChangedFiles() ([]string, error) {
	if base := os.Getenv("GITHUB_BASE_REF"); base != "" {
		return changedFilesAgainstBase(base)
	}

	if before := os.Getenv("GITHUB_EVENT_BEFORE"); before != "" && before != "0000000000000000000000000000000000000000" {
		return changedFilesBetween(before, "HEAD")
	}

	return []string{}, nil
}

func changedFilesAgainstBase(base string) ([]string, error) {
	if err := exec.Command("git", "fetch", "origin", base, "--depth=1").Run(); err != nil {
		return nil, fmt.Errorf("fetch base branch %q: %w", base, err)
	}
	return changedFilesBetween("origin/"+base, "HEAD")
}

func changedFilesBetween(from, to string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=d", from, to)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git diff %s..%s: %w: %s", from, to, err, strings.TrimSpace(out.String()))
	}

	raw := strings.TrimSpace(out.String())
	if raw == "" {
		return []string{}, nil
	}

	lines := strings.Split(raw, "\n")
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}