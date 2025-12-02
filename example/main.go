package main

import (
	"fmt"
	"log"
	"time"

	"github.com/inovacc/git-nerds/pkg/nerds"
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
}

func basicExample() {
	// Open repository at current directory
	repo, err := nerds.Open(".")
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
	repo, err := nerds.Open(".")
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
	repo, err := nerds.Open(".")
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
	repo, err := nerds.Open(".")
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
