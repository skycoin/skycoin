package gui

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/util/utc"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

const (
	// the name of CSRF header
	HeaderName = "X-CSRF-Token"

	// Max-Age in seconds for cookie. 30 seconds default
	MaxAge = time.Duration(30) * time.Second

	// csrf token length
	tokenLength = 64
)

type CSRFToken struct {
	Value     []byte
	ExpiresAt time.Time
}

type CSRFStore struct {
	token *CSRFToken
	sync.Mutex
	once sync.Once
}

// getTokenValue returns a url safe base64 encoded token
func (c *CSRFStore) getTokenValue() string {
	return b64encode(c.token.Value)
}

// setToken sets a new CSRF token
// if the value is changing the expire time should also change
// so there is no explicit method to just set the value of the token
func (c *CSRFStore) setToken(token *CSRFToken) {
	c.Lock()
	defer c.Unlock()
	c.token = token
}

// verifyExpireTime checks if token expiry time is greater than current time
func (c *CSRFStore) verifyExpireTime() bool {
	return utc.UnixNow() < c.token.ExpiresAt.Unix()

}

// verifyToken checks that the given token is same as the internal token
func (c *CSRFStore) verifyToken(headerToken string) bool {
	// check if token values are same
	if headerToken == c.getTokenValue() {
		// make sure token is still valid
		return c.verifyExpireTime()
	}

	return false
}

// getCSRFStore returns a CSRFStore instance
func (c *CSRFStore) getCSRFStore() *CSRFStore {
	c.once.Do(func() {
		// initialize the csrf store
		if c.token == nil {
			// intialize the csrf token
			c.token = generateToken()
		}
	})

	return c
}

// method: GET
// url: /csrf
func getCSRFToken(gateway Gatewayer, store *CSRFStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		// check if token is still valid
		// otherwise generate new token
		if !store.verifyExpireTime() {
			store.setToken(generateToken())
		}

		wh.SendOr404(w, store.getTokenValue())

	}
}

// CSRFCheck verifies X-CSRF-Token header value
func CSRFCheck(handler http.Handler, store *CSRFStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HeaderName)
		if !store.verifyToken(token) {
			wh.Error403Msg(w, "invalid CSRF token")
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// b64encode returns a url safe base64 encoded string
func b64encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// generateToken generates a new CSRF Token
func generateToken() *CSRFToken {
	bytes := make([]byte, tokenLength)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic(err)
	}

	token := CSRFToken{
		bytes,
		utc.Now().Add(MaxAge),
	}

	return &token
}

func checkForPRNG() {
	buf := make([]byte, 1)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

func init() {
	checkForPRNG()
}
