package runner

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

type countingRunner struct {
	mu        sync.Mutex
	active    int
	maxActive int
	seen      map[string]int
}

func (r *countingRunner) Run(
	ctx context.Context,
	task Task,
) Result {
	r.mu.Lock()

	r.active++

	if r.active > r.maxActive {
		r.maxActive = r.active
	}

	r.seen[task.ID]++

	r.mu.Unlock()

	time.Sleep(20 * time.Millisecond)

	r.mu.Lock()
	r.active--
	r.mu.Unlock()

	return Result{
		ID:   task.ID,
		Text: task.ID,
	}
}

func TestRunTasksWithWorkers(t *testing.T) {
	r := &countingRunner{
		seen: make(map[string]int),
	}

	tasks := make([]Task, 10)

	for i := range tasks {
		tasks[i] = Task{
			ID: fmt.Sprintf("task-%02d", i),
		}
	}

	results := RunTasks(
		context.Background(),
		r,
		tasks,
		3,
	)

	if len(results) != len(tasks) {
		t.Fatalf(
			"expected %d results, got %d",
			len(tasks),
			len(results),
		)
	}

	if r.maxActive > 3 {
		t.Fatalf(
			"expected at most 3 concurrent workers, got %d",
			r.maxActive,
		)
	}

	if r.maxActive < 2 {
		t.Fatalf(
			"expected parallel execution, max active=%d",
			r.maxActive,
		)
	}

	for _, task := range tasks {
		if r.seen[task.ID] != 1 {
			t.Fatalf(
				"task %q processed %d times",
				task.ID,
				r.seen[task.ID],
			)
		}
	}
}
func TestRunTasksCancellation(t *testing.T) {
	r := &countingRunner{
		seen: make(map[string]int),
	}

	tasks := make([]Task, 100)

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

	results := RunTasks(
		ctx,
		r,
		tasks,
		2,
	)

	if len(results) >= len(tasks) {
		t.Fatalf(
			"expected cancellation before all tasks completed, got %d results",
			len(results),
		)
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Fatalf(
			"expected DeadlineExceeded, got %v",
			ctx.Err(),
		)
	}
}
