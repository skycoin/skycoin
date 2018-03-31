package visor

import (
	"errors"
	"fmt"
	"sort"

	"time"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	// MaxDropletPrecision represents the decimal precision of droplets
	MaxDropletPrecision uint64 = 3

	//DefaultMaxBlockSize is max block size
	DefaultMaxBlockSize int = 32 * 1024
)

var (
	logger = logging.MustGetLogger("visor")

	// errInvalidDecimals is returned by DropletPrecisionCheck if a coin amount has an invalid number of decimal places
	errInvalidDecimals = errors.New("invalid amount, too many decimal places")

	// maxDropletDivisor represents the modulus divisor when checking droplet precision rules.
	// It is computed from MaxDropletPrecision in init()
	maxDropletDivisor uint64
)

// MaxDropletDivisor represents the modulus divisor when checking droplet precision rules.
func MaxDropletDivisor() uint64 {
	// The value is wrapped in a getter to make it immutable to external packages
	return maxDropletDivisor
}

func init() {
	// Compute maxDropletDivisor from precision
	maxDropletDivisor = calculateDivisor(MaxDropletPrecision)
}

func calculateDivisor(precision uint64) uint64 {
	if precision > droplet.Exponent {
		logger.Panic("precision must be <= droplet.Exponent")
	}

	n := droplet.Exponent - precision
	var i uint64 = 1
	for k := uint64(0); k < n; k++ {
		i = i * 10
	}
	return i
}

// DropletPrecisionCheck checks if an amount of coins is valid given decimal place restrictions
func DropletPrecisionCheck(amount uint64) error {
	if amount%maxDropletDivisor != 0 {
		return errInvalidDecimals
	}
	return nil
}

// BuildInfo represents the build info
type BuildInfo struct {
	Version string `json:"version"` // version number
	Commit  string `json:"commit"`  // git commit id
}

// Config configuration parameters for the Visor
type Config struct {
	// Is this the master blockchain
	IsMaster bool

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
	// How often to check the unconfirmed pool for transactions that become valid
	UnconfirmedRefreshRate time.Duration
	// How often to remove transactions that become permanently invalid from the unconfirmed pool
	UnconfirmedRemoveInvalidRate time.Duration
	// How often to rebroadcast unconfirmed transactions
	UnconfirmedResendPeriod time.Duration
	// Maximum size of a block, in bytes.
	MaxBlockSize int

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
	// open bolt db read-only
	DBReadOnly bool
	// enable arbitrating mode
	Arbitrating bool
	// wallet directory
	WalletDirectory string
	// build info, including version, build time etc.
	BuildInfo BuildInfo
	// disables wallet API
	DisableWalletAPI bool
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

		UnconfirmedCheckInterval:     time.Hour * 2,
		UnconfirmedMaxAge:            time.Hour * 48,
		UnconfirmedRefreshRate:       time.Minute,
		UnconfirmedRemoveInvalidRate: time.Minute,
		UnconfirmedResendPeriod:      time.Minute,
		MaxBlockSize:                 DefaultMaxBlockSize,

		GenesisAddress:    cipher.Address{},
		GenesisSignature:  cipher.Sig{},
		GenesisTimestamp:  0,
		GenesisCoinVolume: 0, //100e12, 100e6 * 10e6
	}

	return c
}

// Verify verifies the configuration
func (c Config) Verify() error {
	if c.IsMaster {
		if c.BlockchainPubkey != cipher.PubKeyFromSecKey(c.BlockchainSeckey) {
			return errors.New("Cannot run in master: invalid seckey for pubkey")
		}
	}

	return nil
}

// historyer is the interface that provides methods for accessing history data that are parsed from blockchain.
type historyer interface {
	GetUxout(uxid cipher.SHA256) (*historydb.UxOut, error)
	ParseBlock(b *coin.Block) error
	GetTransaction(hash cipher.SHA256) (*historydb.Transaction, error)
	GetAddrUxOuts(address cipher.Address) ([]*historydb.UxOut, error)
	GetAddrTxns(address cipher.Address) ([]historydb.Transaction, error)
	ForEach(f func(tx *historydb.Transaction) error) error
	ResetIfNeed() error
	ParsedHeight() int64
}

// Blockchainer is the interface that provides methods for accessing the blockchain data
type Blockchainer interface {
	GetGenesisBlock() *coin.SignedBlock
	GetBlocks(start, end uint64) []coin.SignedBlock
	GetLastBlocks(n uint64) []coin.SignedBlock
	GetBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error)
	GetBlockBySeq(seq uint64) (*coin.SignedBlock, error)
	Unspent() blockdb.UnspentPool
	Len() uint64
	Head() (*coin.SignedBlock, error)
	HeadSeq() uint64
	Time() uint64
	NewBlock(txns coin.Transactions, currentTime uint64) (*coin.Block, error)
	ExecuteBlockWithTx(tx *bolt.Tx, sb *coin.SignedBlock) error
	VerifyBlockTxnConstraints(tx coin.Transaction) error
	VerifySingleTxnHardConstraints(tx coin.Transaction) error
	VerifySingleTxnAllConstraints(tx coin.Transaction, maxSize int) error
	TransactionFee(t *coin.Transaction) (uint64, error)
	Notify(b coin.Block)
	BindListener(bl BlockListener)
	UpdateDB(f func(tx *bolt.Tx) error) error
}

// UnconfirmedTxnPooler is the interface that provides methods for
// accessing the unconfirmed transaction pool
type UnconfirmedTxnPooler interface {
	SetAnnounced(hash cipher.SHA256, t time.Time) error
	InjectTransaction(bc Blockchainer, t coin.Transaction, maxSize int) (bool, *ErrTxnViolatesSoftConstraint, error)
	RawTxns() coin.Transactions
	RemoveTransactions(txns []cipher.SHA256) error
	RemoveTransactionsWithTx(tx *bolt.Tx, txns []cipher.SHA256)
	Refresh(bc Blockchainer, maxBlockSize int) ([]cipher.SHA256, error)
	RemoveInvalid(bc Blockchainer) ([]cipher.SHA256, error)
	FilterKnown(txns []cipher.SHA256) []cipher.SHA256
	GetKnown(txns []cipher.SHA256) coin.Transactions
	RecvOfAddresses(bh coin.BlockHeader, addrs []cipher.Address) (coin.AddressUxOuts, error)
	SpendsOfAddresses(addrs []cipher.Address, unspent blockdb.UnspentGetter) (coin.AddressUxOuts, error)
	GetSpendingOutputs(unspent blockdb.UnspentPool) (coin.UxArray, error)
	GetIncomingOutputs(bh coin.BlockHeader) coin.UxArray
	Get(hash cipher.SHA256) (*UnconfirmedTxn, bool)
	GetTxns(filter func(tx UnconfirmedTxn) bool) []UnconfirmedTxn
	GetTxHashes(filter func(tx UnconfirmedTxn) bool) []cipher.SHA256
	ForEach(f func(cipher.SHA256, *UnconfirmedTxn) error) error
	GetUnspentsOfAddr(addr cipher.Address) coin.UxArray
	Len() int
}

// Visor manages the Blockchain as both a Master and a Normal
type Visor struct {
	Config Config
	// Unconfirmed transactions, held for relay until we get block confirmation
	Unconfirmed UnconfirmedTxnPooler
	Blockchain  Blockchainer
	history     historyer
	bcParser    *BlockchainParser
	wallets     *wallet.Service
	db          *bolt.DB
}

// NewVisor creates a Visor for managing the blockchain database
func NewVisor(c Config, db *bolt.DB) (*Visor, error) {
	logger.Debug("Creating new visor")
	if c.IsMaster {
		logger.Debug("Visor is master")
	}

	if err := c.Verify(); err != nil {
		return nil, err
	}

	db, bc, err := loadBlockchain(db, c.BlockchainPubkey, c.Arbitrating)
	if err != nil {
		return nil, err
	}

	history, err := historydb.New(db)
	if err != nil {
		return nil, err
	}

	// creates blockchain parser instance
	bp := NewBlockchainParser(history, bc)

	bc.BindListener(bp.FeedBlock)

	wltServ, err := wallet.NewService(c.WalletDirectory, c.DisableWalletAPI)
	if err != nil {
		return nil, err
	}

	v := &Visor{
		Config:      c,
		db:          db,
		Blockchain:  bc,
		Unconfirmed: NewUnconfirmedTxnPool(db),
		history:     history,
		bcParser:    bp,
		wallets:     wltServ,
	}

	return v, nil
}

// Run starts the visor
func (vs *Visor) Run() error {
	if err := vs.maybeCreateGenesisBlock(); err != nil {
		return err
	}

	removed, err := vs.RemoveInvalidUnconfirmed()
	if err != nil {
		return err
	}
	logger.Infof("Removed %d invalid txns from pool", len(removed))

	return vs.bcParser.Run()
}

// Shutdown shuts down the visor
func (vs *Visor) Shutdown() {
	defer logger.Info("DB and BlockchainParser closed")

	vs.bcParser.Shutdown()

	if err := vs.db.Close(); err != nil {
		logger.Errorf("db.Close() error: %v", err)
	}
}

// maybeCreateGenesisBlock creates a genesis block if necessary
func (vs *Visor) maybeCreateGenesisBlock() error {
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
		logger.Infof("Genesis block signature=%s", sb.Sig.Hex())
	} else {
		sb = coin.SignedBlock{
			Block: *b,
			Sig:   vs.Config.GenesisSignature,
		}
	}

	return vs.ExecuteSignedBlock(sb)
}

// GenesisPreconditions panics if conditions for genesis block are not met
func (vs *Visor) GenesisPreconditions() {
	if vs.Config.BlockchainSeckey != (cipher.SecKey{}) {
		if vs.Config.BlockchainPubkey != cipher.PubKeyFromSecKey(vs.Config.BlockchainSeckey) {
			logger.Panic("Cannot create genesis block. Invalid secret key for pubkey")
		}
	}
}

// RefreshUnconfirmed checks unconfirmed txns against the blockchain and returns
// all transaction that turn to valid.
func (vs *Visor) RefreshUnconfirmed() ([]cipher.SHA256, error) {
	return vs.Unconfirmed.Refresh(vs.Blockchain, vs.Config.MaxBlockSize)
}

// RemoveInvalidUnconfirmed removes transactions that become permanently invalid
// (by violating hard constraints) from the pool.
// Returns the transaction hashes that were removed.
func (vs *Visor) RemoveInvalidUnconfirmed() ([]cipher.SHA256, error) {
	return vs.Unconfirmed.RemoveInvalid(vs.Blockchain)
}

// CreateBlock creates a SignedBlock from pending transactions
func (vs *Visor) CreateBlock(when uint64) (coin.SignedBlock, error) {
	if !vs.Config.IsMaster {
		logger.Panic("Only master chain can create blocks")
	}

	var sb coin.SignedBlock

	// Gather all unconfirmed transactions
	txns := vs.Unconfirmed.RawTxns()

	if len(txns) == 0 {
		return sb, errors.New("No transactions")
	}

	logger.Infof("Unconfirmed pool has %d transactions pending", len(txns))

	// Filter transactions that violate all constraints
	var filteredTxns coin.Transactions
	for _, txn := range txns {
		if err := vs.Blockchain.VerifySingleTxnAllConstraints(txn, vs.Config.MaxBlockSize); err != nil {
			logger.Warningf("Transaction %s violates constraints: %v", txn.TxIDHex(), err)
		} else {
			filteredTxns = append(filteredTxns, txn)
		}
	}

	nRemoved := len(txns) - len(filteredTxns)
	if nRemoved > 0 {
		logger.Infof("CreateBlock ignored %d transactions violating constraints", nRemoved)
	}

	txns = filteredTxns

	if len(txns) == 0 {
		logger.Info("No transactions after filtering for constraint violations")
		return sb, errors.New("No transactions after filtering for constraint violations")
	}

	// Sort them by highest fee per kilobyte
	txns = coin.SortTransactions(txns, vs.Blockchain.TransactionFee)

	// Apply block size transaction limit
	txns = txns.TruncateBytesTo(vs.Config.MaxBlockSize)

	if len(txns) == 0 {
		logger.Panic("TruncateBytesTo removed all transactions")
	}

	logger.Infof("Creating new block with %d transactions, head time %d", len(txns), when)

	b, err := vs.Blockchain.NewBlock(txns, when)
	if err != nil {
		logger.Warningf("Blockchain.NewBlock failed: %v", err)
		return sb, err
	}

	return vs.SignBlock(*b), nil
}

// CreateAndExecuteBlock creates a SignedBlock from pending transactions and executes it
func (vs *Visor) CreateAndExecuteBlock() (coin.SignedBlock, error) {
	sb, err := vs.CreateBlock(uint64(utc.UnixNow()))
	if err == nil {
		return sb, vs.ExecuteSignedBlock(sb)
	}

	return sb, err
}

// ExecuteSignedBlock adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (vs *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	if err := b.VerifySignature(vs.Config.BlockchainPubkey); err != nil {
		return err
	}

	if err := vs.db.Update(func(tx *bolt.Tx) error {
		if err := vs.Blockchain.ExecuteBlockWithTx(tx, &b); err != nil {
			return err
		}

		// Remove the transactions in the Block from the unconfirmed pool
		txHashes := make([]cipher.SHA256, 0, len(b.Block.Body.Transactions))
		for _, tx := range b.Block.Body.Transactions {
			txHashes = append(txHashes, tx.Hash())
		}
		vs.Unconfirmed.RemoveTransactionsWithTx(tx, txHashes)

		return nil
	}); err != nil {
		return err
	}

	vs.Blockchain.Notify(b.Block)
	return nil
}

// SignBlock signs a block for master.  Will panic if anything is invalid
func (vs *Visor) SignBlock(b coin.Block) coin.SignedBlock {
	if !vs.Config.IsMaster {
		logger.Panic("Only master chain can sign blocks")
	}

	sig := cipher.SignHash(b.HashHeader(), vs.Config.BlockchainSeckey)

	return coin.SignedBlock{
		Block: b,
		Sig:   sig,
	}
}

/*
	Return Data
*/

// GetUnspentOutputs makes local copy and update when block header changes
// update should lock
// isolate effect of threading
// call .Array() to get []UxOut array
func (vs *Visor) GetUnspentOutputs() ([]coin.UxOut, error) {
	return vs.Blockchain.Unspent().GetAll()
}

// UnconfirmedSpendingOutputs returns all spending outputs in unconfirmed tx pool
func (vs *Visor) UnconfirmedSpendingOutputs() (coin.UxArray, error) {
	return vs.Unconfirmed.GetSpendingOutputs(vs.Blockchain.Unspent())
}

// UnconfirmedIncomingOutputs returns all predicted outputs that are in pending tx pool
func (vs *Visor) UnconfirmedIncomingOutputs() (coin.UxArray, error) {
	head, err := vs.Blockchain.Head()
	if err != nil {
		return coin.UxArray{}, err
	}

	return vs.Unconfirmed.GetIncomingOutputs(head.Head), nil
}

// GetSignedBlocksSince returns signed blocks in an inclusive range of [seq+1, seq+ct]
func (vs *Visor) GetSignedBlocksSince(seq, ct uint64) ([]coin.SignedBlock, error) {
	avail := uint64(0)
	head, err := vs.Blockchain.Head()
	if err != nil {
		return nil, err
	}

	headSeq := head.Seq()
	if headSeq > seq {
		avail = headSeq - seq
	}
	if avail < ct {
		ct = avail
	}
	if ct == 0 {
		return nil, nil
	}

	blocks := make([]coin.SignedBlock, 0, ct)
	for j := uint64(0); j < ct; j++ {
		i := seq + 1 + j
		b, err := vs.Blockchain.GetBlockBySeq(i)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, *b)
	}
	return blocks, nil
}

// HeadBkSeq returns the highest BkSeq we know, returns -1 if the chain is empty
func (vs *Visor) HeadBkSeq() uint64 {
	return vs.Blockchain.HeadSeq()
}

// GetBlockchainMetadata returns descriptive Blockchain information
func (vs *Visor) GetBlockchainMetadata() BlockchainMetadata {
	return NewBlockchainMetadata(vs)
}

// GetBlock returns a copy of the block at seq. Returns error if seq out of range
// Move to blockdb
func (vs *Visor) GetBlock(seq uint64) (*coin.SignedBlock, error) {
	var b coin.SignedBlock
	if seq > vs.Blockchain.HeadSeq() {
		return &b, errors.New("Block seq out of range")
	}

	return vs.Blockchain.GetBlockBySeq(seq)
}

// GetBlocks returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
// move to blockdb
func (vs *Visor) GetBlocks(start, end uint64) []coin.SignedBlock {
	return vs.Blockchain.GetBlocks(start, end)
}

// InjectTransaction records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain.
// The bool return value is whether or not the transaction was already in the pool.
// If the transaction violates hard constraints, it is rejected, and error will not be nil.
// If the transaction only violates soft constraints, it is still injected, and the soft constraint violation is returned.
func (vs *Visor) InjectTransaction(txn coin.Transaction) (bool, *ErrTxnViolatesSoftConstraint, error) {
	return vs.Unconfirmed.InjectTransaction(vs.Blockchain, txn, vs.Config.MaxBlockSize)
}

// InjectTransactionStrict records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain.
// The bool return value is whether or not the transaction was already in the pool.
// If the transaction violates hard or soft constraints, it is rejected, and error will not be nil.
func (vs *Visor) InjectTransactionStrict(txn coin.Transaction) (bool, error) {
	if err := vs.Blockchain.VerifySingleTxnAllConstraints(txn, vs.Config.MaxBlockSize); err != nil {
		return false, err
	}

	known, _, err := vs.Unconfirmed.InjectTransaction(vs.Blockchain, txn, vs.Config.MaxBlockSize)
	return known, err
}

// GetAddressTxns returns the Transactions whose unspents give coins to a cipher.Address.
// This includes unconfirmed txns' predicted unspents.
func (vs *Visor) GetAddressTxns(a cipher.Address) ([]Transaction, error) {
	var txns []Transaction

	mxSeq := vs.HeadBkSeq()
	txs, err := vs.history.GetAddrTxns(a)
	if err != nil {
		return []Transaction{}, err
	}

	for _, tx := range txs {
		h := mxSeq - tx.BlockSeq + 1

		bk, err := vs.GetBlockBySeq(tx.BlockSeq)
		if err != nil {
			return []Transaction{}, err
		}

		if bk == nil {
			return []Transaction{}, fmt.Errorf("No block exsit in depth:%d", tx.BlockSeq)
		}

		txns = append(txns, Transaction{
			Txn:    tx.Tx,
			Status: NewConfirmedTransactionStatus(h, tx.BlockSeq),
			Time:   bk.Time(),
		})
	}

	// Look in the unconfirmed pool
	uxs := vs.Unconfirmed.GetUnspentsOfAddr(a)
	for _, ux := range uxs {
		tx, ok := vs.Unconfirmed.Get(ux.Body.SrcTransaction)
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

	return txns, nil
}

// GetTransaction returns a Transaction by hash.
func (vs *Visor) GetTransaction(txHash cipher.SHA256) (*Transaction, error) {
	// Look in the unconfirmed pool
	tx, ok := vs.Unconfirmed.Get(txHash)
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

	headSeq := vs.HeadBkSeq()

	confirms := headSeq - txn.BlockSeq + 1
	b, err := vs.GetBlockBySeq(txn.BlockSeq)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, fmt.Errorf("found no block in seq %v", txn.BlockSeq)
	}

	return &Transaction{
		Txn:    txn.Tx,
		Status: NewConfirmedTransactionStatus(confirms, txn.BlockSeq),
		Time:   b.Time(),
	}, nil
}

// TxFilter transaction filter type
type TxFilter interface {
	// Returns whether the transaction is matched
	Match(*Transaction) bool
}

// baseFilter is a helper struct for generating TxFilter.
type baseFilter struct {
	f func(tx *Transaction) bool
}

func (f baseFilter) Match(tx *Transaction) bool {
	return f.f(tx)
}

// AddrsFilter collects all addresses related transactions.
func AddrsFilter(addrs []cipher.Address) TxFilter {
	return addrsFilter{Addrs: addrs}
}

// addrsFilter
type addrsFilter struct {
	Addrs []cipher.Address
}

// Match implements the TxFilter interface, this actually won't be used, only the 'Addrs' member is used.
func (af addrsFilter) Match(tx *Transaction) bool { return true }

// ConfirmedTxFilter collects the transaction whose 'Confirmed' status matchs the parameter passed in.
func ConfirmedTxFilter(isConfirmed bool) TxFilter {
	return baseFilter{func(tx *Transaction) bool {
		return tx.Status.Confirmed == isConfirmed
	}}
}

// GetTransactions returns transactions that can pass the filters.
// If any 'AddrsFilter' exist, call vs.getTransactionsOfAddrs, cause
// there's an address index of transactions in db which, having address as key and transaction hashes as value.
// If no filters is provided, returns all transactions.
func (vs *Visor) GetTransactions(flts ...TxFilter) ([]Transaction, error) {
	var addrFlts []addrsFilter
	var otherFlts []TxFilter
	// Splits the filters into AddrsFilter and other filters
	for _, f := range flts {
		switch v := f.(type) {
		case addrsFilter:
			addrFlts = append(addrFlts, v)
		default:
			otherFlts = append(otherFlts, f)
		}
	}

	// Accumulates all addresses in address filters
	addrs := accumulateAddressInFilter(addrFlts)

	// Traverses all transactions to do collection if there's no address filter.
	if len(addrs) == 0 {
		return vs.traverseTxns(otherFlts...)
	}

	// Gets addresses related transactions
	txns, err := getTransactionsOfAddrs(vs, addrs)
	if err != nil {
		return nil, err
	}

	// Checks other filters
	var retTxns []Transaction
	f := func(tx *Transaction, flts ...TxFilter) bool {
		for _, flt := range otherFlts {
			if !flt.Match(tx) {
				return false
			}
		}

		return true
	}

	for _, tx := range txns {
		if f(&tx, otherFlts...) {
			retTxns = append(retTxns, tx)
		}
	}

	return retTxns, nil
}

func accumulateAddressInFilter(afs []addrsFilter) []cipher.Address {
	// Accumulate all addresses in address filters
	addrMap := make(map[cipher.Address]struct{}, 0)
	var addrs []cipher.Address
	for _, af := range afs {
		for _, a := range af.Addrs {
			if _, exist := addrMap[a]; exist {
				continue
			}
			addrMap[a] = struct{}{}
			addrs = append(addrs, a)
		}
	}
	return addrs
}

func getTransactionsOfAddrs(vs *Visor, addrs []cipher.Address) ([]Transaction, error) {
	addrTxns, err := vs.getTransactionsOfAddrs(addrs)
	if err != nil {
		return nil, err
	}

	// Converts address transactions map into []Transaction,
	// and remove duplicate txns
	txnMap := make(map[cipher.SHA256]struct{}, 0)
	var txns []Transaction
	for _, txs := range addrTxns {
		for _, tx := range txs {
			if _, exist := txnMap[tx.Txn.Hash()]; exist {
				continue
			}
			txnMap[tx.Txn.Hash()] = struct{}{}
			txns = append(txns, tx)
		}
	}
	return txns, nil
}

// getTransactionsOfAddrs returns all addresses related transactions.
// Including both confirmed and unconfirmed transactions.
func (vs *Visor) getTransactionsOfAddrs(addrs []cipher.Address) (map[cipher.Address][]Transaction, error) {
	// Initialize the address transactions map
	addrTxs := make(map[cipher.Address][]Transaction)

	// Get the head block seq, for caculating the tx status
	headBkSeq := vs.HeadBkSeq()
	for _, a := range addrs {
		var txns []Transaction
		txs, err := vs.history.GetAddrTxns(a)
		if err != nil {
			return nil, err
		}

		for _, tx := range txs {
			h := headBkSeq - tx.BlockSeq + 1

			bk, err := vs.GetBlockBySeq(tx.BlockSeq)
			if err != nil {
				return nil, err
			}

			if bk == nil {
				return nil, fmt.Errorf("block of seq: %d doesn't exist", tx.BlockSeq)
			}

			txns = append(txns, Transaction{
				Txn:    tx.Tx,
				Status: NewConfirmedTransactionStatus(h, tx.BlockSeq),
				Time:   bk.Time(),
			})
		}

		// Look in the unconfirmed pool
		uxs := vs.Unconfirmed.GetUnspentsOfAddr(a)
		for _, ux := range uxs {
			tx, ok := vs.Unconfirmed.Get(ux.Body.SrcTransaction)
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

		addrTxs[a] = txns
	}

	return addrTxs, nil
}

// traverseTxns traverses transactions in historydb and unconfirmed tx pool in db,
// returns transactions that can pass the filters.
func (vs *Visor) traverseTxns(flts ...TxFilter) ([]Transaction, error) {
	headBkSeq := vs.HeadBkSeq()
	var txns []Transaction
	err := vs.history.ForEach(func(tx *historydb.Transaction) error {
		h := headBkSeq - tx.BlockSeq + 1
		bk, err := vs.GetBlockBySeq(tx.BlockSeq)
		if err != nil {
			return fmt.Errorf("get block of seq: %v failed: %v", tx.BlockSeq, err)
		}

		if bk == nil {
			return fmt.Errorf("block of seq: %d doesn't exist", tx.BlockSeq)
		}

		txn := Transaction{
			Txn:    tx.Tx,
			Status: NewConfirmedTransactionStatus(h, tx.BlockSeq),
			Time:   bk.Time(),
		}

		// Checks filters
		for _, f := range flts {
			if !f.Match(&txn) {
				return nil
			}
		}

		txns = append(txns, txn)
		return nil
	})

	if err != nil {
		return nil, err
	}

	txns = sortTxns(txns)

	// Gets all unconfirmed transactions
	unconfirmedTxns := vs.Unconfirmed.GetTxns(func(tx UnconfirmedTxn) bool { return true })
	for _, ux := range unconfirmedTxns {
		tx := Transaction{
			Txn:    ux.Txn,
			Status: NewUnconfirmedTransactionStatus(),
			Time:   uint64(nanoToTime(ux.Received).Unix()),
		}

		// Checks filters
		for _, f := range flts {
			if !f.Match(&tx) {
				continue
			}
			txns = append(txns, tx)
		}
	}
	return txns, nil
}

func txMatchFilters(tx *Transaction, flts ...TxFilter) bool {
	for _, f := range flts {
		if !f.Match(tx) {
			return false
		}
	}
	return true
}

// Sort transactions by block seq, if equal then compare hash
func sortTxns(txns []Transaction) []Transaction {
	sort.Slice(txns, func(i, j int) bool {
		if txns[i].Status.BlockSeq < txns[j].Status.BlockSeq {
			return true
		}

		if txns[i].Status.BlockSeq > txns[j].Status.BlockSeq {
			return false
		}

		// If transactions in the same block, compare the hash string
		return txns[i].Txn.Hash().Hex() < txns[j].Txn.Hash().Hex()
	})
	return txns
}

// AddressBalance computes the total balance for cipher.Addresses and their coin.UxOuts
func (vs *Visor) AddressBalance(auxs coin.AddressUxOuts) (uint64, uint64, error) {
	prevTime := vs.Blockchain.Time()
	var coins uint64
	var hours uint64
	for _, uxs := range auxs {
		for _, ux := range uxs {
			uxHours, err := ux.CoinHours(prevTime)
			if err != nil {
				return 0, 0, err
			}

			coins, err = coin.AddUint64(coins, ux.Body.Coins)
			if err != nil {
				return 0, 0, err
			}

			hours, err = coin.AddUint64(hours, uxHours)
			if err != nil {
				return 0, 0, err
			}
		}
	}
	return coins, hours, nil
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
func (vs *Visor) GetBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error) {
	return vs.Blockchain.GetBlockByHash(hash)
}

// GetBlockBySeq get block of speicific seq, return nil on not found.
func (vs *Visor) GetBlockBySeq(seq uint64) (*coin.SignedBlock, error) {
	return vs.Blockchain.GetBlockBySeq(seq)
}

// GetLastBlocks returns last N blocks
func (vs *Visor) GetLastBlocks(num uint64) []coin.SignedBlock {
	return vs.Blockchain.GetLastBlocks(num)
}

// GetHeadBlock gets head block.
func (vs Visor) GetHeadBlock() (*coin.SignedBlock, error) {
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

// ScanAheadWalletAddresses scans ahead N addresses in a wallet, looking for a non-empty balance
func (vs Visor) ScanAheadWalletAddresses(wltName string, scanN uint64) (wallet.Wallet, error) {
	return vs.wallets.ScanAheadWalletAddresses(wltName, scanN, vs)
}

// GetBalanceOfAddrs returns balance pairs of given addreses
func (vs Visor) GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error) {
	var bps []wallet.BalancePair
	auxs := vs.Blockchain.Unspent().GetUnspentsOfAddrs(addrs)
	spendUxs, err := vs.Unconfirmed.SpendsOfAddresses(addrs, vs.Blockchain.Unspent())
	if err != nil {
		return nil, fmt.Errorf("get unconfirmed spending failed when checking addresses balance: %v", err)
	}

	head, err := vs.Blockchain.Head()
	if err != nil {
		return nil, err
	}

	recvUxs, err := vs.Unconfirmed.RecvOfAddresses(head.Head, addrs)
	if err != nil {
		return nil, fmt.Errorf("get unconfirmed receiving failed when checking addresses balance: %v", err)
	}

	headTime := head.Time()
	for _, addr := range addrs {
		uxs, ok := auxs[addr]
		if !ok {
			bps = append(bps, wallet.BalancePair{})
			continue
		}

		outUxs := spendUxs[addr]
		inUxs := recvUxs[addr]
		predictedUxs := uxs.Sub(outUxs).Add(inUxs)

		coins, err := uxs.Coins()
		if err != nil {
			return nil, fmt.Errorf("uxs.Coins failed: %v", err)
		}

		coinHours, err := uxs.CoinHours(headTime)
		if err != nil {
			switch err {
			case coin.ErrAddEarnedCoinHoursAdditionOverflow:
				coinHours = 0
				err = nil
			default:
				return nil, fmt.Errorf("uxs.CoinHours failed: %v", err)
			}
		}

		pcoins, err := predictedUxs.Coins()
		if err != nil {
			return nil, fmt.Errorf("predictedUxs.Coins failed: %v", err)
		}

		pcoinHours, err := predictedUxs.CoinHours(headTime)
		if err != nil {
			switch err {
			case coin.ErrAddEarnedCoinHoursAdditionOverflow:
				coinHours = 0
				err = nil
			default:
				return nil, fmt.Errorf("predictedUxs.CoinHours failed: %v", err)
			}
		}

		bp := wallet.BalancePair{
			Confirmed: wallet.Balance{
				Coins: coins,
				Hours: coinHours,
			},
			Predicted: wallet.Balance{
				Coins: pcoins,
				Hours: pcoinHours,
			},
		}

		bps = append(bps, bp)
	}

	return bps, nil
}
