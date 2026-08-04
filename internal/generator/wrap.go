package generator

import (
	"strings"

	"golang.org/x/image/font"
)

func WrapText(
	text string,
	face font.Face,
	maxWidth int,
) []string {
	if text == "" {
		return []string{""}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0)
	currentLine := words[0]

	for _, word := range words[1:] {
		candidate := currentLine + " " + word
		width := font.MeasureString(face, candidate).Ceil()

		if width <= maxWidth {
			currentLine = candidate
			continue
		}

		lines = append(lines, currentLine)
		currentLine = word
	}

	lines = append(lines, currentLine)

	return lines
}
