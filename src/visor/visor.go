package visor

import (
	"errors"
	"fmt"

	"time"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	// MaxDropletPrecision represents the precision of droplets
	MaxDropletPrecision = 0
	// MaxDropletDivisor represents the modulus divisor when checking droplet precision rules
	MaxDropletDivisor uint64 = 1e6
)

var (
	logger = logging.MustGetLogger("visor")

	// ErrInvalidDecimals is returned by DropletPrecisionCheck if a coin amount has an invalid number of decimal places
	ErrInvalidDecimals = errors.New("invalid amount, too many decimal places")
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

	// Modulus divisor of coin amount to control decimal place precision
	// Valid values are 1e6, 1e5, 1e4, 1e3, 1e2, 1
	MaxDropletDivisor uint64

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
		// UnconfirmedRefreshRate:   time.Minute * 30,
		UnconfirmedResendPeriod: time.Minute,
		MaxBlockSize:            1024 * 32,
		MaxDropletDivisor:       MaxDropletDivisor,

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

	switch c.MaxDropletDivisor {
	case 1e6, 1e5, 1e4, 1e3, 1e2, 1:
	default:
		return errors.New("MaxDropletDivisor must be 1e6, 1e5, 1e4, 1e3, 1e2 or 1")
	}

	return nil
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
	wallets  *wallet.Service
	db       *bolt.DB
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

	wltServ, err := wallet.NewService(c.WalletDirectory)
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

	if err := vs.processUnconfirmedTxns(); err != nil {
		return err
	}

	return vs.bcParser.Run()
}

// Shutdown shuts down the visor
func (vs *Visor) Shutdown() {
	defer logger.Info("DB and BlockchainParser closed")

	vs.bcParser.Shutdown()

	if err := vs.db.Close(); err != nil {
		logger.Error("db.Close() error: %v", err)
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
		logger.Info("Genesis block signature=%s", sb.Sig.Hex())
	} else {
		sb = coin.SignedBlock{
			Block: *b,
			Sig:   vs.Config.GenesisSignature,
		}
	}

	return vs.ExecuteSignedBlock(sb)
}

// check if there're unconfirmed transactions that are actually
// already executed, and remove them if any
func (vs *Visor) processUnconfirmedTxns() error {
	removeTxs := []cipher.SHA256{}
	vs.Unconfirmed.ForEach(func(hash cipher.SHA256, tx *UnconfirmedTxn) error {
		// check if the tx already executed
		if err := vs.Blockchain.VerifyTransaction(tx.Txn); err != nil {
			removeTxs = append(removeTxs, hash)
		}

		txn, err := vs.history.GetTransaction(hash)
		if err != nil {
			return fmt.Errorf("process unconfirmed txs failed: %v", err)
		}

		if txn != nil {
			removeTxs = append(removeTxs, hash)
		}

		return nil
	})

	if len(removeTxs) > 0 {
		vs.Unconfirmed.RemoveTransactions(removeTxs)
	}

	return nil
}

// GenesisPreconditions panics if conditions for genesis block are not met
func (vs *Visor) GenesisPreconditions() {
	if vs.Config.BlockchainSeckey != (cipher.SecKey{}) {
		if vs.Config.BlockchainPubkey != cipher.PubKeyFromSecKey(vs.Config.BlockchainSeckey) {
			logger.Panicf("Cannot create genesis block. Invalid secret key for pubkey")
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
	if !vs.Config.IsMaster {
		logger.Panic("Only master chain can create blocks")
	}

	var sb coin.SignedBlock
	if vs.Unconfirmed.Len() == 0 {
		return sb, errors.New("No transactions")
	}

	// Gather all unconfirmed transactions
	txns := vs.Unconfirmed.RawTxns()
	logger.Info("Unconfirmed pool has %d transactions pending", len(txns))

	// Sort them by highest fee per kilobyte
	txns = coin.SortTransactions(txns, vs.Blockchain.TransactionFee)

	// Filter transactions that do not obey droplet precision rules
	var filteredTxns coin.Transactions
	for _, txn := range txns {
		skip := false
		for _, o := range txn.Out {
			if err := dropletPrecisionCheck(o.Coins, vs.Config.MaxDropletDivisor); err != nil {
				skip = true
				break
			}
		}

		if !skip {
			filteredTxns = append(filteredTxns, txn)
		}
	}

	nRemoved := len(txns) - len(filteredTxns)
	if nRemoved > 0 {
		logger.Info("CreateBlock ignored %d transactions with too many decimal places", nRemoved)
	}

	txns = filteredTxns

	// Apply block size transaction limit
	txns = txns.TruncateBytesTo(vs.Config.MaxBlockSize)

	logger.Info("Creating new block with %d transactions, head time %d", len(txns), when)

	b, err := vs.Blockchain.NewBlock(txns, when)
	if err != nil {
		logger.Warning("Blockchain.NewBlock failed: %v", err)
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
	if err := vs.verifySignedBlock(&b); err != nil {
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

// GetSignedBlocksSince returns N signed blocks more recent than Seq. Does not return nil.
func (vs *Visor) GetSignedBlocksSince(seq, ct uint64) ([]coin.SignedBlock, error) {
	avail := uint64(0)
	head, err := vs.Blockchain.Head()
	if err != nil {
		return []coin.SignedBlock{}, err
	}

	headSeq := head.Seq()
	if headSeq > seq {
		avail = headSeq - seq
	}
	if avail < ct {
		ct = avail
	}
	if ct == 0 {
		return []coin.SignedBlock{}, nil
	}
	blocks := make([]coin.SignedBlock, 0, ct)
	for j := uint64(0); j < ct; j++ {
		i := seq + 1 + j
		b, err := vs.Blockchain.GetBlockBySeq(i)
		if err != nil {
			return []coin.SignedBlock{}, err
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

// InjectTxn records a coin.Transaction to the UnconfirmedTxnPool if the txn is not
// already in the blockchain
// TODO
// - rename InjectTransaction
// Refactor
// Why do does this return both error and bool
func (vs *Visor) InjectTxn(txn coin.Transaction) (bool, error) {
	// Ignore transactions that do not conform to decimal restrictions
	for _, o := range txn.Out {
		if err := DropletPrecisionCheck(o.Coins); err != nil {
			return false, err
		}
	}

	return vs.Unconfirmed.InjectTxn(vs.Blockchain, txn)
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

// GetLastTxs returns last confirmed transactions, return nil if empty
func (vs *Visor) GetLastTxs() ([]*Transaction, error) {
	ltxs, err := vs.history.GetLastTxs()
	if err != nil {
		return nil, err
	}

	txs := make([]*Transaction, len(ltxs))
	var confirms uint64
	bh := vs.HeadBkSeq()
	var b *coin.SignedBlock
	for i, tx := range ltxs {
		confirms = uint64(bh) - tx.BlockSeq + 1
		b, err = vs.GetBlockBySeq(tx.BlockSeq)
		if err != nil {
			return nil, err
		}

		if b == nil {
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

// LoadAndScanWallet loads wallet from seed and scan ahead N addresses
func (vs Visor) LoadAndScanWallet(wltName string, seed string, scanN uint64, ops ...wallet.Option) (wallet.Wallet, error) {
	return vs.wallets.LoadAndScanWallet(wltName, seed, scanN, vs, ops...)
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

		coins := uxs.Coins()
		coinHours := uxs.CoinHours(headTime)
		pcoins := predictedUxs.Coins()
		pcoinHours := predictedUxs.CoinHours(headTime)
		bp := wallet.BalancePair{
			Confirmed: wallet.Balance{Coins: coins, Hours: coinHours},
			Predicted: wallet.Balance{Coins: pcoins, Hours: pcoinHours},
		}

		bps = append(bps, bp)
	}

	return bps, nil
}

// DropletPrecisionCheck checks if the amount is valid
func DropletPrecisionCheck(amount uint64) error {
	return dropletPrecisionCheck(amount, MaxDropletDivisor)
}

func dropletPrecisionCheck(amount, divisor uint64) error {
	if amount%divisor != 0 {
		return ErrInvalidDecimals
	}

	return nil
}
