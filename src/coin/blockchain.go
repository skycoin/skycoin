package coin

import (
	"errors"
	"fmt"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"gopkg.in/op/go-logging.v1"

	//"time"
	"bytes"
)

var (
	logger      = logging.MustGetLogger("skycoin.coin")
	DebugLevel1 = true //checks for extremely unlikely conditions (10e-40)
	DebugLevel2 = true //enable checks for impossible conditions
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

type Block struct {
	Head BlockHeader
	Body BlockBody
}

type BlockHeader struct {
	Version uint32

	Time  uint64
	BkSeq uint64 //increment every block
	Fee   uint64 //fee in block, used for Proof of Stake

	PrevHash cipher.SHA256 //hash of header of previous block
	BodyHash cipher.SHA256 //hash of transaction block

	UxHash cipher.SHA256 //XOR of sha256 of elements in unspent output set

}

type BlockBody struct {
	Transactions Transactions
}

//TODO: merge header/body and cleanup top level interface

/*
Todo: merge header/body

type Block struct {
    Time  uint64
    BkSeq uint64 //increment every block
    Fee   uint64 //fee in block, used for Proof of Stake

    HashPrevBlock cipher.SHA256 //hash of header of previous block
    BodyHash      cipher.SHA256 //hash of transaction block

    Transactions Transactions
}
*/

func newBlock(prev Block, currentTime uint64, unspent UnspentPool,
	txns Transactions, calc FeeCalculator) Block {
	if len(txns) == 0 {
		log.Panic("Refusing to create block with no transactions")
	}
	fee, err := txns.Fees(calc)
	if err != nil {
		// This should have been caught earlier
		log.Panicf("Invalid transaction fees: %v", err)
	}
	body := BlockBody{txns}
	return Block{
		Head: newBlockHeader(prev.Head, unspent, currentTime, fee, body),
		Body: body,
	}
}

func (self *Block) HashHeader() cipher.SHA256 {
	return self.Head.Hash()
}

func (self *Block) HashBody() cipher.SHA256 {
	return self.Body.Hash()
}

// Returns the size of the Block's Transactions, in bytes
func (self *Block) Size() int {
	return self.Body.Size()
}

func (self *Block) String() string {
	return self.Head.String()
}

// Looks up a Transaction by its Head.Hash.
// Returns the Transaction and whether it was found or not
// TODO -- build a private index on the block, or a global blockchain one
// mapping txns to their block + tx index
// TODO: Deprecate? Utility Function
func (self *Block) GetTransaction(txHash cipher.SHA256) (Transaction, bool) {
	txns := self.Body.Transactions
	for i, _ := range txns {
		if txns[i].Hash() == txHash {
			return txns[i], true
		}
	}
	return Transaction{}, false
}

func newBlockHeader(prev BlockHeader, unspent UnspentPool, currentTime,
	fee uint64, body BlockBody) BlockHeader {
	if currentTime <= prev.Time {
		log.Panic("Time can only move forward")
	}
	prevHash := prev.Hash()
	return BlockHeader{
		BodyHash: body.Hash(),
		Version:  prev.Version,
		PrevHash: prevHash,
		Time:     currentTime,
		BkSeq:    prev.BkSeq + 1,
		Fee:      fee,
		UxHash:   getUxHash(unspent),
	}
}

func (self *BlockHeader) Hash() cipher.SHA256 {
	b1 := encoder.Serialize(*self)
	return cipher.SumSHA256(b1)
}

func (self *BlockHeader) Bytes() []byte {
	return encoder.Serialize(*self)
}

func (self *BlockHeader) String() string {
	return fmt.Sprintf("Version: %d\nTime: %d\nBkSeq: %d\nFee: %d\n"+
		"PrevHash: %s\nBodyHash: %s", self.Version, self.Time, self.BkSeq,
		self.Fee, self.PrevHash.Hex(), self.BodyHash.Hex())
}

// Returns the merkle hash of contained transactions
func (self *BlockBody) Hash() cipher.SHA256 {
	hashes := make([]cipher.SHA256, len(self.Transactions))
	for i, _ := range self.Transactions {
		hashes[i] = self.Transactions[i].Hash()
	}
	// Merkle hash of transactions
	return cipher.Merkle(hashes)
}

// Returns the size of Transactions, in bytes
func (self *BlockBody) Size() int {
	// We can't use length of self.Bytes() because it has a length prefix
	// Need only the sum of transaction sizes
	return self.Transactions.Size()
}

func (self *BlockBody) Bytes() []byte {
	return encoder.Serialize(*self)
}

type Blockchain struct {
	Blocks  []Block
	Unspent UnspentPool
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks:  make([]Block, 0),
		Unspent: NewUnspentPool(),
	}
}

// Creates a genesis block and applies it against chain
// Takes in time as parameter
func (self *Blockchain) CreateGenesisBlock(genesisAddress cipher.Address,
	timestamp uint64, genesisCoins uint64) Block {
	logger.Info("Creating new genesis block with address %s",
		genesisAddress.String())
	if len(self.Blocks) > 0 {
		log.Panic("Genesis block already created")
	}
	// Why is there a transaction in the genesis block?
	// Limits the special casing around genesis blocks:
	//      -Allows assumption that all UxOuts have a SrcTransaction that
	//      -can be found in the blockchain, without checking for genesis UxOut
	txn := Transaction{}
	txn.PushOutput(genesisAddress, genesisCoins, genesisCoins)
	body := BlockBody{Transactions{txn}}
	prevHash := cipher.SHA256{}
	head := BlockHeader{
		Time:     timestamp,
		BodyHash: body.Hash(),
		PrevHash: prevHash,
		BkSeq:    0,
		Version:  0,
		Fee:      0,
		UxHash:   getUxHash(self.Unspent),
	}
	b := Block{
		Head: head,
		Body: body,
	}
	self.Blocks = append(self.Blocks, b)
	// Genesis output
	ux := UxOut{
		Head: UxHead{
			Time:  b.Head.Time,
			BkSeq: 0,
		},
		Body: UxBody{
			SrcTransaction: txn.InnerHash, //user inner hash
			Address:        genesisAddress,
			Coins:          genesisCoins,
			Hours:          genesisCoins, // Allocate 1 coin hour per coin
		},
	}
	self.Unspent.Add(ux)
	return b
}

// Returns the most recent confirmed block
func (self *Blockchain) Head() Block {
	return self.Blocks[len(self.Blocks)-1]
}

// Time returns time of last block
// used as system clock indepedent clock for coin hour calculations
// TODO: Deprecate
func (self *Blockchain) Time() uint64 {
	return self.Head().Head.Time
}

// Creates a Block given an array of Transactions.  It does not verify the
// block; ExecuteBlock will handle verification.  Transactions must be sorted.
func (self *Blockchain) NewBlockFromTransactions(txns Transactions,
	currentTime uint64) (Block, error) {
	if currentTime <= self.Time() {
		log.Panic("Time can only move forward")
	}
	if len(txns) == 0 {
		return Block{}, errors.New("No transactions")
	}
	err := self.verifyTransactions(txns)
	if err != nil {
		return Block{}, err
	}
	b := newBlock(self.Head(), currentTime, self.Unspent, txns,
		self.TransactionFee)
	//make sure block is valid
	if DebugLevel2 == true {
		if err := verifyBlockHeader(self.Head(), b); err != nil {
			log.Panic("Impossible Error: not allowed to fail")
		}
		if err := self.verifyTransactions(b.Body.Transactions); err != nil {
			log.Panic("Impossible Error: not allowed to fail")
		}
	}
	return b, nil
}

// Attempts to append block to blockchain.  Returns the UxOuts created,
// and an error if the block is invalid.
func (self *Blockchain) ExecuteBlock(b Block) (UxArray, error) {
	var uxs UxArray = nil
	err := self.VerifyBlock(b)
	if err != nil {
		return uxs, err
	}
	txns := b.Body.Transactions
	for _, tx := range txns {
		// Remove spent outputs
		self.Unspent.DelMultiple(tx.In)
		// Create new outputs
		txUxs := CreateUnspents(b.Head, tx)
		for i, _ := range txUxs {
			self.Unspent.Add(txUxs[i])
		}
		uxs = append(uxs, txUxs...)
	}

	self.Blocks = append(self.Blocks, b)

	return uxs, nil
}

// Verifies the BlockHeader and BlockBody
func (self *Blockchain) VerifyBlock(b Block) error {
	if err := verifyBlockHeader(self.Head(), b); err != nil {
		return err
	}
	if err := self.verifyTransactions(b.Body.Transactions); err != nil {
		return err
	}
	if err := self.verifyUxHash(b); err != nil {
		return err
	}
	return nil
}

// Compares the state of the current UxHash hash to state of unspent
// output pool.
func (self *Blockchain) verifyUxHash(b Block) error {
	if !bytes.Equal(b.Head.UxHash[:], self.Unspent.XorHash[:]) {
		return errors.New("UxHash does not match")
	}
	return nil
}

// Checks that the inputs to the transaction exist,
// that the transaction does not create or destroy coins and that the
// signatures on the transaction are valid
func (self *Blockchain) VerifyTransaction(tx Transaction) error {
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

	uxIn, err := self.Unspent.GetMultiple(tx.In)
	if err != nil {
		return err
	}
	// Checks whether ux inputs exist,
	// Check that signatures are allowed to spend inputs
	if err := verifyTransactionInputs(tx, uxIn); err != nil {
		return err
	}

	// Get the UxOuts we expect to have when the block is created.
	uxOut := CreateUnspents(self.Head().Head, tx)
	// Check that there are any duplicates within this set
	if uxOut.HasDupes() {
		return errors.New("Duplicate unspent outputs in transaction")
	}
	if DebugLevel1 {
		// Check that new unspents don't collide with existing.  This should
		// also be checked in verifyTransactions
		for i, _ := range uxOut {
			if self.Unspent.Has(uxOut[i].Hash()) {
				return errors.New("New unspent collides with existing unspent")
			}
		}
	}

	// Check that no coins are lost, and sufficient coins and hours are spent
	err = verifyTransactionSpending(self.Time(), tx, uxIn, uxOut)
	if err != nil {
		return err
	}
	return nil
}

// Creates the expected outputs for a transaction.
func CreateUnspents(bh BlockHeader, tx Transaction) UxArray {
	h := tx.Hash()
	uxo := make(UxArray, len(tx.Out))
	for i, _ := range tx.Out {
		uxo[i] = UxOut{
			Head: UxHead{
				Time:  bh.Time,
				BkSeq: bh.BkSeq,
			},
			Body: UxBody{
				SrcTransaction: h,
				Address:        tx.Out[i].Address,
				Coins:          tx.Out[i].Coins,
				Hours:          tx.Out[i].Hours,
			},
		}
	}
	return uxo
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
func (self *Blockchain) processTransactions(txns Transactions,
	arbitrating bool) (Transactions, error) {
	// Transactions need to be sorted by fee and hash before arbitrating
	if arbitrating {
		txns = SortTransactions(txns, self.TransactionFee)
	}
	//TODO: audit
	if len(txns) == 0 {
		if arbitrating {
			return txns, nil
		} else {
			// If there are no transactions, a block should not be made
			return nil, errors.New("No transactions")
		}
	}

	skip := make(map[int]byte)
	uxHashes := make(UxHashSet, len(txns))
	for i, tx := range txns {
		// Check the transaction against itself.  This covers the hash,
		// signature indices and duplicate spends within itself
		err := self.VerifyTransaction(tx)
		if err != nil {
			if arbitrating {
				skip[i] = byte(1)
				continue
			} else {
				return nil, err
			}
		}
		// Check that each pending unspent will be unique
		uxb := UxBody{
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
				if self.Unspent.Has(h) {
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
		newtxns := make(Transactions, len(txns)-len(skip))
		j := 0
		for i, _ := range txns {
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
			for a, _ := range s.In {
				for b, _ := range t.In {
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
		newtxns := make(Transactions, len(txns)-len(skip))
		j := 0
		for i, _ := range txns {
			if _, shouldSkip := skip[i]; !shouldSkip {
				newtxns[j] = txns[i]
				j++
			}
		}
		return newtxns, nil
	} else {
		return txns, nil
	}
}

// Returns an error if any Transaction in txns is invalid
func (self *Blockchain) verifyTransactions(txns Transactions) error {
	_, err := self.processTransactions(txns, false)
	return err
}

// Returns an array of Transactions with invalid ones removed from txns.
// The Transaction hash is used to arbitrate between double spends.
// txns must be sorted by hash.
func (self *Blockchain) ArbitrateTransactions(txns Transactions) Transactions {
	newtxns, err := self.processTransactions(txns, true)
	if err != nil {
		log.Panicf("arbitrateTransactions failed unexpectedly: %v", err)
	}
	return newtxns
}

// Calculates the current transaction fee in coinhours of a Transaction
func (self *Blockchain) TransactionFee(t *Transaction) (uint64, error) {
	headTime := self.Time()
	inHours := uint64(0)
	// Compute input hours
	for i, _ := range t.In {
		in, ok := self.Unspent.Get(t.In[i])
		if !ok {
			return 0, errors.New("Unspent output does not exist")
		}
		inHours += in.CoinHours(headTime)
	}
	// Compute output hours
	outHours := uint64(0)
	for i, _ := range t.Out {
		outHours += t.Out[i].Hours
	}
	if inHours < outHours {
		return 0, errors.New("Insufficient coinhours for transaction outputs")
	}
	return inHours - outHours, nil
}

/* Unassigned operators */

// Validates the inputs to a transaction by checking signatures. Assumes txn
// has valid number of signatures for inputs.
func verifyTransactionInputs(tx Transaction, uxIn UxArray) error {
	if DebugLevel2 {
		if len(tx.In) != len(tx.Sigs) || len(tx.In) != len(uxIn) {
			log.Panic("tx.In != tx.Sigs != uxIn")
		}
		if tx.InnerHash != tx.HashInner() {
			log.Panic("Invalid Tx Header Hash")
		}
	}

	// Check signatures against unspent address
	for i, _ := range tx.In {
		hash := cipher.AddSHA256(tx.InnerHash, tx.In[i]) //use inner hash, not outer hash
		err := cipher.ChkSig(uxIn[i].Body.Address, hash, tx.Sigs[i])
		if err != nil {
			return errors.New("Signature not valid for output being spent")
		}
	}
	if DebugLevel2 {
		// Check that hashes match.
		// This would imply a bug with UnspentPool.GetMultiple
		if len(tx.In) != len(uxIn) {
			log.Panic("tx.In does not match uxIn")
		}
		for i, _ := range tx.In {
			if tx.In[i] != uxIn[i].Hash() {
				log.Panic("impossible error: Ux hash mismatch")
			}
		}
	}
	return nil
}

// Checks that coins will not be destroyed and that enough coins are hours
// are being spent for the outputs
func verifyTransactionSpending(headTime uint64, tx Transaction,
	uxIn UxArray, uxOut UxArray) error {
	coinsIn := uint64(0)
	hoursIn := uint64(0)
	for i, _ := range uxIn {
		coinsIn += uxIn[i].Body.Coins
		hoursIn += uxIn[i].CoinHours(headTime)
	}
	coinsOut := uint64(0)
	hoursOut := uint64(0)
	for i, _ := range uxOut {
		coinsOut += uxOut[i].Body.Coins
		hoursOut += uxOut[i].Body.Hours
	}
	if coinsIn < coinsOut {
		return errors.New("Insufficient coins")
	}
	if coinsIn > coinsOut {
		return errors.New("Transactions may not create or destroy coins")
	}
	if hoursIn < hoursOut {
		return errors.New("Insufficient coin hours")
	}
	return nil
}

// Returns error if the BlockHeader is not valid
func verifyBlockHeader(head Block, b Block) error {
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

// Returns unspent output checksum for the Block. Must be called after Block
// is fully initialized, and before its outputs are added to the unspent pool
func getUxHash(unspent UnspentPool) cipher.SHA256 {
	return unspent.XorHash
}
