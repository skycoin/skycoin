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

const (
    // Timestamp of the genesis block
    GENESIS_TIMESTAMP uint64 = 1391649057
    // Signature of the genesis block
    GENESIS_SIGNATURE = "a1a09bee02a92fddaf34856aedde9c1ef626caaf31ada221fc2acc9212493e61064b32d4cfd92f38948e799f231f8c42428086405bbf42f9e913a149c0ca743f00"
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
    // Where the block signatures are saved
    BlockSigsFile string
    // Master keypair & address
    MasterKeys WalletEntry
    // Genesis block sig
    GenesisSig coin.Sig
}

func NewVisorConfig() VisorConfig {
    return VisorConfig{
        IsMaster:                 false,
        CanSpend:                 true,
        TestNetwork:              true,
        WalletFile:               "",
        WalletSizeMin:            1,
        BlockCreationInterval:    15,
        UnconfirmedCheckInterval: time.Hour * 2,
        UnconfirmedMaxAge:        time.Hour * 48,
        UnconfirmedRefreshRate:   time.Minute * 30,
        TransactionsPerBlock:     1000, // 1000/15 = 66tps. Bitcoin is 7tps
        BlockchainFile:           "",
        BlockSigsFile:            "",
        MasterKeys:               WalletEntry{},
        GenesisSig:               coin.Sig{},
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
    var wallet *Wallet = nil
    if c.IsMaster {
        wallet = createMasterWallet(c.MasterKeys)
    } else {
        wallet = loadWallet(c.WalletFile, c.WalletSizeMin)
    }

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
    err = blockSigs.Verify(c.MasterKeys.Public, blockchain)
    if err != nil {
        log.Panic("Invalid block signatures")
    }

    v := &Visor{
        Config:          c,
        keys:            NewVisorKeys(c.MasterKeys),
        blockchain:      blockchain,
        blockSigs:       blockSigs,
        UnconfirmedTxns: NewUnconfirmedTxnPool(),
        Wallet:          wallet,
    }

    // Load the genesis block and sign it, if we need one
    if len(blockchain.Blocks) == 0 {
        v.CreateGenesisBlock()
    }
    return v
}

// Returns a Visor with minimum initialization necessary for empty blockchain
// access
func NewMinimalVisor(c VisorConfig) *Visor {
    return &Visor{
        Config:          c,
        keys:            NewVisorKeys(c.MasterKeys),
        blockchain:      coin.NewBlockchain(),
        blockSigs:       NewBlockSigs(),
        UnconfirmedTxns: nil,
        Wallet:          nil,
    }
}

// Creates the genesis block as needed
func (self *Visor) CreateGenesisBlock() SignedBlock {
    b := coin.Block{}
    addr := self.Config.MasterKeys.Address
    if self.Config.IsMaster {
        b = self.blockchain.CreateMasterGenesisBlock(addr)
    } else {
        b = self.blockchain.CreateGenesisBlock(addr, GENESIS_TIMESTAMP)
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
            Sig:   coin.MustSigFromHex(GENESIS_SIGNATURE),
        }
    }
    self.blockSigs.record(&sb)
    err := self.blockSigs.Verify(self.Config.MasterKeys.Public, self.blockchain)
    if err != nil {
        log.Panicf("Signed the genesis block, but its invalid: %v", err)
    }
    return sb
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
func (self *Visor) Spend(amt Balance, fee uint64,
    dest coin.Address) (coin.Transaction, error) {
    logger.Info("Attempting to send %d coins, %d hours to %s with %d fee",
        amt.Coins, amt.Hours, dest.String(), fee)
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

// Returns the highest BkSeq we know
func (self *Visor) MostRecentBkSeq() uint64 {
    return self.blockSigs.MaxSeq
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

// Records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
func (self *Visor) RecordTxn(txn coin.Transaction) error {
    addrs := self.Wallet.GetAddresses()
    return self.UnconfirmedTxns.RecordTxn(self.blockchain, txn, addrs)
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

// Creates a wallet with a single master entry
func createMasterWallet(master WalletEntry) *Wallet {
    w := NewWallet()
    if err := w.AddEntry(master); err != nil {
        log.Panic("Master entry already exists in wallet: %v", err)
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

type SignedBlock struct {
    Block coin.Block
    Sig   coin.Sig
}

// Used to serialize the BlockSigs.Sigs map
type BlockSigSerialized struct {
    BkSeq uint64
    Sig   coin.Sig
}

// Used to serialize the BlockSigs.Sigs map
type BlockSigsSerialized struct {
    Sigs []BlockSigSerialized
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
}

func NewBlockSigs() BlockSigs {
    bs := BlockSigs{
        Sigs:   make(map[uint64]coin.Sig),
        MaxSeq: 0,
    }
    return bs
}

func LoadBlockSigs(filename string) (BlockSigs, error) {
    bs := NewBlockSigs()
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return bs, err
    }
    sigs := BlockSigsSerialized{make([]BlockSigSerialized, 0)}
    err = encoder.DeserializeRaw(data, &sigs)
    if err != nil {
        return bs, err
    }
    bs.Sigs = make(map[uint64]coin.Sig, len(sigs.Sigs))
    for _, s := range sigs.Sigs {
        bs.Sigs[s.BkSeq] = s.Sig
        if s.BkSeq > bs.MaxSeq {
            bs.MaxSeq = s.BkSeq
        }
    }
    return bs, nil
}

func (self *BlockSigs) Save(filename string) error {
    // Convert the Sigs map to an array of element
    sigs := make([]BlockSigSerialized, 0, len(self.Sigs))
    for k, v := range self.Sigs {
        sigs = append(sigs, BlockSigSerialized{
            BkSeq: k,
            Sig:   v,
        })
    }
    bss := BlockSigsSerialized{sigs}
    data := encoder.Serialize(bss)
    return util.SaveBinary(filename, data, 0644)
}

// Checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
func (self *BlockSigs) Verify(masterPublic coin.PubKey, bc *coin.Blockchain) error {
    blocks := uint64(len(bc.Blocks))
    if blocks != uint64(len(self.Sigs)) {
        return errors.New("NSigs != NBlocks")
    }
    for k, v := range self.Sigs {
        if k > self.MaxSeq {
            return errors.New("Invalid MaxSeq")
        } else if k > blocks {
            return errors.New("Signature for unknown block")
        }
        err := coin.VerifySignature(masterPublic, v, bc.Blocks[k].HashHeader())
        if err != nil {
            return err
        }
    }
    return nil
}

// Adds a SignedBlock
func (self *BlockSigs) record(sb *SignedBlock) {
    seq := sb.Block.Header.BkSeq
    self.Sigs[seq] = sb.Sig
    if seq > self.MaxSeq {
        self.MaxSeq = seq
    }
}
