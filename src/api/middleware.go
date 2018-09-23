package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/iputil"
)

// ContentSecurityPolicy represents the value of content-security-policy
// header in http response
const ContentSecurityPolicy = "script-src 'self' 127.0.0.1"

// CSPHandler enables CSP
func CSPHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", ContentSecurityPolicy)
		handler.ServeHTTP(w, r)
	})
}

// HostCheck checks that the request's Host header is 127.0.0.1:$port or localhost:$port
// if the HTTP interface host is also a localhost address.
// This prevents DNS rebinding attacks, where an attacker uses a DNS rebinding service
// to bypass CORS checks.
// If the HTTP interface host is not a localhost address,
// the Host header is not checked. This is considered a public interface.
// If the Host header is not set, it is not checked.
// All major browsers send the Host header as required by the HTTP spec.
// hostWhitelist allows additional Host header values to be accepted.
func HostCheck(host string, hostWhitelist []string, handler http.Handler) http.Handler {
	addr := host
	var port uint16
	if strings.Contains(host, ":") {
		var err error
		addr, port, err = iputil.SplitAddr(host)
		if err != nil {
			logger.Panic(err)
		}
	}

	isLocalhost := iputil.IsLocalhost(addr)

	if isLocalhost && port == 0 {
		logger.Panic("localhost with no port specified is unsupported")
	}

	hostWhitelistMap := make(map[string]struct{}, len(hostWhitelist)+2)
	for _, k := range hostWhitelist {
		hostWhitelistMap[k] = struct{}{}
	}
	hostWhitelistMap[fmt.Sprintf("127.0.0.1:%d", port)] = struct{}{}
	hostWhitelistMap[fmt.Sprintf("localhost:%d", port)] = struct{}{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: The "Host" header is not in http.Request.Header, it's put in the http.Request.Host field
		_, isWhitelisted := hostWhitelistMap[r.Host]
		if isLocalhost && r.Host != "" && !isWhitelisted {
			logger.Critical().Errorf("Detected DNS rebind attempt - configured-host=%s header-host=%s", host, r.Host)
			wh.Error403(w, "Invalid Host")
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// OriginRefererCheck checks the Origin header if present, falling back on Referer.
// The Origin or Referer hostname must match the configured host.
// If neither are present, the request is allowed.  All major browsers will set
// at least one of these values. If neither are set, assume it is a request
// from curl/wget.
func OriginRefererCheck(host string, hostWhitelist []string, handler http.Handler) http.Handler {
	hostWhitelistMap := make(map[string]struct{}, len(hostWhitelist)+1)
	for _, k := range hostWhitelist {
		hostWhitelistMap[k] = struct{}{}
	}
	hostWhitelistMap[host] = struct{}{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")

		toCheck := origin
		if toCheck == "" {
			toCheck = referer
		}

		if toCheck != "" {
			u, err := url.Parse(toCheck)
			if err != nil {
				logger.Critical().Errorf("Invalid URL in Origin or Referer header: %s %v", toCheck, err)
				wh.Error403(w, "Invalid URL in Origin or Referer header")
				return
			}

			if _, isWhitelisted := hostWhitelistMap[u.Host]; !isWhitelisted {
				logger.Critical().Errorf("Origin or Referer header value %s does not match host and is not whitelisted", toCheck)
				wh.Error403(w, "Invalid Origin or Referer")
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}
