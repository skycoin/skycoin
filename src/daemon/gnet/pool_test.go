package gnet

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	logging "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)

var (
	addr          = "127.0.0.1:50823"
	addrb         = "127.0.0.1:50824"
	addrc         = "127.0.0.1:50825"
	port          = 50823
	address       = "127.0.0.1"
	listener      net.Listener
	conn          net.Conn
	silenceLogger = false
)

func init() {
	if silenceLogger {
		logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
	}
}

func TestNewConnectionPool(t *testing.T) {
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.MaxConnections = 108
	cfg.DialTimeout = time.Duration(777)
	p := NewConnectionPool(cfg, nil)
	assert.Equal(t, p.Config, cfg)
	assert.Equal(t, p.Config.Port, uint16(cfg.Port))
	assert.Equal(t, p.Config.Address, cfg.Address)
	assert.NotNil(t, p.pool)
	assert.Equal(t, len(p.pool), 0)
	assert.NotNil(t, p.addresses)
	assert.Equal(t, len(p.addresses), 0)
	assert.Equal(t, p.connId, 0)
}

func TestNewConnection(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.ConnectionWriteQueueSize = 101
	p := NewConnectionPool(cfg, nil)
	defer p.Shutdown()
	p.Run()
	wait()
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	wait()
	c := p.addresses[conn.LocalAddr().String()]
	assert.Equal(t, p.pool[p.connId], c)
	assert.Equal(t, p.connId, 1)
	assert.Equal(t, c.Addr(), conn.LocalAddr().String())
	assert.Equal(t, cap(c.WriteQueue), cfg.ConnectionWriteQueueSize)
	assert.NotNil(t, c.Buffer)
	assert.Equal(t, c.Buffer.Len(), 0)
	assert.Equal(t, c.ConnectionPool, p)
	assert.False(t, c.LastSent.IsZero())
	assert.False(t, c.LastReceived.IsZero())
}

func TestNewConnectionAlreadyConnected(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	defer p.Shutdown()
	p.Run()
	wait()
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	wait()
	c := p.addresses[conn.LocalAddr().String()]
	assert.NotNil(t, c)
	_, err = p.NewConnection(c.Conn)
	assert.NotNil(t, err)
}

func TestAcceptConnections(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	called := false
	cfg.ConnectCallback = func(addr string, solicited bool) {
		assert.False(t, solicited)
		// assert.NotNil(t, c)
		called = true
	}
	p := NewConnectionPool(cfg, nil)
	defer p.Shutdown()
	// go func() {
	p.Run()
	assert.NotNil(t, p.listener)
	// }()
	// go handleXConnections(p, 1)
	// go p.AcceptConnections()
	// Make a successful connection
	wait()
	// assert.NotNil(t, p.listener)
	c, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	defer c.Close()
	if err != nil {
		t.Fatalf("Dialing pool failed: %v", err)
	}
	wait()
	// assert.NotNil(t, p.listener)
	assert.Equal(t, len(p.addresses), 1)
	assert.Equal(t, len(p.pool), 1)
	if len(p.pool) == 0 {
		t.Fatalf("Pool empty, would crash")
	}
	assert.Equal(t, c.RemoteAddr().String(),
		p.pool[1].Conn.LocalAddr().String())
	assert.Equal(t, c.LocalAddr().String(),
		p.pool[1].Conn.RemoteAddr().String())
	assert.True(t, called)
}

func TestListeningAddress(t *testing.T) {
	wait()
	t.Run("listening", func(t *testing.T) {
		// cleanupNet()
		cfg := NewConfig()
		cfg.Address = ""
		cfg.Port = 0
		p := NewConnectionPool(cfg, nil)
		defer p.Shutdown()
		// assert.Nil(t, p.StartListen())
		p.Run()
		wait()
		// addr, err := p.ListeningAddress()
		// assert.Nil(t, err)
		// assert.NotNil(t, addr)
		t.Log("ListeningAddress: ", addr)
	})
}

func TestStartListen(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	called := false
	cfg.ConnectCallback = func(addr string, solicited bool) {
		assert.False(t, solicited)
		// assert.NotNil(t, c)
		called = true
	}
	p := NewConnectionPool(cfg, nil)
	defer p.Shutdown()
	p.Run()
	wait()
	_, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	wait()
	assert.True(t, called)
	assert.NotNil(t, p.listener)
	// assert.Nil(t, p.StartListen())
	// p.StopListen()
}

func TestStartListenTwice(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	defer p.Shutdown()
	p.Run()
	wait()
	assert.Panics(t, func() { p.Run() })
}

func TestStartListenFailed(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
	wait()
	q := NewConnectionPool(cfg, nil)
	// // // Can't listen on the same port
	assert.Panics(t, func() { q.Run() })
}

func TestStopListen(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	p.Run()
	wait()
	assert.NotNil(t, p.listener)
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	defer conn.Close()
	wait()
	assert.Equal(t, len(p.pool), 1)
	p.Shutdown()
	wait()
	assert.Nil(t, p.listener)
	assert.Equal(t, len(p.pool), 0)
	assert.Equal(t, len(p.addresses), 0)
	// Listening again should have no error
	p.Run()
	wait()
	p.Shutdown()
	wait()
	assert.Nil(t, p.listener)
	assert.Equal(t, len(p.pool), 0)
}

func TestHandleConnection(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address

	// Unsolicited
	called := false
	cfg.ConnectCallback = func(address string, s bool) {
		assert.False(t, s)
		// assert.Equal(t, address, addr)
		called = true
	}
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
	wait()
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	wait()
	assert.True(t, called)
	assert.True(t, p.IsConnExist(conn.LocalAddr().String()))
	called = false
	delete(p.addresses, conn.LocalAddr().String())
	delete(p.pool, 1)

	// Solicited
	p.Config.ConnectCallback = func(address string, s bool) {
		assert.True(t, s)
		assert.Equal(t, address, addr)
		called = true
	}

	go p.handleConnection(conn, true)
	wait()
	assert.True(t, p.IsConnExist(conn.RemoteAddr().String()))
	assert.True(t, called)
	called = false
	assert.Equal(t, len(p.addresses), 1)
	assert.Equal(t, len(p.pool), 1)
}

func TestConnect(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	// cfg.Port
	p := NewConnectionPool(cfg, nil)
	p.Run()
	wait()
	err := p.Connect(addr)
	wait()
	assert.Nil(t, err)

	// If already connected, should return same connection
	err = p.Connect(addr)
	wait()
	assert.Nil(t, err)

	delete(p.addresses, addr)

	p.Shutdown()
	wait()
	wc := make(chan struct{})
	go func() {
		p.Connect(addr)
		wc <- struct{}{}
	}()

	select {
	case <-wc:
		t.Fatal("Should not connection to a shutdown connection pool")
		return
	default:
		return
	}
}

func TestConnectNoTimeout(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.DialTimeout = 0
	cfg.Port += 1
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
	err := p.Connect(addr)
	wait()
	assert.NotNil(t, err)
}

func TestDisconnect(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
	wait()
	_, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	wait()
	c := p.pool[1]
	assert.NotNil(t, c)
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		assert.Equal(t, addr, c.Addr())
	}

	p.Disconnect(c.Addr(), DisconnectMalformedMessage)

	// Disconnecting a connection that isn't known has no effect
	// c = &Connection{Id: 88}

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		t.Fatal("disconnect unknow connection should not see this")
	}
	p.Disconnect("", nil)
}

func TestConnectionClose(t *testing.T) {
	wait()
	c := &Connection{
		Conn:       NewDummyConn(addr),
		Buffer:     &bytes.Buffer{},
		WriteQueue: make(chan Message),
	}

	c.Buffer.WriteByte(7)
	assert.Equal(t, c.Buffer.Len(), 1)
	c.Close()
	wait()
	assert.Nil(t, c.WriteQueue)
	assert.Equal(t, c.Buffer.Len(), 0)
}

func TestGetConnections(t *testing.T) {
	wait()
	p := NewConnectionPool(NewConfig(), nil)
	c := &Connection{Id: 1}
	d := &Connection{Id: 2}
	e := &Connection{Id: 3}
	p.pool[c.Id] = c
	p.pool[d.Id] = d
	p.pool[e.Id] = e
	p.Run()
	defer p.Shutdown()
	conns := p.GetConnections()
	assert.Equal(t, len(conns), 3)
	m := make(map[int]*Connection, 3)
	for i, c := range conns {
		m[c.Id] = &conns[i]
	}
	assert.Equal(t, len(m), 3)
	for i := 1; i <= 3; i++ {
		assert.Equal(t, m[i], p.pool[i])
	}
}

func TestConnectionReadLoop(t *testing.T) {
	wait()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
	wait()

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// assert.Equal(t, connID, 1)
		assert.Equal(t, reason, DisconnectReadFailed)
	}

	// 1:
	// Use a mock net.Conn that captures SetReadDeadline
	// and throws an error on Read
	reconn := NewReadErrorConn()
	go p.handleConnection(reconn, false)
	wait()
	assert.True(t, reconn.(*ReadErrorConn).ReadDeadlineSet != time.Time{})
	reconn.Close()

	// 2:
	// Use a mock net.Conn that fails on SetReadDeadline
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// assert.Equal(t, connID, 2)
		assert.Equal(t, reason, DisconnectSetReadDeadlineFailed)
	}

	rdfconn := &ReadDeadlineFailedConn{}
	go p.handleConnection(rdfconn, false)
	wait()
	rdfconn.Close()

	// 3:
	// Use a mock net.Conn that returns some bytes on Read
	// Look for these bytes copied into the eventChannel
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// assert.Equal(t, connID, 3)
		assert.Equal(t, reason, DisconnectInvalidMessageLength)
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
		// assert.Equal(t, connID, 4)
		assert.Equal(t, reason, DisconnectReadFailed)
	}
	go p.handleConnection(rnconn, false)
	wait()
	rnconn.stop()
	wait()
	rnconn.Close()
}

func TestProcessConnectionBuffers(t *testing.T) {
	wait()
	resetHandler()
	EraseMessages()
	RegisterMessage(DummyPrefix, DummyMessage{})
	RegisterMessage(ErrorPrefix, ErrorMessage{})
	VerifyMessages()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()

	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
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
	assert.NotEqual(t, c.LastReceived, time.Time{})
	assert.Equal(t, c.Buffer.Len(), 0)
	conn.Write([]byte{5, 0, 0, 0, 0})
	wait()
	assert.Equal(t, c.Buffer.Len(), 5)

	// Push multiple messages, the first causing an error, and confirm that
	// the remaining messages were unprocessed.
	// t.Logf("Pushing multiple messages, first one causing an error")
	c.Buffer.Reset()
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		assert.Equal(t, reason, errors.New("Bad"))
	}

	conn.Write([]byte{4, 0, 0, 0, 'E', 'R', 'R', 0x00})
	wait()
	assert.Equal(t, c.Buffer.Len(), 0)
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		t.Fatal("should not see this")
	}
	conn.Write([]byte{4, 0, 0, 0, 'D', 'U', 'M', 'Y'})
	wait()
	assert.Equal(t, c.Buffer.Len(), 0)

	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		assert.Equal(t, c.Addr(), addr)
		assert.Equal(t, reason, DisconnectInvalidMessageLength)
		assert.Nil(t, p.pool[1])
	}
	// Sending a length of < messagePrefixLength should cause a disconnect
	t.Logf("Pushing message with too small length")
	logger.Critical("666b")
	c.Buffer.Reset()
	logger.Critical("666c")
	conn.Write([]byte{messagePrefixLength - 1, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	wait()

	// Sending a length > MaxMessageLength should cause a disconnect
	conn, err = net.Dial("tcp", addr)
	assert.Nil(t, err)
	c = p.pool[2]
	t.Logf("Pushing message with too large length")
	max := p.Config.MaxMessageLength
	p.Config.MaxMessageLength = 4
	p.Config.DisconnectCallback = func(addr string, r DisconnectReason) {
		assert.Equal(t, DisconnectInvalidMessageLength, r)
		assert.Nil(t, p.pool[2])
	}
	// p.pool[1] = c
	// c.Buffer.Reset()
	conn.Write([]byte{5, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	wait()
	p.Config.MaxMessageLength = max

	// Send a malformed message, where ConvertToMessage fails
	// This is an unknown Message ID
	t.Logf("Pushing message with unknown ID")
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		c = p.addresses[addr]
		c.Buffer.Reset()
		conn.Write([]byte{4, 0, 0, 0, 'Y', 'Y', 'Y', 'Z'})
	}
	conn, err = net.Dial("tcp", addr)
	assert.Nil(t, err)
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		// assert.Equal(t, connID, 3)
		assert.Equal(t, reason, DisconnectMalformedMessage)
		assert.Nil(t, p.pool[3])
	}
	wait()
}

func TestConnectionWriteLoop(t *testing.T) {
	wait()
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()

	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
	wait()
	_, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
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
	assert.Equal(t, sr.Message, m)
	assert.Equal(t, sr.Addr, c.Addr())
	assert.Nil(t, sr.Error)
	assert.False(t, c.LastSent.IsZero())
	c.LastSent = time.Time{}
	// Send a failed message to c
	sendByteMessage = failingSendByteMessage
	p.Config.DisconnectCallback = func(addr string, reason DisconnectReason) {
		assert.Equal(t, reason, DisconnectWriteFailed)
	}
	p.SendMessage(c.Addr(), m)
	wait()
	if len(p.SendResults) == 0 {
		t.Fatalf("No send results, would block")
	}
	sr = <-p.SendResults
	assert.Equal(t, sr.Message, m)
	assert.Equal(t, sr.Addr, c.Addr())
	assert.NotNil(t, sr.Error)
	assert.True(t, c.LastSent.IsZero())
}

func TestPoolSendMessage(t *testing.T) {
	wait()
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	cfg := NewConfig()
	cfg.Address = address
	cfg.Port = uint16(port)
	cfg.WriteTimeout = time.Second
	// cfg.ConnectionWriteQueueSize = 1
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
	wait()
	assert.NotEqual(t, p.Config.ConnectionWriteQueueSize, 0)
	cc := make(chan *Connection)
	p.Config.ConnectCallback = func(addr string, solicited bool) {
		cc <- p.pool[1]
	}

	_, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	wait()

	c := <-cc
	m := NewByteMessage(88)
	p.SendMessage(c.Addr(), m)

	// queue full
	for i := 0; i < cap(c.WriteQueue); i++ {
		c.WriteQueue <- m
	}

	err = p.SendMessage(c.Addr(), m)
	assert.Equal(t, DisconnectWriteQueueFull, err)
}

func TestPoolBroadcastMessage(t *testing.T) {
	wait()
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	VerifyMessages()
	cfg := NewConfig()
	cfg.Address = address
	cfg.Port = uint16(port)
	p := NewConnectionPool(cfg, nil)
	p.Run()
	defer p.Shutdown()
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
		assert.Nil(t, reason)
	}
	_, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	_, err = net.Dial("tcp", addr)
	assert.Nil(t, err)

	<-ready

	m := NewByteMessage(88)
	p.BroadcastMessage(m)
	wait()
}

func TestPoolReceiveMessage(t *testing.T) {
	resetHandler()
	EraseMessages()
	RegisterMessage(BytePrefix, ByteMessage{})
	RegisterMessage(ErrorPrefix, ErrorMessage{})
	VerifyMessages()

	c := &Connection{}
	assert.True(t, c.LastReceived.IsZero())
	p := NewConnectionPool(NewConfig(), nil)

	// Valid message received
	b := make([]byte, 0)
	b = append(b, BytePrefix[:]...)
	b = append(b, byte(7))
	reason, err := p.receiveMessage(c, b)
	assert.Nil(t, err)
	assert.False(t, c.LastReceived.IsZero())
	assert.Nil(t, reason)

	// Invalid byte message received
	b = []byte{1}
	reason, err = p.receiveMessage(c, b)
	assert.NotNil(t, err)
	assert.Nil(t, reason)

	// Valid message, but handler returns a DisconnectReason
	b = make([]byte, 0)
	b = append(b, ErrorPrefix[:]...)
	reason, err = p.receiveMessage(c, b)
	assert.Nil(t, err)
	assert.NotNil(t, reason)
	assert.Equal(t, reason.Error(), "Bad")
}

// /* Helpers */

func wait() {
	time.Sleep(time.Millisecond * 50)
}

type DummyAddr struct {
	addr string
}

func NewDummyAddr(addr string) *DummyAddr {
	return &DummyAddr{
		addr: addr,
	}
}

func (self *DummyAddr) Network() string {
	return self.addr
}

func (self *DummyAddr) String() string {
	return self.Network()
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

type ReadErrorConn struct {
	net.Conn
	ReadDeadlineSet time.Time
}

func NewReadErrorConn() net.Conn {
	return &ReadErrorConn{nil, time.Time{}}
}

func (self *ReadErrorConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (self *ReadErrorConn) SetReadDeadline(t time.Time) error {
	self.ReadDeadlineSet = t
	return nil
}

func (self *ReadErrorConn) Read(b []byte) (int, error) {
	return 0, errors.New("failed")
}

func (self *ReadErrorConn) Close() error {
	return nil
}

type ReadDeadlineFailedConn struct {
	net.Conn
}

func (self *ReadDeadlineFailedConn) Read(b []byte) (int, error) {
	return 0, nil
}

func (self *ReadDeadlineFailedConn) SetReadDeadline(t time.Time) error {
	return errors.New("Failed")
}

func (self *ReadDeadlineFailedConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (self *ReadDeadlineFailedConn) Close() error {
	return nil
}

type ReadAlwaysConn struct {
	net.Conn
	stopReading bool
}

func (self *ReadAlwaysConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (self *ReadAlwaysConn) Close() error {
	return nil
}

func (self *ReadAlwaysConn) Read(b []byte) (int, error) {
	if self.stopReading {
		return 0, errors.New("done")
	}
	if len(b) == 0 {
		return 0, nil
	}
	b[0] = byte(88)
	time.Sleep(time.Millisecond * 2)
	return 1, nil
}

func (self *ReadAlwaysConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (self *ReadAlwaysConn) stop() {
	self.stopReading = true
}

type ReadNothingConn struct {
	net.Conn
	stopReading bool
}

func (self *ReadNothingConn) Read(b []byte) (int, error) {
	if self.stopReading {
		return 0, errors.New("done")
	} else {
		time.Sleep(time.Millisecond * 2)
		return 0, nil
	}
}

func (self *ReadNothingConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (self *ReadNothingConn) RemoteAddr() net.Addr {
	return NewDummyAddr(addr)
}

func (self *ReadNothingConn) Close() error {
	return nil
}

func (self *ReadNothingConn) stop() {
	self.stopReading = true
}
