package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestDownscale50(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 100, 80))

	// Контрастный узор поможет увидеть эффект масштабирования.
	for y := 0; y < 80; y++ {
		for x := 0; x < 100; x++ {
			if (x+y)%2 == 0 {
				source.SetRGBA(x, y, color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 255,
				})
			} else {
				source.SetRGBA(x, y, color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 255,
				})
			}
		}
	}

	first := Downscale50(source)
	second := Downscale50(source)

	if first.Bounds() != source.Bounds() {
		t.Fatalf(
			"expected bounds %v, got %v",
			source.Bounds(),
			first.Bounds(),
		)
	}

	if imagesEqual(source, first) {
		t.Fatal("expected downscaled image to differ from source")
	}

	if !imagesEqual(first, second) {
		t.Fatal("expected repeated transformation to be deterministic")
	}
}
