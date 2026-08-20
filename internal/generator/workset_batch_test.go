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
	if got[0].Width != 20 {
		t.Fatalf("unexpected width: %d", got[0].Width)
	}

	if got[0].Height != 20 {
		t.Fatalf("unexpected height: %d", got[0].Height)
	}

	if got[0].Format != "png" {
		t.Fatalf("unexpected format: %q", got[0].Format)
	}

	if got[0].SHA256 == "" {
		t.Fatal("expected SHA-256 to be set")
	}

	if len(got[0].Tags) != 2 {
		t.Fatalf("unexpected tags: %#v", got[0].Tags)
	}

	if got[0].Tags[0] != "grayscale" ||
		got[0].Tags[1] != "synthetic" {
		t.Fatalf("tags are not sorted: %#v", got[0].Tags)
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
func TestBuildSyntheticWorksetUsesPerRecordSeed(t *testing.T) {
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

	createSource := func(name string) {
		t.Helper()

		path := filepath.Join(root, "real", name)

		img := image.NewRGBA(image.Rect(0, 0, 20, 20))

		for y := 0; y < 20; y++ {
			for x := 0; x < 20; x++ {
				img.Set(x, y, color.RGBA{
					R: uint8(x * 10),
					G: uint8(y * 10),
					B: 100,
					A: 255,
				})
			}
		}

		if err := SavePNG(path, img); err != nil {
			t.Fatalf("save source %q: %v", name, err)
		}
	}

	createSource("page-001.png")
	createSource("page-002.png")

	parents := []corpus.Record{
		{
			ID:         "page-001",
			Image:      "real/page-001.png",
			References: []string{"first"},
			Language:   "ru",
			Task:       "full-page OCR ru",
		},
		{
			ID:         "page-002",
			Image:      "real/page-002.png",
			References: []string{"second"},
			Language:   "ru",
			Task:       "full-page OCR ru",
		},
	}

	first, err := BuildSyntheticWorkset(
		root,
		parents,
		"noise-light",
		42,
	)
	if err != nil {
		t.Fatalf("first BuildSyntheticWorkset: %v", err)
	}

	if len(first) != 2 {
		t.Fatalf("expected 2 records, got %d", len(first))
	}

	if first[0].Transform.Seed == first[1].Transform.Seed {
		t.Fatalf(
			"expected different per-record seeds, got %q",
			first[0].Transform.Seed,
		)
	}

	second, err := BuildSyntheticWorkset(
		root,
		parents,
		"noise-light",
		42,
	)
	if err != nil {
		t.Fatalf("second BuildSyntheticWorkset: %v", err)
	}

	for i := range first {
		if first[i].Transform.Seed != second[i].Transform.Seed {
			t.Fatalf(
				"seed changed between runs for %q: %q != %q",
				first[i].ID,
				first[i].Transform.Seed,
				second[i].Transform.Seed,
			)
		}

		if first[i].SHA256 != second[i].SHA256 {
			t.Fatalf(
				"SHA-256 changed between runs for %q",
				first[i].ID,
			)
		}
	}
}
