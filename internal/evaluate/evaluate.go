package evaluate

import (
	"fmt"
	"strings"

	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/metrics"
	"ocr-quality-toolkit/internal/normalize"
)

type Result struct {
	ID         string  `json:"id"`
	CER        float64 `json:"cer"`
	WER        float64 `json:"wer"`
	Similarity float64 `json:"similarity"`
	ExactMatch bool    `json:"exact_match"`
	Status     Status  `json:"status"`

	ReferenceCharacters    int `json:"reference_characters"`
	CharacterHits          int `json:"character_hits"`
	CharacterSubstitutions int `json:"character_substitutions"`
	CharacterDeletions     int `json:"character_deletions"`
	CharacterInsertions    int `json:"character_insertions"`

	ReferenceWords    int `json:"reference_words"`
	WordHits          int `json:"word_hits"`
	WordSubstitutions int `json:"word_substitutions"`
	WordDeletions     int `json:"word_deletions"`
	WordInsertions    int `json:"word_insertions"`

	Reference  string `json:"reference"`
	Hypothesis string `json:"hypothesis"`
	Image      string `json:"image"`
}

func Evaluate(
	records []corpus.Record,
	hypotheses []corpus.Hypothesis,
	profile normalize.Profile,
) ([]Result, error) {
	manifestIDs := make(map[string]bool, len(records))

	for _, record := range records {
		manifestIDs[record.ID] = true
	}

	hypothesisIDs := make(map[string]bool, len(hypotheses))

	for _, hypothesis := range hypotheses {
		if hypothesisIDs[hypothesis.ID] {
			return nil, fmt.Errorf(
				"duplicate hypothesis ID: %s",
				hypothesis.ID,
			)
		}

		if !manifestIDs[hypothesis.ID] {
			return nil, fmt.Errorf(
				"unknown hypothesis ID: %s",
				hypothesis.ID,
			)
		}

		hypothesisIDs[hypothesis.ID] = true
	}

	results := make([]Result, 0, len(records))

	for _, record := range records {
		var hypothesisText string
		hypothesisFound := false

		result := Result{
			ID: record.ID,
		}

		for _, hypothesis := range hypotheses {
			if hypothesis.ID == record.ID {
				hypothesisText = hypothesis.Text
				hypothesisFound = true
				break
			}
		}

		if !hypothesisFound {
			result.Status = StatusMissingHypothesis
			results = append(results, result)
			continue
		}

		reference, err := chooseBestReference(
			record.References,
			hypothesisText,
			profile,
		)
		result.Reference = reference
		result.Hypothesis = hypothesisText
		result.Image = record.Image

		if err != nil {
			return nil, fmt.Errorf(
				"evaluate record %q: %w",
				record.ID,
				err,
			)
		}

		hypothesis, err := normalize.Normalize(
			hypothesisText,
			profile,
		)
		if err != nil {
			return nil, err
		}

		result.CER = metrics.CER(reference, hypothesis)
		result.WER = metrics.WER(reference, hypothesis)
		referenceSymbols := stringSymbols(reference)
		hypothesisSymbols := stringSymbols(hypothesis)

		characterAlignment := metrics.Align(
			referenceSymbols,
			hypothesisSymbols,
		)

		result.ReferenceCharacters = len(referenceSymbols)
		result.CharacterHits = characterAlignment.Hits
		result.CharacterSubstitutions = characterAlignment.Substitutions
		result.CharacterDeletions = characterAlignment.Deletions
		result.CharacterInsertions = characterAlignment.Insertions

		referenceWords := strings.Fields(reference)
		hypothesisWords := strings.Fields(hypothesis)

		wordAlignment := metrics.Align(
			referenceWords,
			hypothesisWords,
		)

		result.ReferenceWords = len(referenceWords)
		result.WordHits = wordAlignment.Hits
		result.WordSubstitutions = wordAlignment.Substitutions
		result.WordDeletions = wordAlignment.Deletions
		result.WordInsertions = wordAlignment.Insertions

		result.Similarity = metrics.Similarity(
			referenceSymbols,
			hypothesisSymbols,
		)

		result.ExactMatch = reference == hypothesis

		if result.ExactMatch {
			result.Status = StatusSuccess
		} else {
			result.Status = StatusOCRError
		}

		results = append(results, result)
	}

	return results, nil
}

func chooseBestReference(
	references []string,
	hypothesis string,
	profile normalize.Profile,
) (string, error) {
	if len(references) == 0 {
		return "", fmt.Errorf("no references provided")
	}

	normalizedHypothesis, err := normalize.Normalize(
		hypothesis,
		profile,
	)
	if err != nil {
		return "", err
	}

	var bestReference string
	bestCER := -1.0

	for _, reference := range references {
		normalizedReference, err := normalize.Normalize(
			reference,
			profile,
		)
		if err != nil {
			return "", err
		}

		currentCER := metrics.CER(
			normalizedReference,
			normalizedHypothesis,
		)

		if bestCER < 0 || currentCER < bestCER {
			bestReference = normalizedReference
			bestCER = currentCER
		}
	}

	return bestReference, nil
}

func stringSymbols(text string) []string {
	symbols := make([]string, 0, len([]rune(text)))

	for _, r := range text {
		symbols = append(symbols, string(r))
	}

	return symbols
}
