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


/*
Sync is useless
*/

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
    //TransactionRebroadcastRate time.Duration
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
