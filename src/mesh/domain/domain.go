package domain

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
)

type RouteId uuid.UUID
type MessageId uuid.UUID

// Fields must be public (capital first letter) for encoder
type MessageBase struct {
	// If RouteId is unknown, but not cipher.PubKey{}, then the message should be received here
	//  the RouteId can be used to reply back thru the route
	SendId   RouteId
	SendBack bool
	// For sending the reply from the last node in a route
	FromPeer cipher.PubKey
	Reliably bool
	Nonce    [4]byte
}

type UserMessage struct {
	MessageBase
	MessageId MessageId
	Index     uint64
	Count     uint64
	Contents  []byte
}

type SetRouteMessage struct {
	MessageBase
	SetRouteId            RouteId
	ConfirmId             RouteId
	ForwardToPeer         cipher.PubKey
	ForwardRewriteSendId  RouteId
	BackwardToPeer        cipher.PubKey
	BackwardRewriteSendId RouteId
	DurationHint          time.Duration
}

// This allows ExtendRoute() to block so that messages aren't lost while a route is
//  not yet established
type SetRouteReply struct {
	MessageBase
	ConfirmId RouteId
}

// Refreshes the route as it passes thru it
type RefreshRouteMessage struct {
	MessageBase
	DurationHint time.Duration
	ConfirmId    RouteId
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

// Get a node from the network
//type GetNodeMessage struct {
//	MessageBase
//}

// Set up a node from the network
//type SetUpNodeMessage struct {
//	MessageBase
//}

// Delete a node from the network
//type DeleteNodeMessage struct {
//	MessageBase
//}

// Get a route between two nodes from the node manager
//type GetNodeRouteMessage struct {
//	MessageBase
//}

type RouteConfig struct {
	Id    uuid.UUID
	Peers []cipher.PubKey
}

type MessageToSend struct {
	ThruRoute uuid.UUID
	Contents  []byte
	Reliably  bool
}

type MessageToReceive struct {
	Contents      []byte
	Reply         []byte
	ReplyReliably bool
}

type ToConnect struct {
	Peer cipher.PubKey
	Info string
}
