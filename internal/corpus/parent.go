package corpus

import "fmt"

// ValidateParentIDs проверяет, что все parent_id
// ссылаются на существующие записи manifest.
func ValidateParentIDs(records []Record) error {
	ids := make(map[string]struct{}, len(records))

	for _, record := range records {
		ids[record.ID] = struct{}{}
	}

	for _, record := range records {
		if record.ParentID == "" {
			continue
		}

		if _, ok := ids[record.ParentID]; !ok {
			return fmt.Errorf(
				"parent_id %q not found for record %q",
				record.ParentID,
				record.ID,
			)
		}
	}

	return nil
}
