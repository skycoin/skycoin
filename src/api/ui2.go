// Package api src/api/ui2.go c3-coin-api
package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

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

// mountUI2 serves the page handler that build makes at /ui2/. Its client
// calls this mux in process, so the pages go through the same API set, auth
// and CSRF checks as any other client, without a network round trip.
func mountUI2(mux *http.ServeMux, c muxConfig, build func(*Client) http.Handler) {
	client := HandlerClient(mux, c.host)
	client.SetAuth(c.username, c.password)
	h := http.StripPrefix("/ui2", build(client))
	if !c.disableCSP {
		h = CSPHandler(h, ContentSecurityPolicy)
	}
	h = CSRFCheck(apiVersion1, c.disableCSRF, h)
	if !c.disableHeaderCheck {
		h = ui2HeaderCheck(c, h)
	}
	h = basicAuth(apiVersion1, c.username, c.password, "skycoin daemon", h)
	mux.Handle("/ui2/", gzhttp.GzipHandler(h))
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
