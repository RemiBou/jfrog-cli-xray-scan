package commands

import (
	"os"
	"text/template"
)

const defaultTmpl = `{{block "list" .}}{{range .Artifacts}} 
{{.General.Name}} | {{range .Issues}}{{.Severity}}{{range .Components}}{{.FixedVersions}}{{end}}{{end}}
{{end}} {{end}}`

type resultPrinter struct {
	tmpl *template.Template
}

func newPrinter() (*resultPrinter, error) {
	tmpl, err := template.New("default").Parse(defaultTmpl)
	if err != nil {
		return nil, err
	}
	return &resultPrinter{tmpl: tmpl}, nil
}

func (r *resultPrinter) print(result ComponentSummaryResult) error {
	err := r.tmpl.Execute(os.Stdout, result)
	if err != nil {
		return err
	}
	return nil
}
