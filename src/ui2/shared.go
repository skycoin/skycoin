// Package ui2 src/ui2/shared.go c3-ui-wallet
package ui2

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"sync"
)

// StaticFS holds the wallet's static files, among them htmx.min.js, which
// other server rendered pages can serve as their own.
func StaticFS() fs.FS { return staticFS }

// Assets serves static files under URLs that carry a hash of their content,
// so a browser can keep them for good. A name is looked up in each file
// system in turn.
type Assets struct {
	layers []fs.FS
	mu     sync.Mutex
	ver    map[string]string
}

// NewAssets serves the files of layers, the first one holding a name winning.
func NewAssets(layers ...fs.FS) *Assets {
	return &Assets{layers: layers, ver: map[string]string{}}
}

// Open finds name in the first layer that has it.
func (a *Assets) Open(name string) (fs.File, error) {
	var err error
	for _, l := range a.layers {
		var f fs.File
		if f, err = l.Open(name); err == nil {
			return f, nil
		}
	}
	if err == nil {
		err = fs.ErrNotExist
	}
	return nil, err
}

// Version is a short hash of the file's content, empty for a missing file.
func (a *Assets) Version(name string) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if v, ok := a.ver[name]; ok {
		return v
	}
	v := ""
	if b, err := fs.ReadFile(a, name); err == nil {
		sum := sha256.Sum256(b)
		v = hex.EncodeToString(sum[:6])
	}
	a.ver[name] = v
	return v
}

// URL is the address of a file relative to the page's <base>.
func (a *Assets) URL(name string) string {
	return "static/" + name + "?v=" + a.Version(name)
}

// ServeHTTP serves a request for /static/<name>, marking it immutable when
// it names the current version.
func (a *Assets) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/static/")
	if v := r.URL.Query().Get("v"); v != "" && v == a.Version(name) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	}
	http.StripPrefix("/static/", http.FileServer(http.FS(a))).ServeHTTP(w, r)
}

// Catalog holds the translations of a set of pages, one <code>.json file of
// flat keys per language.
type Catalog struct {
	fsys  fs.FS
	langs []Lang
	mu    sync.Mutex
	m     map[string]map[string]string
}

// NewCatalog reads the translations in fsys. langs are the languages offered,
// English among them.
func NewCatalog(fsys fs.FS, langs []Lang) *Catalog {
	return &Catalog{fsys: fsys, langs: langs, m: map[string]map[string]string{}}
}

// Langs are the languages offered.
func (c *Catalog) Langs() []Lang { return c.langs }

// Known reports whether code is a language offered.
func (c *Catalog) Known(code string) bool {
	for _, l := range c.langs {
		if l.Code == code {
			return true
		}
	}
	return false
}

// Pick is the language named by the cookie, else the first offered one in an
// Accept-Language header, else English.
func (c *Catalog) Pick(cookie, accept string) string {
	if c.Known(cookie) {
		return cookie
	}
	for _, part := range strings.Split(accept, ",") {
		tag, _, _ := strings.Cut(strings.TrimSpace(part), ";")
		tag, _, _ = strings.Cut(strings.ToLower(tag), "-")
		if c.Known(tag) {
			return tag
		}
	}
	return "en"
}

// Get is key in lang, else in English, else the key itself.
func (c *Catalog) Get(lang, key string) string {
	if s := c.Load(lang)[key]; s != "" {
		return s
	}
	if s := c.Load("en")[key]; s != "" {
		return s
	}
	return key
}

// T is key in lang with its {{ name }} placeholders filled from args, which
// are name, value pairs.
func (c *Catalog) T(lang, key string, args ...string) string {
	s := c.Get(lang, key)
	for i := 0; i+1 < len(args); i += 2 {
		s = strings.ReplaceAll(s, "{{ "+args[i]+" }}", args[i+1])
	}
	return s
}

// Load is every key of lang, empty for a language not offered.
func (c *Catalog) Load(lang string) map[string]string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if m, ok := c.m[lang]; ok {
		return m
	}
	m := map[string]string{}
	c.m[lang] = m
	if !c.Known(lang) {
		return m
	}
	if b, err := fs.ReadFile(c.fsys, lang+".json"); err == nil {
		_ = json.Unmarshal(b, &m) //nolint:errcheck
	}
	return m
}
