package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestApply(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 8, 8))

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			source.SetRGBA(x, y, color.RGBA{
				R: uint8(x * 20),
				G: uint8(y * 20),
				B: uint8((x + y) * 10),
				A: 255,
			})
		}
	}

	tests := []struct {
		name    string
		profile string
		seed    int64
	}{
		{
			name:    "grayscale",
			profile: "grayscale",
		},
		{
			name:    "low contrast",
			profile: "low-contrast",
		},
		{
			name:    "noise light",
			profile: "noise-light",
			seed:    42,
		},
		{
			name:    "jpeg 70",
			profile: "jpeg-70",
		},
		{
			name:    "downscale 50",
			profile: "downscale-50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Apply(source, tt.profile, tt.seed)
			if err != nil {
				t.Fatalf("Apply returned error: %v", err)
			}

			if result.Bounds() != source.Bounds() {
				t.Fatalf(
					"expected bounds %v, got %v",
					source.Bounds(),
					result.Bounds(),
				)
			}
		})
	}
}

func TestApplyUnknownProfile(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 4, 4))

	_, err := Apply(source, "unknown", 42)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
