package analysis

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/inovacc/git-nerds/internal/git"
	"github.com/inovacc/git-nerds/internal/parse"
)

// BranchAnalyzer provides branch-related analytics
type BranchAnalyzer struct {
	backend git.Backend
	options *git.LogOptions
}

// NewBranchAnalyzer creates a new branch analyzer
func NewBranchAnalyzer(backend git.Backend, options *git.LogOptions) *BranchAnalyzer {
	return &BranchAnalyzer{
		backend: backend,
		options: options,
	}
}

// BranchInfo represents detailed branch information
type BranchInfo struct {
	Name        string
	Hash        string
	LastCommit  time.Time
	CommitCount int
	Author      string
	Age         time.Duration
	IsActive    bool
	IsCurrent   bool
}

// ListBranches returns all branches with basic information
func (b *BranchAnalyzer) ListBranches() ([]string, error) {
	output, err := b.backend.Branches("-a")
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	return parse.ParseBranches(output)
}

// DetailedBranchInfo returns detailed information for all branches
func (b *BranchAnalyzer) DetailedBranchInfo() ([]BranchInfo, error) {
	// Get all branches with their commit info
	// Format: refname, committerdate, objectname, authorname
	output, err := b.backend.ForEachRef(
		"--format=%(refname:short)|%(committerdate:iso)|%(objectname:short)|%(authorname)",
		"refs/heads/",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch info: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	branches := make([]BranchInfo, 0, len(lines))

	currentBranch, _ := b.backend.CurrentBranch()

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}

		name := parts[0]
		dateStr := parts[1]
		hash := parts[2]
		author := parts[3]

		// Parse date
		lastCommit, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err != nil {
			continue
		}

		// Get commit count for this branch
		commitCount, _ := b.getCommitCount(name)

		// Calculate age
		age := time.Since(lastCommit)

		// Consider active if committed in last 30 days
		isActive := age < 30*24*time.Hour

		branches = append(branches, BranchInfo{
			Name:        name,
			Hash:        hash,
			LastCommit:  lastCommit,
			CommitCount: commitCount,
			Author:      author,
			Age:         age,
			IsActive:    isActive,
			IsCurrent:   name == currentBranch,
		})
	}

	return branches, nil
}

// getCommitCount returns the number of commits in a branch
func (b *BranchAnalyzer) getCommitCount(branch string) (int, error) {
	output, err := b.backend.RevList("--count", branch)
	if err != nil {
		return 0, err
	}

	var count int
	fmt.Sscanf(strings.TrimSpace(output), "%d", &count)
	return count, nil
}

// BranchesByDate returns branches sorted by last commit date
func (b *BranchAnalyzer) BranchesByDate() ([]BranchInfo, error) {
	branches, err := b.DetailedBranchInfo()
	if err != nil {
		return nil, err
	}

	// Sort by last commit date, newest first
	sort.Slice(branches, func(i, j int) bool {
		return branches[i].LastCommit.After(branches[j].LastCommit)
	})

	return branches, nil
}

// ActiveBranches returns branches that have been active recently
func (b *BranchAnalyzer) ActiveBranches(days int) ([]BranchInfo, error) {
	branches, err := b.DetailedBranchInfo()
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	active := make([]BranchInfo, 0)

	for _, branch := range branches {
		if branch.LastCommit.After(cutoff) {
			active = append(active, branch)
		}
	}

	// Sort by last commit date
	sort.Slice(active, func(i, j int) bool {
		return active[i].LastCommit.After(active[j].LastCommit)
	})

	return active, nil
}

// StaleBranches returns branches that haven't been active recently
func (b *BranchAnalyzer) StaleBranches(days int) ([]BranchInfo, error) {
	branches, err := b.DetailedBranchInfo()
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	stale := make([]BranchInfo, 0)

	for _, branch := range branches {
		if branch.LastCommit.Before(cutoff) {
			stale = append(stale, branch)
		}
	}

	// Sort by last commit date, oldest first
	sort.Slice(stale, func(i, j int) bool {
		return stale[i].LastCommit.Before(stale[j].LastCommit)
	})

	return stale, nil
}

// BranchTree generates a simple ASCII tree of branches
func (b *BranchAnalyzer) BranchTree() (string, error) {
	// Use git log --graph for visualization
	output, err := b.backend.Log(
		"--graph",
		"--oneline",
		"--all",
		"--decorate",
		"--simplify-by-decoration",
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate branch tree: %w", err)
	}

	return output, nil
}

// MergeStatistics represents merge commit statistics
type MergeStatistics struct {
	TotalMerges     int
	MergesByAuthor  map[string]int
	MergesByMonth   map[string]int
	AverageMergeAge time.Duration
}

// GetMergeStatistics analyzes merge commits
func (b *BranchAnalyzer) GetMergeStatistics() (*MergeStatistics, error) {
	// Get merge commits only
	opts := *b.options
	opts.MergesOnly = true
	opts.Format = "%an|%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=format:%Y-%m"}, args...)

	output, err := b.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge statistics: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")

	stats := &MergeStatistics{
		TotalMerges:    len(lines),
		MergesByAuthor: make(map[string]int),
		MergesByMonth:  make(map[string]int),
	}

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		author := parts[0]
		month := parts[1]

		stats.MergesByAuthor[author]++
		stats.MergesByMonth[month]++
	}

	return stats, nil
}

// CompareWithBranch compares current branch with another branch
func (b *BranchAnalyzer) CompareWithBranch(otherBranch string) (*BranchComparison, error) {
	currentBranch, err := b.backend.CurrentBranch()
	if err != nil {
		return nil, err
	}

	// Commits in current but not in other
	aheadOutput, err := b.backend.RevList("--count", fmt.Sprintf("%s..%s", otherBranch, currentBranch))
	if err != nil {
		return nil, err
	}

	// Commits in other but not in current
	behindOutput, err := b.backend.RevList("--count", fmt.Sprintf("%s..%s", currentBranch, otherBranch))
	if err != nil {
		return nil, err
	}

	var ahead, behind int
	fmt.Sscanf(strings.TrimSpace(aheadOutput), "%d", &ahead)
	fmt.Sscanf(strings.TrimSpace(behindOutput), "%d", &behind)

	return &BranchComparison{
		CurrentBranch: currentBranch,
		OtherBranch:   otherBranch,
		Ahead:         ahead,
		Behind:        behind,
	}, nil
}

// BranchComparison represents comparison between two branches
type BranchComparison struct {
	CurrentBranch string
	OtherBranch   string
	Ahead         int // commits in current but not in other
	Behind        int // commits in other but not in current
}
