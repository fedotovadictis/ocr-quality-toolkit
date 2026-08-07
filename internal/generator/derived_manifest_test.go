package generator

import (
	"testing"

	"ocr-quality-toolkit/internal/corpus"
)

func TestBuildDerivedManifest(t *testing.T) {
	parent := corpus.Record{
		ID:         "synthetic-001",
		Image:      "images/synthetic-001.png",
		References: []string{"Привет, OCR!"},
		Language:   "ru",
		Task:       "full-page OCR",
		Tags:       []string{"synthetic"},
	}

	profiles := []string{
		"grayscale",
		"low-contrast",
		"noise-light",
		"jpeg-70",
		"downscale-50",
	}

	records := BuildDerivedManifest(
		parent,
		profiles,
		"internship-2026",
	)

	if len(records) != len(profiles) {
		t.Fatalf(
			"expected %d records, got %d",
			len(profiles),
			len(records),
		)
	}

	for i, record := range records {
		expectedProfile := profiles[i]

		if record.ParentID != parent.ID {
			t.Fatalf(
				"expected parent_id %q, got %q",
				parent.ID,
				record.ParentID,
			)
		}

		expectedID := parent.ID + "__" + expectedProfile
		if record.ID != expectedID {
			t.Fatalf(
				"expected id %q, got %q",
				expectedID,
				record.ID,
			)
		}

		if record.Transform.Name != expectedProfile {
			t.Fatalf(
				"expected transform %q, got %q",
				expectedProfile,
				record.Transform.Name,
			)
		}

		if record.Transform.Seed != "internship-2026" {
			t.Fatalf(
				"unexpected seed %q",
				record.Transform.Seed,
			)
		}

		if len(record.References) != 1 ||
			record.References[0] != parent.References[0] {
			t.Fatalf(
				"unexpected references: %v",
				record.References,
			)
		}
	}
}
func TestBuildDerivedManifestHasValidParentIDs(t *testing.T) {
	parent := corpus.Record{
		ID:         "synthetic-001",
		Image:      "images/synthetic-001.png",
		References: []string{"Привет, OCR!"},
		Language:   "ru",
		Task:       "full-page OCR",
		Tags:       []string{"synthetic"},
	}

	profiles := []string{
		"grayscale",
		"low-contrast",
		"noise-light",
		"jpeg-70",
		"downscale-50",
	}

	derived := BuildDerivedManifest(
		parent,
		profiles,
		"internship-2026",
	)

	records := make([]corpus.Record, 0, len(derived)+1)
	records = append(records, parent)
	records = append(records, derived...)

	if err := corpus.ValidateParentIDs(records); err != nil {
		t.Fatalf(
			"ValidateParentIDs returned error: %v",
			err,
		)
	}
}
