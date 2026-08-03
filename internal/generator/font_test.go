package generator

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

func TestLoadFont(t *testing.T) {
	tests := []struct {
		name      string
		prepare   func(t *testing.T) string
		wantError bool
	}{
		{
			name: "valid TTF",
			prepare: func(t *testing.T) string {
				t.Helper()

				path := filepath.Join(t.TempDir(), "font.ttf")

				if err := os.WriteFile(path, goregular.TTF, 0o600); err != nil {
					t.Fatalf("write test font: %v", err)
				}

				return path
			},
			wantError: false,
		},
		{
			name: "missing file",
			prepare: func(t *testing.T) string {
				t.Helper()

				return filepath.Join(
					t.TempDir(),
					"missing.ttf",
				)
			},
			wantError: true,
		},
		{
			name: "invalid TTF",
			prepare: func(t *testing.T) string {
				t.Helper()

				path := filepath.Join(t.TempDir(), "broken.ttf")

				if err := os.WriteFile(
					path,
					[]byte("not a font"),
					0o600,
				); err != nil {
					t.Fatalf("write broken font: %v", err)
				}

				return path
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.prepare(t)

			font, err := LoadFont(path)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("LoadFont returned error: %v", err)
			}

			if font == nil {
				t.Fatal("expected loaded font, got nil")
			}
		})
	}
}
