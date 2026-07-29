package metrics

// Distance возвращает расстояние Левенштейна между двумя строками.
//
// Расстояние рассчитывается по Unicode-символам, а не по байтам.
// Допустимые операции: вставка, удаление и замена.
func Distance(reference, hypothesis string) int {
	referenceRunes := []rune(reference)
	hypothesisRunes := []rune(hypothesis)

	if len(referenceRunes) == 0 {
		return len(hypothesisRunes)
	}

	if len(hypothesisRunes) == 0 {
		return len(referenceRunes)
	}

	previous := make([]int, len(hypothesisRunes)+1)
	current := make([]int, len(hypothesisRunes)+1)

	for j := range previous {
		previous[j] = j
	}

	for i, referenceRune := range referenceRunes {
		current[0] = i + 1

		for j, hypothesisRune := range hypothesisRunes {
			substitutionCost := 0
			if referenceRune != hypothesisRune {
				substitutionCost = 1
			}

			deletion := previous[j+1] + 1
			insertion := current[j] + 1
			substitution := previous[j] + substitutionCost

			current[j+1] = min3(
				deletion,
				insertion,
				substitution,
			)
		}

		previous, current = current, previous
	}

	return previous[len(hypothesisRunes)]
}

func min3(first, second, third int) int {
	result := first

	if second < result {
		result = second
	}

	if third < result {
		result = third
	}

	return result
}
