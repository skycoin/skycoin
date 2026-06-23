//go:build tinygo

package api

import "net/http"

// newHTTPClient builds the HTTP client used by Client. TinyGo's
// http.Transport does not expose the Dial or TLSHandshakeTimeout fields, so
// only the overall client timeout is configured.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: httpClientTimeout,
	}
}
