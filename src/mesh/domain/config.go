package domain

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
)

type RouteConfig struct {
	ID    uuid.UUID
	Peers []cipher.PubKey
}

type NodeConfig struct {
	PubKey                        cipher.PubKey
	MaximumForwardingDuration     time.Duration
	RefreshRouteDuration          time.Duration
	ExpireMessagesInterval        time.Duration
	ExpireRoutesInterval          time.Duration
	TimeToAssembleMessage         time.Duration
	TransportMessageChannelLength int
	//ChaCha20Key                   [32]byte
}

type TransportConfig struct {
	SendChannelLength uint32
}
