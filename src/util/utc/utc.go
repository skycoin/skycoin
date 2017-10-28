// Package utc provides helper function to get UTC time
package utc

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
