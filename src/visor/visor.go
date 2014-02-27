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
    "time"
)

var (
    logger = logging.MustGetLogger("skycoin.visor")
)

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
    // How often new blocks are created by the master, in seconds
    BlockCreationInterval uint64
    // How often an unconfirmed txn is checked against the blockchain
    UnconfirmedCheckInterval time.Duration
    // How long we'll hold onto an unconfirmed txn
    UnconfirmedMaxAge time.Duration
    // How often to refresh the unconfirmed pool
    UnconfirmedRefreshRate time.Duration
    // Maximum size of a block, in bytes.
    MaxBlockSize int
    // Divisor of coin hours required as fee. E.g. with hours=100 and factor=4,
    // 25 additional hours are required as a fee.  A value of 0 disables
    // the fee requirement.
    CoinHourBurnFactor uint64
    // Where the blockchain is saved
    BlockchainFile string
    // Where the block signatures are saved
    BlockSigsFile string
    // Master keypair & address
    MasterKeys WalletEntry
    // Genesis block sig
    GenesisSignature coin.Sig
    // Genesis block timestamp
    GenesisTimestamp uint64
}

//Note, put cap on block size, not on transactions/block
//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewVisorConfig() VisorConfig {
    return VisorConfig{
        IsMaster:                 false,
        CanSpend:                 true,
        WalletFile:               "",
        WalletSizeMin:            1,
        BlockCreationInterval:    15,
        UnconfirmedCheckInterval: time.Hour * 2,
        UnconfirmedMaxAge:        time.Hour * 48,
        UnconfirmedRefreshRate:   time.Minute * 30,
        MaxBlockSize:             1024 * 32,
        CoinHourBurnFactor:       2,
        BlockchainFile:           "",
        BlockSigsFile:            "",
        MasterKeys:               WalletEntry{},
        GenesisSignature:         coin.Sig{},
        GenesisTimestamp:         0,
    }
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
    Config VisorConfig
    // Unconfirmed transactions, held for relay until we get block confirmation
    Unconfirmed *UnconfirmedTxnPool
    // Wallet holding our keys for spending
    Wallet Wallet
    // Master & personal keys
    masterKeys WalletEntry
    blockchain *coin.Blockchain
    blockSigs  BlockSigs
}

// Creates a normal Visor given a master's public key
func NewVisor(c VisorConfig) *Visor {
    logger.Debug("Creating new visor")
    // Make sure inputs are correct
    if c.IsMaster {
        logger.Debug("Visor is master")
    }
    if c.IsMaster {
        if err := c.MasterKeys.Verify(); err != nil {
            log.Panicf("Invalid master wallet entry: %v", err)
        }
    } else {
        if err := c.MasterKeys.VerifyPublic(); err != nil {
            log.Panicf("Invalid master address or pubkey: %v", err)
        }
    }

    // Load the wallet
    wallet := Wallet(nil)
    if c.IsMaster {
        wallet = CreateMasterWallet(c.MasterKeys)
    } else {
        wallet = loadSimpleWallet(c.WalletFile, c.WalletSizeMin)
    }

    // Load the blockchain the block signatures
    blockchain := loadBlockchain(c.BlockchainFile, c.MasterKeys.Address)
    blockSigs, err := LoadBlockSigs(c.BlockSigsFile)
    if err != nil {
        if os.IsNotExist(err) {
            logger.Info("BlockSigsFile \"%s\" not found", c.BlockSigsFile)
        } else {
            log.Panicf("Failed to load BlockSigsFile \"%s\"", c.BlockSigsFile)
        }
        blockSigs = NewBlockSigs()
    }

    v := &Visor{
        Config:      c,
        blockchain:  blockchain,
        blockSigs:   blockSigs,
        Unconfirmed: NewUnconfirmedTxnPool(),
        Wallet:      wallet,
    }
    // Load the genesis block and sign it, if we need one
    if len(blockchain.Blocks) == 0 {
        v.CreateGenesisBlock()
    }
    err = blockSigs.Verify(c.MasterKeys.Public, blockchain)
    if err != nil {
        log.Panicf("Invalid block signatures: %v", err)
    }

    return v
}

// Returns a Visor with minimum initialization necessary for empty blockchain
// access
func NewMinimalVisor(c VisorConfig) *Visor {
    return &Visor{
        Config:      c,
        blockchain:  coin.NewBlockchain(),
        blockSigs:   NewBlockSigs(),
        Unconfirmed: NewUnconfirmedTxnPool(),
        Wallet:      nil,
    }
}

// Creates the genesis block as needed
func (self *Visor) CreateGenesisBlock() SignedBlock {
    b := coin.Block{}
    addr := self.Config.MasterKeys.Address
    if self.Config.IsMaster {
        b = self.blockchain.CreateMasterGenesisBlock(addr)
    } else {
        b = self.blockchain.CreateGenesisBlock(addr, self.Config.GenesisTimestamp)
    }
    sb := SignedBlock{}
    if self.Config.IsMaster {
        sb = self.signBlock(b)
    } else {
        sb = SignedBlock{
            Block: b,
            Sig:   self.Config.GenesisSignature,
        }
    }
    self.blockSigs.record(&sb)
    err := self.blockSigs.Verify(self.Config.MasterKeys.Public, self.blockchain)
    if err != nil {
        log.Panicf("Signed the genesis block, but its invalid: %v", err)
    }
    return sb
}

// Checks unconfirmed txns against the blockchain and purges ones too old
func (self *Visor) RefreshUnconfirmed() {
    logger.Debug("Refreshing unconfirmed transactions")
    self.Unconfirmed.Refresh(self.blockchain,
        self.Config.UnconfirmedCheckInterval, self.Config.UnconfirmedMaxAge)
}

// Saves the coin.Blockchain to disk
func (self *Visor) SaveBlockchain() error {
    if self.Config.BlockchainFile == "" {
        return errors.New("No BlockchainFile location set")
    } else {
        return SaveBlockchain(self.blockchain, self.Config.BlockchainFile)
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

// Saves BlockSigs to disk
func (self *Visor) SaveBlockSigs() error {
    if self.Config.BlockSigsFile == "" {
        return errors.New("No BlockSigsFile location set")
    } else {
        return self.blockSigs.Save(self.Config.BlockSigsFile)
    }
}

// Creates and returns a WalletEntry and saves the wallet to disk
func (self *Visor) CreateAddressAndSave() (WalletEntry, error) {
    we := self.Wallet.CreateEntry()
    err := self.SaveWallet()
    if err != nil {
        m := "Failed to save wallet to \"%s\" after creating new address"
        logger.Warning(m, self.Config.WalletFile)
    }
    return we, err
}

// Creates a SignedBlock from pending transactions
func (self *Visor) createBlock() (SignedBlock, error) {
    var sb SignedBlock
    if !self.Config.IsMaster {
        log.Panic("Only master chain can create blocks")
    }
    if len(self.Unconfirmed.Txns) == 0 {
        return sb, errors.New("No transactions")
    }
    txns := self.Unconfirmed.RawTxns()
    b, err := self.blockchain.NewBlockFromTransactions(txns,
        self.Config.BlockCreationInterval, self.Config.MaxBlockSize)
    if err != nil {
        return sb, err
    }
    return self.signBlock(b), nil
}

// Creates a SignedBlock from pending transactions and executes it
func (self *Visor) CreateAndExecuteBlock() (SignedBlock, error) {
    sb, err := self.createBlock()
    if err == nil {
        return sb, self.ExecuteSignedBlock(sb)
    } else {
        return sb, err
    }
}

// Adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (self *Visor) ExecuteSignedBlock(b SignedBlock) error {
    if err := self.verifySignedBlock(&b); err != nil {
        return err
    }
    _, err := self.blockchain.ExecuteBlock(b.Block)
    if err != nil {
        return err
    }
    // TODO -- save them even if out of order, and execute later
    // But make sure all prechecking as possible is done
    // TODO -- check if bitcoin allows blocks to be receiving out of order
    self.blockSigs.record(&b)
    // Remove the transactions in the Block from the unconfirmed pool
    self.Unconfirmed.RemoveTransactions(self.blockchain,
        b.Block.Body.Transactions)
    return nil
}

// Returns N signed blocks more recent than Seq. Does not return nil.
func (self *Visor) GetSignedBlocksSince(seq, ct uint64) []SignedBlock {
    var avail uint64 = 0
    if self.blockSigs.MaxSeq > seq {
        avail = self.blockSigs.MaxSeq - seq
    }
    if avail < ct {
        ct = avail
    }
    if ct == 0 {
        return []SignedBlock{}
    }
    blocks := make([]SignedBlock, 0, ct)
    for j := uint64(0); j < ct; j++ {
        i := seq + 1 + j
        blocks = append(blocks, SignedBlock{
            Sig:   self.blockSigs.Sigs[i],
            Block: self.blockchain.Blocks[i],
        })
    }
    return blocks
}

// Returns the signed genesis block. Panics if signature or block not found
func (self *Visor) GetGenesisBlock() SignedBlock {
    gsig, ok := self.blockSigs.Sigs[0]
    if !ok {
        log.Panic("No genesis signature")
    }
    if len(self.blockchain.Blocks) == 0 {
        log.Panic("No genesis block")
    }
    return SignedBlock{
        Sig:   gsig,
        Block: self.blockchain.Blocks[0],
    }
}

// Returns the highest BkSeq we know
func (self *Visor) MostRecentBkSeq() uint64 {
    h := self.blockchain.Head()
    return h.Head.BkSeq
}

// Returns descriptive coin.Blockchain information
func (self *Visor) GetBlockchainMetadata() BlockchainMetadata {
    return NewBlockchainMetadata(self)
}

// Returns a readable copy of the block at seq. Returns error if seq out of range
func (self *Visor) GetReadableBlock(seq uint64) (ReadableBlock, error) {
    if b, err := self.GetBlock(seq); err == nil {
        return NewReadableBlock(&b), nil
    } else {
        return ReadableBlock{}, err
    }
}

// Returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
func (self *Visor) GetReadableBlocks(start, end uint64) []ReadableBlock {
    blocks := self.GetBlocks(start, end)
    rbs := make([]ReadableBlock, 0, len(blocks))
    for _, b := range blocks {
        rbs = append(rbs, NewReadableBlock(&b))
    }
    return rbs
}

// Returns a copy of the block at seq. Returns error if seq out of range
func (self *Visor) GetBlock(seq uint64) (coin.Block, error) {
    var b coin.Block
    if seq >= uint64(len(self.blockchain.Blocks)) {
        return b, errors.New("Block seq out of range")
    }
    return self.blockchain.Blocks[seq], nil
}

// Returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
func (self *Visor) GetBlocks(start, end uint64) []coin.Block {
    if end > uint64(len(self.blockchain.Blocks)) {
        end = uint64(len(self.blockchain.Blocks))
    }
    var length uint64 = 0
    if start < end {
        length = end - start
    }
    blocks := make([]coin.Block, 0, length)
    for i := start; i < end; i++ {
        blocks = append(blocks, self.blockchain.Blocks[i])
    }
    return blocks
}

// Updates an UnconfirmedTxn's Announce field
func (self *Visor) SetAnnounced(h coin.SHA256, t time.Time) {
    self.Unconfirmed.SetAnnounced(h, t)
}

// Records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
func (self *Visor) RecordTxn(txn coin.Transaction) (error, bool) {
    addrs := self.Wallet.GetAddressSet()
    return self.Unconfirmed.RecordTxn(self.blockchain, txn, addrs,
        self.Config.MaxBlockSize, self.Config.CoinHourBurnFactor)
}

// Returns the Transactions whose unspents give coins to a coin.Address.
// This includes unconfirmed txns' predicted unspents.
func (self *Visor) GetAddressTransactions(a coin.Address) []Transaction {
    txns := make([]Transaction, 0)
    // Look in the blockchain
    uxs := self.blockchain.Unspent.AllForAddress(a)
    mxSeq := self.MostRecentBkSeq()
    for _, ux := range uxs {
        bk := self.blockchain.Blocks[ux.Head.BkSeq]
        tx, ok := bk.GetTransaction(ux.Body.SrcTransaction)
        if ok {
            h := mxSeq - bk.Head.BkSeq + 1
            txns = append(txns, Transaction{
                Txn:    tx,
                Status: NewConfirmedTransactionStatus(h),
            })
        }
    }

    // Look in the unconfirmed pool
    uxs = self.Unconfirmed.Unspent.AllForAddress(a)
    for _, ux := range uxs {
        tx, ok := self.Unconfirmed.Txns[ux.Body.SrcTransaction]
        if !ok {
            logger.Critical("Unconfirmed unspent missing unconfirmed txn")
            continue
        }
        txns = append(txns, Transaction{
            Txn:    tx.Txn,
            Status: NewUnconfirmedTransactionStatus(),
        })
    }

    return txns
}

// Returns a Transaction by hash.
func (self *Visor) GetTransaction(txHash coin.SHA256) Transaction {
    // Look in the unconfirmed pool
    tx, ok := self.Unconfirmed.Txns[txHash]
    if ok {
        return Transaction{
            Txn:    tx.Txn,
            Status: NewUnconfirmedTransactionStatus(),
        }
    }

    // Look in the blockchain
    // TODO -- this is extremely slow as it does a full blockchain scan
    // We need an index from txn hash to block.  At least an index per block
    // to its contained txns
    for _, b := range self.blockchain.Blocks {
        tx, ok := b.GetTransaction(txHash)
        if ok {
            height := self.MostRecentBkSeq() - b.Head.BkSeq + 1
            return Transaction{
                Txn:    tx,
                Status: NewConfirmedTransactionStatus(height),
            }
        }
    }

    // Otherwise unknown
    return Transaction{
        Status: NewUnknownTransactionStatus(),
    }
}

// Creates a transaction spending amt with additional fee.  Fee is in addition
// to the base required fee given amt.Hours.
func (self *Visor) Spend(amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    if !self.Config.CanSpend {
        return coin.Transaction{}, errors.New("Spending disabled")
    }
    tx, err := CreateSpendingTransaction(self.Wallet, self.Unconfirmed,
        &self.blockchain.Unspent, self.blockchain.Time(), amt, fee,
        self.Config.CoinHourBurnFactor, dest)
    if err != nil {
        return tx, err
    }
    if err := VerifyTransaction(self.blockchain, &tx, self.Config.MaxBlockSize,
        self.Config.CoinHourBurnFactor); err != nil {
        log.Panicf("Created invalid spending txn: %v", err)
    }
    if err := self.blockchain.VerifyTransaction(tx); err != nil {
        log.Panicf("Created invalid spending txn: %v", err)
    }
    return tx, err
}

// Returns the balance of the wallet
func (self *Visor) TotalBalance() Balance {
    addrs := self.Wallet.GetAddresses()
    auxs := self.blockchain.Unspent.AllForAddresses(addrs)
    return self.totalBalance(auxs)
}

// // Returns the total balance of the wallet including unconfirmed outputs
// func (self *Visor) TotalBalancePredicted() Balance {
//     auxs := self.getAvailableBalances()
//     return self.totalBalance(auxs)
// }

// Returns the balance for a single address in the Wallet
func (self *Visor) Balance(a coin.Address) Balance {
    uxs := self.blockchain.Unspent.AllForAddress(a)
    return self.balance(uxs)
}

// // Returns the balance for a single address in the Wallet, including
// // unconfirmed outputs
// func (self *Visor) BalancePredicted(a coin.Address) Balance {
//     uxs := self.getAvailableBalance(a)
//     return self.balance(uxs)
// }

// Computes the total balance for coin.Addresses and their coin.UxOuts
func (self *Visor) totalBalance(auxs coin.AddressUxOuts) Balance {
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

// // Returns the total of known Unspents available to us, and our own
// // unconfirmed unspents
// func (self *Visor) getAvailableBalances() coin.AddressUxOuts {
//     addrs := self.Wallet.GetAddresses()
//     auxs := self.blockchain.Unspent.AllForAddresses(addrs)
//     uauxs := self.Unconfirmed.Unspent.AllForAddresses(addrs)
//     logger.Warning("Confirmed unspents: %v\n", auxs)
//     logger.Warning("Unconfirmed unspents: %v\n", uauxs)
//     return auxs.Merge(uauxs, addrs)
// }

// // Returns the total of known unspents available for an address, including
// // unconfirmed requests
// func (self *Visor) getAvailableBalance(a coin.Address) []coin.UxOut {
//     auxs := self.blockchain.Unspent.AllForAddress(a)
//     uauxs := self.Unconfirmed.Unspent.AllForAddress(a)
//     return append(auxs, uauxs...)
// }

// Returns an error if the coin.Sig is not valid for the coin.Block
func (self *Visor) verifySignedBlock(b *SignedBlock) error {
    return coin.VerifySignature(self.Config.MasterKeys.Public, b.Sig,
        b.Block.HashHeader())
}

// Signs a block for master.  Will panic if anything is invalid
func (self *Visor) signBlock(b coin.Block) SignedBlock {
    if !self.Config.IsMaster {
        log.Panic("Only master chain can sign blocks")
    }
    sig := coin.SignHash(b.HashHeader(), self.Config.MasterKeys.Secret)
    sb := SignedBlock{
        Block: b,
        Sig:   sig,
    }
    return sb
}

// Loads a wallet but subdues errors into the logger, or panics
func loadSimpleWallet(filename string, sizeMin int) *SimpleWallet {
    wallet := NewSimpleWallet()
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
    wallet.Populate(sizeMin)
    if filename != "" {
        err := wallet.Save(filename)
        if err == nil {
            logger.Info("Saved wallet file to \"%s\"", filename)
        } else {
            log.Panicf("Failed to save wallet file to \"%s\": %v", filename,
                err)
        }
    }
    return wallet
}

// Creates a wallet with a single master entry
func CreateMasterWallet(master WalletEntry) *SimpleWallet {
    w := NewSimpleWallet()
    if err := w.AddEntry(master); err != nil {
        log.Panic("Failed to add master wallet entry: %v", err)
    }
    return w
}

// Loads a coin.Blockchain from disk
func LoadBlockchain(filename string) (*coin.Blockchain, error) {
    bc := &coin.Blockchain{}
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return bc, err
    }
    err = encoder.DeserializeRaw(data, bc)
    if err != nil {
        return bc, err
    }
    logger.Info("Loaded blockchain from \"%s\"", filename)
    logger.Debug("Rebuilding UnspentPool indices")
    bc.Unspent.Rebuild()
    return bc, nil
}

// Loads a blockchain but subdues errors into the logger, or panics.
// If no blockchain is found, it creates a new empty one
func loadBlockchain(filename string, genAddr coin.Address) *coin.Blockchain {
    bc := &coin.Blockchain{}
    created := false
    if filename != "" {
        var err error
        bc, err = LoadBlockchain(filename)
        if err == nil {
            if len(bc.Blocks) == 0 {
                log.Panic("Loaded empty blockchain")
            }
            loadedGenAddr := bc.Blocks[0].Body.Transactions[0].Out[0].Address
            if loadedGenAddr != genAddr {
                log.Panic("Configured genesis address does not match the " +
                    "address in the blockchain")
            }
            created = true
        } else {
            if os.IsNotExist(err) {
                logger.Info("No blockchain file, will create a new blockchain")
            } else {
                log.Panicf("Failed to load blockchain file \"%s\": %v",
                    filename, err)
            }
        }
    }
    if !created {
        bc = coin.NewBlockchain()
    }
    return bc
}

// Saves blockchain to disk
func SaveBlockchain(bc *coin.Blockchain, filename string) error {
    // TODO -- blockchain file must be forward compatible
    data := encoder.Serialize(bc)
    return util.SaveBinary(filename, data, 0644)
}
