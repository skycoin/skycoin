//go:build tinygo

package commands

// renderDocs: TinyGo's reflect does not implement func-signature introspection
// (reflect.Type.NumOut), which html/template requires to validate its built-in
// escaping functions. Rather than pull in html/template, serve a static notice;
// every other part of the explorer (block/transaction views, API proxy, static
// frontend) works normally. The live API endpoints remain reachable directly.
func renderDocs(endpoints []APIEndpoint) string {
	_ = endpoints
	return `<html><head><title>skycoin explorer</title></head><body>` +
		`<h1>Skycoin blockchain explorer</h1>` +
		`<p>The rendered API documentation page is not available in this ` +
		`(TinyGo) build. The API endpoints themselves are fully functional; ` +
		`query them directly.</p></body></html>`
}
