package metrics

func Similarity(reference, hypothesis []string) float64 {
	maxLength := len(reference)
	if len(hypothesis) > maxLength {
		maxLength = len(hypothesis)
	}

	if maxLength == 0 {
		return 1
	}

	alignment := Align(reference, hypothesis)

	return 1 - float64(alignment.Distance)/float64(maxLength)
}
