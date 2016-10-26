package domain

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
)

type RouteID uuid.UUID
type MessageID uuid.UUID

type MeshMessage struct {
	ReplyTo  ReplyTo
	Contents []byte
}

type ReplyTo struct {
	RouteID    RouteID
	FromPeerID cipher.PubKey
}

// Fields must be public (capital first letter) for encoder
type MessageBase struct {
	// If RouteId is unknown, but not cipher.PubKey{}, then the message should be received here
	//  the RouteId can be used to reply back thru the route
	SendRouteID RouteID
	SendBack    bool
	// For sending the reply from the last node in a route
	FromPeerID cipher.PubKey
	Nonce      [4]byte
}

type UserMessage struct {
	MessageBase
	MessageID MessageID
	Index     uint64
	Count     uint64
	Contents  []byte
}

type SetRouteMessage struct {
	MessageBase
	SetRouteID                 RouteID
	ConfirmRouteID             RouteID
	ForwardToPeerID            cipher.PubKey
	ForwardRewriteSendRouteID  RouteID
	BackwardToPeerID           cipher.PubKey
	BackwardRewriteSendRouteID RouteID
	DurationHint               time.Duration
}

// This allows ExtendRoute() to block so that messages aren't lost while a route is
//  not yet established
type SetRouteReply struct {
	MessageBase
	ConfirmRouteID RouteID
}

// Refreshes the route as it passes thru it
type RefreshRouteMessage struct {
	MessageBase
	DurationHint    time.Duration
	ConfirmRoutedID RouteID
}

// Deletes the route as it passes thru it
type DeleteRouteMessage struct {
	MessageBase
}

// Add a new node to the network
type AddNodeMessage struct {
	MessageBase
	Content []byte
}

type MessageUnderAssembly struct {
	Fragments   map[uint64]UserMessage
	SendRouteID RouteID
	SendBack    bool
	Count       uint64
	Dropped     bool
	ExpiryTime  time.Time
}
