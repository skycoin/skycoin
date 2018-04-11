package httphelper

//  Utilities for sending JSON

import (
	"encoding/json"
	"net/http"
	"time"

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

// Duration JSON type copied from https://github.com/vrischmann/jsonutil, MIT License

// Duration is a wrapper around time.Duration which implements json.Unmarshaler and json.Marshaler.
// It marshals and unmarshals the duration as a string in the format accepted by time.ParseDuration and returned by time.Duration.String.
type Duration struct {
	time.Duration
}

// FromDuration is a convenience factory to create a Duration instance from the
// given time.Duration value.
func FromDuration(d time.Duration) Duration {
	return Duration{d}
}

// MarshalJSON implements the json.Marshaler interface. The duration is a quoted-string in the format accepted by time.ParseDuration and returned by time.Duration.String.
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. The duration is expected to be a quoted-string of a duration in the format accepted by time.ParseDuration.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	tmp, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	d.Duration = tmp

	return nil
}
