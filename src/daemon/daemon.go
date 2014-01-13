package daemon

/* TODO

   Heartbeat message, Cull idle peers, either gnet or this will need to keep
   track of last message time (probably gnet)

    Figure out why messages/connections are stalling
    Setup http://golang.org/pkg/net/http/pprof/ to review goroutine status
*/

import (
    "crypto/sha1"
    "encoding/binary"
    "encoding/hex"
    "errors"
    "fmt"
    "log"
    "math/rand"
    "net"
    "strconv"
    "strings"
    "time"

    "github.com/nictuku/dht"
    "github.com/op/go-logging"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
)

var (
    // Logger
    logger = logging.MustGetLogger("skycoin.daemon")
    // Application version
    version int32 = 1
    // Maximum number of peers to keep account of in the PeerList
    maxPeers = 1000
    // Peer list
    Peers = pex.NewPex(maxPeers)
    // Cull peers after they havent been seen in this much time
    peerExpiration = time.Hour * 24 * 7
    // Cull expired peers on this interval
    cullPeerRate = time.Minute * 10
    // How often to check and initiate an outgoing connection if needed
    outgoingConnectionsRate = time.Second * 5
    // How often to clear expired blacklist entries
    updateBlacklistRate = time.Minute
    // How often to request peers via PEX
    requestPeersRate = time.Minute
    // Connection pool
    Pool *gnet.ConnectionPool = nil
    // Timeout when trying to connect to new peers through the pool
    poolConnectTimeout = time.Second * 30
    // How often to check for stale connections
    clearStaleConnectionsRate = time.Minute
    // How long a connection can idle before considered stale
    idleConnectionLimit = time.Minute * 90
    // How often to check for needed pings
    pingCheckRate = time.Minute
    // How long to wait before sending another ping
    pingRate = idleConnectionLimit / 3
    // How often to process message buffers and generate events
    messageHandlingRate = time.Millisecond * 30
    // Number of outgoing connections to maintain
    outgoingConnectionsMax = 8
    // Separate index of outgoing connections. The pool aggregates all connections
    // TODO -- should this be part of gnet?
    outgoingConnections = make(map[string]*gnet.Connection,
        outgoingConnectionsMax)
    // Maximum number of connections to try at once
    pendingConnectionsMax = 16
    // Number of connections waiting to be formed or timeout
    pendingConnections = make(map[string]*pex.Peer, pendingConnectionsMax)
    // Keep track of unsolicited clients who should notify us of their version
    expectingVersions = make(map[string]time.Time)
    // Keep track of a connection's mirror value, to avoid double
    // connections (one to their listener, and one to our listener)
    // Maps from addr to mirror value
    connectionMirrors = make(map[string]uint32)
    // Maps from mirror value to a map of ip (no port)
    // We use a map of ip as value because multiple peers can have the same
    // mirror (to avoid attacks enabled by our use of mirrors),
    // but only one per base ip
    mirrorConnections = make(map[uint32]map[string]uint16)
    // How long to wait for a version packet
    versionWait = time.Second * 30
    // How often to check for peers that have decided to stop communicating
    cullInvalidRate = time.Second * 3
    // DHT manager
    DHT *dht.DHT = nil
    // Info to be hashed for identifying peers on the skycoin network
    dhtInfo = "skycoin-skycoin-skycoin-skycoin-skycoin-skycoin-skycoin"
    // Hex encoded sha1 sum of dhtInfo
    dhtInfoHash dht.InfoHash = ""
    // Number of peers to try to get via DHT
    dhtDesiredPeers = 20
    // How many local peers, from any source, before we stop requesting DHT peers
    dhtPeerLimit = 100
    // DHT Bootstrap routers
    dhtBootstrapNodes = []string{
        "1.a.magnets.im:6881",
        "router.bittorrent.com:6881",
        "router.utorrent.com:6881",
        "dht.transmissionbt.com:6881",
    }
    // How often to request more peers via DHT
    dhtBootstrapRequestRate = time.Second * 10
    // Magic value for detecting self-connection
    mirrorValue = rand.New(rand.NewSource(time.Now().UTC().UnixNano())).Uint32()
    // DisconnectReasons
    DisconnectInvalidVersion    gnet.DisconnectReason = errors.New("Invalid version")
    DisconnectVersionTimeout    gnet.DisconnectReason = errors.New("Version timeout")
    DisconnectVersionSendFailed gnet.DisconnectReason = errors.New("Version send failed")
    DisconnectIsBlacklisted     gnet.DisconnectReason = errors.New("Blacklisted")
    DisconnectSelf              gnet.DisconnectReason = errors.New("Self connect")
    DisconnectConnectedTwice    gnet.DisconnectReason = errors.New("Already connected")
    DisconnectIdle              gnet.DisconnectReason = errors.New("Idle")
    // Blacklist a peer when they get disconnected for these gnet.DisconnectReason
    BlacklistOffenses = map[gnet.DisconnectReason]time.Duration{
        DisconnectSelf:                      time.Hour * 24,
        DisconnectVersionTimeout:            time.Hour,
        gnet.DisconnectInvalidMessageLength: time.Hour * 8,
        gnet.DisconnectMalformedMessage:     time.Hour * 8,
        gnet.DisconnectUnknownMessage:       time.Hour * 8,
    }
    // Client connection/disconnection callbacks
    onConnectEvent = make(chan ConnectEvent, 8)
    // Connection failure events
    connectionErrors = make(chan ConnectionError, 8)
    // Message handling queue
    messageEvent = make(chan AsyncMessage, gnet.EventChannelBufferSize)
)

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

// Initializes the daemon subsystem.  Data is sent over both TCP and UDP for
// port.  dataDir is where application data is stored. Sending anything to
// the quit channel will stop the daemon.
func Init(port int, dataDir string, quit chan int) {
    RegisterMessages()
    InitDHT(port)
    InitPool(port)
    InitPeers(dataDir)
    BeginPeerAcquisition(Pool)
    go PeersLoop(quit)
}

// Terminates peer subsytem safely
func Shutdown(dataDir string) {
    ShutdownPool()
    ShutdownPeers(dataDir)
}

// Registers our Messages with gnet
func RegisterMessages() {
    gnet.RegisterMessage(IntroductionMessage{})
    gnet.RegisterMessage(GetPeersMessage{})
    gnet.RegisterMessage(GivePeersMessage{})
    gnet.RegisterMessage(PingMessage{})
    gnet.RegisterMessage(PongMessage{})
    gnet.VerifyMessages()
}

// Sets up the DHT node for peer bootstrapping
func InitDHT(port int) {
    var err error
    sum := sha1.Sum([]byte(dhtInfo))
    // Create a hex encoded sha1 sum of a string to be used for DH
    dhtInfoHash, err = dht.DecodeInfoHash(hex.EncodeToString(sum[:]))
    if err != nil {
        log.Panic("Failed to create InfoHash: %v", err)
        return
    }
    DHT, err = dht.NewDHTNode(port, dhtDesiredPeers, true)
    if err != nil {
        log.Panicf("Failed to init DHT: %v", err)
        return
    }
    logger.Info("Init DHT on port %d\n", port)
    go DHT.DoDHT()
}

// Begins listening on port for connections and periodically scanning for
// messages on read_interval
func InitPool(port int) {
    logger.Info("InitPool on port %d\n", port)
    if Pool != nil {
        log.Panic("ConnectionPool is already initialised")
    }
    gnet.ConnectionTimeout = poolConnectTimeout
    Pool = gnet.NewConnectionPool(port)
    Pool.DisconnectCallback = onGnetDisconnect
    Pool.ConnectCallback = onGnetConnect
    go func() {
        err := Pool.StartListen()
        if err != nil {
            log.Panic(err)
        }
    }()
}

// Closes all connections and stops listening
func ShutdownPool() {
    if Pool != nil {
        Pool.StopListen()
    }
    logger.Info("Shutdown pool\n")
}

// Configure the pex.PeerList and load local data
func InitPeers(data_directory string) {
    err := Peers.Load(data_directory)
    if err != nil {
        logger.Notice("Failed to load peer database\n")
        logger.Notice("Reason: %v\n", err)
    }
    logger.Debug("Init peers\n")
}

// Hits bootstrap nodes if cached peers are not found.  Begins making
// connections to some peers and requesting more
func BeginPeerAcquisition(pool *gnet.ConnectionPool) {
    logger.Info("BeginPeerAcquisition\n")
    logger.Debug("Known peers:\n")
    for addr, _ := range Peers.Peerlist {
        logger.Debug("\t%s\n", addr)
    }
    if len(Peers.Peerlist) == 0 {
        logger.Debug("\tNone\n")
    }
}

// Main loop for peer/connection management. Send anything to quit to shut it
// down
func PeersLoop(quit chan int) {
    dhtBootstrapTicker := time.Tick(dhtBootstrapRequestRate)
    cullInvalidTicker := time.Tick(cullInvalidRate)
    outgoingConnectionsTicker := time.Tick(outgoingConnectionsRate)
    clearOldPeersTicker := time.Tick(cullPeerRate)
    requestPeersTicker := time.Tick(requestPeersRate)
    updateBlacklistTicker := time.Tick(updateBlacklistRate)
    messageHandlingTicker := time.Tick(messageHandlingRate)
    clearStaleConnectionsTicker := time.Tick(clearStaleConnectionsRate)
    pingCheckTicker := time.Tick(pingCheckRate)
main:
    for {
        select {
        // Continually make requests to the DHT, if we need peers
        case <-dhtBootstrapTicker:
            if len(Peers.Peerlist) < dhtPeerLimit {
                go DHT.PeersRequest(string(dhtInfoHash), true)
            }
        // Flush expired blacklisted peers
        case <-updateBlacklistTicker:
            Peers.Blacklist.Refresh()
        // Remove connections that failed to complete the handshake
        case <-cullInvalidTicker:
            cullInvalidConnections()
        // Request peers via PEX
        case <-requestPeersTicker:
            Peers.RequestPeers(Pool.GetRawConnections(), NewGetPeersMessage)
        // Remove peers we haven't seen in a while
        case <-clearOldPeersTicker:
            Peers.Peerlist.ClearOld(peerExpiration)
        // Remove connections that haven't said anything in a while
        case <-clearStaleConnectionsTicker:
            clearStaleConnections()
        // Sends pings as needed
        case <-pingCheckTicker:
            sendPings()
        // Fill up our outgoing connections
        case <-outgoingConnectionsTicker:
            if len(outgoingConnections) < outgoingConnectionsMax &&
                len(pendingConnections) < pendingConnectionsMax {
                connectToRandomPeer()
            }
        // Process the connection queue
        case <-messageHandlingTicker:
            Pool.HandleMessages()
        // Process callbacks for when a client connects. No disconnect chan
        // is needed because the callback is triggered by HandleDisconnectEvent
        // which is already select{}ed here
        case r := <-onConnectEvent:
            onConnect(r)
        // Handle connection errors
        case r := <-connectionErrors:
            handleConnectionError(r)
        // Update Peers when DHT reports a new one
        case r := <-DHT.PeersRequestResults:
            receivedDHTPeer(r)
        case r := <-Pool.DisconnectQueue:
            Pool.HandleDisconnectEvent(r)
        // Message handlers
        case m := <-messageEvent:
            m.Process()
        case <-quit:
            break main
        }
    }
}

// Triggered when an gnet.Connection terminates. Disconnect events are not
// pushed to a separate channel, because disconnects are already processed
// by a queue in the PeersLoop() select{}.
func onGnetDisconnect(c *gnet.Connection, reason gnet.DisconnectReason) {
    a := c.Addr()
    logger.Info("%s disconnected because: %v\n", a, reason)
    duration, exists := BlacklistOffenses[reason]
    if exists {
        Peers.AddBlacklistEntry(a, duration)
    }
    delete(outgoingConnections, a)
    delete(expectingVersions, a)
    // Remove peer from the bidirectional mirror map
    ip := strings.Split(a, ":")[0]
    mirror := connectionMirrors[a]
    m := mirrorConnections[mirror]
    if len(m) <= 1 {
        delete(mirrorConnections, mirror)
    } else {
        delete(m, ip)
    }
    delete(connectionMirrors, a)
}

// Triggered when an gnet.Connection is connected
func onGnetConnect(c *gnet.Connection, solicited bool) {
    onConnectEvent <- ConnectEvent{Addr: c.Addr(), Solicited: solicited}
}

// Called when a ConnectEvent is processed off the onConnectEvent channel
func onConnect(e ConnectEvent) {
    a := e.Addr
    if e.Solicited {
        logger.Info("Connected to %s as we requested\n", a)
    } else {
        logger.Info("Received unsolicited connection to %s\n", a)
    }
    delete(pendingConnections, a)
    c := Pool.Addresses[a]
    if c == nil {
        logger.Warning("While processing an onConnect event, no pool " +
            "connection was found")
        return
    }
    blacklisted := Peers.IsBlacklisted(a)
    if blacklisted {
        logger.Info("%s is blacklisted, disconnecting\n", a)
        Pool.Disconnect(c, DisconnectIsBlacklisted)
        return
    }
    logger.Debug("Sending version message to %s\n", a)
    outgoingConnections[a] = c
    expectingVersions[a] = time.Now()
    err := c.Controller.SendMessage(NewIntroductionMessage())
    if err != nil {
        logger.Error("Failed to send introduction message: %v\n", err)
        Pool.Disconnect(c, DisconnectSelf)
        return
    }
}

// Removes connections that have not sent a message in too long
func clearStaleConnections() {
    now := time.Now()
    for _, c := range Pool.Pool {
        if c.LastReceived.Add(idleConnectionLimit).Before(now) {
            Pool.Disconnect(c, DisconnectIdle)
        }
    }
}

// Send a ping if our last message sent was over pingRate ago
func sendPings() {
    now := time.Now()
    for _, c := range Pool.Pool {
        if c.LastSent.Add(pingRate).Before(now) {
            err := c.Controller.SendMessage(&PingMessage{})
            if err != nil {
                logger.Warning("Failed to send ping message to %s\n", c.Addr())
            }
        }
    }
}

// Removes unsolicited connections who haven't sent a version
func cullInvalidConnections() {
    // This method only handles the erroneous people from the DHT, but not
    // malicious nodes
    now := time.Now()
    for a, t := range expectingVersions {
        if Pool.Addresses[a] == nil {
            delete(expectingVersions, a)
            continue
        }
        if t.Add(versionWait).Before(now) {
            logger.Info("Removing %s for not sending a version\n", a)
            delete(expectingVersions, a)
            Pool.Disconnect(Pool.Addresses[a], DisconnectVersionTimeout)
            delete(Peers.Peerlist, a)
        }
    }
}

// Shutdown the PeerList
func ShutdownPeers(data_directory string) {
    err := Peers.Save(data_directory)
    if err != nil {
        logger.Warning("Failed to save peer database\n")
        logger.Warning("Reason: %v\n", err)
    }
    logger.Debug("Shutdown peers\n")
}

// Called when the DHT finds a peer
func receivedDHTPeer(r map[dht.InfoHash][]string) {
    for _, peers := range r {
        for _, p := range peers {
            // Check that p is valid; for some reason the dht library
            // sometimes gives us a peer that panics when its decoded
            if len(p) < 5 {
                continue
            }
            peer := dht.DecodePeerAddress(p)
            logger.Debug("DHT Peer: %s\n", peer)
            _, err := Peers.AddPeer(peer)
            if err != nil {
                logger.Info("Failed to add DHT peer: %v\n", err)
            }
        }
    }
}

func RequestDHTPeers() {
    ih := string(dhtInfoHash)
    if ih == "" {
        log.Panic("dhtInfoHash is not initialized")
        return
    }
    logger.Info("Requesting DHT Peers\n")
    DHT.PeersRequest(ih, true)
}

// Attempts to connect to a random peer. If it fails, the peer is removed
func connectToRandomPeer() {
    // Make a connection to a random peer
    peers := Peers.Peerlist.Random(0)
    for _, p := range peers {
        if Pool.Addresses[p.Addr] == nil && pendingConnections[p.Addr] == nil {
            logger.Debug("Trying to connect to %s\n", p.Addr)
            pendingConnections[p.Addr] = p
            go func() {
                _, err := Pool.Connect(p.Addr)
                if err != nil {
                    connectionErrors <- ConnectionError{p.Addr, err}
                }
            }()
            break
        }
    }
}

// Called when connecting/dialing an address fails
func handleConnectionError(c ConnectionError) {
    if c.Error == nil {
        return
    }
    // Remove a peer if we fail to connect to it
    logger.Debug("Removing %s because failed to connect: %v\n", c.Addr,
        c.Error)
    delete(pendingConnections, c.Addr)
    delete(Peers.Peerlist, c.Addr)
}

/* Messages */

type AsyncMessage interface {
    Process()
}

// Sent to request peers
type GetPeersMessage struct {
    c *gnet.MessageContext `-`
}

func NewGetPeersMessage() pex.GetPeersMessage {
    return &GetPeersMessage{}
}

func (self *GetPeersMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    messageEvent <- self
    return nil
}

func (self *GetPeersMessage) Process() {
    Peers.RespondToGetPeersMessage(self.c.Conn.Conn, NewGivePeersMessage)
}

func (self *GetPeersMessage) Send(c net.Conn) error {
    return gnet.WriteMessage(c, self)
}

// Contains a list of newline delimited peers
// TODO -- use an addr struct and send as integers

type IPAddr struct {
    IP   uint32
    Port uint16
}

type GivePeersMessage struct {
    Peers []IPAddr
    c     *gnet.MessageContext `-`
}

func NewGivePeersMessage(peers []*pex.Peer) pex.GivePeersMessage {
    ipaddrs := make([]IPAddr, 0, len(peers))
    for _, ps := range peers {
        // TODO -- support ipv6
        ipport := strings.Split(ps.Addr, ":")
        ipb := net.ParseIP(ipport[0]).To4()
        if ipb == nil {
            logger.Warning("Ignoring IPv6 address %s\n", ipport[0])
            continue
        }
        ip := binary.BigEndian.Uint32(ipb)
        port, err := strconv.ParseUint(ipport[1], 10, 16)
        if err != nil {
            logger.Error("Invalid port in peer address %s\n", ps.Addr)
            continue
        }
        ipaddrs = append(ipaddrs, IPAddr{IP: ip, Port: uint16(port)})
    }
    return &GivePeersMessage{Peers: ipaddrs}
}

func (self *GivePeersMessage) GetPeers() []string {
    peers := make([]string, 0, len(self.Peers))
    ipb := make([]byte, 4)
    for _, ipaddr := range self.Peers {
        binary.BigEndian.PutUint32(ipb, ipaddr.IP)
        peer := fmt.Sprintf("%s:%d", net.IP(ipb).String(), ipaddr.Port)
        peers = append(peers, peer)
    }
    return peers
}

func (self *GivePeersMessage) Send(c net.Conn) error {
    return gnet.WriteMessage(c, self)
}

func (self *GivePeersMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    messageEvent <- self
    return nil
}

func (self *GivePeersMessage) Process() {
    peers := self.GetPeers()
    if len(peers) != 0 {
        logger.Debug("Got these peers via PEX:\n")
        for _, p := range peers {
            logger.Debug("\t%s\n", p)
        }
    }
    Peers.RespondToGivePeersMessage(self)
}

// An IntroductionMessage is sent on first connect by both parties
type IntroductionMessage struct {
    // Mirror is a random value generated on client startup that is used
    // to identify self-connections
    Mirror uint32
    // Port is the port that this client is listening on
    Port uint16
    // Our client version
    Version int32

    c   *gnet.MessageContext `-`
    // We validate the message in Handle() and cache the result for Process()
    valid bool `-` // skip it during encoding
}

func NewIntroductionMessage() *IntroductionMessage {
    return &IntroductionMessage{
        Mirror:  mirrorValue,
        Version: version,
        Port:    Pool.ListenPort,
    }
}

// Responds to an gnet.Pool event. We implement Handle() here because we
// need to control the DisconnectReason sent back to gnet.  We still implement
// Process(), where we do modifications that are not threadsafe
func (self *IntroductionMessage) Handle(mc *gnet.MessageContext) (err error) {
    addr := mc.Conn.Addr()
    // Disconnect if this is a self connection (we have the same mirror value)
    if self.Mirror == mirrorValue {
        logger.Info("Remote mirror value %v matches ours\n", self.Mirror)
        Pool.Disconnect(mc.Conn, DisconnectSelf)
        err = DisconnectSelf
    }
    // Disconnect if not running the same version
    if self.Version != version {
        logger.Info("%s has different version %d. Disconnecting.\n",
            addr, self.Version)
        Pool.Disconnect(mc.Conn, DisconnectInvalidVersion)
        err = DisconnectInvalidVersion
    } else {
        logger.Info("%s verified for version %d\n", addr, version)
    }
    // Disconnect if connected twice to the same peer (judging by ip:mirror)
    ips := mirrorConnections[self.Mirror]
    if ips != nil {
        ip := strings.Split(addr, ":")[0]
        if port, exists := ips[ip]; exists {
            logger.Info("%s is already connected as %s\n", addr,
                fmt.Sprintf("%s:%d", ip, port))
            Pool.Disconnect(mc.Conn, DisconnectConnectedTwice)
            err = DisconnectConnectedTwice
        }
    }
    self.valid = (err == nil)
    self.c = mc
    messageEvent <- self
    return
}

// Processes an event queued by Handle()
func (self *IntroductionMessage) Process() {
    delete(expectingVersions, self.c.Conn.Addr())
    if !self.valid {
        return
    }
    // Add the remote peer with their chosen listening port
    a := self.c.Conn.Addr()
    ipport := strings.Split(a, ":")
    ip := ipport[0]
    port, err := strconv.ParseUint(ipport[1], 10, 16)
    if err != nil {
        // This should never happen, but the program should still work if it
        // does.
        logger.Error("Invalid port for connection %s\n", a)
    }
    Peers.AddPeer(fmt.Sprintf("%s:%d", ip, self.Port))
    // Record their listener, to avoid double connections
    connectionMirrors[ip] = self.Mirror
    m := mirrorConnections[self.Mirror]
    if m == nil {
        m = make(map[string]uint16, 1)
    }
    m[ip] = uint16(port)
    mirrorConnections[self.Mirror] = m
}

// Sent to keep a connection alive
type PingMessage struct {
    c *gnet.MessageContext `-`
}

func (self *PingMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    messageEvent <- self
    return nil
}

func (self *PingMessage) Process() {
    logger.Debug("Reply to ping from %s\n", self.c.Conn.Addr())
    if self.c.Conn.Controller.SendMessage(&PongMessage{}) != nil {
        logger.Warning("Failed to send PongMessage to %s\n", self.c.Conn.Addr())
    }
}

// Sent in reply to a PingMessage
type PongMessage struct {
}

func (self *PongMessage) Handle(mc *gnet.MessageContext) error {
    // There is nothing to do; gnet updates Connection.LastMessage internally
    // when this is received
    logger.Debug("Received pong from %s\n", mc.Conn.Addr())
    return nil
}

// DHT Event Logger
type DHTLogger struct{}

// Logs a GetPeers event
func (self *DHTLogger) GetPeers(ip *net.UDPAddr, id string,
    _info dht.InfoHash) {
    id = hex.EncodeToString([]byte(id))
    info := hex.EncodeToString([]byte(_info))
    logger.Debug("DHT GetPeers event occured:\n\tid: %s\n\tinfohash: %s\n",
        id, info)
}
