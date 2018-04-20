// Package httphelper HTTP Error Response Helpers
package httphelper

import (
	"fmt"
	"net/http"

	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("gui")
)

// HTTPError wraps http.Error
func HTTPError(w http.ResponseWriter, status int, httpMsg string) {
	msg := fmt.Sprintf("%d %s", status, httpMsg)
	logger.Errorf(msg)
	http.Error(w, msg, status)
}

func httpError(w http.ResponseWriter, status int) {
	HTTPError(w, status, http.StatusText(status))
}

func errorXXXMsg(w http.ResponseWriter, status int, msg string) {
	httpMsg := http.StatusText(status)
	if msg != "" {
		httpMsg = fmt.Sprintf("%s - %s", httpMsg, msg)
	}
	HTTPError(w, status, httpMsg)
}

// Error400 respond with a 400 error
func Error400(w http.ResponseWriter, msg string) {
	errorXXXMsg(w, http.StatusBadRequest, msg)
}

// Error401 response with a 401 error
func Error401(w http.ResponseWriter, auth, msg string) {
	w.Header().Set("WWW-Authenticate", auth)
	httpMsg := http.StatusText(http.StatusUnauthorized)
	if msg != "" {
		httpMsg = fmt.Sprintf("%s - %s", httpMsg, msg)
	}
	HTTPError(w, http.StatusUnauthorized, httpMsg)
}

// Error403 respond with a 403 error
func Error403(w http.ResponseWriter) {
	httpError(w, http.StatusForbidden)
}

// Error403Msg respond with a 403 error and include a message
func Error403Msg(w http.ResponseWriter, msg string) {
	errorXXXMsg(w, http.StatusForbidden, msg)
}

// Error404 respond with a 404 error
func Error404(w http.ResponseWriter) {
	httpError(w, http.StatusNotFound)
}

// Error404Msg respond with a 404 error and include a message
func Error404Msg(w http.ResponseWriter, msg string) {
	errorXXXMsg(w, http.StatusNotFound, msg)
}

// Error405 respond with a 405 error
func Error405(w http.ResponseWriter) {
	httpError(w, http.StatusMethodNotAllowed)
}

// Error415 respond with a 415 error
func Error415(w http.ResponseWriter) {
	httpError(w, http.StatusUnsupportedMediaType)
}

// Error501 respond with a 501 error
func Error501(w http.ResponseWriter) {
	httpError(w, http.StatusNotImplemented)
}

// Error500 respond with a 500 error
func Error500(w http.ResponseWriter) {
	httpError(w, http.StatusInternalServerError)
}

// Error500Msg respond with a 500 error and include a message
func Error500Msg(w http.ResponseWriter, msg string) {
	errorXXXMsg(w, http.StatusInternalServerError, msg)
}

// Error503Msg respond with a 503 error and include a message
func Error503Msg(w http.ResponseWriter, msg string) {
	errorXXXMsg(w, http.StatusServiceUnavailable, msg)
}
