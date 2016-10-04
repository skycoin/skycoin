package domain

import (
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

type Route struct {
	// Forward should never be cipher.PubKey{}
	ForwardToPeer        cipher.PubKey
	ForwardRewriteSendId RouteId

	BackwardToPeer        cipher.PubKey
	BackwardRewriteSendId RouteId

	// time.Unix(0,0) means it lives forever
	ExpiryTime time.Time
}

type LocalRoute struct {
	LastForwardingPeer cipher.PubKey
	TerminatingPeer    cipher.PubKey
	LastHopId          RouteId
	LastConfirmed      time.Time
}

type ReplyTo struct {
	RouteId  RouteId
	FromPeer cipher.PubKey
}
