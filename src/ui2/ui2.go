// Package ui2 src/ui2/ui2.go c3-ui-wallet
//
// Package ui2 is a server rendered wallet GUI built with templ and htmx. It
// sits beside the Angular wallet and reads the node through its REST API
// client, so any host holding an api.Client can serve it.
//
// Never use html/template or text/template here: they panic under TinyGo.
package ui2

import (
	"context"
	"embed"
	"io/fs"
	"strings"
)

//go:generate go run github.com/a-h/templ/cmd/templ@v0.3.1100 generate
//go:generate goimports -w -local github.com/skycoin/skycoin .

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

var assets = NewAssets(staticFS)

// Asset is the URL of a static file, relative to the page's <base>, carrying
// its content hash so it can be cached for good.
func Asset(name string) string {
	return assets.URL(name)
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
	Root    string
	// Coin is the coin the pages are about, and Coins every coin the host
	// offers, when it serves several.
	Coin  *Coin
	Coins []Coin
}

// Gateway is the gateway of the page's coin, empty for a Skycoin or fiber coin.
func (p *Page) Gateway() string {
	if p.Coin == nil {
		return ""
	}
	return p.Coin.Gateway()
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
	return &Page{Base: "./", Lang: "en", Classic: "angular/"}
}

// T translates key into the page language. args are name, value pairs that
// fill {{ name }} placeholders.
func T(ctx context.Context, key string, args ...string) string {
	return catalog.T(From(ctx).Lang, key, args...)
}
