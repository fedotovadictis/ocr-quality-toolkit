package generator

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildSyntheticImage(t *testing.T) {
	dir := t.TempDir()

	sourcePath := filepath.Join(dir, "source.png")
	targetPath := filepath.Join(dir, "synthetic.png")

	source := image.NewRGBA(image.Rect(0, 0, 20, 20))

	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			source.Set(x, y, color.RGBA{
				R: uint8(x * 10),
				G: uint8(y * 10),
				B: 100,
				A: 255,
			})
		}
	}

	if err := SavePNG(sourcePath, source); err != nil {
		t.Fatalf("save source image: %v", err)
	}

	err := BuildSyntheticImage(
		sourcePath,
		targetPath,
		"grayscale",
		42,
	)
	if err != nil {
		t.Fatalf("BuildSyntheticImage returned error: %v", err)
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("synthetic image was not created: %v", err)
	}

	if info.Size() == 0 {
		t.Fatal("synthetic image is empty")
	}
}
