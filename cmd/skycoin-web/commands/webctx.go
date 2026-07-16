// Package commands provides commands for the skycoin web interface.
//
// This file provides a tiny net/http request-context shim used by all of the
// web handlers. The web interface was originally written against
// github.com/gin-gonic/gin, but gin unconditionally pulls in
// github.com/quic-go/quic-go/http3 and github.com/ugorji/go/codec, neither of
// which works under TinyGo (quic-go needs QUIC TLS APIs TinyGo lacks, and
// ugorji indexes a cache by reflect.Kind, whose numbering differs under
// TinyGo). Using only net/http keeps a single code path that builds and runs
// under both the standard Go toolchain and TinyGo.
//
// webCtx mirrors the handful of *gin.Context methods the handlers used
// (Header/Param/Status/String/Data/JSON plus the Request/Writer fields) so the
// handler bodies stay unchanged.
package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H is a shorthand for a map, replacing gin.H.
type H map[string]any

// webCtx is a minimal per-request context over net/http.
type webCtx struct {
	Writer  http.ResponseWriter
	Request *http.Request
	params  map[string]string
}

// newCtx builds a webCtx for a request. Path wildcards captured by the
// ServeMux (via r.PathValue) are copied into params by the caller when needed.
func newCtx(w http.ResponseWriter, r *http.Request) *webCtx {
	return &webCtx{Writer: w, Request: r}
}

// Param returns a routing parameter captured from the URL path.
func (c *webCtx) Param(key string) string {
	return c.params[key]
}

// Query returns a URL query parameter.
func (c *webCtx) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// DefaultQuery returns a URL query parameter, or defaultValue if absent.
func (c *webCtx) DefaultQuery(key, defaultValue string) string {
	if v, ok := c.Request.URL.Query()[key]; ok && len(v) > 0 {
		return v[0]
	}
	return defaultValue
}

// ShouldBindJSON decodes the request body as JSON into obj.
func (c *webCtx) ShouldBindJSON(obj any) error {
	return json.NewDecoder(c.Request.Body).Decode(obj)
}

// Header sets a response header.
func (c *webCtx) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}

// Status writes an HTTP status code with no body.
func (c *webCtx) Status(code int) {
	c.Writer.WriteHeader(code)
}

// String writes a formatted plain-text response.
func (c *webCtx) String(code int, format string, args ...any) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(code)
	if len(args) == 0 {
		_, _ = c.Writer.Write([]byte(format)) //nolint:errcheck
		return
	}
	_, _ = fmt.Fprintf(c.Writer, format, args...) //nolint:errcheck,gosec
}

// Data writes a raw response body with the given content type.
func (c *webCtx) Data(code int, contentType string, data []byte) {
	if contentType != "" {
		c.Writer.Header().Set("Content-Type", contentType)
	}
	c.Writer.WriteHeader(code)
	_, _ = c.Writer.Write(data) //nolint:errcheck
}

// JSON marshals obj and writes it as an application/json response.
func (c *webCtx) JSON(code int, obj any) {
	b, err := json.Marshal(obj)
	if err != nil {
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(c.Writer, "500 Internal Server Error - json marshal: %v", err) //nolint:errcheck
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(code)
	_, _ = c.Writer.Write(b) //nolint:errcheck
}

// recoverMiddleware wraps a handler so a panic becomes a 500 rather than a
// dropped connection, matching gin.Default's Recovery behavior.
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "500 Internal Server Error - %v", rec) //nolint:errcheck
			}
		}()
		next.ServeHTTP(w, r)
	})
}
