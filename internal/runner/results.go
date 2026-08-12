package runner

import (
	"encoding/json"
	"fmt"
	"os"
)

// AppendResult добавляет один OCR-результат в JSONL-файл.
func AppendResult(path string, result Result) error {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return fmt.Errorf("open results %q: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encode result %q: %w", result.ID, err)
	}

	return nil
}
