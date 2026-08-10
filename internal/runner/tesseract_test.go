package runner

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
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
func TestTesseractRunnerCancellation(t *testing.T) {
	dir := t.TempDir()

	sourcePath := filepath.Join(dir, "fake_tesseract.go")

	source := `package main

import "time"

func main() {
	time.Sleep(10 * time.Second)
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()

	result := r.Run(
		ctx,
		Task{
			ID:        "page-timeout",
			ImagePath: "image.png",
		},
	)

	if result.Error == "" {
		t.Fatal("expected cancellation error, got empty string")
	}
	if !strings.Contains(result.Error, "context deadline exceeded") {
		t.Fatalf(
			"expected deadline error, got %q",
			result.Error,
		)
	}
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf(
			"expected DeadlineExceeded, got %v",
			ctx.Err(),
		)
	}

	if result.ID != "page-timeout" {
		t.Fatalf(
			"expected id %q, got %q",
			"page-timeout",
			result.ID,
		)
	}
}
