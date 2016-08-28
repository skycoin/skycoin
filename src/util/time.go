package util

import (
	"time"
)

// Returns the current UTC time
func Now() time.Time {
	return time.Now().UTC()
}

// Returns the current UTC time as unix timestamp
func UnixNow() int64 {
	return Now().Unix()
}

// Returns the zero value time
func ZeroTime() time.Time {
	return time.Time{}
}
