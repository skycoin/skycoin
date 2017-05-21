package visor

import (
	"errors"
	"fmt"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"

	"bytes"
)

var (
	// DebugLevel1 checks for extremely unlikely conditions (10e-40)
	DebugLevel1 = true
	// DebugLevel2 enable checks for impossible conditions
	DebugLevel2 = true

	// ErrUnspentNotExist represents the error of unspent output in a tx does not exist
	ErrUnspentNotExist = errors.New("Unspent output does not exist")
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

// BlockTree provide method for access
type BlockTree interface {
	AddBlock(b *coin.Block) error
	RemoveBlock(b *coin.Block) error
	GetBlock(hash cipher.SHA256) *coin.Block
	GetBlockInDepth(dep uint64, filter func(hps []coin.HashPair) cipher.SHA256) *coin.Block
}

// Walker function for go through blockchain
type Walker func(hps []coin.HashPair) cipher.SHA256

// BlockListener notify the register when new block is appended to the chain
type BlockListener func(b coin.Block)

// Blockchain use blockdb to store the blocks.
type Blockchain struct {
	tree        BlockTree
	walker      Walker
	unspent     coin.UnspentPool
	head        cipher.SHA256
	blkListener []BlockListener
	arbitrating bool
}

// Option represents the option when creating the blockchain
type Option func(*Blockchain)

// NewBlockchain use the walker go through the tree and update the head and unspent outputs.
func NewBlockchain(tree BlockTree, walker Walker, ops ...Option) *Blockchain {
	bc := &Blockchain{
		tree:    tree,
		walker:  walker,
		unspent: coin.NewUnspentPool(),
	}

	for _, op := range ops {
		op(bc)
	}

	bc.walkTree()
	return bc
}

// Arbitrating option to change the mode
func Arbitrating(enable bool) Option {
	return func(bc *Blockchain) {
		bc.arbitrating = enable
	}
}

// GetUnspent returns the unspent output pool.
func (bc *Blockchain) GetUnspent() *coin.UnspentPool {
	return &bc.unspent
}

func (bc *Blockchain) walkTree() {
	var dep uint64
	var preBlock *coin.Block

	for {
		b := bc.tree.GetBlockInDepth(dep, bc.walker)
		if b == nil {
			break
		}
		if dep > 0 {
			if b.PreHashHeader() != preBlock.HashHeader() {
				logger.Panicf("walk tree failed, pre hash header not match")
			}
		}
		preBlock = b
		if err := bc.updateUnspent(*b); err != nil {
			logger.Panicf("update unspent failed, err: %v", err.Error())
		}

		bc.head = b.HashHeader()
		dep++
	}
}

// Len returns the length of current blockchain.
func (bc Blockchain) Len() uint64 {
	return bc.Head().Seq() + 1
}

// GetGenesisBlock get genesis block.
func (bc Blockchain) GetGenesisBlock() *coin.Block {
	return bc.tree.GetBlockInDepth(0, bc.walker)
}

// CreateGenesisBlock creates genesis block in blockchain.
func (bc *Blockchain) CreateGenesisBlock(genesisAddr cipher.Address, genesisCoins, timestamp uint64) coin.Block {
	txn := coin.Transaction{}
	txn.PushOutput(genesisAddr, genesisCoins, genesisCoins)
	body := coin.BlockBody{Transactions: coin.Transactions{txn}}
	prevHash := cipher.SHA256{}
	head := coin.BlockHeader{
		Time:     timestamp,
		BodyHash: body.Hash(),
		PrevHash: prevHash,
		BkSeq:    0,
		Version:  0,
		Fee:      0,
		UxHash:   coin.NewUnspentPool().GetUxHash(),
	}
	b := coin.Block{
		Head: head,
		Body: body,
	}
	// b.Body.Transactions[0].UpdateHeader()
	bc.addBlock(&b)
	bc.head = b.HashHeader()
	ux := coin.UxOut{
		Head: coin.UxHead{
			Time:  timestamp,
			BkSeq: 0,
		},
		Body: coin.UxBody{
			SrcTransaction: txn.InnerHash, //user inner hash
			Address:        genesisAddr,
			Coins:          genesisCoins,
			Hours:          genesisCoins, // Allocate 1 coin hour per coin
		},
	}
	bc.unspent.Add(ux)
	bc.notify(b)
	return b
}

func (bc *Blockchain) addBlock(b *coin.Block) error {
	return bc.tree.AddBlock(b)
}

// GetBlock get block of specific hash in the blockchain, return nil on not found.
func (bc Blockchain) GetBlock(hash cipher.SHA256) *coin.Block {
	return bc.tree.GetBlock(hash)
}

// Head returns the most recent confirmed block
func (bc Blockchain) Head() *coin.Block {
	return bc.GetBlock(bc.head)
}

// Time returns time of last block
// used as system clock indepedent clock for coin hour calculations
// TODO: Deprecate
func (bc *Blockchain) Time() uint64 {
	return bc.Head().Time()
}

// NewBlockFromTransactions creates a Block given an array of Transactions.  It does not verify the
// block; ExecuteBlock will handle verification.  Transactions must be sorted.
func (bc Blockchain) NewBlockFromTransactions(txns coin.Transactions, currentTime uint64) (coin.Block, error) {
	if currentTime <= bc.Time() {
		log.Panic("Time can only move forward")
	}
	if len(txns) == 0 {
		return coin.Block{}, errors.New("No transactions")
	}
	txns, err := bc.verifyTransactions(txns)
	if err != nil {
		return coin.Block{}, err
	}
	b := coin.NewBlock(*bc.Head(), currentTime, bc.unspent, txns, bc.TransactionFee)

	//make sure block is valid
	if DebugLevel2 == true {
		if err := bc.verifyBlockHeader(b); err != nil {
			log.Panic("Impossible Error: not allowed to fail")
		}
		txns, err := bc.verifyTransactions(b.Body.Transactions)
		if err != nil {
			log.Panic("Impossible Error: not allowed to fail")
		}
		b.Body.Transactions = txns
	}
	return b, nil
}

// ExecuteBlock Attempts to append block to blockchain.  Returns the UxOuts created,
// and an error if the block is invalid.
func (bc *Blockchain) ExecuteBlock(b *coin.Block) (coin.UxArray, error) {
	var uxs coin.UxArray
	err := bc.verifyBlock(*b)
	if err != nil {
		return nil, err
	}
	txns := b.Body.Transactions
	for _, tx := range txns {
		// Remove spent outputs
		bc.unspent.DelMultiple(tx.In)
		// Create new outputs
		txUxs := coin.CreateUnspents(b.Head, tx)
		for i := range txUxs {
			bc.unspent.Add(txUxs[i])
		}
		uxs = append(uxs, txUxs...)
	}

	b.Head.PrevHash = bc.head
	bc.addBlock(b)

	// update the head
	bc.head = b.HashHeader()
	bc.notify(*b)
	return uxs, nil
}

func (bc *Blockchain) updateUnspent(b coin.Block) error {
	if err := bc.verifyBlock(b); err != nil {
		return err
	}
	txns := b.Body.Transactions
	for _, tx := range txns {
		// Remove spent outputs
		bc.unspent.DelMultiple(tx.In)
		// Create new outputs
		txUxs := coin.CreateUnspents(b.Head, tx)
		for i := range txUxs {
			bc.unspent.Add(txUxs[i])
		}
	}
	return nil
}

// VerifyBlock verifies the BlockHeader and BlockBody
func (bc Blockchain) verifyBlock(b coin.Block) error {
	gb := bc.GetGenesisBlock()
	if gb.HashHeader() != b.HashHeader() {
		if err := bc.verifyBlockHeader(b); err != nil {
			return err
		}

		txns, err := bc.verifyTransactions(b.Body.Transactions)
		if err != nil {
			return err
		}
		b.Body.Transactions = txns
	}

	if err := bc.verifyUxHash(b); err != nil {
		return err
	}
	return nil
}

// Compares the state of the current UxHash hash to state of unspent
// output pool.
func (bc Blockchain) verifyUxHash(b coin.Block) error {
	if !bytes.Equal(b.Head.UxHash[:], bc.unspent.XorHash[:]) {
		return errors.New("UxHash does not match")
	}
	return nil
}

// VerifyTransaction checks that the inputs to the transaction exist,
// that the transaction does not create or destroy coins and that the
// signatures on the transaction are valid
func (bc Blockchain) VerifyTransaction(tx coin.Transaction) error {
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
	// Check for non 1e6 multiple coin outputs
	// Check for zero coin outputs
	// Check valid looking signatures
	if err := tx.Verify(); err != nil {
		return err
	}

	uxIn, err := bc.unspent.GetMultiple(tx.In)
	if err != nil {
		return err
	}
	// Checks whether ux inputs exist,
	// Check that signatures are allowed to spend inputs
	if err := tx.VerifyInput(uxIn); err != nil {
		return err
	}

	// Get the UxOuts we expect to have when the block is created.
	uxOut := coin.CreateUnspents(bc.Head().Head, tx)
	// Check that there are any duplicates within this set
	if uxOut.HasDupes() {
		return errors.New("Duplicate unspent outputs in transaction")
	}
	if DebugLevel1 {
		// Check that new unspents don't collide with existing.  This should
		// also be checked in verifyTransactions
		for i := range uxOut {
			if bc.unspent.Has(uxOut[i].Hash()) {
				return errors.New("New unspent collides with existing unspent")
			}
		}
	}

	// Check that no coins are lost, and sufficient coins and hours are spent
	err = coin.VerifyTransactionSpending(bc.Time(), uxIn, uxOut)
	if err != nil {
		return err
	}
	return nil
}

// GetBlockInDepth return block whose BkSeq is seq.
func (bc Blockchain) GetBlockInDepth(dep uint64) *coin.Block {
	return bc.tree.GetBlockInDepth(dep, bc.walker)
}

// GetBlocks return blocks whose seq are in the range of start and end.
func (bc Blockchain) GetBlocks(start, end uint64) []coin.Block {
	if start > end {
		return []coin.Block{}
	}

	blocks := []coin.Block{}
	for i := start; i <= end; i++ {
		b := bc.tree.GetBlockInDepth(i, bc.walker)
		if b == nil {
			break
		}
		blocks = append(blocks, *b)
	}
	return blocks
}

// GetLastBlocks return the latest N blocks.
func (bc Blockchain) GetLastBlocks(num uint64) []coin.Block {
	var blocks []coin.Block
	if num == 0 {
		return blocks
	}

	end := bc.Head().Seq()
	start := end - num + 1
	if start < 0 {
		start = 0
	}
	return bc.GetBlocks(start, end)
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
func (bc Blockchain) processTransactions(txns coin.Transactions, arbitrating bool) (coin.Transactions, error) {
	// Transactions need to be sorted by fee and hash before arbitrating
	if arbitrating {
		txns = coin.SortTransactions(txns, bc.TransactionFee)
	}
	//TODO: audit
	if len(txns) == 0 {
		if arbitrating {
			return txns, nil
		}
		// If there are no transactions, a block should not be made
		return nil, errors.New("No transactions")
	}

	skip := make(map[int]byte)
	uxHashes := make(coin.UxHashSet, len(txns))
	for i, tx := range txns {
		// Check the transaction against itself.  This covers the hash,
		// signature indices and duplicate spends within itself
		err := bc.VerifyTransaction(tx)
		if err != nil {
			if arbitrating {
				skip[i] = byte(1)
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
				if arbitrating {
					skip[i] = byte(1)
					continue
				} else {
					m := "Duplicate unspent output across transactions"
					return nil, errors.New(m)
				}
			}
			if DebugLevel1 {
				// Check that the expected unspent is not already in the pool.
				// This should never happen because its a hash collision
				if bc.unspent.Has(h) {
					if arbitrating {
						skip[i] = byte(1)
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
		skip = make(map[int]byte)
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
						if arbitrating {
							// The txn with the highest fee and lowest hash
							// is chosen when attempting a double spend.
							// Since the txns are sorted, we skip the 2nd
							// iterable
							skip[j] = byte(1)
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
		newtxns := make(coin.Transactions, len(txns)-len(skip))
		j := 0
		for i := range txns {
			if _, shouldSkip := skip[i]; !shouldSkip {
				newtxns[j] = txns[i]
				j++
			}
		}
		return newtxns, nil
	}

	return txns, nil
}

// verifyTransactions returns an error if any Transaction in txns is invalid
func (bc Blockchain) verifyTransactions(txns coin.Transactions) (coin.Transactions, error) {
	return bc.processTransactions(txns, bc.arbitrating)
}

// ArbitrateTransactions returns an array of Transactions with invalid ones removed from txns.
// The Transaction hash is used to arbitrate between double spends.
// txns must be sorted by hash.
func (bc Blockchain) ArbitrateTransactions(txns coin.Transactions) coin.Transactions {
	newtxns, err := bc.processTransactions(txns, true)
	if err != nil {
		log.Panicf("arbitrateTransactions failed unexpectedly: %v", err)
	}
	return newtxns
}

// TransactionFee calculates the current transaction fee in coinhours of a Transaction
func (bc Blockchain) TransactionFee(t *coin.Transaction) (uint64, error) {
	headTime := bc.Time()
	inHours := uint64(0)
	// Compute input hours
	for i := range t.In {
		in, ok := bc.unspent.Get(t.In[i])
		if !ok {
			return 0, ErrUnspentNotExist
		}
		inHours += in.CoinHours(headTime)
	}
	// Compute output hours
	outHours := uint64(0)
	for i := range t.Out {
		outHours += t.Out[i].Hours
	}
	if inHours < outHours {
		return 0, errors.New("Insufficient coinhours for transaction outputs")
	}
	return inHours - outHours, nil
}

// VerifySigs checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
func (bc Blockchain) VerifySigs(pubKey cipher.PubKey, sigs *blockdb.BlockSigs) error {
	for i := uint64(0); i <= bc.Head().Seq(); i++ {
		b := bc.GetBlockInDepth(i)
		if b == nil {
			return fmt.Errorf("no block in depth %v", i)
		}
		// get sig
		sig, err := sigs.Get(b.HashHeader())
		if err != nil {
			return err
		}

		if err := cipher.VerifySignature(pubKey, sig, b.HashHeader()); err != nil {
			return err
		}
	}

	return nil
}

// VerifyBlockHeader Returns error if the BlockHeader is not valid
func (bc Blockchain) verifyBlockHeader(b coin.Block) error {
	//check BkSeq
	head := bc.Head()
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

// notifies the listener the new block.
func (bc *Blockchain) notify(b coin.Block) {
	for _, l := range bc.blkListener {
		l(b)
	}
}
