package visor

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/encoder"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "io/ioutil"
    "log"
    "os"
    "sort"
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
    // How often new blocks are created by the master
    BlockCreationInterval uint64
    // How often an unconfirmed txn is checked against the blockchain
    UnconfirmedCheckInterval time.Duration
    // How long we'll hold onto an unconfirmed txn
    UnconfirmedMaxAge time.Duration
    // How often to refresh the unconfirmed pool
    UnconfirmedRefreshRate time.Duration
    // Maximum number of transactions per block, when creating
    TransactionsPerBlock int
    // Where the blockchain is saved
    BlockchainFile string
}

func NewVisorConfig() VisorConfig {
    return VisorConfig{
        IsMaster:                 false,
        CanSpend:                 true,
        TestNetwork:              true,
        WalletFile:               "",
        WalletSizeMin:            100,
        BlockCreationInterval:    15,
        UnconfirmedCheckInterval: time.Hour * 2,
        UnconfirmedMaxAge:        time.Hour * 48,
        UnconfirmedRefreshRate:   time.Minute * 30,
        TransactionsPerBlock:     1000, // 1000/15 = 66tps. Bitcoin is 7tps
        BlockchainFile:           "",
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

    wallet := loadWallet(c.WalletFile, c.WalletSizeMin)
    blockchain := loadBlockchain(c.BlockchainFile, master.Address,
        c.BlockCreationInterval)

    return &Visor{
        Config:          c,
        keys:            NewVisorKeys(master),
        blockchain:      blockchain,
        blockSigs:       NewBlockSigs(),
        UnconfirmedTxns: NewUnconfirmedTxnPool(),
        Wallet:          wallet,
    }
}

// Checks unconfirmed txns against the blockchain and purges ones that too old
func (self *Visor) RefreshUnconfirmed() {
    logger.Debug("Refreshing unconfirmed transactions")
    self.UnconfirmedTxns.Refresh(self.blockchain,
        self.Config.UnconfirmedCheckInterval, self.Config.UnconfirmedMaxAge)
}

// Saves the coin.Blockchain to disk
func (self *Visor) SaveBlockchain() error {
    if self.Config.BlockchainFile == "" {
        return errors.New("No BlockchainFile location set")
    } else {
        // TODO -- blockchain file must be forward compatible
        data := encoder.Serialize(self.blockchain)
        return util.SaveBinary(self.Config.BlockchainFile, data, 0644)
    }
}

// Saves the Wallet to disk
func (self *Visor) SaveWallet() error {
    if self.Config.WalletFile == "" {
        return errors.New("No WalletFile location set")
    } else {
        return self.Wallet.Save(self.Config.WalletFile)
    }
}

// Creates and returns a WalletEntry and saves the wallet to disk
func (self *Visor) CreateAddressAndSave() (WalletEntry, error) {
    we := self.Wallet.CreateAddress()
    err := self.SaveWallet()
    if err != nil {
        m := "Failed to save wallet to \"%s\" after creating new address"
        logger.Warning(m, self.Config.WalletFile)
    }
    return we, err
}

// Creates a SignedBlock from pending transactions
func (self *Visor) CreateBlock() (SignedBlock, error) {
    var sb SignedBlock
    if !self.Config.IsMaster {
        return sb, errors.New("Only master chain can create blocks")
    }
    if len(self.UnconfirmedTxns.Txns) == 0 {
        return sb, errors.New("No transactions")
    }
    txns := self.UnconfirmedTxns.RawTxns()
    // TODO -- transactions should be sorted by tx fee
    sort.Sort(txns)
    nTxns := len(txns)
    if nTxns > self.Config.TransactionsPerBlock {
        txns = txns[:self.Config.TransactionsPerBlock]
    }
    b, err := self.blockchain.NewBlockFromTransactions(txns)
    if err != nil {
        return sb, err
    }
    sb, err = self.signBlock(b)
    if err == nil {
        return sb, self.ExecuteSignedBlock(sb)
    } else {
        return sb, err
    }
}

// Creates a Transaction spending coins and hours from our coins
// TODO -- handle txn fees.  coin.Transaciton does not implement fee support
func (self *Visor) Spend(amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    var txn coin.Transaction
    if !self.Config.CanSpend {
        return txn, errors.New("Spending disabled")
    }
    if amt.IsZero() {
        return txn, errors.New("Zero spend amount")
    }
    needed := amt
    needed.Hours += fee
    auxs := self.getAvailableBalances()
    for a, uxs := range auxs {
        entry, exists := self.Wallet.GetEntry(a)
        if !exists {
            log.Panic("On second thought, the wallet entry does not exist")
        }
        for _, ux := range uxs {
            if needed.IsZero() {
                break
            }
            coinHours := ux.CoinHours(self.blockchain.Time())
            b := NewBalance(ux.Body.Coins, coinHours)
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
    txn.UpdateHeader()

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
    self.UnconfirmedTxns.RemoveTransactions(self.blockchain,
        b.Block.Body.Transactions)
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
    return self.UnconfirmedTxns.RecordTxn(self.blockchain, txn)
}

// Returns the balance of the wallet
func (self *Visor) TotalBalance() Balance {
    addrs := self.Wallet.GetAddresses()
    auxs := self.blockchain.Unspent.AllForAddresses(addrs)
    return self.totalBalance(auxs)
}

// Returns the total balance of the wallet including unconfirmed outputs
func (self *Visor) TotalBalancePredicted() Balance {
    auxs := self.getAvailableBalances()
    return self.totalBalance(auxs)
}

// Returns the balance for a single address in the Wallet
func (self *Visor) Balance(a coin.Address) Balance {
    uxs := self.blockchain.Unspent.AllForAddress(a)
    return self.balance(uxs)
}

// Returns the balance for a single address in the Wallet, including
// unconfirmed outputs
func (self *Visor) BalancePredicted(a coin.Address) Balance {
    uxs := self.getAvailableBalance(a)
    return self.balance(uxs)
}

// Computes the total balance for coin.Addresses and their coin.UxOuts
func (self *Visor) totalBalance(auxs coin.AddressUnspents) Balance {
    prevTime := self.blockchain.Time()
    b := NewBalance(0, 0)
    for _, uxs := range auxs {
        for _, ux := range uxs {
            b = b.Add(NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
        }
    }
    return b
}

// Computes the balance for a coin.Address's coin.UxOuts
func (self *Visor) balance(uxs []coin.UxOut) Balance {
    prevTime := self.blockchain.Time()
    b := NewBalance(0, 0)
    for _, ux := range uxs {
        b = b.Add(NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
    }
    return b
}

// Returns the total of known Unspents available to us, and our own
// unconfirmed unspents
func (self *Visor) getAvailableBalances() coin.AddressUnspents {
    addrs := self.Wallet.GetAddresses()
    auxs := self.blockchain.Unspent.AllForAddresses(addrs)
    uauxs := self.UnconfirmedTxns.Unspent.AllForAddresses(addrs)
    return auxs.Merge(uauxs, addrs)
}

// Returns the total of known unspents available for an address, including
// unconfirmed requests
func (self *Visor) getAvailableBalance(a coin.Address) []coin.UxOut {
    auxs := self.blockchain.Unspent.AllForAddress(a)
    uauxs := self.UnconfirmedTxns.Unspent.AllForAddress(a)
    return append(auxs, uauxs...)
}

// Returns an error if the coin.Sig is not valid for the coin.Block
func (self *Visor) verifySignedBlock(b *SignedBlock) error {
    return coin.VerifySignature(self.keys.Master.Public, b.Sig,
        b.Block.HashHeader())
}

// Signs a block for master
func (self *Visor) signBlock(b coin.Block) (sb SignedBlock, e error) {
    if !self.Config.IsMaster {
        log.Panic("Only master chain can sign blocks")
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

// Loads a wallet but subdues errors into the logger, or panics
func loadWallet(filename string, sizeMin int) *Wallet {
    wallet := NewWallet()
    if filename != "" {
        err := wallet.Load(filename)
        if err != nil {
            if os.IsNotExist(err) {
                logger.Info("Wallet file \"%s\" does not exist", filename)
            } else {
                log.Panicf("Failed to load wallet file: %v", err)
            }
        }
    }
    wallet.populate(sizeMin)
    if filename != "" {
        err := wallet.Save(filename)
        if err == nil {
            logger.Info("Saved wallet file to \"%s\"", filename)
        } else {
            log.Panicf("Failed to save wallet file to \"%s\": ", filename,
                err)
        }
    }
    return wallet
}

// Loads a blockchain but subdues errors into the logger, or panics
func loadBlockchain(filename string, masterAddress coin.Address,
    creationInterval uint64) *coin.Blockchain {
    bc := &coin.Blockchain{}
    created := false
    if filename != "" {
        data, err := ioutil.ReadFile(filename)
        if err == nil {
            err = encoder.DeserializeRaw(data, bc)
            if err == nil {
                created = true
                logger.Info("Loaded blockchain from \"%s\"", filename)
            } else {
                log.Panicf("Failed to deserialize blockfrom from \"%s\"",
                    filename)
            }
        } else {
            if os.IsNotExist(err) {
                logger.Info("No blockchain file, will create a new blockchain")
            }
        }
    }

    // Make sure we are not changing the blockchain configuration from the
    // one we loaded
    if created {
        // TODO -- support changing the block creation interval.  Its used
        // in the blockchain internally
        if bc.CreationInterval != creationInterval {
            log.Panic("Creation interval was changed since the old blockchain")
        }
        logger.Notice("Loaded blockchain's genesis address can't be " +
            "checked against configured genesis address")
        logger.Info("Rebuiling UnspentPool indices")
        bc.Unspent.Rebuild()
    } else {
        bc = coin.NewBlockchain(masterAddress, creationInterval)
    }
    return bc
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
