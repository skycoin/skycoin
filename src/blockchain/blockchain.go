package Blockchain

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
    logger = logging.MustGetLogger("skycoin.Blockchain")
)

// Note: can use testnetpubkey as genesis address

type GenesisBlockCfg {
    GenesisAddress coin.Address
    GenesisSignature coin.Sig
    GenesisTime uint64
    PubKey coin.PubKey
    Coins uint64
}

var TestNet GenesisBlockCfg
var SkyNet GenesisBlockCfg //main blockchain

//testnet config
func init() {
    TestNet.PubKey = coin.MustPubKeyFromHex("025a3b22eb1e132a01f485119ae343342d92ab8599d9ad613a76e3b27f878bca8b")
    //TestNet.GenesisSignature = coin.MustSigFromHex()
    TestNet.GenesisAddress = coin.MustDecodeBase58Address("26HbgWGwrToLZ6aX8VHtQmH4SPj4baQ5S3p")
    TestNet.GenesisTime = 1392584986 //set time
    TestNet.Coins = 1e12 //almost as many as Ripple
    //TestNet.GenesisSignature = coin.MustSigFromHex()

}

//main net config
func init() {
    SkyNet.PubKey = coin.MustPubKeyFromHex("02bb0be2976457d2e30a9aea9b0057b0eb9d1ad6509ef743c25c737f24d6241a99")
    //TestNet.GenesisSignature = coin.MustSigFromHex()
    SkyNet.GenesisAddress = coin.MustDecodeBase58Address("26HbgWGwrToLZ6aX8VHtQmH4SPj4baQ5S3p")
    SkyNet.GenesisTime = 1392584987 //set time
    SkyNet.Coins = 100e6 //100 million
}

// Configuration parameters for the Blockchain
type BlockchainConfig struct {

    Genesis GenesisBlockCfg
    // Is this the master blockchain
    IsMaster bool
    SecKey coin.SecKey //set for writes

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
    //Maximum number of transactions per block, when creating
    TransactionsPerBlock int
    
    // Where the blockchain is saved
    BlockchainFile string    

}

func (self *BlockchainConfig) SetBlockchainSecKey(seed string) {
    pub,sec := coin.GenerateDeterministicKeyPair([]byte(seed))
    if pub != self.PubKey {
        log.Panic("ERROR: pubkey does not correspond to loaded pubkey")
    }
    self.SecKey = sec
}

//Note, put cap on block size, not on transactions/block
//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewBlockchainConfig() BlockchainConfig {
    //set pubkey based upon testnet, mainnet and local
    return BlockchainConfig{

        IsMaster:                 false,
        TestNetwork:              true,

        BlockCreationInterval:    15,
        UnconfirmedCheckInterval: time.Minute * 5,
        UnconfirmedMaxAge:        time.Hour * 48, //drop transaction not executed in 48 hours
        UnconfirmedRefreshRate:   time.Minute * 30,
        TransactionsPerBlock:     150, //10 transactions/second, 1.5 KB/s
        BlockchainFile:           "",
        BlockSigsFile:            "",

        SecKey: coin.SecKey{},
    }
}

//NewTestnetBlockchain Config creates Blockchain for the testnet
func NewTestnetBlockchainConfig() BlockchainConfig {
    VC := NewBlockchainConfig()
    VC.PubKey = coin.MustPubKeyFromHex(testnet_pubkey_hex)
    VC.TestNetwork = true
    VC.Genesis = TestNet
    return VC
}

//NewTestnetBlockchain Config creates Blockchain for the mainnet
func NewMainnetBlockchainConfig() BlockchainConfig {
    VC := NewBlockchainConfig()
    VC.PubKey = coin.MustPubKeyFromHex(mainnet_pubkey_hex)
    VC.TestNetwork = false
    VC.Genesis = MainNet
    return VC
}

//Generate Blockchain configuration for client only Blockchain, not intended to be synced to network
func NewLocalBlockchainConfig() BlockchainConfig {
    pubkey,seckey := coin.GenerateKeyPair() //generate new/random pubkey/private key
    
    VC := NewBlockchainConfig()
    VC.SecKey = seckey
    VC.Genesis.Pubkey = PubKey
    VC.Genesis.GenesisAddress = coin.AddressFromPubKey(pubkey)
    return VC
}


// Manages the Blockchain as both a Master and a Normal
type Blockchain struct {
    Config BlockchainConfig
    // Unconfirmed transactions, held for relay until we get block confirmation
    UnconfirmedTxns *UnconfirmedTxnPool
    //blockchain storag
    blockchain *coin.Blockchain
    blockSigs  BlockSigs
}

// Creates a normal Blockchain given a master's public key
func NewBlockchain(c BlockchainConfig) *Blockchain {
    logger.Debug("Creating new Blockchain")
    // Make sure inputs are correct
    if c.IsMaster {
        logger.Debug("Blockchain is master")
    }
    if c.IsMaster {
        if err := c.SecKey.Verify(); err != nil {
            log.Panicf("Invalid privatekey: %v", err)
        }
        if c.PubKey != coin.PubKeyFromSecKey(c.SecKey) {
            log.Panic("SecKey does not correspond to PubKey")
        }
    } else {
        if err := c.PubKey.Verify(); err != nil {
            log.Panicf("Invalid pubkey: %v", err)
        }

    }

/*
    // Load the blockchain the block signatures
    blockchain := loadBlockchain(c.BlockchainFile)
    blockSigs, err := LoadBlockSigs(c.BlockSigsFile)
    if err != nil {
        if os.IsNotExist(err) {
            logger.Info("BlockSigsFile \"%s\" not found", c.BlockSigsFile)
        } else {
            log.Panicf("Failed to load BlockSigsFile \"%s\"", c.BlockSigsFile)
        }
        blockSigs = NewBlockSigs()
    }

    v := &Blockchain{
        Config:          c,
        blockchain:      blockchain,
        blockSigs:       blockSigs,
        UnconfirmedTxns: NewUnconfirmedTxnPool(),
    }
    // Load the genesis block and sign it, if we need one
    if len(blockchain.Blocks) == 0 {
        v.CreateGenesisBlock()
    }
    err = blockSigs.Verify(c.PubKey, blockchain)
    if err != nil {
        log.Panicf("Invalid block signatures: %v", err)
    }
*/
    return v
}

// Returns a Blockchain with minimum initialization necessary for empty blockchain
// access

/*
func NewMinimalBlockchain(c BlockchainConfig) *Blockchain {
    return &Blockchain{
        Config:          c,
        blockchain:      coin.NewBlockchain(),
        blockSigs:       NewBlockSigs(),
        UnconfirmedTxns: nil,
    }
}
*/

// Creates the genesis block
func (self *Blockchain) PushGenesisBlock() SignedBlock {
    
    self.Config.IsMaster == false {
        log.Panic()
    }

    //b := coin.Block{}
    addr := coin.MustDecodeBase58Address(genesis_address) //genesis address
    b := self.blockchain.CreateMasterGenesisBlock(addr)
    sb = self.signBlock(b)
    self.blockSigs.record(&sb)
    err := self.blockSigs.Verify(self.Config.PubKey, self.blockchain)
    if err != nil {
        log.Panicf("Signed the genesis block, but its invalid: %v", err)
    }
    return sb
}

func (self *Blockchain) CreateGenesisBlock() SignedBlock {
    b := coin.Block{}
    addr := coin.MustDecodeBase58Address(genesis_address) //genesis address
    //addr := coin.AddressFromPubKey(self.Config.PubKey)

    b = self.blockchain.CreateGenesisBlock(addr, self.Config.GenesisTimestamp)

    //sb := SignedBlock{}

    sb := SignedBlock{
        Block: b,
        Sig:   self.Config.GenesisSignature,
    }

    self.blockSigs.record(&sb)
    err := self.blockSigs.Verify(self.Config.PubKey, self.blockchain)
    if err != nil {
        log.Panicf("Signed the genesis block, but its invalid: %v", err)
    }
    return sb
}

// Checks unconfirmed txns against the blockchain and purges ones too old
func (self *Blockchain) RefreshUnconfirmed() {
    logger.Debug("Refreshing unconfirmed transactions")
    self.UnconfirmedTxns.Refresh(self.blockchain,
        self.Config.UnconfirmedCheckInterval, self.Config.UnconfirmedMaxAge)
}

// Saves the coin.Blockchain to disk
func (self *Blockchain) SaveBlockchain() error {
    if self.Config.BlockchainFile == "" {
        return errors.New("No BlockchainFile location set")
    } else {
        return SaveBlockchain(self.blockchain, self.Config.BlockchainFile)
    }
}


// Saves BlockSigs file to disk
func (self *Blockchain) SaveBlockSigs() error {
    if self.Config.BlockSigsFile == "" {
        return errors.New("No BlockSigsFile location set")
    } else {
        return self.blockSigs.Save(self.Config.BlockSigsFile)
    }
}

// Creates a SignedBlock from pending transactions
func (self *Blockchain) createBlock() (SignedBlock, error) {
    //var sb SignedBlock
    if !self.Config.IsMaster {
        log.Panic("Only master chain can create blocks")
    }

    txns := self.UnconfirmedTxns.RawTxns()

    //TODO: sort by arrival time/announce time
    sort.Sort(txns) //sort by arrival time
    nTxns := len(txns)

    //TODO: filter valid first
    if nTxns > self.Config.TransactionsPerBlock {
        txns = txns[:self.Config.TransactionsPerBlock]
    }
    b, err := self.blockchain.NewBlockFromTransactions(txns,
        self.Config.BlockCreationInterval)
    if err != nil {
        return b, err
    }
    return b, err
}

func (self *Blockchain) SignBlock() (block coin.Block) (SignedBlock, error) {
    if self.Config.SecKey == {} {
        log.Panic("Only master chain can create blocks")
    }

    return self.signBlock(b), nil
}

//InjectBLock inputs a new block and applies it against the block chain
// state if it is valid
func (self *Blockchain) InjectBlock(b SignedBlock) (error) {

    if err := self.verifySignedBlock(&b); err != nil {
        return err
    }

    if b.Block.Seq +1 != b.Block.Header.BkSeq {
        return errors.New("Out of Sequence Block")
    }

    //apply block against blockchain
    if err := self.blockchain.ExecuteBlock(b.Block); err != nil {
        return err
    }

    self.blockSigs.record(&b) //save block to disc
    return nil
}


// Should only need 1 block at time
// Returns N signed blocks more recent than Seq. Does not return nil.
func (self *Blockchain) GetSignedBlocksSince(seq, ct uint64) []SignedBlock {
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
func (self *Blockchain) GetGenesisBlock() SignedBlock {
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
func (self *Blockchain) MostRecentBkSeq() uint64 {
    h := self.blockchain.Head()
    return h.Header.BkSeq
}

// Returns descriptive coin.Blockchain information
func (self *Blockchain) GetBlockchainMetadata() BlockchainMetadata {
    return NewBlockchainMetadata(self)
}

// Returns a readable copy of the block at seq. Returns error if seq out of range
func (self *Blockchain) GetReadableBlock(seq uint64) (ReadableBlock, error) {
    if b, err := self.GetBlock(seq); err == nil {
        return NewReadableBlock(&b), nil
    } else {
        return ReadableBlock{}, err
    }
}

// Returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
func (self *Blockchain) GetReadableBlocks(start, end uint64) []ReadableBlock {
    blocks := self.GetBlocks(start, end)
    rbs := make([]ReadableBlock, 0, len(blocks))
    for _, b := range blocks {
        rbs = append(rbs, NewReadableBlock(&b))
    }
    return rbs
}

// Returns a copy of the block at seq. Returns error if seq out of range
func (self *Blockchain) GetBlock(seq uint64) (coin.Block, error) {
    var b coin.Block
    if seq >= uint64(len(self.blockchain.Blocks)) {
        return b, errors.New("Block seq out of range")
    }
    return self.blockchain.Blocks[seq], nil
}

// Returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
func (self *Blockchain) GetBlocks(start, end uint64) []coin.Block {
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
func (self *Blockchain) SetTxnAnnounce(h coin.SHA256, t time.Time) {
    self.UnconfirmedTxns.SetAnnounced(h, t)
}


//InjectTransaction makes the blockchain server aware of raw transactions
//InjectTransaction inserts the transaction into the unconfirmed set
func (self *Blockchain) InjectTransaction(txn coin.Transaction) (error) {

    //should not be doing verification here
    //verification against blockchain will fail if user creates
    //output that spends outputs created by unspent transctions
    //should check signature validity     
    if err := self.blockchain.VerifyTransaction(t); err != nil {
        return err
    }

    self.UnconfirmedTxns.RecordTxn(txn, didAnnounce)
}

func (self *Blockchain) verifySignedBlock(b *SignedBlock) error {
    return coin.VerifySignature(self.Config.PubKey, b.Sig,
        b.Block.HashHeader())
}

// Signs a block for master.  Will panic if anything is invalid
func (self *Blockchain) signBlock(b coin.Block) SignedBlock {
    if !self.Config.IsMaster {
        log.Panic("Only master chain can sign blocks")
    }
    sig := coin.SignHash(b.HashHeader(), self.Config.SecKey)
    sb := SignedBlock{
        Block: b,
        Sig:   sig,
    }
    return sb
}

// Loads a coin.Blockchain from disk
// DANGER, safer to reparse whole block chain

/*
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
*/

/*
// Loads a blockchain but subdues errors into the logger, or panics.
// If no blockchain is found, it creates a new empty one
func loadBlockchain(filename string) *coin.Blockchain {
    bc := &coin.Blockchain{}
    created := false
    if filename != "" {
        var err error
        bc, err = LoadBlockchain(filename)
        if err == nil {
            if len(bc.Blocks) == 0 {
                log.Panic("Loaded empty blockchain")
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
*/

// Saves blockchain to disk
// Safer to reparse whole block chain

/*
func SaveBlockchain(bc *coin.Blockchain, filename string) error {
    // TODO -- blockchain file must be forward compatible
    // No -- If chain 
    data := encoder.Serialize(bc)
    return util.SaveBinary(filename, data, 0644)
}
*/
