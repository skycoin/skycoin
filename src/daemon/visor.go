package daemon

import (
	"errors"
	"fmt"
	"sort"
	"time"

	//"github.com/skycoin/skycoin/src/aether/gnet"
	"github.com/skycoin/skycoin/src/aether/gnet"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor"
	//"github.com/skycoin/skycoin/src/wallet"
)

/*
Visor should not be duplicated
- this should be pushed into /src/visor
*/

type VisorConfig struct {
	Config visor.VisorConfig
	// Disabled the visor completely
	Disabled bool
	// How often to request blocks from peers
	BlocksRequestRate time.Duration
	// How often to announce our blocks to peers
	BlocksAnnounceRate time.Duration
	// How many blocks to respond with to a GetBlocksMessage
	BlocksResponseCount uint64
}

func NewVisorConfig() VisorConfig {
	return VisorConfig{
		Config:              visor.NewVisorConfig(),
		Disabled:            false,
		BlocksRequestRate:   time.Minute * 5,
		BlocksAnnounceRate:  time.Minute * 15,
		BlocksResponseCount: 20,
	}
}

type Visor struct {
	Config VisorConfig
	Visor  *visor.Visor
	// Peer-reported blockchain length.  Use to estimate download progress
	blockchainLengths map[string]uint64
}

func NewVisor(c VisorConfig) *Visor {
	var v *visor.Visor = nil
	if !c.Disabled {
		v = visor.NewVisor(c.Config)
	}
	return &Visor{
		Config:            c,
		Visor:             v,
		blockchainLengths: make(map[string]uint64),
	}
}

// Closes the Wallet, saving it to disk
func (self *Visor) Shutdown() {
	if self.Config.Disabled {
		return
	}
	// Save the wallet
	/*
		errs := self.Visor.SaveWallets()
		if len(errs) == 0 {
			logger.Info("Saved wallets")
		} else {
			logger.Critical("Failed to save wallets: %v", errs)
		}
	*/

	bcFile := self.Config.Config.BlockchainFile
	err := self.Visor.SaveBlockchain()
	if err == nil {
		logger.Info("Saved blockchain to \"%s\"", bcFile)
	} else {
		logger.Critical("Failed to save blockchain to \"%s\"", bcFile)
	}
	bsFile := self.Config.Config.BlockSigsFile
	err = self.Visor.SaveBlockSigs()
	if err == nil {
		logger.Info("Saved block sigs to \"%s\"", bsFile)
	} else {
		logger.Critical("Failed to save block sigs to \"%s\"", bsFile)
	}
}

// Checks unconfirmed txns against the blockchain and purges ones too old
func (self *Visor) RefreshUnconfirmed() {
	if self.Config.Disabled {
		return
	}
	self.Visor.RefreshUnconfirmed()
}

// Sends a GetBlocksMessage to all connections
func (self *Visor) RequestBlocks(pool *Pool) {
	if self.Config.Disabled {
		return
	}
	m := NewGetBlocksMessage(self.Visor.MostRecentBkSeq())
	pool.Pool.BroadcastMessage(m)
}

// Sends an AnnounceBlocksMessage to all connections
func (self *Visor) AnnounceBlocks(pool *Pool) {
	if self.Config.Disabled {
		return
	}
	m := NewAnnounceBlocksMessage(self.Visor.MostRecentBkSeq())
	pool.Pool.BroadcastMessage(m)
}

// Sends a GetBlocksMessage to one connected address
func (self *Visor) RequestBlocksFromAddr(pool *Pool, addr string) error {
	if self.Config.Disabled {
		return errors.New("Visor disabled")
	}
	m := NewGetBlocksMessage(self.Visor.MostRecentBkSeq())
	c := pool.Pool.Addresses[addr]
	if c == nil {
		return fmt.Errorf("Tried to send GetBlocksMessage to %s, but we're "+
			"not connected", addr)
	}
	pool.Pool.SendMessage(c, m)
	return nil
}

// Sets all txns as announced
func (self *Visor) SetTxnsAnnounced(txns []cipher.SHA256) {
	now := util.Now()
	for _, h := range txns {
		self.Visor.Unconfirmed.SetAnnounced(h, now)
	}
}

// Sends a signed block to all connections.
// TODO: deprecate, should only send to clients that request by hash
func (self *Visor) broadcastBlock(sb visor.SignedBlock, pool *Pool) {
	if self.Config.Disabled {
		return
	}
	m := NewGiveBlocksMessage([]visor.SignedBlock{sb})
	pool.Pool.BroadcastMessage(m)
}

// Broadcasts a single transaction to all peers.
func (self *Visor) BroadcastTransaction(t coin.Transaction, pool *Pool) {
	if self.Config.Disabled {
		return
	}
	m := NewGiveTxnsMessage(coin.Transactions{t})
	logger.Debug("Broadcasting GiveTxnsMessage to %d conns",
		len(pool.Pool.Pool))
	pool.Pool.BroadcastMessage(m)
}

// Creates a spend transaction and broadcasts it to the network
// Spend is replaced with transaction injection

/*
func (self *Visor) Spend(walletID wallet.WalletID, amt wallet.Balance,
	fee uint64, dest cipher.Address, pool *Pool) (coin.Transaction, error) {
	if self.Config.Disabled {
		return coin.Transaction{}, errors.New("Visor disabled")
	}
	logger.Info("Attempting to send %d coins, %d hours to %s with %d fee",
		amt.Coins, amt.Hours, dest.String(), fee)
	txn, err := self.Visor.Spend(walletID, amt, fee, dest)
	if err != nil {
		return txn, err
	}
	err, _ = self.Visor.InjectTxn(txn)
	if err == nil {
		self.BroadcastTransaction(txn, pool)
	}
	return txn, err
}
*/

//move into visor
func (self *Visor) InjectTransaction(txn coin.Transaction, pool *Pool) (coin.Transaction, error) {

	//logger.Info("Attempting to send %d coins, %d hours to %s with %d fee",
	//	amt.Coins, amt.Hours, dest.String(), fee)
	//txn, err := self.Visor.Spend(walletID, amt, fee, dest)
	//if err != nil {
	//	return txn, err
	//}

	err := txn.Verify()

	if err != nil {
		return txn, errors.New("Transaction Verification Failed")
	}

	err, _ = self.Visor.InjectTxn(txn)
	if err == nil {
		self.BroadcastTransaction(txn, pool)
	}
	return txn, err
}

// Resends a known UnconfirmedTxn.
func (self *Visor) ResendTransaction(h cipher.SHA256, pool *Pool) {
	if self.Config.Disabled {
		return
	}
	if ut, ok := self.Visor.Unconfirmed.Txns[h]; ok {
		self.BroadcastTransaction(ut.Txn, pool)
	}
	return
}

// Creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a master chain.  Returns creation error and
// whether it was published or not
func (self *Visor) CreateAndPublishBlock(pool *Pool) error {
	if self.Config.Disabled {
		return errors.New("Visor disabled")
	}
	sb, err := self.Visor.CreateAndExecuteBlock()
	if err != nil {
		return err
	}
	self.broadcastBlock(sb, pool)
	return nil
}

// Updates internal state when a connection disconnects
func (self *Visor) RemoveConnection(addr string) {
	delete(self.blockchainLengths, addr)
}

// Saves a peer-reported blockchain length
func (self *Visor) recordBlockchainLength(addr string, bkLen uint64) {
	self.blockchainLengths[addr] = bkLen
}

// Returns the blockchain length estimated from peer reports
func (self *Visor) EstimateBlockchainLength() uint64 {
	ourLen := self.Visor.MostRecentBkSeq() + 1
	if len(self.blockchainLengths) < 2 {
		return ourLen
	}
	lengths := make(BlockchainLengths, len(self.blockchainLengths))
	i := 0
	for _, seq := range self.blockchainLengths {
		lengths[i] = seq
		i++
	}
	sort.Sort(lengths)
	median := len(lengths) / 2
	var val uint64 = 0
	if len(lengths)%2 == 0 {
		val = (lengths[median] + lengths[median-1]) / 2
	} else {
		val = lengths[median]
	}
	if val < ourLen {
		return ourLen
	} else {
		return val
	}
}

// Communication layer for the coin pkg

// Sent to request blocks since LastBlock
type GetBlocksMessage struct {
	LastBlock uint64
	c         *gnet.MessageContext `enc:"-"`
}

func NewGetBlocksMessage(lastBlock uint64) *GetBlocksMessage {
	return &GetBlocksMessage{
		LastBlock: lastBlock,
	}
}

func (self *GetBlocksMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetBlocksMessage) Process(d *Daemon) {
	// TODO -- we need the sig to be sent with the block, but only the master
	// can sign blocks.  Thus the sig needs to be stored with the block.
	// TODO -- move 20 to either Messages.Config or Visor.Config
	if d.Visor.Config.Disabled {
		return
	}
	// Record this as this peer's highest block
	d.Visor.recordBlockchainLength(self.c.Conn.Addr(), self.LastBlock)
	// Fetch and return signed blocks since LastBlock
	blocks := d.Visor.Visor.GetSignedBlocksSince(self.LastBlock,
		d.Visor.Config.BlocksResponseCount)
	logger.Debug("Got %d blocks since %d", len(blocks), self.LastBlock)
	if len(blocks) == 0 {
		return
	}
	m := NewGiveBlocksMessage(blocks)
	d.Pool.Pool.SendMessage(self.c.Conn, m)
}

// Sent in response to GetBlocksMessage, or unsolicited
type GiveBlocksMessage struct {
	Blocks []visor.SignedBlock
	c      *gnet.MessageContext `enc:"-"`
}

func NewGiveBlocksMessage(blocks []visor.SignedBlock) *GiveBlocksMessage {
	return &GiveBlocksMessage{
		Blocks: blocks,
	}
}

func (self *GiveBlocksMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GiveBlocksMessage) Process(d *Daemon) {
	if d.Visor.Config.Disabled {
		logger.Critical("Visor disabled, ignoring GiveBlocksMessage")
		return
	}
	processed := 0
	maxSeq := d.Visor.Visor.MostRecentBkSeq()
	for _, b := range self.Blocks {
		// To minimize waste when receiving multiple responses from peers
		// we only break out of the loop if the block itself is invalid.
		// E.g. if we request 20 blocks since 0 from 2 peers, and one peer
		// replies with 15 and the other 20, if we did not do this check and
		// the reply with 15 was received first, we would toss the one with 20
		// even though we could process it at the time.
		if b.Block.Head.BkSeq <= maxSeq {
			continue
		}
		err := d.Visor.Visor.ExecuteSignedBlock(b)
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
	logger.Critical("Processed %d/%d blocks", processed, len(self.Blocks))
	if processed == 0 {
		return
	}
	// Announce our new blocks to peers
	m := NewAnnounceBlocksMessage(d.Visor.Visor.MostRecentBkSeq())
	d.Pool.Pool.BroadcastMessage(m)
}

// Tells a peer our highest known BkSeq. The receiving peer can choose
// to send GetBlocksMessage in response
type AnnounceBlocksMessage struct {
	MaxBkSeq uint64
	c        *gnet.MessageContext `enc:"-"`
}

func NewAnnounceBlocksMessage(seq uint64) *AnnounceBlocksMessage {
	return &AnnounceBlocksMessage{
		MaxBkSeq: seq,
	}
}

func (self *AnnounceBlocksMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *AnnounceBlocksMessage) Process(d *Daemon) {
	if d.Visor.Config.Disabled {
		return
	}
	bkSeq := d.Visor.Visor.MostRecentBkSeq()
	if bkSeq >= self.MaxBkSeq {
		return
	}
	m := NewGetBlocksMessage(bkSeq)
	d.Pool.Pool.SendMessage(self.c.Conn, m)
}

type SendingTxnsMessage interface {
	GetTxns() []cipher.SHA256
}

// Tells a peer that we have these transactions
type AnnounceTxnsMessage struct {
	Txns []cipher.SHA256
	c    *gnet.MessageContext `enc:"-"`
}

func NewAnnounceTxnsMessage(txns []cipher.SHA256) *AnnounceTxnsMessage {
	return &AnnounceTxnsMessage{
		Txns: txns,
	}
}

func (self *AnnounceTxnsMessage) GetTxns() []cipher.SHA256 {
	return self.Txns
}

func (self *AnnounceTxnsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *AnnounceTxnsMessage) Process(d *Daemon) {
	if d.Visor.Config.Disabled {
		return
	}
	unknown := d.Visor.Visor.Unconfirmed.FilterKnown(self.Txns)
	if len(unknown) == 0 {
		return
	}
	m := NewGetTxnsMessage(unknown)
	d.Pool.Pool.SendMessage(self.c.Conn, m)
}

type GetTxnsMessage struct {
	Txns []cipher.SHA256
	c    *gnet.MessageContext `enc:"-"`
}

func NewGetTxnsMessage(txns []cipher.SHA256) *GetTxnsMessage {
	return &GetTxnsMessage{
		Txns: txns,
	}
}

func (self *GetTxnsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetTxnsMessage) Process(d *Daemon) {
	if d.Visor.Config.Disabled {
		return
	}
	// Locate all txns from the unconfirmed pool
	// reply to sender with GiveTxnsMessage
	known := d.Visor.Visor.Unconfirmed.GetKnown(self.Txns)
	if len(known) == 0 {
		return
	}
	logger.Debug("%d/%d txns known", len(known), len(self.Txns))
	m := NewGiveTxnsMessage(known)
	d.Pool.Pool.SendMessage(self.c.Conn, m)
}

type GiveTxnsMessage struct {
	Txns coin.Transactions
	c    *gnet.MessageContext `enc:"-"`
}

func NewGiveTxnsMessage(txns coin.Transactions) *GiveTxnsMessage {
	return &GiveTxnsMessage{
		Txns: txns,
	}
}

func (self *GiveTxnsMessage) GetTxns() []cipher.SHA256 {
	return self.Txns.Hashes()
}

func (self *GiveTxnsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GiveTxnsMessage) Process(d *Daemon) {
	if d.Visor.Config.Disabled {
		return
	}
	hashes := make([]cipher.SHA256, 0, len(self.Txns))
	// Update unconfirmed pool with these transactions
	for _, txn := range self.Txns {
		// Only announce transactions that are new to us, so that peers can't
		// spam relays
		if err, known := d.Visor.Visor.InjectTxn(txn); err == nil && !known {
			hashes = append(hashes, txn.Hash())
		} else {
			logger.Warning("Failed to record txn: %v", err)
		}
	}
	// Announce these transactions to peers
	if len(hashes) != 0 {
		m := NewAnnounceTxnsMessage(hashes)
		d.Pool.Pool.BroadcastMessage(m)
	}
}

type BlockchainLengths []uint64

func (self BlockchainLengths) Len() int {
	return len(self)
}

func (self BlockchainLengths) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self BlockchainLengths) Less(i, j int) bool {
	return self[i] < self[j]
}
