package daemon

import (
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
)

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
    blocks := d.Visor.GetSignedBlocksSince(self.LastBlock, 20)
    if blocks == nil {
        return
    }
    m := NewGiveBlocksMessage(blocks)
    err := d.Pool.Pool.Dispatcher.SendMessage(self.c.Conn, m)
    if err != nil {
        logger.Warning("Failed to send GiveBlocksMessage: %v", err)
    }
}

type BlockMessage struct {
    Block coin.Block
    // TODO -- its not clear where to store the Sig.  The Sig is created
    // by the master server signing the block, but for others to relay it
    // they need to save it.
    Sig coin.Sig
}

// Sent in response to GetBlocksMessage, or unsolicited
type GiveBlocksMessage struct {
    Blocks []BlockMessage
    c      *gnet.MessageContext `enc:"-"`
}

func NewGiveBlocksMessage(blocks []visor.SignedBlock) *GiveBlocksMessage {
    bms := make([]BlockMessage, len(blocks))
    for _, b := range blocks {
        bms = append(bms, BlockMessage{
            Block: b.Block,
            Sig:   b.Sig,
        })
    }
    return &GiveBlocksMessage{
        Blocks: bms,
    }
}

func (self *GiveBlocksMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GiveBlocksMessage) Process(d *Daemon) {
    for _, b := range self.Blocks {
        err := d.Visor.ExecuteBlock(b)
        if err != nil {
            log.Info("Failed to execute received block: %v", err)
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

func (self *AnnounceBlocksMessage) Handle(mc *gnet.messageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *AnnounceBlocksMessage) Process(d *Daemon) {
    bkSeq := self.Visor.MostRecentBkSeq()
    if bkSeq >= self.MaxBkSeq {
        return
    }
    m := NewGetBlocksMessage(bkSeq)
    err := d.Pool.Pool.Dispatcher.SendMessage(self.c.Conn, m)
    if err != nil {
        logger.Warning("Failed to send GetBlocksMessage: %v", err)
    }
}
