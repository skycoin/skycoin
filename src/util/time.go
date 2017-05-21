package util

import (
	"time"
)

// Now returns the current UTC time
func Now() time.Time {
	return time.Now().UTC()
}

// UnixNow returns the current UTC time as unix timestamp
func UnixNow() int64 {
	return Now().Unix()
}

// ZeroTime returns the zero value time
func ZeroTime() time.Time {
	return time.Time{}
}
