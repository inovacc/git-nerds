package analysis

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/inovacc/git-nerds/internal/git"
	"github.com/inovacc/git-nerds/internal/parse"
)

// TemporalAnalyzer provides time-based analytics
type TemporalAnalyzer struct {
	backend git.Backend
	options *git.LogOptions
}

// NewTemporalAnalyzer creates a new temporal analyzer
func NewTemporalAnalyzer(backend git.Backend, options *git.LogOptions) *TemporalAnalyzer {
	return &TemporalAnalyzer{
		backend: backend,
		options: options,
	}
}

// CommitsByDay returns commits grouped by day
func (t *TemporalAnalyzer) CommitsByDay() (map[string]int, error) {
	opts := *t.options
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=short"}, args...)

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits by day: %w", err)
	}

	return parse.ParseDateCounts(output)
}

// CommitsByMonth returns commits grouped by month (YYYY-MM format)
func (t *TemporalAnalyzer) CommitsByMonth() (map[string]int, error) {
	opts := *t.options
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=format:%Y-%m"}, args...)

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits by month: %w", err)
	}

	return parse.ParseDateCounts(output)
}

// CommitsByYear returns commits grouped by year
func (t *TemporalAnalyzer) CommitsByYear() (map[string]int, error) {
	opts := *t.options
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=format:%Y"}, args...)

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits by year: %w", err)
	}

	return parse.ParseDateCounts(output)
}

// CommitsByWeekday returns commits grouped by weekday
func (t *TemporalAnalyzer) CommitsByWeekday() (map[string]int, error) {
	opts := *t.options
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=format:%A"}, args...) // Full weekday name

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits by weekday: %w", err)
	}

	return parse.ParseDateCounts(output)
}

// CommitsByHour returns commits grouped by hour (0-23)
func (t *TemporalAnalyzer) CommitsByHour() (map[int]int, error) {
	opts := *t.options
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=format:%H"}, args...)

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits by hour: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	result := make(map[int]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var hour int
		if _, err := fmt.Sscanf(line, "%d", &hour); err == nil {
			result[hour]++
		}
	}

	return result, nil
}

// CommitsByTimezone returns commits grouped by timezone
func (t *TemporalAnalyzer) CommitsByTimezone() (map[string]int, error) {
	opts := *t.options
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=format:%z"}, args...) // Timezone offset

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits by timezone: %w", err)
	}

	return parse.ParseDateCounts(output)
}

// ActivityHeatmap represents commit activity heatmap data
type ActivityHeatmap struct {
	Days []DayActivity
}

// DayActivity represents commit activity for a single day
type DayActivity struct {
	Date            time.Time
	Weekday         time.Weekday
	CommitCount     int
	HourlyBreakdown map[int]int
}

// GenerateHeatmap generates a heatmap for the last N days
func (t *TemporalAnalyzer) GenerateHeatmap(days int) (*ActivityHeatmap, error) {
	since := time.Now().AddDate(0, 0, -days)
	opts := *t.options
	opts.Since = since
	opts.Format = "%ad|%H"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=format:%Y-%m-%d|%H"}, args...)

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate heatmap: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Map to store day -> hour -> count
	dayMap := make(map[string]map[int]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		dateStr := parts[0]
		var hour int
		fmt.Sscanf(parts[1], "%d", &hour)

		if _, exists := dayMap[dateStr]; !exists {
			dayMap[dateStr] = make(map[int]int)
		}
		dayMap[dateStr][hour]++
	}

	// Convert to slice
	result := &ActivityHeatmap{
		Days: make([]DayActivity, 0, len(dayMap)),
	}

	for dateStr, hourMap := range dayMap {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		total := 0
		for _, count := range hourMap {
			total += count
		}

		result.Days = append(result.Days, DayActivity{
			Date:            date,
			Weekday:         date.Weekday(),
			CommitCount:     total,
			HourlyBreakdown: hourMap,
		})
	}

	// Sort by date
	sort.Slice(result.Days, func(i, j int) bool {
		return result.Days[i].Date.Before(result.Days[j].Date)
	})

	return result, nil
}

// CalendarData represents a commit calendar for visualization
type CalendarData struct {
	Year   int
	Months []MonthData
}

// MonthData represents commit data for a specific month
type MonthData struct {
	Month      time.Month
	WeekMatrix [][]int // week -> day of week -> commit count
}

// GenerateCalendar generates a calendar view for a specific year and author
func (t *TemporalAnalyzer) GenerateCalendar(year int, author string) (*CalendarData, error) {
	opts := *t.options
	opts.Since = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	opts.Until = time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	if author != "" {
		opts.Author = author
	}
	opts.Format = "%ad"
	args := git.BuildLogArgs(&opts)
	args = append([]string{"--date=short"}, args...)

	output, err := t.backend.Log(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate calendar: %w", err)
	}

	dateCounts, err := parse.ParseDateCounts(output)
	if err != nil {
		return nil, err
	}

	calendar := &CalendarData{
		Year:   year,
		Months: make([]MonthData, 12),
	}

	// Initialize months
	for i := 0; i < 12; i++ {
		calendar.Months[i] = MonthData{
			Month:      time.Month(i + 1),
			WeekMatrix: make([][]int, 0),
		}
	}

	// Populate with commit data
	for dateStr, count := range dateCounts {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil || date.Year() != year {
			continue
		}

		monthIdx := int(date.Month()) - 1
		// Note: Calendar matrix population would require more complex logic
		// This is a simplified version
		_ = monthIdx
		_ = count
	}

	return calendar, nil
}

// CommitTrend represents commit trend over time
type CommitTrend struct {
	Period string
	Count  int
}

// GetCommitTrend returns commit trends sorted by period
func (t *TemporalAnalyzer) GetCommitTrend(period string) ([]CommitTrend, error) {
	var data map[string]int
	var err error

	switch period {
	case "day":
		data, err = t.CommitsByDay()
	case "month":
		data, err = t.CommitsByMonth()
	case "year":
		data, err = t.CommitsByYear()
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	result := make([]CommitTrend, 0, len(data))
	for period, count := range data {
		result = append(result, CommitTrend{
			Period: period,
			Count:  count,
		})
	}

	// Sort by period
	sort.Slice(result, func(i, j int) bool {
		return result[i].Period < result[j].Period
	})

	return result, nil
}
