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

// Note: Wallets are moved into RPC. RPC loads/unloads wallets
// Visor allows for injection of raw transactions and new blocks
// Visor saves blocks and loads them from disc
// Visor checks the signature on blocks
// Visor does not generate transactions, but accepts raw transactions

var (
    logger = logging.MustGetLogger("skycoin.visor")
)

// Holds the master and personal keys
type VisorKeys struct {
    // The master server's key.  The Secret will be empty unless running as
    // a master instance
    //Master WalletEntry
    // // Our personal keys
    // Wallet Wallet
    PubKey coin.PubKey
    SecKey coin.SecKey
}

//GenerateVisorKey generates visor privatekey deterministicly from seed valud
func GenerateVisorKey(seed string) VisorKeys {
    pub,sec := coin.GenerateDeterministicKeyPair([]byte(seed))
    return VisorKeys{
        //Master: master,
        Seckey: sec,
        PubKey: pub,
    }
}

// Configuration parameters for the Visor
type VisorConfig struct {
    // Is this the master blockchain
    IsMaster bool
    // Is allowed to create transactions
    //CanSpend bool
    // Wallet file location
    //WalletFile string
    
    // Minimum number of addresses to keep in the wallet
    //WalletSizeMin int
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
    // Where the block signatures are saved
    BlockSigsFile string
    // Master keypair & address
    MasterKeys coin.Sig
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
        //CanSpend:                 true,
        TestNetwork:              true,
        //WalletFile:               "",
        //WalletSizeMin:            1,
        BlockCreationInterval:    15,
        UnconfirmedCheckInterval: time.Hour * 2,
        UnconfirmedMaxAge:        time.Hour * 48,
        UnconfirmedRefreshRate:   time.Minute * 30,
        TransactionsPerBlock:     1000, // 1000/15 = 66tps. Bitcoin is 7tps
        BlockchainFile:           "",
        BlockSigsFile:            "",
        //MasterKeys:               WalletEntry{},
        GenesisSignature:         coin.Sig{},
        GenesisTimestamp:         0,

        keys       VisorKeys
    }
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
    Config VisorConfig
    // Unconfirmed transactions, held for relay until we get block confirmation
    UnconfirmedTxns *UnconfirmedTxnPool
    // Wallet holding our keys for spending
    //Wallet *Wallet
    // Master & personal keys
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

    //TODO: add privatekey from seed

    // Load the wallet
    //var wallet *Wallet = nil
    //if c.IsMaster {
    //    wallet = createMasterWallet(c.MasterKeys)
    //} else {
    //    wallet = loadWallet(c.WalletFile, c.WalletSizeMin)
    //}

    // Load the blockchain the block signatures
    blockchain := loadBlockchain(c.BlockchainFile)
    blockSigs, err := LoadBlockSigs(c.BlockSigsFile)
    if err != nil {
        if os.IsNotExist(err) {
            logger.Info("BlockSigsFile \"%s\" not found", c.BlockSigsFile)
        } else {
            log.Panic("Failed to load BlockSigsFile \"%s\"", c.BlockSigsFile)
        }
        blockSigs = NewBlockSigs()
    }

    v := &Visor{
        Config:          c,
        keys:            GenerateVisorKey(c.MasterKeys),
        blockchain:      blockchain,
        blockSigs:       blockSigs,
        UnconfirmedTxns: NewUnconfirmedTxnPool(),
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
        Config:          c,
        //keys:            GenerateVisorKey(c.MasterKeys),
        blockchain:      coin.NewBlockchain(),
        blockSigs:       NewBlockSigs(),
        UnconfirmedTxns: nil,
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
        var err error
        sb, err = self.signBlock(b)
        if err != nil {
            log.Panicf("Failed to sign genesis block: %v", err)
        }
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

/*
func (self *Visor) SaveWallet() error {
    if self.Config.WalletFile == "" {
        return errors.New("No WalletFile location set")
    } else {
        return self.Wallet.Save(self.Config.WalletFile)
    }
}
*/

// Saves BlockSigs to disk
func (self *Visor) SaveBlockSigs() error {
    if self.Config.BlockSigsFile == "" {
        return errors.New("No BlockSigsFile location set")
    } else {
        return self.blockSigs.Save(self.Config.BlockSigsFile)
    }
}

// Creates and returns a WalletEntry and saves the wallet to disk
/*
func (self *Visor) CreateAddressAndSave() (WalletEntry, error) {
    we := self.Wallet.CreateEntry()
    err := self.SaveWallet()
    if err != nil {
        m := "Failed to save wallet to \"%s\" after creating new address"
        logger.Warning(m, self.Config.WalletFile)
    }
    return we, err
}
*/

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
    b, err := self.blockchain.NewBlockFromTransactions(txns,
        self.Config.BlockCreationInterval)
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

/*
func (self *Visor) Spend(amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    logger.Info("Attempting to send %d coins, %d hours to %s with %d fee",
        amt.Coins, amt.Hours, dest.String(), fee)
    var txn coin.Transaction
    //if !self.Config.CanSpend {
    //    return txn, errors.New("Spending disabled")
    //}
    if amt.IsZero() {
        return txn, errors.New("Zero spend amount")
    }
    needed := amt
    needed.Hours += fee
    auxs := self.getAvailableBalances()
    toSign := make([]coin.SecKey, 0)

loop:
    for a, uxs := range auxs {
        entry, exists := self.Wallet.GetEntry(a)
        if !exists {
            log.Panic("On second thought, the wallet entry does not exist")
        }
        for _, ux := range uxs {
            if needed.IsZero() {
                break loop
            }
            coinHours := ux.CoinHours(self.blockchain.Time())
            b := NewBalance(ux.Body.Coins, coinHours)
            if needed.GreaterThanOrEqual(b) {
                needed = needed.Sub(b)
                txn.PushInput(ux.Hash())
                toSign = append(toSign, entry.Secret)
            } else {
                change := b.Sub(needed)
                needed = needed.Sub(needed)
                txn.PushInput(ux.Hash())
                toSign = append(toSign, entry.Secret)
                // TODO -- Don't reuse address for change.
                txn.PushOutput(ux.Body.Address, change.Coins, change.Hours)
            }
        }
    }

    txn.PushOutput(dest, amt.Coins, amt.Hours)
    txn.SignInputs(toSign)
    txn.UpdateHeader()

    if needed.IsZero() {
        return txn, nil
    } else {
        return txn, errors.New("Not enough coins or hours")
    }
}
*/

//InjectTransaction injections a raw transaction into unconfirmed/pending transaction set.
//raw transactions are created by wallets
func (self *Visor) InjectTransaction(coin.Transaction) (error) {


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

// Returns N signed blocks more recent than Seq. Does not return nil.
func (self *Visor) GetSignedBlocksSince(seq uint64, ct uint64) []SignedBlock {
    var avail uint64 = 0
    if self.blockSigs.MaxSeq > seq {
        avail = self.blockSigs.MaxSeq - seq
    }
    if avail > ct {
        avail = ct
    }
    if avail == 0 {
        return []SignedBlock{}
    }
    blocks := make([]SignedBlock, 0, avail)
    for i := seq + 1; i <= self.blockSigs.MaxSeq; i++ {
        if sig, exists := self.blockSigs.Sigs[i]; exists {
            blocks = append(blocks, SignedBlock{
                Sig:   sig,
                Block: self.blockchain.Blocks[i],
            })
        }
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
    return h.Header.BkSeq
}

// Returns descriptive coin.Blockchain information
func (self *Visor) GetBlockchainMetadata() BlockchainMetadata {
    return NewBlockchainMetadata(self)
}

// Returns a readable copy of the block at seq. Returns error if seq out of range
func (self *Visor) GetReadableBlock(seq uint64) (ReadableBlock, error) {
    b, err := self.GetBlock(seq)
    if err != nil {
        return ReadableBlock{}, err
    }
    return NewReadableBlock(&b), nil
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
        if i >= uint64(len(self.blockchain.Blocks)) {
            break
        }
        blocks = append(blocks, self.blockchain.Blocks[i])
    }
    return blocks
}

// Updates an UnconfirmedTxn's Announce field
func (self *Visor) SetAnnounced(h coin.SHA256, t time.Time) {
    self.UnconfirmedTxns.SetAnnounced(h, t)
}

// Records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
func (self *Visor) RecordTxn(txn coin.Transaction, didAnnounce bool) error {
    entries := make(map[coin.Address]byte, len(self.Wallet.Entries))
    for a, _ := range self.Wallet.Entries {
        entries[a] = byte(1)
    }
    return self.UnconfirmedTxns.RecordTxn(self.blockchain, txn,
        entries, didAnnounce)
}

// Returns the Transactions associated with a coin.Address. This includes
// unconfirmed txns.
func (self *Visor) GetAddressTransactions(a coin.Address) []Transaction {
    txns := make([]Transaction, 0)
    // Look in the blockchain
    uxs := self.blockchain.Unspent.AllForAddress(a)
    for _, ux := range uxs {
        bk := self.blockchain.Blocks[ux.Head.BkSeq]
        tx, ok := bk.GetTransaction(ux.Body.SrcTransaction)
        if ok {
            txns = append(txns, Transaction{
                Txn:    tx,
                Status: NewUnconfirmedTransactionStatus(),
            })
        }
    }

    // Look in the unconfirmed pool
    uxs = self.UnconfirmedTxns.Unspent.AllForAddress(a)
    for _, ux := range uxs {
        tx, ok := self.UnconfirmedTxns.Txns[ux.Body.SrcTransaction]
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
    tx, ok := self.UnconfirmedTxns.Txns[txHash]
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
            height := self.MostRecentBkSeq() - b.Header.BkSeq + 1
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

// Returns the balance of the wallet

/*
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
*/

//TODO - return UxOut

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

/*
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
*/

/*
// Creates a wallet with a single master entry
func createMasterWallet(master WalletEntry) *Wallet {
    w := NewWallet()
    if err := w.AddEntry(master); err != nil {
        log.Panic("Master entry already exists in wallet: %v", err)
    }
    return w
}
*/

// Loads a coin.Blockchain from disk
func LoadBlockchain(filename string) (*coin.Blockchain, error) {
    bc := &coin.Blockchain{}
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return bc, err
    }
    return bc, encoder.DeserializeRaw(data, bc)
}

// Loads a blockchain but subdues errors into the logger, or panics
func loadBlockchain(filename string) *coin.Blockchain {
    bc := &coin.Blockchain{}
    created := false
    if filename != "" {
        var err error
        bc, err = LoadBlockchain(filename)
        if err == nil {
            created = true
            logger.Info("Loaded blockchain from \"%s\"", filename)
            logger.Notice("Loaded blockchain's genesis address can't be " +
                "checked against configured genesis address")
            logger.Info("Rebuiling UnspentPool indices")
            bc.Unspent.Rebuild()
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
