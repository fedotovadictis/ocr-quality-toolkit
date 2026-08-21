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
	Error      string  `json:"error,omitempty"`

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

	Reference            string              `json:"reference"`
	Hypothesis           string              `json:"hypothesis"`
	NormalizedReference  string              `json:"normalized_reference"`
	NormalizedHypothesis string              `json:"normalized_hypothesis"`
	SelectedReference    int                 `json:"selected_reference"`
	Image                string              `json:"image"`
	Alignment            []metrics.Operation `json:"alignment"`
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
		var matchedHypothesis corpus.Hypothesis
		hypothesisFound := false

		result := Result{
			ID:                record.ID,
			Image:             record.Image,
			SelectedReference: -1,
		}

		for _, hypothesis := range hypotheses {
			if hypothesis.ID == record.ID {
				matchedHypothesis = hypothesis
				hypothesisFound = true
				break
			}
		}

		if !hypothesisFound {
			result.Status = StatusMissingHypothesis
			results = append(results, result)
			continue
		}

		if matchedHypothesis.Error != "" {
			result.Status = StatusEngineError
			result.Error = matchedHypothesis.Error
			result.Hypothesis = matchedHypothesis.Text

			results = append(results, result)
			continue
		}

		hypothesisText := matchedHypothesis.Text

		rawReference, normalizedReference, selectedIndex, err :=
			chooseBestReference(
				record.References,
				hypothesisText,
				profile,
			)
		if err != nil {
			return nil, fmt.Errorf(
				"evaluate record %q: %w",
				record.ID,
				err,
			)
		}

		result.Reference = rawReference
		result.Hypothesis = hypothesisText
		result.NormalizedReference = normalizedReference
		result.SelectedReference = selectedIndex

		normalizedHypothesis, err := normalize.Normalize(
			hypothesisText,
			profile,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"normalize hypothesis for record %q: %w",
				record.ID,
				err,
			)
		}

		result.NormalizedHypothesis = normalizedHypothesis

		result.CER = metrics.CER(
			normalizedReference,
			normalizedHypothesis,
		)

		result.WER = metrics.WER(
			normalizedReference,
			normalizedHypothesis,
		)

		referenceSymbols := stringSymbols(normalizedReference)
		hypothesisSymbols := stringSymbols(normalizedHypothesis)

		characterAlignment := metrics.Align(
			referenceSymbols,
			hypothesisSymbols,
		)

		result.Alignment = append(
			[]metrics.Operation(nil),
			characterAlignment.Operations...,
		)

		result.ReferenceCharacters = len(referenceSymbols)
		result.CharacterHits = characterAlignment.Hits
		result.CharacterSubstitutions = characterAlignment.Substitutions
		result.CharacterDeletions = characterAlignment.Deletions
		result.CharacterInsertions = characterAlignment.Insertions

		referenceWords := strings.Fields(normalizedReference)
		hypothesisWords := strings.Fields(normalizedHypothesis)

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

		result.ExactMatch =
			normalizedReference == normalizedHypothesis

		result.Status = StatusSuccess

		results = append(results, result)
	}

	return results, nil
}

func chooseBestReference(
	references []string,
	hypothesis string,
	profile normalize.Profile,
) (
	rawReference string,
	normalizedReference string,
	selectedIndex int,
	err error,
) {
	if len(references) == 0 {
		return "", "", -1, fmt.Errorf("no references provided")
	}

	normalizedHypothesis, err := normalize.Normalize(
		hypothesis,
		profile,
	)
	if err != nil {
		return "", "", -1, err
	}

	bestCER := -1.0
	bestIndex := -1
	bestRaw := ""
	bestNormalized := ""

	for index, reference := range references {
		normalized, err := normalize.Normalize(
			reference,
			profile,
		)
		if err != nil {
			return "", "", -1, err
		}

		currentCER := metrics.CER(
			normalized,
			normalizedHypothesis,
		)

		if bestCER < 0 || currentCER < bestCER {
			bestCER = currentCER
			bestIndex = index
			bestRaw = reference
			bestNormalized = normalized
		}
	}

	return bestRaw, bestNormalized, bestIndex, nil
}

func stringSymbols(text string) []string {
	symbols := make([]string, 0, len([]rune(text)))

	for _, r := range text {
		symbols = append(symbols, string(r))
	}

	return symbols
}
