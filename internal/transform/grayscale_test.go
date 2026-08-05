package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestGrayscale(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 2, 2))

	source.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	source.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	source.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	source.Set(1, 1, color.RGBA{R: 120, G: 80, B: 40, A: 255})

	result := Grayscale(source)

	if result.Bounds() != source.Bounds() {
		t.Fatalf(
			"expected bounds %v, got %v",
			source.Bounds(),
			result.Bounds(),
		)
	}

	changed := false

	for y := result.Bounds().Min.Y; y < result.Bounds().Max.Y; y++ {
		for x := result.Bounds().Min.X; x < result.Bounds().Max.X; x++ {
			r, g, b, _ := result.At(x, y).RGBA()

			if r != g || g != b {
				t.Fatalf(
					"pixel (%d,%d) is not grayscale: r=%d g=%d b=%d",
					x,
					y,
					r,
					g,
					b,
				)
			}

			if result.At(x, y) != source.At(x, y) {
				changed = true
			}
		}
	}

	if !changed {
		t.Fatal("expected grayscale image to differ from source")
	}
}
