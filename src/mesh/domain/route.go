package domain

import (
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

type Route struct {
	// Forward should never be cipher.PubKey{}
	ForwardToPeer             cipher.PubKey
	ForwardRewriteSendRouteID RouteID

	BackwardToPeer             cipher.PubKey
	BackwardRewriteSendRouteID RouteID

	// time.Unix(0,0) means it lives forever
	ExpiryTime time.Time
}

type LocalRoute struct {
	LastForwardingPeer cipher.PubKey
	TerminatingPeer    cipher.PubKey
	LastHopRouteID     RouteID
	LastConfirmed      time.Time
}

type ReplyTo struct {
	RouteID  RouteID
	FromPeer cipher.PubKey
}
