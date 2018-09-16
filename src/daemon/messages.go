package daemon

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util/iputil"
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
		Mirror: rand.New(rand.NewSource(time.Now().UTC().UnixNano())).Uint32(),
	}
}

// IPAddr compact representation of IP:Port
type IPAddr struct {
	IP   uint32
	Port uint16
}

// NewIPAddr returns an IPAddr from an ip:port string.
func NewIPAddr(addr string) (ipaddr IPAddr, err error) {
	ips, port, err := iputil.SplitAddr(addr)
	if err != nil {
		return
	}

	// TODO -- support ipv6
	ipb := net.ParseIP(ips).To4()
	if ipb == nil {
		err = errors.New("Ignoring IPv6 address")
		return
	}

	ip := binary.BigEndian.Uint32(ipb)
	ipaddr.IP = ip
	ipaddr.Port = port
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
	Process(d Daemoner)
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
func (gpm *GetPeersMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	// self.connID = mc.ConnID
	gpm.addr = mc.Addr
	return daemon.(Daemoner).RecordMessageEvent(gpm, mc)
}

// Process Notifies the Pex instance that peers were requested
func (gpm *GetPeersMessage) Process(d Daemoner) {
	if d.PexConfig().Disabled {
		return
	}

	peers := d.RandomExchangeable(d.PexConfig().ReplyCount)
	if len(peers) == 0 {
		logger.Debug("We have no peers to send in reply")
		return
	}

	m := NewGivePeersMessage(peers)
	if err := d.SendMessage(gpm.addr, m); err != nil {
		logger.Errorf("Send GivePeersMessage to %s failed: %v", gpm.addr, err)
	}
}

// GivePeersMessage sent in response to GetPeersMessage
type GivePeersMessage struct {
	Peers []IPAddr
	c     *gnet.MessageContext `enc:"-"`
}

// NewGivePeersMessage []*pex.Peer is converted to []IPAddr for binary transmission
func NewGivePeersMessage(peers []pex.Peer) *GivePeersMessage {
	ipaddrs := make([]IPAddr, 0, len(peers))
	for _, ps := range peers {
		ipaddr, err := NewIPAddr(ps.Addr)
		if err != nil {
			logger.Warningf("GivePeersMessage skipping address %s", ps.Addr)
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
	return daemon.(Daemoner).RecordMessageEvent(gpm, mc)
}

// Process Notifies the Pex instance that peers were received
func (gpm *GivePeersMessage) Process(d Daemoner) {
	if d.PexConfig().Disabled {
		return
	}
	peers := gpm.GetPeers()
	logger.Debugf("Got these peers via PEX: %s", strings.Join(peers, ", "))

	d.AddPeers(peers)
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
	c       *gnet.MessageContext `enc:"-"`
	// We validate the message in Handle() and cache the result for Process()
	valid bool `enc:"-"` // skip it during encoding
	// Extra would be parsed as blockchain pubkey if it's not empty
	Extra []byte `enc:",omitempty"`
}

// NewIntroductionMessage creates introduction message
func NewIntroductionMessage(mirror uint32, version int32, port uint16, extra []byte) *IntroductionMessage {
	return &IntroductionMessage{
		Mirror:  mirror,
		Version: version,
		Port:    port,
		Extra:   extra,
	}
}

// Handle Responds to an gnet.Pool event. We implement Handle() here because we
// need to control the DisconnectReason sent back to gnet.  We still implement
// Process(), where we do modifications that are not threadsafe
func (intro *IntroductionMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	d := daemon.(Daemoner)

	err := func() error {
		// Disconnect if this is a self connection (we have the same mirror value)
		if intro.Mirror == d.Mirror() {
			logger.Infof("Remote mirror value %v matches ours", intro.Mirror)
			if err := d.Disconnect(mc.Addr, ErrDisconnectSelf); err != nil {
				logger.WithError(err).WithField("addr", mc.Addr).Warning("Disconnect")
			}
			return ErrDisconnectSelf
		}

		// Disconnect if not running the same version
		if intro.Version != d.DaemonConfig().Version {
			logger.Infof("%s has different version %d. Disconnecting.",
				mc.Addr, intro.Version)
			if err := d.Disconnect(mc.Addr, ErrDisconnectInvalidVersion); err != nil {
				logger.WithError(err).WithField("addr", mc.Addr).Warning("Disconnect")
			}
			return ErrDisconnectInvalidVersion
		}

		logger.Infof("%s verified for version %d", mc.Addr, intro.Version)

		// Checks the genesis hash if not empty
		if len(intro.Extra) > 0 {
			var bcPubKey cipher.PubKey
			if len(intro.Extra) < len(bcPubKey) {
				logger.Infof("Extra data length does not meet the minimum requirement")
				if err := d.Disconnect(mc.Addr, ErrDisconnectInvalidExtraData); err != nil {
					logger.WithError(err).WithField("addr", mc.Addr).Warning("Disconnect")
				}
				return ErrDisconnectInvalidExtraData
			}
			copy(bcPubKey[:], intro.Extra[:len(bcPubKey)])

			if d.BlockchainPubkey() != bcPubKey {
				logger.Infof("Blockchain pubkey does not match, local: %s, remote: %s", d.BlockchainPubkey().Hex(), bcPubKey.Hex())
				if err := d.Disconnect(mc.Addr, ErrDisconnectBlockchainPubkeyNotMatched); err != nil {
					logger.WithError(err).WithField("addr", mc.Addr).Warning("Disconnect")
				}
				return ErrDisconnectBlockchainPubkeyNotMatched
			}
		}

		// only solicited connection can be added to exchange peer list, cause accepted
		// connection may not have incomming  port.
		ip, port, err := iputil.SplitAddr(mc.Addr)
		if err != nil {
			// This should never happen, but the program should still work if it
			// does.
			logger.Errorf("Invalid Addr() for connection: %s", mc.Addr)
			if err := d.Disconnect(mc.Addr, ErrDisconnectOtherError); err != nil {
				logger.WithError(err).WithField("addr", mc.Addr).Warning("Disconnect")
			}
			return ErrDisconnectOtherError
		}

		// Checks if the introduction message is from outgoing connection.
		// It's outgoing connection if port == intro.Port, as the incoming
		// connection's port is a random port, it's different from the port
		// in introduction message.
		if port == intro.Port {
			if err := d.SetHasIncomingPort(mc.Addr); err != nil {
				logger.Errorf("Failed to set peer has incoming port status, %v", err)
			}
		} else {
			if err := d.AddPeer(fmt.Sprintf("%s:%d", ip, intro.Port)); err != nil {
				logger.Errorf("Failed to add peer: %v", err)
			}
		}

		// Disconnect if connected twice to the same peer (judging by ip:mirror)
		knownPort, exists := d.GetMirrorPort(mc.Addr, intro.Mirror)
		if exists {
			logger.Infof("%s is already connected on port %d", mc.Addr, knownPort)
			if err := d.Disconnect(mc.Addr, ErrDisconnectConnectedTwice); err != nil {
				logger.WithError(err).WithField("addr", mc.Addr).Warning("Disconnect")
			}
			return ErrDisconnectConnectedTwice
		}
		return nil
	}()

	intro.valid = (err == nil)
	intro.c = mc

	if err != nil {
		d.IncreaseRetryTimes(mc.Addr)
		d.RemoveFromExpectingIntroductions(mc.Addr)
		return err
	}

	err = d.RecordMessageEvent(intro, mc)
	d.ResetRetryTimes(mc.Addr)
	return err
}

// Process an event queued by Handle()
func (intro *IntroductionMessage) Process(d Daemoner) {
	d.RemoveFromExpectingIntroductions(intro.c.Addr)
	if !intro.valid {
		return
	}
	// Add the remote peer with their chosen listening port
	a := intro.c.Addr

	// Record their listener, to avoid double connections
	err := d.RecordConnectionMirror(a, intro.Mirror)
	if err != nil {
		// This should never happen, but the program should not allow itself
		// to be corrupted in case it does
		logger.Errorf("Invalid port for connection %s", a)
		if err := d.Disconnect(intro.c.Addr, ErrDisconnectOtherError); err != nil {
			logger.WithError(err).WithField("addr", intro.c.Addr).Warning("Disconnect")
		}
		return
	}

	// Request blocks immediately after they're confirmed
	err = d.RequestBlocksFromAddr(intro.c.Addr)
	if err == nil {
		logger.Debugf("Successfully requested blocks from %s", intro.c.Addr)
	} else {
		logger.Warning(err)
	}

	// Announce unconfirmed txns
	if err := d.AnnounceAllTxns(); err != nil {
		logger.WithError(err).Warning("AnnounceAllTxns failed")
	}
}

// PingMessage Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
	c *gnet.MessageContext `enc:"-"`
}

// Handle implements the Messager interface
func (ping *PingMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	ping.c = mc
	return daemon.(Daemoner).RecordMessageEvent(ping, mc)
}

// Process Sends a PongMessage to the sender of PingMessage
func (ping *PingMessage) Process(d Daemoner) {
	if d.DaemonConfig().LogPings {
		logger.Debugf("Reply to ping from %s", ping.c.Addr)
	}
	if err := d.SendMessage(ping.c.Addr, &PongMessage{}); err != nil {
		logger.Errorf("Send PongMessage to %s failed: %v", ping.c.Addr, err)
	}
}

// PongMessage Sent in reply to a PingMessage.  No action is taken when this is received.
type PongMessage struct {
}

// Handle handles message
func (pong *PongMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	// There is nothing to do; gnet updates Connection.LastMessage internally
	// when this is received
	if daemon.(*Daemon).Config.LogPings {
		logger.Debugf("Received pong from %s", mc.Addr)
	}
	return nil
}

// GetBlocksMessage sent to request blocks since LastBlock
type GetBlocksMessage struct {
	LastBlock       uint64
	RequestedBlocks uint64
	c               *gnet.MessageContext `enc:"-"`
}

// NewGetBlocksMessage creates GetBlocksMessage
func NewGetBlocksMessage(lastBlock uint64, requestedBlocks uint64) *GetBlocksMessage {
	return &GetBlocksMessage{
		LastBlock:       lastBlock,
		RequestedBlocks: requestedBlocks, // count of blocks requested
	}
}

// Handle handles message
func (gbm *GetBlocksMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	gbm.c = mc
	return daemon.(Daemoner).RecordMessageEvent(gbm, mc)
}

// Process should send number to be requested, with request
func (gbm *GetBlocksMessage) Process(d Daemoner) {
	// TODO -- we need the sig to be sent with the block, but only the master
	// can sign blocks.  Thus the sig needs to be stored with the block.
	if d.DaemonConfig().DisableNetworking {
		return
	}
	// Record this as this peer's highest block
	d.RecordPeerHeight(gbm.c.Addr, gbm.LastBlock)
	// Fetch and return signed blocks since LastBlock
	blocks, err := d.GetSignedBlocksSince(gbm.LastBlock, gbm.RequestedBlocks)
	if err != nil {
		logger.Infof("Get signed blocks failed: %v", err)
		return
	}

	if len(blocks) == 0 {
		return
	}

	logger.Debugf("Got %d blocks since %d", len(blocks), gbm.LastBlock)

	m := NewGiveBlocksMessage(blocks)
	if err := d.SendMessage(gbm.c.Addr, m); err != nil {
		logger.Errorf("Send GiveBlocksMessage to %s failed: %v", gbm.c.Addr, err)
	}
}

// GiveBlocksMessage sent in response to GetBlocksMessage, or unsolicited
type GiveBlocksMessage struct {
	Blocks []coin.SignedBlock
	c      *gnet.MessageContext `enc:"-"`
}

// NewGiveBlocksMessage creates GiveBlocksMessage
func NewGiveBlocksMessage(blocks []coin.SignedBlock) *GiveBlocksMessage {
	return &GiveBlocksMessage{
		Blocks: blocks,
	}
}

// Handle handle message
func (gbm *GiveBlocksMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gbm.c = mc
	return daemon.(Daemoner).RecordMessageEvent(gbm, mc)
}

// Process process message
func (gbm *GiveBlocksMessage) Process(d Daemoner) {
	if d.DaemonConfig().DisableNetworking {
		logger.Critical().Info("Visor disabled, ignoring GiveBlocksMessage")
		return
	}

	// These DB queries are not performed in a transaction for performance reasons.
	// It is not necessary that the blocks be executed together in a single transaction.

	processed := 0
	maxSeq, ok, err := d.HeadBkSeq()
	if err != nil {
		logger.WithError(err).Error("visor.HeadBkSeq failed")
		return
	}
	if !ok {
		logger.Error("No HeadBkSeq found, cannot execute blocks")
		return
	}

	for _, b := range gbm.Blocks {
		// To minimize waste when receiving multiple responses from peers
		// we only break out of the loop if the block itself is invalid.
		// E.g. if we request 20 blocks since 0 from 2 peers, and one peer
		// replies with 15 and the other 20, if we did not do this check and
		// the reply with 15 was received first, we would toss the one with 20
		// even though we could process it at the time.
		if b.Seq() <= maxSeq {
			continue
		}

		err := d.ExecuteSignedBlock(b)
		if err == nil {
			logger.Critical().Infof("Added new block %d", b.Block.Head.BkSeq)
			processed++
		} else {
			logger.Critical().Errorf("Failed to execute received block %d: %v", b.Block.Head.BkSeq, err)
			// Blocks must be received in order, so if one fails its assumed
			// the rest are failing
			break
		}
	}
	if processed == 0 {
		return
	}

	headBkSeq, ok, err := d.HeadBkSeq()
	if err != nil {
		logger.WithError(err).Error("visor.HeadBkSeq failed")
		return
	}
	if !ok {
		logger.Error("No HeadBkSeq found after executing blocks, will not announce blocks")
		return
	}

	if headBkSeq < maxSeq {
		logger.Critical().Warning("HeadBkSeq decreased after executing blocks")
	} else if headBkSeq-maxSeq != uint64(processed) {
		logger.Critical().Warning("HeadBkSeq increased by %d but we processed %s blocks", headBkSeq-maxSeq, processed)
	}

	// Announce our new blocks to peers
	m1 := NewAnnounceBlocksMessage(headBkSeq)
	if err := d.BroadcastMessage(m1); err != nil {
		logger.WithError(err).Warning("Broadcast AnnounceBlocksMessage failed")
	}

	// Request more blocks
	m2 := NewGetBlocksMessage(headBkSeq, d.DaemonConfig().BlocksResponseCount)
	if err := d.BroadcastMessage(m2); err != nil {
		logger.WithError(err).Warning("Broadcast GetBlocksMessage failed")
	}
}

// AnnounceBlocksMessage tells a peer our highest known BkSeq. The receiving peer can choose
// to send GetBlocksMessage in response
type AnnounceBlocksMessage struct {
	MaxBkSeq uint64
	c        *gnet.MessageContext `enc:"-"`
}

// NewAnnounceBlocksMessage creates message
func NewAnnounceBlocksMessage(seq uint64) *AnnounceBlocksMessage {
	return &AnnounceBlocksMessage{
		MaxBkSeq: seq,
	}
}

// Handle handles message
func (abm *AnnounceBlocksMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	abm.c = mc
	return daemon.(Daemoner).RecordMessageEvent(abm, mc)
}

// Process process message
func (abm *AnnounceBlocksMessage) Process(d Daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}

	headBkSeq, ok, err := d.HeadBkSeq()
	if err != nil {
		logger.WithError(err).Error("AnnounceBlocksMessage Visor.HeadBkSeq failed")
		return
	}
	if !ok {
		logger.Error("AnnounceBlocksMessage no head block, cannot process AnnounceBlocksMessage")
		return
	}

	if headBkSeq >= abm.MaxBkSeq {
		return
	}

	// TODO: Should this be block get request for current sequence?
	// If client is not caught up, won't attempt to get block
	m := NewGetBlocksMessage(headBkSeq, d.DaemonConfig().BlocksResponseCount)
	if err := d.SendMessage(abm.c.Addr, m); err != nil {
		logger.Errorf("Send GetBlocksMessage to %s failed: %v", abm.c.Addr, err)
	}
}

// SendingTxnsMessage send transaction message interface
type SendingTxnsMessage interface {
	GetFiltered() []cipher.SHA256
}

// AnnounceTxnsMessage tells a peer that we have these transactions
type AnnounceTxnsMessage struct {
	Transactions []cipher.SHA256
	c            *gnet.MessageContext `enc:"-"`
}

// NewAnnounceTxnsMessage creates announce txns message
func NewAnnounceTxnsMessage(txns []cipher.SHA256) *AnnounceTxnsMessage {
	return &AnnounceTxnsMessage{
		Transactions: txns,
	}
}

// GetFiltered returns txns
func (atm *AnnounceTxnsMessage) GetFiltered() []cipher.SHA256 {
	return atm.Transactions
}

// Handle handle message
func (atm *AnnounceTxnsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	atm.c = mc
	return daemon.(Daemoner).RecordMessageEvent(atm, mc)
}

// Process process message
func (atm *AnnounceTxnsMessage) Process(d Daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}

	unknown, err := d.GetUnconfirmedUnknown(atm.Transactions)
	if err != nil {
		logger.WithError(err).Error("AnnounceTxnsMessage Visor.GetUnconfirmedUnknown failed")
		return
	}

	if len(unknown) == 0 {
		return
	}

	m := NewGetTxnsMessage(unknown)
	if err := d.SendMessage(atm.c.Addr, m); err != nil {
		logger.Errorf("Send GetTxnsMessage to %s failed: %v", atm.c.Addr, err)
	}
}

// GetTxnsMessage request transactions of given hash
type GetTxnsMessage struct {
	Transactions []cipher.SHA256
	c            *gnet.MessageContext `enc:"-"`
}

// NewGetTxnsMessage creates GetTxnsMessage
func NewGetTxnsMessage(txns []cipher.SHA256) *GetTxnsMessage {
	return &GetTxnsMessage{
		Transactions: txns,
	}
}

// Handle handle message
func (gtm *GetTxnsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gtm.c = mc
	return daemon.(Daemoner).RecordMessageEvent(gtm, mc)
}

// Process process message
func (gtm *GetTxnsMessage) Process(d Daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}

	// Locate all txns from the unconfirmed pool
	known, err := d.GetUnconfirmedKnown(gtm.Transactions)
	if err != nil {
		logger.WithError(err).Error("GetTxnsMessage Visor.GetUnconfirmedKnown failed")
		return
	}
	if len(known) == 0 {
		return
	}

	// Reply to sender with GiveTxnsMessage
	m := NewGiveTxnsMessage(known)
	if err := d.SendMessage(gtm.c.Addr, m); err != nil {
		logger.Errorf("Send GiveTxnsMessage to %s failed: %v", gtm.c.Addr, err)
	}
}

// GiveTxnsMessage tells the transaction of given hashes
type GiveTxnsMessage struct {
	Transactions coin.Transactions
	c            *gnet.MessageContext `enc:"-"`
}

// NewGiveTxnsMessage creates GiveTxnsMessage
func NewGiveTxnsMessage(txns coin.Transactions) *GiveTxnsMessage {
	return &GiveTxnsMessage{
		Transactions: txns,
	}
}

// GetFiltered returns transactions hashes
func (gtm *GiveTxnsMessage) GetFiltered() []cipher.SHA256 {
	return gtm.Transactions.Hashes()
}

// Handle handle message
func (gtm *GiveTxnsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gtm.c = mc
	return daemon.(Daemoner).RecordMessageEvent(gtm, mc)
}

// Process process message
func (gtm *GiveTxnsMessage) Process(d Daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}

	hashes := make([]cipher.SHA256, 0, len(gtm.Transactions))
	// Update unconfirmed pool with these transactions
	for _, txn := range gtm.Transactions {
		// Only announce transactions that are new to us, so that peers can't spam relays
		// It is not necessary to inject all of the transactions inside a database transaction,
		// since each is independent
		known, softErr, err := d.InjectTransaction(txn)
		if err != nil {
			logger.Warningf("Failed to record transaction %s: %v", txn.Hash().Hex(), err)
			continue
		} else if softErr != nil {
			logger.Warningf("Transaction soft violation: %v", err)
			continue
		} else if known {
			logger.Warningf("Duplicate Transaction: %s", txn.Hash().Hex())
			continue
		}

		hashes = append(hashes, txn.Hash())
	}

	// Announce these transactions to peers
	if len(hashes) != 0 {
		logger.Debugf("Announce %d transactions", len(hashes))
		m := NewAnnounceTxnsMessage(hashes)
		if err := d.BroadcastMessage(m); err != nil {
			logger.WithError(err).Warning("Broadcast AnnounceTxnsMessage failed")
		}
	}
}
