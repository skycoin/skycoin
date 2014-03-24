package visor

import (
    "errors"
    "fmt"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/skycoin/skycoin/src/wallet"
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
    // Wallet files location
    WalletDirectory string
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
    MasterKeys wallet.WalletEntry
    // Genesis block sig
    GenesisSignature coin.Sig
    // Genesis block timestamp
    GenesisTimestamp uint64
    // Number of coins in genesis block
    GenesisCoinVolume uint64
    // Function that creates a new Wallet
    WalletConstructor wallet.WalletConstructor
    // Default type of wallet to create
    WalletTypeDefault wallet.WalletType
}

//Note, put cap on block size, not on transactions/block
//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewVisorConfig() VisorConfig {
    return VisorConfig{
        IsMaster:                 false,
        CanSpend:                 true,
        WalletDirectory:          "",
        BlockCreationInterval:    15,
        UnconfirmedCheckInterval: time.Hour * 2,
        UnconfirmedMaxAge:        time.Hour * 48,
        UnconfirmedRefreshRate:   time.Minute * 30,
        MaxBlockSize:             1024 * 32,
        CoinHourBurnFactor:       2,
        BlockchainFile:           "",
        BlockSigsFile:            "",
        MasterKeys:               wallet.WalletEntry{},
        GenesisSignature:         coin.Sig{},
        GenesisTimestamp:         0,
        GenesisCoinVolume:        100e6,
        WalletConstructor:        wallet.NewSimpleWallet,
        WalletTypeDefault:        wallet.SimpleWalletType,
    }
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
    Config VisorConfig
    // Unconfirmed transactions, held for relay until we get block confirmation
    Unconfirmed *UnconfirmedTxnPool
    // Wallets holding our keys for spending
    Wallets wallet.Wallets
    // Master & personal keys
    masterKeys wallet.WalletEntry
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

    // Load the wallets
    wallets := wallet.Wallets{}
    if c.IsMaster {
        wallets = wallet.Wallets{CreateMasterWallet(c.MasterKeys)}
    } else {
        if c.WalletDirectory != "" {
            w, err := wallet.LoadWallets(c.WalletDirectory)
            if err != nil {
                log.Panicf("Failed to load all wallets: %v", err)
            }
            wallets = w
        }
        if len(wallets) == 0 {
            wallets.Add(c.WalletConstructor())
            if c.WalletDirectory != "" {
                errs := wallets.Save(c.WalletDirectory)
                if len(errs) != 0 {
                    log.Panicf("Failed to save wallets: %v", errs)
                }
            }
        }
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
        Wallets:     wallets,
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
        Wallets:     nil,
    }
}

func (self *Visor) CreateFreshGenesisBlock() (SignedBlock, error) {
    if len(self.blockchain.Blocks) != 0 || len(self.blockSigs.Sigs) != 0 {
        log.Panic("Blockchain already has genesis")
    }
    gb := self.blockchain.CreateGenesisBlock(self.Config.MasterKeys.Address,
        uint64(util.UnixNow()), self.Config.GenesisCoinVolume)
    sb := self.SignBlock(gb)
    if err := self.verifySignedBlock(&sb); err != nil {
        log.Panic("Signed a fresh genesis block, but its invalid: %v", err)
    }
    self.blockSigs.record(&sb)
    return sb, nil
}

// Creates the genesis block as needed
func (self *Visor) CreateGenesisBlock() SignedBlock {
    if len(self.blockchain.Blocks) != 0 || len(self.blockSigs.Sigs) != 0 {
        log.Panic("Blockchain already has genesis")
    }
    addr := self.Config.MasterKeys.Address
    b := self.blockchain.CreateGenesisBlock(addr, self.Config.GenesisTimestamp,
        self.Config.GenesisCoinVolume)
    sb := SignedBlock{
        Block: b,
        Sig:   self.Config.GenesisSignature,
    }
    self.blockSigs.record(&sb)
    err := self.blockSigs.Verify(self.Config.MasterKeys.Public,
        self.blockchain)
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

func (self *Visor) CreateWallet() wallet.Wallet {
    w := self.Config.WalletConstructor()
    self.Wallets.Add(w)
    return w
}

func (self *Visor) SaveWallet(walletID wallet.WalletID) error {
    w := self.Wallets.Get(walletID)
    if w == nil {
        return fmt.Errorf("Unknown wallet %s", walletID)
    }
    return w.Save(self.Config.WalletDirectory)
}

func (self *Visor) SaveWallets() map[wallet.WalletID]error {
    return self.Wallets.Save(self.Config.WalletDirectory)
}

// Loads & unloads wallets based on WalletDirectory contents
func (self *Visor) ReloadWallets() error {
    wallets, err := wallet.LoadWallets(self.Config.WalletDirectory)
    if err != nil {
        return err
    }
    self.Wallets = wallets
    return nil
}

// Saves BlockSigs to disk
func (self *Visor) SaveBlockSigs() error {
    if self.Config.BlockSigsFile == "" {
        return errors.New("No BlockSigsFile location set")
    } else {
        return self.blockSigs.Save(self.Config.BlockSigsFile)
    }
}

// Creates a SignedBlock from pending transactions
func (self *Visor) CreateBlock(when uint64) (SignedBlock, error) {
    var sb SignedBlock
    if !self.Config.IsMaster {
        log.Panic("Only master chain can create blocks")
    }
    if len(self.Unconfirmed.Txns) == 0 {
        return sb, errors.New("No transactions")
    }
    txns := self.Unconfirmed.RawTxns()
    txns = coin.SortTransactions(txns, self.blockchain.TransactionFee)
    txns = txns.TruncateBytesTo(self.Config.MaxBlockSize)
    b, err := self.blockchain.NewBlockFromTransactions(txns, when)
    if err != nil {
        return sb, err
    }
    return self.SignBlock(b), nil
}

// Creates a SignedBlock from pending transactions and executes it
func (self *Visor) CreateAndExecuteBlock() (SignedBlock, error) {
    sb, err := self.CreateBlock(uint64(util.UnixNow()))
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
    addrs := self.Wallets.GetAddressSet()
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
func (self *Visor) Spend(walletID wallet.WalletID, amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    if !self.Config.CanSpend {
        return coin.Transaction{}, errors.New("Spending disabled")
    }
    wallet := self.Wallets.Get(walletID)
    if wallet == nil {
        return coin.Transaction{}, fmt.Errorf("Unknown wallet %v", walletID)
    }
    tx, err := CreateSpendingTransaction(wallet, self.Unconfirmed,
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

// Returns the confirmed & predicted balance for a single address
func (self *Visor) AddressBalance(addr coin.Address) BalancePair {
    auxs := self.blockchain.Unspent.AllForAddress(addr)
    puxs := self.Unconfirmed.SpendsForAddress(&self.blockchain.Unspent, addr)
    confirmed := self.balance(auxs)
    predicted := self.balance(auxs.Sub(puxs))
    return BalancePair{confirmed, predicted}
}

// Returns the confirmed & predicted balance for a Wallet
func (self *Visor) WalletBalance(walletID wallet.WalletID) BalancePair {
    wallet := self.Wallets.Get(walletID)
    if wallet == nil {
        return BalancePair{}
    }
    auxs := self.blockchain.Unspent.AllForAddresses(wallet.GetAddresses())
    puxs := self.Unconfirmed.SpendsForAddresses(&self.blockchain.Unspent,
        wallet.GetAddressSet())
    confirmed := self.totalBalance(auxs)
    predicted := self.totalBalance(auxs.Sub(puxs))
    return BalancePair{confirmed, predicted}
}

// Return the total balance of all loaded wallets
func (self *Visor) TotalBalance() BalancePair {
    b := BalancePair{}
    for _, w := range self.Wallets {
        c := self.WalletBalance(w.GetID())
        b.Confirmed = b.Confirmed.Add(c.Confirmed)
        b.Predicted = b.Confirmed.Add(c.Predicted)
    }
    return b
}

// Computes the total balance for a coin.Address's coin.UxOuts
func (self *Visor) balance(uxs coin.UxArray) Balance {
    prevTime := self.blockchain.Time()
    b := NewBalance(0, 0)
    for _, ux := range uxs {
        b = b.Add(NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
    }
    return b
}

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

// Returns an error if the coin.Sig is not valid for the coin.Block
func (self *Visor) verifySignedBlock(b *SignedBlock) error {
    return coin.VerifySignature(self.Config.MasterKeys.Public, b.Sig,
        b.Block.HashHeader())
}

// Signs a block for master.  Will panic if anything is invalid
func (self *Visor) SignBlock(b coin.Block) SignedBlock {
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

// Creates a wallet with a single master entry
func CreateMasterWallet(master wallet.WalletEntry) wallet.Wallet {
    w := wallet.NewEmptySimpleWallet()
    // The master wallet shouldn't be saved to disk so we clear its filename
    w.SetFilename("")
    if err := w.AddEntry(master); err != nil {
        log.Panicf("Failed to add master wallet entry: %v", err)
    }
    return w
}
