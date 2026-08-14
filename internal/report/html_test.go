package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ocr-quality-toolkit/internal/evaluate"
)

func TestWriteHTML(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"report.html",
	)

	report := Report{
		Version: "1",
		Overall: EvaluationStats{
			Total:    1,
			Coverage: 1,
			CER:      0.1,
			WER:      0.2,
		},
		Results: []evaluate.Result{
			{
				ID:         "page-001",
				Status:     evaluate.StatusSuccess,
				CER:        0.1,
				WER:        0.2,
				Reference:  "Привет мир",
				Hypothesis: "Привет мр",
				Image:      "images/page-001.png",
			},
		},
	}

	if err := WriteHTML(path, report); err != nil {
		t.Fatalf("WriteHTML returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read HTML report: %v", err)
	}

	html := string(data)

	expected := []string{
		"page-001",
		"Привет мир",
		"Привет мр",
		"images/page-001.png",
		"CER",
		"WER",
	}

	for _, value := range expected {
		if !strings.Contains(html, value) {
			t.Fatalf(
				"expected HTML to contain %q",
				value,
			)
		}
	}
}
