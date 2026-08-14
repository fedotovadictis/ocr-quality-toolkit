package report

import (
	"strings"

	"ocr-quality-toolkit/internal/metrics"
)

type AlignmentItem struct {
	Type       string
	Reference  string
	Hypothesis string
}

func BuildAlignment(
	reference string,
	hypothesis string,
) []AlignmentItem {
	referenceWords := strings.Fields(reference)
	hypothesisWords := strings.Fields(hypothesis)

	alignment := metrics.Align(
		referenceWords,
		hypothesisWords,
	)

	items := make([]AlignmentItem, 0, len(alignment.Operations))

	for _, operation := range alignment.Operations {
		items = append(items, AlignmentItem{
			Type:       string(operation.Type),
			Reference:  operation.Reference,
			Hypothesis: operation.Hypothesis,
		})
	}

	return items
}
