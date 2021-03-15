package daemon

import (
	"time"

	"github.com/skycoin/skycoin/src/daemon/gnet"
)

// PoolConfig pool config
type PoolConfig struct {
	// Timeout when trying to connect to new peers through the pool
	DialTimeout time.Duration
	// How often to process message buffers and generate events
	MessageHandlingRate time.Duration
	// How long to wait before sending another ping
	PingRate time.Duration
	// How long a connection can idle before considered stale
	IdleLimit time.Duration
	// How often to check for needed pings
	IdleCheckRate time.Duration
	// How often to check for stale connections
	ClearStaleRate time.Duration
	// Buffer size for gnet.ConnectionPool's network Read events
	EventChannelSize int
	// Maximum number of connections
	MaxConnections int
	// Maximum number of outgoing connections
	MaxOutgoingConnections int
	// Maximum number of incoming connections
	MaxIncomingConnections int
	// Maximum number of outgoing connections to peers in the DefaultConnections list to maintain
	MaxDefaultPeerOutgoingConnections int
	// Default "trusted" peers
	DefaultConnections []string
	// Maximum length of incoming messages in bytes
	MaxIncomingMessageLength int
	// Maximum length of outgoing messages in bytes
	MaxOutgoingMessageLength int
	// These should be assigned by the controlling daemon
	address string
	port    int
}

// NewPoolConfig creates pool config
func NewPoolConfig() PoolConfig {
	return PoolConfig{
		port:                              6677,
		address:                           "",
		DialTimeout:                       time.Second * 30,
		MessageHandlingRate:               time.Millisecond * 50,
		PingRate:                          5 * time.Second,
		IdleLimit:                         60 * time.Second,
		IdleCheckRate:                     1 * time.Second,
		ClearStaleRate:                    1 * time.Second,
		EventChannelSize:                  4096,
		MaxConnections:                    128,
		MaxOutgoingConnections:            8,
		MaxIncomingConnections:            120,
		MaxDefaultPeerOutgoingConnections: 2,
		MaxOutgoingMessageLength:          256 * 1024,
		MaxIncomingMessageLength:          1024 * 1024,
	}
}

// Pool maintains config and pool
type Pool struct {
	Config PoolConfig
	Pool   *gnet.ConnectionPool
}

// NewPool creates pool
func NewPool(cfg PoolConfig, d *Daemon) (*Pool, error) {
	gnetCfg := gnet.NewConfig()
	gnetCfg.DialTimeout = cfg.DialTimeout
	gnetCfg.Port = uint16(cfg.port)
	gnetCfg.Address = cfg.address
	gnetCfg.ConnectCallback = d.onGnetConnect
	gnetCfg.DisconnectCallback = d.onGnetDisconnect
	gnetCfg.ConnectFailureCallback = d.onGnetConnectFailure
	gnetCfg.MaxConnections = cfg.MaxConnections
	gnetCfg.MaxOutgoingConnections = cfg.MaxOutgoingConnections
	gnetCfg.MaxIncomingConnections = cfg.MaxIncomingConnections
	gnetCfg.MaxDefaultPeerOutgoingConnections = cfg.MaxDefaultPeerOutgoingConnections
	gnetCfg.DefaultConnections = cfg.DefaultConnections
	gnetCfg.MaxIncomingMessageLength = cfg.MaxIncomingMessageLength
	gnetCfg.MaxOutgoingMessageLength = cfg.MaxOutgoingMessageLength

	pool, err := gnet.NewConnectionPool(gnetCfg, d)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Config: cfg,
		Pool:   pool,
	}, nil
}

// Shutdown closes all connections and stops listening
func (pool *Pool) Shutdown() {
	if pool == nil {
		return
	}
	pool.Pool.Shutdown()
}

// Run starts listening on the configured Port
func (pool *Pool) Run() error {
	logger.Infof("daemon.Pool listening on port %d", pool.Config.port)
	return pool.Pool.Run()
}

// RunOffline runs the pool without a listener. This is necessary to process strand requests.
func (pool *Pool) RunOffline() error {
	return pool.Pool.RunOffline()
}

// sendPings send a ping if our last message sent was over pingRate ago
func (pool *Pool) sendPings() {
	if err := pool.Pool.SendPings(pool.Config.PingRate, &PingMessage{}); err != nil {
		logger.WithError(err).Error("sendPings failed")
	}
}

// getStaleConnections returns connections that have been idle for longer than idleLimit
func (pool *Pool) getStaleConnections() ([]string, error) {
	return pool.Pool.GetStaleConnections(pool.Config.IdleLimit)
}

// IsMaxOutgoingDefaultConnectionsReached returns whether max outgoing default connections reached
func (pool *Pool) IsMaxOutgoingDefaultConnectionsReached() bool {
	return pool.Pool.IsMaxOutgoingDefaultConnectionsReached()
}
