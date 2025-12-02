package main

import (
	"fmt"
	"log"

	gitnerds "github.com/inovacc/git-nerds"
)

func main() {
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
