package main

import (
	"errors"
	"fmt"
	"os"

	"ocr-quality-toolkit/internal/app"
)

func main() {
	err := app.Run(
		os.Args[1:],
		os.Stdout,
		os.Stderr,
	)

	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)

	if errors.Is(err, app.ErrRegression) {
		os.Exit(1)
	}

	os.Exit(2)
}
