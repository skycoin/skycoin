// Package mathutil provides math utilities
package mathutil

import (
	"errors"
	"math"
)

var (
	// ErrUint64MultOverflow is returned if when multiplying uint64 values would overflow uint64
	ErrUint64MultOverflow = errors.New("uint64 multiplication overflow")
	// ErrUint64AddOverflow is returned if when adding uint64 values would overflow uint64
	ErrUint64AddOverflow = errors.New("uint64 addition overflow")
	// ErrUint32AddOverflow is returned if when adding uint32 values would overflow uint32
	ErrUint32AddOverflow = errors.New("uint32 addition overflow")
	// ErrUint64OverflowsInt64 is returned if when converting a uint64 to an int64 would overflow int64
	ErrUint64OverflowsInt64 = errors.New("uint64 overflows int64")
	// ErrInt64UnderflowsUint64 is returned if when converting an int64 to a uint64 would underflow uint64
	ErrInt64UnderflowsUint64 = errors.New("int64 underflows uint64")
	// ErrIntUnderflowsUint32 is returned if when converting an int to a uint32 would underflow uint32
	ErrIntUnderflowsUint32 = errors.New("int underflows uint32")
	// ErrIntOverflowsUint32 is returned if when converting an int to a uint32 would overflow uint32
	ErrIntOverflowsUint32 = errors.New("int overflows uint32")
)

// MultUint64 multiplies a by b, returning an error if the values would overflow
func MultUint64(a, b uint64) (uint64, error) {
	c := a * b
	if a != 0 && c/a != b {
		return 0, ErrUint64MultOverflow
	}
	return c, nil
}

// AddUint64 adds a and b, returning an error if the values would overflow
func AddUint64(a, b uint64) (uint64, error) {
	c := a + b
	if c < a || c < b {
		return 0, ErrUint64AddOverflow
	}
	return c, nil
}

// AddUint32 adds a and b, returning an error if the values would overflow
func AddUint32(a, b uint32) (uint32, error) {
	c := a + b
	if c < a || c < b {
		return 0, ErrUint32AddOverflow
	}
	return c, nil
}

// Uint64ToInt64 converts a uint64 to an int64, returning an error if the uint64 value overflows int64
func Uint64ToInt64(a uint64) (int64, error) {
	b := int64(a)
	if b < 0 {
		return 0, ErrUint64OverflowsInt64
	}
	return b, nil
}

// Int64ToUint64 converts an int64 to a uint64, returning an error if the int64 value underflows uint64
func Int64ToUint64(a int64) (uint64, error) {
	if a < 0 {
		return 0, ErrInt64UnderflowsUint64
	}
	return uint64(a), nil
}

// IntToUint32 converts int to uint32, returning an error if the int value is negative or overflows uint32
func IntToUint32(a int) (uint32, error) {
	if a < 0 {
		return 0, ErrIntUnderflowsUint32
	}

	if uint64(a) > math.MaxUint32 {
		return 0, ErrIntOverflowsUint32
	}

	return uint32(a), nil
}
