package gnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"io"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util"
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

	// Logger
	logger = util.MustGetLogger("gnet")
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
	// member variables access channel
	memChannel chan func()
	// quit channel
	Quit chan struct{}
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
		memChannel:   make(chan func()),
	}

	return pool
}

// Run starts the connection pool
func (pool *ConnectionPool) Run() {
	pool.Quit = make(chan struct{})
	go func() {
		for {
			select {
			case memActionFunc := <-pool.memChannel:
				// this goroutine will handle all member variable's writing and reading actions.
				memActionFunc()
			case <-pool.Quit:
				return
			}
		}
	}()

	// start the connection accept loop
	addr := fmt.Sprintf("%s:%v", pool.Config.Address, pool.Config.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	pool.listener = ln

	go func() {
		for {
			logger.Info("Listening for connections...")
			conn, err := ln.Accept()
			if err != nil {
				// When Accept() returns with a non-nill error, we check the quit
				// channel to see if we should continue or quit. If quit, then we quit.
				// Otherwise we continue
				select {
				case <-pool.Quit:
					return
				default:
					// without the default case the select will block.
				}

				continue
			}

			go pool.handleConnection(conn, false)
		}
	}()
}

// Shutdown gracefully shutdown the connection pool
func (pool *ConnectionPool) Shutdown() {
	close(pool.Quit)
	pool.listener.Close()
	pool.listener = nil
	pool.pool = make(map[int]*Connection)
	pool.addresses = make(map[string]*Connection)
}

// strand ensures all read and write action of pool's member variable are in one thread.
func (pool *ConnectionPool) strand(f func()) {
	q := make(chan struct{})
	pool.memChannel <- func() {
		defer close(q)
		f()
	}
	<-q
}

// NewConnection creates a new Connection around a net.Conn.  Trying to make a connection
// to an address that is already connected will panic.
func (pool *ConnectionPool) NewConnection(conn net.Conn, solicited bool) (*Connection, error) {
	a := conn.RemoteAddr().String()
	var nc *Connection
	var err error
	pool.strand(func() {
		if pool.addresses[a] != nil {
			err = fmt.Errorf("Already connected to %s", a)
			return
		}
		pool.connID++
		nc = NewConnection(pool, pool.connID, conn,
			pool.Config.ConnectionWriteQueueSize, solicited)

		pool.pool[nc.ID] = nc
		pool.addresses[a] = nc
	})

	return nc, err
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
	a := conn.RemoteAddr().String()
	if pool.IsConnExist(a) {
		log.Panicf("Connection %s already exists", a)
	}

	var c *Connection
	reason := ErrDisconnectUnexpectedError
	defer func() {
		logger.Debug("End connection handler of %s", conn.RemoteAddr())
		// notify to exist the receive message loop
		if c != nil {
			pool.Disconnect(c.Addr(), reason)
		}
	}()

	c, err := pool.NewConnection(conn, solicited)
	if err != nil {
		log.Panic(err)
	}

	if pool.Config.ConnectCallback != nil {
		pool.Config.ConnectCallback(c.Addr(), solicited)
	}

	msgChan := make(chan []byte, 10)
	errChan := make(chan error)
	go readLoop(c, pool.Config.ReadTimeout, pool.Config.MaxMessageLength, msgChan, errChan)
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
				reason = ErrDisconnectWriteFailed
				return
			}
			pool.updateLastSent(c.Addr(), Now())
		case msg := <-msgChan:
			dc, err := pool.receiveMessage(c, msg)
			if err != nil {
				reason = ErrDisconnectMalformedMessage
				return
			}
			if dc != nil {
				reason = dc
				return
			}
		case err := <-errChan:
			if err != nil {
				reason = err
				return
			}
		}
	}
}

func readLoop(conn *Connection, timeout time.Duration, maxMsgLen int, msgChan chan []byte, errChan chan error) {
	// read data from connection
	defer func() {
		logger.Debug("End readLoop of %s", conn.Addr())
	}()
	reader := bufio.NewReader(conn.Conn)
	buf := make([]byte, 1024)
	var rerr error
	for {
		deadline := time.Time{}
		if timeout != 0 {
			deadline = time.Now().Add(timeout)
		}
		if err := conn.Conn.SetReadDeadline(deadline); err != nil {
			rerr = ErrDisconnectSetReadDeadlineFailed
			break
		}

		data, err := readData(reader, buf)
		if err != nil {
			rerr = err
			break
		}

		if data == nil {
			continue
		}

		// write date to buffer.
		conn.Buffer.Write(data)

		// decode data
		datas, err := decodeData(conn.Buffer, maxMsgLen)
		if err != nil {
			rerr = err
			break
		}

		for _, d := range datas {
			// use select to avoid the goroutine leak, cause if msgChan has no receiver, this goroutine
			// will leak
			select {
			case msgChan <- d:
			default:
				return
			}
		}
	}

	if rerr != nil {
		select {
		case errChan <- rerr:
		default:
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
		log.Panic("Failed to copy all the bytes")
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
func (pool *ConnectionPool) IsConnExist(addr string) bool {
	var exist bool
	pool.strand(func() {
		if _, ok := pool.addresses[addr]; ok {
			exist = true
		}
	})
	return exist
}

func (pool *ConnectionPool) updateLastSent(addr string, t time.Time) {
	pool.strand(func() {
		if conn, ok := pool.addresses[addr]; ok {
			conn.LastSent = t
		}
	})
}

func (pool *ConnectionPool) updateLastRecv(addr string, t time.Time) {
	pool.strand(func() {
		if conn, ok := pool.addresses[addr]; ok {
			conn.LastReceived = t
		}
	})
}

// GetConnection returns a connection copy if exist
func (pool *ConnectionPool) GetConnection(addr string) *Connection {
	var conn Connection
	var exist bool
	pool.strand(func() {
		if c, ok := pool.addresses[addr]; ok {
			// copy connection
			conn = *c
			exist = true
		}
	})
	if exist {
		return &conn
	}
	return nil
}

// Connect to an address
func (pool *ConnectionPool) Connect(address string) error {
	if pool.IsConnExist(address) {
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
func (pool *ConnectionPool) Disconnect(addr string, r DisconnectReason) {
	var exist bool
	pool.strand(func() {
		if conn, ok := pool.addresses[addr]; ok {
			exist = true
			delete(pool.pool, conn.ID)
			delete(pool.addresses, addr)
			conn.Close()
		}
	})

	if pool.Config.DisconnectCallback != nil && exist {
		pool.Config.DisconnectCallback(addr, r)
	}
}

// GetConnections returns an copy of pool connections
func (pool *ConnectionPool) GetConnections() []Connection {
	conns := []Connection{}
	pool.strand(func() {
		for _, conn := range pool.pool {
			conns = append(conns, *conn)
		}
	})
	return conns
}

// Size returns the pool size
func (pool *ConnectionPool) Size() int {
	var l int
	pool.strand(func() {
		l = len(pool.pool)
	})
	return l
}

// SendMessage sends a Message to a Connection and pushes the result onto the
// SendResults channel.
func (pool *ConnectionPool) SendMessage(addr string, msg Message) error {
	if pool.Config.DebugPrint {
		logger.Debug("Send, Msg Type: %s", reflect.TypeOf(msg))
	}
	var msgQueueFull bool
	pool.strand(func() {
		if conn, ok := pool.addresses[addr]; ok {
			select {
			case conn.WriteQueue <- msg:
			default:
				msgQueueFull = true
			}
		}
	})

	if msgQueueFull {
		return ErrDisconnectWriteQueueFull
	}

	return nil
}

// BroadcastMessage sends a Message to all connections in the Pool.
func (pool *ConnectionPool) BroadcastMessage(msg Message) (err error) {
	if pool.Config.DebugPrint {
		logger.Debug("Broadcast, Msg Type: %s", reflect.TypeOf(msg))
	}

	fullWriteQueue := []string{}
	pool.strand(func() {
		if len(pool.pool) == 0 {
			err = errors.New("Connection pool is empty")
			return
		}

		for _, conn := range pool.pool {
			select {
			case conn.WriteQueue <- msg:
			case <-time.After(5 * time.Second):
				fullWriteQueue = append(fullWriteQueue, conn.Addr())
			}
		}
		if len(fullWriteQueue) == len(pool.pool) {
			err = errors.New("There's no available connection in pool")
		}
	})

	for _, addr := range fullWriteQueue {
		pool.Disconnect(addr, ErrDisconnectWriteQueueFull)
	}
	return
}

// Unpacks incoming bytes to a Message and calls the message handler.  If
// the bytes cannot be converted to a Message, the error is returned as the
// first return value.  Otherwise, error will be nil and DisconnectReason will
// be the value returned from the message handler.
func (pool *ConnectionPool) receiveMessage(c *Connection, msg []byte) (DisconnectReason, error) {
	m, err := convertToMessage(c.ID, msg, pool.Config.DebugPrint)
	if err != nil {
		return nil, err
	}
	pool.updateLastRecv(c.Addr(), Now())
	return m.Handle(NewMessageContext(c), pool.messageState), nil
}

// SendPings sends a ping if our last message sent was over pingRate ago
func (pool *ConnectionPool) SendPings(rate time.Duration, msg Message) {
	now := util.Now()
	var addrs []string
	pool.strand(func() {
		for _, conn := range pool.pool {
			if conn.LastSent.Add(rate).Before(now) {
				addrs = append(addrs, conn.Addr())
			}
		}
	})

	for _, a := range addrs {
		pool.SendMessage(a, msg)
	}
}

// ClearStaleConnections removes connections that have not sent a message in too long
func (pool *ConnectionPool) ClearStaleConnections(idleLimit time.Duration, reason DisconnectReason) {
	now := Now()
	idleConns := []string{}
	pool.strand(func() {
		for _, conn := range pool.pool {
			if conn.LastReceived.Add(idleLimit).Before(now) {
				idleConns = append(idleConns, conn.Addr())
			}
		}
	})

	for _, a := range idleConns {
		pool.Disconnect(a, reason)
	}
}

// Now returns the current UTC time
func Now() time.Time {
	return time.Now().UTC()
}
