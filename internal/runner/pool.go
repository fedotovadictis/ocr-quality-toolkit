package runner

import (
	"context"
	"sync"
)

func RunTasks(
	ctx context.Context,
	runner Runner,
	tasks []Task,
	workers int,
) []Result {
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
				select {
				case <-ctx.Done():
					return
				default:
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
		results = append(results, result)
	}

	return results
}
