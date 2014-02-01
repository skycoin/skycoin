package daemon

import (
    "github.com/skycoin/gnet"
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
    for _, b := range self.Blocks {
        err := d.Visor.Visor.ExecuteSignedBlock(b)
        if err != nil {
            logger.Info("Failed to execute received block: %v", err)
            // Blocks must be received in order, so if one fails its assumed
            // the rest are failing
            break
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
