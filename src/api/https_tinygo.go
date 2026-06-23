//go:build tinygo

package api

import "errors"

// ErrHTTPSUnsupported is returned when HTTPS is requested in a TinyGo build.
// TinyGo's crypto/tls does not provide tls.Listen, so the HTTPS web interface
// is unavailable; use the plain HTTP interface instead.
var ErrHTTPSUnsupported = errors.New("HTTPS web interface is not supported in TinyGo builds")

// CreateHTTPS is unsupported under TinyGo and always returns
// ErrHTTPSUnsupported.
func CreateHTTPS(host string, c Config, gateway Gatewayer, certFile, keyFile string) (*Server, error) {
	return nil, ErrHTTPSUnsupported
}
