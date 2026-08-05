package transform

import (
	"image"
	"image/color"
)

// Grayscale преобразует изображение в оттенки серого,
// сохраняя исходные размеры.
func Grayscale(source image.Image) image.Image {
	bounds := source.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(source.At(x, y))
			result.Set(x, y, gray)
		}
	}

	return result
}
