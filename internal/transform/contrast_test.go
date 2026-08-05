package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestLowContrast(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 2, 1))

	source.Set(0, 0, color.RGBA{R: 20, G: 20, B: 20, A: 255})
	source.Set(1, 0, color.RGBA{R: 240, G: 240, B: 240, A: 255})

	result := LowContrast(source)

	if result.Bounds() != source.Bounds() {
		t.Fatalf(
			"expected bounds %v, got %v",
			source.Bounds(),
			result.Bounds(),
		)
	}

	darkBefore := color.RGBAModel.Convert(source.At(0, 0)).(color.RGBA)
	darkAfter := color.RGBAModel.Convert(result.At(0, 0)).(color.RGBA)

	if darkAfter.R <= darkBefore.R {
		t.Fatalf(
			"expected dark pixel to become lighter: before=%v after=%v",
			darkBefore,
			darkAfter,
		)
	}

	lightBefore := color.RGBAModel.Convert(source.At(1, 0)).(color.RGBA)
	lightAfter := color.RGBAModel.Convert(result.At(1, 0)).(color.RGBA)

	if lightAfter.R >= lightBefore.R {
		t.Fatalf(
			"expected light pixel to become darker: before=%v after=%v",
			lightBefore,
			lightAfter,
		)
	}

	if darkAfter == darkBefore && lightAfter == lightBefore {
		t.Fatal("expected low-contrast image to differ from source")
	}
}
