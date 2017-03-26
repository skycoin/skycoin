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

	logger = util.MustGetLogger("daemon")
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
		OutgoingMax:                16,
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
	outgoingConnections *OutgoingConnections
	// Number of connections waiting to be formed or timeout
	pendingConnections *PendingConnections
	// Keep track of unsolicited clients who should notify us of their version
	expectingIntroductions *ExpectIntroductions
	// Keep track of a connection's mirror value, to avoid double
	// connections (one to their listener, and one to our listener)
	// Maps from addr to mirror value
	connectionMirrors *ConnectionMirrors
	// Maps from mirror value to a map of ip (no port)
	// We use a map of ip as value because multiple peers can have the same
	// mirror (to avoid attacks enabled by our use of mirrors),
	// but only one per base ip
	mirrorConnections *MirrorConnections
	// Client connection/disconnection callbacks
	onConnectEvent chan ConnectEvent
	// Connection failure events
	connectionErrors chan ConnectionError
	// Tracking connections from the same base IP.  Multiple connections
	// from the same base IP are allowed but limited.
	ipCounts *IPCount
	// Message handling queue
	messageEvents chan MessageEvent
	// channel for reading and writing member variable thread safly.
	memChannel chan func()
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

		expectingIntroductions: NewExpectIntroductions(),
		connectionMirrors:      NewConnectionMirrors(),
		mirrorConnections:      NewMirrorConnections(),
		ipCounts:               NewIPCount(),
		// TODO -- if there are performance problems from blocking chans,
		// Its because we are connecting to more things than OutgoingMax
		// if we have private peers
		onConnectEvent: make(chan ConnectEvent,
			config.Daemon.OutgoingMax),
		connectionErrors: make(chan ConnectionError,
			config.Daemon.OutgoingMax),
		outgoingConnections: NewOutgoingConnections(config.Daemon.OutgoingMax),
		pendingConnections:  NewPendingConnections(config.Daemon.PendingMax),
		messageEvents: make(chan MessageEvent,
			config.Pool.EventChannelSize),
		memChannel: make(chan func()),
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

// Shutdown Terminates all subsystems safely.  To stop the Daemon run loop, send a value
// over the quit channel provided to Init.  The Daemon run loop must be stopped
// before calling this function.
func (dm *Daemon) Shutdown() {
	dm.Pool.Shutdown()
	dm.Peers.Shutdown()
	dm.Visor.Shutdown()
	gnet.EraseMessages()
}

// Start main loop for peer/connection management. Send anything to quit to shut it
// down
func (dm *Daemon) Start(quit chan int) {
	if !dm.Config.DisableIncomingConnections {
		dm.Pool.Start()
	}

	// TODO -- run blockchain stuff in its own goroutine
	blockInterval := time.Duration(dm.Visor.Config.Config.BlockCreationInterval)
	// blockchainBackupTicker := time.Tick(self.Visor.Config.BlockchainBackupRate)
	blockCreationTicker := time.NewTicker(time.Second * blockInterval)
	if !dm.Visor.Config.Config.IsMaster {
		blockCreationTicker.Stop()
	}

	unconfirmedRefreshTicker := time.Tick(dm.Visor.Config.Config.UnconfirmedRefreshRate)
	blocksRequestTicker := time.Tick(dm.Visor.Config.BlocksRequestRate)
	blocksAnnounceTicker := time.Tick(dm.Visor.Config.BlocksAnnounceRate)

	privateConnectionsTicker := time.Tick(dm.Config.PrivateRate)
	cullInvalidTicker := time.Tick(dm.Config.CullInvalidRate)
	outgoingConnectionsTicker := time.Tick(dm.Config.OutgoingRate)
	clearOldPeersTicker := time.Tick(dm.Peers.Config.CullRate)
	requestPeersTicker := time.Tick(dm.Peers.Config.RequestRate)
	// updateBlacklistTicker := time.Tick(dm.Peers.Config.UpdateBlacklistRate)
	clearStaleConnectionsTicker := time.Tick(dm.Pool.Config.ClearStaleRate)
	idleCheckTicker := time.Tick(dm.Pool.Config.IdleCheckRate)

	// connecto to trusted peers
	if !dm.Config.DisableOutgoingConnections {
		go dm.connectToTrustPeer()
	}

main:
	for {
		select {
		// Flush expired blacklisted peers
		// case <-updateBlacklistTicker:
		// 	if !dm.Peers.Config.Disabled {
		// 		dm.Peers.Peers.Blacklist.Refresh()
		// 	}
		// Remove connections that failed to complete the handshake
		case <-cullInvalidTicker:
			if !dm.Config.DisableNetworking {
				dm.cullInvalidConnections()
			}
		// Request peers via PEX
		case <-requestPeersTicker:
			dm.Peers.requestPeers(dm.Pool)
		// Remove peers we haven't seen in a while
		case <-clearOldPeersTicker:
			if !dm.Peers.Config.Disabled {
				dm.Peers.Peers.ClearOld(dm.Peers.Config.Expiration)
			}
		// Remove connections that haven't said anything in a while
		case <-clearStaleConnectionsTicker:
			if !dm.Config.DisableNetworking {
				dm.Pool.clearStaleConnections()
			}
		// Sends pings as needed
		case <-idleCheckTicker:
			if !dm.Config.DisableNetworking {
				dm.Pool.sendPings()
			}
		// Fill up our outgoing connections
		case <-outgoingConnectionsTicker:
			trustPeerNum := len(dm.Peers.Peers.GetAllTrustedPeers())
			if !dm.Config.DisableOutgoingConnections &&
				dm.outgoingConnections.Len() < (dm.Config.OutgoingMax+trustPeerNum) &&
				dm.pendingConnections.Len() < dm.Config.PendingMax {
				dm.connectToRandomPeer()
			}
		// Always try to stay connected to our private peers
		// TODO (also, connect to all of them on start)
		case <-privateConnectionsTicker:
			if !dm.Config.DisableOutgoingConnections {
				dm.makePrivateConnections()
			}
		// Process callbacks for when a client connects. No disconnect chan
		// is needed because the callback is triggered by HandleDisconnectEvent
		// which is already select{}ed here
		case r := <-dm.onConnectEvent:
			if dm.Config.DisableNetworking {
				log.Panic("There should be no connect events")
			}
			dm.onConnect(r)
		// Handle connection errors
		case r := <-dm.connectionErrors:
			if dm.Config.DisableNetworking {
				log.Panic("There should be no connection errors")
			}
			dm.handleConnectionError(r)
		// Process message sending results
		case r := <-dm.Pool.Pool.SendResults:
			if dm.Config.DisableNetworking {
				log.Panic("There should be nothing in SendResults")
			}
			dm.handleMessageSendResult(r)
		// Message handlers
		case m := <-dm.messageEvents:
			if dm.Config.DisableNetworking {
				log.Panic("There should be no message events")
			}
			dm.processMessageEvent(m)
		// Process any pending RPC requests
		case req := <-dm.Gateway.Requests:
			req()

		// TODO -- run these in the Visor
		// Create blocks, if master chain
		case <-blockCreationTicker.C:
			if dm.Visor.Config.Config.IsMaster {
				err := dm.Visor.CreateAndPublishBlock(dm.Pool)
				if err != nil {
					logger.Error("Failed to create block: %v", err)
				} else {
					// Not a critical error, but we want it visible in logs
					logger.Critical("Created and published a new block")
				}
			}
		case <-unconfirmedRefreshTicker:
			dm.Visor.RefreshUnconfirmed()
		case <-blocksRequestTicker:
			dm.Visor.RequestBlocks(dm.Pool)
		case <-blocksAnnounceTicker:
			dm.Visor.AnnounceBlocks(dm.Pool)
		case f := <-dm.memChannel:
			f()
		case <-quit:
			break main
		}
	}
}

// GetListenPort returns the ListenPort for a given address.  If no port is found, 0 is
// returned
func (dm *Daemon) GetListenPort(addr string) uint16 {
	m, ok := dm.connectionMirrors.Get(addr)
	if !ok {
		return 0
	}

	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Error("GetListenPort received invalid addr: %v", err)
		return 0
	}

	p, ok := dm.mirrorConnections.Get(m, ip)
	if !ok {
		return 0
	}
	return p
}

// Connects to a given peer.  Returns an error if no connection attempt was
// made.  If the connection attempt itself fails, the error is sent to
// the connectionErrors channel.
func (dm *Daemon) connectToPeer(p *pex.Peer) error {
	if dm.Config.DisableOutgoingConnections {
		return errors.New("Outgoing connections disabled")
	}
	a, _, err := SplitAddr(p.Addr)
	if err != nil {
		logger.Warning("PEX gave us an invalid peer: %v", err)
		return errors.New("Invalid peer")
	}
	if dm.Config.LocalhostOnly && !IsLocalhost(a) {
		return errors.New("Not localhost")
	}

	if dm.Pool.Pool.IsConnExist(p.Addr) {
		return errors.New("Already connected")
	}

	if _, ok := dm.pendingConnections.Get(p.Addr); ok {
		return errors.New("Connection is pending")
	}
	cnt, ok := dm.ipCounts.Get(a)
	if !dm.Config.LocalhostOnly && ok && cnt != 0 {
		return errors.New("Already connected to a peer with this base IP")
	}
	logger.Debug("Trying to connect to %s", p.Addr)
	dm.pendingConnections.Add(p.Addr, p)
	go func() {
		if err := dm.Pool.Pool.Connect(p.Addr); err != nil {
			dm.connectionErrors <- ConnectionError{p.Addr, err}
		}
	}()
	return nil
}

// Connects to all private peers
func (dm *Daemon) makePrivateConnections() {
	if dm.Config.DisableOutgoingConnections {
		return
	}
	addrs := dm.Peers.Peers.GetPrivateAddresses()
	for _, addr := range addrs {
		p, exist := dm.Peers.Peers.GetPeerByAddr(addr)
		if exist {
			logger.Info("Private peer attempt: %s", p.Addr)
			if err := dm.connectToPeer(&p); err != nil {
				logger.Debug("Did not connect to private peer: %v", err)
			}
		}
	}
}

func (dm *Daemon) connectToTrustPeer() {
	if dm.Config.DisableIncomingConnections {
		return
	}

	logger.Info("connect to trusted peers")
	// make connections to all trusted peers
	peers := dm.Peers.Peers.GetPublicTrustPeers()
	for _, p := range peers {
		dm.connectToPeer(p)
		// if dm.connectToPeer(p) == nil {
		// 	break
		// }
	}
}

// Attempts to connect to a random peer. If it fails, the peer is removed
func (dm *Daemon) connectToRandomPeer() {
	if dm.Config.DisableOutgoingConnections {
		return
	}
	// Make a connection to a random (public) peer
	peers := dm.Peers.Peers.RandomPublic(0)
	for _, p := range peers {
		dm.connectToPeer(p)
		// if dm.connectToPeer(p) == nil {
		// 	break
		// }
	}
}

// We remove a peer from the Pex if we failed to connect
// Failure to connect
// Use exponential backoff, not peer list
func (dm *Daemon) handleConnectionError(c ConnectionError) {
	logger.Debug("Failed to connect to %s with error: %v", c.Addr, c.Error)

	dm.pendingConnections.Remove(c.Addr)

	if dm.Peers.Config.Disabled != true {
		dm.Peers.RemovePeer(c.Addr)
	}

	dm.Peers.Peers.IncreaseRetryTimes(c.Addr)
	//use exponential backoff

	/*
		duration, exists := BlacklistOffenses[ConnectFailed]
		if exists {
			self.Peers.Peers.AddBlacklistEntry(c.Addr, duration)
		}
	*/
}

// Removes unsolicited connections who haven't sent a version
func (dm *Daemon) cullInvalidConnections() {
	// This method only handles the erroneous people from the DHT, but not
	// malicious nodes
	now := util.Now()
	addrs := dm.expectingIntroductions.CullInvalidConns(func(addr string, t time.Time) bool {
		if !dm.Pool.Pool.IsConnExist(addr) {
			return true
		}

		if t.Add(dm.Config.IntroductionWait).Before(now) {
			return true
		}
		return false
	})

	for _, a := range addrs {
		if dm.Pool.Pool.IsConnExist(a) {
			logger.Info("Removing %s for not sending a version", a)
			dm.Pool.Pool.Disconnect(a, DisconnectIntroductionTimeout)
			dm.Peers.RemovePeer(a)
		}
	}
}

// Records an AsyncMessage to the messageEvent chan.  Do not access
// messageEvent directly.
func (dm *Daemon) recordMessageEvent(m AsyncMessage,
	c *gnet.MessageContext) error {
	dm.messageEvents <- MessageEvent{m, c}
	return nil
}

// check if the connection needs introduction message
func (dm *Daemon) needsIntro(addr string) bool {
	_, exist := dm.expectingIntroductions.Get(addr)
	return exist
}

// Processes a queued AsyncMessage.
func (dm *Daemon) processMessageEvent(e MessageEvent) {
	// The first message received must be an Introduction
	// We have to check at process time and not record time because
	// Introduction message does not update ExpectingIntroductions until its
	// Process() is called
	// _, needsIntro := self.expectingIntroductions[e.Context.Addr]
	// if needsIntro {
	if dm.needsIntro(e.Context.Addr) {
		_, isIntro := e.Message.(*IntroductionMessage)
		if !isIntro {
			dm.Pool.Pool.Disconnect(e.Context.Addr, DisconnectNoIntroduction)
		}
	}
	e.Message.Process(dm)
}

// Called when a ConnectEvent is processed off the onConnectEvent channel
func (dm *Daemon) onConnect(e ConnectEvent) {
	a := e.Addr

	if e.Solicited {
		logger.Info("Connected to %s as we requested", a)
	} else {
		logger.Info("Received unsolicited connection to %s", a)
	}

	dm.pendingConnections.Remove(a)

	if !dm.Pool.Pool.IsConnExist(a) {
		logger.Warning("While processing an onConnect event, no pool " +
			"connection was found")
		return
	}

	// blacklisted := dm.Peers.Peers.IsBlacklisted(a)
	// if blacklisted {
	// 	logger.Info("%s is blacklisted, disconnecting", a)
	// 	dm.Pool.Pool.Disconnect(a, DisconnectIsBlacklisted)
	// 	return
	// }

	if dm.ipCountMaxed(a) {
		logger.Info("Max connections for %s reached, disconnecting", a)
		dm.Pool.Pool.Disconnect(a, DisconnectIPLimitReached)
		return
	}

	dm.recordIPCount(a)

	if e.Solicited {
		dm.outgoingConnections.Add(a)
	}

	dm.expectingIntroductions.Add(a, util.Now())
	logger.Debug("Sending introduction message to %s", a)
	m := NewIntroductionMessage(dm.Messages.Mirror, dm.Config.Version,
		dm.Pool.Pool.Config.Port)
	dm.Pool.Pool.SendMessage(a, m)
}

// Triggered when an gnet.Connection terminates. Disconnect events are not
// pushed to a separate channel, because disconnects are already processed
// by a queue in the daemon.Run() select{}.
func (dm *Daemon) onGnetDisconnect(addr string, reason gnet.DisconnectReason) {
	// a := c.Addr()
	logger.Info("%s disconnected because: %v", addr, reason)
	// duration, exists := BlacklistOffenses[reason]
	// if exists {
	// 	dm.Peers.Peers.AddBlacklistEntry(addr, duration)
	// }

	dm.outgoingConnections.Remove(addr)
	dm.expectingIntroductions.Remove(addr)
	dm.Visor.RemoveConnection(addr)
	dm.removeIPCount(addr)
	dm.removeConnectionMirror(addr)
}

// Triggered when an gnet.Connection is connected
func (dm *Daemon) onGnetConnect(addr string, solicited bool) {
	dm.onConnectEvent <- ConnectEvent{Addr: addr, Solicited: solicited}
}

// Returns whether the ipCount maximum has been reached
func (dm *Daemon) ipCountMaxed(addr string) bool {
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("ipCountMaxed called with invalid addr: %v", err)
		return true
	}

	if cnt, ok := dm.ipCounts.Get(ip); ok {
		return cnt >= dm.Config.IPCountsMax
	}
	return false
}

// Adds base IP to ipCount or returns error if max is reached
func (dm *Daemon) recordIPCount(addr string) {
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("recordIPCount called with invalid addr: %v", err)
		return
	}
	dm.ipCounts.Increase(ip)
}

// Removes base IP from ipCount
func (dm *Daemon) removeIPCount(addr string) {
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("removeIPCount called with invalid addr: %v", err)
		return
	}
	dm.ipCounts.Decrease(ip)
}

// Adds addr + mirror to the connectionMirror mappings
func (dm *Daemon) recordConnectionMirror(addr string, mirror uint32) error {
	ip, port, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("recordConnectionMirror called with invalid addr: %v",
			err)
		return err
	}
	dm.connectionMirrors.Add(addr, mirror)
	dm.mirrorConnections.Add(mirror, ip, port)
	return nil
}

// Removes an addr from the connectionMirror mappings
func (dm *Daemon) removeConnectionMirror(addr string) {
	mirror, ok := dm.connectionMirrors.Get(addr)
	if !ok {
		return
	}
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("removeConnectionMirror called with invalid addr: %v",
			err)
		return
	}

	// remove ip from specific mirror
	dm.mirrorConnections.Remove(mirror, ip)

	dm.connectionMirrors.Remove(addr)
}

// Returns whether an addr+mirror's port and whether the port exists
func (dm *Daemon) getMirrorPort(addr string, mirror uint32) (uint16, bool) {
	ip, _, err := SplitAddr(addr)
	if err != nil {
		logger.Warning("getMirrorPort called with invalid addr: %v", err)
		return 0, false
	}
	return dm.mirrorConnections.Get(mirror, ip)
}

// When an async message send finishes, its result is handled by this
func (dm *Daemon) handleMessageSendResult(r gnet.SendResult) {
	if r.Error != nil {
		logger.Warning("Failed to send %s to %s: %v",
			reflect.TypeOf(r.Message).Name(), r.Addr, r.Error)
		return
	}
	switch r.Message.(type) {
	case SendingTxnsMessage:
		dm.Visor.SetTxnsAnnounced(r.Message.(SendingTxnsMessage).GetTxns())
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
