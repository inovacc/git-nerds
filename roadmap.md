# Git Nerds - Module Roadmap

> **Vision**: Build a comprehensive, reusable Go module that provides deep intelligence about Git repositories for integration into any application.

---

## 0) Vision

- **Product**: A Go library that extracts, analyzes, and exposes Git repository metrics (history, branches, authors, files, merges, releases)
- **Delivery**: Clean Go API with optional CLI for testing/validation
- **Inspiration**: Provide the same insights as git-quick-stats but as a consumable module
- **Use Cases**: Analytics tools, CI/CD systems, code review tools, IDE plugins, team dashboards

## 1) Design Principles

- **Module-First**: Library designed for consumption, not a standalone app
- **Clean API**: Simple, intuitive Go interfaces following best practices
- **Backend Agnostic**: Support multiple Git implementations (exec, go-git)
- **Zero Dependencies**: Minimal external dependencies for core functionality
- **Performance**: Streaming, caching, and concurrent processing where beneficial
- **Testable**: Comprehensive test coverage with fixtures and golden tests

---

## 2) Roadmap by Phase

### Phase 0 — Project Foundation
- [x] Repository setup with Go modules
- [ ] Conventional commits and commitlint configuration
- [ ] Git hooks (pre-commit: fmt/vet/test, pre-push: race tests)
- [ ] Semantic versioning with git tags (vX.Y.Z)
- [ ] GitHub Actions: build, test, lint, release
- [ ] Issue and PR templates
- [ ] Code of conduct and contributing guide

### Phase 1 — Core Git Backend
- [ ] **Backend Interface**: Define clean abstraction for Git operations
  - `internal/git/backend.go` - Backend interface
  - `internal/git/exec.go` - Git CLI implementation (default)
  - `internal/git/gogit.go` - go-git implementation (future)
- [ ] **Commit Walker**: Iterate through repository history
  - Parse `git log` output with custom formats
  - Support date ranges, branch filtering, path filtering
- [ ] **Parsers**: Robust parsing for Git command output
  - `internal/parse/` - JSON and custom format parsers
  - Handle edge cases (merge commits, empty commits, etc.)

### Phase 2 — Domain Model
- [ ] **Core Types** (`pkg/nerds/types.go`):
  - `Repository` - Main entry point
  - `Commit` - Commit data
  - `Author` - Author/contributor information
  - `Branch` - Branch metadata
  - `Tag` - Tag information
  - `FileChange` - File modification data
- [ ] **Options** (`pkg/nerds/options.go`):
  - Date ranges (Since, Until)
  - Branch filtering
  - Path specifications (include/exclude)
  - Author filtering
  - Merge commit handling

### Phase 3 — Analysis Engines
- [ ] **Author Analytics** (`internal/analysis/authors.go`):
  - Commits per author
  - Lines added/deleted per author
  - Files modified per author
  - Contribution timeline
  - New contributors detection
  - Reviewer suggestions based on file history
- [ ] **Temporal Analysis** (`internal/analysis/temporal.go`):
  - Commits by day/week/month/year
  - Commits by hour of day
  - Commits by day of week
  - Commits by timezone
  - Activity heatmaps
  - Calendar views
- [ ] **Branch Analysis** (`internal/analysis/branches.go`):
  - Branch tree visualization
  - Branches by date
  - Branch age calculation
  - Merge frequency analysis
- [ ] **File Analysis** (`internal/analysis/files.go`):
  - File churn (add/delete frequency)
  - Hotspot detection (frequently modified files)
  - File size over time
  - Ownership by file

### Phase 4 — Public API
- [ ] **Repository** (`pkg/nerds/repository.go`):
  ```go
  type Repository interface {
      // Core Stats
      DetailedStats() (*Stats, error)
      StatsByBranch(branch string) (*Stats, error)

      // Authors
      Contributors() ([]Contributor, error)
      NewContributors(since time.Time) ([]Contributor, error)
      CommitsPerAuthor() (map[string]int, error)
      SuggestReviewers(file string) ([]string, error)

      // Temporal
      CommitsByDay() (map[string]int, error)
      CommitsByMonth() (map[string]int, error)
      CommitsByYear() (map[string]int, error)
      CommitsByWeekday() (map[string]int, error)
      CommitsByHour() (map[int]int, error)

      // Branches
      BranchTree() (*Tree, error)
      BranchesByDate() ([]Branch, error)

      // Export
      ExportJSON() (string, error)
      ExportCSV() (string, error)
      ExportMarkdown() (string, error)
  }
  ```
- [ ] **Factory Function** (`pkg/nerds/repository.go`):
  ```go
  func Open(path string, opts ...*Options) (*Repository, error)
  ```

### Phase 5 — Export Formats
- [ ] **JSON Export** (`pkg/nerds/export.go`):
  - Structured JSON with all statistics
  - Machine-readable format for downstream tools
- [ ] **CSV Export**:
  - Author statistics
  - Commit timeline
  - File changes
- [ ] **Markdown Export**:
  - Human-readable reports
  - Suitable for PR comments, documentation
  - Generated changelogs

### Phase 6 — Changelog Generation
- [ ] **Changelog Builder** (`internal/analysis/changelog.go`):
  - Overall repository changelog
  - Per-author changelogs
  - Between tags/releases
  - Conventional commit format support
  - Grouping by type (feat, fix, docs, etc.)

### Phase 7 — Performance & Optimization
- [ ] **Caching Strategy**:
  - Optional in-memory caching
  - Incremental analysis (avoid re-parsing)
- [ ] **Concurrency**:
  - Worker pools for parallel analysis
  - Stream processing for large repos
- [ ] **Benchmarking**:
  - Performance benchmarks for all analysis functions
  - Memory profiling
  - Optimization based on real-world repos

### Phase 8 — Testing & Quality
- [ ] **Unit Tests**:
  - 80%+ code coverage
  - Test all public API methods
- [ ] **Integration Tests**:
  - Test with real Git repositories (testdata/)
  - Golden tests for output formats
- [ ] **Fixture Repositories**:
  - Small test repos with known characteristics
  - Edge cases (empty repos, single commit, etc.)
- [ ] **Documentation**:
  - GoDoc for all public APIs
  - Usage examples
  - Integration guides

### Phase 9 — Optional CLI
- [ ] **CLI Wrapper** (`cmd/git-nerds/main.go`):
  - Thin wrapper around pkg/nerds
  - Primarily for testing and demonstration
  - Non-interactive command-line interface
  - Output to stdout (JSON/CSV/Markdown)

### Phase 10 — Advanced Features
- [ ] **Submodule Support**:
  - Detect submodules
  - Aggregate submodule statistics
- [ ] **Advanced Git Features**:
  - `git blame` integration
  - `git bisect` support for finding regressions
  - Worktree awareness
- [ ] **Custom Analyzers**:
  - Plugin system for custom analysis
  - Registry pattern for extensibility

---

## 3) Feature Priority Matrix

| Priority | Feature | Phase | Status |
|----------|---------|-------|--------|
| P0 | Git backend (exec) | 1 | In Progress |
| P0 | Core domain model | 2 | Pending |
| P0 | Public API design | 4 | Pending |
| P0 | JSON export | 5 | Pending |
| P1 | Author analytics | 3 | Pending |
| P1 | Temporal analysis | 3 | Pending |
| P1 | Branch analysis | 3 | Pending |
| P1 | Markdown export | 5 | Pending |
| P2 | File hotspot detection | 3 | Pending |
| P2 | Changelog generation | 6 | Pending |
| P2 | CSV export | 5 | Pending |
| P2 | Caching & optimization | 7 | Pending |
| P3 | CLI wrapper | 9 | Pending |
| P3 | Submodule support | 10 | Pending |
| P3 | Advanced Git features | 10 | Pending |

---

## 4) Module Structure

```
git-nerds/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── CONTRIBUTING.md
│
├── pkg/nerds/              # PUBLIC API
│   ├── repository.go       # Main Repository type and Open()
│   ├── options.go          # Configuration options
│   ├── types.go            # Public types (Stats, Author, etc.)
│   ├── export.go           # Export functions
│   └── errors.go           # Public error types
│
├── internal/               # INTERNAL IMPLEMENTATION
│   ├── git/
│   │   ├── backend.go      # Backend interface
│   │   ├── exec.go         # Git CLI backend
│   │   └── gogit.go        # go-git backend (future)
│   │
│   ├── analysis/
│   │   ├── authors.go      # Author analytics
│   │   ├── temporal.go     # Time-based analysis
│   │   ├── branches.go     # Branch analysis
│   │   ├── files.go        # File analysis
│   │   └── changelog.go    # Changelog generation
│   │
│   ├── parse/
│   │   ├── log.go          # Parse git log output
│   │   ├── diff.go         # Parse git diff output
│   │   └── branch.go       # Parse branch info
│   │
│   └── unixcompat/         # Unix tool replacements
│       ├── sort.go
│       └── filter.go
│
├── cmd/git-nerds/          # OPTIONAL CLI
│   └── main.go             # Thin wrapper for testing
│
└── testdata/               # TEST FIXTURES
    ├── repos/              # Sample Git repositories
    └── golden/             # Expected output files
```

---

## 5) API Design Examples

### Basic Usage
```go
package main

import (
    "fmt"
    "github.com/yourusername/git-nerds/pkg/nerds"
)

func main() {
    repo, err := nerds.Open(".")
    if err != nil {
        panic(err)
    }

    authors, err := repo.Contributors()
    if err != nil {
        panic(err)
    }

    for _, author := range authors {
        fmt.Printf("%s: %d commits\n", author.Name, author.Commits)
    }
}
```

### Advanced Usage with Options
```go
package main

import (
    "time"
    "github.com/yourusername/git-nerds/pkg/nerds"
)

func main() {
    opts := &nerds.Options{
        Since:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
        Until:         time.Now(),
        Branch:        "main",
        PathSpec:      []string{":!vendor", ":!*.generated.go"},
        IgnoreAuthors: []string{"bot@.*", ".*\\[bot\\]"},
        IncludeMerges: false,
    }

    repo, err := nerds.Open("/path/to/repo", opts)
    if err != nil {
        panic(err)
    }

    // Export to JSON
    json, err := repo.ExportJSON()
    // ... use json
}
```

### Integration Example
```go
// Code review tool integration
func suggestReviewersForPR(repoPath string, changedFiles []string) ([]string, error) {
    repo, err := nerds.Open(repoPath)
    if err != nil {
        return nil, err
    }

    reviewers := make(map[string]int)
    for _, file := range changedFiles {
        suggestions, err := repo.SuggestReviewers(file)
        if err != nil {
            continue
        }
        for _, reviewer := range suggestions {
            reviewers[reviewer]++
        }
    }

    // Return top 3 reviewers
    // ... sort and return
    return topReviewers, nil
}
```

---

## 6) Metrics Definitions

### Author Metrics
- **Commits**: Total commits by author
- **Lines Added**: Sum of insertions
- **Lines Deleted**: Sum of deletions
- **Lines Changed**: Insertions + Deletions
- **Files Modified**: Unique files touched
- **Active Days**: Days with at least one commit
- **First Seen**: Date of first commit
- **Last Seen**: Date of most recent commit

### Temporal Metrics
- **Commits per Day/Week/Month/Year**: Count by time period
- **Commits by Hour**: Distribution across 24 hours
- **Commits by Weekday**: Distribution Mon-Sun
- **Commits by Timezone**: Based on commit timestamp

### Branch Metrics
- **Branch Count**: Total branches
- **Branch Age**: Time since branch creation
- **Merge Frequency**: Merges per time period
- **Active Branches**: Branches with recent commits

### File Metrics
- **Churn**: Frequency of changes
- **Hotspots**: Files with highest churn
- **Size**: Current file size
- **Ownership**: Primary author by commit count

---

## 7) Testing Strategy

### Unit Tests
- Test each analysis function independently
- Mock Git backend for isolated testing
- Cover edge cases and error conditions

### Integration Tests
- Use real Git repositories in testdata/
- Verify end-to-end functionality
- Compare output with known-good results (golden tests)

### Performance Tests
- Benchmark critical paths
- Test with large repositories
- Memory usage profiling

### Example Test
```go
func TestAuthorStats(t *testing.T) {
    repo, err := nerds.Open("testdata/repos/simple")
    if err != nil {
        t.Fatal(err)
    }

    authors, err := repo.Contributors()
    if err != nil {
        t.Fatal(err)
    }

    // Verify expected authors
    if len(authors) != 3 {
        t.Errorf("expected 3 authors, got %d", len(authors))
    }
}
```

---

## 8) Release Plan

### v0.1.0 - Core Foundation
- Git backend (exec)
- Basic domain model
- Author analytics
- Simple JSON export
- Unit tests

### v0.2.0 - Temporal Analysis
- Commits by day/week/month/year
- Commits by hour/weekday
- Activity heatmaps
- Markdown export

### v0.3.0 - Branch & File Analysis
- Branch statistics
- File hotspot detection
- Reviewer suggestions
- CSV export

### v0.4.0 - Changelogs
- Changelog generation
- Per-author changelogs
- Release diff
- Conventional commits support

### v0.5.0 - Performance & Polish
- Caching implementation
- Performance optimizations
- Comprehensive documentation
- CLI wrapper

### v1.0.0 - Production Ready
- go-git backend support
- Extensive test coverage (>80%)
- Production-hardened
- Stable API

---

## 9) Documentation Requirements

### GoDoc
- [ ] All public types documented
- [ ] All public functions documented
- [ ] Package-level documentation
- [ ] Usage examples in doc comments

### README
- [x] Overview and features
- [x] Installation instructions
- [x] Quick start guide
- [x] API reference
- [x] Configuration examples
- [x] Use cases

### Additional Docs
- [ ] Integration guide
- [ ] Architecture overview
- [ ] Performance tuning guide
- [ ] Contributing guide

---

## 10) Success Criteria

### Functional
- ✅ Can analyze any Git repository
- ✅ Provides accurate statistics
- ✅ Exports in multiple formats
- ✅ Easy to integrate

### Non-Functional
- ✅ Zero dependencies for core features
- ✅ Fast (< 1s for typical repos)
- ✅ Memory efficient (< 100MB for typical repos)
- ✅ Well-tested (> 80% coverage)
- ✅ Well-documented

### Adoption
- ✅ Used by at least 3 different applications
- ✅ Positive community feedback
- ✅ Active maintenance and issue resolution

---

## 11) Risks & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Git output format changes | High | Use stable format options, version tests |
| Large repository performance | Medium | Implement caching, streaming |
| Complex merge history parsing | Medium | Comprehensive test fixtures |
| Cross-platform compatibility | Low | Test on Windows, Linux, macOS |
| Dependency bloat | Low | Minimize external dependencies |

---

## 12) Sprint Plan (2-week sprints)

- **Sprint 1**: Foundation + Git backend + domain model
- **Sprint 2**: Author analytics + basic export (JSON)
- **Sprint 3**: Temporal analysis + markdown export
- **Sprint 4**: Branch & file analysis + CSV export
- **Sprint 5**: Changelog generation + polish
- **Sprint 6**: Performance optimization + documentation + v0.5.0

---

## 13) v0.1.0 Launch Checklist

Core Features:
- [ ] Git backend interface implemented
- [ ] Exec backend working
- [ ] Repository.Open() function
- [ ] Author statistics (commits, lines, files)
- [ ] JSON export
- [ ] Basic configuration options

Quality:
- [ ] Unit tests for all public APIs
- [ ] Integration test with sample repo
- [ ] GoDoc documentation
- [ ] README with examples
- [ ] GitHub Actions CI

Release:
- [ ] Semantic versioning (v0.1.0)
- [ ] Git tag created
- [ ] CHANGELOG.md
- [ ] GitHub release notes
