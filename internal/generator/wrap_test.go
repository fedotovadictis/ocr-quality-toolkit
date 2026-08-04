package generator

import (
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func TestWrapText(t *testing.T) {
	parsedFont, err := opentype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("parse test font: %v", err)
	}

	face, err := opentype.NewFace(
		parsedFont,
		&opentype.FaceOptions{
			Size:    20,
			DPI:     72,
			Hinting: font.HintingFull,
		},
	)
	if err != nil {
		t.Fatalf("create test font face: %v", err)
	}
	defer face.Close()

	tests := []struct {
		name     string
		text     string
		maxWidth int
		want     []string
	}{
		{
			name:     "short text",
			text:     "Привет",
			maxWidth: 500,
			want:     []string{"Привет"},
		},
		{
			name:     "long text",
			text:     "Привет мир это длинная строка",
			maxWidth: 120,
			want: []string{
				"Привет мир",
				"это длинная",
				"строка",
			},
		},
		{
			name:     "empty text",
			text:     "",
			maxWidth: 120,
			want:     []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapText(tt.text, face, tt.maxWidth)

			if len(got) != len(tt.want) {
				t.Fatalf(
					"expected %d lines, got %d: %#v",
					len(tt.want),
					len(got),
					got,
				)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf(
						"line %d: expected %q, got %q",
						i,
						tt.want[i],
						got[i],
					)
				}
			}
		})
	}
}
