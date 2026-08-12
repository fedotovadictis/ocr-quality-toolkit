package report

import "ocr-quality-toolkit/internal/runner"

type Stats struct {
	Total             int
	Successful        int
	Failed            int
	TotalDurationMS   int64
	AverageDurationMS int64
}

func CalculateStats(results []runner.Result) Stats {
	stats := Stats{
		Total: len(results),
	}

	for _, result := range results {
		if result.Error != "" {
			stats.Failed++
		} else {
			stats.Successful++
		}

		stats.TotalDurationMS += result.DurationMS
	}

	if stats.Total > 0 {
		stats.AverageDurationMS =
			stats.TotalDurationMS / int64(stats.Total)
	}

	return stats
}
