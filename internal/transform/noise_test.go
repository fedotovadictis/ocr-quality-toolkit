package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestNoiseLight(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 4, 4))

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			source.SetRGBA(x, y, color.RGBA{
				R: 120,
				G: 130,
				B: 140,
				A: 255,
			})
		}
	}

	first := NoiseLight(source, 42)
	second := NoiseLight(source, 42)
	differentSeed := NoiseLight(source, 43)

	if first.Bounds() != source.Bounds() {
		t.Fatalf(
			"expected bounds %v, got %v",
			source.Bounds(),
			first.Bounds(),
		)
	}

	if imagesEqual(source, first) {
		t.Fatal("expected noisy image to differ from source")
	}

	if !imagesEqual(first, second) {
		t.Fatal("same seed must produce identical noise")
	}

	if imagesEqual(first, differentSeed) {
		t.Fatal("different seeds must produce different noise")
	}
}

func imagesEqual(first, second image.Image) bool {
	if first.Bounds() != second.Bounds() {
		return false
	}

	bounds := first.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if first.At(x, y) != second.At(x, y) {
				return false
			}
		}
	}

	return true
}
