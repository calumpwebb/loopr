package git

import (
	"os/exec"
	"strings"
)

// CurrentBranch returns the current git branch name
func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Push pushes the current branch to origin with -u flag fallback
func Push(branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	err := cmd.Run()
	if err != nil {
		// Try with -u flag (create remote branch)
		cmd = exec.Command("git", "push", "-u", "origin", branch)
		return cmd.Run()
	}
	return nil
}

// IsGitRepo checks if the current directory is a git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}
