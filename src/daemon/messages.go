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

//go:generate skyencoder -unexported -struct IntroductionMessage
//go:generate skyencoder -unexported -struct GivePeersMessage
//go:generate skyencoder -unexported -struct GetBlocksMessage
//go:generate skyencoder -unexported -struct GiveBlocksMessage
//go:generate skyencoder -unexported -struct AnnounceBlocksMessage
//go:generate skyencoder -unexported -struct GetTxnsMessage
//go:generate skyencoder -unexported -struct GiveTxnsMessage
//go:generate skyencoder -unexported -struct AnnounceTxnsMessage
//go:generate skyencoder -unexported -struct DisconnectMessage
//go:generate skyencoder -unexported -struct IPAddr
//go:generate skyencoder -unexported -output-path . -package daemon -struct SignedBlock github.com/skycoin/skycoin/src/coin
//go:generate skyencoder -unexported -output-path . -package daemon -struct Transaction github.com/skycoin/skycoin/src/coin

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

// EncodeSize implements gnet.Serializer
func (gpm *GetPeersMessage) EncodeSize() uint64 {
	return 0
}

// Encode implements gnet.Serializer
func (gpm *GetPeersMessage) Encode(buf []byte) error {
	return nil
}

// Decode implements gnet.Serializer
func (gpm *GetPeersMessage) Decode(buf []byte) (uint64, error) {
	return 0, nil
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
	Peers []IPAddr             `enc:",maxlen=512"`
	c     *gnet.MessageContext `enc:"-"`
}

// NewGivePeersMessage []*pex.Peer is converted to []IPAddr for binary transmission
// If the size of the message would exceed maxMsgLength, the IPAddr slice is truncated.
func NewGivePeersMessage(peers []pex.Peer, maxMsgLength uint64) *GivePeersMessage {
	if len(peers) > 512 {
		peers = peers[:512]
	}

	ipaddrs := make([]IPAddr, 0, len(peers))
	for _, ps := range peers {
		ipaddr, err := NewIPAddr(ps.Addr)
		if err != nil {
			logger.WithError(err).WithField("addr", ps.Addr).Warning("GivePeersMessage skipping invalid address")
			continue
		}
		ipaddrs = append(ipaddrs, ipaddr)
	}

	m := &GivePeersMessage{
		Peers: ipaddrs,
	}
	truncateGivePeersMessage(m, maxMsgLength)
	return m
}

// truncateGivePeersMessage truncates the blocks in GivePeersMessage to fit inside of MaxOutgoingMessageLength
func truncateGivePeersMessage(m *GivePeersMessage, maxMsgLength uint64) {
	// The message length will include a 4 byte message type prefix.
	// Panic if the prefix can't fit, otherwise we can't adjust the uint64 safely
	if maxMsgLength < 4 {
		logger.Panic("maxMsgLength must be >= 4")
	}

	maxMsgLength -= 4

	// Measure the current message size, if it fits, return
	n := m.EncodeSize()
	if n <= maxMsgLength {
		return
	}

	// Measure the size of an empty message
	var mm GivePeersMessage
	size := mm.EncodeSize()

	// Measure the size of the peers, advancing the slice index until it reaches capacity
	index := -1
	for i, ip := range m.Peers {
		x := encodeSizeIPAddr(&ip)
		if size+x > maxMsgLength {
			break
		}
		size += x
		index = i
	}

	m.Peers = m.Peers[:index+1]

	if len(m.Peers) == 0 {
		logger.Critical().Error("truncateGivePeersMessage truncated peers to an empty slice")
	}
}

// EncodeSize implements gnet.Serializer
func (gpm *GivePeersMessage) EncodeSize() uint64 {
	return encodeSizeGivePeersMessage(gpm)
}

// Encode implements gnet.Serializer
func (gpm *GivePeersMessage) Encode(buf []byte) error {
	return encodeGivePeersMessageToBuffer(buf, gpm)
}

// Decode implements gnet.Serializer
func (gpm *GivePeersMessage) Decode(buf []byte) (uint64, error) {
	return decodeGivePeersMessage(buf, gpm)
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
	UserAgent            useragent.Data       `enc:"-"`
	UnconfirmedVerifyTxn params.VerifyTxn     `enc:"-"`
	GenesisHash          cipher.SHA256        `enc:"-"`

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
	// GenesisHash         cipher.SHA256 // genesis block hash
	Extra []byte `enc:",omitempty"`
}

// NewIntroductionMessage creates introduction message
func NewIntroductionMessage(mirror uint32, version int32, port uint16, pubkey cipher.PubKey, userAgent string, verifyParams params.VerifyTxn, genesisHash cipher.SHA256) *IntroductionMessage {
	return &IntroductionMessage{
		Mirror:          mirror,
		ProtocolVersion: version,
		ListenPort:      port,
		Extra:           newIntroductionMessageExtra(pubkey, userAgent, verifyParams, genesisHash),
	}
}

func newIntroductionMessageExtra(pubkey cipher.PubKey, userAgent string, verifyParams params.VerifyTxn, genesisHash cipher.SHA256) []byte {
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

	extra := make([]byte, len(pubkey)+len(userAgentSerialized)+len(verifyParamsSerialized)+len(genesisHash))

	copy(extra[:len(pubkey)], pubkey[:])
	i := len(pubkey)
	copy(extra[i:], verifyParamsSerialized)
	i += len(verifyParamsSerialized)
	copy(extra[i:], userAgentSerialized)
	i += len(userAgentSerialized)
	copy(extra[i:i+len(genesisHash)], genesisHash[:])

	return extra
}

// EncodeSize implements gnet.Serializer
func (intro *IntroductionMessage) EncodeSize() uint64 {
	return encodeSizeIntroductionMessage(intro)
}

// Encode implements gnet.Serializer
func (intro *IntroductionMessage) Encode(buf []byte) error {
	return encodeIntroductionMessageToBuffer(buf, intro)
}

// Decode implements gnet.Serializer
func (intro *IntroductionMessage) Decode(buf []byte) (uint64, error) {
	return decodeIntroductionMessage(buf, intro)
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

	if err := intro.Verify(d.DaemonConfig(), logrus.Fields{
		"addr":   addr,
		"gnetID": intro.c.ConnID,
	}); err != nil {
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

// Verify checks if the introduction message is valid returning the appropriate error
func (intro *IntroductionMessage) Verify(dc DaemonConfig, logFields logrus.Fields) error {
	// Disconnect if this is a self connection (we have the same mirror value)
	if intro.Mirror == dc.Mirror {
		logger.WithFields(logFields).WithField("mirror", intro.Mirror).Info("Remote mirror value matches ours")
		return ErrDisconnectSelf
	}

	// Disconnect if peer version is not within the supported range
	if intro.ProtocolVersion < dc.MinProtocolVersion {
		logger.WithFields(logFields).WithFields(logrus.Fields{
			"protocolVersion":    intro.ProtocolVersion,
			"minProtocolVersion": dc.MinProtocolVersion,
		}).Info("protocol version below minimum supported protocol version")
		return ErrDisconnectVersionNotSupported
	}

	logger.WithFields(logFields).WithField("protocolVersion", intro.ProtocolVersion).Debug("Peer protocol version accepted")

	// v24 does not send blockchain pubkey or user agent
	// v25 sends blockchain pubkey and user agent
	// v24 and v25 check the blockchain pubkey and user agent, would accept message with no Pubkey and user agent
	// v26 would check the blockchain pubkey and reject if not matched or not provided, and parses a user agent
	// v26 adds genesis hash
	// v27 would require and check the genesis hash
	extraLen := len(intro.Extra)
	if extraLen == 0 {
		logger.WithFields(logFields).Warning("Blockchain pubkey is not provided")
		return ErrDisconnectBlockchainPubkeyNotProvided
	}

	var bcPubKey cipher.PubKey
	if extraLen < len(bcPubKey) {
		logger.WithFields(logFields).Warning("Extra data length does not meet the minimum requirement")
		return ErrDisconnectInvalidExtraData
	}
	copy(bcPubKey[:], intro.Extra[:len(bcPubKey)])

	if dc.BlockchainPubkey != bcPubKey {
		logger.WithFields(logFields).WithFields(logrus.Fields{
			"pubkey":       bcPubKey.Hex(),
			"daemonPubkey": dc.BlockchainPubkey.Hex(),
		}).Warning("Blockchain pubkey does not match")
		return ErrDisconnectBlockchainPubkeyNotMatched
	}

	i := len(bcPubKey)
	if extraLen < i+9 {
		logger.WithFields(logFields).Warning("IntroductionMessage transaction verification parameters could not be deserialized: not enough data")
		return ErrDisconnectInvalidExtraData
	}
	if err := encoder.DeserializeRawExact(intro.Extra[i:i+9], &intro.UnconfirmedVerifyTxn); err != nil {
		// This should not occur due to the previous length check
		logger.Critical().WithError(err).WithFields(logFields).Warning("unconfirmedVerifyTxn params could not be deserialized")
		return ErrDisconnectInvalidExtraData
	}
	i += 9

	if err := intro.UnconfirmedVerifyTxn.Validate(); err != nil {
		logger.WithError(err).WithFields(logFields).WithFields(logrus.Fields{
			"burnFactor":          intro.UnconfirmedVerifyTxn.BurnFactor,
			"maxTransactionSize":  intro.UnconfirmedVerifyTxn.MaxTransactionSize,
			"maxDropletPrecision": intro.UnconfirmedVerifyTxn.MaxDropletPrecision,
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

	userAgentSerialized := intro.Extra[i:]
	userAgent, userAgentLen, err := encoder.DeserializeString(userAgentSerialized, useragent.MaxLen)
	if err != nil {
		logger.WithError(err).WithFields(logFields).Warning("Extra data user agent string could not be deserialized")
		return ErrDisconnectInvalidExtraData
	}

	intro.UserAgent, err = useragent.Parse(useragent.Sanitize(userAgent))
	if err != nil {
		logger.WithError(err).WithFields(logFields).WithField("userAgent", userAgent).Warning("User agent is invalid")
		return ErrDisconnectInvalidUserAgent
	}
	i += int(userAgentLen)

	remainingLen := extraLen - i
	if remainingLen > 0 && remainingLen < len(intro.GenesisHash) {
		logger.WithFields(logFields).Warning("Extra data genesis hash could not be deserialized: not enough data")
		return ErrDisconnectInvalidExtraData
	}
	copy(intro.GenesisHash[:], intro.Extra[i:])

	return nil
}

// PingMessage Sent to keep a connection alive. A PongMessage is sent in reply.
type PingMessage struct {
	c *gnet.MessageContext `enc:"-"`
}

// EncodeSize implements gnet.Serializer
func (ping *PingMessage) EncodeSize() uint64 {
	return 0
}

// Encode implements gnet.Serializer
func (ping *PingMessage) Encode(buf []byte) error {
	return nil
}

// Decode implements gnet.Serializer
func (ping *PingMessage) Decode(buf []byte) (uint64, error) {
	return 0, nil
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

	if d.DaemonConfig().LogPings {
		logger.WithFields(fields).Debug("Replying to ping")
	}
	if err := d.sendMessage(ping.c.Addr, &PongMessage{}); err != nil {
		logger.WithFields(fields).WithError(err).Error("Send PongMessage failed")
	}
}

// PongMessage Sent in reply to a PingMessage.  No action is taken when this is received.
type PongMessage struct {
}

// EncodeSize implements gnet.Serializer
func (pong *PongMessage) EncodeSize() uint64 {
	return 0
}

// Encode implements gnet.Serializer
func (pong *PongMessage) Encode(buf []byte) error {
	return nil
}

// Decode implements gnet.Serializer
func (pong *PongMessage) Decode(buf []byte) (uint64, error) {
	return 0, nil
}

// Handle handles message
func (pong *PongMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	// There is nothing to do; gnet updates Connection.LastMessage internally
	// when this is received
	if daemon.(daemoner).DaemonConfig().LogPings {
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

// EncodeSize implements gnet.Serializer
func (dm *DisconnectMessage) EncodeSize() uint64 {
	return encodeSizeDisconnectMessage(dm)
}

// Encode implements gnet.Serializer
func (dm *DisconnectMessage) Encode(buf []byte) error {
	return encodeDisconnectMessageToBuffer(buf, dm)
}

// Decode implements gnet.Serializer
func (dm *DisconnectMessage) Decode(buf []byte) (uint64, error) {
	return decodeDisconnectMessage(buf, dm)
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
func NewGetBlocksMessage(lastBlock, requestedBlocks uint64) *GetBlocksMessage {
	return &GetBlocksMessage{
		LastBlock:       lastBlock,
		RequestedBlocks: requestedBlocks,
	}
}

// EncodeSize implements gnet.Serializer
func (gbm *GetBlocksMessage) EncodeSize() uint64 {
	return encodeSizeGetBlocksMessage(gbm)
}

// Encode implements gnet.Serializer
func (gbm *GetBlocksMessage) Encode(buf []byte) error {
	return encodeGetBlocksMessageToBuffer(buf, gbm)
}

// Decode implements gnet.Serializer
func (gbm *GetBlocksMessage) Decode(buf []byte) (uint64, error) {
	return decodeGetBlocksMessage(buf, gbm)
}

// Handle handles message
func (gbm *GetBlocksMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gbm.c = mc
	return daemon.(daemoner).recordMessageEvent(gbm, mc)
}

// process should send number to be requested, with request
func (gbm *GetBlocksMessage) process(d daemoner) {
	dc := d.DaemonConfig()
	if dc.DisableNetworking {
		return
	}

	fields := logrus.Fields{
		"addr":   gbm.c.Addr,
		"gnetID": gbm.c.ConnID,
	}

	// Record this as this peer's highest block
	d.recordPeerHeight(gbm.c.Addr, gbm.c.ConnID, gbm.LastBlock)

	// Cap the number of requested blocks (TODO - necessary since we have size limits enforced later?)
	requestedBlocks := gbm.RequestedBlocks
	if requestedBlocks > dc.MaxGetBlocksResponseCount {
		logger.WithFields(logrus.Fields{
			"requestedBlocks":    requestedBlocks,
			"maxRequestedBlocks": dc.MaxGetBlocksResponseCount,
		}).WithFields(fields).Debug("GetBlocksMessage.RequestedBlocks value exceeds configured limit, reducing")
		requestedBlocks = dc.MaxGetBlocksResponseCount
	}

	// Fetch and return signed blocks since LastBlock
	blocks, err := d.getSignedBlocksSince(gbm.LastBlock, requestedBlocks)
	if err != nil {
		logger.WithFields(fields).WithError(err).Error("getSignedBlocksSince failed")
		return
	}

	if len(blocks) == 0 {
		return
	}

	logger.WithFields(fields).Debugf("GetBlocksMessage: replying with %d blocks after block %d", len(blocks), gbm.LastBlock)

	m := NewGiveBlocksMessage(blocks, dc.MaxOutgoingMessageLength)
	if len(m.Blocks) != len(blocks) {
		logger.WithField("startBlockSeq", blocks[0].Head.BkSeq).WithFields(fields).Warningf("NewGiveBlocksMessage truncated %d blocks to %d blocks", len(blocks), len(m.Blocks))
	}

	if err := d.sendMessage(gbm.c.Addr, m); err != nil {
		logger.WithFields(fields).WithError(err).Error("Send GiveBlocksMessage failed")
	}
}

// GiveBlocksMessage sent in response to GetBlocksMessage, or unsolicited
type GiveBlocksMessage struct {
	Blocks []coin.SignedBlock   `enc:",maxlen=128"`
	c      *gnet.MessageContext `enc:"-"`
}

// NewGiveBlocksMessage creates GiveBlocksMessage.
// If the size of message would exceed maxMsgLength, the block slice is truncated.
func NewGiveBlocksMessage(blocks []coin.SignedBlock, maxMsgLength uint64) *GiveBlocksMessage {
	if len(blocks) > 128 {
		blocks = blocks[:128]
	}
	m := &GiveBlocksMessage{
		Blocks: blocks,
	}
	truncateGiveBlocksMessage(m, maxMsgLength)
	return m
}

// truncateGiveBlocksMessage truncates the blocks in GiveBlocksMessage to fit inside of MaxOutgoingMessageLength
func truncateGiveBlocksMessage(m *GiveBlocksMessage, maxMsgLength uint64) {
	// The message length will include a 4 byte message type prefix.
	// Panic if the prefix can't fit, otherwise we can't adjust the uint64 safely
	if maxMsgLength < 4 {
		logger.Panic("maxMsgLength must be >= 4")
	}

	maxMsgLength -= 4

	// Measure the current message size, if it fits, return
	n := m.EncodeSize()
	if n <= maxMsgLength {
		return
	}

	// Measure the size of an empty message
	var mm GiveBlocksMessage
	size := mm.EncodeSize()

	// Measure the size of the blocks, advancing the slice index until it reaches capacity
	index := -1
	for i, b := range m.Blocks {
		x := encodeSizeSignedBlock(&b)
		if size+x > maxMsgLength {
			break
		}
		size += x
		index = i
	}

	m.Blocks = m.Blocks[:index+1]

	if len(m.Blocks) == 0 {
		logger.Critical().Error("truncateGiveBlocksMessage truncated blocks to an empty slice")
	}
}

// EncodeSize implements gnet.Serializer
func (m *GiveBlocksMessage) EncodeSize() uint64 {
	return encodeSizeGiveBlocksMessage(m)
}

// Encode implements gnet.Serializer
func (m *GiveBlocksMessage) Encode(buf []byte) error {
	return encodeGiveBlocksMessageToBuffer(buf, m)
}

// Decode implements gnet.Serializer
func (m *GiveBlocksMessage) Decode(buf []byte) (uint64, error) {
	return decodeGiveBlocksMessage(buf, m)
}

// Handle handle message
func (m *GiveBlocksMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	m.c = mc
	return daemon.(daemoner).recordMessageEvent(m, mc)
}

// process process message
func (m *GiveBlocksMessage) process(d daemoner) {
	if d.DaemonConfig().DisableNetworking {
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
	gbm := NewGetBlocksMessage(headBkSeq, d.DaemonConfig().GetBlocksRequestCount)
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

// EncodeSize implements gnet.Serializer
func (abm *AnnounceBlocksMessage) EncodeSize() uint64 {
	return encodeSizeAnnounceBlocksMessage(abm)
}

// Encode implements gnet.Serializer
func (abm *AnnounceBlocksMessage) Encode(buf []byte) error {
	return encodeAnnounceBlocksMessageToBuffer(buf, abm)
}

// Decode implements gnet.Serializer
func (abm *AnnounceBlocksMessage) Decode(buf []byte) (uint64, error) {
	return decodeAnnounceBlocksMessage(buf, abm)
}

// Handle handles message
func (abm *AnnounceBlocksMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	abm.c = mc
	return daemon.(daemoner).recordMessageEvent(abm, mc)
}

// process process message
func (abm *AnnounceBlocksMessage) process(d daemoner) {
	if d.DaemonConfig().DisableNetworking {
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
	m := NewGetBlocksMessage(headBkSeq, d.DaemonConfig().GetBlocksRequestCount)
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

// NewAnnounceTxnsMessage creates announce txns message.
// If the size of the message would exceed maxMsgLength, the hashes slice is truncated.
func NewAnnounceTxnsMessage(txns []cipher.SHA256, maxMsgLength uint64) *AnnounceTxnsMessage {
	if len(txns) > 256 {
		txns = txns[:256]
	}
	m := &AnnounceTxnsMessage{
		Transactions: txns,
	}
	hashes := truncateAnnounceTxnsHashes(m, maxMsgLength)
	m.Transactions = hashes
	return m
}

// truncateAnnounceTxnsHashes truncates the hashes in AnnounceTxnsMessage to fit inside of MaxOutgoingMessageLength
func truncateAnnounceTxnsHashes(m *AnnounceTxnsMessage, maxMsgLength uint64) []cipher.SHA256 {
	// The message length will include a 4 byte message type prefix.
	// Panic if the prefix can't fit, otherwise we can't adjust the uint64 safely
	if maxMsgLength < 4 {
		logger.Panic("maxMsgLength must be >= 4")
	}

	maxMsgLength -= 4

	// Measure the current message size, if it fits, return
	n := m.EncodeSize()
	if n <= maxMsgLength {
		return m.Transactions
	}

	// Measure the size of an empty message
	var mm AnnounceTxnsMessage
	size := mm.EncodeSize()

	if maxMsgLength < size {
		logger.Panic("maxMsgLength must be <= 4 + sizeof(empty AnnounceTxnsMessage)")
	}

	maxMsgLength -= size

	hashes := truncateSHA256Slice(m.Transactions, maxMsgLength)

	if len(hashes) == 0 {
		logger.Critical().Error("truncateAnnounceTxnsHashes truncated hashes to an empty slice")
	}

	return hashes
}

func truncateSHA256Slice(hashes []cipher.SHA256, maxLength uint64) []cipher.SHA256 {
	if len(hashes) == 0 {
		return hashes
	}

	size := len(hashes[0])

	n := maxLength / uint64(size)

	if n > uint64(len(hashes)) {
		return hashes
	}

	return hashes[:n]
}

// EncodeSize implements gnet.Serializer
func (atm *AnnounceTxnsMessage) EncodeSize() uint64 {
	return encodeSizeAnnounceTxnsMessage(atm)
}

// Encode implements gnet.Serializer
func (atm *AnnounceTxnsMessage) Encode(buf []byte) error {
	return encodeAnnounceTxnsMessageToBuffer(buf, atm)
}

// Decode implements gnet.Serializer
func (atm *AnnounceTxnsMessage) Decode(buf []byte) (uint64, error) {
	return decodeAnnounceTxnsMessage(buf, atm)
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
	dc := d.DaemonConfig()
	if dc.DisableNetworking {
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

	m := NewGetTxnsMessage(unknown, dc.MaxOutgoingMessageLength)
	if len(m.Transactions) != len(unknown) {
		logger.Warningf("NewGetTxnsMessage truncated %d hashes to %d hashes", len(unknown), len(m.Transactions))
	}

	if err := d.sendMessage(atm.c.Addr, m); err != nil {
		logger.WithFields(fields).WithError(err).Error("Send GetTxnsMessage failed")
	}
}

// GetTxnsMessage request transactions of given hash
type GetTxnsMessage struct {
	Transactions []cipher.SHA256      `enc:",maxlen=256"`
	c            *gnet.MessageContext `enc:"-"`
}

// NewGetTxnsMessage creates GetTxnsMessage.
// If the size of the message would exceed maxMsgLength, the hashes slice is truncated.
func NewGetTxnsMessage(txns []cipher.SHA256, maxMsgLength uint64) *GetTxnsMessage {
	if len(txns) > 256 {
		txns = txns[:256]
	}
	m := &GetTxnsMessage{
		Transactions: txns,
	}
	hashes := truncateGetTxnsHashes(m, maxMsgLength)
	m.Transactions = hashes
	return m
}

// truncateGetTxnsHashes truncates the hashes in GetTxnsMessage to fit inside of MaxOutgoingMessageLength
func truncateGetTxnsHashes(m *GetTxnsMessage, maxMsgLength uint64) []cipher.SHA256 {
	// The message length will include a 4 byte message type prefix.
	// Panic if the prefix can't fit, otherwise we can't adjust the uint64 safely
	if maxMsgLength < 4 {
		logger.Panic("maxMsgLength must be >= 4")
	}

	maxMsgLength -= 4

	// Measure the current message size, if it fits, return
	n := m.EncodeSize()
	if n <= maxMsgLength {
		return m.Transactions
	}

	// Measure the size of an empty message
	var mm GetTxnsMessage
	size := mm.EncodeSize()

	if maxMsgLength < size {
		logger.Panic("maxMsgLength must be <= 4 + sizeof(empty GetTxnsMessage)")
	}

	maxMsgLength -= size

	hashes := truncateSHA256Slice(m.Transactions, maxMsgLength)

	if len(hashes) == 0 {
		logger.Critical().Error("truncateGetTxnsHashes truncated hashes to an empty slice")
	}

	return hashes
}

// EncodeSize implements gnet.Serializer
func (gtm *GetTxnsMessage) EncodeSize() uint64 {
	return encodeSizeGetTxnsMessage(gtm)
}

// Encode implements gnet.Serializer
func (gtm *GetTxnsMessage) Encode(buf []byte) error {
	return encodeGetTxnsMessageToBuffer(buf, gtm)
}

// Decode implements gnet.Serializer
func (gtm *GetTxnsMessage) Decode(buf []byte) (uint64, error) {
	return decodeGetTxnsMessage(buf, gtm)
}

// Handle handle message
func (gtm *GetTxnsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gtm.c = mc
	return daemon.(daemoner).recordMessageEvent(gtm, mc)
}

// process process message
func (gtm *GetTxnsMessage) process(d daemoner) {
	dc := d.DaemonConfig()
	if dc.DisableNetworking {
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
	m := NewGiveTxnsMessage(known, dc.MaxOutgoingMessageLength)
	if len(m.Transactions) != len(known) {
		logger.Warningf("NewGiveTxnsMessage truncated %d hashes to %d hashes", len(known), len(m.Transactions))
	}

	if err := d.sendMessage(gtm.c.Addr, m); err != nil {
		logger.WithError(err).WithFields(fields).Error("Send GiveTxnsMessage")
	}
}

// GiveTxnsMessage tells the transaction of given hashes
type GiveTxnsMessage struct {
	Transactions []coin.Transaction   `enc:",maxlen=256"`
	c            *gnet.MessageContext `enc:"-"`
}

// NewGiveTxnsMessage creates GiveTxnsMessage.
// If the size of the message would exceed maxMsgLength, the transactions slice is truncated.
func NewGiveTxnsMessage(txns []coin.Transaction, maxMsgLength uint64) *GiveTxnsMessage {
	if len(txns) > 256 {
		txns = txns[:256]
	}
	m := &GiveTxnsMessage{
		Transactions: txns,
	}
	truncateGiveTxnsMessage(m, maxMsgLength)
	return m
}

// truncateGiveTxnsMessage truncates the transactions in GiveTxnsMessage to fit inside of MaxOutgoingMessageLength
func truncateGiveTxnsMessage(m *GiveTxnsMessage, maxMsgLength uint64) {
	// The message length will include a 4 byte message type prefix.
	// Panic if the prefix can't fit, otherwise we can't adjust the uint64 safely
	if maxMsgLength < 4 {
		logger.Panic("maxMsgLength must be >= 4")
	}

	maxMsgLength -= 4

	// Measure the current message size, if it fits, return
	n := m.EncodeSize()
	if n <= maxMsgLength {
		return
	}

	// Measure the size of an empty message
	var mm GiveTxnsMessage
	size := mm.EncodeSize()

	// Measure the size of the txns, advancing the slice index until it reaches capacity
	index := -1
	for i, txn := range m.Transactions {
		x := encodeSizeTransaction(&txn)
		if size+x > maxMsgLength {
			break
		}
		size += x
		index = i
	}

	m.Transactions = m.Transactions[:index+1]

	if len(m.Transactions) == 0 {
		logger.Critical().Error("truncateGiveTxnsMessage truncated txns to an empty slice")
	}
}

// EncodeSize implements gnet.Serializer
func (gtm *GiveTxnsMessage) EncodeSize() uint64 {
	return encodeSizeGiveTxnsMessage(gtm)
}

// Encode implements gnet.Serializer
func (gtm *GiveTxnsMessage) Encode(buf []byte) error {
	return encodeGiveTxnsMessageToBuffer(buf, gtm)
}

// Decode implements gnet.Serializer
func (gtm *GiveTxnsMessage) Decode(buf []byte) (uint64, error) {
	return decodeGiveTxnsMessage(buf, gtm)
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
	dc := d.DaemonConfig()
	if dc.DisableNetworking {
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
			logger.WithError(softErr).WithField("txid", txn.Hash().Hex()).Warning("Transaction soft violation")
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
	m := NewAnnounceTxnsMessage(hashes, dc.MaxOutgoingMessageLength)
	if len(m.Transactions) != len(hashes) {
		logger.Warningf("NewAnnounceTxnsMessage truncated %d hashes to %d hashes", len(hashes), len(m.Transactions))
	}

	if ids, err := d.broadcastMessage(m); err != nil {
		logger.WithError(err).Warning("Broadcast AnnounceTxnsMessage failed")
	} else {
		logger.Debugf("Announced %d transactions to %d peers", len(hashes), len(ids))
	}
}
