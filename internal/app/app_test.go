package app

import (
	"bytes"
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
