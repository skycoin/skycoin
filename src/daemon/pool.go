package daemon

import (
    "github.com/skycoin/gnet"
    "log"
    "time"
)

type PoolConfig struct {
    Port int
    // Timeout when trying to connect to new peers through the pool
    DialTimeout time.Duration
    // How often to process message buffers and generate events
    MessageHandlingRate time.Duration
    // How long to wait before sending another ping
    PingRate time.Duration
    // How long a connection can idle before considered stale
    IdleLimit time.Duration
}

func NewPoolConfig() *PoolConfig {
    defIdleLimit := time.Minute * 90
    return &PoolConfig{
        Port:                6677,
        DialTimeout:         time.Second * 30,
        MessageHandlingRate: time.Millisecond * 30,
        PingRate:            defIdleLimit / 3,
        IdleLimit:           defIdleLimit,
    }
}

type Pool struct {
    Config *PoolConfig
    Pool   *gnet.ConnectionPool
}

func NewPool(c *PoolConfig) *Pool {
    return &Pool{
        Config: c,
        Pool:   nil,
    }
}

// Begins listening on port for connections and periodically scanning for
// messages on read_interval
func (self *Pool) Init(d *Daemon) {
    logger.Info("InitPool on port %d", self.Config.Port)
    gnet.DialTimeout = self.Config.DialTimeout
    pool := gnet.NewConnectionPool(self.Config.Port, d)
    pool.ConnectCallback = d.onGnetConnect
    pool.DisconnectCallback = d.onGnetDisconnect
    self.Pool = pool
}

// Closes all connections and stops listening
func (self *Pool) Shutdown() {
    if self.Pool != nil {
        self.Pool.StopListen()
        logger.Info("Shutdown pool")
    }
}

func (self *Pool) Start() {
    err := self.Pool.StartListen()
    if err != nil {
        log.Panic(err)
    }
}

// Send a ping if our last message sent was over pingRate ago
func (self *Pool) sendPings() {
    now := time.Now().UTC()
    for _, c := range self.Pool.Pool {
        if c.LastSent.Add(self.Config.PingRate).Before(now) {
            err := self.Pool.Dispatcher.SendMessage(c, &PingMessage{})
            if err != nil {
                logger.Warning("Failed to send ping message to %s", c.Addr())
            }
        }
    }
}

// Removes connections that have not sent a message in too long
func (self *Pool) clearStaleConnections() {
    now := time.Now().UTC()
    for _, c := range self.Pool.Pool {
        if c.LastReceived.Add(self.Config.IdleLimit).Before(now) {
            self.Pool.Disconnect(c, DisconnectIdle)
        }
    }
}

// Requests peers from our connections
// TODO -- batching all peer requests at once may cause performance issues
func (self *Pool) requestPeers() {
    for _, c := range self.Pool.Pool {
        m := NewGetPeersMessage()
        err := self.Pool.Dispatcher.SendMessage(c, m)
        if err != nil {
            logger.Warning("Failed to request peers from %s", c.Addr())
        }
    }
}
