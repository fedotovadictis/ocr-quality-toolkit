package generator

import (
	"fmt"
	"sort"
	"strings"

	"ocr-quality-toolkit/internal/corpus"
)

// TextInput описывает одну исходную текстовую запись
// для генерации синтетической страницы.
type TextInput struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// ReadTextInputs читает JSONL с исходными текстами
// и возвращает записи в стабильном порядке по ID.
func ReadTextInputs(path string) ([]TextInput, error) {
	inputs, err := corpus.ReadJSONL[TextInput](path)
	if err != nil {
		return nil, fmt.Errorf("read generation texts: %w", err)
	}

	seen := make(map[string]struct{}, len(inputs))

	for i, input := range inputs {
		if strings.TrimSpace(input.ID) == "" {
			return nil, fmt.Errorf(
				"text record %d: empty id",
				i+1,
			)
		}

		if _, exists := seen[input.ID]; exists {
			return nil, fmt.Errorf(
				"text record %d: duplicate id %q",
				i+1,
				input.ID,
			)
		}

		seen[input.ID] = struct{}{}
	}

	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].ID < inputs[j].ID
	})

	return inputs, nil
}
