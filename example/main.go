package main

import (
	"fmt"
	"log"
	"time"

	gitnerds "github.com/inovacc/git-nerds"
)

func main() {
	// Example 1: Basic usage - open current directory
	fmt.Println("=== Example 1: Basic Repository Stats ===")
	basicExample()

	fmt.Println("\n=== Example 2: Author Analytics ===")
	authorExample()

	fmt.Println("\n=== Example 3: Temporal Analysis ===")
	temporalExample()

	fmt.Println("\n=== Example 4: Export to JSON ===")
	exportExample()

	fmt.Println("=== GIT NERDS MODULE OUTPUT ===\n")

	repo, err := gitnerds.Open("..")
	if err != nil {
		log.Fatal(err)
	}

	// 1. Detailed Stats
	fmt.Println("1. DETAILED REPOSITORY STATS")
	fmt.Println("─────────────────────────────")
	stats, err := repo.DetailedStats()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Total Commits:    %d\n", stats.TotalCommits)
		fmt.Printf("Total Authors:    %d\n", stats.TotalAuthors)
		fmt.Printf("Lines Added:      %d\n", stats.LinesAdded)
		fmt.Printf("Lines Deleted:    %d\n", stats.LinesDeleted)
		fmt.Printf("Lines Changed:    %d\n", stats.LinesChanged)
		fmt.Printf("First Commit:     %s\n", stats.FirstCommitAt.Format("2006-01-02 15:04:05 -0700"))
		fmt.Printf("Last Commit:      %s\n", stats.LastCommitAt.Format("2006-01-02 15:04:05 -0700"))
		fmt.Println()
	}

	// 2. Contributors
	fmt.Println("2. CONTRIBUTORS")
	fmt.Println("───────────────")
	contributors, err := repo.Contributors()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for i, c := range contributors {
			fmt.Printf("%d. %s <%s>\n", i+1, c.Name, c.Email)
			fmt.Printf("   Commits: %d, Since: %s\n", c.Commits, c.Since.Format("2006-01-02"))
		}
		fmt.Println()
	}

	// 3. Commits Per Author
	fmt.Println("3. COMMITS PER AUTHOR")
	fmt.Println("─────────────────────")
	commitsPerAuthor, err := repo.CommitsPerAuthor()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for author, count := range commitsPerAuthor {
			fmt.Printf("%-30s %d commits\n", author, count)
		}
		fmt.Println()
	}

	// 4. Commits by Weekday
	fmt.Println("4. COMMITS BY WEEKDAY")
	fmt.Println("─────────────────────")
	byWeekday, err := repo.CommitsByWeekday()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		for _, day := range weekdays {
			if count, ok := byWeekday[day]; ok {
				fmt.Printf("%-12s %d commits\n", day, count)
			}
		}
		fmt.Println()
	}

	// 5. Commits by Hour
	fmt.Println("5. COMMITS BY HOUR (Top 5)")
	fmt.Println("──────────────────────────")
	byHour, err := repo.CommitsByHour()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		type hourCount struct {
			hour  int
			count int
		}
		hours := make([]hourCount, 0)
		for h, c := range byHour {
			hours = append(hours, hourCount{h, c})
		}
		// Simple bubble sort for top 5
		for i := 0; i < len(hours); i++ {
			for j := i + 1; j < len(hours); j++ {
				if hours[j].count > hours[i].count {
					hours[i], hours[j] = hours[j], hours[i]
				}
			}
		}
		limit := 5
		if len(hours) < limit {
			limit = len(hours)
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("%02d:00        %d commits\n", hours[i].hour, hours[i].count)
		}
		fmt.Println()
	}

	// 6. Commits by Month
	fmt.Println("6. COMMITS BY MONTH")
	fmt.Println("───────────────────")
	byMonth, err := repo.CommitsByMonth()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for month, count := range byMonth {
			fmt.Printf("%-10s %d commits\n", month, count)
		}
		fmt.Println()
	}

	// 7. Commits by Year
	fmt.Println("7. COMMITS BY YEAR")
	fmt.Println("──────────────────")
	byYear, err := repo.CommitsByYear()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for year, count := range byYear {
			fmt.Printf("%-6s %d commits\n", year, count)
		}
		fmt.Println()
	}

	// 8. Author Details
	if stats != nil && len(stats.Authors) > 0 {
		fmt.Println("8. DETAILED AUTHOR STATS")
		fmt.Println("────────────────────────")
		for i, author := range stats.Authors {
			fmt.Printf("%d. %s <%s>\n", i+1, author.Name, author.Email)
			fmt.Printf("   Commits:       %d\n", author.Commits)
			fmt.Printf("   Lines Added:   %d\n", author.LinesAdded)
			fmt.Printf("   Lines Deleted: %d\n", author.LinesDeleted)
			fmt.Printf("   Lines Changed: %d\n", author.LinesChanged)
			fmt.Printf("   Files Changed: %d\n", author.FilesChanged)
			fmt.Printf("   First Commit:  %s\n", author.FirstCommit.Format("2006-01-02"))
			fmt.Printf("   Last Commit:   %s\n", author.LastCommit.Format("2006-01-02"))
			fmt.Println()
		}
	}

	// 9. Branches
	fmt.Println("9. BRANCHES")
	fmt.Println("───────────")
	branches, err := repo.BranchesByDate()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for i, branch := range branches {
			fmt.Printf("%d. %s (updated: %s)\n", i+1, branch.Name, branch.UpdatedAt.Format("2006-01-02"))
		}
		fmt.Println()
	}
}

func basicExample() {
	// Open repository at current directory
	repo, err := gitnerds.Open(".")
	if err != nil {
		log.Printf("Failed to open repository: %v", err)
		return
	}

	// Get detailed statistics
	stats, err := repo.DetailedStats()
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
		return
	}

	fmt.Printf("Total Commits: %d\n", stats.TotalCommits)
	fmt.Printf("Total Authors: %d\n", stats.TotalAuthors)
	fmt.Printf("Lines Added: %d\n", stats.LinesAdded)
	fmt.Printf("Lines Deleted: %d\n", stats.LinesDeleted)
	fmt.Printf("Lines Changed: %d\n", stats.LinesChanged)

	// Show top 5 contributors
	fmt.Println("\nTop 5 Contributors:")
	limit := 5
	if len(stats.Authors) < limit {
		limit = len(stats.Authors)
	}
	for i := 0; i < limit; i++ {
		author := stats.Authors[i]
		fmt.Printf("  %d. %s - %d commits, %d lines changed\n",
			i+1, author.Name, author.Commits, author.LinesChanged)
	}
}

func authorExample() {
	repo, err := gitnerds.Open(".")
	if err != nil {
		log.Printf("Failed to open repository: %v", err)
		return
	}

	// Get all contributors
	contributors, err := repo.Contributors()
	if err != nil {
		log.Printf("Failed to get contributors: %v", err)
		return
	}

	fmt.Printf("Total Contributors: %d\n", len(contributors))

	// Get commits per author
	commitsPerAuthor, err := repo.CommitsPerAuthor()
	if err != nil {
		log.Printf("Failed to get commits per author: %v", err)
		return
	}

	fmt.Println("\nCommits per author:")
	count := 0
	for author, commits := range commitsPerAuthor {
		if count >= 5 { // Show top 5
			break
		}
		fmt.Printf("  %s: %d commits\n", author, commits)
		count++
	}

	// Get new contributors in the last 90 days
	since := time.Now().AddDate(0, 0, -90)
	newContributors, err := repo.NewContributors(since)
	if err != nil {
		log.Printf("Failed to get new contributors: %v", err)
		return
	}

	if len(newContributors) > 0 {
		fmt.Printf("\nNew contributors in last 90 days: %d\n", len(newContributors))
		for _, c := range newContributors {
			fmt.Printf("  %s (since %s)\n", c.Name, c.Since.Format("2006-01-02"))
		}
	}
}

func temporalExample() {
	repo, err := gitnerds.Open(".")
	if err != nil {
		log.Printf("Failed to open repository: %v", err)
		return
	}

	// Commits by weekday
	byWeekday, err := repo.CommitsByWeekday()
	if err != nil {
		log.Printf("Failed to get commits by weekday: %v", err)
		return
	}

	fmt.Println("Commits by weekday:")
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range weekdays {
		if count, ok := byWeekday[day]; ok {
			fmt.Printf("  %s: %d commits\n", day, count)
		}
	}

	// Commits by hour (show peak hours)
	byHour, err := repo.CommitsByHour()
	if err != nil {
		log.Printf("Failed to get commits by hour: %v", err)
		return
	}

	var peakHour int
	var peakCount int
	for hour, count := range byHour {
		if count > peakCount {
			peakCount = count
			peakHour = hour
		}
	}
	fmt.Printf("\nPeak commit hour: %02d:00 (%d commits)\n", peakHour, peakCount)

	// Commits by month (last 6 months)
	byMonth, err := repo.CommitsByMonth()
	if err != nil {
		log.Printf("Failed to get commits by month: %v", err)
		return
	}

	fmt.Println("\nRecent monthly activity:")
	count := 0
	for month, commits := range byMonth {
		if count >= 6 {
			break
		}
		fmt.Printf("  %s: %d commits\n", month, commits)
		count++
	}
}

func exportExample() {
	repo, err := gitnerds.Open(".")
	if err != nil {
		log.Printf("Failed to open repository: %v", err)
		return
	}

	// Export to JSON
	json, err := repo.ExportJSON()
	if err != nil {
		log.Printf("Failed to export JSON: %v", err)
		return
	}

	fmt.Println("JSON export (first 500 chars):")
	if len(json) > 500 {
		fmt.Println(json[:500] + "...")
	} else {
		fmt.Println(json)
	}

	// Export to Markdown
	markdown, err := repo.ExportMarkdown()
	if err != nil {
		log.Printf("Failed to export Markdown: %v", err)
		return
	}

	fmt.Println("\nMarkdown export (first 500 chars):")
	if len(markdown) > 500 {
		fmt.Println(markdown[:500] + "...")
	} else {
		fmt.Println(markdown)
	}
}
