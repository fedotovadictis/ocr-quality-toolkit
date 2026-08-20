package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/evaluate"
	"ocr-quality-toolkit/internal/generator"
	"ocr-quality-toolkit/internal/report"
	"ocr-quality-toolkit/internal/runner"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
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

func TestRunCorpusImportMWSMissingSource(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"corpus",
			"import-mws",
		},
		&stdout,
		&stderr,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"missing required flag: -source",
	) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunCorpusImportMWSMissingOutput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"corpus",
			"import-mws",
			"-source", "testdata/mws",
		},
		&stdout,
		&stderr,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"missing required flag: -out",
	) {
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

	var got report.Report
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode report: %v", err)
	}

	if len(got.Results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(got.Results),
		)
	}

	if got.Overall.Total != 1 {
		t.Fatalf(
			"expected overall total 1, got %d",
			got.Overall.Total,
		)
	}

	result := got.Results[0]

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
func TestRunCompareSuccess(t *testing.T) {
	dir := t.TempDir()

	baselinePath := filepath.Join(dir, "baseline.json")
	currentPath := filepath.Join(dir, "current.json")

	baseline := report.Report{
		Overall: report.EvaluationStats{
			CER:      0.10,
			WER:      0.20,
			Coverage: 0.90,
		},
	}

	current := report.Report{
		Overall: report.EvaluationStats{
			CER:      0.11,
			WER:      0.19,
			Coverage: 0.92,
		},
	}

	if err := report.WriteJSON(baselinePath, baseline); err != nil {
		t.Fatalf("write baseline: %v", err)
	}

	if err := report.WriteJSON(currentPath, current); err != nil {
		t.Fatalf("write current: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"compare",
			"-baseline", baselinePath,
			"-current", currentPath,
			"-max-cer-increase", "0.02",
			"-max-wer-increase", "0.02",
			"-max-coverage-decrease", "0.05",
		},
		&stdout,
		&stderr,
	)

	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !strings.Contains(stdout.String(), "CER delta") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestRunCompareRegression(t *testing.T) {
	dir := t.TempDir()

	baselinePath := filepath.Join(dir, "baseline.json")
	currentPath := filepath.Join(dir, "current.json")

	baseline := report.Report{
		Overall: report.EvaluationStats{
			CER:      0.10,
			WER:      0.20,
			Coverage: 0.95,
		},
	}

	current := report.Report{
		Overall: report.EvaluationStats{
			CER:      0.20,
			WER:      0.20,
			Coverage: 0.95,
		},
	}

	if err := report.WriteJSON(baselinePath, baseline); err != nil {
		t.Fatalf("write baseline: %v", err)
	}

	if err := report.WriteJSON(currentPath, current); err != nil {
		t.Fatalf("write current: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"compare",
			"-baseline", baselinePath,
			"-current", currentPath,
			"-max-cer-increase", "0.02",
			"-max-wer-increase", "0.02",
			"-max-coverage-decrease", "0.05",
		},
		&stdout,
		&stderr,
	)

	if !errors.Is(err, ErrRegression) {
		t.Fatalf(
			"expected ErrRegression, got %v",
			err,
		)
	}
}
func TestRunBuildWorkset(t *testing.T) {
	dir := t.TempDir()

	realDir := filepath.Join(dir, "real")
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatalf("create real dir: %v", err)
	}

	sourcePath := filepath.Join(realDir, "page.png")

	source := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			source.Set(x, y, color.RGBA{
				R: uint8(x * 10),
				G: uint8(y * 10),
				B: 100,
				A: 255,
			})
		}
	}

	if err := generator.SavePNG(sourcePath, source); err != nil {
		t.Fatalf("save source image: %v", err)
	}

	realManifestPath := filepath.Join(dir, "real-manifest.jsonl")
	syntheticManifestPath := filepath.Join(dir, "synthetic-manifest.jsonl")
	worksetManifestPath := filepath.Join(dir, "manifest.jsonl")

	realRecords := []corpus.Record{
		{
			ID:         "page-001",
			Image:      "real/page.png",
			References: []string{"Пример текста"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Tags:       []string{"real"},
		},
	}

	if err := corpus.WriteJSONL(realManifestPath, realRecords); err != nil {
		t.Fatalf("write real manifest: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"build-workset",
			"-root", dir,
			"-real-manifest", realManifestPath,
			"-synthetic-manifest", syntheticManifestPath,
			"-out", worksetManifestPath,
			"-profile", "grayscale",
			"-seed", "42",
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf(
			"build-workset failed: %v\nstderr: %s",
			err,
			stderr.String(),
		)
	}

	workset, err := corpus.ReadManifest(worksetManifestPath)
	if err != nil {
		t.Fatalf("read workset manifest: %v", err)
	}

	if len(workset) != 2 {
		t.Fatalf("expected 2 records, got %d", len(workset))
	}

	if workset[0].ID != "page-001" {
		t.Fatalf("unexpected real id: %q", workset[0].ID)
	}

	if workset[1].ParentID != "page-001" {
		t.Fatalf(
			"unexpected synthetic parent id: %q",
			workset[1].ParentID,
		)
	}

	syntheticPath := filepath.Join(
		dir,
		"synthetic",
		"page-001__grayscale.png",
	)

	if _, err := os.Stat(syntheticPath); err != nil {
		t.Fatalf("synthetic image not created: %v", err)
	}
}
func TestRunCorpusValidateSuccess(t *testing.T) {
	dir := t.TempDir()

	imageDir := filepath.Join(dir, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		t.Fatalf("create image dir: %v", err)
	}

	imagePath := filepath.Join(imageDir, "page.png")

	img := image.NewRGBA(image.Rect(0, 0, 32, 24))

	file, err := os.Create(imagePath)
	if err != nil {
		t.Fatalf("create image: %v", err)
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatalf("encode image: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close image: %v", err)
	}

	data, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("read image: %v", err)
	}

	sum := sha256.Sum256(data)

	manifestPath := filepath.Join(dir, "manifest.jsonl")

	records := []corpus.Record{
		{
			ID:         "page-001",
			Image:      "images/page.png",
			References: []string{"Тест"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      32,
			Height:     24,
			Format:     "png",
			Tags:       []string{"real"},
			SHA256:     hex.EncodeToString(sum[:]),
		},
	}

	if err := corpus.WriteJSONL(manifestPath, records); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err = Run(
		[]string{
			"corpus",
			"validate",
			"-manifest",
			manifestPath,
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf(
			"corpus validate failed: %v\nstderr: %s",
			err,
			stderr.String(),
		)
	}

	if !strings.Contains(stdout.String(), "Corpus is valid") {
		t.Fatalf(
			"unexpected stdout: %s",
			stdout.String(),
		)
	}

	if !strings.Contains(stdout.String(), "Records: 1") {
		t.Fatalf(
			"unexpected stdout: %s",
			stdout.String(),
		)
	}
}
func TestRunCorpusStats(t *testing.T) {
	dir := t.TempDir()

	manifestPath := filepath.Join(dir, "manifest.jsonl")

	records := []corpus.Record{
		{
			ID:         "real-001",
			Image:      "images/real.jpg",
			References: []string{"text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Format:     "jpeg",
		},
		{
			ID:         "synthetic-001",
			ParentID:   "real-001",
			Image:      "images/synthetic.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Format:     "png",
			Transform: corpus.Transform{
				Name: "grayscale",
			},
		},
	}

	if err := corpus.WriteJSONL(manifestPath, records); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"corpus",
			"stats",
			"-manifest",
			manifestPath,
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf(
			"corpus stats failed: %v\nstderr: %s",
			err,
			stderr.String(),
		)
	}

	output := stdout.String()

	expected := []string{
		"Records: 2",
		"Real: 1",
		"Synthetic: 1",
		"ru: 2",
		"jpeg: 1",
		"png: 1",
		"grayscale: 1",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Fatalf(
				"expected output to contain %q, got:\n%s",
				want,
				output,
			)
		}
	}
}

func TestRunImageTransform(t *testing.T) {
	dir := t.TempDir()

	imageDir := filepath.Join(dir, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		t.Fatalf("create image dir: %v", err)
	}

	sourcePath := filepath.Join(imageDir, "page.png")

	img := image.NewRGBA(image.Rect(0, 0, 20, 20))

	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 10),
				G: uint8(y * 10),
				B: 100,
				A: 255,
			})
		}
	}

	if err := generator.SavePNG(sourcePath, img); err != nil {
		t.Fatalf("save source image: %v", err)
	}

	manifestPath := filepath.Join(dir, "manifest.jsonl")
	outputDir := filepath.Join(dir, "transformed")

	records := []corpus.Record{
		{
			ID:         "page-001",
			Image:      "images/page.png",
			References: []string{"Пример текста"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      20,
			Height:     20,
			Format:     "png",
			Tags:       []string{"synthetic"},
		},
	}

	if err := corpus.WriteJSONL(manifestPath, records); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"image",
			"transform",
			"-manifest", manifestPath,
			"-profiles", "grayscale,noise-light",
			"-seed", "42",
			"-out", outputDir,
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf(
			"image transform failed: %v\nstderr: %s",
			err,
			stderr.String(),
		)
	}

	outputManifest := filepath.Join(
		outputDir,
		"manifest.jsonl",
	)

	got, err := corpus.ReadManifest(outputManifest)
	if err != nil {
		t.Fatalf("read transformed manifest: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf(
			"expected 2 transformed records, got %d",
			len(got),
		)
	}

	expectedProfiles := []string{
		"grayscale",
		"noise-light",
	}

	for i, profile := range expectedProfiles {
		record := got[i]

		expectedID := "page-001__" + profile

		if record.ID != expectedID {
			t.Fatalf(
				"expected id %q, got %q",
				expectedID,
				record.ID,
			)
		}

		if record.ParentID != "page-001" {
			t.Fatalf(
				"unexpected parent_id: %q",
				record.ParentID,
			)
		}

		if record.Transform.Name != profile {
			t.Fatalf(
				"expected transform %q, got %q",
				profile,
				record.Transform.Name,
			)
		}

		if record.Transform.Seed != "42" {
			t.Fatalf(
				"unexpected seed: %q",
				record.Transform.Seed,
			)
		}

		if record.Width != 20 {
			t.Fatalf(
				"unexpected width for %q: %d",
				record.ID,
				record.Width,
			)
		}

		if record.Height != 20 {
			t.Fatalf(
				"unexpected height for %q: %d",
				record.ID,
				record.Height,
			)
		}

		if record.Format != "png" {
			t.Fatalf(
				"unexpected format for %q: %q",
				record.ID,
				record.Format,
			)
		}

		if record.SHA256 == "" {
			t.Fatalf(
				"empty SHA-256 for %q",
				record.ID,
			)
		}

		transformedPath := filepath.Join(
			outputDir,
			filepath.FromSlash(record.Image),
		)

		if _, err := os.Stat(transformedPath); err != nil {
			t.Fatalf(
				"transformed image %q not created: %v",
				record.ID,
				err,
			)
		}
	}

	if !strings.Contains(
		stdout.String(),
		"transformed records: 2",
	) {
		t.Fatalf(
			"unexpected stdout:\n%s",
			stdout.String(),
		)
	}
}
func TestRunImageTransformUnknownProfile(t *testing.T) {
	dir := t.TempDir()

	imageDir := filepath.Join(dir, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		t.Fatalf("create image dir: %v", err)
	}

	sourcePath := filepath.Join(imageDir, "page.png")

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if err := generator.SavePNG(sourcePath, img); err != nil {
		t.Fatalf("save source image: %v", err)
	}

	manifestPath := filepath.Join(dir, "manifest.jsonl")

	records := []corpus.Record{
		{
			ID:         "page-001",
			Image:      "images/page.png",
			References: []string{"text"},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      10,
			Height:     10,
			Format:     "png",
		},
	}

	if err := corpus.WriteJSONL(manifestPath, records); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"image",
			"transform",
			"-manifest", manifestPath,
			"-profiles", "unknown-profile",
			"-seed", "42",
			"-out", filepath.Join(dir, "out"),
		},
		&stdout,
		&stderr,
	)

	if err == nil {
		t.Fatal("expected unknown profile error, got nil")
	}

	if !strings.Contains(err.Error(), "unknown transform profile") {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestRunCorpusImportMWSUnsupportedTask(t *testing.T) {
	dir := t.TempDir()

	metadataPath := filepath.Join(dir, "metadata.jsonl")
	if err := os.WriteFile(metadataPath, []byte(""), 0o600); err != nil {
		t.Fatalf("write metadata: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"corpus",
			"import-mws",
			"-source", dir,
			"-task", "reasoning VQA ru",
			"-out", filepath.Join(dir, "manifest.jsonl"),
		},
		&stdout,
		&stderr,
	)

	if err == nil {
		t.Fatal("expected unsupported task error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"unsupported MWS task",
	) {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestRunCorpusGenerate(t *testing.T) {
	dir := t.TempDir()

	textsPath := filepath.Join(dir, "texts.jsonl")
	fontDir := filepath.Join(dir, "fonts")
	outputDir := filepath.Join(dir, "synthetic")

	if err := os.MkdirAll(fontDir, 0o755); err != nil {
		t.Fatalf("create font dir: %v", err)
	}

	fontPath := filepath.Join(fontDir, "regular.ttf")
	if err := os.WriteFile(fontPath, goregular.TTF, 0o600); err != nil {
		t.Fatalf("write font: %v", err)
	}

	texts := `{"id":"text-002","text":"Hello OCR"}
{"id":"text-001","text":"Счёт № 12345"}
`

	if err := os.WriteFile(textsPath, []byte(texts), 0o600); err != nil {
		t.Fatalf("write texts: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{
			"corpus",
			"generate",
			"-texts", textsPath,
			"-font-dir", fontDir,
			"-pages", "1",
			"-seed", "42",
			"-out", outputDir,
		},
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatalf(
			"corpus generate failed: %v\nstderr: %s",
			err,
			stderr.String(),
		)
	}

	manifestPath := filepath.Join(
		outputDir,
		"manifest.jsonl",
	)

	records, err := corpus.ReadManifest(manifestPath)
	if err != nil {
		t.Fatalf("read generated manifest: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf(
			"expected 1 generated record, got %d",
			len(records),
		)
	}

	record := records[0]

	if record.ID == "" {
		t.Fatal("expected generated id")
	}

	if len(record.References) != 1 {
		t.Fatalf(
			"unexpected references: %#v",
			record.References,
		)
	}

	if record.References[0] != "Счёт № 12345" {
		t.Fatalf(
			"unexpected reference: %q",
			record.References[0],
		)
	}

	if record.Width != 1200 {
		t.Fatalf(
			"unexpected width: %d",
			record.Width,
		)
	}

	if record.Height != 1600 {
		t.Fatalf(
			"unexpected height: %d",
			record.Height,
		)
	}

	if record.Format != "png" {
		t.Fatalf(
			"unexpected format: %q",
			record.Format,
		)
	}

	if record.SHA256 == "" {
		t.Fatal("expected SHA-256")
	}

	if len(record.Tags) != 1 ||
		record.Tags[0] != "synthetic" {
		t.Fatalf(
			"unexpected tags: %#v",
			record.Tags,
		)
	}

	imagePath := filepath.Join(
		outputDir,
		filepath.FromSlash(record.Image),
	)

	if _, err := os.Stat(imagePath); err != nil {
		t.Fatalf(
			"generated image not found: %v",
			err,
		)
	}

	if !strings.Contains(
		stdout.String(),
		"generated records: 1",
	) {
		t.Fatalf(
			"unexpected stdout:\n%s",
			stdout.String(),
		)
	}
}
