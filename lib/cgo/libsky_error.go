package main

import (
	"errors"
)

const (
	// SKY_ERROR generic error condition
	SKY_ERROR = 0xFFFFFFFF
	// SKY_BAD_HANDLE invalid handle argument
	SKY_BAD_HANDLE = 0xFF000001
	// SKY_OK error code is used to report success
	SKY_OK = 0
)

var (
	ErrorBadHandle = errors.New("Invalid or unknown handle value")
	ErrorUnknown   = errors.New("Unexpected error")

	errorToCodeMap = map[error]uint32{
		ErrorBadHandle: SKY_BAD_HANDLE,
		ErrorUnknown:   SKY_ERROR,
	}
)

func libErrorCode(err error) uint32 {
	if err == nil {
		return SKY_OK
	}
	if errcode, isKnownError := errorToCodeMap[err]; isKnownError {
		return errcode
	}
	return SKY_ERROR
}

// Catch panic signals emitted by internal implementation
// of API methods. This function is mainly used in defer statements
// exceuted immediately before returning from API calls.
//
// @param errcode error status in function body
// @param err			`recover()` result
//
func catchApiPanic(errcode uint32, err interface{}) uint32 {
	if errcode != SKY_OK {
		// Error already detected in function body
		// Return right away
		return errcode
	}
	if err != nil {
		// TODO: Fix to be like retVal = libErrorCode(err)
		return SKY_ERROR
	}
	return SKY_OK
}
