package daemon

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor"
)

//TODO
//- download block headers
//- request blocks individually across multiple peers

//TODO
//- use CXO for blocksync

/*
Visor should not be duplicated
- this should be pushed into /src/visor
*/

// VisorConfig represents the configuration of visor
type VisorConfig struct {
	Config visor.Config
	// Disabled the visor completely
	Disabled bool
	// How often to request blocks from peers
	BlocksRequestRate time.Duration
	// How often to announce our blocks to peers
	BlocksAnnounceRate time.Duration
	// How many blocks to respond with to a GetBlocksMessage
	BlocksResponseCount uint64
	//how long between saving copies of the blockchain
	BlockchainBackupRate time.Duration
	// Max announce txns hash number
	MaxTxnAnnounceNum int
	// How often to announce our unconfirmed txns to peers
	TxnsAnnounceRate time.Duration
}

// NewVisorConfig creates default visor config
func NewVisorConfig() VisorConfig {
	return VisorConfig{
		Config:               visor.NewVisorConfig(),
		Disabled:             false,
		BlocksRequestRate:    time.Second * 60, //backup, could be disabled
		BlocksAnnounceRate:   time.Second * 60, //backup, could be disabled
		BlocksResponseCount:  20,
		BlockchainBackupRate: time.Second * 30,
		MaxTxnAnnounceNum:    16,
		TxnsAnnounceRate:     time.Minute,
	}
}

// Visor struct
type Visor struct {
	Config VisorConfig
	v      *visor.Visor
	// Peer-reported blockchain length.  Use to estimate download progress
	blockchainLengths map[string]uint64
	reqC              chan reqFunc // all request will go through this channel, to keep writing and reading member variable thread safe.
	Shutdown          context.CancelFunc
}

type reqFunc func(context.Context)

// NewVisor creates visor instance
func NewVisor(c VisorConfig) (*Visor, error) {
	if c.Disabled {
		return &Visor{
			Config:            c,
			blockchainLengths: make(map[string]uint64),
			reqC:              make(chan reqFunc, 100),
		}, nil
	}

	var v *visor.Visor
	v, closeVs, err := visor.NewVisor(c.Config)
	if err != nil {
		return nil, err
	}

	vs := &Visor{
		Config:            c,
		v:                 v,
		blockchainLengths: make(map[string]uint64),
		reqC:              make(chan reqFunc, 100),
	}

	vs.Shutdown = func() {
		// close the visor
		closeVs()
	}

	return vs, nil
}

// Run starts the visor
func (vs *Visor) Run() error {
	defer logger.Info("Visor closed")
	errC := make(chan error, 1)
	go func() {
		// vs.Shutdown will notify the vs.v.Run to return.
		errC <- vs.v.Run()
	}()

	for {
		select {
		case err := <-errC:
			return err
		case req := <-vs.reqC:
			func() {
				cxt, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
				defer cancel()
				req(cxt)
			}()
		}
	}
}

// the callback function must not be blocked.
func (vs *Visor) strand(f func()) {
	done := make(chan struct{})
	vs.reqC <- func(cxt context.Context) {
		defer close(done)
		c := make(chan struct{})
		go func() {
			defer close(c)
			f()
		}()
		select {
		case <-cxt.Done():
			logger.Error("%v", cxt.Err())
			return
		case <-c:
			return
		}
	}
	<-done
}

// RefreshUnconfirmed checks unconfirmed txns against the blockchain and purges ones too old
func (vs *Visor) RefreshUnconfirmed() (hashes []cipher.SHA256) {
	if vs.Config.Disabled {
		return
	}
	vs.strand(func() {
		hashes = vs.v.RefreshUnconfirmed()
	})
	return
}

// RequestBlocks Sends a GetBlocksMessage to all connections
func (vs *Visor) RequestBlocks(pool *Pool) {
	if vs.Config.Disabled {
		return
	}
	vs.strand(func() {
		m := NewGetBlocksMessage(vs.v.HeadBkSeq(), vs.Config.BlocksResponseCount)
		pool.Pool.BroadcastMessage(m)
	})
}

// AnnounceBlocks sends an AnnounceBlocksMessage to all connections
func (vs *Visor) AnnounceBlocks(pool *Pool) {
	if vs.Config.Disabled {
		return
	}
	vs.strand(func() {
		m := NewAnnounceBlocksMessage(vs.v.HeadBkSeq())
		pool.Pool.BroadcastMessage(m)
	})
}

// AnnounceAllTxns announces local unconfirmed transactions
func (vs *Visor) AnnounceAllTxns(pool *Pool) {
	if vs.Config.Disabled {
		return
	}
	vs.strand(func() {
		// get local unconfirmed transaction hashes.
		hashes := vs.v.GetAllValidUnconfirmedTxHashes()
		// filter all thoses invalid txns
		hashesSet := divideHashes(hashes, vs.Config.MaxTxnAnnounceNum)
		for _, hs := range hashesSet {
			m := NewAnnounceTxnsMessage(hs)
			if err := pool.Pool.BroadcastMessage(m); err != nil {
				logger.Debug("Broadcast AnnounceTxnsMessage failed, err:%v", err)
				return
			}
		}
	})
}

// AnnounceTxns announce given transaction hashes.
func (vs *Visor) AnnounceTxns(pool *Pool, txns []cipher.SHA256) {
	if vs.Config.Disabled {
		return
	}
	if len(txns) <= 0 {
		return
	}

	vs.strand(func() {
		m := NewAnnounceTxnsMessage(txns)
		if err := pool.Pool.BroadcastMessage(m); err != nil {
			logger.Debug("Broadcast AnnounceTxnsMessage failed, err:%v", err)
		}
	})
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
	if vs.Config.Disabled {
		return errors.New("Visor disabled")
	}
	var err error
	vs.strand(func() {
		m := NewGetBlocksMessage(vs.v.HeadBkSeq(), vs.Config.BlocksResponseCount)
		var exist bool
		exist, err = pool.Pool.IsConnExist(addr)
		if err != nil {
			return
		}

		if !exist {
			err = fmt.Errorf("Tried to send GetBlocksMessage to %s, but we're "+
				"not connected", addr)
			return
		}
		err = pool.Pool.SendMessage(addr, m)
	})
	return err
}

// SetTxnsAnnounced sets all txns as announced
func (vs *Visor) SetTxnsAnnounced(txns []cipher.SHA256) {
	vs.strand(func() {
		now := utc.Now()
		for _, h := range txns {
			vs.v.Unconfirmed.SetAnnounced(h, now)
		}
	})
}

// Sends a signed block to all connections.
// TODO: deprecate, should only send to clients that request by hash
func (vs *Visor) broadcastBlock(sb coin.SignedBlock, pool *Pool) {
	if vs.Config.Disabled {
		return
	}
	m := NewGiveBlocksMessage([]coin.SignedBlock{sb})
	pool.Pool.BroadcastMessage(m)
}

// broadcastTransaction broadcasts a single transaction to all peers.
func (vs *Visor) broadcastTransaction(t coin.Transaction, pool *Pool) {
	if vs.Config.Disabled {
		logger.Debug("broadcast tx disabled")
		return
	}
	m := NewGiveTxnsMessage(coin.Transactions{t})
	l, err := pool.Pool.Size()
	if err != nil {
		logger.Error("Broadcast GivenTxnsMessage failed: %v", err)
		return
	}

	logger.Debug("Broadcasting GiveTxnsMessage to %d conns", l)
	pool.Pool.BroadcastMessage(m)
}

// InjectTransaction injects transaction to the unconfirmed pool and broadcasts it
// The transaction must have a valid fee, be well-formed and not spend timelocked outputs.
func (vs *Visor) InjectTransaction(txn coin.Transaction, pool *Pool) error {
	var err error
	vs.strand(func() {
		if err = vs.injectTransaction(txn, pool); err != nil {
			return
		}

		vs.broadcastTransaction(txn, pool)
	})
	return err
}

func (vs *Visor) injectTransaction(txn coin.Transaction, pool *Pool) error {
	if err := vs.verifyTransaction(txn); err != nil {
		return err
	}

	_, err := vs.v.InjectTxn(txn)
	return err
}

func (vs *Visor) verifyTransaction(txn coin.Transaction) error {
	inUxs, err := vs.v.Blockchain.Unspent().GetArray(txn.In)
	if err != nil {
		return err
	}

	fee, err := visor.TransactionFee(&txn, vs.v.Blockchain.Time(), inUxs)
	if err != nil {
		return err
	}

	if err := visor.VerifyTransactionFee(&txn, fee); err != nil {
		return err
	}

	if visor.TransactionIsLocked(inUxs) {
		return errors.New("Transaction has locked address inputs")
	}

	if err := txn.Verify(); err != nil {
		return fmt.Errorf("Transaction Verification Failed, %v", err)
	}

	// valid the spending coins
	for _, out := range txn.Out {
		if err := DropletPrecisionCheck(out.Coins); err != nil {
			return err
		}
	}

	return nil
}

// ResendTransaction resends a known UnconfirmedTxn.
func (vs *Visor) ResendTransaction(h cipher.SHA256, pool *Pool) {
	if vs.Config.Disabled {
		return
	}
	vs.strand(func() {
		if ut, ok := vs.v.Unconfirmed.Get(h); ok {
			vs.broadcastTransaction(ut.Txn, pool)
		}
	})
	return
}

// ResendUnconfirmedTxns resents all unconfirmed transactions
func (vs *Visor) ResendUnconfirmedTxns(pool *Pool) []cipher.SHA256 {
	var txids []cipher.SHA256
	if vs.Config.Disabled {
		return txids
	}
	vs.strand(func() {
		txns := vs.v.GetAllUnconfirmedTxns()

		for i := range txns {
			logger.Debugf("Rebroadcast tx %s", txns[i].Hash().Hex())
			vs.broadcastTransaction(txns[i].Txn, pool)
			txids = append(txids, txns[i].Txn.Hash())
		}
	})
	return txids
}

// CreateAndPublishBlock creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a master chain.  Returns creation error and
// whether it was published or not
func (vs *Visor) CreateAndPublishBlock(pool *Pool) error {
	if vs.Config.Disabled {
		return errors.New("Visor disabled")
	}
	var err error
	vs.strand(func() {
		var sb coin.SignedBlock
		sb, err = vs.v.CreateAndExecuteBlock()
		if err != nil {
			return
		}
		vs.broadcastBlock(sb, pool)
	})
	return err
}

// RemoveConnection updates internal state when a connection disconnects
func (vs *Visor) RemoveConnection(addr string) {
	vs.strand(func() {
		delete(vs.blockchainLengths, addr)
	})
}

// RecordBlockchainLength saves a peer-reported blockchain length
func (vs *Visor) RecordBlockchainLength(addr string, bkLen uint64) {
	vs.strand(func() {
		vs.blockchainLengths[addr] = bkLen
	})
}

// EstimateBlockchainLength returns the blockchain length estimated from peer reports
// Deprecate. Should not need. Just report time of last block
func (vs *Visor) EstimateBlockchainLength() uint64 {
	var maxLen uint64
	vs.strand(func() {
		ourLen := vs.v.HeadBkSeq()
		if len(vs.blockchainLengths) < 2 {
			maxLen = ourLen
			return
		}
		for _, seq := range vs.blockchainLengths {
			if maxLen < seq {
				maxLen = seq
			}
		}
	})
	return maxLen
}

// HeadBkSeq returns the head sequence
func (vs *Visor) HeadBkSeq() uint64 {
	var seq uint64
	vs.strand(func() {
		seq = vs.v.HeadBkSeq()
	})
	return seq
}

// ExecuteSignedBlock executes signed block
func (vs *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	var err error
	vs.strand(func() {
		err = vs.v.ExecuteSignedBlock(b)
	})
	return err
}

// GetSignedBlocksSince returns numbers of signed blocks since seq.
func (vs *Visor) GetSignedBlocksSince(seq uint64, num uint64) (sbs []coin.SignedBlock, err error) {
	vs.strand(func() {
		sbs, err = vs.v.GetSignedBlocksSince(seq, num)
	})
	return
}

// UnConfirmFilterKnown returns all unknow transaction hashes
func (vs *Visor) UnConfirmFilterKnown(txns []cipher.SHA256) []cipher.SHA256 {
	var ts []cipher.SHA256
	vs.strand(func() {
		ts = vs.v.Unconfirmed.FilterKnown(txns)
	})
	return ts
}

// UnConfirmKnow returns all know tansactions
func (vs *Visor) UnConfirmKnow(hashes []cipher.SHA256) (txns coin.Transactions) {
	vs.strand(func() {
		txns = vs.v.Unconfirmed.GetKnown(hashes)
	})
	return
}

// InjectTxn only try to append transaction into local blockchain, don't broadcast it.
func (vs *Visor) InjectTxn(tx coin.Transaction) (know bool, err error) {
	vs.strand(func() {
		know, err = vs.v.InjectTxn(tx)
	})
	return
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
		RequestedBlocks: requestedBlocks, //count of blocks requested
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
	// TODO -- move 20 to either Messages.Config or Visor.Config
	if d.Visor.Config.Disabled {
		return
	}
	// Record this as this peer's highest block
	d.Visor.RecordBlockchainLength(gbm.c.Addr, gbm.LastBlock)
	// Fetch and return signed blocks since LastBlock
	blocks, err := d.Visor.GetSignedBlocksSince(gbm.LastBlock, gbm.RequestedBlocks)
	if err != nil {
		logger.Info("Get signed blocks failed: %v", err)
		return
	}

	logger.Debug("Got %d blocks since %d", len(blocks), gbm.LastBlock)
	if len(blocks) == 0 {
		return
	}
	m := NewGiveBlocksMessage(blocks)
	d.Pool.Pool.SendMessage(gbm.c.Addr, m)
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
	if d.Visor.Config.Disabled {
		logger.Critical("Visor disabled, ignoring GiveBlocksMessage")
		return
	}
	processed := 0
	maxSeq := d.Visor.HeadBkSeq()
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
			logger.Critical("Added new block %d", b.Block.Head.BkSeq)
			processed++
		} else {
			logger.Critical("Failed to execute received block: %v", err)
			// Blocks must be received in order, so if one fails its assumed
			// the rest are failing
			break
		}
	}
	if processed == 0 {
		return
	}

	headBkSeq := d.Visor.HeadBkSeq()
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
	if d.Visor.Config.Disabled {
		return
	}
	headBkSeq := d.Visor.HeadBkSeq()
	if headBkSeq >= abm.MaxBkSeq {
		return
	}
	//should this be block get request for current sequence?
	//if client is not caught up, wont attempt to get block
	m := NewGetBlocksMessage(headBkSeq, d.Visor.Config.BlocksResponseCount)
	d.Pool.Pool.SendMessage(abm.c.Addr, m)
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
	if d.Visor.Config.Disabled {
		return
	}
	unknown := d.Visor.UnConfirmFilterKnown(atm.Txns)
	if len(unknown) == 0 {
		return
	}
	m := NewGetTxnsMessage(unknown)
	d.Pool.Pool.SendMessage(atm.c.Addr, m)
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
func (gtm *GetTxnsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	gtm.c = mc
	return daemon.(*Daemon).recordMessageEvent(gtm, mc)
}

// Process process message
func (gtm *GetTxnsMessage) Process(d *Daemon) {
	if d.Visor.Config.Disabled {
		return
	}
	// Locate all txns from the unconfirmed pool
	// reply to sender with GiveTxnsMessage
	known := d.Visor.UnConfirmKnow(gtm.Txns)
	if len(known) == 0 {
		return
	}
	logger.Debug("%d/%d txns known", len(known), len(gtm.Txns))
	m := NewGiveTxnsMessage(known)
	d.Pool.Pool.SendMessage(gtm.c.Addr, m)
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
	if d.Visor.Config.Disabled {
		return
	}
	if len(gtm.Txns) > 32 {
		logger.Warning("More than 32 transactions in pool. Implement breaking transactions transmission into multiple packets")
	}

	hashes := make([]cipher.SHA256, 0, len(gtm.Txns))
	// Update unconfirmed pool with these transactions
	for _, txn := range gtm.Txns {
		// Only announce transactions that are new to us, so that peers can't
		// spam relays
		known, err := d.Visor.InjectTxn(txn)
		if err != nil {
			logger.Warning("Failed to record transaction %s: %v", txn.Hash().Hex(), err)
			continue
		}

		if known {
			logger.Warning("Duplicate Transaction: %s", txn.Hash().Hex())
		} else {
			hashes = append(hashes, txn.Hash())
		}
	}
	// Announce these transactions to peers
	if len(hashes) != 0 {
		logger.Debugf("Announce %d transactions", len(hashes))
		m := NewAnnounceTxnsMessage(hashes)
		d.Pool.Pool.BroadcastMessage(m)
	}
}

// BlockchainLengths an array of uint64
type BlockchainLengths []uint64

// Len for sorting
func (bcl BlockchainLengths) Len() int {
	return len(bcl)
}

// Swap for sorting
func (bcl BlockchainLengths) Swap(i, j int) {
	bcl[i], bcl[j] = bcl[j], bcl[i]
}

// Less for sorting
func (bcl BlockchainLengths) Less(i, j int) bool {
	return bcl[i] < bcl[j]
}

type byTxnRecvTime []visor.UnconfirmedTxn

func (txs byTxnRecvTime) Len() int {
	return len(txs)
}

func (txs byTxnRecvTime) Swap(i, j int) {
	txs[i], txs[j] = txs[j], txs[i]
}

func (txs byTxnRecvTime) Less(i, j int) bool {
	return txs[i].Received < txs[j].Received
}
