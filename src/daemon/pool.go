package daemon

import (
	"time"

	//"github.com/skycoin/skycoin/src/daemon/gnet"
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
	// Maximum number of connections to maintain
	MaxConnections int
	// These should be assigned by the controlling daemon
	address string
	port    int
}

// NewPoolConfig creates pool config
func NewPoolConfig() PoolConfig {
	//defIdleLimit := time.Minute
	return PoolConfig{
		port:                6677,
		address:             "",
		DialTimeout:         time.Second * 30,
		MessageHandlingRate: time.Millisecond * 50,
		PingRate:            5 * time.Second,
		IdleLimit:           60 * time.Second,
		IdleCheckRate:       1 * time.Second,
		ClearStaleRate:      1 * time.Second,
		EventChannelSize:    4096,
		MaxConnections:      128,
	}
}

// Pool maintains config and pool
type Pool struct {
	Config PoolConfig
	Pool   *gnet.ConnectionPool
}

// NewPool creates pool
func NewPool(c PoolConfig, d *Daemon) *Pool {
	pool := &Pool{
		Config: c,
		Pool:   nil,
	}

	cfg := gnet.NewConfig()
	cfg.DialTimeout = pool.Config.DialTimeout
	cfg.Port = uint16(pool.Config.port)
	cfg.Address = pool.Config.address
	cfg.ConnectCallback = d.onGnetConnect
	cfg.DisconnectCallback = d.onGnetDisconnect
	cfg.MaxConnections = pool.Config.MaxConnections

	pool.Pool = gnet.NewConnectionPool(cfg, d)

	return pool
}

// Shutdown closes all connections and stops listening
func (pool *Pool) Shutdown() {
	if pool.Pool != nil {
		pool.Pool.Shutdown()
	}
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

// Send a ping if our last message sent was over pingRate ago
func (pool *Pool) sendPings() {
	pool.Pool.SendPings(pool.Config.PingRate, &PingMessage{})
}

// Removes connections that have not sent a message in too long
func (pool *Pool) clearStaleConnections() {
	pool.Pool.ClearStaleConnections(pool.Config.IdleLimit, ErrDisconnectIdle)
}
