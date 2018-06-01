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

var errorCodeByError map[error]uint16

var initErrorCodeMap = func() {
	errorCodeByError = make(map[error]uint16)
	for i, err := range errorByCode {
		errorCodeByError[err] = uint16(i)
	}
}

// Unexpected error condition detected
const ErrorCodeUnknown = 0xFFFF

// Success error code
const Success = 0

// Retrieve error object by corresponding error code
func GetError(code uint16) error {
	return errorByCode[code]
}

// Retrieve error code representing corresponding error object
func GetErrorCode(err error) uint16 {
	if initErrorCodeMap != nil {
		initErrorCodeMap()
		initErrorCodeMap = nil
	}
	return errorCodeByError[err]
}
