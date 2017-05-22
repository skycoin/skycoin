// Package httphelper HTTP Error Response Helpers
package httphelper

import (
	"net/http"
	"strings"
)

// HTTPError wraps http.Error
func HTTPError(w http.ResponseWriter, status int, defaultMsg string,
	messages []string) {
	message := defaultMsg
	if len(messages) != 0 {
		message = strings.Join(messages, "<br>")
	}
	http.Error(w, message, status)
}

// Error400 response 400 error
func Error400(w http.ResponseWriter, messages ...string) {
	HTTPError(w, http.StatusBadRequest, "Bad request", messages)
}

// Error404 response 404 error
func Error404(w http.ResponseWriter, messages ...string) {
	HTTPError(w, http.StatusNotFound, "Not found", messages)
}

// Error405 response 405
func Error405(w http.ResponseWriter, messages ...string) {
	HTTPError(w, http.StatusMethodNotAllowed, "Method not allowed", messages)
}

// Error501 response 501
func Error501(w http.ResponseWriter, messages ...string) {
	HTTPError(w, http.StatusNotImplemented, "Not implemented", messages)
}

// Error500 response 500
func Error500(w http.ResponseWriter, messages ...string) {
	HTTPError(w, http.StatusInternalServerError, "Internal server error",
		messages)
}
