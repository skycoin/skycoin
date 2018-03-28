package httphelper

//  Utilities for sending JSON

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/skycoin/src/util/logging"
)

// SendJSON emits JSON to an http response
func SendJSON(w http.ResponseWriter, m interface{}) error {
	out, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")

	if _, err := w.Write(out); err != nil {
		return err
	}

	return nil
}

// SendJSONOr500 writes an object as JSON, writing a 500 error if it fails
func SendJSONOr500(log *logging.Logger, w http.ResponseWriter, m interface{}) {
	if err := SendJSON(w, m); err != nil {
		log.Errorf("%v", err)
		Error500(w)
	}
}
