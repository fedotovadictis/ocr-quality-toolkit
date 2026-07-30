package metrics

func CER(reference, hypothesis string) float64 {
	if len(reference) == 0 {
		if len(hypothesis) == 0 {
			return 0
		}
		return 1
	}

	referenceSymbols := make([]string, 0, len([]rune(reference)))
	for _, r := range reference {
		referenceSymbols = append(referenceSymbols, string(r))
	}
	hypothesisSymbols := make([]string, 0, len([]rune(hypothesis)))
	for _, r := range hypothesis {
		hypothesisSymbols = append(hypothesisSymbols, string(r))
	}

	alignment := Align(referenceSymbols, hypothesisSymbols)
	errors := alignment.Substitutions +
		alignment.Deletions +
		alignment.Insertions

	return float64(errors) / float64(len(referenceSymbols))
}
