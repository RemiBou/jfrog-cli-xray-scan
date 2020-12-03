package commands

import (
	"io"
	"text/template"
)

const defaultTmpl = `|{{printf "%-30v" "component"}}|{{printf "%-10v" "issues (high/medium/low)"}}|{{printf "%-10v" "min fix version"}}|{{printf "%-10v" "licences"}}
{{range $val := .}} 
{{printf "%-30v" .Component}}|{{.HighCount}}/{{.MediumCount}}/{{printf "%-30v" .LowCount}}|{{printf "%-30v" .MinVersion}}|{{printf "%-30v" .License}}
{{end}}`

type resultPrinter struct {
	tmpl *template.Template
}
type resultLineSummary struct {
	Component   string
	HighCount   int
	MediumCount int
	LowCount    int
	MinVersion  string
	License     string
}

func newPrinter() (*resultPrinter, error) {
	tmpl, err := template.New("default").Parse(defaultTmpl)
	if err != nil {
		return nil, err
	}
	return &resultPrinter{tmpl: tmpl}, nil
}

func (r *resultPrinter) print(result ComponentSummaryResult, writer io.Writer) error {
	lines := make([]resultLineSummary, 0, len(result.Artifacts))
	for _, artifact := range result.Artifacts {
		lines = append(lines, resultLineSummary{Component: artifact.General.ComponentID})
	}
	err := r.tmpl.Execute(writer, lines)
	if err != nil {
		return err
	}
	return nil
}
