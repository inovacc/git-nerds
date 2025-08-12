package stats

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Helper function to create a dummy git repository for testing
func createTestRepo(t *testing.T) string {
	path := t.TempDir()
	if _, err := git.PlainInit(path, false); err != nil {
		t.Fatalf("failed to initialize test repo: %v", err)
	}

	return path
}

func TestFindGit(t *testing.T) {
	// Case 1: Git binary exists and is functional
	// This test assumes 'git' is installed on the system where the test is run.
	t.Run("GitExists", func(t *testing.T) {
		gitPath, err := FindGit()
		if err != nil {
			t.Errorf("FindGit() failed with error: %v", err)
		}

		if gitPath == "" {
			t.Error("FindGit() returned an empty path when git is expected to exist")
		}

		t.Logf("Found git at: %s", gitPath)
	})

	// Case 2: Git binary is not in PATH
	t.Run("GitNotFound", func(t *testing.T) {
		// Save the current PATH
		oldPath := os.Getenv("PATH")
		defer func(key, value string) {
			_ = os.Setenv(key, value)
		}("PATH", oldPath)

		// Set a temporary empty PATH to simulate git not being found
		_ = os.Setenv("PATH", "")

		gitPath, err := FindGit()
		if err == nil {
			t.Errorf("FindGit() succeeded when git was not in PATH, returned path: %s", gitPath)
		}

		if gitPath != "" {
			t.Errorf("FindGit() returned non-empty path '%s' when git was not expected to be found", gitPath)
		}

		if !strings.Contains(err.Error(), "local git binary not found") {
			t.Errorf("expected error message to contain 'local git binary not found', got: %v", err)
		}
	})
}

func TestOpenRepo(t *testing.T) {
	// Case 1: Open a valid repository
	t.Run("ValidRepo", func(t *testing.T) {
		repoPath := createTestRepo(t)

		repo, err := OpenRepo(repoPath)
		if err != nil {
			t.Errorf("OpenRepo() failed with error: %v", err)
		}

		if repo == nil {
			t.Error("OpenRepo() returned nil repo for a valid path")
		}

		// Verify the repository can be used
		if _, err = repo.CommitObjects(); err != nil {
			t.Errorf("failed to get commits from opened repo: %v", err)
		}
	})

	// Case 2: Open a non-existent repository
	t.Run("InvalidPath", func(t *testing.T) {
		invalidPath := filepath.Join(t.TempDir(), "non-existent-repo")

		repo, err := OpenRepo(invalidPath)
		if err == nil {
			t.Error("OpenRepo() succeeded for a non-existent path")
		}

		if repo != nil {
			t.Errorf("OpenRepo() returned a non-nil repo for an invalid path: %v", repo)
		}

		if !strings.Contains(err.Error(), "failed to open repo with go-git") {
			t.Errorf("expected error message to contain 'failed to open repo with go-git', got: %v", err)
		}
	})
}

func TestCloneWithLocalGit(t *testing.T) {
	// This test uses a dummy remote repository created locally.
	// You might want to use a known public repository for more reliable testing.
	// Example: "[https://github.com/go-git/go-git.git](https://github.com/go-git/go-git.git)"

	// Create a dummy "remote" repository
	remotePath := createTestRepo(t)
	remoteRepo, _ := git.PlainOpen(remotePath)
	wt, _ := remoteRepo.Worktree()
	file, _ := os.Create(filepath.Join(remotePath, "README.md"))
	_ = file.Close()
	_, _ = wt.Add("README.md")
	_, _ = wt.Commit("initial commit", &git.CommitOptions{Author: &object.Signature{}})

	// Case 1: Local git binary found, should clone using exec.Command
	t.Run("CloneWithLocalGit", func(t *testing.T) {
		cloneDir := filepath.Join(t.TempDir(), "local-clone")
		if err := CloneWithLocalGit(remotePath, cloneDir); err != nil {
			t.Errorf("CloneWithLocalGit() failed with error: %v", err)
		}

		if _, err := os.Stat(filepath.Join(cloneDir, "README.md")); os.IsNotExist(err) {
			t.Error("file from cloned repo not found")
		}
	})

	// Case 2: Local git binary not found, should fall back to go-git
	t.Run("CloneWithGoGitFallback", func(t *testing.T) {
		oldPath := os.Getenv("PATH")
		defer func(key, value string) {
			_ = os.Setenv(key, value)
		}("PATH", oldPath)

		_ = os.Setenv("PATH", "")

		cloneDir := filepath.Join(t.TempDir(), "go-git-clone")
		if err := CloneWithLocalGit(remotePath, cloneDir); err != nil {
			t.Errorf("CloneWithLocalGit() failed with fallback error: %v", err)
		}

		if _, err := os.Stat(filepath.Join(cloneDir, "README.md")); os.IsNotExist(err) {
			t.Error("file from cloned repo not found with go-git fallback")
		}
	})
}
