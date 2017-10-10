package daemon

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
)

/*
Todo
- verify that minimum/maximum connections are working
- keep max connections
- maintain minimum number of outgoing connections per server?


*/
var (
	// ErrDisconnectReasons invalid version
	ErrDisconnectInvalidVersion gnet.DisconnectReason = errors.New("Invalid version")
	// ErrDisconnectIntroductionTimeout timeout
	ErrDisconnectIntroductionTimeout gnet.DisconnectReason = errors.New("Version timeout")
	// ErrDisconnectVersionSendFailed version send failed
	ErrDisconnectVersionSendFailed gnet.DisconnectReason = errors.New("Version send failed")
	// ErrDisconnectIsBlacklisted is blacklisted
	ErrDisconnectIsBlacklisted gnet.DisconnectReason = errors.New("Blacklisted")
	// ErrDisconnectSelf self connnect
	ErrDisconnectSelf gnet.DisconnectReason = errors.New("Self connect")
	// ErrDisconnectConnectedTwice connect twice
	ErrDisconnectConnectedTwice gnet.DisconnectReason = errors.New("Already connected")
	// ErrDisconnectIdle idle
	ErrDisconnectIdle gnet.DisconnectReason = errors.New("Idle")
	// ErrDisconnectNoIntroduction no introduction
	ErrDisconnectNoIntroduction gnet.DisconnectReason = errors.New("First message was not an Introduction")
	// ErrDisconnectIPLimitReached ip limit reached
	ErrDisconnectIPLimitReached gnet.DisconnectReason = errors.New("Maximum number of connections for this IP was reached")
	// ErrDisconnectOtherError this is returned when a seemingly impossible error is encountered
	// e.g. net.Conn.Addr() returns an invalid ip:port
	ErrDisconnectOtherError gnet.DisconnectReason = errors.New("Incomprehensible error")

	logger = logging.MustGetLogger("daemon")
)

const (
	// MaxDropletPrecision represents the precision of droplets
	MaxDropletPrecision = 1
	MaxDropletDivisor   = 1e6
)

// Config subsystem configurations
type Config struct {
	Daemon   DaemonConfig
	Messages MessagesConfig
	Pool     PoolConfig
	Peers    PeersConfig
	Gateway  GatewayConfig
	Visor    VisorConfig
}

// NewConfig returns a Config with defaults set
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

// preprocess preprocess for config
func (cfg *Config) preprocess() Config {
	config := *cfg
	if config.Daemon.LocalhostOnly {
		if config.Daemon.Address == "" {
			local, err := LocalhostIP()
			if err != nil {
				logger.Panicf("Failed to obtain localhost IP: %v", err)
			}
			config.Daemon.Address = local
		} else {
			if !IsLocalhost(config.Daemon.Address) {
				logger.Panicf("Invalid address for localhost-only: %s",
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

// DaemonConfig configuration for the Daemon
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
	// Log ping and pong messages
	LogPings bool
}

// NewDaemonConfig creates daemon config
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
		LogPings:                   true,
	}
}

// Daemon stateful properties of the daemon
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
	// Client connection callbacks
	onConnectEvent chan ConnectEvent
	// Client disconnection callbacks
	onDisconnectEvent chan DisconnectEvent
	// Connection failure events
	connectionErrors chan ConnectionError
	// Tracking connections from the same base IP.  Multiple connections
	// from the same base IP are allowed but limited.
	ipCounts *IPCount
	// Message handling queue
	messageEvents chan MessageEvent
	// quit channel
	quitC chan chan struct{}
}

// NewDaemon returns a Daemon with primitives allocated
func NewDaemon(config Config) (*Daemon, error) {
	config = config.preprocess()
	vs, err := NewVisor(config.Visor)
	if err != nil {
		return nil, err
	}

	peers, err := NewPeers(config.Peers)
	if err != nil {
		return nil, err
	}

	d := &Daemon{
		Config:   config.Daemon,
		Messages: NewMessages(config.Messages),
		Peers:    peers,
		Visor:    vs,

		DefaultConnections: DefaultConnections, //passed in from top level

		expectingIntroductions: NewExpectIntroductions(),
		connectionMirrors:      NewConnectionMirrors(),
		mirrorConnections:      NewMirrorConnections(),
		ipCounts:               NewIPCount(),
		// TODO -- if there are performance problems from blocking chans,
		// Its because we are connecting to more things than OutgoingMax
		// if we have private peers
		onConnectEvent:      make(chan ConnectEvent, config.Daemon.OutgoingMax),
		onDisconnectEvent:   make(chan DisconnectEvent, config.Daemon.OutgoingMax),
		connectionErrors:    make(chan ConnectionError, config.Daemon.OutgoingMax),
		outgoingConnections: NewOutgoingConnections(config.Daemon.OutgoingMax),
		pendingConnections:  NewPendingConnections(config.Daemon.PendingMax),
		messageEvents:       make(chan MessageEvent, config.Pool.EventChannelSize),
		quitC:               make(chan chan struct{}),
	}

	d.Gateway = NewGateway(config.Gateway, d)
	d.Messages.Config.Register()
	d.Pool = NewPool(config.Pool, d)

	return d, nil
}

// ConnectEvent generated when a client connects
type ConnectEvent struct {
	Addr      string
	Solicited bool
}

// DisconnectEvent generated when a connection terminated
type DisconnectEvent struct {
	Addr   string
	Reason gnet.DisconnectReason
}

// ConnectionError represent a failure to connect/dial a connection, with context
type ConnectionError struct {
	Addr  string
	Error error
}

// MessageEvent encapsulates a deserialized message from the network
type MessageEvent struct {
	Message AsyncMessage
	Context *gnet.MessageContext
}

// Shutdown Terminates all subsystems safely.  To stop the Daemon run loop, send a value
// over the quit channel provided to Init.  The Daemon run loop must be stopped
// before calling this function.
func (dm *Daemon) Shutdown() {
	// close the daemon loop first
	close(dm.quitC)

	if !dm.Config.DisableNetworking {
		dm.Pool.Shutdown()
	}

	dm.Peers.Shutdown()
	dm.Visor.Shutdown()
}

// Run main loop for peer/connection management. Send anything to quit to shut it
// down
func (dm *Daemon) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("recover:%v\n stack:%v", r, string(debug.Stack()))
		}

		logger.Info("Daemon closed")
	}()

	errC := make(chan error)

	// start visor
	go func() {
		errC <- dm.Visor.Run()
	}()

	if !dm.Config.DisableIncomingConnections {
		go func() {
			errC <- dm.Pool.Run()
		}()
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
	clearStaleConnectionsTicker := time.Tick(dm.Pool.Config.ClearStaleRate)
	idleCheckTicker := time.Tick(dm.Pool.Config.IdleCheckRate)

	// connecto to trusted peers
	if !dm.Config.DisableOutgoingConnections {
		go dm.connectToTrustPeer()
	}

	for {
		select {
		case err = <-errC:
			return
		case <-dm.quitC:
			return
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
				logger.Error("There should be no connect events")
				return
			}
			dm.onConnect(r)
		case de := <-dm.onDisconnectEvent:
			if dm.Config.DisableNetworking {
				logger.Error("There should be no disconnect events")
				return
			}
			dm.onDisconnect(de)
		// Handle connection errors
		case r := <-dm.connectionErrors:
			if dm.Config.DisableNetworking {
				logger.Error("There should be no connection errors")
				return
			}
			dm.handleConnectionError(r)
		// Process message sending results
		case r := <-dm.Pool.Pool.SendResults:
			if dm.Config.DisableNetworking {
				logger.Error("There should be nothing in SendResults")
				return
			}
			dm.handleMessageSendResult(r)
		// Message handlers
		case m := <-dm.messageEvents:
			if dm.Config.DisableNetworking {
				logger.Error("There should be no message events")
				return
			}
			dm.processMessageEvent(m)
		// Process any pending RPC requests
		case req := <-dm.Gateway.requests:
			req()
		// TODO -- run these in the Visor
		// Create blocks, if master chain
		case <-blockCreationTicker.C:
			if dm.Visor.Config.Config.IsMaster {
				err := dm.Visor.CreateAndPublishBlock(dm.Pool)
				if err != nil {
					logger.Error("Failed to create block: %v", err)
					continue
				}

				// Not a critical error, but we want it visible in logs
				logger.Critical("Created and published a new block")
			}
		case <-unconfirmedRefreshTicker:
			// get the transactions that turn to valid
			validTxns := dm.Visor.RefreshUnconfirmed()
			// announce this transactions
			dm.Visor.AnnounceTxns(dm.Pool, validTxns)
		case <-blocksRequestTicker:
			dm.Visor.RequestBlocks(dm.Pool)
		case <-blocksAnnounceTicker:
			dm.Visor.AnnounceBlocks(dm.Pool)
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

	conned, err := dm.Pool.Pool.IsConnExist(p.Addr)
	if err != nil {
		return err
	}

	if conned {
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
		// check if the peer has public port
		if p.HasIncomePort {
			// try to connect the peer if it's ip:mirror does not exist
			if _, exist := dm.getMirrorPort(p.Addr, dm.Messages.Mirror); !exist {
				dm.connectToPeer(p)
				continue
			}
		} else {
			// try to connect to the peer if we don't know whether the peer have public port
			dm.connectToPeer(p)
		}
	}

	if len(peers) == 0 {
		// reset the retry times of all peers
		dm.Peers.Peers.ResetAllRetryTimes()
	}
}

// We remove a peer from the Pex if we failed to connect
// Failure to connect
// Use exponential backoff, not peer list
func (dm *Daemon) handleConnectionError(c ConnectionError) {
	logger.Debug("Failed to connect to %s with error: %v", c.Addr, c.Error)

	dm.pendingConnections.Remove(c.Addr)

	dm.Peers.Peers.IncreaseRetryTimes(c.Addr)
}

// Removes unsolicited connections who haven't sent a version
func (dm *Daemon) cullInvalidConnections() {
	// This method only handles the erroneous people from the DHT, but not
	// malicious nodes
	now := utc.Now()
	addrs, err := dm.expectingIntroductions.CullInvalidConns(func(addr string, t time.Time) (bool, error) {
		conned, err := dm.Pool.Pool.IsConnExist(addr)
		if err != nil {
			return false, err
		}

		if !conned {
			return true, nil
		}

		if t.Add(dm.Config.IntroductionWait).Before(now) {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		logger.Error("expectingIntroduction cull invalid connections failed: %v", err)
		return
	}

	for _, a := range addrs {
		exist, err := dm.Pool.Pool.IsConnExist(a)
		if err != nil {
			logger.Error("%v", err)
			return
		}

		if exist {
			logger.Info("Removing %s for not sending a version", a)
			if err := dm.Pool.Pool.Disconnect(a, ErrDisconnectIntroductionTimeout); err != nil {
				logger.Error("%v", err)
				return
			}
			dm.Peers.RemovePeer(a)
		}
	}
}

// Records an AsyncMessage to the messageEvent chan.  Do not access
// messageEvent directly.
func (dm *Daemon) recordMessageEvent(m AsyncMessage, c *gnet.MessageContext) error {
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
			dm.Pool.Pool.Disconnect(e.Context.Addr, ErrDisconnectNoIntroduction)
		}
	}
	e.Message.Process(dm)
}

// Called when a ConnectEvent is processed off the onConnectEvent channel
func (dm *Daemon) onConnect(e ConnectEvent) {
	a := e.Addr

	if e.Solicited {
		logger.Info("Connected to peer: %s (outgoing)", a)
	} else {
		logger.Info("Connected to peer: %s (incoming)", a)
	}

	dm.pendingConnections.Remove(a)

	exist, err := dm.Pool.Pool.IsConnExist(a)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	if !exist {
		logger.Warning("While processing an onConnect event, no pool " +
			"connection was found")
		return
	}

	if dm.ipCountMaxed(a) {
		logger.Info("Max connections for %s reached, disconnecting", a)
		dm.Pool.Pool.Disconnect(a, ErrDisconnectIPLimitReached)
		return
	}

	dm.recordIPCount(a)

	if e.Solicited {
		dm.outgoingConnections.Add(a)
	}

	dm.expectingIntroductions.Add(a, utc.Now())
	logger.Debug("Sending introduction message to %s, mirror:%d", a, dm.Messages.Mirror)
	m := NewIntroductionMessage(dm.Messages.Mirror, dm.Config.Version,
		dm.Pool.Pool.Config.Port)
	dm.Pool.Pool.SendMessage(a, m)
}

func (dm *Daemon) onDisconnect(e DisconnectEvent) {
	logger.Info("%s disconnected because: %v", e.Addr, e.Reason)

	dm.outgoingConnections.Remove(e.Addr)
	dm.expectingIntroductions.Remove(e.Addr)
	dm.Visor.RemoveConnection(e.Addr)
	dm.removeIPCount(e.Addr)
	dm.removeConnectionMirror(e.Addr)
}

// Triggered when an gnet.Connection terminates
func (dm *Daemon) onGnetDisconnect(addr string, reason gnet.DisconnectReason) {
	e := DisconnectEvent{
		Addr:   addr,
		Reason: reason,
	}
	select {
	case dm.onDisconnectEvent <- e:
	default:
		logger.Info("onDisconnectEvent channel is full")
	}
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

// LocalhostIP returns the address for localhost on the machine
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

// IsLocalhost returns true if addr is a localhost address
func IsLocalhost(addr string) bool {
	return net.ParseIP(addr).IsLoopback()
}

// SplitAddr splits an ip:port string to ip, port
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

// DropletPrecisionCheck checks if the amount is valid
func DropletPrecisionCheck(amount uint64) error {
	if amount%MaxDropletDivisor != 0 {
		return fmt.Errorf("invalid amount, too many decimal place")
	}

	return nil
}
