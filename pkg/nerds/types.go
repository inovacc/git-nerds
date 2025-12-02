package nerds

import "time"

// Stats represents comprehensive repository statistics
type Stats struct {
	TotalCommits  int
	TotalAuthors  int
	TotalFiles    int
	LinesAdded    int
	LinesDeleted  int
	LinesChanged  int
	FirstCommitAt time.Time
	LastCommitAt  time.Time
	ActiveDays    int
	Authors       []Author
	Files         []File
	Branches      []Branch
}

// Author represents a contributor's statistics
type Author struct {
	Name         string
	Email        string
	Commits      int
	LinesAdded   int
	LinesDeleted int
	LinesChanged int
	FilesChanged int
	FirstCommit  time.Time
	LastCommit   time.Time
	ActiveDays   int
}

// Commit represents a single commit
type Commit struct {
	Hash      string
	Author    string
	Email     string
	Date      time.Time
	Message   string
	Files     []string
	Additions int
	Deletions int
}

// Branch represents a git branch
type Branch struct {
	Name      string
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Age       time.Duration
	IsActive  bool
}

// File represents file statistics
type File struct {
	Path         string
	Changes      int
	Additions    int
	Deletions    int
	Authors      []string
	LastModified time.Time
}

// Changelog represents changelog entries
type Changelog struct {
	Version string
	Date    time.Time
	Author  string
	Changes []ChangeEntry
}

// ChangeEntry represents a single changelog entry
type ChangeEntry struct {
	Type    string // feat, fix, docs, refactor, etc.
	Scope   string
	Message string
	Hash    string
}

// Tree represents a branch tree visualization
type Tree struct {
	Root     string
	Branches []TreeNode
}

// TreeNode represents a node in the branch tree
type TreeNode struct {
	Name     string
	Parent   string
	Children []string
	Commits  int
}

// Calendar represents a commit calendar heatmap
type Calendar struct {
	Year   int
	Months []MonthData
}

// MonthData represents commit data for a month
type MonthData struct {
	Month int
	Days  map[int]int // day of month -> commit count
}

// Heatmap represents commit activity heatmap
type Heatmap struct {
	Days  []HeatmapDay
	Hours []HeatmapHour
}

// HeatmapDay represents commits for a single day
type HeatmapDay struct {
	Date    time.Time
	Commits int
}

// HeatmapHour represents commits for a specific hour
type HeatmapHour struct {
	Hour    int
	Commits int
}

// Contributor represents a simplified contributor view
type Contributor struct {
	Name    string
	Email   string
	Commits int
	Since   time.Time
}
