//go:build !tinygo

package api

import (
	"net"
	"net/http"
)

// newHTTPClient builds the HTTP client used by Client with explicit dial and
// TLS handshake timeouts.
func newHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: dialTimeout,
		}).DialContext,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   httpClientTimeout,
	}
}
