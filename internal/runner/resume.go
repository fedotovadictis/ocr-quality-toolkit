package runner

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
