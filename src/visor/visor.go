package visor

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/coin"
    "log"
    "os"
    "time"
)

var (
    logger = logging.MustGetLogger("skycoin.visor")
)

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

// Configuration parameters for the Visor
type VisorConfig struct {
    // Is this the master blockchain
    IsMaster bool
    // Is allowed to create transactions
    CanSpend bool
    // Wallet file location
    WalletFile string
    // Minimum number of addresses to keep in the wallet
    WalletSizeMin int
    // Use test network addresses
    TestNetwork bool
}

func NewVisorConfig() VisorConfig {
    return VisorConfig{
        IsMaster:      false,
        CanSpend:      true,
        TestNetwork:   true,
        WalletFile:    "",
        WalletSizeMin: 100,
    }
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
    Config VisorConfig
    // Unconfirmed transactions, held for relay until we get block confirmation
    UnconfirmedTxns *UnconfirmedTxnPool
    // Wallet holding our keys for spending
    Wallet *Wallet
    // Master & personal keys
    keys       VisorKeys
    blockchain *coin.Blockchain
    blockSigs  BlockSigs
}

// Creates a normal Visor given a master's public key
func NewVisor(c VisorConfig, master WalletEntry) *Visor {
    logger.Debug("Creating new visor")
    if c.IsMaster {
        logger.Debug("Visor is master")
    }
    err := master.Verify(c.IsMaster)
    if err != nil {
        log.Panicf("Invalid master wallet entry: %v", err)
    }

    wallet := NewWallet()
    if c.WalletFile != "" {
        err := wallet.Load(c.WalletFile)
        if os.IsNotExist(err) {
            logger.Info("Wallet file \"%s\" does not exist", c.WalletFile)
        } else {
            log.Panicf("Failed to load wallet file: %v", err)
        }
    }
    wallet.Populate(c.WalletSizeMin)
    if c.WalletFile != "" {
        err := wallet.Save(c.WalletFile)
        if err != nil {
            log.Panicf("Failed to save wallet file to \"%s\": ", c.WalletFile,
                err)
        }
    }

    return &Visor{
        Config:          c,
        keys:            NewVisorKeys(master),
        blockchain:      coin.NewBlockchain(master.Address),
        blockSigs:       NewBlockSigs(),
        UnconfirmedTxns: NewUnconfirmedTxnPool(),
        Wallet:          wallet,
    }
}

// Signs a block for master
func (self *Visor) SignBlock(b coin.Block) (sb SignedBlock, e error) {
    if !self.Config.IsMaster {
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

// Creates a Transaction spending coins and hours from our coins
// TODO -- handle txn fees.  coin.Transaciton does not implement fee support
func (self *Visor) Spend(amt Balance,
    dest coin.Address) (coin.Transaction, error) {
    var txn coin.Transaction
    if !self.Config.CanSpend {
        return txn, errors.New("Spending disabled")
    }
    needed := amt
    // needed = needed.Add(fee)
    t := uint64(time.Now().UTC().Unix())
    auxs := self.blockchain.Unspent.AllForAddresses(self.Wallet.GetAddresses())
    for a, uxs := range auxs {
        entry, exists := self.Wallet.GetEntry(a)
        if !exists {
            log.Panic("On second thought, the wallet entry does not exist")
        }
        for _, ux := range uxs {
            if needed.IsZero() {
                break
            }
            b := NewBalance(ux.Body.Coins, ux.CoinHours(t))
            if needed.GreaterThanOrEqual(b) {
                needed = needed.Sub(b)
                txn.PushInput(ux.Hash(), entry.Secret)
            } else {
                change := b.Sub(needed)
                needed = needed.Sub(needed)
                txn.PushInput(ux.Hash(), entry.Secret)
                txn.PushOutput(ux.Body.Address, change.Coins, change.Hours)
            }
        }
    }

    txn.PushOutput(dest, amt.Coins, amt.Hours)

    if needed.IsZero() {
        return txn, nil
    } else {
        return txn, errors.New("Not enough coins or hours")
    }
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

// Returns the total balance for addresses in the Wallet
func (self *Visor) TotalBalance() Balance {
    return self.Wallet.TotalBalance(self.blockchain.Unspent)
}

// Returns the balance for a single address in the Wallet
func (self *Visor) Balance(a coin.Address) Balance {
    return self.Wallet.Balance(self.blockchain.Unspent, a)
}

// Returns an error if the coin.Sig is not valid for the coin.Block
func (self *Visor) verifySignedBlock(b *SignedBlock) error {
    return coin.VerifySignature(self.keys.Master.Public, b.Sig,
        b.Block.HashHeader())
}
