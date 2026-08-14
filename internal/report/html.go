package report

import (
	"fmt"
	"html/template"
	"os"
)

const reportHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>OCR Quality Report</title>
	<style>
		body {
			font-family: sans-serif;
			margin: 2rem;
		}

		table {
			border-collapse: collapse;
			width: 100%;
			margin-bottom: 2rem;
		}

		th, td {
			border: 1px solid #ccc;
			padding: 0.5rem;
			text-align: left;
			vertical-align: top;
		}

		img {
			max-width: 300px;
			height: auto;
		}

		.text {
			white-space: pre-wrap;
		}

		.alignment {
			margin-top: 0.5rem;
			line-height: 1.8;
		}

		.diff {
			font-weight: bold;
			text-decoration: underline;
		}
	</style>
</head>
<body>
	<h1>OCR Quality Report</h1>

	<h2>Overall statistics</h2>

	<table>
		<tr>
			<th>Total</th>
			<th>Coverage</th>
			<th>CER</th>
			<th>WER</th>
		</tr>
		<tr>
			<td>{{.Overall.Total}}</td>
			<td>{{printf "%.4f" .Overall.Coverage}}</td>
			<td>{{printf "%.4f" .Overall.CER}}</td>
			<td>{{printf "%.4f" .Overall.WER}}</td>
		</tr>
	</table>

	<h2>Results</h2>

	{{range .Results}}
	<section>
		<h3>{{.ID}}</h3>

		<p>Status: {{.Status}}</p>
		<p>CER: {{printf "%.4f" .CER}}</p>
		<p>WER: {{printf "%.4f" .WER}}</p>

		{{if .Image}}
		<p>
			<img src="{{.Image}}" alt="{{.ID}}">
		</p>
		{{end}}

		<h4>Reference</h4>
		<div class="text">{{.Reference}}</div>

		<h4>Hypothesis</h4>
		<div class="text">{{.Hypothesis}}</div>

		<h4>Alignment</h4>
		<div class="alignment">
			{{range .Alignment}}
				{{if eq .Type "equal"}}
					<span>{{.Reference}} </span>
				{{else if eq .Type "substitute"}}
					<span class="diff">
						[{{.Reference}} → {{.Hypothesis}}]
					</span>
				{{else if eq .Type "delete"}}
					<span class="diff">
						[-{{.Reference}}]
					</span>
				{{else if eq .Type "insert"}}
					<span class="diff">
						[+{{.Hypothesis}}]
					</span>
				{{end}}
			{{end}}
		</div>
	</section>
	{{end}}
</body>
</html>
`

type HTMLResult struct {
	ID         string
	Status     string
	CER        float64
	WER        float64
	Reference  string
	Hypothesis string
	Image      string
	Alignment  []AlignmentItem
}

type HTMLReport struct {
	Overall EvaluationStats
	Results []HTMLResult
}

func WriteHTML(path string, report Report) error {
	view := HTMLReport{
		Overall: report.Overall,
		Results: make([]HTMLResult, 0, len(report.Results)),
	}

	for _, result := range report.Results {
		view.Results = append(view.Results, HTMLResult{
			ID:         result.ID,
			Status:     string(result.Status),
			CER:        result.CER,
			WER:        result.WER,
			Reference:  result.Reference,
			Hypothesis: result.Hypothesis,
			Image:      result.Image,
			Alignment: BuildAlignment(
				result.Reference,
				result.Hypothesis,
			),
		})
	}

	tmpl, err := template.New("report").Parse(reportHTMLTemplate)
	if err != nil {
		return fmt.Errorf("parse HTML template: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create HTML report %q: %w", path, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, view); err != nil {
		return fmt.Errorf("execute HTML template: %w", err)
	}

	return nil
}
