package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/skycoin/skycoin/src/util/iputil"
)

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
	return hostCheck(host, hostWhitelist, handler)
}

func hostCheck(host string, hostWhitelist []string, handler http.Handler) http.Handler {
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

	isWhiteListed := func(host string) bool {
		if _, ok := hostWhitelistMap[host]; ok {
			return true
		}

		if corsRegex.MatchString(host) {
			return true
		}

		return false
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: The "Host" header is not in http.Request.Header, it's put in the http.Request.Host field
		if isLocalhost && r.Host != "" && !isWhiteListed(r.Host) {
			logger.Critical().Errorf("Detected DNS rebind attempt - configured-host=%s header-host=%s", host, r.Host)
			resp := NewHTTPErrorResponse(http.StatusForbidden, "Invalid Host")
			writeHTTPResponse(w, resp)
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
	return originRefererCheck(host, hostWhitelist, handler)
}

func originRefererCheck(host string, hostWhitelist []string, handler http.Handler) http.Handler {
	hostWhitelistMap := make(map[string]struct{}, len(hostWhitelist)+2)
	for _, k := range hostWhitelist {
		hostWhitelistMap[k] = struct{}{}
	}

	if addr, port, _ := iputil.SplitAddr(host); iputil.IsLocalhost(addr) { // nolint: errcheck
		hostWhitelistMap[fmt.Sprintf("127.0.0.1:%d", port)] = struct{}{}
		hostWhitelistMap[fmt.Sprintf("localhost:%d", port)] = struct{}{}
		hostWhitelistMap["staging.wallet.skycoin.net"] = struct{}{}
		hostWhitelistMap["wallet.skycoin.net"] = struct{}{}

	} else {
		hostWhitelistMap[host] = struct{}{}
	}

	isWhiteListed := func(host string) bool {
		if _, ok := hostWhitelistMap[host]; ok {
			return true
		}

		if corsRegex.MatchString(host) {
			return true
		}

		return false
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		referer := r.Header.Get("Referer")
		toCheck := origin
		toCheckHeader := "Origin"
		if toCheck == "" {
			toCheck = referer
			toCheckHeader = "Referer"
		}

		if toCheck != "" {
			u, err := url.Parse(toCheck)
			if err != nil {
				logger.Critical().Errorf("Invalid URL in %s header: %s %v", toCheckHeader, toCheck, err)
				resp := NewHTTPErrorResponse(http.StatusForbidden, "Invalid URL in Origin or Referer header")
				writeHTTPResponse(w, resp)
				return
			}

			if !isWhiteListed(u.Host) {
				logger.Critical().Errorf("%s header value %s does not match host and is not whitelisted", toCheckHeader, toCheck)
				resp := NewHTTPErrorResponse(http.StatusForbidden, "Invalid Origin or Referer")
				writeHTTPResponse(w, resp)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}
