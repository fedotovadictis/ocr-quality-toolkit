package evaluate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteReport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.json")

	want := []Result{
		{
			ID:         "1",
			CER:        0,
			WER:        0,
			Similarity: 1,
			ExactMatch: true,
			Status:     StatusSuccess,
		},
	}

	if err := WriteReport(path, want); err != nil {
		t.Fatalf("WriteReport returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}

	var got []Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode report: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected report:\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestWriteReportInvalidPath(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"missing-directory",
		"report.json",
	)

	err := WriteReport(path, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
