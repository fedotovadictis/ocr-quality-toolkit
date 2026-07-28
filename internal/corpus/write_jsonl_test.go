package corpus

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteJSONL(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "manifest.jsonl")

	want := []Record{
		{
			ID:         "record-1",
			Image:      "images/one.png",
			References: []string{"first text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      100,
			Height:     200,
			Format:     "png",
			Tags:       []string{"dataset-a"},
			SHA256:     "hash-1",
		},
		{
			ID:         "record-2",
			Image:      "images/two.png",
			References: []string{"second text", "another text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      300,
			Height:     400,
			Format:     "png",
			Tags:       []string{"dataset-b"},
			SHA256:     "hash-2",
		},
	}

	if err := WriteJSONL(outputPath, want); err != nil {
		t.Fatalf("WriteJSONL returned error: %v", err)
	}

	got, err := ReadJSONL[Record](outputPath)
	if err != nil {
		t.Fatalf("ReadJSONL returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected records:\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestWriteJSONLInvalidPath(t *testing.T) {
	tempDir := t.TempDir()

	outputPath := filepath.Join(
		tempDir,
		"missing-directory",
		"manifest.jsonl",
	)

	err := WriteJSONL(outputPath, []Record{
		{ID: "record-1"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
