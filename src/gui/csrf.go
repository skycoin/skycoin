package gui

import (
	"encoding/base64"
	"errors"
	"net/http"
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

	// check if token values are same
	if headerToken != c.token.String() {
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

		wh.SendOr404(w, &map[string]string{"csrf_token": store.getTokenValue()})
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
