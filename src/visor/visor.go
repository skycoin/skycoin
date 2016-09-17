package visor

import (
	"errors"

	"gopkg.in/op/go-logging.v1"
	//"fmt"
	"log"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	//"github.com/skycoin/skycoin/src/wallet"
)

var (
	logger = logging.MustGetLogger("skycoin.visor")
)

// Configuration parameters for the Visor
type VisorConfig struct {
	// Is this the master blockchain
	IsMaster bool

	//WalletDirectory string //move out

	//Public key of blockchain authority
	BlockchainPubkey cipher.PubKey

	//Secret key of blockchain authority (if master)
	BlockchainSeckey cipher.SecKey

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
	//CoinHourBurnFactor uint64

	// Where the blockchain is saved
	BlockchainFile string
	// Where the block signatures are saved
	BlockSigsFile string

	//address for genesis
	GenesisAddress cipher.Address
	// Genesis block sig
	GenesisSignature cipher.Sig
	// Genesis block timestamp
	GenesisTimestamp uint64
	// Number of coins in genesis block
	GenesisCoinVolume uint64
	// Function that creates a new Wallet
	//WalletConstructor wallet.WalletConstructor
	// Default type of wallet to create
	//WalletTypeDefault wallet.WalletType
}

//Note, put cap on block size, not on transactions/block
//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewVisorConfig() VisorConfig {
	c := VisorConfig{
		IsMaster: false,

		//move wallet management out
		//WalletDirectory: "",

		//WalletConstructor: wallet.NewSimpleWallet,
		//WalletTypeDefault: wallet.SimpleWalletType,

		BlockchainPubkey: cipher.PubKey{},
		BlockchainSeckey: cipher.SecKey{},

		BlockCreationInterval: 10,
		//BlockCreationForceInterval: 120, //create block if no block within this many seconds

		UnconfirmedCheckInterval: time.Hour * 2,
		UnconfirmedMaxAge:        time.Hour * 48,
		UnconfirmedRefreshRate:   time.Minute * 30,
		MaxBlockSize:             1024 * 32,

		BlockchainFile: "",
		BlockSigsFile:  "",

		GenesisAddress:    cipher.Address{},
		GenesisSignature:  cipher.Sig{},
		GenesisTimestamp:  0,
		GenesisCoinVolume: 0, //100e12, 100e6 * 10e6
	}

	return c
}

// Manages the Blockchain as both a Master and a Normal
type Visor struct {
	Config VisorConfig
	// Unconfirmed transactions, held for relay until we get block confirmation
	Unconfirmed *UnconfirmedTxnPool
	// Wallets holding our keys for spending
	//Wallets wallet.Wallets
	// Master & personal keys
	Blockchain *coin.Blockchain
	blockSigs  *blockdb.BlockSigs
	txns       *blockdb.Transactions
}

// Creates a normal Visor given a master's public key
// func NewVisor2(c VisorConfig) *Visor {
// 	logger.Debug("Creating new visor")
// 	// Make sure inputs are correct
// 	if c.IsMaster {
// 		logger.Debug("Visor is master")
// 	}
// 	if c.IsMaster {
// 		if c.BlockchainPubkey != cipher.PubKeyFromSecKey(c.BlockchainSeckey) {
// 			log.Panicf("Cannot run in master: invalid seckey for pubkey")
// 		}
// 	}

// 	// Load the blockchain the block signatures
// 	blockchain := loadBlockchain(c.BlockchainFile, c.GenesisAddress)
// 	blockSigs, err := LoadBlockSigs(c.BlockSigsFile)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			logger.Info("BlockSigsFile \"%s\" not found", c.BlockSigsFile)
// 		} else {
// 			log.Panicf("Failed to load BlockSigsFile \"%s\"", c.BlockSigsFile)
// 		}
// 		blockSigs = NewBlockSigs()
// 	}

// 	v := &Visor{
// 		Config:      c,
// 		Blockchain:  blockchain,
// 		blockSigs:   blockSigs,
// 		Unconfirmed: NewUnconfirmedTxnPool(),
// 		//Wallets:     wallets,
// 	}
// 	// Load the genesis block and sign it, if we need one
// 	if len(blockchain.Blocks) == 0 {
// 		if (c.BlockchainSeckey == cipher.SecKey{}) || (c.IsMaster == false) {
// 			v.CreateGenesisBlock()
// 		} else {
// 			v.CreateGenesisBlockInit()
// 		}
// 	}

// 	err = blockSigs.Verify(c.BlockchainPubkey, blockchain)
// 	if err != nil {
// 		log.Panicf("Invalid block signatures: %v", err)
// 	}

// 	return v
// }

func walker(hps []coin.HashPair) cipher.SHA256 {
	return hps[0].Hash
}

// NewVisor Creates a normal Visor given a master's public key
func NewVisor(c VisorConfig) *Visor {
	logger.Debug("Creating new visor")
	// Make sure inputs are correct
	if c.IsMaster {
		logger.Debug("Visor is master")
		if c.BlockchainPubkey != cipher.PubKeyFromSecKey(c.BlockchainSeckey) {
			log.Panicf("Cannot run in master: invalid seckey for pubkey")
		}
	}

	tree := blockdb.NewBlockTree()
	bc := coin.NewBlockchain(tree, walker)
	v := &Visor{
		Config:      c,
		Blockchain:  bc,
		blockSigs:   blockdb.NewBlockSigs(),
		Unconfirmed: NewUnconfirmedTxnPool(),
		txns:        blockdb.NewTransactions(),
	}
	gb := bc.GetGenesisBlock()
	if gb == nil {
		v.GenesisPreconditions()
		b := v.Blockchain.CreateGenesisBlock(c.GenesisAddress, c.GenesisCoinVolume, c.GenesisTimestamp)
		gb = &b
		logger.Debug("create genesis block")

		// record the signature of genesis block
		if c.IsMaster {
			sb := v.SignBlock(*gb)
			v.blockSigs.Add(&sb)
		} else {
			v.blockSigs.Add(&coin.SignedBlock{
				Block: *gb,
				Sig:   c.GenesisSignature,
			})
		}

		// record transaction
		v.txns.Add(&blockdb.TransactionBlock{
			TxHash:    gb.Body.Transactions[0].Hash(),
			BlockHash: gb.HashHeader(),
		})
	}

	if err := v.blockSigs.Verify(c.BlockchainPubkey, v.Blockchain); err != nil {
		log.Panicf("Invalid block signatures: %v", err)
	}
	// gb := v.Blockchain.CreateGenesisBlock(c.GenesisAddress, c.GenesisTimestamp, c.GenesisCoinVolume)
	// var sb SignedBlock
	// if c.IsMaster {
	// 	sb = v.SignBlock(gb)
	// } else {
	// 	sb = SignedBlock{
	// 		Block: gb,
	// 		Sig:   c.GenesisSignature,
	// 	}
	// }
	// v.blockSigs.record(&sb)

	// // check if the genesis block does exist in blockdb.
	// block := blockdb.GetBlock(gb.HashHeader())
	// if block == nil {
	// 	// record the genesis block into blockdb.
	// 	dbBlock := blockdb.Block{
	// 		Block: gb,
	// 	}
	// 	if err := blockdb.SetBlock(dbBlock); err != nil {
	// 		log.Panicf("write block into blockdb failed:%v", err)
	// 	}

	// 	// record the genesis block signature into blockdb.
	// 	if err := blockdb.SetBlockSignature(gb.HashHeader(), gb.Head.PrevHash, c.GenesisSignature, gb.Head.BkSeq); err != nil {
	// 		log.Panicf("write block signature into blockdb failed:%v", err)
	// 	}
	// 	return v
	// }

	// // restore blocks from blockdb
	// var emptyHash cipher.SHA256
	// nxtHash := block.NextHash
	// for {
	// 	if nxtHash == emptyHash {
	// 		break
	// 	}

	// 	// get next block.
	// 	b := blockdb.GetBlock(nxtHash)
	// 	v.Blockchain.ExecuteBlock(b.Block)

	// 	// get next block signature.
	// 	bs := blockdb.GetBlockSignature(nxtHash)
	// 	sb := SignedBlock{
	// 		Block: b.Block,
	// 		Sig:   bs.Sig,
	// 	}
	// 	v.blockSigs.record(&sb)

	// 	nxtHash = b.NextHash
	// }

	// if err := v.blockSigs.Verify(c.BlockchainPubkey, v.Blockchain); err != nil {
	// 	log.Panicf("Invalid block signatures: %v", err)
	// }
	return v
}

// Returns a Visor with minimum initialization necessary for empty blockchain
// access
func NewMinimalVisor(c VisorConfig) *Visor {
	return &Visor{
		Config:      c,
		blockSigs:   blockdb.NewBlockSigs(),
		Unconfirmed: NewUnconfirmedTxnPool(),
		//Wallets:     nil,
	}
}

//panics if conditions for genesis block are not met
func (self *Visor) GenesisPreconditions() {

	//if len(self.Blockchain.Blocks) != 0 || len(self.blockSigs.Sigs) != 0 {
	//	log.Panic("Blockchain already has genesis")
	//}

	//if seckey is set
	if self.Config.BlockchainSeckey != (cipher.SecKey{}) {
		if self.Config.BlockchainPubkey != cipher.PubKeyFromSecKey(self.Config.BlockchainSeckey) {
			log.Panicf("Cannot create genesis block. Invalid secret key for pubkey")
		}
	}

}

// func (self *Visor) CreateGenesisBlockInit() (SignedBlock, error) {
// 	self.GenesisPreconditions()

// 	if len(self.Blockchain.Blocks) != 0 || len(self.blockSigs.Sigs) != 0 {
// 		log.Panic("Blockchain already has genesis")
// 	}
// 	if self.Config.BlockchainPubkey != cipher.PubKeyFromSecKey(self.Config.BlockchainSeckey) {
// 		log.Panicf("Cannot create genesis block. Invalid secret key for pubkey")
// 	}

// 	gb := self.Blockchain.CreateGenesisBlock(self.Config.GenesisAddress,
// 		self.Config.GenesisTimestamp, self.Config.GenesisCoinVolume)
// 	sb := self.SignBlock(gb)
// 	if err := self.verifySignedBlock(&sb); err != nil {
// 		log.Panicf("Signed a fresh genesis block, but its invalid: %v", err)
// 	}
// 	self.blockSigs.record(&sb)

// 	log.Printf("New Genesis:")
// 	log.Printf("genesis_time= %v", sb.Block.Head.Time)
// 	log.Printf("genesis_address= %v", self.Config.GenesisAddress.String())
// 	log.Printf("genesis_signature= %v", sb.Sig.Hex())

// 	return sb, nil
// }

// Creates the genesis block as needed
// func (self *Visor) CreateGenesisBlock() SignedBlock {
// 	self.GenesisPreconditions()

// 	if len(self.Blockchain.Blocks) != 0 || len(self.blockSigs.Sigs) != 0 {
// 		log.Panic("Blockchain already has genesis")
// 	}
// 	//addr := self.Config.GenesisAddress
// 	b := self.Blockchain.CreateGenesisBlock(self.Config.GenesisAddress, self.Config.GenesisTimestamp,
// 		self.Config.GenesisCoinVolume)
// 	sb := SignedBlock{
// 		Block: b,
// 		Sig:   self.Config.GenesisSignature,
// 	}
// 	self.blockSigs.record(&sb)

// 	err := self.blockSigs.Verify(self.Config.BlockchainPubkey,
// 		self.Blockchain)
// 	if err != nil {
// 		log.Panicf("Cannot create genesis block, signature verification failed: %v", err)
// 	}
// 	return sb
// }

// Checks unconfirmed txns against the blockchain and purges ones too old
func (self *Visor) RefreshUnconfirmed() {
	//logger.Debug("Refreshing unconfirmed transactions")
	self.Unconfirmed.Refresh(self.Blockchain,
		self.Config.UnconfirmedCheckInterval, self.Config.UnconfirmedMaxAge)
}

// Saves the coin.Blockchain to disk
// func (self *Visor) SaveBlockchain() error {
// 	if self.Config.BlockchainFile == "" {
// 		return errors.New("No BlockchainFile location set")
// 	} else {
// 		return SaveBlockchain(self.Blockchain, self.Config.BlockchainFile)
// 	}
// }

// Saves BlockSigs to disk
// func (self *Visor) SaveBlockSigs() error {
// 	if self.Config.BlockSigsFile == "" {
// 		return errors.New("No BlockSigsFile location set")
// 	} else {
// 		return self.blockSigs.A(self.Config.BlockSigsFile)
// 	}
// }

// Creates a SignedBlock from pending transactions
func (self *Visor) CreateBlock(when uint64) (coin.SignedBlock, error) {
	var sb coin.SignedBlock
	if !self.Config.IsMaster {
		log.Panic("Only master chain can create blocks")
	}
	if len(self.Unconfirmed.Txns) == 0 {
		return sb, errors.New("No transactions")
	}
	txns := self.Unconfirmed.RawTxns()
	txns = coin.SortTransactions(txns, self.Blockchain.TransactionFee)
	txns = txns.TruncateBytesTo(self.Config.MaxBlockSize)
	b, err := self.Blockchain.NewBlockFromTransactions(txns, when)
	if err != nil {
		return sb, err
	}
	return self.SignBlock(b), nil
}

// Creates a SignedBlock from pending transactions and executes it
func (self *Visor) CreateAndExecuteBlock() (coin.SignedBlock, error) {
	sb, err := self.CreateBlock(uint64(util.UnixNow()))
	if err == nil {
		return sb, self.ExecuteSignedBlock(sb)
	} else {
		return sb, err
	}
}

// Adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (self *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	if err := self.verifySignedBlock(&b); err != nil {
		return err
	}

	if _, err := self.Blockchain.ExecuteBlock(b.Block); err != nil {
		return err
	}
	// TODO -- save them even if out of order, and execute later
	// But make sure all prechecking as possible is done
	// TODO -- check if bitcoin allows blocks to be receiving out of order
	self.blockSigs.Add(&b)

	// add transactions in the block to blockdb
	for _, tx := range b.Block.Body.Transactions {
		self.txns.Add(&blockdb.TransactionBlock{
			TxHash:    tx.Hash(),
			BlockHash: b.Block.HashHeader(),
		})
	}

	// Remove the transactions in the Block from the unconfirmed pool
	self.Unconfirmed.RemoveTransactions(self.Blockchain,
		b.Block.Body.Transactions)
	return nil
}

// Returns an error if the cipher.Sig is not valid for the coin.Block
func (vs *Visor) verifySignedBlock(b *coin.SignedBlock) error {
	return cipher.VerifySignature(vs.Config.BlockchainPubkey, b.Sig, b.Block.HashHeader())
}

// Signs a block for master.  Will panic if anything is invalid
func (vs *Visor) SignBlock(b coin.Block) coin.SignedBlock {
	if !vs.Config.IsMaster {
		log.Panic("Only master chain can sign blocks")
	}
	sig := cipher.SignHash(b.HashHeader(), vs.Config.BlockchainSeckey)
	sb := coin.SignedBlock{
		Block: b,
		Sig:   sig,
	}
	return sb
}

/*
	Return Data
*/

//Make local copy and update when block header changes
// update should lock
// isolate effect of threading
// call .Array() to get []UxOut array
func (self *Visor) GetUnspentOutputs() []coin.UxOut {
	uxs := self.Blockchain.GetUnspent()
	return uxs.Array()
}

func (self *Visor) GetUnspentOutputsMap() coin.UnspentPool {
	uxs := self.Blockchain.GetUnspent()
	return uxs
}

func (self *Visor) GetUnspentOutputReadables() []ReadableOutput {
	uxs := self.GetUnspentOutputs()
	rx_readables := make([]ReadableOutput, len(uxs))
	for i, ux := range uxs {
		rx_readables[i] = NewReadableOutput(ux)
	}
	return rx_readables
}

// Returns N signed blocks more recent than Seq. Does not return nil.
func (self *Visor) GetSignedBlocksSince(seq, ct uint64) []coin.SignedBlock {
	avail := uint64(0)
	headSeq := self.Blockchain.Head().Seq()
	if headSeq > seq {
		avail = headSeq - seq
	}
	if avail < ct {
		ct = avail
	}
	if ct == 0 {
		return []coin.SignedBlock{}
	}
	blocks := make([]coin.SignedBlock, 0, ct)
	for j := uint64(0); j < ct; j++ {
		i := seq + 1 + j
		b := self.Blockchain.GetBlockInDepth(i)
		if b == nil {
			return []coin.SignedBlock{}
		}
		sig, err := self.blockSigs.Get(b.HashHeader())
		if err != nil {
			return []coin.SignedBlock{}
		}

		blocks = append(blocks, coin.SignedBlock{
			Block: *b,
			Sig:   sig,
		})
	}
	return blocks
}

// Returns the signed genesis block. Panics if signature or block not found
func (self *Visor) GetGenesisBlock() coin.SignedBlock {
	b := self.Blockchain.GetGenesisBlock()
	if b == nil {
		log.Panic("No genesis signature")
	}

	sig, err := self.blockSigs.Get(b.HashHeader())
	if err != nil {
		log.Panic(err)
	}

	return coin.SignedBlock{
		Sig:   sig,
		Block: *b,
	}
}

// Returns the highest BkSeq we know
func (self *Visor) HeadBkSeq() uint64 {
	return self.Blockchain.Head().Seq()
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
// Move to blockdb
func (self *Visor) GetBlock(seq uint64) (coin.Block, error) {
	var b coin.Block
	if seq > self.Blockchain.Head().Head.BkSeq {
		return b, errors.New("Block seq out of range")
	}

	return *self.Blockchain.GetBlockInDepth(seq), nil
}

// Returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
// move to blockdb
func (self *Visor) GetBlocks(start, end uint64) []coin.Block {
	return self.Blockchain.GetBlocks(start, end)
}

// Updates an UnconfirmedTxn's Announce field
//func (self *Visor) SetAnnounced(h cipher.SHA256, t time.Time) {
//	self.Unconfirmed.SetAnnounced(h, t)
//}

// Records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
// TODO
// - rename InjectTransaction
// Refactor
// Why do does this return both error and bool
func (self *Visor) InjectTxn(txn coin.Transaction) (error, bool) {
	//addrs := self.Wallets.GetAddressSet()
	return self.Unconfirmed.InjectTxn(self.Blockchain, txn)
}

// Returns the Transactions whose unspents give coins to a cipher.Address.
// This includes unconfirmed txns' predicted unspents.
func (self *Visor) GetAddressTransactions(a cipher.Address) []Transaction {
	txns := make([]Transaction, 0)
	// Look in the blockchain
	uxs := self.Blockchain.GetUnspent().AllForAddress(a)
	mxSeq := self.HeadBkSeq()
	for _, ux := range uxs {
		bk := self.Blockchain.GetBlockInDepth(ux.Head.BkSeq)
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
func (self *Visor) GetTransaction(txHash cipher.SHA256) Transaction {
	// Look in the unconfirmed pool
	tx, ok := self.Unconfirmed.Txns[txHash]
	if ok {
		return Transaction{
			Txn:    tx.Txn,
			Status: NewUnconfirmedTransactionStatus(),
		}
	}

	for {
		// Look in the blockchain
		txBlock := self.txns.Get(txHash)
		if txBlock == nil {
			break
		}
		// get block
		b := self.Blockchain.GetBlock(txBlock.BlockHash)
		if b == nil {
			break
		}

		if txn, ok := b.GetTransaction(txBlock.TxHash); ok {
			return Transaction{
				Txn:    txn,
				Status: NewConfirmedTransactionStatus(self.HeadBkSeq() - b.Seq() + 1),
			}
		}
		break
	}

	// TODO -- this is extremely slow as it does a full blockchain scan
	// We need an index from txn hash to block.  At least an index per block
	// to its contained txns

	// b := self.Blockchain.FindBlock(func(key, value []byte) bool {
	// 	var block coin.Block
	// 	if err := encoder.DeserializeRaw(value, &block); err != nil {
	// 		return false
	// 	}
	// 	if _, ok := block.GetTransaction(txHash); ok {
	// 		return true
	// 	}
	// 	return false
	// })

	// if b != nil {
	// 	tx, _ := b.GetTransaction(txHash)
	// 	height := self.HeadBkSeq() - b.Head.BkSeq + 1
	// 	return Transaction{
	// 		Txn:    tx,
	// 		Status: NewConfirmedTransactionStatus(height),
	// 	}
	// }

	// Otherwise unknown
	return Transaction{
		Status: NewUnknownTransactionStatus(),
	}
}

// Computes the total balance for cipher.Addresses and their coin.UxOuts
func (self *Visor) AddressBalance(auxs coin.AddressUxOuts) (uint64, uint64) {
	prevTime := self.Blockchain.Time()
	//b := wallet.NewBalance(0, 0)
	var coins uint64 = 0
	var hours uint64 = 0
	for _, uxs := range auxs {
		for _, ux := range uxs {
			coins += ux.Body.Coins
			hours += ux.CoinHours(prevTime)
			// FIXME
			//b = b.Add(wallet.NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
		}
	}
	return coins, hours
}

func (self *Visor) GetWalletTransactions(addresses []cipher.Address) []ReadableUnconfirmedTxn {

	ret := make([]ReadableUnconfirmedTxn, 0)

	for _, unconfirmedTxn := range self.Unconfirmed.Txns {
		isRelatedTransaction := false

		for _, out := range unconfirmedTxn.Txn.Out {
			for _, address := range addresses {
				if out.Address == address {
					isRelatedTransaction = true
				}
				if isRelatedTransaction {
					break
				}
			}
		}

		if isRelatedTransaction == true {
			ret = append(ret, NewReadableUnconfirmedTxn(&unconfirmedTxn))
		}
	}

	return ret
}
