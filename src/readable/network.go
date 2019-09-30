package readable

import (
	"github.com/SkycoinProject/skycoin/src/daemon"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/util/useragent"
)

// Connection a connection's state within the daemon
type Connection struct {
	GnetID               uint64                 `json:"id"`
	Addr                 string                 `json:"address"`
	LastSent             int64                  `json:"last_sent"`
	LastReceived         int64                  `json:"last_received"`
	ConnectedAt          int64                  `json:"connected_at"`
	Outgoing             bool                   `json:"outgoing"`
	State                daemon.ConnectionState `json:"state"`
	Mirror               uint32                 `json:"mirror"`
	ListenPort           uint16                 `json:"listen_port"`
	Height               uint64                 `json:"height"`
	UserAgent            useragent.Data         `json:"user_agent"`
	IsTrustedPeer        bool                   `json:"is_trusted_peer"`
	UnconfirmedVerifyTxn VerifyTxn              `json:"unconfirmed_verify_transaction"`
}

// NewConnection copies daemon.Connection to a struct with json tags
func NewConnection(c *daemon.Connection) Connection {
	var lastSent int64
	var lastReceived int64
	var connectedAt int64

	if !c.Gnet.LastSent.IsZero() {
		lastSent = c.Gnet.LastSent.Unix()
	}
	if !c.Gnet.LastReceived.IsZero() {
		lastReceived = c.Gnet.LastReceived.Unix()
	}
	if !c.ConnectedAt.IsZero() {
		connectedAt = c.ConnectedAt.Unix()
	}

	return Connection{
		GnetID:               c.Gnet.ID,
		Addr:                 c.Addr,
		LastSent:             lastSent,
		LastReceived:         lastReceived,
		ConnectedAt:          connectedAt,
		Outgoing:             c.Outgoing,
		State:                c.State,
		Mirror:               c.Mirror,
		ListenPort:           c.ListenPort,
		Height:               c.Height,
		UserAgent:            c.UserAgent,
		IsTrustedPeer:        c.Pex.Trusted,
		UnconfirmedVerifyTxn: NewVerifyTxn(c.UnconfirmedVerifyTxn),
	}
}

// VerifyTxn transaction verification parameters
type VerifyTxn struct {
	BurnFactor          uint32 `json:"burn_factor"`
	MaxTransactionSize  uint32 `json:"max_transaction_size"`
	MaxDropletPrecision uint8  `json:"max_decimals"`
}

// NewVerifyTxn converts params.VerifyTxn to VerifyTxn
func NewVerifyTxn(p params.VerifyTxn) VerifyTxn {
	return VerifyTxn{
		BurnFactor:          p.BurnFactor,
		MaxTransactionSize:  p.MaxTransactionSize,
		MaxDropletPrecision: p.MaxDropletPrecision,
	}
}
