package corpus

import (
	"ocr-quality-toolkit/internal/hash"
	"ocr-quality-toolkit/internal/imageinfo"
	"os"
	"path/filepath"
	"sort"
)

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
	InvalidImages   int
	EmptyReferences int
}

func ImportMWSMetadata(path string) ([]Record, ImportStats, error) {

	datasetDir := filepath.Dir(path)
	rows, err := ReadJSONL[mwsRecord](path)
	if err != nil {
		return nil, ImportStats{}, err
	}

	stats := ImportStats{
		TotalLines: len(rows),
	}

	recordsByImage := make(map[string]*Record)
	for _, row := range rows {
		if row.Type != "full-page OCR ru" {
			continue
		}
		stats.MatchingTasks++

		references := nonEmptyUniqueStrings(row.Answers)
		if len(references) == 0 {
			stats.EmptyReferences++
			continue
		}

		imagePath := filepath.Join(datasetDir, row.FileName)

		info, err := os.Stat(imagePath)
		if err != nil {
			if os.IsNotExist(err) {
				stats.MissingImages++
				continue
			}
			return nil, stats, err
		}
		if info.IsDir() {
			stats.MissingImages++
			continue
		}
		format, width, height, err := imageinfo.Read(imagePath)
		if err != nil {
			stats.InvalidImages++
			continue
		}

		sha, err := hash.FileSHA256(imagePath)
		if err != nil {
			stats.InvalidImages++
			continue
		}

		normalizedImagePath := filepath.ToSlash(filepath.Clean(row.FileName))

		if existing, ok := recordsByImage[normalizedImagePath]; ok {
			existing.References = appendUnique(existing.References, references...)
			existing.Tags = appendUnique(existing.Tags, row.DatasetName)
			continue
		}

		record := &Record{
			ID:         makeMWSID(row.ID, normalizedImagePath),
			Image:      normalizedImagePath,
			References: references,
			Language:   "ru",
			Task:       row.Type,
			Width:      width,
			Height:     height,
			Format:     format,
			Tags:       appendUnique(nil, row.DatasetName),
			SHA256:     sha,
		}

		recordsByImage[normalizedImagePath] = record
		stats.Imported++
	}
	result := make([]Record, 0, len(recordsByImage))

	for _, record := range recordsByImage {
		result = append(result, *record)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	for i := range result {
		sort.Strings(result[i].Tags)
	}
	return result, stats, nil
}
