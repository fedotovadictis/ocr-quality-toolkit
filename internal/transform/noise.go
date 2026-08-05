package transform

import (
	"image"
	"image/color"
	"math/rand"
)

// NoiseLight добавляет небольшой случайный шум к каждому RGB-каналу.
func NoiseLight(source image.Image, seed int64) image.Image {
	bounds := source.Bounds()
	result := image.NewRGBA(bounds)

	random := rand.New(rand.NewSource(seed))

	const maxNoise = 10

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := color.RGBAModel.Convert(source.At(x, y)).(color.RGBA)

			result.SetRGBA(x, y, color.RGBA{
				R: addNoise(pixel.R, random, maxNoise),
				G: addNoise(pixel.G, random, maxNoise),
				B: addNoise(pixel.B, random, maxNoise),
				A: pixel.A,
			})
		}
	}

	return result
}

func addNoise(value uint8, random *rand.Rand, maxNoise int) uint8 {
	delta := random.Intn(2*maxNoise+1) - maxNoise

	result := int(value) + delta

	if result < 0 {
		result = 0
	}
	if result > 255 {
		result = 255
	}

	return uint8(result)
}
