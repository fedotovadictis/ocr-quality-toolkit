package runner

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

type staticRunner struct{}

func (staticRunner) Run(
	ctx context.Context,
	task Task,
) Result {
	return Result{
		ID:   task.ID,
		Text: "done",
	}
}

func TestRunTasksAndSave(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"results.jsonl",
	)

	tasks := []Task{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	results, err := RunTasksAndSave(
		context.Background(),
		staticRunner{},
		tasks,
		2,
		path,
	)

	if err != nil {
		t.Fatalf("RunTasksAndSave returned error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf(
			"expected 3 results, got %d",
			len(results),
		)
	}

	saved, err := ReadResults(path)
	if err != nil {
		t.Fatalf("ReadResults: %v", err)
	}

	if len(saved) != 3 {
		t.Fatalf(
			"expected 3 saved results, got %d",
			len(saved),
		)
	}
}
func TestRunTasksAndSaveKeepsPartialResultsOnCancellation(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"results.jsonl",
	)

	r := &countingRunner{
		seen: make(map[string]int),
	}

	tasks := make([]Task, 20)

	for i := range tasks {
		tasks[i] = Task{
			ID: fmt.Sprintf("task-%02d", i),
		}
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		50*time.Millisecond,
	)
	defer cancel()

	results, err := RunTasksAndSave(
		ctx,
		r,
		tasks,
		2,
		path,
	)

	if err != nil {
		t.Fatalf("RunTasksAndSave returned error: %v", err)
	}

	if len(results) >= len(tasks) {
		t.Fatalf(
			"expected partial results after cancellation, got %d",
			len(results),
		)
	}

	saved, err := ReadResults(path)
	if err != nil {
		t.Fatalf("ReadResults: %v", err)
	}

	if len(saved) != len(results) {
		t.Fatalf(
			"expected %d saved results, got %d",
			len(results),
			len(saved),
		)
	}

	if len(saved) == 0 {
		t.Fatal("expected at least one partial result to be saved")
	}
}

type mixedRunner struct{}

func (mixedRunner) Run(
	ctx context.Context,
	task Task,
) Result {
	if task.ID == "page-error" {
		return Result{
			ID:    task.ID,
			Error: "ocr failed",
		}
	}

	return Result{
		ID:   task.ID,
		Text: "done",
	}
}

func TestRunTasksAndSaveContinuesAfterError(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"results.jsonl",
	)

	tasks := []Task{
		{ID: "page-001"},
		{ID: "page-error"},
		{ID: "page-003"},
	}

	results, err := RunTasksAndSave(
		context.Background(),
		mixedRunner{},
		tasks,
		2,
		path,
	)

	if err != nil {
		t.Fatalf("RunTasksAndSave returned error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf(
			"expected 3 results, got %d",
			len(results),
		)
	}

	saved, err := ReadResults(path)
	if err != nil {
		t.Fatalf("ReadResults: %v", err)
	}

	if len(saved) != 3 {
		t.Fatalf(
			"expected 3 saved results, got %d",
			len(saved),
		)
	}

	foundError := false
	foundSuccess := 0

	for _, result := range saved {
		if result.ID == "page-error" {
			if result.Error != "ocr failed" {
				t.Fatalf(
					"unexpected error result: %q",
					result.Error,
				)
			}

			foundError = true
			continue
		}

		if result.Text == "done" {
			foundSuccess++
		}
	}

	if !foundError {
		t.Fatal("expected OCR error result to be saved")
	}

	if foundSuccess != 2 {
		t.Fatalf(
			"expected 2 successful results, got %d",
			foundSuccess,
		)
	}
}
