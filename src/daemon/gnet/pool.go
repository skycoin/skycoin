/*
Package gnet is the core networking library
*/
package gnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	"io"

	"github.com/sirupsen/logrus"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/SkycoinProject/skycoin/src/daemon/strand"
	"github.com/SkycoinProject/skycoin/src/util/elapse"
	"github.com/SkycoinProject/skycoin/src/util/logging"
)

// DisconnectReason is passed to ConnectionPool's DisconnectCallback
type DisconnectReason error

const (
	receiveMessageDurationThreshold = 500 * time.Millisecond
	readLoopDurationThreshold       = 10 * time.Second
	sendInMsgChanDurationThreshold  = 5 * time.Second
	sendLoopDurationThreshold       = 500 * time.Millisecond
)

var (
	// ErrDisconnectSetReadDeadlineFailed set read deadline failed
	ErrDisconnectSetReadDeadlineFailed DisconnectReason = errors.New("SetReadDeadline failed")
	// ErrDisconnectInvalidMessageLength invalid message length
	ErrDisconnectInvalidMessageLength DisconnectReason = errors.New("Invalid message length")
	// ErrDisconnectMalformedMessage malformed message
	ErrDisconnectMalformedMessage DisconnectReason = errors.New("Malformed message body")
	// ErrDisconnectUnknownMessage unknown message
	ErrDisconnectUnknownMessage DisconnectReason = errors.New("Unknown message ID")
	// ErrDisconnectShutdown shutting down the client
	ErrDisconnectShutdown DisconnectReason = errors.New("Shutdown")
	// ErrDisconnectMessageDecodeUnderflow message data did not fully decode to a message object
	ErrDisconnectMessageDecodeUnderflow DisconnectReason = errors.New("Message data did not fully decode to a message object")
	// ErrDisconnectTruncatedMessageID message data was too short to contain a message ID
	ErrDisconnectTruncatedMessageID DisconnectReason = errors.New("Message data was too short to contain a message ID")

	// ErrConnectionPoolClosed error message indicates the connection pool is closed
	ErrConnectionPoolClosed = errors.New("Connection pool is closed")
	// ErrWriteQueueFull write queue is full
	ErrWriteQueueFull = errors.New("Write queue full")
	// ErrNoReachableConnections when broadcasting a message, no connections were available to send a message to
	ErrNoReachableConnections = errors.New("All pool connections are unreachable at this time")
	// ErrNoMatchingConnections when broadcasting a message, no connections were found for the provided addresses
	ErrNoMatchingConnections = errors.New("No connections found for broadcast addresses")
	// ErrPoolEmpty when broadcasting a message, the connection pool was empty
	ErrPoolEmpty = errors.New("Connection pool is empty after filtering connections")
	// ErrConnectionExists connection exists
	ErrConnectionExists = errors.New("Connection exists")
	// ErrMaxIncomingConnectionsReached max incoming connections reached
	ErrMaxIncomingConnectionsReached = errors.New("Max incoming connections reached")
	// ErrMaxOutgoingConnectionsReached max outgoing connections reached
	ErrMaxOutgoingConnectionsReached = errors.New("Max outgoing connections reached")
	// ErrMaxOutgoingDefaultConnectionsReached max outgoing default connections reached
	ErrMaxOutgoingDefaultConnectionsReached = errors.New("Max outgoing default connections reached")
	// ErrNoAddresses no addresses were provided to BroadcastMessage
	ErrNoAddresses = errors.New("No addresses provided")

	// Logger
	logger = logging.MustGetLogger("gnet")
)

// ReadError connection read error
type ReadError struct {
	Err error
}

func (e ReadError) Error() string {
	return fmt.Sprintf("read failed: %v", e.Err)
}

// WriteError connection read error
type WriteError struct {
	Err error
}

func (e WriteError) Error() string {
	return fmt.Sprintf("write failed: %v", e.Err)
}

// Config gnet config
type Config struct {
	// Address to listen on. Leave empty for arbitrary assignment
	Address string
	// Port to listen on. Set to 0 for arbitrary assignment
	Port uint16
	// Maximum total connections. Must be >= MaxOutgoingConnections + MaxIncomingConnections.
	MaxConnections int
	// Maximum outgoing connections
	MaxOutgoingConnections int
	// Maximum incoming connections
	MaxIncomingConnections int
	// Maximum allowed default outgoing connection number
	MaxDefaultPeerOutgoingConnections int
	// Messages greater than length are rejected and the sender disconnected
	MaxIncomingMessageLength int
	// Messages greater than length are not sent and an error is reported in a SendResult
	MaxOutgoingMessageLength int
	// Timeout is the timeout for dialing new connections.  Use a
	// timeout of 0 to ignore timeout.
	DialTimeout time.Duration
	// Timeout for reading from a connection. Set to 0 to default to the
	// system's timeout
	ReadTimeout time.Duration
	// Timeout for writing to a connection. Set to 0 to default to the
	// system's timeout
	WriteTimeout time.Duration
	// Message sent event buffers
	SendResultsSize int
	// Individual connections' send queue size.  This should be increased
	// if send volume per connection is high, so as not to block
	ConnectionWriteQueueSize int
	// Triggered on client disconnect
	DisconnectCallback DisconnectCallback
	// Triggered on client connect
	ConnectCallback ConnectCallback
	// Triggered on client connect failure
	ConnectFailureCallback ConnectFailureCallback
	// Print debug logs
	DebugPrint bool
	// Default "trusted" peers
	DefaultConnections []string
	// Default connections map
	defaultConnections map[string]struct{}
}

// NewConfig returns a Config with defaults set
func NewConfig() Config {
	return Config{
		Address:                           "",
		Port:                              0,
		MaxConnections:                    128,
		MaxOutgoingMessageLength:          256 * 1024,
		MaxIncomingMessageLength:          1024 * 1024,
		MaxDefaultPeerOutgoingConnections: 2,
		DialTimeout:                       time.Second * 30,
		ReadTimeout:                       time.Second * 30,
		WriteTimeout:                      time.Second * 30,
		SendResultsSize:                   2048,
		ConnectionWriteQueueSize:          128,
		DisconnectCallback:                nil,
		ConnectCallback:                   nil,
		DebugPrint:                        false,
		defaultConnections:                make(map[string]struct{}),
	}
}

const (
	// Byte size of the length prefix in message, sizeof(int32)
	messageLengthPrefixSize = 4
)

// Connection is stored by the ConnectionPool
type Connection struct {
	// Key in ConnectionPool.Pool
	ID uint64
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
func NewConnection(pool *ConnectionPool, id uint64, conn net.Conn, writeQueueSize int, solicited bool) *Connection {
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
func (conn *Connection) Close() error {
	err := conn.Conn.Close()
	close(conn.WriteQueue)
	conn.Buffer = &bytes.Buffer{}
	return err
}

// DisconnectCallback triggered on client disconnect
type DisconnectCallback func(addr string, id uint64, reason DisconnectReason)

// ConnectCallback triggered on client connect
type ConnectCallback func(addr string, id uint64, solicited bool)

// ConnectFailureCallback trigger on client connect failure
type ConnectFailureCallback func(addr string, solicited bool, err error)

// ConnectionPool connection pool
type ConnectionPool struct {
	// Configuration parameters
	Config Config
	// Channel for async message sending
	SendResults chan SendResult
	// All connections, indexed by ConnId
	pool map[uint64]*Connection
	// All connections, indexed by address
	addresses map[string]*Connection
	// connected default peer connections
	defaultOutgoingConnections map[string]struct{}
	// connected outgoing connections
	outgoingConnections map[string]struct{}
	// connected incoming connections
	incomingConnections map[string]struct{}
	// User-defined state to be passed into message handlers
	messageState interface{}
	// Connection ID counter
	connID uint64
	// Listening connection
	listener     net.Listener
	listenerLock sync.Mutex
	// operations channel
	reqC chan strand.Request
	// quit channel
	quit       chan struct{}
	done       chan struct{}
	strandDone chan struct{}
	wg         sync.WaitGroup
}

// NewConnectionPool creates a new ConnectionPool that will listen on
// Config.Port upon StartListen. State is an application defined object that
// will be passed to a Message's Handle().
func NewConnectionPool(c Config, state interface{}) (*ConnectionPool, error) {
	for _, p := range c.DefaultConnections {
		c.defaultConnections[p] = struct{}{}
	}
	if c.MaxConnections < c.MaxOutgoingConnections+c.MaxIncomingConnections {
		return nil, errors.New("MaxConnections must be >= MaxOutgoingConnections + MaxIncomingConnections")
	}

	return &ConnectionPool{
		Config:                     c,
		pool:                       make(map[uint64]*Connection),
		addresses:                  make(map[string]*Connection),
		defaultOutgoingConnections: make(map[string]struct{}),
		outgoingConnections:        make(map[string]struct{}),
		incomingConnections:        make(map[string]struct{}),
		SendResults:                make(chan SendResult, c.SendResultsSize),
		messageState:               state,
		quit:                       make(chan struct{}),
		done:                       make(chan struct{}),
		strandDone:                 make(chan struct{}),
		reqC:                       make(chan strand.Request),
	}, nil
}

// Run starts the connection pool
func (pool *ConnectionPool) Run() error {
	defer close(pool.done)
	defer logger.Info("Connection pool closed")

	// The strand processing goroutine must be started before any error can be
	// returned from Run(), otherwise the Shutdown() call will block if an error occurred
	pool.wg.Add(1)
	go func() {
		defer pool.wg.Done()
		pool.processStrand()
	}()

	// start the connection accept loop
	addr := fmt.Sprintf("%s:%v", pool.Config.Address, pool.Config.Port)
	logger.Infof("Listening for connections on %s...", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	pool.listenerLock.Lock()
	pool.listener = ln
	pool.listenerLock.Unlock()

loop:
	for {
		conn, err := ln.Accept()
		if err != nil {
			// When Accept() returns with a non-nil error, we check the quit
			// channel to see if we should continue or quit
			select {
			case <-pool.quit:
				break loop
			default:
				// without the default case the select will block.
				logger.Error(err.Error())
				continue
			}
		}

		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			if err := pool.handleConnection(conn, false); err != nil {
				logger.WithFields(logrus.Fields{
					"addr":     conn.RemoteAddr(),
					"outgoing": false,
				}).WithError(err).Error("pool.handleConnection")
			}
		}()
	}
	pool.wg.Wait()
	return nil
}

// RunOffline runs the pool in offline mode. No connections will be accepted,
// but strand requests are processed.
func (pool *ConnectionPool) RunOffline() error {
	defer close(pool.done)
	pool.processStrand()
	return nil
}

func (pool *ConnectionPool) processStrand() {
	defer close(pool.strandDone)
	for {
		select {
		case <-pool.quit:
			return
		case req := <-pool.reqC:
			if err := req.Func(); err != nil {
				logger.WithField("operation", req.Name).WithError(err).Error("strand req.Func failed")
			}
		}
	}
}

// Shutdown gracefully shutdown the connection pool
func (pool *ConnectionPool) Shutdown() {
	logger.Info("ConnectionPool.Shutdown called")
	close(pool.quit)
	logger.Info("ConnectionPool.Shutdown closed pool.quit")

	// Wait for all strand() calls to finish
	logger.Info("ConnectionPool.Shutdown waiting for strandDone")
	<-pool.strandDone

	logger.Info("ConnectionPool.Shutdown closing the listener")

	// Close to listener to prevent new connections
	pool.listenerLock.Lock()
	if pool.listener != nil {
		if err := pool.listener.Close(); err != nil {
			logger.WithError(err).Warning("pool.listener.Close error")
		}
	}
	pool.listener = nil
	pool.listenerLock.Unlock()

	logger.Info("ConnectionPool.Shutdown disconnecting all connections")

	// In readData, reader.Read() sometimes blocks instead of returning an error when the
	// listener is closed.
	// Directly close all connections before closing the listener.
	// TODO -- could conn.Close() block too?
	pool.disconnectAll()

	if len(pool.pool) != 0 {
		logger.Critical().Warning("pool.pool is not empty after calling pool.disconnectAll()")
	}
	if len(pool.addresses) != 0 {
		logger.Critical().Warning("pool.addresses is not empty after calling pool.disconnectAll()")
	}

	logger.Info("ConnectionPool.Shutdown waiting for done")

	<-pool.done
}

// strand ensures all read and write action of pool's member variable are in one thread
func (pool *ConnectionPool) strand(name string, f func() error) error {
	name = fmt.Sprintf("daemon.gnet.ConnectionPool.%s", name)
	return strand.Strand(logger, pool.reqC, name, f, pool.quit, ErrConnectionPoolClosed)
}

// ListeningAddress returns the address on which the ConnectionPool listens on
func (pool *ConnectionPool) ListeningAddress() (net.Addr, error) {
	if pool.listener == nil {
		return nil, errors.New("Not listening, call StartListen first")
	}
	return pool.listener.Addr(), nil
}

func (pool *ConnectionPool) canConnect(a string, solicited bool) error {
	if pool.isConnExist(a) {
		return ErrConnectionExists
	}

	if solicited {
		if _, ok := pool.Config.defaultConnections[a]; ok && pool.IsMaxOutgoingDefaultConnectionsReached() {
			return ErrMaxOutgoingDefaultConnectionsReached
		} else if pool.isMaxOutgoingConnectionsReached() {
			return ErrMaxOutgoingConnectionsReached
		}
	} else if pool.isMaxIncomingConnectionsReached() {
		return ErrMaxIncomingConnectionsReached
	}

	return nil
}

// newConnection creates a new Connection around a net.Conn. Trying to make a connection
// to an address that is already connected will failed.
func (pool *ConnectionPool) newConnection(conn net.Conn, solicited bool) (*Connection, error) {
	a := conn.RemoteAddr().String()

	if err := pool.canConnect(a, solicited); err != nil {
		return nil, err
	}

	if solicited {
		pool.outgoingConnections[a] = struct{}{}

		if _, ok := pool.Config.defaultConnections[a]; ok {
			pool.defaultOutgoingConnections[a] = struct{}{}
			l := len(pool.defaultOutgoingConnections)
			logger.WithField("addr", a).Debugf("%d/%d outgoing default connections in use", l, pool.Config.MaxDefaultPeerOutgoingConnections)
		}
	} else {
		pool.incomingConnections[a] = struct{}{}
	}

	// ID must start at 1; in case connID overflows back to 0, force it to 1
	pool.connID++
	if pool.connID == 0 {
		pool.connID = 1
	}

	nc := NewConnection(pool, pool.connID, conn, pool.Config.ConnectionWriteQueueSize, solicited)

	pool.pool[nc.ID] = nc
	pool.addresses[a] = nc

	return nc, nil
}

// Creates a Connection and begins its read and write loop
func (pool *ConnectionPool) handleConnection(conn net.Conn, solicited bool) error {
	defer logger.WithField("addr", conn.RemoteAddr()).Debug("Connection closed")
	addr := conn.RemoteAddr().String()

	c, err := func() (c *Connection, err error) {
		// TODO -- when limits in newConnection() are reached, should we allow the peer
		// to be added anyway, so that we can disconnect it normally and send a disconnect packet?
		// It would have to allowed to complete the handshake, otherwise the DISC packet will be ignored
		// Or we would have to permit a DISC packet before an INTR
		// But the read/write loop would still need to be started
		// A ConnectEvent would need to be triggered, or else the DisconnectEvent gnet ID will not match the
		// pending connection's zero gnet ID.
		defer func() {
			if err != nil {
				if closeErr := conn.Close(); closeErr != nil {
					logger.WithError(closeErr).WithField("addr", addr).Error("handleConnection conn.Close")
				}
			}
		}()

		err = pool.strand("handleConnection", func() error {
			var err error
			c, err = pool.newConnection(conn, solicited)
			if err != nil {
				return err
			}

			if pool.Config.ConnectCallback != nil {
				pool.Config.ConnectCallback(c.Addr(), c.ID, solicited)
			}

			return nil
		})

		return
	}()

	// TODO -- this error is not fully propagated back to a caller of Connect() so the daemon state
	// can get stuck in pending
	if err != nil {
		logger.WithError(err).WithField("addr", conn.RemoteAddr()).Debug("handleConnection: newConnection failed")
		if pool.Config.ConnectFailureCallback != nil {
			pool.Config.ConnectFailureCallback(addr, solicited, err)
		}
		return err
	}

	msgC := make(chan []byte, 32)

	type methodErr struct {
		method string
		err    error
	}
	errC := make(chan methodErr, 3)

	var wg sync.WaitGroup
	wg.Add(1)
	qc := make(chan struct{})
	go func() {
		defer wg.Done()
		if err := pool.readLoop(c, msgC, qc); err != nil {
			errC <- methodErr{
				method: "readLoop",
				err:    err,
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := pool.sendLoop(c, pool.Config.WriteTimeout, pool.Config.MaxOutgoingMessageLength, qc); err != nil {
			errC <- methodErr{
				method: "sendLoop",
				err:    err,
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		elapser := elapse.NewElapser(receiveMessageDurationThreshold, logger)
		defer elapser.CheckForDone()

		for msg := range msgC {
			elapser.Register(fmt.Sprintf("pool.receiveMessage address=%s", addr))
			if err := pool.receiveMessage(c, msg); err != nil {
				errC <- methodErr{
					method: "receiveMessage",
					err:    err,
				}
				return
			}
			elapser.CheckForDone()
		}
	}()

	select {
	case <-pool.quit:
		if err := conn.Close(); err != nil {
			logger.WithError(err).WithField("addr", addr).Error("conn.Close")
		}
	case mErr := <-errC:
		err = mErr.err
		logger.WithError(mErr.err).WithFields(logrus.Fields{
			"addr":   addr,
			"method": mErr.method,
		}).Error("handleConnection failure")

		// This Disconnect does not send a DISC packet because it is inside gnet.
		// A DISC packet is not useful at this point, because the error is most likely
		// that the connection is unreachable.
		// However, it may also be that the connection sent data that could not be deserialized
		// to a message.
		if err := pool.Disconnect(c.Addr(), mErr.err); err != nil {
			logger.WithError(err).WithField("addr", addr).Error("Disconnect")
		}
	}
	close(qc)

	wg.Wait()

	return err
}

func (pool *ConnectionPool) readLoop(conn *Connection, msgChan chan []byte, qc chan struct{}) error {
	defer close(msgChan)
	// read data from connection
	reader := bufio.NewReader(conn.Conn)
	buf := make([]byte, 1024)

	elapser := elapse.NewElapser(readLoopDurationThreshold, logger)
	sendInMsgChanElapser := elapse.NewElapser(sendInMsgChanDurationThreshold, logger)

	defer elapser.CheckForDone()
	defer sendInMsgChanElapser.CheckForDone()

	for {
		elapser.Register(fmt.Sprintf("readLoop addr=%s", conn.Addr()))
		deadline := time.Time{}
		if pool.Config.ReadTimeout != 0 {
			deadline = time.Now().Add(pool.Config.ReadTimeout)
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

		// write data to buffer
		if _, err := conn.Buffer.Write(data); err != nil {
			return err
		}
		// decode data
		datas, err := decodeData(conn.Buffer, pool.Config.MaxIncomingMessageLength)
		if err != nil {
			return err
		}
		for _, d := range datas {
			// use select to avoid the goroutine leak,
			// because if msgChan has no receiver this goroutine will leak
			select {
			case <-qc:
				return nil
			case <-pool.quit:
				return nil
			case msgChan <- d:
			default:
				return errors.New("readLoop msgChan is closed or full")
			}
		}
		sendInMsgChanElapser.CheckForDone()
	}
}

func (pool *ConnectionPool) sendLoop(conn *Connection, timeout time.Duration, maxMsgLength int, qc chan struct{}) error {
	elapser := elapse.NewElapser(sendLoopDurationThreshold, logger)
	defer elapser.CheckForDone()

	for {
		elapser.CheckForDone()
		select {
		case <-pool.quit:
			return nil
		case <-qc:
			return nil
		case m := <-conn.WriteQueue:
			elapser.Register(fmt.Sprintf("conn.WriteQueue address=%s", conn.Addr()))
			if m == nil {
				continue
			}

			err := sendMessage(conn.Conn, m, timeout, maxMsgLength)

			// Update last sent before writing to SendResult,
			// this allows a write to SendResult to be used as a sync marker,
			// since no further action in this block will happen after the write.
			if err == nil {
				if err := pool.updateLastSent(conn.Addr(), Now()); err != nil {
					logger.WithField("addr", conn.Addr()).WithError(err).Warning("updateLastSent failed")
				}
			}

			sr := newSendResult(conn.Addr(), m, err)
			select {
			case <-qc:
				return nil
			case pool.SendResults <- sr:
			default:
				logger.WithField("addr", conn.Addr()).Warning("SendResults queue full")
			}

			if err != nil {
				return err
			}
		}
	}
}

func readData(reader io.Reader, buf []byte) ([]byte, error) {
	c, err := reader.Read(buf)
	if err != nil {
		return nil, &ReadError{
			Err: err,
		}
	}
	if c == 0 {
		return nil, nil
	}
	data := make([]byte, c)
	copy(data, buf)
	return data, nil
}

// decode data from buffer.
func decodeData(buf *bytes.Buffer, maxMsgLength int) ([][]byte, error) {
	dataArray := [][]byte{}
	for buf.Len() > messageLengthPrefixSize {
		prefix := buf.Bytes()[:messageLengthPrefixSize]
		// decode message length
		tmpLength, _, err := encoder.DeserializeUint32(prefix)
		if err != nil {
			// encoder.DeserializeAtomic should only return an error if there wasn't
			// enough data in buf to read the integer, but the prefix buf length
			// is already ensured to be long enough
			logger.Panicf("encoder.DeserializeUint32 failed unexpectedly: %v", err)
		}

		length := int(tmpLength)

		// Disconnect if we received an invalid length
		if length < messagePrefixLength {
			logger.WithFields(logrus.Fields{
				"length":              length,
				"messagePrefixLength": messagePrefixLength,
			}).Warningf("decodeData: length < messagePrefixLength")
			return [][]byte{}, ErrDisconnectInvalidMessageLength
		}

		if length > maxMsgLength {
			logger.WithFields(logrus.Fields{
				"length":       length,
				"maxMsgLength": maxMsgLength,
			}).Warning("decodeData: length > maxMsgLength")
			return [][]byte{}, ErrDisconnectInvalidMessageLength
		}

		if buf.Len()-messageLengthPrefixSize < length {
			return [][]byte{}, nil
		}

		buf.Next(messageLengthPrefixSize) // strip the length prefix
		data := make([]byte, length)
		_, err = buf.Read(data)
		if err != nil {
			return [][]byte{}, err
		}

		dataArray = append(dataArray, data)
	}
	return dataArray, nil
}

// isConnExist check if the connection of address does exist
func (pool *ConnectionPool) isConnExist(addr string) bool {
	_, ok := pool.addresses[addr]
	return ok
}

func (pool *ConnectionPool) isMaxIncomingConnectionsReached() bool {
	return len(pool.incomingConnections) >= pool.Config.MaxIncomingConnections
}

func (pool *ConnectionPool) isMaxOutgoingConnectionsReached() bool {
	return len(pool.outgoingConnections) >= pool.Config.MaxOutgoingConnections
}

// IsMaxOutgoingDefaultConnectionsReached checks whether the max outgoing default connections reached
func (pool *ConnectionPool) IsMaxOutgoingDefaultConnectionsReached() bool {
	return len(pool.defaultOutgoingConnections) >= pool.Config.MaxDefaultPeerOutgoingConnections
}

func (pool *ConnectionPool) updateLastSent(addr string, t time.Time) error {
	return pool.strand("updateLastSent", func() error {
		if conn, ok := pool.addresses[addr]; ok {
			conn.LastSent = t
		}
		return nil
	})
}

func (pool *ConnectionPool) updateLastRecv(addr string, t time.Time) error {
	return pool.strand("updateLastRecv", func() error {
		if conn, ok := pool.addresses[addr]; ok {
			conn.LastReceived = t
		}
		return nil
	})
}

// GetConnection returns a connection copy if exist
func (pool *ConnectionPool) GetConnection(addr string) (*Connection, error) {
	var conn *Connection
	if err := pool.strand("GetConnection", func() error {
		if c, ok := pool.addresses[addr]; ok {
			// copy connection
			cc := *c
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
	if err := pool.strand("canConnect", func() error {
		return pool.canConnect(address, true)
	}); err != nil {
		return err
	}

	logger.WithField("addr", address).Debugf("Making TCP connection")
	conn, err := net.DialTimeout("tcp", address, pool.Config.DialTimeout)
	if err != nil {
		return err
	}

	pool.wg.Add(1)
	go func() {
		defer pool.wg.Done()
		if err := pool.handleConnection(conn, true); err != nil {
			logger.WithFields(logrus.Fields{
				"addr":     conn.RemoteAddr(),
				"outgoing": true,
			}).WithError(err).Error("pool.handleConnection")
		}
	}()
	return nil
}

// Disconnect removes a connection from the pool by address and invokes DisconnectCallback
func (pool *ConnectionPool) Disconnect(addr string, r DisconnectReason) error {
	return pool.strand("Disconnect", func() error {
		logger.WithFields(logrus.Fields{
			"addr":   addr,
			"reason": r,
		}).Debug("Disconnecting")

		// checks if the address is default node address
		isDefaultOutgoingConn := false
		if _, ok := pool.Config.defaultConnections[addr]; ok {
			if _, ok := pool.outgoingConnections[addr]; ok {
				isDefaultOutgoingConn = true
			}
		}

		conn := pool.disconnect(addr, r)

		if conn == nil {
			return errors.New("Disconnect: connection does not exist")
		}

		if isDefaultOutgoingConn {
			l := len(pool.defaultOutgoingConnections)
			logger.Debugf("%d/%d default connections in use", l, pool.Config.MaxDefaultPeerOutgoingConnections)
		}

		if pool.Config.DisconnectCallback != nil {
			pool.Config.DisconnectCallback(addr, conn.ID, r)
		}

		return nil
	})
}

func (pool *ConnectionPool) disconnect(addr string, r DisconnectReason) *Connection {
	conn, ok := pool.addresses[addr]
	if !ok {
		return nil
	}

	fields := logrus.Fields{
		"addr": addr,
		"id":   conn.ID,
	}

	delete(pool.pool, conn.ID)
	delete(pool.addresses, addr)
	delete(pool.defaultOutgoingConnections, addr)
	delete(pool.outgoingConnections, addr)
	delete(pool.incomingConnections, addr)
	if err := conn.Close(); err != nil {
		logger.WithError(err).WithFields(fields).Error("conn.Close")
	}

	logger.WithFields(fields).WithField("reason", r).Debug("Closed connection and removed from pool")

	return conn
}

// disconnectAll disconnects all connections. Only safe to call in Shutdown()
func (pool *ConnectionPool) disconnectAll() {
	for _, conn := range pool.pool {
		addr := conn.Addr()
		pool.disconnect(addr, ErrDisconnectShutdown)
	}
}

// GetConnections returns an copy of pool connections
func (pool *ConnectionPool) GetConnections() ([]Connection, error) {
	conns := []Connection{}
	if err := pool.strand("GetConnections", func() error {
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
	err = pool.strand("Size", func() error {
		l = len(pool.pool)
		return nil
	})
	return
}

// SendMessage sends a Message to a Connection
func (pool *ConnectionPool) SendMessage(addr string, msg Message) error {
	if pool.Config.DebugPrint {
		logger.WithField("msgType", reflect.TypeOf(msg)).Debug("SendMessage")
	}

	return pool.strand("SendMessage", func() error {
		if conn, ok := pool.addresses[addr]; ok {
			select {
			case conn.WriteQueue <- msg:
			default:
				logger.Critical().WithField("addr", addr).Info("Write queue full")
				return ErrWriteQueueFull
			}
		} else {
			return fmt.Errorf("Tried to send %T to %s, but we are not connected", msg, addr)
		}
		return nil
	})
}

// BroadcastMessage sends a Message to all connections specified in addrs.
// If a connection does not exist for a given address, it is skipped.
// If no messages were written to any connection, an error is returned.
// Returns the gnet IDs of connections that the message was queued for sending to.
// Note that actual sending can still fail later, if the connection drops before the message is sent.
func (pool *ConnectionPool) BroadcastMessage(msg Message, addrs []string) ([]uint64, error) {
	if pool.Config.DebugPrint {
		logger.WithField("msgType", reflect.TypeOf(msg)).Debug("BroadcastMessage")
	}

	if len(addrs) == 0 {
		return nil, ErrNoAddresses
	}

	queuedConns := make([]uint64, 0, len(addrs))

	if err := pool.strand("BroadcastMessage", func() error {
		if len(pool.pool) == 0 {
			return ErrPoolEmpty
		}

		fullWriteQueue := 0
		foundConns := 0

		for _, addr := range addrs {
			if conn, ok := pool.addresses[addr]; ok {
				foundConns++
				select {
				case conn.WriteQueue <- msg:
					queuedConns = append(queuedConns, conn.ID)
				default:
					logger.Critical().WithFields(logrus.Fields{
						"addr": conn.Addr(),
						"id":   conn.ID,
					}).Info("Write queue full")
					fullWriteQueue++
				}
			}
		}

		if foundConns == 0 {
			return ErrNoMatchingConnections
		}

		if fullWriteQueue == foundConns {
			return ErrNoReachableConnections
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return queuedConns, nil
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
	now := time.Now().UTC()
	var addrs []string
	if err := pool.strand("SendPings", func() error {
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

// GetStaleConnections returns connections that have been idle for longer than idleLimit
func (pool *ConnectionPool) GetStaleConnections(idleLimit time.Duration) ([]string, error) {
	now := Now()
	var idleConns []string
	if err := pool.strand("GetStaleConnections", func() error {
		for _, conn := range pool.pool {
			if conn.LastReceived.Add(idleLimit).Before(now) {
				idleConns = append(idleConns, conn.Addr())
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return idleConns, nil
}

// Now returns the current UTC time
func Now() time.Time {
	return time.Now().UTC()
}
