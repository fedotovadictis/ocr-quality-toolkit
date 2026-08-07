package transform

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestTransformGolden(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 32, 32))

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			source.SetRGBA(x, y, color.RGBA{
				R: uint8(x * 7),
				G: uint8(y * 7),
				B: uint8((x + y) * 3),
				A: 255,
			})
		}
	}

	profiles := []string{
		"grayscale",
		"low-contrast",
		"noise-light",
		"jpeg-70",
		"downscale-50",
	}

	expected := map[string]string{
		"grayscale":    "f7a1bd257ab132d9cc4c3bf6b3d0b1480b8781aa02791929d07c1c1ba1c4cf41",
		"low-contrast": "ee49218cc2c42b6e38bd0c652329e7c4c18529148a8fde0fa6722f9c631e7286",
		"noise-light":  "82eb8bc9a6748584ccc9bd134150dc39dddc69995e519c7e9401bf612696ed49",
		"jpeg-70":      "69d3e4d04452a999ad30b0db5f4f434e18c12356835c19292ce932b20be94e32",
		"downscale-50": "f4fc17eca90f368dc217afa6ce2bd3f0c5f400266b05eb8f46db5b592746e851",
	}

	for _, profile := range profiles {
		t.Run(profile, func(t *testing.T) {
			first, err := Apply(source, profile, 42)
			if err != nil {
				t.Fatalf("Apply returned error: %v", err)
			}

			second, err := Apply(source, profile, 42)
			if err != nil {
				t.Fatalf("Apply returned error: %v", err)
			}

			firstHash := imageHash(t, first)
			secondHash := imageHash(t, second)

			if firstHash != secondHash {
				t.Fatalf(
					"transform %q is not deterministic: %s != %s",
					profile,
					firstHash,
					secondHash,
				)
			}

			if firstHash != expected[profile] {
				t.Fatalf(
					"golden mismatch for %q: expected %s, got %s",
					profile,
					expected[profile],
					firstHash,
				)
			}
		})
	}
}

func imageHash(t *testing.T, img image.Image) string {
	t.Helper()

	var buffer bytes.Buffer

	if err := png.Encode(&buffer, img); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}

	sum := sha256.Sum256(buffer.Bytes())

	return hex.EncodeToString(sum[:])
}
