package metrics

import (
	"math"
	"testing"
)

func TestSimilarity(t *testing.T) {
	tests := []struct {
		name       string
		reference  []string
		hypothesis []string
		want       float64
	}{
		{
			name:       "both empty",
			reference:  []string{},
			hypothesis: []string{},
			want:       1,
		},
		{
			name:       "equal",
			reference:  []string{"a", "b", "c"},
			hypothesis: []string{"a", "b", "c"},
			want:       1,
		},
		{
			name:       "one substitution",
			reference:  []string{"a", "b", "c"},
			hypothesis: []string{"a", "x", "c"},
			want:       2.0 / 3.0,
		},
		{
			name:       "one deletion",
			reference:  []string{"a", "b", "c"},
			hypothesis: []string{"a", "b"},
			want:       2.0 / 3.0,
		},
		{
			name:       "one insertion",
			reference:  []string{"a", "b"},
			hypothesis: []string{"a", "b", "c"},
			want:       2.0 / 3.0,
		},
		{
			name:       "completely different",
			reference:  []string{"a", "b"},
			hypothesis: []string{"x", "y"},
			want:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Similarity(tt.reference, tt.hypothesis)

			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf(
					"Similarity(%v, %v) = %v, want %v",
					tt.reference,
					tt.hypothesis,
					got,
					tt.want,
				)
			}
		})
	}
}
