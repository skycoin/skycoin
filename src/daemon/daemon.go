package daemon

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
    "log"
    "strconv"
    "strings"
    "time"
)

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
    DisconnectFailedSend gnet.DisconnectReason = errors.New(
        "Failed to send data to this connection")
    DisconnectNoIntroduction gnet.DisconnectReason = errors.New(
        "First message was not an Introduction")
    DisconnectIPLimitReached gnet.DisconnectReason = errors.New(
        "Maximum number of connections for this IP was reached")
    // This is returned when a seemingly impossible error is encountered
    // e.g. net.Conn.Addr() returns an invalid ip:port
    DisconnectOtherError gnet.DisconnectReason = errors.New(
        "Incomprehensible error")

    // Blacklist a peer when they get disconnected for these
    // DisconnectReasons
    BlacklistOffenses = map[gnet.DisconnectReason]time.Duration{
        DisconnectSelf:                      time.Hour * 24,
        DisconnectIntroductionTimeout:       time.Hour,
        DisconnectNoIntroduction:            time.Hour * 8,
        gnet.DisconnectInvalidMessageLength: time.Hour * 8,
        gnet.DisconnectMalformedMessage:     time.Hour * 8,
        gnet.DisconnectUnknownMessage:       time.Hour * 8,
    }

    logger = logging.MustGetLogger("skycoin.daemon")
)

// Subsystem configurations
type Config struct {
    Daemon   DaemonConfig
    Messages MessagesConfig
    Pool     PoolConfig
    Peers    PeersConfig
    DHT      DHTConfig
    RPC      RPCConfig
}

// Returns a Config with defaults set
func NewConfig() *Config {
    return &Config{
        Daemon:   NewDaemonConfig(),
        Pool:     NewPoolConfig(),
        Peers:    NewPeersConfig(),
        DHT:      NewDHTConfig(),
        RPC:      NewRPCConfig(),
        Messages: NewMessagesConfig(),
    }
}

// Configuration for the Daemon
type DaemonConfig struct {
    // Application version. TODO -- manage version better
    Version int32
    // TCP/UDP port for connections and DHT
    Port int
    // Directory where application data is stored
    DataDirectory string
    // How often to check and initiate an outgoing connection if needed
    OutgoingRate time.Duration
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
}

func NewDaemonConfig() DaemonConfig {
    return DaemonConfig{
        Version:          1,
        Port:             6677,
        OutgoingRate:     time.Second * 5,
        OutgoingMax:      8,
        PendingMax:       16,
        IntroductionWait: time.Second * 30,
        CullInvalidRate:  time.Second * 3,
        IPCountsMax:      3,
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
    DHT      *DHT
    RPC      *RPC

    // Separate index of outgoing connections. The pool aggregates all
    // connections.
    outgoingConnections map[string]*gnet.Connection
    // Number of connections waiting to be formed or timeout
    pendingConnections map[string]*pex.Peer
    // Keep track of unsolicited clients who should notify us of their version
    expectingIntroductions map[string]time.Time
    // Keep track of a connection's mirror value, to avoid double
    // connections (one to their listener, and one to our listener)
    // Maps from addr to mirror value
    connectionMirrors map[string]uint32
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
}

// Returns a Daemon with primitives allocated
func NewDaemon(config *Config) *Daemon {
    d := &Daemon{
        Config:   config.Daemon,
        Messages: NewMessages(config.Messages),
        Pool:     NewPool(config.Pool),
        Peers:    NewPeers(config.Peers),
        DHT:      NewDHT(config.DHT),
        expectingIntroductions: make(map[string]time.Time),
        connectionMirrors:      make(map[string]uint32),
        mirrorConnections:      make(map[uint32]map[string]uint16),
        ipCounts:               make(map[string]int),
        onConnectEvent: make(chan ConnectEvent,
            config.Daemon.OutgoingMax),
        connectionErrors: make(chan ConnectionError,
            config.Daemon.OutgoingMax),
        outgoingConnections: make(map[string]*gnet.Connection,
            config.Daemon.OutgoingMax),
        pendingConnections: make(map[string]*pex.Peer,
            config.Daemon.PendingMax),
    }
    d.RPC = NewRPC(config.RPC, d)
    d.Messages.Config.Register()
    d.Pool.Init(d)
    d.Peers.Init()
    d.DHT.Init()
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

// Terminates all subsystems safely.  To stop the Daemon run loop, send a value
// over the quit channel provided to Init.  The Daemon run lopp must be stopped
// before calling this function.
func (self *Daemon) Shutdown() {
    self.DHT.Shutdown()
    self.Pool.Shutdown()
    self.Peers.Shutdown()
    gnet.EraseMessages()
}

// Main loop for peer/connection management. Send anything to quit to shut it
// down
func (self *Daemon) Start(quit chan int) {
    go self.Pool.Start()
    go self.DHT.Start()

    dhtBootstrapTicker := time.Tick(self.DHT.Config.BootstrapRequestRate)
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
        // Continually make requests to the DHT, if we need peers
        case <-dhtBootstrapTicker:
            if len(self.Peers.Peers.Peerlist) < self.DHT.Config.PeerLimit {
                go self.DHT.RequestPeers()
            }
        // Flush expired blacklisted peers
        case <-updateBlacklistTicker:
            self.Peers.Peers.Blacklist.Refresh()
        // Remove connections that failed to complete the handshake
        case <-cullInvalidTicker:
            self.cullInvalidConnections()
        // Request peers via PEX
        case <-requestPeersTicker:
            if !self.Peers.Peers.Full() {
                self.Pool.requestPeers()
            }
        // Remove peers we haven't seen in a while
        case <-clearOldPeersTicker:
            self.Peers.Peers.Peerlist.ClearOld(self.Peers.Config.Expiration)
        // Remove connections that haven't said anything in a while
        case <-clearStaleConnectionsTicker:
            self.Pool.clearStaleConnections()
        // Sends pings as needed
        case <-idleCheckTicker:
            self.Pool.sendPings()
        // Fill up our outgoing connections
        case <-outgoingConnectionsTicker:
            if len(self.outgoingConnections) < self.Config.OutgoingMax &&
                len(self.pendingConnections) < self.Config.PendingMax {
                self.connectToRandomPeer()
            }
        // Process the connection queue
        case <-messageHandlingTicker:
            self.Pool.Pool.HandleMessages()
        // Process callbacks for when a client connects. No disconnect chan
        // is needed because the callback is triggered by HandleDisconnectEvent
        // which is already select{}ed here
        case r := <-self.onConnectEvent:
            self.onConnect(r)
        // Handle connection errors
        case r := <-self.connectionErrors:
            self.handleConnectionError(r)
        // Update Peers when DHT reports a new one
        case r := <-self.DHT.DHT.PeersRequestResults:
            self.DHT.ReceivePeers(r, self.Peers.Peers)
        case r := <-self.Pool.Pool.DisconnectQueue:
            self.Pool.Pool.HandleDisconnectEvent(r)
        // Message handlers
        case m := <-self.Messages.Events:
            m.Process(self)
        // Process any pending API requests
        case fn := <-self.RPC.requests:
            self.RPC.responses <- fn()
        case <-quit:
            break main
        }
    }
}

// Returns the ListenPort for a given address.  If no port is found, 0 is
// returned
func (self *Daemon) getListenPort(addr string) uint16 {
    m, ok := self.connectionMirrors[addr]
    if !ok {
        return 0
    }
    mc := self.mirrorConnections[m]
    if mc == nil {
        log.Panic("mirrorConnections map does not exist, but mirror does")
    }
    return mc[strings.Split(addr, ":")[0]]
}

// Attempts to connect to a random peer. If it fails, the peer is removed
func (self *Daemon) connectToRandomPeer() {
    // Make a connection to a random peer
    peers := self.Peers.Peers.Peerlist.Random(0)
    for _, p := range peers {
        isConnected := self.Pool.Pool.Addresses[p.Addr] != nil
        isPending := self.pendingConnections[p.Addr] != nil
        ipInUse := self.ipCounts[strings.Split(p.Addr, ":")[0]] != 0
        if !isConnected && !isPending && !ipInUse {
            logger.Debug("Trying to connect to %s", p.Addr)
            self.pendingConnections[p.Addr] = p
            go func() {
                _, err := self.Pool.Pool.Connect(p.Addr)
                if err != nil {
                    self.connectionErrors <- ConnectionError{p.Addr, err}
                }
            }()
            break
        }
    }
}

// We remove a peer from the Pex if we failed to connect
func (self *Daemon) handleConnectionError(c ConnectionError) {
    logger.Debug("Removing %s because failed to connect: %v", c.Addr,
        c.Error)
    delete(self.pendingConnections, c.Addr)
    delete(self.Peers.Peers.Peerlist, c.Addr)
}

// Removes unsolicited connections who haven't sent a version
func (self *Daemon) cullInvalidConnections() {
    // This method only handles the erroneous people from the DHT, but not
    // malicious nodes
    now := time.Now().UTC()
    for a, t := range self.expectingIntroductions {
        // Forget about anyone that already disconnected
        if self.Pool.Pool.Addresses[a] == nil {
            delete(self.expectingIntroductions, a)
            continue
        }
        // Remove anyone that fails to send a version within introductionWait time
        if t.Add(self.Config.IntroductionWait).Before(now) {
            logger.Info("Removing %s for not sending a version", a)
            delete(self.expectingIntroductions, a)
            self.Pool.Pool.Disconnect(self.Pool.Pool.Addresses[a],
                DisconnectIntroductionTimeout)
            delete(self.Peers.Peers.Peerlist, a)
        }
    }
}

// Records an AsyncMessage to the messageEvent chan.  Do not access
// messageEvent directly.
func (self *Daemon) recordMessageEvent(m AsyncMessage,
    c *gnet.MessageContext) error {
    // The first message received must be an Introduction
    _, needsIntro := self.expectingIntroductions[c.Conn.Addr()]
    if needsIntro {
        _, isIntro := m.(*IntroductionMessage)
        if !isIntro {
            self.Pool.Pool.Disconnect(c.Conn, DisconnectNoIntroduction)
            return DisconnectNoIntroduction
        }
    }
    self.Messages.Events <- m
    return nil
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
        self.outgoingConnections[a] = c
    }
    self.expectingIntroductions[a] = time.Now().UTC()
    logger.Debug("Sending introduction message to %s", a)
    m := NewIntroductionMessage(self.Messages.Mirror, self.Config.Version,
        self.Pool.Pool.ListenPort)
    err := self.Pool.Pool.Dispatcher.SendMessage(c, m)
    if err != nil {
        logger.Error("Failed to send introduction message: %v", err)
        self.Pool.Pool.Disconnect(c, DisconnectFailedSend)
        return
    }
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
    delete(self.outgoingConnections, a)
    delete(self.expectingIntroductions, a)
    self.removeIPCount(a)
    self.removeConnectionMirror(a)
}

// Triggered when an gnet.Connection is connected
func (self *Daemon) onGnetConnect(c *gnet.Connection, solicited bool) {
    self.onConnectEvent <- ConnectEvent{Addr: c.Addr(), Solicited: solicited}
}

// Returns whether the ipCount maximum has been reached
func (self *Daemon) ipCountMaxed(addr string) bool {
    ip := strings.Split(addr, ":")[0]
    return self.ipCounts[ip] >= self.Config.IPCountsMax
}

// Adds base IP to ipCount or returns error if max is reached
func (self *Daemon) recordIPCount(addr string) {
    ip := strings.Split(addr, ":")[0]
    _, hasCount := self.ipCounts[ip]
    if !hasCount {
        self.ipCounts[ip] = 0
    }
    self.ipCounts[ip] += 1
}

// Removes base IP from ipCount
func (self *Daemon) removeIPCount(addr string) {
    ip := strings.Split(addr, ":")[0]
    if self.ipCounts[ip] <= 1 {
        delete(self.ipCounts, ip)
    } else {
        self.ipCounts[ip] -= 1
    }
}

// Adds addr + mirror to the connectionMirror mappings
func (self *Daemon) recordConnectionMirror(addr string, mirror uint32) error {
    ipPort := strings.Split(addr, ":")
    ip := ipPort[0]
    sport := ipPort[1]
    port64, err := strconv.ParseUint(sport, 10, 16)
    if err != nil {
        return err
    }
    port := uint16(port64)
    self.connectionMirrors[addr] = mirror
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
    mirror, ok := self.connectionMirrors[addr]
    if !ok {
        return
    }
    ip := strings.Split(addr, ":")[0]
    m := self.mirrorConnections[mirror]
    if len(m) <= 1 {
        delete(self.mirrorConnections, mirror)
    } else {
        delete(m, ip)
    }
    delete(self.connectionMirrors, addr)
}

// Returns whether an addr+mirror's port and whether the port exists
func (self *Daemon) getMirrorPort(addr string, mirror uint32) (uint16, bool) {
    ips := self.mirrorConnections[mirror]
    if ips == nil {
        return 0, false
    }
    ip := strings.Split(addr, ":")[0]
    port, exists := ips[ip]
    return port, exists
}
