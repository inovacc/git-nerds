package stats

import (
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repo struct {
	Path          string
	Stats         *RepoStats
	Authors       []string
	Commits       []string
	Branches      []string
	Tags          []string
	HEAD          string
	Files         []string
	FirstCommit   *object.Commit
	LastCommit    *object.Commit
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Contributors  map[string]int
	RemoteURL     string
	DefaultBranch string
	IsBare        bool
}

type FileStat struct {
	Name  string
	Count int
}

type RepoStats struct {
	CommitsByUser     map[string]int
	FileModifications map[string]int
	LinesAdded        int
	LinesDeleted      int
	CommitsByWeekday  map[time.Weekday]int
}

// NewRepo initializes a Repo struct for the given path and loads repo data.
func NewRepo(path string) (*Repo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// Optionally: check if .git exists in path
	return &Repo{Path: absPath}, nil
}

// RunGit executes a git command in the repo and returns output.
func (r *Repo) RunGit(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path

	return cmd.Output()
}

// CalculateStats computes repository statistics similar to case1/main.go
func CalculateStats(repoPath string) (*RepoStats, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	commitIterator, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	stats := &RepoStats{
		CommitsByUser:     make(map[string]int),
		FileModifications: make(map[string]int),
		CommitsByWeekday:  make(map[time.Weekday]int),
	}

	commitIterator.ForEach(func(commit *object.Commit) error {
		if commit.NumParents() == 0 {
			return nil
		}

		stats.CommitsByUser[commit.Author.Email]++
		stats.CommitsByWeekday[commit.Author.When.Weekday()]++

		parent, err := commit.Parent(0)
		if err != nil {
			return err
		}

		patch, err := parent.Patch(commit)
		if err != nil {
			return err
		}

		for _, fs := range patch.Stats() {
			stats.FileModifications[fs.Name]++
			stats.LinesAdded += fs.Addition
			stats.LinesDeleted += fs.Deletion
		}

		return nil
	})

	return stats, nil
}

// TopModifiedFiles returns the top N most modified files
func (s *RepoStats) TopModifiedFiles(n int) []FileStat {
	fileStats := make([]FileStat, 0, len(s.FileModifications))
	for name, count := range s.FileModifications {
		fileStats = append(fileStats, FileStat{Name: name, Count: count})
	}

	sort.Slice(fileStats, func(i, j int) bool {
		return fileStats[i].Count > fileStats[j].Count
	})

	if n > len(fileStats) {
		n = len(fileStats)
	}

	return fileStats[:n]
}
