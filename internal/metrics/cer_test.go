package metrics

import (
	"math"
	"testing"
)

func TestCER(t *testing.T) {
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
			hypothesis: "кот",
			want:       1,
		},
		{
			name:       "equal",
			reference:  "кот",
			hypothesis: "кот",
			want:       0,
		},
		{
			name:       "one substitution",
			reference:  "кот",
			hypothesis: "кит",
			want:       1.0 / 3.0,
		},
		{
			name:       "one deletion",
			reference:  "крот",
			hypothesis: "кот",
			want:       1.0 / 4.0,
		},
		{
			name:       "one insertion",
			reference:  "кот",
			hypothesis: "крот",
			want:       1.0 / 3.0,
		},
		{
			name:       "unicode characters",
			reference:  "ёж",
			hypothesis: "еж",
			want:       1.0 / 2.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CER(tt.reference, tt.hypothesis)

			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("CER(%q, %q) = %v, want %v",
					tt.reference,
					tt.hypothesis,
					got,
					tt.want)
			}
		})
	}
}
