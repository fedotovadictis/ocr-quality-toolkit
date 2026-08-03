package generator

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

// SavePNG сохраняет изображение в формате PNG.
func SavePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create PNG %q: %w", path, err)
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		return fmt.Errorf("encode PNG %q: %w", path, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close PNG %q: %w", path, err)
	}

	return nil
}
