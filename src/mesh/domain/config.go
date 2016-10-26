package domain

import (
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

type NodeConfig struct {
	PubKey                        cipher.PubKey
	MaximumForwardingDuration     time.Duration
	RefreshRouteDuration          time.Duration
	ExpireRoutesInterval          time.Duration
	TransportMessageChannelLength int
	//ChaCha20Key                   [32]byte
}
