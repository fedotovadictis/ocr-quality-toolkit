package generator

import (
	"fmt"
	"os"

	"golang.org/x/image/font/opentype"
)

func LoadFont(path string) (*opentype.Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read font %q: %w", path, err)
	}

	font, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse font %q: %w", path, err)
	}

	return font, nil
}
