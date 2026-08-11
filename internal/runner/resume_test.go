package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilterPendingTasks(t *testing.T) {
	tasks := []Task{
		{ID: "1", ImagePath: "1.png"},
		{ID: "2", ImagePath: "2.png"},
		{ID: "3", ImagePath: "3.png"},
	}

	existing := []Result{
		{ID: "1", Text: "done"},
		{ID: "3", Error: "ocr error"},
	}

	pending := FilterPendingTasks(tasks, existing)

	if len(pending) != 1 {
		t.Fatalf(
			"expected 1 pending task, got %d",
			len(pending),
		)
	}

	if pending[0].ID != "2" {
		t.Fatalf(
			"expected task %q, got %q",
			"2",
			pending[0].ID,
		)
	}
}
func TestReadResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "results.jsonl")

	data := "" +
		`{"id":"1","text":"hello","error":"","stderr":"","duration_ms":10}` + "\n" +
		`{"id":"2","text":"","error":"ocr error","stderr":"failed","duration_ms":20}` + "\n"

	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write results: %v", err)
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

	if results[0].ID != "1" {
		t.Fatalf(
			"expected first id %q, got %q",
			"1",
			results[0].ID,
		)
	}

	if results[0].Text != "hello" {
		t.Fatalf(
			"expected text %q, got %q",
			"hello",
			results[0].Text,
		)
	}

	if results[1].ID != "2" {
		t.Fatalf(
			"expected second id %q, got %q",
			"2",
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
