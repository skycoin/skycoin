package daemon

import (
    "github.com/skycoin/gnet"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestInitPool(t *testing.T) {
    Pool = gnet.NewConnectionPool(port)
    assert.Panics(t, func() { InitPool(port) })
    Pool = nil

    assert.NotPanics(t, func() { InitPool(port) })
    assert.Equal(t, gnet.DialTimeout, poolDialTimeout)
    assert.Equal(t, Pool.DisconnectCallback, onDisconnect)
    assert.Equal(t, Pool.ConnectCallback, onGnetConnect)
    wait()
    // A second call to start listen will panic due to the pool already
    // listening
    assert.Panics(t, func() { assert.Nil(t, Pool.StartListen()) })
    Pool.StopListen()
    wait()
    ShutdownPool()
    wait()
    assert.Nil(t, Pool)
}

func SetupPoolShutdown(t *testing.T) *gnet.ConnectionPool {
    Pool = gnet.NewConnectionPool(port)
    pool := Pool
    assert.Equal(t, len(pool.DisconnectQueue), 0)
    pool.DisconnectQueue <- gnet.DisconnectEvent{
        ConnId: 1,
        Reason: DisconnectOtherError,
    }
    return pool
}

func ConfirmPoolShutdown(t *testing.T, pool *gnet.ConnectionPool) {
    wait()
    assert.Nil(t, Pool)
    // ShutdownPool() should call Pool.StopListen() which closes and resets
    // the DisconnectQueue channel
    assert.Equal(t, len(pool.DisconnectQueue), 0)
}

func TestShutdownPool(t *testing.T) {
    // Shutdown with nil pool is safe
    Pool = nil
    assert.NotPanics(t, ShutdownPool)
    assert.Nil(t, Pool)
    // Valid shutdown
    pool := SetupPoolShutdown(t)
    assert.NotPanics(t, ShutdownPool)
    ConfirmPoolShutdown(t, pool)
}

func testOnGnetConnectSolicitation(t *testing.T, c *gnet.Connection,
    addr string, sol bool) {
    onGnetConnect(c, sol)
    assert.Equal(t, len(onConnectEvent), 1)
    if len(onConnectEvent) == 0 {
        t.Fatalf("onConnectEvent is not empty, would block")
    }
    ce := <-onConnectEvent
    assert.Equal(t, ce.Addr, addr)
    assert.Equal(t, ce.Solicited, sol)
}

func TestOnGnetConnect(t *testing.T) {
    addr := "11.22.33.44:5555"
    c := gnetConnection(addr)
    assert.Equal(t, len(onConnectEvent), 0)
    testOnGnetConnectSolicitation(t, c, addr, false)
    assert.Equal(t, len(onConnectEvent), 0)
    testOnGnetConnectSolicitation(t, c, addr, true)
    assert.Equal(t, len(onConnectEvent), 0)
}

/* Helpers */

func wait() {
    time.Sleep(time.Millisecond * 50)
}
