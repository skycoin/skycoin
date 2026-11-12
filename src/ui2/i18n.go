// Package ui2 src/ui2/i18n.go c3-ui-wallet
package ui2

import (
	"embed"
	"io/fs"
)

// Lang is a language the pages offer.
type Lang struct {
	Code string
	Name string
}

// Langs follows the languages of the Angular wallet
// (src/skycoin-web/src/app/app.config.ts).
var Langs = []Lang{
	{"en", "English"},
	{"es", "Español"},
}

// LangCookie holds the chosen language. ui2.js keeps it in step with the
// Angular wallet's "lang" localStorage entry.
const LangCookie = "ui2-lang"

// PickLang is the language named by the cookie, else the first supported one
// in an Accept-Language header, else English.
func PickLang(cookie, accept string) string { return catalog.Pick(cookie, accept) }

//go:embed i18n
var i18nFS embed.FS

var catalog = NewCatalog(func() fs.FS {
	sub, err := fs.Sub(i18nFS, "i18n")
	if err != nil {
		panic(err)
	}
	return sub
}(), Langs)
