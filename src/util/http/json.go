package httphelper

//  Utilities for sending JSON

import (
	"encoding/json"
	"net/http"
)

// JSONMessage json message
type JSONMessage interface{}

// JSONResponse simple JSON response wrapper
type JSONResponse struct {
	Message string
}

// NewJSONResponse returns a JSONResponse conforming to JSONMessage
func NewJSONResponse(message string) JSONMessage {
	return &JSONResponse{Message: message}
}

// SendJSON emits JSON to an http response
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

// SendOr404 sends an interface as JSON if its not nil (404) or fails (500)
func SendOr404(w http.ResponseWriter, m interface{}) {
	if m == nil {
		Error404(w)
	} else if SendJSON(w, m) != nil {
		Error500(w)
	}
}

// SendOr500 sends an interface as JSON if its not nil (500) or fails (500)
func SendOr500(w http.ResponseWriter, m interface{}) {
	if m == nil {
		Error500(w)
	} else if SendJSON(w, m) != nil {
		Error500(w)
	}
}
