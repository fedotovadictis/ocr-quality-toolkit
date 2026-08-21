package evaluate

import (
	"testing"

	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/metrics"
	"ocr-quality-toolkit/internal/normalize"
)

func TestEvaluateSingleSuccess(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID:         "1",
			References: []string{"кот"},
		},
	}

	hypotheses := []corpus.Hypothesis{
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
			ID:         "1",
			Image:      "images/page.png",
			References: []string{"кот"},
		},
	}

	results, err := Evaluate(
		manifest,
		nil,
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

	if result.Status != StatusMissingHypothesis {
		t.Fatalf(
			"expected status %q, got %q",
			StatusMissingHypothesis,
			result.Status,
		)
	}

	if result.Image != "images/page.png" {
		t.Fatalf(
			"unexpected image: %q",
			result.Image,
		)
	}

	if result.SelectedReference != -1 {
		t.Fatalf(
			"expected selected reference -1, got %d",
			result.SelectedReference,
		)
	}
}

func TestEvaluateRecognizedTextWithErrors(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID:         "1",
			References: []string{"кот"},
		},
	}

	hypotheses := []corpus.Hypothesis{
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

	if result.Status != StatusSuccess {
		t.Fatalf(
			"expected status %q, got %q",
			StatusSuccess,
			result.Status,
		)
	}

	if result.ExactMatch {
		t.Fatal("expected exact match to be false")
	}

	if result.CER == 0 {
		t.Fatal("expected CER to be greater than zero")
	}

	if len(result.Alignment) == 0 {
		t.Fatal("expected non-empty alignment")
	}
}

func TestEvaluateEmptyHypothesisText(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID:         "1",
			References: []string{"кот"},
		},
	}

	hypotheses := []corpus.Hypothesis{
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

	if result.Status != StatusSuccess {
		t.Fatalf(
			"expected status %q, got %q",
			StatusSuccess,
			result.Status,
		)
	}

	if result.CER != 1 {
		t.Fatalf(
			"expected CER=1, got %v",
			result.CER,
		)
	}

	if result.ExactMatch {
		t.Fatal("expected exact match to be false")
	}
}

func TestEvaluateUnknownHypothesisID(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID:         "1",
			References: []string{"кот"},
		},
	}

	hypotheses := []corpus.Hypothesis{
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
			ID:         "1",
			References: []string{"кот"},
		},
	}

	hypotheses := []corpus.Hypothesis{
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

	hypotheses := []corpus.Hypothesis{
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

	if result.SelectedReference != 1 {
		t.Fatalf(
			"expected selected reference 1, got %d",
			result.SelectedReference,
		)
	}

	if result.Reference != "кот" {
		t.Fatalf(
			"unexpected selected raw reference: %q",
			result.Reference,
		)
	}
}

func TestEvaluatePreservesRawAndNormalizedText(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID: "1",
			References: []string{
				"  КОТ  ",
				"Собака",
			},
		},
	}

	hypotheses := []corpus.Hypothesis{
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

	if result.Reference != "  КОТ  " {
		t.Fatalf(
			"unexpected raw reference: %q",
			result.Reference,
		)
	}

	if result.Hypothesis != "кот" {
		t.Fatalf(
			"unexpected raw hypothesis: %q",
			result.Hypothesis,
		)
	}

	if result.NormalizedReference != "кот" {
		t.Fatalf(
			"unexpected normalized reference: %q",
			result.NormalizedReference,
		)
	}

	if result.NormalizedHypothesis != "кот" {
		t.Fatalf(
			"unexpected normalized hypothesis: %q",
			result.NormalizedHypothesis,
		)
	}

	if result.SelectedReference != 0 {
		t.Fatalf(
			"expected selected reference 0, got %d",
			result.SelectedReference,
		)
	}

	if result.CER != 0 {
		t.Fatalf(
			"expected CER=0, got %v",
			result.CER,
		)
	}
}

func TestEvaluateStoresCharacterAlignment(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID:         "1",
			References: []string{"кот"},
		},
	}

	hypotheses := []corpus.Hypothesis{
		{
			ID:   "1",
			Text: "кит",
		},
	}

	results, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfileStrict,
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

	alignment := results[0].Alignment

	if len(alignment) != 3 {
		t.Fatalf(
			"expected 3 alignment operations, got %d: %#v",
			len(alignment),
			alignment,
		)
	}

	if alignment[0].Type != metrics.OperationEqual {
		t.Fatalf(
			"unexpected first operation: %#v",
			alignment[0],
		)
	}

	if alignment[1].Type != metrics.OperationSubstitute {
		t.Fatalf(
			"expected substitution, got %#v",
			alignment[1],
		)
	}

	if alignment[1].Reference != "о" {
		t.Fatalf(
			"unexpected substitution reference: %q",
			alignment[1].Reference,
		)
	}

	if alignment[1].Hypothesis != "и" {
		t.Fatalf(
			"unexpected substitution hypothesis: %q",
			alignment[1].Hypothesis,
		)
	}

	if alignment[2].Type != metrics.OperationEqual {
		t.Fatalf(
			"unexpected last operation: %#v",
			alignment[2],
		)
	}
}

func TestEvaluateEngineError(t *testing.T) {
	manifest := []corpus.Record{
		{
			ID:         "1",
			Image:      "images/page.png",
			References: []string{"кот"},
		},
	}

	hypotheses := []corpus.Hypothesis{
		{
			ID:    "1",
			Error: "process timeout",
		},
	}

	results, err := Evaluate(
		manifest,
		hypotheses,
		normalize.ProfileStrict,
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

	if result.Status != StatusEngineError {
		t.Fatalf(
			"expected status %q, got %q",
			StatusEngineError,
			result.Status,
		)
	}

	if result.Error != "process timeout" {
		t.Fatalf(
			"unexpected engine error: %q",
			result.Error,
		)
	}

	if result.Image != "images/page.png" {
		t.Fatalf(
			"unexpected image: %q",
			result.Image,
		)
	}

	if result.ReferenceCharacters != 0 {
		t.Fatalf(
			"engine error must not contribute reference characters, got %d",
			result.ReferenceCharacters,
		)
	}

	if result.ReferenceWords != 0 {
		t.Fatalf(
			"engine error must not contribute reference words, got %d",
			result.ReferenceWords,
		)
	}

	if len(result.Alignment) != 0 {
		t.Fatalf(
			"engine error must not have alignment: %#v",
			result.Alignment,
		)
	}

	if result.SelectedReference != -1 {
		t.Fatalf(
			"expected selected reference -1, got %d",
			result.SelectedReference,
		)
	}
}
