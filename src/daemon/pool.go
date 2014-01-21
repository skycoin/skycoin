package daemon

import (
    "github.com/skycoin/gnet"
    "log"
    "time"
)

var (
    // Connection pool
    Pool *gnet.ConnectionPool = nil
    // Timeout when trying to connect to new peers through the pool
    poolDialTimeout = time.Second * 30
)

// Begins listening on port for connections and periodically scanning for
// messages on read_interval
func InitPool(port int) {
    logger.Info("InitPool on port %d", port)
    if Pool != nil {
        log.Panic("ConnectionPool is already initialised")
    }
    gnet.DialTimeout = poolDialTimeout
    Pool = gnet.NewConnectionPool(port)
    Pool.DisconnectCallback = onDisconnect
    Pool.ConnectCallback = onGnetConnect
    go func() {
        err := Pool.StartListen()
        if err != nil {
            log.Panic(err)
        }
    }()
}

// Closes all connections and stops listening
func ShutdownPool() {
    if Pool != nil {
        Pool.StopListen()
    }
    Pool = nil
    logger.Info("Shutdown pool")
}

// Triggered when an gnet.Connection is connected
func onGnetConnect(c *gnet.Connection, solicited bool) {
    onConnectEvent <- ConnectEvent{Addr: c.Addr(), Solicited: solicited}
}
