package main

const (
	SKY_OK    = 0
	SKY_ERROR = 0xFFFFFFFF
)

func libErrorCode(err error) uint32 {
	if err == nil {
		return SKY_OK
	}
	// TODO: Implement error codes
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
