package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"ocr-quality-toolkit/internal/generator"
	"os"
	"path/filepath"

	"ocr-quality-toolkit/internal/corpus"
	"ocr-quality-toolkit/internal/evaluate"
	"ocr-quality-toolkit/internal/normalize"
	"ocr-quality-toolkit/internal/report"
	"ocr-quality-toolkit/internal/runner"
)

var ErrRegression = errors.New("quality regression detected")

func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New(
			"usage: ocrq <import-mws|evaluate|run-tesseract|compare> [options]",
		)
	}

	switch args[0] {
	case "evaluate":
		return runEvaluate(args[1:], stdout, stderr)
	case "import-mws":
		return runImportMWS(args[1:], stdout, stderr)
	case "build-workset":
		return runBuildWorkset(args[1:], stdout, stderr)
	case "run-tesseract":
		return runTesseract(args[1:], stdout, stderr)
	case "compare":
		return runCompare(args[1:], stdout, stderr)
	case "report":
		return runReport(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runImportMWS(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("import-mws", flag.ContinueOnError)
	flags.SetOutput(stderr)

	metadataPath := flags.String(
		"metadata",
		"",
		"path to MWS metadata JSONL",
	)

	outputPath := flags.String(
		"output",
		"",
		"path to output manifest JSONL",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *metadataPath == "" {
		return errors.New("missing required flag: -metadata")
	}

	if *outputPath == "" {
		return errors.New("missing required flag: -output")
	}

	records, stats, err := corpus.ImportMWSMetadata(*metadataPath)
	if err != nil {
		return fmt.Errorf("import MWS metadata: %w", err)
	}

	if err := corpus.WriteJSONL(*outputPath, records); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}

	fmt.Fprintf(stdout, "total lines: %d\n", stats.TotalLines)
	fmt.Fprintf(stdout, "matching tasks: %d\n", stats.MatchingTasks)
	fmt.Fprintf(stdout, "imported: %d\n", stats.Imported)
	fmt.Fprintf(stdout, "missing images: %d\n", stats.MissingImages)
	fmt.Fprintf(stdout, "invalid images: %d\n", stats.InvalidImages)
	fmt.Fprintf(stdout, "empty references: %d\n", stats.EmptyReferences)

	return nil
}

func runEvaluate(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("evaluate", flag.ContinueOnError)
	flags.SetOutput(stderr)

	manifestPath := flags.String(
		"manifest",
		"",
		"path to manifest JSONL",
	)

	hypothesesPath := flags.String(
		"hypotheses",
		"",
		"path to hypotheses JSONL",
	)

	normalization := flags.String(
		"normalization",
		"",
		"normalization profile: strict or plain-text-ru",
	)

	outputPath := flags.String(
		"out",
		"",
		"path to output JSON report",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *manifestPath == "" {
		return errors.New("missing required flag: -manifest")
	}

	if *hypothesesPath == "" {
		return errors.New("missing required flag: -hypotheses")
	}

	if *normalization == "" {
		return errors.New("missing required flag: -normalization")
	}

	if *outputPath == "" {
		return errors.New("missing required flag: -out")
	}

	profile := normalize.Profile(*normalization)

	records, err := corpus.ReadManifest(*manifestPath)
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}

	hypotheses, err := corpus.ReadHypotheses(*hypothesesPath)
	if err != nil {
		return fmt.Errorf("read hypotheses: %w", err)
	}

	results, err := evaluate.Evaluate(
		records,
		hypotheses,
		profile,
	)
	if err != nil {
		return fmt.Errorf("evaluate: %w", err)
	}

	fullReport := report.BuildReport(
		records,
		results,
	)

	if err := report.WriteJSON(*outputPath, fullReport); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	fmt.Fprintf(stdout, "evaluated: %d\n", len(results))
	fmt.Fprintf(stdout, "report: %s\n", *outputPath)

	return nil
}

func runTesseract(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("run-tesseract", flag.ContinueOnError)
	flags.SetOutput(stderr)

	manifestPath := flags.String(
		"manifest",
		"",
		"path to manifest JSONL",
	)

	binaryPath := flags.String(
		"binary",
		"tesseract",
		"path to tesseract binary",
	)

	languages := flags.String(
		"languages",
		"rus+eng",
		"tesseract languages",
	)

	psm := flags.Int(
		"psm",
		3,
		"tesseract page segmentation mode",
	)

	workers := flags.Int(
		"workers",
		1,
		"number of parallel OCR workers",
	)

	resume := flags.Bool(
		"resume",
		false,
		"continue OCR run using existing results",
	)

	outputPath := flags.String(
		"out",
		"",
		"path to output results JSONL",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *manifestPath == "" {
		return errors.New("missing required flag: -manifest")
	}

	if *outputPath == "" {
		return errors.New("missing required flag: -out")
	}

	if *workers < 1 {
		return errors.New("workers must be greater than zero")
	}

	records, err := corpus.ReadManifest(*manifestPath)
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}

	tasks := make([]runner.Task, 0, len(records))

	for _, record := range records {
		tasks = append(tasks, runner.Task{
			ID:        record.ID,
			ImagePath: record.Image,
		})
	}

	if *resume {
		existing, err := runner.ReadResults(*outputPath)
		if err == nil {
			tasks = runner.FilterPendingTasks(
				tasks,
				existing,
			)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf(
				"read existing results: %w",
				err,
			)
		}
	}

	tesseract := runner.NewTesseractRunner(
		*binaryPath,
		*languages,
		*psm,
	)

	results, err := runner.RunTasksAndSave(
		context.Background(),
		tesseract,
		tasks,
		*workers,
		*outputPath,
	)
	if err != nil {
		return fmt.Errorf("run OCR: %w", err)
	}

	fmt.Fprintf(stdout, "processed: %d\n", len(results))

	return nil
}

func runCompare(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("compare", flag.ContinueOnError)
	flags.SetOutput(stderr)

	baselinePath := flags.String(
		"baseline",
		"",
		"path to baseline JSON report",
	)

	currentPath := flags.String(
		"current",
		"",
		"path to current JSON report",
	)

	maxCERIncrease := flags.Float64(
		"max-cer-increase",
		0,
		"maximum allowed CER increase",
	)

	maxWERIncrease := flags.Float64(
		"max-wer-increase",
		0,
		"maximum allowed WER increase",
	)

	maxCoverageDecrease := flags.Float64(
		"max-coverage-decrease",
		0,
		"maximum allowed coverage decrease",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *baselinePath == "" {
		return errors.New("missing required flag: -baseline")
	}

	if *currentPath == "" {
		return errors.New("missing required flag: -current")
	}

	baseline, err := report.ReadJSON(*baselinePath)
	if err != nil {
		return fmt.Errorf("read baseline: %w", err)
	}

	current, err := report.ReadJSON(*currentPath)
	if err != nil {
		return fmt.Errorf("read current: %w", err)
	}

	comparison := report.CompareReports(
		baseline,
		current,
	)

	fmt.Fprintf(
		stdout,
		"CER delta: %.6f\n",
		comparison.CERDelta,
	)

	fmt.Fprintf(
		stdout,
		"WER delta: %.6f\n",
		comparison.WERDelta,
	)

	fmt.Fprintf(
		stdout,
		"coverage delta: %.6f\n",
		comparison.CoverageDelta,
	)

	thresholds := report.Thresholds{
		MaxCERIncrease:      *maxCERIncrease,
		MaxWERIncrease:      *maxWERIncrease,
		MaxCoverageDecrease: *maxCoverageDecrease,
	}

	if report.HasRegression(
		comparison,
		thresholds,
	) {
		return ErrRegression
	}

	return nil
}
func runBuildWorkset(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("build-workset", flag.ContinueOnError)
	flags.SetOutput(stderr)

	root := flags.String(
		"root",
		"",
		"workset root directory",
	)

	realManifestPath := flags.String(
		"real-manifest",
		"",
		"path to real records manifest JSONL",
	)

	syntheticManifestPath := flags.String(
		"synthetic-manifest",
		"",
		"path to output synthetic manifest JSONL",
	)

	outputPath := flags.String(
		"out",
		"",
		"path to combined workset manifest JSONL",
	)

	profile := flags.String(
		"profile",
		"grayscale",
		"synthetic transform profile",
	)

	seed := flags.Int64(
		"seed",
		42,
		"deterministic transform seed",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *root == "" {
		return errors.New("missing required flag: -root")
	}

	if *realManifestPath == "" {
		return errors.New("missing required flag: -real-manifest")
	}

	if *syntheticManifestPath == "" {
		return errors.New("missing required flag: -synthetic-manifest")
	}

	if *outputPath == "" {
		return errors.New("missing required flag: -out")
	}

	realRecords, err := corpus.ReadManifest(*realManifestPath)
	if err != nil {
		return fmt.Errorf("read real manifest: %w", err)
	}

	if err := corpus.ValidateCorpusIntegrity(realRecords); err != nil {
		return fmt.Errorf("validate real corpus: %w", err)
	}

	if err := corpus.ValidateImageFiles(*root, realRecords); err != nil {
		return fmt.Errorf("validate real images: %w", err)
	}

	syntheticDir := filepath.Join(*root, "synthetic")

	if err := os.MkdirAll(syntheticDir, 0o755); err != nil {
		return fmt.Errorf("create synthetic directory: %w", err)
	}

	syntheticRecords, err := generator.BuildSyntheticWorkset(
		*root,
		realRecords,
		*profile,
		*seed,
	)
	if err != nil {
		return fmt.Errorf("build synthetic workset: %w", err)
	}

	if err := corpus.WriteJSONL(
		*syntheticManifestPath,
		syntheticRecords,
	); err != nil {
		return fmt.Errorf("write synthetic manifest: %w", err)
	}

	workset, err := corpus.BuildWorkset(
		realRecords,
		syntheticRecords,
	)
	if err != nil {
		return fmt.Errorf("build combined workset: %w", err)
	}

	if err := corpus.ValidateImageFiles(*root, workset); err != nil {
		return fmt.Errorf("validate workset images: %w", err)
	}

	if err := corpus.WriteJSONL(*outputPath, workset); err != nil {
		return fmt.Errorf("write workset manifest: %w", err)
	}

	fmt.Fprintf(stdout, "real records: %d\n", len(realRecords))
	fmt.Fprintf(stdout, "synthetic records: %d\n", len(syntheticRecords))
	fmt.Fprintf(stdout, "total records: %d\n", len(workset))
	fmt.Fprintf(stdout, "manifest: %s\n", *outputPath)

	return nil
}
func runReport(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("report", flag.ContinueOnError)
	flags.SetOutput(stderr)

	inputPath := flags.String(
		"input",
		"",
		"path to input JSON report",
	)

	format := flags.String(
		"format",
		"html",
		"report format",
	)

	outputPath := flags.String(
		"out",
		"",
		"path to output report",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *inputPath == "" {
		return errors.New("missing required flag: -input")
	}

	if *outputPath == "" {
		return errors.New("missing required flag: -out")
	}

	if *format != "html" {
		return fmt.Errorf("unsupported report format %q", *format)
	}

	inputReport, err := report.ReadJSON(*inputPath)
	if err != nil {
		return fmt.Errorf("read report: %w", err)
	}

	if err := report.WriteHTML(*outputPath, inputReport); err != nil {
		return fmt.Errorf("write HTML report: %w", err)
	}

	fmt.Fprintf(stdout, "report: %s\n", *outputPath)

	return nil
}
