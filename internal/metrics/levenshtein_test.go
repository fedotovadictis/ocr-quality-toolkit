package metrics

import "testing"

func TestDistance(t *testing.T) {
	tests := []struct {
		name       string
		reference  string
		hypothesis string
		want       int
	}{
		{
			name:       "both empty",
			reference:  "",
			hypothesis: "",
			want:       0,
		},
		{
			name:       "equal text",
			reference:  "текст",
			hypothesis: "текст",
			want:       0,
		},
		{
			name:       "empty reference",
			reference:  "",
			hypothesis: "три",
			want:       3,
		},
		{
			name:       "empty hypothesis",
			reference:  "три",
			hypothesis: "",
			want:       3,
		},
		{
			name:       "one substitution",
			reference:  "кот",
			hypothesis: "кит",
			want:       1,
		},
		{
			name:       "one deletion",
			reference:  "счёт",
			hypothesis: "сёт",
			want:       1,
		},
		{
			name:       "one insertion",
			reference:  "кот",
			hypothesis: "крот",
			want:       1,
		},
		{
			name:       "latin",
			reference:  "kitten",
			hypothesis: "sitting",
			want:       3,
		},
		{
			name:       "unicode runes not bytes",
			reference:  "ё",
			hypothesis: "е",
			want:       1,
		},
		{
			name:       "mixed languages",
			reference:  "OCR текст",
			hypothesis: "OCR тест",
			want:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Distance(tt.reference, tt.hypothesis)

			if got != tt.want {
				t.Fatalf(
					"Distance(%q, %q) = %d, want %d",
					tt.reference,
					tt.hypothesis,
					got,
					tt.want,
				)
			}
		})
	}
}
