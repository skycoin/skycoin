package daemon

import (
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
)

var errorByCode = [...]error{
	nil,
	ErrDisconnectInvalidVersion,
	ErrDisconnectIntroductionTimeout,
	ErrDisconnectVersionSendFailed,
	ErrDisconnectIsBlacklisted,
	ErrDisconnectSelf,
	ErrDisconnectConnectedTwice,
	ErrDisconnectIdle,
	ErrDisconnectNoIntroduction,
	ErrDisconnectIPLimitReached,
	ErrDisconnectOtherError,
	gnet.ErrDisconnectReadFailed,
	gnet.ErrDisconnectWriteFailed,
	gnet.ErrDisconnectSetReadDeadlineFailed,
	gnet.ErrDisconnectInvalidMessageLength,
	gnet.ErrDisconnectMalformedMessage,
	gnet.ErrDisconnectUnknownMessage,
	nil, //	gnet.ErrDisconnectWriteQueueFull,
	gnet.ErrDisconnectUnexpectedError,
	gnet.ErrConnectionPoolClosed,
	pex.ErrPeerlistFull,
	pex.ErrInvalidAddress,
	pex.ErrNoLocalhost,
	pex.ErrNotExternalIP,
	pex.ErrPortTooLow,
	pex.ErrBlacklistedAddress}

var errorCodeByError map[error]uint32

var initErrorCodeMap = func() {
	errorCodeByError = make(map[error]uint32)
	for i, err := range errorByCode {
		errorCodeByError[err] = uint32(i)
	}
}

// ErrorCodeUnknown is used on unexpected error condition detected
const ErrorCodeUnknown = 0xFFFFFFFF

// Success error code
const Success = 0

// GetError Retrieve error object by corresponding error code
func GetError(code uint32) error {
	if code < uint32(len(errorByCode)) {
		return errorByCode[code]
	}
	return nil
}

// GetErrorCode Retrieve error code representing corresponding error object
func GetErrorCode(err error) uint32 {
	if code, exists := errorCodeByError[err]; exists {
		return code
	}
	return ErrorCodeUnknown
}

func init() {
	initErrorCodeMap()
	initErrorCodeMap = nil
}
