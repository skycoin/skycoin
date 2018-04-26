package coin

import "errors"

var (
	// ErrUint64MultOverflow is returned when multiplying uint64 values would overflow uint64
	ErrUint64MultOverflow = errors.New("uint64 multiplication overflow")
	// ErrUint64AddOverflow is returned when adding uint64 values would overflow uint64
	ErrUint64AddOverflow = errors.New("uint64 addition overflow")
	// ErrUint32AddOverflow is returned when adding uint32 values would overflow uint32
	ErrUint32AddOverflow = errors.New("uint32 addition overflow")
	// ErrUint64OverflowsInt64 is returned when converting a uint64 to an int64 would overflow int64
	ErrUint64OverflowsInt64 = errors.New("uint64 overflows int64")
	// ErrInt64UnderflowsUint64 is returned when converting an int64 to a uint64 would underflow uint64
	ErrInt64UnderflowsUint64 = errors.New("int64 underflows uint64")
)

func multUint64(a, b uint64) (uint64, error) {
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

func addUint32(a, b uint32) (uint32, error) {
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
