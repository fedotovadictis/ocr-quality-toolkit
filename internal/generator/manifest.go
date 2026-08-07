package generator

import "ocr-quality-toolkit/internal/corpus"

// MakeManifestRecord создаёт запись manifest для производного изображения.
func MakeManifestRecord(
	id string,
	parentID string,
	imagePath string,
	reference string,
	transformName string,
	seed string,
) corpus.Record {
	return corpus.Record{
		ID:         id,
		ParentID:   parentID,
		Image:      imagePath,
		References: []string{reference},
		Language:   "ru",
		Task:       "full-page OCR",
		Tags: []string{
			"synthetic",
			transformName,
		},
		Transform: corpus.Transform{
			Name: transformName,
			Seed: seed,
		},
	}
}
