package report

import (
	"testing"

	"ocr-quality-toolkit/internal/runner"
)

func TestCalculateStats(t *testing.T) {
	results := []runner.Result{
		{
			ID:         "page-001",
			Text:       "text",
			DurationMS: 100,
		},
		{
			ID:         "page-002",
			Text:       "text",
			DurationMS: 200,
		},
		{
			ID:         "page-003",
			Error:      "ocr failed",
			DurationMS: 300,
		},
	}

	stats := CalculateStats(results)

	if stats.Total != 3 {
		t.Fatalf("expected total 3, got %d", stats.Total)
	}

	if stats.Successful != 2 {
		t.Fatalf(
			"expected 2 successful, got %d",
			stats.Successful,
		)
	}

	if stats.Failed != 1 {
		t.Fatalf(
			"expected 1 failed, got %d",
			stats.Failed,
		)
	}

	if stats.TotalDurationMS != 600 {
		t.Fatalf(
			"expected total duration 600, got %d",
			stats.TotalDurationMS,
		)
	}

	if stats.AverageDurationMS != 200 {
		t.Fatalf(
			"expected average duration 200, got %d",
			stats.AverageDurationMS,
		)
	}
}
func TestCalculateStatsEmpty(t *testing.T) {
	stats := CalculateStats(nil)

	if stats.Total != 0 {
		t.Fatalf("expected total 0, got %d", stats.Total)
	}

	if stats.Successful != 0 {
		t.Fatalf(
			"expected successful 0, got %d",
			stats.Successful,
		)
	}

	if stats.Failed != 0 {
		t.Fatalf(
			"expected failed 0, got %d",
			stats.Failed,
		)
	}

	if stats.TotalDurationMS != 0 {
		t.Fatalf(
			"expected total duration 0, got %d",
			stats.TotalDurationMS,
		)
	}

	if stats.AverageDurationMS != 0 {
		t.Fatalf(
			"expected average duration 0, got %d",
			stats.AverageDurationMS,
		)
	}
}
