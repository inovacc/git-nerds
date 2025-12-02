package analysis

import (
	"testing"
	"time"

	"github.com/inovacc/git-nerds/internal/git"
)

func setupTestBackend(t *testing.T) git.Backend {
	backend, err := git.NewExecBackend("../../..")
	if err != nil {
		t.Skip("Cannot create git backend:", err)
	}
	return backend
}

func TestNewAuthorAnalyzer(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{}

	analyzer := NewAuthorAnalyzer(backend, opts)
	if analyzer == nil {
		t.Fatal("NewAuthorAnalyzer() returned nil")
	}
}

func TestCommitsPerAuthor(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewAuthorAnalyzer(backend, opts)
	commits, err := analyzer.CommitsPerAuthor()
	if err != nil {
		t.Fatalf("CommitsPerAuthor() error = %v", err)
	}

	if commits == nil {
		t.Fatal("CommitsPerAuthor() returned nil")
	}

	// Validate results
	for author, count := range commits {
		if author == "" {
			t.Error("Found author with empty name")
		}
		if count <= 0 {
			t.Errorf("Author %s has non-positive commit count: %d", author, count)
		}
	}
}

func TestDetailedAuthorStats(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 50,
	}

	analyzer := NewAuthorAnalyzer(backend, opts)
	authors, err := analyzer.DetailedAuthorStats()
	if err != nil {
		t.Fatalf("DetailedAuthorStats() error = %v", err)
	}

	if authors == nil {
		t.Fatal("DetailedAuthorStats() returned nil")
	}

	// Validate results
	for i, author := range authors {
		if author.Name == "" {
			t.Errorf("Author %d has empty name", i)
		}
		if author.Commits < 0 {
			t.Errorf("Author %s has negative commit count", author.Name)
		}
		if author.LinesAdded < 0 {
			t.Errorf("Author %s has negative lines added", author.Name)
		}
		if author.LinesDeleted < 0 {
			t.Errorf("Author %s has negative lines deleted", author.Name)
		}
	}

	// Check sorting (should be sorted by commits descending)
	for i := 1; i < len(authors); i++ {
		if authors[i].Commits > authors[i-1].Commits {
			t.Error("Authors are not sorted by commits descending")
			break
		}
	}
}

func TestNewContributors(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{}

	analyzer := NewAuthorAnalyzer(backend, opts)

	// Get new contributors from last 365 days
	since := time.Now().AddDate(-1, 0, 0)
	newContribs, err := analyzer.NewContributors(since)
	if err != nil {
		t.Fatalf("NewContributors() error = %v", err)
	}

	// Validate that all contributors joined after the since date
	for _, c := range newContribs {
		if c.FirstCommit.Before(since) {
			t.Errorf("Contributor %s first commit %v is before since %v",
				c.Name, c.FirstCommit, since)
		}
	}
}

func TestSuggestReviewers(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{}

	analyzer := NewAuthorAnalyzer(backend, opts)

	// Try with README.md (likely to exist)
	reviewers, err := analyzer.SuggestReviewers("README.md", 5)
	if err != nil {
		t.Logf("SuggestReviewers() error = %v (file may not exist)", err)
		return
	}

	if len(reviewers) > 5 {
		t.Errorf("Expected max 5 reviewers, got %d", len(reviewers))
	}

	for i, r := range reviewers {
		if r == "" {
			t.Errorf("Reviewer %d is empty", i)
		}
	}
}

func TestTopContributors(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewAuthorAnalyzer(backend, opts)

	top, err := analyzer.TopContributors(5)
	if err != nil {
		t.Fatalf("TopContributors() error = %v", err)
	}

	if len(top) > 5 {
		t.Errorf("Expected max 5 contributors, got %d", len(top))
	}

	// Should be sorted by commits descending
	for i := 1; i < len(top); i++ {
		if top[i].Commits > top[i-1].Commits {
			t.Error("Top contributors are not sorted by commits descending")
			break
		}
	}
}

// Benchmark tests
func BenchmarkCommitsPerAuthor(b *testing.B) {
	backend, err := git.NewExecBackend("../../..")
	if err != nil {
		b.Skip("Cannot create git backend")
	}

	opts := &git.LogOptions{Limit: 100}
	analyzer := NewAuthorAnalyzer(backend, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.CommitsPerAuthor()
	}
}

func BenchmarkDetailedAuthorStats(b *testing.B) {
	backend, err := git.NewExecBackend("../../..")
	if err != nil {
		b.Skip("Cannot create git backend")
	}

	opts := &git.LogOptions{Limit: 50}
	analyzer := NewAuthorAnalyzer(backend, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.DetailedAuthorStats()
	}
}
