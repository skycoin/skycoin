
package main

import "C"

const (
	NO_ERROR = 1
	ERR_UNKNOWN = 0
)

func libErrorCode(err error) uint32 {
	if err == nil {
		return NO_ERROR
	}
	// TODO: Implement error codes
	return ERR_UNKNOWN
}

