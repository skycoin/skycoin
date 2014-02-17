package daemon

import (
    "errors"
    "fmt"
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/skycoin/skycoin/src/visor"
    "sort"
    "time"
)

type VisorConfig struct {
    Config visor.VisorConfig
    // Disabled the visor completely
    Disabled bool
    // Location of master keys
    MasterKeysFile string
    // How often to request blocks from peers
    BlocksRequestRate time.Duration
    // How often to announce our blocks to peers
    BlocksAnnounceRate time.Duration
    // How many blocks to respond with to a GetBlocksMessage
    BlocksResponseCount uint64
    // How often to rebroadcast txns that we are a party to
    TransactionRebroadcastRate time.Duration
}

func NewVisorConfig() VisorConfig {
    return VisorConfig{
        Config:                     visor.NewVisorConfig(),
        Disabled:                   false,
        MasterKeysFile:             "",
        BlocksRequestRate:          time.Minute * 5,
        BlocksAnnounceRate:         time.Minute * 15,
        BlocksResponseCount:        20,
        TransactionRebroadcastRate: time.Minute * 5,
    }
}

func (self *VisorConfig) LoadMasterKeys() {
    if self.Disabled {
        return
    }
    self.Config.MasterKeys = visor.MustLoadWalletEntry(self.MasterKeysFile)
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
    // Save the wallet, as long as we're not a master chain.  Master chains
    // don't have a wallet, they have a single genesis wallet entry which is
    // loaded in a different path
    if !self.Config.Config.IsMaster {
        walletFile := self.Config.Config.WalletFile
        err := self.Visor.SaveWallet()
        if err == nil {
            logger.Info("Saved wallet to \"%s\"", walletFile)
        } else {
            logger.Critical("Failed to save wallet to \"%s\": %v", walletFile, err)
        }
    }
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

// Broadcasts any txn that we are a party to
func (self *Visor) BroadcastOurTransactions(pool *Pool) {
    if self.Config.Disabled {
        return
    }
    since := self.Config.TransactionRebroadcastRate * 2
    since = (since * 9) / 10
    txns := self.Visor.UnconfirmedTxns.GetOldOwnedTransactions(since)
    if len(txns) == 0 {
        return
    }
    hashes := make([]coin.SHA256, len(txns))
    for _, tx := range txns {
        hashes = append(hashes, tx.Txn.Hash())
    }
    m := NewAnnounceTxnsMessage(hashes)
    pool.Pool.BroadcastMessage(m)
}

// Sets all txns as announced
func (self *Visor) SetTxnsAnnounced(txns []coin.SHA256) {
    now := util.Now()
    for _, h := range txns {
        self.Visor.UnconfirmedTxns.SetAnnounced(h, now)
    }
}

// Sends a signed block to all connections.
func (self *Visor) broadcastBlock(sb visor.SignedBlock, pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewGiveBlocksMessage([]visor.SignedBlock{sb})
    pool.Pool.BroadcastMessage(m)
}

// Broadcasts a single transaction to all peers.
func (self *Visor) broadcastTransaction(t coin.Transaction, pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewGiveTxnsMessage([]coin.Transaction{t})
    pool.Pool.BroadcastMessage(m)
}

// Creates a spend transaction and broadcasts it to the network
func (self *Visor) Spend(amt visor.Balance, fee uint64,
    dest coin.Address, pool *Pool) (coin.Transaction, error) {
    if self.Config.Disabled {
        return coin.Transaction{}, errors.New("Visor disabled")
    }
    txn, err := self.Visor.Spend(amt, fee, dest)
    if err != nil {
        return txn, err
    }
    self.broadcastTransaction(txn, pool)
    return txn, self.Visor.RecordTxn(txn)
}

// Resends a known UnconfirmedTxn.
func (self *Visor) ResendTransaction(h coin.SHA256, pool *Pool) {
    if self.Config.Disabled {
        return
    }
    if ut, ok := self.Visor.UnconfirmedTxns.Txns[h]; ok {
        self.broadcastTransaction(ut.Txn, pool)
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
    lengths := make(BlockchainLengths, 0, len(self.blockchainLengths))
    for _, seq := range self.blockchainLengths {
        lengths = append(lengths, seq)
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
            logger.Debug("Added new block %d", b.Block.Head.BkSeq)
            processed++
        } else {
            logger.Info("Failed to execute received block: %v", err)
            // Blocks must be received in order, so if one fails its assumed
            // the rest are failing
            break
        }
    }
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

// Tells a peer that we have these transactions
type AnnounceTxnsMessage struct {
    Txns []coin.SHA256
    c    *gnet.MessageContext `enc:"-"`
}

func NewAnnounceTxnsMessage(txns []coin.SHA256) *AnnounceTxnsMessage {
    return &AnnounceTxnsMessage{
        Txns: txns,
    }
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
    unknown := d.Visor.Visor.UnconfirmedTxns.FilterKnown(self.Txns)
    if len(unknown) == 0 {
        return
    }
    m := NewGetTxnsMessage(unknown)
    d.Pool.Pool.SendMessage(self.c.Conn, m)
}

type GetTxnsMessage struct {
    Txns []coin.SHA256
    c    *gnet.MessageContext `enc:"-"`
}

func NewGetTxnsMessage(txns []coin.SHA256) *GetTxnsMessage {
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
    known := d.Visor.Visor.UnconfirmedTxns.GetKnown(self.Txns)
    if len(known) == 0 {
        return
    }
    logger.Debug("%d/%d txns known", len(known), len(self.Txns))
    m := NewGiveTxnsMessage(known)
    d.Pool.Pool.SendMessage(self.c.Conn, m)
}

type GiveTxnsMessage struct {
    Txns []coin.Transaction
    c    *gnet.MessageContext `enc:"-"`
}

func NewGiveTxnsMessage(txns []coin.Transaction) *GiveTxnsMessage {
    return &GiveTxnsMessage{
        Txns: txns,
    }
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
    hashes := make([]coin.SHA256, 0, len(self.Txns))
    // Update unconfirmed pool with these transactions
    for _, txn := range self.Txns {
        if err := d.Visor.Visor.RecordTxn(txn); err == nil {
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
    t := self[i]
    self[i] = self[j]
    self[j] = t
}
func (self BlockchainLengths) Less(i, j int) bool {
    return self[i] < self[j]
}
