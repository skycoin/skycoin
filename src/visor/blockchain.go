package visor

import (
	"bytes"
	"errors"
	"sync"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

const (
	// DebugLevel1 checks for extremely unlikely conditions (10e-40)
	DebugLevel1 = true
	// DebugLevel2 enable checks for impossible conditions
	DebugLevel2 = true

	// SigVerifyTheadNum  signature verifycation goroutine number
	SigVerifyTheadNum = 4
)

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
	Head(*bolt.Tx) (*coin.SignedBlock, error)
	HeadSeq(*bolt.Tx) (uint64, bool, error)
	Len(*bolt.Tx) (uint64, error)
	AddBlock(*bolt.Tx, *coin.SignedBlock) error
	GetBlockByHash(*bolt.Tx, cipher.SHA256) (*coin.Block, error)
	GetSignedBlockByHash(*bolt.Tx, cipher.SHA256) (*coin.SignedBlock, error)
	GetSignedBlockBySeq(*bolt.Tx, uint64) (*coin.SignedBlock, error)
	UnspentPool() blockdb.UnspentPool
	GetGenesisBlock(*bolt.Tx) (*coin.SignedBlock, error)
	GetBlockSignature(*bolt.Tx, *coin.Block) (cipher.Sig, bool, error)
	ForEachBlock(*bolt.Tx, func(*coin.Block) error) error
	ForEachSignature(*bolt.Tx, func(cipher.SHA256, cipher.Sig) error) error
}

// BlockListener notify the register when new block is appended to the chain
type BlockListener func(b coin.Block)

// Blockchain maintains blockchain and provides apis for accessing the chain.
type Blockchain struct {
	db          *dbutil.DB
	pubkey      cipher.PubKey
	blkListener []BlockListener

	// arbitrating mode, if in arbitrating mode, when master node execute blocks,
	// the invalid transaction will be skipped and continue the next; otherwise,
	// node will throw the error and return.
	arbitrating bool
	store       chainStore
}

// Option represents the option when creating the blockchain
type Option func(*Blockchain)

// DefaultWalker default blockchain walker
func DefaultWalker(tx *bolt.Tx, hps []coin.HashPair) (cipher.SHA256, bool) {
	if len(hps) == 0 {
		return cipher.SHA256{}, false
	}
	return hps[0].Hash, true
}

// NewBlockchain use the walker go through the tree and update the head and unspent outputs.
func NewBlockchain(db *dbutil.DB, pubkey cipher.PubKey, ops ...Option) (*Blockchain, error) {
	chainstore, err := blockdb.NewBlockchain(db, DefaultWalker)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		db:     db,
		pubkey: pubkey,
		store:  chainstore,
	}

	for _, op := range ops {
		op(bc)
	}

	// verify signature
	if err := db.View(func(tx *bolt.Tx) error {
		return bc.verifySigs(tx, SigVerifyTheadNum)
	}); err != nil {
		return nil, err
	}

	return bc, nil
}

// Arbitrating option to change the mode
func Arbitrating(enable bool) Option {
	return func(bc *Blockchain) {
		bc.arbitrating = enable
	}
}

// GetGenesisBlock returns genesis block
func (bc *Blockchain) GetGenesisBlock(tx *bolt.Tx) (*coin.SignedBlock, error) {
	return bc.store.GetGenesisBlock(tx)
}

// GetSignedBlockByHash returns block of given hash
func (bc *Blockchain) GetSignedBlockByHash(tx *bolt.Tx, hash cipher.SHA256) (*coin.SignedBlock, error) {
	return bc.store.GetSignedBlockByHash(tx, hash)
}

// GetSignedBlockBySeq returns block of given seq
func (bc *Blockchain) GetSignedBlockBySeq(tx *bolt.Tx, seq uint64) (*coin.SignedBlock, error) {
	return bc.store.GetSignedBlockBySeq(tx, seq)
}

// Head returns the most recent confirmed block
func (bc Blockchain) Head(tx *bolt.Tx) (*coin.SignedBlock, error) {
	return bc.store.Head(tx)
}

// Unspent returns the unspent outputs pool
func (bc *Blockchain) Unspent() blockdb.UnspentPool {
	return bc.store.UnspentPool()
}

// Len returns the length of current blockchain.
func (bc Blockchain) Len(tx *bolt.Tx) (uint64, error) {
	return bc.store.Len(tx)
}

// HeadSeq returns the sequence of head block
func (bc *Blockchain) HeadSeq(tx *bolt.Tx) (uint64, bool, error) {
	return bc.store.HeadSeq(tx)
}

// Time returns time of last block
// used as system clock indepedent clock for coin hour calculations
// TODO: Deprecate
func (bc *Blockchain) Time(tx *bolt.Tx) (uint64, error) {
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
func (bc Blockchain) NewBlock(tx *bolt.Tx, txns coin.Transactions, currentTime uint64) (*coin.Block, error) {
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
	if DebugLevel2 == true {
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

func (bc *Blockchain) processBlock(tx *bolt.Tx, b coin.SignedBlock) (coin.SignedBlock, error) {
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

// ExecuteBlock attempts to append block to blockchain with *bolt.Tx
func (bc *Blockchain) ExecuteBlock(tx *bolt.Tx, sb *coin.SignedBlock) error {
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

// isGenesisBlock checks if the block is genesis block
func (bc Blockchain) isGenesisBlock(tx *bolt.Tx, b coin.Block) (bool, error) {
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
func (bc Blockchain) verifyUxHash(tx *bolt.Tx, b coin.Block) error {
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
func (bc Blockchain) VerifyBlockTxnConstraints(tx *bolt.Tx, txn coin.Transaction) error {
	// NOTE: Unspent().GetArray() returns an error if not all txn.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx, txn.In)
	if err != nil {
		switch err.(type) {
		case blockdb.ErrUnspentNotExist:
			return NewErrTxnViolatesHardConstraint(err)
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

func (bc Blockchain) verifyBlockTxnHardConstraints(tx *bolt.Tx, txn coin.Transaction, head *coin.SignedBlock, uxIn coin.UxArray) error {
	if err := VerifyBlockTxnConstraints(txn, head, uxIn); err != nil {
		return err
	}

	if DebugLevel1 {
		// Check that new unspents don't collide with existing.
		// This should not occur but is a sanity check.
		// NOTE: this is not in the top-level VerifyBlockTxnConstraints
		// because it relies on the unspent pool to check for existence.
		// For remote callers such as the CLI, they'd need to download the whole
		// unspent pool or make a separate API call to check for duplicate unspents.
		uxOut := coin.CreateUnspents(head.Head, txn)
		for i := range uxOut {
			if contains, err := bc.Unspent().Contains(tx, uxOut[i].Hash()); err != nil {
				return err
			} else if contains {
				err := errors.New("New unspent collides with existing unspent")
				return NewErrTxnViolatesHardConstraint(err)
			}
		}
	}

	return nil
}

// VerifySingleTxnHardConstraints checks that the transaction does not violate hard constraints.
// for transactions that are not included in a block.
func (bc Blockchain) VerifySingleTxnHardConstraints(tx *bolt.Tx, txn coin.Transaction) error {
	// NOTE: Unspent().GetArray() returns an error if not all txn.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx, txn.In)
	if err != nil {
		switch err.(type) {
		case blockdb.ErrUnspentNotExist:
			return NewErrTxnViolatesHardConstraint(err)
		default:
			return err
		}
	}

	head, err := bc.Head(tx)
	if err != nil {
		return err
	}

	return bc.verifySingleTxnHardConstraints(tx, txn, head, uxIn)
}

// VerifySingleTxnAllConstraints checks that the transaction does not violate hard or soft constraints,
// for transactions that are not included in a block.
// Hard constraints are checked before soft constraints.
func (bc Blockchain) VerifySingleTxnAllConstraints(tx *bolt.Tx, txn coin.Transaction, maxSize int) error {
	// NOTE: Unspent().GetArray() returns an error if not all txn.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx, txn.In)
	if err != nil {
		return NewErrTxnViolatesHardConstraint(err)
	}

	head, err := bc.Head(tx)
	if err != nil {
		return err
	}

	// Hard constraints must be checked before soft constraints
	if err := bc.verifySingleTxnHardConstraints(tx, txn, head, uxIn); err != nil {
		return err
	}

	return VerifySingleTxnSoftConstraints(txn, head.Time(), uxIn, maxSize)
}

func (bc Blockchain) verifySingleTxnHardConstraints(tx *bolt.Tx, txn coin.Transaction, head *coin.SignedBlock, uxIn coin.UxArray) error {
	if err := VerifySingleTxnHardConstraints(txn, head, uxIn); err != nil {
		return err
	}

	if DebugLevel1 {
		// Check that new unspents don't collide with existing.
		// This should not occur but is a sanity check.
		// NOTE: this is not in the top-level VerifySingleTxnHardConstraints
		// because it relies on the unspent pool to check for existence.
		// For remote callers such as the CLI, they'd need to download the whole
		// unspent pool or make a separate API call to check for duplicate unspents.
		uxOut := coin.CreateUnspents(head.Head, txn)
		for i := range uxOut {
			if contains, err := bc.Unspent().Contains(tx, uxOut[i].Hash()); err != nil {
				return err
			} else if contains {
				err := errors.New("New unspent collides with existing unspent")
				return NewErrTxnViolatesHardConstraint(err)
			}
		}
	}

	return nil
}

// GetBlocks return blocks whose seq are in the range of start and end.
func (bc Blockchain) GetBlocks(tx *bolt.Tx, start, end uint64) ([]coin.SignedBlock, error) {
	if start > end {
		return nil, nil
	}

	var blocks []coin.SignedBlock
	for i := start; i <= end; i++ {
		b, err := bc.store.GetBlockBySeq(tx, i)
		if err != nil {
			logger.WithError(err).Error("bc.store.GetBlockBySeq failed")
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
func (bc Blockchain) GetLastBlocks(tx *bolt.Tx, num uint64) ([]coin.SignedBlock, error) {
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

	return bc.GetBlocks(tx, uint64(start), end)
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
func (bc Blockchain) processTransactions(tx *bolt.Tx, txs coin.Transactions) (coin.Transactions, error) {
	// copy txs so that the following code won't modify the original txns
	txns := make(coin.Transactions, len(txs))
	copy(txns, txs)

	// Transactions need to be sorted by fee and hash before arbitrating
	if bc.arbitrating {
		txns = coin.SortTransactions(txns, bc.TransactionFee(tx, head.Time()))
	}

	//TODO: audit
	if len(txns) == 0 {
		if bc.arbitrating {
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
			case ErrTxnViolatesHardConstraint, ErrTxnViolatesSoftConstraint:
				if bc.arbitrating {
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
				if bc.arbitrating {
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
					if bc.arbitrating {
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
					return nil, errors.New("Duplicate transaction")
				}
			}
			for a := range s.In {
				for b := range t.In {
					if s.In[a] == t.In[b] {
						if bc.arbitrating {
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
func (bc Blockchain) TransactionFee(tx *bolt.Tx, headTime uint64) coin.FeeCalculator {
	return func(t *coin.Transaction) (uint64, error) {
		inUxs, err := bc.Unspent().GetArray(tx, t.In)
		if err != nil {
			return 0, err
		}

		return TransactionFee(t, headTime, inUxs)
	}
}

type sigHash struct {
	sig  cipher.Sig
	hash cipher.SHA256
}

// verifySigs checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
func (bc *Blockchain) verifySigs(tx *bolt.Tx, workers int) error {
	if length, err := bc.Len(tx); err != nil {
		return err
	} else if length == 0 {
		return nil
	}

	sigHashes := make(chan sigHash, 100)
	errC := make(chan error, 100)
	stop := make(chan struct{})
	done := make(chan struct{})

	// Verify block signatures in a worker pool
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
		loop:
			for {
				select {
				case sh := <-sigHashes:
					if err := cipher.VerifySignature(bc.pubkey, sh.sig, sh.hash); err != nil {
						logger.Error("Signature verification failed: %v", err)
						select {
						case errC <- err:
						default:
						}
					}
				case <-stop:
					break loop
				}
			}
		}()
	}

	// Iterate all blocks stored in the "blocks" bucket
	// * Detect if a corresponding signature is missing from the signatures bucket
	// * Verify the signature for the block
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(done)

		errStopped := errors.New("goroutine was stopped")

		if err := bc.store.ForEachBlock(tx, func(block *coin.Block) error {
			sig, ok, err := bc.store.GetBlockSignature(tx, block)
			if err != nil {
				return err
			}
			if !ok {
				return blockdb.NewErrSignatureLost(block)
			}

			select {
			case sigHashes <- sigHash{
				sig:  sig,
				hash: block.HashHeader(),
			}:
				return nil
			case <-stop:
				return errStopped
			}
		}); err != nil && err != errStopped {
			logger.Error("bc.store.ForEachBlock failed: %v", err)
			select {
			case errC <- err:
			default:
			}
		}
	}()

	var foundErr error
loop:
	for {
		select {
		case err := <-errC:
			if err != nil && foundErr == nil {
				foundErr = err
				break loop
			}
		case <-done:
			break loop
		}
	}

	close(stop)
	wg.Wait()

	return foundErr
}

// VerifyBlockHeader Returns error if the BlockHeader is not valid
func (bc Blockchain) verifyBlockHeader(tx *bolt.Tx, b coin.Block) error {
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
	if b.HashBody() != b.Head.BodyHash {
		return errors.New("Computed body hash does not match")
	}
	return nil
}

// BindListener register the listener to blockchain, when new block appended, the listener will be invoked.
func (bc *Blockchain) BindListener(ls BlockListener) {
	bc.blkListener = append(bc.blkListener, ls)
}

// Notify notifies the listener the new block.
func (bc *Blockchain) Notify(b coin.Block) {
	for _, l := range bc.blkListener {
		l(b)
	}
}
