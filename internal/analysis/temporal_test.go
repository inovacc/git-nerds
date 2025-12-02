package analysis

import (
	"testing"

	"github.com/inovacc/git-nerds/internal/git"
)

func TestNewTemporalAnalyzer(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{}

	analyzer := NewTemporalAnalyzer(backend, opts)
	if analyzer == nil {
		t.Fatal("NewTemporalAnalyzer() returned nil")
	}
}

func TestCommitsByDay(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewTemporalAnalyzer(backend, opts)
	byDay, err := analyzer.CommitsByDay()
	if err != nil {
		t.Fatalf("CommitsByDay() error = %v", err)
	}

	if byDay == nil {
		t.Fatal("CommitsByDay() returned nil")
	}

	// Validate date format and counts
	for date, count := range byDay {
		if count <= 0 {
			t.Errorf("Date %s has non-positive count: %d", date, count)
		}
		// Date should be in YYYY-MM-DD format (length 10)
		if len(date) != 10 {
			t.Errorf("Invalid date format: %s", date)
		}
	}
}

func TestCommitsByMonth(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewTemporalAnalyzer(backend, opts)
	byMonth, err := analyzer.CommitsByMonth()
	if err != nil {
		t.Fatalf("CommitsByMonth() error = %v", err)
	}

	if byMonth == nil {
		t.Fatal("CommitsByMonth() returned nil")
	}

	// Validate month format and counts
	for month, count := range byMonth {
		if count <= 0 {
			t.Errorf("Month %s has non-positive count: %d", month, count)
		}
		// Month should be in YYYY-MM format (length 7)
		if len(month) != 7 {
			t.Errorf("Invalid month format: %s", month)
		}
	}
}

func TestCommitsByYear(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewTemporalAnalyzer(backend, opts)
	byYear, err := analyzer.CommitsByYear()
	if err != nil {
		t.Fatalf("CommitsByYear() error = %v", err)
	}

	if byYear == nil {
		t.Fatal("CommitsByYear() returned nil")
	}

	// Validate year format and counts
	for year, count := range byYear {
		if count <= 0 {
			t.Errorf("Year %s has non-positive count: %d", year, count)
		}
		// Year should be 4 digits
		if len(year) != 4 {
			t.Errorf("Invalid year format: %s", year)
		}
	}
}

func TestCommitsByWeekday(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewTemporalAnalyzer(backend, opts)
	byWeekday, err := analyzer.CommitsByWeekday()
	if err != nil {
		t.Fatalf("CommitsByWeekday() error = %v", err)
	}

	if byWeekday == nil {
		t.Fatal("CommitsByWeekday() returned nil")
	}

	validWeekdays := map[string]bool{
		"Monday": true, "Tuesday": true, "Wednesday": true,
		"Thursday": true, "Friday": true, "Saturday": true, "Sunday": true,
	}

	// Validate weekdays and counts
	for day, count := range byWeekday {
		if !validWeekdays[day] {
			t.Errorf("Invalid weekday: %s", day)
		}
		if count <= 0 {
			t.Errorf("Weekday %s has non-positive count: %d", day, count)
		}
	}
}

func TestCommitsByHour(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewTemporalAnalyzer(backend, opts)
	byHour, err := analyzer.CommitsByHour()
	if err != nil {
		t.Fatalf("CommitsByHour() error = %v", err)
	}

	if byHour == nil {
		t.Fatal("CommitsByHour() returned nil")
	}

	// Validate hours and counts
	for hour, count := range byHour {
		if hour < 0 || hour > 23 {
			t.Errorf("Invalid hour: %d", hour)
		}
		if count <= 0 {
			t.Errorf("Hour %d has non-positive count: %d", hour, count)
		}
	}
}

func TestCommitsByTimezone(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewTemporalAnalyzer(backend, opts)
	byTimezone, err := analyzer.CommitsByTimezone()
	if err != nil {
		t.Fatalf("CommitsByTimezone() error = %v", err)
	}

	if byTimezone == nil {
		t.Fatal("CommitsByTimezone() returned nil")
	}

	// Validate timezones and counts
	for timezone, count := range byTimezone {
		if timezone == "" {
			t.Error("Found empty timezone")
		}
		if count <= 0 {
			t.Errorf("Timezone %s has non-positive count: %d", timezone, count)
		}
	}
}

func TestGenerateHeatmap(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{}

	analyzer := NewTemporalAnalyzer(backend, opts)
	heatmap, err := analyzer.GenerateHeatmap(30)
	if err != nil {
		t.Fatalf("GenerateHeatmap() error = %v", err)
	}

	if heatmap == nil {
		t.Fatal("GenerateHeatmap() returned nil")
	}

	// Validate heatmap data
	for i, day := range heatmap.Days {
		if day.Date.IsZero() {
			t.Errorf("Day %d has zero date", i)
		}
		if day.CommitCount < 0 {
			t.Errorf("Day %d has negative commit count", i)
		}
	}

	// Days should be sorted chronologically
	for i := 1; i < len(heatmap.Days); i++ {
		if heatmap.Days[i].Date.Before(heatmap.Days[i-1].Date) {
			t.Error("Heatmap days are not sorted chronologically")
			break
		}
	}
}

func TestGenerateCalendar(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{}

	analyzer := NewTemporalAnalyzer(backend, opts)
	calendar, err := analyzer.GenerateCalendar(2024, "")
	if err != nil {
		t.Fatalf("GenerateCalendar() error = %v", err)
	}

	if calendar == nil {
		t.Fatal("GenerateCalendar() returned nil")
	}

	if calendar.Year != 2024 {
		t.Errorf("Expected year 2024, got %d", calendar.Year)
	}

	if len(calendar.Months) != 12 {
		t.Errorf("Expected 12 months, got %d", len(calendar.Months))
	}
}

func TestGetCommitTrend(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{
		Limit: 100,
	}

	analyzer := NewTemporalAnalyzer(backend, opts)

	tests := []struct {
		name   string
		period string
	}{
		{"by day", "day"},
		{"by month", "month"},
		{"by year", "year"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trend, err := analyzer.GetCommitTrend(tt.period)
			if err != nil {
				t.Fatalf("GetCommitTrend(%s) error = %v", tt.period, err)
			}

			if trend == nil {
				t.Fatal("GetCommitTrend() returned nil")
			}

			// Validate trend data
			for i, item := range trend {
				if item.Period == "" {
					t.Errorf("Trend item %d has empty period", i)
				}
				if item.Count < 0 {
					t.Errorf("Trend item %d has negative count", i)
				}
			}

			// Check sorting
			for i := 1; i < len(trend); i++ {
				if trend[i].Period < trend[i-1].Period {
					t.Error("Trend is not sorted by period")
					break
				}
			}
		})
	}
}

func TestGetCommitTrendInvalidPeriod(t *testing.T) {
	backend := setupTestBackend(t)
	opts := &git.LogOptions{}

	analyzer := NewTemporalAnalyzer(backend, opts)
	_, err := analyzer.GetCommitTrend("invalid")
	if err == nil {
		t.Error("Expected error for invalid period, got nil")
	}
}

// Benchmark tests
func BenchmarkCommitsByDay(b *testing.B) {
	backend, err := git.NewExecBackend("../../..")
	if err != nil {
		b.Skip("Cannot create git backend")
	}

	opts := &git.LogOptions{Limit: 100}
	analyzer := NewTemporalAnalyzer(backend, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.CommitsByDay()
	}
}

func BenchmarkCommitsByWeekday(b *testing.B) {
	backend, err := git.NewExecBackend("../../..")
	if err != nil {
		b.Skip("Cannot create git backend")
	}

	opts := &git.LogOptions{Limit: 100}
	analyzer := NewTemporalAnalyzer(backend, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.CommitsByWeekday()
	}
}

func BenchmarkGenerateHeatmap(b *testing.B) {
	backend, err := git.NewExecBackend("../../..")
	if err != nil {
		b.Skip("Cannot create git backend")
	}

	opts := &git.LogOptions{}
	analyzer := NewTemporalAnalyzer(backend, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = analyzer.GenerateHeatmap(30)
	}
}
