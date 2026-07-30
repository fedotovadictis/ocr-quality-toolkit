package metrics

import "strings"

func WER(reference, hypothesis string) float64 {
	referenceWords := strings.Fields(reference)
	hypothesisWords := strings.Fields(hypothesis)

	if len(referenceWords) == 0 {
		if len(hypothesisWords) == 0 {
			return 0
		}
		return 1
	}

	alignment := Align(referenceWords, hypothesisWords)

	errors := alignment.Substitutions +
		alignment.Deletions +
		alignment.Insertions

	return float64(errors) / float64(len(referenceWords))

}
