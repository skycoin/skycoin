package gnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"

	"io"

	"github.com/skycoin/skycoin/src/cipher/encoder"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
)

// DisconnectReason is passed to ConnectionPool's DisconnectCallback
type DisconnectReason error

var (
	// ErrDisconnectReadFailed also includes a remote closed socket
	ErrDisconnectReadFailed DisconnectReason = errors.New("Read failed")
	// ErrDisconnectWriteFailed write faile
	ErrDisconnectWriteFailed DisconnectReason = errors.New("Write failed")
	// ErrDisconnectSetReadDeadlineFailed set read deadline failed
	ErrDisconnectSetReadDeadlineFailed = errors.New("SetReadDeadline failed")
	// ErrDisconnectInvalidMessageLength invalid message length
	ErrDisconnectInvalidMessageLength DisconnectReason = errors.New("Invalid message length")
	// ErrDisconnectMalformedMessage malformed message
	ErrDisconnectMalformedMessage DisconnectReason = errors.New("Malformed message body")
	// ErrDisconnectUnknownMessage unknow message
	ErrDisconnectUnknownMessage DisconnectReason = errors.New("Unknown message ID")
	// ErrDisconnectWriteQueueFull write queue is full
	ErrDisconnectWriteQueueFull DisconnectReason = errors.New("Write queue full")
	// ErrDisconnectUnexpectedError  unexpected error
	ErrDisconnectUnexpectedError DisconnectReason = errors.New("Unexpected error encountered")
	// ErrConnectionPoolClosed error message indicates the connection pool is closed
	ErrConnectionPoolClosed = errors.New("Connection pool is closed")
	// Logger
	logger = logging.MustGetLogger("gnet")
)

// Config gnet config
type Config struct {
	// Address to listen on. Leave empty for arbitrary assignment
	Address string
	// Port to listen on. Set to 0 for arbitrary assignment
	Port uint16
	// Connection limits
	MaxConnections int
	// Messages greater than length are rejected and the sender disconnected
	MaxMessageLength int
	// Timeout is the timeout for dialing new connections.  Use a
	// timeout of 0 to ignore timeout.
	DialTimeout time.Duration
	// Timeout for reading from a connection. Set to 0 to default to the
	// system's timeout
	ReadTimeout time.Duration
	// Timeout for writing to a connection. Set to 0 to default to the
	// system's timeout
	WriteTimeout time.Duration
	// Broadcast result buffers
	BroadcastResultSize int
	// Individual connections' send queue size.  This should be increased
	// if send volume per connection is high, so as not to block
	ConnectionWriteQueueSize int
	// Triggered on client disconnect
	DisconnectCallback DisconnectCallback
	// Triggered on client connect
	ConnectCallback ConnectCallback
	// Print debug logs
	DebugPrint bool
}

// NewConfig returns a Config with defaults set
func NewConfig() Config {
	return Config{
		Address:          "",
		Port:             0,
		MaxConnections:   128,
		MaxMessageLength: 256 * 1024,
		DialTimeout:      time.Minute,
		ReadTimeout:      time.Minute,
		WriteTimeout:     time.Minute,
		// EventChannelSize:         4096,
		BroadcastResultSize:      16,
		ConnectionWriteQueueSize: 32,
		DisconnectCallback:       nil,
		ConnectCallback:          nil,
		DebugPrint:               false,
	}
}

const (
	// Byte size of the length prefix in message, sizeof(int32)
	messageLengthSize = 4
)

// Connection is stored by the ConnectionPool
type Connection struct {
	// Key in ConnectionPool.Pool
	ID int
	// TCP connection
	Conn net.Conn
	// Message buffer
	Buffer *bytes.Buffer
	// Reference back to ConnectionPool container
	ConnectionPool *ConnectionPool
	// Last time a message was fully parsed and handled
	LastReceived time.Time
	// Last time a message was sent to the connection
	LastSent time.Time
	// Message send queue.
	WriteQueue chan Message
	Solicited  bool
}

// NewConnection creates a new Connection tied to a ConnectionPool
func NewConnection(pool *ConnectionPool, id int, conn net.Conn, writeQueueSize int, solicited bool) *Connection {
	return &Connection{
		ID:             id,
		Conn:           conn,
		Buffer:         &bytes.Buffer{},
		ConnectionPool: pool,
		LastReceived:   Now(),
		LastSent:       Now(),
		WriteQueue:     make(chan Message, writeQueueSize),
		Solicited:      solicited,
	}
}

// Addr returns remote address
func (conn *Connection) Addr() string {
	return conn.Conn.RemoteAddr().String()
}

// String returns connection address
func (conn *Connection) String() string {
	return conn.Addr()
}

// Close close the connection and write queue
func (conn *Connection) Close() {
	conn.Conn.Close()
	close(conn.WriteQueue)
	conn.WriteQueue = nil
	conn.Buffer = &bytes.Buffer{}
}

// DisconnectCallback triggered on client disconnect
type DisconnectCallback func(addr string, reason DisconnectReason)

// ConnectCallback triggered on client connect
type ConnectCallback func(addr string, solicited bool)

// ConnectionPool connection pool
type ConnectionPool struct {
	// Configuration parameters
	Config Config
	// Channel for async message sending
	SendResults chan SendResult
	// All connections, indexed by ConnId
	pool map[int]*Connection
	// All connections, indexed by address
	addresses map[string]*Connection
	// User-defined state to be passed into message handlers
	messageState interface{}
	// Connection ID counter
	connID int
	// Listening connection
	listener net.Listener
	// operations channel
	ops chan func()
	// quit channel
	quit chan struct{}
}

// NewConnectionPool creates a new ConnectionPool that will listen on Config.Port upon
// StartListen.  State is an application defined object that will be
// passed to a Message's Handle().
func NewConnectionPool(c Config, state interface{}) *ConnectionPool {
	pool := &ConnectionPool{
		Config:       c,
		pool:         make(map[int]*Connection),
		addresses:    make(map[string]*Connection),
		SendResults:  make(chan SendResult, c.BroadcastResultSize),
		messageState: state,
	}

	return pool
}

// Run starts the connection pool
func (pool *ConnectionPool) Run() error {
	// init the quit and operations channel here, in case run this pool again.
	pool.quit = make(chan struct{})
	pool.ops = make(chan func())

	go func() {
		for op := range pool.ops {
			op()
		}

		logger.Info("Connection pool closed")
	}()

	// start the connection accept loop
	addr := fmt.Sprintf("%s:%v", pool.Config.Address, pool.Config.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	pool.listener = ln

	logger.Info("Listening for connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			// When Accept() returns with a non-nill error, we check the quit
			// channel to see if we should continue or quit . If quit, then we quit.
			// Otherwise we continue
			select {
			case <-pool.quit:
				close(pool.ops)
				return nil
			default:
				// without the default case the select will block.
				logger.Error("%v", err)
				continue
			}
		}

		go pool.handleConnection(conn, false)
	}
}

// Shutdown gracefully shutdown the connection pool
func (pool *ConnectionPool) Shutdown() {
	pool.strand(func() error {
		pool.addresses = map[string]*Connection{}
		pool.pool = map[int]*Connection{}
		return nil
	})

	close(pool.quit)

	if pool.listener != nil {
		pool.listener.Close()
	}

	pool.listener = nil
}

// strand ensures all read and write action of pool's member variable are in one thread.
func (pool *ConnectionPool) strand(f func() error) (err error) {
	defer func() {
		// send on closed operation channel will panic.
		if r := recover(); r != nil {
			err = ErrConnectionPoolClosed
		}
	}()

	q := make(chan struct{})
	pool.ops <- func() {
		defer close(q)
		err = f()
	}
	<-q
	return
}

// NewConnection creates a new Connection around a net.Conn.  Trying to make a connection
// to an address that is already connected will failed.
func (pool *ConnectionPool) NewConnection(conn net.Conn, solicited bool) (*Connection, error) {
	a := conn.RemoteAddr().String()
	var nc *Connection
	if err := pool.strand(func() error {
		if pool.addresses[a] != nil {
			return fmt.Errorf("Already connected to %s", a)
		}
		pool.connID++
		nc = NewConnection(pool, pool.connID, conn,
			pool.Config.ConnectionWriteQueueSize, solicited)

		pool.pool[nc.ID] = nc
		pool.addresses[a] = nc
		return nil
	}); err != nil {
		return nil, err
	}

	return nc, nil
}

// ListeningAddress returns address, on which the ConnectionPool
// listening on. It returns nil, and error if the ConnectionPool
// is not listening
func (pool *ConnectionPool) ListeningAddress() (net.Addr, error) {
	if pool.listener == nil {
		return nil, errors.New("Not listening, call StartListen first")
	}
	return pool.listener.Addr(), nil
}

// Creates a Connection and begins its read and write loop
func (pool *ConnectionPool) handleConnection(conn net.Conn, solicited bool) {
	addr := conn.RemoteAddr().String()
	exist, err := pool.IsConnExist(addr)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	if exist {
		logger.Error("Connection %s already exists", addr)
		return
	}

	c, err := pool.NewConnection(conn, solicited)
	if err != nil {
		logger.Error("Create connection failed: %v", err)
		return
	}

	if pool.Config.ConnectCallback != nil {
		pool.Config.ConnectCallback(c.Addr(), solicited)
	}

	msgC := make(chan []byte, 10)
	errC := make(chan error, 1)

	go func() {
		errC <- readLoop(c, pool.Config.ReadTimeout, pool.Config.MaxMessageLength, msgC)
	}()

	qc := make(chan chan struct{})
	go func() {
		for {
			select {
			case m := <-c.WriteQueue:
				if m == nil {
					continue
				}
				err := sendMessage(conn, m, pool.Config.WriteTimeout)
				sr := newSendResult(c.Addr(), m, err)
				pool.SendResults <- sr
				if err != nil {
					errC <- err
					return
				}

				if err := pool.updateLastSent(c.Addr(), Now()); err != nil {
					errC <- err
					return
				}
			case msg := <-msgC:
				if err := pool.receiveMessage(c, msg); err != nil {
					errC <- err
					return
				}
			case q := <-qc:
				q <- struct{}{}
				return
			}
		}
	}()

	e := <-errC
	q := make(chan struct{}, 1)
	qc <- q
	<-q

	if err := pool.Disconnect(c.Addr(), e); err != nil {
		logger.Error("Disconnect failed: %v", err)
	}
}

func readLoop(conn *Connection, timeout time.Duration, maxMsgLen int, msgChan chan []byte) error {
	// read data from connection
	reader := bufio.NewReader(conn.Conn)
	buf := make([]byte, 1024)
	for {
		deadline := time.Time{}
		if timeout != 0 {
			deadline = time.Now().Add(timeout)
		}
		if err := conn.Conn.SetReadDeadline(deadline); err != nil {
			return ErrDisconnectSetReadDeadlineFailed
		}

		data, err := readData(reader, buf)
		if err != nil {
			return err
		}

		if data == nil {
			continue
		}

		// write date to buffer.
		if _, err := conn.Buffer.Write(data); err != nil {
			return err
		}

		// decode data
		datas, err := decodeData(conn.Buffer, maxMsgLen)
		if err != nil {
			return err
		}

		for _, d := range datas {
			// use select to avoid the goroutine leak, cause if msgChan has no receiver, this goroutine
			// will leak
			select {
			case msgChan <- d:
			default:
				return errors.New("The msgChan has no receiver")
			}
		}
	}
}

func readData(reader io.Reader, buf []byte) ([]byte, error) {
	c, err := reader.Read(buf)
	if err != nil {
		return nil, ErrDisconnectReadFailed
	}
	if c == 0 {
		return nil, nil
	}

	data := make([]byte, c)
	n := copy(data, buf)
	if n != c {
		// I don't believe this can ever occur
		return nil, errors.New("Failed to copy all the bytes")
	}
	return data, nil
}

// decode data from buffer.
func decodeData(buf *bytes.Buffer, maxMsgLength int) ([][]byte, error) {
	dataArray := [][]byte{}
	for buf.Len() > messageLengthSize {
		//logger.Debug("There is data in the buffer, extracting")
		prefix := buf.Bytes()[:messageLengthSize]
		// decode message length
		tmpLength := uint32(0)
		encoder.DeserializeAtomic(prefix, &tmpLength)
		length := int(tmpLength)
		// logger.Debug("Length is %d", length)
		// Disconnect if we received an invalid length.
		if length < messagePrefixLength ||
			length > maxMsgLength {
			return [][]byte{}, ErrDisconnectInvalidMessageLength
		}

		if buf.Len()-messageLengthSize < length {
			// logger.Debug("Skipping, not enough data to read this")
			return [][]byte{}, nil
		}

		buf.Next(messageLengthSize) // strip the length prefix
		data := make([]byte, length)
		_, err := buf.Read(data)
		if err != nil {
			return [][]byte{}, err
		}

		dataArray = append(dataArray, data)
	}
	return dataArray, nil
}

// IsConnExist check if the connection of address does exist
func (pool *ConnectionPool) IsConnExist(addr string) (bool, error) {
	var exist bool
	if err := pool.strand(func() error {
		if _, ok := pool.addresses[addr]; ok {
			exist = true
		}
		return nil
	}); err != nil {
		return false, fmt.Errorf("Check connection existence failed: %v ", err)
	}

	return exist, nil
}

func (pool *ConnectionPool) updateLastSent(addr string, t time.Time) error {
	return pool.strand(func() error {
		if conn, ok := pool.addresses[addr]; ok {
			conn.LastSent = t
		}
		return nil
	})
}

func (pool *ConnectionPool) updateLastRecv(addr string, t time.Time) error {
	return pool.strand(func() error {
		if conn, ok := pool.addresses[addr]; ok {
			conn.LastReceived = t
		}
		return nil
	})
}

// GetConnection returns a connection copy if exist
func (pool *ConnectionPool) GetConnection(addr string) (*Connection, error) {
	var conn *Connection
	if err := pool.strand(func() error {
		if c, ok := pool.addresses[addr]; ok {
			// copy connection
			var cc = *c
			conn = &cc
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return conn, nil
}

// Connect to an address
func (pool *ConnectionPool) Connect(address string) error {
	exist, err := pool.IsConnExist(address)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	logger.Debug("Making TCP Connection to %s", address)
	conn, err := net.DialTimeout("tcp", address, pool.Config.DialTimeout)
	if err != nil {
		return err
	}

	go pool.handleConnection(conn, true)
	return nil
}

// Disconnect removes a connection from the pool by address, and passes a Disconnection to
// the DisconnectCallback
func (pool *ConnectionPool) Disconnect(addr string, r DisconnectReason) error {
	var exist bool
	if err := pool.strand(func() error {
		if conn, ok := pool.addresses[addr]; ok {
			exist = true
			delete(pool.pool, conn.ID)
			delete(pool.addresses, addr)
			conn.Close()
		}
		return nil
	}); err != nil {
		return err
	}

	if pool.Config.DisconnectCallback != nil && exist {
		pool.Config.DisconnectCallback(addr, r)
	}
	return nil
}

// GetConnections returns an copy of pool connections
func (pool *ConnectionPool) GetConnections() ([]Connection, error) {
	conns := []Connection{}
	if err := pool.strand(func() error {
		for _, conn := range pool.pool {
			conns = append(conns, *conn)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return conns, nil
}

// Size returns the pool size
func (pool *ConnectionPool) Size() (l int, err error) {
	err = pool.strand(func() error {
		l = len(pool.pool)
		return nil
	})
	return
}

// SendMessage sends a Message to a Connection and pushes the result onto the
// SendResults channel.
func (pool *ConnectionPool) SendMessage(addr string, msg Message) error {
	if pool.Config.DebugPrint {
		logger.Debug("Send, Msg Type: %s", reflect.TypeOf(msg))
	}
	var msgQueueFull bool
	if err := pool.strand(func() error {
		if conn, ok := pool.addresses[addr]; ok {
			select {
			case conn.WriteQueue <- msg:
			default:
				msgQueueFull = true
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if msgQueueFull {
		return ErrDisconnectWriteQueueFull
	}

	return nil
}

// BroadcastMessage sends a Message to all connections in the Pool.
func (pool *ConnectionPool) BroadcastMessage(msg Message) error {
	if pool.Config.DebugPrint {
		logger.Debug("Broadcast, Msg Type: %s", reflect.TypeOf(msg))
	}

	fullWriteQueue := []string{}
	if err := pool.strand(func() error {
		if len(pool.pool) == 0 {
			return errors.New("Connection pool is empty")
		}

		for _, conn := range pool.pool {
			select {
			case conn.WriteQueue <- msg:
			case <-time.After(5 * time.Second):
				fullWriteQueue = append(fullWriteQueue, conn.Addr())
			}
		}
		if len(fullWriteQueue) == len(pool.pool) {
			return errors.New("There's no available connection in pool")
		}

		return nil
	}); err != nil {
		return err
	}

	for _, addr := range fullWriteQueue {
		if err := pool.Disconnect(addr, ErrDisconnectWriteQueueFull); err != nil {
			return err
		}
	}
	return nil
}

// Unpacks incoming bytes to a Message and calls the message handler.  If
// the bytes cannot be converted to a Message, the error is returned as the
// first return value.  Otherwise, error will be nil and DisconnectReason will
// be the value returned from the message handler.
func (pool *ConnectionPool) receiveMessage(c *Connection, msg []byte) error {
	m, err := convertToMessage(c.ID, msg, pool.Config.DebugPrint)
	if err != nil {
		return err
	}
	if err := pool.updateLastRecv(c.Addr(), Now()); err != nil {
		return err
	}
	return m.Handle(NewMessageContext(c), pool.messageState)
}

// SendPings sends a ping if our last message sent was over pingRate ago
func (pool *ConnectionPool) SendPings(rate time.Duration, msg Message) error {
	now := utc.Now()
	var addrs []string
	if err := pool.strand(func() error {
		for _, conn := range pool.pool {
			if conn.LastSent.Add(rate).Before(now) {
				addrs = append(addrs, conn.Addr())
			}
		}
		return nil
	}); err != nil {
		return err
	}

	for _, a := range addrs {
		if err := pool.SendMessage(a, msg); err != nil {
			return err
		}
	}

	return nil
}

// ClearStaleConnections removes connections that have not sent a message in too long
func (pool *ConnectionPool) ClearStaleConnections(idleLimit time.Duration, reason DisconnectReason) error {
	now := Now()
	idleConns := []string{}
	if err := pool.strand(func() error {
		for _, conn := range pool.pool {
			if conn.LastReceived.Add(idleLimit).Before(now) {
				idleConns = append(idleConns, conn.Addr())
			}
		}
		return nil
	}); err != nil {
		return err
	}

	for _, a := range idleConns {
		pool.Disconnect(a, reason)
	}
	return nil
}

// Now returns the current UTC time
func Now() time.Time {
	return utc.Now()
}
