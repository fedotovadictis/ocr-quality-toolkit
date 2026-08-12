package app

import (
	"bytes"
	"encoding/json"
	"ocr-quality-toolkit/internal/evaluate"
	"ocr-quality-toolkit/internal/runner"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunNoArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "usage") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"unknown"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunImportMWSMissingMetadata(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"import-mws"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "missing required flag: -metadata") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunImportMWSMissingOutput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"import-mws",
			"-metadata", "metadata.jsonl",
		},
		&stdout,
		&stderr,
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "missing required flag: -output") {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestRunEvaluate(t *testing.T) {
	tempDir := t.TempDir()

	manifestPath := filepath.Join(tempDir, "manifest.jsonl")
	hypothesesPath := filepath.Join(tempDir, "hypotheses.jsonl")
	reportPath := filepath.Join(tempDir, "report.json")

	manifest := `{"id":"1","references":["Кот!"]}
`
	hypotheses := `{"id":"1","text":"кот","engine":"test","model":"test"}
`

	if err := os.WriteFile(
		manifestPath,
		[]byte(manifest),
		0o600,
	); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	if err := os.WriteFile(
		hypothesesPath,
		[]byte(hypotheses),
		0o600,
	); err != nil {
		t.Fatalf("write hypotheses: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"evaluate",
			"-manifest", manifestPath,
			"-hypotheses", hypothesesPath,
			"-normalization", "plain-text-ru",
			"-out", reportPath,
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}

	var results []evaluate.Result
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatalf("decode report: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	result := results[0]

	if result.ID != "1" {
		t.Fatalf("unexpected ID: %q", result.ID)
	}

	if !result.ExactMatch {
		t.Fatal("expected exact match")
	}

	if result.Status != evaluate.StatusSuccess {
		t.Fatalf(
			"expected status %q, got %q",
			evaluate.StatusSuccess,
			result.Status,
		)
	}

	if result.CER != 0 ||
		result.WER != 0 ||
		result.Similarity != 1 {

		t.Fatalf(
			"unexpected metrics: CER=%v WER=%v similarity=%v",
			result.CER,
			result.WER,
			result.Similarity,
		)
	}
}
func TestRunTesseractInvalidWorkers(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"run-tesseract",
			"-manifest", "manifest.jsonl",
			"-out", "results.jsonl",
			"-workers", "0",
		},
		&stdout,
		&stderr,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"workers must be greater than zero",
	) {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestRunTesseractResumeSkipsCompletedRecords(t *testing.T) {
	dir := t.TempDir()

	manifestPath := filepath.Join(dir, "manifest.jsonl")
	resultsPath := filepath.Join(dir, "results.jsonl")

	manifest := `{"id":"page-001","image":"image.png","references":["text"],"language":"rus","task":"ocr","width":100,"height":100,"format":"png","tags":[],"sha256":"abc"}` + "\n"

	if err := os.WriteFile(
		manifestPath,
		[]byte(manifest),
		0o600,
	); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	existing := `{"id":"page-001","text":"already done","error":"","stderr":"","duration_ms":10}` + "\n"

	if err := os.WriteFile(
		resultsPath,
		[]byte(existing),
		0o600,
	); err != nil {
		t.Fatalf("write results: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"run-tesseract",
			"-manifest", manifestPath,
			"-out", resultsPath,
			"-binary", "this-binary-must-not-run",
			"-workers", "1",
			"-resume",
		},
		&stdout,
		&stderr,
	)

	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !strings.Contains(stdout.String(), "processed: 0") {
		t.Fatalf(
			"expected no records to be processed, got %q",
			stdout.String(),
		)
	}

	results, err := runner.ReadResults(resultsPath)
	if err != nil {
		t.Fatalf("ReadResults: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result after resume, got %d",
			len(results),
		)
	}

	if results[0].ID != "page-001" {
		t.Fatalf(
			"expected id %q, got %q",
			"page-001",
			results[0].ID,
		)
	}
}
