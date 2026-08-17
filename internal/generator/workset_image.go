package generator

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"ocr-quality-toolkit/internal/transform"
)

// BuildSyntheticImage читает исходное изображение,
// применяет профиль преобразования и сохраняет результат в PNG
func BuildSyntheticImage(
	sourcePath string,
	targetPath string,
	profile string,
	seed int64,
) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source image: %w", err)
	}
	defer file.Close()

	source, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("decode source image: %w", err)
	}

	result, err := transform.Apply(source, profile, seed)
	if err != nil {
		return fmt.Errorf("apply transform: %w", err)
	}

	if err := SavePNG(targetPath, result); err != nil {
		return fmt.Errorf("save synthetic image: %w", err)
	}

	return nil
}
