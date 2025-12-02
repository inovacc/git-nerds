package nerds

import "errors"

var (
	// ErrInvalidPath is returned when the repository path is invalid
	ErrInvalidPath = errors.New("invalid repository path")

	// ErrNotARepository is returned when the path is not a git repository
	ErrNotARepository = errors.New("not a git repository")

	// ErrInvalidOptions is returned when options are invalid
	ErrInvalidOptions = errors.New("invalid options")

	// ErrGitNotFound is returned when git binary is not found
	ErrGitNotFound = errors.New("git binary not found")

	// ErrGitCommandFailed is returned when a git command fails
	ErrGitCommandFailed = errors.New("git command failed")

	// ErrNoCommits is returned when repository has no commits
	ErrNoCommits = errors.New("no commits found")

	// ErrInvalidBranch is returned when specified branch doesn't exist
	ErrInvalidBranch = errors.New("invalid branch")

	// ErrParseError is returned when parsing git output fails
	ErrParseError = errors.New("failed to parse git output")
)
