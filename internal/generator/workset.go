package generator

import (
	"fmt"

	"ocr-quality-toolkit/internal/corpus"
)

func BuildSyntheticRecord(
	parent corpus.Record,
	imagePath string,
	profile string,
	seed string,
) (corpus.Record, error) {
	if parent.ID == "" {
		return corpus.Record{}, fmt.Errorf("parent id is empty")
	}

	if profile == "" {
		return corpus.Record{}, fmt.Errorf("transform profile is empty")
	}

	id := parent.ID + "__" + profile

	return corpus.Record{
		ID:         id,
		ParentID:   parent.ID,
		Image:      imagePath,
		References: append([]string(nil), parent.References...),
		Language:   parent.Language,
		Task:       parent.Task,
		Tags: []string{
			"synthetic",
			profile,
		},
		Transform: corpus.Transform{
			Name: profile,
			Seed: seed,
		},
	}, nil
}
