package gnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
	This is a connection pool for real/physical TCP/ip connections.

	It sends and receives length prefixed byte messages over a channel (uint16). A
	dispatcher Object handles the message serialization/deserialization and passes
	the message to the service object.

	Channel 0 is the control channel.

	Connection pool triggers ConnectCallback on client connect
	Connection pool trigger DisconnectCallback on client disconnection
	Connection pool triggers MessageCallback on receiving message
*/

/*
   The new connection pool sends length prefixed byte messages and receives
   length prefixed byte messages. Each message has a "channel".

   Each message may correspond to a UDP datagram or individual messages.

   - Channel zero is the "control channel"
   - Messages within a channel are ordered.
   - Messages within a channel must be received in the order they are sent.
   - Channel zero messages should be prioritized
   - Current implementation uses TCP but future will use UDP

    Design:
    - send length (uint32) prefixed messages with (uint16) channel prefix
    - all writes should go through pool
    - no writes into connection object
*/

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
	logger = logging.MustGetLogger("gnet")
)

// Triggered on client disconnect
type DisconnectCallback func(c *Connection, reason DisconnectReason)

// Triggered on client connect
type ConnectCallback func(c *Connection, solicited bool)

// Triggered on incoming message
type MessageCallback func(c *Connection, channel uint16, msg []byte) error

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
	// Triggered on client receiving data
	MessageCallback MessageCallback
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

/*
type ConnectionId {
    Id int
    Addr string
}
*/

//type BlobCallback func([]byte) BlobCallbackResponse

type BinMessage struct {
	Channel uint16
	Message []byte
}

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
	WriteQueue chan BinMessage
	// Synchronizers for read/write loop termination
	writeLoopDone chan bool
	readLoopDone  chan bool
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
		WriteQueue:     make(chan BinMessage, writeQueueSize),
		writeLoopDone:  make(chan bool, 1),
		readLoopDone:   make(chan bool, 1),
	}
	c.writeLoopDone <- false
	c.readLoopDone <- false
	return c
}

func (self *Connection) Addr() string {
	return self.Conn.RemoteAddr().String()
}

func (self *Connection) String() string {
	return self.Addr()
}

func (self *Connection) Close() {
	self.Conn.Close()
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
	//SendResults chan SendResult
	// User-defined state to be passed into message handlers
	//messageState interface{}
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

//receiveMessage(c *Connection, channel uint16, msg []byte) (error, DisconnectReason)

// Creates a new ConnectionPool that will listen on Config.Port upon
// StartListen.
func NewConnectionPool(c Config) *ConnectionPool {
	pool := &ConnectionPool{
		Config:          c,
		Pool:            make(map[int]*Connection),
		Addresses:       make(map[string]*Connection),
		DisconnectQueue: make(chan DisconnectEvent, c.MaxConnections),
		//SendResults:     make(chan SendResult, c.BroadcastResultSize),
		eventChannel: make(chan dataEvent, c.EventChannelSize),
		//messageState: state,
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
	//note, spawning two goroutine per connection...
	go self.connectionReadLoop(c)
	go self.ConnectionWriteLoop(c)
	return c
}

// Connect to an address. Blocks until connected
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

// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.

// Sends []byte over a net.Conn
// Why does this set deadline every packet
var sendByteMessage = func(conn net.Conn, channel uint16, msg []byte,
	timeout time.Duration) error {

	//message length and channel id
	bLen := encoder.SerializeAtomic(uint32(len(msg)))
	chanByte := encoder.SerializeAtomic(uint16(channel))

	//log.Printf("len= %v blen= %v chans= %v \n", len(msg), len(bLen), len(chanByte))
	d := make([]byte, 0, len(msg)+6)
	d = append(d, bLen...)     //length prefix
	d = append(d, chanByte...) //channel id
	d = append(d, msg...)      //message data

	//log.Printf("len2= %v \n ", len(d))

	deadline := time.Time{}
	if timeout != 0 {
		deadline = time.Now().Add(timeout)
	}
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return err
	}
	if _, err := conn.Write(d); err != nil {
		return err
	}
	return nil
}

// Writes message to a client socket in a goroutine.
// This is only public because its very helpful for testing applications
// that use this module.  Don't call it from non-test code.
// TODO: Why is there a write que?
func (self *ConnectionPool) ConnectionWriteLoop(c *Connection) {
	<-c.writeLoopDone
	for {
		m, ok := <-c.WriteQueue
		if !ok {
			break
		}
		//err := sendMessage(c.Conn, m, self.Config.WriteTimeout)
		//err := sendMessage(c.Conn, m.Message, self.Config.WriteTimeout)

		err := sendByteMessage(c.Conn, m.Channel, m.Message, self.Config.WriteTimeout)
		if err != nil {
			self.Disconnect(c, DisconnectWriteFailed)
			break
		}
		//sr := newSendResult(c, m)
		//self.SendResults <- sr

		c.LastSent = Now()
	}
	c.writeLoopDone <- true
}

// Reads data from socket into channel in a goroutine for each connection
func (self *ConnectionPool) connectionReadLoop(conn *Connection) {
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
		logger.Debug("Read %d bytes from %s", c, conn.Addr())

		data := make([]byte, c)
		n := copy(data, buf)
		if n != c {
			// I don't believe this can ever occur
			log.Panic("Failed to copy all the bytes")
		}
		// Write data to channel for processing
		self.eventChannel <- dataEvent{ConnId: conn.Id, Data: data}

		// Write data to buffer
		//n, _ = conn.Buffer.Write(data)
		logger.Debug("Received Data: addr= %s, %d bytes", conn.Addr(), len(data))

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
		// There is no need to check this call's return values
		// "The return value n is the length of p; err is always nil.
		// If the buffer becomes too large, Write will panic with ErrTooLarge."

		n, _ := c.Buffer.Write(event.Data)
		logger.Debug("Received Data: addr= %s, %d bytes", c.Addr(), n)
	}
}

const (
	// uint32 size prefix
	// uint16 channel id prefix
	messageLengthSize = 4 + 2
)

// Converts a client's connection buffer to byte messages
// Keep extracting message events until we dont have enough bytes to read in
func (self *ConnectionPool) processConnectionBuffer(c *Connection) {
	for c.Buffer.Len() >= messageLengthSize {
		//logger.Debug("There is data in the buffer, extracting")
		prefix := c.Buffer.Bytes()[:messageLengthSize]
		// decode message length
		tmpLength := uint32(0)
		encoder.DeserializeAtomic(prefix[0:4], &tmpLength)
		length := int(tmpLength)

		channel := uint16(0)
		encoder.DeserializeAtomic(prefix[4:6], &channel)

		logger.Debug("Extracting message: addr= %s, len(msg)=%d bytes", c.Addr(), length)
		// Disconnect if we received an invalid length.
		if length < messagePrefixLength ||
			length > self.Config.MaxMessageLength {
			logger.Debug("Invalid message length %d received from %s; message size mismatch",
				length, c.Addr())
			self.Disconnect(c, DisconnectInvalidMessageLength)
			break
		}

		if c.Buffer.Len()-messageLengthSize < length {
			logger.Debug("Skipping, not enough data to read this")
			break
		}

		c.Buffer.Next(messageLengthSize) // strip the length prefix
		data := c.Buffer.Next(length)    // read the message contents

		//logger.Debug("Telling the message unpacker about this message")
		c.LastReceived = Now()
		//err, dc := self.receiveMessage(c, channel, data)

		err := self.Config.MessageCallback(c, channel, data)

		if err != nil {
			logger.Debug("Error with the event: %v", err)
			self.Disconnect(c, DisconnectMalformedMessage)
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
	// incoming event data... write to buffer

	self.processEvents()

	// Process all messages from the client buffer
	// trigger callbacks from buffer data
	self.processConnectionBuffers()
}

// Sends a Message to a Connection and pushes the result onto the
// SendResults channel.
func (self *ConnectionPool) SendMessage(c *Connection, channel uint16, msg []byte) {
	select {
	case c.WriteQueue <- BinMessage{Channel: channel, Message: msg}:
	default:
		logger.Debug("Warning: disconnecting client because write queue full")
		self.Disconnect(c, DisconnectWriteQueueFull)
	}
}

// Sends a Message to all connections in the Pool.
func (self *ConnectionPool) BroadcastMessage(channel uint16, msg []byte) {
	for _, c := range self.Pool {
		self.SendMessage(c, channel, msg)
	}
}

// Returns the current UTC time
func Now() time.Time {
	return time.Now().UTC()
}

//stops listening and closes all connections
func (self *ConnectionPool) Shutdown() {
	self.StopListen() //have to do anything?
	for _, con := range self.Addresses {
		con.Close()
	}
}
