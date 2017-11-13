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

var (
	// DebugLevel1 checks for extremely unlikely conditions (10e-40)
	DebugLevel1 = true
	// DebugLevel2 enable checks for impossible conditions
	DebugLevel2 = true

	// ErrUnspentNotExist represents the error of unspent output in a tx does not exist
	ErrUnspentNotExist = errors.New("Unspent output does not exist")
)

const (
	// sigVerifyTheadNum is the number of goroutines to verify block signatures in
	sigVerifyTheadNum = 4
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
	Head(*bolt.Tx) (*coin.SignedBlock, error) // returns head block
	HeadSeq() uint64                          // returns head block sequence
	Len() uint64                              // returns blockchain length
	AddBlock(*bolt.Tx, *coin.SignedBlock) error
	GetBlockByHash(*bolt.Tx, cipher.SHA256) (*coin.Block, error)
	GetSignedBlockByHash(*bolt.Tx, cipher.SHA256) (*coin.SignedBlock, error)
	GetSignedBlockBySeq(*bolt.Tx, uint64) (*coin.SignedBlock, error)
	UnspentPool() blockdb.UnspentPool
	GetGenesisBlock() *coin.SignedBlock
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

// BlockchainOptions represents the option when creating the blockchain
type BlockchainOptions struct {
	Arbitrating                  bool
	SignatureVerificationWorkers int
}

// DefaultWalker default blockchain walker
func DefaultWalker(tx *bolt.Tx, hps []coin.HashPair) (cipher.SHA256, bool) {
	if len(hps) == 0 {
		return cipher.SHA256{}, false
	}
	return hps[0].Hash, true
}

// NewBlockchain use the walker go through the tree and update the head and unspent outputs.
func NewBlockchain(db *dbutil.DB, pubkey cipher.PubKey, opts BlockchainOptions) (*Blockchain, error) {
	chainstore, err := blockdb.NewBlockchain(db, DefaultWalker)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		db:          db,
		pubkey:      pubkey,
		store:       chainstore,
		arbitrating: opts.Arbitrating,
	}

	// verify signature
	if err := db.View(func(tx *bolt.Tx) error {
		workers := opts.SignatureVerificationWorkers
		if workers == 0 {
			workers = sigVerifyTheadNum
		}

		return bc.verifySigs(tx, workers)
	}); err != nil {
		return nil, err
	}

	return bc, nil
}

// GetGenesisBlock returns genesis block
func (bc *Blockchain) GetGenesisBlock() *coin.SignedBlock {
	return bc.store.GetGenesisBlock()
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
func (bc Blockchain) Len() uint64 {
	return bc.store.Len()
}

// HeadSeq returns the sequence of head block
func (bc *Blockchain) HeadSeq() uint64 {
	return bc.store.HeadSeq()
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

// NewBlock creates a Block given an array of Transactions.  It does not verify the
// block; ExecuteBlock will handle verification.  Transactions must be sorted.
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

	txns, err = bc.processTransactions(head, txns)
	if err != nil {
		return nil, err
	}

	uxHash := bc.Unspent().GetUxHash()

	b, err := coin.NewBlock(head.Block, currentTime, uxHash, txns, bc.TransactionFee(head.Time()))
	if err != nil {
		return nil, err
	}

	//make sure block is valid
	if DebugLevel2 == true {
		if err := bc.verifyBlockHeader(tx, *b); err != nil {
			return nil, err
		}
		txns, err := bc.processTransactions(head, b.Body.Transactions)
		if err != nil {
			logger.Panic("Impossible Error: not allowed to fail")
		}
		b.Body.Transactions = txns
	}

	return b, nil
}

// ExecuteBlock attempts to append block to blockchain with *bolt.Tx
func (bc *Blockchain) ExecuteBlock(tx *bolt.Tx, sb *coin.SignedBlock) error {
	var head *coin.SignedBlock
	if bc.Len() > 0 {
		var err error
		if head, err = bc.store.Head(tx); err != nil {
			return err
		}

		sb.Head.PrevHash = head.HashHeader()
	}

	nb, err := bc.processBlock(tx, head, *sb)
	if err != nil {
		return err
	}

	if err := bc.store.AddBlock(tx, &nb); err != nil {
		return err
	}

	return nil
}

func (bc *Blockchain) processBlock(tx *bolt.Tx, head *coin.SignedBlock, b coin.SignedBlock) (coin.SignedBlock, error) {
	if bc.Len() > 0 {
		if bc.isGenesisBlock(b.Block) {
			err := errors.New("Attempt to process genesis block after blockchain has genesis block")
			logger.Panic(err.Error())
			return coin.SignedBlock{}, err
		}

		if err := bc.verifyBlockHeader(tx, b.Block); err != nil {
			return coin.SignedBlock{}, err
		}

		txns, err := bc.processTransactions(head, b.Body.Transactions)
		if err != nil {
			return coin.SignedBlock{}, err
		}

		b.Body.Transactions = txns

		if err := bc.verifyUxHash(b.Block); err != nil {
			return coin.SignedBlock{}, err
		}
	}

	return b, nil
}

// isGenesisBlock checks if the block is genesis block
func (bc Blockchain) isGenesisBlock(b coin.Block) bool {
	gb := bc.store.GetGenesisBlock()
	if gb == nil {
		return false
	}

	return gb.HashHeader() == b.HashHeader()
}

// Compares the state of the current UxHash hash to state of unspent
// output pool.
func (bc Blockchain) verifyUxHash(b coin.Block) error {
	uxHash := bc.Unspent().GetUxHash()

	if !bytes.Equal(b.Head.UxHash[:], uxHash[:]) {
		return errors.New("UxHash does not match")
	}
	return nil
}

// VerifyTransaction checks that the inputs to the transaction exist,
// that the transaction does not create or destroy coins and that the
// signatures on the transaction are valid
func (bc Blockchain) VerifyTransaction(head *coin.SignedBlock, txn coin.Transaction) error {
	//CHECKLIST: DONE: check for duplicate ux inputs/double spending
	//CHECKLIST: DONE: check that inputs of transaction have not been spent
	//CHECKLIST: DONE: check there are no duplicate outputs

	// Q: why are coin hours based on last block time and not
	// current time?
	// A: no two computers will agree on system time. Need system clock
	// indepedent timing that everyone agrees on. fee values would depend on
	// local clock

	// Check transaction type and length
	// Check for duplicate outputs
	// Check for duplicate inputs
	// Check for invalid hash
	// Check for no inputs
	// Check for no outputs
	// Check for zero coin outputs
	// Check valid looking signatures
	if err := txn.Verify(); err != nil {
		return err
	}

	uxIn, err := bc.Unspent().GetArray(txn.In)
	if err != nil {
		return err
	}
	// Checks whether ux inputs exist,
	// Check that signatures are allowed to spend inputs
	if err := txn.VerifyInput(uxIn); err != nil {
		return err
	}

	uxOut := coin.CreateUnspents(head.Head, txn)
	// Check that there are any duplicates within this set
	if uxOut.HasDupes() {
		return errors.New("Duplicate unspent outputs in transaction")
	}
	if DebugLevel1 {
		// Check that new unspents don't collide with existing.  This should
		// also be checked in verifyTransactions
		for i := range uxOut {
			if bc.Unspent().Contains(uxOut[i].Hash()) {
				return errors.New("New unspent collides with existing unspent")
			}
		}
	}

	// Check that no coins are lost, and sufficient coins and hours are spent
	err = coin.VerifyTransactionSpending(head.Time(), uxIn, uxOut)
	if err != nil {
		return err
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
		b, err := bc.store.GetSignedBlockBySeq(tx, i)
		if err != nil {
			logger.Error("GetSignedBlockBySeq failed: %v", err)
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

	end := bc.HeadSeq()
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
func (bc Blockchain) processTransactions(head *coin.SignedBlock, txs coin.Transactions) (coin.Transactions, error) {
	// copy txs so that the following code won't modify the origianl txs
	txns := make(coin.Transactions, len(txs))
	copy(txns, txs)

	// Transactions need to be sorted by fee and hash before arbitrating
	if bc.arbitrating {
		txns = coin.SortTransactions(txns, bc.TransactionFee(head.Time()))
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
	for i, tx := range txns {
		// Check the transaction against itself.  This covers the hash,
		// signature indices and duplicate spends within itself
		err := bc.VerifyTransaction(head, tx)
		if err != nil {
			if bc.arbitrating {
				skip[i] = struct{}{}
				continue
			} else {
				return nil, err
			}
		}

		// Check that each pending unspent will be unique
		uxb := coin.UxBody{
			SrcTransaction: tx.Hash(),
		}
		for _, to := range tx.Out {
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
					m := "Duplicate unspent output across transactions"
					return nil, errors.New(m)
				}
			}
			if DebugLevel1 {
				// Check that the expected unspent is not already in the pool.
				// This should never happen because its a hash collision
				if bc.Unspent().Contains(h) {
					if bc.arbitrating {
						skip[i] = struct{}{}
						continue
					} else {
						m := "Output hash is in the UnspentPool"
						return nil, errors.New(m)
					}
				}
			}
			uxHashes[h] = byte(1)
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
	// within a single Transaction are already checked by VerifyTransaction
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
							m := "Cannot spend output twice in the same block"
							return nil, errors.New(m)
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
func (bc Blockchain) TransactionFee(headTime uint64) coin.FeeCalculator {
	return func(t *coin.Transaction) (uint64, error) {
		inUxs, err := bc.Unspent().GetArray(t.In)
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
	if bc.Len() == 0 {
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
			logger.Error("bc.store.ForEachSignature failed: %v", err)
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
	//check BkSeq
	head, err := bc.store.Head(tx)
	if err != nil {
		return err
	}

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
	logger.Debug("visor.Blockchain.Notify")
	defer logger.Debug("visor.Blockchain.Notify complete")

	for _, l := range bc.blkListener {
		l(b)
	}
}
