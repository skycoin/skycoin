package daemon

/* TODO

   Heartbeat message, Cull idle peers, either gnet or this will need to keep
   track of last message time (probably gnet)

    Figure out why messages/connections are stalling
    Setup http://golang.org/pkg/net/http/pprof/ to review goroutine status
*/

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
    "strings"
    "time"
)

// Meta configuration
const (
    // Application version. TODO -- manage version better
    version int32 = 1
)

var (
    logger = logging.MustGetLogger("skycoin.daemon")
)

// Behavioural configuration
const (
    // How often to check and initiate an outgoing connection if needed
    outgoingConnectionsRate = time.Second * 5
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
    // Maximum number of connections to try at once
    pendingConnectionsMax = 16
    // How long to wait for a version packet
    versionWait = time.Second * 30
    // How often to check for peers that have decided to stop communicating
    cullInvalidRate = time.Second * 3
)

var (
    // DisconnectReasons
    DisconnectInvalidVersion gnet.DisconnectReason = errors.New(
        "Invalid version")
    DisconnectVersionTimeout gnet.DisconnectReason = errors.New(
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
    // This is returned when a seemingly impossible error is encountered
    // e.g. net.Conn.Addr() returns an invalid ip:port
    DisconnectOtherError gnet.DisconnectReason = errors.New(
        "Incomprehensible error")

    // Blacklist a peer when they get disconnected for these
    // DisconnectReasons
    BlacklistOffenses = map[gnet.DisconnectReason]time.Duration{
        DisconnectSelf:                      time.Hour * 24,
        DisconnectVersionTimeout:            time.Hour,
        gnet.DisconnectInvalidMessageLength: time.Hour * 8,
        gnet.DisconnectMalformedMessage:     time.Hour * 8,
        gnet.DisconnectUnknownMessage:       time.Hour * 8,
    }
)

// Global state
var (
    // Separate index of outgoing connections. The pool aggregates all
    // connections.
    // TODO -- should this be part of gnet?
    outgoingConnections = make(map[string]*gnet.Connection,
        outgoingConnectionsMax)
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
    // Client connection/disconnection callbacks
    onConnectEvent = make(chan ConnectEvent, 8)
    // Connection failure events
    connectionErrors = make(chan ConnectionError, 8)
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
    go DaemonLoop(quit)
}

// Terminates peer subsytem safely
func Shutdown(dataDir string) {
    ShutdownPool()
    ShutdownPeers(dataDir)
}

// Main loop for peer/connection management. Send anything to quit to shut it
// down
func DaemonLoop(quit chan int) {
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
                logger.Warning("Failed to send ping message to %s", c.Addr())
            }
        }
    }
}

// Attempts to connect to a random peer. If it fails, the peer is removed
func connectToRandomPeer() {
    // Make a connection to a random peer
    peers := Peers.Peerlist.Random(0)
    for _, p := range peers {
        if Pool.Addresses[p.Addr] == nil && pendingConnections[p.Addr] == nil {
            logger.Debug("Trying to connect to %s", p.Addr)
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
    logger.Debug("Removing %s because failed to connect: %v", c.Addr,
        c.Error)
    delete(pendingConnections, c.Addr)
    delete(Peers.Peerlist, c.Addr)
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
            logger.Info("Removing %s for not sending a version", a)
            delete(expectingVersions, a)
            Pool.Disconnect(Pool.Addresses[a], DisconnectVersionTimeout)
            delete(Peers.Peerlist, a)
        }
    }
}

// Called when a ConnectEvent is processed off the onConnectEvent channel
func onConnect(e ConnectEvent) {
    a := e.Addr
    if e.Solicited {
        logger.Info("Connected to %s as we requested", a)
    } else {
        logger.Info("Received unsolicited connection to %s", a)
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
        logger.Info("%s is blacklisted, disconnecting", a)
        Pool.Disconnect(c, DisconnectIsBlacklisted)
        return
    }
    logger.Debug("Sending version message to %s", a)
    outgoingConnections[a] = c
    expectingVersions[a] = time.Now()
    err := c.Controller.SendMessage(NewIntroductionMessage())
    if err != nil {
        logger.Error("Failed to send introduction message: %v", err)
        Pool.Disconnect(c, DisconnectSelf)
        return
    }
}

// Triggered when an gnet.Connection terminates. Disconnect events are not
// pushed to a separate channel, because disconnects are already processed
// by a queue in the DaemonLoop() select{}.
func onDisconnect(c *gnet.Connection, reason gnet.DisconnectReason) {
    a := c.Addr()
    logger.Info("%s disconnected because: %v", a, reason)
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
