package corpus

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestImportMWSMetadataMergesDuplicateImages(t *testing.T) {
	tempDir := t.TempDir()

	imagePath := filepath.Join(tempDir, "document.png")
	createTestPNG(t, imagePath, 10, 20)

	metadataPath := filepath.Join(tempDir, "metadata.jsonl")

	content := `{"file_name":"document.png","id":"1","type":"full-page OCR ru","dataset_name":"dataset-a","answers":["первый текст"]}
{"file_name":"document.png","id":"2","type":"full-page OCR ru","dataset_name":"dataset-a","answers":["второй текст","первый текст"]}
`

	if err := os.WriteFile(metadataPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write metadata: %v", err)
	}

	records, stats, err := ImportMWSMetadata(metadataPath)
	if err != nil {
		t.Fatalf("ImportMWSMetadata returned error: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	record := records[0]

	if stats.Imported != 1 {
		t.Fatalf("expected Imported=1, got %d", stats.Imported)
	}

	if stats.MatchingTasks != 2 {
		t.Fatalf("expected MatchingTasks=2, got %d", stats.MatchingTasks)
	}

	expectedReferences := []string{
		"первый текст",
		"второй текст",
	}

	if !reflect.DeepEqual(record.References, expectedReferences) {
		t.Fatalf(
			"unexpected references: got %#v, want %#v",
			record.References,
			expectedReferences,
		)
	}

	expectedTags := []string{"dataset-a"}

	if !reflect.DeepEqual(record.Tags, expectedTags) {
		t.Fatalf(
			"unexpected tags: got %#v, want %#v",
			record.Tags,
			expectedTags,
		)
	}
}
func createTestPNG(t *testing.T, path string, width, height int) {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create test PNG: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		t.Fatalf("encode test PNG: %v", err)
	}
}
