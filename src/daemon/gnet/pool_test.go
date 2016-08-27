package gnet

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/op/go-logging.v1"
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

func newNetConn() net.Conn {
	if conn != nil {
		stopConn()
	}
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	return conn
}

func listen() {
	if listener != nil {
		stopListen()
	}
	var err error
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
}

func stopListen() {
	if listener != nil {
		if err := listener.Close(); err != nil {
			log.Panic(err)
		}
	}
	listener = nil
}

func stopConn() {
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Panic(err)
		}
	}
	conn = nil
}

func cleanupNet() {
	stopListen()
	stopConn()
}

func TestConnectionAddr(t *testing.T) {
	cleanupNet()
	listen()
	c := &Connection{Conn: newNetConn()}
	assert.Equal(t, c.Addr(), addr)
	assert.Equal(t, c.String(), addr)
}

func TestNewConnectionPool(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.MaxConnections = 108
	cfg.EventChannelSize = 101
	cfg.DialTimeout = time.Duration(777)
	p := NewConnectionPool(cfg, nil)
	assert.Equal(t, p.Config, cfg)
	assert.Equal(t, p.Config.Port, uint16(cfg.Port))
	assert.Equal(t, p.Config.Address, cfg.Address)
	assert.NotNil(t, p.Pool)
	assert.Equal(t, len(p.Pool), 0)
	assert.NotNil(t, p.Addresses)
	assert.Equal(t, len(p.Addresses), 0)
	assert.NotNil(t, p.eventChannel)
	assert.Equal(t, cap(p.eventChannel), 101)
	assert.NotNil(t, p.DisconnectQueue)
	assert.Equal(t, cap(p.DisconnectQueue), 108)
	assert.NotNil(t, p.connectionQueue)
	assert.Equal(t, cap(p.connectionQueue), 0)
	assert.Equal(t, p.connId, 0)
	cleanupNet()
}

func TestNewConnection(t *testing.T) {
	cleanupNet()
	listen()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.ConnectionWriteQueueSize = 101
	p := NewConnectionPool(cfg, nil)
	c := newPoolConnection(p)
	assert.Equal(t, p.Addresses[addr], c)
	assert.Equal(t, p.Pool[p.connId], c)
	assert.Equal(t, p.connId, 1)
	assert.Equal(t, c.Addr(), addr)
	assert.Equal(t, cap(c.WriteQueue), cfg.ConnectionWriteQueueSize)
	assert.NotNil(t, c.Buffer)
	assert.Equal(t, c.Buffer.Len(), 0)
	assert.Equal(t, c.ConnectionPool, p)
	assert.False(t, c.LastSent.IsZero())
	assert.False(t, c.LastReceived.IsZero())
}

func TestNewConnectionAlreadyConnected(t *testing.T) {
	cleanupNet()
	listen()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	c := newPoolConnection(p)
	assert.Panics(t, func() { p.NewConnection(c.Conn) })
}

func TestAcceptConnections(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	called := false
	cfg.ConnectCallback = func(c *Connection, solicited bool) {
		assert.False(t, solicited)
		assert.NotNil(t, c)
		called = true
	}
	p := NewConnectionPool(cfg, nil)
	defer p.StopListen()
	go func() {
		assert.Nil(t, p.StartListen())
		assert.NotNil(t, p.listener)
	}()
	go handleXConnections(p, 1)
	go p.AcceptConnections()
	// Make a successful connection
	wait()
	assert.NotNil(t, p.listener)
	c, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	defer c.Close()
	if err != nil {
		t.Fatalf("Dialing pool failed: %v", err)
	}
	wait()
	assert.NotNil(t, p.listener)
	assert.Equal(t, len(p.Addresses), 1)
	assert.Equal(t, len(p.Pool), 1)
	if len(p.Pool) == 0 {
		t.Fatalf("Pool empty, would crash")
	}
	assert.Equal(t, c.RemoteAddr().String(),
		p.Pool[1].Conn.LocalAddr().String())
	assert.Equal(t, c.LocalAddr().String(),
		p.Pool[1].Conn.RemoteAddr().String())
	assert.True(t, called)
}

func TestStartListen(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	called := false
	cfg.ConnectCallback = func(c *Connection, solicited bool) {
		assert.False(t, solicited)
		assert.NotNil(t, c)
		called = true
	}
	p := NewConnectionPool(cfg, nil)
	assert.Nil(t, p.listener)
	assert.Nil(t, p.StartListen())
	p.StopListen()
}

func TestStartListenTwice(t *testing.T) {
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	go func() {
		assert.Nil(t, p.StartListen())
	}()
	wait()
	// Listening twice will fail
	assert.Panics(t, func() { assert.Nil(t, p.StartListen()) })
	p.StopListen()
}

func TestStartListenFailed(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	q := NewConnectionPool(cfg, nil)
	go func() {
		assert.Nil(t, p.StartListen())
	}()
	wait()
	// Can't listen on the same port
	assert.NotNil(t, q.StartListen())
	p.StopListen()
}

func TestStopListen(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	go func() {
		assert.Nil(t, p.StartListen())
	}()
	go handleXConnections(p, 1)
	go p.AcceptConnections()
	wait()
	assert.NotNil(t, p.listener)
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	defer conn.Close()
	wait()
	assert.Equal(t, len(p.Pool), 1)
	p.StopListen()
	wait()
	assert.Nil(t, p.listener)
	assert.Equal(t, len(p.Pool), 0)
	assert.Equal(t, len(p.Addresses), 0)
	// Listening again should have no error
	assert.NotPanics(t, func() {
		assert.Nil(t, p.StartListen())
	})
	p.StopListen()
	wait()
	assert.Nil(t, p.listener)
	assert.Equal(t, len(p.Pool), 0)
}

func TestHandleConnection(t *testing.T) {
	cleanupNet()
	listen()
	conn := newNetConn()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address

	// Unsolicited
	called := false
	cfg.ConnectCallback = func(c *Connection, s bool) {
		assert.False(t, s)
		assert.Equal(t, c.Conn.RemoteAddr().String(), addr)
		called = true
	}
	p := NewConnectionPool(cfg, nil)
	go handleXConnections(p, 2)
	var c *Connection
	go func() {
		assert.NotNil(t, conn)
		c = p.handleConnection(conn, false)
	}()
	wait()
	assert.True(t, called)
	assert.Equal(t, c.Addr(), conn.RemoteAddr().String())
	called = false
	delete(p.Addresses, c.Addr())
	delete(p.Pool, 1)

	// Solicited
	p.Config.ConnectCallback = func(c *Connection, s bool) {
		assert.True(t, s)
		assert.Equal(t, c.Conn.RemoteAddr().String(), addr)
		called = true
	}
	go func() {
		c = p.handleConnection(conn, true)
	}()
	wait()
	assert.Equal(t, c.Addr(), conn.RemoteAddr().String())
	assert.True(t, called)
	called = false
	assert.Equal(t, len(p.Addresses), 1)
	assert.Equal(t, len(p.Pool), 1)

	// Already known panics
	assert.Panics(t, func() { p.handleConnection(conn, false) })
	stopConn()
}

func TestConnect(t *testing.T) {
	cleanupNet()
	listen()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.Port += 1
	p := NewConnectionPool(cfg, nil)
	go handleXConnections(p, 1)
	c, err := p.Connect(addr)
	wait()
	assert.Nil(t, err)
	assert.NotNil(t, c)
	if c == nil {
		t.Fatalf("No connection made")
	}

	// If already connected, should return same connection
	newC, err := p.Connect(addr)
	wait()
	assert.Nil(t, err)
	assert.Equal(t, c, newC)
	c.Conn.Close()
	delete(p.Addresses, addr)

	// Failed dial should return error
	stopListen()
	c, err = p.Connect(addr)
	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestConnectNoTimeout(t *testing.T) {
	cleanupNet()
	listen()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.DialTimeout = 0
	cfg.Port += 1
	p := NewConnectionPool(cfg, nil)
	go handleXConnections(p, 1)
	c, err := p.Connect(addr)
	wait()
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestDisconnect(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	assert.Nil(t, p.StartListen())
	go handleXConnections(p, 1)
	go p.AcceptConnections()
	wait()
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	wait()
	c := p.Pool[1]
	assert.NotNil(t, c)
	p.Disconnect(c, DisconnectMalformedMessage)
	de := <-p.DisconnectQueue
	assert.Equal(t, de.ConnId, 1)
	assert.Equal(t, de.Reason, DisconnectMalformedMessage)
	conn.Close()

	// Disconnecting a connection that isn't known has no effect
	c = &Connection{Id: 88}
	p.Disconnect(c, nil)
	assert.Equal(t, len(p.DisconnectQueue), 0)

	// Disconnecting a nil connection has no effect
	p.Disconnect(nil, nil)
	assert.Equal(t, len(p.DisconnectQueue), 0)
	p.StopListen()
}

func TestConnectionClose(t *testing.T) {
	c := &Connection{
		Conn:          NewDummyConn(addr),
		Buffer:        &bytes.Buffer{},
		WriteQueue:    make(chan Message),
		readLoopDone:  make(chan bool, 1),
		writeLoopDone: make(chan bool, 1),
	}
	c.readLoopDone <- false
	c.writeLoopDone <- false
	p := NewConnectionPool(NewConfig(), nil)
	go p.ConnectionWriteLoop(c)
	wait()
	c.Buffer.WriteByte(7)
	assert.Equal(t, c.Buffer.Len(), 1)
	c.Close()
	wait()
	assert.Nil(t, c.WriteQueue)
	assert.Equal(t, c.Buffer.Len(), 0)
}

func TestHandleDisconnectionEvent(t *testing.T) {
	cleanupNet()
	listen()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	cfg.Port += 1
	called := false
	expectReason := DisconnectReadFailed
	cfg.DisconnectCallback = func(c *Connection, r DisconnectReason) {
		called = true
		assert.Equal(t, r, expectReason)
		assert.Equal(t, c.Id, 1)
	}
	p := NewConnectionPool(cfg, nil)
	go handleXConnections(p, 1)
	wait()
	c, err := p.Connect(addr)
	wait()
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, len(p.Pool), 1)
	assert.Equal(t, len(p.Addresses), 1)

	p.Disconnect(c, expectReason)
	assert.Equal(t, len(p.DisconnectQueue), 1)
	if len(p.DisconnectQueue) != 1 {
		t.Fatal("There is nothing in the DisconnectQueue")
	}
	de := <-p.DisconnectQueue
	assert.Equal(t, de.ConnId, 1)
	assert.Equal(t, de.Reason, expectReason)

	p.HandleDisconnectEvent(de)
	assert.Equal(t, len(p.Pool), 0)
	assert.Equal(t, len(p.Addresses), 0)
	assert.True(t, called)
	called = false

	// Handling event with unknown ConnId is ignored
	p.HandleDisconnectEvent(DisconnectEvent{ConnId: 88})
	assert.False(t, called)
}

func TestGetConnections(t *testing.T) {
	cleanupNet()
	p := NewConnectionPool(NewConfig(), nil)
	c := &Connection{Id: 1}
	d := &Connection{Id: 2}
	e := &Connection{Id: 3}
	p.Pool[c.Id] = c
	p.Pool[d.Id] = d
	p.Pool[e.Id] = e

	conns := p.GetConnections()
	assert.Equal(t, len(conns), 3)
	m := make(map[int]*Connection, 3)
	for _, c := range conns {
		m[c.Id] = c
	}
	assert.Equal(t, len(m), 3)
	for i := 1; i <= 3; i++ {
		assert.Equal(t, m[i], p.Pool[i])
	}
}

func TestGetRawConnections(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	assert.Nil(t, p.StartListen())
	go handleXConnections(p, 2)
	go p.AcceptConnections()
	wait()
	conns := p.GetRawConnections()
	assert.Equal(t, len(conns), 0)

	c, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	wait()

	conns = p.GetRawConnections()
	assert.Equal(t, len(conns), 1)
	assert.Equal(t, conns[0].RemoteAddr().String(), c.LocalAddr().String())

	d, err := net.Dial("tcp", addr)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	wait()

	conns = p.GetRawConnections()
	assert.Equal(t, len(conns), 2)
	if conns[0].RemoteAddr().String() == c.LocalAddr().String() {
		assert.Equal(t, conns[0].RemoteAddr().String(), c.LocalAddr().String())
		assert.Equal(t, conns[1].RemoteAddr().String(), d.LocalAddr().String())
	} else {
		assert.Equal(t, conns[0].RemoteAddr().String(), d.LocalAddr().String())
		assert.Equal(t, conns[1].RemoteAddr().String(), c.LocalAddr().String())
	}

	p.StopListen()
	c.Close()
	d.Close()
}

func TestConnectionReadLoop(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	wait()

	// 1:
	// Use a mock net.Conn that captures SetReadDeadline
	// and throws an error on Read
	reconn := NewReadErrorConn()
	c := &Connection{
		Id:            1,
		Conn:          reconn,
		Buffer:        &bytes.Buffer{},
		WriteQueue:    make(chan Message),
		readLoopDone:  make(chan bool, 1),
		writeLoopDone: make(chan bool, 1),
	}
	c.writeLoopDone <- false
	c.readLoopDone <- false
	p.Pool[1] = c
	go p.connectionReadLoop(c)
	wait()
	assert.True(t, reconn.(*ReadErrorConn).ReadDeadlineSet != time.Time{})
	assert.Equal(t, len(c.readLoopDone), 1)
	assert.Equal(t, len(p.DisconnectQueue), 1)
	if len(p.DisconnectQueue) == 0 {
		t.Fatalf("There is nothing in the DisconnectQueue")
	}
	de := <-p.DisconnectQueue
	assert.Equal(t, de.ConnId, 1)
	assert.Equal(t, de.Reason, DisconnectReadFailed)
	p.HandleDisconnectEvent(de)
	reconn.Close()

	// 2:
	// Use a mock net.Conn that fails on SetReadDeadline
	rdfconn := &ReadDeadlineFailedConn{}
	c = &Connection{
		Id:            2,
		Conn:          rdfconn,
		Buffer:        &bytes.Buffer{},
		WriteQueue:    make(chan Message),
		readLoopDone:  make(chan bool, 1),
		writeLoopDone: make(chan bool, 1),
	}
	c.writeLoopDone <- false
	c.readLoopDone <- false
	p.Pool[2] = c
	go p.connectionReadLoop(c)
	wait()
	assert.Equal(t, len(c.readLoopDone), 1)
	assert.Equal(t, len(p.DisconnectQueue), 1)
	if len(p.DisconnectQueue) != 1 {
		t.Fatalf("There is nothing in the DisconnectQueue")
	}
	de = <-p.DisconnectQueue
	assert.Equal(t, de.ConnId, 2)
	assert.Equal(t, de.Reason, DisconnectSetReadDeadlineFailed)
	p.HandleDisconnectEvent(de)
	rdfconn.Close()

	// 3:
	// Use a mock net.Conn that returns some bytes on Read
	// Look for these bytes copied into the eventChannel
	raconn := &ReadAlwaysConn{}
	c = &Connection{
		Id:            3,
		Conn:          raconn,
		Buffer:        &bytes.Buffer{},
		WriteQueue:    make(chan Message),
		readLoopDone:  make(chan bool, 1),
		writeLoopDone: make(chan bool, 1),
	}
	c.writeLoopDone <- false
	c.readLoopDone <- false
	p.Pool[3] = c
	go p.connectionReadLoop(c)
	wait()
	raconn.stop()
	wait()
	assert.Equal(t, len(c.readLoopDone), 1)
	if len(p.DisconnectQueue) == 0 {
		t.Fatalf("There is nothing in the DisconnectQueue")
	}
	de = <-p.DisconnectQueue
	assert.Equal(t, de.ConnId, 3)
	assert.Equal(t, de.Reason, DisconnectReadFailed)
	p.HandleDisconnectEvent(de)
	assert.True(t, len(p.eventChannel) > 0)
	for len(p.eventChannel) > 0 {
		e := <-p.eventChannel
		assert.Equal(t, e.ConnId, 3)
		assert.Equal(t, len(e.Data), 1)
		assert.Equal(t, e.Data[0], byte(88))
	}
	raconn.Close()
	assert.Equal(t, len(p.eventChannel), 0)

	// 4: Use a mock net.Conn that successfully returns 0 bytes when read
	rnconn := &ReadNothingConn{}
	c = &Connection{
		Id:            4,
		Conn:          rnconn,
		Buffer:        &bytes.Buffer{},
		WriteQueue:    make(chan Message),
		readLoopDone:  make(chan bool, 1),
		writeLoopDone: make(chan bool, 1),
	}
	c.writeLoopDone <- false
	c.readLoopDone <- false
	p.Pool[4] = c
	go p.connectionReadLoop(c)
	wait()
	rnconn.stop()
	wait()
	assert.Equal(t, len(c.readLoopDone), 1)
	if len(p.DisconnectQueue) == 0 {
		t.Fatalf("There is nothing in the DisconnectQueue")
	}
	de = <-p.DisconnectQueue
	assert.Equal(t, de.ConnId, 4)
	assert.Equal(t, de.Reason, DisconnectReadFailed)
	p.HandleDisconnectEvent(de)
	assert.Equal(t, len(p.eventChannel), 0)
	rnconn.Close()

	p.StopListen()
}

func TestHandleConnectionQueue(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	c := &Connection{Id: 1, Conn: &ReadAlwaysConn{}}
	go func() {
		p.connectionQueue <- c
	}()
	wait()
	d := p.handleConnectionQueue()
	assert.Equal(t, c, d)
	assert.Equal(t, len(p.Pool), 1)
	assert.Equal(t, len(p.Addresses), 1)
	assert.NotNil(t, p.Pool[1])
	assert.NotNil(t, p.Addresses[c.Addr()])

	d = p.handleConnectionQueue()
	assert.Nil(t, d)
}

func TestProcessEvents(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	c := &Connection{
		Id:     1,
		Buffer: &bytes.Buffer{},
		Conn:   &ReadAlwaysConn{},
	}
	p.Pool[1] = c
	p.eventChannel <- dataEvent{ConnId: 1, Data: []byte{1}}
	p.eventChannel <- dataEvent{ConnId: 1, Data: []byte{2}}
	p.eventChannel <- dataEvent{ConnId: 1, Data: []byte{3}}
	// Intentionally unknown ConnId
	p.eventChannel <- dataEvent{ConnId: 2, Data: []byte{3, 0, 0, 0}}

	p.processEvents()
	assert.Equal(t, len(p.eventChannel), 0)
	assert.True(t, bytes.Equal(c.Buffer.Bytes(), []byte{1, 2, 3}))
}

func TestProcessConnectionBuffers(t *testing.T) {
	cleanupNet()
	RegisterMessage(DummyPrefix, DummyMessage{})
	RegisterMessage(ErrorPrefix, ErrorMessage{})
	VerifyMessages()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	var checkDisconnect = func(reason DisconnectReason) bool {
		if len(p.DisconnectQueue) > 0 {
			de := <-p.DisconnectQueue
			assert.Equal(t, reason, de.Reason)
			assert.Equal(t, de.ConnId, 1)
			logger.Critical("handle disc event")
			p.HandleDisconnectEvent(de)
			return true
		}
		return false
	}

	// No connections, so nothing should happen
	assert.NotPanics(t, p.processConnectionBuffers)

	c := &Connection{
		Id:            1,
		Buffer:        &bytes.Buffer{},
		Conn:          &ReadAlwaysConn{},
		LastReceived:  time.Time{},
		readLoopDone:  make(chan bool, 1),
		writeLoopDone: make(chan bool, 1),
	}
	c.readLoopDone <- false
	c.writeLoopDone <- false
	p.Pool[1] = c
	// No data in buffer
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.False(t, checkDisconnect(nil))

	c.Buffer.Write([]byte{0})
	// Not enough data in buffer to read length
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.False(t, checkDisconnect(nil))

	c.Buffer.Reset()
	c.Buffer.Write([]byte{4, 0, 0, 0})
	// Enough data for length, but no content data
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.False(t, checkDisconnect(nil))

	// A DummyMessage should have been processed
	c.Buffer.Write([]byte{'D', 'U', 'M', 'Y'})
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.NotEqual(t, c.LastReceived, time.Time{})
	assert.Equal(t, c.Buffer.Len(), 0)
	assert.False(t, checkDisconnect(nil))

	// Enough data for length but not enough to process the length
	c.Buffer.Reset()
	c.Buffer.Write([]byte{5, 0, 0, 0, 0})
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.Equal(t, c.Buffer.Len(), 5)
	assert.False(t, checkDisconnect(nil))

	// Push multiple messages, the first causing an error, and confirm that
	// the remaining messages were unprocessed.
	t.Logf("Pushing multiple messages, first one causing an error")
	c.Buffer.Reset()
	c.Buffer.Write([]byte{4, 0, 0, 0, 'E', 'R', 'R', 0x00})
	c.Buffer.Write([]byte{4, 0, 0, 0, 'D', 'U', 'M', 'Y'})
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.Equal(t, c.Buffer.Len(), 8)
	expect := []byte{4, 0, 0, 0, 'D', 'U', 'M', 'Y'}
	assert.True(t, bytes.Equal(c.Buffer.Bytes(), expect))
	assert.False(t, checkDisconnect(nil))

	// Sending a length of < messagePrefixLength should cause a disconnect
	t.Logf("Pushing message with too small length")
	logger.Critical("666b")
	c.Buffer.Reset()
	logger.Critical("666c")
	c.Buffer.Write([]byte{messagePrefixLength - 1, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	logger.Critical("666d")
	assert.NotPanics(t, p.processConnectionBuffers)
	logger.Critical("666e")
	assert.True(t, checkDisconnect(DisconnectInvalidMessageLength))
	logger.Critical("666f")
	assert.Nil(t, p.Pool[1])

	// Sending a length > MaxMessageLength should cause a disconnect
	c.readLoopDone <- false
	c.writeLoopDone <- false
	t.Logf("Pushing message with too large length")
	max := p.Config.MaxMessageLength
	p.Config.MaxMessageLength = 4
	p.Pool[1] = c
	c.Buffer.Reset()
	c.Buffer.Write([]byte{5, 0, 0, 0, 'B', 'Y', 'T', 'E'})
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.True(t, checkDisconnect(DisconnectInvalidMessageLength))
	assert.Nil(t, p.Pool[1])
	p.Config.MaxMessageLength = max

	// Send a malformed message, where ConvertToMessage fails
	// This is an unknown Message ID
	c.readLoopDone <- false
	c.writeLoopDone <- false
	t.Logf("Pushing message with unknown ID")
	p.Pool[1] = c
	c.Buffer.Reset()
	c.Buffer.Write([]byte{4, 0, 0, 0, 'Y', 'Y', 'Y', 'Z'})
	assert.NotPanics(t, p.processConnectionBuffers)
	assert.True(t, checkDisconnect(DisconnectMalformedMessage))
	assert.Nil(t, p.Pool[1])

	p.StopListen()
}

func TestConnectionWriteLoop(t *testing.T) {
	resetHandler()
	sendByteMessage = noopSendByteMessage
	encodeMessage = noopEncodeMessage
	p, b, c := setupTwoConnections()
	go p.ConnectionWriteLoop(b)
	go p.ConnectionWriteLoop(c)
	wait()

	m := NewByteMessage(88)
	// Send a successful message to b
	p.SendMessage(b, m)
	wait()
	if len(p.SendResults) == 0 {
		t.Fatalf("No send results, would block")
	}
	sr := <-p.SendResults
	assert.Equal(t, sr.Message, m)
	assert.Equal(t, sr.Connection, b)
	assert.Nil(t, sr.Error)
	assert.False(t, b.LastSent.IsZero())
	// Send a failed message to c
	sendByteMessage = failingSendByteMessage
	p.SendMessage(c, m)
	wait()
	if len(p.SendResults) == 0 {
		t.Fatalf("No send results, would block")
	}
	sr = <-p.SendResults
	assert.Equal(t, sr.Message, m)
	assert.Equal(t, sr.Connection, c)
	assert.NotNil(t, sr.Error)
	assert.True(t, c.LastSent.IsZero())
	if len(p.DisconnectQueue) == 0 {
		t.Fatalf("DisconnectQueue empty, would block")
	}
	de := <-p.DisconnectQueue
	assert.Equal(t, de.Reason, DisconnectWriteFailed)
	assert.Equal(t, de.ConnId, c.Id)
}

func TestHandleMessages(t *testing.T) {
	cleanupNet()
	cfg := NewConfig()
	cfg.Port = uint16(port)
	cfg.Address = address
	p := NewConnectionPool(cfg, nil)
	assert.NotPanics(t, p.HandleMessages)
}

func TestPoolSendMessage(t *testing.T) {
	resetHandler()
	c := &Connection{LastSent: time.Time{}, WriteQueue: make(chan Message, 1)}
	cfg := NewConfig()
	cfg.WriteTimeout = time.Second
	p := NewConnectionPool(cfg, nil)
	assert.NotEqual(t, p.Config.ConnectionWriteQueueSize, 0)
	sendByteMessage = noopSendByteMessage
	encodeMessage = noopEncodeMessage
	m := NewByteMessage(88)
	p.SendMessage(c, m)
	assert.Equal(t, len(c.WriteQueue), 1)
	if len(c.WriteQueue) == 0 {
		t.Fatal("c.WriteQueue empty, would block")
	}
	m2 := <-c.WriteQueue
	assert.Equal(t, m, m2)

	// queue full
	for i := 0; i < cap(c.WriteQueue); i++ {
		c.WriteQueue <- m
	}
	p.Pool[c.Id] = c
	p.SendMessage(c, m)
	assert.Equal(t, len(p.DisconnectQueue), 1)
	if len(p.DisconnectQueue) == 0 {
		t.Fatal("DisconnectQueue empty, would block")
	}
	de := <-p.DisconnectQueue
	assert.Equal(t, de.Reason, DisconnectWriteQueueFull)
	assert.Equal(t, de.ConnId, c.Id)
}

func TestPoolBroadcastMessage(t *testing.T) {
	resetHandler()
	p, b, c := setupTwoConnections()
	sendByteMessage = noopSendByteMessage
	encodeMessage = noopEncodeMessage
	m := NewByteMessage(88)
	p.BroadcastMessage(m)
	wait()
	assert.Equal(t, len(b.WriteQueue), 1)
	assert.Equal(t, len(c.WriteQueue), 1)
	if len(b.WriteQueue) == 0 {
		t.Fatal("b.WriteQueue empty, would block")
	}
	if len(c.WriteQueue) == 0 {
		t.Fatal("c.WriteQueue empty, would block")
	}
	m2 := <-b.WriteQueue
	assert.Equal(t, m, m2)
	m2 = <-c.WriteQueue
	assert.Equal(t, m, m2)
}

func TestPoolReceiveMessage(t *testing.T) {
	EraseMessages()
	resetHandler()
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
	err, reason := p.receiveMessage(c, b)
	assert.Nil(t, err)
	assert.False(t, c.LastReceived.IsZero())
	assert.Nil(t, reason)

	// Invalid byte message received
	b = []byte{1}
	err, reason = p.receiveMessage(c, b)
	assert.NotNil(t, err)
	assert.Nil(t, reason)

	// Valid message, but handler returns a DisconnectReason
	b = make([]byte, 0)
	b = append(b, ErrorPrefix[:]...)
	err, reason = p.receiveMessage(c, b)
	assert.Nil(t, err)
	assert.NotNil(t, reason)
	assert.Equal(t, reason.Error(), "Bad")
}

/* Helpers */

func wait() {
	time.Sleep(time.Millisecond * 50)
}

func newPoolConnection(p *ConnectionPool) (c *Connection) {
	go p.NewConnection(newNetConn())
	for {
		c = p.handleConnectionQueue()
		if c != nil {
			break
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	return
}

func handleXConnections(p *ConnectionPool, x int) {
	n := 0
	for {
		if p.handleConnectionQueue() != nil {
			n++
			logger.Critical("Got a connection off queue")
		}
		if n >= x {
			break
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	logger.Critical("Got %d queue connections", n)
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

func makeConnection(id int, addr string) *Connection {
	c := &Connection{
		Id:            id,
		Conn:          NewDummyConn(addr),
		LastSent:      time.Time{},
		LastReceived:  time.Time{},
		Buffer:        &bytes.Buffer{},
		WriteQueue:    make(chan Message, 10),
		readLoopDone:  make(chan bool, 1),
		writeLoopDone: make(chan bool, 1),
	}
	c.readLoopDone <- false
	c.writeLoopDone <- false
	return c
}

func setupTwoConnections() (*ConnectionPool, *Connection, *Connection) {
	b := makeConnection(1, addr)
	c := makeConnection(2, addrb)
	p := NewConnectionPool(NewConfig(), nil)
	p.Pool[b.Id] = b
	p.Pool[c.Id] = c
	p.Addresses[b.Conn.RemoteAddr().String()] = b
	p.Addresses[c.Conn.RemoteAddr().String()] = c
	return p, b, c
}
