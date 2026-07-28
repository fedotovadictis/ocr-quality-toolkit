package corpus

import (
	"ocr-quality-toolkit/internal/hash"
	"ocr-quality-toolkit/internal/imageinfo"
	"os"
	"path/filepath"
	"strings"
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

	var result []Record
	for _, row := range rows {
		if row.Type != "full-page OCR ru" {
			continue
		}
		stats.MatchingTasks++

		var reference string

		for _, answer := range row.Answers {
			if strings.TrimSpace(answer) != "" {
				reference = answer
				break
			}
		}
		if reference == "" {
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

		record := Record{
			ID:         row.ID,
			Image:      row.FileName,
			References: []string{reference},
			Language:   "ru",
			Task:       row.Type,
			Width:      width,
			Height:     height,
			Format:     format,
			Tags:       []string{},
			SHA256:     sha,
		}

		result = append(result, record)
		stats.Imported++
	}
	return result, stats, nil
}
