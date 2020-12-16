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
	// ErrDisconnectUnexpectedError this is returned when a seemingly impossible error is encountered, e.g. net.Conn.Addr() returns an invalid ip:port
	ErrDisconnectUnexpectedError gnet.DisconnectReason = errors.New("Unexpected error")
	// ErrDisconnectMaxOutgoingConnectionsReached is returned when connection pool size is greater than the maximum allowed
	ErrDisconnectMaxOutgoingConnectionsReached gnet.DisconnectReason = errors.New("Maximum outgoing connections was reached")
	// ErrDisconnectBlockchainPubkeyNotMatched is returned when the blockchain pubkey in introduction does not match
	ErrDisconnectBlockchainPubkeyNotMatched gnet.DisconnectReason = errors.New("Blockchain pubkey does not match")
	// ErrDisconnectBlockchainPubkeyNotProvided is returned when the blockchain pubkey in introduction is not provided
	ErrDisconnectBlockchainPubkeyNotProvided gnet.DisconnectReason = errors.New("Blockchain pubkey is not provided")
	// ErrDisconnectInvalidExtraData is returned when extra field can't be parsed
	ErrDisconnectInvalidExtraData gnet.DisconnectReason = errors.New("Invalid extra data in message")
	// ErrDisconnectReceivedDisconnect received a DisconnectMessage
	ErrDisconnectReceivedDisconnect gnet.DisconnectReason = errors.New("Received DisconnectMessage")
	// ErrDisconnectInvalidUserAgent is returned if the peer provides an invalid user agent
	ErrDisconnectInvalidUserAgent gnet.DisconnectReason = errors.New("Invalid user agent")
	// ErrDisconnectRequestedByOperator the operator of the node requested a disconnect
	ErrDisconnectRequestedByOperator gnet.DisconnectReason = errors.New("Disconnect requested by the node operator")
	// ErrDisconnectPeerlistFull the peerlist is full
	ErrDisconnectPeerlistFull gnet.DisconnectReason = errors.New("Peerlist is full")
	// ErrDisconnectInvalidBurnFactor invalid burn factor in introduction message
	ErrDisconnectInvalidBurnFactor gnet.DisconnectReason = errors.New("Invalid burn factor in introduction message")
	// ErrDisconnectInvalidMaxTransactionSize invalid max transaction size in introduction message
	ErrDisconnectInvalidMaxTransactionSize gnet.DisconnectReason = errors.New("Invalid max transaction size in introduction message")
	// ErrDisconnectInvalidMaxDropletPrecision invalid max droplet precision in introduction message
	ErrDisconnectInvalidMaxDropletPrecision gnet.DisconnectReason = errors.New("Invalid max droplet precision in introduction message")

	// ErrDisconnectUnknownReason used when mapping an unknown reason code to an error. Is not sent over the network.
	ErrDisconnectUnknownReason gnet.DisconnectReason = errors.New("Unknown DisconnectReason")

	disconnectReasonCodes = map[gnet.DisconnectReason]uint16{
		ErrDisconnectUnknownReason: 0,

		ErrDisconnectVersionNotSupported:           1,
		ErrDisconnectIntroductionTimeout:           2,
		ErrDisconnectIsBlacklisted:                 3,
		ErrDisconnectSelf:                          4,
		ErrDisconnectConnectedTwice:                5,
		ErrDisconnectIdle:                          6,
		ErrDisconnectNoIntroduction:                7,
		ErrDisconnectIPLimitReached:                8,
		ErrDisconnectUnexpectedError:               9,
		ErrDisconnectMaxOutgoingConnectionsReached: 10,
		ErrDisconnectBlockchainPubkeyNotMatched:    11,
		ErrDisconnectInvalidExtraData:              12,
		ErrDisconnectReceivedDisconnect:            13,
		ErrDisconnectInvalidUserAgent:              14,
		ErrDisconnectRequestedByOperator:           15,
		ErrDisconnectPeerlistFull:                  16,
		ErrDisconnectInvalidBurnFactor:             17,
		ErrDisconnectInvalidMaxTransactionSize:     18,
		ErrDisconnectInvalidMaxDropletPrecision:    19,

		// gnet codes are registered here, but they are not sent in a DISC
		// message by gnet. Only daemon sends a DISC packet.
		// If gnet chooses to disconnect it will not send a DISC packet.
		gnet.ErrDisconnectSetReadDeadlineFailed:  1001,
		gnet.ErrDisconnectInvalidMessageLength:   1002,
		gnet.ErrDisconnectMalformedMessage:       1003,
		gnet.ErrDisconnectUnknownMessage:         1004,
		gnet.ErrDisconnectShutdown:               1005,
		gnet.ErrDisconnectMessageDecodeUnderflow: 1006,
		gnet.ErrDisconnectTruncatedMessageID:     1007,
	}

	disconnectCodeReasons map[uint16]gnet.DisconnectReason
)

func init() {
	disconnectCodeReasons = make(map[uint16]gnet.DisconnectReason, len(disconnectReasonCodes))

	for r, c := range disconnectReasonCodes {
		disconnectCodeReasons[c] = r
	}
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
