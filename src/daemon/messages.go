package daemon

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util/utc"
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

// MessageConfig config contains a gnet.Message's 4byte prefix and a
// reference interface
type MessageConfig struct {
	Prefix  gnet.MessagePrefix
	Message interface{}
}

// NewMessageConfig creates message config
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

// MessagesConfig slice of MessageConfig
type MessagesConfig struct {
	// Message ID prefices
	Messages []MessageConfig
}

// NewMessagesConfig creates messages config
func NewMessagesConfig() MessagesConfig {
	return MessagesConfig{
		Messages: getMessageConfigs(),
	}
}

// Register registers our Messages with gnet
func (msc *MessagesConfig) Register() {
	for _, mc := range msc.Messages {
		gnet.RegisterMessage(mc.Prefix, mc.Message)
	}
	gnet.VerifyMessages()
}

// Messages messages struct
type Messages struct {
	Config MessagesConfig
	// Magic value for detecting self-connection
	Mirror uint32
}

// NewMessages creates Messages
func NewMessages(c MessagesConfig) *Messages {
	return &Messages{
		Config: c,
		Mirror: rand.New(rand.NewSource(utc.Now().UnixNano())).Uint32(),
	}
}

// IPAddr compact representation of IP:Port
type IPAddr struct {
	IP   uint32
	Port uint16
}

// NewIPAddr returns an IPAddr from an ip:port string.  If ipv6 or invalid, error is
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

// String returns IPAddr as "ip:port"
func (ipa IPAddr) String() string {
	ipb := make([]byte, 4)
	binary.BigEndian.PutUint32(ipb, ipa.IP)
	return fmt.Sprintf("%s:%d", net.IP(ipb).String(), ipa.Port)
}

// AsyncMessage messages that perform an action when received must implement this interface.
// Process() is called after the message is pulled off of messageEvent channel.
// Messages should place themselves on the messageEvent channel in their
// Handle() method required by gnet.
type AsyncMessage interface {
	Process(d *Daemon)
}

// GetPeersMessage sent to request peers
type GetPeersMessage struct {
	// c *gnet.MessageContext `enc:"-"`
	// connID int    `enc:"-"`
	addr string `enc:"-"`
}

// NewGetPeersMessage creates GetPeersMessage
func NewGetPeersMessage() *GetPeersMessage {
	return &GetPeersMessage{}
}

// Handle handles message
func (gpm *GetPeersMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	// self.connID = mc.ConnID
	gpm.addr = mc.Addr
	return daemon.(*Daemon).recordMessageEvent(gpm, mc)
}

// Process Notifies the Pex instance that peers were requested
func (gpm *GetPeersMessage) Process(d *Daemon) {
	if d.Peers.Config.Disabled {
		return
	}

	peers := d.Peers.Peers.RandomExchgPublic(d.Peers.Config.ReplyCount)
	if len(peers) == 0 {
		logger.Debug("We have no peers to send in reply")
		return
	}

	// logger.Info(fmt.Sprintf("give exchange peers:%+v", peers))

	m := NewGivePeersMessage(peers)
	d.Pool.Pool.SendMessage(gpm.addr, m)
}

// GivePeersMessage sent in response to GetPeersMessage
type GivePeersMessage struct {
	Peers []IPAddr
	c     *gnet.MessageContext `enc:"-"`
}

// NewGivePeersMessage []*pex.Peer is converted to []IPAddr for binary transmission
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
func (gpm *GivePeersMessage) GetPeers() []string {
	peers := make([]string, len(gpm.Peers))
	for i, ipaddr := range gpm.Peers {
		peers[i] = ipaddr.String()
	}
	return peers
}

// Handle handle message
func (gpm *GivePeersMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gpm.c = mc
	return daemon.(*Daemon).recordMessageEvent(gpm, mc)
}

// Process Notifies the Pex instance that peers were received
func (gpm *GivePeersMessage) Process(d *Daemon) {
	if d.Peers.Config.Disabled {
		return
	}
	peers := gpm.GetPeers()
	if len(peers) != 0 {
		logger.Debug("Got these peers via PEX: %s", strings.Join(peers, ", "))
	}
	d.Peers.Peers.AddPeers(peers)
}

// IntroductionMessage jan IntroductionMessage is sent on first connect by both parties
type IntroductionMessage struct {
	// Mirror is a random value generated on client startup that is used
	// to identify self-connections
	Mirror uint32
	// Port is the port that this client is listening on
	Port uint16
	// Our client version
	Version int32

	c *gnet.MessageContext `enc:"-"`
	// We validate the message in Handle() and cache the result for Process()
	valid bool `enc:"-"` // skip it during encoding
}

// NewIntroductionMessage creates introduction message
func NewIntroductionMessage(mirror uint32, version int32, port uint16) *IntroductionMessage {
	return &IntroductionMessage{
		Mirror:  mirror,
		Version: version,
		Port:    port,
	}
}

// Handle Responds to an gnet.Pool event. We implement Handle() here because we
// need to control the DisconnectReason sent back to gnet.  We still implement
// Process(), where we do modifications that are not threadsafe
func (intro *IntroductionMessage) Handle(mc *gnet.MessageContext, daemon interface{}) (err error) {
	d := daemon.(*Daemon)
	addr := mc.Addr
	// Disconnect if this is a self connection (we have the same mirror value)
	if intro.Mirror == d.Messages.Mirror {
		logger.Info("Remote mirror value %v matches ours", intro.Mirror)
		d.Pool.Pool.Disconnect(mc.Addr, ErrDisconnectSelf)
		err = ErrDisconnectSelf
	}
	// Disconnect if not running the same version
	if intro.Version != d.Config.Version {
		logger.Info("%s has different version %d. Disconnecting.",
			addr, intro.Version)
		d.Pool.Pool.Disconnect(mc.Addr, ErrDisconnectInvalidVersion)
		err = ErrDisconnectInvalidVersion
	} else {
		logger.Info("%s verified for version %d", addr, intro.Version)
	}

	// only solicited connection can be added to exchange peer list, cause accepted
	// connection may not have incomming  port.
	ip, port, err := SplitAddr(mc.Addr)
	if err != nil {
		// This should never happen, but the program should still work if it
		// does.
		logger.Error("Invalid Addr() for connection: %s", mc.Addr)
		d.Pool.Pool.Disconnect(intro.c.Addr, ErrDisconnectOtherError)
		err = ErrDisconnectOtherError
	}

	if port == intro.Port {
		if err := d.Peers.Peers.SetPeerHasInPort(mc.Addr, true); err != nil {
			logger.Error("Failed to set peer hasInPort statue, %v", err)
		}
	} else {
		_, err = d.Peers.Peers.AddPeer(fmt.Sprintf("%s:%d", ip, intro.Port))
		if err != nil {
			logger.Error("Failed to add peer: %v", err)
		}
	}

	// Disconnect if connected twice to the same peer (judging by ip:mirror)
	knownPort, exists := d.getMirrorPort(addr, intro.Mirror)
	if exists {
		logger.Info("%s is already connected on port %d", addr, knownPort)
		d.Pool.Pool.Disconnect(mc.Addr, ErrDisconnectConnectedTwice)
		err = ErrDisconnectConnectedTwice
	}

	intro.valid = (err == nil)
	intro.c = mc
	if err == nil {
		err = d.recordMessageEvent(intro, mc)
		d.Peers.Peers.ResetRetryTimes(mc.Addr)
	} else {
		d.Peers.Peers.IncreaseRetryTimes(mc.Addr)
		d.expectingIntroductions.Remove(mc.Addr)
	}
	return
}

// Process an event queued by Handle()
func (intro *IntroductionMessage) Process(d *Daemon) {
	d.expectingIntroductions.Remove(intro.c.Addr)
	if !intro.valid {
		return
	}
	// Add the remote peer with their chosen listening port
	a := intro.c.Addr

	// Record their listener, to avoid double connections
	err := d.recordConnectionMirror(a, intro.Mirror)
	if err != nil {
		// This should never happen, but the program should not allow itself
		// to be corrupted in case it does
		logger.Error("Invalid port for connection %s", a)
		d.Pool.Pool.Disconnect(intro.c.Addr, ErrDisconnectOtherError)
		return
	}

	// Request blocks immediately after they're confirmed
	err = d.Visor.RequestBlocksFromAddr(d.Pool, intro.c.Addr)
	if err == nil {
		logger.Debug("Successfully requested blocks from %s", intro.c.Addr)
	} else {
		logger.Warning("%v", err)
	}

	// Anounce unconfirmed know txns
	d.Visor.AnnounceAllTxns(d.Pool)
}

// PingMessage Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
	c *gnet.MessageContext `enc:"-"`
}

// Handle implements the Messager interface
func (ping *PingMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	ping.c = mc
	return daemon.(*Daemon).recordMessageEvent(ping, mc)
}

// Process Sends a PongMessage to the sender of PingMessage
func (ping *PingMessage) Process(d *Daemon) {
	if d.Config.LogPings {
		logger.Debug("Reply to ping from %s", ping.c.Addr)
	}
	d.Pool.Pool.SendMessage(ping.c.Addr, &PongMessage{})
}

// PongMessage Sent in reply to a PingMessage.  No action is taken when this is received.
type PongMessage struct {
}

// Handle handles message
func (pong *PongMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	// There is nothing to do; gnet updates Connection.LastMessage internally
	// when this is received
	if daemon.(*Daemon).Config.LogPings {
		logger.Debug("Received pong from %s", mc.Addr)
	}
	return nil
}
