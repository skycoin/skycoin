package visor

import (
	"errors"
	"fmt"

	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

var (
	logger = util.MustGetLogger("visor")
)

// Config configuration parameters for the Visor
type Config struct {
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
	// How often to rebroadcast unconfirmed transactions
	UnconfirmedResendPeriod time.Duration
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
	DBPath      string
	Arbitrating bool // enable arbitrating
}

// NewVisorConfig put cap on block size, not on transactions/block
//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewVisorConfig() Config {
	c := Config{
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
		UnconfirmedRefreshRate:   time.Minute,
		// UnconfirmedRefreshRate:   time.Minute * 30,
		UnconfirmedResendPeriod: time.Minute,
		MaxBlockSize:            1024 * 32,

		GenesisAddress:    cipher.Address{},
		GenesisSignature:  cipher.Sig{},
		GenesisTimestamp:  0,
		GenesisCoinVolume: 0, //100e12, 100e6 * 10e6
	}

	return c
}

// Visor manages the Blockchain as both a Master and a Normal
type Visor struct {
	Config Config
	// Unconfirmed transactions, held for relay until we get block confirmation
	Unconfirmed *UnconfirmedTxnPool
	Blockchain  *Blockchain
	blockSigs   *blockdb.BlockSigs
	history     *historydb.HistoryDB
	bcParser    *BlockchainParser
}

func walker(hps []coin.HashPair) cipher.SHA256 {
	return hps[0].Hash
}

// open the blockdb.
func openDB(dbFile string) (*bolt.DB, func()) {
	// dbFile := filepath.Join(util.DataDir, dbpath)
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		panic(fmt.Errorf("Open boltdb failed, err:%v", err))
	}
	return db, func() {
		db.Close()
	}
}

// VsClose visor close function
type VsClose func()

// NewVisor Creates a normal Visor given a master's public key
func NewVisor(c Config) (*Visor, VsClose) {
	logger.Debug("Creating new visor")
	// Make sure inputs are correct
	if c.IsMaster {
		logger.Debug("Visor is master")
		if c.BlockchainPubkey != cipher.PubKeyFromSecKey(c.BlockchainSeckey) {
			log.Panicf("Cannot run in master: invalid seckey for pubkey")
		}
	}

	db, closeDB := openDB(c.DBPath)
	history, err := historydb.New(db)
	if err != nil {
		log.Panic(err)
	}

	tree := blockdb.NewBlockTree(db)
	bc := NewBlockchain(tree, walker, Arbitrating(c.Arbitrating))
	bp := NewBlockchainParser(history, bc)

	bc.BindListener(bp.BlockListener)

	bp.Start()

	v := &Visor{
		Config:      c,
		Blockchain:  bc,
		blockSigs:   blockdb.NewBlockSigs(db),
		Unconfirmed: NewUnconfirmedTxnPool(db),
		history:     history,
		bcParser:    bp,
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
			logger.Info("genesis block signature=%s", sb.Sig.Hex())
		} else {
			v.blockSigs.Add(&coin.SignedBlock{
				Block: *gb,
				Sig:   c.GenesisSignature,
			})
		}
	}

	if err := v.Blockchain.VerifySigs(c.BlockchainPubkey, v.blockSigs); err != nil {
		log.Panicf("Invalid block signatures: %v", err)
	}

	return v, func() {
		closeDB()
	}
}

// NewMinimalVisor returns a Visor with minimum initialization necessary for empty blockchain
// access
// func NewMinimalVisor(c VisorConfig) (*Visor {
// 	db, _ := openDB(c.DBPath)
// 	return &Visor{
// 		Config:      c,
// 		blockSigs:   blockdb.NewBlockSigs(db),
// 		Unconfirmed: NewUnconfirmedTxnPool(db),
// 		//Wallets:     nil,
// 	}
// }

// GenesisPreconditions panics if conditions for genesis block are not met
func (vs *Visor) GenesisPreconditions() {
	//if seckey is set
	if vs.Config.BlockchainSeckey != (cipher.SecKey{}) {
		if vs.Config.BlockchainPubkey != cipher.PubKeyFromSecKey(vs.Config.BlockchainSeckey) {
			log.Panicf("Cannot create genesis block. Invalid secret key for pubkey")
		}
	}
}

// RefreshUnconfirmed checks unconfirmed txns against the blockchain and returns
// all transaction that turn to valid.
func (vs *Visor) RefreshUnconfirmed() []cipher.SHA256 {
	return vs.Unconfirmed.Refresh(vs.Blockchain)
}

// CreateBlock creates a SignedBlock from pending transactions
func (vs *Visor) CreateBlock(when uint64) (coin.SignedBlock, error) {
	var sb coin.SignedBlock
	if !vs.Config.IsMaster {
		log.Panic("Only master chain can create blocks")
	}
	if vs.Unconfirmed.Txns.len() == 0 {
		return sb, errors.New("No transactions")
	}
	txns := vs.Unconfirmed.RawTxns()
	txns = coin.SortTransactions(txns, vs.Blockchain.TransactionFee)
	txns = txns.TruncateBytesTo(vs.Config.MaxBlockSize)
	b, err := vs.Blockchain.NewBlockFromTransactions(txns, when)
	if err != nil {
		return sb, err
	}
	return vs.SignBlock(b), nil
}

// CreateAndExecuteBlock creates a SignedBlock from pending transactions and executes it
func (vs *Visor) CreateAndExecuteBlock() (coin.SignedBlock, error) {
	sb, err := vs.CreateBlock(uint64(util.UnixNow()))
	if err == nil {
		return sb, vs.ExecuteSignedBlock(sb)
	}

	return sb, err
}

// ExecuteSignedBlock adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (vs *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	if err := vs.verifySignedBlock(&b); err != nil {
		return err
	}

	if _, err := vs.Blockchain.ExecuteBlock(&b.Block); err != nil {
		return err
	}
	// TODO -- save them even if out of order, and execute later
	// But make sure all prechecking as possible is done
	// TODO -- check if bitcoin allows blocks to be receiving out of order
	vs.blockSigs.Add(&b)

	// Remove the transactions in the Block from the unconfirmed pool
	vs.Unconfirmed.RemoveTransactions(vs.Blockchain, b.Block.Body.Transactions)
	return nil
}

// Returns an error if the cipher.Sig is not valid for the coin.Block
func (vs *Visor) verifySignedBlock(b *coin.SignedBlock) error {
	return cipher.VerifySignature(vs.Config.BlockchainPubkey, b.Sig, b.Block.HashHeader())
}

// SignBlock signs a block for master.  Will panic if anything is invalid
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

// GetUnspentOutputs makes local copy and update when block header changes
// update should lock
// isolate effect of threading
// call .Array() to get []UxOut array
func (vs *Visor) GetUnspentOutputs() []coin.UxOut {
	uxs := vs.Blockchain.GetUnspent()
	return uxs.Array()
}

// GetUnspentOutputsMap return unspent output map
func (vs *Visor) GetUnspentOutputsMap() coin.UnspentPool {
	uxs := vs.Blockchain.GetUnspent()
	return *uxs
}

// GetUnspentOutputReadables returns readable unspent outputs
func (vs *Visor) GetUnspentOutputReadables() []ReadableOutput {
	uxs := vs.GetUnspentOutputs()
	rxReadables := make([]ReadableOutput, len(uxs))
	for i, ux := range uxs {
		rxReadables[i] = NewReadableOutput(ux)
	}
	return rxReadables
}

// AllSpendsOutputs returns all spending outputs in unconfirmed tx pool
func (vs *Visor) AllSpendsOutputs() []ReadableOutput {
	return vs.Unconfirmed.AllSpendsOutputs(vs.Blockchain.GetUnspent())
}

// AllIncommingOutputs returns all predicted outputs that are in pending tx pool
func (vs *Visor) AllIncommingOutputs() []ReadableOutput {
	return vs.Unconfirmed.AllIncommingOutputs(vs.Blockchain.Head().Head)
}

// GetSignedBlocksSince returns N signed blocks more recent than Seq. Does not return nil.
func (vs *Visor) GetSignedBlocksSince(seq, ct uint64) []coin.SignedBlock {
	avail := uint64(0)
	headSeq := vs.Blockchain.Head().Seq()
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
		b := vs.Blockchain.GetBlockInDepth(i)
		if b == nil {
			return []coin.SignedBlock{}
		}
		sig, err := vs.blockSigs.Get(b.HashHeader())
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

// GetGenesisBlock returns the signed genesis block. Panics if signature or block not found
func (vs *Visor) GetGenesisBlock() coin.SignedBlock {
	b := vs.Blockchain.GetGenesisBlock()
	if b == nil {
		log.Panic("No genesis signature")
	}

	sig, err := vs.blockSigs.Get(b.HashHeader())
	if err != nil {
		log.Panic(err)
	}

	return coin.SignedBlock{
		Sig:   sig,
		Block: *b,
	}
}

// HeadBkSeq returns the highest BkSeq we know
func (vs *Visor) HeadBkSeq() uint64 {
	return vs.Blockchain.Head().Seq()
}

// GetBlockchainMetadata returns descriptive Blockchain information
func (vs *Visor) GetBlockchainMetadata() BlockchainMetadata {
	return NewBlockchainMetadata(vs)
}

// GetReadableBlock returns a readable copy of the block at seq. Returns error if seq out of range
func (vs *Visor) GetReadableBlock(seq uint64) (ReadableBlock, error) {
	b, err := vs.GetBlock(seq)
	if err != nil {
		return ReadableBlock{}, err
	}

	return NewReadableBlock(&b), nil
}

// GetReadableBlocks returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
func (vs *Visor) GetReadableBlocks(start, end uint64) []ReadableBlock {
	blocks := vs.GetBlocks(start, end)
	rbs := make([]ReadableBlock, 0, len(blocks))
	for _, b := range blocks {
		rbs = append(rbs, NewReadableBlock(&b))
	}
	return rbs
}

// GetBlock returns a copy of the block at seq. Returns error if seq out of range
// Move to blockdb
func (vs *Visor) GetBlock(seq uint64) (coin.Block, error) {
	var b coin.Block
	if seq > vs.Blockchain.Head().Head.BkSeq {
		return b, errors.New("Block seq out of range")
	}

	return *vs.Blockchain.GetBlockInDepth(seq), nil
}

// GetBlocks returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
// move to blockdb
func (vs *Visor) GetBlocks(start, end uint64) []coin.Block {
	return vs.Blockchain.GetBlocks(start, end)
}

// InjectTxn records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
// TODO
// - rename InjectTransaction
// Refactor
// Why do does this return both error and bool
func (vs *Visor) InjectTxn(txn coin.Transaction) (bool, error) {
	//addrs := self.Wallets.GetAddressSet()
	return vs.Unconfirmed.InjectTxn(vs.Blockchain, txn)
}

// GetAddressTransactions returns the Transactions whose unspents give coins to a cipher.Address.
// This includes unconfirmed txns' predicted unspents.
func (vs *Visor) GetAddressTransactions(a cipher.Address) []Transaction {
	var txns []Transaction
	// Look in the blockchain
	uxs := vs.Blockchain.GetUnspent().AllForAddress(a)
	mxSeq := vs.HeadBkSeq()
	var bk *coin.Block
	for _, ux := range uxs {
		if bk = vs.GetBlockBySeq(ux.Head.BkSeq); bk == nil {
			return txns
		}

		tx, ok := bk.GetTransaction(ux.Body.SrcTransaction)
		if ok {
			h := mxSeq - bk.Head.BkSeq + 1
			txns = append(txns, Transaction{
				Txn:    tx,
				Status: NewConfirmedTransactionStatus(h, bk.Head.BkSeq),
				Time:   bk.Time(),
			})
		}
	}

	// Look in the unconfirmed pool
	uxs = vs.Unconfirmed.Unspent.getAllForAddress(a)
	for _, ux := range uxs {
		tx, ok := vs.Unconfirmed.Txns.get(ux.Body.SrcTransaction)
		if !ok {
			logger.Critical("Unconfirmed unspent missing unconfirmed txn")
			continue
		}
		txns = append(txns, Transaction{
			Txn:    tx.Txn,
			Status: NewUnconfirmedTransactionStatus(),
			Time:   uint64(nanoToTime(tx.Received).Unix()),
		})
	}

	return txns
}

// GetTransaction returns a Transaction by hash.
func (vs *Visor) GetTransaction(txHash cipher.SHA256) (*Transaction, error) {
	// Look in the unconfirmed pool
	tx, ok := vs.Unconfirmed.Txns.get(txHash)
	if ok {
		return &Transaction{
			Txn:    tx.Txn,
			Status: NewUnconfirmedTransactionStatus(),
			Time:   uint64(nanoToTime(tx.Received).Unix()),
		}, nil
	}

	txn, err := vs.history.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}

	if txn == nil {
		return nil, nil
	}

	confirms := vs.GetHeadBlock().Seq() - txn.BlockSeq + 1
	b := vs.GetBlockBySeq(txn.BlockSeq)
	if b == nil {
		return nil, fmt.Errorf("found no block in seq %v", txn.BlockSeq)
	}

	return &Transaction{
		Txn:    txn.Tx,
		Status: NewConfirmedTransactionStatus(confirms, txn.BlockSeq),
		Time:   b.Time(),
	}, nil
}

// AddressBalance computes the total balance for cipher.Addresses and their coin.UxOuts
func (vs *Visor) AddressBalance(auxs coin.AddressUxOuts) (uint64, uint64) {
	prevTime := vs.Blockchain.Time()
	//b := wallet.NewBalance(0, 0)
	var coins uint64
	var hours uint64
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

// GetUnconfirmedTxns gets all confirmed transactions of specific addresses
func (vs *Visor) GetUnconfirmedTxns(filter func(UnconfirmedTxn) bool) []UnconfirmedTxn {
	return vs.Unconfirmed.GetTxns(filter)
}

// ToAddresses represents a filter that check if tx has output to the given addresses
func ToAddresses(addresses []cipher.Address) func(UnconfirmedTxn) bool {
	return func(tx UnconfirmedTxn) (isRelated bool) {
		for _, out := range tx.Txn.Out {
			for _, address := range addresses {
				if out.Address == address {
					isRelated = true
					return
				}
			}
		}
		return
	}
}

// GetAllUnconfirmedTxns returns all unconfirmed transactions
func (vs *Visor) GetAllUnconfirmedTxns() []UnconfirmedTxn {
	return vs.Unconfirmed.GetTxns(All)
}

// GetAllValidUnconfirmedTxHashes returns all valid unconfirmed transaction hashes
func (vs *Visor) GetAllValidUnconfirmedTxHashes() []cipher.SHA256 {
	return vs.Unconfirmed.GetTxHashes(IsValid)
}

// GetBlockByHash get block of specific hash header, return nil on not found.
func (vs *Visor) GetBlockByHash(hash cipher.SHA256) *coin.Block {
	return vs.Blockchain.GetBlock(hash)
}

// GetBlockBySeq get block of speicific seq, return nil on not found.
func (vs *Visor) GetBlockBySeq(seq uint64) *coin.Block {
	return vs.Blockchain.GetBlockInDepth(seq)
}

// GetLastTxs returns last confirmed transactions, return nil if empty
func (vs *Visor) GetLastTxs() ([]*Transaction, error) {
	ltxs, err := vs.history.GetLastTxs()
	if err != nil {
		return nil, err
	}

	txs := make([]*Transaction, len(ltxs))
	var confirms uint64
	bh := vs.GetHeadBlock().Seq()
	var b *coin.Block
	for i, tx := range ltxs {
		confirms = bh - tx.BlockSeq + 1
		if b = vs.GetBlockBySeq(tx.BlockSeq); b == nil {
			return nil, fmt.Errorf("found no block in seq %v", tx.BlockSeq)
		}

		txs[i] = &Transaction{
			Txn:    tx.Tx,
			Status: NewConfirmedTransactionStatus(confirms, tx.BlockSeq),
			Time:   b.Time(),
		}
	}
	return txs, nil
}

// GetHeadBlock gets head block.
func (vs Visor) GetHeadBlock() *coin.Block {
	return vs.Blockchain.Head()
}

// GetUxOutByID gets UxOut by hash id.
func (vs Visor) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	return vs.history.GetUxout(id)
}

// GetAddrUxOuts gets all the address affected UxOuts.
func (vs Visor) GetAddrUxOuts(address cipher.Address) ([]*historydb.UxOut, error) {
	return vs.history.GetAddrUxOuts(address)
}
