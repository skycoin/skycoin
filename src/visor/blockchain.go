package visor

import (
	"bytes"
	"errors"
	"sync"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor/blockdb"
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
	Head() (*coin.SignedBlock, error) // returns head block
	HeadSeq() uint64                  // returns head block sequence
	Len() uint64                      // returns blockchain lenght
	AddBlockWithTx(tx *bolt.Tx, b *coin.SignedBlock) error
	GetBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error)
	GetBlockBySeq(seq uint64) (*coin.SignedBlock, error)
	UnspentPool() blockdb.UnspentPool
	GetGenesisBlock() *coin.SignedBlock
}

// BlockListener notify the register when new block is appended to the chain
type BlockListener func(b coin.Block)

// Blockchain maintains blockchain and provides apis for accessing the chain.
type Blockchain struct {
	db          *bolt.DB
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
func DefaultWalker(hps []coin.HashPair) cipher.SHA256 {
	return hps[0].Hash
}

// NewBlockchain use the walker go through the tree and update the head and unspent outputs.
func NewBlockchain(db *bolt.DB, pubkey cipher.PubKey, ops ...Option) (*Blockchain, error) {
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
	if err := bc.verifySigs(); err != nil {
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
func (bc *Blockchain) GetGenesisBlock() *coin.SignedBlock {
	return bc.store.GetGenesisBlock()
}

// GetBlockByHash returns block of given hash
func (bc *Blockchain) GetBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error) {
	return bc.store.GetBlockByHash(hash)
}

// GetBlockBySeq returns block of given seq
func (bc *Blockchain) GetBlockBySeq(seq uint64) (*coin.SignedBlock, error) {
	return bc.store.GetBlockBySeq(seq)
}

func (bc *Blockchain) processBlock(b coin.SignedBlock) (coin.SignedBlock, error) {
	if bc.Len() > 0 {
		if !bc.isGenesisBlock(b.Block) {
			if err := bc.verifyBlockHeader(b.Block); err != nil {
				return coin.SignedBlock{}, err
			}

			txns, err := bc.processTransactions(b.Body.Transactions)
			if err != nil {
				return coin.SignedBlock{}, err
			}
			b.Body.Transactions = txns

			if err := bc.verifyUxHash(b.Block); err != nil {
				return coin.SignedBlock{}, err
			}

		}
	}

	return b, nil
}

// Unspent returns the unspent outputs pool
func (bc *Blockchain) Unspent() blockdb.UnspentPool {
	return bc.store.UnspentPool()
}

// Len returns the length of current blockchain.
func (bc Blockchain) Len() uint64 {
	return bc.store.Len()
}

// Head returns the most recent confirmed block
func (bc Blockchain) Head() (*coin.SignedBlock, error) {
	return bc.store.Head()
}

// HeadSeq returns the sequence of head block
func (bc *Blockchain) HeadSeq() uint64 {
	return bc.store.HeadSeq()
}

// Time returns time of last block
// used as system clock indepedent clock for coin hour calculations
// TODO: Deprecate
func (bc *Blockchain) Time() uint64 {
	b, err := bc.Head()
	if err != nil {
		return 0
	}

	return b.Time()
}

// NewBlock creates a Block given an array of Transactions.
// Only hard constraints are applied to transactions in the block.
// The caller of this function should apply any additional soft constraints,
// and choose which transactions to place into the block.
func (bc Blockchain) NewBlock(txns coin.Transactions, currentTime uint64) (*coin.Block, error) {
	if currentTime <= bc.Time() {
		return nil, errors.New("Time can only move forward")
	}

	if len(txns) == 0 {
		return nil, errors.New("No transactions")
	}
	txns, err := bc.processTransactions(txns)
	if err != nil {
		return nil, err
	}
	uxHash := bc.Unspent().GetUxHash()

	head, err := bc.Head()
	if err != nil {
		return nil, err
	}

	b, err := coin.NewBlock(head.Block, currentTime, uxHash, txns, bc.TransactionFee)
	if err != nil {
		return nil, err
	}

	//make sure block is valid
	if DebugLevel2 == true {
		if err := bc.verifyBlockHeader(*b); err != nil {
			return nil, err
		}
		txns, err := bc.processTransactions(b.Body.Transactions)
		if err != nil {
			logger.Panic("Impossible Error: not allowed to fail")
		}
		b.Body.Transactions = txns
	}
	return b, nil
}

// ExecuteBlockWithTx attempts to append block to blockchain with *bolt.Tx
func (bc *Blockchain) ExecuteBlockWithTx(tx *bolt.Tx, sb *coin.SignedBlock) error {
	if bc.Len() > 0 {
		head, err := bc.Head()
		if err != nil {
			return err
		}

		sb.Head.PrevHash = head.HashHeader()
	}

	nb, err := bc.processBlock(*sb)
	if err != nil {
		return err
	}

	if err := bc.store.AddBlockWithTx(tx, &nb); err != nil {
		return err
	}

	return nil
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

// VerifyBlockTxnConstraints checks that the transaction does not violate hard constraints,
// for transactions that are already included in a block.
func (bc Blockchain) VerifyBlockTxnConstraints(tx coin.Transaction) error {
	// NOTE: Unspent().GetArray() returns an error if not all tx.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx.In)
	if err != nil {
		switch err.(type) {
		case blockdb.ErrUnspentNotExist:
			return NewErrTxnViolatesHardConstraint(err)
		default:
			return err
		}
	}

	head, err := bc.Head()
	if err != nil {
		return err
	}

	return bc.verifyBlockTxnHardConstraints(tx, head, uxIn)
}

func (bc Blockchain) verifyBlockTxnHardConstraints(tx coin.Transaction, head *coin.SignedBlock, uxIn coin.UxArray) error {
	if err := VerifyBlockTxnConstraints(tx, head, uxIn); err != nil {
		return err
	}

	if DebugLevel1 {
		// Check that new unspents don't collide with existing.
		// This should not occur but is a sanity check.
		// NOTE: this is not in the top-level VerifyBlockTxnConstraints
		// because it relies on the unspent pool to check for existence.
		// For remote callers such as the CLI, they'd need to download the whole
		// unspent pool or make a separate API call to check for duplicate unspents.
		uxOut := coin.CreateUnspents(head.Head, tx)
		for i := range uxOut {
			if bc.Unspent().Contains(uxOut[i].Hash()) {
				err := errors.New("New unspent collides with existing unspent")
				return NewErrTxnViolatesHardConstraint(err)
			}
		}
	}

	return nil
}

// VerifySingleTxnHardConstraints checks that the transaction does not violate hard constraints.
// for transactions that are not included in a block.
func (bc Blockchain) VerifySingleTxnHardConstraints(tx coin.Transaction) error {
	// NOTE: Unspent().GetArray() returns an error if not all tx.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx.In)
	if err != nil {
		switch err.(type) {
		case blockdb.ErrUnspentNotExist:
			return NewErrTxnViolatesHardConstraint(err)
		default:
			return err
		}
	}

	head, err := bc.Head()
	if err != nil {
		return err
	}

	return bc.verifySingleTxnHardConstraints(tx, head, uxIn)
}

// VerifySingleTxnAllConstraints checks that the transaction does not violate hard or soft constraints,
// for transactions that are not included in a block.
// Hard constraints are checked before soft constraints.
func (bc Blockchain) VerifySingleTxnAllConstraints(tx coin.Transaction, maxSize int) error {
	// NOTE: Unspent().GetArray() returns an error if not all tx.In can be found
	// This prevents double spends
	uxIn, err := bc.Unspent().GetArray(tx.In)
	if err != nil {
		return NewErrTxnViolatesHardConstraint(err)
	}

	head, err := bc.Head()
	if err != nil {
		return err
	}

	// Hard constraints must be checked before soft constraints
	if err := bc.verifySingleTxnHardConstraints(tx, head, uxIn); err != nil {
		return err
	}

	return VerifySingleTxnSoftConstraints(tx, head.Time(), uxIn, maxSize)
}

func (bc Blockchain) verifySingleTxnHardConstraints(tx coin.Transaction, head *coin.SignedBlock, uxIn coin.UxArray) error {
	if err := VerifySingleTxnHardConstraints(tx, head, uxIn); err != nil {
		return err
	}

	if DebugLevel1 {
		// Check that new unspents don't collide with existing.
		// This should not occur but is a sanity check.
		// NOTE: this is not in the top-level VerifySingleTxnHardConstraints
		// because it relies on the unspent pool to check for existence.
		// For remote callers such as the CLI, they'd need to download the whole
		// unspent pool or make a separate API call to check for duplicate unspents.
		uxOut := coin.CreateUnspents(head.Head, tx)
		for i := range uxOut {
			if bc.Unspent().Contains(uxOut[i].Hash()) {
				err := errors.New("New unspent collides with existing unspent")
				return NewErrTxnViolatesHardConstraint(err)
			}
		}
	}

	return nil
}

// GetBlocks return blocks whose seq are in the range of start and end.
func (bc Blockchain) GetBlocks(start, end uint64) ([]coin.SignedBlock, error) {
	if start > end {
		return nil, nil
	}

	var blocks []coin.SignedBlock
	for i := start; i <= end; i++ {
		b, err := bc.store.GetBlockBySeq(i)
		if err != nil {
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
func (bc Blockchain) GetLastBlocks(num uint64) ([]coin.SignedBlock, error) {
	if num == 0 {
		return nil, nil
	}

	end := bc.HeadSeq()
	start := int(end-num) + 1
	if start < 0 {
		start = 0
	}
	return bc.GetBlocks(uint64(start), end)
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
func (bc Blockchain) processTransactions(txs coin.Transactions) (coin.Transactions, error) {
	// copy txs so that the following code won't modify the origianl txs
	txns := make(coin.Transactions, len(txs))
	copy(txns, txs)

	// Transactions need to be sorted by fee and hash before arbitrating
	if bc.arbitrating {
		txns = coin.SortTransactions(txns, bc.TransactionFee)
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
		err := bc.VerifyBlockTxnConstraints(tx)
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
func (bc Blockchain) TransactionFee(t *coin.Transaction) (uint64, error) {
	headTime := bc.Time()
	inUxs, err := bc.Unspent().GetArray(t.In)
	if err != nil {
		return 0, err
	}

	return fee.TransactionFee(t, headTime, inUxs)
}

// verifySigs checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
func (bc *Blockchain) verifySigs() error {
	if bc.Len() == 0 {
		return nil
	}

	head, err := bc.Head()
	if err != nil {
		return err
	}

	seqC := make(chan uint64)

	shutdown, errC := bc.sigVerifier(seqC)

	for i := uint64(0); i <= head.Seq(); i++ {
		seqC <- i
	}

	shutdown()

	return <-errC
}

// signature verifier will get block seq from seqC channel,
// and have multiple thread to do signature verification.
func (bc *Blockchain) sigVerifier(seqC chan uint64) (func(), <-chan error) {
	quitC := make(chan struct{})
	wg := sync.WaitGroup{}
	errC := make(chan error, 1)
	for i := 0; i < SigVerifyTheadNum; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case seq := <-seqC:
					if err := bc.verifyBlockSig(seq); err != nil {
						errC <- err
						return
					}
				case <-quitC:
					return
				}
			}
		}(i)
	}

	return func() {
		close(quitC)
		wg.Wait()
		select {
		case errC <- nil:
			// no error
		default:
			// already has error in errC
		}
	}, errC
}

func (bc *Blockchain) verifyBlockSig(seq uint64) error {
	sb, err := bc.store.GetBlockBySeq(seq)
	if err != nil {
		return err
	}

	return sb.VerifySignature(bc.pubkey)
}

// VerifyBlockHeader Returns error if the BlockHeader is not valid
func (bc Blockchain) verifyBlockHeader(b coin.Block) error {
	//check BkSeq
	head, err := bc.Head()
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
	for _, l := range bc.blkListener {
		l(b)
	}
}

// UpdateDB updates db with given func
func (bc *Blockchain) UpdateDB(f func(t *bolt.Tx) error) error {
	return bc.db.Update(f)
}
