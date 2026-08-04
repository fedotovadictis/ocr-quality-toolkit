package generator

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

func TestGeneratePage(t *testing.T) {
	fontPath := writeTestFont(t)

	options := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   fontPath,
		FontSize:   24,
		LineHeight: 32,
		Seed:       42,
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
func TestGeneratePageWrapsLongText(t *testing.T) {
	fontPath := writeTestFont(t)

	options := PageOptions{
		Width:      180,
		Height:     300,
		Margin:     20,
		FontPath:   fontPath,
		FontSize:   20,
		LineHeight: 30,
	}

	page, err := GeneratePage(
		"Привет мир это длинная строка для проверки переноса",
		options,
	)
	if err != nil {
		t.Fatalf("GeneratePage returned error: %v", err)
	}

	secondLineStart := options.Margin + int(options.FontSize) + 5
	secondLineEnd := secondLineStart + options.LineHeight

	foundInk := false

	for y := secondLineStart; y < secondLineEnd && !foundInk; y++ {
		for x := options.Margin; x < options.Width-options.Margin; x++ {
			pixel := color.RGBAModel.Convert(page.At(x, y)).(color.RGBA)

			if pixel != (color.RGBA{R: 255, G: 255, B: 255, A: 255}) {
				foundInk = true
				break
			}
		}
	}

	if !foundInk {
		t.Fatal("expected text to be drawn on a second line")
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
func writeTestFont(t *testing.T) string {
	t.Helper()

	fontPath := filepath.Join(t.TempDir(), "font.ttf")

	if err := os.WriteFile(fontPath, goregular.TTF, 0o600); err != nil {
		t.Fatalf("write test font: %v", err)
	}

	return fontPath
}
func TestGeneratePageSameInputProducesSamePNG(t *testing.T) {
	fontPath := writeTestFont(t)

	options := PageOptions{
		Width:      600,
		Height:     800,
		Margin:     40,
		FontPath:   fontPath,
		FontSize:   24,
		LineHeight: 32,
		Seed:       42,
	}

	firstPage, err := GeneratePage("Привет, OCR!", options)
	if err != nil {
		t.Fatalf("first GeneratePage returned error: %v", err)
	}

	secondPage, err := GeneratePage("Привет, OCR!", options)
	if err != nil {
		t.Fatalf("second GeneratePage returned error: %v", err)
	}

	firstPNG := encodePNG(t, firstPage)
	secondPNG := encodePNG(t, secondPage)

	if !bytes.Equal(firstPNG, secondPNG) {
		t.Fatal("same input and seed must produce identical PNG bytes")
	}
}
func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()

	var buffer bytes.Buffer

	if err := png.Encode(&buffer, img); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}

	return buffer.Bytes()
}
