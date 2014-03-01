package daemon

import (
    "encoding/binary"
    "errors"
    "fmt"
    "github.com/skycoin/gnet"
    "github.com/skycoin/pex"
    "github.com/skycoin/skycoin/src/util"
    "math/rand"
    "net"
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

// Message config contains a gnet.Message's 4byte prefix and a
// reference interface
type MessageConfig struct {
    Prefix  gnet.MessagePrefix
    Message interface{}
}

func NewMessageConfig(prefix string, m interface{}) MessageConfig {
    return MessageConfig{
        Message: m,
        Prefix:  gnet.MessagePrefixFromString(prefix),
    }
}

// Creates and populates the message configs
func getMessageConfigs() []MessageConfig {
    return []MessageConfig{
        NewMessageConfig("INTR", IntroductionMessage{}),
        NewMessageConfig("GETP", GetPeersMessage{}),
        NewMessageConfig("GIVP", GivePeersMessage{}),
        NewMessageConfig("PING", PingMessage{}),
        NewMessageConfig("PONG", PongMessage{}),
        NewMessageConfig("GETB", GetBlocksMessage{}),
        NewMessageConfig("GIVB", GiveBlocksMessage{}),
        NewMessageConfig("ANNB", AnnounceBlocksMessage{}),
        NewMessageConfig("GETT", GetTxnsMessage{}),
        NewMessageConfig("GIVT", GiveTxnsMessage{}),
        NewMessageConfig("ANNT", AnnounceTxnsMessage{}),
    }
}

type MessagesConfig struct {
    // Message ID prefices
    Messages []MessageConfig
}

func NewMessagesConfig() MessagesConfig {
    return MessagesConfig{
        Messages: getMessageConfigs(),
    }
}

// Registers our Messages with gnet
func (self *MessagesConfig) Register() {
    for _, mc := range self.Messages {
        gnet.RegisterMessage(mc.Prefix, mc.Message)
    }
    gnet.VerifyMessages()
}

type Messages struct {
    Config MessagesConfig
    // Magic value for detecting self-connection
    Mirror uint32
}

func NewMessages(c MessagesConfig) *Messages {
    return &Messages{
        Config: c,
        Mirror: rand.New(rand.NewSource(util.Now().UnixNano())).Uint32(),
    }
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
    ips, port, err := SplitAddr(addr)
    if err != nil {
        return
    }
    ipb := net.ParseIP(ips).To4()
    if ipb == nil {
        err = errors.New("Ignoring IPv6 address")
        return
    }
    ip := binary.BigEndian.Uint32(ipb)
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
    Process(d *Daemon)
}

// Sent to request peers
type GetPeersMessage struct {
    c *gnet.MessageContext `enc:"-"`
}

func NewGetPeersMessage() *GetPeersMessage {
    return &GetPeersMessage{}
}

func (self *GetPeersMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

// Notifies the Pex instance that peers were requested
func (self *GetPeersMessage) Process(d *Daemon) {
    if d.Peers.Config.Disabled {
        return
    }
    peers := d.Peers.Peers.Peerlist.RandomPublic(d.Peers.Config.ReplyCount)
    if len(peers) == 0 {
        logger.Debug("We have no peers to send in reply")
        return
    }
    m := NewGivePeersMessage(peers)
    d.Pool.Pool.SendMessage(self.c.Conn, m)
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
    peers := make([]string, len(self.Peers))
    for i, ipaddr := range self.Peers {
        peers[i] = ipaddr.String()
    }
    return peers
}

func (self *GivePeersMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

// Notifies the Pex instance that peers were received
func (self *GivePeersMessage) Process(d *Daemon) {
    if d.Peers.Config.Disabled {
        return
    }
    peers := self.GetPeers()
    if len(peers) != 0 {
        logger.Debug("Got these peers via PEX:")
        for _, p := range peers {
            logger.Debug("\t%s", p)
        }
    }
    d.Peers.Peers.AddPeers(peers)
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

func NewIntroductionMessage(mirror uint32, version int32,
    port uint16) *IntroductionMessage {
    return &IntroductionMessage{
        Mirror:  mirror,
        Version: version,
        Port:    port,
    }
}

// Responds to an gnet.Pool event. We implement Handle() here because we
// need to control the DisconnectReason sent back to gnet.  We still implement
// Process(), where we do modifications that are not threadsafe
func (self *IntroductionMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) (err error) {
    d := daemon.(*Daemon)
    addr := mc.Conn.Addr()
    // Disconnect if this is a self connection (we have the same mirror value)
    if self.Mirror == d.Messages.Mirror {
        logger.Info("Remote mirror value %v matches ours", self.Mirror)
        d.Pool.Pool.Disconnect(mc.Conn, DisconnectSelf)
        err = DisconnectSelf
    }
    // Disconnect if not running the same version
    if self.Version != d.Config.Version {
        logger.Info("%s has different version %d. Disconnecting.",
            addr, self.Version)
        d.Pool.Pool.Disconnect(mc.Conn, DisconnectInvalidVersion)
        err = DisconnectInvalidVersion
    } else {
        logger.Info("%s verified for version %d", addr, self.Version)
    }
    // Disconnect if connected twice to the same peer (judging by ip:mirror)
    knownPort, exists := d.getMirrorPort(addr, self.Mirror)
    if exists {
        logger.Info("%s is already connected on port %d", addr, knownPort)
        d.Pool.Pool.Disconnect(mc.Conn, DisconnectConnectedTwice)
        err = DisconnectConnectedTwice
    }

    self.valid = (err == nil)
    self.c = mc
    if err == nil {
        err = d.recordMessageEvent(self, mc)
    }
    return
}

// Processes an event queued by Handle()
func (self *IntroductionMessage) Process(d *Daemon) {
    delete(d.ExpectingIntroductions, self.c.Conn.Addr())
    if !self.valid {
        return
    }
    // Add the remote peer with their chosen listening port
    a := self.c.Conn.Addr()
    ip, _, err := SplitAddr(a)
    if err != nil {
        // This should never happen, but the program should still work if it
        // does.
        logger.Error("Invalid Addr() for connection: %s", a)
        d.Pool.Pool.Disconnect(self.c.Conn, DisconnectOtherError)
        return
    }
    // Record their listener, to avoid double connections
    err = d.recordConnectionMirror(a, self.Mirror)
    if err != nil {
        // This should never happen, but the program should not allow itself
        // to be corrupted in case it does
        logger.Error("Invalid port for connection %s", a)
        d.Pool.Pool.Disconnect(self.c.Conn, DisconnectOtherError)
        return
    }
    _, err = d.Peers.Peers.AddPeer(fmt.Sprintf("%s:%d", ip, self.Port))
    if err != nil {
        logger.Error("Failed to add peer: %v", err)
    }

    // Request blocks immediately after they're confirmed
    err = d.Visor.RequestBlocksFromAddr(d.Pool, self.c.Conn.Addr())
    if err == nil {
        logger.Debug("Successfully requested blocks from %s",
            self.c.Conn.Addr())
    } else {
        logger.Warning("%v", err)
    }
}

// Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
    c *gnet.MessageContext `enc:"-"`
}

func (self *PingMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

// Sends a PongMessage to the sender of PingMessage
func (self *PingMessage) Process(d *Daemon) {
    logger.Debug("Reply to ping from %s", self.c.Conn.Addr())
    d.Pool.Pool.SendMessage(self.c.Conn, &PongMessage{})
}

// Sent in reply to a PingMessage.  No action is taken when this is received.
type PongMessage struct {
}

func (self *PongMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    // There is nothing to do; gnet updates Connection.LastMessage internally
    // when this is received
    logger.Debug("Received pong from %s", mc.Conn.Addr())
    return nil
}
