package visor

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "log"
)

var (
    logger = logging.MustGetLogger("skycoin.visor")
)

type WalletEntry struct {
    Address coin.Address
    Public  coin.PubKey
    Secret  coin.SecKey
}

func WalletEntryFromReadable(w *ReadableWalletEntry) WalletEntry {
    // Wallet entries are shared as a form of identification, the secret key
    // is not required
    var s coin.SecKey
    if w.Secret != "" {
        s = coin.SecKeyFromHex(w.Secret)
    }
    return WalletEntry{
        Address: coin.DecodeBase58Address(w.Address),
        Public:  coin.PubKeyFromHex(w.Public),
        Secret:  s,
    }
}

// Checks that the public key is derivable from the secret key if present,
// and that the public key is associated with the address
func (self *WalletEntry) Verify(isMaster bool) error {
    var emptySecret coin.SecKey
    if self.Secret == emptySecret {
        if isMaster {
            return errors.New("WalletEntry is master, but has no secret key")
        }
    } else {
        if coin.PubKeyFromSecKey(self.Secret) != self.Public {
            return errors.New("Invalid public key for secret key")
        }
    }
    return self.Address.Verify(self.Public)
}

type ReadableWalletEntry struct {
    Address string `json:"address"`
    Public  string `json:"public_key"`
    Secret  string `json:"secret_key"`
}

func NewReadableWalletEntry(w *WalletEntry) ReadableWalletEntry {
    return ReadableWalletEntry{
        Address: w.Address.String(),
        Public:  w.Public.Hex(),
        Secret:  w.Secret.Hex(),
    }
}

// Loads a WalletEntry from filename, where the file contains a
// ReadableWalletEntry
func LoadWalletEntry(filename string) (WalletEntry, error) {
    w := &ReadableWalletEntry{}
    err := util.LoadJSON(filename, w)
    if err != nil {
        return WalletEntry{}, err
    } else {
        return WalletEntryFromReadable(w), nil
    }
}

// Holds the master and personal keys
type VisorKeys struct {
    // The master server's key.  The Secret will be empty unless running as
    // a master instance
    Master WalletEntry
    // // Our personal keys
    // Wallet Wallet
}

func NewVisorKeys(master WalletEntry) VisorKeys {
    return VisorKeys{
        Master: master,
        // TODO -- use a deterministic wallet.  However, how do we know
        // how many addresses we need to generate from the deterministic
        // wallet? e.g. user creates 10,000 addresses with it, has balance on
        // half of them including the 10,000th, loses wallet db and has to
        // recreate from seed
        // Wallet: NewWallet(),
    }
}

type SignedBlock struct {
    Block coin.Block
    Sig   coin.Sig
}

// Manages known BlockSigs as received.
// TODO -- support out of order blocks.  This requires a change to the
// message protocol to support ranges similar to bitcoin's locator hashes.
// We also need to keep track of whether a block has been executed so that
// as continuity is established we can execute chains of blocks.
// TODO -- Since we will need to hold blocks that cannot be verified
// immediately against the blockchain, we need to be able to hold multiple
// BlockSigs per BkSeq, or use hashes as keys.  For now, this is not a
// problem assuming the signed blocks created from master are valid blocks,
// because we can check the signature independently of the blockchain.
type BlockSigs struct {
    Sigs   map[uint64]coin.Sig
    MaxSeq uint64
    MinSeq uint64
}

func NewBlockSigs() BlockSigs {
    return BlockSigs{
        Sigs:   make(map[uint64]coin.Sig),
        MaxSeq: 0,
        MinSeq: 0,
    }
}

// Adds a SignedBlock
func (self *BlockSigs) record(sb *SignedBlock) {
    seq := sb.Block.Header.BkSeq
    self.Sigs[seq] = sb.Sig
    if seq > self.MaxSeq {
        self.MaxSeq = seq
    }
    if seq < self.MinSeq {
        self.MinSeq = seq
    }
}

// Holds unconfirmed transactions
type UnconfirmedTxnPool struct {
    Txns map[coin.SHA256]coin.Transaction
}

func NewUnconfirmedTxnPool() *UnconfirmedTxnPool {
    return &UnconfirmedTxnPool{
        Txns: make(map[coin.SHA256]coin.Transaction),
    }
}

// Returns txn hashes with known ones removed
func (self *UnconfirmedTxnPool) FilterKnown(txns []coin.SHA256) []coin.SHA256 {
    unknown := make([]coin.SHA256, 0)
    for _, h := range txns {
        _, known := self.Txns[h]
        if !known {
            unknown = append(unknown, h)
        }
    }
    return unknown
}

// Returns all known coin.Transactions from the pool, given hashes to select
func (self *UnconfirmedTxnPool) GetKnown(txns []coin.SHA256) []coin.Transaction {
    known := make([]coin.Transaction, 0)
    for _, h := range txns {
        txn, unknown := self.Txns[h]
        if !unknown {
            known = append(known, txn)
        }
    }
    return known
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
    // Is this the master blockchain
    IsMaster bool
    // Master & personal keys
    keys            VisorKeys
    blockchain      *coin.Blockchain
    blockSigs       BlockSigs
    UnconfirmedTxns *UnconfirmedTxnPool
}

// Creates a normal Visor given a master's public key
func NewVisor(master WalletEntry, isMaster bool) Visor {
    logger.Debug("Creating new visor")
    if isMaster {
        logger.Debug("Visor is master")
    }
    err := master.Verify(isMaster)
    if err != nil {
        log.Panicf("Invalid master wallet entry: %v", err)
    }
    return Visor{
        IsMaster:        isMaster,
        keys:            NewVisorKeys(master),
        blockchain:      coin.NewBlockchain(master.Address),
        blockSigs:       NewBlockSigs(),
        UnconfirmedTxns: NewUnconfirmedTxnPool(),
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
    // TODO -- save them even if out of order, and execute later
    // But make sure all prechecking as possible is done
    // TODO -- check if bitcoin allows blocks to be receiving out of order
    self.blockSigs.record(&b)
    // Remove the transactions in the Block from the unconfirmed pool
    for _, tx := range b.Block.Body.Transactions {
        delete(self.UnconfirmedTxns.Txns, tx.Header.Hash)
    }
    return nil
}

// Returns N signed blocks more recent than Seq. Returns nil if no blocks
func (self *Visor) GetSignedBlocksSince(seq uint64, ct int) []SignedBlock {
    if seq < self.blockSigs.MinSeq {
        seq = self.blockSigs.MinSeq
    }
    if seq >= self.blockSigs.MaxSeq {
        return nil
    }
    blocks := make([]SignedBlock, 0, ct)
    for i := seq; i < self.blockSigs.MaxSeq; i++ {
        if sig, exists := self.blockSigs.Sigs[i]; exists {
            blocks = append(blocks, SignedBlock{
                Sig:   sig,
                Block: self.blockchain.Blocks[i],
            })
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
    return self.blockSigs.MaxSeq
}

// Records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
func (self *Visor) RecordTxn(txn coin.Transaction) error {
    if err := txn.Verify(); err != nil {
        return err
    }
    if err := self.blockchain.VerifyTransaction(txn); err != nil {
        return err
    }
    self.UnconfirmedTxns.Txns[txn.Header.Hash] = txn
    return nil
}

// Returns an error if the coin.Sig is not valid for the coin.Block
func (self *Visor) verifySignedBlock(b *SignedBlock) error {
    return coin.VerifySignature(self.keys.Master.Public, b.Sig,
        b.Block.HashHeader())
}
