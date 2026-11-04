// Package ui2 is a server rendered wallet GUI built with templ and htmx. It
// sits beside the Angular wallet and reads the node through its REST API
// client, so any host holding an api.Client can serve it.
//
// Never use html/template or text/template here: they panic under TinyGo.
package ui2

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"strings"
)

//go:generate go run github.com/a-h/templ/cmd/templ@v0.3.1070 generate

// HTMXVersion is the vendored static/htmx.min.js, the same file skywire's ui2
// serves.
const HTMXVersion = "2.0.11"

//go:embed static
var static embed.FS

var staticFS = func() fs.FS {
	sub, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	return sub
}()

// assetVer maps each static file to a short hash of its content, so a page
// can name a URL that changes whenever the file does.
var assetVer = func() map[string]string {
	m := map[string]string{}
	_ = fs.WalkDir(staticFS, ".", func(p string, d fs.DirEntry, err error) error { //nolint:errcheck
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(staticFS, p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		m[p] = hex.EncodeToString(sum[:6])
		return nil
	})
	return m
}()

// Asset is the URL of a static file, relative to the page's <base>, carrying
// its content hash so it can be cached for good.
func Asset(name string) string {
	return "static/" + name + "?v=" + assetVer[name]
}

// BaseFor is the relative path from a request path, as seen below the mount
// point, back to the mount point. Every URL in a page is relative to it.
func BaseFor(path string) string {
	n := strings.Count(path, "/") - 1
	if n <= 0 {
		return "./"
	}
	return strings.Repeat("../", n)
}

// Page is what every page and fragment needs to render.
type Page struct {
	Base    string
	Lang    string
	Classic string
}

type pageKey struct{}

func with(ctx context.Context, p Page) context.Context {
	return context.WithValue(ctx, pageKey{}, &p)
}

// From is the Page in ctx, or an English page rooted at "./".
func From(ctx context.Context) *Page {
	if p, ok := ctx.Value(pageKey{}).(*Page); ok {
		return p
	}
	return &Page{Base: "./", Lang: "en", Classic: "../"}
}

// T translates key into the page language. args are name, value pairs that
// fill {{ name }} placeholders.
func T(ctx context.Context, key string, args ...string) string {
	s := catalog.get(From(ctx).Lang, key)
	for i := 0; i+1 < len(args); i += 2 {
		s = strings.ReplaceAll(s, "{{ "+args[i]+" }}", args[i+1])
	}
	return s
}
