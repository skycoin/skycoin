package domain

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
)

type RouteID uuid.UUID
type MessageID uuid.UUID

// Fields must be public (capital first letter) for encoder
type MessageBase struct {
	// If RouteId is unknown, but not cipher.PubKey{}, then the message should be received here
	//  the RouteId can be used to reply back thru the route
	SendRouteID RouteID
	SendBack    bool
	// For sending the reply from the last node in a route
	FromPeer cipher.PubKey
	Nonce    [4]byte
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
	ForwardToPeer              cipher.PubKey
	ForwardRewriteSendRouteID  RouteID
	BackwardToPeer             cipher.PubKey
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

type MessageToSend struct {
	ThruRoute uuid.UUID
	Contents  []byte
}

type MessageToReceive struct {
	Contents []byte
	Reply    []byte
}

type MessageUnderAssembly struct {
	Fragments   map[uint64]UserMessage
	SendRouteID RouteID
	SendBack    bool
	Count       uint64
	Dropped     bool
	ExpiryTime  time.Time
}

type MeshMessage struct {
	ReplyTo  ReplyTo
	Contents []byte
}
