package corpus

import "fmt"

// BuildWorkset объединяет реальные и синтетические записи
// в единый рабочий корпус и проверяет его целостность.
func BuildWorkset(real, synthetic []Record) ([]Record, error) {
	records := make([]Record, 0, len(real)+len(synthetic))
	records = append(records, real...)
	records = append(records, synthetic...)

	if err := ValidateCorpusIntegrity(records); err != nil {
		return nil, fmt.Errorf("validate workset: %w", err)
	}

	return records, nil
}
