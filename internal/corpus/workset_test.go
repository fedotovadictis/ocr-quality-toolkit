package corpus

import "testing"

func TestBuildWorkset(t *testing.T) {
	real := []Record{
		{
			ID:         "real-001",
			Image:      "real/page-001.png",
			References: []string{"Привет мир"},
			Language:   "ru",
			Task:       "ocr",
			Tags:       []string{"real"},
		},
	}

	synthetic := []Record{
		{
			ID:         "synthetic-001",
			ParentID:   "real-001",
			Image:      "synthetic/page-001-grayscale.png",
			References: []string{"Привет мир"},
			Language:   "ru",
			Task:       "ocr",
			Tags:       []string{"synthetic"},
			Transform: Transform{
				Name: "grayscale",
				Seed: "42",
			},
		},
	}

	got, err := BuildWorkset(real, synthetic)
	if err != nil {
		t.Fatalf("BuildWorkset returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}

	if got[0].ID != "real-001" {
		t.Fatalf("unexpected first record: %q", got[0].ID)
	}

	if got[1].ID != "synthetic-001" {
		t.Fatalf("unexpected second record: %q", got[1].ID)
	}
}
func TestBuildWorksetDuplicateID(t *testing.T) {
	real := []Record{
		{
			ID:         "page-001",
			Image:      "real/page-001.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "ocr",
		},
	}

	synthetic := []Record{
		{
			ID:         "page-001",
			ParentID:   "page-001",
			Image:      "synthetic/page-001.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "ocr",
			Transform: Transform{
				Name: "grayscale",
			},
		},
	}

	_, err := BuildWorkset(real, synthetic)
	if err == nil {
		t.Fatal("expected duplicate ID error, got nil")
	}
}
