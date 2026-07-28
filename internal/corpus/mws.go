package corpus

import "strings"

type mwsRecord struct {
	FileName    string   `json:"file_name"`
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	DatasetName string   `json:"dataset_name"`
	Question    string   `json:"question"`
	Answers     []string `json:"answers"`
}
type ImportStats struct {
	TotalLines      int
	MatchingTasks   int
	Imported        int
	MissingImages   int
	EmptyReferences int
}

func ImportMWSMetadata(path string) ([]mwsRecord, ImportStats, error) {
	rows, err := ReadJSONL[mwsRecord](path)
	if err != nil {
		return nil, ImportStats{}, err
	}

	stats := ImportStats{
		TotalLines: len(rows),
	}

	var result []mwsRecord
	for _, row := range rows {
		if row.Type != "full-page OCR ru" {
			continue
		}
		stats.MatchingTasks += 1

		var reference string

		for _, answer := range row.Answers {
			if strings.TrimSpace(answer) != "" {
				reference = answer
				break
			}
		}
		if reference == "" {
			stats.EmptyReferences += 1
			continue
		}
		result = append(result, row)
	}
	return result, stats, nil
}
