package runner

import (
	"testing"
)

func TestFilterPendingTasks(t *testing.T) {
	tasks := []Task{
		{ID: "1", ImagePath: "1.png"},
		{ID: "2", ImagePath: "2.png"},
		{ID: "3", ImagePath: "3.png"},
	}

	existing := []Result{
		{ID: "1", Text: "done"},
		{ID: "3", Error: "ocr error"},
	}

	pending := FilterPendingTasks(tasks, existing)

	if len(pending) != 1 {
		t.Fatalf(
			"expected 1 pending task, got %d",
			len(pending),
		)
	}

	if pending[0].ID != "2" {
		t.Fatalf(
			"expected task %q, got %q",
			"2",
			pending[0].ID,
		)
	}
}
