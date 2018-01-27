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
	gnet.ErrDisconnectWriteQueueFull,
	gnet.ErrDisconnectUnexpectedError,
	gnet.ErrConnectionPoolClosed,
	pex.ErrPeerlistFull,
	pex.ErrInvalidAddress,
	pex.ErrNoLocalhost,
	pex.ErrNotExternalIP,
	pex.ErrPortTooLow,
	pex.ErrBlacklistedAddress}

var errorCodeByError map[error]uint8

var initErrorCodeMap = func() {
	errorCodeByError = make(map[error]uint8)
	for i, err := range errorByCode {
		errorCodeByError[err] = uint8(i)
	}
}

const ErrorCodeNone = 0xFF

func GetError(code uint8) error {
	return errorByCode[code]
}

func GetErrorCode(err error) uint8 {
	if initErrorCodeMap != nil {
		initErrorCodeMap()
		initErrorCodeMap = nil
	}
	return errorCodeByError[err]
}
