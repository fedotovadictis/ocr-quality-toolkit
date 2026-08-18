package corpus

import "sort"

// ValidateCorpus выполняет полную проверку корпуса.
func ValidateCorpus(root string, records []Record) []error {
	var errs []error

	errs = append(errs, ValidateRecords(records)...)
	errs = append(errs, ValidateRecordFiles(root, records)...)

	if err := ValidateParentIDs(records); err != nil {
		errs = append(errs, err)
	}

	sort.SliceStable(
		errs,
		func(i, j int) bool {
			return errs[i].Error() < errs[j].Error()
		},
	)

	return errs
}
