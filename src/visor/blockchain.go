package visor

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

const (
	// DebugLevel1 checks for extremely unlikely conditions (10e-40)
	DebugLevel1 = true
	// DebugLevel2 enable checks for impossible conditions
	DebugLevel2 = true
)

var (
	// ErrVerifyStopped is returned when database verification is interrupted
	ErrVerifyStopped = errors.New("database verification stopped")
)

// ErrBlockNotExist may be returned if a block is not found
type ErrBlockNotExist struct {
	Seq uint64
}

// NewErrBlockNotExist creates an ErrBlockNotExist based on an unknown block sequence
func NewErrBlockNotExist(seq uint64) ErrBlockNotExist {
	return ErrBlockNotExist{
		Seq: seq,
	}
}

func (e ErrBlockNotExist) Error() string {
	return fmt.Sprintf("block does not exist seq=%d", e.Seq)
}

//Warning: 10e6 is 10 million, 1e6 is 1 million

// Note: DebugLevel1 adds additional checks for hash collisions that
// are unlikely to occur. DebugLevel2 adds checks for conditions that
// can only occur through programmer error and malice.

// Note: a droplet is the base coin unit. Each Skycoin is one million droplets

//Termonology:
// UXTO - unspent transaction outputs
// UX - outputs10
// TX - transactions

//Notes:
// transactions (TX) consume outputs (UX) and produce new outputs (UX)
// Tx.Uxi() - set of outputs consumed by transaction
// Tx.Uxo() - set of outputs created by transaction

// chainStore
type chainStore interface {
	Head(*dbutil.Tx) (*coin.SignedBlock, error)
	HeadSeq(*dbutil.Tx) (uint64, bool, error)
	Len(*dbutil.Tx) (uint64, error)
	AddBlock(*dbutil.Tx, *coin.SignedBlock) error
	GetBlockByHash(*dbutil.Tx, cipher.SHA256) (*coin.Block, error)
	GetSignedBlockByHash(*dbutil.Tx, cipher.SHA256) (*coin.SignedBlock, error)
	GetSignedBlockBySeq(*dbutil.Tx, uint64) (*coin.SignedBlock, error)
	UnspentPool() blockdb.UnspentPooler
	GetGenesisBlock(*dbutil.Tx) (*coin.SignedBlock, error)
	GetBlockSignature(*dbutil.Tx, *coin.Block) (cipher.Sig, bool, error)
	ForEachBlock(*dbutil.Tx, func(*coin.Block) error) error
}

// DefaultWalker default blockchain walker
func DefaultWalker(tx *dbutil.Tx, hps []coin.HashPair) (cipher.SHA256, bool) {
	if len(hps) == 0 {
		return cipher.SHA256{}, false
	}
	return hps[0].Hash, true
}

// CreateBuckets creates the buckets used by the blockdb
func CreateBuckets(db *dbutil.DB) error {
	return db.Update("CreateBuckets", func(tx *dbutil.Tx) error {
		if err := historydb.CreateBuckets(tx); err != nil {
			return err
		}

		if err := blockdb.CreateBuckets(tx); err != nil {
			return err
		}

		return dbutil.CreateBuckets(tx, [][]byte{
			UnconfirmedTxnsBkt,
			UnconfirmedUnspentsBkt,
		})
	})
}

// BlockchainConfig configures Blockchain options
type BlockchainConfig struct {
	// Arbitrating mode: if in arbitrating mode, when block publishing node execute blocks,
	// the invalid transaction will be skipped and continue the next; otherwise,
	// node will throw the error and return.
	Arbitrating bool
	Pubkey      cipher.PubKey
}

// Blockchain maintains blockchain and provides apis for accessing the chain.
type Blockchain struct {
	db    *dbutil.DB
	cfg   BlockchainConfig
	store chainStore
}

// NewBlockchain creates a Blockchain
func NewBlockchain(db *dbutil.DB, cfg BlockchainConfig) (*Blockchain, error) {
	chainstore, err := blockdb.NewBlockchain(db, DefaultWalker)
	if err != nil {
		return nil, err
	}

	return &Blockchain{
		cfg:   cfg,
		db:    db,
		store: chainstore,
	}, nil
}

// GetGenesisBlock returns genesis block
func (bc *Blockchain) GetGenesisBlock(tx *dbutil.Tx) (*coin.SignedBlock, error) {
	return bc.store.GetGenesisBlock(tx)
}

// GetSignedBlockByHash returns block of given hash
func (bc *Blockchain) GetSignedBlockByHash(tx *dbutil.Tx, hash cipher.SHA256) (*coin.SignedBlock, error) {
	return bc.store.GetSignedBlockByHash(tx, hash)
}

// GetSignedBlockBySeq returns block of given seq
func (bc *Blockchain) GetSignedBlockBySeq(tx *dbutil.Tx, seq uint64) (*coin.SignedBlock, error) {
	return bc.store.GetSignedBlockBySeq(tx, seq)
}

// Head returns the most recent confirmed block
func (bc Blockchain) Head(tx *dbutil.Tx) (*coin.SignedBlock, error) {
	return bc.store.Head(tx)
}

// Unspent returns the unspent outputs pool
func (bc *Blockchain) Unspent() blockdb.UnspentPooler {
	return bc.store.UnspentPool()
}

// Len returns the length of current blockchain.
func (bc Blockchain) Len(tx *dbutil.Tx) (uint64, error) {
	return bc.store.Len(tx)
}

// HeadSeq returns the sequence of head block
func (bc *Blockchain) HeadSeq(tx *dbutil.Tx) (uint64, bool, error) {
	return bc.store.HeadSeq(tx)
}

// Time returns time of last block
// used as system clock indepedent clock for coin hour calculations
// TODO: Deprecate
func (bc *Blockchain) Time(tx *dbutil.Tx) (uint64, error) {
	b, err := bc.Head(tx)
	if err != nil {
		if err == blockdb.ErrNoHeadBlock {
			return 0, nil
		}
		return 0, err
	}

	return b.Time(), nil
}

// NewBlock creates a Block given an array of Transactions.
// Only hard constraints are applied to transactions in the block.
// The caller of this function should apply any additional soft constraints,
// and choose which transactions to place into the block.
func (bc Blockchain) NewBlock(tx *dbutil.Tx, txns coin.Transactions, currentTime uint64) (*coin.Block, error) {
	if len(txns) == 0 {
		return nil, errors.New("No transactions")
	}

	head, err := bc.store.Head(tx)
	if err != nil {
		return nil, err
	}

	if currentTime <= head.Time() {
		return nil, errors.New("Time can only move forward")
	}

	txns, err = bc.processTransactions(tx, txns)
	if err != nil {
		return nil, err
	}

	uxHash, err := bc.Unspent().GetUxHash(tx)
	if err != nil {
		return nil, err
	}

	feeCalc := bc.TransactionFee(tx, head.Time())

	b, err := coin.NewBlock(head.Block, currentTime, uxHash, txns, feeCalc)
	if err != nil {
		return nil, err
	}

	// make sure block is valid
	if DebugLevel2 {
		if err := bc.verifyBlockHeader(tx, *b); err != nil {
			return nil, err
		}
		txns, err := bc.processTransactions(tx, b.Body.Transactions)
		if err != nil {
			logger.Panicf("bc.processTransactions second verification call failed: %v", err)
		}
		b.Body.Transactions = txns
	}
	return b, nil
}

func (bc *Blockchain) processBlock(tx *dbutil.Tx, b coin.SignedBlock) (coin.SignedBlock, error) {
	length, err := bc.Len(tx)
	if err != nil {
		return coin.SignedBlock{}, err
	}

	if length > 0 {
		if isGenesis, err := bc.isGenesisBlock(tx, b.Block); err != nil {
			return coin.SignedBlock{}, err
		} else if isGenesis {
			err := errors.New("Attempted to process genesis block after blockchain has genesis block")
			logger.Warning(err.Error())
			return coin.SignedBlock{}, err
		} else {
			if err := bc.verifyBlockHeader(tx, b.Block); err != nil {
				return coin.SignedBlock{}, err
			}

			txns, err := bc.processTransactions(tx, b.Body.Transactions)
			if err != nil {
				return coin.SignedBlock{}, err
			}
			b.Body.Transactions = txns

			if err := bc.verifyUxHash(tx, b.Block); err != nil {
				return coin.SignedBlock{}, err
			}

		}
	}

	return b, nil
}

// ExecuteBlock attempts to append block to blockchain with *dbutil.Tx
func (bc *Blockchain) ExecuteBlock(tx *dbutil.Tx, sb *coin.SignedBlock) error {
	length, err := bc.Len(tx)
	if err != nil {
		return err
	}

	if length > 0 {
		head, err := bc.Head(tx)
		if err != nil {
			return err
		}

		// TODO -- why do we modify the block here?
		sb.Head.PrevHash = head.HashHeader()
	}

	nb, err := bc.processBlock(tx, *sb)
	if err != nil {
		return err
	}

	if err := bc.store.AddBlock(tx, &nb); err != nil {
		return err
	}

	return nil
}

// VerifyBlock verifies specified block against current state of blockchain.
func (bc *Blockchain) VerifyBlock(tx *dbutil.Tx, sb *coin.SignedBlock) error {
	_, err := bc.processBlock(tx, *sb)

	return err
}

// isGenesisBlock checks if the block is genesis block
func (bc Blockchain) isGenesisBlock(tx *dbutil.Tx, b coin.Block) (bool, error) {
	gb, err := bc.store.GetGenesisBlock(tx)
	if err != nil {
		return false, err
	}
	if gb == nil {
		return false, nil
	}

	return gb.HashHeader() == b.HashHeader(), nil
}

// Compares the state of the current UxHash hash to state of unspent
// output pool.
func (bc Blockchain) verifyUxHash(tx *dbutil.Tx, b coin.Block) error {
	uxHash, err := bc.Unspent().GetUxHash(tx)
	if err != nil {
		return err
	}

	if !bytes.Equal(b.Head.UxHash[:], uxHash[:]) {
		return errors.New("UxHash does not match")
	}

	return nil
}

// VerifyBlockTxnConstraints checks that the transaction does not violate hard constraints,
// for transactions that are already included in a block.
func (bc Blockchain) VerifyBlockTxnConstraints(tx *dbutil.Tx, txn coin.Transaction) error {
	// NOTE: Unspent().GetArray() returns an error if not all txn.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx, txn.In)
	if err != nil {
		switch err.(type) {
		case blockdb.ErrUnspentNotExist:
			return transaction.NewErrTxnViolatesHardConstraint(err)
		default:
			return err
		}
	}

	head, err := bc.Head(tx)
	if err != nil {
		return err
	}

	return bc.verifyBlockTxnHardConstraints(tx, txn, head, uxIn)
}

func (bc Blockchain) verifyBlockTxnHardConstraints(tx *dbutil.Tx, txn coin.Transaction, head *coin.SignedBlock, uxIn coin.UxArray) error {
	if err := transaction.VerifyBlockTxnConstraints(txn, head.Head, uxIn); err != nil {
		return err
	}

	if DebugLevel1 {
		// Check that new unspents don't collide with existing.
		// This should not occur but is a sanity check.
		// NOTE: this is not in the top-level VerifyBlockTxnConstraints
		// because it relies on the unspent pool to check for existence.
		// For remote callers such as the CLI, they'd need to download the whole
		// unspent pool or make a separate API call to check for duplicate unspents.
		uxOuts := coin.CreateUnspents(head.Head, txn)
		for i := range uxOuts {
			if contains, err := bc.Unspent().Contains(tx, uxOuts[i].Hash()); err != nil {
				return err
			} else if contains {
				err := errors.New("New unspent collides with existing unspent")
				return transaction.NewErrTxnViolatesHardConstraint(err)
			}
		}
	}

	return nil
}

// VerifySingleTxnHardConstraints checks that the transaction does not violate hard constraints.
// for transactions that are not included in a block.
func (bc Blockchain) VerifySingleTxnHardConstraints(tx *dbutil.Tx, txn coin.Transaction, signed transaction.TxnSignedFlag) error {
	// NOTE: Unspent().GetArray() returns an error if not all txn.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx, txn.In)
	if err != nil {
		switch err.(type) {
		case blockdb.ErrUnspentNotExist:
			return transaction.NewErrTxnViolatesHardConstraint(err)
		default:
			return err
		}
	}

	head, err := bc.Head(tx)
	if err != nil {
		return err
	}

	return bc.verifySingleTxnHardConstraints(tx, txn, head, uxIn, signed)
}

// VerifySingleTxnSoftHardConstraints checks that the transaction does not violate hard or soft constraints,
// for transactions that are not included in a block.
// Hard constraints are checked before soft constraints.
func (bc Blockchain) VerifySingleTxnSoftHardConstraints(tx *dbutil.Tx, txn coin.Transaction, distParams params.Distribution, verifyParams params.VerifyTxn, signed transaction.TxnSignedFlag) (*coin.SignedBlock, coin.UxArray, error) {
	// NOTE: Unspent().GetArray() returns an error if not all txn.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx, txn.In)
	if err != nil {
		return nil, nil, transaction.NewErrTxnViolatesHardConstraint(err)
	}

	head, err := bc.Head(tx)
	if err != nil {
		return nil, nil, err
	}

	// Hard constraints must be checked before soft constraints
	if err := bc.verifySingleTxnHardConstraints(tx, txn, head, uxIn, signed); err != nil {
		return nil, nil, err
	}

	if err := transaction.VerifySingleTxnSoftConstraints(txn, head.Time(), uxIn, distParams, verifyParams); err != nil {
		return nil, nil, err
	}

	return head, uxIn, nil
}

func (bc Blockchain) verifySingleTxnHardConstraints(tx *dbutil.Tx, txn coin.Transaction, head *coin.SignedBlock, uxIn coin.UxArray, signed transaction.TxnSignedFlag) error {
	if err := transaction.VerifySingleTxnHardConstraints(txn, head.Head, uxIn, signed); err != nil {
		return err
	}

	if DebugLevel1 {
		// Check that new unspents don't collide with existing.
		// This should not occur but is a sanity check.
		// NOTE: this is not in the top-level VerifySingleTxnHardConstraints
		// because it relies on the unspent pool to check for existence.
		// For remote callers such as the CLI, they'd need to download the whole
		// unspent pool or make a separate API call to check for duplicate unspents.
		uxOuts := coin.CreateUnspents(head.Head, txn)
		for i := range uxOuts {
			if contains, err := bc.Unspent().Contains(tx, uxOuts[i].Hash()); err != nil {
				return err
			} else if contains {
				err := errors.New("New unspent collides with existing unspent")
				return transaction.NewErrTxnViolatesHardConstraint(err)
			}
		}
	}

	return nil
}

// GetBlocks returns blocks matching seqs. If any block is not found, returns an error.
func (bc Blockchain) GetBlocks(tx *dbutil.Tx, seqs []uint64) ([]coin.SignedBlock, error) {
	blocks := make([]coin.SignedBlock, len(seqs))

	for i, s := range seqs {
		b, err := bc.store.GetSignedBlockBySeq(tx, s)
		if err != nil {
			return nil, err
		}

		if b == nil {
			return nil, NewErrBlockNotExist(s)
		}

		blocks[i] = *b
	}

	return blocks, nil
}

// GetBlocksInRange return blocks whose seq are in the range of start and end.
func (bc Blockchain) GetBlocksInRange(tx *dbutil.Tx, start, end uint64) ([]coin.SignedBlock, error) {
	if start > end {
		return nil, nil
	}

	var blocks []coin.SignedBlock
	for i := start; i <= end; i++ {
		b, err := bc.store.GetSignedBlockBySeq(tx, i)
		if err != nil {
			logger.WithError(err).Error("bc.store.GetSignedBlockBySeq failed")
			return nil, err
		}

		if b == nil {
			break
		}

		blocks = append(blocks, *b)
	}

	return blocks, nil
}

// GetLastBlocks return the latest N blocks.
func (bc Blockchain) GetLastBlocks(tx *dbutil.Tx, num uint64) ([]coin.SignedBlock, error) {
	if num == 0 {
		return nil, nil
	}

	end, ok, err := bc.HeadSeq(tx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	start := int(end-num) + 1
	if start < 0 {
		start = 0
	}

	return bc.GetBlocksInRange(tx, uint64(start), end)
}

/* Private */

// Validates a set of Transactions, individually, against each other and
// against the Blockchain.  If firstFail is true, it will return an error
// as soon as it encounters one.  Else, it will return an array of
// Transactions that are valid as a whole.  It may return an error if
// firstFalse is false, if there is no way to filter the txns into a valid
// array, i.e. processTransactions(processTransactions(txn, false), true)
// should not result in an error, unless all txns are invalid.
// TODO:
//  - move arbitration to visor
//  - blockchain should have strict checking
func (bc Blockchain) processTransactions(tx *dbutil.Tx, txs coin.Transactions) (coin.Transactions, error) {
	// copy txs so that the following code won't modify the original txns
	txns := make(coin.Transactions, len(txs))
	copy(txns, txs)

	head, err := bc.store.Head(tx)
	if err != nil {
		return nil, err
	}

	// Transactions need to be sorted by fee and hash before arbitrating
	if bc.cfg.Arbitrating {
		txns, err = coin.SortTransactions(txns, bc.TransactionFee(tx, head.Time()))
		if err != nil {
			logger.Critical().WithError(err).Error("processTransactions: coin.SortTransactions failed")
			return nil, err
		}
	}

	//TODO: audit
	if len(txns) == 0 {
		if bc.cfg.Arbitrating {
			return txns, nil
		}

		// If there are no transactions, a block should not be made
		return nil, errors.New("No transactions")
	}

	skip := make(map[int]struct{})
	uxHashes := make(coin.UxHashSet, len(txns))
	for i, txn := range txns {
		// Check the transaction against itself.  This covers the hash,
		// signature indices and duplicate spends within itself
		if err := bc.VerifyBlockTxnConstraints(tx, txn); err != nil {
			switch err.(type) {
			case transaction.ErrTxnViolatesSoftConstraint:
				logger.Critical().WithError(err).Panic("bc.VerifyBlockTxnConstraints should not return a ErrTxnViolatesSoftConstraint error")
			case transaction.ErrTxnViolatesHardConstraint:
				if bc.cfg.Arbitrating {
					skip[i] = struct{}{}
					continue
				}
			}

			return nil, err
		}

		// Check that each pending unspent will be unique
		uxb := coin.UxBody{
			SrcTransaction: txn.Hash(),
		}

		for _, to := range txn.Out {
			uxb.Coins = to.Coins
			uxb.Hours = to.Hours
			uxb.Address = to.Address

			h := uxb.Hash()
			_, exists := uxHashes[h]
			if exists {
				if bc.cfg.Arbitrating {
					skip[i] = struct{}{}
					continue
				} else {
					return nil, errors.New("Duplicate unspent output across transactions")
				}
			}

			if DebugLevel1 {
				// Check that the expected unspent is not already in the pool.
				// This should never happen because its a hash collision
				if contains, err := bc.Unspent().Contains(tx, h); err != nil {
					return nil, err
				} else if contains {
					if bc.cfg.Arbitrating {
						skip[i] = struct{}{}
						continue
					} else {
						return nil, errors.New("Output hash is in the UnspentPool")
					}
				}
			}

			uxHashes[h] = struct{}{}
		}
	}

	// Filter invalid transactions before arbitrating between colliding ones
	if len(skip) > 0 {
		newtxns := make(coin.Transactions, len(txns)-len(skip))
		j := 0
		for i := range txns {
			if _, shouldSkip := skip[i]; !shouldSkip {
				newtxns[j] = txns[i]
				j++
			}
		}
		txns = newtxns
		skip = make(map[int]struct{})
	}

	// Check to ensure that there are no duplicate spends in the entire block,
	// and that we aren't creating duplicate outputs.  Duplicate outputs
	// within a single Transaction are already checked by VerifyBlockTxnConstraints
	hashes := txns.Hashes()
	for i := 0; i < len(txns)-1; i++ {
		s := txns[i]
		for j := i + 1; j < len(txns); j++ {
			t := txns[j]
			if DebugLevel1 {
				if hashes[i] == hashes[j] {
					// This is a non-recoverable error for filtering, and
					// should never occur.  It indicates a hash collision
					// amongst different txns. Duplicate transactions are
					// caught earlier, when duplicate expected outputs are
					// checked for, and will not trigger this.
					return nil, errors.New("Unexpected duplicate transaction")
				}
			}
			for a := range s.In {
				for b := range t.In {
					if s.In[a] == t.In[b] {
						if bc.cfg.Arbitrating {
							// The txn with the highest fee and lowest hash
							// is chosen when attempting a double spend.
							// Since the txns are sorted, we skip the 2nd
							// iterable
							skip[j] = struct{}{}
						} else {
							return nil, errors.New("Cannot spend output twice in the same block")
						}
					}
				}
			}
		}
	}

	// Filter the final results, if necessary
	if len(skip) > 0 {
		newtxns := make(coin.Transactions, 0, len(txns)-len(skip))
		for i := range txns {
			if _, shouldSkip := skip[i]; !shouldSkip {
				newtxns = append(newtxns, txns[i])
			}
		}
		return newtxns, nil
	}

	return txns, nil
}

// TransactionFee calculates the current transaction fee in coinhours of a Transaction
func (bc Blockchain) TransactionFee(tx *dbutil.Tx, headTime uint64) coin.FeeCalculator {
	return func(txn *coin.Transaction) (uint64, error) {
		inUxs, err := bc.Unspent().GetArray(tx, txn.In)
		if err != nil {
			return 0, err
		}

		return fee.TransactionFee(txn, headTime, inUxs)
	}
}

// VerifySignature checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
func (bc *Blockchain) VerifySignature(block *coin.SignedBlock) error {
	err := block.VerifySignature(bc.cfg.Pubkey)
	if err != nil {
		logger.Errorf("Blockchain signature verification failed for block %d: %v", block.Head.BkSeq, err)
	}
	return err
}

// WalkChain walk through the blockchain concurrently
// The quit channel is optional and if closed, this method still stop.
func (bc *Blockchain) WalkChain(workers int, f func(*dbutil.Tx, *coin.SignedBlock) error, quit chan struct{}) error {
	if quit == nil {
		quit = make(chan struct{})
	}

	signedBlockC := make(chan *coin.SignedBlock, 100)
	errC := make(chan error, 100)
	interrupt := make(chan struct{})
	verifyDone := make(chan struct{})

	// Verify block signatures in a worker pool
	var workerWg sync.WaitGroup
	workerWg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer workerWg.Done()
			if err := bc.db.View("WalkChain verify blocks", func(tx *dbutil.Tx) error {
				for b := range signedBlockC {
					if err := f(tx, b); err != nil {
						// if err := cipher.VerifyPubKeySignedHash(bc.cfg.Pubkey, sh.sig, sh.hash); err != nil {
						// logger.Errorf("Signature verification failed: %v", err)
						select {
						case errC <- err:
						default:
						}
					}
				}
				return nil
			}); err != nil {
				logger.WithError(err).Error("WalkChain verify blocks db transaction failed")
			}
		}()
	}

	// Wait for verification worker goroutines to finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		workerWg.Wait()
		close(verifyDone)
	}()

	// Iterate all blocks stored in the "blocks" bucket
	// * Detect if a corresponding signature is missing from the signatures bucket
	// * Verify the signature for the block
	wg.Add(1)
	go func() {
		if err := bc.db.View("WalkChain get blocks", func(tx *dbutil.Tx) error {
			if length, err := bc.Len(tx); err != nil {
				return err
			} else if length == 0 {
				return nil
			}
			defer wg.Done()
			defer close(signedBlockC)

			errInterrupted := errors.New("goroutine was stopped")

			if err := bc.store.ForEachBlock(tx, func(block *coin.Block) error {
				sig, ok, err := bc.store.GetBlockSignature(tx, block)
				if err != nil {
					return err
				}
				if !ok {
					return blockdb.NewErrMissingSignature(block)
				}

				signedBlock := &coin.SignedBlock{
					Sig:   sig,
					Block: *block,
				}

				select {
				case signedBlockC <- signedBlock:
					return nil
				case <-quit:
					return errInterrupted
				case <-interrupt:
					return errInterrupted
				}
			}); err != nil && err != errInterrupted {
				switch err.(type) {
				case blockdb.ErrMissingSignature:
				default:
					logger.Errorf("bc.store.ForEachBlock failed: %v", err)
				}
				select {
				case errC <- err:
				default:
				}
			}
			return nil
		}); err != nil {
			logger.WithError(err).Error("WalkChain get blocks db transaction failed")
		}
	}()

	var err error
	select {
	case err = <-errC:
		if err != nil {
			break
		}
	case <-quit:
		err = ErrVerifyStopped
		break
	case <-verifyDone:
		break
	}

	close(interrupt)
	wg.Wait()
	return err
}

// VerifyBlockHeader Returns error if the BlockHeader is not valid
func (bc Blockchain) verifyBlockHeader(tx *dbutil.Tx, b coin.Block) error {
	head, err := bc.Head(tx)
	if err != nil {
		return err
	}

	//check BkSeq
	if b.Head.BkSeq != head.Head.BkSeq+1 {
		return errors.New("BkSeq invalid")
	}
	//check Time, only requirement is that its monotonely increasing
	if b.Head.Time <= head.Head.Time {
		return errors.New("Block time must be > head time")
	}
	// Check block hash against previous head
	if b.Head.PrevHash != head.HashHeader() {
		return errors.New("PrevHash does not match current head")
	}

	if b.Body.Hash() != b.Head.BodyHash {
		return errors.New("Computed body hash does not match")
	}
	return nil
}
