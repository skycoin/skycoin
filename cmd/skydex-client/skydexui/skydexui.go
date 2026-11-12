// Package skydexui cmd/skydex-client/skydexui/skydexui.go c3-ui-dex
//
// Package skydexui is the SkyDEX trading client rendered on the server with
// templ, with htmx for the parts the React client polls. It is served at the
// root, with the React client at /angular/.
//
// Never use html/template or text/template here: they panic under TinyGo.
package skydexui

import (
	"embed"
	"errors"
	"io/fs"

	"github.com/skycoin/skycoin/src/ui2"
)

//go:generate go run github.com/a-h/templ/cmd/templ@v0.3.1100 generate
//go:generate goimports -w -local github.com/skycoin/skycoin .

//go:embed static
var embedded embed.FS

func sub(dir string) fs.FS {
	s, err := fs.Sub(embedded, dir)
	if err != nil {
		panic(err)
	}
	return s
}

// The client's own files come first, htmx comes from ui2.
var assets = ui2.NewAssets(sub("static"), ui2.StaticFS())

// Asset is the URL of a static file relative to the page's <base>.
func Asset(name string) string { return assets.URL(name) }

// Market is the client's connection to one market, as the local control API
// in cmd/skydex-client/commands holds it.
type Market interface {
	// Default is the market reference to pre-fill on the connect screen.
	Default() string
	// Status reports the connection and the market's reference and name.
	Status() (connected bool, ref, name string)
	// Connect dials ref and returns the market's display name.
	Connect(ref string) (name string, err error)
	Disconnect()
	// Call sends one request to the connected market and decodes its reply
	// into resp. A market error comes back as a *MarketError.
	Call(msgType string, req, resp any) error
}

// MarketError is an error reply from the market.
type MarketError struct {
	Code    string
	Message string
}

func (e *MarketError) Error() string { return e.Message }

// ErrNotConnected is returned by Call with no market connected.
var ErrNotConnected = errors.New("not connected to a market")
