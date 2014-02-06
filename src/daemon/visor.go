package daemon

import (
    "errors"
    "fmt"
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
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

func (self *VisorConfig) LoadMasterKeys() error {
    if self.Disabled {
        return nil
    }
    keys, err := visor.MustLoadWalletEntry(self.MasterKeysFile)
    if err != nil {
        return err
    }
    self.Config.MasterKeys = keys
    return nil
}

type Visor struct {
    Config VisorConfig
    Visor  *visor.Visor
}

func NewVisor(c VisorConfig) *Visor {
    var v *visor.Visor = nil
    if !c.Disabled {
        v = visor.NewVisor(c.Config)
    }
    return &Visor{
        Config: c,
        Visor:  v,
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

// Sends a GetBlocksMessage to all connections
func (self *Visor) RequestBlocks(pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewGetBlocksMessage(self.Visor.MostRecentBkSeq())
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    for a, _ := range errs {
        logger.Error("Failed to send GetBlocksMessage to %s\n", a)
    }
}

// Sends an AnnounceBlocksMessage to all connections
func (self *Visor) AnnounceBlocks(pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewAnnounceBlocksMessage(self.Visor.MostRecentBkSeq())
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    for a, _ := range errs {
        logger.Error("Failed to send AnnounceBlocksMessage to %s\n", a)
    }
}

// Sends a GetBlocksMessage to one connection
func (self *Visor) RequestBlocksFromConn(pool *Pool, addr string) error {
    if self.Config.Disabled {
        return nil
    }
    m := NewGetBlocksMessage(self.Visor.MostRecentBkSeq())
    c := pool.Pool.Addresses[addr]
    if c == nil {
        return fmt.Errorf("Tried to send GetBlocksMessage to %s, but we're "+
            "not connected", addr)
    }
    err := pool.Pool.Dispatcher.SendMessage(c, m)
    if err == nil {
        return nil
    } else {
        return fmt.Errorf("Failed to send GetBlocksMessage to %s: %v\n",
            c.Addr(), err)
    }
}

// Broadcasts any txn that we are a party to
func (self *Visor) BroadcastOurTransactions(pool *Pool) {
    if self.Config.Disabled {
        return
    }
    since := (self.Config.TransactionRebroadcastRate * 2) - (time.Second * 30)
    txns := self.Visor.UnconfirmedTxns.GetOwnedTransactionsSince(since)
    hashes := make([]coin.SHA256, len(txns))
    for _, tx := range txns {
        hashes = append(hashes, tx.Txn.Header.Hash)
    }
    m := NewAnnounceTxnsMessages(hashes)
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    if len(errs) != len(pool.Pool.Pool) {
        now := time.Now().UTC()
        for _, tx := range txns {
            tx.Announced = now
        }
    }
}

// Sends a signed block to all connections
func (self *Visor) broadcastBlock(sb visor.SignedBlock, pool *Pool) error {
    if self.Config.Disabled {
        return errors.New("Visor disabled")
    }
    m := NewGiveBlocksMessage([]visor.SignedBlock{sb})
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    if len(errs) == len(pool.Pool.Pool) {
        logger.Warning("Failed to give blocks to anyone")
    }
    return nil
}

// Broadcasts a single transaction to all peers
func (self *Visor) broadcastTransaction(t coin.Transaction, pool *Pool) error {
    if self.Config.Disabled {
        return errors.New("Visor disabled")
    }
    m := NewGiveTxnsMessage([]coin.Transaction{t})
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    if len(errs) == len(pool.Pool.Pool) {
        logger.Warning("Failed to give transaction to anyone")
        return errors.New("Did not broadcast to anyone")
    }
    return nil
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
    didAnnounce := self.broadcastTransaction(txn, pool) == nil
    err = self.Visor.RecordTxn(txn, didAnnounce)
    return txn, err
}

// Creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a master chain.
func (self *Visor) CreateAndPublishBlock(pool *Pool) error {
    if self.Config.Disabled {
        return errors.New("Visor disabled")
    }
    sb, err := self.Visor.CreateBlock()
    if err == nil {
        return self.broadcastBlock(sb, pool)
    } else {
        return err
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
    blocks := d.Visor.Visor.GetSignedBlocksSince(self.LastBlock,
        d.Visor.Config.BlocksResponseCount)
    logger.Debug("Got %d blocks since %d", len(blocks), self.LastBlock)
    if len(blocks) == 0 {
        return
    }
    m := NewGiveBlocksMessage(blocks)
    err := d.Pool.Pool.Dispatcher.SendMessage(self.c.Conn, m)
    if err != nil {
        logger.Warning("Failed to send GiveBlocksMessage: %v", err)
    }
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
    if len(self.Blocks) == 0 {
        return
    }
    processed := 0
    for i, b := range self.Blocks {
        err := d.Visor.Visor.ExecuteSignedBlock(b)
        if err == nil {
            logger.Debug("Added new block %d", b.Block.Header.BkSeq)
        } else {
            logger.Info("Failed to execute received block: %v", err)
            // Blocks must be received in order, so if one fails its assumed
            // the rest are failing
            break
        }
        processed = i + 1
    }
    if processed == 0 {
        return
    }
    // Announce our new blocks to peers
    m := NewAnnounceBlocksMessage(d.Visor.Visor.MostRecentBkSeq())
    errs := d.Pool.Pool.Dispatcher.BroadcastMessage(m)
    for a, _ := range errs {
        logger.Warning("Failed to announce blocks to %s", a)
    }
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
    err := d.Pool.Pool.Dispatcher.SendMessage(self.c.Conn, m)
    if err != nil {
        logger.Warning("Failed to send GetBlocksMessage: %v", err)
    }
}

// Tells a peer that we have these transactions
type AnnounceTxnsMessage struct {
    Txns []coin.SHA256
    c    *gnet.MessageContext `enc:"-"`
}

func NewAnnounceTxnsMessages(txns []coin.SHA256) *AnnounceTxnsMessage {
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
    err := d.Pool.Pool.Dispatcher.SendMessage(self.c.Conn, m)
    if err != nil {
        logger.Warning("Failed to send GetTxnsMessage to %s",
            self.c.Conn.Addr())
    }
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
    m := NewGiveTxnsMessage(known)
    err := d.Pool.Pool.Dispatcher.SendMessage(self.c.Conn, m)
    if err != nil {
        logger.Warning("Failed to send GiveTxnsMessage to %s",
            self.c.Conn.Addr())
    }
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
        err := d.Visor.Visor.RecordTxn(txn, false)
        if err == nil {
            hashes = append(hashes, txn.Header.Hash)
        } else {
            logger.Warning("Failed to record txn: %v", err)
        }
    }
    // Announce these transactions to peers
    if len(hashes) != 0 {
        now := time.Now().UTC()
        m := NewAnnounceTxnsMessages(hashes)
        errs := d.Pool.Pool.Dispatcher.BroadcastMessage(m)
        if len(errs) != len(d.Pool.Pool.Pool) {
            for _, h := range hashes {
                d.Visor.Visor.SetAnnounced(h, now)
            }
        }
    }
}
