package daemon

import (
	"errors"

	"github.com/skycoin/skycoin/src/daemon/gnet"
)

var (
	// ErrDisconnectVersionNotSupported version is below minimum supported version
	ErrDisconnectVersionNotSupported gnet.DisconnectReason = errors.New("Version is below minimum supported version")
	// ErrDisconnectIntroductionTimeout timeout
	ErrDisconnectIntroductionTimeout gnet.DisconnectReason = errors.New("Introduction timeout")
	// ErrDisconnectIsBlacklisted is blacklisted
	ErrDisconnectIsBlacklisted gnet.DisconnectReason = errors.New("Blacklisted")
	// ErrDisconnectSelf self connnect
	ErrDisconnectSelf gnet.DisconnectReason = errors.New("Self connect")
	// ErrDisconnectConnectedTwice connect twice
	ErrDisconnectConnectedTwice gnet.DisconnectReason = errors.New("Already connected")
	// ErrDisconnectIdle idle
	ErrDisconnectIdle gnet.DisconnectReason = errors.New("Idle")
	// ErrDisconnectNoIntroduction no introduction
	ErrDisconnectNoIntroduction gnet.DisconnectReason = errors.New("First message was not an Introduction")
	// ErrDisconnectIPLimitReached ip limit reached
	ErrDisconnectIPLimitReached gnet.DisconnectReason = errors.New("Maximum number of connections for this IP was reached")
	// ErrDisconnectIncomprehensibleError this is returned when a seemingly impossible error is encountered
	// e.g. net.Conn.Addr() returns an invalid ip:port
	ErrDisconnectIncomprehensibleError gnet.DisconnectReason = errors.New("Incomprehensible error")
	// ErrDisconnectMaxOutgoingConnectionsReached is returned when connection pool size is greater than the maximum allowed
	ErrDisconnectMaxOutgoingConnectionsReached gnet.DisconnectReason = errors.New("Maximum outgoing connections was reached")
	// ErrDisconnectBlockchainPubkeyNotMatched is returned when the blockchain pubkey in introduction does not match
	ErrDisconnectBlockchainPubkeyNotMatched gnet.DisconnectReason = errors.New("Blockchain pubkey does not match")
	// ErrDisconnectInvalidExtraData is returned when extra field can't be parsed as specific data type.
	// e.g. ExtraData length in IntroductionMessage is not the same as cipher.PubKey
	ErrDisconnectInvalidExtraData gnet.DisconnectReason = errors.New("Invalid extra data in message")
	// ErrDisconnectPeerlistFull no space in peers pool
	ErrDisconnectPeerlistFull gnet.DisconnectReason = errors.New("Peerlist is full")
	// ErrDisconnectReceivedDisconnect received a DisconnectMessage
	ErrDisconnectReceivedDisconnect gnet.DisconnectReason = errors.New("Received DisconnectMessage")
	// ErrDisconnectUnknownReason used when mapping an unknown reason code to an error. Is not sent over the network.
	ErrDisconnectUnknownReason gnet.DisconnectReason = errors.New("Unknown DisconnectReason")

	disconnectReasonCodes = map[gnet.DisconnectReason]uint16{
		ErrDisconnectUnknownReason: 0,

		ErrDisconnectVersionNotSupported:           1,
		ErrDisconnectIntroductionTimeout:           2,
		ErrDisconnectIsBlacklisted:                 4,
		ErrDisconnectSelf:                          5,
		ErrDisconnectConnectedTwice:                6,
		ErrDisconnectIdle:                          7,
		ErrDisconnectNoIntroduction:                8,
		ErrDisconnectIPLimitReached:                9,
		ErrDisconnectIncomprehensibleError:         10,
		ErrDisconnectMaxOutgoingConnectionsReached: 11,
		ErrDisconnectBlockchainPubkeyNotMatched:    12,
		ErrDisconnectInvalidExtraData:              13,
		ErrDisconnectPeerlistFull:                  14,
		ErrDisconnectReceivedDisconnect:            16,

		gnet.ErrDisconnectReadFailed:            1001,
		gnet.ErrDisconnectWriteFailed:           1002,
		gnet.ErrDisconnectSetReadDeadlineFailed: 1003,
		gnet.ErrDisconnectInvalidMessageLength:  1004,
		gnet.ErrDisconnectMalformedMessage:      1005,
		gnet.ErrDisconnectUnknownMessage:        1006,
		gnet.ErrDisconnectUnexpectedError:       1007,
	}

	disconnectCodeReasons map[uint16]gnet.DisconnectReason
)

func init() {
	disconnectCodeReasons = make(map[uint16]gnet.DisconnectReason, len(disconnectReasonCodes))
}

// DisconnectReasonToCode maps a gnet.DisconnectReason to a 16-byte code
func DisconnectReasonToCode(r gnet.DisconnectReason) uint16 {
	return disconnectReasonCodes[r]
}

// DisconnectCodeToReason maps a disconnect code to a gnet.DisconnectReason
func DisconnectCodeToReason(c uint16) gnet.DisconnectReason {
	r, ok := disconnectCodeReasons[c]
	if !ok {
		return ErrDisconnectUnknownReason
	}
	return r
}
