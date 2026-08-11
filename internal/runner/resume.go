package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// FilterPendingTasks возвращает только те задачи,
// для которых ещё нет результата предыдущего запуска.
func FilterPendingTasks(
	tasks []Task,
	existing []Result,
) []Task {
	done := make(map[string]struct{}, len(existing))

	for _, result := range existing {
		done[result.ID] = struct{}{}
	}

	pending := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		if _, ok := done[task.ID]; ok {
			continue
		}

		pending = append(pending, task)
	}

	return pending
}
func ReadResults(path string) ([]Result, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open results %q: %w", path, err)
	}
	defer file.Close()

	var results []Result

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++

		var result Result

		if err := json.Unmarshal(scanner.Bytes(), &result); err != nil {
			return nil, fmt.Errorf(
				"results line %d: %w",
				lineNum,
				err,
			)
		}

		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read results %q: %w", path, err)
	}

	return results, nil
}
