package daemon

import (
    "encoding/binary"
    "errors"
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
    mirrorValue = rand.New(rand.NewSource(
        time.Now().UTC().UnixNano())).Uint32()
    // Message handling queue
    messageEvent = make(chan AsyncMessage, gnet.EventChannelBufferSize)
    // Message ID prefices
    IntroductionPrefix = gnet.MessagePrefix{'I', 'N', 'T', 'R'}
    GetPeersPrefix     = gnet.MessagePrefix{'G', 'E', 'T', 'P'}
    GivePeersPrefix    = gnet.MessagePrefix{'G', 'I', 'V', 'P'}
    PingPrefix         = gnet.MessagePrefix{'P', 'I', 'N', 'G'}
    PongPrefix         = gnet.MessagePrefix{'P', 'O', 'N', 'G'}
)

// Registers our Messages with gnet
func RegisterMessages() {
    gnet.RegisterMessage(IntroductionPrefix, IntroductionMessage{})
    gnet.RegisterMessage(GetPeersPrefix, GetPeersMessage{})
    gnet.RegisterMessage(GivePeersPrefix, GivePeersMessage{})
    gnet.RegisterMessage(PingPrefix, PingMessage{})
    gnet.RegisterMessage(PongPrefix, PongMessage{})
    gnet.VerifyMessages()
}

// Compact representation of IP:Port
type IPAddr struct {
    IP   uint32
    Port uint16
}

// Returns an IPAddr from an ip:port string.  If ipv6 or invalid, error is
// returned
func NewIPAddr(addr string) (ipaddr IPAddr, err error) {
    // TODO -- support ipv6
    ipport := strings.Split(addr, ":")
    if len(ipport) != 2 {
        err = errors.New("Invalid ip:port string")
        return
    }
    ipb := net.ParseIP(ipport[0]).To4()
    if ipb == nil {
        err = errors.New("Ignoring IPv6 address")
        return
    }
    ip := binary.BigEndian.Uint32(ipb)
    port, err := strconv.ParseUint(ipport[1], 10, 16)
    if err != nil {
        err = errors.New("Invalid port")
        return
    }
    ipaddr.IP = ip
    ipaddr.Port = uint16(port)
    return
}

// Returns IPAddr as "ip:port"
func (self IPAddr) String() string {
    ipb := make([]byte, 4)
    binary.BigEndian.PutUint32(ipb, self.IP)
    return fmt.Sprintf("%s:%d", net.IP(ipb).String(), self.Port)
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
    c *gnet.MessageContext `enc:"-"`
}

func NewGetPeersMessage() *GetPeersMessage {
    return &GetPeersMessage{}
}

func (self *GetPeersMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    return recordMessageEvent(self, mc)
}

// Notifies the Pex instance that peers were requested
func (self *GetPeersMessage) Process() {
    peers := Peers.Peerlist.Random(peerReplyCount)
    if len(peers) == 0 {
        logger.Debug("We have no peers to send in reply")
        return
    }
    m := NewGivePeersMessage(peers)
    err := self.c.Conn.Controller.SendMessage(m)
    if err != nil {
        logger.Warning("Failed to send GivePeersMessage: %v", err)
    }
}

// Sent in response to GetPeersMessage
type GivePeersMessage struct {
    Peers []IPAddr
    c     *gnet.MessageContext `enc:"-"`
}

// []*pex.Peer is converted to []IPAddr for binary transmission
func NewGivePeersMessage(peers []*pex.Peer) *GivePeersMessage {
    ipaddrs := make([]IPAddr, 0, len(peers))
    for _, ps := range peers {
        ipaddr, err := NewIPAddr(ps.Addr)
        if err != nil {
            logger.Warning("GivePeersMessage skipping address %s", ps.Addr)
            logger.Warning(err.Error())
            continue
        }
        ipaddrs = append(ipaddrs, ipaddr)
    }
    return &GivePeersMessage{Peers: ipaddrs}
}

// GetPeers is required by the pex.GivePeersMessage interface.
// It returns the peers contained in the message as an array of "ip:port"
// strings.
func (self *GivePeersMessage) GetPeers() []string {
    peers := make([]string, 0, len(self.Peers))
    for _, ipaddr := range self.Peers {
        peers = append(peers, ipaddr.String())
    }
    return peers
}

func (self *GivePeersMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    return recordMessageEvent(self, mc)
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
    Peers.AddPeers(peers)
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

    c   *gnet.MessageContext `enc:"-"`
    // We validate the message in Handle() and cache the result for Process()
    valid bool `enc:"-"` // skip it during encoding
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
    if err == nil {
        err = recordMessageEvent(self, mc)
    }
    return
}

// Processes an event queued by Handle()
func (self *IntroductionMessage) Process() {
    delete(expectingIntroductions, self.c.Conn.Addr())
    if !self.valid {
        return
    }
    // Add the remote peer with their chosen listening port
    a := self.c.Conn.Addr()
    ipport := strings.Split(a, ":")
    if len(ipport) != 2 {
        // This should never happen, but the program should still work if it
        // does.
        logger.Error("Invalid Addr() for connection: %s", a)
        Pool.Disconnect(self.c.Conn, DisconnectOtherError)
        return
    }
    ip := ipport[0]
    port, err := strconv.ParseUint(ipport[1], 10, 16)
    if err != nil {
        // This should never happen, but the program should still work if it
        // does.
        logger.Error("Invalid port for connection %s", a)
        Pool.Disconnect(self.c.Conn, DisconnectOtherError)
        return
    }
    Peers.AddPeer(fmt.Sprintf("%s:%d", ip, self.Port))
    // Record their listener, to avoid double connections
    connectionMirrors[a] = self.Mirror
    m := mirrorConnections[self.Mirror]
    if m == nil {
        m = make(map[string]uint16, 1)
    }
    m[ip] = uint16(port)
    mirrorConnections[self.Mirror] = m
}

// Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
    c *gnet.MessageContext `enc:"-"`
}

func (self *PingMessage) Handle(mc *gnet.MessageContext) error {
    self.c = mc
    return recordMessageEvent(self, mc)
}

// Sends a PongMessage to the sender of PingMessage
func (self *PingMessage) Process() {
    logger.Debug("Reply to ping from %s", self.c.Conn.Addr())
    err := self.c.Conn.Controller.SendMessage(&PongMessage{})
    if err != nil {
        logger.Warning("Failed to send PongMessage to %s", self.c.Conn.Addr())
        logger.Warning("Reason: %v", err)
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
