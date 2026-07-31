package evaluate

import (
	"testing"

	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/normalize"
)

func TestEvaluateSingleSuccess(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"кот",
			},
		},
	}

	hypotheses := []Hypothesis{
		{
			ID:   "1",
			Text: "кот",
		},
	}

	results, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfilePlainTextRU,
	)

	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}

	result := results[0]

	if result.ID != "1" {
		t.Fatalf(
			"unexpected id: %s",
			result.ID,
		)
	}

	if result.CER != 0 {
		t.Fatalf(
			"expected CER=0, got %v",
			result.CER,
		)
	}

	if result.WER != 0 {
		t.Fatalf(
			"expected WER=0, got %v",
			result.WER,
		)
	}

	if result.Similarity != 1 {
		t.Fatalf(
			"expected similarity=1, got %v",
			result.Similarity,
		)
	}

	if !result.ExactMatch {
		t.Fatal("expected exact match")
	}

	if result.Status != StatusSuccess {
		t.Fatalf(
			"expected status %q, got %q",
			StatusSuccess,
			result.Status,
		)
	}
}
func TestEvaluateMissingHypothesis(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"кот",
			},
		},
	}

	hypotheses := []Hypothesis{}

	results, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfilePlainTextRU,
	)

	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}

	if results[0].Status != StatusMissingHypothesis {
		t.Fatalf(
			"expected status %q, got %q",
			StatusMissingHypothesis,
			results[0].Status,
		)
	}
}

func TestEvaluateOCRError(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"кот",
			},
		},
	}

	hypotheses := []Hypothesis{
		{
			ID:   "1",
			Text: "кит",
		},
	}

	results, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfilePlainTextRU,
	)

	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	result := results[0]

	if result.Status != StatusOCRError {
		t.Fatalf(
			"expected status %q, got %q",
			StatusOCRError,
			result.Status,
		)
	}

	if result.ExactMatch {
		t.Fatal("expected exact match to be false")
	}

	if result.CER == 0 {
		t.Fatal("expected CER to be greater than zero")
	}
}

func TestEvaluateEmptyHypothesisText(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"кот",
			},
		},
	}

	hypotheses := []Hypothesis{
		{
			ID:   "1",
			Text: "",
		},
	}

	results, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfilePlainTextRU,
	)

	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	result := results[0]

	if result.Status != StatusOCRError {
		t.Fatalf(
			"expected status %q, got %q",
			StatusOCRError,
			result.Status,
		)
	}
}
func TestEvaluateUnknownHypothesisID(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"кот",
			},
		},
	}

	hypotheses := []Hypothesis{
		{
			ID:   "999",
			Text: "кот",
		},
	}

	_, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfilePlainTextRU,
	)

	if err == nil {
		t.Fatal("expected error for unknown hypothesis ID")
	}
}
func TestEvaluateDuplicateHypothesisID(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"кот",
			},
		},
	}

	hypotheses := []Hypothesis{
		{
			ID:   "1",
			Text: "кот",
		},
		{
			ID:   "1",
			Text: "кит",
		},
	}

	_, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfilePlainTextRU,
	)

	if err == nil {
		t.Fatal("expected error for duplicate hypothesis ID")
	}
}
func TestEvaluateChoosesBestReference(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"кит",
				"кот",
			},
		},
	}

	hypotheses := []Hypothesis{
		{
			ID:   "1",
			Text: "кот",
		},
	}

	results, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfilePlainTextRU,
	)

	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	result := results[0]

	if result.CER != 0 {
		t.Fatalf(
			"expected best reference with CER=0, got %v",
			result.CER,
		)
	}

	if !result.ExactMatch {
		t.Fatal("expected exact match with best reference")
	}
}
