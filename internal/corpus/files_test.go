package corpus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateImageFiles(t *testing.T) {
	root := t.TempDir()

	imagePath := filepath.Join(root, "images", "page-001.png")

	if err := os.MkdirAll(filepath.Dir(imagePath), 0o755); err != nil {
		t.Fatalf("create image dir: %v", err)
	}

	if err := os.WriteFile(imagePath, []byte("fake"), 0o600); err != nil {
		t.Fatalf("write image: %v", err)
	}

	records := []Record{
		{
			ID:    "page-001",
			Image: "images/page-001.png",
		},
	}

	if err := ValidateImageFiles(root, records); err != nil {
		t.Fatalf("ValidateImageFiles returned error: %v", err)
	}
}

func TestValidateImageFilesMissing(t *testing.T) {
	root := t.TempDir()

	records := []Record{
		{
			ID:    "page-001",
			Image: "images/missing.png",
		},
	}

	err := ValidateImageFiles(root, records)
	if err == nil {
		t.Fatal("expected missing image error, got nil")
	}
}
func TestValidateImageFilesDirectory(t *testing.T) {
	root := t.TempDir()

	imageDir := filepath.Join(root, "images", "page-001.png")

	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		t.Fatalf("create directory: %v", err)
	}

	records := []Record{
		{
			ID:    "page-001",
			Image: "images/page-001.png",
		},
	}

	err := ValidateImageFiles(root, records)
	if err == nil {
		t.Fatal("expected directory error, got nil")
	}
}
