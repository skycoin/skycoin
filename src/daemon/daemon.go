/*
Package daemon controls the networking layer of the skycoin daemon
*/
package daemon

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/elapse"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/iputil"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/useragent"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// ErrNetworkingDisabled is returned if networking is disabled
	ErrNetworkingDisabled = errors.New("Networking is disabled")
	// ErrNoPeerAcceptsTxn is returned if no peer will propagate a transaction broadcasted with BroadcastUserTransaction
	ErrNoPeerAcceptsTxn = errors.New("No peer will propagate this transaction")

	logger = logging.MustGetLogger("daemon")
)

// IsBroadcastFailure returns true if an error indicates that a broadcast operation failed
func IsBroadcastFailure(err error) bool {
	switch err {
	case ErrNetworkingDisabled,
		gnet.ErrPoolEmpty,
		gnet.ErrNoMatchingConnections,
		gnet.ErrNoReachableConnections,
		gnet.ErrNoAddresses:
		return true
	default:
		return false
	}
}

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
		Visor:    visor.NewConfig(),
	}
}

// preprocess preprocess for config
func (cfg *Config) preprocess() (Config, error) {
	config := *cfg
	if config.Daemon.LocalhostOnly {
		if config.Daemon.Address == "" {
			local, err := iputil.LocalhostIP()
			if err != nil {
				logger.WithError(err).Panic("Failed to obtain localhost IP")
			}
			config.Daemon.Address = local
		} else {
			if !iputil.IsLocalhost(config.Daemon.Address) {
				logger.WithField("addr", config.Daemon.Address).Panic("Invalid address for localhost-only")
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

	if config.Daemon.MaxConnections < config.Daemon.MaxOutgoingConnections {
		return Config{}, errors.New("MaxOutgoingConnections cannot be more than MaxConnections")
	}

	if config.Daemon.MaxPendingConnections > config.Daemon.MaxOutgoingConnections {
		config.Daemon.MaxPendingConnections = config.Daemon.MaxOutgoingConnections
	}

	config.Pool.MaxConnections = config.Daemon.MaxConnections
	config.Pool.MaxOutgoingConnections = config.Daemon.MaxOutgoingConnections

	userAgent, err := config.Daemon.UserAgent.Build()
	if err != nil {
		return Config{}, err
	}
	if userAgent == "" {
		return Config{}, errors.New("user agent is required")
	}
	config.Daemon.userAgent = userAgent

	return config, nil
}

// DaemonConfig configuration for the Daemon
type DaemonConfig struct { // nolint: golint
	// Protocol version. TODO -- manage version better
	ProtocolVersion int32
	// Minimum accepted protocol version
	MinProtocolVersion int32
	// IP Address to serve on. Leave empty for automatic assignment
	Address string
	// BlockchainPubkey blockchain pubkey string
	BlockchainPubkey cipher.PubKey
	// TCP/UDP port for connections
	Port int
	// Directory where application data is stored
	DataDirectory string
	// How often to check and initiate an outgoing connection to a trusted connection if needed
	OutgoingTrustedRate time.Duration
	// How often to check and initiate an outgoing connection if needed
	OutgoingRate time.Duration
	// How often to re-attempt to fill any missing private (aka required)  connections
	PrivateRate time.Duration
	// Maximum number of connections
	MaxConnections int
	// Number of outgoing connections to maintain
	MaxOutgoingConnections int
	// Maximum number of connections to try at once
	MaxPendingConnections int
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
	// Default "trusted" peers
	DefaultConnections []string
	// User agent (sent in introduction messages)
	UserAgent useragent.Data
	userAgent string // parsed from UserAgent in preprocess()
	// Transaction verification parameters for unconfirmed transactions
	UnconfirmedVerifyTxn params.VerifyTxn
	// Random nonce value for detecting self-connection in introduction messages
	Mirror uint32
}

// NewDaemonConfig creates daemon config
func NewDaemonConfig() DaemonConfig {
	return DaemonConfig{
		ProtocolVersion:              2,
		MinProtocolVersion:           2,
		Address:                      "",
		Port:                         6677,
		OutgoingRate:                 time.Second * 5,
		OutgoingTrustedRate:          time.Millisecond * 100,
		PrivateRate:                  time.Second * 5,
		MaxConnections:               128,
		MaxOutgoingConnections:       8,
		MaxPendingConnections:        8,
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
		Mirror:                       rand.New(rand.NewSource(time.Now().UTC().UnixNano())).Uint32(),
		UnconfirmedVerifyTxn:         params.UserVerifyTxn,
	}
}

//go:generate go install
//go:generate mockery -name daemoner -case underscore -inpkg -testonly

// daemoner Daemon interface
type daemoner interface {
	Disconnect(addr string, r gnet.DisconnectReason) error
	sendMessage(addr string, msg gnet.Message) error
	broadcastMessage(msg gnet.Message) ([]uint64, error)
	disconnectNow(addr string, r gnet.DisconnectReason) error
	addPeers(addrs []string) int
	recordPeerHeight(addr string, gnetID, height uint64)
	getSignedBlocksSince(seq, count uint64) ([]coin.SignedBlock, error)
	headBkSeq() (uint64, bool, error)
	executeSignedBlock(b coin.SignedBlock) error
	filterKnownUnconfirmed(txns []cipher.SHA256) ([]cipher.SHA256, error)
	getKnownUnconfirmed(txns []cipher.SHA256) (coin.Transactions, error)
	requestBlocksFromAddr(addr string) error
	announceAllValidTxns() error
	daemonConfig() DaemonConfig
	pexConfig() pex.Config
	injectTransaction(txn coin.Transaction) (bool, *visor.ErrTxnViolatesSoftConstraint, error)
	recordMessageEvent(m asyncMessage, c *gnet.MessageContext) error
	connectionIntroduced(addr string, gnetID uint64, m *IntroductionMessage) (*connection, error)
	sendRandomPeers(addr string) error
}

// Daemon stateful properties of the daemon
type Daemon struct {
	// Daemon configuration
	Config DaemonConfig

	// Components
	Messages *Messages
	pool     *Pool
	pex      *pex.Pex
	Gateway  *Gateway
	visor    *visor.Visor

	// Cache of announced transactions that are flushed to the database periodically
	announcedTxns *announcedTxnsCache
	// Cache of connection metadata
	connections *Connections
	// connect, disconnect, message, error events channel
	events chan interface{}
	// quit channel
	quit chan struct{}
	// done channel
	done chan struct{}
}

// NewDaemon returns a Daemon with primitives allocated
func NewDaemon(config Config, db *dbutil.DB) (*Daemon, error) {
	config, err := config.preprocess()
	if err != nil {
		return nil, err
	}

	vs, err := visor.NewVisor(config.Visor, db)
	if err != nil {
		return nil, err
	}

	pex, err := pex.New(config.Pex)
	if err != nil {
		return nil, err
	}

	d := &Daemon{
		Config:   config.Daemon,
		Messages: NewMessages(config.Messages),
		pex:      pex,
		visor:    vs,

		announcedTxns: newAnnouncedTxnsCache(),
		connections:   NewConnections(),
		events:        make(chan interface{}, config.Pool.EventChannelSize),
		quit:          make(chan struct{}),
		done:          make(chan struct{}),
	}

	d.pool, err = NewPool(config.Pool, d)
	if err != nil {
		return nil, err
	}

	d.Gateway = NewGateway(config.Gateway, d)
	d.Messages.Config.Register()

	return d, nil
}

// ConnectEvent generated when a client connects
type ConnectEvent struct {
	GnetID    uint64
	Addr      string
	Solicited bool
}

// DisconnectEvent generated when a connection terminated
type DisconnectEvent struct {
	GnetID uint64
	Addr   string
	Reason gnet.DisconnectReason
}

// ConnectFailureEvent represent a failure to connect/dial a connection, with context
type ConnectFailureEvent struct {
	Addr      string
	Solicited bool
	Error     error
}

// messageEvent encapsulates a deserialized message from the network
type messageEvent struct {
	Message asyncMessage
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
	dm.pool.Shutdown()

	logger.Info("Shutting down Gateway")
	dm.Gateway.Shutdown()

	logger.Info("Shutting down Pex")
	dm.pex.Shutdown()

	<-dm.done
}

// Init prepares daemon before Run()
func (dm *Daemon) Init() error {
	if err := dm.visor.Init(); err != nil {
		logger.WithError(err).Error("visor.Visor.Init failed")
		return err
	}

	return nil
}

// Run main loop for peer/connection management
func (dm *Daemon) Run() error {
	defer logger.Info("Daemon closed")
	defer close(dm.done)

	logger.Infof("Daemon UserAgent is %s", dm.Config.userAgent)
	logger.Info("Daemon unconfirmed BurnFactor is %d", dm.Config.UnconfirmedVerifyTxn.BurnFactor)
	logger.Info("Daemon unconfirmed MaxTransactionSize is %d", dm.Config.UnconfirmedVerifyTxn.MaxTransactionSize)
	logger.Info("Daemon unconfirmed MaxDropletPrecision is %d", dm.Config.UnconfirmedVerifyTxn.MaxDropletPrecision)

	errC := make(chan error, 5)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := dm.pex.Run(); err != nil {
			logger.WithError(err).Error("daemon.Pex.Run failed")
			errC <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if dm.Config.DisableIncomingConnections {
			if err := dm.pool.RunOffline(); err != nil {
				logger.WithError(err).Error("daemon.Pool.RunOffline failed")
				errC <- err
			}
		} else {
			if err := dm.pool.Run(); err != nil {
				logger.WithError(err).Error("daemon.Pool.Run failed")
				errC <- err
			}
		}
	}()

	blockInterval := time.Duration(dm.Config.BlockCreationInterval)
	blockCreationTicker := time.NewTicker(time.Second * blockInterval)
	if !dm.visor.Config.IsBlockPublisher {
		blockCreationTicker.Stop()
	}

	unconfirmedRefreshTicker := time.NewTicker(dm.Config.UnconfirmedRefreshRate)
	defer unconfirmedRefreshTicker.Stop()
	unconfirmedRemoveInvalidTicker := time.NewTicker(dm.Config.UnconfirmedRemoveInvalidRate)
	defer unconfirmedRemoveInvalidTicker.Stop()
	blocksRequestTicker := time.NewTicker(dm.Config.BlocksRequestRate)
	defer blocksRequestTicker.Stop()
	blocksAnnounceTicker := time.NewTicker(dm.Config.BlocksAnnounceRate)
	defer blocksAnnounceTicker.Stop()

	// outgoingTrustedConnectionsTicker is used to maintain at least one connection to a trusted peer.
	// This may be configured at a very frequent rate, so if no trusted connections could be reached,
	// there could be a lot of churn.
	// The additional outgoingTrustedConnectionsTicker parameters are used to
	// skip ticks of the outgoingTrustedConnectionsTicker in the event of total failure.
	// outgoingTrustedConnectionsTickerSkipDuration is the minimum time to wait between
	// ticks in the event of total failure.
	outgoingTrustedConnectionsTicker := time.NewTicker(dm.Config.OutgoingTrustedRate)
	defer outgoingTrustedConnectionsTicker.Stop()
	outgoingTrustedConnectionsTickerSkipDuration := time.Second * 5
	outgoingTrustedConnectionsTickerSkip := false
	var outgoingTrustedConnectionsTickerSkipStart time.Time

	privateConnectionsTicker := time.NewTicker(dm.Config.PrivateRate)
	defer privateConnectionsTicker.Stop()
	cullInvalidTicker := time.NewTicker(dm.Config.CullInvalidRate)
	defer cullInvalidTicker.Stop()
	outgoingConnectionsTicker := time.NewTicker(dm.Config.OutgoingRate)
	defer outgoingConnectionsTicker.Stop()
	requestPeersTicker := time.NewTicker(dm.pex.Config.RequestRate)
	defer requestPeersTicker.Stop()
	clearStaleConnectionsTicker := time.NewTicker(dm.pool.Config.ClearStaleRate)
	defer clearStaleConnectionsTicker.Stop()
	idleCheckTicker := time.NewTicker(dm.pool.Config.IdleCheckRate)
	defer idleCheckTicker.Stop()

	flushAnnouncedTxnsTicker := time.NewTicker(dm.Config.FlushAnnouncedTxnsRate)
	defer flushAnnouncedTxnsTicker.Stop()

	// Connect to all trusted peers on startup to try to ensure a connection establishes quickly.
	// The number of connections to default peers is restricted;
	// if multiple connections succeed, extra connections beyond the limit will be disconnected.
	if !dm.Config.DisableOutgoingConnections {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dm.connectToTrustedPeers()
		}()
	}

	var setupErr error
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

			case r := <-dm.pool.Pool.SendResults:
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

		case <-cullInvalidTicker.C:
			// Remove connections that failed to complete the handshake
			elapser.Register("cullInvalidTicker")
			if !dm.Config.DisableNetworking {
				dm.cullInvalidConnections()
			}

		case <-requestPeersTicker.C:
			// Request peers via PEX
			elapser.Register("requestPeersTicker")
			if dm.pex.Config.Disabled {
				continue
			}

			if dm.pex.IsFull() {
				continue
			}

			m := NewGetPeersMessage()
			if _, err := dm.broadcastMessage(m); err != nil {
				logger.WithError(err).Error("Broadcast GetPeersMessage failed")
				continue
			}

		case <-clearStaleConnectionsTicker.C:
			// Remove connections that haven't said anything in a while
			elapser.Register("clearStaleConnectionsTicker")
			if !dm.Config.DisableNetworking {
				conns, err := dm.pool.getStaleConnections()
				if err != nil {
					logger.WithError(err).Error("getStaleConnections failed")
					continue
				}

				for _, addr := range conns {
					if err := dm.Disconnect(addr, ErrDisconnectIdle); err != nil {
						logger.WithError(err).WithField("addr", addr).Error("Disconnect")
					}
				}
			}

		case <-idleCheckTicker.C:
			// Sends pings as needed
			elapser.Register("idleCheckTicker")
			if !dm.Config.DisableNetworking {
				dm.pool.sendPings()
			}

		case <-outgoingConnectionsTicker.C:
			// Fill up our outgoing connections
			elapser.Register("outgoingConnectionsTicker")
			dm.connectToRandomPeer()

		case <-outgoingTrustedConnectionsTicker.C:
			// Try to maintain at least one trusted connection
			elapser.Register("outgoingTrustedConnectionsTicker")
			// If connecting to a trusted peer totally fails, make sure to wait longer between further attempts
			if outgoingTrustedConnectionsTickerSkip {
				if time.Since(outgoingTrustedConnectionsTickerSkipStart) < outgoingTrustedConnectionsTickerSkipDuration {
					continue
				}
			}

			if err := dm.maybeConnectToTrustedPeer(); err != nil && err != ErrNetworkingDisabled {
				logger.Critical().WithError(err).Error("maybeConnectToTrustedPeer")
				outgoingTrustedConnectionsTickerSkip = true
				outgoingTrustedConnectionsTickerSkipStart = time.Now()
			} else {
				outgoingTrustedConnectionsTickerSkip = false
			}

		case <-privateConnectionsTicker.C:
			// Always try to stay connected to our private peers
			// TODO (also, connect to all of them on start)
			elapser.Register("privateConnectionsTicker")
			if !dm.Config.DisableOutgoingConnections {
				dm.makePrivateConnections()
			}

		case r := <-dm.events:
			elapser.Register("dm.event")
			if dm.Config.DisableNetworking {
				logger.Critical().Error("Networking is disabled, there should be no events")
			} else {
				dm.handleEvent(r)
			}

		case <-flushAnnouncedTxnsTicker.C:
			elapser.Register("flushAnnouncedTxnsTicker")
			txns := dm.announcedTxns.flush()

			if err := dm.visor.SetTransactionsAnnounced(txns); err != nil {
				logger.WithError(err).Error("Failed to set unconfirmed txn announce time")
			}

		case req := <-dm.Gateway.requests:
			// Process any pending RPC requests
			elapser.Register("dm.Gateway.requests")
			if err := req.Func(); err != nil {
				logger.WithError(err).Error()
			}

		case <-blockCreationTicker.C:
			// Create blocks, if block publisher
			elapser.Register("blockCreationTicker.C")
			if dm.visor.Config.IsBlockPublisher {
				sb, err := dm.createAndPublishBlock()
				if err != nil {
					logger.WithError(err).Error("Failed to create and publish block")
					continue
				}

				// Not a critical error, but we want it visible in logs
				head := sb.Block.Head
				logger.Critical().WithFields(logrus.Fields{
					"version": head.Version,
					"seq":     head.BkSeq,
					"time":    head.Time,
				}).Info("Created and published a new block")
			}

		case <-unconfirmedRefreshTicker.C:
			elapser.Register("unconfirmedRefreshTicker")
			// Get the transactions that turn to valid
			validTxns, err := dm.visor.RefreshUnconfirmed()
			if err != nil {
				logger.WithError(err).Error("dm.Visor.RefreshUnconfirmed failed")
				continue
			}
			// Announce these transactions
			if err := dm.announceTxns(validTxns); err != nil {
				logger.WithError(err).Warning("announceTxns failed")
			}

		case <-unconfirmedRemoveInvalidTicker.C:
			elapser.Register("unconfirmedRemoveInvalidTicker")
			// Remove transactions that become invalid (violating hard constraints)
			removedTxns, err := dm.visor.RemoveInvalidUnconfirmed()
			if err != nil {
				logger.WithError(err).Error("dm.Visor.RemoveInvalidUnconfirmed failed")
				continue
			}
			if len(removedTxns) > 0 {
				logger.Infof("Remove %d txns from pool that began violating hard constraints", len(removedTxns))
			}

		case <-blocksRequestTicker.C:
			elapser.Register("blocksRequestTicker")
			if err := dm.requestBlocks(); err != nil {
				logger.WithError(err).Warning("requestBlocks failed")
			}

		case <-blocksAnnounceTicker.C:
			elapser.Register("blocksAnnounceTicker")
			if err := dm.announceBlocks(); err != nil {
				logger.WithError(err).Warning("announceBlocks failed")
			}

		case setupErr = <-errC:
			logger.WithError(setupErr).Error("read from errc")
			break loop
		}
	}

	if setupErr != nil {
		return setupErr
	}

	wg.Wait()

	return nil
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
		logger.Critical().WithField("addr", p.Addr).WithError(err).Warning("PEX gave us an invalid peer")
		return errors.New("Invalid peer")
	}

	if dm.Config.LocalhostOnly && !iputil.IsLocalhost(a) {
		return errors.New("Not localhost")
	}

	if c := dm.connections.get(p.Addr); c != nil {
		return errors.New("Already connected to this peer")
	}

	cnt := dm.connections.IPCount(a)
	if !dm.Config.LocalhostOnly && cnt != 0 {
		return errors.New("Already connected to a peer with this base IP")
	}

	logger.WithField("addr", p.Addr).Debug("Establishing outgoing connection")

	if _, err := dm.connections.pending(p.Addr); err != nil {
		logger.Critical().WithError(err).WithField("addr", p.Addr).Error("dm.connections.pending failed")
		return err
	}

	go func() {
		if err := dm.pool.Pool.Connect(p.Addr); err != nil {
			dm.events <- ConnectFailureEvent{
				Addr:      p.Addr,
				Solicited: true,
				Error:     err,
			}
		}
	}()
	return nil
}

// Connects to all private peers
func (dm *Daemon) makePrivateConnections() {
	if dm.Config.DisableOutgoingConnections {
		return
	}

	peers := dm.pex.Private()
	for _, p := range peers {
		logger.WithField("addr", p.Addr).Info("Private peer attempt")
		if err := dm.connectToPeer(p); err != nil {
			logger.WithField("addr", p.Addr).WithError(err).Debug("Did not connect to private peer")
		}
	}
}

// connectToTrustedPeers tries to connect to all trusted peers
func (dm *Daemon) connectToTrustedPeers() {
	if dm.Config.DisableOutgoingConnections {
		return
	}

	logger.Info("Connect to trusted peers")
	// Make connections to all trusted peers to try to ensure a connection
	// MaxDefaultPeerOutgoingConnections limits will be enforced in gnet
	// after connections have been established, so not all trusted peers will be connected to.
	peers := dm.pex.TrustedPublic()
	for _, p := range peers {
		if err := dm.connectToPeer(p); err != nil {
			logger.WithError(err).WithField("addr", p.Addr).Warning("connect to trusted peer failed")
		}
	}
}

// maybeConnectToTrustedPeer tries to connect to one trusted peer if there are no trusted connections
func (dm *Daemon) maybeConnectToTrustedPeer() error {
	if dm.Config.DisableOutgoingConnections {
		return ErrNetworkingDisabled
	}

	peers := dm.pex.TrustedPublic()
	for _, p := range peers {
		// Don't make a connection if we have a trusted peer connection
		if len(dm.connections.getByListenAddr(p.Addr)) != 0 {
			return nil
		}
	}

	connected := false
	for _, p := range peers {
		if err := dm.connectToPeer(p); err != nil {
			logger.WithError(err).WithField("addr", p.Addr).Warning("maybeConnectToTrustedPeer: connectToPeer failed")
			continue
		}
		connected = true
		break
	}

	if !connected {
		return errors.New("Could not connect to any trusted peer")
	}

	return nil
}

// connectToRandomPeer attempts to connect to a random peer. If it fails, the peer is removed.
func (dm *Daemon) connectToRandomPeer() {
	if dm.Config.DisableOutgoingConnections {
		return
	}
	if dm.connections.OutgoingLen() >= dm.Config.MaxOutgoingConnections {
		return
	}
	if dm.connections.PendingLen() >= dm.Config.MaxPendingConnections {
		return
	}
	if dm.connections.Len() >= dm.Config.MaxConnections {
		return
	}

	// Make a connection to a random (public) peer
	peers := dm.pex.RandomPublic(dm.Config.MaxOutgoingConnections - dm.connections.OutgoingLen())
	for _, p := range peers {
		if err := dm.connectToPeer(p); err != nil {
			logger.WithError(err).WithField("addr", p.Addr).Warning("connectToPeer failed")
		}
	}

	// TODO -- don't reset if not needed?
	if len(peers) == 0 {
		dm.pex.ResetAllRetryTimes()
	}
}

// Removes connections who haven't sent a version after connecting
func (dm *Daemon) cullInvalidConnections() {
	now := time.Now().UTC()
	for _, c := range dm.connections.all() {
		if c.State != ConnectionStateConnected {
			continue
		}

		if c.ConnectedAt.Add(dm.Config.IntroductionWait).Before(now) {
			logger.WithField("addr", c.Addr).Info("Disconnecting peer for not sending a version")
			if err := dm.Disconnect(c.Addr, ErrDisconnectIntroductionTimeout); err != nil {
				logger.WithError(err).WithField("addr", c.Addr).Error("Disconnect")
			}
		}
	}
}

func (dm *Daemon) isTrustedPeer(addr string) bool {
	peer, ok := dm.pex.GetPeer(addr)
	if !ok {
		return false
	}

	return peer.Trusted
}

// recordMessageEvent records an asyncMessage to the messageEvent chan.  Do not access
// messageEvent directly.
func (dm *Daemon) recordMessageEvent(m asyncMessage, c *gnet.MessageContext) error {
	dm.events <- messageEvent{
		Message: m,
		Context: c,
	}
	return nil
}

func (dm *Daemon) handleEvent(e interface{}) {
	switch x := e.(type) {
	case messageEvent:
		dm.onMessageEvent(x)
	case ConnectEvent:
		dm.onConnectEvent(x)
	case DisconnectEvent:
		dm.onDisconnectEvent(x)
	case ConnectFailureEvent:
		dm.onConnectFailure(x)
	default:
		logger.WithFields(logrus.Fields{
			"type":  fmt.Sprintf("%T", e),
			"value": fmt.Sprintf("%+v", e),
		}).Panic("Invalid object in events queue")
	}
}

func (dm *Daemon) onMessageEvent(e messageEvent) {
	// If the connection does not exist or the gnet ID is different, abort message processing
	// This can occur because messageEvents for a given connection may occur
	// after that connection has disconnected.
	c := dm.connections.get(e.Context.Addr)
	if c == nil {
		logger.WithFields(logrus.Fields{
			"addr":        e.Context.Addr,
			"messageType": fmt.Sprintf("%T", e.Message),
		}).Info("onMessageEvent no connection found")
		return
	}

	if c.gnetID != e.Context.ConnID {
		logger.WithFields(logrus.Fields{
			"addr":          e.Context.Addr,
			"connGnetID":    c.gnetID,
			"contextGnetID": e.Context.ConnID,
			"messageType":   fmt.Sprintf("%T", e.Message),
		}).Info("onMessageEvent connection gnetID does not match")
		return
	}

	// The first message received must be INTR, DISC or GIVP
	if !c.HasIntroduced() {
		switch e.Message.(type) {
		case *IntroductionMessage, *DisconnectMessage, *GivePeersMessage:
		default:
			logger.WithFields(logrus.Fields{
				"addr":        e.Context.Addr,
				"messageType": fmt.Sprintf("%T", e.Message),
			}).Info("needsIntro but first message is not INTR, DISC or GIVP")
			if err := dm.Disconnect(e.Context.Addr, ErrDisconnectNoIntroduction); err != nil {
				logger.WithError(err).WithField("addr", e.Context.Addr).Error("Disconnect")
			}
			return
		}
	}

	e.Message.process(dm)
}

func (dm *Daemon) onConnectEvent(e ConnectEvent) {
	fields := logrus.Fields{
		"addr":     e.Addr,
		"outgoing": e.Solicited,
		"gnetID":   e.GnetID,
	}
	logger.WithFields(fields).Info("onConnectEvent")

	// Update the connections state machine first
	c, err := dm.connections.connected(e.Addr, e.GnetID)
	if err != nil {
		logger.Critical().WithError(err).WithFields(fields).Error("connections.Connected failed")
		if err := dm.Disconnect(e.Addr, ErrDisconnectUnexpectedError); err != nil {
			logger.WithError(err).WithFields(fields).Error("Disconnect")
		}
		return
	}

	// The connection should already be known as outgoing/solicited due to an earlier connections.pending call.
	// If they do not match, there is e.Addr flaw in the concept or implementation of the state machine.
	if c.Outgoing != e.Solicited {
		logger.Critical().WithFields(fields).Warning("Connection.Outgoing does not match ConnectEvent.Solicited state")
	}

	if dm.ipCountMaxed(e.Addr) {
		logger.WithFields(fields).Info("Max connections for this IP address reached, disconnecting")
		if err := dm.Disconnect(e.Addr, ErrDisconnectIPLimitReached); err != nil {
			logger.WithError(err).WithFields(fields).Error("Disconnect")
		}
		return
	}

	logger.WithFields(fields).Debug("Sending introduction message")

	if err := dm.sendMessage(e.Addr, NewIntroductionMessage(
		dm.Config.Mirror,
		dm.Config.ProtocolVersion,
		dm.pool.Pool.Config.Port,
		dm.Config.BlockchainPubkey,
		dm.Config.userAgent,
		dm.Config.UnconfirmedVerifyTxn,
	)); err != nil {
		logger.WithFields(fields).WithError(err).Error("Send IntroductionMessage failed")
		return
	}
}

func (dm *Daemon) onDisconnectEvent(e DisconnectEvent) {
	fields := logrus.Fields{
		"addr":   e.Addr,
		"reason": e.Reason,
		"gnetID": e.GnetID,
	}
	logger.WithFields(fields).Info("onDisconnectEvent")

	if err := dm.connections.remove(e.Addr, e.GnetID); err != nil {
		logger.WithError(err).WithFields(fields).Error("connections.Remove failed")
		return
	}

	// TODO -- blacklist peer for certain reasons, not just remove
	switch e.Reason {
	case ErrDisconnectIntroductionTimeout,
		ErrDisconnectBlockchainPubkeyNotMatched,
		ErrDisconnectInvalidExtraData,
		ErrDisconnectInvalidUserAgent:
		if !dm.isTrustedPeer(e.Addr) {
			dm.pex.RemovePeer(e.Addr)
		}
	case ErrDisconnectNoIntroduction,
		ErrDisconnectVersionNotSupported,
		ErrDisconnectSelf:
		dm.pex.IncreaseRetryTimes(e.Addr)
	default:
		switch e.Reason.Error() {
		case "read failed: EOF":
			dm.pex.IncreaseRetryTimes(e.Addr)
		}
	}
}

func (dm *Daemon) onConnectFailure(c ConnectFailureEvent) {
	// Remove the pending connection from connections and update the retry times in pex
	logger.WithField("addr", c.Addr).WithError(c.Error).Debug("onConnectFailure")

	// onConnectFailure should only trigger for "pending" connections which have gnet ID 0;
	// connections in any other state will have a nonzero gnet ID.
	// if the connection is in a different state, the gnet ID will not match, the connection
	// won't be removed and we'll receive an error.
	// If this happens, it is a bug, and the connections state may be corrupted.
	if err := dm.connections.remove(c.Addr, 0); err != nil {
		logger.Critical().WithField("addr", c.Addr).WithError(err).Error("connections.remove")
	}

	if strings.HasSuffix(c.Error.Error(), "connect: connection refused") {
		dm.pex.IncreaseRetryTimes(c.Addr)
	}
}

// onGnetDisconnect triggered when a gnet.Connection terminates
func (dm *Daemon) onGnetDisconnect(addr string, gnetID uint64, reason gnet.DisconnectReason) {
	dm.events <- DisconnectEvent{
		GnetID: gnetID,
		Addr:   addr,
		Reason: reason,
	}
}

// onGnetConnect Triggered when a gnet.Connection connects
func (dm *Daemon) onGnetConnect(addr string, gnetID uint64, solicited bool) {
	dm.events <- ConnectEvent{
		GnetID:    gnetID,
		Addr:      addr,
		Solicited: solicited,
	}
}

// onGnetConnectFailure triggered when a gnet.Connection fails to connect
func (dm *Daemon) onGnetConnectFailure(addr string, solicited bool, err error) {
	dm.events <- ConnectFailureEvent{
		Addr:      addr,
		Solicited: solicited,
		Error:     err,
	}
}

// Returns whether the ipCount maximum has been reached.
// Always false when using LocalhostOnly config.
func (dm *Daemon) ipCountMaxed(addr string) bool {
	ip, _, err := iputil.SplitAddr(addr)
	if err != nil {
		logger.Critical().WithField("addr", addr).Error("ipCountMaxed called with invalid addr")
		return true
	}

	return !dm.Config.LocalhostOnly && dm.connections.IPCount(ip) >= dm.Config.IPCountsMax
}

// When an async message send finishes, its result is handled by this.
// This method must take care to perform only thread-safe actions, since it is called
// outside of the daemon run loop
func (dm *Daemon) handleMessageSendResult(r gnet.SendResult) {
	if r.Error != nil {
		logger.WithError(r.Error).WithFields(logrus.Fields{
			"addr":    r.Addr,
			"msgType": reflect.TypeOf(r.Message),
		}).Warning("Failed to send message")
		return
	}

	if m, ok := r.Message.(SendingTxnsMessage); ok {
		dm.announcedTxns.add(m.GetFiltered())
	}

	if m, ok := r.Message.(*DisconnectMessage); ok {
		if err := dm.disconnectNow(r.Addr, m.reason); err != nil {
			logger.WithError(err).WithField("addr", r.Addr).Warning("disconnectNow")
		}
	}
}

// requestBlocks sends a GetBlocksMessage to all connections
func (dm *Daemon) requestBlocks() error {
	if dm.Config.DisableNetworking {
		return ErrNetworkingDisabled
	}

	headSeq, ok, err := dm.visor.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot request blocks, there is no head block")
	}

	m := NewGetBlocksMessage(headSeq, dm.Config.BlocksResponseCount)

	if _, err := dm.broadcastMessage(m); err != nil {
		logger.WithError(err).Debug("Broadcast GetBlocksMessage failed")
		return err
	}

	return nil
}

// announceBlocks sends an AnnounceBlocksMessage to all connections
func (dm *Daemon) announceBlocks() error {
	if dm.Config.DisableNetworking {
		return ErrNetworkingDisabled
	}

	headSeq, ok, err := dm.visor.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot announce blocks, there is no head block")
	}

	m := NewAnnounceBlocksMessage(headSeq)

	if _, err := dm.broadcastMessage(m); err != nil {
		logger.WithError(err).Debug("Broadcast AnnounceBlocksMessage failed")
		return err
	}

	return nil
}

// announceTxns announces given transaction hashes
func (dm *Daemon) announceTxns(txns []cipher.SHA256) error {
	if dm.Config.DisableNetworking {
		return ErrNetworkingDisabled
	}

	if len(txns) == 0 {
		return nil
	}

	m := NewAnnounceTxnsMessage(txns)

	if _, err := dm.broadcastMessage(m); err != nil {
		logger.WithError(err).Debug("Broadcast AnnounceTxnsMessage failed")
		return err
	}

	return nil
}

// createAndPublishBlock creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a block publisher.
// Will not create a block if outgoing connections are disabled.
// If the block was created but the broadcast failed, the error will be non-nil but the
// SignedBlock value will not be empty.
// TODO -- refactor this method -- it should either always create a block and maybe broadcast it,
// or use a database transaction to rollback block publishing if broadcast failed (however, this will cause a slow DB write)
func (dm *Daemon) createAndPublishBlock() (*coin.SignedBlock, error) {
	if dm.Config.DisableNetworking {
		return nil, ErrNetworkingDisabled
	}

	sb, err := dm.visor.CreateAndExecuteBlock()
	if err != nil {
		return nil, err
	}

	err = dm.broadcastBlock(sb)

	return &sb, err
}

// ResendUnconfirmedTxns resends all unconfirmed transactions and returns the hashes that were successfully rebroadcast.
// It does not return an error if broadcasting fails.
func (dm *Daemon) ResendUnconfirmedTxns() ([]cipher.SHA256, error) {
	if dm.Config.DisableNetworking {
		return nil, ErrNetworkingDisabled
	}

	txns, err := dm.visor.GetAllUnconfirmedTransactions()
	if err != nil {
		return nil, err
	}

	var txids []cipher.SHA256
	for i := range txns {
		logger.WithField("txid", txns[i].Hash().Hex()).Debug("Rebroadcast transaction")
		if _, err := dm.BroadcastTransaction(txns[i].Transaction); err == nil {
			txids = append(txids, txns[i].Transaction.Hash())
		}
	}

	return txids, nil
}

// BroadcastTransaction broadcasts a single transaction to all peers.
func (dm *Daemon) BroadcastTransaction(txn coin.Transaction) ([]uint64, error) {
	if dm.Config.DisableNetworking {
		return nil, ErrNetworkingDisabled
	}

	m := NewGiveTxnsMessage(coin.Transactions{txn})
	ids, err := dm.broadcastMessage(m)
	if err != nil {
		logger.WithError(err).Error("Broadcast GiveTxnsMessage failed")
		return nil, err
	}

	logger.Debugf("BroadcastTransaction to %d conns", len(ids))

	return ids, nil
}

// BroadcastUserTransaction broadcasts a single transaction to all peers.
// Returns an error if no peers that would propagate the transaction could be reached.
func (dm *Daemon) BroadcastUserTransaction(txn coin.Transaction, head *coin.SignedBlock, inputs coin.UxArray) error {
	ids, err := dm.BroadcastTransaction(txn)
	if err != nil {
		return err
	}

	accepts, err := checkBroadcastTxnRecipients(dm.connections, ids, txn, head, inputs)

	if err != nil {
		logger.WithError(err).Error("BroadcastUserTransaction")
		return err
	}

	logger.Debugf("BroadcastUserTransaction transaction propagated by %d/%d conns", accepts, len(ids))

	return nil
}

// checkBroadcastTxnRecipients checks whether or not the recipients of a txn broadcast would accept the transaction as valid,
// based upon their reported txn verification parameters.
// If no recipient would accept the txn, an error is returned.
// The number of recipients that claim to accept the transaction is returned.
func checkBroadcastTxnRecipients(connections *Connections, ids []uint64, txn coin.Transaction, head *coin.SignedBlock, inputs coin.UxArray) (int, error) {
	// Check if the connections will accept our transaction as valid.
	// Clients v24 and earlier do not propagate soft-invalid transactions.
	// Clients v24 and earlier do not advertise a user agent.
	// Clients v24 and earlier do not advertise their transaction verification parameters,
	// but will use defaults of BurnFactor=2, MaxTransactionSize=32768, MaxDropletPrecision=3.
	// If none of the connections will propagate our transaction, return an error.
	accepts := 0

	for _, id := range ids {
		c := connections.getByGnetID(id)
		if c == nil {
			continue
		}

		if !c.HasIntroduced() {
			continue
		}

		// If the peer has not set their user agent, they are v24 or earlier.
		// v24 and earlier will not propagate a transaction that does not pass soft-validation.
		// Check if our transaction would pass their soft-validation, using the hardcoded defaults
		// that are used by v24 and earlier.
		if c.UserAgent.Empty() {
			if err := verifyUserTxnAgainstPeer(txn, head, inputs, params.VerifyTxn{
				BurnFactor:          2,
				MaxTransactionSize:  32 * 1024,
				MaxDropletPrecision: 3,
			}); err != nil {
				logger.WithFields(logrus.Fields{
					"addr":   c.Addr,
					"gnetID": c.gnetID,
				}).Debug("Peer will not propagate this transaction")
				continue
			}
		}

		accepts++
	}

	if accepts == 0 {
		return 0, ErrNoPeerAcceptsTxn
	}

	return accepts, nil
}

// verifyUserTxnAgainstPeer returns an error if a user-created transaction would not pass soft-validation
// according to a peer's reported verification parameters
func verifyUserTxnAgainstPeer(txn coin.Transaction, head *coin.SignedBlock, inputs coin.UxArray, verifyParams params.VerifyTxn) error {
	// Check the droplet precision
	for _, o := range txn.Out {
		if err := params.DropletPrecisionCheck(verifyParams.MaxDropletPrecision, o.Coins); err != nil {
			return err
		}
	}

	// Check the txn size
	txnSize, err := txn.Size()
	if err != nil {
		logger.Critical().WithError(err).Error("txn.Size failed unexpectedly")
		return err
	}

	if txnSize > verifyParams.MaxTransactionSize {
		return visor.ErrTxnExceedsMaxBlockSize
	}

	// Check the coinhour burn fee
	f, err := fee.TransactionFee(&txn, head.Time(), inputs)
	if err != nil {
		return err
	}

	if err := fee.VerifyTransactionFee(&txn, f, verifyParams.BurnFactor); err != nil {
		return err
	}

	return nil
}

// Disconnect sends a DisconnectMessage to a peer. After the DisconnectMessage is sent, the peer is disconnected.
// This allows all pending messages to be sent. Any message queued after a DisconnectMessage is unlikely to be sent
// to the peer (but possible).
func (dm *Daemon) Disconnect(addr string, r gnet.DisconnectReason) error {
	logger.WithFields(logrus.Fields{
		"addr":   addr,
		"reason": r,
	}).Debug("Sending DisconnectMessage")
	return dm.sendMessage(addr, NewDisconnectMessage(r))
}

// Implements private daemoner interface methods:

// requestBlocksFromAddr sends a GetBlocksMessage to one connected address
func (dm *Daemon) requestBlocksFromAddr(addr string) error {
	if dm.Config.DisableNetworking {
		return ErrNetworkingDisabled
	}

	headSeq, ok, err := dm.visor.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot request blocks from addr, there is no head block")
	}

	m := NewGetBlocksMessage(headSeq, dm.Config.BlocksResponseCount)
	return dm.sendMessage(addr, m)
}

// broadcastBlock sends a signed block to all connections
func (dm *Daemon) broadcastBlock(sb coin.SignedBlock) error {
	if dm.Config.DisableNetworking {
		return ErrNetworkingDisabled
	}

	m := NewGiveBlocksMessage([]coin.SignedBlock{sb})
	_, err := dm.broadcastMessage(m)
	return err
}

// daemonConfig returns the daemon config
func (dm *Daemon) daemonConfig() DaemonConfig {
	return dm.Config
}

// connectionIntroduced transfers a connection to the "introduced" state in the connections state machine
// and updates other state
func (dm *Daemon) connectionIntroduced(addr string, gnetID uint64, m *IntroductionMessage) (*connection, error) {
	c, err := dm.connections.introduced(addr, gnetID, m)
	if err != nil {
		return nil, err
	}

	listenAddr := c.ListenAddr()

	fields := logrus.Fields{
		"addr":       addr,
		"gnetID":     m.c.ConnID,
		"connGnetID": c.gnetID,
		"listenPort": m.ListenPort,
		"listenAddr": listenAddr,
	}

	if c.Outgoing {
		// For successful outgoing connections, mark the peer as having an incoming port in the pex peerlist
		// The peer should already be in the peerlist, since we use the peerlist to choose an outgoing connection to make
		if err := dm.pex.SetHasIncomingPort(listenAddr, true); err != nil {
			logger.Critical().WithError(err).WithFields(fields).Error("pex.SetHasIncomingPort failed")
			return nil, err
		}
	} else {
		// For successful incoming connections, add the peer to the peer list, with their self-reported listen port
		if err := dm.pex.AddPeer(listenAddr); err != nil {
			logger.Critical().WithError(err).WithFields(fields).Error("pex.AddPeer failed")
			return nil, err
		}
	}

	if err := dm.pex.SetUserAgent(listenAddr, c.UserAgent); err != nil {
		logger.Critical().WithError(err).WithFields(fields).Error("pex.SetUserAgent failed")
		return nil, err
	}

	dm.pex.ResetRetryTimes(listenAddr)

	return c, nil
}

// sendRandomPeers sends a random sample of peers to another peer
func (dm *Daemon) sendRandomPeers(addr string) error {
	peers := dm.pex.RandomExchangeable(dm.pex.Config.ReplyCount)
	if len(peers) == 0 {
		logger.Debug("sendRandomPeers: no peers to send in reply")
		return errors.New("No peers available")
	}

	m := NewGivePeersMessage(peers)
	return dm.sendMessage(addr, m)
}

// announceAllValidTxns broadcasts valid unconfirmed transactions
func (dm *Daemon) announceAllValidTxns() error {
	if dm.Config.DisableNetworking {
		return ErrNetworkingDisabled
	}

	// Get valid unconfirmed transaction hashes
	hashes, err := dm.visor.GetAllValidUnconfirmedTxHashes()
	if err != nil {
		return err
	}

	// Divide hashes into multiple sets of max size
	hashesSet := divideHashes(hashes, dm.Config.MaxTxnAnnounceNum)

	for _, hs := range hashesSet {
		m := NewAnnounceTxnsMessage(hs)
		if _, err := dm.broadcastMessage(m); err != nil {
			logger.WithError(err).Debug("Broadcast AnnounceTxnsMessage failed")
			return err
		}
	}

	return nil
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

// sendMessage sends a Message to a Connection and pushes the result onto the SendResults channel.
func (dm *Daemon) sendMessage(addr string, msg gnet.Message) error {
	return dm.pool.Pool.SendMessage(addr, msg)
}

// broadcastMessage sends a Message to all introduced connections in the Pool.
// Returns the gnet IDs of connections that broadcast succeeded for.
// Note that a connection could still fail to receive the message under certain network conditions,
// there is no guarantee that a message was broadcast.
func (dm *Daemon) broadcastMessage(msg gnet.Message) ([]uint64, error) {
	if dm.Config.DisableNetworking {
		return nil, ErrNetworkingDisabled
	}

	conns := dm.connections.all()
	var addrs []string
	for _, c := range conns {
		if c.HasIntroduced() {
			addrs = append(addrs, c.Addr)
		}
	}

	return dm.pool.Pool.BroadcastMessage(msg, addrs)
}

// disconnectNow disconnects from a peer immediately without sending a DisconnectMessage. Any pending messages
// will not be sent to the peer.
func (dm *Daemon) disconnectNow(addr string, r gnet.DisconnectReason) error {
	return dm.pool.Pool.Disconnect(addr, r)
}

// pexConfig returns the pex config
func (dm *Daemon) pexConfig() pex.Config {
	return dm.pex.Config
}

// addPeers adds peers to the pex
func (dm *Daemon) addPeers(addrs []string) int {
	return dm.pex.AddPeers(addrs)
}

// recordPeerHeight records the height of specific peer
func (dm *Daemon) recordPeerHeight(addr string, gnetID, height uint64) {
	if err := dm.connections.SetHeight(addr, gnetID, height); err != nil {
		logger.Critical().WithError(err).WithField("addr", addr).Error("connections.SetHeight failed")
	}
}

// getSignedBlocksSince returns N signed blocks since given seq
func (dm *Daemon) getSignedBlocksSince(seq, count uint64) ([]coin.SignedBlock, error) {
	return dm.visor.GetSignedBlocksSince(seq, count)
}

// headBkSeq returns the head block sequence
func (dm *Daemon) headBkSeq() (uint64, bool, error) {
	return dm.visor.HeadBkSeq()
}

// executeSignedBlock executes the signed block
func (dm *Daemon) executeSignedBlock(b coin.SignedBlock) error {
	return dm.visor.ExecuteSignedBlock(b)
}

// filterKnownUnconfirmed returns unconfirmed txn hashes with known ones removed
func (dm *Daemon) filterKnownUnconfirmed(txns []cipher.SHA256) ([]cipher.SHA256, error) {
	return dm.visor.FilterKnownUnconfirmed(txns)
}

// getKnownUnconfirmed returns unconfirmed txn hashes with known ones removed
func (dm *Daemon) getKnownUnconfirmed(txns []cipher.SHA256) (coin.Transactions, error) {
	return dm.visor.GetKnownUnconfirmed(txns)
}

// injectTransaction records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain.
// The bool return value is whether or not the transaction was already in the pool.
// If the transaction violates hard constraints, it is rejected, and error will not be nil.
// If the transaction only violates soft constraints, it is still injected, and the soft constraint violation is returned.
func (dm *Daemon) injectTransaction(txn coin.Transaction) (bool, *visor.ErrTxnViolatesSoftConstraint, error) {
	return dm.visor.InjectForeignTransaction(txn)
}
