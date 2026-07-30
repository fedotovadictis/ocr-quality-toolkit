package metrics

import (
	"math"
	"testing"
)

func TestWER(t *testing.T) {
	tests := []struct {
		name       string
		reference  string
		hypothesis string
		want       float64
	}{
		{
			name:       "both empty",
			reference:  "",
			hypothesis: "",
			want:       0,
		},
		{
			name:       "empty reference",
			reference:  "",
			hypothesis: "hello",
			want:       1,
		},
		{
			name:       "equal",
			reference:  "hello world",
			hypothesis: "hello world",
			want:       0,
		},
		{
			name:       "one substitution",
			reference:  "hello world",
			hypothesis: "hello there",
			want:       0.5,
		},
		{
			name:       "one deletion",
			reference:  "hello beautiful world",
			hypothesis: "hello world",
			want:       1.0 / 3.0,
		},
		{
			name:       "one insertion",
			reference:  "hello world",
			hypothesis: "hello amazing world",
			want:       0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WER(tt.reference, tt.hypothesis)

			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("WER(%q, %q) = %v, want %v",
					tt.reference,
					tt.hypothesis,
					got,
					tt.want,
				)
			}
		})
	}
}
