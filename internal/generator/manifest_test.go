package generator

import "testing"

func TestMakeManifestRecord(t *testing.T) {
	record := MakeManifestRecord(
		"page-001__noise-light",
		"page-001",
		"images/page-001__noise-light.png",
		"Привет, OCR!",
		"noise-light",
		"internship-2026",
	)

	if record.ID != "page-001__noise-light" {
		t.Fatalf("unexpected id: %q", record.ID)
	}

	if record.ParentID != "page-001" {
		t.Fatalf(
			"expected parent_id %q, got %q",
			"page-001",
			record.ParentID,
		)
	}

	if record.Image != "images/page-001__noise-light.png" {
		t.Fatalf("unexpected image: %q", record.Image)
	}

	if len(record.References) != 1 {
		t.Fatalf(
			"expected 1 reference, got %d",
			len(record.References),
		)
	}

	if record.References[0] != "Привет, OCR!" {
		t.Fatalf("unexpected reference: %q", record.References[0])
	}

	if record.Transform.Name != "noise-light" {
		t.Fatalf(
			"expected transform %q, got %q",
			"noise-light",
			record.Transform.Name,
		)
	}

	if record.Transform.Seed != "internship-2026" {
		t.Fatalf(
			"expected seed %q, got %q",
			"internship-2026",
			record.Transform.Seed,
		)
	}
}
