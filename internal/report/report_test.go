package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/evaluate"
)

func TestBuildReport(t *testing.T) {
	records := []corpus.Record{
		{
			ID:       "1",
			Language: "ru",
			Task:     "ocr",
			Tags:     []string{"synthetic"},
			Transform: corpus.Transform{
				Name: "grayscale",
			},
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
	}

	report := BuildReport(records, results)

	if report.Version == "" {
		t.Fatal("expected non-empty report version")
	}

	if report.Overall.Total != 1 {
		t.Fatalf(
			"expected total 1, got %d",
			report.Overall.Total,
		)
	}

	if report.Overall.CER != 0.1 {
		t.Fatalf(
			"expected CER 0.1, got %v",
			report.Overall.CER,
		)
	}

	if report.ByLanguage["ru"].Total != 1 {
		t.Fatal("expected ru language group")
	}

	if report.ByTask["ocr"].Total != 1 {
		t.Fatal("expected ocr task group")
	}

	if report.ByTags["synthetic"].Total != 1 {
		t.Fatal("expected synthetic tag group")
	}

	if report.ByTransform["grayscale"].Total != 1 {
		t.Fatal("expected grayscale transform group")
	}

	if len(report.Results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(report.Results),
		)
	}
}
func TestWriteJSON(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"report.json",
	)

	report := Report{
		Version: "1",
		Overall: EvaluationStats{
			Total:      2,
			Successful: 1,
			Missing:    1,
			Coverage:   0.5,
			CER:        0.1,
			WER:        0.2,
		},
		ByLanguage: map[string]EvaluationStats{
			"ru": {
				Total: 2,
			},
		},
		ByTask:      map[string]EvaluationStats{},
		ByTags:      map[string]EvaluationStats{},
		ByTransform: map[string]EvaluationStats{},
		Results: []evaluate.Result{
			{
				ID:     "page-001",
				Status: evaluate.StatusSuccess,
			},
			{
				ID:     "page-002",
				Status: evaluate.StatusMissingHypothesis,
			},
		},
	}

	if err := WriteJSON(path, report); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}

	var decoded Report

	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode report: %v", err)
	}

	if decoded.Version != "1" {
		t.Fatalf(
			"expected version %q, got %q",
			"1",
			decoded.Version,
		)
	}

	if decoded.Overall.Total != 2 {
		t.Fatalf(
			"expected total 2, got %d",
			decoded.Overall.Total,
		)
	}

	if decoded.Overall.Coverage != 0.5 {
		t.Fatalf(
			"expected coverage 0.5, got %v",
			decoded.Overall.Coverage,
		)
	}

	if decoded.ByLanguage["ru"].Total != 2 {
		t.Fatal("expected ru language group")
	}

	if len(decoded.Results) != 2 {
		t.Fatalf(
			"expected 2 results, got %d",
			len(decoded.Results),
		)
	}
}
func TestBuildReportEmpty(t *testing.T) {
	report := BuildReport(nil, nil)

	if report.Version == "" {
		t.Fatal("expected non-empty version")
	}

	if report.Overall.Total != 0 {
		t.Fatalf(
			"expected total 0, got %d",
			report.Overall.Total,
		)
	}

	if report.Overall.Coverage != 0 {
		t.Fatalf(
			"expected coverage 0, got %v",
			report.Overall.Coverage,
		)
	}

	if len(report.Results) != 0 {
		t.Fatalf(
			"expected 0 results, got %d",
			len(report.Results),
		)
	}
}
func TestBuildReportPartialResults(t *testing.T) {
	records := []corpus.Record{
		{ID: "1", Language: "ru"},
		{ID: "2", Language: "ru"},
		{ID: "3", Language: "ru"},
	}

	results := []evaluate.Result{
		{
			ID:                  "1",
			Status:              evaluate.StatusSuccess,
			ReferenceCharacters: 10,
			ReferenceWords:      2,
		},
		{
			ID:     "2",
			Status: evaluate.StatusOCRError,
		},
		{
			ID:     "3",
			Status: evaluate.StatusMissingHypothesis,
		},
	}

	report := BuildReport(records, results)

	if report.Overall.Total != 3 {
		t.Fatalf(
			"expected total 3, got %d",
			report.Overall.Total,
		)
	}

	if report.Overall.Successful != 2 {
		t.Fatalf(
			"expected 2 covered results, got %d",
			report.Overall.Successful,
		)
	}

	if report.Overall.Missing != 1 {
		t.Fatalf(
			"expected 1 missing result, got %d",
			report.Overall.Missing,
		)
	}

	expectedCoverage := float64(2) / float64(3)

	if report.Overall.Coverage != expectedCoverage {
		t.Fatalf(
			"expected coverage %v, got %v",
			expectedCoverage,
			report.Overall.Coverage,
		)
	}

	if len(report.Results) != 3 {
		t.Fatalf(
			"expected 3 results, got %d",
			len(report.Results),
		)
	}
}
