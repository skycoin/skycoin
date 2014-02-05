package daemon

import (
    "errors"
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
    "log"
    "time"
)

type VisorConfig struct {
    Config visor.VisorConfig
    // Location of master keys
    MasterKeysFile string
    // Master public/secret key and genesis address
    MasterKeys visor.WalletEntry
    // How often to request blocks from peers
    BlocksRequestRate time.Duration
    // How often to announce our blocks to peers
    BlocksAnnounceRate time.Duration
}

func NewVisorConfig() VisorConfig {
    return VisorConfig{
        Config:             visor.NewVisorConfig(),
        MasterKeysFile:     "",
        MasterKeys:         visor.WalletEntry{},
        BlocksRequestRate:  time.Minute * 15,
        BlocksAnnounceRate: time.Minute * 30,
    }
}

func (self *VisorConfig) LoadMasterKeys() error {
    keys, err := visor.LoadWalletEntry(self.MasterKeysFile)
    if err != nil {
        return err
    }
    if err := keys.Verify(); err != nil {
        log.Panicf("Invalid master keys: %v", err)
    }
    self.MasterKeys = keys
    return nil
}

type Visor struct {
    Config VisorConfig
    Visor  *visor.Visor
}

func NewVisor(c VisorConfig) *Visor {
    v := visor.NewVisor(c.Config, c.MasterKeys)
    return &Visor{
        Config: c,
        Visor:  v,
    }
}

// Closes the Wallet, saving it to disk
func (self *Visor) Shutdown() {
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
}

// Sends a GetBlocksMessage to all connections
func (self *Visor) RequestBlocks(pool *Pool) {
    m := NewGetBlocksMessage(self.Visor.MostRecentBkSeq())
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    for a, _ := range errs {
        logger.Error("Failed to send GetBlocksMessage to %s\n", a)
    }
}

// Sends an AnnounceBlocksMessage to all connections
func (self *Visor) AnnounceBlocks(pool *Pool) {
    m := NewAnnounceBlocksMessage(self.Visor.MostRecentBkSeq())
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    for a, _ := range errs {
        logger.Error("Failed to send AnnounceBlocksMessage to %s\n", a)
    }
}

// Sends a GetBlocksMessage to one connection
func (self *Visor) RequestBlocksFromConn(pool *Pool, addr string) {
    m := NewGetBlocksMessage(self.Visor.MostRecentBkSeq())
    c := pool.Pool.Addresses[addr]
    if c == nil {
        logger.Warning("Tried to send GetBlocksMessage to %s, but we're "+
            "not connected", addr)
        return
    }
    err := pool.Pool.Dispatcher.SendMessage(c, m)
    if err != nil {
        logger.Error("Failed to send GetBlocksMessage to %s\n", c.Addr())
    }
}

// Sends a signed block to all connections
func (self *Visor) broadcastBlock(sb visor.SignedBlock, pool *Pool) error {
    m := NewGiveBlocksMessage([]visor.SignedBlock{sb})
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    if len(errs) == len(pool.Pool.Pool) {
        return errors.New("Failed to give blocks to anyone")
    } else {
        return nil
    }
}

// Broadcasts a single transaction to all peers
func (self *Visor) broadcastTransaction(t coin.Transaction, pool *Pool) error {
    m := NewGiveTxnsMessage([]coin.Transaction{t})
    errs := pool.Pool.Dispatcher.BroadcastMessage(m)
    if len(errs) == len(pool.Pool.Pool) {
        return errors.New("Failed to give transaction to anyone")
    } else {
        return nil
    }
}

// Creates a spend transaction and broadcasts it to the network
func (self *Visor) Spend(amt visor.Balance, fee uint64,
    dest coin.Address, pool *Pool) (coin.Transaction, error) {
    txn, err := self.Visor.Spend(amt, fee, dest)
    if err != nil {
        return txn, err
    }
    err = self.Visor.RecordTxn(txn)
    if err != nil {
        return txn, err
    }
    err = self.broadcastTransaction(txn, pool)
    return txn, err
}

// Creates a block from unconfirmed transactions and sends it to the network.
// Will panic if not running as a master chain.
func (self *Visor) CreateAndPublishBlock(pool *Pool) error {
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
    blocks := d.Visor.Visor.GetSignedBlocksSince(self.LastBlock, 20)
    if blocks == nil {
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
    if len(self.Blocks) == 0 {
        return
    }
    processed := 0
    for i, b := range self.Blocks {
        err := d.Visor.Visor.ExecuteSignedBlock(b)
        if err != nil {
            logger.Info("Failed to execute received block: %v", err)
            // Blocks must be received in order, so if one fails its assumed
            // the rest are failing
        }
        processed = i + 1
    }

    // Announce our new blocks to peers, if we got any
    if processed > 0 {
        m := NewAnnounceBlocksMessage(d.Visor.Visor.MostRecentBkSeq())
        errs := d.Pool.Pool.Dispatcher.BroadcastMessage(m)
        for a, _ := range errs {
            logger.Warning("Failed to announce blocks to %s", a)
        }
    }

    // Send a new GetBlocksMessage, in case we aren't finished yet
    // This also helps in case we receive out of order - if we got blocks
    // but couldn't insert them because they were not sequential for us,
    // requesting the blocks will allow us to catch up
    bkSeq := d.Visor.Visor.MostRecentBkSeq()
    n := NewGetBlocksMessage(bkSeq)
    errs := d.Pool.Pool.Dispatcher.BroadcastMessage(n)
    for a, _ := range errs {
        logger.Warning("Failed to send GetBlocksMessage to %s", a)
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
    // TODO
    // check if we have these txns already
    // look in unconfirmed pool
    // look in Blockchain (need datastructure for blockchain)
    // if we don't have these txns already, send a GetTxnsMessage

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
    // Update unconfirmed pool with these transactions
    for _, txn := range self.Txns {
        err := d.Visor.Visor.RecordTxn(txn)
        if err != nil {
            logger.Warning("Failed to record txn: %v", err)
        }
    }
}
