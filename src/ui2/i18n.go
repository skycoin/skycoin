package ui2

import (
	"embed"
	"encoding/json"
	"strings"
	"sync"
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
func PickLang(cookie, accept string) string {
	if known(cookie) {
		return cookie
	}
	for _, part := range strings.Split(accept, ",") {
		tag, _, _ := strings.Cut(strings.TrimSpace(part), ";")
		tag, _, _ = strings.Cut(strings.ToLower(tag), "-")
		if known(tag) {
			return tag
		}
	}
	return "en"
}

func known(code string) bool {
	for _, l := range Langs {
		if l.Code == code {
			return true
		}
	}
	return false
}

//go:embed i18n
var i18nFS embed.FS

type catalogT struct {
	mu    sync.Mutex
	langs map[string]map[string]string
}

var catalog = &catalogT{langs: map[string]map[string]string{}}

// get is key in lang, else in English, else the key itself.
func (c *catalogT) get(lang, key string) string {
	if s := c.load(lang)[key]; s != "" {
		return s
	}
	if s := c.load("en")[key]; s != "" {
		return s
	}
	return key
}

func (c *catalogT) load(lang string) map[string]string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if m, ok := c.langs[lang]; ok {
		return m
	}
	m := map[string]string{}
	c.langs[lang] = m
	if !known(lang) {
		return m
	}
	if b, err := i18nFS.ReadFile("i18n/" + lang + ".json"); err == nil {
		_ = json.Unmarshal(b, &m) //nolint:errcheck
	}
	return m
}
