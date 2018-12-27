package daemon

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/iputil"
	"github.com/skycoin/skycoin/src/util/useragent"
)

// Message represent a packet to be serialized over the network by
// the gnet encoder.
// They must implement the gnet.Message interface
// All concurrent daemon write operations are synchronized by the daemon's
// DaemonLoop().
// Message do this by caching the gnet.MessageContext received in Handle()
// and placing itself on the messageEvent channel.
// When the message is retrieved from the messageEvent channel, its process()
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
		NewMessageConfig("DISC", DisconnectMessage{}),
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
}

// NewMessages creates Messages
func NewMessages(c MessagesConfig) *Messages {
	return &Messages{
		Config: c,
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

// asyncMessage messages that perform an action when received must implement this interface.
// process() is called after the message is pulled off of messageEvent channel.
// Messages should place themselves on the messageEvent channel in their
// Handle() method required by gnet.
type asyncMessage interface {
	process(d daemoner)
}

// GetPeersMessage sent to request peers
type GetPeersMessage struct {
	addr string `enc:"-"`
}

// NewGetPeersMessage creates GetPeersMessage
func NewGetPeersMessage() *GetPeersMessage {
	return &GetPeersMessage{}
}

// Handle handles message
func (gpm *GetPeersMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gpm.addr = mc.Addr
	return daemon.(daemoner).recordMessageEvent(gpm, mc)
}

// process Notifies the Pex instance that peers were requested
func (gpm *GetPeersMessage) process(d daemoner) {
	if d.pexConfig().Disabled {
		return
	}

	if err := d.sendRandomPeers(gpm.addr); err != nil {
		logger.WithField("addr", gpm.addr).WithError(err).Error("SendRandomPeers failed")
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
			logger.WithError(err).WithField("addr", ps.Addr).Warning("GivePeersMessage skipping invalid address")
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
	return daemon.(daemoner).recordMessageEvent(gpm, mc)
}

// process Notifies the Pex instance that peers were received
func (gpm *GivePeersMessage) process(d daemoner) {
	if d.pexConfig().Disabled {
		return
	}

	peers := gpm.GetPeers()

	if len(peers) == 0 {
		return
	}

	// Cap the number of peers printed in the log to prevent log spam abuse
	peersToFmt := peers
	if len(peersToFmt) > 30 {
		peersToFmt = peersToFmt[:30]
	}
	peersStr := strings.Join(peersToFmt, ", ")
	if len(peers) != len(peersToFmt) {
		peersStr += fmt.Sprintf(" and %d more", len(peers)-len(peersToFmt))
	}

	logger.WithFields(logrus.Fields{
		"addr":   gpm.c.Addr,
		"gnetID": gpm.c.ConnID,
		"peers":  peersStr,
		"count":  len(peers),
	}).Debug("Received peers via PEX")

	d.addPeers(peers)
}

// IntroductionMessage is sent on first connect by both parties
type IntroductionMessage struct {
	c                    *gnet.MessageContext `enc:"-"`
	userAgent            useragent.Data       `enc:"-"`
	unconfirmedVerifyTxn params.VerifyTxn     `enc:"-"`

	// Mirror is a random value generated on client startup that is used to identify self-connections
	Mirror uint32
	// ListenPort is the port that this client is listening on
	ListenPort uint16
	// Protocol version
	ProtocolVersion int32

	// Extra is extra bytes added to the struct to accommodate multiple versions of this packet.
	// Currently it contains the blockchain pubkey and user agent but will accept a client that does not provide it.
	// If any of this data is provided, it must include a valid blockchain pubkey and a valid user agent string (maxlen=256).
	// Contents of extra:
	// ExtraByte           uint32 // length prefix of []byte
	// Pubkey              cipher.Pubkey // blockchain pubkey
	// BurnFactor          uint32 // burn factor for announced txns
	// MaxTxnSize          uint32 // max txn size for announced txns
	// MaxDropletPrecision uint8 // maximum number of decimal places for announced txns
	// UserAgent           string `enc:",maxlen=256"`
	Extra []byte `enc:",omitempty"`
}

// NewIntroductionMessage creates introduction message
func NewIntroductionMessage(mirror uint32, version int32, port uint16, pubkey cipher.PubKey, userAgent string, verifyParams params.VerifyTxn) *IntroductionMessage {
	return &IntroductionMessage{
		Mirror:          mirror,
		ProtocolVersion: version,
		ListenPort:      port,
		Extra:           newIntroductionMessageExtra(pubkey, userAgent, verifyParams),
	}
}

func newIntroductionMessageExtra(pubkey cipher.PubKey, userAgent string, verifyParams params.VerifyTxn) []byte {
	if len(userAgent) > useragent.MaxLen {
		logger.WithFields(logrus.Fields{
			"userAgent": userAgent,
			"maxLen":    useragent.MaxLen,
		}).Panic("user agent exceeds max len")
	}
	if userAgent == "" {
		logger.Panic("user agent is required")
	}
	useragent.MustParse(userAgent)

	if err := verifyParams.Validate(); err != nil {
		logger.Panic(err)
	}

	userAgentSerialized := encoder.SerializeString(userAgent)
	verifyParamsSerialized := encoder.Serialize(verifyParams)

	extra := make([]byte, len(pubkey)+len(userAgentSerialized)+len(verifyParamsSerialized))

	copy(extra[:len(pubkey)], pubkey[:])
	i := len(pubkey)
	copy(extra[i:], verifyParamsSerialized)
	i += len(verifyParamsSerialized)
	copy(extra[i:], userAgentSerialized)

	return extra
}

// Handle records message event in daemon
func (intro *IntroductionMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	intro.c = mc
	return daemon.(daemoner).recordMessageEvent(intro, mc)
}

// process an event queued by Handle()
func (intro *IntroductionMessage) process(d daemoner) {
	addr := intro.c.Addr

	fields := logrus.Fields{
		"addr":       addr,
		"gnetID":     intro.c.ConnID,
		"listenPort": intro.ListenPort,
	}

	logger.WithFields(fields).Debug("IntroductionMessage.process")

	if err := intro.verify(d); err != nil {
		if err := d.Disconnect(addr, err); err != nil {
			logger.WithError(err).WithFields(fields).Warning("Disconnect")
		}
		return
	}

	if _, err := d.connectionIntroduced(addr, intro.c.ConnID, intro); err != nil {
		logger.WithError(err).WithFields(fields).Warning("connectionIntroduced failed")
		var reason gnet.DisconnectReason
		switch err {
		// It is hypothetically possible that a message would get processed after
		// a disconnect event for a given connection.
		// In this case, drop the packet.
		// Do not perform a disconnect, since this would operate on the new connection.
		// This should be prevented by an earlier check in daemon.onMessageEvent()
		case ErrConnectionGnetIDMismatch, ErrConnectionStateNotConnected, ErrConnectionAlreadyIntroduced:
			logger.Critical().WithError(err).WithFields(fields).Warning("IntroductionMessage.process connection state out of order")
			return
		case ErrConnectionNotExist:
			return
		case ErrConnectionIPMirrorExists:
			reason = ErrDisconnectConnectedTwice
		case pex.ErrPeerlistFull:
			reason = ErrDisconnectPeerlistFull
			// Send more peers before disconnecting
			logger.WithFields(fields).Debug("Sending peers before disconnecting due to peer list full")
			if err := d.sendRandomPeers(addr); err != nil {
				logger.WithError(err).WithFields(fields).Warning("sendRandomPeers failed")
			}
		default:
			reason = ErrDisconnectUnexpectedError
		}

		if err := d.Disconnect(addr, reason); err != nil {
			logger.WithError(err).WithFields(fields).Warning("Disconnect")
		}

		return
	}

	// Request blocks immediately after they're confirmed
	if err := d.requestBlocksFromAddr(addr); err != nil {
		logger.WithError(err).WithFields(fields).Warning("requestBlocksFromAddr")
	} else {
		logger.WithFields(fields).Debug("Requested blocks")
	}

	// Announce unconfirmed txns
	if err := d.announceAllValidTxns(); err != nil {
		logger.WithError(err).Warning("announceAllValidTxns failed")
	}
}

func (intro *IntroductionMessage) verify(d daemoner) error {
	addr := intro.c.Addr

	fields := logrus.Fields{
		"addr":   addr,
		"gnetID": intro.c.ConnID,
	}

	dc := d.daemonConfig()

	// Disconnect if this is a self connection (we have the same mirror value)
	if intro.Mirror == dc.Mirror {
		logger.WithFields(fields).WithField("mirror", intro.Mirror).Info("Remote mirror value matches ours")
		return ErrDisconnectSelf
	}

	// Disconnect if peer version is not within the supported range
	if intro.ProtocolVersion < dc.MinProtocolVersion {
		logger.WithFields(fields).WithFields(logrus.Fields{
			"protocolVersion":    intro.ProtocolVersion,
			"minProtocolVersion": dc.MinProtocolVersion,
		}).Info("protocol version below minimum supported protocol version")
		return ErrDisconnectVersionNotSupported
	}

	logger.WithFields(fields).WithField("protocolVersion", intro.ProtocolVersion).Debug("Peer protocol version accepted")

	// v24 does not send blockchain pubkey or user agent
	// v25 sends blockchain pubkey and user agent
	// v24 and v25 check the blockchain pubkey and user agent, would accept message with no Pubkey and user agent
	// v26 would check the blockchain pubkey and reject if not matched or not provided, and parses a user agent
	if len(intro.Extra) > 0 {
		var bcPubKey cipher.PubKey
		if len(intro.Extra) < len(bcPubKey) {
			logger.WithFields(fields).Warning("Extra data length does not meet the minimum requirement")
			return ErrDisconnectInvalidExtraData
		}
		copy(bcPubKey[:], intro.Extra[:len(bcPubKey)])

		if dc.BlockchainPubkey != bcPubKey {
			logger.WithFields(fields).WithFields(logrus.Fields{
				"pubkey":       bcPubKey.Hex(),
				"daemonPubkey": dc.BlockchainPubkey.Hex(),
			}).Warning("Blockchain pubkey does not match")
			return ErrDisconnectBlockchainPubkeyNotMatched
		}

		i := len(bcPubKey)
		if len(intro.Extra) < i+9 {
			logger.WithFields(fields).Warning("IntroductionMessage transaction verification parameters could not be deserialized: not enough data")
			return ErrDisconnectInvalidExtraData
		}
		if err := encoder.DeserializeRaw(intro.Extra[i:i+9], &intro.unconfirmedVerifyTxn); err != nil {
			// This should not occur due to the previous length check
			logger.Critical().WithError(err).WithFields(fields).Warning("unconfirmedVerifyTxn params could not be deserialized")
			return ErrDisconnectInvalidExtraData
		}

		if err := intro.unconfirmedVerifyTxn.Validate(); err != nil {
			logger.WithError(err).WithFields(fields).WithFields(logrus.Fields{
				"burnFactor":          intro.unconfirmedVerifyTxn.BurnFactor,
				"maxTransactionSize":  intro.unconfirmedVerifyTxn.MaxTransactionSize,
				"maxDropletPrecision": intro.unconfirmedVerifyTxn.MaxDropletPrecision,
			}).Warning("Invalid unconfirmedVerifyTxn params")
			switch err {
			case params.ErrInvalidBurnFactor:
				return ErrDisconnectInvalidBurnFactor
			case params.ErrInvalidMaxTransactionSize:
				return ErrDisconnectInvalidMaxTransactionSize
			case params.ErrInvalidMaxDropletPrecision:
				return ErrDisconnectInvalidMaxDropletPrecision
			default:
				return ErrDisconnectUnexpectedError
			}
		}

		userAgentSerialized := intro.Extra[len(bcPubKey)+9:]
		userAgent, _, err := encoder.DeserializeString(userAgentSerialized, useragent.MaxLen)
		if err != nil {
			logger.WithError(err).WithFields(fields).Warning("Extra data user agent string could not be deserialized")
			return ErrDisconnectInvalidExtraData
		}

		intro.userAgent, err = useragent.Parse(useragent.Sanitize(userAgent))
		if err != nil {
			logger.WithError(err).WithFields(fields).WithField("userAgent", userAgent).Warning("User agent is invalid")
			return ErrDisconnectInvalidUserAgent
		}
	}

	return nil
}

// PingMessage Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
	c *gnet.MessageContext `enc:"-"`
}

// Handle implements the Messager interface
func (ping *PingMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	ping.c = mc
	return daemon.(daemoner).recordMessageEvent(ping, mc)
}

// process Sends a PongMessage to the sender of PingMessage
func (ping *PingMessage) process(d daemoner) {
	fields := logrus.Fields{
		"addr":   ping.c.Addr,
		"gnetID": ping.c.ConnID,
	}

	if d.daemonConfig().LogPings {
		logger.WithFields(fields).Debug("Replying to ping")
	}
	if err := d.sendMessage(ping.c.Addr, &PongMessage{}); err != nil {
		logger.WithFields(fields).WithError(err).Error("Send PongMessage failed")
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
		logger.WithFields(logrus.Fields{
			"addr":   mc.Addr,
			"gnetID": mc.ConnID,
		}).Debug("Received pong")
	}
	return nil
}

// DisconnectMessage sent to a peer before disconnecting, indicating the reason for disconnect
type DisconnectMessage struct {
	c      *gnet.MessageContext  `enc:"-"`
	reason gnet.DisconnectReason `enc:"-"`

	// Error code
	ReasonCode uint16

	// Reserved for future use
	Reserved []byte
}

// NewDisconnectMessage creates message sent to reject previously received message
func NewDisconnectMessage(reason gnet.DisconnectReason) *DisconnectMessage {
	return &DisconnectMessage{
		reason:     reason,
		ReasonCode: DisconnectReasonToCode(reason),
		Reserved:   nil,
	}
}

// Handle an event queued by Handle()
func (dm *DisconnectMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	dm.c = mc
	return daemon.(daemoner).recordMessageEvent(dm, mc)
}

// process disconnect message by reflexively disconnecting
func (dm *DisconnectMessage) process(d daemoner) {
	logger.WithFields(logrus.Fields{
		"addr":   dm.c.Addr,
		"gnetID": dm.c.ConnID,
		"code":   dm.ReasonCode,
		"reason": DisconnectCodeToReason(dm.ReasonCode),
	}).Infof("DisconnectMessage received")

	if err := d.disconnectNow(dm.c.Addr, ErrDisconnectReceivedDisconnect); err != nil {
		logger.WithError(err).WithField("addr", dm.c.Addr).Warning("disconnectNow")
	}
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
func (gbm *GetBlocksMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gbm.c = mc
	return daemon.(daemoner).recordMessageEvent(gbm, mc)
}

// process should send number to be requested, with request
func (gbm *GetBlocksMessage) process(d daemoner) {
	if d.daemonConfig().DisableNetworking {
		return
	}

	fields := logrus.Fields{
		"addr":   gbm.c.Addr,
		"gnetID": gbm.c.ConnID,
	}

	// Record this as this peer's highest block
	d.recordPeerHeight(gbm.c.Addr, gbm.c.ConnID, gbm.LastBlock)
	// Fetch and return signed blocks since LastBlock
	blocks, err := d.getSignedBlocksSince(gbm.LastBlock, gbm.RequestedBlocks)
	if err != nil {
		logger.WithError(err).Error("Get signed blocks failed")
		return
	}

	if len(blocks) == 0 {
		return
	}

	logger.WithFields(fields).Debugf("Got %d blocks since %d", len(blocks), gbm.LastBlock)

	m := NewGiveBlocksMessage(blocks)
	if err := d.sendMessage(gbm.c.Addr, m); err != nil {
		logger.WithFields(fields).WithError(err).Error("Send GiveBlocksMessage failed")
	}
}

// GiveBlocksMessage sent in response to GetBlocksMessage, or unsolicited
type GiveBlocksMessage struct {
	Blocks []coin.SignedBlock   `enc:",maxlen=128"`
	c      *gnet.MessageContext `enc:"-"`
}

// NewGiveBlocksMessage creates GiveBlocksMessage
func NewGiveBlocksMessage(blocks []coin.SignedBlock) *GiveBlocksMessage {
	return &GiveBlocksMessage{
		Blocks: blocks,
	}
}

// Handle handle message
func (m *GiveBlocksMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	m.c = mc
	return daemon.(daemoner).recordMessageEvent(m, mc)
}

// process process message
func (m *GiveBlocksMessage) process(d daemoner) {
	if d.daemonConfig().DisableNetworking {
		logger.Critical().Info("Visor disabled, ignoring GiveBlocksMessage")
		return
	}

	// These DB queries are not performed in a transaction for performance reasons.
	// It is not necessary that the blocks be executed together in a single transaction.

	processed := 0
	maxSeq, ok, err := d.headBkSeq()
	if err != nil {
		logger.WithError(err).Error("d.headBkSeq failed")
		return
	}
	if !ok {
		logger.Error("No HeadBkSeq found, cannot execute blocks")
		return
	}

	for _, b := range m.Blocks {
		// To minimize waste when receiving multiple responses from peers
		// we only break out of the loop if the block itself is invalid.
		// E.g. if we request 20 blocks since 0 from 2 peers, and one peer
		// replies with 15 and the other 20, if we did not do this check and
		// the reply with 15 was received first, we would toss the one with 20
		// even though we could process it at the time.
		if b.Seq() <= maxSeq {
			continue
		}

		err := d.executeSignedBlock(b)
		if err == nil {
			logger.Critical().WithField("seq", b.Block.Head.BkSeq).Info("Added new block")
			processed++
		} else {
			logger.Critical().WithError(err).WithField("seq", b.Block.Head.BkSeq).Error("Failed to execute received block")
			// Blocks must be received in order, so if one fails its assumed
			// the rest are failing
			break
		}
	}
	if processed == 0 {
		return
	}

	headBkSeq, ok, err := d.headBkSeq()
	if err != nil {
		logger.WithError(err).Error("d.headBkSeq failed")
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
	abm := NewAnnounceBlocksMessage(headBkSeq)
	if _, err := d.broadcastMessage(abm); err != nil {
		logger.WithError(err).Warning("Broadcast AnnounceBlocksMessage failed")
	}

	// Request more blocks
	gbm := NewGetBlocksMessage(headBkSeq, d.daemonConfig().BlocksResponseCount)
	if _, err := d.broadcastMessage(gbm); err != nil {
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
	return daemon.(daemoner).recordMessageEvent(abm, mc)
}

// process process message
func (abm *AnnounceBlocksMessage) process(d daemoner) {
	if d.daemonConfig().DisableNetworking {
		return
	}

	fields := logrus.Fields{
		"addr":   abm.c.Addr,
		"gnetID": abm.c.ConnID,
	}

	headBkSeq, ok, err := d.headBkSeq()
	if err != nil {
		logger.WithError(err).Error("AnnounceBlocksMessage d.headBkSeq failed")
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
	m := NewGetBlocksMessage(headBkSeq, d.daemonConfig().BlocksResponseCount)
	if err := d.sendMessage(abm.c.Addr, m); err != nil {
		logger.WithError(err).WithFields(fields).Error("Send GetBlocksMessage")
	}
}

// SendingTxnsMessage send transaction message interface
type SendingTxnsMessage interface {
	GetFiltered() []cipher.SHA256
}

// AnnounceTxnsMessage tells a peer that we have these transactions
type AnnounceTxnsMessage struct {
	Transactions []cipher.SHA256      `enc:",maxlen=256"`
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
	return daemon.(daemoner).recordMessageEvent(atm, mc)
}

// process process message
func (atm *AnnounceTxnsMessage) process(d daemoner) {
	if d.daemonConfig().DisableNetworking {
		return
	}

	fields := logrus.Fields{
		"addr":   atm.c.Addr,
		"gnetID": atm.c.ConnID,
	}

	unknown, err := d.filterKnownUnconfirmed(atm.Transactions)
	if err != nil {
		logger.WithError(err).Error("AnnounceTxnsMessage d.filterKnownUnconfirmed failed")
		return
	}

	if len(unknown) == 0 {
		return
	}

	m := NewGetTxnsMessage(unknown)
	if err := d.sendMessage(atm.c.Addr, m); err != nil {
		logger.WithFields(fields).WithError(err).Error("Send GetTxnsMessage failed")
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
	return daemon.(daemoner).recordMessageEvent(gtm, mc)
}

// process process message
func (gtm *GetTxnsMessage) process(d daemoner) {
	if d.daemonConfig().DisableNetworking {
		return
	}

	fields := logrus.Fields{
		"addr":   gtm.c.Addr,
		"gnetID": gtm.c.ConnID,
	}

	// Locate all txns from the unconfirmed pool
	known, err := d.getKnownUnconfirmed(gtm.Transactions)
	if err != nil {
		logger.WithError(err).Error("GetTxnsMessage d.getKnownUnconfirmed failed")
		return
	}
	if len(known) == 0 {
		return
	}

	// Reply to sender with GiveTxnsMessage
	m := NewGiveTxnsMessage(known)
	if err := d.sendMessage(gtm.c.Addr, m); err != nil {
		logger.WithError(err).WithFields(fields).Error("Send GiveTxnsMessage")
	}
}

// GiveTxnsMessage tells the transaction of given hashes
type GiveTxnsMessage struct {
	Transactions []coin.Transaction   `enc:",maxlen=256"`
	c            *gnet.MessageContext `enc:"-"`
}

// NewGiveTxnsMessage creates GiveTxnsMessage
func NewGiveTxnsMessage(txns []coin.Transaction) *GiveTxnsMessage {
	return &GiveTxnsMessage{
		Transactions: txns,
	}
}

// GetFiltered returns transactions hashes
func (gtm *GiveTxnsMessage) GetFiltered() []cipher.SHA256 {
	return coin.Transactions(gtm.Transactions).Hashes()
}

// Handle handle message
func (gtm *GiveTxnsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gtm.c = mc
	return daemon.(daemoner).recordMessageEvent(gtm, mc)
}

// process process message
func (gtm *GiveTxnsMessage) process(d daemoner) {
	if d.daemonConfig().DisableNetworking {
		return
	}

	hashes := make([]cipher.SHA256, 0, len(gtm.Transactions))
	// Update unconfirmed pool with these transactions
	for _, txn := range gtm.Transactions {
		// Only announce transactions that are new to us, so that peers can't spam relays
		// It is not necessary to inject all of the transactions inside a database transaction,
		// since each is independent
		known, softErr, err := d.injectTransaction(txn)
		if err != nil {
			logger.WithError(err).WithField("txid", txn.Hash().Hex()).Warning("Failed to record transaction")
			continue
		} else if softErr != nil {
			logger.WithError(err).WithField("txid", txn.Hash().Hex()).Warning("Transaction soft violation")
			// Allow soft txn violations to rebroadcast
		} else if known {
			logger.WithField("txid", txn.Hash().Hex()).Debug("Duplicate transaction")
			continue
		}

		hashes = append(hashes, txn.Hash())
	}

	if len(hashes) == 0 {
		return
	}

	// Announce these transactions to peers
	m := NewAnnounceTxnsMessage(hashes)
	if ids, err := d.broadcastMessage(m); err != nil {
		logger.WithError(err).Warning("Broadcast AnnounceTxnsMessage failed")
	} else {
		logger.Debugf("Announced %d transactions to %d peers", len(hashes), len(ids))
	}
}
