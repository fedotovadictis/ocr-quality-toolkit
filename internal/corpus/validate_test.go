package corpus

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateRecordsCollectsMultipleErrors(t *testing.T) {
	records := []Record{
		{
			ID:         "b",
			Image:      "../outside.png",
			References: []string{""},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Tags:       []string{"real", "real"},
		},
		{
			ID:         "a",
			Image:      "",
			References: []string{"text"},
			Language:   "",
			Task:       "",
		},
	}

	errs := ValidateRecords(records)

	if len(errs) < 5 {
		t.Fatalf(
			"expected multiple validation errors, got %d: %v",
			len(errs),
			errs,
		)
	}

	joined := make([]string, 0, len(errs))
	for _, err := range errs {
		joined = append(joined, err.Error())
	}

	text := strings.Join(joined, "\n")

	expected := []string{
		"manifest is not sorted by id",
		"unsafe image path",
		"empty references",
		"duplicate tag",
		"empty image path",
		"empty language",
		"empty task",
	}

	for _, want := range expected {
		if !strings.Contains(text, want) {
			t.Errorf(
				"expected error containing %q, got:\n%s",
				want,
				text,
			)
		}
	}
}
func TestValidateRecordsRejectsInvalidUTF8Reference(t *testing.T) {
	invalid := string([]byte{0xff, 0xfe, 0xfd})

	records := []Record{
		{
			ID:         "page-001",
			Image:      "images/page.png",
			References: []string{invalid},
			Language:   "ru",
			Task:       "full-page OCR ru",
		},
	}

	errs := ValidateRecords(records)

	found := false
	for _, err := range errs {
		if strings.Contains(err.Error(), "invalid UTF-8 reference") {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected invalid UTF-8 error, got: %v", errs)
	}
}

func TestValidateRecordsRejectsUnsafePaths(t *testing.T) {
	tests := []string{
		"../outside.png",
		"..\\outside.png",
		"C:\\outside.png",
	}

	for _, imagePath := range tests {
		t.Run(imagePath, func(t *testing.T) {
			records := []Record{
				{
					ID:         "page-001",
					Image:      imagePath,
					References: []string{"text"},
					Language:   "ru",
					Task:       "full-page OCR ru",
				},
			}

			errs := ValidateRecords(records)

			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), "unsafe image path") {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf(
					"expected unsafe image path error for %q, got: %v",
					imagePath,
					errs,
				)
			}
		})
	}
}
func TestValidateRecordImageMetadata(t *testing.T) {
	root := t.TempDir()

	imageDir := filepath.Join(root, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		t.Fatalf("create image dir: %v", err)
	}

	imagePath := filepath.Join(imageDir, "page.png")

	img := image.NewRGBA(image.Rect(0, 0, 32, 24))

	file, err := os.Create(imagePath)
	if err != nil {
		t.Fatalf("create image: %v", err)
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatalf("encode png: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close image: %v", err)
	}

	data, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("read image: %v", err)
	}

	sum := sha256.Sum256(data)

	records := []Record{
		{
			ID:         "page-001",
			Image:      "images/page.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      32,
			Height:     24,
			Format:     "png",
			SHA256:     hex.EncodeToString(sum[:]),
		},
	}

	errs := ValidateRecordFiles(root, records)

	if len(errs) != 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}
}
func TestValidateRecordFilesDetectsMetadataErrors(t *testing.T) {
	root := t.TempDir()

	imageDir := filepath.Join(root, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		t.Fatalf("create image dir: %v", err)
	}

	imagePath := filepath.Join(imageDir, "page.png")

	img := image.NewRGBA(image.Rect(0, 0, 32, 24))

	file, err := os.Create(imagePath)
	if err != nil {
		t.Fatalf("create image: %v", err)
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatalf("encode png: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close image: %v", err)
	}

	records := []Record{
		{
			ID:         "page-001",
			Image:      "images/page.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      100,
			Height:     200,
			Format:     "jpeg",
			SHA256:     "wrong-sha256",
		},
	}

	errs := ValidateRecordFiles(root, records)

	var messages []string
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	text := strings.Join(messages, "\n")

	expected := []string{
		"image checksum mismatch",
		"width mismatch",
		"height mismatch",
		"format mismatch",
	}

	for _, want := range expected {
		if !strings.Contains(text, want) {
			t.Errorf(
				"expected error containing %q, got:\n%s",
				want,
				text,
			)
		}
	}
}
func TestValidateRecordFilesRejectsInvalidImage(t *testing.T) {
	root := t.TempDir()

	imageDir := filepath.Join(root, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		t.Fatalf("create image dir: %v", err)
	}

	imagePath := filepath.Join(imageDir, "broken.png")

	if err := os.WriteFile(
		imagePath,
		[]byte("not an image"),
		0o600,
	); err != nil {
		t.Fatalf("write broken image: %v", err)
	}

	records := []Record{
		{
			ID:         "broken",
			Image:      "images/broken.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
		},
	}

	errs := ValidateRecordFiles(root, records)

	found := false
	for _, err := range errs {
		if strings.Contains(err.Error(), "decode image") {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected decode image error, got: %v", errs)
	}
}
func TestValidateCorpusCollectsStructuralFileAndParentErrors(t *testing.T) {
	root := t.TempDir()

	records := []Record{
		{
			ID:         "b",
			ParentID:   "missing-parent",
			Image:      "images/missing.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
		},
		{
			ID:         "a",
			Image:      "../outside.png",
			References: []string{""},
			Language:   "ru",
			Task:       "full-page OCR ru",
		},
	}

	errs := ValidateCorpus(root, records)

	if len(errs) < 4 {
		t.Fatalf(
			"expected multiple validation errors, got %d: %v",
			len(errs),
			errs,
		)
	}

	var messages []string
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	text := strings.Join(messages, "\n")

	expected := []string{
		"manifest is not sorted by id",
		"unsafe image path",
		"empty references",
		"missing-parent",
	}

	for _, want := range expected {
		if !strings.Contains(text, want) {
			t.Errorf(
				"expected error containing %q, got:\n%s",
				want,
				text,
			)
		}
	}
}
