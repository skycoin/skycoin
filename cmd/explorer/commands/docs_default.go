//go:build !tinygo

// Package commands cmd/explorer/commands/docs_default.go c4-ui-explorer
package commands

import (
	"bytes"
	"fmt"
	"html/template"
)

// renderDocs renders the API documentation page from docTemplate.
func renderDocs(endpoints []APIEndpoint) (string, error) {
	t, err := template.New("docs").Parse(docTemplate)
	if err != nil {
		return "", fmt.Errorf("parse docs template: %w", err)
	}
	b := &bytes.Buffer{}
	if err := t.Execute(b, endpoints); err != nil {
		return "", fmt.Errorf("render docs template: %w", err)
	}
	return b.String(), nil
}
