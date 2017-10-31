// build ignore

package gnet

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/stretchr/testify/require"
)

const (
	addr          = "127.0.0.1:50823"
	addrb         = "127.0.0.1:50824"
	addrc         = "127.0.0.1:50825"
	port          = 50823
	address       = "127.0.0.1"
	silenceLogger = false
)

var (
	listener net.Listener
	conn     net.Conn
)

func init() {
	if silenceLogger {
		logging.Disable()
	}
}

func newTestConfig() Config {
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	return cfg
}

func TestNewConnectionPool(t *testing.T) {
	cfg := newTestConfig()
	cfg.MaxConnections = 108
	cfg.DialTimeout = time.Duration(777)

	p := NewConnectionPool(cfg, nil)
	require.Equal(t, p.Config, cfg)
	require.Equal(t, p.Config.Port, uint16(cfg.Port))
	require.Equal(t, p.Config.Address, cfg.Address)
	require.NotNil(t, p.pool)
	require.Equal(t, len(p.pool), 0)
	require.NotNil(t, p.addresses)
	require.Equal(t, len(p.addresses), 0)
	require.Equal(t, p.connID, 0)
}

func TestNewConnection(t *testing.T) {
	cfg := newTestConfig()
	cfg.ConnectionWriteQueueSize = 101
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()

	wait()
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()

	c := p.addresses[conn.LocalAddr().String()]
	require.Equal(t, p.pool[p.connID], c)
	require.Equal(t, p.connID, 1)
	require.Equal(t, c.Addr(), conn.LocalAddr().String())
	require.Equal(t, cap(c.WriteQueue), cfg.ConnectionWriteQueueSize)
	require.NotNil(t, c.Buffer)
	require.Equal(t, c.Buffer.Len(), 0)
	require.Equal(t, c.ConnectionPool, p)
	require.False(t, c.LastSent.IsZero())
	require.False(t, c.LastReceived.IsZero())

	p.Shutdown()
	<-q
}

func TestNewConnectionAlreadyConnected(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()

	c := p.addresses[conn.LocalAddr().String()]
	require.NotNil(t, c)

	_, err = p.NewConnection(c.Conn, true)
	require.Error(t, err)

	p.Shutdown()
	<-q
}

func TestAcceptConnections(t *testing.T) {
	cfg := newTestConfig()
	called := false
	cfg.ConnectCallback = func(addr string, solicited bool) {
		require.False(t, solicited)
		// require.NotNil(t, c)
		called = true
	}
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	require.NotNil(t, p.listener)

	c, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer c.Close()
	wait()

	require.Len(t, p.addresses, 1)
	require.Len(t, p.pool, 1)

	require.Equal(t, c.RemoteAddr().String(), p.pool[1].Conn.LocalAddr().String())
	require.Equal(t, c.LocalAddr().String(), p.pool[1].Conn.RemoteAddr().String())
	require.True(t, called)

	p.Shutdown()
	<-q
}

func TestStartListen(t *testing.T) {
	cfg := newTestConfig()
	called := false
	cfg.ConnectCallback = func(addr string, solicited bool) {
		require.False(t, solicited)
		// require.NotNil(t, c)
		called = true
	}
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()

	require.True(t, called)
	require.NotNil(t, p.listener)

	p.Shutdown()
	<-q
}

func TestStartListenFailed(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)
	qc := make(chan struct{})
	go func() {
		defer close(qc)
		p.Run()
	}()
	wait()
	q := NewConnectionPool(cfg, nil)
	require.NotNil(t, q.Run())
	p.Shutdown()
	<-qc
}

func TestStopListen(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	require.NotNil(t, p.listener)
	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()

	require.Equal(t, len(p.pool), 1)
	p.Shutdown()
	<-q

	require.Nil(t, p.listener)
}

func TestHandleConnection(t *testing.T) {
	cfg := newTestConfig()

	// Unsolicited
	called := false
	cfg.ConnectCallback = func(address string, s bool) {
		require.False(t, s)
		// require.Equal(t, address, addr)
		called = true
	}
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()

	require.True(t, called)

	exist, err := p.IsConnExist(conn.LocalAddr().String())
	require.NoError(t, err)
	require.True(t, exist)

	called = false
	delete(p.addresses, conn.LocalAddr().String())
	delete(p.pool, 1)

	// Solicited
	p.Config.ConnectCallback = func(address string, s bool) {
		require.True(t, s)
		require.Equal(t, address, addr)
		called = true
	}

	go p.handleConnection(conn, true)
	wait()

	exist, err = p.IsConnExist(conn.RemoteAddr().String())
	require.NoError(t, err)
	require.True(t, exist)
	require.True(t, called)

	called = false
	require.Equal(t, len(p.addresses), 1)
	require.Equal(t, len(p.pool), 1)

	p.Shutdown()
	<-q
}

func TestConnect(t *testing.T) {
	cfg := newTestConfig()
	// cfg.Port
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	err := p.Connect(addr)
	require.NoError(t, err)
	wait()

	// If already connected, should return same connection
	err = p.Connect(addr)
	require.NoError(t, err)
	wait()

	delete(p.addresses, addr)

	p.Shutdown()
	<-q

	wait()
	wc := make(chan struct{})
	go func() {
		p.Connect(addr)
		wc <- struct{}{}
	}()

	select {
	case <-wc:
		t.Fatal("Should not connect to an already-shutdown connection pool")
		return
	default:
		return
	}
}

func TestConnectNoTimeout(t *testing.T) {
	cfg := newTestConfig()
	cfg.DialTimeout = 0
	cfg.Port++

	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	p.Shutdown()
	<-q

	err := p.Connect(addr)
	wait()

	require.Error(t, err)
}

func TestDisconnect(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()

	c := p.pool[1]
	require.NotNil(t, c)
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		require.Equal(t, addr, c.Addr())
	}

	p.Disconnect(c.Addr(), ErrDisconnectMalformedMessage)

	// Disconnecting a connection that isn't known has no effect
	// c = &Connection{Id: 88}

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		t.Fatal("disconnect unknown connection should not see this")
	}
	p.Disconnect("", nil)

	p.Shutdown()
	<-q
}

func TestConnectionClose(t *testing.T) {
	c := &Connection{
		Conn:       NewDummyConn(addr),
		Buffer:     &bytes.Buffer{},
		WriteQueue: make(chan Message),
	}

	c.Buffer.WriteByte(7)
	require.Equal(t, c.Buffer.Len(), 1)
	c.Close()
	wait()

	require.Nil(t, c.WriteQueue)
	require.Equal(t, c.Buffer.Len(), 0)
}

func TestGetConnections(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)
	c := &Connection{ID: 1}
	d := &Connection{ID: 2}
	e := &Connection{ID: 3}
	p.pool[c.ID] = c
	p.pool[d.ID] = d
	p.pool[e.ID] = e

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	conns, err := p.GetConnections()
	require.NoError(t, err)
	require.Equal(t, len(conns), 3)

	m := make(map[int]*Connection, 3)
	for i, c := range conns {
		m[c.ID] = &conns[i]
	}

	require.Equal(t, len(m), 3)
	for i := 1; i <= 3; i++ {
		require.Equal(t, m[i], p.pool[i])
	}

	p.Shutdown()
	<-q
}

func TestConnectionReadLoop(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()

	wait()

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// require.Equal(t, connID, 1)
		require.Equal(t, reason, errors.New("read data failed: failed"))
	}

	// 1:
	// Use a mock net.Conn that captures SetReadDeadline
	// and throws an error on Read
	reconn := NewReadErrorConn()
	go p.handleConnection(reconn, false)
	wait()
	require.True(t, reconn.(*ReadErrorConn).ReadDeadlineSet != time.Time{})
	reconn.Close()

	// 2:
	// Use a mock net.Conn that fails on SetReadDeadline
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// require.Equal(t, connID, 2)
		require.Equal(t, reason, ErrDisconnectSetReadDeadlineFailed)
	}

	rdfconn := &ReadDeadlineFailedConn{}
	go p.handleConnection(rdfconn, false)
	wait()
	rdfconn.Close()

	// 3:
	// Use a mock net.Conn that returns some bytes on Read
	// Look for these bytes copied into the eventChannel
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// require.Equal(t, connID, 3)
		require.Equal(t, reason, ErrDisconnectInvalidMessageLength)
	}
	raconn := &ReadAlwaysConn{}
	go p.handleConnection(raconn, false)
	wait()
	raconn.stop()
	wait()
	raconn.Close()

	// 4: Use a mock net.Conn that successfully returns 0 bytes when read
	rnconn := &ReadNothingConn{}
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// require.Equal(t, connID, 4)
		require.Equal(t, reason, errors.New("read data failed: done"))
	}
	go p.handleConnection(rnconn, false)
	wait()
	rnconn.stop()
	wait()
	rnconn.Close()

	p.Shutdown()
	<-q
}

func TestProcessConnectionBuffers(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	RegisterMessage(ErrorPrefix, ErrorMessage{})
	VerifyMessages()
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()
	c := p.pool[1]

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		fmt.Println(reason)
		t.Fatal("should not see this")
	}

	conn.Write([]byte{4, 0, 0, 0})

	// A DummyMessage should have been processed
	conn.Write([]byte{'D', 'U', 'M', 'Y'})
	wait()
	require.NotEqual(t, c.LastReceived, time.Time{})
	require.Equal(t, c.Buffer.Len(), 0)
	conn.Write([]byte{5, 0, 0, 0, 0})
	wait()
	require.Equal(t, c.Buffer.Len(), 5)

	// Push multiple messages, the first causing an error, and confirm that
	// the remaining messages were unprocessed.
	t.Logf("Pushing multiple messages, first one causing an error")
	c.Buffer.Reset()
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		require.Equal(t, reason, errors.New("Bad"))
	}

	conn.Write([]byte{4, 0, 0, 0, 'E', 'R', 'R', 0x00})
	wait()
	require.Equal(t, c.Buffer.Len(), 0)

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		fmt.Println(reason)
		t.Fatal("should not see this")
	}
	conn.Write([]byte{4, 0, 0, 0, 'D', 'U', 'M', 'Y'})
	wait()
	require.Equal(t, c.Buffer.Len(), 0)

	conn, err = net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()
	c = p.pool[2]

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		require.Equal(t, c.Addr(), addr)
		require.Equal(t, reason, ErrDisconnectInvalidMessageLength)
		require.Nil(t, p.pool[1])
	}

	// Sending a length of < messagePrefixLength should cause a disconnect
	t.Logf("Pushing message with too small length")
	c.Buffer.Reset()
	conn.Write([]byte{messagePrefixLength - 1, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	wait()

	// Sending a length > MaxMessageLength should cause a disconnect
	conn, err = net.Dial("tcp", addr)
	require.NoError(t, err)
	c = p.pool[3]
	t.Logf("Pushing message with too large length")
	max := p.Config.MaxMessageLength
	p.Config.MaxMessageLength = 4
	p.Config.DisconnectCallback = func(addr string, r DisconnectReason) {
		require.Equal(t, ErrDisconnectInvalidMessageLength, r)
		require.Nil(t, p.pool[3])
	}
	conn.Write([]byte{5, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	wait()
	p.Config.MaxMessageLength = max

	// Send a malformed message, where ConvertToMessage fails
	// This is an unknown Message ID
	// t.Logf("Pushing message with unknown ID")
	// p.Config.ConnectCallback = func(addr string, solicited bool) {
	// 	c = p.addresses[addr]
	// 	c.Buffer.Reset()
	// 	conn.Write([]byte{4, 0, 0, 0, 'Y', 'Y', 'Y', 'Z'})
	// }
	// conn, err = net.Dial("tcp", addr)
	// require.NoError(t, err)
	// p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
	// 	// require.Equal(t, connID, 3)
	// 	require.Equal(t, reason, ErrDisconnectMalformedMessage)
	// 	require.Nil(t, p.pool[3])
	// }
	// wait()
	p.Shutdown()
	<-q
}

func TestConnectionWriteLoop(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()

	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)
	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()

	wait()
	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()
	c := p.pool[1]

	m := NewByteMessage(88)
	// Send a successful message to b
	p.SendMessage(c.Addr(), m)
	wait()

	if len(p.SendResults) == 0 {
		t.Fatalf("No send results, would block")
	}
	sr := <-p.SendResults
	require.Equal(t, sr.Message, m)
	require.Equal(t, sr.Addr, c.Addr())
	require.Nil(t, sr.Error)
	require.False(t, c.LastSent.IsZero())
	c.LastSent = time.Time{}
	// Send a failed message to c
	sendByteMessage = failingSendByteMessage
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		require.Equal(t, reason, errors.New("failed"))
	}
	p.SendMessage(c.Addr(), m)
	wait()

	if len(p.SendResults) == 0 {
		t.Fatalf("No send results, would block")
	}
	sr = <-p.SendResults
	require.Equal(t, sr.Message, m)
	require.Equal(t, sr.Addr, c.Addr())
	require.NotNil(t, sr.Error)
	require.True(t, c.LastSent.IsZero())

	p.Shutdown()
	<-q
}

func TestPoolSendMessage(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	cfg := newTestConfig()
	cfg.WriteTimeout = time.Second
	cfg.BroadcastResultSize = 1
	// cfg.ConnectionWriteQueueSize = 1
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	require.NotEqual(t, p.Config.ConnectionWriteQueueSize, 0)
	cc := make(chan *Connection)
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		cc <- p.pool[1]
	}

	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	wait()

	c := <-cc
	m := NewByteMessage(88)
	p.SendMessage(c.Addr(), m)

	// queue full
	for i := 0; i < cap(c.WriteQueue)+1; i++ {
		c.WriteQueue <- m
	}

	fmt.Printf("%v\n", len(c.WriteQueue))
	err = p.SendMessage(c.Addr(), m)
	require.Equal(t, ErrDisconnectWriteQueueFull, err)

	p.Shutdown()
	<-q
}

func TestPoolBroadcastMessage(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	ready := make(chan struct{})
	var i int
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		i++
		if i == 2 {
			ready <- struct{}{}
		}
	}

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		require.Nil(t, reason)
	}
	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	_, err = net.Dial("tcp", addr)
	require.NoError(t, err)

	<-ready

	m := NewByteMessage(88)
	err = p.BroadcastMessage(m)
	require.NoError(t, err)

	wait()

	p.Shutdown()
	<-q
}

func TestPoolReceiveMessage(t *testing.T) {
	wait()
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	RegisterMessage(ErrorPrefix, ErrorMessage{})
	VerifyMessages()

	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	c := NewConnection(p, 1, NewDummyConn(addr), 10, true)

	// Valid message received
	b := make([]byte, 0)
	b = append(b, BytePrefix[:]...)
	b = append(b, byte(7))
	err := p.receiveMessage(c, b)
	require.NoError(t, err)
	require.False(t, c.LastReceived.IsZero())

	// Invalid byte message received
	b = []byte{1}
	err = p.receiveMessage(c, b)
	require.Error(t, err)

	// Valid message, but handler returns a DisconnectReason
	b = make([]byte, 0)
	b = append(b, ErrorPrefix[:]...)
	err = p.receiveMessage(c, b)
	require.Equal(t, err.Error(), "Bad")

	p.Shutdown()
	<-q
}

// Helpers

func wait() {
	time.Sleep(time.Millisecond * 100)
}

type DummyAddr struct {
	addr string
}

func NewDummyAddr(addr string) *DummyAddr {
	return &DummyAddr{
		addr: addr,
	}
}

func (da *DummyAddr) Network() string {
	return da.addr
}

func (da *DummyAddr) String() string {
	return da.Network()
}

type DummyConn struct {
	net.Conn
	addr string
}

func NewDummyConn(addr string) net.Conn {
	return &DummyConn{addr: addr}
}

func (dc *DummyConn) RemoteAddr() net.Addr {
	return NewDummyAddr(dc.addr)
}

func (dc *DummyConn) LocalAddr() net.Addr {
	return dc.RemoteAddr()
}

func (dc *DummyConn) Close() error {
	return nil
}

func (dc *DummyConn) Read(b []byte) (int, error) {
	return 0, nil
}

func (dc *DummyConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (dc *DummyConn) Write(b []byte) (int, error) {
	return len(b), nil
}

type ReadErrorConn struct {
	net.Conn
	ReadDeadlineSet time.Time
}

func NewReadErrorConn() net.Conn {
	return &ReadErrorConn{nil, time.Time{}}
}

func (rec *ReadErrorConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (rec *ReadErrorConn) SetReadDeadline(t time.Time) error {
	rec.ReadDeadlineSet = t
	return nil
}

func (rec *ReadErrorConn) Read(b []byte) (int, error) {
	return 0, errors.New("failed")
}

func (rec *ReadErrorConn) Close() error {
	return nil
}

type ReadDeadlineFailedConn struct {
	net.Conn
}

func (c *ReadDeadlineFailedConn) Read(b []byte) (int, error) {
	return 0, nil
}

func (c *ReadDeadlineFailedConn) SetReadDeadline(t time.Time) error {
	return errors.New("Failed")
}

func (c *ReadDeadlineFailedConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (c *ReadDeadlineFailedConn) Close() error {
	return nil
}

type ReadAlwaysConn struct {
	net.Conn
	stopReading bool
}

func (c *ReadAlwaysConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (c *ReadAlwaysConn) Close() error {
	return nil
}

func (c *ReadAlwaysConn) Read(b []byte) (int, error) {
	if c.stopReading {
		return 0, errors.New("done")
	}
	if len(b) == 0 {
		return 0, nil
	}
	b[0] = byte(88)
	time.Sleep(time.Millisecond * 2)
	return 1, nil
}

func (c *ReadAlwaysConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *ReadAlwaysConn) stop() {
	c.stopReading = true
}

type ReadNothingConn struct {
	net.Conn
	stopReading bool
}

func (c *ReadNothingConn) Read(b []byte) (int, error) {
	if c.stopReading {
		return 0, errors.New("done")
	}
	time.Sleep(time.Millisecond * 2)
	return 0, nil
}

func (c *ReadNothingConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *ReadNothingConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (c *ReadNothingConn) Close() error {
	return nil
}

func (c *ReadNothingConn) stop() {
	c.stopReading = true
}
