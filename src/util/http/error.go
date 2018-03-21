// Package httphelper HTTP Error Response Helpers
package httphelper

import (
	"fmt"
	"net/http"
)

// HTTPError wraps http.Error
func HTTPError(w http.ResponseWriter, status int, httpMsg string) {
	msg := fmt.Sprintf("%d %s", status, httpMsg)
	http.Error(w, msg, status)
}

// Error400 response 400 error
func Error400(w http.ResponseWriter, msg string) {
	httpMsg := "Bad Request"
	if msg != "" {
		httpMsg = fmt.Sprintf("%s - %s", httpMsg, msg)
	}
	HTTPError(w, http.StatusBadRequest, httpMsg)
}

// Error403 response 404 error
func Error403(w http.ResponseWriter) {
	HTTPError(w, http.StatusForbidden, "Forbidden")
}

func Error403Msg(w http.ResponseWriter, msg string) {
	httpMsg := "Forbidden"
	if msg != "" {
		httpMsg = fmt.Sprintf("%s - %s", httpMsg, msg)
	}
	HTTPError(w, http.StatusForbidden, httpMsg)
}

// Error404 response 404 error
func Error404(w http.ResponseWriter) {
	HTTPError(w, http.StatusNotFound, "Not Found")
}

// Error405 response 405
func Error405(w http.ResponseWriter) {
	HTTPError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}

// Error501 response 501
func Error501(w http.ResponseWriter) {
	HTTPError(w, http.StatusNotImplemented, "Not Implemented")
}

// Error500 response 500
func Error500(w http.ResponseWriter) {
	HTTPError(w, http.StatusInternalServerError, "Internal Server Error")
}

// Error500Msg response 500 with custom message
func Error500Msg(w http.ResponseWriter, msg string) {
	httpMsg := "Internal Server Error"
	if msg != "" {
		httpMsg = fmt.Sprintf("%s - %s", httpMsg, msg)
	}
	HTTPError(w, http.StatusInternalServerError, httpMsg)
}
