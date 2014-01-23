package daemon

import (
    "errors"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
    "github.com/stretchr/testify/assert"
    "strings"
    "testing"
    "time"
)

func TestGetListenPort(t *testing.T) {
    // No connectionMirror found
    assert.Equal(t, getListenPort(addr), uint16(0))
    // No mirrorConnection map exists
    connectionMirrors[addr] = uint32(4)
    assert.Panics(t, func() { getListenPort(addr) })
    // Everything is good
    m := make(map[string]uint16)
    mirrorConnections[uint32(4)] = m
    m[addrIP] = uint16(6667)
    assert.Equal(t, getListenPort(addr), uint16(6667))

    // cleanup
    delete(mirrorConnections, uint32(4))
    delete(connectionMirrors, addr)
}

func TestInit(t *testing.T) {
    // TODO --
    // Need to be able to stop the DHT
}

func TestShutdown(t *testing.T) {
    pool := SetupPoolShutdown(t)
    SetupPeersShutdown(t)
    assert.NotPanics(t, func() { Shutdown("./") })
    ConfirmPeersShutdown(t)
    ConfirmPoolShutdown(t, pool)
}

func TestDaemonLoop(t *testing.T) {
    // TODO
    // Make sure every possible event is handled??
}

func TestRequestPeers(t *testing.T) {
    gnet.EraseMessages()
    RegisterMessages()
    Pool = nil
    _m := maxPeers
    maxPeers = 1
    InitPeers("./")
    Peers.AddPeer(addr)
    // Nothing should happen if the peer list is full. It would have a nil
    // dereference of Pool if it continued further
    assert.NotPanics(t, requestPeers)

    Peers = nil
    InitPeers("./")

    Pool = gnet.NewConnectionPool(port)
    c := gnetConnection(addr)
    Pool.Pool[1] = c
    assert.NotPanics(t, requestPeers)
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))

    // Failing send should not panic
    c.Conn = NewFailingConn(addr)
    c.LastSent = time.Unix(0, 0)
    assert.NotPanics(t, requestPeers)
    assert.Equal(t, c.LastSent, time.Unix(0, 0))

    gnet.EraseMessages()
    ShutdownPool()
    maxPeers = _m
}

func TestClearStaleConnections(t *testing.T) {
    Pool = gnet.NewConnectionPool(port)
    c := gnetConnection(addr)
    d := gnetConnection(addrb)
    c.LastReceived = time.Unix(0, 0)
    d.LastReceived = time.Now()
    Pool.Pool[1] = c
    Pool.Pool[2] = d
    assert.NotPanics(t, clearStaleConnections)
    assert.Equal(t, len(Pool.DisconnectQueue), 1)
    if len(Pool.DisconnectQueue) == 0 {
        t.Fatalf("Empty DisconnectQueue, would block")
    }
    de := <-Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 1)
    assert.Equal(t, de.Reason, DisconnectIdle)
    ShutdownPool()
}

func TestSendPings(t *testing.T) {
    RegisterMessages()
    Pool = gnet.NewConnectionPool(port)
    c := gnetConnection(addr)
    Pool.Pool[1] = c
    assert.NotPanics(t, sendPings)
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    lastSent := c.LastSent
    assert.NotPanics(t, sendPings)
    assert.Equal(t, c.LastSent, lastSent)

    // Failing write should not panic
    c.Conn = NewFailingConn(addr)
    c.LastSent = time.Unix(0, 0)
    assert.NotPanics(t, sendPings)
    assert.Equal(t, c.LastSent, time.Unix(0, 0))

    ShutdownPool()
    gnet.EraseMessages()
}

func TestConnectToRandomPeer(t *testing.T) {
    Pool = gnet.NewConnectionPool(port)
    Peers = pex.NewPex(maxPeers)

    dt := gnet.DialTimeout
    gnet.DialTimeout = 1 // nanosecond
    // Valid attempt to connect
    Peers.AddPeer(addr)
    assert.NotPanics(t, connectToRandomPeer)
    wait()
    assert.Equal(t, len(pendingConnections), 1)
    assert.Equal(t, len(connectionErrors), 1)
    if len(connectionErrors) == 0 {
        t.Fatalf("connectionErrors empty, would block")
    }
    ce := <-connectionErrors
    assert.Equal(t, ce.Addr, addr)
    assert.NotNil(t, ce.Error)
    delete(pendingConnections, addr)

    // Two peers, one successful connect attempt and one skipped
    Peers.AddPeer(addrb)
    assert.NotPanics(t, connectToRandomPeer)
    wait()
    assert.Equal(t, len(pendingConnections), 1)
    assert.Equal(t, len(connectionErrors), 1)
    if len(connectionErrors) == 0 {
        t.Fatalf("connectionErrors empty, would block")
    }
    ce = <-connectionErrors
    assert.True(t, (ce.Addr == addr) || (ce.Addr == addrb))
    assert.NotNil(t, ce.Error)
    delete(pendingConnections, addr)
    delete(pendingConnections, addrb)
    delete(Peers.Peerlist, addrb)

    // Already connected, skip
    Peers.AddPeer(addr)
    Pool.Addresses[addr] = gnetConnection(addr)
    assert.NotPanics(t, connectToRandomPeer)
    assert.Equal(t, len(pendingConnections), 0)
    assert.Equal(t, len(connectionErrors), 0)
    delete(Pool.Addresses, addr)

    // Pending connection, skip
    pendingConnections[addr] = pex.NewPeer(addr)
    assert.NotPanics(t, connectToRandomPeer)
    assert.Equal(t, len(pendingConnections), 1)
    assert.Equal(t, len(connectionErrors), 0)
    gnet.DialTimeout = dt

    resetPeers()
    ShutdownPool()
}

func TestHandleConnectionError(t *testing.T) {
    Peers = pex.NewPex(maxPeers)
    p, _ := Peers.AddPeer(addr)
    pendingConnections[addr] = p
    assert.NotPanics(t, func() {
        handleConnectionError(ConnectionError{addr, nil})
    })
    assert.Equal(t, len(pendingConnections), 0)
    assert.Equal(t, len(Peers.Peerlist), 0)
    p, _ = Peers.AddPeer(addr)
    pendingConnections[addr] = p
    assert.NotPanics(t, func() {
        handleConnectionError(ConnectionError{addr, errors.New("bad")})
    })
    assert.Equal(t, len(pendingConnections), 0)
    assert.Equal(t, len(Peers.Peerlist), 0)
    resetPeers()
}

func TestCullInvalidConnections(t *testing.T) {
    Peers = pex.NewPex(maxPeers)
    Pool = gnet.NewConnectionPool(port)
    // Is fine
    expectingIntroductions[addr] = time.Now()
    // Is expired
    expectingIntroductions[addrb] = time.Unix(0, 0)
    // Is not in pool
    expectingIntroductions[addrc] = time.Unix(0, 0)
    Peers.AddPeer(addr)
    Peers.AddPeer(addrb)
    Peers.AddPeer(addrc)
    Pool.Addresses[addr] = gnetConnection(addr)
    Pool.Addresses[addrb] = gnetConnection(addrb)
    Pool.Addresses[addrb].Id = 2
    Pool.Pool[1] = Pool.Addresses[addr]
    Pool.Pool[2] = Pool.Addresses[addrb]

    assert.NotPanics(t, cullInvalidConnections)

    assert.Equal(t, len(expectingIntroductions), 1)
    assert.Equal(t, len(Peers.Peerlist), 2)
    assert.Equal(t, len(Pool.DisconnectQueue), 1)
    if len(Pool.DisconnectQueue) == 0 {
        t.Fatal("Pool.DisconnectQueue not empty, would block")
    }
    de := <-Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 2)
    assert.Equal(t, de.Reason, DisconnectVersionTimeout)

    resetPeers()
    ShutdownPool()
}

func TestOnConnect(t *testing.T) {
    RegisterMessages()
    Peers = pex.NewPex(maxPeers)
    Pool = gnet.NewConnectionPool(port)

    // Test a valid connection, unsolicited
    e := ConnectEvent{addr, false}
    p, _ := Peers.AddPeer(addr)
    c := gnetConnection(addr)
    pendingConnections[addr] = p
    Pool.Pool[1] = c
    Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(pendingConnections), 0)
    // This is not an outgoing connection, we did not solicit it
    assert.Equal(t, len(outgoingConnections), 0)
    // We should be expecting its version
    assert.Equal(t, len(expectingIntroductions), 1)
    _, exists := expectingIntroductions[addr]
    assert.True(t, exists)
    // An introduction should have been sent
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    // Cleanup
    delete(expectingIntroductions, addr)

    // Test a valid connection, solicited
    e = ConnectEvent{addr, true}
    c = gnetConnection(addr)
    pendingConnections[addr] = p
    Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(pendingConnections), 0)
    // We should mark this as an outgoing connection since we solicited it
    assert.Equal(t, len(outgoingConnections), 1)
    assert.NotNil(t, outgoingConnections[addr])
    // We should be expecting its version
    assert.Equal(t, len(expectingIntroductions), 1)
    _, exists = expectingIntroductions[addr]
    assert.True(t, exists)
    // An introduction should have been sent
    assert.NotEqual(t, c.LastSent, time.Unix(0, 0))
    // Cleanup
    delete(expectingIntroductions, addr)
    delete(outgoingConnections, addr)

    // Test a valid connection, but failing to send a message
    e = ConnectEvent{addr, true}
    c = gnetConnection(addr)
    c.Conn = NewFailingConn(addr)
    pendingConnections[addr] = p
    Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { onConnect(e) })
    wait()
    // This connection should no longer be pending
    assert.Equal(t, len(pendingConnections), 0)
    // We should mark this as an outgoing connection since we solicited it
    assert.Equal(t, len(outgoingConnections), 1)
    assert.NotNil(t, outgoingConnections[addr])
    // We should be expecting its version
    assert.Equal(t, len(expectingIntroductions), 1)
    _, exists = expectingIntroductions[addr]
    assert.True(t, exists)
    // An introduction should not have been sent, it failed
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    // We should be looking to disconnect this client
    assert.Equal(t, len(Pool.DisconnectQueue), 1)
    if len(Pool.DisconnectQueue) == 0 {
        t.Fatal("Pool.DisconnectQueue is empty, would block")
    }
    de := <-Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 1)
    assert.Equal(t, de.Reason, DisconnectFailedSend)
    // Cleanup
    delete(expectingIntroductions, addr)
    delete(outgoingConnections, addr)

    // Test a connection that is not connected by the time of processing
    e = ConnectEvent{addr, true}
    delete(Pool.Addresses, addr)
    pendingConnections[addr] = p
    assert.NotPanics(t, func() { onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(pendingConnections), 0)
    // No message should have been sent
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    // We should not be expecting its version
    assert.Equal(t, len(expectingIntroductions), 0)

    // Test a connection that is blacklisted
    e = ConnectEvent{addr, true}
    c = gnetConnection(addr)
    Peers.AddBlacklistEntry(addr, time.Hour)
    pendingConnections[addr] = p
    Pool.Addresses[addr] = c
    assert.NotPanics(t, func() { onConnect(e) })
    // This connection should no longer be pending
    assert.Equal(t, len(pendingConnections), 0)
    // No message should have been sent
    assert.Equal(t, c.LastSent, time.Unix(0, 0))
    // We should not be expecting its version
    assert.Equal(t, len(expectingIntroductions), 0)
    // We should be looking to disconnect this client
    assert.Equal(t, len(Pool.DisconnectQueue), 1)
    if len(Pool.DisconnectQueue) == 0 {
        t.Fatal("Pool.DisconnectQueue is empty, would block")
    }
    de = <-Pool.DisconnectQueue
    assert.Equal(t, de.ConnId, 1)
    assert.Equal(t, de.Reason, DisconnectIsBlacklisted)
    // Cleanup
    delete(Peers.Blacklist, addr)

    resetPeers()
    ShutdownPool()
    gnet.EraseMessages()
}

func setupTestOnDisconnect(c *gnet.Connection, mirror uint32) {
    outgoingConnections[addr] = c
    expectingIntroductions[addr] = time.Now()
    mirrorConnections[mirror] = make(map[string]uint16)
    mirrorConnections[mirror][addrIP] = addrPort
    connectionMirrors[addr] = mirror
}

func TestOnDisconnect(t *testing.T) {
    Peers = pex.NewPex(maxPeers)
    c := gnetConnection(addr)
    var mirror uint32 = 100

    // Not blacklistable
    reason := DisconnectFailedSend
    setupTestOnDisconnect(c, mirror)
    assert.NotPanics(t, func() { onDisconnect(c, reason) })
    // Should not be in blacklist
    assert.Equal(t, len(Peers.Blacklist), 0)
    // Should no longer be in outgoingConnections
    assert.Equal(t, len(outgoingConnections), 0)
    // Should no longer be in expectingIntroductions
    assert.Equal(t, len(expectingIntroductions), 0)
    // Should be removed from the mirror, and the mirror dict for this ip
    // should be removed
    assert.Equal(t, len(mirrorConnections), 0)
    assert.Equal(t, len(connectionMirrors), 0)

    // Blacklistable
    reason = DisconnectVersionTimeout
    setupTestOnDisconnect(c, mirror)
    assert.NotPanics(t, func() { onDisconnect(c, reason) })
    assert.Equal(t, len(Peers.Blacklist), 1)
    assert.NotNil(t, Peers.Blacklist[addr])
    // Should be in blacklist
    assert.Equal(t, len(Peers.Blacklist), 1)
    assert.NotNil(t, Peers.Blacklist[addr])
    // Should no longer be in outgoingConnections
    assert.Equal(t, len(outgoingConnections), 0)
    // Should no longer be in expectingIntroductions
    assert.Equal(t, len(expectingIntroductions), 0)
    // Should be removed from the mirror, and the mirror dict for this ip
    // should be removed
    assert.Equal(t, len(mirrorConnections), 0)
    assert.Equal(t, len(connectionMirrors), 0)
    // Cleanup
    delete(Peers.Blacklist, addr)

    // mirrorConnections should retain a submap if there are other ports
    // inside
    reason = DisconnectFailedSend
    setupTestOnDisconnect(c, mirror)
    mirrorConnections[mirror][strings.Split(addrb, ":")[0]] = addrPort
    assert.NotPanics(t, func() { onDisconnect(c, reason) })
    // Should not be in blacklist
    assert.Equal(t, len(Peers.Blacklist), 0)
    // Should no longer be in outgoingConnections
    assert.Equal(t, len(outgoingConnections), 0)
    // Should no longer be in expectingIntroductions
    assert.Equal(t, len(expectingIntroductions), 0)
    // Should be removed from the mirror, and the mirror dict for this ip
    // should be removed
    assert.Equal(t, len(mirrorConnections), 1)
    assert.Equal(t, len(mirrorConnections[mirror]), 1)
    assert.Equal(t, len(connectionMirrors), 0)

    resetPeers()
}
