package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendResult(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.jsonl")

	first := Result{
		ID:         "page-001",
		Text:       "first text",
		DurationMS: 10,
	}

	second := Result{
		ID:         "page-002",
		Error:      "ocr error",
		Stderr:     "failed",
		DurationMS: 20,
	}

	if err := AppendResult(path, first); err != nil {
		t.Fatalf("append first result: %v", err)
	}

	if err := AppendResult(path, second); err != nil {
		t.Fatalf("append second result: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read results: %v", err)
	}

	lines := strings.Split(
		strings.TrimSpace(string(data)),
		"\n",
	)

	if len(lines) != 2 {
		t.Fatalf(
			"expected 2 lines, got %d",
			len(lines),
		)
	}

	results, err := ReadResults(path)
	if err != nil {
		t.Fatalf("ReadResults: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 results, got %d",
			len(results),
		)
	}

	if results[0].ID != "page-001" {
		t.Fatalf(
			"expected first id %q, got %q",
			"page-001",
			results[0].ID,
		)
	}

	if results[1].ID != "page-002" {
		t.Fatalf(
			"expected second id %q, got %q",
			"page-002",
			results[1].ID,
		)
	}

	if results[1].Error != "ocr error" {
		t.Fatalf(
			"expected error %q, got %q",
			"ocr error",
			results[1].Error,
		)
	}
}
