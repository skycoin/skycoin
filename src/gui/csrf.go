package gui

import (
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

const (
	// the name of CSRF header
	CSRFHeaderName = "X-CSRF-Token"

	// Max-Age in seconds for cookie. 30 seconds default
	CSRFMaxAge = time.Duration(30) * time.Second

	// csrf token length
	csrfTokenLength = 64
)

type CSRFToken struct {
	Value     []byte
	ExpiresAt time.Time
}

type CSRFStore struct {
	token   *CSRFToken
	Enabled bool
	sync.RWMutex
}

// getTokenValue returns a url safe base64 encoded token
func (c *CSRFStore) getTokenValue() string {
	c.RLock()
	defer c.RUnlock()
	return base64.RawURLEncoding.EncodeToString(c.token.Value)
}

// setToken sets a new CSRF token
// if the value is changing the expire time should also change
// so there is no explicit method to just set the value of the token
func (c *CSRFStore) setToken(token CSRFToken) {
	c.Lock()
	defer c.Unlock()
	c.token = &token
}

// verifyExpireTime checks if token expiry time is greater than current time
func (c *CSRFStore) verifyExpireTime() bool {
	return c.token.ExpiresAt.After(time.Now())

}

// verifyToken checks that the given token is same as the internal token
func (c *CSRFStore) verifyToken(headerToken string) bool {
	c.RLock()
	defer c.RUnlock()

	// check if token is initialized
	if c.token == nil {
		return false
	}

	// check if token values are same
	if headerToken == c.getTokenValue() {
		// make sure token is still valid
		return c.verifyExpireTime()
	}

	return false
}

// method: GET
// url: /csrf
func getCSRFToken(gateway Gatewayer, store *CSRFStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || !store.Enabled {
			wh.Error404(w)
			return
		}

		// generate a new token
		store.setToken(generateToken())

		wh.SendOr404(w, &map[string]string{"csrf_token": store.getTokenValue()})
	}
}

// CSRFCheck verifies X-CSRF-Token header value
func CSRFCheck(handler http.Handler, store *CSRFStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && store.Enabled {
			token := r.Header.Get(CSRFHeaderName)
			if !store.verifyToken(token) {
				wh.Error403Msg(w, "invalid CSRF token")
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

// generateToken generates a new CSRF Token
func generateToken() CSRFToken {
	bytes := cipher.RandByte(csrfTokenLength)

	token := CSRFToken{
		bytes,
		time.Now().Add(CSRFMaxAge),
	}

	return token
}
