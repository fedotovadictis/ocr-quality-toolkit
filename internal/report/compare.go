package report

type Comparison struct {
	CERDelta      float64 `json:"cer_delta"`
	WERDelta      float64 `json:"wer_delta"`
	CoverageDelta float64 `json:"coverage_delta"`
}

func CompareReports(
	baseline Report,
	current Report,
) Comparison {
	return Comparison{
		CERDelta:      current.Overall.CER - baseline.Overall.CER,
		WERDelta:      current.Overall.WER - baseline.Overall.WER,
		CoverageDelta: current.Overall.Coverage - baseline.Overall.Coverage,
	}

}

type Thresholds struct {
	MaxCERIncrease      float64
	MaxWERIncrease      float64
	MaxCoverageDecrease float64
}

func HasRegression(
	comparison Comparison,
	thresholds Thresholds,
) bool {
	if comparison.CERDelta > thresholds.MaxCERIncrease {
		return true
	}

	if comparison.WERDelta > thresholds.MaxWERIncrease {
		return true
	}

	if comparison.CoverageDelta < -thresholds.MaxCoverageDecrease {
		return true
	}

	return false
}
