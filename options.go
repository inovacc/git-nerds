package git_nerds

import "time"

// Options configures repository analysis behavior
type Options struct {
	// Time range filters
	Since time.Time
	Until time.Time

	// Branch filtering
	Branch string

	// Path specifications (include/exclude patterns)
	// Use gitignore-style patterns, e.g., ":!vendor", ":!*.generated.go"
	PathSpec []string

	// Author filtering (regex patterns)
	// Example: []string{"bot@.*", ".*\\[bot\\]"}
	IgnoreAuthors []string

	// Merge commit handling
	IncludeMerges bool // if false, excludes merge commits
	OnlyMerges    bool // if true, shows only merge commits

	// Result limiting
	Limit int // limit number of results (0 = no limit)

	// Sorting options
	SortBy    string // "name", "commits", "lines", etc.
	SortOrder string // "asc" or "desc"

	// Additional git log options
	LogOptions []string
}

// DefaultOptions returns sensible default options
func DefaultOptions() *Options {
	return &Options{
		Since:         time.Time{}, // beginning of repo
		Until:         time.Now(),
		Branch:        "",         // current branch
		PathSpec:      []string{}, // no exclusions
		IgnoreAuthors: []string{}, // no ignores
		IncludeMerges: false,
		OnlyMerges:    false,
		Limit:         0,
		SortBy:        "commits",
		SortOrder:     "desc",
		LogOptions:    []string{},
	}
}

// Validate checks if options are valid
func (o *Options) Validate() error {
	// Add validation logic here if needed
	return nil
}
