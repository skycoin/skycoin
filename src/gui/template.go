// Template helpers
package gui

import (
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

func LoadTemplate(html_file string) (*template.Template, error) {
	const template_dir = "./src/gui/static"
	if !strings.HasPrefix(html_file, template_dir) {
		html_file = filepath.Join(template_dir, html_file)
	}

	t, err := template.ParseFiles(html_file, "./src/gui/static/common.html")
	return t, err
}

func ShowTemplate(w http.ResponseWriter, html_file string, p interface{}) {
	t, err := LoadTemplate(html_file)
	if err != nil {
		Error500(w, err.Error())
	}
	err = t.Execute(w, p)
	if err != nil {
		Error500(w, err.Error())
	}
}
