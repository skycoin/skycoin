package commands

import (
	"net/http"
	"net/url"
)

// sameOriginWrites refuses requests that can change state when a browser says
// they come from another site, so no other page can send coins or edit wallets.
// It follows net/http's CrossOriginProtection, which TinyGo's net/http lacks.
func sameOriginWrites(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !crossOriginWrite(r) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "cross-origin request refused", http.StatusForbidden)
	})
}

func crossOriginWrite(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	}
	switch r.Header.Get("Sec-Fetch-Site") {
	case "same-origin", "none":
		return false
	case "":
	default:
		return true
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Not a browser, or one too old to say: only a browser can be tricked.
		return false
	}
	u, err := url.Parse(origin)
	return err != nil || u.Host != r.Host
}
