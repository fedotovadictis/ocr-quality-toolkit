package transform

import (
	"image"
	"image/color"
)

// LowContrast уменьшает контраст изображения,
// сохраняя исходные размеры.
func LowContrast(source image.Image) image.Image {
	bounds := source.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := color.RGBAModel.Convert(source.At(x, y)).(color.RGBA)

			result.SetRGBA(x, y, color.RGBA{
				R: reduceContrast(pixel.R),
				G: reduceContrast(pixel.G),
				B: reduceContrast(pixel.B),
				A: pixel.A,
			})
		}
	}

	return result
}

func reduceContrast(value uint8) uint8 {
	const midpoint = 128.0
	const factor = 0.5

	adjusted := midpoint + (float64(value)-midpoint)*factor

	if adjusted < 0 {
		return 0
	}
	if adjusted > 255 {
		return 255
	}

	return uint8(adjusted)
}
