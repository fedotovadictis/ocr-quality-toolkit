package corpus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func ReadJSONL[T any](path string) ([]T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []T
	lineNum := 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineNum++

		var record T
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
func ReadManifest(path string) ([]Record, error) {
	return ReadJSONL[Record](path)
}

func ReadHypotheses(path string) ([]Hypothesis, error) {
	return ReadJSONL[Hypothesis](path)
}
