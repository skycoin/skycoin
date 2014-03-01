package daemon

import (
    "errors"
    "fmt"
    "github.com/skycoin/gnet"
    //"github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    //"github.com/skycoin/skycoin/src/Sync"
    "sort"
    "time"
)


type SyncConfig struct {
    //Config Sync.SyncConfig
    // Disabled the Sync completely
    Disabled bool
    // Location of master keys
    //MasterKeysFile string
    // How often to request blocks from peers
    BlocksRequestRate time.Duration
    // How often to announce our blocks to peers
    BlocksAnnounceRate time.Duration
    // How many blocks to respond with to a GetBlocksMessage
    BlocksResponseCount uint64
    // How often to rebroadcast txns that we are a party to
    TransactionRebroadcastRate time.Duration
}

func NewSyncConfig() SyncConfig {
    return SyncConfig{
        //Config:                     Sync.NewSyncConfig(),
        Disabled:                   false,
        //MasterKeysFile:             "",
        BlocksRequestRate:          time.Minute * 5,
        BlocksAnnounceRate:         time.Minute * 15,
        BlocksResponseCount:        20,
        TransactionRebroadcastRate: time.Minute * 5,
    }
}

type Sync struct {
    Config SyncConfig
    //Sync  *Sync.Sync
    // Peer-reported blockchain length.  Use to estimate download progress
    
    //blockchainLengths map[string]uint64

    TxReplicator BlobReplicator //transactions are blobs with flood replication
}

func NewSync(c SyncConfig) *Sync {
    var v *Sync.Sync = nil
    if !c.Disabled {
        v = Sync.NewSync(c.Config)
    }
    return &Sync{
        Config:            c,
        Sync:             v,
        //blockchainLengths: make(map[string]uint64),
    }
}

//save peers to disc?
func (self *Sync) Shutdown() {

}

// Sends a GetBlocksMessage to all connections
//
func (self *Sync) RequestBlocks(pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewGetBlocksMessage(self.Sync.MostRecentBkSeq())
    pool.Pool.BroadcastMessage(m)
}

// Sends an AnnounceBlocksMessage to all connections
func (self *Sync) AnnounceBlocks(pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewAnnounceBlocksMessage(self.Sync.MostRecentBkSeq())
    pool.Pool.BroadcastMessage(m)
}

// Sends a GetBlocksMessage to one connected address
func (self *Sync) RequestBlocksFromAddr(pool *Pool, addr string) error {
    if self.Config.Disabled {
        return errors.New("Sync disabled")
    }
    m := NewGetBlocksMessage(self.Sync.MostRecentBkSeq())
    c := pool.Pool.Addresses[addr]
    if c == nil {
        return fmt.Errorf("Tried to send GetBlocksMessage to %s, but we're "+
            "not connected", addr)
    }
    pool.Pool.SendMessage(c, m)
    return nil
}

// Updates internal state when a connection disconnects
func (self *Sync) RemoveConnection(addr string) {
    //delete(self.blockchainLengths, addr)
}

// Saves a peer-reported blockchain length
func (self *Sync) recordBlockchainLength(addr string, bkLen uint64) {
    //self.blockchainLengths[addr] = bkLen
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
    // TODO -- move 20 to either Messages.Config or Sync.Config
    if d.Sync.Config.Disabled {
        return
    }
    // Record this as this peer's highest block
    d.Sync.recordBlockchainLength(self.c.Conn.Addr(), self.LastBlock)
    // Fetch and return signed blocks since LastBlock
    blocks := d.Sync.Sync.GetSignedBlocksSince(self.LastBlock,
        d.Sync.Config.BlocksResponseCount)
    logger.Debug("Got %d blocks since %d", len(blocks), self.LastBlock)
    if len(blocks) == 0 {
        return
    }
    m := NewGiveBlocksMessage(blocks)
    d.Pool.Pool.SendMessage(self.c.Conn, m)
}

// Sent in response to GetBlocksMessage, or unsolicited
type GiveBlocksMessage struct {
    Blocks []Sync.SignedBlock
    c      *gnet.MessageContext `enc:"-"`
}

func NewGiveBlocksMessage(blocks []Sync.SignedBlock) *GiveBlocksMessage {
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
    logger.Critical("Sync disabled, ignoring GiveBlocksMessage")
    if d.Sync.Config.Disabled {
        logger.Critical("Sync disabled, ignoring GiveBlocksMessage")
        return
    }
    processed := 0
    maxSeq := d.Sync.Sync.MostRecentBkSeq()
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
        err := d.Sync.Sync.ExecuteSignedBlock(b)
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
    m := NewAnnounceBlocksMessage(d.Sync.Sync.MostRecentBkSeq())
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
    if d.Sync.Config.Disabled {
        return
    }
    bkSeq := d.Sync.Sync.MostRecentBkSeq()
    if bkSeq >= self.MaxBkSeq {
        return
    }
    m := NewGetBlocksMessage(bkSeq)
    d.Pool.Pool.SendMessage(self.c.Conn, m)
}

// Sends a signed block to all connections.
// Should only send if person requests
func (self *Sync) broadcastBlock(sb Sync.SignedBlock, pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewGiveBlocksMessage([]Sync.SignedBlock{sb})
    pool.Pool.BroadcastMessage(m)
}

/*
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
*/
