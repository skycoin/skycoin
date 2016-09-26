// Utilities for sending JSON
package util_http

import (
	"encoding/json"
	"net/http"
)

type JSONMessage interface{}

// Simple JSON response wrapper
type JSONResponse struct {
	Message string
}

// Returns a JSONResponse conforming to JSONMessage
func NewJSONResponse(message string) JSONMessage {
	return &JSONResponse{Message: message}
}

// Emits JSON to an http response
func SendJSON(w http.ResponseWriter, message JSONMessage) error {
	out, err := json.MarshalIndent(message, "", "    ")
	if err == nil {
		_, err := w.Write(out)
		if err != nil {
			return err
		}
	}
	return err
}

// Sends an interface as JSON if its not nil (404) or fails (500)
func SendOr404(w http.ResponseWriter, m interface{}) {
	if m == nil {
		Error404(w)
	} else if SendJSON(w, m) != nil {
		Error500(w)
	}
}

// Sends an interface as JSON if its not nil (500) or fails (500)
func SendOr500(w http.ResponseWriter, m interface{}) {
	if m == nil {
		Error500(w)
	} else if SendJSON(w, m) != nil {
		Error500(w)
	}
}
