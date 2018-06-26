package cipher

import (
	"errors"
)

var (
	// ErrInvalidLength Unexpected size of address bytes buffer
	ErrInvalidLength = errors.New("Invalid address length")
	// ErrInvalidChecksum Computed checksum did not match expected value
	ErrInvalidChecksum = errors.New("Invalid checksum")
	// ErrInvalidVersion Unsupported address version value
	ErrInvalidVersion = errors.New("Invalid version")
	// ErrInvalidPubKey Public key invalid for address
	ErrInvalidPubKey = errors.New("Public key invalid for address")
	// ErrInvalidFirstByte Invalid first byte in wallet import format string
	ErrInvalidFirstByte = errors.New("first byte invalid")
	// ErrInvalidLastByte 33rd byte in wallet import format string is invalid
	ErrInvalidLastByte = errors.New("invalid 33rd byte")
)
