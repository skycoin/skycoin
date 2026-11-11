// Package explorerui is the blockchain explorer rendered on the server with
// templ, with htmx only for paging and refreshing in place. Every page works
// without JavaScript, so crawlers see the same content people do. It shares
// the asset, language and paging helpers of the wallet's ui2.
//
// Never use html/template or text/template here: they panic under TinyGo.
package explorerui

import (
	"context"
	"embed"
	"io/fs"

	"github.com/skycoin/skycoin/src/ui2"
)

//go:generate go run github.com/a-h/templ/cmd/templ@v0.3.1100 generate
//go:generate goimports -w -local github.com/skycoin/skycoin .

//go:embed static i18n
var embedded embed.FS

func sub(dir string) fs.FS {
	s, err := fs.Sub(embedded, dir)
	if err != nil {
		panic(err)
	}
	return s
}

// The explorer's own files come first, htmx comes from ui2.
var assets = ui2.NewAssets(sub("static"), ui2.StaticFS())

// Asset is the URL of a static file relative to the page's <base>.
func Asset(name string) string { return assets.URL(name) }

// Langs follows the Angular explorer (explorer/src/app/app.config.ts).
var Langs = []ui2.Lang{
	{Code: "en", Name: "English"},
	{Code: "es", Name: "Español"},
	{Code: "zh", Name: "中文"},
}

var catalog = ui2.NewCatalog(sub("i18n"), Langs)

// LangCookie holds the language picked with the language links.
const LangCookie = "explorer-lang"

// Page is what every page needs to render.
type Page struct {
	// Base leads from the page back to where the handler is mounted.
	Base string
	// Root leads from the mount point to the explorer's root, where the API
	// docs and the Angular explorer are served.
	Root string
	// Path is the request path below the mount point.
	Path string
	// Origin is scheme://host of the public site, for canonical links.
	Origin string
	Lang   string
	Coin   Coin
	// Syncing is a node still catching up with the network.
	Syncing bool
	// Classic is the same page in the Angular explorer, empty when there is
	// none.
	Classic string
}

type pageKey struct{}

func with(ctx context.Context, p *Page) context.Context {
	return context.WithValue(ctx, pageKey{}, p)
}

// From is the Page in ctx, or an English page rooted at "./".
func From(ctx context.Context) *Page {
	if p, ok := ctx.Value(pageKey{}).(*Page); ok {
		return p
	}
	return &Page{Base: "./", Root: "./", Lang: "en", Coin: defaultCoin}
}

// T translates key into the page language. args are name, value pairs that
// fill {{ name }} placeholders.
func T(ctx context.Context, key string, args ...string) string {
	return catalog.T(From(ctx).Lang, key, args...)
}
