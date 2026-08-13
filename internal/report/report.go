package report

import (
	"encoding/json"
	"fmt"
	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/evaluate"
	"os"
)

type Report struct {
	Version     string                     `json:"version"`
	Overall     EvaluationStats            `json:"overall"`
	ByLanguage  map[string]EvaluationStats `json:"by_language"`
	ByTask      map[string]EvaluationStats `json:"by_task"`
	ByTags      map[string]EvaluationStats `json:"by_tags"`
	ByTransform map[string]EvaluationStats `json:"by_transform"`
	Results     []evaluate.Result          `json:"results"`
}

func BuildReport(
	records []corpus.Record,
	results []evaluate.Result,
) Report {
	return Report{
		Version:     "1",
		Overall:     CalculateEvaluationStats(results),
		ByLanguage:  GroupByLanguage(records, results),
		ByTask:      GroupByTask(records, results),
		ByTags:      GroupByTags(records, results),
		ByTransform: GroupByTransform(records, results),
		Results:     results,
	}
}
func WriteJSON(path string, report Report) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create report %q: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("encode report: %w", err)
	}

	return nil
}
