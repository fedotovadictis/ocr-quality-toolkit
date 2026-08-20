package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCorpus(t *testing.T) {
	dir := t.TempDir()
	fontPath := writeTestFont(t)

	inputs := []TextInput{
		{
			ID:   "text-002",
			Text: "Hello OCR",
		},
		{
			ID:   "text-001",
			Text: "Счёт № 12345",
		},
	}

	options := PageOptions{
		Width:      300,
		Height:     400,
		Margin:     20,
		FontPath:   fontPath,
		FontSize:   20,
		LineHeight: 28,
		Seed:       42,
	}

	got, err := GenerateCorpus(
		inputs,
		options,
		2,
		dir,
	)
	if err != nil {
		t.Fatalf("GenerateCorpus returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}

	for _, record := range got {
		if record.ID == "" {
			t.Fatal("expected non-empty id")
		}

		if len(record.References) != 1 {
			t.Fatalf(
				"unexpected references: %#v",
				record.References,
			)
		}

		if record.Width != 300 {
			t.Fatalf(
				"unexpected width: %d",
				record.Width,
			)
		}

		if record.Height != 400 {
			t.Fatalf(
				"unexpected height: %d",
				record.Height,
			)
		}

		if record.Format != "png" {
			t.Fatalf(
				"unexpected format: %q",
				record.Format,
			)
		}

		if record.SHA256 == "" {
			t.Fatal("expected SHA-256")
		}

		if len(record.Tags) != 1 ||
			record.Tags[0] != "synthetic" {
			t.Fatalf(
				"unexpected tags: %#v",
				record.Tags,
			)
		}

		imagePath := filepath.Join(
			dir,
			filepath.FromSlash(record.Image),
		)

		if _, err := os.Stat(imagePath); err != nil {
			t.Fatalf(
				"generated image does not exist: %v",
				err,
			)
		}
	}

	manifestPath := filepath.Join(
		dir,
		"manifest.jsonl",
	)

	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf(
			"manifest not created: %v",
			err,
		)
	}
}

func TestGenerateCorpusPageLimit(t *testing.T) {
	dir := t.TempDir()
	fontPath := writeTestFont(t)

	inputs := []TextInput{
		{ID: "text-001", Text: "one"},
		{ID: "text-002", Text: "two"},
		{ID: "text-003", Text: "three"},
	}

	options := PageOptions{
		Width:      300,
		Height:     400,
		Margin:     20,
		FontPath:   fontPath,
		FontSize:   20,
		LineHeight: 28,
		Seed:       42,
	}

	got, err := GenerateCorpus(
		inputs,
		options,
		2,
		dir,
	)
	if err != nil {
		t.Fatalf("GenerateCorpus returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf(
			"expected 2 records, got %d",
			len(got),
		)
	}
}
