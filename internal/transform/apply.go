package transform

import (
	"fmt"
	"image"
)

// Apply применяет преобразование изображения по имени профиля.
func Apply(
	source image.Image,
	profile string,
	seed int64,
) (image.Image, error) {
	switch profile {
	case "clean":
		return source, nil

	case "grayscale":
		return Grayscale(source), nil

	case "low-contrast":
		return LowContrast(source), nil

	case "noise-light":
		return NoiseLight(source, seed), nil

	case "jpeg-70":
		return JPEG70(source)

	case "downscale-50":
		return Downscale50(source), nil

	default:
		return nil, fmt.Errorf(
			"unknown transform profile %q",
			profile,
		)
	}
}
