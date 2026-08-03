package generator

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

func TestGeneratePage(t *testing.T) {
	fontPath := filepath.Join(t.TempDir(), "font.ttf")

	if err := os.WriteFile(fontPath, goregular.TTF, 0o600); err != nil {
		t.Fatalf("write test font: %v", err)
	}

	options := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   fontPath,
		FontSize:   24,
		LineHeight: 32,
	}

	page, err := GeneratePage("Привет, OCR!", options)
	if err != nil {
		t.Fatalf("GeneratePage returned error: %v", err)
	}

	bounds := page.Bounds()

	if bounds.Dx() != options.Width {
		t.Fatalf(
			"expected width %d, got %d",
			options.Width,
			bounds.Dx(),
		)
	}

	if bounds.Dy() != options.Height {
		t.Fatalf(
			"expected height %d, got %d",
			options.Height,
			bounds.Dy(),
		)
	}

	if isImageWhite(page) {
		t.Fatal("expected page with drawn text, got completely white image")
	}
}

func TestGeneratePageMissingFont(t *testing.T) {
	options := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   filepath.Join(t.TempDir(), "missing.ttf"),
		FontSize:   24,
		LineHeight: 32,
	}

	_, err := GeneratePage("Привет", options)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func isImageWhite(img interface {
	Bounds() image.Rectangle
	At(x, y int) color.Color
}) bool {
	bounds := img.Bounds()

	white := color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			got := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			if got != white {
				return false
			}
		}
	}

	return true
}
