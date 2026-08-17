package generator

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"ocr-quality-toolkit/internal/corpus"
)

func TestBuildSyntheticWorkset(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(
		filepath.Join(root, "real"),
		0o755,
	); err != nil {
		t.Fatalf("create real dir: %v", err)
	}

	if err := os.MkdirAll(
		filepath.Join(root, "synthetic"),
		0o755,
	); err != nil {
		t.Fatalf("create synthetic dir: %v", err)
	}

	sourcePath := filepath.Join(root, "real", "page.png")

	source := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			source.Set(x, y, color.RGBA{
				R: uint8(x * 10),
				G: uint8(y * 10),
				B: 100,
				A: 255,
			})
		}
	}

	if err := SavePNG(sourcePath, source); err != nil {
		t.Fatalf("save source: %v", err)
	}

	parents := []corpus.Record{
		{
			ID:         "page-001",
			Image:      "real/page.png",
			References: []string{"Пример текста"},
			Language:   "ru",
			Task:       "full-page OCR ru",
		},
	}

	got, err := BuildSyntheticWorkset(
		root,
		parents,
		"grayscale",
		42,
	)
	if err != nil {
		t.Fatalf("BuildSyntheticWorkset returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}

	if got[0].ParentID != "page-001" {
		t.Fatalf("unexpected parent id: %q", got[0].ParentID)
	}

	target := filepath.Join(
		root,
		"synthetic",
		"page-001__grayscale.png",
	)

	if _, err := os.Stat(target); err != nil {
		t.Fatalf("synthetic image not created: %v", err)
	}
}
