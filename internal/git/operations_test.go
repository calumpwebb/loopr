package git

import (
	"testing"
)

func TestIsGitRepo(t *testing.T) {
	// Since we're running this in the loopr directory which is a git repo
	if !IsGitRepo() {
		t.Error("Expected IsGitRepo to return true in a git repository")
	}
}

func TestCurrentBranch(t *testing.T) {
	branch, err := CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch failed: %v", err)
	}
	if branch == "" {
		t.Error("Expected non-empty branch name")
	}
}
