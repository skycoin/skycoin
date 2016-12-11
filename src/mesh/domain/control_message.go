package domain

import (
	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
)

type SetControlChannelMessage struct {
	FromPeerID cipher.PubKey
}

type SetControlChannelResponseMessage struct {
	FromPeerID cipher.PubKey
	ChannelID  uuid.UUID
}

type SetRouteControlMessage struct {
	ChannelID uuid.UUID
	RequestID uuid.UUID

	FromPeerID      cipher.PubKey
	RouteID         RouteID
	ForwardRouteID  RouteID
	ForwardPeerID   cipher.PubKey
	BackwardRouteID RouteID
	BackwardPeerID  cipher.PubKey
}

type RemoveRouteControlMessage struct {
	FromPeerID cipher.PubKey
	ChannelID  uuid.UUID
	RequestID  uuid.UUID

	RouteID RouteID
}

type HealthCheckControlMessage struct {
	FromPeerID cipher.PubKey
	ChannelID  uuid.UUID
	RequestID  uuid.UUID

	RouteID RouteID
}

type ResponseMessage struct {
	FromPeerID cipher.PubKey
	RequestID  uuid.UUID

	Result bool
}
