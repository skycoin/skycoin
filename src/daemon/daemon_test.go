package daemon

import (
    "errors"
    "github.com/nictuku/dht"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
    "github.com/stretchr/testify/assert"
    "strings"
    "testing"
    "time"
)

func newDefaultDaemon() *Daemon {
    cleanupPeers()
    c := NewConfig()
    c.Visor.Disabled = true
    return NewDaemon(c)
}

func TestGetListenPort(t *testing.T) {
    d := newDefaultDaemon()
    // No connectionMirror found
    assert.Equal(t, d.getListenPort(addr), uint16(0))
    // No mirrorConnection map exists
    d.connectionMirrors[addr] = uint32(4)
    assert.Panics(t, func() { d.getListenPort(addr) })
    // Everything is good
    m := make(map[string]uint16)
    d.mirrorConnections[uint32(4)] = m
    m[addrIP] = uint16(6667)
    assert.Equal(t, d.getListenPort(addr), uint16(6667))
    shutdown(d)
}

func TestStart(t *testing.T) {
    gnet.EraseMessages()
    defer cleanupPeers()
    d := newDefaultDaemon()
    quit := make(chan int)
    go d.Start(quit)
    assert.NotEqual(t, len(gnet.MessageIdMap), 0)
    assert.NotNil(t, d.Pool)
    assert.NotNil(t, d.Peers)
    assert.NotNil(t, d.DHT)
    assert.NotNil(t, d.Messages)
    assert.NotNil(t, d.RPC)
    quit <- 1
    wait()
    shutdown(d)
}

func TestShutdown(t *testing.T) {
    cleanupPeers()
    d := newDefaultDaemon()
    d.Peers.Peers.AddPeer(addr)
    d.Pool.Pool.DisconnectQueue <- gnet.DisconnectEvent{
        ConnId: 1,
        Reason: DisconnectOtherError,
    }
    assert.NotPanics(t, func() { d.Shutdown() })
    confirmPeersShutdown(t)
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 0)
    assert.Nil(t, d.DHT.DHT)
    cleanupPeers()
}

func setupDaemonLoop() (*Daemon, chan int) {
    d := newDefaultDaemon()
    quit := make(chan int)
    return d, quit
}

func TestDaemonLoopQuit(t *testing.T) {
    d, quit := setupDaemonLoop()
    done := false
    go func() {
        d.Start(quit)
        done = true
    }()
    wait()
    quit <- 1
    wait()
    assert.True(t, done)
    shutdown(d)
}

func TestDaemonLoopApiRequest(t *testing.T) {
    d, quit := setupDaemonLoop()
    go d.Start(quit)
    d.RPC.requests <- func() interface{} { return &Connection{Id: 7} }
    resp := <-d.RPC.responses
    assert.Equal(t, resp.(*Connection).Id, 7)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopOnConnectEvent(t *testing.T) {
    d, quit := setupDaemonLoop()
    go d.Start(quit)
    d.pendingConnections[addr] = pex.NewPeer(addr)
    d.onConnectEvent <- ConnectEvent{addr, false}
    wait()
    assert.Equal(t, len(d.pendingConnections), 0)
    assert.Nil(t, d.pendingConnections[addr])
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopConnectionErrors(t *testing.T) {
    d, quit := setupDaemonLoop()
    go d.Start(quit)
    d.pendingConnections[addr] = pex.NewPeer(addr)
    d.connectionErrors <- ConnectionError{addr, errors.New("failed")}
    wait()
    assert.Equal(t, len(d.pendingConnections), 0)
    assert.Nil(t, d.pendingConnections[addr])
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopDisconnectQueue(t *testing.T) {
    d, quit := setupDaemonLoop()
    go d.Start(quit)
    d.Pool.Pool.Pool[1] = gnetConnection(addr)
    e := gnet.DisconnectEvent{ConnId: 1, Reason: DisconnectIdle}
    d.Pool.Pool.DisconnectQueue <- e
    wait()
    assert.Equal(t, len(d.Pool.Pool.Pool), 0)
    quit <- 1
    wait()
    shutdown(d)
}

type DummyAsyncMessage struct {
    fn func()
}

func (self *DummyAsyncMessage) Process(d *Daemon) {
    self.fn()
}

func TestDaemonLoopMessageEvent(t *testing.T) {
    d, quit := setupDaemonLoop()
    go d.Start(quit)
    called := false
    m := &DummyAsyncMessage{fn: func() { called = true }}
    d.messageEvents <- MessageEvent{m, messageContext(addr)}
    wait()
    assert.True(t, called)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopDHTResults(t *testing.T) {
    d, quit := setupDaemonLoop()
    assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
    go d.Start(quit)
    m := make(map[dht.InfoHash][]string, 1)
    m[d.DHT.InfoHash] = []string{"abcdef"}
    d.DHT.DHT.PeersRequestResults <- m
    wait()
    assert.Equal(t, len(d.Peers.Peers.Peerlist), 1)
    assert.NotNil(t, d.Peers.Peers.Peerlist["97.98.99.100:25958"])
    quit <- 1
    wait()
    shutdown(d)
}

// TODO -- how to test tickers?
// TODO -- override rate to something very fast?

func TestDaemonLoopDHTBootstrapTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    d.DHT.Config.BootstrapRequestRate = time.Millisecond * 10
    go d.Start(quit)
    // Can't really test DHT internals, but we'll know if it crashes or not
    time.Sleep(time.Millisecond * 15)
    d.DHT.Config.PeerLimit = 0
    time.Sleep(time.Millisecond * 15)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopBlacklistTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    assert.Equal(t, len(d.Peers.Peers.Blacklist), 0)
    d.Peers.Peers.AddBlacklistEntry(addr, time.Millisecond)
    assert.Equal(t, len(d.Peers.Peers.Blacklist), 1)
    d.Peers.Config.UpdateBlacklistRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    assert.Equal(t, len(d.Peers.Peers.Blacklist), 0)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopCullInvalidTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    d.expectingIntroductions[addr] = time.Now().Add(-time.Hour)
    d.Config.CullInvalidRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    assert.Equal(t, len(d.expectingIntroductions), 0)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopRequestPeersTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    c := gnetConnection(addr)
    d.Pool.Pool.Pool[1] = c
    d.Pool.Pool.Addresses[c.Addr()] = c
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    d.Peers.Config.RequestRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopClearOldPeersTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    p := pex.NewPeer(addr)
    p.LastSeen = time.Unix(0, 0)
    d.Peers.Peers.Peerlist[addr] = p
    d.Peers.Config.CullRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopClearStaleConnectionsTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    c := gnetConnection(addr)
    c.LastReceived = time.Unix(0, 0)
    d.Pool.Pool.Pool[c.Id] = c
    d.Pool.Config.ClearStaleRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    assert.Equal(t, len(d.Pool.Pool.Pool), 0)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopPingCheckTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    c := gnetConnection(addr)
    c.LastSent = time.Unix(0, 0)
    d.Pool.Pool.Pool[c.Id] = c
    d.Pool.Config.IdleCheckRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopOutgoingConnectionsTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    d.Pool.Pool.Config.DialTimeout = 1 // nanosecond
    d.Peers.Peers.AddPeer(addr)
    d.Config.OutgoingRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    // Should have made a connection attempt, timed out, put an error
    // the queue, handled by d.Run, resulting in the peer being removed
    assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
    quit <- 1
    wait()
    shutdown(d)
}

func TestDaemonLoopMessageHandlingTicker(t *testing.T) {
    d, quit := setupDaemonLoop()
    d.Pool.Config.MessageHandlingRate = time.Millisecond * 10
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    // Can't test Pool internals from here
    quit <- 1
    wait()
    shutdown(d)
}

func TestRequestPeers(t *testing.T) {
    d := newDefaultDaemon()
    d.Peers.Config.Max = 1
    d.Peers.Peers.AddPeer(addr)
    // Nothing should happen if the peer list is full. It would have a nil
    // dereference of Pool if it continued further
    assert.NotPanics(t, func() { d.Peers.requestPeers(d.Pool) })

    c := gnetConnection(addr)
    d.Pool.Pool.Pool[1] = c
    d.Pool.Pool.Addresses[c.Addr()] = c
    assert.NotPanics(t, func() { d.Peers.requestPeers(d.Pool) })
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))

    // Failing send should not panic
    c.Conn = NewFailingConn(addr)
    c.LastSent = time.Unix(0, 0)
    assert.NotPanics(t, func() { d.Peers.requestPeers(d.Pool) })
    assert.Equal(t, c.LastSent, time.Unix(0, 0))

    shutdown(d)
}

func TestClearStaleConnections(t *testing.T) {
    dm := newDefaultDaemon()
    c := gnetConnection(addr)
    d := gnetConnection(addrb)
    c.LastReceived = time.Unix(0, 0)
    d.LastReceived = time.Now()
    dm.Pool.Pool.Pool[1] = c
    dm.Pool.Pool.Pool[2] = d
    assert.NotPanics(t, dm.Pool.clearStaleConnections)
    assert.Equal(t, len(dm.Pool.Pool.DisconnectQueue), 1)
    if len(dm.Pool.Pool.DisconnectQueue) == 0 {
        t.Fatalf("Empty DisconnectQueue, would block")
    }
    de := <-dm.Pool.Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 1)
    assert.Equal(t, de.Reason, DisconnectIdle)
    shutdown(dm)
}

func TestSendPings(t *testing.T) {
    d := newDefaultDaemon()
    c := gnetConnection(addr)
    d.Pool.Pool.Pool[1] = c
    assert.NotPanics(t, d.Pool.sendPings)
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    lastSent := c.LastSent
    assert.NotPanics(t, d.Pool.sendPings)
    assert.Equal(t, c.LastSent, lastSent)

    // Failing write should not panic
    c.Conn = NewFailingConn(addr)
    c.LastSent = time.Unix(0, 0)
    assert.NotPanics(t, d.Pool.sendPings)
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    shutdown(d)
}

func TestConnectToRandomPeer(t *testing.T) {
    d := newDefaultDaemon()
    d.Pool.Pool.Config.DialTimeout = 1 // nanosecond
    // Valid attempt to connect
    d.Peers.Peers.AddPeer(addr)
    assert.NotPanics(t, d.connectToRandomPeer)
    wait()
    assert.Equal(t, len(d.pendingConnections), 1)
    assert.Equal(t, len(d.connectionErrors), 1)
    if len(d.connectionErrors) == 0 {
        t.Fatalf("connectionErrors empty, would block")
    }
    ce := <-d.connectionErrors
    assert.Equal(t, ce.Addr, addr)
    assert.NotNil(t, ce.Error)
    delete(d.pendingConnections, addr)

    // Two peers, one successful connect attempt and one skipped
    d.Peers.Peers.AddPeer(addrb)
    assert.NotPanics(t, d.connectToRandomPeer)
    wait()
    assert.Equal(t, len(d.pendingConnections), 1)
    assert.Equal(t, len(d.connectionErrors), 1)
    if len(d.connectionErrors) == 0 {
        t.Fatalf("connectionErrors empty, would block")
    }
    ce = <-d.connectionErrors
    assert.True(t, (ce.Addr == addr) || (ce.Addr == addrb))
    assert.NotNil(t, ce.Error)
    delete(d.pendingConnections, addr)
    delete(d.pendingConnections, addrb)
    delete(d.Peers.Peers.Peerlist, addrb)

    // Already connected, skip
    d.Peers.Peers.AddPeer(addr)
    d.Pool.Pool.Addresses[addr] = gnetConnection(addr)
    assert.NotPanics(t, d.connectToRandomPeer)
    assert.Equal(t, len(d.pendingConnections), 0)
    assert.Equal(t, len(d.connectionErrors), 0)
    delete(d.Pool.Pool.Addresses, addr)

    // Pending connection, skip
    d.pendingConnections[addr] = pex.NewPeer(addr)
    assert.NotPanics(t, d.connectToRandomPeer)
    assert.Equal(t, len(d.pendingConnections), 1)
    assert.Equal(t, len(d.connectionErrors), 0)
    delete(d.pendingConnections, addr)

    // Already connected to this base IP at least once, skip
    d.ipCounts[addrIP] = 1
    assert.NotPanics(t, d.connectToRandomPeer)
    assert.Equal(t, len(d.ipCounts), 1)
    assert.Equal(t, d.ipCounts[addrIP], 1)
    assert.Equal(t, len(d.pendingConnections), 0)
    assert.Equal(t, len(d.connectionErrors), 0)
    delete(d.ipCounts, addrIP)

    shutdown(d)
    cleanupPeers()
}

func TestHandleConnectionError(t *testing.T) {
    d := newDefaultDaemon()
    p, _ := d.Peers.Peers.AddPeer(addr)
    d.pendingConnections[addr] = p
    assert.NotPanics(t, func() {
        d.handleConnectionError(ConnectionError{addr, nil})
    })
    assert.Equal(t, len(d.pendingConnections), 0)
    assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
    p, _ = d.Peers.Peers.AddPeer(addr)
    d.pendingConnections[addr] = p
    assert.NotPanics(t, func() {
        d.handleConnectionError(ConnectionError{addr, errors.New("bad")})
    })
    assert.Equal(t, len(d.pendingConnections), 0)
    assert.Equal(t, len(d.Peers.Peers.Peerlist), 0)
    shutdown(d)
}

func TestCullInvalidConnections(t *testing.T) {
    d := newDefaultDaemon()
    // Is fine
    d.expectingIntroductions[addr] = time.Now()
    // Is expired
    d.expectingIntroductions[addrb] = time.Unix(0, 0)
    // Is not in pool
    d.expectingIntroductions[addrc] = time.Unix(0, 0)
    d.Peers.Peers.AddPeer(addr)
    d.Peers.Peers.AddPeer(addrb)
    d.Peers.Peers.AddPeer(addrc)
    d.Pool.Pool.Addresses[addr] = gnetConnection(addr)
    d.Pool.Pool.Addresses[addrb] = gnetConnection(addrb)
    d.Pool.Pool.Addresses[addrb].Id = 2
    d.Pool.Pool.Pool[1] = d.Pool.Pool.Addresses[addr]
    d.Pool.Pool.Pool[2] = d.Pool.Pool.Addresses[addrb]

    assert.NotPanics(t, d.cullInvalidConnections)

    assert.Equal(t, len(d.expectingIntroductions), 1)
    assert.Equal(t, len(d.Peers.Peers.Peerlist), 2)
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
    if len(d.Pool.Pool.DisconnectQueue) == 0 {
        t.Fatal("pool.Pool.DisconnectQueue not empty, would block")
    }
    de := <-d.Pool.Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 2)
    assert.Equal(t, de.Reason, DisconnectIntroductionTimeout)
    shutdown(d)
}

func TestRecordMessageEventValid(t *testing.T) {
    d := newDefaultDaemon()
    // Valid message, not expecting Introduction
    assert.Equal(t, len(d.messageEvents), 0)
    delete(d.expectingIntroductions, addr)
    m := &PingMessage{}
    m.c = messageContext(addr)
    err := d.recordMessageEvent(m, m.c)
    assert.Nil(t, err)
    assert.Equal(t, len(d.messageEvents), 1)
    if len(d.messageEvents) == 0 {
        t.Fatal("d.messageEvents empty, would block")
    }
    me := <-d.messageEvents
    _, ok := me.Message.(*PingMessage)
    assert.True(t, ok)
    shutdown(d)
}

func TestRecordMessageEventIsIntroduction(t *testing.T) {
    // Needs Introduction and thats what it has received
    d := newDefaultDaemon()
    d.expectingIntroductions[addr] = time.Now().UTC()
    assert.Equal(t, len(d.messageEvents), 0)
    m := NewIntroductionMessage(d.Messages.Mirror, d.Config.Version,
        d.Pool.Pool.Config.Port)
    m.c = messageContext(addr)
    err := d.recordMessageEvent(m, m.c)
    assert.Nil(t, err)
    assert.Equal(t, len(d.messageEvents), 1)
    if len(d.messageEvents) == 0 {
        t.Fatal("d.messageEvents empty, would block")
    }
    me := <-d.messageEvents
    _, ok := me.Message.(*IntroductionMessage)
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 0)
    assert.True(t, ok)
    delete(d.expectingIntroductions, addr)
    shutdown(d)
}

func TestRecordMessageEventNeedsIntroduction(t *testing.T) {
    // Needs Introduction but didn't get it first
    d := newDefaultDaemon()
    m := &PingMessage{}
    m.c = messageContext(addr)
    d.Pool.Pool.Addresses[addr] = m.c.Conn
    d.Pool.Pool.Pool[m.c.Conn.Id] = m.c.Conn
    assert.Equal(t, len(d.messageEvents), 0)
    d.expectingIntroductions[addr] = time.Now().UTC()
    d.processMessageEvent(MessageEvent{m, m.c})
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
    if len(d.Pool.Pool.DisconnectQueue) == 0 {
        t.Fatal("DisconnectQueue empty, would block")
    }
    de := <-d.Pool.Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, m.c.Conn.Id)
    assert.Equal(t, de.Reason, DisconnectNoIntroduction)
    delete(d.expectingIntroductions, addr)
    shutdown(d)
}

func TestOnConnect(t *testing.T) {
    d := newDefaultDaemon()

    // Test a valid connection, unsolicited
    e := ConnectEvent{addr, false}
    p, _ := d.Peers.Peers.AddPeer(addr)
    c := gnetConnection(addr)
    d.pendingConnections[addr] = p
    d.Pool.Pool.Pool[1] = c
    d.Pool.Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { d.onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(d.pendingConnections), 0)
    // This is not an outgoing connection, we did not solicit it
    assert.Equal(t, len(d.outgoingConnections), 0)
    // We should be expecting its version
    assert.Equal(t, len(d.expectingIntroductions), 1)
    _, exists := d.expectingIntroductions[addr]
    assert.True(t, exists)
    // An introduction should have been sent
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    // d.ipCounts should be 1
    assert.Equal(t, d.ipCounts[addrIP], 1)
    // Cleanup
    delete(d.ipCounts, addrIP)
    delete(d.expectingIntroductions, addr)

    // Test a valid connection, solicited
    e = ConnectEvent{addr, true}
    c = gnetConnection(addr)
    d.pendingConnections[addr] = p
    d.Pool.Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { d.onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(d.pendingConnections), 0)
    // We should mark this as an outgoing connection since we solicited it
    assert.Equal(t, len(d.outgoingConnections), 1)
    assert.NotNil(t, d.outgoingConnections[addr])
    // We should be expecting its version
    assert.Equal(t, len(d.expectingIntroductions), 1)
    _, exists = d.expectingIntroductions[addr]
    assert.True(t, exists)
    // An introduction should have been sent
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    // d.ipCounts should be 1
    assert.Equal(t, d.ipCounts[addrIP], 1)
    // Cleanup
    delete(d.expectingIntroductions, addr)
    delete(d.outgoingConnections, addr)
    delete(d.ipCounts, addrIP)

    // Test a valid connection, but failing to send a message
    e = ConnectEvent{addr, true}
    c = gnetConnection(addr)
    c.Conn = NewFailingConn(addr)
    d.pendingConnections[addr] = p
    d.Pool.Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { d.onConnect(e) })
    wait()
    // This connection should no longer be pending
    assert.Equal(t, len(d.pendingConnections), 0)
    // We should mark this as an outgoing connection since we solicited it
    assert.Equal(t, len(d.outgoingConnections), 1)
    assert.NotNil(t, d.outgoingConnections[addr])
    // We should be expecting its version
    assert.Equal(t, len(d.expectingIntroductions), 1)
    _, exists = d.expectingIntroductions[addr]
    assert.True(t, exists)
    // An introduction should not have been sent, it failed
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    // We should be looking to disconnect this client
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
    if len(d.Pool.Pool.DisconnectQueue) == 0 {
        t.Fatal("Pool.DisconnectQueue is empty, would block")
    }
    de := <-d.Pool.Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 1)
    assert.Equal(t, de.Reason, DisconnectFailedSend)
    // d.ipCounts should be 1, since we haven't processed the disconnect yet
    assert.Equal(t, d.ipCounts[addrIP], 1)
    // Cleanup
    delete(d.expectingIntroductions, addr)
    delete(d.outgoingConnections, addr)
    delete(d.ipCounts, addrIP)

    // Test a connection that is not connected by the time of processing
    e = ConnectEvent{addr, true}
    delete(d.Pool.Pool.Addresses, addr)
    d.pendingConnections[addr] = p
    assert.NotPanics(t, func() { d.onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(d.pendingConnections), 0)
    // No message should have been sent
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    // We should not be expecting its version
    assert.Equal(t, len(d.expectingIntroductions), 0)
    // We should not have recorded it to ipCount
    assert.Equal(t, d.ipCounts[addrIP], 0)

    // Test a connection that is blacklisted
    e = ConnectEvent{addr, true}
    c = gnetConnection(addr)
    d.Peers.Peers.AddBlacklistEntry(addr, time.Hour)
    d.pendingConnections[addr] = p
    d.Pool.Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { d.onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(d.pendingConnections), 0)
    // No message should have been sent
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    // We should not be expecting its version
    assert.Equal(t, len(d.expectingIntroductions), 0)
    // We should not have recorded its ipCount
    assert.Equal(t, d.ipCounts[addrIP], 0)
    // We should be looking to disconnect this client
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
    if len(d.Pool.Pool.DisconnectQueue) == 0 {
        t.Fatal("pool.Pool.DisconnectQueue is empty, would block")
    }
    de = <-d.Pool.Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 1)
    assert.Equal(t, de.Reason, DisconnectIsBlacklisted)
    // Cleanup
    delete(d.Peers.Peers.Blacklist, addr)

    // Test a connection that has reached maxed ipCount
    e = ConnectEvent{addr, true}
    c = gnetConnection(addr)
    d.ipCounts[addrIP] = d.Config.IPCountsMax
    d.pendingConnections[addr] = p
    d.Pool.Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { d.onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(d.pendingConnections), 0)
    // No message should have been sent
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    // We should not be expecting its version
    assert.Equal(t, len(d.expectingIntroductions), 0)
    // d.ipCounts should be unchanged
    assert.Equal(t, d.ipCounts[addrIP], d.Config.IPCountsMax)
    // We should be looking to disconnect this client
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 1)
    if len(d.Pool.Pool.DisconnectQueue) == 0 {
        t.Fatal("pool.Pool.DisconnectQueue is empty, would block")
    }
    de = <-d.Pool.Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 1)
    assert.Equal(t, de.Reason, DisconnectIPLimitReached)
    // Cleanup
    delete(d.ipCounts, addrIP)
    gnet.EraseMessages()
    shutdown(d)
}

func setupTestOnDisconnect(d *Daemon, c *gnet.Connection, mirror uint32) {
    d.outgoingConnections[addr] = c
    d.expectingIntroductions[addr] = time.Now()
    d.mirrorConnections[mirror] = make(map[string]uint16)
    d.mirrorConnections[mirror][addrIP] = addrPort
    d.connectionMirrors[addr] = mirror
}

func TestOnDisconnect(t *testing.T) {
    d := newDefaultDaemon()
    c := gnetConnection(addr)
    var mirror uint32 = 100

    // Not blacklistable
    reason := DisconnectFailedSend
    setupTestOnDisconnect(d, c, mirror)
    assert.NotPanics(t, func() { d.onGnetDisconnect(c, reason) })
    // Should not be in blacklist
    assert.Equal(t, len(d.Peers.Peers.Blacklist), 0)
    // Should no longer be in outgoingConnections
    assert.Equal(t, len(d.outgoingConnections), 0)
    // Should no longer be in d.expectingIntroductions
    assert.Equal(t, len(d.expectingIntroductions), 0)
    // Should be removed from the mirror, and the mirror dict for this ip
    // should be removed
    assert.Equal(t, len(d.mirrorConnections), 0)
    assert.Equal(t, len(d.connectionMirrors), 0)

    // Blacklistable
    reason = DisconnectIntroductionTimeout
    setupTestOnDisconnect(d, c, mirror)
    assert.NotPanics(t, func() { d.onGnetDisconnect(c, reason) })
    assert.Equal(t, len(d.Peers.Peers.Blacklist), 1)
    assert.NotNil(t, d.Peers.Peers.Blacklist[addr])
    // Should be in blacklist
    assert.Equal(t, len(d.Peers.Peers.Blacklist), 1)
    assert.NotNil(t, d.Peers.Peers.Blacklist[addr])
    // Should no longer be in outgoingConnections
    assert.Equal(t, len(d.outgoingConnections), 0)
    // Should no longer be in d.expectingIntroductions
    assert.Equal(t, len(d.expectingIntroductions), 0)
    // Should be removed from the mirror, and the mirror dict for this ip
    // should be removed
    assert.Equal(t, len(d.mirrorConnections), 0)
    assert.Equal(t, len(d.connectionMirrors), 0)
    // Cleanup
    delete(d.Peers.Peers.Blacklist, addr)

    // d.mirrorConnections should retain a submap if there are other ports
    // inside
    reason = DisconnectFailedSend
    setupTestOnDisconnect(d, c, mirror)
    d.mirrorConnections[mirror][strings.Split(addrb, ":")[0]] = addrPort
    assert.NotPanics(t, func() { d.onGnetDisconnect(c, reason) })
    // Should not be in blacklist
    assert.Equal(t, len(d.Peers.Peers.Blacklist), 0)
    // Should no longer be in outgoingConnections
    assert.Equal(t, len(d.outgoingConnections), 0)
    // Should no longer be in d.expectingIntroductions
    assert.Equal(t, len(d.expectingIntroductions), 0)
    // Should be removed from the mirror, and the mirror dict for this ip
    // should be removed
    assert.Equal(t, len(d.mirrorConnections), 1)
    assert.Equal(t, len(d.mirrorConnections[mirror]), 1)
    assert.Equal(t, len(d.connectionMirrors), 0)
    shutdown(d)
}

func TestIPCountMaxed(t *testing.T) {
    d := newDefaultDaemon()
    assert.Equal(t, d.ipCounts[addrIP], 0)
    d.ipCounts[addrIP] = d.Config.IPCountsMax
    assert.True(t, d.ipCountMaxed(addrIP))
    d.ipCounts[addrIP] = 1
    assert.False(t, d.ipCountMaxed(addrIP))
    delete(d.ipCounts, addrIP)
    assert.False(t, d.ipCountMaxed(addrIP))
    shutdown(d)
}

func TestRecordIPCount(t *testing.T) {
    d := newDefaultDaemon()
    assert.Equal(t, d.ipCounts[addrIP], 0)
    d.recordIPCount(addr)
    assert.Equal(t, d.ipCounts[addrIP], 1)
    d.recordIPCount(addr)
    assert.Equal(t, d.ipCounts[addrIP], 2)
    delete(d.ipCounts, addrIP)
    shutdown(d)
}

func TestRemoveIPCount(t *testing.T) {
    d := newDefaultDaemon()
    assert.Equal(t, d.ipCounts[addrIP], 0)
    d.removeIPCount(addr)
    assert.Equal(t, d.ipCounts[addrIP], 0)
    d.ipCounts[addrIP] = 7
    d.removeIPCount(addr)
    assert.Equal(t, d.ipCounts[addrIP], 6)
    delete(d.ipCounts, addrIP)
    shutdown(d)
}

func assertConnectMirrors(t *testing.T, d *Daemon) {
    m := d.connectionMirrors[addr]
    assert.Equal(t, m, d.Messages.Mirror)
    assert.NotEqual(t, m, 0)
    assert.Equal(t, len(d.connectionMirrors), 1)
    assert.Equal(t, len(d.mirrorConnections), 1)
    assert.Equal(t, len(d.mirrorConnections[d.Messages.Mirror]), 1)
    p, exists := d.mirrorConnections[d.Messages.Mirror][addrIP]
    assert.True(t, exists)
    assert.Equal(t, p, addrPort)
    shutdown(d)
}

func TestRecordConnectionMirror(t *testing.T) {
    d := newDefaultDaemon()
    assert.Equal(t, len(d.connectionMirrors), 0)
    assert.Equal(t, len(d.mirrorConnections), 0)
    d.recordConnectionMirror(addr, d.Messages.Mirror)
    assertConnectMirrors(t, d)

    // 2nd attempt should be a noop
    d.recordConnectionMirror(addr, d.Messages.Mirror)
    assertConnectMirrors(t, d)

    delete(d.connectionMirrors, addr)
    delete(d.mirrorConnections, d.Messages.Mirror)
    shutdown(d)
}

func TestRemoveConnectionMirror(t *testing.T) {
    d := newDefaultDaemon()
    // No recorded addr should be noop
    assert.Equal(t, len(d.connectionMirrors), 0)
    assert.Equal(t, len(d.mirrorConnections), 0)
    assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
    assert.Equal(t, len(d.connectionMirrors), 0)
    assert.Equal(t, len(d.mirrorConnections), 0)

    // With no connectionMirror recorded, we can't clean up the
    // d.mirrorConnections
    d.mirrorConnections[d.Messages.Mirror] = make(map[string]uint16)
    d.mirrorConnections[d.Messages.Mirror][addrIP] = addrPort
    assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
    assert.Equal(t, len(d.connectionMirrors), 0)
    assert.Equal(t, len(d.mirrorConnections), 1)
    assert.Equal(t, len(d.mirrorConnections[d.Messages.Mirror]), 1)

    // Should clean up if all valid
    d.connectionMirrors[addr] = d.Messages.Mirror
    assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
    assert.Equal(t, len(d.connectionMirrors), 0)
    assert.Equal(t, len(d.mirrorConnections), 0)

    // Cleaning up should leave d.mirrorConnections[addr] intact if multiple
    d.connectionMirrors[addr] = d.Messages.Mirror
    d.mirrorConnections[d.Messages.Mirror] = make(map[string]uint16)
    d.mirrorConnections[d.Messages.Mirror][addrIP] = addrPort
    d.mirrorConnections[d.Messages.Mirror][addrbIP] = addrbPort

    assert.NotPanics(t, func() { d.removeConnectionMirror(addr) })
    assert.Equal(t, len(d.connectionMirrors), 0)
    assert.Equal(t, len(d.mirrorConnections), 1)
    assert.Equal(t, d.mirrorConnections[d.Messages.Mirror][addrbIP], addrbPort)
    delete(d.mirrorConnections, d.Messages.Mirror)
    shutdown(d)
}

func TestGetMirrorPort(t *testing.T) {
    d := newDefaultDaemon()
    p, exists := d.getMirrorPort(addr, d.Messages.Mirror)
    assert.Equal(t, p, uint16(0))
    assert.False(t, exists)
    d.mirrorConnections[d.Messages.Mirror] = make(map[string]uint16)
    d.mirrorConnections[d.Messages.Mirror][addrIP] = addrPort
    p, exists = d.getMirrorPort(addr, d.Messages.Mirror)
    assert.Equal(t, p, addrPort)
    assert.True(t, exists)
    delete(d.mirrorConnections, d.Messages.Mirror)
    shutdown(d)
}

/* Helpers */

func shutdown(d *Daemon) {
    d.Shutdown()
    wait()
    cleanupPeers()
}
