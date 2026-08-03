package generator

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestSavePNG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "page.png")

	img := image.NewRGBA(image.Rect(0, 0, 120, 80))

	for y := 0; y < 80; y++ {
		for x := 0; x < 120; x++ {
			img.Set(x, y, color.White)
		}
	}

	img.Set(10, 10, color.Black)

	if err := SavePNG(path, img); err != nil {
		t.Fatalf("SavePNG returned error: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open saved PNG: %v", err)
	}
	defer file.Close()

	decoded, err := png.Decode(file)
	if err != nil {
		t.Fatalf("decode saved PNG: %v", err)
	}

	bounds := decoded.Bounds()

	if bounds.Dx() != 120 {
		t.Fatalf("expected width 120, got %d", bounds.Dx())
	}

	if bounds.Dy() != 80 {
		t.Fatalf("expected height 80, got %d", bounds.Dy())
	}
}

func TestSavePNGInvalidPath(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"missing-directory",
		"page.png",
	)

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	err := SavePNG(path, img)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
