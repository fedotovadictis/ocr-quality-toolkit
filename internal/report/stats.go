package report

import (
	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/evaluate"
	"ocr-quality-toolkit/internal/runner"
)

type Stats struct {
	Total               int
	Successful          int
	Failed              int
	TotalDurationMS     int64
	AverageDurationMS   int64
	ReferenceCharacters int
	CharacterErrors     int
	CER                 float64

	ReferenceWords int
	WordErrors     int
	WER            float64

	Coverage float64
}

func CalculateStats(results []runner.Result) Stats {
	stats := Stats{
		Total: len(results),
	}

	for _, result := range results {
		if result.Error != "" {
			stats.Failed++
		} else {
			stats.Successful++
		}

		stats.TotalDurationMS += result.DurationMS
	}

	if stats.Total > 0 {
		stats.AverageDurationMS =
			stats.TotalDurationMS / int64(stats.Total)
	}

	return stats
}

type EvaluationStats struct {
	Total               int     `json:"total"`
	Successful          int     `json:"successful"`
	EngineErrors        int     `json:"engine_errors"`
	Missing             int     `json:"missing"`
	Coverage            float64 `json:"coverage"`
	ReferenceCharacters int     `json:"reference_characters"`
	CharacterErrors     int     `json:"character_errors"`
	CER                 float64 `json:"cer"`
	ReferenceWords      int     `json:"reference_words"`
	WordErrors          int     `json:"word_errors"`
	WER                 float64 `json:"wer"`
}

func CalculateEvaluationStats(
	results []evaluate.Result,
) EvaluationStats {
	stats := EvaluationStats{
		Total: len(results),
	}

	for _, result := range results {
		switch result.Status {
		case evaluate.StatusSuccess:
			stats.Successful++

			stats.ReferenceCharacters +=
				result.ReferenceCharacters

			stats.CharacterErrors +=
				result.CharacterSubstitutions +
					result.CharacterDeletions +
					result.CharacterInsertions

			stats.ReferenceWords +=
				result.ReferenceWords

			stats.WordErrors +=
				result.WordSubstitutions +
					result.WordDeletions +
					result.WordInsertions

		case evaluate.StatusEngineError:
			stats.EngineErrors++

		case evaluate.StatusMissingHypothesis:
			stats.Missing++
		}
	}

	if stats.ReferenceCharacters > 0 {
		stats.CER =
			float64(stats.CharacterErrors) /
				float64(stats.ReferenceCharacters)
	}

	if stats.ReferenceWords > 0 {
		stats.WER =
			float64(stats.WordErrors) /
				float64(stats.ReferenceWords)
	}

	if stats.Total > 0 {
		stats.Coverage =
			float64(stats.Successful) /
				float64(stats.Total)
	}

	return stats
}

func GroupByLanguage(
	records []corpus.Record,
	results []evaluate.Result,
) map[string]EvaluationStats {
	return groupBy(
		records,
		results,
		func(record corpus.Record) []string {
			return []string{record.Language}
		},
	)
}

func GroupByTask(
	records []corpus.Record,
	results []evaluate.Result,
) map[string]EvaluationStats {
	return groupBy(
		records,
		results,
		func(record corpus.Record) []string {
			return []string{record.Task}
		},
	)
}

func GroupByTags(
	records []corpus.Record,
	results []evaluate.Result,
) map[string]EvaluationStats {
	return groupBy(
		records,
		results,
		func(record corpus.Record) []string {
			return record.Tags
		},
	)
}

func GroupByTransform(
	records []corpus.Record,
	results []evaluate.Result,
) map[string]EvaluationStats {
	return groupBy(
		records,
		results,
		func(record corpus.Record) []string {
			return []string{record.Transform.Name}
		},
	)
}

func groupBy(
	records []corpus.Record,
	results []evaluate.Result,
	keyFn func(corpus.Record) []string,
) map[string]EvaluationStats {
	resultByID := make(
		map[string]evaluate.Result,
		len(results),
	)

	for _, result := range results {
		resultByID[result.ID] = result
	}

	groupedResults := make(
		map[string][]evaluate.Result,
	)

	for _, record := range records {
		result, ok := resultByID[record.ID]
		if !ok {
			continue
		}

		for _, key := range keyFn(record) {
			if key == "" {
				continue
			}

			groupedResults[key] = append(
				groupedResults[key],
				result,
			)
		}
	}

	groups := make(
		map[string]EvaluationStats,
		len(groupedResults),
	)

	for key, groupResults := range groupedResults {
		groups[key] =
			CalculateEvaluationStats(groupResults)
	}

	return groups
}
