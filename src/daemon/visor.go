package daemon

import (
	"errors"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/strand"
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
	// Internal request buffer size
	RequestBufferSize int
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
		RequestBufferSize:    100,
	}
}

// Visor struct
type Visor struct {
	Config VisorConfig
	v      *visor.Visor
	// Peer-reported blockchain height.  Use to estimate download progress
	blockchainHeights map[string]uint64
	// all request will go through this channel, to keep writing and reading member variable thread safe.
	reqC chan strand.Request
	quit chan struct{}
	done chan struct{}
}

// NewVisor creates visor instance
func NewVisor(c VisorConfig) (*Visor, error) {
	vs := &Visor{
		Config:            c,
		blockchainHeights: make(map[string]uint64),
		reqC:              make(chan strand.Request, c.RequestBufferSize),
		quit:              make(chan struct{}),
		done:              make(chan struct{}),
	}

	var v *visor.Visor
	v, err := visor.NewVisor(c.Config)
	if err != nil {
		return nil, err
	}

	vs.v = v

	return vs, nil
}

// Run starts the visor
func (vs *Visor) Run() error {
	defer close(vs.done)
	defer logger.Info("Visor closed")
	errC := make(chan error, 1)
	go func() {
		err := vs.v.Run()
		errC <- err
	}()

	for {
		select {
		case <-vs.quit:
			return nil
		case err := <-errC:
			return err
		case req := <-vs.reqC:
			if err := req.Func(); err != nil {
				logger.Error("Visor request func failed: %v", err)
			}
		}
	}
}

// Shutdown shuts down the visor
func (vs *Visor) Shutdown() {
	close(vs.quit)

	// Must wait for Run() to finish first, so that any pending calls to
	// visor.Visor have completed. These calls may use a bolt.DB transaction.
	// vs.v.Shutdown() will close the database, and all transactions must
	// be completed before closing the database.
	<-vs.done

	vs.v.Shutdown()
}

func (vs *Visor) strand(name string, f func() error) error {
	name = fmt.Sprintf("daemon.Visor.%s", name)
	return strand.Strand(logger, vs.reqC, name, f, vs.quit, nil)
}

// RefreshUnconfirmed checks unconfirmed txns against the blockchain and purges ones too old
func (vs *Visor) RefreshUnconfirmed() ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256

	if err := vs.strand("RefreshUnconfirmed", func() error {
		var err error
		hashes, err = vs.v.RefreshUnconfirmed()
		return err
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// RequestBlocks Sends a GetBlocksMessage to all connections
func (vs *Visor) RequestBlocks(pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	err := vs.strand("RequestBlocks", func() error {
		m := NewGetBlocksMessage(vs.v.HeadBkSeq(), vs.Config.BlocksResponseCount)
		return pool.Pool.BroadcastMessage(m)
	})

	if err != nil {
		logger.Debug("Broadcast GetBlocksMessage failed: %v", err)
	}

	return err
}

// AnnounceBlocks sends an AnnounceBlocksMessage to all connections
func (vs *Visor) AnnounceBlocks(pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	err := vs.strand("AnnounceBlocks", func() error {
		m := NewAnnounceBlocksMessage(vs.v.HeadBkSeq())
		return pool.Pool.BroadcastMessage(m)
	})

	if err != nil {
		logger.Debug("Broadcast AnnounceBlocksMessage failed: %v", err)
	}

	return err
}

// AnnounceAllTxns announces local unconfirmed transactions
func (vs *Visor) AnnounceAllTxns(pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	return vs.strand("AnnounceAllTxns", func() error {
		// Get local unconfirmed transaction hashes.
		hashes, err := vs.v.GetAllValidUnconfirmedTxHashes()
		if err != nil {
			return err
		}

		logger.Debug("Visor.AnnounceAllTxns: announcing %d txns", len(hashes))

		// Divide hashes into multiple sets of max size
		hashesSet := divideHashes(hashes, vs.Config.MaxTxnAnnounceNum)

		for _, hs := range hashesSet {
			m := NewAnnounceTxnsMessage(hs)
			if err := pool.Pool.BroadcastMessage(m); err != nil {
				return err
			}
		}

		return nil
	})
}

// AnnounceTxns announces given transaction hashes.
func (vs *Visor) AnnounceTxns(pool *Pool, txns []cipher.SHA256) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	if len(txns) == 0 {
		return nil
	}

	err := vs.strand("AnnounceTxns", func() error {
		m := NewAnnounceTxnsMessage(txns)
		return pool.Pool.BroadcastMessage(m)
	})

	if err != nil {
		logger.Debug("Broadcast AnnounceTxnsMessage failed: %v", err)
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

	err := vs.strand("RequestBlocksFromAddr", func() error {
		m := NewGetBlocksMessage(vs.v.HeadBkSeq(), vs.Config.BlocksResponseCount)
		exist, err := pool.Pool.IsConnExist(addr)
		if err != nil {
			return err
		}

		if !exist {
			return fmt.Errorf("Tried to send GetBlocksMessage to %s, but we are not connected", addr)
		}

		return pool.Pool.SendMessage(addr, m)
	})

	return err
}

// SetTxnsAnnounced sets all txns as announced
func (vs *Visor) SetTxnsAnnounced(txns []cipher.SHA256) {
	vs.strand("SetTxnsAnnounced", func() error {
		now := utc.Now()
		if err := vs.v.SetTxnsAnnounced(txns, now); err != nil {
			logger.Error("Failed to set unconfirmed txns announce time")
			return err
		}

		return nil
	})
}

// InjectTransaction injects transaction to the unconfirmed pool and broadcasts it
// The transaction must have a valid fee, be well-formed and not spend timelocked outputs.
func (vs *Visor) InjectTransaction(txn coin.Transaction, pool *Pool) error {
	return vs.strand("InjectTransaction", func() error {
		if err := vs.verifyInjectTransaction(txn); err != nil {
			return err
		}

		return vs.broadcastTransaction(txn, pool)
	})
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

	logger.Debug("Broadcasting GiveTxnsMessage to %d conns", l)

	err = pool.Pool.BroadcastMessage(m)
	if err != nil {
		logger.Error("Broadcast GivenTxnsMessage failed: %v", err)
	}

	return err
}

func (vs *Visor) verifyInjectTransaction(txn coin.Transaction) error {
	if err := vs.verifyTransaction(txn); err != nil {
		return err
	}

	_, err := vs.v.InjectTransaction(txn)
	return err
}

func (vs *Visor) verifyTransaction(txn coin.Transaction) error {
	// TODO -- when Unspent becomes a bolt.DB,
	// then the result from Blockchain.Unspent().GetArray()
	// and the result from GetHeadBlockTime()
	// need to be merged into a single method in visor.Visor,
	// so that a bolt.Tx can be used properly
	inUxs, err := vs.v.Blockchain.Unspent().GetArray(txn.In)
	if err != nil {
		return err
	}

	headTime, err := vs.v.GetHeadBlockTime()
	if err != nil {
		return err
	}

	fee, err := visor.TransactionFee(&txn, headTime, inUxs)
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
func (vs *Visor) ResendTransaction(h cipher.SHA256, pool *Pool) error {
	if vs.Config.DisableNetworking {
		return nil
	}

	return vs.strand("ResendTransaction", func() error {
		if ut, err := vs.v.GetUnconfirmedTxn(h); err != nil {
			return err
		} else if ut != nil {
			return vs.broadcastTransaction(ut.Txn, pool)
		}
		return nil
	})
}

// ResendUnconfirmedTxns resents all unconfirmed transactions
func (vs *Visor) ResendUnconfirmedTxns(pool *Pool) ([]cipher.SHA256, error) {
	if vs.Config.DisableNetworking {
		return nil, nil
	}

	var txids []cipher.SHA256
	if err := vs.strand("ResendUnconfirmedTxns", func() error {
		txns, err := vs.v.GetAllUnconfirmedTxns()
		if err != nil {
			return err
		}

		for i := range txns {
			logger.Debugf("Rebroadcast tx %s", txns[i].Hash().Hex())
			if err := vs.broadcastTransaction(txns[i].Txn, pool); err != nil {
				logger.Warningf("Failed to rebroadcast transaction %s", txns[i].Hash().Hex())
			} else {
				txids = append(txids, txns[i].Txn.Hash())
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return txids, nil
}

// CreateAndPublishBlock creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a master chain.  Returns creation error and
// whether it was published or not
func (vs *Visor) CreateAndPublishBlock(pool *Pool) (coin.SignedBlock, error) {
	if vs.Config.DisableNetworking {
		return coin.SignedBlock{}, errors.New("Visor disabled")
	}

	var sb coin.SignedBlock
	err := vs.strand("CreateAndPublishBlock", func() error {
		var err error
		sb, err = vs.v.CreateAndExecuteBlock()
		if err != nil {
			return err
		}

		return vs.broadcastBlock(sb, pool)
	})

	return sb, err
}

// RemoveConnection updates internal state when a connection disconnects
func (vs *Visor) RemoveConnection(addr string) {
	vs.strand("RemoveConnection", func() error {
		delete(vs.blockchainHeights, addr)
		return nil
	})
}

// RecordBlockchainHeight saves a peer-reported blockchain length
func (vs *Visor) RecordBlockchainHeight(addr string, bkLen uint64) {
	vs.strand("RecordBlockchainHeight", func() error {
		vs.blockchainHeights[addr] = bkLen
		return nil
	})
}

// EstimateBlockchainHeight returns the blockchain length estimated from peer reports
// Deprecate. Should not need. Just report time of last block
func (vs *Visor) EstimateBlockchainHeight() uint64 {
	var maxLen uint64
	vs.strand("EstimateBlockchainHeight", func() error {
		ourLen := vs.v.HeadBkSeq()
		if len(vs.blockchainHeights) < 2 {
			maxLen = ourLen
			return nil
		}

		for _, seq := range vs.blockchainHeights {
			if maxLen < seq {
				maxLen = seq
			}
		}

		return nil
	})
	return maxLen
}

// PeerBlockchainHeight is a peer's IP address with their reported blockchain height
type PeerBlockchainHeight struct {
	Address string
	Height  uint64
}

// GetPeerBlockchainHeights returns recorded peers' blockchain heights as an array.
func (vs *Visor) GetPeerBlockchainHeights() []PeerBlockchainHeight {
	var peerHeights []PeerBlockchainHeight
	vs.strand("GetPeerBlockchainHeights", func() error {
		if len(vs.blockchainHeights) == 0 {
			return nil
		}

		peerHeights = make([]PeerBlockchainHeight, 0, len(peerHeights))
		for addr, height := range vs.blockchainHeights {
			peerHeights = append(peerHeights, PeerBlockchainHeight{
				Address: addr,
				Height:  height,
			})
		}

		return nil
	})

	return peerHeights
}

// HeadBkSeq returns the head sequence
func (vs *Visor) HeadBkSeq() uint64 {
	var seq uint64
	vs.strand("HeadBkSeq", func() error {
		seq = vs.v.HeadBkSeq()
		return nil
	})
	return seq
}

// ExecuteSignedBlock executes signed block
func (vs *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	return vs.strand("ExecuteSignedBlock", func() error {
		return vs.v.ExecuteSignedBlock(b)
	})
}

// GetSignedBlocksSince returns numbers of signed blocks since seq.
func (vs *Visor) GetSignedBlocksSince(seq uint64, num uint64) ([]coin.SignedBlock, error) {
	var sbs []coin.SignedBlock
	err := vs.strand("GetSignedBlocksSince", func() error {
		var err error
		sbs, err = vs.v.GetSignedBlocksSince(seq, num)
		return err
	})
	return sbs, err
}

// UnconfirmedUnknown returns all unknown transaction hashes
func (vs *Visor) UnconfirmedUnknown(txns []cipher.SHA256) ([]cipher.SHA256, error) {
	var ts []cipher.SHA256

	if err := vs.strand("UnconfirmedUnknown", func() error {
		var err error
		ts, err = vs.v.GetUnconfirmedUnknown(txns)
		return err
	}); err != nil {
		return nil, err
	}

	return ts, nil
}

// UnconfirmedKnown returns all know tansactions
func (vs *Visor) UnconfirmedKnown(hashes []cipher.SHA256) (coin.Transactions, error) {
	var txns coin.Transactions

	if err := vs.strand("UnconfirmedKnown", func() error {
		var err error
		txns, err = vs.v.GetUnconfirmedKnown(hashes)
		return err
	}); err != nil {
		return nil, err
	}

	return txns, nil
}

// InjectTransactionNoBroadcast only try to append transaction into local blockchain, don't broadcast it.
func (vs *Visor) InjectTransactionNoBroadcast(txn coin.Transaction) (bool, error) {
	var known bool
	err := vs.strand("InjectTransactionNoBroadcast", func() error {
		var err error
		known, err = vs.v.InjectTransaction(txn)
		return err
	})
	return known, err
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
func (gbm *GetBlocksMessage) Process(d *Daemon) error {
	// TODO -- we need the sig to be sent with the block, but only the master
	// can sign blocks.  Thus the sig needs to be stored with the block.
	// TODO -- move to either Messages.Config or Visor.Config
	if d.Visor.Config.DisableNetworking {
		return nil
	}
	// Record this as this peer's highest block
	d.Visor.RecordBlockchainHeight(gbm.c.Addr, gbm.LastBlock)
	// Fetch and return signed blocks since LastBlock
	blocks, err := d.Visor.GetSignedBlocksSince(gbm.LastBlock, gbm.RequestedBlocks)
	if err != nil {
		logger.Info("Get signed blocks failed: %v", err)
		return err
	}

	logger.Debug("Got %d blocks since %d", len(blocks), gbm.LastBlock)
	if len(blocks) == 0 {
		return nil
	}
	m := NewGiveBlocksMessage(blocks)

	err = d.Pool.Pool.SendMessage(gbm.c.Addr, m)
	if err != nil {
		logger.Error("Send GiveBlocksMessage to %s failed: %v", gbm.c.Addr, err)
	}
	return err
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
func (gbm *GiveBlocksMessage) Process(d *Daemon) error {
	if d.Visor.Config.DisableNetworking {
		logger.Critical("Visor disabled, ignoring GiveBlocksMessage")
		return nil
	}

	logger.Debug("daemon.GiveBlocksMessage.Process has %d blocks", len(gbm.Blocks))

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
		return nil
	}

	headBkSeq := d.Visor.HeadBkSeq()

	// Announce our new blocks to peers
	m1 := NewAnnounceBlocksMessage(headBkSeq)
	if err := d.Pool.Pool.BroadcastMessage(m1); err != nil {
		logger.Warning("Broadcast AnnounceBlocksMessage failed")
	}

	// Request more blocks.
	m2 := NewGetBlocksMessage(headBkSeq, d.Visor.Config.BlocksResponseCount)
	return d.Pool.Pool.BroadcastMessage(m2)
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
func (abm *AnnounceBlocksMessage) Process(d *Daemon) error {
	if d.Visor.Config.DisableNetworking {
		return nil
	}

	headBkSeq := d.Visor.HeadBkSeq()
	if headBkSeq >= abm.MaxBkSeq {
		return nil
	}

	// TODO: Should this be block get request for current sequence?
	// If client is not caught up, won't attempt to get block
	m := NewGetBlocksMessage(headBkSeq, d.Visor.Config.BlocksResponseCount)
	err := d.Pool.Pool.SendMessage(abm.c.Addr, m)
	if err != nil {
		logger.Error("Send GetBlocksMessage to %s failed: %v", abm.c.Addr, err)
	}
	return err
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
func (atm *AnnounceTxnsMessage) Process(d *Daemon) error {
	if d.Visor.Config.DisableNetworking {
		return nil
	}

	unknown, err := d.Visor.UnconfirmedUnknown(atm.Txns)
	if err != nil {
		return err
	}

	if len(unknown) == 0 {
		return nil
	}

	m := NewGetTxnsMessage(unknown)
	err = d.Pool.Pool.SendMessage(atm.c.Addr, m)
	if err != nil {
		logger.Error("Send GetTxnsMessage to %s failed: %v", atm.c.Addr, err)
	}
	return err
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
func (gtm *GetTxnsMessage) Process(d *Daemon) error {
	if d.Visor.Config.DisableNetworking {
		return nil
	}

	// Locate all txns from the unconfirmed pool
	known, err := d.Visor.UnconfirmedKnown(gtm.Txns)
	if err != nil {
		return err
	}

	if len(known) == 0 {
		return nil
	}

	// Reply to sender with GiveTxnsMessage
	logger.Debug("GetTxnsMessage.Process: %d/%d txns known", len(known), len(gtm.Txns))
	m := NewGiveTxnsMessage(known)
	err = d.Pool.Pool.SendMessage(gtm.c.Addr, m)
	if err != nil {
		logger.Error("Send GiveTxnsMessage to %s failed: %v", gtm.c.Addr, err)
	}
	return err
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
func (gtm *GiveTxnsMessage) Process(d *Daemon) error {
	if d.Visor.Config.DisableNetworking {
		return nil
	}

	hashes := make([]cipher.SHA256, 0, len(gtm.Txns))
	// Update unconfirmed pool with these transactions
	for _, txn := range gtm.Txns {
		// Only announce transactions that are new to us, so that peers can't spam relays
		known, err := d.Visor.InjectTransactionNoBroadcast(txn)
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
	if len(hashes) == 0 {
		return nil
	}

	logger.Debugf("Announce %d transactions", len(hashes))
	m := NewAnnounceTxnsMessage(hashes)
	return d.Pool.Pool.BroadcastMessage(m)
}
