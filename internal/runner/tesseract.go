package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

type TesseractRunner struct {
	binary    string
	languages string
	psm       int
}

func NewTesseractRunner(
	binary string,
	languages string,
	psm int,
) *TesseractRunner {
	return &TesseractRunner{
		binary:    binary,
		languages: languages,
		psm:       psm,
	}
}

func (r *TesseractRunner) Run(
	ctx context.Context,
	task Task,
) Result {
	start := time.Now()

	result := Result{
		ID: task.ID,
	}

	args := []string{
		task.ImagePath,
		"stdout",
		"-l", r.languages,
		"--psm", strconv.Itoa(r.psm),
	}

	cmd := exec.CommandContext(
		ctx,
		r.binary,
		args...,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf(
			"run tesseract: %v",
			err,
		)
		result.Stderr = stderr.String()
		result.DurationMS = time.Since(start).Milliseconds()

		return result
	}

	result.Text = stdout.String()
	result.Stderr = stderr.String()
	result.DurationMS = time.Since(start).Milliseconds()

	return result
}
