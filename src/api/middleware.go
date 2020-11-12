package api

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/iputil"
)

// ContentSecurityPolicy represents the value of content-security-policy
// header in http response
const ContentSecurityPolicy = "default-src 'self'" +
	"; connect-src 'self' https://api.coinpaprika.com https://swaplab.cc https://version.skycoin.com https://downloads.skycoin.com http://127.0.0.1:9510" +
	"; img-src 'self' 'unsafe-inline' data:" +
	"; style-src 'self' 'unsafe-inline'" +
	"; object-src	'none'" +
	"; form-action 'none'" +
	"; frame-ancestors 'none'" +
	"; block-all-mixed-content" +
	"; base-uri 'self'"

// CSPHandler sets the Content-Security-Policy header
func CSPHandler(handler http.Handler, policy string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", policy)
		handler.ServeHTTP(w, r)
	})
}

// ContentTypeJSONRequired enforces Content-Type: application/json in a POST request.
// Return 415 Unsupported Media Type if the Content-Type is not application/json,
// in the V2 error format.
func ContentTypeJSONRequired(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			contentType := r.Header.Get("Content-Type")
			if !isContentTypeJSON(contentType) {
				resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
				writeHTTPResponse(w, resp)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

// isContentTypeJSON returns true if the content type is application/json,
// allowing the content-type string to include extra parameters like charset=utf-8,
// for example `Content-Type: application/json; charset=utf-8` will return true.
func isContentTypeJSON(contentType string) bool {
	return contentType == ContentTypeJSON || strings.HasPrefix(contentType, ContentTypeJSON+";")
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
	return hostCheck(apiVersion1, host, hostWhitelist, handler)
}

func hostCheck(apiVersion, host string, hostWhitelist []string, handler http.Handler) http.Handler {
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
			writeError(w, apiVersion, http.StatusForbidden, "Invalid Host")
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
	return originRefererCheck(apiVersion1, host, hostWhitelist, handler)
}

func originRefererCheck(apiVersion, host string, hostWhitelist []string, handler http.Handler) http.Handler {
	hostWhitelistMap := make(map[string]struct{}, len(hostWhitelist)+2)
	for _, k := range hostWhitelist {
		hostWhitelistMap[k] = struct{}{}
	}

	if addr, port, _ := iputil.SplitAddr(host); iputil.IsLocalhost(addr) { //nolint:errcheck
		hostWhitelistMap[fmt.Sprintf("127.0.0.1:%d", port)] = struct{}{}
		hostWhitelistMap[fmt.Sprintf("localhost:%d", port)] = struct{}{}
	} else {
		hostWhitelistMap[host] = struct{}{}
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
				writeError(w, apiVersion, http.StatusForbidden, "Invalid URL in Origin or Referer header")
				return
			}

			if _, isWhitelisted := hostWhitelistMap[u.Host]; !isWhitelisted {
				logger.Critical().Errorf("%s header value %s does not match host and is not whitelisted", toCheckHeader, toCheck)
				writeError(w, apiVersion, http.StatusForbidden, "Invalid Origin or Referer")
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

func basicAuth(apiVersion, username, password, realm string, f http.Handler) http.HandlerFunc {
	needsAuth := username != "" || password != ""
	usernamePasswordHash := cipher.SumSHA256(append([]byte(username), []byte(password)...))
	authHeader := fmt.Sprintf("Basic realm=%q", realm)

	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if needsAuth {
			if !ok {
				w.Header().Set("WWW-Authenticate", authHeader)
				writeError(w, apiVersion, http.StatusUnauthorized, "")
				return
			}

			userPassHash := cipher.SumSHA256(append([]byte(user), []byte(pass)...))

			if subtle.ConstantTimeCompare(userPassHash[:], usernamePasswordHash[:]) != 1 {
				w.Header().Set("WWW-Authenticate", authHeader)
				writeError(w, apiVersion, http.StatusUnauthorized, "")
				return
			}
		} else {
			// If auth is not configured but the request provides auth, reject
			// This will avoid a mistake where the daemon is not configured with auth,
			// but the client is, and does not realize the daemon is not configured with auth
			// because all requests are accepted
			if user != "" || pass != "" {
				w.Header().Set("WWW-Authenticate", authHeader)
				writeError(w, apiVersion, http.StatusUnauthorized, "")
				return
			}
		}

		f.ServeHTTP(w, r)
	}
}

func writeError(w http.ResponseWriter, apiVersion string, code int, msg string) {
	switch apiVersion {
	case apiVersion1:
		wh.ErrorXXX(w, code, msg)
	case apiVersion2:
		writeHTTPResponse(w, NewHTTPErrorResponse(code, msg))
	default:
		wh.Error500(w, "Invalid internal API version")
	}
}
