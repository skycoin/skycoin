package daemon

import (
    "fmt"
    "github.com/op/go-logging"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "net"
    "testing"
    "time"
)

var (
    poolPort             = 6688
    addrIP               = "111.22.33.44"
    addrPort      uint16 = 5555
    addr                 = "111.22.33.44:5555"
    addrb                = "112.22.33.44:6666"
    badAddrPort          = "111.22.44.33:x"
    badAddrNoPort        = "111.22.44.33"
    silenceLogger        = true
)

func init() {
    if silenceLogger {
        logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
    }
}

func TestRegisterMessages(t *testing.T) {
    assert.NotPanics(t, RegisterMessages)
}

func TestNewIPAddr(t *testing.T) {
    i, err := NewIPAddr(addr)
    assert.Nil(t, err)
    assert.Equal(t, i.IP, uint32(1863721260))
    assert.Equal(t, i.Port, uint16(5555))

    bad := []string{"", "127.0.0.1", "127.0.0.1:x", "x:7777", ":",
        "127.0.0.1:7777:7777", "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
        "[1fff:0:a88:85a3::ac1f]:8001"}
    for _, b := range bad {
        _, err := NewIPAddr(b)
        assert.NotNil(t, err)
    }
}

func TestIPAddrString(t *testing.T) {
    i, err := NewIPAddr(addr)
    assert.Nil(t, err)
    assert.Equal(t, addr, i.String())
}

func testSimpleMessageHandler(t *testing.T, m gnet.Message) {
    assert.Nil(t, m.Handle(messageContext(addr)))
    assert.Equal(t, len(messageEvent), 1)
    if len(messageEvent) != 1 {
        t.Fatal("messageEvent is empty")
    }
    <-messageEvent
}

func TestGetPeersMessage(t *testing.T) {
    m := NewGetPeersMessage()
    testSimpleMessageHandler(t, m)
    Peers.AddPeer(addr)
    m.c = messageContext(addr)
    m.Process()
    assert.NotEqual(t, m.c.Conn.LastSent, time.Unix(0, 0))
    resetPeers()
}

func TestGivePeersMessage(t *testing.T) {
    addrs := []string{addr, addrb, "7"}
    peers := make([]*pex.Peer, 0, 3)
    for _, addr := range addrs {
        peers = append(peers, &pex.Peer{Addr: addr})
    }
    m := NewGivePeersMessage(peers)
    assert.Equal(t, len(m.GetPeers()), 2)
    testSimpleMessageHandler(t, m)
    assert.Equal(t, m.GetPeers()[0], addrs[0])
    assert.Equal(t, m.GetPeers()[1], addrs[1])
    // Peers should be added to the pex when processed
    m.Process()
    assert.Equal(t, len(Peers.Peerlist), 2)
    resetPeers()
}

func TestIntroductionMessageHandle(t *testing.T) {
    Pool = gnet.NewConnectionPool(poolPort)
    mc := messageContext(addr)
    m := NewIntroductionMessage()

    // Test valid handling
    m.Mirror = mirrorValue + 1
    err := m.Handle(mc)
    assert.Nil(t, err)
    if len(messageEvent) == 0 {
        t.Fatalf("messageEvent is empty")
    }
    <-messageEvent
    assert.True(t, m.valid)
    m.valid = false

    // Test matching mirror
    m.Mirror = mirrorValue
    err = m.Handle(mc)
    assert.Equal(t, err, DisconnectSelf)
    m.Mirror = mirrorValue + 1
    assert.False(t, m.valid)

    // Test mismatched version
    m.Version = version + 1
    err = m.Handle(mc)
    assert.Equal(t, err, DisconnectInvalidVersion)
    assert.False(t, m.valid)

    // Test already connected
    mirrorConnections[m.Mirror] = make(map[string]uint16)
    mirrorConnections[m.Mirror][addrIP] = addrPort + 1
    err = m.Handle(mc)
    assert.Equal(t, err, DisconnectConnectedTwice)
    delete(mirrorConnections, m.Mirror)
    assert.False(t, m.valid)
    Pool = nil

    for len(messageEvent) > 0 {
        <-messageEvent
    }
}

func TestIntroductionMessageProcess(t *testing.T) {
    Pool = gnet.NewConnectionPool(poolPort)
    m := NewIntroductionMessage()
    m.c = messageContext(addr)
    Pool.Pool[1] = m.c.Conn

    // Test invalid
    m.valid = false
    expectingVersions[addr] = time.Now()
    m.Process()
    // expectingVersions should get updated
    _, x := expectingVersions[addr]
    assert.False(t, x)
    // mirrorConnections should not have an entry
    _, x = mirrorConnections[m.Mirror]
    assert.False(t, x)
    assert.Equal(t, len(Peers.Peerlist), 0)

    // Test valid
    m.valid = true
    expectingVersions[addr] = time.Now()
    m.Process()
    // expectingVersions should get updated
    _, x = expectingVersions[addr]
    assert.False(t, x)
    assert.Equal(t, len(Peers.Peerlist), 1)
    assert.Equal(t, connectionMirrors[addrIP], m.Mirror)
    assert.NotNil(t, mirrorConnections[m.Mirror])
    assert.Equal(t, mirrorConnections[m.Mirror][addrIP], addrPort)
    peerAddr := fmt.Sprintf("%s:%d", addrIP, poolPort)
    assert.NotNil(t, Peers.Peerlist[peerAddr])
    resetPeers()

    // Handle impossibly bad ip:port returned from conn.Addr()
    // User should be disconnected
    m.valid = true
    m.c = messageContext(badAddrPort)
    m.Process()
    if len(Pool.DisconnectQueue) != 1 {
        t.Fatalf("DisconnectQueue empty")
    }
    <-Pool.DisconnectQueue

    m.valid = true
    m.c = messageContext(badAddrNoPort)
    m.Process()
    if len(Pool.DisconnectQueue) != 1 {
        t.Fatalf("DisconnectQueue empty")
    }
    <-Pool.DisconnectQueue

    Pool = nil
}

func TestPingMessage(t *testing.T) {
    m := &PingMessage{}
    testSimpleMessageHandler(t, m)

    m.c = messageContext(addr)
    m.Process()
    // A pong message should have been sent
    assert.NotEqual(t, m.c.Conn.LastSent, time.Unix(0, 0))
}

func TestPongMessage(t *testing.T) {
    m := &PongMessage{}
    // Pongs dont do anything
    assert.Nil(t, m.Handle(messageContext(addr)))
}

/* Helpers */

func messageContext(addr string) *gnet.MessageContext {
    conn := &gnet.Connection{
        Id:       1,
        Conn:     NewDummyConn(addr),
        LastSent: time.Unix(0, 0),
    }
    conn.Controller = gnet.NewConnectionController(conn)
    return &gnet.MessageContext{
        Conn: conn,
    }
}

func resetPeers() {
    Peers = pex.NewPex(maxPeers)
}

type DummyGivePeersMessage struct {
    peers []*pex.Peer
}

func (self *DummyGivePeersMessage) Send(c net.Conn) error {
    return nil
}

func (self *DummyGivePeersMessage) GetPeers() []string {
    p := make([]string, 0, len(self.peers))
    for _, ps := range self.peers {
        p = append(p, ps.Addr)
    }
    return p
}

type DummyAddr struct {
    addr string
}

func NewDummyAddr(addr string) net.Addr {
    return &DummyAddr{addr: addr}
}

func (self *DummyAddr) String() string {
    return self.addr
}

func (self *DummyAddr) Network() string {
    return self.addr
}

type DummyConn struct {
    net.Conn
    addr string
}

func NewDummyConn(addr string) net.Conn {
    return &DummyConn{addr: addr}
}

func (self *DummyConn) RemoteAddr() net.Addr {
    return NewDummyAddr(self.addr)
}

func (self *DummyConn) LocalAddr() net.Addr {
    return self.RemoteAddr()
}

func (self *DummyConn) Close() error {
    return nil
}

func (self *DummyConn) Read(b []byte) (int, error) {
    return 0, nil
}

func (self *DummyConn) SetWriteDeadline(t time.Time) error {
    return nil
}

func (self *DummyConn) Write(b []byte) (int, error) {
    return len(b), nil
}
