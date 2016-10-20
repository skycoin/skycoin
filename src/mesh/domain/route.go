package domain

import (
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

type Route struct {
	// Forward should never be cipher.PubKey{}
	ForwardToPeerID           cipher.PubKey
	ForwardRewriteSendRouteID RouteID

	BackwardToPeerID           cipher.PubKey
	BackwardRewriteSendRouteID RouteID

	// time.Unix(0,0) means it lives forever
	ExpiryTime time.Time
}

type LocalRoute struct {
	LastForwardingPeerID cipher.PubKey
	TerminatingPeerID    cipher.PubKey
	LastHopRouteID       RouteID
	LastConfirmed        time.Time
}
