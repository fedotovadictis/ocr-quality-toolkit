package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestTesseractRunnerSuccess(t *testing.T) {
	dir := t.TempDir()

	sourcePath := filepath.Join(dir, "fake_tesseract.go")

	source := `package main

import "fmt"

func main() {
	fmt.Print("Распознанный текст\n")
}
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write fake source: %v", err)
	}

	binaryName := "fake-tesseract"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	binaryPath := filepath.Join(dir, binaryName)

	build := exec.Command(
		"go",
		"build",
		"-o",
		binaryPath,
		sourcePath,
	)

	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf(
			"build fake tesseract: %v\n%s",
			err,
			output,
		)
	}

	r := NewTesseractRunner(
		binaryPath,
		"rus+eng",
		3,
	)

	result := r.Run(
		context.Background(),
		Task{
			ID:        "page-001",
			ImagePath: "image.png",
		},
	)

	if result.Error != "" {
		t.Fatalf("unexpected error: %s", result.Error)
	}

	if result.ID != "page-001" {
		t.Fatalf(
			"expected id %q, got %q",
			"page-001",
			result.ID,
		)
	}

	if result.Text != "Распознанный текст\n" {
		t.Fatalf(
			"unexpected text: %q",
			result.Text,
		)
	}
	if result.DurationMS < 0 {
		t.Fatalf(
			"expected non-negative duration, got %d",
			result.DurationMS,
		)
	}
}

func TestTesseractRunnerError(t *testing.T) {
	dir := t.TempDir()

	sourcePath := filepath.Join(dir, "fake_tesseract.go")

	source := `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprint(os.Stderr, "fake tesseract error")
	os.Exit(1)
}
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write fake source: %v", err)
	}

	binaryName := "fake-tesseract"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	binaryPath := filepath.Join(dir, binaryName)

	build := exec.Command(
		"go",
		"build",
		"-o",
		binaryPath,
		sourcePath,
	)

	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf(
			"build fake tesseract: %v\n%s",
			err,
			output,
		)
	}

	r := NewTesseractRunner(
		binaryPath,
		"rus+eng",
		3,
	)

	result := r.Run(
		context.Background(),
		Task{
			ID:        "page-error",
			ImagePath: "image.png",
		},
	)

	if result.Error == "" {
		t.Fatal("expected error, got empty string")
	}

	if result.Stderr != "fake tesseract error" {
		t.Fatalf(
			"unexpected stderr: %q",
			result.Stderr,
		)
	}

	if result.ID != "page-error" {
		t.Fatalf(
			"expected id %q, got %q",
			"page-error",
			result.ID,
		)
	}
}
