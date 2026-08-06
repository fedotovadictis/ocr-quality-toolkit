package transform

import (
	"image"

	xdraw "golang.org/x/image/draw"
)

// Downscale50 уменьшает изображение в два раза,
// затем масштабирует его обратно до исходного размера.
func Downscale50(source image.Image) image.Image {
	bounds := source.Bounds()

	width := bounds.Dx()
	height := bounds.Dy()

	downWidth := width / 2
	downHeight := height / 2

	if downWidth < 1 {
		downWidth = 1
	}
	if downHeight < 1 {
		downHeight = 1
	}

	small := image.NewRGBA(
		image.Rect(0, 0, downWidth, downHeight),
	)

	xdraw.ApproxBiLinear.Scale(
		small,
		small.Bounds(),
		source,
		bounds,
		xdraw.Over,
		nil,
	)

	result := image.NewRGBA(
		image.Rect(
			bounds.Min.X,
			bounds.Min.Y,
			bounds.Max.X,
			bounds.Max.Y,
		),
	)

	xdraw.ApproxBiLinear.Scale(
		result,
		result.Bounds(),
		small,
		small.Bounds(),
		xdraw.Over,
		nil,
	)

	return result
}
