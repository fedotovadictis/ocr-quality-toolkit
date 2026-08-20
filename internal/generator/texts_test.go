package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadTextInputs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "texts.jsonl")

	data := `{"id":"text-002","text":"Hello world"}
{"id":"text-001","text":"Счёт № 12345"}
`

	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write texts: %v", err)
	}

	got, err := ReadTextInputs(path)
	if err != nil {
		t.Fatalf("ReadTextInputs returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 texts, got %d", len(got))
	}

	if got[0].ID != "text-001" {
		t.Fatalf("unexpected first id: %q", got[0].ID)
	}

	if got[0].Text != "Счёт № 12345" {
		t.Fatalf("unexpected first text: %q", got[0].Text)
	}

	if got[1].ID != "text-002" {
		t.Fatalf("unexpected second id: %q", got[1].ID)
	}
}

func TestReadTextInputsDuplicateID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "texts.jsonl")

	data := `{"id":"text-001","text":"first"}
{"id":"text-001","text":"second"}
`

	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write texts: %v", err)
	}

	_, err := ReadTextInputs(path)
	if err == nil {
		t.Fatal("expected duplicate id error, got nil")
	}
}

func TestReadTextInputsEmptyID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "texts.jsonl")

	data := `{"id":"","text":"text"}`

	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write texts: %v", err)
	}

	_, err := ReadTextInputs(path)
	if err == nil {
		t.Fatal("expected empty id error, got nil")
	}
}

func TestReadTextInputsAllowsEmptyText(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "texts.jsonl")

	data := `{"id":"text-001","text":""}`

	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("write texts: %v", err)
	}

	got, err := ReadTextInputs(path)
	if err != nil {
		t.Fatalf("ReadTextInputs returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}
}
