package corpus

import "strings"

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func appendUnique(values []string, additions ...string) []string {
	for _, addition := range additions {
		if addition == "" {
			continue
		}

		if containsString(values, addition) {
			continue
		}

		values = append(values, addition)
	}

	return values
}

func nonEmptyUniqueStrings(values []string) []string {
	result := make([]string, 0, len(values))

	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}

		result = appendUnique(result, value)
	}

	return result
}
