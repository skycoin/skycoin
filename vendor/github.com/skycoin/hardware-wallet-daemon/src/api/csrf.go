package api

import (
	"net/http"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

const (
	// CSRFHeaderName is the name of the CSRF header
	CSRFHeaderName = "X-CSRF-Token"

	// CSRFMaxAge is the lifetime of a CSRF token in seconds
	CSRFMaxAge = time.Second * 30

	csrfSecretLength = 64

	csrfNonceLength = 64
)

var (
	// ErrCSRFInvalid is returned when the the CSRF token is in invalid format
	ErrCSRFInvalid = errors.New("invalid CSRF token")
	// ErrCSRFInvalidSignature is returned when the signature of the csrf token is invalid
	ErrCSRFInvalidSignature = errors.New("invalid CSRF token signature")
	// ErrCSRFExpired is returned when the csrf token has expired
	ErrCSRFExpired = errors.New("csrf token expired")
)

var csrfSecretKey []byte

func init() {
	csrfSecretKey = cipher.RandByte(csrfSecretLength)
}

// CSRFToken csrf token
type CSRFToken struct {
	Nonce     []byte
	ExpiresAt time.Time
}

// newCSRFToken generates a new CSRF Token
func newCSRFToken() (string, error) {
	return newCSRFTokenWithTime(time.Now().Add(CSRFMaxAge))
}

func newCSRFTokenWithTime(expiresAt time.Time) (string, error) {
	token := &CSRFToken{
		Nonce:     cipher.RandByte(csrfNonceLength),
		ExpiresAt: expiresAt,
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, csrfSecretKey)
	_, err = h.Write([]byte(tokenJSON))
	if err != nil {
		return "", err
	}

	sig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	signingString := base64.RawURLEncoding.EncodeToString(tokenJSON)

	return strings.Join([]string{signingString, sig}, "."), nil
}

// verifyCSRFToken checks validity of the given token
func verifyCSRFToken(headerToken string) error {
	tokenParts := strings.Split(headerToken, ".")
	if len(tokenParts) != 2 {
		return ErrCSRFInvalid
	}

	signingString, err := base64.RawURLEncoding.DecodeString(tokenParts[0])
	if err != nil {
		return err
	}

	h := hmac.New(sha256.New, csrfSecretKey)
	_, err = h.Write([]byte(signingString))
	if err != nil {
		return err
	}

	sig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if sig != tokenParts[1] {
		return ErrCSRFInvalidSignature
	}

	var csrfToken CSRFToken
	err = json.Unmarshal(signingString, &csrfToken)
	if err != nil {
		return err
	}

	if time.Now().After(csrfToken.ExpiresAt) {
		return ErrCSRFExpired
	}

	return nil
}

// Creates a new CSRF token. Previous CSRF tokens are invalidated by this call.
// URI: /api/v1/csrf
// Method: GET
// Response:
//  csrf_token: CSRF token to use in POST requests
func getCSRFToken(enabled bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		if !enabled {
			logger.Warning("CSRF check disabled")
			resp := NewHTTPErrorResponse(http.StatusNotFound, "")
			writeHTTPResponse(w, resp)
			return
		}

		// generate a new token
		csrfToken, err := newCSRFToken()
		if err != nil {
			logger.Error(err)
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to create a csrf token: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: csrfToken,
		})
	}
}

// CSRFCheck verifies X-CSRF-Token header value
func CSRFCheck(enabled bool, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enabled {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodDelete:
				token := r.Header.Get(CSRFHeaderName)
				if err := verifyCSRFToken(token); err != nil {
					logger.Errorf("CSRF token invalid: %v", err)
					resp := NewHTTPErrorResponse(http.StatusForbidden, err.Error())
					writeHTTPResponse(w, resp)
					return
				}
			}
		}

		handler.ServeHTTP(w, r)
	})
}
