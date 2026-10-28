package daemon

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/visor"
)

// multiSignerProtocolVersion is the first protocol version that understands
// ProposeBlockMessage and GiveBlockCosigsMessage. Older peers never get one.
const multiSignerProtocolVersion int32 = 5

const (
	// maxBlockCosigsPerMessage keeps a GiveBlockCosigsMessage near 26 KB
	maxBlockCosigsPerMessage = 256
	// blockCosigCacheSize bounds the blocks whose cosignatures are held
	blockCosigCacheSize = 4096
	// maxProposalsPerHeight bounds the proposals held for the next block
	maxProposalsPerHeight = 64
	// maxProposalClockAhead is how far ahead of our clock a proposal's time may be
	// (docs/multi-signer-fiber.md, Rounds and proposers, step 2)
	maxProposalClockAhead = time.Minute
	// cosigRequestInterval limits how often a node asks peers again for blocks
	// it holds without enough signatures
	cosigRequestInterval = 5 * time.Second
	// signerTickInterval is how often a signer checks its round
	signerTickInterval = time.Second
)

// BlockCosigEntry is one signer's signature over the header of the block with Hash.
type BlockCosigEntry struct {
	Hash  cipher.SHA256
	Cosig coin.BlockCosig
}

const blockCosigEntrySize = 32 + 4 + 65

// GiveBlockCosigsMessage carries block cosignatures. A signer broadcasts its own,
// nodes relay new ones, and a peer sends those of blocks right before the
// GiveBlocksMessage holding them.
type GiveBlockCosigsMessage struct {
	Cosigs []BlockCosigEntry
	c      *gnet.MessageContext `enc:"-"`
}

// EncodeSize implements gnet.Serializer
func (m *GiveBlockCosigsMessage) EncodeSize() uint64 {
	return 4 + uint64(len(m.Cosigs))*blockCosigEntrySize
}

// Encode implements gnet.Serializer
func (m *GiveBlockCosigsMessage) Encode(buf []byte) error {
	if uint64(len(buf)) < m.EncodeSize() {
		return encoder.ErrBufferUnderflow
	}
	if len(m.Cosigs) > maxBlockCosigsPerMessage {
		return encoder.ErrMaxLenExceeded
	}
	e := &encoder.Encoder{Buffer: buf}
	e.Uint32(uint32(len(m.Cosigs))) //nolint:gosec // bounded by maxBlockCosigsPerMessage
	for _, c := range m.Cosigs {
		e.CopyBytes(c.Hash[:])
		e.Uint32(c.Cosig.Signer)
		e.CopyBytes(c.Cosig.Sig[:])
	}
	return nil
}

// Decode implements gnet.Serializer
func (m *GiveBlockCosigsMessage) Decode(buf []byte) (uint64, error) {
	d := &encoder.Decoder{Buffer: buf}
	n, err := d.Uint32()
	if err != nil {
		return 0, err
	}
	if n > maxBlockCosigsPerMessage {
		return 0, encoder.ErrMaxLenExceeded
	}
	if uint64(len(d.Buffer)) < uint64(n)*blockCosigEntrySize {
		return 0, encoder.ErrBufferUnderflow
	}
	m.Cosigs = make([]BlockCosigEntry, n)
	for i := range m.Cosigs {
		copy(m.Cosigs[i].Hash[:], d.Buffer[:32])
		d.Buffer = d.Buffer[32:]
		if m.Cosigs[i].Cosig.Signer, err = d.Uint32(); err != nil {
			return 0, err
		}
		copy(m.Cosigs[i].Cosig.Sig[:], d.Buffer[:65])
		d.Buffer = d.Buffer[65:]
	}
	return uint64(len(buf) - len(d.Buffer)), nil //nolint:gosec // d.Buffer is a suffix of buf
}

// Handle handle message
func (m *GiveBlockCosigsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	m.c = mc
	return daemon.(daemoner).recordMessageEvent(m, mc)
}

func (m *GiveBlockCosigsMessage) process(d daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}
	d.onBlockCosigs(m.c.Addr, m.Cosigs)
}

// ProposeBlockMessage carries a signer's proposal for the next block of a
// multi-signer chain, for the round in which the signer is the proposer.
type ProposeBlockMessage struct {
	Round uint64
	Block coin.SignedBlock
	c     *gnet.MessageContext `enc:"-"`
}

// EncodeSize implements gnet.Serializer
func (m *ProposeBlockMessage) EncodeSize() uint64 {
	return 8 + encodeSizeSignedBlock(&m.Block)
}

// Encode implements gnet.Serializer
func (m *ProposeBlockMessage) Encode(buf []byte) error {
	if uint64(len(buf)) < m.EncodeSize() {
		return encoder.ErrBufferUnderflow
	}
	e := &encoder.Encoder{Buffer: buf}
	e.Uint64(m.Round)
	return encodeSignedBlockToBuffer(buf[8:], &m.Block)
}

// Decode implements gnet.Serializer
func (m *ProposeBlockMessage) Decode(buf []byte) (uint64, error) {
	d := &encoder.Decoder{Buffer: buf}
	var err error
	if m.Round, err = d.Uint64(); err != nil {
		return 0, err
	}
	n, err := decodeSignedBlock(buf[8:], &m.Block)
	if err != nil {
		return 0, err
	}
	return 8 + n, nil
}

// Handle handle message
func (m *ProposeBlockMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	m.c = mc
	return daemon.(daemoner).recordMessageEvent(m, mc)
}

func (m *ProposeBlockMessage) process(d daemoner) {
	if d.DaemonConfig().DisableNetworking {
		return
	}
	d.onBlockProposal(m.c.Addr, m)
}

type blockProposal struct {
	block coin.SignedBlock
	round uint64
}

// multiSignerState is the daemon's view of the next block of a multi-signer
// chain: the proposals for it and the cosignatures received for any block.
type multiSignerState struct {
	mu sync.Mutex

	cosigs     map[cipher.SHA256]map[uint32]cipher.Sig
	cosigOrder []cipher.SHA256

	// height is the block the proposals are for, the one after the head
	height    uint64
	proposals map[cipher.SHA256]blockProposal
	// lastRound is the round of the last signer tick, to rebroadcast on a new round
	lastRound     uint64
	lastCosigsAsk time.Time
}

func newMultiSignerState() *multiSignerState {
	return &multiSignerState{
		cosigs:    make(map[cipher.SHA256]map[uint32]cipher.Sig),
		proposals: make(map[cipher.SHA256]blockProposal),
	}
}

// addCosig keeps a cosignature and reports whether it is new.
func (s *multiSignerState) addCosig(hash cipher.SHA256, c coin.BlockCosig) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	set, ok := s.cosigs[hash]
	if !ok {
		for len(s.cosigOrder) >= blockCosigCacheSize {
			delete(s.cosigs, s.cosigOrder[0])
			s.cosigOrder = s.cosigOrder[1:]
		}
		set = make(map[uint32]cipher.Sig)
		s.cosigs[hash] = set
		s.cosigOrder = append(s.cosigOrder, hash)
	}
	if _, ok := set[c.Signer]; ok {
		return false
	}
	set[c.Signer] = c.Sig
	return true
}

func (s *multiSignerState) cosigsFor(hash cipher.SHA256) []coin.BlockCosig {
	s.mu.Lock()
	defer s.mu.Unlock()
	set := s.cosigs[hash]
	if len(set) == 0 {
		return nil
	}
	cs := make([]coin.BlockCosig, 0, len(set))
	for i, sig := range set {
		cs = append(cs, coin.BlockCosig{Signer: i, Sig: sig})
	}
	return cs
}

// setHeight forgets the proposals once the head has moved.
func (s *multiSignerState) setHeight(h uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.height != h {
		s.height = h
		s.proposals = make(map[cipher.SHA256]blockProposal)
	}
}

// addProposal keeps a proposal for the next block and reports whether it is new.
func (s *multiSignerState) addProposal(b coin.SignedBlock, round uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if b.Seq() != s.height {
		return false
	}
	h := b.HashHeader()
	if _, ok := s.proposals[h]; ok || len(s.proposals) >= maxProposalsPerHeight {
		return false
	}
	s.proposals[h] = blockProposal{block: b, round: round}
	return true
}

func (s *multiSignerState) proposal(hash cipher.SHA256) (blockProposal, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.proposals[hash]
	return p, ok
}

// otherVotes returns the hashes of the other proposals at the current height
// that signer i has signed, as proposer or cosigner.
func (s *multiSignerState) otherVotes(rule coin.MultiSignerRule, hash cipher.SHA256, i uint32) []cipher.SHA256 {
	s.mu.Lock()
	defer s.mu.Unlock()
	var hs []cipher.SHA256
	for h, p := range s.proposals {
		if h == hash {
			continue
		}
		_, cosigned := s.cosigs[h][i]
		if cosigned || rule.ProposerIndex(p.block) == int(i) {
			hs = append(hs, h)
		}
	}
	return hs
}

func (dm *Daemon) multiSignerRule() coin.MultiSignerRule {
	return dm.visor.Config.MultiSigner
}

// blockRoundTimeout is how long a round lasts, three block intervals by default
func (dm *Daemon) blockRoundTimeout() time.Duration {
	if dm.config.BlockRoundTimeout > 0 {
		return dm.config.BlockRoundTimeout
	}
	return 3 * time.Second * time.Duration(dm.config.BlockCreationInterval) //nolint:gosec
}

// blockRound returns the round for the block after head at now. Round 0 starts
// at the head block's time, so every node with the head counts the same rounds.
func (dm *Daemon) blockRound(head *coin.SignedBlock, now time.Time) uint64 {
	elapsed := now.Sub(time.Unix(int64(head.Time()), 0)) //nolint:gosec
	if elapsed <= 0 {
		return 0
	}
	return uint64(elapsed / dm.blockRoundTimeout()) //nolint:gosec // elapsed > 0
}

// multiSignerHead returns the head block and moves the proposal state to the block after it.
func (dm *Daemon) multiSignerHead() (*coin.SignedBlock, error) {
	head, err := dm.visor.GetHeadBlock()
	if err != nil {
		return nil, err
	}
	dm.multiSigner.setHeight(head.Seq() + 1)
	return head, nil
}

// signerTick runs a block signer's part of a round: the proposer builds, signs
// and sends a block, and on a new round each signer resends what it signed.
func (dm *Daemon) signerTick() {
	rule := dm.multiSignerRule()
	me := dm.visor.SignerIndex()
	if me < 0 {
		return
	}
	head, err := dm.multiSignerHead()
	if err != nil {
		logger.WithError(err).Error("signerTick: head block")
		return
	}
	h := head.Seq() + 1
	if !rule.Applies(h) {
		return
	}

	now := time.Now()
	round := dm.blockRound(head, now)
	newRound := round != dm.multiSigner.lastRound
	dm.multiSigner.lastRound = round

	voted, hasVote, err := dm.visor.BlockVote(h)
	if err != nil {
		logger.WithError(err).Error("signerTick: read block vote")
		return
	}
	if hasVote {
		if newRound {
			dm.resendVote(h, voted)
		}
		return
	}

	if rule.Proposer(h, round) != me {
		return
	}
	if now.Before(time.Unix(int64(head.Time()), 0).Add(time.Second * time.Duration(dm.config.BlockCreationInterval))) { //nolint:gosec
		return
	}

	sb, err := dm.visor.CreateBlockProposal(now)
	if err != nil {
		logger.WithError(err).Debug("signerTick: no block to propose")
		return
	}
	hash := sb.HashHeader()
	// Recorded before the signature leaves the node (docs/multi-signer-fiber.md, Safety)
	if err := dm.visor.RecordBlockVote(h, hash); err != nil {
		logger.WithError(err).Error("signerTick: record block vote")
		return
	}
	dm.multiSigner.addProposal(sb, round)
	logger.Critical().WithFields(logrus.Fields{
		"seq":   h,
		"round": round,
		"hash":  hash.Hex(),
	}).Info("Proposed block")

	dm.relayProposal(&ProposeBlockMessage{Round: round, Block: sb}, "")
	dm.tryExecuteProposal(hash)
}

// resendVote resends the proposal this signer signed at h and the cosignatures
// it holds for it, so signers that missed them can still sign.
func (dm *Daemon) resendVote(h uint64, hash cipher.SHA256) {
	if p, ok := dm.multiSigner.proposal(hash); ok {
		dm.relayProposal(&ProposeBlockMessage{Round: p.round, Block: p.block}, "")
	}
	var entries []BlockCosigEntry
	for _, c := range dm.multiSigner.cosigsFor(hash) {
		entries = append(entries, BlockCosigEntry{Hash: hash, Cosig: c})
	}
	dm.relayBlockCosigs(entries, "")
	logger.WithField("seq", h).Debug("Resent block vote for a new round")
}

// onBlockProposal checks a proposal, relays it, cosigns it if this node is a
// signer that has not signed at its height, and executes it once final.
func (dm *Daemon) onBlockProposal(addr string, m *ProposeBlockMessage) {
	rule := dm.multiSignerRule()
	if !rule.Enabled() {
		return
	}
	b := m.Block
	fields := logrus.Fields{"addr": addr, "seq": b.Seq(), "round": m.Round}

	head, err := dm.multiSignerHead()
	if err != nil {
		logger.WithError(err).Error("onBlockProposal: head block")
		return
	}
	switch {
	case b.Seq() <= head.Seq():
		return
	case b.Seq() > head.Seq()+1:
		if err := dm.requestBlocksFromAddr(addr); err != nil {
			logger.WithFields(fields).WithError(err).Debug("onBlockProposal: request blocks")
		}
		return
	case !rule.Applies(b.Seq()):
		return
	}

	p := rule.ProposerIndex(b)
	if p < 0 || p != rule.Proposer(b.Seq(), m.Round) {
		logger.WithFields(fields).Info("Ignoring a proposal not signed by the proposer of its round")
		return
	}
	now := time.Now()
	if m.Round > dm.blockRound(head, now)+1 {
		logger.WithFields(fields).Info("Ignoring a proposal for a future round")
		return
	}

	for j := range b.Body.Transactions {
		dm.attachPQWitness(&b.Body.Transactions[j])
	}
	if !dm.multiSigner.addProposal(b, m.Round) {
		return
	}
	hash := b.HashHeader()
	for _, other := range dm.multiSigner.otherVotes(rule, hash, uint32(p)) { //nolint:gosec
		logger.Critical().WithFields(fields).Warningf("Signer %s signed two blocks at height %d: %s and %s", rule.Signers[p].Hex(), b.Seq(), hash.Hex(), other.Hex())
	}
	dm.relayProposal(&ProposeBlockMessage{Round: m.Round, Block: b}, addr)

	if dm.visor.SignerIndex() >= 0 {
		dm.maybeCosign(b, head, now)
	}
	dm.tryExecuteProposal(hash)
}

// maybeCosign signs a proposal if it is valid and this signer has not signed
// another block at its height (docs/multi-signer-fiber.md, steps 2 and 3).
func (dm *Daemon) maybeCosign(b coin.SignedBlock, head *coin.SignedBlock, now time.Time) {
	fields := logrus.Fields{"seq": b.Seq(), "hash": b.HashHeader().Hex()}
	if _, voted, err := dm.visor.BlockVote(b.Seq()); err != nil || voted {
		return
	}
	if b.Time() <= head.Time() {
		logger.WithFields(fields).Info("Not signing a proposal whose time is not after the head")
		return
	}
	if time.Unix(int64(b.Time()), 0).After(now.Add(maxProposalClockAhead)) { //nolint:gosec
		logger.WithFields(fields).Info("Not signing a proposal whose time is ahead of our clock")
		return
	}
	if err := dm.visor.VerifyBlock(b); err != nil {
		logger.WithFields(fields).WithError(err).Info("Not signing an invalid proposal")
		return
	}
	hash := b.HashHeader()
	if err := dm.visor.RecordBlockVote(b.Seq(), hash); err != nil {
		var already visor.ErrAlreadyVoted
		if !errors.As(err, &already) {
			logger.WithFields(fields).WithError(err).Error("maybeCosign: record block vote")
		}
		return
	}
	c, err := dm.visor.SignBlockCosig(b.Block)
	if err != nil {
		logger.WithFields(fields).WithError(err).Error("maybeCosign: sign")
		return
	}
	dm.multiSigner.addCosig(hash, c)
	logger.Critical().WithFields(fields).Info("Signed block proposal")
	dm.relayBlockCosigs([]BlockCosigEntry{{Hash: hash, Cosig: c}}, "")
}

// onBlockCosigs keeps the valid new cosignatures, relays them, and executes any
// block they make final.
func (dm *Daemon) onBlockCosigs(addr string, entries []BlockCosigEntry) {
	rule := dm.multiSignerRule()
	if !rule.Enabled() {
		return
	}
	var fresh []BlockCosigEntry
	hashes := make(map[cipher.SHA256]struct{})
	for _, e := range entries {
		if rule.VerifyCosig(e.Hash, e.Cosig) != nil {
			continue
		}
		if dm.multiSigner.addCosig(e.Hash, e.Cosig) {
			fresh = append(fresh, e)
			hashes[e.Hash] = struct{}{}
			if _, ok := dm.multiSigner.proposal(e.Hash); ok {
				for _, other := range dm.multiSigner.otherVotes(rule, e.Hash, e.Cosig.Signer) {
					logger.Critical().Warningf("Signer %s signed two blocks at one height: %s and %s", rule.Signers[e.Cosig.Signer].Hex(), e.Hash.Hex(), other.Hex())
				}
			}
		}
	}
	if len(fresh) == 0 {
		return
	}
	dm.relayBlockCosigs(fresh, addr)
	if _, err := dm.multiSignerHead(); err != nil {
		logger.WithError(err).Error("onBlockCosigs: head block")
		return
	}
	for h := range hashes {
		dm.tryExecuteProposal(h)
	}
}

// tryExecuteProposal executes a held block once it has k signatures, then
// relays it with its signatures.
func (dm *Daemon) tryExecuteProposal(hash cipher.SHA256) {
	p, ok := dm.multiSigner.proposal(hash)
	if !ok {
		return
	}
	b := p.block
	b.Cosigs = dm.multiSigner.cosigsFor(hash)
	if _, err := dm.multiSignerRule().Verify(b); err != nil {
		return
	}
	if err := dm.visor.ExecuteSignedBlock(b); err != nil {
		logger.WithError(err).WithField("seq", b.Seq()).Error("Failed to execute a final multi-signer block")
		return
	}
	signers := dm.multiSignerRule().BlockSigners(b)
	keys := make([]string, len(signers))
	for i, pk := range signers {
		keys[i] = pk.Hex()
	}
	logger.Critical().WithFields(logrus.Fields{
		"seq":     b.Seq(),
		"hash":    hash.Hex(),
		"signers": keys,
	}).Info("Added new multi-signer block")

	if _, err := dm.multiSignerHead(); err != nil {
		logger.WithError(err).Error("tryExecuteProposal: head block")
	}
	if err := dm.broadcastBlock(b); err != nil {
		logger.WithError(err).Warning("tryExecuteProposal: broadcast block")
	}
}

// holdBlock keeps a block that arrived without enough signatures, as the next
// block, until they arrive, and asks peers for its signatures again.
func (dm *Daemon) holdBlock(b coin.SignedBlock) {
	rule := dm.multiSignerRule()
	if !rule.Applies(b.Seq()) || rule.ProposerIndex(b) < 0 {
		return
	}
	if _, err := dm.multiSignerHead(); err != nil {
		return
	}
	b.Cosigs = nil
	dm.multiSigner.addProposal(b, 0)
	dm.tryExecuteProposal(b.HashHeader())

	s := dm.multiSigner
	s.mu.Lock()
	ask := time.Since(s.lastCosigsAsk) >= cosigRequestInterval
	if ask {
		s.lastCosigsAsk = time.Now()
	}
	s.mu.Unlock()
	if ask {
		if err := dm.requestBlocks(); err != nil {
			logger.WithError(err).Debug("holdBlock: request blocks")
		}
	}
}

// cosigsFor returns the cosignatures received for a block
func (dm *Daemon) cosigsFor(hash cipher.SHA256) []coin.BlockCosig {
	return dm.multiSigner.cosigsFor(hash)
}

func blockCosigEntries(blocks []coin.SignedBlock) []BlockCosigEntry {
	var entries []BlockCosigEntry
	for _, b := range blocks {
		h := b.HashHeader()
		for _, c := range b.Cosigs {
			entries = append(entries, BlockCosigEntry{Hash: h, Cosig: c})
		}
	}
	return entries
}

func newGiveBlockCosigsMessages(entries []BlockCosigEntry) []*GiveBlockCosigsMessage {
	var msgs []*GiveBlockCosigsMessage
	for len(entries) > 0 {
		n := min(len(entries), maxBlockCosigsPerMessage)
		msgs = append(msgs, &GiveBlockCosigsMessage{Cosigs: entries[:n]})
		entries = entries[n:]
	}
	return msgs
}

// sendBlockCosigs sends the cosignatures of blocks to addr ahead of the blocks
// themselves, if the peer understands them.
func (dm *Daemon) sendBlockCosigs(addr string, blocks []coin.SignedBlock) error {
	msgs := newGiveBlockCosigsMessages(blockCosigEntries(blocks))
	if len(msgs) == 0 {
		return nil
	}
	c := dm.connections.get(addr)
	if c == nil {
		return errors.New("sendBlockCosigs: no connection")
	}
	if c.ProtocolVersion < multiSignerProtocolVersion {
		return nil
	}
	for _, m := range msgs {
		if err := dm.sendMessage(addr, m); err != nil {
			return err
		}
	}
	return nil
}

// multiSignerPeers returns the introduced peers that understand the
// multi-signer messages, except skip.
func (dm *Daemon) multiSignerPeers(skip string) []string {
	var addrs []string
	for _, c := range dm.connections.all() {
		if c.Addr != skip && c.HasIntroduced() && c.ProtocolVersion >= multiSignerProtocolVersion {
			addrs = append(addrs, c.Addr)
		}
	}
	return addrs
}

// relayBlockCosigs sends cosignatures to every capable peer except skip
func (dm *Daemon) relayBlockCosigs(entries []BlockCosigEntry, skip string) {
	addrs := dm.multiSignerPeers(skip)
	if len(addrs) == 0 {
		return
	}
	for _, m := range newGiveBlockCosigsMessages(entries) {
		if _, err := dm.out.BroadcastMessage(m, addrs); err != nil {
			logger.WithError(err).Warning("Broadcast GiveBlockCosigsMessage failed")
			return
		}
	}
}

// relayProposal sends a proposal, with the PQ witnesses of its transactions,
// to every capable peer except skip
func (dm *Daemon) relayProposal(m *ProposeBlockMessage, skip string) {
	addrs := dm.multiSignerPeers(skip)
	if len(addrs) == 0 {
		return
	}
	for _, w := range newGivePQWitnessesMessages(m.Block.Body.Transactions) {
		if _, err := dm.out.BroadcastMessage(w, addrs); err != nil {
			logger.WithError(err).Warning("Broadcast GivePQWitnessesMessage failed")
		}
	}
	if _, err := dm.out.BroadcastMessage(m, addrs); err != nil {
		logger.WithError(err).Warning("Broadcast ProposeBlockMessage failed")
	}
}

// publishesNextBlock reports whether this block publisher makes the next block
// alone. On a multi-signer chain that holds only before the fork height, and
// only for the blockchain key.
func (dm *Daemon) publishesNextBlock() bool {
	rule := dm.multiSignerRule()
	if !rule.Enabled() {
		return true
	}
	pk, err := cipher.PubKeyFromSecKey(dm.visor.Config.BlockchainSeckey)
	if err != nil || pk != dm.visor.Config.BlockchainPubkey {
		return false
	}
	seq, ok, err := dm.visor.HeadBkSeq()
	return err == nil && ok && !rule.Applies(seq+1)
}
