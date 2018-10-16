// Package httphelper provides HTTP related utility methods
package httphelper

import (
	"fmt"
	"net/http"
)

// ErrorXXX writes an error message with status code
func ErrorXXX(w http.ResponseWriter, status int, msg string) {
	httpMsg := fmt.Sprintf("%d %s", status, http.StatusText(status))
	if msg != "" {
		httpMsg = fmt.Sprintf("%s - %s", httpMsg, msg)
	}

	http.Error(w, httpMsg, status)
}

// Error400 respond with a 400 error and include a message
func Error400(w http.ResponseWriter, msg string) {
	ErrorXXX(w, http.StatusBadRequest, msg)
}

// Error401 respond with a 401 error
func Error401(w http.ResponseWriter, auth, msg string) {
	w.Header().Set("WWW-Authenticate", auth)
	ErrorXXX(w, http.StatusUnauthorized, msg)
}

// Error403 respond with a 403 error and include a message
func Error403(w http.ResponseWriter, msg string) {
	ErrorXXX(w, http.StatusForbidden, msg)
}

// Error404 respond with a 404 error and include a message
func Error404(w http.ResponseWriter, msg string) {
	ErrorXXX(w, http.StatusNotFound, msg)
}

// Error405 respond with a 405 error
func Error405(w http.ResponseWriter) {
	ErrorXXX(w, http.StatusMethodNotAllowed, "")
}

// Error415 respond with a 415 error
func Error415(w http.ResponseWriter) {
	ErrorXXX(w, http.StatusUnsupportedMediaType, "")
}

// Error422 response with a 422 error and include a message
func Error422(w http.ResponseWriter, msg string) {
	ErrorXXX(w, http.StatusUnprocessableEntity, msg)
}

// Error500 respond with a 500 error and include a message
func Error500(w http.ResponseWriter, msg string) {
	ErrorXXX(w, http.StatusInternalServerError, msg)
}

// Error501 respond with a 501 error
func Error501(w http.ResponseWriter) {
	ErrorXXX(w, http.StatusNotImplemented, "")
}

// Error503 respond with a 503 error and include a message
func Error503(w http.ResponseWriter, msg string) {
	ErrorXXX(w, http.StatusServiceUnavailable, msg)
}
