package git

import (
	"strings"
	"testing"
	"time"
)

func TestNewExecBackend(t *testing.T) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		t.Fatalf("NewExecBackend() error = %v", err)
	}

	if backend == nil {
		t.Fatal("NewExecBackend() returned nil")
	}

	if backend.RootPath() == "" {
		t.Error("RootPath() returned empty string")
	}
}

func TestExecBackendLog(t *testing.T) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		t.Skip("Git not available or not a repository")
	}

	output, err := backend.Log("--oneline", "-n", "5")
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}

	if output == "" {
		t.Error("Log() returned empty output")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 5 {
		t.Errorf("Expected max 5 lines, got %d", len(lines))
	}
}

func TestExecBackendLogPretty(t *testing.T) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		t.Skip("Git not available or not a repository")
	}

	format := "%H|%an|%ae"
	output, err := backend.LogPretty(format, "-n", "3")
	if err != nil {
		t.Fatalf("LogPretty() error = %v", err)
	}

	if output == "" {
		t.Error("LogPretty() returned empty output")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "|") {
			t.Errorf("Line doesn't contain separator: %s", line)
		}
	}
}

func TestExecBackendBranches(t *testing.T) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		t.Skip("Git not available or not a repository")
	}

	output, err := backend.Branches()
	if err != nil {
		t.Fatalf("Branches() error = %v", err)
	}

	if output == "" {
		t.Error("Branches() returned empty output")
	}
}

func TestExecBackendCurrentBranch(t *testing.T) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		t.Skip("Git not available or not a repository")
	}

	branch, err := backend.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch() error = %v", err)
	}

	if branch == "" {
		t.Error("CurrentBranch() returned empty string")
	}
}

func TestExecBackendRevList(t *testing.T) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		t.Skip("Git not available or not a repository")
	}

	output, err := backend.RevList("--count", "HEAD")
	if err != nil {
		t.Fatalf("RevList() error = %v", err)
	}

	if output == "" {
		t.Error("RevList() returned empty output")
	}
}

func TestExecBackendShortlog(t *testing.T) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		t.Skip("Git not available or not a repository")
	}

	output, err := backend.Shortlog("-s", "-n")
	if err != nil {
		t.Fatalf("Shortlog() error = %v", err)
	}

	if output == "" {
		t.Error("Shortlog() returned empty output")
	}
}

func TestBuildLogArgs(t *testing.T) {
	tests := []struct {
		name string
		opts *LogOptions
		want []string
	}{
		{
			name: "empty options",
			opts: &LogOptions{},
			want: []string{},
		},
		{
			name: "with since and until",
			opts: &LogOptions{
				Since: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Until: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			want: []string{"--since=2024-01-01", "--until=2024-12-31"},
		},
		{
			name: "with author",
			opts: &LogOptions{
				Author: "John Doe",
			},
			want: []string{"--author=John Doe"},
		},
		{
			name: "with no merges",
			opts: &LogOptions{
				NoMerges: true,
			},
			want: []string{"--no-merges"},
		},
		{
			name: "with merges only",
			opts: &LogOptions{
				MergesOnly: true,
			},
			want: []string{"--merges"},
		},
		{
			name: "with limit",
			opts: &LogOptions{
				Limit: 10,
			},
			want: []string{"--max-count=10"},
		},
		{
			name: "with branch",
			opts: &LogOptions{
				Branch: "main",
			},
			want: []string{"main"},
		},
		{
			name: "with pathspec",
			opts: &LogOptions{
				PathSpec: []string{":!vendor", ":!node_modules"},
			},
			want: []string{"--", ":!vendor", ":!node_modules"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildLogArgs(tt.opts)

			// Check that all expected args are present
			for _, wantArg := range tt.want {
				found := false
				for _, gotArg := range got {
					if gotArg == wantArg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("BuildLogArgs() missing argument %q, got %v", wantArg, got)
				}
			}
		})
	}
}

func TestBuildLogArgsWithFormat(t *testing.T) {
	opts := &LogOptions{
		Format: "%H|%an|%ae",
	}

	args := BuildLogArgs(opts)

	found := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "--pretty=format:") {
			found = true
			if !strings.Contains(arg, "%H|%an|%ae") {
				t.Error("Format not properly included in args")
			}
		}
	}

	if !found {
		t.Error("Format argument not found in args")
	}
}

func TestBuildLogArgsWithExtraArgs(t *testing.T) {
	opts := &LogOptions{
		ExtraArgs: []string{"--all", "--graph"},
	}

	args := BuildLogArgs(opts)

	hasAll := false
	hasGraph := false
	for _, arg := range args {
		if arg == "--all" {
			hasAll = true
		}
		if arg == "--graph" {
			hasGraph = true
		}
	}

	if !hasAll || !hasGraph {
		t.Error("Extra args not included properly")
	}
}

// Benchmark tests
func BenchmarkExecBackendLog(b *testing.B) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		b.Skip("Git not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = backend.Log("--oneline", "-n", "10")
	}
}

func BenchmarkExecBackendShortlog(b *testing.B) {
	backend, err := NewExecBackend("../..")
	if err != nil {
		b.Skip("Git not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = backend.Shortlog("-s", "-n")
	}
}

func BenchmarkBuildLogArgs(b *testing.B) {
	opts := &LogOptions{
		Since:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:     time.Now(),
		Author:    "test",
		Format:    "%H|%an",
		Branch:    "main",
		NoMerges:  true,
		Limit:     100,
		PathSpec:  []string{":!vendor"},
		ExtraArgs: []string{"--all"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildLogArgs(opts)
	}
}
