package runner

import (
	"context"
	"testing"
)

type fakeRunner struct {
	result Result
}

func (f fakeRunner) Run(
	ctx context.Context,
	task Task,
) Result {
	result := f.result
	result.ID = task.ID

	return result
}

func TestRunnerInterface(t *testing.T) {
	var r Runner = fakeRunner{
		result: Result{
			Text: "Привет, OCR!",
		},
	}

	result := r.Run(
		context.Background(),
		Task{
			ID:        "1",
			ImagePath: "image.png",
		},
	)

	if result.ID != "1" {
		t.Fatalf("expected id %q, got %q", "1", result.ID)
	}

	if result.Text != "Привет, OCR!" {
		t.Fatalf("unexpected text: %q", result.Text)
	}
}
