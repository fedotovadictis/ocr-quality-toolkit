package corpus

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

// ValidateRecords проверяет структурную целостность записей manifest
// и возвращает все найденные ошибки в стабильном порядке.
func ValidateRecords(records []Record) []error {
	var errs []error

	if !isSortedByID(records) {
		errs = append(
			errs,
			fmt.Errorf("manifest is not sorted by id"),
		)
	}

	seenIDs := make(map[string]struct{}, len(records))

	for index, record := range records {
		position := index + 1

		if strings.TrimSpace(record.ID) == "" {
			errs = append(
				errs,
				fmt.Errorf("record %d: empty id", position),
			)
		} else {
			if _, exists := seenIDs[record.ID]; exists {
				errs = append(
					errs,
					fmt.Errorf(
						"record %q: duplicate id",
						record.ID,
					),
				)
			}

			seenIDs[record.ID] = struct{}{}
		}

		imagePath := strings.TrimSpace(record.Image)

		if imagePath == "" {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: empty image path",
					record.ID,
				),
			)
		} else if !isSafeRelativePath(imagePath) {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: unsafe image path %q",
					record.ID,
					record.Image,
				),
			)
		}

		if len(record.References) == 0 || !hasNonEmptyReference(record.References) {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: empty references",
					record.ID,
				),
			)
		}
		for _, reference := range record.References {
			if !utf8.ValidString(reference) {
				errs = append(
					errs,
					fmt.Errorf(
						"record %q: invalid UTF-8 reference",
						record.ID,
					),
				)
			}
		}

		if strings.TrimSpace(record.Language) == "" {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: empty language",
					record.ID,
				),
			)
		}

		if strings.TrimSpace(record.Task) == "" {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: empty task",
					record.ID,
				),
			)
		}

		for _, tagErr := range validateTags(record.ID, record.Tags) {
			errs = append(errs, tagErr)
		}
	}

	sort.SliceStable(
		errs,
		func(i, j int) bool {
			return errs[i].Error() < errs[j].Error()
		},
	)

	return errs
}

func isSortedByID(records []Record) bool {
	for i := 1; i < len(records); i++ {
		if records[i-1].ID > records[i].ID {
			return false
		}
	}

	return true
}

func isSafeRelativePath(path string) bool {
	if path == "" {
		return false
	}

	normalized := strings.ReplaceAll(path, "\\", "/")

	if filepath.IsAbs(path) {
		return false
	}

	if len(normalized) >= 2 && normalized[1] == ':' {
		return false
	}

	clean := filepath.ToSlash(filepath.Clean(normalized))

	if clean == ".." || strings.HasPrefix(clean, "../") {
		return false
	}

	return true
}

func hasNonEmptyReference(references []string) bool {
	for _, reference := range references {
		if strings.TrimSpace(reference) != "" {
			return true
		}
	}

	return false
}

func validateTags(id string, tags []string) []error {
	var errs []error

	seen := make(map[string]struct{}, len(tags))

	for _, tag := range tags {
		if _, exists := seen[tag]; exists {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: duplicate tag %q",
					id,
					tag,
				),
			)
		}

		seen[tag] = struct{}{}
	}

	if !sort.StringsAreSorted(tags) {
		errs = append(
			errs,
			fmt.Errorf(
				"record %q: tags are not sorted",
				id,
			),
		)
	}

	return errs
}
