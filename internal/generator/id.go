package generator

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strings"
)

func MakePageID(
	sourceID string,
	text string,
	options PageOptions,
) string {
	fontName := filepath.Base(
		strings.ReplaceAll(options.FontPath, `\`, `/`),
	)

	input := fmt.Sprintf(
		"%s|%s|%d|%d|%d|%s|%.4f|%d|%d",
		sourceID,
		text,
		options.Width,
		options.Height,
		options.Margin,
		fontName,
		options.FontSize,
		options.LineHeight,
		options.Seed,
	)

	hash := sha256.Sum256([]byte(input))

	return fmt.Sprintf(
		"synthetic-%s-%x",
		sourceID,
		hash[:6],
	)
}
