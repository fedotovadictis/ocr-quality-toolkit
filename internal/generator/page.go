package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type PageOptions struct {
	Width      int
	Height     int
	Margin     int
	FontPath   string
	FontSize   float64
	LineHeight int
	Seed       int64
}

func GeneratePage(
	text string,
	options PageOptions,
) (image.Image, error) {
	loadedFont, err := LoadFont(options.FontPath)
	if err != nil {
		return nil, fmt.Errorf("load font: %w", err)
	}

	face, err := opentype.NewFace(
		loadedFont,
		&opentype.FaceOptions{
			Size:    options.FontSize,
			DPI:     72,
			Hinting: font.HintingFull,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create font face: %w", err)
	}
	defer face.Close()

	page := image.NewRGBA(
		image.Rect(
			0,
			0,
			options.Width,
			options.Height,
		),
	)

	draw.Draw(
		page,
		page.Bounds(),
		&image.Uniform{C: color.White},
		image.Point{},
		draw.Src,
	)

	drawer := font.Drawer{
		Dst:  page,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot: fixed.P(
			options.Margin,
			options.Margin+int(options.FontSize),
		),
	}

	maxWidth := options.Width - 2*options.Margin
	lines := WrapText(text, face, maxWidth)

	for i, line := range lines {
		drawer.Dot = fixed.P(
			options.Margin,
			options.Margin+int(options.FontSize)+i*options.LineHeight,
		)
		drawer.DrawString(line)
	}

	return page, nil
}
