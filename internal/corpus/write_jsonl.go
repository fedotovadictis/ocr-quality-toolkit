package corpus

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJSONL[T any](path string, values []T) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for _, value := range values {
		if err := encoder.Encode(value); err != nil {
			return fmt.Errorf("write JSONL record: %w", err)
		}
	}

	return nil
}
