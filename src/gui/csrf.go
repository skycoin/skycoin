package gui

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

const (
	// CSRFHeaderName is the name of the CSRF header
	CSRFHeaderName = "X-CSRF-Token"

	// CSRFMaxAge is the lifetime of a CSRF token in seconds
	CSRFMaxAge = time.Second * 30

	csrfTokenLength = 64
)

// CSRFToken csrf token
type CSRFToken struct {
	Value     []byte
	ExpiresAt time.Time
}

// newCSRFToken generates a new CSRF Token
func newCSRFToken() CSRFToken {
	return CSRFToken{
		Value:     cipher.RandByte(csrfTokenLength),
		ExpiresAt: time.Now().Add(CSRFMaxAge),
	}
}

// String returns the token in base64 URL-safe encoded format
func (c *CSRFToken) String() string {
	return base64.RawURLEncoding.EncodeToString(c.Value)
}

// CSRFStore encapsulates a single CSRFToken
type CSRFStore struct {
	token   *CSRFToken
	Enabled bool
	sync.RWMutex
}

// getTokenValue returns a url safe base64 encoded token
func (c *CSRFStore) getTokenValue() string {
	c.RLock()
	defer c.RUnlock()
	return c.token.String()
}

// setToken sets a new CSRF token
// if the value is changing the expire time should also change
// so there is no explicit method to just set the value of the token
func (c *CSRFStore) setToken(token CSRFToken) {
	c.Lock()
	defer c.Unlock()
	c.token = &token
}

// expired checks if token expiry time is greater than current time
func (c *CSRFStore) expired() bool {
	return c.token == nil || time.Now().After(c.token.ExpiresAt)

}

// verifyToken checks that the given token is same as the internal token
func (c *CSRFStore) verifyToken(headerToken string) error {
	c.RLock()
	defer c.RUnlock()

	// check if token is initialized
	if c.token == nil || len(c.token.Value) == 0 {
		return errors.New("token not initialized")
	}

	a, err := base64.RawURLEncoding.DecodeString(headerToken)
	if err != nil {
		return err
	}

	// check if token values are same, using a constant time comparison
	if subtle.ConstantTimeCompare(a, c.token.Value) != 1 {
		return errors.New("invalid token")
	}

	// make sure token is still valid
	if c.expired() {
		return errors.New("token has expired")
	}

	return nil
}

// Creates a new CSRF token. Previous CSRF tokens are invalidated by this call.
// URI: /csrf
// Method: GET
// Response:
//  csrf_token: CSRF token to use in POST requests
func getCSRFToken(gateway Gatewayer, store *CSRFStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		if !store.Enabled {
			logger.Warning("CSRF check disabled")
			wh.Error404(w)
			return
		}

		// generate a new token
		store.setToken(newCSRFToken())

		wh.SendJSONOr500(logger, w, &map[string]string{"csrf_token": store.getTokenValue()})
	}
}

// CSRFCheck verifies X-CSRF-Token header value
func CSRFCheck(store *CSRFStore, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if store.Enabled {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodDelete:
				token := r.Header.Get(CSRFHeaderName)
				if err := store.verifyToken(token); err != nil {
					logger.Errorf("CSRF token invalid: %v", err)
					wh.Error403Msg(w, "invalid CSRF token")
					return
				}
			}
		}

		handler.ServeHTTP(w, r)
	})
}

// OriginRefererCheck checks the Origin header if present, falling back on Referer.
// The Origin or Referer hostname must match the configured host.
// If neither are present, the request is allowed.  All major browsers will set
// at least one of these values. If neither are set, assume it is a request
// from curl/wget.
func OriginRefererCheck(host string, handler http.Handler) http.Handler {
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
				logger.Critical("Invalid URL in Origin or Referer header: %s %v", toCheck, err)
				wh.Error403(w)
				return
			}

			if u.Host != host {
				logger.Critical("Origin or Referer header value %s does not match host", toCheck)
				wh.Error403(w)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}
