package importer

import (
	"os/exec"
	"strings"
)

// GitHead returns the current HEAD commit hash for the repo at dir.
func GitHead(dir string) string {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

// GitPull runs git pull --ff-only in the given directory.
func GitPull(dir string) error {
	cmd := exec.Command("git", "-C", dir, "pull", "--ff-only", "origin", "main")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// GitChangedFiles returns the list of files changed between two commits.
func GitChangedFiles(dir, fromCommit, toCommit string) ([]string, error) {
	out, err := exec.Command("git", "-C", dir, "diff", "--name-only", fromCommit, toCommit).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
	for _, l := range lines {
		if l != "" {
			result = append(result, l)
		}
	}
	return result, nil
}
