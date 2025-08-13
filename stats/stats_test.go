package stats

import (
	"os"
	"testing"
)

func TestCalculateStats(t *testing.T) {
	testRepo := "../.."
	if _, err := os.Stat(testRepo); os.IsNotExist(err) {
		t.Skip("Test repo not found: ", testRepo)
	}

	stats, err := CalculateStats(testRepo)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}

	if len(stats.CommitsByUser) == 0 {
		t.Errorf("Expected commits by user, got 0")
	}

	if len(stats.FileModifications) == 0 {
		t.Errorf("Expected file modifications, got 0")
	}

	if stats.LinesAdded < 0 || stats.LinesDeleted < 0 {
		t.Errorf("Lines added/deleted should be >= 0")
	}

	if len(stats.CommitsByWeekday) == 0 {
		t.Errorf("Expected commits by weekday, got 0")
	}

	topFiles := stats.TopModifiedFiles(3)
	if len(topFiles) == 0 {
		t.Errorf("Expected top modified files, got 0")
	}
}
