package domain

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher"
)

type RouteConfig struct {
	Id    uuid.UUID
	Peers []cipher.PubKey
}

type NodeConfig struct {
	PubKey cipher.PubKey
	//ChaCha20Key                   [32]byte
	MaximumForwardingDuration     time.Duration
	RefreshRouteDuration          time.Duration
	ExpireMessagesInterval        time.Duration
	ExpireRoutesInterval          time.Duration
	TimeToAssembleMessage         time.Duration
	TransportMessageChannelLength int
}
