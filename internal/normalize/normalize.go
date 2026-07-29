package normalize

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Profile определяет набор правил нормализации текста.
type Profile string

const (
	// ProfileStrict сохраняет регистр и пунктуацию.
	ProfileStrict Profile = "strict"

	// ProfilePlainTextRU приводит текст к упрощённому виду
	// для сравнения русского OCR-текста.
	ProfilePlainTextRU Profile = "plain-text-ru"
)

// Normalize нормализует текст в соответствии с выбранным профилем.
func Normalize(text string, profile Profile) (string, error) {
	switch profile {
	case ProfileStrict:
		return normalizeStrict(text), nil
	case ProfilePlainTextRU:
		return normalizePlainTextRU(text), nil
	default:
		return "", fmt.Errorf("unknown normalization profile %q", profile)
	}
}

func normalizeStrict(text string) string {
	text = norm.NFC.String(text)
	text = normalizeLineEndings(text)

	lines := strings.Split(text, "\n")

	for i := range lines {
		lines[i] = strings.TrimRightFunc(lines[i], unicode.IsSpace)
	}

	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	return strings.Join(lines[start:end], "\n")
}

func normalizePlainTextRU(text string) string {
	text = norm.NFC.String(text)
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "ё", "е")
	text = normalizeLineEndings(text)

	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		switch {
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			builder.WriteRune(' ')
		case unicode.IsSpace(r):
			builder.WriteRune(' ')
		default:
			builder.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(builder.String()), " ")
}

func normalizeLineEndings(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(text, "\r", "\n")
}
