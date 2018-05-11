package daemon

import (
	"errors"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

//TODO
//- download block headers
//- request blocks individually across multiple peers

//TODO
//- use CXO for blocksync

// VisorConfig represents the configuration of visor
type VisorConfig struct {
	Config visor.Config
	// Disable visor networking
	DisableNetworking bool
	// How often to request blocks from peers
	BlocksRequestRate time.Duration
	// How often to announce our blocks to peers
	BlocksAnnounceRate time.Duration
	// How many blocks to respond with to a GetBlocksMessage
	BlocksResponseCount uint64
	// How long between saving copies of the blockchain
	BlockchainBackupRate time.Duration
	// Max announce txns hash number
	MaxTxnAnnounceNum int
	// How often to announce our unconfirmed txns to peers
	TxnsAnnounceRate time.Duration
	// How long to wait for Visor request to process
	RequestDeadline time.Duration
}

// NewVisorConfig creates default visor config
func NewVisorConfig() VisorConfig {
	return VisorConfig{
		Config:               visor.NewVisorConfig(),
		DisableNetworking:    false,
		BlocksRequestRate:    time.Second * 60,
		BlocksAnnounceRate:   time.Second * 60,
		BlocksResponseCount:  20,
		BlockchainBackupRate: time.Second * 30,
		MaxTxnAnnounceNum:    16,
		TxnsAnnounceRate:     time.Minute,
		RequestDeadline:      time.Second * 3,
	}
}

// Visor struct
type Visor struct {
	Config VisorConfig
	v      *visor.Visor

	// Peer-reported blockchain height.  Use to estimate download progress
	blockchainHeights    map[string]uint64
	blockchanHeightsLock sync.Mutex
}

// NewVisor creates visor instance
func NewVisor(c VisorConfig, db *dbutil.DB) (*Visor, error) {
	vs := &Visor{
		Config:            c,
		blockchainHeights: make(map[string]uint64),
	}

	v, err := visor.NewVisor(c.Config, db)
	if err != nil {
		return nil, err
	}

	vs.v = v

	return vs, nil
}

// RefreshUnconfirmed checks unconfirmed txns against the blockchain and marks
// and returns those that become valid
func (vs *Visor) RefreshUnconfirmed() ([]cipher.SHA256, error) {
	return vs.v.RefreshUnconfirmed()
}

// RemoveInvalidUnconfirmed checks unconfirmed txns against the blockchain and
// purges those that become permanently invalid, violating hard constraints
func (vs *Visor) RemoveInvalidUnconfirmed() ([]cipher.SHA256, error) {
	return vs.v.RemoveInvalidUnconfirmed()
}

// RequestBlocks Sends a GetBlocksMessage to all connections
func (vs *Visor) RequestBlocks(pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	headSeq, ok, err := vs.v.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot request blocks, there is no head block")
	}

	m := NewGetBlocksMessage(headSeq, vs.Config.BlocksResponseCount)

	err = pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Debugf("Broadcast GetBlocksMessage failed: %v", err)
	}

	return err
}

// AnnounceBlocks sends an AnnounceBlocksMessage to all connections
func (vs *Visor) AnnounceBlocks(pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	headSeq, ok, err := vs.v.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot announce blocks, there is no head block")
	}

	m := NewAnnounceBlocksMessage(headSeq)

	err = pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Debugf("Broadcast AnnounceBlocksMessage failed: %v", err)
	}

	return err
}

// AnnounceAllTxns announces local unconfirmed transactions
func (vs *Visor) AnnounceAllTxns(pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	// Get local unconfirmed transaction hashes.
	hashes, err := vs.v.GetAllValidUnconfirmedTxHashes()
	if err != nil {
		return err
	}

	// Divide hashes into multiple sets of max size
	hashesSet := divideHashes(hashes, vs.Config.MaxTxnAnnounceNum)

	for _, hs := range hashesSet {
		m := NewAnnounceTxnsMessage(hs)
		if err = pool.Pool.BroadcastMessage(m); err != nil {
			break
		}
	}

	if err != nil {
		logger.Debugf("Broadcast AnnounceTxnsMessage failed, err: %v", err)
	}

	return err
}

// AnnounceTxns announces given transaction hashes.
func (vs *Visor) AnnounceTxns(pool *Pool, txns []cipher.SHA256) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	if len(txns) == 0 {
		return nil
	}

	m := NewAnnounceTxnsMessage(txns)

	err := pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Debugf("Broadcast AnnounceTxnsMessage failed: %v", err)
	}

	return err
}

func divideHashes(hashes []cipher.SHA256, n int) [][]cipher.SHA256 {
	if len(hashes) == 0 {
		return [][]cipher.SHA256{}
	}

	var j int
	var hashesArray [][]cipher.SHA256

	if len(hashes) > n {
		for i := range hashes {
			if len(hashes[j:i]) == n {
				hs := make([]cipher.SHA256, n)
				copy(hs, hashes[j:i])
				hashesArray = append(hashesArray, hs)
				j = i
			}
		}
	}

	hs := make([]cipher.SHA256, len(hashes)-j)
	copy(hs, hashes[j:])
	hashesArray = append(hashesArray, hs)
	return hashesArray
}

// RequestBlocksFromAddr sends a GetBlocksMessage to one connected address
func (vs *Visor) RequestBlocksFromAddr(pool *Pool, addr string) error {
	if vs.Config.DisableNetworking {
		return errors.New("Visor disabled")
	}

	headSeq, ok, err := vs.v.HeadBkSeq()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Cannot request blocks from addr, there is no head block")
	}

	m := NewGetBlocksMessage(headSeq, vs.Config.BlocksResponseCount)

	return pool.Pool.SendMessage(addr, m)
}

// SetTxnsAnnounced sets all txns as announced
func (vs *Visor) SetTxnsAnnounced(txns map[cipher.SHA256]int64) error {
	if err := vs.v.SetTxnsAnnounced(txns); err != nil {
		logger.WithError(err).Error("Failed to set unconfirmed txn announce time")
		return err
	}

	return nil
}

// InjectBroadcastTransaction injects transaction to the unconfirmed pool and broadcasts it.
// If the transaction violates either hard or soft constraints, it is not broadcast.
// This method is to be used by user-initiated transaction injections.
// For transactions received over the network, use InjectTransaction and check the result to
// decide on repropagation.
func (vs *Visor) InjectBroadcastTransaction(txn coin.Transaction, pool *Pool) error {
	if _, err := vs.v.InjectTransactionStrict(txn); err != nil {
		return err
	}

	return vs.broadcastTransaction(txn, pool)
}

// InjectTransaction adds a transaction to the unconfirmed txn pool if it does not violate hard constraints.
// The transaction is added to the pool if it only violates soft constraints.
// If a soft constraint is violated, the specific error is returned separately.
func (vs *Visor) InjectTransaction(txn coin.Transaction) (bool, *visor.ErrTxnViolatesSoftConstraint, error) {
	return vs.v.InjectTransaction(txn)
}

// Sends a signed block to all connections.
// TODO: deprecate, should only send to clients that request by hash
func (vs *Visor) broadcastBlock(sb coin.SignedBlock, pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	m := NewGiveBlocksMessage([]coin.SignedBlock{sb})
	return pool.Pool.BroadcastMessage(m)
}

// broadcastTransaction broadcasts a single transaction to all peers.
func (vs *Visor) broadcastTransaction(t coin.Transaction, pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	m := NewGiveTxnsMessage(coin.Transactions{t})
	l, err := pool.Pool.Size()
	if err != nil {
		return err
	}

	logger.Debugf("Broadcasting GiveTxnsMessage to %d conns", l)

	err = pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Errorf("Broadcast GivenTxnsMessage failed: %v", err)
	}

	return err
}

// ResendTransaction resends a known UnconfirmedTxn
func (vs *Visor) ResendTransaction(h cipher.SHA256, pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	ut, err := vs.v.GetUnconfirmedTxn(h)
	if err != nil {
		return err
	}

	if ut != nil {
		return vs.broadcastTransaction(ut.Txn, pool)
	}

	return nil
}

// ResendUnconfirmedTxns resends all unconfirmed transactions and returns the hashes that were successfully rebroadcast
func (vs *Visor) ResendUnconfirmedTxns(pool *Pool) ([]cipher.SHA256, error) {
	if vs.Config.DisableNetworking {
		return nil, nil
	}

	txns, err := vs.v.GetAllUnconfirmedTxns()
	if err != nil {
		return nil, err
	}

	var txids []cipher.SHA256
	for i := range txns {
		logger.Debugf("Rebroadcast tx %s", txns[i].Hash().Hex())
		if err := vs.broadcastTransaction(txns[i].Txn, pool); err == nil {
			txids = append(txids, txns[i].Txn.Hash())
		}
	}

	return txids, nil
}

// CreateAndPublishBlock creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a master chain.
// Even if the error is not nil, the block may have been created an is returned.
// The error can be non-nil if the creation failed or if the broadcast failed.
func (vs *Visor) CreateAndPublishBlock(pool *Pool) (*coin.SignedBlock, error) {
	if vs.Config.DisableNetworking {
		return nil, errors.New("Visor disabled")
	}

	sb, err := vs.v.CreateAndExecuteBlock()
	if err != nil {
		return nil, err
	}

	err = vs.broadcastBlock(sb, pool)

	return &sb, err
}

// RemoveConnection updates internal state when a connection disconnects
func (vs *Visor) RemoveConnection(addr string) {
	vs.blockchainHeightsMutex.Lock()
	defer vs.blockchainHeightsMutex.Unlock()

	delete(vs.blockchainHeights, addr)
}

// RecordBlockchainHeight saves a peer-reported blockchain length
func (vs *Visor) RecordBlockchainHeight(addr string, bkLen uint64) {
	vs.blockchainHeightsMutex.Lock()
	defer vs.blockchainHeightsMutex.Unlock()

	vs.blockchainHeights[addr] = bkLen
}

// EstimateBlockchainHeight returns the blockchain length estimated from peer reports
// Deprecate. Should not need. Just report time of last block
func (vs *Visor) EstimateBlockchainHeight() (uint64, error) {
	maxLen, _, err := vs.v.HeadBkSeq()
	if err != nil {
		return 0, err
	}

	vs.blockchainHeightsMutex.Lock()
	defer vs.blockchainHeightsMutex.Unlock()

	for _, seq := range vs.blockchainHeights {
		if maxLen < seq {
			maxLen = seq
		}
	}

	return maxLen, nil
}

// PeerBlockchainHeight is a peer's IP address with their reported blockchain height
type PeerBlockchainHeight struct {
	Address string
	Height  uint64
}

// GetPeerBlockchainHeights returns recorded peers' blockchain heights as an array.
func (vs *Visor) GetPeerBlockchainHeights() []PeerBlockchainHeight {
	vs.blockchainHeightsMutex.Lock()
	defer vs.blockchainHeightsMutex.Unlock()

	if len(vs.blockchainHeights) == 0 {
		return nil
	}

	peerHeights := make([]PeerBlockchainHeight, 0, len(vs.blockchainHeights))
	for addr, height := range vs.blockchainHeights {
		peerHeights = append(peerHeights, PeerBlockchainHeight{
			Address: addr,
			Height:  height,
		})
	}

	return peerHeights
}

// HeadBkSeq returns the head sequence
func (vs *Visor) HeadBkSeq() (uint64, bool, error) {
	return vs.v.HeadBkSeq()
}

// ExecuteSignedBlock executes signed blocks
func (vs *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	return vs.v.ExecuteSignedBlock(b)
}

// GetSignedBlocksSince returns signed blocks in an inclusive range of [seq+1, seq+ct]
func (vs *Visor) GetSignedBlocksSince(seq uint64, ct uint64) ([]coin.SignedBlock, error) {
	return vs.v.GetSignedBlocksSince(seq, ct)
}

// UnconfirmedUnknown returns all unknown transaction hashes
func (vs *Visor) UnconfirmedUnknown(txns []cipher.SHA256) ([]cipher.SHA256, error) {
	return vs.v.GetUnconfirmedUnknown(txns)
}

// UnconfirmedKnown returns all know tansactions
func (vs *Visor) UnconfirmedKnown(hashes []cipher.SHA256) (coin.Transactions, error) {
	return vs.v.GetUnconfirmedKnown(hashes)
}

// Communication layer for the coin pkg

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
	return daemon.(*Daemon).recordMessageEvent(gbm, mc)
}

// Process should send number to be requested, with request
func (gbm *GetBlocksMessage) Process(d *Daemon) {
	// TODO -- we need the sig to be sent with the block, but only the master
	// can sign blocks.  Thus the sig needs to be stored with the block.
	// TODO -- move to either Messages.Config or Visor.Config
	if d.Visor.Config.DisableNetworking {
		return
	}
	// Record this as this peer's highest block
	d.Visor.RecordBlockchainHeight(gbm.c.Addr, gbm.LastBlock)
	// Fetch and return signed blocks since LastBlock
	blocks, err := d.Visor.GetSignedBlocksSince(gbm.LastBlock, gbm.RequestedBlocks)
	if err != nil {
		logger.Infof("Get signed blocks failed: %v", err)
		return
	}

	if len(blocks) == 0 {
		return
	}

	logger.Debugf("Got %d blocks since %d", len(blocks), gbm.LastBlock)

	m := NewGiveBlocksMessage(blocks)
	if err := d.Pool.Pool.SendMessage(gbm.c.Addr, m); err != nil {
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
func (gbm *GiveBlocksMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	gbm.c = mc
	return daemon.(*Daemon).recordMessageEvent(gbm, mc)
}

// Process process message
func (gbm *GiveBlocksMessage) Process(d *Daemon) {
	if d.Visor.Config.DisableNetworking {
		logger.Critical().Info("Visor disabled, ignoring GiveBlocksMessage")
		return
	}

	// These DB queries are not performed in a transaction for performance reasons.
	// It is not necessary that the blocks be executed together in a single transaction.

	processed := 0
	maxSeq, ok, err := d.Visor.HeadBkSeq()
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

		err := d.Visor.ExecuteSignedBlock(b)
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

	headBkSeq, ok, err := d.Visor.HeadBkSeq()
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
	d.Pool.Pool.BroadcastMessage(m1)
	//request more blocks.
	m2 := NewGetBlocksMessage(headBkSeq, d.Visor.Config.BlocksResponseCount)
	d.Pool.Pool.BroadcastMessage(m2)
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
func (abm *AnnounceBlocksMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	abm.c = mc
	return daemon.(*Daemon).recordMessageEvent(abm, mc)
}

// Process process message
func (abm *AnnounceBlocksMessage) Process(d *Daemon) {
	if d.Visor.Config.DisableNetworking {
		return
	}

	headBkSeq, ok, err := d.Visor.HeadBkSeq()
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
	m := NewGetBlocksMessage(headBkSeq, d.Visor.Config.BlocksResponseCount)
	if err := d.Pool.Pool.SendMessage(abm.c.Addr, m); err != nil {
		logger.Errorf("Send GetBlocksMessage to %s failed: %v", abm.c.Addr, err)
	}
}

// SendingTxnsMessage send transaction message interface
type SendingTxnsMessage interface {
	GetTxns() []cipher.SHA256
}

// AnnounceTxnsMessage tells a peer that we have these transactions
type AnnounceTxnsMessage struct {
	Txns []cipher.SHA256
	c    *gnet.MessageContext `enc:"-"`
}

// NewAnnounceTxnsMessage creates announce txns message
func NewAnnounceTxnsMessage(txns []cipher.SHA256) *AnnounceTxnsMessage {
	return &AnnounceTxnsMessage{
		Txns: txns,
	}
}

// GetTxns returns txns
func (atm *AnnounceTxnsMessage) GetTxns() []cipher.SHA256 {
	return atm.Txns
}

// Handle handle message
func (atm *AnnounceTxnsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	atm.c = mc
	return daemon.(*Daemon).recordMessageEvent(atm, mc)
}

// Process process message
func (atm *AnnounceTxnsMessage) Process(d *Daemon) {
	if d.Visor.Config.DisableNetworking {
		return
	}

	unknown, err := d.Visor.UnconfirmedUnknown(atm.Txns)
	if err != nil {
		logger.WithError(err).Error("AnnounceTxnsMessage Visor.UnconfirmedUnknown failed")
		return
	}

	if len(unknown) == 0 {
		return
	}

	m := NewGetTxnsMessage(unknown)
	if err := d.Pool.Pool.SendMessage(atm.c.Addr, m); err != nil {
		logger.Errorf("Send GetTxnsMessage to %s failed: %v", atm.c.Addr, err)
	}
}

// GetTxnsMessage request transactions of given hash
type GetTxnsMessage struct {
	Txns []cipher.SHA256
	c    *gnet.MessageContext `enc:"-"`
}

// NewGetTxnsMessage creates GetTxnsMessage
func NewGetTxnsMessage(txns []cipher.SHA256) *GetTxnsMessage {
	return &GetTxnsMessage{
		Txns: txns,
	}
}

// Handle handle message
func (gtm *GetTxnsMessage) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	gtm.c = mc
	return daemon.(*Daemon).recordMessageEvent(gtm, mc)
}

// Process process message
func (gtm *GetTxnsMessage) Process(d *Daemon) {
	if d.Visor.Config.DisableNetworking {
		return
	}

	// Locate all txns from the unconfirmed pool
	known, err := d.Visor.UnconfirmedKnown(gtm.Txns)
	if err != nil {
		logger.WithError(err).Error("GetTxnsMessage Visor.UnconfirmedKnown failed")
		return
	}
	if len(known) == 0 {
		return
	}

	// Reply to sender with GiveTxnsMessage
	m := NewGiveTxnsMessage(known)
	if err := d.Pool.Pool.SendMessage(gtm.c.Addr, m); err != nil {
		logger.Errorf("Send GiveTxnsMessage to %s failed: %v", gtm.c.Addr, err)
	}
}

// GiveTxnsMessage tells the transaction of given hashes
type GiveTxnsMessage struct {
	Txns coin.Transactions
	c    *gnet.MessageContext `enc:"-"`
}

// NewGiveTxnsMessage creates GiveTxnsMessage
func NewGiveTxnsMessage(txns coin.Transactions) *GiveTxnsMessage {
	return &GiveTxnsMessage{
		Txns: txns,
	}
}

// GetTxns returns transactions hashes
func (gtm *GiveTxnsMessage) GetTxns() []cipher.SHA256 {
	return gtm.Txns.Hashes()
}

// Handle handle message
func (gtm *GiveTxnsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	gtm.c = mc
	return daemon.(*Daemon).recordMessageEvent(gtm, mc)
}

// Process process message
func (gtm *GiveTxnsMessage) Process(d *Daemon) {
	if d.Visor.Config.DisableNetworking {
		return
	}

	hashes := make([]cipher.SHA256, 0, len(gtm.Txns))
	// Update unconfirmed pool with these transactions
	for _, txn := range gtm.Txns {
		// Only announce transactions that are new to us, so that peers can't spam relays
		known, softErr, err := d.Visor.InjectTransaction(txn)
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
		d.Pool.Pool.BroadcastMessage(m)
	}
}
