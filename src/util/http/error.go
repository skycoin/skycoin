// HTTP Error Response Helpers
package util_http

import (
	"net/http"
	"strings"
)

func HttpError(w http.ResponseWriter, status int, default_message string,
	messages []string) {
	message := default_message
	if len(messages) != 0 {
		message = strings.Join(messages, "<br>")
	}
	http.Error(w, message, status)
}

func Error400(w http.ResponseWriter, messages ...string) {
	HttpError(w, http.StatusBadRequest, "Bad request", messages)
}

func Error404(w http.ResponseWriter, messages ...string) {
	HttpError(w, http.StatusNotFound, "Not found", messages)
}

func Error405(w http.ResponseWriter, messages ...string) {
	HttpError(w, http.StatusMethodNotAllowed, "Method not allowed", messages)
}

func Error501(w http.ResponseWriter, messages ...string) {
	HttpError(w, http.StatusNotImplemented, "Not implemented", messages)
}

func Error500(w http.ResponseWriter, messages ...string) {
	HttpError(w, http.StatusInternalServerError, "Internal server error",
		messages)
}
