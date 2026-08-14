package report

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

func TestCompareReports(t *testing.T) {
	baseline := Report{
		Overall: EvaluationStats{
			CER:      0.10,
			WER:      0.20,
			Coverage: 0.90,
		},
	}

	current := Report{
		Overall: EvaluationStats{
			CER:      0.12,
			WER:      0.18,
			Coverage: 0.95,
		},
	}

	comparison := CompareReports(
		baseline,
		current,
	)

	if !almostEqual(comparison.CERDelta, 0.02) {
		t.Fatalf(
			"expected CER delta 0.02, got %v",
			comparison.CERDelta,
		)
	}

	if !almostEqual(comparison.WERDelta, -0.02) {
		t.Fatalf(
			"expected WER delta -0.02, got %v",
			comparison.WERDelta,
		)
	}

	if !almostEqual(comparison.CoverageDelta, 0.05) {
		t.Fatalf(
			"expected coverage delta 0.05, got %v",
			comparison.CoverageDelta,
		)
	}
}

func TestHasRegression(t *testing.T) {
	tests := []struct {
		name       string
		comparison Comparison
		want       bool
	}{
		{
			name: "no regression",
			comparison: Comparison{
				CERDelta:      0.01,
				WERDelta:      -0.02,
				CoverageDelta: 0.03,
			},
			want: false,
		},
		{
			name: "CER regression",
			comparison: Comparison{
				CERDelta: 0.03,
			},
			want: true,
		},
		{
			name: "WER regression",
			comparison: Comparison{
				WERDelta: 0.04,
			},
			want: true,
		},
		{
			name: "coverage regression",
			comparison: Comparison{
				CoverageDelta: -0.06,
			},
			want: true,
		},
	}

	thresholds := Thresholds{
		MaxCERIncrease:      0.02,
		MaxWERIncrease:      0.02,
		MaxCoverageDecrease: 0.05,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasRegression(
				tt.comparison,
				thresholds,
			)

			if got != tt.want {
				t.Fatalf(
					"expected regression=%v, got %v",
					tt.want,
					got,
				)
			}
		})
	}
}
