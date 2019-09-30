package httphelper

//  Utilities for sending JSON

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/droplet"
	"github.com/SkycoinProject/skycoin/src/util/logging"
)

// SendJSONOr500 writes an object as JSON, writing a 500 error if it fails
func SendJSONOr500(log *logging.Logger, w http.ResponseWriter, m interface{}) {
	out, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		Error500(w, "json.MarshalIndent failed")
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if _, err := w.Write(out); err != nil {
		log.WithError(err).Error("http Write failed")
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

// Address is a wrapper around cipher.Address which implements json.Unmarshaler and json.Marshaler.
// It marshals and unmarshals the address as a string
type Address struct {
	cipher.Address
}

// UnmarshalJSON unmarshals a string address to a cipher.Address
func (a *Address) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, err := cipher.DecodeBase58Address(s)
	if err != nil {
		return fmt.Errorf("invalid address: %v", err)
	}

	a.Address = tmp

	return nil
}

// MarshalJSON marshals a cipher.Address in its string representation
func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.Address.String() + `"`), nil
}

// SHA256 is a wrapper around cipher.SHA256 which implements json.Unmarshaler and json.Marshaler.
// It marshals and unmarshals the address as a string
type SHA256 struct {
	cipher.SHA256
}

// UnmarshalJSON unmarshals a string address to a cipher.SHA256
func (a *SHA256) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, err := cipher.SHA256FromHex(s)
	if err != nil {
		return fmt.Errorf("invalid SHA256 hash: %v", err)
	}

	a.SHA256 = tmp

	return nil
}

// MarshalJSON marshals a cipher.SHA256 in its string representation
func (a SHA256) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.SHA256.Hex() + `"`), nil
}

// Coins is a wrapper around uint64 which implements json.Unmarshaler and json.Marshaler.
// It unmarshals a fixed-point decimal string to droplets and vice versa
type Coins uint64

// UnmarshalJSON unmarshals a fixed-point decimal string to droplets
func (c *Coins) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, err := droplet.FromString(s)
	if err != nil {
		return err
	}

	*c = Coins(tmp)

	return nil
}

// MarshalJSON marshals droplets to a fixed-point decimal string
func (c Coins) MarshalJSON() ([]byte, error) {
	s, err := droplet.ToString(uint64(c))
	if err != nil {
		return nil, err
	}

	return []byte(`"` + s + `"`), nil
}

// Value returns the underlying uint64 value
func (c Coins) Value() uint64 {
	return uint64(c)
}

// Hours is a wrapper around uint64 which implements json.Unmarshaler and json.Marshaler.
// It unmarshals a fixed-point decimal string to droplets and vice versa
type Hours uint64

// UnmarshalJSON unmarshals a fixed-point decimal string to droplets
func (h *Hours) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tmp, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid hours value: %v", err)
	}

	*h = Hours(tmp)

	return nil
}

// MarshalJSON marshals droplets to a fixed-point decimal string
func (h Hours) MarshalJSON() ([]byte, error) {
	s := fmt.Sprint(h)
	return []byte(`"` + s + `"`), nil
}

// Value returns the underlying uint64 value
func (h Hours) Value() uint64 {
	return uint64(h)
}
