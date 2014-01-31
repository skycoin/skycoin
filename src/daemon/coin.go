package daemon

import (
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
)

// Communication layer for the coin pkg

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
    // TODO
}

type blockMessage struct {
    Block coin.Block
    Sig   coin.Sig
}

type GiveBlocksMessage struct {
    Blocks []blockMessage
    c      *gnet.MessageContext `enc:"-"`
}

func NetGiveBlocksMessage(blocks []blockMessage) *GiveBlocksMessage {
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
    // TODO -- where is the global blockchain stored?
    // The blockchain needs to either be a global in the daemon
    // or passed into the daemon from the controlling program
    // and then passed into every Process() call
}
