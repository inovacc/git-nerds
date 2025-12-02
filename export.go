package git_nerds

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ExportJSON exports repository statistics to JSON format
func (r *Repository) ExportJSON() (string, error) {
	stats, err := r.DetailedStats()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(data), nil
}

// ExportCSV exports repository statistics to CSV format
func (r *Repository) ExportCSV() (string, error) {
	stats, err := r.DetailedStats()
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write author statistics
	writer.Write([]string{"Author", "Email", "Commits", "Lines Added", "Lines Deleted", "Files Changed"})
	for _, author := range stats.Authors {
		writer.Write([]string{
			author.Name,
			author.Email,
			fmt.Sprintf("%d", author.Commits),
			fmt.Sprintf("%d", author.LinesAdded),
			fmt.Sprintf("%d", author.LinesDeleted),
			fmt.Sprintf("%d", author.FilesChanged),
		})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to write CSV: %w", err)
	}

	return buf.String(), nil
}

// ExportMarkdown exports repository statistics to Markdown format
func (r *Repository) ExportMarkdown() (string, error) {
	stats, err := r.DetailedStats()
	if err != nil {
		return "", err
	}

	var md strings.Builder

	md.WriteString("# Repository Statistics\n\n")
	md.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Overall stats
	md.WriteString("## Overview\n\n")
	md.WriteString(fmt.Sprintf("- **Total Commits:** %d\n", stats.TotalCommits))
	md.WriteString(fmt.Sprintf("- **Total Authors:** %d\n", stats.TotalAuthors))
	md.WriteString(fmt.Sprintf("- **Total Files:** %d\n", stats.TotalFiles))
	md.WriteString(fmt.Sprintf("- **Lines Added:** %d\n", stats.LinesAdded))
	md.WriteString(fmt.Sprintf("- **Lines Deleted:** %d\n", stats.LinesDeleted))
	md.WriteString(fmt.Sprintf("- **Active Days:** %d\n", stats.ActiveDays))
	md.WriteString("\n")

	// Author statistics
	md.WriteString("## Contributors\n\n")
	md.WriteString("| Author | Commits | Lines Added | Lines Deleted | Files Changed |\n")
	md.WriteString("|--------|---------|-------------|---------------|---------------|\n")
	for _, author := range stats.Authors {
		md.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d |\n",
			author.Name,
			author.Commits,
			author.LinesAdded,
			author.LinesDeleted,
			author.FilesChanged,
		))
	}
	md.WriteString("\n")

	// Top files
	if len(stats.Files) > 0 {
		md.WriteString("## Most Modified Files\n\n")
		md.WriteString("| File | Changes | Additions | Deletions |\n")
		md.WriteString("|------|---------|-----------|------------|\n")
		limit := 10
		if len(stats.Files) < limit {
			limit = len(stats.Files)
		}
		for i := 0; i < limit; i++ {
			file := stats.Files[i]
			md.WriteString(fmt.Sprintf("| %s | %d | %d | %d |\n",
				file.Path,
				file.Changes,
				file.Additions,
				file.Deletions,
			))
		}
		md.WriteString("\n")
	}

	return md.String(), nil
}
