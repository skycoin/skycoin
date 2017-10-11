package httphelper

//  Utilities for sending JSON

import (
	"encoding/json"
	"net/http"
)

// SendJSON emits JSON to an http response
func SendJSON(w http.ResponseWriter, m interface{}) error {
	out, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}

	if _, err := w.Write(out); err != nil {
		return err
	}

	return nil
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
