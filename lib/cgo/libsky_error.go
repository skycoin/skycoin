package main

import "C"

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

func catchApiPanic() (retVal uint32) {
	retVal = SKY_OK
	if err := recover(); err != nil {
		// TODO: Fix to be like retVal = libErrorCode(err)
		retVal = SKY_ERROR
	}
	return
}
