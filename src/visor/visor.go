package visor

import (
    "encoding/hex"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "log"
)

// Holds a Secret/Public keypair
type Keypair struct {
    Secret coin.SecKey
    Public coin.PubKey
}

// Holds the master and personal keys
type VisorKeys struct {
    // The master server's key.  The Secret will be empty unless running as
    // a master instance
    MasterKey Keypair
    // Our personal keys
    Keys []Keypair
}

// Returns the VisorKeys for the Master server
func NewMasterVisorKeys(hexKey string) VisorKeys {
    secret := coin.SecKeyFromHex(hexKey)
    return VisorKeys{
        MasterKey: Keypair{
            Secret: secret,
            Public: coin.PubKeyFromSecKey(secret),
        },
        Keys: make([]Keypair, 0),
    }
}

func NewVisorKeys(pubMaster coin.PubKey) VisorKeys {
    return VisorKeys{
        MasterKey: Keypair{
            Public: pubMaster,
        },
        Keys: make([]Keypair, 0),
    }
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
    // Is this the master blockchain
    IsMaster bool
    // Master & personal keys
    Keys       VisorKeys
    Blockchain *coin.Blockchain
}

func NewMasterVisor(hexKey string) *Blockchain {
    return &Blockchain{
        IsMaster:   true,
        Keys:       NewMasterVisorKeys(hexKey),
        Blockchain: coin.NewBlockchain(),
    }
}

func NewVisor(c VisorKeys) *Blockchain {
    return &Blockchain{
        IsMaster:   false,
        Config:     NewVisorKeys(),
        Blockchain: coin.NewBlockchain(),
    }
}

// Returns an error if the coin.Sig is not valid for the coin.Block
func (self *Visor) VerifySignedBlock(b coin.Block, s coin.Sig) error {
    return coin.VerifySignature(self.Keys.Master.Public, s, b)
}

func (self *Visor) SignBlock(b coin.Bloc) (coin.Sig, error) {
    if !self.IsMaster {
        log.Panic("You cannot sign blocks")
    }
    return coin.SignHash(block.HashHeader(), self.Keys.Master.Secret)
}

// Adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (self *Visor) ExecuteBlock(b coin.Block, s coin.Sig) error {
    err := self.VerifySignedBlock(b, s)
    if err != nil {
        return err
    }
    err = self.Blockchain.ExecuteBlock(block)
    if err != nil {
        return err
    }
    return nil
}
