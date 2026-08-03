package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"ocr-quality-toolkit/internal/evaluate"
	"ocr-quality-toolkit/internal/normalize"

	"ocr-quality-toolkit/internal/corpus"
)

func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New(
			"usage: ocrq <import-mws|evaluate> [options]",
		)
	}

	switch args[0] {
	case "evaluate":
		return runEvaluate(args[1:], stdout, stderr)
	case "import-mws":
		return runImportMWS(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}

}

func runImportMWS(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("import-mws", flag.ContinueOnError)
	flags.SetOutput(stderr)

	metadataPath := flags.String("metadata", "", "path to MWS metadata JSONL")
	outputPath := flags.String("output", "", "path to output manifest JSONL")

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

	if err := evaluate.WriteReport(*outputPath, results); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	fmt.Fprintf(stdout, "evaluated: %d\n", len(results))
	fmt.Fprintf(stdout, "report: %s\n", *outputPath)

	return nil
}
