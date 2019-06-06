/*
Package timeutil provides time related utility methods
*/
package timeutil

import "time"

// NanoToTime converts nanoseconds to time.Time
func NanoToTime(n int64) time.Time {
	zeroTime := time.Time{}
	if n == zeroTime.UnixNano() {
		return zeroTime
	}
	return time.Unix(n/int64(time.Second), n%int64(time.Second))
}
