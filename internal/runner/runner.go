package runner

import "context"

// Task описывает одну задачу OCR.
type Task struct {
	ID        string
	ImagePath string
}

// Result описывает результат OCR одной задачи.
type Result struct {
	ID         string
	Text       string
	Error      string
	Stderr     string
	DurationMS int64
}

// Runner выполняет OCR одной задачи.
type Runner interface {
	Run(ctx context.Context, task Task) Result
}
