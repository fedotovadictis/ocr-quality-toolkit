package runner

import (
	"context"
	"fmt"
	"sync"
)

// RunTasksAndSave запускает OCR-задачи параллельно
// и сохраняет каждый результат сразу после обработки.
func RunTasksAndSave(
	ctx context.Context,
	runner Runner,
	tasks []Task,
	workers int,
	path string,
) ([]Result, error) {
	if workers < 1 {
		workers = 1
	}

	taskCh := make(chan Task)
	resultCh := make(chan Result)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range taskCh {
				if ctx.Err() != nil {
					return
				}

				result := runner.Run(ctx, task)

				select {
				case resultCh <- result:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		defer close(taskCh)

		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := make([]Result, 0, len(tasks))

	for result := range resultCh {
		if err := AppendResult(path, result); err != nil {
			return results, fmt.Errorf(
				"save result %q: %w",
				result.ID,
				err,
			)
		}

		results = append(results, result)
	}

	return results, nil
}
