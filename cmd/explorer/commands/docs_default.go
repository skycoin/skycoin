//go:build !tinygo

package commands

import (
	"bytes"
	"html/template"
	"log"
)

// renderDocs renders the API documentation page from docTemplate.
func renderDocs(endpoints []APIEndpoint) string {
	t := template.Must(template.New("docs").Parse(docTemplate))
	b := &bytes.Buffer{}
	if err := t.Execute(b, endpoints); err != nil {
		log.Panic(err)
	}
	return b.String()
}
