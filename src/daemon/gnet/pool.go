package gnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util"
)

var DebugPrint bool = true //disable to disable printing

// TODO -- parameterize configuration per pool

// DisconnectReason is passed to ConnectionPool's DisconnectCallback
type DisconnectReason error

var (
	// DisconnectReadFailed also includes a remote closed socket
	DisconnectReadFailed DisconnectReason = errors.New(
		"Read failed")
	DisconnectWriteFailed DisconnectReason = errors.New(
		"Write failed")
	DisconnectSetReadDeadlineFailed = errors.New(
		"SetReadDeadline failed")
	DisconnectInvalidMessageLength DisconnectReason = errors.New(
		"Invalid message length")
	DisconnectMalformedMessage DisconnectReason = errors.New(
		"Malformed message body")
	DisconnectUnknownMessage DisconnectReason = errors.New(
		"Unknown message ID")
	DisconnectWriteQueueFull DisconnectReason = errors.New(
		"Write queue full")
	DisconnectUnexpectedError DisconnectReason = errors.New(
		"Unexpected error encountered")

	// Logger
	logger = util.MustGetLogger("gnet")
)

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
	// Event channel buffering
	EventChannelSize int
	// Broadcast result buffers
	BroadcastResultSize int
	// Individual connections' send queue size.  This should be increased
	// if send volume per connection is high, so as not to block
	ConnectionWriteQueueSize int
	// Triggered on client disconnect
	DisconnectCallback DisconnectCallback
	// Triggered on client connect
	ConnectCallback ConnectCallback
}

// Returns a Config with defaults set
func NewConfig() Config {
	return Config{
		Address:                  "",
		Port:                     0,
		MaxConnections:           128,
		MaxMessageLength:         256 * 1024,
		DialTimeout:              time.Minute,
		ReadTimeout:              time.Minute,
		WriteTimeout:             time.Minute,
		EventChannelSize:         4096,
		BroadcastResultSize:      16,
		ConnectionWriteQueueSize: 32,
		DisconnectCallback:       nil,
		ConnectCallback:          nil,
	}
}

const (
	// Byte size of the length prefix in message, sizeof(int32)
	messageLengthSize = 4
)

// Connection is stored by the ConnectionPool
type Connection struct {
	// Key in ConnectionPool.Pool
	Id int
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
	// Synchronizers for read/write loop termination
	writeLoopDone chan bool
	readLoopDone  chan bool
	// acquire and release to use raw incoming connection
	acquire chan struct{}
	// close acquire once
	acquireClose *sync.Once
}

// Creates a new Connection tied to a ConnectionPool
func NewConnection(pool *ConnectionPool, id int, conn net.Conn,
	writeQueueSize int) *Connection {
	c := &Connection{
		Id:             id,
		Conn:           conn,
		Buffer:         &bytes.Buffer{},
		ConnectionPool: pool,
		LastReceived:   Now(),
		LastSent:       Now(),
		WriteQueue:     make(chan Message, writeQueueSize),
		writeLoopDone:  make(chan bool, 1),
		readLoopDone:   make(chan bool, 1),
		acquire:        make(chan struct{}),
		acquireClose:   new(sync.Once),
	}
	c.Release() // make it released
	c.writeLoopDone <- false
	c.readLoopDone <- false
	return c
}

// Acquire can be used only from inside (*ConnectionPool).ConnectCallback.
// The role of the Acquire is to get raw incoming connection to perform some
// work (like raw handskae). You must to call Release when done. When you
// acquire a connection, the read-write loop of the connetion can't start
// before you relase it
func (self *Connection) Acquire() {
	// If you want to perform some work on an incoming connection but
	// if you do it inside ConnectCallback it will block execution flow.
	// The Accquire and Release method are used to prevent starting read- and
	// write-loops untill we release the connection. This way we can do the wrok
	// in separate goroutine, that doesn't block main execution flow

	// when acquire was blocked for reading,
	// the connection is acquired
	self.acquire = make(chan struct{})
	// close it once
	self.acquireClose = new(sync.Once)
}

// Release used after Acquire when you done. It's unsafe to call Release if you
// don't call Acquire. It's safe to call it many times. The design is following
//
//     myPoolConfig.ConnectCallback = func(gc *Connection, outgoing bool) {
//         if !outgoing { // if the connection is incoming
//             gc.Acquire()
//             go func() {
//                 defer gc.Release()
//                 performRawHandshake(gc)
//             } ()
//             return
//         }
//     }
//
// It's possible to use Acquire/Release for outgoing connections too. A
// connection you use will be closed when you call StopListen. If you call
// Acquire on some connection and then StopListen the connection will be
// released and closed
func (self *Connection) Release() {
	self.acquireClose.Do(func() {
		close(self.acquire)
	})
}

func (self *Connection) Addr() string {
	return self.Conn.RemoteAddr().String()
}

func (self *Connection) String() string {
	return self.Addr()
}

func (self *Connection) Close() {
	self.Conn.Close()
	self.Release() // release possible acquired connections
	if self.WriteQueue != nil {
		close(self.WriteQueue)
	}
	<-self.writeLoopDone
	<-self.readLoopDone

	self.Buffer.Reset()
	self.WriteQueue = nil
}

// Channel event for incoming socket data
type dataEvent struct {
	ConnId int    // Id for the Connection that produced this event
	Data   []byte // Data sent by the Connection
}

// Event generated when user disconnects
type DisconnectEvent struct {
	ConnId int
	Reason DisconnectReason
}

// Triggered on client disconnect
type DisconnectCallback func(c *Connection, reason DisconnectReason)

// Triggered on client connect
type ConnectCallback func(c *Connection, solicited bool)

type ConnectionPool struct {
	// Configuration parameters
	Config Config
	// All connections, indexed by ConnId
	Pool map[int]*Connection
	// All connections, indexed by address
	Addresses map[string]*Connection
	// Channel for buffered disconnects
	DisconnectQueue chan DisconnectEvent
	// Channel for async message sending
	SendResults chan SendResult
	// User-defined state to be passed into message handlers
	messageState interface{}
	// Connection ID counter
	connId int
	// Channel for this pool
	eventChannel chan dataEvent
	// Listening connection
	listener net.Listener
	// connectionQueue, for eliminating contention of the Pool/Addresses
	connectionQueue chan *Connection
	// channel for synchronizing teardown
	acceptLock chan bool
}

// Creates a new ConnectionPool that will listen on Config.Port upon
// StartListen.  State is an application defined object that will be
// passed to a Message's Handle().
func NewConnectionPool(c Config, state interface{}) *ConnectionPool {
	pool := &ConnectionPool{
		Config:          c,
		Pool:            make(map[int]*Connection),
		Addresses:       make(map[string]*Connection),
		DisconnectQueue: make(chan DisconnectEvent, c.MaxConnections),
		SendResults:     make(chan SendResult, c.BroadcastResultSize),
		eventChannel:    make(chan dataEvent, c.EventChannelSize),
		messageState:    state,
		// connectionQueue must not be buffered to guarantee the Pool is
		// updated before processing is done
		connectionQueue: make(chan *Connection),
		acceptLock:      make(chan bool, 1),
	}
	pool.acceptLock <- false
	return pool
}

// Creates a new Connection around a net.Conn.  Trying to make a connection
// to an address that is already connected will panic.
func (self *ConnectionPool) NewConnection(conn net.Conn) *Connection {
	a := conn.RemoteAddr().String()
	if self.Addresses[a] != nil {
		log.Panicf("Already connected to %s", a)
	}
	self.connId++
	c := NewConnection(self, self.connId, conn,
		self.Config.ConnectionWriteQueueSize)
	self.connectionQueue <- c
	return c
}

// Accepts incoming connections and sends them off for processing.
// This function blocks, you will want to run it in a goroutine.
func (self *ConnectionPool) AcceptConnections() {
	logger.Debug("Accepting connections...")
	if self.listener == nil {
		log.Panic("Not listening, call StartListen first")
	}
	defer func() {
		// If Accept fails for a reason other than StopListen closing listener,
		// this could block indefinitely
		self.acceptLock <- true
	}()
	// Grab the lock
	<-self.acceptLock
	// If we raced with StopListen and it beat us, listener will be nil and
	// the intent was for us to stop.
	if self.listener == nil {
		return
	}
	for {
		// this blocks until connection or error
		conn, err := self.listener.Accept()
		if err != nil {
			logger.Debug("No longer accepting connections")
			break
		}
		self.handleConnection(conn, false)
	}
}

// ListeningAddress returns address, on which the ConnectionPool
// listening on. It returns nil, and error if the ConnectionPool
// is not listening
func (self *ConnectionPool) ListeningAddress() (net.Addr, error) {
	if self.listener == nil {
		return nil, errors.New("Not listening, call StartListen first")
	}
	return self.listener.Addr(), nil
}

// Begin listening on the port the ConnectionPool is configured for.
// Calling StartListen() twice with no intermediate StopListen() will panic.
func (self *ConnectionPool) StartListen() error {
	if self.listener != nil {
		log.Panic("Already listening, call StopListen first")
	}
	addr := fmt.Sprintf("%s:%v", self.Config.Address, self.Config.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Debug("Listening on %s", addr)
	self.listener = ln
	return nil
}

// Stop accepting connections, close all connections, empty the Pool and
// reset all event channels
func (self *ConnectionPool) StopListen() {
	logger.Debug("Closing %d connections", len(self.Pool))
	if self.listener != nil {
		self.listener.Close()
	}
	// Synchronize with AcceptConnections:
	// Grab the acceptLock. Either:
	//  -Waits for AcceptConnections to stop due to listener being closed
	//  -AcceptConnections is not running at all
	//  -AcceptConnections is running but has not tried to grab its lock yet
	<-self.acceptLock
	// Set listener to nil now, in case we raced with AcceptConnections
	// and beat it to the lock.  We assume that in a race, StopListen
	// takes priority, since you shouldn't be calling StopListen in a
	// goroutine but you always call AcceptConnections in one.
	self.listener = nil
	// Unlock for future AcceptConnections
	self.acceptLock <- false
	for i, c := range self.Pool {
		c.Conn.Close()
		delete(self.Pool, i)
		delete(self.Addresses, c.Addr())
	}
	// Channels must be closed after connections are closed, so that their
	// RW goroutines can terminate.  Otherwise they might try to send on a
	// closed queue.
	close(self.eventChannel)
	close(self.DisconnectQueue)
	close(self.connectionQueue)
	self.connectionQueue = make(chan *Connection)
	self.eventChannel = make(chan dataEvent, self.Config.EventChannelSize)
	self.DisconnectQueue = make(chan DisconnectEvent,
		self.Config.MaxConnections)
}

// Creates a Connection and begins its read and write loop
func (self *ConnectionPool) handleConnection(conn net.Conn,
	solicited bool) *Connection {
	a := conn.RemoteAddr().String()
	c := self.Addresses[a]
	if c != nil {
		log.Panicf("Connection %s already exists", a)
	}
	c = self.NewConnection(conn)
	if self.Config.ConnectCallback != nil {
		self.Config.ConnectCallback(c, solicited)
	}
	go self.connectionReadLoop(c)
	go self.ConnectionWriteLoop(c)
	return c
}

// Connect to an address
func (self *ConnectionPool) Connect(address string) (*Connection, error) {
	if c := self.Addresses[address]; c != nil {
		return c, nil
	}
	logger.Debug("Making TCP Connection to %s", address)
	conn, err := net.DialTimeout("tcp", address, self.Config.DialTimeout)
	if err != nil {
		return nil, err
	}
	return self.handleConnection(conn, true), nil
}

// Removes a Connection from the pool, and passes a DisconnectReason to the
// DisconnectCallback
func (self *ConnectionPool) Disconnect(c *Connection, r DisconnectReason) {
	if c == nil {
		return
	}
	_, exists := self.Pool[c.Id]
	if !exists {
		return
	}
	self.DisconnectQueue <- DisconnectEvent{
		ConnId: c.Id,
		Reason: r,
	}
}

// Removes address found in a DisconnectEvent.  This needs to be called when
// your application pulls a DisconnectEvent off of the DisconnectQueue
func (self *ConnectionPool) HandleDisconnectEvent(e DisconnectEvent) {
	c := self.Pool[e.ConnId]
	if c == nil {
		return
	}
	c.Close()
	delete(self.Pool, c.Id)
	delete(self.Addresses, c.Addr())
	if self.Config.DisconnectCallback != nil {
		self.Config.DisconnectCallback(c, e.Reason)
	}
}

// Returns the pool as an array of *Connections
func (self *ConnectionPool) GetConnections() []*Connection {
	conns := make([]*Connection, len(self.Pool))
	i := 0
	for _, c := range self.Pool {
		conns[i] = c
		i++
	}
	return conns
}

// Returns the underlying net.Conn in the Pool
func (self *ConnectionPool) GetRawConnections() []net.Conn {
	conns := make([]net.Conn, len(self.Pool))
	i := 0
	for _, c := range self.Pool {
		conns[i] = c.Conn
		i++
	}
	return conns
}

// Writes message to a client socket in a goroutine.
// This is only public because its very helpful for testing applications
// that use this module.  Don't call it from non-test code.
func (self *ConnectionPool) ConnectionWriteLoop(c *Connection) {
	<-c.acquire
	<-c.writeLoopDone
	for {
		m, ok := <-c.WriteQueue
		if !ok {
			break
		}
		err := sendMessage(c.Conn, m, self.Config.WriteTimeout)
		sr := newSendResult(c, m, err)
		self.SendResults <- sr
		if err != nil {
			self.Disconnect(c, DisconnectWriteFailed)
			break
		}
		c.LastSent = Now()
	}
	c.writeLoopDone <- true
}

// Reads data from socket into channel in a goroutine for each connection
func (self *ConnectionPool) connectionReadLoop(conn *Connection) {
	<-conn.acquire
	<-conn.readLoopDone
	reason := DisconnectUnexpectedError
	defer func() {
		self.Disconnect(conn, reason)
	}()
	reader := bufio.NewReader(conn.Conn)
	buf := make([]byte, 1024)
	for {
		deadline := time.Time{}
		if self.Config.ReadTimeout != 0 {
			deadline = time.Now().Add(self.Config.ReadTimeout)
		}
		if err := conn.Conn.SetReadDeadline(deadline); err != nil {
			logger.Error("Failed to set read deadline for %s", conn.Addr())
			reason = DisconnectSetReadDeadlineFailed
			break
		}
		c, err := reader.Read(buf)
		if err != nil {
			logger.Debug("Failed to read from %s: %v", conn.Addr(), err)
			reason = DisconnectReadFailed
			break
		}
		if c == 0 {
			continue
		}
		//logger.Debug("Read %d bytes from %s", c, conn.Addr())

		data := make([]byte, c)
		n := copy(data, buf)
		if n != c {
			// I don't believe this can ever occur
			log.Panic("Failed to copy all the bytes")
		}
		// Write data to channel for processing
		self.eventChannel <- dataEvent{ConnId: conn.Id, Data: data}
	}
	conn.readLoopDone <- true
}

// Processes a new queued connection
func (self *ConnectionPool) handleConnectionQueue() *Connection {
	select {
	case c := <-self.connectionQueue:
		logger.Debug("New connection to %s pulled off connectionQueue",
			c.Addr())
		self.Pool[c.Id] = c
		self.Addresses[c.Addr()] = c
		return c
	default:
		return nil
	}
}

// Pulls all available event data from eventChannel and copies it to each
// client's buffer
// TODO -- this might not need to be separate, the client can maybe
// write to their buffer directly?
func (self *ConnectionPool) processEvents() {
	for len(self.eventChannel) > 0 {
		event := <-self.eventChannel
		c, worked := self.Pool[event.ConnId]
		if !worked {
			// The client must have been disconnected before we got to its
			// message here
			logger.Debug("No connection found for the event")
			continue
		}
		//logger.Debug("Got an event from %s", c.Addr())

		// There is no need to check this call's return values
		// "The return value n is the length of p; err is always nil.
		// If the buffer becomes too large, Write will panic with ErrTooLarge."
		n, _ := c.Buffer.Write(event.Data)
		_ = n
		//logger.Debug("Wrote %d bytes to the client buffer", n)
	}
}

// Converts a client's connection buffer to byte messages
func (self *ConnectionPool) processConnectionBuffer(c *Connection) {
	for c.Buffer.Len() > messageLengthSize {
		//logger.Debug("There is data in the buffer, extracting")
		prefix := c.Buffer.Bytes()[:messageLengthSize]
		// decode message length
		tmpLength := uint32(0)
		encoder.DeserializeAtomic(prefix, &tmpLength)
		length := int(tmpLength)
		//logger.Debug("Length is %d", length)
		// Disconnect if we received an invalid length.
		if length < messagePrefixLength ||
			length > self.Config.MaxMessageLength {
			logger.Debug("Invalid message length %d received from %s",
				length, c.Addr())
			self.Disconnect(c, DisconnectInvalidMessageLength)
			break
		}

		if c.Buffer.Len()-messageLengthSize < length {
			//logger.Debug("Skipping, not enough data to read this")
			break
		}

		c.Buffer.Next(messageLengthSize) // strip the length prefix
		data := c.Buffer.Next(length)    // read the message contents

		//logger.Debug("Telling the message unpacker about this message")
		err, dc := self.receiveMessage(c, data)
		if err != nil {
			logger.Debug("Error with the event: %v", err)
			self.Disconnect(c, DisconnectMalformedMessage)
			break
		}
		if dc != nil {
			// The handler disconnected the connection, stop processing
			break
		}
	}
}

// Converts clients' connection buffers to messages as data is available
func (self *ConnectionPool) processConnectionBuffers() {
	for _, c := range self.Pool {
		self.processConnectionBuffer(c)
	}
}

// Processes and clears pending messages
func (self *ConnectionPool) HandleMessages() {
	// Update the Pool for new connections. We do this here so that there is
	// no contention for RW access to the pool or addresses (assuming
	// HandleMessages() is in the same select as the DisconnectQueue
	self.handleConnectionQueue()

	// Copy event data to the client's buffer
	self.processEvents()

	// Process all messages from the client buffer
	self.processConnectionBuffers()
}

// Sends a Message to a Connection and pushes the result onto the
// SendResults channel.
func (self *ConnectionPool) SendMessage(c *Connection, msg Message) {

	if DebugPrint {
		logger.Debug("Send, Msg Type: %s", reflect.TypeOf(msg))
	}

	select {
	case c.WriteQueue <- msg:
	default:
		self.Disconnect(c, DisconnectWriteQueueFull)
	}
}

//var MessageIdMap = make(map[reflect.Type]MessagePrefix)
//t := reflect.TypeOf(msg)

// Sends a Message to all connections in the Pool.
func (self *ConnectionPool) BroadcastMessage(msg Message) {
	if DebugPrint {
		logger.Debug("Broadcast, Msg Type: %s", reflect.TypeOf(msg))
	}

	for _, c := range self.Pool {
		self.SendMessage(c, msg)
	}
}

// Unpacks incoming bytes to a Message and calls the message handler.  If
// the bytes cannot be converted to a Message, the error is returned as the
// first return value.  Otherwise, error will be nil and DisconnectReason will
// be the value returned from the message handler.
func (self *ConnectionPool) receiveMessage(c *Connection,
	msg []byte) (error, DisconnectReason) {
	m, err := convertToMessage(c, msg)
	if err != nil {
		return err, nil
	}
	c.LastReceived = Now()
	return nil, m.Handle(NewMessageContext(c), self.messageState)
	//return nil, m.

	//deprecate messageState

	/*
		For message prefix, look up the function
		- then call reflect.Call(m map[string]interface{}, name string, params ... interface{}) (result []reflect.Value, err error)
		- then convert the type

	*/

}

// Returns the current UTC time
func Now() time.Time {
	return time.Now().UTC()
}
