// Package nerds provides comprehensive Git repository analysis and statistics.
//
// Git Nerds is a Go module that extracts detailed insights from Git repositories,
// including author statistics, temporal analysis, branch information, and more.
//
// # Basic Usage
//
// Open a repository and get basic statistics:
//
//	repo, err := nerds.Open(".")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	stats, err := repo.DetailedStats()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Total commits: %d\n", stats.TotalCommits)
//	fmt.Printf("Total authors: %d\n", stats.TotalAuthors)
//
// # Configuration
//
// Configure analysis with Options:
//
//	opts := &nerds.Options{
//		Since:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//		Until:         time.Now(),
//		Branch:        "main",
//		PathSpec:      []string{":!vendor", ":!node_modules"},
//		IgnoreAuthors: []string{"bot@.*"},
//		IncludeMerges: false,
//	}
//
//	repo, err := nerds.Open("/path/to/repo", opts)
//
// # Author Analytics
//
// Get contributor information:
//
//	contributors, err := repo.Contributors()
//	for _, c := range contributors {
//		fmt.Printf("%s: %d commits\n", c.Name, c.Commits)
//	}
//
//	// Get commits per author
//	counts, err := repo.CommitsPerAuthor()
//
//	// Suggest reviewers for a file
//	reviewers, err := repo.SuggestReviewers("main.go")
//
// # Temporal Analysis
//
// Analyze commit patterns over time:
//
//	// Commits by day
//	byDay, err := repo.CommitsByDay()
//
//	// Commits by month
//	byMonth, err := repo.CommitsByMonth()
//
//	// Commits by weekday
//	byWeekday, err := repo.CommitsByWeekday()
//
//	// Commits by hour
//	byHour, err := repo.CommitsByHour()
//
//	// Generate heatmap for last 30 days
//	heatmap, err := repo.CommitsHeatmap(30)
//
// # Branch Analysis
//
// Analyze branches:
//
//	branches, err := repo.BranchesByDate()
//	for _, b := range branches {
//		fmt.Printf("%s: last commit %s\n", b.Name, b.UpdatedAt)
//	}
//
//	tree, err := repo.BranchTree()
//
// # Export
//
// Export statistics in various formats:
//
//	// JSON export
//	json, err := repo.ExportJSON()
//	fmt.Println(json)
//
//	// Markdown report
//	markdown, err := repo.ExportMarkdown()
//	fmt.Println(markdown)
//
//	// CSV export
//	csv, err := repo.ExportCSV()
//	fmt.Println(csv)
//
// # Integration Example
//
// Integrate with a code review system:
//
//	func suggestReviewersForPR(repoPath string, files []string) ([]string, error) {
//		repo, err := nerds.Open(repoPath)
//		if err != nil {
//			return nil, err
//		}
//
//		reviewers := make(map[string]int)
//		for _, file := range files {
//			suggestions, _ := repo.SuggestReviewers(file)
//			for _, r := range suggestions {
//				reviewers[r]++
//			}
//		}
//
//		// Return top reviewers
//		// ...
//		return topReviewers, nil
//	}
package git_nerds
