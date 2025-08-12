package stats

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// FindGit looks for the 'git' binary on the system's PATH.
// It returns the path to the binary if found and working, or an error otherwise.
func FindGit() (string, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		// Return the actual error from LookPath
		return "", fmt.Errorf("local git binary not found: %w", err)
	}

	// Verify the binary is a working git client
	if err := exec.Command(gitPath, "--version").Run(); err != nil {
		// Return the actual error from the command execution
		return "", fmt.Errorf("found binary is not a valid git client: %w", err)
	}

	return gitPath, nil
}

// OpenRepo tries to open a git repo using the local git binary, falling back to go-git if not found.
// Note: This function now only uses the go-git library to return a *git.Repository object.
// The presence of a local git binary is checked to signal its availability, but not used directly
// in this function to open the repo, as go-git is the primary library for repo object returns.
func OpenRepo(path string) (*git.Repository, error) {
	// First, check if a local git binary exists.
	// We'll return this path to the caller if it's found, indicating that
	// commands could be run with the local binary if needed.
	// The function signature has been simplified to return only the repo and error.
	if _, err := FindGit(); err != nil {
		// Local git not found, fall back to go-git
		repo, err := git.PlainOpen(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open repo with go-git: %w", err)
		}

		return repo, nil
	}

	// If a local binary exists, we can still use go-git to open the repository.
	// go-git is a powerful library and may be preferred for a pure Go implementation.
	// If you need to run specific commands with the local binary, you would need a different function.
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo with go-git even though local binary exists: %w", err)
	}

	return repo, nil
}

// CloneWithLocalGit A new function to demonstrate how you might use a local git binary.
// This is not part of the original code but shows a more robust approach.
func CloneWithLocalGit(url string, dir string) error {
	gitPath, err := FindGit()
	if err != nil {
		// Local git binary not found, try cloning with go-git
		// The clone options are kept minimal to avoid go-git trying to use an
		// external process, which was causing the test to fail.
		cloneOpts := &git.CloneOptions{
			URL:      url,
			Progress: nil,
		}

		// Only add HTTP auth if it's an HTTP/S URL.
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			cloneOpts.Auth = &http.BasicAuth{}
		}

		if _, err := git.PlainClone(dir, false, cloneOpts); err != nil {
			return fmt.Errorf("failed to clone with go-git: %w", err)
		}
	} else {
		// Local git binary found, use it to clone
		if err := exec.Command(gitPath, "clone", url, dir).Run(); err != nil {
			return fmt.Errorf("failed to clone with local git: %w", err)
		}
	}

	return nil
}
