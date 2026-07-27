package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadJSONL(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantCount int
		wantIDs   []string
	}{
		{
			name:      "one record",
			data:      `{"id":"1","text":"hello"}`,
			wantCount: 1,
			wantIDs:   []string{"1"},
		},
		{
			name: "multiple records",
			data: `{"id":"1","text":"hello"}
{"id":"2","text":"world"}`,
			wantCount: 2,
			wantIDs:   []string{"1", "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "test.jsonl")

			if err := os.WriteFile(path, []byte(tt.data), 0644); err != nil {
				t.Fatal(err)
			}

			records, err := ReadJSONL[Hypothesis](path)
			if err != nil {
				t.Fatalf("ReadJSONL returned error: %v", err)
			}

			if len(records) != tt.wantCount {
				t.Fatalf(
					"expected %d records, got %d",
					tt.wantCount,
					len(records),
				)
			}

			for i, wantID := range tt.wantIDs {
				if records[i].ID != wantID {
					t.Errorf(
						"record %d: expected ID %q, got %q",
						i,
						wantID,
						records[i].ID,
					)
				}
			}
		})
	}
}
func TestReadJSONLInvalidJSONIncludesLineNumber(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.jsonl")

	data := `{"id":"1","text":"hello"}
{"id":"2","text":}`

	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ReadJSONL[Hypothesis](path)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("expected error to contain line number 2, got %q", err.Error())
	}
}
