package daemon

// import (
// 	"errors"
// 	"os"
// 	"os/signal"
// 	"runtime/pprof"
// 	"strings"
// 	"syscall"
// 	"testing"
// 	"time"

// 	//"github.com/skycoin/skycoin/src/daemon/gnet"
// 	"github.com/skycoin/skycoin/src/daemon/gnet"

// 	"github.com/skycoin/skycoin/src/daemon/pex"
// 	//"github.com/skycoin/skycoin/src/coin"
// 	"github.com/skycoin/skycoin/src/util"
// 	//"github.com/skycoin/skycoin/src/wallet"
// 	"github.com/stretchr/testify/assert"
// )

// func catchSigusr1() {
// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGUSR1)
// 	for {
// 		select {
// 		case <-sigchan:
// 			pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
// 		}
// 	}
// }

// func setupExistingPool(p *Pool) *gnet.Connection {
// 	gc := gnetConnection(addr)
// 	p.Pool.Pool[gc.Id] = gc
// 	p.Pool.Addresses[gc.Addr()] = gc
// 	return gc
// }

// func setupPool() (*Pool, *gnet.Connection) {
// 	m := NewMessagesConfig()
// 	m.Register()
// 	p := NewPool(NewPoolConfig())
// 	p.Init(nil)
// 	return p, setupExistingPool(p)
// }

// func newDefaultDaemon() *Daemon {
// 	cleanupPeers()
// 	c := NewConfig()
// 	//we := wallet.NewWalletEntry()
// 	//c.Visor.Config.MasterKeys = we
// 	//c.Visor.Config.GenesisSignature = createGenesisSignature(we)
// 	c.Visor.Disabled = true
// 	// c.DHT.Disabled = true
// 	return NewDaemon(c)
// }

// func newDHTDaemon() *Daemon {
// 	cleanupPeers()
// 	c := NewConfig()
// 	c.Visor.Disabled = true
// 	// c.DHT.Disabled = false
// 	return NewDaemon(c)
// }

// func setupDaemonLoop() (*Daemon, chan int) {
// 	return newDefaultDaemon(), make(chan int)
// }

// func setupDaemonLoopDHT() (*Daemon, chan int) {
// 	return newDHTDaemon(), make(chan int)
// }

// func closeDaemon(d *Daemon, quit chan int) {
// 	wait()
// 	logger.Critical("sending quit")
// 	quit <- 1
// 	logger.Critical("quit sent")
// 	shutdown(d)
// }

// func shutdown(d *Daemon) {
// 	d.Shutdown()
// 	wait()
// 	cleanupPeers()
// 	gnet.EraseMessages()
// }

// func TestConfigPreprocess(t *testing.T) {
// 	c := NewConfig()
// 	a := "127.0.0.1"
// 	p := 7779
// 	// Test that addr, port are copied to subconfigs
// 	c.Daemon.Port = p
// 	c.Daemon.Address = a
// 	d := c.preprocess()
// 	assert.Equal(t, d.Pool.port, p)
// 	assert.Equal(t, d.Pool.address, a)
// 	// assert.Equal(t, d.DHT.port, p)

// 	// Test localhost only with localhost addr
// 	c = NewConfig()
// 	c.Daemon.LocalhostOnly = true
// 	c.Daemon.Address = a
// 	assert.NotPanics(t, func() { c.preprocess() })
// 	d = c.preprocess()
// 	// assert.True(t, d.DHT.Disabled)
// 	assert.Equal(t, d.Pool.address, a)
// 	assert.True(t, d.Peers.AllowLocalhost)

// 	// Test localhost only with unassigned addr
// 	c = NewConfig()
// 	c.Daemon.LocalhostOnly = true
// 	c.Daemon.Address = ""
// 	assert.NotPanics(t, func() { c.preprocess() })
// 	d = c.preprocess()
// 	assert.True(t, IsLocalhost(d.Daemon.Address))
// 	assert.True(t, IsLocalhost(d.Pool.address))
// 	assert.True(t, d.Peers.AllowLocalhost)

// 	// Test localhost only with nonlocal addr
// 	c = NewConfig()
// 	c.Daemon.LocalhostOnly = true
// 	c.Daemon.Address = "11.22.33.44"
// 	assert.Panics(t, func() { c.preprocess() })

// 	// Test disable networking disables all
// 	c = NewConfig()
// 	c.Daemon.DisableNetworking = true
// 	d = c.preprocess()
// 	assert.True(t, d.Daemon.DisableNetworking)
// 	assert.True(t, d.Daemon.DisableOutgoingConnections)
// 	assert.True(t, d.Daemon.DisableIncomingConnections)
// 	// assert.True(t, d.DHT.Disabled)
// 	assert.True(t, d.Peers.Disabled)

// 	// Test coverage for logging statements
// 	c = NewConfig()
// 	c.Daemon.DisableNetworking = false
// 	c.Daemon.DisableIncomingConnections = true
// 	c.Daemon.DisableOutgoingConnections = true
// 	assert.NotPanics(t, func() { c.preprocess() })
// }

// func TestGetListenPort(t *testing.T) {
// 	d := newDefaultDaemon()
// 	// No connectionMirror found
// 	assert.Equal(t, d.GetListenPort(addr), uint16(0))
// 	// No mirrorConnection map exists
// 	d.ConnectionMirrors[addr] = uint32(4)
// 	assert.Panics(t, func() { d.GetListenPort(addr) })
// 	// Everything is good
// 	m := make(map[string]uint16)
// 	d.mirrorConnections[uint32(4)] = m
// 	m[addrIP] = uint16(6667)
// 	assert.Equal(t, d.GetListenPort(addr), uint16(6667))
// 	shutdown(d)
// 	// Bad addr
// 	d.ConnectionMirrors["xcasca"] = uint32(4)
// 	d.mirrorConnections[uint32(4)] = m
// 	assert.Equal(t, d.GetListenPort("xcasca"), uint16(0))
// }

// func TestStart(t *testing.T) {
// 	gnet.EraseMessages()
// 	defer cleanupPeers()
// 	d, quit := setupDaemonLoopDHT()
// 	defer closeDaemon(d, quit)
// 	assert.NotNil(t, d)
// 	assert.NotNil(t, d.Pool)
// 	// assert.NotNil(t, d.DHT)
// 	go d.Start(quit)
// 	wait()
// 	assert.NotEqual(t, len(gnet.MessageIdMap), 0)
// 	assert.NotNil(t, d.Pool)
// 	assert.NotNil(t, d.Peers)
// 	// assert.NotNil(t, d.DHT)
// 	assert.NotNil(t, d.Messages)
// 	assert.NotNil(t, d.Gateway)
// }

// func TestShutdown(t *testing.T) {
// 	cleanupPeers()
// 	d := newDHTDaemon()
// 	d.Peers.Peers.AddPeer(addr)
// 	d.Pool.Pool.DisconnectQueue <- gnet.DisconnectEvent{
// 		ConnId: 1,
// 		Reason: DisconnectOtherError,
// 	}
// 	assert.NotPanics(t, func() { d.Shutdown() })
// 	confirmPeersShutdown(t)
// 	assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 0)
// 	// assert.Nil(t, d.DHT.DHT)
// 	cleanupPeers()
// }

// func TestDaemonLoopDisabledPanics(t *testing.T) {
// 	c := NewConfig()
// 	c.Daemon.DisableNetworking = true
// 	c.Visor.Disabled = true
// 	d := NewDaemon(c)
// 	quit := make(chan int)
// 	done := make(chan bool)
// 	panics := func() {
// 		assert.Panics(t, func() { d.Start(quit) })
// 		done <- true
// 	}

// 	d.onConnectEvent <- ConnectEvent{}
// 	go panics()
// 	<-done

// 	d.connectionErrors <- ConnectionError{}
// 	go panics()
// 	<-done

// 	// d.DHT.DHT.PeersRequestResults <- make(map[dht.InfoHash][]string)
// 	// go panics()
// 	// <-done

// 	d.Pool.Pool.DisconnectQueue <- gnet.DisconnectEvent{}
// 	go panics()
// 	<-done

// 	d.messageEvents <- MessageEvent{}
// 	go panics()
// 	<-done

// 	d.Pool.Pool.SendResults <- gnet.SendResult{}
// 	go panics()
// 	<-done

// 	shutdown(d)
// }

// func TestDaemonLoopQuit(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	done := false
// 	go func() {
// 		d.Start(quit)
// 		done = true
// 	}()
// 	wait()
// 	quit <- 1
// 	wait()
// 	assert.True(t, done)
// 	shutdown(d)
// }

// func TestDaemonLoopApiRequest(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)

// 	conn := make(chan interface{})
// 	d.Gateway.Requests <- func() {
// 		conn <- &Connection{Id: 7}
// 	}

// 	resp := <-conn

// 	assert.Equal(t, resp.(*Connection).Id, 7)
// }

// func TestDaemonLoopOnConnectEvent(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)
// 	d.pendingConnections[addr] = pex.NewPeer(addr)
// 	d.onConnectEvent <- ConnectEvent{addr, false}
// 	wait()
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Nil(t, d.pendingConnections[addr])
// }

// func TestDaemonLoopConnectionErrors(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)
// 	d.pendingConnections[addr] = pex.NewPeer(addr)
// 	d.connectionErrors <- ConnectionError{addr, errors.New("failed")}
// 	wait()
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Nil(t, d.pendingConnections[addr])
// }

// func TestDaemonLoopDisconnectQueue(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)
// 	d.Pool.Pool.Pool[1] = gnetConnection(addr)
// 	e := gnet.DisconnectEvent{ConnId: 1, Reason: DisconnectIdle}
// 	d.Pool.Pool.DisconnectQueue <- e
// 	wait()
// 	assert.Equal(t, len(d.Pool.Pool.Pool), 0)
// }

// /*
// //networking loop and response gets moved into visor
// func TestDaemonLoopSendResults(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)
// 	c := gnetConnection(addr)
// 	d.Pool.Pool.Pool[1] = c
// 	vc := NewVisorConfig()
// 	v := NewVisor(vc)
// 	d.Visor = v
// 	txn := addUnconfirmedTxn(d.Visor)
// 	ut := d.Visor.Visor.Unconfirmed.Txns[txn.Hash()]
// 	assert.True(t, ut.Announced.IsZero())
// 	txns := coin.Transactions{txn.Txn}
// 	m := NewAnnounceTxnsMessage(txns.Hashes())
// 	sr := gnet.SendResult{Connection: c, Error: nil, Message: m}
// 	d.Pool.Pool.SendResults <- sr
// 	wait()
// 	ut = d.Visor.Visor.Unconfirmed.Txns[txn.Hash()]
// 	assert.False(t, ut.Announced.IsZero())
// }
// */

// type DummyAsyncMessage struct {
// 	fn func()
// }

// func (self *DummyAsyncMessage) Process(d *Daemon) {
// 	self.fn()
// }

// func TestDaemonLoopMessageEvent(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)
// 	called := false
// 	m := &DummyAsyncMessage{fn: func() { called = true }}
// 	d.messageEvents <- MessageEvent{m, messageContext(addr)}
// 	wait()
// 	assert.True(t, called)
// }

// // func TestDaemonLoopDHTResults(t *testing.T) {
// // 	d, quit := setupDaemonLoopDHT()
// // 	defer closeDaemon(d, quit)
// // 	assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
// // 	go d.Start(quit)
// // 	m := make(map[dht.InfoHash][]string, 1)
// // 	m[d.DHT.InfoHash] = []string{"abcdef"}
// // 	d.DHT.DHT.PeersRequestResults <- m
// // 	wait()
// // 	assert.Equal(t, len(d.Peers.Peers.Peerlist), 1)
// // 	assert.NotNil(t, d.Peers.Peers.Peerlist["97.98.99.100:25958"])
// // }

// // func testDaemonLoopDHTBootstrapTicker(t *testing.T, d *Daemon, quit chan int) {
// // 	d.DHT.Config.BootstrapRequestRate = time.Millisecond * 10
// // 	go d.Start(quit)
// // 	// Can't really test DHT internals, but we'll know if it crashes or not
// // 	time.Sleep(time.Millisecond * 15)
// // 	d.DHT.Config.PeerLimit = 0
// // 	time.Sleep(time.Millisecond * 15)
// // }

// // func TestDaemonLoopDHTBootstrapTicker(t *testing.T) {
// // 	d, quit := setupDaemonLoopDHT()
// // 	defer closeDaemon(d, quit)
// // 	testDaemonLoopDHTBootstrapTicker(t, d, quit)
// // }

// // func TestDaemonLoopDHTBootstrapTickerDisabled(t *testing.T) {
// // 	d, quit := setupDaemonLoopDHT()
// // 	defer closeDaemon(d, quit)
// // 	d.DHT.Config.Disabled = true
// // 	testDaemonLoopDHTBootstrapTicker(t, d, quit)
// // }

// func testDaemonLoopBlacklistTicker(t *testing.T, d *Daemon, quit chan int,
// 	count int) {
// 	assert.Equal(t, len(d.Peers.Peers.Blacklist), 0)
// 	d.Peers.Peers.AddBlacklistEntry(addr, time.Millisecond)
// 	assert.Equal(t, len(d.Peers.Peers.Blacklist), 1)
// 	d.Peers.Config.UpdateBlacklistRate = time.Millisecond * 10
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	assert.Equal(t, len(d.Peers.Peers.Blacklist), count)
// }

// func TestDaemonLoopBlacklistTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopBlacklistTicker(t, d, quit, 0)
// }

// func TestDaemonLoopBlacklistTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Peers.Config.Disabled = true
// 	testDaemonLoopBlacklistTicker(t, d, quit, 1)
// }

// func testDaemonLoopCullInvalidTicker(t *testing.T, d *Daemon, quit chan int,
// 	count int) {
// 	d.ExpectingIntroductions[addr] = time.Now().Add(-time.Hour)
// 	d.Config.CullInvalidRate = time.Millisecond * 10
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	assert.Equal(t, len(d.ExpectingIntroductions), count)
// }

// func TestDaemonLoopCullInvalidTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopCullInvalidTicker(t, d, quit, 0)
// }

// func TestDaemonLoopCullInvalidTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Config.DisableNetworking = true
// 	testDaemonLoopCullInvalidTicker(t, d, quit, 1)
// }

// func testDaemonLoopRequestPeersTicker(t *testing.T, d *Daemon, quit chan int,
// 	sent bool) {
// 	c := gnetConnection(addr)
// 	go d.Pool.Pool.ConnectionWriteLoop(c)
// 	d.Pool.Pool.Pool[1] = c
// 	d.Pool.Pool.Addresses[c.Addr()] = c
// 	assert.True(t, c.LastSent.IsZero())
// 	d.Peers.Config.RequestRate = time.Millisecond * 10
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	wait()
// 	if sent {
// 		assert.False(t, c.LastSent.IsZero())
// 	} else {
// 		assert.True(t, c.LastSent.IsZero())
// 	}
// }

// func TestDaemonLoopRequestPeersTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopRequestPeersTicker(t, d, quit, true)
// }

// func TestDaemonLoopRequestPeersTickerFull(t *testing.T) {
// 	cfg := NewConfig()
// 	cfg.Visor.Disabled = true
// 	cfg.Peers.Max = 1
// 	d := NewDaemon(cfg)
// 	d.Peers.Peers.AddPeer(addr)
// 	quit := make(chan int)
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopRequestPeersTicker(t, d, quit, false)
// }

// func TestDaemonLoopRequestPeersTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Peers.Config.Disabled = true
// 	testDaemonLoopRequestPeersTicker(t, d, quit, false)
// }

// func testDaemonLoopClearOldPeersTicker(t *testing.T, d *Daemon, quit chan int,
// 	count int) {
// 	p := pex.NewPeer(addr)
// 	p.LastSeen = time.Time()
// 	d.Peers.Peers.Peerlist[addr] = p
// 	d.Peers.Config.CullRate = time.Millisecond * 10
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	assert.Equal(t, len(d.Peers.Peers.Peerlist), count)
// }

// func TestDaemonLoopClearOldPeersTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopClearOldPeersTicker(t, d, quit, 0)
// }

// func TestDaemonLoopClearOldPeersTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Peers.Config.Disabled = true
// 	testDaemonLoopClearOldPeersTicker(t, d, quit, 1)
// }

// func testDaemonLoopClearStaleConnectionsTicker(t *testing.T, d *Daemon,
// 	quit chan int, poolCount int) {
// 	c := gnetConnection(addr)
// 	c.LastReceived = time.Time{}
// 	d.Pool.Pool.Pool[c.Id] = c
// 	d.Pool.Config.ClearStaleRate = time.Millisecond * 10
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	assert.Equal(t, len(d.Pool.Pool.Pool), poolCount)
// }

// func TestDaemonLoopClearStaleConnectionsTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopClearStaleConnectionsTicker(t, d, quit, 0)
// }

// func TestDaemonLoopClearStaleConnectionsTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Config.DisableNetworking = true
// 	testDaemonLoopClearStaleConnectionsTicker(t, d, quit, 1)
// }

// func testDaemonLoopPingCheckTicker(t *testing.T, d *Daemon, quit chan int,
// 	sent bool) {
// 	c := gnetConnection(addr)
// 	go d.Pool.Pool.ConnectionWriteLoop(c)
// 	c.LastSent = time.Time{}
// 	d.Pool.Pool.Pool[c.Id] = c
// 	d.Pool.Config.IdleCheckRate = time.Millisecond * 10
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	if sent {
// 		assert.False(t, c.LastSent.IsZero())
// 	} else {
// 		assert.True(t, c.LastSent.IsZero())
// 	}
// 	c.Close()
// 	wait()
// }

// func TestDaemonLoopPingCheckTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopPingCheckTicker(t, d, quit, true)
// }

// func TestDaemonLoopPingCheckTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Config.DisableNetworking = true
// 	testDaemonLoopPingCheckTicker(t, d, quit, false)
// }

// func testDaemonLoopOutgoingConnectionsTicker(t *testing.T, d *Daemon,
// 	quit chan int, peerCount int) {
// 	d.Pool.Pool.Config.DialTimeout = 1 // nanosecond
// 	d.Config.OutgoingRate = time.Millisecond * 10
// 	d.Peers.Peers.AddPeer(addr)
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	// Should have made a connection attempt, timed out, put an error
// 	// the queue, handled by d.Run, resulting in the peer being removed
// 	assert.Equal(t, len(d.Peers.Peers.Peerlist), peerCount)
// }

// func TestDaemonLoopOutgoingConnectionsTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopOutgoingConnectionsTicker(t, d, quit, 0)
// }

// func TestDaemonLoopOutgoingConnectionsTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Config.DisableOutgoingConnections = true
// 	testDaemonLoopOutgoingConnectionsTicker(t, d, quit, 1)
// }

// func TestDaemonLoopOutgoingConnectionsTickerOutgoingMax(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Config.OutgoingMax = 0
// 	testDaemonLoopOutgoingConnectionsTicker(t, d, quit, 1)
// }

// func TestDaemonLoopOutgoingConnectionsTickerPendingMax(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Config.PendingMax = 0
// 	testDaemonLoopOutgoingConnectionsTicker(t, d, quit, 1)
// }

// func TestDaemonLoopPrivateConnectionsTicker(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	d.Config.DisableOutgoingConnections = false
// 	d.Config.PrivateRate = time.Millisecond * 10
// 	d.Config.OutgoingRate = time.Hour
// 	// Add private peer
// 	d.Peers.Peers.AllowLocalhost = true
// 	addr := "192.168.1.155:45339"
// 	p, err := d.Peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	p.Private = true
// 	// Add public peer. This one shouldn't get sent to
// 	r, err := d.Peers.Peers.AddPeer(addrb)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, r)
// 	r.Private = false
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	assert.Equal(t, len(d.pendingConnections), 1)
// 	assert.NotNil(t, d.pendingConnections[addr])
// }

// func TestDaemonLoopPrivateConnectionsTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	d.Config.DisableOutgoingConnections = true
// 	d.Config.PrivateRate = time.Millisecond * 10
// 	// Add private peer
// 	p, err := d.Peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	p.Private = true
// 	// Add public peer. This one shouldn't get sent to
// 	r, err := d.Peers.Peers.AddPeer(addrb)
// 	assert.Nil(t, err)
// 	r.Private = false
// 	defer closeDaemon(d, quit)
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	// Nothing should happen if disabled
// 	assert.Equal(t, len(d.pendingConnections), 0)
// }

// func testDaemonLoopMessageHandlingTicker(t *testing.T, d *Daemon,
// 	quit chan int) {
// 	d.Pool.Config.MessageHandlingRate = time.Millisecond * 10
// 	go d.Start(quit)
// 	time.Sleep(time.Millisecond * 15)
// 	// Can't test Pool internals from here, just see that it doesn't crash
// }

// func TestDaemonLoopMessageHandlingTickerD(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	testDaemonLoopMessageHandlingTicker(t, d, quit)
// }

// func TestDaemonLoopMessageHandlingTickerDisabled(t *testing.T) {
// 	d, quit := setupDaemonLoop()
// 	defer closeDaemon(d, quit)
// 	d.Config.DisableNetworking = true
// 	testDaemonLoopMessageHandlingTicker(t, d, quit)
// }

// func TestClearStaleConnections(t *testing.T) {
// 	dm := newDefaultDaemon()
// 	defer shutdown(dm)
// 	c := gnetConnection(addr)
// 	d := gnetConnection(addrb)
// 	c.LastReceived = time.Time{}
// 	d.LastReceived = time.Now()
// 	dm.Pool.Pool.Pool[1] = c
// 	dm.Pool.Pool.Pool[2] = d
// 	assert.NotPanics(t, dm.Pool.clearStaleConnections)
// 	assert.Equal(t, len(dm.Pool.Pool.DisconnectQueue), 1)
// 	if len(dm.Pool.Pool.DisconnectQueue) == 0 {
// 		t.Fatalf("Empty DisconnectQueue, would block")
// 	}
// 	de := <-dm.Pool.Pool.DisconnectQueue
// 	assert.Equal(t, de.ConnId, 1)
// 	assert.Equal(t, de.Reason, DisconnectIdle)
// }

// func TestSendPings(t *testing.T) {
// 	d := newDefaultDaemon()
// 	defer shutdown(d)
// 	c := gnetConnection(addr)
// 	go d.Pool.Pool.ConnectionWriteLoop(c)
// 	d.Pool.Pool.Pool[1] = c
// 	assert.NotPanics(t, d.Pool.sendPings)
// 	wait()
// 	assert.Equal(t, len(d.Pool.Pool.SendResults), 1)
// 	if len(d.Pool.Pool.SendResults) == 0 {
// 		t.Fatal("SendResults empty, would block")
// 	}
// 	sr := <-d.Pool.Pool.SendResults
// 	assert.Equal(t, sr.Connection, c)
// 	assert.Nil(t, sr.Error)
// 	_, ok := sr.Message.(*PingMessage)
// 	assert.True(t, ok)
// 	assert.False(t, c.LastSent.IsZero())

// 	// No pings should be sent, since we just pinged
// 	lastSent := c.LastSent
// 	assert.NotPanics(t, d.Pool.sendPings)
// 	wait()
// 	assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
// 	assert.Equal(t, c.LastSent, lastSent)
// }

// func TestMakePrivateConnections(t *testing.T) {
// 	d := newDefaultDaemon()
// 	defer shutdown(d)
// 	d.Pool.Config.DialTimeout = time.Hour
// 	addr := "192.168.1.198:43071"
// 	addrb := "192.168.1.197:43072"
// 	// Two privates
// 	p, err := d.Peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	p.Private = true
// 	p, err = d.Peers.Peers.AddPeer(addrb)
// 	assert.Nil(t, err)
// 	p.Private = true
// 	// Not private
// 	p, err = d.Peers.Peers.AddPeer(addrc)
// 	assert.Nil(t, err)
// 	p.Private = false

// 	// Disabled
// 	d.Config.DisableOutgoingConnections = true
// 	d.makePrivateConnections()
// 	assert.Equal(t, len(d.pendingConnections), 0)

// 	// Enabled
// 	d.Config.DisableOutgoingConnections = false
// 	d.makePrivateConnections()
// 	assert.Equal(t, len(d.pendingConnections), 2)
// 	assert.NotNil(t, d.pendingConnections[addr])
// 	assert.NotNil(t, d.pendingConnections[addrb])
// 	assert.Nil(t, d.pendingConnections[addrc])
// }

// func TestConnectToPeer(t *testing.T) {
// 	d := newDefaultDaemon()
// 	defer shutdown(d)
// 	defer cleanupPeers()
// 	d.Pool.Pool.Config.DialTimeout = 1 // nanosecond

// 	assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)

// 	// PEX somehow has invalid peer
// 	p := &pex.Peer{Addr: "xcasca"}
// 	d.Peers.Peers.Peerlist["xcasca"] = p
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Equal(t, len(d.connectionErrors), 0)
// 	delete(d.Peers.Peers.Peerlist, "xcasca")

// 	// Disabled outgoing conns
// 	d.Config.DisableOutgoingConnections = true
// 	p, err := d.Peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Equal(t, len(d.connectionErrors), 0)
// 	d.Config.DisableOutgoingConnections = false

// 	// Localhost only, and peer isn't
// 	d.Config.LocalhostOnly = true
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Equal(t, len(d.connectionErrors), 0)
// 	delete(d.Peers.Peers.Peerlist, addr)

// 	// Valid attempt to connect to localhost
// 	d.Peers.Peers.AllowLocalhost = true
// 	p, err = d.Peers.Peers.AddPeer(localAddr)
// 	assert.Nil(t, err)
// 	if err != nil {
// 		t.Logf("Error with local addr: %v\n", err)
// 	}
// 	t.Logf("Peerlist: %v\n", d.Peers.Peers.Peerlist)
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	wait()
// 	assert.Equal(t, len(d.pendingConnections), 1)
// 	assert.Equal(t, len(d.connectionErrors), 1)
// 	if len(d.connectionErrors) == 0 {
// 		shutdown(d)
// 		cleanupPeers()
// 		t.Fatalf("connectionErrors empty, would block")
// 	}
// 	ce := <-d.connectionErrors
// 	assert.Equal(t, ce.Addr, localAddr)
// 	assert.NotNil(t, ce.Error)
// 	delete(d.pendingConnections, localAddr)
// 	delete(d.Peers.Peers.Peerlist, localAddr)
// 	d.Config.LocalhostOnly = false
// 	d.Peers.Peers.AllowLocalhost = false

// 	// Valid attempt to connect
// 	p, err = d.Peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	wait()
// 	assert.Equal(t, len(d.pendingConnections), 1)
// 	assert.Equal(t, len(d.connectionErrors), 1)
// 	if len(d.connectionErrors) == 0 {
// 		t.Fatalf("connectionErrors empty, would block")
// 	}
// 	ce = <-d.connectionErrors
// 	assert.Equal(t, ce.Addr, addr)
// 	assert.NotNil(t, ce.Error)
// 	delete(d.pendingConnections, addr)

// 	// Two peers, one successful connect attempt and one skipped
// 	p, err = d.Peers.Peers.AddPeer(addrb)
// 	assert.Nil(t, err)
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	wait()
// 	assert.Equal(t, len(d.pendingConnections), 1)
// 	assert.Equal(t, len(d.connectionErrors), 1)
// 	if len(d.connectionErrors) == 0 {
// 		t.Fatalf("connectionErrors empty, would block")
// 	}
// 	ce = <-d.connectionErrors
// 	assert.True(t, (ce.Addr == addr) || (ce.Addr == addrb))
// 	assert.NotNil(t, ce.Error)
// 	delete(d.pendingConnections, addr)
// 	delete(d.pendingConnections, addrb)
// 	delete(d.Peers.Peers.Peerlist, addrb)

// 	// Already connected, skip
// 	p, err = d.Peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	d.Pool.Pool.Addresses[addr] = gnetConnection(addr)
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Equal(t, len(d.connectionErrors), 0)
// 	delete(d.Pool.Pool.Addresses, addr)

// 	// Pending connection, skip
// 	p = pex.NewPeer(addr)
// 	d.pendingConnections[addr] = p
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	assert.Equal(t, len(d.pendingConnections), 1)
// 	assert.Equal(t, len(d.connectionErrors), 0)
// 	delete(d.pendingConnections, addr)

// 	// Already connected to this base IP at least once, skip
// 	d.ipCounts[addrIP] = 1
// 	assert.NotPanics(t, func() { d.connectToPeer(p) })
// 	assert.Equal(t, len(d.ipCounts), 1)
// 	assert.Equal(t, d.ipCounts[addrIP], 1)
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Equal(t, len(d.connectionErrors), 0)
// 	delete(d.ipCounts, addrIP)
// }

// func TestConnectToRandomPeer(t *testing.T) {
// 	d := newDefaultDaemon()
// 	defer shutdown(d)
// 	defer cleanupPeers()

// 	// Disabled
// 	d.Config.DisableOutgoingConnections = true
// 	assert.NotPanics(t, d.connectToRandomPeer)
// 	assert.Equal(t, len(d.pendingConnections), 0)

// 	// Enabled, but only private peers
// 	addr := "192.168.1.196:30954"
// 	p, err := d.Peers.Peers.AddPeer(addr)
// 	assert.Nil(t, err)
// 	p.Private = true
// 	d.Config.DisableOutgoingConnections = false
// 	assert.NotPanics(t, d.connectToRandomPeer)
// 	assert.Equal(t, len(d.pendingConnections), 0)

// 	// Enabled, and we have a public peer
// 	p.Private = false
// 	assert.NotPanics(t, d.connectToRandomPeer)
// 	assert.Equal(t, len(d.pendingConnections), 1)
// 	assert.NotNil(t, d.pendingConnections[addr])
// }

// func TestHandleConnectionError(t *testing.T) {
// 	d := newDefaultDaemon()
// 	p, _ := d.Peers.Peers.AddPeer(addr)
// 	d.pendingConnections[addr] = p
// 	assert.NotPanics(t, func() {
// 		d.handleConnectionError(ConnectionError{addr, nil})
// 	})
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
// 	p, _ = d.Peers.Peers.AddPeer(addr)
// 	d.pendingConnections[addr] = p
// 	assert.NotPanics(t, func() {
// 		d.handleConnectionError(ConnectionError{addr, errors.New("bad")})
// 	})
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
// 	shutdown(d)
// }

// // func TestCullInvalidConnections(t *testing.T) {
// // 	d := newDefaultDaemon()
// // 	// Is fine
// // 	d.ExpectingIntroductions[addr] = time.Now()
// // 	// Is expired
// // 	d.ExpectingIntroductions[addrb] = time.Time{}
// // 	// Is not in pool
// // 	d.ExpectingIntroductions[addrc] = time.Time{}
// // 	d.Peers.Peers.AddPeer(addr)
// // 	d.Peers.Peers.AddPeer(addrb)
// // 	d.Peers.Peers.AddPeer(addrc)
// // 	d.Pool.Pool.Addresses[addr] = gnetConnection(addr)
// // 	d.Pool.Pool.Addresses[addrb] = gnetConnection(addrb)
// // 	d.Pool.Pool.Addresses[addrb].Id = 2
// // 	d.Pool.Pool.Pool[1] = d.Pool.Pool.Addresses[addr]
// // 	d.Pool.Pool.Pool[2] = d.Pool.Pool.Addresses[addrb]

// // 	assert.NotPanics(t, d.cullInvalidConnections)

// // 	assert.Equal(t, len(d.ExpectingIntroductions), 1)
// // 	assert.Equal(t, len(d.Peers.Peers.Peerlist), 2)
// // 	assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
// // 	if len(d.Pool.Pool.DisconnectQueue) == 0 {
// // 		t.Fatal("pool.Pool.DisconnectQueue not empty, would block")
// // 	}
// // 	de := <-d.Pool.Pool.DisconnectQueue
// // 	assert.Equal(t, de.ConnId, 2)
// // 	assert.Equal(t, de.Reason, DisconnectIntroductionTimeout)
// // 	shutdown(d)
// // }

// func TestRecordMessageEventValid(t *testing.T) {
// 	d := newDefaultDaemon()
// 	// Valid message, not expecting Introduction
// 	assert.Equal(t, len(d.messageEvents), 0)
// 	delete(d.ExpectingIntroductions, addr)
// 	m := &PingMessage{}
// 	m.c = messageContext(addr)
// 	err := d.recordMessageEvent(m, m.c)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(d.messageEvents), 1)
// 	if len(d.messageEvents) == 0 {
// 		t.Fatal("d.messageEvents empty, would block")
// 	}
// 	me := <-d.messageEvents
// 	_, ok := me.Message.(*PingMessage)
// 	assert.True(t, ok)
// 	shutdown(d)
// }

// func TestRecordMessageEventIsIntroduction(t *testing.T) {
// 	// Needs Introduction and thats what it has received
// 	d := newDefaultDaemon()
// 	d.ExpectingIntroductions[addr] = utc.Now()
// 	assert.Equal(t, len(d.messageEvents), 0)
// 	m := NewIntroductionMessage(d.Messages.Mirror, d.Config.Version,
// 		d.Pool.Pool.Config.Port)
// 	m.c = messageContext(addr)
// 	err := d.recordMessageEvent(m, m.c)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(d.messageEvents), 1)
// 	if len(d.messageEvents) == 0 {
// 		t.Fatal("d.messageEvents empty, would block")
// 	}
// 	me := <-d.messageEvents
// 	_, ok := me.Message.(*IntroductionMessage)
// 	assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 0)
// 	assert.True(t, ok)
// 	delete(d.ExpectingIntroductions, addr)
// 	shutdown(d)
// }

// func TestRecordMessageEventNeedsIntroduction(t *testing.T) {
// 	// Needs Introduction but didn't get it first
// 	d := newDefaultDaemon()
// 	m := &PingMessage{}
// 	m.c = messageContext(addr)
// 	d.Pool.Pool.Addresses[addr] = m.c.Conn
// 	d.Pool.Pool.Pool[m.c.Conn.Id] = m.c.Conn
// 	assert.Equal(t, len(d.messageEvents), 0)
// 	d.ExpectingIntroductions[addr] = utc.Now()
// 	d.processMessageEvent(MessageEvent{m, m.c})
// 	assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
// 	if len(d.Pool.Pool.DisconnectQueue) == 0 {
// 		t.Fatal("DisconnectQueue empty, would block")
// 	}
// 	de := <-d.Pool.Pool.DisconnectQueue
// 	assert.Equal(t, de.ConnId, m.c.Conn.Id)
// 	assert.Equal(t, de.Reason, DisconnectNoIntroduction)
// 	delete(d.ExpectingIntroductions, addr)
// 	shutdown(d)
// }

// func TestOnConnect(t *testing.T) {
// 	d := newDefaultDaemon()

// 	// Test a valid connection, unsolicited
// 	e := ConnectEvent{addr, false}
// 	p, _ := d.Peers.Peers.AddPeer(addr)

// 	c := setupExistingPool(d.Pool)

// 	go d.Pool.Pool.ConnectionWriteLoop(c)
// 	d.pendingConnections[addr] = p
// 	assert.NotPanics(t, func() { d.onConnect(e) })
// 	wait()
// 	assert.Equal(t, len(d.Pool.Pool.SendResults), 1)
// 	if len(d.Pool.Pool.SendResults) == 0 {
// 		t.Fatalf("SendResults empty, would block")
// 	}
// 	sr := <-d.Pool.Pool.SendResults
// 	assert.Equal(t, sr.Connection, c)
// 	assert.Nil(t, sr.Error)
// 	_, ok := sr.Message.(*IntroductionMessage)
// 	assert.True(t, ok)
// 	// This connection should no longer be pending
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	// This is not an outgoing connection, we did not solicit it
// 	assert.Equal(t, len(d.OutgoingConnections), 0)
// 	// We should be expecting its version
// 	assert.Equal(t, len(d.ExpectingIntroductions), 1)
// 	_, exists := d.ExpectingIntroductions[addr]
// 	assert.True(t, exists)
// 	// An introduction should have been sent
// 	assert.False(t, c.LastSent.IsZero())
// 	// d.ipCounts should be 1
// 	assert.Equal(t, d.ipCounts[addrIP], 1)
// 	// Cleanup
// 	delete(d.ipCounts, addrIP)
// 	delete(d.ExpectingIntroductions, addr)
// 	c.Close()

// 	// Test a valid connection, solicited
// 	e = ConnectEvent{addr, true}
// 	c = gnetConnection(addr)
// 	go d.Pool.Pool.ConnectionWriteLoop(c)
// 	d.pendingConnections[addr] = p
// 	d.Pool.Pool.Addresses[addr] = c
// 	d.Pool.Pool.Pool[c.Id] = c
// 	assert.NotPanics(t, func() { d.onConnect(e) })
// 	wait()
// 	assert.Equal(t, len(d.Pool.Pool.SendResults), 1)
// 	if len(d.Pool.Pool.SendResults) == 0 {
// 		t.Fatalf("SendResults empty, would block")
// 	}
// 	sr = <-d.Pool.Pool.SendResults
// 	assert.Equal(t, sr.Connection, c)
// 	assert.Nil(t, sr.Error)
// 	_, ok = sr.Message.(*IntroductionMessage)
// 	assert.True(t, ok)
// 	// This connection should no longer be pending
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	// We should mark this as an outgoing connection since we solicited it
// 	assert.Equal(t, len(d.OutgoingConnections), 1)
// 	assert.NotNil(t, d.OutgoingConnections[addr])
// 	// We should be expecting its version
// 	assert.Equal(t, len(d.ExpectingIntroductions), 1)
// 	_, exists = d.ExpectingIntroductions[addr]
// 	assert.True(t, exists)
// 	// An introduction should have been sent
// 	assert.False(t, c.LastSent.IsZero())
// 	// d.ipCounts should be 1
// 	assert.Equal(t, d.ipCounts[addrIP], 1)
// 	// Cleanup
// 	c.Close()
// 	delete(d.ExpectingIntroductions, addr)
// 	delete(d.OutgoingConnections, addr)
// 	delete(d.ipCounts, addrIP)

// 	// Test a connection that is not connected by the time of processing
// 	c.LastSent = time.Time{}
// 	e = ConnectEvent{addr, true}
// 	delete(d.Pool.Pool.Addresses, addr)
// 	d.pendingConnections[addr] = p
// 	assert.NotPanics(t, func() { d.onConnect(e) })
// 	wait()
// 	assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
// 	// This connection should no longer be pending
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	// No message should have been sent
// 	assert.True(t, c.LastSent.IsZero())
// 	// We should not be expecting its version
// 	assert.Equal(t, len(d.ExpectingIntroductions), 0)
// 	// We should not have recorded it to ipCount
// 	assert.Equal(t, d.ipCounts[addrIP], 0)

// 	// Test a connection that is blacklisted
// 	e = ConnectEvent{addr, true}
// 	c = gnetConnection(addr)
// 	go d.Pool.Pool.ConnectionWriteLoop(c)
// 	d.Peers.Peers.AddBlacklistEntry(addr, time.Hour)
// 	d.pendingConnections[addr] = p
// 	d.Pool.Pool.Addresses[addr] = c
// 	d.Pool.Pool.Pool[c.Id] = c
// 	assert.NotPanics(t, func() { d.onConnect(e) })
// 	wait()
// 	// No introduction should have been sent
// 	assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
// 	// This connection should no longer be pending
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	// No message should have been sent
// 	assert.True(t, c.LastSent.IsZero())
// 	// We should not be expecting its version
// 	assert.Equal(t, len(d.ExpectingIntroductions), 0)
// 	// We should not have recorded its ipCount
// 	assert.Equal(t, d.ipCounts[addrIP], 0)
// 	// We should be looking to disconnect this client
// 	assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
// 	if len(d.Pool.Pool.DisconnectQueue) == 0 {
// 		t.Fatal("pool.Pool.DisconnectQueue is empty, would block")
// 	}
// 	de := <-d.Pool.Pool.DisconnectQueue
// 	assert.Equal(t, de.ConnId, 1)
// 	assert.Equal(t, de.Reason, DisconnectIsBlacklisted)
// 	// Cleanup
// 	c.Close()
// 	delete(d.Peers.Peers.Blacklist, addr)

// 	// Test a connection that has reached maxed ipCount
// 	e = ConnectEvent{addr, true}
// 	c = gnetConnection(addr)
// 	go d.Pool.Pool.ConnectionWriteLoop(c)
// 	d.ipCounts[addrIP] = d.Config.IPCountsMax
// 	d.pendingConnections[addr] = p
// 	d.Pool.Pool.Addresses[addr] = c
// 	d.Pool.Pool.Pool[c.Id] = c
// 	assert.NotPanics(t, func() { d.onConnect(e) })
// 	wait()
// 	// No introduction should have been sent
// 	assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
// 	// This connection should no longer be pending
// 	assert.Equal(t, len(d.pendingConnections), 0)
// 	// No message should have been sent
// 	assert.True(t, c.LastSent.IsZero())
// 	// We should not be expecting its version
// 	assert.Equal(t, len(d.ExpectingIntroductions), 0)
// 	// d.ipCounts should be unchanged
// 	assert.Equal(t, d.ipCounts[addrIP], d.Config.IPCountsMax)
// 	// We should be looking to disconnect this client
// 	assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
// 	if len(d.Pool.Pool.DisconnectQueue) == 0 {
// 		t.Fatal("pool.Pool.DisconnectQueue is empty, would block")
// 	}
// 	de = <-d.Pool.Pool.DisconnectQueue
// 	assert.Equal(t, de.ConnId, 1)
// 	assert.Equal(t, de.Reason, DisconnectIPLimitReached)
// 	// Cleanup
// 	c.Close()
// 	delete(d.ipCounts, addrIP)
// 	gnet.EraseMessages()
// 	shutdown(d)
// }

// func setupTestOnDisconnect(d *Daemon, c *gnet.Connection, mirror uint32) {
// 	d.OutgoingConnections[addr] = c
// 	d.ExpectingIntroductions[addr] = time.Now()
// 	d.mirrorConnections[mirror] = make(map[string]uint16)
// 	d.mirrorConnections[mirror][addrIP] = addrPort
// 	d.ConnectionMirrors[addr] = mirror
// }

// func TestOnDisconnect(t *testing.T) {
// 	gnet.EraseMessages()
// 	d := newDefaultDaemon()
// 	c := gnetConnection(addr)
// 	var mirror uint32 = 100

// 	// Not blacklistable
// 	reason := gnet.DisconnectWriteFailed
// 	setupTestOnDisconnect(d, c, mirror)
// 	assert.NotPanics(t, func() { d.onGnetDisconnect(c, reason) })
// 	// Should not be in blacklist
// 	assert.Equal(t, len(d.Peers.Peers.Blacklist), 0)
// 	// Should no longer be in OutgoingConnections
// 	assert.Equal(t, len(d.OutgoingConnections), 0)
// 	// Should no longer be in d.ExpectingIntroductions
// 	assert.Equal(t, len(d.ExpectingIntroductions), 0)
// 	// Should be removed from the mirror, and the mirror dict for this ip
// 	// should be removed
// 	assert.Equal(t, len(d.mirrorConnections), 0)
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)

// 	// Blacklistable
// 	reason = DisconnectIntroductionTimeout
// 	setupTestOnDisconnect(d, c, mirror)
// 	assert.NotPanics(t, func() { d.onGnetDisconnect(c, reason) })
// 	assert.Equal(t, len(d.Peers.Peers.Blacklist), 1)
// 	assert.NotNil(t, d.Peers.Peers.Blacklist[addr])
// 	// Should be in blacklist
// 	assert.Equal(t, len(d.Peers.Peers.Blacklist), 1)
// 	assert.NotNil(t, d.Peers.Peers.Blacklist[addr])
// 	// Should no longer be in OutgoingConnections
// 	assert.Equal(t, len(d.OutgoingConnections), 0)
// 	// Should no longer be in d.ExpectingIntroductions
// 	assert.Equal(t, len(d.ExpectingIntroductions), 0)
// 	// Should be removed from the mirror, and the mirror dict for this ip
// 	// should be removed
// 	assert.Equal(t, len(d.mirrorConnections), 0)
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	// Cleanup
// 	delete(d.Peers.Peers.Blacklist, addr)

// 	// d.mirrorConnections should retain a submap if there are other ports
// 	// inside
// 	reason = gnet.DisconnectWriteFailed
// 	setupTestOnDisconnect(d, c, mirror)
// 	d.mirrorConnections[mirror][strings.Split(addrb, ":")[0]] = addrPort
// 	assert.NotPanics(t, func() { d.onGnetDisconnect(c, reason) })
// 	// Should not be in blacklist
// 	assert.Equal(t, len(d.Peers.Peers.Blacklist), 0)
// 	// Should no longer be in OutgoingConnections
// 	assert.Equal(t, len(d.OutgoingConnections), 0)
// 	// Should no longer be in d.ExpectingIntroductions
// 	assert.Equal(t, len(d.ExpectingIntroductions), 0)
// 	// Should be removed from the mirror, and the mirror dict for this ip
// 	// should be removed
// 	assert.Equal(t, len(d.mirrorConnections), 1)
// 	assert.Equal(t, len(d.mirrorConnections[mirror]), 1)
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	shutdown(d)
// }

// func TestIPCountMaxed(t *testing.T) {
// 	d := newDefaultDaemon()
// 	assert.Equal(t, d.ipCounts[addrIP], 0)
// 	d.ipCounts[addrIP] = d.Config.IPCountsMax
// 	assert.True(t, d.ipCountMaxed(addr))
// 	d.ipCounts[addrIP] = 1
// 	assert.False(t, d.ipCountMaxed(addr))
// 	delete(d.ipCounts, addrIP)
// 	assert.False(t, d.ipCountMaxed(addr))
// 	// Invalid addr
// 	assert.True(t, d.ipCountMaxed("xcasca"))
// 	shutdown(d)
// }

// func TestRecordIPCount(t *testing.T) {
// 	d := newDefaultDaemon()
// 	assert.Equal(t, d.ipCounts[addrIP], 0)
// 	d.recordIPCount(addr)
// 	assert.Equal(t, d.ipCounts[addrIP], 1)
// 	d.recordIPCount(addr)
// 	assert.Equal(t, d.ipCounts[addrIP], 2)
// 	delete(d.ipCounts, addrIP)
// 	// Invalid addr
// 	d.recordIPCount("xcasca")
// 	assert.Equal(t, len(d.ipCounts), 0)
// 	shutdown(d)
// }

// func TestRemoveIPCount(t *testing.T) {
// 	d := newDefaultDaemon()
// 	assert.Equal(t, d.ipCounts[addrIP], 0)
// 	d.removeIPCount(addr)
// 	assert.Equal(t, d.ipCounts[addrIP], 0)
// 	d.ipCounts[addrIP] = 7
// 	d.removeIPCount(addr)
// 	assert.Equal(t, d.ipCounts[addrIP], 6)
// 	// Invalid addr
// 	d.ipCounts["xcasca"] = 1
// 	d.removeIPCount("xcasca")
// 	assert.Equal(t, d.ipCounts[addrIP], 6)
// 	assert.Equal(t, d.ipCounts["xcasca"], 1)
// 	assert.Equal(t, len(d.ipCounts), 2)
// 	delete(d.ipCounts, addrIP)
// 	shutdown(d)
// }

// func assertConnectMirrors(t *testing.T, d *Daemon) {
// 	m := d.ConnectionMirrors[addr]
// 	assert.Equal(t, m, d.Messages.Mirror)
// 	assert.NotEqual(t, m, 0)
// 	assert.Equal(t, len(d.ConnectionMirrors), 1)
// 	assert.Equal(t, len(d.mirrorConnections), 1)
// 	assert.Equal(t, len(d.mirrorConnections[d.Messages.Mirror]), 1)
// 	p, exists := d.mirrorConnections[d.Messages.Mirror][addrIP]
// 	assert.True(t, exists)
// 	assert.Equal(t, p, addrPort)
// 	shutdown(d)
// }

// func TestRecordConnectionMirror(t *testing.T) {
// 	d := newDefaultDaemon()
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	assert.Equal(t, len(d.mirrorConnections), 0)
// 	assert.Nil(t, d.recordConnectionMirror(addr, d.Messages.Mirror))
// 	assertConnectMirrors(t, d)

// 	// 2nd attempt should be a noop
// 	assert.Nil(t, d.recordConnectionMirror(addr, d.Messages.Mirror))
// 	assertConnectMirrors(t, d)

// 	assert.NotNil(t, d.recordConnectionMirror("xcasca", d.Messages.Mirror))
// 	assertConnectMirrors(t, d)

// 	delete(d.ConnectionMirrors, addr)
// 	delete(d.mirrorConnections, d.Messages.Mirror)
// 	shutdown(d)
// }

// func TestRemoveConnectionMirror(t *testing.T) {
// 	d := newDefaultDaemon()

// 	// No recorded addr should be noop
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	assert.Equal(t, len(d.mirrorConnections), 0)
// 	assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	assert.Equal(t, len(d.mirrorConnections), 0)

// 	// Invalid addr should be rejected
// 	d.ConnectionMirrors["xcasca"] = d.Messages.Mirror
// 	d.mirrorConnections[d.Messages.Mirror] = make(map[string]uint16)
// 	d.mirrorConnections[d.Messages.Mirror]["xcasca"] = 0
// 	assert.NotPanics(t, func() { d.removeConnectionMirror("xcasca") })
// 	assert.Equal(t, len(d.ConnectionMirrors), 1)
// 	assert.Equal(t, len(d.mirrorConnections), 1)
// 	assert.Equal(t, len(d.mirrorConnections[d.Messages.Mirror]), 1)
// 	delete(d.mirrorConnections, d.Messages.Mirror)
// 	delete(d.ConnectionMirrors, "xcasca")

// 	// With no connectionMirror recorded, we can't clean up the
// 	// d.mirrorConnections
// 	d.mirrorConnections[d.Messages.Mirror] = make(map[string]uint16)
// 	d.mirrorConnections[d.Messages.Mirror][addrIP] = addrPort
// 	assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	assert.Equal(t, len(d.mirrorConnections), 1)
// 	assert.Equal(t, len(d.mirrorConnections[d.Messages.Mirror]), 1)

// 	// Should clean up if all valid
// 	d.ConnectionMirrors[addr] = d.Messages.Mirror
// 	assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	assert.Equal(t, len(d.mirrorConnections), 0)

// 	// Cleaning up should leave d.mirrorConnections[addr] intact if multiple
// 	d.ConnectionMirrors[addr] = d.Messages.Mirror
// 	d.mirrorConnections[d.Messages.Mirror] = make(map[string]uint16)
// 	d.mirrorConnections[d.Messages.Mirror][addrIP] = addrPort
// 	d.mirrorConnections[d.Messages.Mirror][addrbIP] = addrbPort

// 	assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
// 	assert.Equal(t, len(d.ConnectionMirrors), 0)
// 	assert.Equal(t, len(d.mirrorConnections), 1)
// 	assert.Equal(t, d.mirrorConnections[d.Messages.Mirror][addrbIP], addrbPort)
// 	delete(d.mirrorConnections, d.Messages.Mirror)
// 	shutdown(d)
// }

// func TestGetMirrorPort(t *testing.T) {
// 	d := newDefaultDaemon()
// 	p, exists := d.getMirrorPort(addr, d.Messages.Mirror)
// 	assert.Equal(t, p, uint16(0))
// 	assert.False(t, exists)
// 	d.mirrorConnections[d.Messages.Mirror] = make(map[string]uint16)
// 	d.mirrorConnections[d.Messages.Mirror][addrIP] = addrPort
// 	p, exists = d.getMirrorPort(addr, d.Messages.Mirror)
// 	assert.Equal(t, p, addrPort)
// 	assert.True(t, exists)
// 	// Invalid addr
// 	p, exists = d.getMirrorPort("xcasca", d.Messages.Mirror)
// 	assert.Equal(t, p, uint16(0))
// 	assert.False(t, exists)
// 	delete(d.mirrorConnections, d.Messages.Mirror)
// 	shutdown(d)
// }

// //should not require visor
// /*
// func TestHandleMessageSendResult(t *testing.T) {
// 	d := newDefaultDaemon()
// 	defer shutdown(d)

// 	// Nothing happens: Message successfully sent and isnt a SendingTxnsMessage
// 	m := NewGetBlocksMessage(6)
// 	sr := gnet.SendResult{
// 		Message:    m,
// 		Connection: nil,
// 		Error:      nil,
// 	}
// 	assert.NotPanics(t, func() { d.handleMessageSendResult(sr) })

// 	// Add a txn for txn announce update testing

// 	vc, _ := setupVisor()
// 	v := NewVisor(vc)
// 	tx := addUnconfirmedTxn(v)
// 	assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 1)
// 	ut := v.Visor.Unconfirmed.Txns[tx.Hash()]
// 	assert.True(t, ut.Announced.IsZero())
// 	txns := coin.Transactions{tx.Txn}
// 	m2 := NewAnnounceTxnsMessage(txns.Hashes())

// 	// Logs a warning, and exits
// 	sr.Message = m2
// 	sr.Error = errors.New("Failed")
// 	sr.Connection = gnetConnection(addr)
// 	assert.NotPanics(t, func() { d.handleMessageSendResult(sr) })
// 	ut = v.Visor.Unconfirmed.Txns[tx.Hash()]
// 	assert.True(t, ut.Announced.IsZero())

// 	// Updates announcement
// 	sr.Error = nil
// 	sr.Message = m2
// 	d.Visor = v
// 	assert.NotPanics(t, func() {
// 		d.handleMessageSendResult(sr)
// 	})
// 	ut = v.Visor.Unconfirmed.Txns[tx.Hash()]
// 	assert.False(t, ut.Announced.IsZero())
// }
// */

// func TestIsLocalhost(t *testing.T) {
// 	assert.True(t, IsLocalhost("127.0.0.1"))
// 	assert.False(t, IsLocalhost(addrIP))
// }

// func TestLocalhostIP(t *testing.T) {
// 	ip, err := LocalhostIP()
// 	assert.Nil(t, err)
// 	assert.True(t, strings.HasPrefix(ip, "127"))
// }

// func TestSplitAddr(t *testing.T) {
// 	a, p, err := SplitAddr(addr)
// 	assert.Nil(t, err)
// 	assert.Equal(t, a, addrIP)
// 	assert.Equal(t, p, addrPort)
// 	a, p, err = SplitAddr(addrIP)
// 	assert.NotNil(t, err)
// }
