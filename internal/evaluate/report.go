package evaluate

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteReport(path string, results []Result) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create report %q: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("encode report: %w", err)
	}

	return nil
}
