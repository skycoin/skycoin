package daemon

import (
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
)

type VisorConfig struct {
    // Is a master visor
    IsMaster bool
    // Location of master keys
    MasterKeysFile string
    // Master public/secret key and genesis address
    MasterKeys visor.WalletEntry
    // Is running on the test network
    TestNetwork bool
}

func NewVisorConfig() VisorConfig {
    return VisorConfig{
        IsMaster:       false,
        MasterKeysFile: "",
        MasterKeys:     visor.WalletEntry{},
        TestNetwork:    true,
    }
}

func (self *VisorConfig) LoadMasterKeys() error {
    keys, err := visor.LoadWalletEntry(self.MasterKeysFile)
    if err != nil {
        return err
    }
    self.MasterKeys = keys
    return nil
}

type Visor struct {
    Config VisorConfig
    Visor  visor.Visor
}

func NewVisor(c VisorConfig) *Visor {
    v := visor.NewVisor(c.MasterKeys, c.IsMaster)
    return &Visor{
        Config: c,
        Visor:  v,
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

    // Announce our new blocks to peers
    if processed == 0 {
        return
    }
    m := NewAnnounceBlocksMessage(d.Visor.Visor.MostRecentBkSeq())
    for _, c := range d.Pool.Pool.Pool {
        err := d.Pool.Pool.Dispatcher.SendMessage(c, m)
        if err != nil {
            logger.Warning("Failed to announce blocks to %s", c.Addr())
        }
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
        d.Visor.Visor.RecordTxn(txn)
    }
}
