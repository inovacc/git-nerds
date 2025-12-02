package analysis

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/inovacc/git-nerds/internal/git"
	"github.com/inovacc/git-nerds/internal/parse"
)

// AuthorAnalyzer provides author-related analytics
type AuthorAnalyzer struct {
	backend git.Backend
	options *git.LogOptions
}

// NewAuthorAnalyzer creates a new author analyzer
func NewAuthorAnalyzer(backend git.Backend, options *git.LogOptions) *AuthorAnalyzer {
	return &AuthorAnalyzer{
		backend: backend,
		options: options,
	}
}

// CommitsPerAuthor returns commit counts grouped by author
func (a *AuthorAnalyzer) CommitsPerAuthor() (map[string]int, error) {
	// Use git shortlog for efficient counting
	args := git.BuildLogArgs(a.options)
	args = append([]string{"-s", "-n", "-e"}, args...) // summary, numbered, email

	output, err := a.backend.Shortlog(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get shortlog: %w", err)
	}

	stats, err := parse.ParseAuthorStats(output)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, stat := range stats {
		result[stat.Name] = stat.Count
	}

	return result, nil
}

// AuthorDetails represents detailed author information
type AuthorDetails struct {
	Name         string
	Email        string
	Commits      int
	LinesAdded   int
	LinesDeleted int
	FilesChanged int
	FirstCommit  time.Time
	LastCommit   time.Time
	ActiveDays   int
}

// DetailedAuthorStats returns comprehensive statistics for all authors
func (a *AuthorAnalyzer) DetailedAuthorStats() ([]AuthorDetails, error) {
	// Get commits with detailed information
	format := "%H|%an|%ae|%ad|%s"
	args := git.BuildLogArgs(a.options)
	args = append([]string{"--date=iso", "--numstat"}, args...)

	output, err := a.backend.Log(append([]string{"--pretty=format:" + format}, args...)...)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	// Parse commits
	commits, err := parse.ParseCommitLog(output)
	if err != nil {
		return nil, err
	}

	// Aggregate by author
	authorMap := make(map[string]*AuthorDetails)

	for _, commit := range commits {
		key := commit.Email
		if _, exists := authorMap[key]; !exists {
			authorMap[key] = &AuthorDetails{
				Name:        commit.Author,
				Email:       commit.Email,
				FirstCommit: commit.Date,
				LastCommit:  commit.Date,
			}
		}

		author := authorMap[key]
		author.Commits++
		author.LinesAdded += commit.Additions
		author.LinesDeleted += commit.Deletions
		author.FilesChanged += len(commit.Files)

		if commit.Date.Before(author.FirstCommit) {
			author.FirstCommit = commit.Date
		}
		if commit.Date.After(author.LastCommit) {
			author.LastCommit = commit.Date
		}
	}

	// Calculate active days for each author
	for email, author := range authorMap {
		activeDays, err := a.calculateActiveDays(email)
		if err == nil {
			author.ActiveDays = activeDays
		}
	}

	// Convert map to slice
	result := make([]AuthorDetails, 0, len(authorMap))
	for _, author := range authorMap {
		result = append(result, *author)
	}

	// Sort by commits descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Commits > result[j].Commits
	})

	return result, nil
}

// calculateActiveDays counts unique days an author has committed
func (a *AuthorAnalyzer) calculateActiveDays(email string) (int, error) {
	opts := *a.options
	opts.Author = email
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=short"}, args...)

	output, err := a.backend.Log(args...)
	if err != nil {
		return 0, err
	}

	// Count unique dates
	dates := make(map[string]bool)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			dates[line] = true
		}
	}

	return len(dates), nil
}

// NewContributors returns contributors who joined after a given date
func (a *AuthorAnalyzer) NewContributors(since time.Time) ([]AuthorDetails, error) {
	all, err := a.DetailedAuthorStats()
	if err != nil {
		return nil, err
	}

	result := make([]AuthorDetails, 0)
	for _, author := range all {
		if author.FirstCommit.After(since) {
			result = append(result, author)
		}
	}

	// Sort by first commit date
	sort.Slice(result, func(i, j int) bool {
		return result[i].FirstCommit.Before(result[j].FirstCommit)
	})

	return result, nil
}

// SuggestReviewers suggests reviewers for a file based on commit history
func (a *AuthorAnalyzer) SuggestReviewers(file string, limit int) ([]string, error) {
	// Get commits that modified this file
	opts := *a.options
	opts.PathSpec = []string{file}
	args := git.BuildLogArgs(&opts)

	output, err := a.backend.Log(append([]string{"--pretty=format:%ae"}, args...)...)
	if err != nil {
		return nil, fmt.Errorf("failed to get log for file: %w", err)
	}

	authors := parse.ParseAuthors(output)

	// Count occurrences
	counts := make(map[string]int)
	for _, author := range authors {
		counts[author]++
	}

	// Sort by count
	type authorCount struct {
		email string
		count int
	}
	sorted := make([]authorCount, 0, len(counts))
	for email, count := range counts {
		sorted = append(sorted, authorCount{email, count})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	// Return top N
	if limit > len(sorted) {
		limit = len(sorted)
	}

	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = sorted[i].email
	}

	return result, nil
}

// TopContributors returns the top N contributors by commit count
func (a *AuthorAnalyzer) TopContributors(limit int) ([]AuthorDetails, error) {
	all, err := a.DetailedAuthorStats()
	if err != nil {
		return nil, err
	}

	if limit > len(all) {
		limit = len(all)
	}

	return all[:limit], nil
}
