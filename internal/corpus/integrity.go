package corpus

import (
	"fmt"
	"strings"
)

// ValidateCorpusIntegrity проверяет базовую целостность корпуса.
func ValidateCorpusIntegrity(records []Record) error {
	seenIDs := make(map[string]struct{}, len(records))

	for i, record := range records {
		if strings.TrimSpace(record.ID) == "" {
			return fmt.Errorf("record %d: empty id", i+1)
		}

		if _, exists := seenIDs[record.ID]; exists {
			return fmt.Errorf("record %d: duplicate id %q", i+1, record.ID)
		}

		seenIDs[record.ID] = struct{}{}

		if strings.TrimSpace(record.Image) == "" {
			return fmt.Errorf(
				"record %q: empty image path",
				record.ID,
			)
		}

		if len(record.References) == 0 {
			return fmt.Errorf(
				"record %q: no references",
				record.ID,
			)
		}

		hasReference := false

		for _, reference := range record.References {
			if strings.TrimSpace(reference) != "" {
				hasReference = true
				break
			}
		}

		if !hasReference {
			return fmt.Errorf(
				"record %q: empty references",
				record.ID,
			)
		}

		if strings.TrimSpace(record.Language) == "" {
			return fmt.Errorf(
				"record %q: empty language",
				record.ID,
			)
		}

		if strings.TrimSpace(record.Task) == "" {
			return fmt.Errorf(
				"record %q: empty task",
				record.ID,
			)
		}
	}

	if err := ValidateParentIDs(records); err != nil {
		return fmt.Errorf("validate parent ids: %w", err)
	}

	return nil
}
