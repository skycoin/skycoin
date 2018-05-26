package daemon

import (
	"bytes"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"

	"github.com/skycoin/skycoin/src/util/elapse"
	"github.com/skycoin/skycoin/src/util/iputil"
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
	ErrDisconnectOtherError                  gnet.DisconnectReason = errors.New("Incomprehensible error")
	ErrDisconnectMaxDefaultConnectionReached                       = errors.New("Maximum number of default connections was reached")
	// ErrDisconnectMaxOutgoingConnectionsReached is returned when connection pool size is greater than the maximum allowed
	ErrDisconnectMaxOutgoingConnectionsReached gnet.DisconnectReason = errors.New("Maximum outgoing connections was reached")

	logger = logging.MustGetLogger("daemon")
)

const (
	daemonRunDurationThreshold = time.Millisecond * 200
)

// Config subsystem configurations
type Config struct {
	Daemon   DaemonConfig
	Messages MessagesConfig
	Pool     PoolConfig
	Pex      pex.Config
	Gateway  GatewayConfig
	Visor    visor.Config
}

// NewConfig returns a Config with defaults set
func NewConfig() Config {
	return Config{
		Daemon:   NewDaemonConfig(),
		Pool:     NewPoolConfig(),
		Pex:      pex.NewConfig(),
		Gateway:  NewGatewayConfig(),
		Messages: NewMessagesConfig(),
		Visor:    visor.NewVisorConfig(),
	}
}

// preprocess preprocess for config
func (cfg *Config) preprocess() Config {
	config := *cfg
	if config.Daemon.LocalhostOnly {
		if config.Daemon.Address == "" {
			local, err := iputil.LocalhostIP()
			if err != nil {
				logger.Panicf("Failed to obtain localhost IP: %v", err)
			}
			config.Daemon.Address = local
		} else {
			if !iputil.IsLocalhost(config.Daemon.Address) {
				logger.Panicf("Invalid address for localhost-only: %s", config.Daemon.Address)
			}
		}
		config.Pex.AllowLocalhost = true
	}
	config.Pool.port = config.Daemon.Port
	config.Pool.address = config.Daemon.Address

	if config.Daemon.DisableNetworking {
		logger.Info("Networking is disabled")
		config.Pex.Disabled = true
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
type DaemonConfig struct { // nolint: golint
	// Protocol version. TODO -- manage version better
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
	// How often to update the database with transaction announcement timestamps
	FlushAnnouncedTxnsRate time.Duration
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
	// How often to request blocks from peers
	BlocksRequestRate time.Duration
	// How often to announce our blocks to peers
	BlocksAnnounceRate time.Duration
	// How many blocks to respond with to a GetBlocksMessage
	BlocksResponseCount uint64
	// Max announce txns hash number
	MaxTxnAnnounceNum int
	// How often new blocks are created by the signing node, in seconds
	BlockCreationInterval uint64
	// How often to check the unconfirmed pool for transactions that become valid
	UnconfirmedRefreshRate time.Duration
	// How often to remove transactions that become permanently invalid from the unconfirmed pool
	UnconfirmedRemoveInvalidRate time.Duration
}

// NewDaemonConfig creates daemon config
func NewDaemonConfig() DaemonConfig {
	return DaemonConfig{
		Version:                      2,
		Address:                      "",
		Port:                         6677,
		OutgoingRate:                 time.Second * 5,
		PrivateRate:                  time.Second * 5,
		OutgoingMax:                  8,
		PendingMax:                   8,
		IntroductionWait:             time.Second * 30,
		CullInvalidRate:              time.Second * 3,
		FlushAnnouncedTxnsRate:       time.Second * 3,
		IPCountsMax:                  3,
		DisableNetworking:            false,
		DisableOutgoingConnections:   false,
		DisableIncomingConnections:   false,
		LocalhostOnly:                false,
		LogPings:                     true,
		BlocksRequestRate:            time.Second * 60,
		BlocksAnnounceRate:           time.Second * 60,
		BlocksResponseCount:          20,
		MaxTxnAnnounceNum:            16,
		BlockCreationInterval:        10,
		UnconfirmedRefreshRate:       time.Minute,
		UnconfirmedRemoveInvalidRate: time.Minute,
	}
}

// Daemon stateful properties of the daemon
type Daemon struct {
	// Daemon configuration
	Config DaemonConfig

	// Components
	Messages *Messages
	Pool     *Pool
	Pex      *pex.Pex
	Gateway  *Gateway
	Visor    *visor.Visor

	DefaultConnections []string

	// Cache of announced transactions that are flushed to the database periodically
	announcedTxns *announcedTxnsCache
	// Cache of reported peer blockchain heights
	Heights *peerBlockchainHeights

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
	quit chan struct{}
	// done channel
	done chan struct{}
	// log buffer
	LogBuff bytes.Buffer
}

// NewDaemon returns a Daemon with primitives allocated
func NewDaemon(config Config, db *dbutil.DB, defaultConns []string) (*Daemon, error) {
	config = config.preprocess()
	vs, err := visor.NewVisor(config.Visor, db)
	if err != nil {
		return nil, err
	}

	pex, err := pex.New(config.Pex, defaultConns)
	if err != nil {
		return nil, err
	}

	d := &Daemon{
		Config:   config.Daemon,
		Messages: NewMessages(config.Messages),
		Pex:      pex,
		Visor:    vs,

		DefaultConnections: defaultConns,

		announcedTxns: newAnnouncedTxnsCache(),
		Heights:       newPeerBlockchainHeights(),

		expectingIntroductions: NewExpectIntroductions(),
		connectionMirrors:      NewConnectionMirrors(),
		mirrorConnections:      NewMirrorConnections(),
		ipCounts:               NewIPCount(),
		// TODO -- if there are performance problems from blocking chans,
		// Its because we are connecting to more things than OutgoingMax
		// if we have private peers
		onConnectEvent:      make(chan ConnectEvent, config.Pool.MaxConnections*2),
		onDisconnectEvent:   make(chan DisconnectEvent, config.Pool.MaxConnections*2),
		connectionErrors:    make(chan ConnectionError, config.Pool.MaxConnections*2),
		outgoingConnections: NewOutgoingConnections(config.Daemon.OutgoingMax),
		pendingConnections:  NewPendingConnections(config.Daemon.PendingMax),
		messageEvents:       make(chan MessageEvent, config.Pool.EventChannelSize),
		quit:                make(chan struct{}),
		done:                make(chan struct{}),
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
	defer logger.Info("Daemon shutdown complete")

	// close daemon run loop first to avoid creating new connection after
	// the connection pool is shutdown.
	logger.Info("Stopping the daemon run loop")
	close(dm.quit)

	logger.Info("Shutting down Pool")
	dm.Pool.Shutdown()

	logger.Info("Shutting down Gateway")
	dm.Gateway.Shutdown()

	logger.Info("Shutting down Pex")
	dm.Pex.Shutdown()

	<-dm.done
}

// Run main loop for peer/connection management.
// Send anything to the quit channel to shut it down.
func (dm *Daemon) Run() error {
	defer logger.Info("Daemon closed")
	defer close(dm.done)

	if err := dm.Visor.Init(); err != nil {
		logger.WithError(err).Error("visor.Visor.Init failed")
		return err
	}

	errC := make(chan error, 5)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := dm.Pex.Run(); err != nil {
			logger.WithError(err).Error("daemon.Pex.Run failed")
			errC <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if dm.Config.DisableIncomingConnections {
			if err := dm.Pool.RunOffline(); err != nil {
				logger.WithError(err).Error("daemon.Pool.RunOffline failed")
				errC <- err
			}
		} else {
			if err := dm.Pool.Run(); err != nil {
				logger.WithError(err).Error("daemon.Pool.Run failed")
				errC <- err
			}
		}
	}()

	blockInterval := time.Duration(dm.Config.BlockCreationInterval)
	blockCreationTicker := time.NewTicker(time.Second * blockInterval)
	if !dm.Visor.Config.IsMaster {
		blockCreationTicker.Stop()
	}

	unconfirmedRefreshTicker := time.Tick(dm.Config.UnconfirmedRefreshRate)
	unconfirmedRemoveInvalidTicker := time.Tick(dm.Config.UnconfirmedRemoveInvalidRate)
	blocksRequestTicker := time.Tick(dm.Config.BlocksRequestRate)
	blocksAnnounceTicker := time.Tick(dm.Config.BlocksAnnounceRate)

	privateConnectionsTicker := time.Tick(dm.Config.PrivateRate)
	cullInvalidTicker := time.Tick(dm.Config.CullInvalidRate)
	outgoingConnectionsTicker := time.Tick(dm.Config.OutgoingRate)
	requestPeersTicker := time.Tick(dm.Pex.Config.RequestRate)
	clearStaleConnectionsTicker := time.Tick(dm.Pool.Config.ClearStaleRate)
	idleCheckTicker := time.Tick(dm.Pool.Config.IdleCheckRate)

	flushAnnouncedTxnsTicker := time.Tick(dm.Config.FlushAnnouncedTxnsRate)

	// Connect to trusted peers
	if !dm.Config.DisableOutgoingConnections {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dm.connectToTrustPeer()
		}()
	}

	var err error
	elapser := elapse.NewElapser(daemonRunDurationThreshold, logger)

	// Process SendResults in a separate goroutine, otherwise SendResults
	// will fill up much faster than can be processed by the daemon run loop
	// dm.handleMessageSendResult must take care not to perform any operation
	// that would violate thread safety, since it is not serialized by the daemon run loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		elapser := elapse.NewElapser(daemonRunDurationThreshold, logger)
	loop:
		for {
			elapser.CheckForDone()
			select {
			case <-dm.quit:
				break loop

			case r := <-dm.Pool.Pool.SendResults:
				// Process message sending results
				elapser.Register("dm.Pool.Pool.SendResults")
				if dm.Config.DisableNetworking {
					logger.Error("There should be nothing in SendResults")
					return
				}
				dm.handleMessageSendResult(r)
			}
		}
	}()

loop:
	for {
		elapser.CheckForDone()
		select {
		case <-dm.quit:
			break loop

		case <-cullInvalidTicker:
			// Remove connections that failed to complete the handshake
			elapser.Register("cullInvalidTicker")
			if !dm.Config.DisableNetworking {
				dm.cullInvalidConnections()
			}

		case <-requestPeersTicker:
			// Request peers via PEX
			elapser.Register("requestPeersTicker")
			if dm.Pex.Config.Disabled {
				continue
			}

			if dm.Pex.IsFull() {
				continue
			}

			m := NewGetPeersMessage()
			if err := dm.Pool.Pool.BroadcastMessage(m); err != nil {
				logger.Error(err)
			}

		case <-clearStaleConnectionsTicker:
			// Remove connections that haven't said anything in a while
			elapser.Register("clearStaleConnectionsTicker")
			if !dm.Config.DisableNetworking {
				dm.Pool.clearStaleConnections()
			}

		case <-idleCheckTicker:
			// Sends pings as needed
			elapser.Register("idleCheckTicker")
			if !dm.Config.DisableNetworking {
				dm.Pool.sendPings()
			}

		case <-outgoingConnectionsTicker:
			// Fill up our outgoing connections
			elapser.Register("outgoingConnectionsTicker")
			trustPeerNum := len(dm.Pex.Trusted())
			if !dm.Config.DisableOutgoingConnections &&
				dm.outgoingConnections.Len() < (dm.Config.OutgoingMax+trustPeerNum) &&
				dm.pendingConnections.Len() < dm.Config.PendingMax {
				dm.connectToRandomPeer()
			}

		case <-privateConnectionsTicker:
			// Always try to stay connected to our private peers
			// TODO (also, connect to all of them on start)
			elapser.Register("privateConnectionsTicker")
			if !dm.Config.DisableOutgoingConnections {
				dm.makePrivateConnections()
			}

		case r := <-dm.onConnectEvent:
			// Process callbacks for when a client connects. No disconnect chan
			// is needed because the callback is triggered by HandleDisconnectEvent
			// which is already select{}ed here
			elapser.Register("dm.onConnectEvent")
			if dm.Config.DisableNetworking {
				logger.Error("There should be no connect events")
				return nil
			}
			dm.onConnect(r)

		case de := <-dm.onDisconnectEvent:
			elapser.Register("dm.onDisconnectEvent")
			if dm.Config.DisableNetworking {
				logger.Error("There should be no disconnect events")
				return nil
			}
			dm.onDisconnect(de)

		case r := <-dm.connectionErrors:
			// Handle connection errors
			elapser.Register("dm.connectionErrors")
			if dm.Config.DisableNetworking {
				logger.Error("There should be no connection errors")
				return nil
			}
			dm.handleConnectionError(r)

		case <-flushAnnouncedTxnsTicker:
			elapser.Register("flushAnnouncedTxnsTicker")
			txns := dm.announcedTxns.flush()

			if err := dm.Visor.SetTxnsAnnounced(txns); err != nil {
				logger.WithError(err).Error("Failed to set unconfirmed txn announce time")
				return err
			}

		case m := <-dm.messageEvents:
			// Message handlers
			elapser.Register("dm.messageEvents")
			if dm.Config.DisableNetworking {
				logger.Error("There should be no message events")
				return nil
			}
			dm.processMessageEvent(m)

		case req := <-dm.Gateway.requests:
			// Process any pending RPC requests
			elapser.Register("dm.Gateway.requests")
			req.Func()

		case <-blockCreationTicker.C:
			// Create blocks, if master chain
			elapser.Register("blockCreationTicker.C")
			if dm.Visor.Config.IsMaster {
				sb, err := dm.CreateAndPublishBlock()
				if err != nil {
					logger.Errorf("Failed to create and publish block: %v", err)
					continue
				}

				// Not a critical error, but we want it visible in logs
				head := sb.Block.Head
				logger.Critical().Infof("Created and published a new block, version=%d seq=%d time=%d", head.Version, head.BkSeq, head.Time)
			}

		case <-unconfirmedRefreshTicker:
			elapser.Register("unconfirmedRefreshTicker")
			// Get the transactions that turn to valid
			validTxns, err := dm.Visor.RefreshUnconfirmed()
			if err != nil {
				logger.Errorf("dm.Visor.RefreshUnconfirmed failed: %v", err)
				continue
			}
			// Announce these transactions
			dm.AnnounceTxns(validTxns)

		case <-unconfirmedRemoveInvalidTicker:
			elapser.Register("unconfirmedRemoveInvalidTicker")
			// Remove transactions that become invalid (violating hard constraints)
			removedTxns, err := dm.Visor.RemoveInvalidUnconfirmed()
			if err != nil {
				logger.Errorf("dm.Visor.RemoveInvalidUnconfirmed failed: %v", err)
				continue
			}
			if len(removedTxns) > 0 {
				logger.Infof("Remove %d txns from pool that began violating hard constraints", len(removedTxns))
			}

		case <-blocksRequestTicker:
			elapser.Register("blocksRequestTicker")
			dm.RequestBlocks()

		case <-blocksAnnounceTicker:
			elapser.Register("blocksAnnounceTicker")
			dm.AnnounceBlocks()

		case err = <-errC:
			break loop
		}
	}

	wg.Wait()

	return err
}

// GetListenPort returns the ListenPort for a given address.
// If no port is found, 0 is returned.
func (dm *Daemon) GetListenPort(addr string) uint16 {
	m, ok := dm.connectionMirrors.Get(addr)
	if !ok {
		return 0
	}

	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Errorf("GetListenPort received invalid addr: %v", err)
		return 0
	}

	p, ok := dm.mirrorConnections.Get(m, ip)
	if !ok {
		return 0
	}
	return p
}

// Connects to a given peer. Returns an error if no connection attempt was
// made. If the connection attempt itself fails, the error is sent to
// the connectionErrors channel.
func (dm *Daemon) connectToPeer(p pex.Peer) error {
	if dm.Config.DisableOutgoingConnections {
		return errors.New("Outgoing connections disabled")
	}

	a, _, err := iputil.SplitAddr(p.Addr)
	if err != nil {
		logger.Warningf("PEX gave us an invalid peer: %v", err)
		return errors.New("Invalid peer")
	}
	if dm.Config.LocalhostOnly && !iputil.IsLocalhost(a) {
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

	logger.Debugf("Trying to connect to %s", p.Addr)
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

	peers := dm.Pex.Private()
	for _, p := range peers {
		logger.Infof("Private peer attempt: %s", p.Addr)
		if err := dm.connectToPeer(p); err != nil {
			logger.Debugf("Did not connect to private peer: %v", err)
		}
	}
}

func (dm *Daemon) connectToTrustPeer() {
	if dm.Config.DisableIncomingConnections {
		return
	}

	logger.Info("Connect to trusted peers")
	// Make connections to all trusted peers
	peers := dm.Pex.TrustedPublic()
	for _, p := range peers {
		dm.connectToPeer(p)
	}
}

// Attempts to connect to a random peer. If it fails, the peer is removed.
func (dm *Daemon) connectToRandomPeer() {
	if dm.Config.DisableOutgoingConnections {
		return
	}

	// Make a connection to a random (public) peer
	peers := dm.Pex.RandomPublic(dm.Config.OutgoingMax)
	for _, p := range peers {
		// Check if the peer has public port
		if p.HasIncomingPort {
			// Try to connect the peer if it's ip:mirror does not exist
			if _, exist := dm.getMirrorPort(p.Addr, dm.Messages.Mirror); !exist {
				dm.connectToPeer(p)
				continue
			}
		} else {
			// Try to connect to the peer if we don't know whether the peer have public port
			dm.connectToPeer(p)
		}
	}

	if len(peers) == 0 {
		// Reset the retry times of all peers,
		dm.Pex.ResetAllRetryTimes()
	}
}

// We remove a peer from the Pex if we failed to connect
// TODO - On failure to connect, use exponential backoff, not peer list
func (dm *Daemon) handleConnectionError(c ConnectionError) {
	logger.Debugf("Failed to connect to %s with error: %v", c.Addr, c.Error)
	dm.pendingConnections.Remove(c.Addr)

	dm.Pex.IncreaseRetryTimes(c.Addr)
}

// Removes unsolicited connections who haven't sent a version
func (dm *Daemon) cullInvalidConnections() {
	// This method only handles the erroneous people from the DHT, but not
	// malicious nodes
	now := utc.Now()
	addrs, err := dm.expectingIntroductions.CullInvalidConns(
		func(addr string, t time.Time) (bool, error) {
			conned, err := dm.Pool.Pool.IsConnExist(addr)
			if err != nil {
				return false, err
			}

			// Do not remove trusted peers
			if dm.isTrustedPeer(addr) {
				return false, nil
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
		logger.Errorf("expectingIntroduction cull invalid connections failed: %v", err)
		return
	}

	for _, a := range addrs {
		exist, err := dm.Pool.Pool.IsConnExist(a)
		if err != nil {
			logger.Error(err)
			return
		}

		if exist {
			logger.Infof("Removing %s for not sending a version", a)
			if err := dm.Pool.Pool.Disconnect(a, ErrDisconnectIntroductionTimeout); err != nil {
				logger.Error(err)
				return
			}
			dm.Pex.RemovePeer(a)
		}
	}
}

func (dm *Daemon) isTrustedPeer(addr string) bool {
	peer, ok := dm.Pex.GetPeerByAddr(addr)
	if !ok {
		return false
	}

	return peer.Trusted
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
		logger.Infof("Connected to peer: %s (outgoing)", a)
	} else {
		logger.Infof("Connected to peer: %s (incoming)", a)
	}

	dm.pendingConnections.Remove(a)

	exist, err := dm.Pool.Pool.IsConnExist(a)
	if err != nil {
		logger.Error(err)
		return
	}

	if !exist {
		logger.Warning("While processing an onConnect event, no pool connection was found")
		return
	}

	if dm.ipCountMaxed(a) {
		logger.Infof("Max connections for %s reached, disconnecting", a)
		dm.Pool.Pool.Disconnect(a, ErrDisconnectIPLimitReached)
		return
	}

	dm.recordIPCount(a)

	if e.Solicited {
		// Disconnect if the max outgoing connections is reached
		n, err := dm.Pool.Pool.OutgoingConnectionsNum()
		if err != nil {
			logger.WithError(err).Error("get outgoing connections number failed")
			return
		}

		if n > dm.Config.OutgoingMax {
			logger.Warningf("max outgoing connections is reached, disconnecting %v", a)
			dm.Pool.Pool.Disconnect(a, ErrDisconnectMaxOutgoingConnectionsReached)
			return
		}

		dm.outgoingConnections.Add(a)
	}

	dm.expectingIntroductions.Add(a, utc.Now())
	logger.Debugf("Sending introduction message to %s, mirror:%d", a, dm.Messages.Mirror)
	m := NewIntroductionMessage(dm.Messages.Mirror, dm.Config.Version, dm.Pool.Pool.Config.Port)
	if err := dm.Pool.Pool.SendMessage(a, m); err != nil {
		logger.Errorf("Send IntroductionMessage to %s failed: %v", a, err)
	}
}

func (dm *Daemon) onDisconnect(e DisconnectEvent) {
	logger.Infof("%s disconnected because: %v", e.Addr, e.Reason)

	dm.outgoingConnections.Remove(e.Addr)
	dm.expectingIntroductions.Remove(e.Addr)
	dm.Heights.Remove(e.Addr)
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
		logger.Warning("onDisconnectEvent channel is full")
	}
}

// Triggered when an gnet.Connection is connected
func (dm *Daemon) onGnetConnect(addr string, solicited bool) {
	dm.onConnectEvent <- ConnectEvent{Addr: addr, Solicited: solicited}
}

// Returns whether the ipCount maximum has been reached
func (dm *Daemon) ipCountMaxed(addr string) bool {
	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Warningf("ipCountMaxed called with invalid addr: %v", err)
		return true
	}

	if cnt, ok := dm.ipCounts.Get(ip); ok {
		return cnt >= dm.Config.IPCountsMax
	}
	return false
}

// Adds base IP to ipCount or returns error if max is reached
func (dm *Daemon) recordIPCount(addr string) {
	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Warningf("recordIPCount called with invalid addr: %v", err)
		return
	}
	dm.ipCounts.Increase(ip)
}

// Removes base IP from ipCount
func (dm *Daemon) removeIPCount(addr string) {
	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Warningf("removeIPCount called with invalid addr: %v", err)
		return
	}
	dm.ipCounts.Decrease(ip)
}

// Adds addr + mirror to the connectionMirror mappings
func (dm *Daemon) recordConnectionMirror(addr string, mirror uint32) error {
	ip, port, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Warningf("recordConnectionMirror called with invalid addr: %v", err)
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
	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Warningf("removeConnectionMirror called with invalid addr: %v", err)
		return
	}

	// remove ip from specific mirror
	dm.mirrorConnections.Remove(mirror, ip)

	dm.connectionMirrors.Remove(addr)
}

// Returns whether an addr+mirror's port and whether the port exists
func (dm *Daemon) getMirrorPort(addr string, mirror uint32) (uint16, bool) {
	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Warningf("getMirrorPort called with invalid addr: %v", err)
		return 0, false
	}
	return dm.mirrorConnections.Get(mirror, ip)
}

// When an async message send finishes, its result is handled by this.
// This method must take care to perform only thread-safe actions, since it is called
// outside of the daemon run loop
func (dm *Daemon) handleMessageSendResult(r gnet.SendResult) {
	if r.Error != nil {
		logger.Warningf("Failed to send %s to %s: %v", reflect.TypeOf(r.Message), r.Addr, r.Error)
		return
	}
	switch r.Message.(type) {
	case SendingTxnsMessage:
		dm.announcedTxns.add(r.Message.(SendingTxnsMessage).GetTxns())
	default:
	}
}

// RequestBlocks Sends a GetBlocksMessage to all connections
func (dm *Daemon) RequestBlocks() error {
	if dm.Config.DisableOutgoingConnections {
		return nil
	}

	headSeq, ok, err := dm.Visor.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot request blocks, there is no head block")
	}

	m := NewGetBlocksMessage(headSeq, dm.Config.BlocksResponseCount)

	err = dm.Pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Debugf("Broadcast GetBlocksMessage failed: %v", err)
	}

	return err
}

// AnnounceBlocks sends an AnnounceBlocksMessage to all connections
func (dm *Daemon) AnnounceBlocks() error {
	if dm.Config.DisableOutgoingConnections {
		return nil
	}

	headSeq, ok, err := dm.Visor.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot announce blocks, there is no head block")
	}

	m := NewAnnounceBlocksMessage(headSeq)

	err = dm.Pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Debugf("Broadcast AnnounceBlocksMessage failed: %v", err)
	}

	return err
}

// AnnounceAllTxns announces local unconfirmed transactions
func (dm *Daemon) AnnounceAllTxns() error {
	if dm.Config.DisableOutgoingConnections {
		return nil
	}

	// Get local unconfirmed transaction hashes.
	hashes, err := dm.Visor.GetAllValidUnconfirmedTxHashes()
	if err != nil {
		return err
	}

	// Divide hashes into multiple sets of max size
	hashesSet := divideHashes(hashes, dm.Config.MaxTxnAnnounceNum)

	for _, hs := range hashesSet {
		m := NewAnnounceTxnsMessage(hs)
		if err = dm.Pool.Pool.BroadcastMessage(m); err != nil {
			break
		}
	}

	if err != nil {
		logger.Debugf("Broadcast AnnounceTxnsMessage failed, err: %v", err)
	}

	return err
}

func divideHashes(hashes []cipher.SHA256, n int) [][]cipher.SHA256 {
	if len(hashes) == 0 {
		return [][]cipher.SHA256{}
	}

	var j int
	var hashesArray [][]cipher.SHA256

	if len(hashes) > n {
		for i := range hashes {
			if len(hashes[j:i]) == n {
				hs := make([]cipher.SHA256, n)
				copy(hs, hashes[j:i])
				hashesArray = append(hashesArray, hs)
				j = i
			}
		}
	}

	hs := make([]cipher.SHA256, len(hashes)-j)
	copy(hs, hashes[j:])
	hashesArray = append(hashesArray, hs)
	return hashesArray
}

// AnnounceTxns announces given transaction hashes.
func (dm *Daemon) AnnounceTxns(txns []cipher.SHA256) error {
	if dm.Config.DisableOutgoingConnections {
		return nil
	}

	if len(txns) == 0 {
		return nil
	}

	m := NewAnnounceTxnsMessage(txns)

	err := dm.Pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Debugf("Broadcast AnnounceTxnsMessage failed: %v", err)
	}

	return err
}

// RequestBlocksFromAddr sends a GetBlocksMessage to one connected address
func (dm *Daemon) RequestBlocksFromAddr(addr string) error {
	if dm.Config.DisableOutgoingConnections {
		return errors.New("Outgoing connections disabled")
	}

	headSeq, ok, err := dm.Visor.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot request blocks from addr, there is no head block")
	}

	m := NewGetBlocksMessage(headSeq, dm.Config.BlocksResponseCount)

	return dm.Pool.Pool.SendMessage(addr, m)
}

// InjectBroadcastTransaction injects transaction to the unconfirmed pool and broadcasts it.
// If the transaction violates either hard or soft constraints, it is not broadcast.
// This method is to be used by user-initiated transaction injections.
// For transactions received over the network, use InjectTransaction and check the result to
// decide on repropagation.
func (dm *Daemon) InjectBroadcastTransaction(txn coin.Transaction) error {
	if _, err := dm.Visor.InjectTransactionStrict(txn); err != nil {
		return err
	}

	return dm.broadcastTransaction(txn)
}

// ResendUnconfirmedTxns resends all unconfirmed transactions and returns the hashes that were successfully rebroadcast
func (dm *Daemon) ResendUnconfirmedTxns() ([]cipher.SHA256, error) {
	if dm.Config.DisableOutgoingConnections {
		return nil, nil
	}

	txns, err := dm.Visor.GetAllUnconfirmedTxns()
	if err != nil {
		return nil, err
	}

	var txids []cipher.SHA256
	for i := range txns {
		logger.Debugf("Rebroadcast tx %s", txns[i].Hash().Hex())
		if err := dm.broadcastTransaction(txns[i].Txn); err == nil {
			txids = append(txids, txns[i].Txn.Hash())
		}
	}

	return txids, nil
}

// broadcastTransaction broadcasts a single transaction to all peers.
func (dm *Daemon) broadcastTransaction(t coin.Transaction) error {
	if dm.Config.DisableOutgoingConnections {
		return nil
	}

	m := NewGiveTxnsMessage(coin.Transactions{t})
	l, err := dm.Pool.Pool.Size()
	if err != nil {
		return err
	}

	logger.Debugf("Broadcasting GiveTxnsMessage to %d conns", l)

	err = dm.Pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Errorf("Broadcast GivenTxnsMessage failed: %v", err)
	}

	return err
}

// CreateAndPublishBlock creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a master chain.
// Will not create a block if outgoing connections are disabled.
// If the block was created but the broadcast failed, the error will be non-nil but the
// SignedBlock value will not be empty.
// TODO -- refactor this method -- it should either always create a block and maybe broadcast it,
// or use a database transaction to rollback block publishing if broadcast failed (however, this will cause a slow DB write)
func (dm *Daemon) CreateAndPublishBlock() (*coin.SignedBlock, error) {
	if dm.Config.DisableOutgoingConnections {
		return nil, errors.New("Outgoing connections disabled")
	}

	sb, err := dm.Visor.CreateAndExecuteBlock()
	if err != nil {
		return nil, err
	}

	err = dm.broadcastBlock(sb)

	return &sb, err
}

// Sends a signed block to all connections.
func (dm *Daemon) broadcastBlock(sb coin.SignedBlock) error {
	if dm.Config.DisableOutgoingConnections {
		return nil
	}

	m := NewGiveBlocksMessage([]coin.SignedBlock{sb})
	return dm.Pool.Pool.BroadcastMessage(m)
}
