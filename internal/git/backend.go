package git

import "time"

// Backend defines the interface for Git operations
// This allows for multiple implementations (exec, go-git, mocks)
type Backend interface {
	// Log executes git log with the given arguments
	Log(args ...string) (string, error)

	// LogPretty executes git log with a custom pretty format
	LogPretty(format string, args ...string) (string, error)

	// Branches lists branches
	Branches(args ...string) (string, error)

	// Tags lists tags
	Tags(args ...string) (string, error)

	// Diff shows differences
	Diff(args ...string) (string, error)

	// Show shows various types of objects
	Show(args ...string) (string, error)

	// RevList lists commit objects in reverse chronological order
	RevList(args ...string) (string, error)

	// ForEachRef iterates over references
	ForEachRef(args ...string) (string, error)

	// Shortlog summarizes git log output
	Shortlog(args ...string) (string, error)

	// CurrentBranch returns the current branch name
	CurrentBranch() (string, error)

	// RootPath returns the repository root path
	RootPath() string
}

// LogOptions provides structured options for git log queries
type LogOptions struct {
	Since         time.Time
	Until         time.Time
	Author        string
	Format        string
	Branch        string
	PathSpec      []string
	NoMerges      bool
	MergesOnly    bool
	Limit         int
	IgnoreAuthors []string
	ExtraArgs     []string
}

// BranchInfo represents branch information
type BranchInfo struct {
	Name      string
	Hash      string
	IsCurrent bool
	IsRemote  bool
}

// TagInfo represents tag information
type TagInfo struct {
	Name string
	Hash string
	Date time.Time
}
