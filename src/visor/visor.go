package visor

import (
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/coin"
    "log"
)

var (
    logger = logging.MustGetLogger("skycoin.visor")
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
    Master Keypair
    // // Our personal keys
    // Wallet Wallet
}

// Returns the VisorKeys for the master client
func NewMasterVisorKeys(masterSecret string) VisorKeys {
    secret := coin.SecKeyFromHex(masterSecret)
    return VisorKeys{
        Master: Keypair{
            Secret: secret,
            Public: coin.PubKeyFromSecKey(secret),
        },
        // Wallet: NewWallet(),
    }
}

// Returns the VisorKeys for the normal client
func NewVisorKeys(masterPublic string) VisorKeys {
    pub := coin.PubKeyFromHex(masterPublic)
    return VisorKeys{
        Master: Keypair{
            Public: pub,
        },
        // Wallet: NewWallet(),
    }
}

// Holds references to signed blocks outside of the blockchain
type SignedBlock struct {
    Block coin.Block
    // Block signature
    Sig coin.Sig
}

// Manages known SignedBlocks as received.
// TODO -- support out of order blocks.  This requires a change to the
// message protocol to support ranges similar to bitcoin's locator hashes.
// We also need to keep track of whether a block has been executed so that
// as continuity is established we can execute chains of blocks.
// TODO -- Since we will need to hold blocks that cannot be verified
// immediately against the blockchain, we need to be able to hold multiple
// SignedBlocks per BkSeq, or use hashes as keys.  For now, this is not a
// problem assuming the signed blocks created from master are valid blocks,
// because we can check the signature independently of the blockchain.
type SignedBlocks struct {
    Blocks map[uint64]SignedBlock
    MaxSeq uint64
    MinSeq uint64
}

func NewSignedBlocks() SignedBlocks {
    return SignedBlocks{
        Blocks: make(map[uint64]SignedBlock),
        MaxSeq: 0,
        MinSeq: 0,
    }
}

// Adds a SignedBlock
func (self *SignedBlocks) record(sb *SignedBlock) {
    seq := sb.Block.Header.BkSeq
    self.Blocks[seq] = *sb
    if seq > self.MaxSeq {
        self.MaxSeq = seq
    }
    if seq < self.MinSeq {
        self.MinSeq = seq
    }
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
    // Is this the master blockchain
    IsMaster bool
    // Master & personal keys
    keys         VisorKeys
    blockchain   *coin.Blockchain
    signedBlocks SignedBlocks
}

// Creates a master Visor with a hex encoded secret key
func NewMasterVisor(masterSecret string) *Visor {
    return &Visor{
        IsMaster:     true,
        keys:         NewMasterVisorKeys(masterSecret),
        blockchain:   coin.NewBlockchain(),
        signedBlocks: NewSignedBlocks(),
    }
}

// Creates a normal Visor given a master's public key
func NewVisor(masterPublic string) *Visor {
    return &Visor{
        IsMaster:     false,
        keys:         NewVisorKeys(masterPublic),
        blockchain:   coin.NewBlockchain(),
        signedBlocks: NewSignedBlocks(),
    }
}

// Signs a block for master
func (self *Visor) SignBlock(b coin.Block) (sb SignedBlock, e error) {
    if !self.IsMaster {
        log.Panic("You cannot sign blocks")
    }
    sig, err := coin.SignHash(b.HashHeader(), self.keys.Master.Secret)
    if err != nil {
        e = err
        return
    }
    sb = SignedBlock{
        Block: b,
        Sig:   sig,
    }
    return
}

// Adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (self *Visor) ExecuteSignedBlock(b SignedBlock) error {
    err := self.verifySignedBlock(&b)
    if err != nil {
        return err
    }
    err = self.blockchain.ExecuteBlock(b.Block)
    if err != nil {
        return err
    }
    self.signedBlocks.record(&b)
    return nil
}

// Returns N signed blocks more recent than Seq. Returns nil if no blocks
func (self *Visor) GetSignedBlocksSince(seq uint64, ct int) []SignedBlock {
    if seq < self.signedBlocks.MinSeq {
        seq = self.signedBlocks.MinSeq
    }
    if seq >= self.signedBlocks.MaxSeq {
        return nil
    }
    blocks := make([]SignedBlock, 0, ct)
    for i := seq; i < self.signedBlocks.MaxSeq; i++ {
        if b, exists := self.signedBlocks.Blocks[i]; exists {
            blocks = append(blocks, b)
        }
    }
    if len(blocks) == 0 {
        return nil
    } else {
        return blocks
    }
}

// Returns the highest BkSeq we know
func (self *Visor) MostRecentBkSeq() uint64 {
    return self.signedBlocks.MaxSeq
}

// Returns an error if the coin.Sig is not valid for the coin.Block
func (self *Visor) verifySignedBlock(b *SignedBlock) error {
    return coin.VerifySignature(self.keys.Master.Public, b.Sig,
        b.Block.HashHeader())
}
