package stats

import (
	"os/exec"
	"path/filepath"
)

type Repo struct {
	Path string
	// Add fields for stats as needed, e.g. Authors, Commits, Branches, etc.
}

// NewRepo initializes a Repo struct for the given path and loads repo data.
func NewRepo(path string) (*Repo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// Optionally: check if .git exists in path
	return &Repo{Path: absPath}, nil
}

// RunGit executes a git command in the repo and returns output.
func (r *Repo) RunGit(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path

	return cmd.Output()
}

// Add methods to populate and return stats as needed.
