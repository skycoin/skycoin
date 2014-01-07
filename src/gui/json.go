// Utilities for sending JSON
package gui

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
    out, err := json.Marshal(message)
    if err == nil {
        _, err := w.Write(out)
        if err != nil {
            return err
        }
    }
    return err
}
