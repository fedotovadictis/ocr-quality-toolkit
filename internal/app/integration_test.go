package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"ocr-quality-toolkit/internal/report"
)

func TestFullEvaluationReportFlow(t *testing.T) {
	dir := t.TempDir()

	manifestPath := filepath.Join(dir, "manifest.jsonl")
	hypothesesPath := filepath.Join(dir, "hypotheses.jsonl")
	jsonReportPath := filepath.Join(dir, "report.json")
	htmlReportPath := filepath.Join(dir, "report.html")
	baselinePath := filepath.Join(dir, "baseline.json")

	manifest := `{"id":"page-001","image":"images/page-001.png","references":["Привет мир"],"language":"ru","task":"ocr","tags":["synthetic"],"transform":{"name":"grayscale"}}` + "\n"

	hypotheses := `{"id":"page-001","text":"Привет мр","engine":"tesseract","model":"test"}` + "\n"

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

	err := Run(
		[]string{
			"evaluate",
			"-manifest", manifestPath,
			"-hypotheses", hypothesesPath,
			"-normalization", "strict",
			"-out", jsonReportPath,
		},
		os.Stdout,
		os.Stderr,
	)
	if err != nil {
		t.Fatalf("evaluate failed: %v", err)
	}

	data, err := os.ReadFile(jsonReportPath)
	if err != nil {
		t.Fatalf("read JSON report: %v", err)
	}

	var got report.Report

	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode JSON report: %v", err)
	}

	if got.Overall.Total != 1 {
		t.Fatalf(
			"expected total 1, got %d",
			got.Overall.Total,
		)
	}

	if len(got.Results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(got.Results),
		)
	}

	if got.Results[0].Reference != "Привет мир" {
		t.Fatalf(
			"unexpected reference: %q",
			got.Results[0].Reference,
		)
	}

	if got.Results[0].Hypothesis != "Привет мр" {
		t.Fatalf(
			"unexpected hypothesis: %q",
			got.Results[0].Hypothesis,
		)
	}

	if err := report.WriteHTML(
		htmlReportPath,
		got,
	); err != nil {
		t.Fatalf("write HTML report: %v", err)
	}

	if _, err := os.Stat(htmlReportPath); err != nil {
		t.Fatalf("HTML report not created: %v", err)
	}

	if err := report.WriteJSON(
		baselinePath,
		got,
	); err != nil {
		t.Fatalf("write baseline: %v", err)
	}

	err = Run(
		[]string{
			"compare",
			"-baseline", baselinePath,
			"-current", jsonReportPath,
			"-max-cer-increase", "0",
			"-max-wer-increase", "0",
			"-max-coverage-decrease", "0",
		},
		os.Stdout,
		os.Stderr,
	)
	if err != nil {
		t.Fatalf("compare failed: %v", err)
	}
}
