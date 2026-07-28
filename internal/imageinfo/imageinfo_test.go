package imageinfo

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestReadDetectsFormatByContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "image.jpg")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create test image: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 10, 20))

	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatalf("encode test image: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close test image: %v", err)
	}

	format, width, height, err := Read(path)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if format != "png" {
		t.Errorf("format = %q, want %q", format, "png")
	}

	if width != 10 {
		t.Errorf("width = %d, want %d", width, 10)
	}

	if height != 20 {
		t.Errorf("height = %d, want %d", height, 20)
	}
}
