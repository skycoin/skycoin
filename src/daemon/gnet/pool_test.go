// build ignore

package gnet

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	addr          = "127.0.0.1:50823"
	port          = 50823
	address       = "127.0.0.1"
	silenceLogger = false
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
	p := NewConnectionPool(cfg, nil)

	cc := make(chan *Connection, 1)
	var wasSolicited *bool
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		wasSolicited = &solicited
		require.False(t, solicited)
		cc <- p.pool[1]
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	require.NotNil(t, p.listener)

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	c := <-cc
	require.NotNil(t, c)

	require.Len(t, p.addresses, 1)
	require.Len(t, p.pool, 1)

	require.Equal(t, conn.RemoteAddr().String(), c.Conn.LocalAddr().String())
	require.Equal(t, conn.LocalAddr().String(), c.Conn.RemoteAddr().String())

	require.NotNil(t, wasSolicited)
	require.False(t, *wasSolicited)

	p.Shutdown()
	<-q
}

func TestStartListenFailed(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)
	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	pp := NewConnectionPool(cfg, nil)
	err := pp.Run()
	require.Error(t, err)

	p.Shutdown()
	<-q
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

	p := NewConnectionPool(cfg, nil)

	// Unsolicited
	cc := make(chan *Connection, 1)
	var wasSolicited *bool
	p.Config.ConnectCallback = func(address string, solicited bool) {
		wasSolicited = &solicited
		cc <- p.pool[1]
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	c := <-cc
	require.NotNil(t, c)

	exist, err := p.IsConnExist(conn.LocalAddr().String())
	require.NoError(t, err)
	require.True(t, exist)

	delete(p.addresses, conn.LocalAddr().String())
	delete(p.pool, 1)

	require.NotNil(t, wasSolicited)
	require.False(t, *wasSolicited)

	// Solicited
	p.Config.ConnectCallback = func(address string, s bool) {
		wasSolicited = &s
		cc <- p.pool[2]
	}

	go p.handleConnection(conn, true)

	c = <-cc
	require.NotNil(t, c)
	require.Equal(t, addr, c.Addr())

	require.NotNil(t, wasSolicited)
	require.True(t, *wasSolicited)

	exist, err = p.IsConnExist(conn.RemoteAddr().String())
	require.NoError(t, err)
	require.True(t, exist)

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

	// Pool is shutdown, connect should fail
	wc := make(chan struct{})
	var connectErr error
	go func() {
		defer close(wc)
		connectErr = p.Connect(addr)
	}()

	<-wc

	require.Error(t, connectErr)
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

	// Setup a callback to capture the connection pointer so we can get the address
	cc := make(chan *Connection, 1)
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		cc <- p.pool[1]
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	c := <-cc
	require.NotNil(t, c)

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		require.Equal(t, addr, c.Addr())
	}

	p.Disconnect(c.Addr(), ErrDisconnectMalformedMessage)

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

type fakeConn struct {
	net.Conn
	addr string
}

func (f fakeConn) RemoteAddr() net.Addr {
	return fakeAddr{
		addr: f.addr,
	}
}

type fakeAddr struct {
	net.Addr
	addr string
}

func (f fakeAddr) String() string {
	return f.addr
}

func TestGetConnections(t *testing.T) {
	cfg := newTestConfig()
	p := NewConnectionPool(cfg, nil)

	c := &Connection{
		ID:   1,
		Conn: fakeConn{addr: "1.2.3.4"},
	}
	d := &Connection{
		ID:   2,
		Conn: fakeConn{addr: "2.3.4.5"},
	}
	e := &Connection{
		ID:   3,
		Conn: fakeConn{addr: "3.4.5.6"},
	}

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

	// Setup a callback to capture the connection pointer so we can get the address
	cc := make(chan *Connection, 1)
	i := 1
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		cc <- p.pool[i]
		i++
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	c := <-cc

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		t.Fatalf("Unexpected disconnect address=%s reason=%v", addr, reason)
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

	c = <-cc
	require.NotNil(t, c)

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		require.Equal(t, c.Addr(), addr)
		require.Equal(t, reason, ErrDisconnectInvalidMessageLength)
		require.Nil(t, p.pool[1])
		require.Nil(t, p.pool[2])
	}

	// Sending a length of < messagePrefixLength should cause a disconnect
	t.Logf("Pushing message with too small length")
	c.Buffer.Reset()
	conn.Write([]byte{messagePrefixLength - 1, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	wait()

	// Sending a length > MaxMessageLength should cause a disconnect
	conn, err = net.Dial("tcp", addr)
	require.NoError(t, err)

	c = <-cc
	require.NotNil(t, c)

	t.Logf("Pushing message with too large length")
	max := p.Config.MaxMessageLength
	p.Config.MaxMessageLength = 4
	p.Config.DisconnectCallback = func(addr string, r DisconnectReason) {
		require.Equal(t, ErrDisconnectInvalidMessageLength, r)
		require.Nil(t, p.pool[1])
		require.Nil(t, p.pool[2])
		require.Nil(t, p.pool[3])
	}
	conn.Write([]byte{5, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	wait()
	p.Config.MaxMessageLength = max

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

	// Setup a callback to capture the connection pointer so we can get the address
	cc := make(chan *Connection, 1)
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		cc <- p.pool[1]
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()

	wait()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	wait()

	// Go's internals seem to be blocking on network read, write something to
	// the connection to hopefully avoid this
	_, err = conn.Write([]byte{0})
	require.NoError(t, err)

	c := <-cc
	require.NotNil(t, c)

	m := NewByteMessage(88)
	// Send a successful message to b
	err = p.SendMessage(c.Addr(), m)
	require.NoError(t, err)

	var sr SendResult
	select {
	case sr = <-p.SendResults:
	case <-time.After(time.Second * 2):
		t.Fatal("No send results, would block")
	}

	require.Equal(t, sr.Message, m)
	require.Equal(t, sr.Addr, c.Addr())
	require.Nil(t, sr.Error)

	wait()

	var c2 *Connection
	p.strand("", func() error {
		fmt.Println("len(p.pool)", len(p.pool))
		c2 = p.pool[1]
		return nil
	})

	require.NotNil(t, c2)

	require.False(t, c2.LastSent.IsZero())
	fmt.Println("c2.LastSent", c2.LastSent)
	c2.LastSent = time.Time{}
	require.True(t, c2.LastSent.IsZero())
	fmt.Println("Reset c2.LastSent to zero")

	// Send a failed message to c2
	sendByteMessage = failingSendByteMessage

	disconnectErr := make(chan DisconnectReason, 1)
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		fmt.Printf("DisconnectCallback called, address=%s reason=%v\n", addr, reason)
		disconnectErr <- reason
	}
	err = p.SendMessage(c2.Addr(), m)
	require.NoError(t, err)

	select {
	case sr = <-p.SendResults:
	case <-time.After(time.Second * 2):
		t.Fatal("No send results, would block")
	}
	require.Equal(t, sr.Message, m)
	require.Equal(t, sr.Addr, c2.Addr())
	require.NotNil(t, sr.Error)

	reason := <-disconnectErr
	require.NotNil(t, reason)
	require.Equal(t, errors.New("send byte message failed"), reason)

	fmt.Println("Final c2.LastSent", c2.LastSent)
	require.True(t, c2.LastSent.IsZero())

	p.Shutdown()
	<-q
}

func TestPoolSendMessageOK(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()

	cfg := newTestConfig()
	cfg.WriteTimeout = time.Second
	cfg.SendResultsSize = 1
	cfg.ConnectionWriteQueueSize = 8
	p := NewConnectionPool(cfg, nil)

	// Setup a callback to capture the connection pointer so we can get the address
	cc := make(chan *Connection, 1)
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		cc <- p.pool[1]
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	c := <-cc
	m := NewByteMessage(88)
	err = p.SendMessage(c.Addr(), m)
	require.NoError(t, err)

	p.Shutdown()
	<-q
}

func TestPoolSendMessageWriteQueueFull(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()

	cfg := newTestConfig()
	cfg.WriteTimeout = time.Second
	cfg.SendResultsSize = 1
	cfg.ConnectionWriteQueueSize = 0
	p := NewConnectionPool(cfg, nil)

	// Setup a callback to capture the connection pointer so we can get the address
	cc := make(chan *Connection, 1)
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		cc <- p.pool[1]
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	_, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	c := <-cc

	// Send messages faster than can be processed to trigger ErrWriteQueueFull
	attempts := 100
	gotErr := false
	var once sync.Once
	m := NewByteMessage(88)
	addr := c.Addr()
	var wg sync.WaitGroup
	wg.Add(attempts)
	for i := 0; i < attempts; i++ {
		go func() {
			defer wg.Done()
			err := p.SendMessage(addr, m)
			if err == ErrWriteQueueFull {
				once.Do(func() {
					gotErr = true
				})
			}
		}()
	}

	wg.Wait()

	require.True(t, gotErr)

	p.Shutdown()
	<-q
}

func TestPoolBroadcastMessage(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()

	cfg := newTestConfig()
	cfg.ConnectionWriteQueueSize = 1
	p := NewConnectionPool(cfg, nil)

	ready := make(chan struct{})
	var i int
	var counterLock sync.Mutex
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		counterLock.Lock()
		defer counterLock.Unlock()
		i++
		if i == 2 {
			close(ready)
		}
	}

	q := make(chan struct{})
	go func() {
		defer close(q)
		p.Run()
	}()
	wait()

	conn1, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	// Go's internals seem to be blocking on network read, write something to
	// the connection to hopefully avoid this
	_, err = conn1.Write([]byte{0})
	require.NoError(t, err)

	conn2, err := net.Dial("tcp", addr)
	require.NoError(t, err)

	// Go's internals seem to be blocking on network read, write something to
	// the connection to hopefully avoid this
	_, err = conn2.Write([]byte{0})
	require.NoError(t, err)

	<-ready

	m := NewByteMessage(88)
	err = p.BroadcastMessage(m)
	require.NoError(t, err)

	attempts := 100
	gotErr := false
	var once sync.Once
	var wg sync.WaitGroup
	wg.Add(attempts)
	for i := 0; i < attempts; i++ {
		go func() {
			defer wg.Done()
			err := p.BroadcastMessage(m)
			if err == ErrNoReachableConnections {
				once.Do(func() {
					gotErr = true
				})
			}
		}()
	}

	wg.Wait()

	require.True(t, gotErr)

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
