package generator

import (
	"testing"

	"ocr-quality-toolkit/internal/corpus"
)

func TestBuildSyntheticRecord(t *testing.T) {
	parent := corpus.Record{
		ID:         "mws-001",
		Image:      "real/image_001.jpg",
		References: []string{"Пример текста"},
		Language:   "ru",
		Task:       "full-page OCR ru",
		Tags:       []string{"personal"},
	}

	got, err := BuildSyntheticRecord(
		parent,
		"synthetic/mws-001__grayscale.png",
		"grayscale",
		"42",
	)
	if err != nil {
		t.Fatalf("BuildSyntheticRecord returned error: %v", err)
	}

	if got.ID != "mws-001__grayscale" {
		t.Fatalf("unexpected id: %q", got.ID)
	}

	if got.ParentID != parent.ID {
		t.Fatalf("unexpected parent id: %q", got.ParentID)
	}

	if got.Image != "synthetic/mws-001__grayscale.png" {
		t.Fatalf("unexpected image path: %q", got.Image)
	}

	if got.Language != parent.Language {
		t.Fatalf("unexpected language: %q", got.Language)
	}

	if got.Task != parent.Task {
		t.Fatalf("unexpected task: %q", got.Task)
	}

	if len(got.References) != 1 || got.References[0] != "Пример текста" {
		t.Fatalf("unexpected references: %#v", got.References)
	}

	if got.Transform.Name != "grayscale" {
		t.Fatalf("unexpected transform: %q", got.Transform.Name)
	}

	if got.Transform.Seed != "42" {
		t.Fatalf("unexpected seed: %q", got.Transform.Seed)
	}
}
