package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecBackend implements Backend using the git CLI
type ExecBackend struct {
	repoPath string
	gitPath  string
}

// NewExecBackend creates a new exec-based backend
func NewExecBackend(repoPath string) (*ExecBackend, error) {
	// Find git binary
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("git binary not found: %w", err)
	}

	// Verify git binary works
	cmd := exec.Command(gitPath, "--version")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git binary not functional: %w", err)
	}

	return &ExecBackend{
		repoPath: repoPath,
		gitPath:  gitPath,
	}, nil
}

// runGit executes a git command and returns the output
func (b *ExecBackend) runGit(args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(b.gitPath, args...)
	cmd.Dir = b.repoPath
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s failed: %w\nstderr: %s", strings.Join(args, " "), err, stderr.String())
	}

	return stdout.String(), nil
}

// Log executes git log with the given arguments
func (b *ExecBackend) Log(args ...string) (string, error) {
	fullArgs := append([]string{"log"}, args...)
	return b.runGit(fullArgs...)
}

// LogPretty executes git log with a custom pretty format
func (b *ExecBackend) LogPretty(format string, args ...string) (string, error) {
	fullArgs := append([]string{"log", "--pretty=format:" + format}, args...)
	return b.runGit(fullArgs...)
}

// Branches lists branches
func (b *ExecBackend) Branches(args ...string) (string, error) {
	fullArgs := append([]string{"branch"}, args...)
	return b.runGit(fullArgs...)
}

// Tags lists tags
func (b *ExecBackend) Tags(args ...string) (string, error) {
	fullArgs := append([]string{"tag"}, args...)
	return b.runGit(fullArgs...)
}

// Diff shows differences
func (b *ExecBackend) Diff(args ...string) (string, error) {
	fullArgs := append([]string{"diff"}, args...)
	return b.runGit(fullArgs...)
}

// Show shows various types of objects
func (b *ExecBackend) Show(args ...string) (string, error) {
	fullArgs := append([]string{"show"}, args...)
	return b.runGit(fullArgs...)
}

// RevList lists commit objects in reverse chronological order
func (b *ExecBackend) RevList(args ...string) (string, error) {
	fullArgs := append([]string{"rev-list"}, args...)
	return b.runGit(fullArgs...)
}

// ForEachRef iterates over references
func (b *ExecBackend) ForEachRef(args ...string) (string, error) {
	fullArgs := append([]string{"for-each-ref"}, args...)
	return b.runGit(fullArgs...)
}

// Shortlog summarizes git log output
func (b *ExecBackend) Shortlog(args ...string) (string, error) {
	fullArgs := append([]string{"shortlog"}, args...)
	return b.runGit(fullArgs...)
}

// CurrentBranch returns the current branch name
func (b *ExecBackend) CurrentBranch() (string, error) {
	output, err := b.runGit("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// RootPath returns the repository root path
func (b *ExecBackend) RootPath() string {
	return b.repoPath
}

// BuildLogArgs builds git log arguments from LogOptions
func BuildLogArgs(opts *LogOptions) []string {
	args := []string{}

	// Date range
	if !opts.Since.IsZero() {
		args = append(args, fmt.Sprintf("--since=%s", opts.Since.Format("2006-01-02")))
	}
	if !opts.Until.IsZero() {
		args = append(args, fmt.Sprintf("--until=%s", opts.Until.Format("2006-01-02")))
	}

	// Author filter
	if opts.Author != "" {
		args = append(args, fmt.Sprintf("--author=%s", opts.Author))
	}

	// Format
	if opts.Format != "" {
		args = append(args, fmt.Sprintf("--pretty=format:%s", opts.Format))
	}

	// Merge handling
	if opts.NoMerges {
		args = append(args, "--no-merges")
	} else if opts.MergesOnly {
		args = append(args, "--merges")
	}

	// Limit
	if opts.Limit > 0 {
		args = append(args, fmt.Sprintf("--max-count=%d", opts.Limit))
	}

	// Branch
	if opts.Branch != "" {
		args = append(args, opts.Branch)
	}

	// Extra arguments
	args = append(args, opts.ExtraArgs...)

	// Path specifications
	if len(opts.PathSpec) > 0 {
		args = append(args, "--")
		args = append(args, opts.PathSpec...)
	}

	return args
}
