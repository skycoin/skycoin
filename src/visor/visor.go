package visor

import (
	"errors"
	"fmt"
	"sort"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/util/logging"
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
	sanityCheck()
	// Compute maxDropletDivisor from precision
	maxDropletDivisor = calculateDivisor(MaxDropletPrecision)
}

func sanityCheck() {
	if InitialUnlockedCount > DistributionAddressesTotal {
		logger.Panic("unlocked addresses > total distribution addresses")
	}

	if uint64(len(distributionAddresses)) != DistributionAddressesTotal {
		logger.Panic("available distribution addresses > total allowed distribution addresses")
	}

	if DistributionAddressInitialBalance*DistributionAddressesTotal > MaxCoinSupply {
		logger.Panic("total balance in distribution addresses > max coin supply")
	}
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
	Branch  string `json:"branch"`  // git branch name
}

// Config configuration parameters for the Visor
type Config struct {
	// Is this the master blockchain
	IsMaster bool

	//Public key of blockchain authority
	BlockchainPubkey cipher.PubKey

	//Secret key of blockchain authority (if master)
	BlockchainSeckey cipher.SecKey

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
	// enable arbitrating mode
	Arbitrating bool
	// wallet directory
	WalletDirectory string
	// build info, including version, build time etc.
	BuildInfo BuildInfo
	// enables wallet API
	EnableWalletAPI bool
	// enables seed API
	EnableSeedAPI bool
	// wallet crypto type
	WalletCryptoType wallet.CryptoType
}

// NewVisorConfig put cap on block size, not on transactions/block
//Skycoin transactions are smaller than Bitcoin transactions so skycoin has
//a higher transactions per second for the same block size
func NewVisorConfig() Config {
	c := Config{
		IsMaster: false,

		BlockchainPubkey: cipher.PubKey{},
		BlockchainSeckey: cipher.SecKey{},

		MaxBlockSize: DefaultMaxBlockSize,

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

//go:generate go install
//go:generate goautomock -template=testify Historyer

// Historyer is the interface that provides methods for accessing history data that are parsed from blockchain.
type Historyer interface {
	GetUxOuts(tx *dbutil.Tx, uxids []cipher.SHA256) ([]*historydb.UxOut, error)
	ParseBlock(tx *dbutil.Tx, b coin.Block) error
	GetTransaction(tx *dbutil.Tx, hash cipher.SHA256) (*historydb.Transaction, error)
	GetAddrUxOuts(tx *dbutil.Tx, address cipher.Address) ([]*historydb.UxOut, error)
	GetAddressTxns(tx *dbutil.Tx, address cipher.Address) ([]historydb.Transaction, error)
	NeedsReset(tx *dbutil.Tx) (bool, error)
	Erase(tx *dbutil.Tx) error
	ParsedHeight(tx *dbutil.Tx) (uint64, bool, error)
	ForEachTxn(tx *dbutil.Tx, f func(cipher.SHA256, *historydb.Transaction) error) error
}

// Blockchainer is the interface that provides methods for accessing the blockchain data
type Blockchainer interface {
	GetGenesisBlock(tx *dbutil.Tx) (*coin.SignedBlock, error)
	GetBlocks(tx *dbutil.Tx, start, end uint64) ([]coin.SignedBlock, error)
	GetLastBlocks(tx *dbutil.Tx, n uint64) ([]coin.SignedBlock, error)
	GetSignedBlockByHash(tx *dbutil.Tx, hash cipher.SHA256) (*coin.SignedBlock, error)
	GetSignedBlockBySeq(tx *dbutil.Tx, seq uint64) (*coin.SignedBlock, error)
	Unspent() blockdb.UnspentPooler
	Len(tx *dbutil.Tx) (uint64, error)
	Head(tx *dbutil.Tx) (*coin.SignedBlock, error)
	HeadSeq(tx *dbutil.Tx) (uint64, bool, error)
	Time(tx *dbutil.Tx) (uint64, error)
	NewBlock(tx *dbutil.Tx, txns coin.Transactions, currentTime uint64) (*coin.Block, error)
	ExecuteBlock(tx *dbutil.Tx, sb *coin.SignedBlock) error
	VerifyBlockTxnConstraints(tx *dbutil.Tx, txn coin.Transaction) error
	VerifySingleTxnHardConstraints(tx *dbutil.Tx, txn coin.Transaction) error
	VerifySingleTxnSoftHardConstraints(tx *dbutil.Tx, txn coin.Transaction, maxSize int) error
	TransactionFee(tx *dbutil.Tx, hours uint64) coin.FeeCalculator
}

// UnconfirmedTxnPooler is the interface that provides methods for
// accessing the unconfirmed transaction pool
type UnconfirmedTxnPooler interface {
	SetTxnsAnnounced(tx *dbutil.Tx, hashes map[cipher.SHA256]int64) error
	InjectTransaction(tx *dbutil.Tx, bc Blockchainer, t coin.Transaction, maxSize int) (bool, *ErrTxnViolatesSoftConstraint, error)
	RawTxns(tx *dbutil.Tx) (coin.Transactions, error)
	RemoveTransactions(tx *dbutil.Tx, txns []cipher.SHA256) error
	Refresh(tx *dbutil.Tx, bc Blockchainer, maxBlockSize int) ([]cipher.SHA256, error)
	RemoveInvalid(tx *dbutil.Tx, bc Blockchainer) ([]cipher.SHA256, error)
	GetUnknown(tx *dbutil.Tx, txns []cipher.SHA256) ([]cipher.SHA256, error)
	GetKnown(tx *dbutil.Tx, txns []cipher.SHA256) (coin.Transactions, error)
	RecvOfAddresses(tx *dbutil.Tx, bh coin.BlockHeader, addrs []cipher.Address) (coin.AddressUxOuts, error)
	GetIncomingOutputs(tx *dbutil.Tx, bh coin.BlockHeader) (coin.UxArray, error)
	Get(tx *dbutil.Tx, hash cipher.SHA256) (*UnconfirmedTxn, error)
	GetTxns(tx *dbutil.Tx, filter func(tx UnconfirmedTxn) bool) ([]UnconfirmedTxn, error)
	GetTxHashes(tx *dbutil.Tx, filter func(tx UnconfirmedTxn) bool) ([]cipher.SHA256, error)
	ForEach(tx *dbutil.Tx, f func(cipher.SHA256, UnconfirmedTxn) error) error
	GetUnspentsOfAddr(tx *dbutil.Tx, addr cipher.Address) (coin.UxArray, error)
	Len(tx *dbutil.Tx) (uint64, error)
}

// Visor manages the Blockchain as both a Master and a Normal
type Visor struct {
	Config Config
	DB     *dbutil.DB
	// Unconfirmed transactions, held for relay until we get block confirmation
	Unconfirmed UnconfirmedTxnPooler
	Blockchain  Blockchainer
	Wallets     *wallet.Service
	StartedAt   time.Time

	history Historyer
}

// NewVisor creates a Visor for managing the blockchain database
func NewVisor(c Config, db *dbutil.DB) (*Visor, error) {
	logger.Info("Creating new visor")
	if c.IsMaster {
		logger.Info("Visor is master")
	}

	if err := c.Verify(); err != nil {
		return nil, err
	}

	// Loads wallet
	wltServConfig := wallet.Config{
		WalletDir:       c.WalletDirectory,
		CryptoType:      c.WalletCryptoType,
		EnableWalletAPI: c.EnableWalletAPI,
		EnableSeedAPI:   c.EnableSeedAPI,
	}

	wltServ, err := wallet.NewService(wltServConfig)
	if err != nil {
		return nil, err
	}

	if !db.IsReadOnly() {
		if err := CreateBuckets(db); err != nil {
			logger.WithError(err).Error("CreateBuckets failed")
			return nil, err
		}
	}

	bc, err := NewBlockchain(db, BlockchainConfig{
		Pubkey:      c.BlockchainPubkey,
		Arbitrating: c.Arbitrating,
	})
	if err != nil {
		return nil, err
	}

	history := historydb.New()

	if !db.IsReadOnly() {
		if err := db.Update("build unspent indexes and init history", func(tx *dbutil.Tx) error {
			headSeq, _, err := bc.HeadSeq(tx)
			if err != nil {
				return err
			}

			if err := bc.Unspent().MaybeBuildIndexes(tx, headSeq); err != nil {
				return err
			}

			return initHistory(tx, bc, history)
		}); err != nil {
			return nil, err
		}
	}

	utp, err := NewUnconfirmedTxnPool(db)
	if err != nil {
		return nil, err
	}

	v := &Visor{
		Config:      c,
		DB:          db,
		Blockchain:  bc,
		Unconfirmed: utp,
		history:     history,
		Wallets:     wltServ,
		StartedAt:   time.Now(),
	}

	return v, nil
}

// Init initializes starts the visor
func (vs *Visor) Init() error {
	logger.Info("Visor init")

	if vs.DB.IsReadOnly() {
		return nil
	}

	return vs.DB.Update("visor init", func(tx *dbutil.Tx) error {
		if err := vs.maybeCreateGenesisBlock(tx); err != nil {
			return err
		}

		removed, err := vs.Unconfirmed.RemoveInvalid(tx, vs.Blockchain)
		if err != nil {
			return err
		}
		logger.Infof("Removed %d invalid txns from pool", len(removed))

		return nil
	})
}

func initHistory(tx *dbutil.Tx, bc *Blockchain, history *historydb.HistoryDB) error {
	logger.Info("Visor initHistory")

	shouldReset, err := history.NeedsReset(tx)
	if err != nil {
		return err
	}

	if !shouldReset {
		return nil
	}

	logger.Info("Resetting historyDB")

	if err := history.Erase(tx); err != nil {
		return err
	}

	// Reparse the history up to the blockchain head
	headSeq, _, err := bc.HeadSeq(tx)
	if err != nil {
		return err
	}

	if err := parseHistoryTo(tx, history, bc, headSeq); err != nil {
		logger.WithError(err).Error("parseHistoryTo failed")
		return err
	}

	return nil
}

func parseHistoryTo(tx *dbutil.Tx, history *historydb.HistoryDB, bc *Blockchain, height uint64) error {
	logger.Info("Visor parseHistoryTo")

	parsedHeight, _, err := history.ParsedHeight(tx)
	if err != nil {
		return err
	}

	for i := uint64(0); i < height-parsedHeight; i++ {
		b, err := bc.GetSignedBlockBySeq(tx, parsedHeight+i+1)
		if err != nil {
			return err
		}

		if b == nil {
			return fmt.Errorf("no block exists in depth: %d", parsedHeight+i+1)
		}

		if err := history.ParseBlock(tx, b.Block); err != nil {
			return err
		}
	}

	return nil
}

// maybeCreateGenesisBlock creates a genesis block if necessary
func (vs *Visor) maybeCreateGenesisBlock(tx *dbutil.Tx) error {
	logger.Info("Visor maybeCreateGenesisBlock")
	gb, err := vs.Blockchain.GetGenesisBlock(tx)
	if err != nil {
		return err
	}
	if gb != nil {
		return nil
	}

	logger.Info("Create genesis block")
	vs.GenesisPreconditions()
	b, err := coin.NewGenesisBlock(vs.Config.GenesisAddress, vs.Config.GenesisCoinVolume, vs.Config.GenesisTimestamp)
	if err != nil {
		return err
	}

	var sb coin.SignedBlock
	// record the signature of genesis block
	if vs.Config.IsMaster {
		sb = vs.signBlock(*b)
		logger.Infof("Genesis block signature=%s", sb.Sig.Hex())
	} else {
		sb = coin.SignedBlock{
			Block: *b,
			Sig:   vs.Config.GenesisSignature,
		}
	}

	return vs.executeSignedBlock(tx, sb)
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
	var hashes []cipher.SHA256
	if err := vs.DB.Update("RefreshUnconfirmed", func(tx *dbutil.Tx) error {
		var err error
		hashes, err = vs.Unconfirmed.Refresh(tx, vs.Blockchain, vs.Config.MaxBlockSize)
		return err
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// RemoveInvalidUnconfirmed removes transactions that become permanently invalid
// (by violating hard constraints) from the pool.
// Returns the transaction hashes that were removed.
func (vs *Visor) RemoveInvalidUnconfirmed() ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256
	if err := vs.DB.Update("RemoveInvalidUnconfirmed", func(tx *dbutil.Tx) error {
		var err error
		hashes, err = vs.Unconfirmed.RemoveInvalid(tx, vs.Blockchain)
		return err
	}); err != nil {
		return nil, err
	}

	return hashes, nil
}

// CreateBlock creates a SignedBlock from pending transactions
func (vs *Visor) createBlock(tx *dbutil.Tx, when uint64) (coin.SignedBlock, error) {
	if !vs.Config.IsMaster {
		logger.Panic("Only master chain can create blocks")
	}

	// Gather all unconfirmed transactions
	txns, err := vs.Unconfirmed.RawTxns(tx)
	if err != nil {
		return coin.SignedBlock{}, err
	}

	if len(txns) == 0 {
		return coin.SignedBlock{}, errors.New("No transactions")
	}

	logger.Infof("Unconfirmed pool has %d transactions pending", len(txns))

	// Filter transactions that violate all constraints
	var filteredTxns coin.Transactions
	for _, txn := range txns {
		if err := vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, txn, vs.Config.MaxBlockSize); err != nil {
			switch err.(type) {
			case ErrTxnViolatesHardConstraint, ErrTxnViolatesSoftConstraint:
				logger.Warningf("Transaction %s violates constraints: %v", txn.TxIDHex(), err)
			default:
				return coin.SignedBlock{}, err
			}
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
		return coin.SignedBlock{}, errors.New("No transactions after filtering for constraint violations")
	}

	head, err := vs.Blockchain.Head(tx)
	if err != nil {
		return coin.SignedBlock{}, err
	}

	// Sort them by highest fee per kilobyte
	txns = coin.SortTransactions(txns, vs.Blockchain.TransactionFee(tx, head.Time()))

	// Apply block size transaction limit
	txns = txns.TruncateBytesTo(vs.Config.MaxBlockSize)

	if len(txns) == 0 {
		logger.Panic("TruncateBytesTo removed all transactions")
	}

	logger.Infof("Creating new block with %d transactions, head time %d", len(txns), when)

	b, err := vs.Blockchain.NewBlock(tx, txns, when)
	if err != nil {
		logger.Warningf("Blockchain.NewBlock failed: %v", err)
		return coin.SignedBlock{}, err
	}

	return vs.signBlock(*b), nil
}

// CreateAndExecuteBlock creates a SignedBlock from pending transactions and executes it
func (vs *Visor) CreateAndExecuteBlock() (coin.SignedBlock, error) {
	var sb coin.SignedBlock

	err := vs.DB.Update("CreateAndExecuteBlock", func(tx *dbutil.Tx) error {
		var err error
		sb, err = vs.createBlock(tx, uint64(utc.UnixNow()))
		if err != nil {
			return err
		}

		return vs.executeSignedBlock(tx, sb)
	})

	return sb, err
}

// ExecuteSignedBlock adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (vs *Visor) ExecuteSignedBlock(b coin.SignedBlock) error {
	return vs.DB.Update("ExecuteSignedBlock", func(tx *dbutil.Tx) error {
		return vs.executeSignedBlock(tx, b)
	})
}

// executeSignedBlock adds a block to the blockchain, or returns error.
// Blocks must be executed in sequence, and be signed by the master server
func (vs *Visor) executeSignedBlock(tx *dbutil.Tx, b coin.SignedBlock) error {
	if err := b.VerifySignature(vs.Config.BlockchainPubkey); err != nil {
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

	// Update the HistoryDB
	return vs.history.ParseBlock(tx, b.Block)
}

// signBlock signs a block for master.  Will panic if anything is invalid
func (vs *Visor) signBlock(b coin.Block) coin.SignedBlock {
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

// GetAllUnspentOutputs returns all unspent outputs
func (vs *Visor) GetAllUnspentOutputs() (coin.UxArray, error) {
	var ux []coin.UxOut
	if err := vs.DB.View("GetAllUnspentOutputs", func(tx *dbutil.Tx) error {
		var err error
		ux, err = vs.Blockchain.Unspent().GetAll(tx)
		return err
	}); err != nil {
		return nil, err
	}

	return ux, nil
}

// GetUnspentOutputs returns unspent outputs from the pool, queried by hashes.
// If any do not exist, ErrUnspentNotExist is returned
func (vs *Visor) GetUnspentOutputs(hashes []cipher.SHA256) (coin.UxArray, error) {
	var outputs coin.UxArray
	if err := vs.DB.View("GetUnspentOutputs", func(tx *dbutil.Tx) error {
		var err error
		outputs, err = vs.Blockchain.Unspent().GetArray(tx, hashes)
		return err
	}); err != nil {
		return nil, err
	}

	return outputs, nil
}

// UnconfirmedSpendingOutputs returns all spending outputs in unconfirmed tx pool
func (vs *Visor) UnconfirmedSpendingOutputs() (coin.UxArray, error) {
	var uxa coin.UxArray

	if err := vs.DB.View("UnconfirmedSpendingOutputs", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("UnconfirmedIncomingOutputs", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetSignedBlocksSince", func(tx *dbutil.Tx) error {
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

// HeadBkSeq returns the highest BkSeq we know, returns false in the 2nd return value
// if the blockchain is empty
func (vs *Visor) HeadBkSeq() (uint64, bool, error) {
	var headSeq uint64
	var ok bool

	if err := vs.DB.View("HeadBkSeq", func(tx *dbutil.Tx) error {
		var err error
		headSeq, ok, err = vs.Blockchain.HeadSeq(tx)
		return err
	}); err != nil {
		return 0, false, err
	}

	return headSeq, ok, nil
}

// GetBlockchainMetadata returns descriptive Blockchain information
func (vs *Visor) GetBlockchainMetadata() (*BlockchainMetadata, error) {
	var head *coin.SignedBlock
	var unconfirmedLen, unspentsLen uint64

	if err := vs.DB.View("GetBlockchainMetadata", func(tx *dbutil.Tx) error {
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
func (vs *Visor) GetBlock(seq uint64) (*coin.SignedBlock, error) {
	var b *coin.SignedBlock

	if err := vs.DB.View("GetBlock", func(tx *dbutil.Tx) error {
		headSeq, ok, err := vs.Blockchain.HeadSeq(tx)
		if err != nil {
			return err
		}

		if !ok || seq > headSeq {
			return errors.New("Block seq out of range")
		}

		b, err = vs.Blockchain.GetSignedBlockBySeq(tx, seq)
		return err
	}); err != nil {
		return nil, err
	}

	return b, nil
}

// GetBlocks returns multiple blocks between start and end (not including end). Returns
// empty slice if unable to fulfill request, it does not return nil.
func (vs *Visor) GetBlocks(start, end uint64) ([]coin.SignedBlock, error) {
	var blocks []coin.SignedBlock

	if err := vs.DB.View("GetBlocks", func(tx *dbutil.Tx) error {
		var err error
		blocks, err = vs.Blockchain.GetBlocks(tx, start, end)
		return err
	}); err != nil {
		return nil, err
	}

	return blocks, nil
}

// InjectTransaction records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain.
// The bool return value is whether or not the transaction was already in the pool.
// If the transaction violates hard constraints, it is rejected, and error will not be nil.
// If the transaction only violates soft constraints, it is still injected, and the soft constraint violation is returned.
func (vs *Visor) InjectTransaction(txn coin.Transaction) (bool, *ErrTxnViolatesSoftConstraint, error) {
	var known bool
	var softErr *ErrTxnViolatesSoftConstraint

	if err := vs.DB.Update("InjectTransaction", func(tx *dbutil.Tx) error {
		var err error
		known, softErr, err = vs.Unconfirmed.InjectTransaction(tx, vs.Blockchain, txn, vs.Config.MaxBlockSize)
		return err
	}); err != nil {
		return false, nil, err
	}

	return known, softErr, nil
}

// InjectTransactionStrict records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain.
// The bool return value is whether or not the transaction was already in the pool.
// If the transaction violates hard or soft constraints, it is rejected, and error will not be nil.
func (vs *Visor) InjectTransactionStrict(txn coin.Transaction) (bool, error) {
	if err := VerifySingleTxnUserConstraints(txn); err != nil {
		return false, err
	}

	var known bool

	if err := vs.DB.Update("InjectTransactionStrict", func(tx *dbutil.Tx) error {
		err := vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, txn, vs.Config.MaxBlockSize)
		if err != nil {
			return err
		}

		known, _, err = vs.Unconfirmed.InjectTransaction(tx, vs.Blockchain, txn, vs.Config.MaxBlockSize)
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

	if err := vs.DB.View("GetAddressTxns", func(tx *dbutil.Tx) error {
		txs, err := vs.history.GetAddressTxns(tx, a)
		if err != nil {
			return err
		}

		mxSeq, ok, err := vs.Blockchain.HeadSeq(tx)
		if err != nil {
			return err
		} else if !ok {
			if len(txns) > 0 {
				return fmt.Errorf("Found %d txns for addresses but block head seq is missing", len(txns))
			}
			return nil
		}

		for _, txn := range txs {
			if mxSeq < txn.BlockSeq {
				return fmt.Errorf("Blockchain head seq %d is earlier than history txn seq %d", mxSeq, txn.BlockSeq)
			}
			h := mxSeq - txn.BlockSeq + 1

			bk, err := vs.Blockchain.GetSignedBlockBySeq(tx, txn.BlockSeq)
			if err != nil {
				return err
			}

			if bk == nil {
				return fmt.Errorf("No block exists in depth: %d", txn.BlockSeq)
			}

			txns = append(txns, Transaction{
				Txn:    txn.Tx,
				Status: NewConfirmedTransactionStatus(h, txn.BlockSeq),
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
				logger.Critical().Error("Unconfirmed unspent missing unconfirmed txn")
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

	if err := vs.DB.View("GetTransaction", func(tx *dbutil.Tx) error {
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

		headSeq, ok, err := vs.Blockchain.HeadSeq(tx)
		if err != nil {
			return err
		} else if !ok {
			return errors.New("Blockchain is empty but history has transactions")
		}

		b, err := vs.Blockchain.GetSignedBlockBySeq(tx, htxn.BlockSeq)
		if err != nil {
			return err
		}

		if b == nil {
			return fmt.Errorf("found no block in seq %v", htxn.BlockSeq)
		}

		if headSeq < htxn.BlockSeq {
			return fmt.Errorf("Blockchain head seq %d is earlier than history txn seq %d", headSeq, htxn.BlockSeq)
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
		var txns []Transaction
		if err := vs.DB.View("GetTransactions traverseTxns", func(tx *dbutil.Tx) error {
			var err error
			txns, err = vs.traverseTxns(tx, otherFlts...)
			return err
		}); err != nil {
			return nil, err
		}
		return txns, nil
	}

	// Gets addresses related transactions
	var addrTxns map[cipher.Address][]Transaction
	if err := vs.DB.View("GetTransactions getTransactionsOfAddrs", func(tx *dbutil.Tx) error {
		var err error
		addrTxns, err = vs.getTransactionsOfAddrs(tx, addrs)
		return err
	}); err != nil {
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

// getTransactionsOfAddrs returns all addresses related transactions.
// Including both confirmed and unconfirmed transactions.
func (vs *Visor) getTransactionsOfAddrs(tx *dbutil.Tx, addrs []cipher.Address) (map[cipher.Address][]Transaction, error) {
	// Initialize the address transactions map
	addrTxs := make(map[cipher.Address][]Transaction)

	// Get the head block seq, for calculating the tx status
	headBkSeq, ok, err := vs.Blockchain.HeadSeq(tx)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("No head block seq")
	}

	for _, a := range addrs {
		var txns []Transaction
		addrTxns, err := vs.history.GetAddressTxns(tx, a)
		if err != nil {
			return nil, err
		}

		for _, txn := range addrTxns {
			if headBkSeq < txn.BlockSeq {
				err := errors.New("Transaction block sequence is less than the head block sequence")
				logger.Critical().WithError(err).WithFields(logrus.Fields{
					"headBkSeq":  headBkSeq,
					"txBlockSeq": txn.BlockSeq,
				}).Error()
				return nil, err
			}
			h := headBkSeq - txn.BlockSeq + 1

			bk, err := vs.Blockchain.GetSignedBlockBySeq(tx, txn.BlockSeq)
			if err != nil {
				return nil, err
			}

			if bk == nil {
				return nil, fmt.Errorf("block of seq: %d doesn't exist", txn.BlockSeq)
			}

			txns = append(txns, Transaction{
				Txn:    txn.Tx,
				Status: NewConfirmedTransactionStatus(h, txn.BlockSeq),
				Time:   bk.Time(),
			})
		}

		// Look in the unconfirmed pool
		uxs, err := vs.Unconfirmed.GetUnspentsOfAddr(tx, a)
		if err != nil {
			return nil, err
		}

		for _, ux := range uxs {
			txn, err := vs.Unconfirmed.Get(tx, ux.Body.SrcTransaction)
			if err != nil {
				return nil, err
			}

			if txn == nil {
				logger.Critical().Error("Unconfirmed unspent missing unconfirmed txn")
				continue
			}

			txns = append(txns, Transaction{
				Txn:    txn.Txn,
				Status: NewUnconfirmedTransactionStatus(),
				Time:   uint64(nanoToTime(txn.Received).Unix()),
			})
		}

		addrTxs[a] = txns
	}

	return addrTxs, nil
}

// traverseTxns traverses transactions in historydb and unconfirmed tx pool in db,
// returns transactions that can pass the filters.
func (vs *Visor) traverseTxns(tx *dbutil.Tx, flts ...TxFilter) ([]Transaction, error) {
	// Get the head block seq, for calculating the tx status
	headBkSeq, ok, err := vs.Blockchain.HeadSeq(tx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("No head block seq")
	}

	var txns []Transaction

	if err := vs.history.ForEachTxn(tx, func(_ cipher.SHA256, hTxn *historydb.Transaction) error {
		if headBkSeq < hTxn.BlockSeq {
			err := errors.New("Transaction block sequence is less than the head block sequence")
			logger.Critical().WithError(err).WithFields(logrus.Fields{
				"headBkSeq":  headBkSeq,
				"txBlockSeq": hTxn.BlockSeq,
			}).Error()
			return err
		}

		h := headBkSeq - hTxn.BlockSeq + 1

		bk, err := vs.Blockchain.GetSignedBlockBySeq(tx, hTxn.BlockSeq)
		if err != nil {
			return fmt.Errorf("get block of seq: %v failed: %v", hTxn.BlockSeq, err)
		}

		if bk == nil {
			return fmt.Errorf("block of seq: %d doesn't exist", hTxn.BlockSeq)
		}

		txn := Transaction{
			Txn:    hTxn.Tx,
			Status: NewConfirmedTransactionStatus(h, hTxn.BlockSeq),
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
	}); err != nil {
		return nil, err
	}

	txns = sortTxns(txns)

	// Gets all unconfirmed transactions
	unconfirmedTxns, err := vs.Unconfirmed.GetTxns(tx, func(txn UnconfirmedTxn) bool {
		return true
	})
	if err != nil {
		return nil, err
	}

	for _, ux := range unconfirmedTxns {
		txn := Transaction{
			Txn:    ux.Txn,
			Status: NewUnconfirmedTransactionStatus(),
			Time:   uint64(nanoToTime(ux.Received).Unix()),
		}

		// Checks filters
		for _, f := range flts {
			if !f.Match(&txn) {
				continue
			}
			txns = append(txns, txn)
		}
	}
	return txns, nil
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
func (vs *Visor) AddressBalance(head *coin.SignedBlock, auxs coin.AddressUxOuts) (uint64, uint64, error) {
	prevTime := head.Time()
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
func (vs *Visor) GetUnconfirmedTxns(filter func(UnconfirmedTxn) bool) ([]UnconfirmedTxn, error) {
	var txns []UnconfirmedTxn

	if err := vs.DB.View("GetUnconfirmedTxns", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetAllUnconfirmedTxns", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetAllValidUnconfirmedTxHashes", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetSignedBlockByHash", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetSignedBlockBySeq", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetLastBlocks", func(tx *dbutil.Tx) error {
		var err error
		blocks, err = vs.Blockchain.GetLastBlocks(tx, num)
		return err
	}); err != nil {
		return nil, err
	}

	return blocks, nil
}

// GetHeadBlock gets head block.
func (vs Visor) GetHeadBlock() (*coin.SignedBlock, error) {
	var b *coin.SignedBlock

	if err := vs.DB.View("GetHeadBlock", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetHeadBlockTime", func(tx *dbutil.Tx) error {
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
	var outs []*historydb.UxOut

	if err := vs.DB.View("GetUxOutByID", func(tx *dbutil.Tx) error {
		var err error
		outs, err = vs.history.GetUxOuts(tx, []cipher.SHA256{id})
		return err
	}); err != nil {
		return nil, err
	}

	if len(outs) == 0 {
		return nil, nil
	}

	return outs[0], nil
}

// GetAddrUxOuts gets all the address affected UxOuts.
func (vs Visor) GetAddrUxOuts(address cipher.Address) ([]*historydb.UxOut, error) {
	var out []*historydb.UxOut

	if err := vs.DB.View("GetAddrUxOuts", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("RecvOfAddresses", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetIncomingOutputs", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetUnconfirmedTxn", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetUnconfirmedUnknown", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("GetUnconfirmedKnown", func(tx *dbutil.Tx) error {
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

	if err := vs.DB.View("UnconfirmedSpendsOfAddresses", func(tx *dbutil.Tx) error {
		var err error
		outs, err = vs.unconfirmedSpendsOfAddresses(tx, addrs)
		return err
	}); err != nil {
		return nil, err
	}

	return outs, nil
}

// unconfirmedSpendsOfAddresses returns all unconfirmed coin.UxOut spends of addresses
func (vs *Visor) unconfirmedSpendsOfAddresses(tx *dbutil.Tx, addrs []cipher.Address) (coin.AddressUxOuts, error) {
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
func (vs *Visor) SetTxnsAnnounced(hashes map[cipher.SHA256]int64) error {
	if len(hashes) == 0 {
		return nil
	}

	return vs.DB.Update("SetTxnsAnnounced", func(tx *dbutil.Tx) error {
		return vs.Unconfirmed.SetTxnsAnnounced(tx, hashes)
	})
}

// GetBalanceOfAddrs returns balance pairs of given addreses
func (vs Visor) GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error) {
	if len(addrs) == 0 {
		return nil, nil
	}

	auxs := make(coin.AddressUxOuts, len(addrs))
	recvUxs := make(coin.AddressUxOuts, len(addrs))
	var uxa coin.UxArray
	var head *coin.SignedBlock

	if err := vs.DB.View("GetBalanceOfAddrs", func(tx *dbutil.Tx) error {
		var err error
		head, err = vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		// Get all transactions from the unconfirmed pool
		txns, err := vs.Unconfirmed.RawTxns(tx)
		if err != nil {
			return err
		}

		// Create predicted unspent outputs from the unconfirmed transactions
		recvUxs, err = txnOutputsForAddrs(head.Head, addrs, txns)
		if err != nil {
			return err
		}

		var inputs []cipher.SHA256
		for _, txn := range txns {
			inputs = append(inputs, txn.In...)
		}

		// Get unspents for the inputs being spent
		uxa, err = vs.Blockchain.Unspent().GetArray(tx, inputs)
		if err != nil {
			return fmt.Errorf("GetArray failed when checking addresses balance: %v", err)
		}

		// Get unspents owned by the addresses
		auxs, err = vs.Blockchain.Unspent().GetUnspentsOfAddrs(tx, addrs)
		if err != nil {
			return fmt.Errorf("GetUnspentsOfAddrs failed when checking addresses balance: %v", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Build all unconfirmed transaction inputs that are associated with the addresses
	spendUxs := make(coin.AddressUxOuts, len(addrs))

	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		addrm[addr] = struct{}{}
	}

	for _, ux := range uxa {
		if _, ok := addrm[ux.Body.Address]; ok {
			spendUxs[ux.Body.Address] = append(spendUxs[ux.Body.Address], ux)
		}
	}

	var bps []wallet.BalancePair

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

// GetUnspentsOfAddrs returns unspent outputs of multiple addresses
func (vs *Visor) GetUnspentsOfAddrs(addrs []cipher.Address) (coin.AddressUxOuts, error) {
	var uxa coin.AddressUxOuts

	if err := vs.DB.View("GetUnspentsOfAddrs", func(tx *dbutil.Tx) error {
		var err error
		uxa, err = vs.Blockchain.Unspent().GetUnspentsOfAddrs(tx, addrs)
		return err
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}

// VerifyTxnVerbose verifies a transaction, it returns transaction's input uxouts, whether the
// transaction is confirmed, and error if any
func (vs *Visor) VerifyTxnVerbose(txn *coin.Transaction) ([]wallet.UxBalance, bool, error) {
	var uxa coin.UxArray
	var head *coin.SignedBlock
	var isTxnConfirmed bool
	err := vs.DB.View("VerifyTxnVerbose", func(tx *dbutil.Tx) error {
		var err error
		head, err = vs.Blockchain.Head(tx)
		if err != nil {
			return err
		}

		uxa, err = vs.Blockchain.Unspent().GetArray(tx, txn.In)
		switch err.(type) {
		case nil:
		case blockdb.ErrUnspentNotExist:
			uxid := err.(blockdb.ErrUnspentNotExist).UxID
			// Gets uxouts of txn.In from historydb
			outs, err := vs.history.GetUxOuts(tx, txn.In)
			if err != nil {
				return err
			}

			if len(outs) == 0 {
				err = fmt.Errorf("transaction input of %s does not exist in either unspent pool or historydb", uxid)
				return NewErrTxnViolatesHardConstraint(err)
			}

			uxa = coin.UxArray{}
			for _, out := range outs {
				uxa = append(uxa, out.Out)
			}

			// Checks if the transaction is confirmed
			txnHash := txn.Hash()
			historyTxn, err := vs.history.GetTransaction(tx, txnHash)
			if err != nil {
				return fmt.Errorf("get transaction of %v from historydb failed: %v", txnHash, err)
			}

			if historyTxn != nil {
				// Transaction is confirmed
				isTxnConfirmed = true
			}

			return nil
		default:
			return err
		}

		if err := VerifySingleTxnUserConstraints(*txn); err != nil {
			return err
		}

		if err := VerifySingleTxnSoftConstraints(*txn, head.Time(), uxa, vs.Config.MaxBlockSize); err != nil {
			return err
		}

		return VerifySingleTxnHardConstraints(*txn, head, uxa)
	})

	// If we were able to query the inputs, return the verbose inputs to the caller
	// even if the transaction failed validation
	var uxs []wallet.UxBalance
	if len(uxa) != 0 {
		var otherErr error
		uxs, otherErr = wallet.NewUxBalances(head.Time(), uxa)
		if otherErr != nil {
			return nil, isTxnConfirmed, otherErr
		}
	}

	return uxs, isTxnConfirmed, err
}

// AddressCount returns the total number of addresses with unspents
func (vs *Visor) AddressCount() (uint64, error) {
	var count uint64
	if err := vs.DB.View("AddressCount", func(tx *dbutil.Tx) error {
		var err error
		count, err = vs.Blockchain.Unspent().AddressCount(tx)
		return err
	}); err != nil {
		return 0, err
	}

	return count, nil
}

// CreateTransactionDeprecated creates a transaction using an entire wallet,
// specifying only coins and one destination
func (vs *Visor) CreateTransactionDeprecated(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	w, err := vs.Wallets.GetWallet(wltID)
	if err != nil {
		logger.WithError(err).Error("Wallets.GetWallet failed")
		return nil, err
	}

	// Get all addresses from the wallet for checking params against
	addrs := w.GetAddresses()

	var auxs coin.AddressUxOuts
	var head *coin.SignedBlock

	if err := vs.DB.View("CreateTransactionDeprecated", func(tx *dbutil.Tx) error {
		head, err = vs.Blockchain.Head(tx)
		if err != nil {
			logger.Errorf("Blockchain.Head failed: %v", err)
			return err
		}

		// Get unspent outputs, while checking that there are no unconfirmed outputs
		auxs, err = vs.getUnspentsForSpending(tx, addrs, false)
		if err != nil {
			if err != wallet.ErrSpendingUnconfirmed {
				logger.WithError(err).Error("getUnspentsForSpending failed")
			}
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Create and sign transaction
	var txn *coin.Transaction
	if err := vs.Wallets.ViewWallet(w, password, func(w *wallet.Wallet) error {
		var err error
		txn, err = w.CreateAndSignTransaction(auxs, head.Time(), coins, dest)
		return err
	}); err != nil {
		logger.WithError(err).Error("CreateAndSignTransaction failed")
		return nil, err
	}

	// The wallet can create transactions that would not pass all validation, such as the decimal restriction,
	// because the wallet is not aware of visor-level constraints.
	// Check that the transaction is valid before returning it to the caller.
	// NOTE: this isn't inside the database transaction, but it's safe,
	// if a racing database write caused this transaction to be invalid, it would be caught here
	if err := VerifySingleTxnUserConstraints(*txn); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction constraints")
		return nil, err
	}
	if err := vs.DB.View("VerifySingleTxnSoftHardConstraints", func(tx *dbutil.Tx) error {
		return vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, vs.Config.MaxBlockSize)
	}); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction constraints")
		return nil, err
	}

	return txn, nil
}

// CreateTransaction creates a transaction based upon the parameters in wallet.CreateTransactionParams
func (vs *Visor) CreateTransaction(params wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error) {
	if err := params.Validate(); err != nil {
		return nil, nil, err
	}

	w, err := vs.Wallets.GetWallet(params.Wallet.ID)
	if err != nil {
		logger.WithError(err).Error("Wallets.GetWallet failed")
		return nil, nil, err
	}

	// Get all addresses from the wallet for checking params against
	allAddrs := w.GetAddresses()

	var auxs coin.AddressUxOuts
	var head *coin.SignedBlock

	if err := vs.DB.View("CreateTransaction", func(tx *dbutil.Tx) error {
		var err error
		head, err = vs.Blockchain.Head(tx)
		if err != nil {
			logger.WithError(err).Error("Blockchain.Head failed")
			return err
		}

		auxs, err = vs.getCreateTransactionAuxs(tx, params, allAddrs)
		return err
	}); err != nil {
		return nil, nil, err
	}

	// Create and sign transaction
	var txn *coin.Transaction
	var inputs []wallet.UxBalance
	if err := vs.Wallets.ViewWallet(w, params.Wallet.Password, func(w *wallet.Wallet) error {
		var err error
		txn, inputs, err = w.CreateAndSignTransactionAdvanced(params, auxs, head.Time())
		return err
	}); err != nil {
		logger.WithError(err).Error("CreateAndSignTransactionAdvanced failed")
		return nil, nil, err
	}

	// The wallet can create transactions that would not pass all validation, such as the decimal restriction,
	// because the wallet is not aware of visor-level constraints.
	// Check that the transaction is valid before returning it to the caller.
	// NOTE: this isn't inside the database transaction, but it's safe,
	// if a racing database write caused this transaction to be invalid, it would be caught here
	if err := VerifySingleTxnUserConstraints(*txn); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction constraints")
		return nil, nil, err
	}
	if err := vs.DB.View("VerifySingleTxnSoftHardConstraints", func(tx *dbutil.Tx) error {
		return vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, vs.Config.MaxBlockSize)
	}); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction constraints")
		return nil, nil, err
	}

	return txn, inputs, nil
}

func (vs *Visor) getCreateTransactionAuxs(tx *dbutil.Tx, params wallet.CreateTransactionParams, allAddrs []cipher.Address) (coin.AddressUxOuts, error) {
	allAddrsMap := make(map[cipher.Address]struct{}, len(allAddrs))
	for _, a := range allAddrs {
		allAddrsMap[a] = struct{}{}
	}

	var auxs coin.AddressUxOuts
	if len(params.Wallet.UxOuts) != 0 {
		// Check if any of the outputs are in an unconfirmed spend
		hashesMap := make(map[cipher.SHA256]struct{}, len(params.Wallet.UxOuts))
		for _, h := range params.Wallet.UxOuts {
			hashesMap[h] = struct{}{}
		}

		// Get all unconfirmed spending uxouts
		unconfirmedTxns, err := vs.Unconfirmed.RawTxns(tx)
		if err != nil {
			return nil, err
		}

		var unconfirmedSpends []cipher.SHA256
		for _, txn := range unconfirmedTxns {
			unconfirmedSpends = append(unconfirmedSpends, txn.In...)
		}

		if params.IgnoreUnconfirmed {
			// Filter unconfirmed spends
			prevLen := len(hashesMap)
			for _, h := range unconfirmedSpends {
				delete(hashesMap, h)
			}

			if prevLen != len(hashesMap) {
				params.Wallet.UxOuts = make([]cipher.SHA256, 0, len(hashesMap))
				for h := range hashesMap {
					params.Wallet.UxOuts = append(params.Wallet.UxOuts, h)
				}
			}
		} else {
			for _, h := range unconfirmedSpends {
				if _, ok := hashesMap[h]; ok {
					return nil, wallet.ErrSpendingUnconfirmed
				}
			}
		}

		// Retrieve the uxouts from the pool.
		// An error is returned if any do not exist
		uxouts, err := vs.Blockchain.Unspent().GetArray(tx, params.Wallet.UxOuts)
		if err != nil {
			return nil, err
		}

		// Build coin.AddressUxOuts map, and check that the address is in the wallets
		auxs = make(coin.AddressUxOuts)
		for _, o := range uxouts {
			if _, ok := allAddrsMap[o.Body.Address]; !ok {
				return nil, wallet.ErrUnknownUxOut
			}
			auxs[o.Body.Address] = append(auxs[o.Body.Address], o)
		}

	} else {
		addrs := params.Wallet.Addresses
		if len(addrs) == 0 {
			addrs = allAddrs
		} else {
			// Check that requested addresses are in the wallet
			for _, a := range addrs {
				if _, ok := allAddrsMap[a]; !ok {
					return nil, wallet.ErrUnknownAddress
				}
			}
		}

		// Get unspent outputs, while checking that there are no unconfirmed outputs
		var err error
		auxs, err = vs.getUnspentsForSpending(tx, addrs, params.IgnoreUnconfirmed)
		if err != nil {
			return nil, err
		}
	}

	return auxs, nil
}

// getUnspentsForSpending returns the unspent outputs for a set of addresses,
// but returns an error if any of the unspents are in the unconfirmed outputs pool
func (vs *Visor) getUnspentsForSpending(tx *dbutil.Tx, addrs []cipher.Address, ignoredUnconfirmed bool) (coin.AddressUxOuts, error) {
	unconfirmedAuxs, err := vs.unconfirmedSpendsOfAddresses(tx, addrs)
	if err != nil {
		err = fmt.Errorf("UnconfirmedSpendsOfAddresses failed: %v", err)
		return nil, err
	}

	if !ignoredUnconfirmed {
		// Check that this is not trying to spend unconfirmed outputs
		if len(unconfirmedAuxs) > 0 {
			return nil, wallet.ErrSpendingUnconfirmed
		}
	}

	auxs, err := vs.Blockchain.Unspent().GetUnspentsOfAddrs(tx, addrs)
	if err != nil {
		err = fmt.Errorf("GetUnspentsOfAddrs failed: %v", err)
		return nil, err
	}

	// Filter unconfirmed
	if ignoredUnconfirmed && len(unconfirmedAuxs) > 0 {
		auxs = auxs.Sub(unconfirmedAuxs)
	}

	return auxs, nil
}
