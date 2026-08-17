package corpus

import (
	"testing"
)

func TestValidateCorpusIntegrity(t *testing.T) {
	records := []Record{
		{
			ID:         "real-001",
			Image:      "images/real-001.png",
			References: []string{"Пример текста"},
			Language:   "ru",
			Task:       "ocr",
			Tags:       []string{"real"},
		},
		{
			ID:         "synthetic-001",
			ParentID:   "real-001",
			Image:      "images/synthetic-001.png",
			References: []string{"Пример текста"},
			Language:   "ru",
			Task:       "ocr",
			Tags:       []string{"synthetic"},
			Transform: Transform{
				Name: "grayscale",
				Seed: "42",
			},
		},
	}

	if err := ValidateCorpusIntegrity(records); err != nil {
		t.Fatalf(
			"ValidateCorpusIntegrity returned error: %v",
			err,
		)
	}
}
func TestValidateCorpusIntegrityDuplicateID(t *testing.T) {
	records := []Record{
		{
			ID:         "same-id",
			Image:      "images/1.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "ocr",
		},
		{
			ID:         "same-id",
			Image:      "images/2.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "ocr",
		},
	}

	err := ValidateCorpusIntegrity(records)
	if err == nil {
		t.Fatal("expected duplicate ID error, got nil")
	}
}

func TestValidateCorpusIntegrityEmptyReference(t *testing.T) {
	records := []Record{
		{
			ID:         "page-001",
			Image:      "images/page-001.png",
			References: []string{"   "},
			Language:   "ru",
			Task:       "ocr",
		},
	}

	err := ValidateCorpusIntegrity(records)
	if err == nil {
		t.Fatal("expected empty reference error, got nil")
	}
}

func TestValidateCorpusIntegrityEmptyImage(t *testing.T) {
	records := []Record{
		{
			ID:         "page-001",
			Image:      "",
			References: []string{"text"},
			Language:   "ru",
			Task:       "ocr",
		},
	}

	err := ValidateCorpusIntegrity(records)
	if err == nil {
		t.Fatal("expected empty image path error, got nil")
	}
}
