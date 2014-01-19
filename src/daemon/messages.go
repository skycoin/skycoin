package daemon

import (
    "encoding/binary"
    "fmt"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
    "math/rand"
    "net"
    "strconv"
    "strings"
    "time"
)

// Message represent a packet to be serialized over the network by
// the gnet encoder.
// They must implement the gnet.Message interface
// All concurrent daemon write operations are synchronized by the daemon's
// DaemonLoop().
// Message do this by caching the gnet.MessageContext received in Handle()
// and placing itself on the messageEvent channel.
// When the message is retrieved from the messageEvent channel, its Process()
// method is called.

var (
    // Magic value for detecting self-connection
    mirrorValue = rand.New(rand.NewSource(time.Now().UTC().UnixNano())).Uint32()
    // Message handling queue
    messageEvent = make(chan AsyncMessage, gnet.EventChannelBufferSize)
)

// Registers our Messages with gnet
func RegisterMessages() {
    gnet.RegisterMessage(IntroductionMessage{})
    gnet.RegisterMessage(GetPeersMessage{})
    gnet.RegisterMessage(GivePeersMessage{})
    gnet.RegisterMessage(PingMessage{})
    gnet.RegisterMessage(PongMessage{})
    gnet.VerifyMessages()
}

// Compact representation of IP:Port
type IPAddr struct {
    IP   uint32
    Port uint16
}

// Messages that perform an action when received must implement this interface.
// Process() is called after the message is pulled off of messageEvent channel.
// Messages should place themselves on the messageEvent channel in their
// Handle() method required by gnet.
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

// Notifies the Pex instance that peers were requested
func (self *GetPeersMessage) Process() {
    Peers.RespondToGetPeersMessage(self.c.Conn.Conn, NewGivePeersMessage)
}

// Send is required by the pex.GetPeersMessage interface
func (self *GetPeersMessage) Send(c net.Conn) error {
    return gnet.WriteMessage(c, self)
}

// Sent in response to GetPeersMessage
type GivePeersMessage struct {
    Peers []IPAddr
    c     *gnet.MessageContext `-`
}

// []*pex.Peer is converted to []IPAddr for binary transmission
func NewGivePeersMessage(peers []*pex.Peer) pex.GivePeersMessage {
    ipaddrs := make([]IPAddr, 0, len(peers))
    for _, ps := range peers {
        // TODO -- support ipv6
        ipport := strings.Split(ps.Addr, ":")
        ipb := net.ParseIP(ipport[0]).To4()
        if ipb == nil {
            logger.Warning("Ignoring IPv6 address %s", ipport[0])
            continue
        }
        ip := binary.BigEndian.Uint32(ipb)
        port, err := strconv.ParseUint(ipport[1], 10, 16)
        if err != nil {
            logger.Error("Invalid port in peer address %s", ps.Addr)
            continue
        }
        ipaddrs = append(ipaddrs, IPAddr{IP: ip, Port: uint16(port)})
    }
    return &GivePeersMessage{Peers: ipaddrs}
}

// GetPeers is required by the pex.GivePeersMessage interface.
// It returns the peers contained in the message as an array of "ip:port"
// strings.
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

// Send is required by the pex.GivePeersMessage interface
func (self *GivePeersMessage) Send(c net.Conn) error {
    return gnet.WriteMessage(c, self)
}

func (self *GivePeersMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    messageEvent <- self
    return nil
}

// Notifies the Pex instance that peers were received
func (self *GivePeersMessage) Process() {
    peers := self.GetPeers()
    if len(peers) != 0 {
        logger.Debug("Got these peers via PEX:")
        for _, p := range peers {
            logger.Debug("\t%s", p)
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
        logger.Info("Remote mirror value %v matches ours", self.Mirror)
        Pool.Disconnect(mc.Conn, DisconnectSelf)
        err = DisconnectSelf
    }
    // Disconnect if not running the same version
    if self.Version != version {
        logger.Info("%s has different version %d. Disconnecting.",
            addr, self.Version)
        Pool.Disconnect(mc.Conn, DisconnectInvalidVersion)
        err = DisconnectInvalidVersion
    } else {
        logger.Info("%s verified for version %d", addr, version)
    }
    // Disconnect if connected twice to the same peer (judging by ip:mirror)
    ips := mirrorConnections[self.Mirror]
    if ips != nil {
        ip := strings.Split(addr, ":")[0]
        if port, exists := ips[ip]; exists {
            logger.Info("%s is already connected as %s", addr,
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
        logger.Error("Invalid port for connection %s", a)
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

// Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
    c *gnet.MessageContext `-`
}

func (self *PingMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    messageEvent <- self
    return nil
}

// Sends a PongMessage to the sender of PingMessage
func (self *PingMessage) Process() {
    logger.Debug("Reply to ping from %s", self.c.Conn.Addr())
    if self.c.Conn.Controller.SendMessage(&PongMessage{}) != nil {
        logger.Warning("Failed to send PongMessage to %s", self.c.Conn.Addr())
    }
}

// Sent in reply to a PingMessage.  No action is taken when this is received.
type PongMessage struct {
}

func (self *PongMessage) Handle(mc *gnet.MessageContext) error {
    // There is nothing to do; gnet updates Connection.LastMessage internally
    // when this is received
    logger.Debug("Received pong from %s", mc.Conn.Addr())
    return nil
}
