package daemon

import (
	"time"
	//"github.com/skycoin/skycoin/src/daemon/gnet"
)

// func TestInitPool(t *testing.T) {
// 	d := newDefaultDaemon()
// 	pool := NewPool(NewPoolConfig())
// 	assert.Nil(t, pool.Pool)
// 	assert.NotPanics(t, func() { pool.Init(d) })
// 	assert.Equal(t, pool.Pool.Config.DialTimeout, pool.Config.DialTimeout)
// 	assert.NotNil(t, pool.Pool.Config.DisconnectCallback)
// 	assert.NotNil(t, pool.Pool.Config.ConnectCallback)
// 	wait()
// 	go func() {
// 		assert.NotPanics(t, pool.Start)
// 	}()
// 	wait()
// 	// A second call to start listen will panic due to the pool already
// 	// listening
// 	assert.Panics(t, pool.Start)
// 	pool.Pool.StopListen()
// 	wait()
// 	pool.Shutdown()
// 	shutdown(d)
// }

// func TestShutdownPool(t *testing.T) {
// 	// Shutting down should flush the DisconnectQueue, among other things
// 	d := newDefaultDaemon()
// 	pool := d.Pool
// 	assert.Equal(t, len(pool.Pool.DisconnectQueue), 0)
// 	pool.Pool.DisconnectQueue <- gnet.DisconnectEvent{
// 		ConnId: 1,
// 		Reason: DisconnectOtherError,
// 	}
// 	assert.NotPanics(t, d.Pool.Shutdown)
// 	wait()
// 	// pool.Shutdown() should call Pool.StopListen() which closes and resets
// 	// the DisconnectQueue channel
// 	assert.Equal(t, len(pool.Pool.DisconnectQueue), 0)
// 	shutdown(d)
// }

// func testOnGnetConnectSolicitation(t *testing.T, d *Daemon,
// 	c *gnet.Connection, addr string, sol bool) {
// 	d.onGnetConnect(c, sol)
// 	assert.Equal(t, len(d.onConnectEvent), 1)
// 	if len(d.onConnectEvent) == 0 {
// 		t.Fatalf("onConnectEvent is not empty, would block")
// 	}
// 	ce := <-d.onConnectEvent
// 	assert.Equal(t, ce.Addr, addr)
// 	assert.Equal(t, ce.Solicited, sol)
// }

// func TestOnGnetConnect(t *testing.T) {
// 	d := newDefaultDaemon()
// 	addr := "11.22.33.44:5555"
// 	c := gnetConnection(addr)
// 	assert.Equal(t, len(d.onConnectEvent), 0)
// 	testOnGnetConnectSolicitation(t, d, c, addr, false)
// 	assert.Equal(t, len(d.onConnectEvent), 0)
// 	testOnGnetConnectSolicitation(t, d, c, addr, true)
// 	assert.Equal(t, len(d.onConnectEvent), 0)
// 	shutdown(d)
// }

/* Helpers */

func wait() {
	time.Sleep(time.Millisecond * 50)
}
