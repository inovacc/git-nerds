# Git Nerds

A Go module that provides comprehensive intelligence about Git repositories.

---

## Overview

**Git Nerds** is a reusable Go library that extracts detailed statistics and insights from any local Git repository. Designed to be consumed by other applications, it provides a clean API for
repository analysis, contribution metrics, and historical insights.

This project was inspired by the original [`git-quick-stats`](https://github.com/git-quick-stats/git-quick-stats) Bash script, reimagined as a modern Go module for easy integration into any
application.

## Features

- **Author Analytics**: Commits, insertions, deletions, lines changed, files modified per contributor
- **Changelogs**: Generate overall and per-author changelogs
- **Temporal Analysis**: Daily, monthly, yearly, weekday, and hourly commit patterns
- **Visualizations**: Calendar heatmaps and activity visualizations
- **Branch Intelligence**: Branch trees, branch age, and timeline analysis
- **Contributor Insights**: Identify contributors, new contributors, and suggest reviewers
- **Multiple Backends**: Native Git command execution with optional go-git support
- **Export Options**: JSON, CSV, and structured data formats
- **Flexible API**: Easy-to-use Go interfaces for custom integrations

## Installation

```sh
go get github.com/inovacc/git-nerds
```

## Quick Start

```go
package main

import (
  "fmt"
  nerds "github.com/inovacc/git-nerds"
)

func main() {
  // Initialize repository analyzer
  repo, err := nerds.Open("/path/to/repository")
  if err != nil {
    panic(err)
  }

  // Get contributor statistics
  contributors, err := repo.Contributors()
  if err != nil {
    panic(err)
  }

  for _, c := range contributors {
    fmt.Printf("%s: %d commits\n", c.Name, c.Commits)
  }

  // Get commit activity by day
  _, err = repo.CommitsByDay()
  if err != nil {
    panic(err)
  }

  // Export to JSON
  json, err := repo.ExportJSON()
  if err != nil {
    panic(err)
  }
  fmt.Println(json)
}
```

## API Reference

### Core Methods

#### Repository Analysis

```go
repo.DetailedStats() (*Stats, error)                 // Comprehensive repository statistics
repo.StatsByBranch(branch string) (*Stats, error)    // Stats for a specific branch
repo.Changelogs() ([]Changelog, error)               // Generate changelogs
repo.ChangelogsByAuthor(author string) ([]Changelog, error) // Author-specific changelogs
```

#### Author Analytics

```go
repo.Contributors() ([]Contributor, error)                      // List all contributors
repo.NewContributors(since time.Time) ([]Contributor, error)    // New contributors since date
repo.CommitsPerAuthor() (map[string]int, error)                 // Commit count by author
repo.SuggestReviewers(file string) ([]string, error)            // Suggest reviewers for a file
```

#### Temporal Analysis

```go
repo.CommitsByDay() (map[string]int, error)        // Commits per day
repo.CommitsByMonth() (map[string]int, error)      // Commits per month
repo.CommitsByYear() (map[string]int, error)       // Commits per year
repo.CommitsByWeekday() (map[string]int, error)    // Commits per weekday
repo.CommitsByHour() (map[int]int, error)          // Commits per hour
repo.CommitsByTimezone() (map[string]int, error)   // Commits per timezone
```

#### Branch Analysis

```go
repo.BranchTree() (*Tree, error)           // ASCII graph of branch history
repo.BranchesByDate() ([]Branch, error)    // Branches sorted by date
```

#### Visualization

```go
repo.CommitsCalendar(author string) (*Calendar, error) // Calendar heatmap
repo.CommitsHeatmap(days int) (*Heatmap, error)        // Activity heatmap
```

#### Export

```go
repo.ExportJSON() (string, error)     // Export to JSON
repo.ExportCSV() (string, error)      // Export to CSV
repo.ExportMarkdown() (string, error) // Export as Markdown report
```

## Configuration

Configure repository analysis with flexible options:

```go
package main

import (
  "time"
  nerds "github.com/inovacc/git-nerds"
)

func main() {
  // Open repository with options
  repo, err := nerds.Open("/path/to/repo", &nerds.Options{
    Since:         time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
    Until:         time.Now(),
    Branch:        "main",
    Limit:         50,
    PathSpec:      []string{":!vendor", ":!node_modules"}, // Exclude paths
    IgnoreAuthors: []string{"bot@.*"},                      // Regex patterns
    IncludeMerges: true,
  })
  if err != nil {
    panic(err)
  }

  // Use the configured repository
  _, err = repo.DetailedStats()
  if err != nil {
    panic(err)
  }
}
```

## Project Structure

```
git-nerds/
├── README.md
├── go.mod
├── doc.go                 # Package documentation and examples
├── repository.go          # Public API: Repository type and core methods
├── options.go             # Configuration options
├── types.go               # Public types (Stats, Author, etc.)
├── export.go              # Export functionality
├── example/
│   └── main.go            # End-to-end usage examples
├── internal/
│   ├── git/               # Git backend implementations
│   │   ├── backend.go     # Backend interface
│   │   └── exec.go        # Git CLI execution (default)
│   ├── analysis/          # Analysis engines
│   ├── parse/             # Git output parsers
│   └── stats/             # Stats helpers and aggregations
└── ...

```

## Notes & Design Principles

- **Module-First**: Designed as a library, not a standalone application
- **Clean API**: Simple, intuitive Go interfaces
- **Backend**: Uses native Git CLI execution backend by default (internal/git/exec). A go-git backend may be added later.
- **Zero Dependencies**: Core functionality with minimal external deps
- **Testable**: Comprehensive test coverage with fixtures

## Use Cases

- **Code Review Tools**: Suggest reviewers based on file history
- **Analytics Dashboards**: Build custom Git analytics UIs
- **CI/CD Pipelines**: Generate automated reports and insights
- **Developer Tools**: Integrate repository intelligence into IDEs
- **Team Metrics**: Track team contributions and patterns
- **Documentation**: Auto-generate contributor lists and changelogs

## See Also

- [git-quick-stats (original Bash)](https://github.com/git-quick-stats/git-quick-stats) - Original inspiration

## License

[Add your license here]

## Contributing

Contributions are welcome! This is a module-first library designed for easy integration into any Go application.

