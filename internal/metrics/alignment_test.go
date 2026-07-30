package metrics

import (
	"reflect"
	"testing"
)

func TestAlign(t *testing.T) {
	tests := []struct {
		name       string
		reference  []string
		hypothesis []string
		want       Alignment
	}{
		{
			name:       "both empty",
			reference:  nil,
			hypothesis: nil,
			want:       Alignment{},
		},
		{
			name:       "equal",
			reference:  []string{"к", "о", "т"},
			hypothesis: []string{"к", "о", "т"},
			want: Alignment{
				Distance: 0,
				Hits:     3,
				Operations: []Operation{
					{Type: OperationEqual, Reference: "к", Hypothesis: "к"},
					{Type: OperationEqual, Reference: "о", Hypothesis: "о"},
					{Type: OperationEqual, Reference: "т", Hypothesis: "т"},
				},
			},
		},
		{
			name:       "substitution",
			reference:  []string{"к", "о", "т"},
			hypothesis: []string{"к", "и", "т"},
			want: Alignment{
				Distance:      1,
				Hits:          2,
				Substitutions: 1,
				Operations: []Operation{
					{Type: OperationEqual, Reference: "к", Hypothesis: "к"},
					{Type: OperationSubstitute, Reference: "о", Hypothesis: "и"},
					{Type: OperationEqual, Reference: "т", Hypothesis: "т"},
				},
			},
		},
		{
			name:       "deletion",
			reference:  []string{"к", "р", "о", "т"},
			hypothesis: []string{"к", "о", "т"},
			want: Alignment{
				Distance:  1,
				Hits:      3,
				Deletions: 1,
				Operations: []Operation{
					{Type: OperationEqual, Reference: "к", Hypothesis: "к"},
					{Type: OperationDelete, Reference: "р"},
					{Type: OperationEqual, Reference: "о", Hypothesis: "о"},
					{Type: OperationEqual, Reference: "т", Hypothesis: "т"},
				},
			},
		},
		{
			name:       "insertion",
			reference:  []string{"к", "о", "т"},
			hypothesis: []string{"к", "р", "о", "т"},
			want: Alignment{
				Distance:   1,
				Hits:       3,
				Insertions: 1,
				Operations: []Operation{
					{Type: OperationEqual, Reference: "к", Hypothesis: "к"},
					{Type: OperationInsert, Hypothesis: "р"},
					{Type: OperationEqual, Reference: "о", Hypothesis: "о"},
					{Type: OperationEqual, Reference: "т", Hypothesis: "т"},
				},
			},
		},
		{
			name:       "empty reference",
			reference:  nil,
			hypothesis: []string{"О", "0"},
			want: Alignment{
				Distance:   2,
				Insertions: 2,
				Operations: []Operation{
					{Type: OperationInsert, Hypothesis: "О"},
					{Type: OperationInsert, Hypothesis: "0"},
				},
			},
		},
		{
			name:       "empty hypothesis",
			reference:  []string{"З", "3"},
			hypothesis: nil,
			want: Alignment{
				Distance:  2,
				Deletions: 2,
				Operations: []Operation{
					{Type: OperationDelete, Reference: "З"},
					{Type: OperationDelete, Reference: "3"},
				},
			},
		},
		{
			name:       "deterministic substitution priority",
			reference:  []string{"a", "b"},
			hypothesis: []string{"b", "a"},
			want: Alignment{
				Distance:      2,
				Substitutions: 2,
				Operations: []Operation{
					{Type: OperationSubstitute, Reference: "a", Hypothesis: "b"},
					{Type: OperationSubstitute, Reference: "b", Hypothesis: "a"},
				},
			},
		},
		{
			name:       "cyrillic substitution",
			reference:  []string{"с", "ч", "ё", "т"},
			hypothesis: []string{"с", "ч", "е", "т"},
			want: Alignment{
				Distance:      1,
				Hits:          3,
				Substitutions: 1,
				Operations: []Operation{
					{Type: OperationEqual, Reference: "с", Hypothesis: "с"},
					{Type: OperationEqual, Reference: "ч", Hypothesis: "ч"},
					{Type: OperationSubstitute, Reference: "ё", Hypothesis: "е"},
					{Type: OperationEqual, Reference: "т", Hypothesis: "т"},
				},
			},
		},
		{
			name:       "ocr confusable characters",
			reference:  []string{"0", "3"},
			hypothesis: []string{"О", "З"},
			want: Alignment{
				Distance:      2,
				Hits:          0,
				Substitutions: 2,
				Operations: []Operation{
					{Type: OperationSubstitute, Reference: "0", Hypothesis: "О"},
					{Type: OperationSubstitute, Reference: "3", Hypothesis: "З"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Align(tt.reference, tt.hypothesis)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf(
					"unexpected alignment:\ngot:  %#v\nwant: %#v",
					got,
					tt.want,
				)
			}
		})
	}
}
