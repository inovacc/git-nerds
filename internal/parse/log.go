package parse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CommitInfo represents parsed commit information
type CommitInfo struct {
	Hash      string
	Author    string
	Email     string
	Date      time.Time
	Subject   string
	Body      string
	Additions int
	Deletions int
	Files     []string
}

// AuthorInfo represents parsed author information
type AuthorInfo struct {
	Name  string
	Email string
	Count int
}

// ParseCommitLog parses git log output with custom format
// Expected format: hash|author|email|date|subject
func ParseCommitLog(output string) ([]CommitInfo, error) {
	if output == "" {
		return []CommitInfo{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	commits := make([]CommitInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue // Skip malformed lines
		}

		// Parse date
		date, err := time.Parse("2006-01-02 15:04:05 -0700", parts[3])
		if err != nil {
			// Try alternative formats
			date, err = time.Parse(time.RFC3339, parts[3])
			if err != nil {
				continue // Skip if date parsing fails
			}
		}

		commits = append(commits, CommitInfo{
			Hash:    parts[0],
			Author:  parts[1],
			Email:   parts[2],
			Date:    date,
			Subject: parts[4],
		})
	}

	return commits, nil
}

// ParseAuthors parses author names from git log output
func ParseAuthors(output string) []string {
	if output == "" {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	authors := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			authors = append(authors, line)
		}
	}

	return authors
}

// ParseAuthorStats parses author statistics from git shortlog output
// Expected format: "   123  Author Name"
func ParseAuthorStats(output string) ([]AuthorInfo, error) {
	if output == "" {
		return []AuthorInfo{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	stats := make([]AuthorInfo, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split on first tab or multiple spaces
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		count, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		// Rest is author name
		author := strings.Join(parts[1:], " ")

		stats = append(stats, AuthorInfo{
			Name:  author,
			Count: count,
		})
	}

	return stats, nil
}

// ParseDateCounts parses date-based commit counts
func ParseDateCounts(output string) (map[string]int, error) {
	if output == "" {
		return map[string]int{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	counts := make(map[string]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			counts[line]++
		}
	}

	return counts, nil
}

// ParseNumstat parses git log --numstat output
// Format: additions\tdeletions\tfilename
func ParseNumstat(output string) ([]FileStats, error) {
	if output == "" {
		return []FileStats{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	stats := make([]FileStats, 0)

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}

		additions, err1 := strconv.Atoi(parts[0])
		deletions, err2 := strconv.Atoi(parts[1])

		// Handle binary files (marked as "-")
		if err1 != nil || err2 != nil {
			additions = 0
			deletions = 0
		}

		stats = append(stats, FileStats{
			File:      parts[2],
			Additions: additions,
			Deletions: deletions,
		})
	}

	return stats, nil
}

// FileStats represents file change statistics
type FileStats struct {
	File      string
	Additions int
	Deletions int
}

// ParseBranches parses git branch output
func ParseBranches(output string) ([]string, error) {
	if output == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		// Remove leading * and whitespace
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
		}

		if line != "" {
			branches = append(branches, line)
		}
	}

	return branches, nil
}

// ParseTags parses git tag output
func ParseTags(output string) ([]string, error) {
	if output == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	tags := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			tags = append(tags, line)
		}
	}

	return tags, nil
}

// ParseWeekday extracts weekday from date string
func ParseWeekday(dateStr string) (time.Weekday, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Sunday, fmt.Errorf("invalid date format: %w", err)
	}
	return date.Weekday(), nil
}

// ParseHour extracts hour from datetime string
func ParseHour(datetimeStr string) (int, error) {
	date, err := time.Parse("2006-01-02 15:04:05", datetimeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid datetime format: %w", err)
	}
	return date.Hour(), nil
}
