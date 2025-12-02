package git_nerds

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "current directory",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "parent directory",
			path:    "../..",
			wantErr: false,
		},
		{
			name:    "non-existent path",
			path:    "/nonexistent/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && repo == nil {
				t.Error("Open() returned nil repository")
			}
		})
	}
}

func TestOpenWithOptions(t *testing.T) {
	opts := &Options{
		Since:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:  time.Now(),
		Branch: "main",
		Limit:  10,
	}

	repo, err := Open("../..", opts)
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	if repo.options.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", repo.options.Limit)
	}
}

func TestRepositoryPath(t *testing.T) {
	repo, err := Open(".")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	path := repo.Path()
	if path == "" {
		t.Error("Path() returned empty string")
	}

	if !filepath.IsAbs(path) {
		t.Error("Path() should return absolute path")
	}
}

func TestRepositoryOptions(t *testing.T) {
	opts := DefaultOptions()
	opts.Limit = 5

	repo, err := Open("../..", opts)
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	returnedOpts := repo.Options()
	if returnedOpts.Limit != 5 {
		t.Errorf("Expected limit 5, got %d", returnedOpts.Limit)
	}
}

func TestDetailedStats(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	stats, err := repo.DetailedStats()
	if err != nil {
		t.Fatalf("DetailedStats() error = %v", err)
	}

	if stats == nil {
		t.Fatal("DetailedStats() returned nil")
	}

	if stats.TotalCommits < 0 {
		t.Error("TotalCommits should be non-negative")
	}

	if stats.TotalAuthors < 0 {
		t.Error("TotalAuthors should be non-negative")
	}

	if len(stats.Authors) != stats.TotalAuthors {
		t.Errorf("Authors slice length %d != TotalAuthors %d", len(stats.Authors), stats.TotalAuthors)
	}
}

func TestContributors(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	contributors, err := repo.Contributors()
	if err != nil {
		t.Fatalf("Contributors() error = %v", err)
	}

	if contributors == nil {
		t.Fatal("Contributors() returned nil")
	}

	for i, c := range contributors {
		if c.Name == "" {
			t.Errorf("Contributor %d has empty name", i)
		}
		if c.Commits < 0 {
			t.Errorf("Contributor %s has negative commit count", c.Name)
		}
	}
}

func TestCommitsPerAuthor(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	commits, err := repo.CommitsPerAuthor()
	if err != nil {
		t.Fatalf("CommitsPerAuthor() error = %v", err)
	}

	if commits == nil {
		t.Fatal("CommitsPerAuthor() returned nil")
	}

	totalCommits := 0
	for author, count := range commits {
		if author == "" {
			t.Error("Found author with empty name")
		}
		if count < 0 {
			t.Errorf("Author %s has negative commit count", author)
		}
		totalCommits += count
	}

	if totalCommits == 0 && len(commits) > 0 {
		t.Error("Commits map is not empty but total is 0")
	}
}

func TestCommitsByDay(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	byDay, err := repo.CommitsByDay()
	if err != nil {
		t.Fatalf("CommitsByDay() error = %v", err)
	}

	if byDay == nil {
		t.Fatal("CommitsByDay() returned nil")
	}

	for date, count := range byDay {
		// Validate date format (YYYY-MM-DD)
		if _, err := time.Parse("2006-01-02", date); err != nil {
			t.Errorf("Invalid date format: %s", date)
		}
		if count < 0 {
			t.Errorf("Date %s has negative count", date)
		}
	}
}

func TestCommitsByWeekday(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	byWeekday, err := repo.CommitsByWeekday()
	if err != nil {
		t.Fatalf("CommitsByWeekday() error = %v", err)
	}

	if byWeekday == nil {
		t.Fatal("CommitsByWeekday() returned nil")
	}

	validWeekdays := map[string]bool{
		"Monday": true, "Tuesday": true, "Wednesday": true,
		"Thursday": true, "Friday": true, "Saturday": true, "Sunday": true,
	}

	for day, count := range byWeekday {
		if !validWeekdays[day] {
			t.Errorf("Invalid weekday: %s", day)
		}
		if count < 0 {
			t.Errorf("Weekday %s has negative count", day)
		}
	}
}

func TestCommitsByHour(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	byHour, err := repo.CommitsByHour()
	if err != nil {
		t.Fatalf("CommitsByHour() error = %v", err)
	}

	if byHour == nil {
		t.Fatal("CommitsByHour() returned nil")
	}

	for hour, count := range byHour {
		if hour < 0 || hour > 23 {
			t.Errorf("Invalid hour: %d", hour)
		}
		if count < 0 {
			t.Errorf("Hour %d has negative count", hour)
		}
	}
}

func TestBranchesByDate(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	branches, err := repo.BranchesByDate()
	if err != nil {
		t.Fatalf("BranchesByDate() error = %v", err)
	}

	if branches == nil {
		t.Fatal("BranchesByDate() returned nil")
	}

	// Check that branches are sorted (newest first)
	for i := 1; i < len(branches); i++ {
		if branches[i].UpdatedAt.After(branches[i-1].UpdatedAt) {
			t.Error("Branches are not sorted by date (newest first)")
			break
		}
	}
}

func TestSuggestReviewers(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	// Try to suggest reviewers for README.md (likely to exist)
	reviewers, err := repo.SuggestReviewers("README.md")
	if err != nil {
		t.Logf("SuggestReviewers() error = %v (file may not exist)", err)
		return
	}

	if reviewers != nil {
		for _, r := range reviewers {
			if r == "" {
				t.Error("Found empty reviewer")
			}
		}
	}
}

func TestNewContributors(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	// Get contributors from last 365 days
	since := time.Now().AddDate(-1, 0, 0)
	newContribs, err := repo.NewContributors(since)
	if err != nil {
		t.Fatalf("NewContributors() error = %v", err)
	}

	// Validate that all contributors joined after the since date
	for _, c := range newContribs {
		if c.Since.Before(since) {
			t.Errorf("Contributor %s joined before since date", c.Name)
		}
	}
}

func TestExportJSON(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	json, err := repo.ExportJSON()
	if err != nil {
		t.Fatalf("ExportJSON() error = %v", err)
	}

	if json == "" {
		t.Error("ExportJSON() returned empty string")
	}

	// Basic JSON validation
	if json[0] != '{' {
		t.Error("JSON should start with {")
	}
}

func TestExportMarkdown(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	md, err := repo.ExportMarkdown()
	if err != nil {
		t.Fatalf("ExportMarkdown() error = %v", err)
	}

	if md == "" {
		t.Error("ExportMarkdown() returned empty string")
	}

	// Should contain markdown headers
	if len(md) > 0 && md[0] != '#' {
		t.Error("Markdown should start with # header")
	}
}

func TestExportCSV(t *testing.T) {
	repo, err := Open("../..")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	csv, err := repo.ExportCSV()
	if err != nil {
		t.Fatalf("ExportCSV() error = %v", err)
	}

	if csv == "" {
		t.Error("ExportCSV() returned empty string")
	}
}

// Benchmark tests
func BenchmarkOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Open("../..")
	}
}

func BenchmarkDetailedStats(b *testing.B) {
	repo, err := Open("../..")
	if err != nil {
		b.Skip("Cannot open repository")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.DetailedStats()
	}
}

func BenchmarkCommitsPerAuthor(b *testing.B) {
	repo, err := Open("../..")
	if err != nil {
		b.Skip("Cannot open repository")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.CommitsPerAuthor()
	}
}

// Test helper to check if we're in a git repository
func isGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}
