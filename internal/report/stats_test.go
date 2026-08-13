package report

import (
	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/evaluate"
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
func TestCalculateEvaluationStats(t *testing.T) {
	results := []evaluate.Result{
		{
			ID:                     "1",
			Status:                 evaluate.StatusSuccess,
			ReferenceCharacters:    10,
			CharacterSubstitutions: 1,
			CharacterDeletions:     1,
			CharacterInsertions:    0,
			ReferenceWords:         4,
			WordSubstitutions:      1,
			WordDeletions:          0,
			WordInsertions:         0,
		},
		{
			ID:                     "2",
			Status:                 evaluate.StatusOCRError,
			ReferenceCharacters:    20,
			CharacterSubstitutions: 2,
			CharacterDeletions:     0,
			CharacterInsertions:    1,
			ReferenceWords:         6,
			WordSubstitutions:      1,
			WordDeletions:          1,
			WordInsertions:         0,
		},
		{
			ID:     "3",
			Status: evaluate.StatusMissingHypothesis,
		},
	}

	stats := CalculateEvaluationStats(results)

	if stats.Total != 3 {
		t.Fatalf("expected total 3, got %d", stats.Total)
	}

	if stats.Successful != 2 {
		t.Fatalf(
			"expected 2 successful, got %d",
			stats.Successful,
		)
	}

	if stats.Missing != 1 {
		t.Fatalf(
			"expected 1 missing, got %d",
			stats.Missing,
		)
	}

	expectedCER := float64(5) / float64(30)
	if stats.CER != expectedCER {
		t.Fatalf(
			"expected CER %v, got %v",
			expectedCER,
			stats.CER,
		)
	}

	expectedWER := float64(3) / float64(10)
	if stats.WER != expectedWER {
		t.Fatalf(
			"expected WER %v, got %v",
			expectedWER,
			stats.WER,
		)
	}

	expectedCoverage := float64(2) / float64(3)
	if stats.Coverage != expectedCoverage {
		t.Fatalf(
			"expected coverage %v, got %v",
			expectedCoverage,
			stats.Coverage,
		)
	}
}
func TestGroupStatsByLanguage(t *testing.T) {
	records := []corpus.Record{
		{
			ID:       "1",
			Language: "ru",
		},
		{
			ID:       "2",
			Language: "ru",
		},
		{
			ID:       "3",
			Language: "en",
		},
	}

	results := []evaluate.Result{
		{
			ID:                     "1",
			Status:                 evaluate.StatusSuccess,
			ReferenceCharacters:    10,
			CharacterSubstitutions: 1,
			ReferenceWords:         2,
			WordSubstitutions:      1,
		},
		{
			ID:                  "2",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
			CharacterDeletions:  1,
			ReferenceWords:      2,
		},
		{
			ID:                  "3",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 20,
			CharacterInsertions: 2,
			ReferenceWords:      4,
			WordInsertions:      1,
		},
	}

	groups := GroupByLanguage(records, results)

	ru, ok := groups["ru"]
	if !ok {
		t.Fatal("expected ru group")
	}

	if ru.Total != 2 {
		t.Fatalf("expected ru total 2, got %d", ru.Total)
	}

	if ru.CER != 0.1 {
		t.Fatalf("expected ru CER 0.1, got %v", ru.CER)
	}

	en, ok := groups["en"]
	if !ok {
		t.Fatal("expected en group")
	}

	if en.Total != 1 {
		t.Fatalf("expected en total 1, got %d", en.Total)
	}

	if en.CER != 0.1 {
		t.Fatalf("expected en CER 0.1, got %v", en.CER)
	}
}
func TestGroupStatsByTask(t *testing.T) {
	records := []corpus.Record{
		{ID: "1", Task: "printed"},
		{ID: "2", Task: "printed"},
		{ID: "3", Task: "handwritten"},
	}

	results := []evaluate.Result{
		{
			ID:                  "1",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
		},
		{
			ID:                     "2",
			Status:                 evaluate.StatusSuccess,
			ReferenceCharacters:    10,
			CharacterSubstitutions: 2,
		},
		{
			ID:                  "3",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
			CharacterDeletions:  1,
		},
	}

	groups := GroupByTask(records, results)

	printed, ok := groups["printed"]
	if !ok {
		t.Fatal("expected printed group")
	}

	if printed.Total != 2 {
		t.Fatalf("expected printed total 2, got %d", printed.Total)
	}

	if printed.CER != 0.1 {
		t.Fatalf("expected printed CER 0.1, got %v", printed.CER)
	}

	handwritten, ok := groups["handwritten"]
	if !ok {
		t.Fatal("expected handwritten group")
	}

	if handwritten.Total != 1 {
		t.Fatalf(
			"expected handwritten total 1, got %d",
			handwritten.Total,
		)
	}

	if handwritten.CER != 0.1 {
		t.Fatalf(
			"expected handwritten CER 0.1, got %v",
			handwritten.CER,
		)
	}
}

func TestGroupStatsByTags(t *testing.T) {
	records := []corpus.Record{
		{
			ID:   "1",
			Tags: []string{"scan", "low-quality"},
		},
		{
			ID:   "2",
			Tags: []string{"scan"},
		},
		{
			ID:   "3",
			Tags: []string{"synthetic"},
		},
	}

	results := []evaluate.Result{
		{
			ID:                  "1",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
			CharacterDeletions:  1,
		},
		{
			ID:                  "2",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
		},
		{
			ID:                     "3",
			Status:                 evaluate.StatusSuccess,
			ReferenceCharacters:    10,
			CharacterSubstitutions: 2,
		},
	}

	groups := GroupByTags(records, results)

	scan, ok := groups["scan"]
	if !ok {
		t.Fatal("expected scan group")
	}

	if scan.Total != 2 {
		t.Fatalf("expected scan total 2, got %d", scan.Total)
	}

	if scan.CER != 0.05 {
		t.Fatalf("expected scan CER 0.05, got %v", scan.CER)
	}

	lowQuality, ok := groups["low-quality"]
	if !ok {
		t.Fatal("expected low-quality group")
	}

	if lowQuality.Total != 1 {
		t.Fatalf(
			"expected low-quality total 1, got %d",
			lowQuality.Total,
		)
	}

	synthetic, ok := groups["synthetic"]
	if !ok {
		t.Fatal("expected synthetic group")
	}

	if synthetic.CER != 0.2 {
		t.Fatalf(
			"expected synthetic CER 0.2, got %v",
			synthetic.CER,
		)
	}
}
func TestGroupStatsByTransform(t *testing.T) {
	records := []corpus.Record{
		{
			ID: "1",
			Transform: corpus.Transform{
				Name: "grayscale",
			},
		},
		{
			ID: "2",
			Transform: corpus.Transform{
				Name: "grayscale",
			},
		},
		{
			ID: "3",
			Transform: corpus.Transform{
				Name: "jpeg-70",
			},
		},
	}

	results := []evaluate.Result{
		{
			ID:                  "1",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
			CharacterDeletions:  1,
		},
		{
			ID:                  "2",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
		},
		{
			ID:                     "3",
			Status:                 evaluate.StatusSuccess,
			ReferenceCharacters:    10,
			CharacterSubstitutions: 2,
		},
	}

	groups := GroupByTransform(records, results)

	grayscale, ok := groups["grayscale"]
	if !ok {
		t.Fatal("expected grayscale group")
	}

	if grayscale.Total != 2 {
		t.Fatalf(
			"expected grayscale total 2, got %d",
			grayscale.Total,
		)
	}

	if grayscale.CER != 0.05 {
		t.Fatalf(
			"expected grayscale CER 0.05, got %v",
			grayscale.CER,
		)
	}

	jpeg, ok := groups["jpeg-70"]
	if !ok {
		t.Fatal("expected jpeg-70 group")
	}

	if jpeg.Total != 1 {
		t.Fatalf(
			"expected jpeg-70 total 1, got %d",
			jpeg.Total,
		)
	}

	if jpeg.CER != 0.2 {
		t.Fatalf(
			"expected jpeg-70 CER 0.2, got %v",
			jpeg.CER,
		)
	}
}
