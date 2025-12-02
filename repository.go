package git_nerds

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	analysis2 "github.com/inovacc/git-nerds/internal/analysis"
	git2 "github.com/inovacc/git-nerds/internal/git"
)

// Repository provides access to Git repository statistics and analysis
type Repository struct {
	path    string
	options *Options
	backend git2.Backend
}

// Open opens a Git repository at the specified path
func Open(path string, opts ...*Options) (*Repository, error) {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPath, err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPath, absPath)
	}

	// Check if it's a git repository
	gitDir := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		// Maybe it's the .git directory itself
		if filepath.Base(absPath) != ".git" {
			return nil, fmt.Errorf("%w: %s", ErrNotARepository, absPath)
		}
		absPath = filepath.Dir(absPath)
	}

	// Use provided options or defaults
	var options *Options
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	} else {
		options = DefaultOptions()
	}

	// Validate options
	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOptions, err)
	}

	// Create backend
	backend, err := git2.NewExecBackend(absPath)
	if err != nil {
		return nil, err
	}

	return &Repository{
		path:    absPath,
		options: options,
		backend: backend,
	}, nil
}

// Path returns the repository path
func (r *Repository) Path() string {
	return r.path
}

// Options returns the repository options
func (r *Repository) Options() *Options {
	return r.options
}

// toLogOptions converts Options to git.LogOptions
func (r *Repository) toLogOptions() *git2.LogOptions {
	return &git2.LogOptions{
		Since:         r.options.Since,
		Until:         r.options.Until,
		Author:        "",
		Format:        "",
		Branch:        r.options.Branch,
		PathSpec:      r.options.PathSpec,
		NoMerges:      !r.options.IncludeMerges && !r.options.OnlyMerges,
		MergesOnly:    r.options.OnlyMerges,
		Limit:         r.options.Limit,
		IgnoreAuthors: r.options.IgnoreAuthors,
		ExtraArgs:     r.options.LogOptions,
	}
}

// DetailedStats returns comprehensive repository statistics
func (r *Repository) DetailedStats() (*Stats, error) {
	logOpts := r.toLogOptions()

	// Create analyzers
	authorAnalyzer := analysis2.NewAuthorAnalyzer(r.backend, logOpts)
	branchAnalyzer := analysis2.NewBranchAnalyzer(r.backend, logOpts)

	// Get author details
	authors, err := authorAnalyzer.DetailedAuthorStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get author stats: %w", err)
	}

	// Get branch details
	branches, err := branchAnalyzer.DetailedBranchInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get branch info: %w", err)
	}

	// Calculate totals
	stats := &Stats{
		TotalAuthors: len(authors),
		Authors:      make([]Author, len(authors)),
		Branches:     make([]Branch, len(branches)),
	}

	for i, a := range authors {
		stats.Authors[i] = Author{
			Name:         a.Name,
			Email:        a.Email,
			Commits:      a.Commits,
			LinesAdded:   a.LinesAdded,
			LinesDeleted: a.LinesDeleted,
			LinesChanged: a.LinesAdded + a.LinesDeleted,
			FilesChanged: a.FilesChanged,
			FirstCommit:  a.FirstCommit,
			LastCommit:   a.LastCommit,
			ActiveDays:   a.ActiveDays,
		}
		stats.TotalCommits += a.Commits
		stats.LinesAdded += a.LinesAdded
		stats.LinesDeleted += a.LinesDeleted
	}

	stats.LinesChanged = stats.LinesAdded + stats.LinesDeleted

	if len(authors) > 0 {
		stats.FirstCommitAt = authors[len(authors)-1].FirstCommit
		stats.LastCommitAt = authors[0].LastCommit
	}

	for i, b := range branches {
		stats.Branches[i] = Branch{
			Name:      b.Name,
			Hash:      b.Hash,
			UpdatedAt: b.LastCommit,
			Age:       b.Age,
			IsActive:  b.IsActive,
		}
	}

	return stats, nil
}

// StatsByBranch returns statistics for a specific branch
func (r *Repository) StatsByBranch(branch string) (*Stats, error) {
	// TODO: Implement branch-specific stats
	return &Stats{}, fmt.Errorf("not implemented yet")
}

// Contributors returns all contributors
func (r *Repository) Contributors() ([]Contributor, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewAuthorAnalyzer(r.backend, logOpts)

	authors, err := analyzer.DetailedAuthorStats()
	if err != nil {
		return nil, err
	}

	result := make([]Contributor, len(authors))
	for i, a := range authors {
		result[i] = Contributor{
			Name:    a.Name,
			Email:   a.Email,
			Commits: a.Commits,
			Since:   a.FirstCommit,
		}
	}

	return result, nil
}

// NewContributors returns new contributors since a given date
func (r *Repository) NewContributors(since time.Time) ([]Contributor, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewAuthorAnalyzer(r.backend, logOpts)

	authors, err := analyzer.NewContributors(since)
	if err != nil {
		return nil, err
	}

	result := make([]Contributor, len(authors))
	for i, a := range authors {
		result[i] = Contributor{
			Name:    a.Name,
			Email:   a.Email,
			Commits: a.Commits,
			Since:   a.FirstCommit,
		}
	}

	return result, nil
}

// CommitsPerAuthor returns commit counts by author
func (r *Repository) CommitsPerAuthor() (map[string]int, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewAuthorAnalyzer(r.backend, logOpts)
	return analyzer.CommitsPerAuthor()
}

// SuggestReviewers suggests reviewers for a file based on history
func (r *Repository) SuggestReviewers(file string) ([]string, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewAuthorAnalyzer(r.backend, logOpts)
	return analyzer.SuggestReviewers(file, 5) // Top 5 reviewers
}

// CommitsByDay returns commits grouped by day
func (r *Repository) CommitsByDay() (map[string]int, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)
	return analyzer.CommitsByDay()
}

// CommitsByMonth returns commits grouped by month
func (r *Repository) CommitsByMonth() (map[string]int, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)
	return analyzer.CommitsByMonth()
}

// CommitsByYear returns commits grouped by year
func (r *Repository) CommitsByYear() (map[string]int, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)
	return analyzer.CommitsByYear()
}

// CommitsByWeekday returns commits grouped by weekday
func (r *Repository) CommitsByWeekday() (map[string]int, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)
	return analyzer.CommitsByWeekday()
}

// CommitsByHour returns commits grouped by hour
func (r *Repository) CommitsByHour() (map[int]int, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)
	return analyzer.CommitsByHour()
}

// CommitsByTimezone returns commits grouped by timezone
func (r *Repository) CommitsByTimezone() (map[string]int, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)
	return analyzer.CommitsByTimezone()
}

// BranchTree returns the branch tree structure
func (r *Repository) BranchTree() (*Tree, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewBranchAnalyzer(r.backend, logOpts)

	treeOutput, err := analyzer.BranchTree()
	if err != nil {
		return nil, err
	}

	// Return a simple tree with the output
	return &Tree{
		Root:     r.path,
		Branches: []TreeNode{{Name: treeOutput}},
	}, nil
}

// BranchesByDate returns branches sorted by date
func (r *Repository) BranchesByDate() ([]Branch, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewBranchAnalyzer(r.backend, logOpts)

	branches, err := analyzer.BranchesByDate()
	if err != nil {
		return nil, err
	}

	result := make([]Branch, len(branches))
	for i, b := range branches {
		result[i] = Branch{
			Name:      b.Name,
			Hash:      b.Hash,
			UpdatedAt: b.LastCommit,
			Age:       b.Age,
			IsActive:  b.IsActive,
		}
	}

	return result, nil
}

// CommitsCalendar returns a calendar heatmap of commits
func (r *Repository) CommitsCalendar(author string) (*Calendar, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)

	year := time.Now().Year()
	calData, err := analyzer.GenerateCalendar(year, author)
	if err != nil {
		return nil, err
	}

	// Convert to public type
	result := &Calendar{
		Year:   calData.Year,
		Months: make([]MonthData, len(calData.Months)),
	}

	for i, m := range calData.Months {
		result.Months[i] = MonthData{
			Month: int(m.Month),
			Days:  make(map[int]int),
		}
	}

	return result, nil
}

// CommitsHeatmap returns a heatmap of commits for the last N days
func (r *Repository) CommitsHeatmap(days int) (*Heatmap, error) {
	logOpts := r.toLogOptions()
	analyzer := analysis2.NewTemporalAnalyzer(r.backend, logOpts)

	heatData, err := analyzer.GenerateHeatmap(days)
	if err != nil {
		return nil, err
	}

	// Convert to public type
	result := &Heatmap{
		Days:  make([]HeatmapDay, len(heatData.Days)),
		Hours: make([]HeatmapHour, 0),
	}

	for i, d := range heatData.Days {
		result.Days[i] = HeatmapDay{
			Date:    d.Date,
			Commits: d.CommitCount,
		}
	}

	return result, nil
}

// Changelogs generates changelogs
func (r *Repository) Changelogs() ([]Changelog, error) {
	// TODO: Implement changelog generation
	return nil, fmt.Errorf("not implemented yet")
}

// ChangelogsByAuthor generates changelogs for a specific author
func (r *Repository) ChangelogsByAuthor(author string) ([]Changelog, error) {
	// TODO: Implement author-specific changelogs
	return nil, fmt.Errorf("not implemented yet")
}
