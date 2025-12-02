package parse

import (
	"reflect"
	"testing"
	"time"
)

func TestParseCommitLog(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int // number of commits expected
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "single commit",
			input:   "abc123|John Doe|john@example.com|2024-01-01 10:00:00 +0000|Initial commit",
			want:    1,
			wantErr: false,
		},
		{
			name: "multiple commits",
			input: `abc123|John Doe|john@example.com|2024-01-01 10:00:00 +0000|Initial commit
def456|Jane Smith|jane@example.com|2024-01-02 11:00:00 +0000|Add feature
ghi789|Bob Johnson|bob@example.com|2024-01-03 12:00:00 +0000|Fix bug`,
			want:    3,
			wantErr: false,
		},
		{
			name: "malformed line",
			input: `abc123|John Doe|john@example.com|2024-01-01 10:00:00 +0000|Initial commit
invalid line
def456|Jane Smith|jane@example.com|2024-01-02 11:00:00 +0000|Add feature`,
			want:    2, // Should skip invalid line
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommitLog(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommitLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.want {
				t.Errorf("ParseCommitLog() returned %d commits, want %d", len(got), tt.want)
			}

			// Validate first commit if present
			if len(got) > 0 {
				commit := got[0]
				if commit.Hash == "" {
					t.Error("Commit hash is empty")
				}
				if commit.Author == "" {
					t.Error("Commit author is empty")
				}
				if commit.Date.IsZero() {
					t.Error("Commit date is zero")
				}
			}
		})
	}
}

func TestParseAuthors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty input",
			input: "",
			want:  []string{},
		},
		{
			name:  "single author",
			input: "John Doe",
			want:  []string{"John Doe"},
		},
		{
			name: "multiple authors",
			input: `John Doe
Jane Smith
Bob Johnson`,
			want: []string{"John Doe", "Jane Smith", "Bob Johnson"},
		},
		{
			name: "with empty lines",
			input: `John Doe

Jane Smith

Bob Johnson`,
			want: []string{"John Doe", "Jane Smith", "Bob Johnson"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseAuthors(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAuthors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAuthorStats(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int // number of authors expected
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "single author",
			input:   "   100  John Doe",
			want:    1,
			wantErr: false,
		},
		{
			name: "multiple authors",
			input: `   100  John Doe
    50  Jane Smith
    25  Bob Johnson`,
			want:    3,
			wantErr: false,
		},
		{
			name:    "tab separated",
			input:   "100	John Doe",
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAuthorStats(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAuthorStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.want {
				t.Errorf("ParseAuthorStats() returned %d authors, want %d", len(got), tt.want)
			}

			// Validate first author if present
			if len(got) > 0 {
				author := got[0]
				if author.Name == "" {
					t.Error("Author name is empty")
				}
				if author.Count <= 0 {
					t.Error("Author count should be positive")
				}
			}
		})
	}
}

func TestParseDateCounts(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]int
	}{
		{
			name:  "empty input",
			input: "",
			want:  map[string]int{},
		},
		{
			name:  "single date",
			input: "2024-01-01",
			want:  map[string]int{"2024-01-01": 1},
		},
		{
			name: "multiple dates",
			input: `2024-01-01
2024-01-01
2024-01-02
2024-01-03`,
			want: map[string]int{
				"2024-01-01": 2,
				"2024-01-02": 1,
				"2024-01-03": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateCounts(tt.input)
			if err != nil {
				t.Errorf("ParseDateCounts() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDateCounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNumstat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int // number of file stats expected
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "single file",
			input:   "10	5	main.go",
			want:    1,
			wantErr: false,
		},
		{
			name:    "multiple files",
			input:   "10	5	main.go\n20	8	utils.go\n5	2	README.md",
			want:    3,
			wantErr: false,
		},
		{
			name:    "binary file",
			input:   "-	-	image.png",
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNumstat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNumstat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.want {
				t.Errorf("ParseNumstat() returned %d file stats, want %d", len(got), tt.want)
			}

			// Validate first file stat if present
			if len(got) > 0 {
				stat := got[0]
				if stat.File == "" {
					t.Error("File name is empty")
				}
				if stat.Additions < 0 || stat.Deletions < 0 {
					t.Error("Additions/deletions should be non-negative")
				}
			}
		})
	}
}

func TestParseBranches(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty input",
			input: "",
			want:  []string{},
		},
		{
			name:  "single branch",
			input: "main",
			want:  []string{"main"},
		},
		{
			name:  "current branch (with asterisk)",
			input: "* main",
			want:  []string{"main"},
		},
		{
			name: "multiple branches",
			input: `* main
  develop
  feature/new-feature`,
			want: []string{"main", "develop", "feature/new-feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBranches(tt.input)
			if err != nil {
				t.Errorf("ParseBranches() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty input",
			input: "",
			want:  []string{},
		},
		{
			name:  "single tag",
			input: "v1.0.0",
			want:  []string{"v1.0.0"},
		},
		{
			name: "multiple tags",
			input: `v1.0.0
v1.1.0
v2.0.0`,
			want: []string{"v1.0.0", "v1.1.0", "v2.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTags(tt.input)
			if err != nil {
				t.Errorf("ParseTags() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWeekday(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Weekday
		wantErr bool
	}{
		{
			name:    "Monday",
			input:   "2024-01-01", // Monday
			want:    time.Monday,
			wantErr: false,
		},
		{
			name:    "Friday",
			input:   "2024-01-05", // Friday
			want:    time.Friday,
			wantErr: false,
		},
		{
			name:    "invalid date",
			input:   "invalid",
			want:    time.Sunday,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWeekday(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWeekday() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseWeekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseHour(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "morning",
			input:   "2024-01-01 09:30:00",
			want:    9,
			wantErr: false,
		},
		{
			name:    "afternoon",
			input:   "2024-01-01 15:45:00",
			want:    15,
			wantErr: false,
		},
		{
			name:    "midnight",
			input:   "2024-01-01 00:00:00",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid datetime",
			input:   "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHour(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseHour() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkParseCommitLog(b *testing.B) {
	input := `abc123|John Doe|john@example.com|2024-01-01 10:00:00 +0000|Initial commit
def456|Jane Smith|jane@example.com|2024-01-02 11:00:00 +0000|Add feature
ghi789|Bob Johnson|bob@example.com|2024-01-03 12:00:00 +0000|Fix bug`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseCommitLog(input)
	}
}

func BenchmarkParseAuthorStats(b *testing.B) {
	input := `   100  John Doe
    50  Jane Smith
    25  Bob Johnson`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseAuthorStats(input)
	}
}

func BenchmarkParseDateCounts(b *testing.B) {
	input := `2024-01-01
2024-01-01
2024-01-02
2024-01-03
2024-01-03
2024-01-03`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseDateCounts(input)
	}
}
