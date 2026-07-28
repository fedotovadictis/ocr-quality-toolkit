package hash

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSHA256(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	err := os.WriteFile(path, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("write test file: %v", err)
	}

	got, err := FileSHA256(path)
	if err != nil {
		t.Fatalf("FileSHA256() error: %v", err)
	}

	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	if got != want {
		t.Errorf("FileSHA256() = %q, want %q", got, want)
	}
}
