package daemon

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/op/go-logging.v1"

	//"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util"
)

/*
Todo
- verify that minimum/maximum connections are working
- keep max connections
- maintain minimum number of outgoing connections per server?


*/
var (
	// DisconnectReasons
	DisconnectInvalidVersion gnet.DisconnectReason = errors.New(
		"Invalid version")
	DisconnectIntroductionTimeout gnet.DisconnectReason = errors.New(
		"Version timeout")
	DisconnectVersionSendFailed gnet.DisconnectReason = errors.New(
		"Version send failed")
	DisconnectIsBlacklisted gnet.DisconnectReason = errors.New(
		"Blacklisted")
	DisconnectSelf gnet.DisconnectReason = errors.New(
		"Self connect")
	DisconnectConnectedTwice gnet.DisconnectReason = errors.New(
		"Already connected")
	DisconnectIdle gnet.DisconnectReason = errors.New(
		"Idle")
	DisconnectNoIntroduction gnet.DisconnectReason = errors.New(
		"First message was not an Introduction")
	DisconnectIPLimitReached gnet.DisconnectReason = errors.New(
		"Maximum number of connections for this IP was reached")
	// This is returned when a seemingly impossible error is encountered
	// e.g. net.Conn.Addr() returns an invalid ip:port
	DisconnectOtherError gnet.DisconnectReason = errors.New(
		"Incomprehensible error")

	//Use exponential backoff for connections
	//ConnectFailed gnet.DisconnectReason = errors.New(
	//	"Could Not Connect Error")

	// Blacklist a peer when they get disconnected for these
	// DisconnectReasons

	BlacklistOffenses = map[gnet.DisconnectReason]time.Duration{
		//DisconnectSelf:                      time.Second * 1,
		//DisconnectIntroductionTimeout:       time.Second * 1,
		DisconnectNoIntroduction:            time.Minute * 20,
		gnet.DisconnectInvalidMessageLength: time.Minute * 20,
		gnet.DisconnectMalformedMessage:     time.Minute * 20,
		gnet.DisconnectUnknownMessage:       time.Minute * 20,
		//ConnectFailed:                       time.Minute * 60,
	}

	logger = logging.MustGetLogger("skycoin.daemon")
)

// Subsystem configurations
type Config struct {
	Daemon   DaemonConfig
	Messages MessagesConfig
	Pool     PoolConfig
	Peers    PeersConfig
	Gateway  GatewayConfig
	Visor    VisorConfig
}

// Returns a Config with defaults set
func NewConfig() Config {
	return Config{
		Daemon:   NewDaemonConfig(),
		Pool:     NewPoolConfig(),
		Peers:    NewPeersConfig(),
		Gateway:  NewGatewayConfig(),
		Messages: NewMessagesConfig(),
		Visor:    NewVisorConfig(),
	}
}

func (self *Config) preprocess() Config {
	config := *self
	if config.Daemon.LocalhostOnly {
		if config.Daemon.Address == "" {
			local, err := LocalhostIP()
			if err != nil {
				log.Panicf("Failed to obtain localhost IP: %v", err)
			}
			config.Daemon.Address = local
		} else {
			if !IsLocalhost(config.Daemon.Address) {
				log.Panicf("Invalid address for localhost-only: %s",
					config.Daemon.Address)
			}
		}
		config.Peers.AllowLocalhost = true
	}
	config.Pool.port = config.Daemon.Port
	config.Pool.address = config.Daemon.Address

	if config.Daemon.DisableNetworking {
		config.Peers.Disabled = true
		config.Daemon.DisableIncomingConnections = true
		config.Daemon.DisableOutgoingConnections = true
	} else {
		if config.Daemon.DisableIncomingConnections {
			logger.Info("Incoming connections are disabled.")
		}
		if config.Daemon.DisableOutgoingConnections {
			logger.Info("Outgoing connections are disabled.")
		}
	}

	return config
}

// Configuration for the Daemon
type DaemonConfig struct {
	// Application version. TODO -- manage version better
	Version int32
	// IP Address to serve on. Leave empty for automatic assignment
	Address string
	// TCP/UDP port for connections
	Port int
	// Directory where application data is stored
	DataDirectory string
	// How often to check and initiate an outgoing connection if needed
	OutgoingRate time.Duration
	// How often to re-attempt to fill any missing private (aka required)
	// connections
	PrivateRate time.Duration
	// Number of outgoing connections to maintain
	OutgoingMax int
	// Maximum number of connections to try at once
	PendingMax int
	// How long to wait for a version packet
	IntroductionWait time.Duration
	// How often to check for peers that have decided to stop communicating
	CullInvalidRate time.Duration
	// How many connections are allowed from the same base IP
	IPCountsMax int
	// Disable all networking activity
	DisableNetworking bool
	// Don't make outgoing connections
	DisableOutgoingConnections bool
	// Don't allow incoming connections
	DisableIncomingConnections bool
	// Run on localhost and only connect to localhost peers
	LocalhostOnly bool
}

func NewDaemonConfig() DaemonConfig {
	return DaemonConfig{
		Version:                    2,
		Address:                    "",
		Port:                       6677,
		OutgoingRate:               time.Second * 5,
		PrivateRate:                time.Second * 5,
		OutgoingMax:                32,
		PendingMax:                 16,
		IntroductionWait:           time.Second * 30,
		CullInvalidRate:            time.Second * 3,
		IPCountsMax:                3,
		DisableNetworking:          false,
		DisableOutgoingConnections: false,
		DisableIncomingConnections: false,
		LocalhostOnly:              false,
	}
}

// Stateful properties of the daemon
type Daemon struct {
	// Daemon configuration
	Config DaemonConfig

	// Components
	Messages *Messages
	Pool     *Pool
	Peers    *Peers
	Gateway  *Gateway
	Visor    *Visor

	DefaultConnections []string

	// Separate index of outgoing connections. The pool aggregates all
	// connections.
	OutgoingConnections map[string]*gnet.Connection
	// Number of connections waiting to be formed or timeout
	pendingConnections map[string]*pex.Peer
	// Keep track of unsolicited clients who should notify us of their version
	ExpectingIntroductions map[string]time.Time
	// Keep track of a connection's mirror value, to avoid double
	// connections (one to their listener, and one to our listener)
	// Maps from addr to mirror value
	ConnectionMirrors map[string]uint32
	// Maps from mirror value to a map of ip (no port)
	// We use a map of ip as value because multiple peers can have the same
	// mirror (to avoid attacks enabled by our use of mirrors),
	// but only one per base ip
	mirrorConnections map[uint32]map[string]uint16
	// Client connection/disconnection callbacks
	onConnectEvent chan ConnectEvent
	// Connection failure events
	connectionErrors chan ConnectionError
	// Tracking connections from the same base IP.  Multiple connections
	// from the same base IP are allowed but limited.
	ipCounts map[string]int
	// Message handling queue
	messageEvents chan MessageEvent
}

// Returns a Daemon with primitives allocated
func NewDaemon(config Config) *Daemon {
	config = config.preprocess()
	d := &Daemon{
		Config:   config.Daemon,
		Messages: NewMessages(config.Messages),
		Pool:     NewPool(config.Pool),
		Peers:    NewPeers(config.Peers),
		Visor:    NewVisor(config.Visor),

		DefaultConnections: DefaultConnections, //passed in from top level

		ExpectingIntroductions: make(map[string]time.Time),
		ConnectionMirrors:      make(map[string]uint32),
		mirrorConnections:      make(map[uint32]map[string]uint16),
		ipCounts:               make(map[string]int),
		// TODO -- if there are performance problems from blocking chans,
		// Its because we are connecting to more things than OutgoingMax
		// if we have private peers
		onConnectEvent: make(chan ConnectEvent,
			config.Daemon.OutgoingMax),
		connectionErrors: make(chan ConnectionError,
			config.Daemon.OutgoingMax),
		OutgoingConnections: make(map[string]*gnet.Connection,
			config.Daemon.OutgoingMax),
		pendingConnections: make(map[string]*pex.Peer,
			config.Daemon.PendingMax),
		messageEvents: make(chan MessageEvent,
			config.Pool.EventChannelSize),
	}
	d.Gateway = NewGateway(config.Gateway, d)
	d.Messages.Config.Register()
	d.Pool.Init(d)
	d.Peers.Init()
	return d
}

// Generated when a client connects
type ConnectEvent struct {
	Addr      string
	Solicited bool
}

// Represent a failure to connect/dial a connection, with context
type ConnectionError struct {
	Addr  string
	Error error
}

// Encapsulates a deserialized message from the network
type MessageEvent struct {
	Message AsyncMessage
	Context *gnet.MessageContext
}

// Terminates all subsystems safely.  To stop the Daemon run loop, send a value
// over the quit channel provided to Init.  The Daemon run loop must be stopped
// before calling this function.
func (self *Daemon) Shutdown() {
	self.Pool.Shutdown()
	self.Peers.Shutdown()
	self.Visor.Shutdown()
	gnet.EraseMessages()
}

// Main loop for peer/connection management. Send anything to quit to shut it
// down
func (self *Daemon) Start(quit chan int) {
	if !self.Config.DisableIncomingConnections {
		self.Pool.StartListen() //no goroutine
		go self.Pool.AcceptConnections()
	}

	// TODO -- run blockchain stuff in its own goroutine
	blockInterval := time.Duration(self.Visor.Config.Config.BlockCreationInterval)
	// blockchainBackupTicker := time.Tick(self.Visor.Config.BlockchainBackupRate)
	blockCreationTicker := time.NewTicker(time.Second * blockInterval)
	if !self.Visor.Config.Config.IsMaster {
		blockCreationTicker.Stop()
	}

	unconfirmedRefreshTicker := time.Tick(self.Visor.Config.Config.UnconfirmedRefreshRate)
	blocksRequestTicker := time.Tick(self.Visor.Config.BlocksRequestRate)
	blocksAnnounceTicker := time.Tick(self.Visor.Config.BlocksAnnounceRate)

	privateConnectionsTicker := time.Tick(self.Config.PrivateRate)
	cullInvalidTicker := time.Tick(self.Config.CullInvalidRate)
	outgoingConnectionsTicker := time.Tick(self.Config.OutgoingRate)
	clearOldPeersTicker := time.Tick(self.Peers.Config.CullRate)
	requestPeersTicker := time.Tick(self.Peers.Config.RequestRate)
	updateBlacklistTicker := time.Tick(self.Peers.Config.UpdateBlacklistRate)
	messageHandlingTicker := time.Tick(self.Pool.Config.MessageHandlingRate)
	clearStaleConnectionsTicker := time.Tick(self.Pool.Config.ClearStaleRate)
	idleCheckTicker := time.Tick(self.Pool.Config.IdleCheckRate)

main:
	for {
		select {
		// Flush expired blacklisted peers
		case <-updateBlacklistTicker:
			if !self.Peers.Config.Disabled {
				self.Peers.Peers.Blacklist.Refresh()
			}
		// Remove connections that failed to complete the handshake
		case <-cullInvalidTicker:
			if !self.Config.DisableNetworking {
				self.cullInvalidConnections()
			}
		// Request peers via PEX
		case <-requestPeersTicker:
			self.Peers.requestPeers(self.Pool)
		// Remove peers we haven't seen in a while
		case <-clearOldPeersTicker:
			if !self.Peers.Config.Disabled {
				self.Peers.Peers.Peerlist.ClearOld(self.Peers.Config.Expiration)
			}
		// Remove connections that haven't said anything in a while
		case <-clearStaleConnectionsTicker:
			if !self.Config.DisableNetworking {
				self.Pool.clearStaleConnections()
			}
		// Sends pings as needed
		case <-idleCheckTicker:
			if !self.Config.DisableNetworking {
				self.Pool.sendPings()
			}
		// Fill up our outgoing connections
		case <-outgoingConnectionsTicker:
			if !self.Config.DisableOutgoingConnections &&
				len(self.OutgoingConnections) < self.Config.OutgoingMax &&
				len(self.pendingConnections) < self.Config.PendingMax {
				self.connectToRandomPeer()
			}
		// Always try to stay connected to our private peers
		// TODO (also, connect to all of them on start)
		case <-privateConnectionsTicker:
			if !self.Config.DisableOutgoingConnections {
				self.makePrivateConnections()
			}
		// Process the connection queue
		case <-messageHandlingTicker:
			if !self.Config.DisableNetworking {
				self.Pool.Pool.HandleMessages()
			}
		// Process callbacks for when a client connects. No disconnect chan
		// is needed because the callback is triggered by HandleDisconnectEvent
		// which is already select{}ed here
		case r := <-self.onConnectEvent:
			if self.Config.DisableNetworking {
				log.Panic("There should be no connect events")
			}
			self.onConnect(r)
		// Handle connection errors
		case r := <-self.connectionErrors:
			if self.Config.DisableNetworking {
				log.Panic("There should be no connection errors")
			}
			self.handleConnectionError(r)
		// Process disconnections
		case r := <-self.Pool.Pool.DisconnectQueue:
			if self.Config.DisableNetworking {
				log.Panic("There should be nothing in the DisconnectQueue")
			}
			self.Pool.Pool.HandleDisconnectEvent(r)
		// Process message sending results
		case r := <-self.Pool.Pool.SendResults:
			if self.Config.DisableNetworking {
				log.Panic("There should be nothing in SendResults")
			}
			self.handleMessageSendResult(r)
		// Message handlers
		case m := <-self.messageEvents:
			if self.Config.DisableNetworking {
				log.Panic("There should be no message events")
			}
			self.processMessageEvent(m)
		// Process any pending RPC requests
		case req := <-self.Gateway.Requests:
			req.Response <- req.Handle()

		//save blockchain periodically
		// case <-blockchainBackupTicker:
		// self.Visor.SaveBlockchain()

		// TODO -- run these in the Visor
		// Create blocks, if master chain
		case <-blockCreationTicker.C:
			if self.Visor.Config.Config.IsMaster {
				err := self.Visor.CreateAndPublishBlock(self.Pool)
				if err != nil {
					logger.Error("Failed to create block: %v", err)
				} else {
					// Not a critical error, but we want it visible in logs
					logger.Critical("Created and published a new block")
				}
			}
		case <-unconfirmedRefreshTicker:
			self.Visor.RefreshUnconfirmed()
		case <-blocksRequestTicker:
			self.Visor.RequestBlocks(self.Pool)
		case <-blocksAnnounceTicker:
			self.Visor.AnnounceBlocks(self.Pool)

		case <-quit:
			break main
		}
	}
}

// Returns the ListenPort for a given address.  If no port is found, 0 is
// returned
func (self *Daemon) GetListenPort(addr string) uint16 {
	m, ok := self.ConnectionMirrors[addr]
	if !ok {
		return 0
	}
	mc := self.mirrorConnections[m]
	if mc == nil {
		log.Panic("mirrorConnections map does not exist, but mirror does")
	}
	a, _, err := SplitAddr(addr)
	if err != nil {
		logger.Error("GetListenPort received invalid addr: %v", err)
		return 0
	} else {
		return mc[a]
	}
}

// Connects to a given peer.  Returns an error if no connection attempt was
// made.  If the connection attempt itself fails, the error is sent to
// the connectionErrors channel.
func (self *Daemon) connectToPeer(p *pex.Peer) error {
	if self.Config.DisableOutgoingConnections {
		return errors.New("Outgoing connections disabled")
	}
	a, _, err := SplitAddr(p.Addr)
	if err != nil {
		logger.Warning("PEX gave us an invalid peer: %v", err)
		return errors.New("Invalid peer")
	}
	if self.Config.LocalhostOnly && !IsLocalhost(a) {
		return errors.New("Not localhost")
	}
	if self.Pool.Pool.Addresses[p.Addr] != nil {
		return errors.New("Already connected")
	}
	if self.pendingConnections[p.Addr] != nil {
		return errors.New("Connection is pending")
	}
	if !self.Config.LocalhostOnly && self.ipCounts[a] != 0 {
		return errors.New("Already connected to a peer with this base IP")
	}
	logger.Debug("Trying to connect to %s", p.Addr)
	self.pendingConnections[p.Addr] = p
	go func() {
		_, err := self.Pool.Pool.Connect(p.Addr)
		if err != nil {
			self.connectionErrors <- ConnectionError{p.Addr, err}
		}
	}()
	return nil
}

// Connects to all private peers
func (self *Daemon) makePrivateConnections() {
	if self.Config.DisableOutgoingConnections {
		return
	}
	for _, p := range self.Peers.Peers.Peerlist {
		if p.Private {
			logger.Info("Private peer attempt: %s", p.Addr)
			if err := self.connectToPeer(p); err != nil {
				logger.Debug("Did not connect to private peer: %v", err)
			}
		}
	}
}

// Attempts to connect to a random peer. If it fails, the peer is removed
func (self *Daemon) connectToRandomPeer() {
	if self.Config.DisableOutgoingConnections {
		return
	}
	// Make a connection to a random (public) peer
	peers := self.Peers.Peers.Peerlist.RandomPublic(0)
	for _, p := range peers {
		if self.connectToPeer(p) == nil {
			break
		}
	}
}

// We remove a peer from the Pex if we failed to connect
// Failure to connect
// Use exponential backoff, not peer list
func (self *Daemon) handleConnectionError(c ConnectionError) {
	logger.Debug("Failed to connect to %s with error: %v", c.Addr,
		c.Error)
	delete(self.pendingConnections, c.Addr)

	if self.Peers.Config.Disabled != true {
		self.Peers.RemovePeer(c.Addr)
	}

	//use exponential backoff

	/*
		duration, exists := BlacklistOffenses[ConnectFailed]
		if exists {
			self.Peers.Peers.AddBlacklistEntry(c.Addr, duration)
		}
	*/
}

// Removes unsolicited connections who haven't sent a version
func (self *Daemon) cullInvalidConnections() {
	// This method only handles the erroneous people from the DHT, but not
	// malicious nodes
	now := util.Now()
	for a, t := range self.ExpectingIntroductions {
		// Forget about anyone that already disconnected
		if self.Pool.Pool.Addresses[a] == nil {
			delete(self.ExpectingIntroductions, a)
			continue
		}
		// Remove anyone that fails to send a version within introductionWait time
		if t.Add(self.Config.IntroductionWait).Before(now) {
			logger.Info("Removing %s for not sending a version", a)
			delete(self.ExpectingIntroductions, a)
			self.Pool.Pool.Disconnect(self.Pool.Pool.Addresses[a],
				DisconnectIntroductionTimeout)
			self.Peers.RemovePeer(a)
		}
	}
}

// Records an AsyncMessage to the messageEvent chan.  Do not access
// messageEvent directly.
func (self *Daemon) recordMessageEvent(m AsyncMessage,
	c *gnet.MessageContext) error {
	self.messageEvents <- MessageEvent{m, c}
	return nil
}

// Processes a queued AsyncMessage.
func (self *Daemon) processMessageEvent(e MessageEvent) {
	// The first message received must be an Introduction
	// We have to check at process time and not record time because
	// Introduction message does not update ExpectingIntroductions until its
	// Process() is called
	_, needsIntro := self.ExpectingIntroductions[e.Context.Conn.Addr()]
	if needsIntro {
		_, isIntro := e.Message.(*IntroductionMessage)
		if !isIntro {
			self.Pool.Pool.Disconnect(e.Context.Conn, DisconnectNoIntroduction)
		}
	}
	e.Message.Process(self)
}

// Called when a ConnectEvent is processed off the onConnectEvent channel
func (self *Daemon) onConnect(e ConnectEvent) {
	a := e.Addr

	if e.Solicited {
		logger.Info("Connected to %s as we requested", a)
	} else {
		logger.Info("Received unsolicited connection to %s", a)
	}

	delete(self.pendingConnections, a)

	c := self.Pool.Pool.Addresses[a]
	if c == nil {
		logger.Warning("While processing an onConnect event, no pool " +
			"connection was found")
		return
	}

	blacklisted := self.Peers.Peers.IsBlacklisted(a)
	if blacklisted {
		logger.Info("%s is blacklisted, disconnecting", a)
		self.Pool.Pool.Disconnect(c, DisconnectIsBlacklisted)
		return
	}

	if self.ipCountMaxed(a) {
		logger.Info("Max connections for %s reached, disconnecting", a)
		self.Pool.Pool.Disconnect(c, DisconnectIPLimitReached)
		return
	}

	self.recordIPCount(a)

	if e.Solicited {
		self.OutgoingConnections[a] = c
	}
	self.ExpectingIntroductions[a] = util.Now()
	logger.Debug("Sending introduction message to %s", a)
	m := NewIntroductionMessage(self.Messages.Mirror, self.Config.Version,
		self.Pool.Pool.Config.Port)
	self.Pool.Pool.SendMessage(c, m)
}

// Triggered when an gnet.Connection terminates. Disconnect events are not
// pushed to a separate channel, because disconnects are already processed
// by a queue in the daemon.Run() select{}.
func (self *Daemon) onGnetDisconnect(c *gnet.Connection,
	reason gnet.DisconnectReason) {
	a := c.Addr()
	logger.Info("%s disconnected because: %v", a, reason)
	duration, exists := BlacklistOffenses[reason]
	if exists {
		self.Peers.Peers.AddBlacklistEntry(a, duration)
	}
	delete(self.OutgoingConnections, a)
	delete(self.ExpectingIntroductions, a)
	self.Visor.RemoveConnection(a)
	self.removeIPCount(a)
	self.removeConnectionMirror(a)
}

// Triggered when an gnet.Connection is connected
func (self *Daemon) onGnetConnect(c *gnet.Connection, solicited bool) {
	self.onConnectEvent <- ConnectEvent{Addr: c.Addr(), Solicited: solicited}
}

// Returns whether the ipCount maximum has been reached
func (self *Daemon) ipCountMaxed(addr string) bool {
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("ipCountMaxed called with invalid addr: %v", err)
		return true
	}
	return self.ipCounts[ip] >= self.Config.IPCountsMax
}

// Adds base IP to ipCount or returns error if max is reached
func (self *Daemon) recordIPCount(addr string) {
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("recordIPCount called with invalid addr: %v", err)
		return
	}
	_, hasCount := self.ipCounts[ip]
	if !hasCount {
		self.ipCounts[ip] = 0
	}
	self.ipCounts[ip] += 1
}

// Removes base IP from ipCount
func (self *Daemon) removeIPCount(addr string) {
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("removeIPCount called with invalid addr: %v", err)
		return
	}
	if self.ipCounts[ip] <= 1 {
		delete(self.ipCounts, ip)
	} else {
		self.ipCounts[ip] -= 1
	}
}

// Adds addr + mirror to the connectionMirror mappings
func (self *Daemon) recordConnectionMirror(addr string, mirror uint32) error {
	ip, port, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("recordConnectionMirror called with invalid addr: %v",
			err)
		return err
	}
	self.ConnectionMirrors[addr] = mirror
	m := self.mirrorConnections[mirror]
	if m == nil {
		m = make(map[string]uint16, 1)
	}
	m[ip] = port
	self.mirrorConnections[mirror] = m
	return nil
}

// Removes an addr from the connectionMirror mappings
func (self *Daemon) removeConnectionMirror(addr string) {
	mirror, ok := self.ConnectionMirrors[addr]
	if !ok {
		return
	}
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("removeConnectionMirror called with invalid addr: %v",
			err)
		return
	}
	m := self.mirrorConnections[mirror]
	if len(m) <= 1 {
		delete(self.mirrorConnections, mirror)
	} else {
		delete(m, ip)
	}
	delete(self.ConnectionMirrors, addr)
}

// Returns whether an addr+mirror's port and whether the port exists
func (self *Daemon) getMirrorPort(addr string, mirror uint32) (uint16, bool) {
	ips := self.mirrorConnections[mirror]
	if ips == nil {
		return 0, false
	}
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("getMirrorPort called with invalid addr: %v", err)
		return 0, false
	}
	port, exists := ips[ip]
	return port, exists
}

// When an async message send finishes, its result is handled by this
func (self *Daemon) handleMessageSendResult(r gnet.SendResult) {
	if r.Error != nil {
		logger.Warning("Failed to send %s to %s: %v",
			reflect.TypeOf(r.Message).Name(), r.Connection.Addr(), r.Error)
		return
	}
	switch r.Message.(type) {
	case SendingTxnsMessage:
		self.Visor.SetTxnsAnnounced(r.Message.(SendingTxnsMessage).GetTxns())
	default:
	}
}

// Returns the address for localhost on the machine
func LocalhostIP() (string, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return "", err
		}
		for _, a := range aa {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("No local IP found")
}

// Returns true if addr is a localhost address
func IsLocalhost(addr string) bool {
	return net.ParseIP(addr).IsLoopback()
}

// Splits an ip:port string to ip, port
func SplitAddr(addr string) (string, uint16, error) {
	pts := strings.Split(addr, ":")
	if len(pts) != 2 {
		return pts[0], 0, fmt.Errorf("Invalid addr %s", addr)
	}
	port64, err := strconv.ParseUint(pts[1], 10, 16)
	if err != nil {
		return pts[0], 0, fmt.Errorf("Invalid port in %s", addr)
	}
	return pts[0], uint16(port64), nil
}
