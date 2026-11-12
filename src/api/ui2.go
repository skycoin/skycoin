// Package api src/api/ui2.go c3-coin-api
package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/klauspost/compress/gzhttp"
)

// HandlerClient is a client of the API that h serves, calling h in process
// instead of over the network. host is the Host its requests carry.
func HandlerClient(h http.Handler, host string) *Client {
	return &Client{
		Addr:       "http://" + host + "/",
		HTTPClient: &http.Client{Transport: localTransport{h}},
	}
}

// mountUI2 serves the page handler that build makes at the site root. Its
// client calls this mux in process, so the pages go through the same API
// set, auth and CSRF checks as any other client, without a network round
// trip.
func mountUI2(mux *http.ServeMux, c muxConfig, build func(*Client) http.Handler) {
	client := HandlerClient(mux, c.host)
	client.SetAuth(c.username, c.password)
	h := build(client)
	if !c.disableCSP {
		h = CSPHandler(h, ContentSecurityPolicy)
	}
	h = CSRFCheck(apiVersion1, c.disableCSRF, h)
	if !c.disableHeaderCheck {
		h = ui2HeaderCheck(c, h)
	}
	h = basicAuth(apiVersion1, c.username, c.password, "skycoin daemon", h)
	mux.Handle("/", gzhttp.GzipHandler(h))
	mux.Handle("/ui2/", MovedPrefix("/ui2", ""))
}

// AngularPrefix is where the Angular wallet is served. Its build expects
// this base (src/skycoin-web/angular.json).
const AngularPrefix = "/angular"

// MountAngular serves the Angular wallet h at AngularPrefix. Its assets
// and favicon stay at the root too, where the server rendered pages take
// them from. withUI2 says the root belongs to those pages, else it
// redirects to the Angular wallet.
func MountAngular(mux *http.ServeMux, h http.Handler, withUI2 bool) {
	mux.Handle(AngularPrefix+"/", http.StripPrefix(AngularPrefix, h))
	mux.Handle("/assets/", h)
	mux.Handle("/favicon.ico", h)
	mux.Handle("/ngsw-worker.js", SafetyWorker(h))
	if !withUI2 {
		mux.Handle("/{$}", http.RedirectHandler(AngularPrefix+"/", http.StatusFound))
	}
}

// SafetyWorker answers for the service worker the Angular wallet registered
// at the root before it moved, with the Angular build's safety-worker.js,
// which unregisters it and drops its caches. h serves the Angular build.
func SafetyWorker(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		r2.URL.Path, r2.URL.RawPath = "/safety-worker.js", ""
		h.ServeHTTP(w, r2)
	})
}

// MovedPrefix redirects a request below from to the same path below to,
// keeping the query.
func MovedPrefix(from, to string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Clean(to + "/" + strings.TrimPrefix(r.URL.Path, from))
		if r.URL.RawQuery != "" {
			p += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, p, http.StatusMovedPermanently) //nolint:gosec // p is a cleaned path on this host
	})
}

// ui2HeaderCheck always checks the Host header against DNS rebinding, and
// checks Origin and Referer on writes only, so a link from another site can
// still open a page.
func ui2HeaderCheck(c muxConfig, h http.Handler) http.Handler {
	hosted := hostCheck(apiVersion1, c.host, c.hostWhitelist, h)
	checked := originRefererCheck(apiVersion1, c.host, c.hostWhitelist, hosted)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			hosted.ServeHTTP(w, r)
			return
		}
		checked.ServeHTTP(w, r)
	})
}

// localTransport answers a request by calling a handler in process.
type localTransport struct{ h http.Handler }

func (t localTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	if r.Body == nil {
		r.Body = http.NoBody
	}
	rec := &recorder{header: http.Header{}, status: http.StatusOK}
	t.h.ServeHTTP(rec, r)
	return &http.Response{
		Status:        fmt.Sprintf("%d %s", rec.status, http.StatusText(rec.status)),
		StatusCode:    rec.status,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        rec.header,
		Body:          io.NopCloser(&rec.body),
		ContentLength: int64(rec.body.Len()),
		Request:       req,
	}, nil
}

type recorder struct {
	header      http.Header
	body        bytes.Buffer
	status      int
	wroteHeader bool
}

func (r *recorder) Header() http.Header { return r.header }

func (r *recorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status, r.wroteHeader = code, true
	}
}

func (r *recorder) Write(b []byte) (int, error) {
	r.wroteHeader = true
	return r.body.Write(b)
}
