package visor

import (
	"errors"
	"fmt"

	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("visor")
)

// BuildInfo represents the build info
type BuildInfo struct {
	Version string `json:"version"` // version number
	Commit  string `json:"commit"`  // git commit id
}

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
	// bolt db file path
	DBPath string
	// enable arbitrating mode
	Arbitrating bool
	// wallet directory
	WalletDirectory string
	// build info, including version, build time etc.
	BuildInfo BuildInfo
}

// NewVisorConfig put cap on block size, not on transactions/block
//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewVisorConfig() Config {
	c := Config{
		IsMaster: false,

		BlockchainPubkey: cipher.PubKey{},
		BlockchainSeckey: cipher.SecKey{},

		BlockCreationInterval: 10,
		//BlockCreationForceInterval: 120, //create block if no block within this many seconds

		UnconfirmedCheckInterval: time.Hour * 2,
		UnconfirmedMaxAge:        time.Hour * 48,
		UnconfirmedRefreshRate:   time.Minute,
		UnconfirmedResendPeriod:  time.Minute,
		MaxBlockSize:             1024 * 32,

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
	// blockSigs   *blockdb.BlockSigs
	history  *historydb.HistoryDB
	bcParser *BlockchainParser
	Wallets  *wallet.Service
	DB       *dbutil.DB

	done chan struct{}
}

// NewVisor Creates a normal Visor given a master's public key
func NewVisor(c Config) (*Visor, error) {
	logger.Debug("Creating new visor")
	// Make sure inputs are correct
	if c.IsMaster {
		logger.Debug("Visor is master")
		if c.BlockchainPubkey != cipher.PubKeyFromSecKey(c.BlockchainSeckey) {
			// logger.Panicf("Cannot run in master: invalid seckey for pubkey")
			return nil, errors.New("Cannot run in master: invalid seckey for pubkey")
		}
	}

	db, bc, err := loadBlockchain(c.DBPath, c.BlockchainPubkey, BlockchainOptions{
		Arbitrating: c.Arbitrating,
	})
	if err != nil {
		return nil, err
	}

	history, err := historydb.New(db)
	if err != nil {
		return nil, err
	}

	// creates blockchain parser instance
	bp := NewBlockchainParser(db, history, bc)

	bc.BindListener(bp.FeedBlock)

	wltServ, err := wallet.NewService(c.WalletDirectory)
	if err != nil {
		return nil, err
	}

	unconfirmed, err := NewUnconfirmedTxnPool(db)
	if err != nil {
		return nil, err
	}

	v := &Visor{
		Config:      c,
		DB:          db,
		Blockchain:  bc,
		Unconfirmed: unconfirmed,
		history:     history,
		bcParser:    bp,
		Wallets:     wltServ,
		done:        make(chan struct{}),
	}

	return v, nil
}

// Run starts the visor
func (vs *Visor) Run() error {
	defer close(vs.done)

	if err := vs.DB.Update(func(tx *bolt.Tx) error {
		logger.Debug("Visor.Run: maybeCreateGenesisBlock")
		if err := vs.maybeCreateGenesisBlock(tx); err != nil {
			return err
		}

		logger.Debug("Visor.Run: processUnconfirmedTxns")
		if err := vs.processUnconfirmedTxns(tx); err != nil {
			return err
		}

		logger.Debug("Visor.Run: BlockchainParser.Init")
		return vs.bcParser.Init(tx)
	}); err != nil {
		return err
	}

	return vs.bcParser.Run()
}

// Shutdown shuts down the visor
func (vs *Visor) Shutdown() {
	defer logger.Info("DB and BlockchainParser closed")

	vs.bcParser.Shutdown()

	<-vs.done

	logger.Info("Closing bolt.DB")
	if err := vs.DB.Close(); err != nil {
		logger.Error("db.Close() error: %v", err)
	}
}

// maybeCreateGenesisBlock creates a genesis block if necessary
func (vs *Visor) maybeCreateGenesisBlock(tx *bolt.Tx) error {
	if vs.Blockchain.GetGenesisBlock() != nil {
		return nil
	}

	logger.Debug("Create genesis block")
	vs.GenesisPreconditions()
	b, err := coin.NewGenesisBlock(vs.Config.GenesisAddress, vs.Config.GenesisCoinVolume, vs.Config.GenesisTimestamp)
	if err != nil {
		return err
	}

	var sb coin.SignedBlock
	// record the signature of genesis block
	if vs.Config.IsMaster {
		sb = vs.SignBlock(*b)
		logger.Info("Genesis block signature=%s", sb.Sig.Hex())
	} else {
		sb = coin.SignedBlock{
			Block: *b,
			Sig:   vs.Config.GenesisSignature,
		}
	}

	return vs.executeSignedBlock(tx, sb)
}

// check if there're unconfirmed transactions that are actually
// already executed, and remove them if any
func (vs *Visor) processUnconfirmedTxns(tx *bolt.Tx) error {
	var removeTxs []cipher.SHA256

	if err := vs.Unconfirmed.ForEach(tx, func(hash cipher.SHA256, utxn UnconfirmedTxn) error {
		// check if the tx already executed
		if err := vs.Blockchain.VerifyTransaction(tx, utxn.Txn); err != nil {
			removeTxs = append(removeTxs, hash)
		}

		txn, err := vs.history.GetTransaction(tx, hash)
		if err != nil {
			return err
		}

		if txn != nil {
			removeTxs = append(removeTxs, hash)
		}

		return nil
	}); err != nil {
		return err
	}

	if len(removeTxs) > 0 {
		return vs.Unconfirmed.RemoveTransactions(tx, removeTxs)
	}

	return nil
}

// GenesisPreconditions panics if conditions for genesis block are not met
func (vs *Visor) GenesisPreconditions() {
	//if seckey is set
	if vs.Config.BlockchainSeckey != (cipher.SecKey{}) {
		if vs.Config.BlockchainPubkey != cipher.PubKeyFromSecKey(vs.Config.BlockchainSeckey) {
			logger.Panicf("Cannot create genesis block. Invalid secret key for pubkey")
		}
	}
}

// RefreshUnconfirmed checks unconfirmed txns against the blockchain and returns
// all transaction that turn to valid.
func (vs *Visor) RefreshUnconfirmed() ([]cipher.SHA256, error) {
	logger.Debug("visor.Visor.RefreshUnconfirmed")
	var hashes []cipher.SHA256

	if err := vs.DB.Update(func(tx *bolt.Tx) error {
		var err error
		hashes, err = vs.Unconfirmed.Refresh(tx, vs.Blockchain)
		return err
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// CreateAndExecuteBlock creates a SignedBlock from pending transactions and executes it
func (vs *Visor) CreateAndExecuteBlock() (coin.SignedBlock, error) {
	var sb coin.SignedBlock

	if err := vs.DB.Update(func(tx *bolt.Tx) error {
		var err error
		sb, err = vs.createBlock(tx, uint64(utc.UnixNow()))
		if err != nil {
			return err
		}

		return vs.executeSignedBlock(tx, sb)
	}); err != nil {
		return coin.SignedBlock{}, err
	}

	return sb, nil
}

// ExecuteSignedBlock adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (vs *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	return vs.DB.Update(func(tx *bolt.Tx) error {
		return vs.executeSignedBlock(tx, b)
	})
}

// createBlock creates a SignedBlock from pending transactions
func (vs *Visor) createBlock(tx *bolt.Tx, when uint64) (coin.SignedBlock, error) {
	if !vs.Config.IsMaster {
		logger.Panic("Only master chain can create blocks")
	}

	length, err := vs.Unconfirmed.Len(tx)
	if err != nil {
		return coin.SignedBlock{}, err
	}

	if length == 0 {
		return coin.SignedBlock{}, errors.New("No transactions")
	}

	txns, err := vs.Unconfirmed.RawTxns(tx)
	if err != nil {
		return coin.SignedBlock{}, err
	}

	head, err := vs.Blockchain.Head(tx)
	if err != nil {
		return coin.SignedBlock{}, err
	}

	txns = coin.SortTransactions(txns, vs.Blockchain.TransactionFee(tx, head.Time()))
	txns = txns.TruncateBytesTo(vs.Config.MaxBlockSize)

	b, err := vs.Blockchain.NewBlock(tx, txns, when)
	if err != nil {
		return coin.SignedBlock{}, err
	}

	return vs.SignBlock(*b), nil
}

// executeSignedBlock adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (vs *Visor) executeSignedBlock(tx *bolt.Tx, b coin.SignedBlock) error {
	if err := vs.verifySignedBlock(&b); err != nil {
		return err
	}

	if err := vs.Blockchain.ExecuteBlock(tx, &b); err != nil {
		return err
	}

	// Remove the transactions in the Block from the unconfirmed pool
	txHashes := make([]cipher.SHA256, 0, len(b.Block.Body.Transactions))
	for _, tx := range b.Block.Body.Transactions {
		txHashes = append(txHashes, tx.Hash())
	}

	if err := vs.Unconfirmed.RemoveTransactions(tx, txHashes); err != nil {
		return err
	}

	vs.Blockchain.Notify(b.Block)
	return nil
}

// Returns an error if the cipher.Sig is not valid for the coin.Block
func (vs *Visor) verifySignedBlock(b *coin.SignedBlock) error {
	return cipher.VerifySignature(vs.Config.BlockchainPubkey, b.Sig, b.Block.HashHeader())
}

// SignBlock signs a block for master.  Will panic if anything is invalid
func (vs *Visor) SignBlock(b coin.Block) coin.SignedBlock {
	if !vs.Config.IsMaster {
		logger.Panic("Only master chain can sign blocks")
	}
	sig := cipher.SignHash(b.HashHeader(), vs.Config.BlockchainSeckey)
	sb := coin.SignedBlock{
		Block: b,
		Sig:   sig,
	}
	return sb
}

// GetUnspentOutputs returns all unspent outputs from the unspent pool
func (vs *Visor) GetUnspentOutputs() (coin.UxArray, error) {
	var uxa coin.UxArray

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		uxa, err = vs.Blockchain.Unspent().GetAll(tx)
		return err
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}

// GetUnspentOutput returns an unspent outputs from the unspent pool by hash
func (vs *Visor) GetUnspentOutput(h cipher.SHA256) (*coin.UxOut, error) {
	var ux *coin.UxOut

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		ux, err = vs.Blockchain.Unspent().Get(tx, h)
		return err
	}); err != nil {
		return nil, err
	}

	return ux, nil
}

// UnconfirmedSpendingOutputs returns all spending outputs in unconfirmed tx pool
func (vs *Visor) UnconfirmedSpendingOutputs() (coin.UxArray, error) {
	var uxa coin.UxArray

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var inputs []cipher.SHA256
		txns, err := vs.Unconfirmed.RawTxns(tx)
		if err != nil {
			return err
		}

		for _, txn := range txns {
			inputs = append(inputs, txn.In...)
		}

		uxa, err = vs.Blockchain.Unspent().GetArray(tx, inputs)
		return err
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}

// UnconfirmedIncomingOutputs returns all predicted outputs that are in pending tx pool
func (vs *Visor) UnconfirmedIncomingOutputs() (coin.UxArray, error) {
	var uxa coin.UxArray

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		head, err := vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		uxa, err = vs.Unconfirmed.GetIncomingOutputs(tx, head.Head)
		return err
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}

// GetSignedBlocksSince returns N signed blocks more recent than Seq. Does not return nil.
func (vs *Visor) GetSignedBlocksSince(seq, ct uint64) ([]coin.SignedBlock, error) {
	var blocks []coin.SignedBlock

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		avail := uint64(0)
		head, err := vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		headSeq := head.Seq()
		if headSeq > seq {
			avail = headSeq - seq
		}
		if avail < ct {
			ct = avail
		}
		if ct == 0 {
			return nil
		}

		blocks = make([]coin.SignedBlock, 0, ct)
		for j := uint64(0); j < ct; j++ {
			i := seq + 1 + j
			b, err := vs.Blockchain.GetSignedBlockBySeq(tx, i)
			if err != nil {
				return err
			}

			blocks = append(blocks, *b)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return blocks, nil
}

// HeadBkSeq returns the highest BkSeq we know, returns -1 if the chain is empty
func (vs *Visor) HeadBkSeq() uint64 {
	return vs.Blockchain.HeadSeq()
}

// GetBlockchainMetadata returns descriptive Blockchain information
func (vs *Visor) GetBlockchainMetadata() (*BlockchainMetadata, error) {
	var head *coin.SignedBlock
	var unconfirmedLen, unspentsLen uint64

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		head, err = vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		unconfirmedLen, err = vs.Unconfirmed.Len(tx)
		if err != nil {
			return err
		}

		unspentsLen, err = vs.Blockchain.Unspent().Len(tx)
		return err
	}); err != nil {
		return nil, err
	}

	return NewBlockchainMetadata(head, unconfirmedLen, unspentsLen)
}

// GetBlock returns a copy of the block at seq. Returns error if seq out of range
// Move to blockdb
func (vs *Visor) GetBlock(seq uint64) (*coin.SignedBlock, error) {
	if seq > vs.Blockchain.HeadSeq() {
		return nil, errors.New("Block seq out of range")
	}

	var b *coin.SignedBlock

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		b, err = vs.Blockchain.GetSignedBlockBySeq(tx, seq)
		return err
	}); err != nil {
		return nil, err
	}

	return b, nil
}

// GetBlocks returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
// move to blockdb
func (vs *Visor) GetBlocks(start, end uint64) ([]coin.SignedBlock, error) {
	var blocks []coin.SignedBlock

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		blocks, err = vs.Blockchain.GetBlocks(tx, start, end)
		return err
	}); err != nil {
		return nil, err
	}

	return blocks, nil
}

// InjectTransaction records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
func (vs *Visor) InjectTransaction(txn coin.Transaction) (bool, error) {
	var known bool

	if err := vs.DB.Update(func(tx *bolt.Tx) error {
		var err error
		known, err = vs.Unconfirmed.InjectTransaction(tx, vs.Blockchain, txn)
		return err
	}); err != nil {
		return false, err
	}

	return known, nil
}

// GetAddressTxns returns the Transactions whose unspents give coins to a cipher.Address.
// This includes unconfirmed txns' predicted unspents.
func (vs *Visor) GetAddressTxns(a cipher.Address) ([]Transaction, error) {
	var txns []Transaction

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		mxSeq := vs.HeadBkSeq()
		txs, err := vs.history.GetAddrTxns(tx, a)
		if err != nil {
			return err
		}

		for _, tx := range txs {
			h := mxSeq - tx.BlockSeq + 1

			bk, err := vs.GetSignedBlockBySeq(tx.BlockSeq)
			if err != nil {
				return err
			}

			if bk == nil {
				return fmt.Errorf("No block exists in depth: %d", tx.BlockSeq)
			}

			txns = append(txns, Transaction{
				Txn:    tx.Tx,
				Status: NewConfirmedTransactionStatus(h, tx.BlockSeq),
				Time:   bk.Time(),
			})
		}

		// Look in the unconfirmed pool
		uxs, err := vs.Unconfirmed.GetUnspentsOfAddr(tx, a)
		if err != nil {
			return err
		}

		for _, ux := range uxs {
			utxn, err := vs.Unconfirmed.Get(tx, ux.Body.SrcTransaction)
			if err != nil {
				return err
			}

			if utxn == nil {
				logger.Critical("Unconfirmed unspent missing unconfirmed txn")
				continue
			}

			txns = append(txns, Transaction{
				Txn:    utxn.Txn,
				Status: NewUnconfirmedTransactionStatus(),
				Time:   uint64(nanoToTime(utxn.Received).Unix()),
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return txns, nil
}

// GetTransaction returns a Transaction by hash.
func (vs *Visor) GetTransaction(txHash cipher.SHA256) (*Transaction, error) {
	var txn *Transaction

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		// Look in the unconfirmed pool
		utxn, err := vs.Unconfirmed.Get(tx, txHash)
		if err != nil {
			return err
		}

		if utxn != nil {
			txn = &Transaction{
				Txn:    utxn.Txn,
				Status: NewUnconfirmedTransactionStatus(),
				Time:   uint64(nanoToTime(utxn.Received).Unix()),
			}
			return nil
		}

		htxn, err := vs.history.GetTransaction(tx, txHash)
		if err != nil {
			return err
		}

		if htxn == nil {
			return nil
		}

		headSeq := vs.HeadBkSeq()

		b, err := vs.GetSignedBlockBySeq(htxn.BlockSeq)
		if err != nil {
			return err
		}

		if b == nil {
			return fmt.Errorf("found no block in seq %v", htxn.BlockSeq)
		}

		confirms := headSeq - htxn.BlockSeq + 1
		txn = &Transaction{
			Txn:    htxn.Tx,
			Status: NewConfirmedTransactionStatus(confirms, htxn.BlockSeq),
			Time:   b.Time(),
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return txn, nil
}

// AddressBalance computes the total balance for cipher.Addresses and their coin.UxOuts
func (vs *Visor) AddressBalance(head *coin.SignedBlock, auxs coin.AddressUxOuts) wallet.Balance {
	prevTime := head.Time()

	var coins, hours uint64
	for _, uxs := range auxs {
		for _, ux := range uxs {
			coins += ux.Body.Coins
			hours += ux.CoinHours(prevTime)
		}
	}

	return wallet.Balance{
		Coins: coins,
		Hours: hours,
	}
}

// GetUnconfirmedTxns gets all confirmed transactions of specific addresses
func (vs *Visor) GetUnconfirmedTxns(filter func(UnconfirmedTxn) bool) ([]UnconfirmedTxn, error) {
	var txns []UnconfirmedTxn

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		txns, err = vs.Unconfirmed.GetTxns(tx, filter)
		return err
	}); err != nil {
		return nil, err
	}

	return txns, nil
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
func (vs *Visor) GetAllUnconfirmedTxns() ([]UnconfirmedTxn, error) {
	var txns []UnconfirmedTxn

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		txns, err = vs.Unconfirmed.GetTxns(tx, All)
		return err
	}); err != nil {
		return nil, err
	}

	return txns, nil
}

// GetAllValidUnconfirmedTxHashes returns all valid unconfirmed transaction hashes
func (vs *Visor) GetAllValidUnconfirmedTxHashes() ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		hashes, err = vs.Unconfirmed.GetTxHashes(tx, IsValid)
		return err
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// GetSignedBlockByHash get block of specific hash header, return nil on not found.
func (vs *Visor) GetSignedBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error) {
	var sb *coin.SignedBlock

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		sb, err = vs.Blockchain.GetSignedBlockByHash(tx, hash)
		return err
	}); err != nil {
		return nil, err
	}

	return sb, nil
}

// GetSignedBlockBySeq get block of specific seq, return nil on not found.
func (vs *Visor) GetSignedBlockBySeq(seq uint64) (*coin.SignedBlock, error) {
	var b *coin.SignedBlock

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		b, err = vs.Blockchain.GetSignedBlockBySeq(tx, seq)
		return err
	}); err != nil {
		return nil, err
	}

	return b, nil
}

// GetLastBlocks returns last N blocks
func (vs *Visor) GetLastBlocks(num uint64) ([]coin.SignedBlock, error) {
	var blocks []coin.SignedBlock

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		blocks, err = vs.Blockchain.GetLastBlocks(tx, num)
		return err
	}); err != nil {
		return nil, err
	}

	return blocks, nil
}

// GetLastTxs returns last confirmed transactions, return nil if empty
func (vs *Visor) GetLastTxs() ([]*Transaction, error) {
	var txns []*Transaction

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		ltxns, err := vs.history.GetLastTxs(tx)
		if err != nil {
			return err
		}

		bh := vs.HeadBkSeq()
		var confirms uint64
		txns = make([]*Transaction, len(ltxns))

		for i, txn := range ltxns {
			confirms = uint64(bh) - txn.BlockSeq + 1
			b, err := vs.Blockchain.GetSignedBlockBySeq(tx, txn.BlockSeq)
			if err != nil {
				return err
			}

			if b == nil {
				return fmt.Errorf("found no block in seq %v", txn.BlockSeq)
			}

			txns[i] = &Transaction{
				Txn:    txn.Tx,
				Status: NewConfirmedTransactionStatus(confirms, txn.BlockSeq),
				Time:   b.Time(),
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return txns, nil
}

// GetHeadBlock gets head block.
func (vs Visor) GetHeadBlock() (*coin.SignedBlock, error) {
	var b *coin.SignedBlock

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		b, err = vs.Blockchain.Head(tx)
		return err
	}); err != nil {
		return nil, err
	}

	return b, nil
}

// GetHeadBlockTime returns the time of the head block.
func (vs Visor) GetHeadBlockTime() (uint64, error) {
	var t uint64

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		t, err = vs.Blockchain.Time(tx)
		return err
	}); err != nil {
		return 0, err
	}

	return t, nil
}

// GetUxOutByID gets UxOut by hash id.
func (vs Visor) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	var out *historydb.UxOut

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		out, err = vs.history.GetUxout(tx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return out, nil
}

// GetAddrUxOuts gets all the address affected UxOuts.
func (vs Visor) GetAddrUxOuts(address cipher.Address) ([]*historydb.UxOut, error) {
	var out []*historydb.UxOut

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		out, err = vs.history.GetAddrUxOuts(tx, address)
		return err
	}); err != nil {
		return nil, err
	}

	return out, nil
}

// RecvOfAddresses returns unconfirmed receiving uxouts of addresses
func (vs *Visor) RecvOfAddresses(addrs []cipher.Address) (coin.AddressUxOuts, error) {
	var uxouts coin.AddressUxOuts

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		head, err := vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		uxouts, err = vs.Unconfirmed.RecvOfAddresses(tx, head.Head, addrs)
		return err
	}); err != nil {
		return nil, err
	}

	return uxouts, nil
}

// GetIncomingOutputs returns all predicted outputs that are in pending tx pool
func (vs *Visor) GetIncomingOutputs() (coin.UxArray, error) {
	var uxa coin.UxArray

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		head, err := vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		uxa, err = vs.Unconfirmed.GetIncomingOutputs(tx, head.Head)
		return err
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}

// GetUnconfirmedTxn gets an unconfirmed transaction from the DB
func (vs *Visor) GetUnconfirmedTxn(hash cipher.SHA256) (*UnconfirmedTxn, error) {
	var txn *UnconfirmedTxn

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		txn, err = vs.Unconfirmed.Get(tx, hash)
		return err
	}); err != nil {
		return nil, err
	}

	return txn, nil
}

// GetUnconfirmedUnknown returns unconfirmed txn hashes with known ones removed
func (vs *Visor) GetUnconfirmedUnknown(txns []cipher.SHA256) ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		hashes, err = vs.Unconfirmed.GetUnknown(tx, txns)
		return err
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// GetUnconfirmedKnown returns unconfirmed txn hashes with known ones removed
func (vs *Visor) GetUnconfirmedKnown(txns []cipher.SHA256) (coin.Transactions, error) {
	var hashes coin.Transactions

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		hashes, err = vs.Unconfirmed.GetKnown(tx, txns)
		return err
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// UnconfirmedSpendsOfAddresses returns all unconfirmed coin.UxOut spends of addresses
func (vs *Visor) UnconfirmedSpendsOfAddresses(addrs []cipher.Address) (coin.AddressUxOuts, error) {
	var outs coin.AddressUxOuts

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		outs, err = vs.unconfirmedSpendsOfAddresses(tx, addrs)
		return err
	}); err != nil {
		return nil, err
	}

	return outs, nil
}

// unconfirmedSpendsOfAddresses returns all unconfirmed coin.UxOut spends of addresses
func (vs *Visor) unconfirmedSpendsOfAddresses(tx *bolt.Tx, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	txns, err := vs.Unconfirmed.RawTxns(tx)
	if err != nil {
		return nil, err
	}

	var inputs []cipher.SHA256
	for _, txn := range txns {
		inputs = append(inputs, txn.In...)
	}

	uxa, err := vs.Blockchain.Unspent().GetArray(tx, inputs)
	if err != nil {
		return nil, err
	}

	outs := make(coin.AddressUxOuts, len(addrs))

	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	for _, ux := range uxa {
		if _, ok := addrm[ux.Body.Address]; ok {
			outs[ux.Body.Address] = append(outs[ux.Body.Address], ux)
		}
	}

	return outs, nil
}

// SetTxnsAnnounced updates announced time of specific tx
func (vs *Visor) SetTxnsAnnounced(txns []cipher.SHA256, t time.Time) error {
	return vs.DB.Update(func(tx *bolt.Tx) error {
		return vs.Unconfirmed.SetTxnsAnnounced(tx, txns, t)
	})
}

// GetBalanceOfAddrs returns balance pair for multiple addresses
func (vs *Visor) GetBalanceOfAddrs(addrs []cipher.Address) (wallet.BalancePair, error) {
	var head *coin.SignedBlock
	var allUxs, spendUxs, recvUxs coin.AddressUxOuts

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		head, err = vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		allUxs, err = vs.Blockchain.Unspent().GetUnspentsOfAddrs(tx, addrs)
		if err != nil {
			return err
		}

		spendUxs, err = vs.unconfirmedSpendsOfAddresses(tx, addrs)
		if err != nil {
			err = fmt.Errorf("vs.unconfirmedSpendsOfAddresses failed: %v", err)
			return err
		}

		recvUxs, err = vs.Unconfirmed.RecvOfAddresses(tx, head.Head, addrs)
		if err != nil {
			err = fmt.Errorf("vs.recvOfAddresses failed: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return wallet.BalancePair{}, err
	}

	confirmed := vs.AddressBalance(head, allUxs)
	predicted := vs.AddressBalance(head, allUxs.Sub(spendUxs).Add(recvUxs))

	return wallet.BalancePair{
		Confirmed: confirmed,
		Predicted: predicted,
	}, nil
}

// GetTxnVerificationData returns data necessary for verifying a transaction
func (vs *Visor) GetTxnVerificationData(txn coin.Transaction) (coin.UxArray, uint64, error) {
	var uxa coin.UxArray
	var headTime uint64

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		uxa, err = vs.Blockchain.Unspent().GetArray(tx, txn.In)
		if err != nil {
			return err
		}

		headTime, err = vs.Blockchain.Time(tx)
		return err
	}); err != nil {
		return nil, 0, err
	}

	return uxa, headTime, nil
}

// GetUnspentsOfAddrs returns unspent outputs of multiple addresses
func (vs *Visor) GetUnspentsOfAddrs(addrs []cipher.Address) (coin.AddressUxOuts, error) {
	var uxa coin.AddressUxOuts

	if err := vs.DB.View(func(tx *bolt.Tx) error {
		var err error
		uxa, err = vs.Blockchain.Unspent().GetUnspentsOfAddrs(tx, addrs)
		return err
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}
