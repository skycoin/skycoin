package blockchain

import (
    "errors"
    "github.com/op/go-logging"
    //"github.com/skycoin/encoder"
    "github.com/skycoin/skycoin/src/coin"
    //"github.com/skycoin/skycoin/src/util"
    //"io/ioutil"
    "log"
    //"os"
    //"sort"
    "fmt"
    "time"
)

var (
    logger = logging.MustGetLogger("skycoin.Blockchain")
)

// Note: can use testnetpubkey as genesis address

type GenesisBlockCfg struct {
    GenesisAddress   cipher.Address
    GenesisSignature cipher.Sig
    GenesisTime      uint64
    GenesisCoins     uint64
    PubKey           cipher.PubKey
}

var TestNet GenesisBlockCfg
var MainNet GenesisBlockCfg //main blockchain

//testnet config
func init() {
    TestNet.PubKey = coin.MustPubKeyFromHex("025a3b22eb1e132a01f485119ae343342d92ab8599d9ad613a76e3b27f878bca8b")
    //TestNet.GenesisSignature = coin.MustSigFromHex()
    TestNet.GenesisAddress = cipher.MustDecodeBase58Address("26HbgWGwrToLZ6aX8VHtQmH4SPj4baQ5S3p")
    TestNet.GenesisTime = 1392584986 //set time
    TestNet.GenesisCoins = 1e12      //almost as many as Ripple
    //TestNet.GenesisSignature = coin.MustSigFromHex()
}

//main net config
func init() {
    MainNet.PubKey = coin.MustPubKeyFromHex("02bb0be2976457d2e30a9aea9b0057b0eb9d1ad6509ef743c25c737f24d6241a99")
    //TestNet.GenesisSignature = coin.MustSigFromHex()
    MainNet.GenesisAddress = cipher.MustDecodeBase58Address("26HbgWGwrToLZ6aX8VHtQmH4SPj4baQ5S3p")
    MainNet.GenesisTime = 1392584987 //set time
    MainNet.GenesisCoins = 100e6     //100 million
}

//var (
var BlockCreationInterval int = 15
var MaxTransactionSize int = 16 * 1024
var MaxBlockSize int = 32 * 1024
var MaxTransactionsPerBlock int = 1024

//)

// Configuration parameters for the Blockchain
type Blockchain struct {
    Genesis GenesisBlockCfg
    // Is this the master blockchain
    IsMaster bool
    SecKey   cipher.SecKey //set for writes
    // Use test network addresses
    TestNetwork bool

    Unconfirmed *UnconfirmedTxnPool
    blockchain  *coin.Blockchain
    Blocks      []SignedBlock
}

func (self *Blockchain) SetBlockchainSecKey(seed string) {
    pub, sec := coin.GenerateDeterministicKeyPair([]byte(seed))
    if pub != self.Genesis.PubKey {
        log.Panic("ERROR: pubkey does not correspond to loaded pubkey")
    }
    self.SecKey = sec
}

//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewBlockchain() Blockchain {
    //set pubkey based upon testnet, mainnet and local
    bc := Blockchain{

        IsMaster:    false, //writes blocks
        TestNetwork: true,

        //BlockCreationInterval:    15,
        //UnconfirmedCheckInterval: time.Minute * 5,
        //UnconfirmedMaxAge:        time.Hour * 48, //drop transaction not executed in 48 hours
        //UnconfirmedRefreshRate:   time.Minute * 30,
        //TransactionsPerBlock:     150, //10 transactions/second, 1.5 KB/s

        //BlockchainFile:           "",
        //BlockSigsFile:            "",

        SecKey: cipher.SecKey{},
    }
    bc.Blocks = make([]SignedBlock, 0) //save blocks to disc and only store head block
    bc.blockchain = coin.NewBlockchain()
    bc.Unconfirmed = NewUnconfirmedTxnPool()
    return bc
}

//NewTestnetBlockchain Config creates Blockchain for the testnet
func NewTestnetBlockchain() Blockchain {
    VC := NewBlockchain()
    VC.TestNetwork = true
    VC.Genesis = TestNet
    VC.InjectGenesisBlock()
    return VC
}

//NewTestnetBlockchain Config creates Blockchain for the mainnet
func NewMainnetBlockchain() Blockchain {
    VC := NewBlockchain()
    VC.TestNetwork = false
    VC.Genesis = MainNet
    VC.InjectGenesisBlock()
    return VC
}

//Generate Blockchain configuration for client only Blockchain, not intended to be synced to network
func NewLocalBlockchain() Blockchain {
    pubkey, seckey := cipher.GenerateKeyPair() //generate new/random pubkey/private key

    fmt.Printf("NewLocalBlockchain: genesis address seckey= %v \n", seckey.Hex())
    VC := NewBlockchain()
    VC.SecKey = seckey
    VC.Genesis.GenesisAddress = cipher.AddressFromPubKey(pubkey)
    VC.Genesis.GenesisTime = uint64(time.Now().Unix())
    VC.InjectGenesisBlock()
    return VC
}

/*
   Note: the genesis block does not need to be saved to disc
   Note: the genesis does not have signature because its implicit
   Note: the genesis block is part of block chain initialization
*/
func (self *Blockchain) InjectGenesisBlock() {
    var block coin.Block = self.blockchain.CreateGenesisBlock(self.Genesis.GenesisAddress, self.Genesis.GenesisTime, self.Genesis.GenesisCoins)
    _ = block //genesis block is automaticly applied to chain
}

// Checks unconfirmed txns against the blockchain and purges ones too old
func (self *Blockchain) RefreshUnconfirmed() {
    logger.Debug("Refreshing unconfirmed transactions")
    self.Unconfirmed.Refresh(self.blockchain,
        int(time.Minute*1)) // refresh every minute
}

//InjectTransaction makes the blockchain server aware of raw transactions
//InjectTransaction inserts the transaction into the unconfirmed set
// TODO: lock for thread safety
func (self *Blockchain) InjectTransaction(txn coin.Transaction) error {
    //strict filter would disallow transactions that cant be executed from unspent output set
    if txn.Size() > MaxTransactionSize { //16 KB/max size
        return errors.New("transaction over size limit")
    }
    if err := self.blockchain.VerifyTransaction(txn); err != nil {
        return err
    }
    self.Unconfirmed.RecordTxn(txn)
    return nil
}

func (self *Blockchain) PendingTransactions() bool {
    if len(self.Unconfirmed.RawTxns()) == 0 {
        return false
    }
    return true
}

// Creates a SignedBlock from pending transactions
// Applies transaction limit constraint
// Applies max block size constraint
// Should order transactions by priority
func (self *Blockchain) CreateBlock() (coin.Block, error) {
    //var sb SignedBlock
    if self.SecKey == (cipher.SecKey{}) {
        log.Panic("Only master chain can create blocks")
    }

    txns := self.Unconfirmed.RawTxns()

    //sort
    //arbritrate
    //truncate

    //TODO: sort by arrival time/announce time
    //TODO: filter valid first

    txns = coin.SortTransactions(txns, self.Blockchain.TransactionFee)
    txns = self.blockchain.ArbitrateTransactions(txns)
    nTxns := len(txns)
    if nTxns > MaxTransactionsPerBlock {
        txns = txns[:MaxTransactionsPerBlock]
    }
    txns = txns.TruncateBytesTo(MaxBlockSize) //cap at 32 KB
    //TODO: ERROR< NewBlockFromTransactions arbritates!
    b, err := self.blockchain.NewBlockFromTransactions(txns,
        uint64(time.Now().Unix()))
    //remove creation interval, from new block
    if err != nil {
        return b, err
    }
    return b, err
}

// Signs a block for master.  Will panic if anything is invalid
func (self *Blockchain) SignBlock(b coin.Block) SignedBlock {
    if !self.IsMaster {
        log.Panic("Only master chain can sign blocks")
    }
    sig := cipher.SignHash(b.HashHeader(), self.SecKey)
    sb := SignedBlock{
        Block: b,
        Sig:   sig,
    }
    return sb
}

//InjectBLock inputs a new block and applies it against the block chain
// state if it is valid
// TODO: lock for thread safety
func (self *Blockchain) InjectBlock(b SignedBlock) error {
    if err := self.verifySignedBlock(&b); err != nil {
        return err
    }

    if b.Block.Head.BkSeq+1 != self.blockchain.Head().Head.BkSeq {
        return errors.New("Out of Sequence Block")
    }

    //apply block against blockchain
    //this should not fail if signature is valid
    if _, err := self.blockchain.ExecuteBlock(b.Block); err != nil {
        return err
    }

    //self.blockSigs.record(&b) //save block to disc
    //self.Save()

    self.Blocks = append(self.Blocks, b)

    return nil
}

// Returns the highest BkSeq we know
// Replace with GetHead
func (self *Blockchain) MostRecentBkSeq() uint64 { //alread in meta
    h := self.blockchain.Head()
    return h.Head.BkSeq
}

// Returns a copy of the block at seq. Returns error if seq out of range
func (self *Blockchain) GetBlock(seq uint64) (coin.Block, error) {

    if seq == 0 {
        return self.blockchain.Blocks[0], nil
    }

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
func (self *Blockchain) SetTxnAnnounce(h cipher.SHA256) {
    self.Unconfirmed.SetAnnounced(h, time.Now().Unix())
}

func (self *Blockchain) verifySignedBlock(b *SignedBlock) error {
    return cipher.VerifySignature(self.Genesis.PubKey, b.Sig,
        b.Block.HashHeader())
}
