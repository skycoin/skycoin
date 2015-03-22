package coin

import (
	"bytes"
	"errors"
	"log"
	"math"
	"sort"

	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
)

/*
Transaction with N inputs, M ouputs is
- 32 bytes constant
- 32+65 bytes per input
- 21+8+8 bytes per output

Skycoin Transactions are
- 97 bytes per input +  37 bytes per output + 37 bytes
Bitcoin Transactions are
- 180 bytes per input + 34 bytes per output + 10 bytes

Sigs is the array of signatures
- the Nth signature is the authorization to spend the Nth output consumed in transaction
- the hash signed is SHA256sum of transaction inner hash and the hash of output being spent

The inner hash is SHA256 hash of the serialization of Input and Ouput array
The outer hash is the hash of the whole transaction serialization
*/

type Transaction struct {
	Length    uint32        //length prefix
	Type      uint8         //transaction type
	InnerHash cipher.SHA256 //inner hash SHA256 of In[],Out[]

	Sigs []cipher.Sig        //list of signatures, 64+1 bytes each
	In   []cipher.SHA256     //ouputs being spent
	Out  []TransactionOutput //ouputs being created
}

//hash output/name is function of Hash
type TransactionOutput struct {
	Address cipher.Address //address to send to
	Coins   uint64         //amount to be sent in coins
	Hours   uint64         //amount to be sent in coin hours
}

// Verify attempts to determine if the transaction is well formed
// Verify cannot check transaction signatures, it needs the address from unspents
// Verify cannot check if outputs being spent exist
// Verify cannot check if the transaction would create or destroy coins
// or if the inputs have the required coin base
func (self *Transaction) Verify() error {

	h := self.HashInner()
	if h != self.InnerHash {
		return errors.New("Invalid header hash")
	}

	if len(self.In) == 0 {
		return errors.New("No inputs")
	}
	if len(self.Out) == 0 {
		return errors.New("No outputs")
	}

	// Check signature index fields
	if len(self.Sigs) != len(self.In) {
		return errors.New("Invalid number of signatures")
	}
	if len(self.Sigs) >= math.MaxUint16 {
		return errors.New("Too many signatures and inputs")
	}

	// Check duplicate inputs
	uxOuts := make(map[cipher.SHA256]int, len(self.In))
	for i, _ := range self.In {
		uxOuts[self.In[i]] = 1
	}
	if len(uxOuts) != len(self.In) {
		return errors.New("Duplicate spend")
	}

	if self.Type != 0 {
		return errors.New("transaction type invalid")
	}
	if self.Length != uint32(self.Size()) {
		return errors.New("transaction size prefix invalid")
	}

	// Check for duplicate potential outputs
	outputs := make(map[cipher.SHA256]int, len(self.Out))
	uxb := UxBody{
		SrcTransaction: self.Hash(),
	}
	for _, to := range self.Out {
		uxb.Coins = to.Coins
		uxb.Hours = to.Hours
		uxb.Address = to.Address
		outputs[uxb.Hash()] = 1
	}
	if len(outputs) != len(self.Out) {
		return errors.New("Duplicate output in transaction")
	}

	// Validate signature
	for i, sig := range self.Sigs {
		hash := cipher.AddSHA256(self.InnerHash, self.In[i])
		if err := cipher.VerifySignedHash(sig, hash); err != nil {
			return err
		}
	}

	// Artificial restriction to prevent spam
	// Must spend only multiples of 1e6
	for _, txo := range self.Out {
		if txo.Coins == 0 {
			return errors.New("Zero coin output")
		}
		if txo.Coins%1e6 != 0 {
			return errors.New("Transaction outputs must be multiple of 1e6 " +
				"base units")
		}
	}

	return nil
}

// Adds a UxArray to the Transaction given the hash of a UxOut.
// Returns the signature index for later signing
func (self *Transaction) PushInput(uxOut cipher.SHA256) uint16 {
	if len(self.In) >= math.MaxUint16 {
		log.Panic("Max transaction inputs reached")
	}
	self.In = append(self.In, uxOut)
	return uint16(len(self.In) - 1)
}

// Adds a TransactionOutput, sending coins & hours to an Address
func (self *Transaction) PushOutput(dst cipher.Address, coins, hours uint64) {
	to := TransactionOutput{
		Address: dst,
		Coins:   coins,
		Hours:   hours,
	}
	self.Out = append(self.Out, to)
}

// Signs all inputs in the transaction
func (self *Transaction) SignInputs(keys []cipher.SecKey) {
	self.InnerHash = self.HashInner() //update hash

	if len(self.Sigs) != 0 {
		log.Panic("Transaction has been signed")
	}
	if len(keys) != len(self.In) {
		log.Panic("Invalid number of keys")
	}
	if len(keys) > math.MaxUint16 {
		log.Panic("Too many key")
	}
	if len(keys) == 0 {
		log.Panic("No keys")
	}
	sigs := make([]cipher.Sig, len(self.In))
	inner_hash := self.HashInner()
	for i, k := range keys {
		h := cipher.AddSHA256(inner_hash, self.In[i]) //hash to sign
		sigs[i] = cipher.SignHash(h, k)
	}
	self.Sigs = sigs
}

// Returns the encoded byte size of the transaction
func (self *Transaction) Size() int {
	return len(self.Serialize())
}

// Hashes an entire Transaction struct, including the TransactionHeader
func (self *Transaction) Hash() cipher.SHA256 {
	b := self.Serialize()
	return cipher.SumSHA256((b)
}

// Returns the encoded size and the hash of it (avoids duplicate encoding)
func (self *Transaction) SizeHash() (int, cipher.SHA256) {
	b := self.Serialize()
	return len(b), cipher.SumSHA256((b)
}

// Saves the txn body hash to TransactionHeader.Hash
func (self *Transaction) UpdateHeader() {
	self.Length = uint32(self.Size())
	self.Type = byte(0x00)
	self.InnerHash = self.HashInner()
}

// Hashes only the Transaction Inputs & Outputs
// This is what is signed
// Client hashes the inner hash with hash of output being spent and signs it with private key
func (self *Transaction) HashInner() cipher.SHA256 {
	b1 := encoder.Serialize(self.In)
	b2 := encoder.Serialize(self.Out)
	b3 := append(b1, b2...)
	return cipher.SumSHA256(b3)
}

func (self *Transaction) Serialize() []byte {
	return encoder.Serialize(*self)
}

func TransactionDeserialize(b []byte) Transaction {
	t := Transaction{}
	if err := encoder.DeserializeRaw(b, &t); err != nil {
		log.Panic("Failed to deserialize transaction")
	}
	return t
}

// Returns the coin hours sent as outputs. This does not include the fee.
func (self *Transaction) OutputHours() uint64 {
	hours := uint64(0)
	for i, _ := range self.Out {
		hours += self.Out[i].Hours
	}
	return hours
}

type Transactions []Transaction

// Calculates all the fees in Transactions
func (self Transactions) Fees(calc FeeCalculator) (uint64, error) {
	total := uint64(0)
	for i, _ := range self {
		fee, err := calc(&self[i])
		if err != nil {
			return 0, err
		}
		total += fee
	}
	return total, nil
}

func (self Transactions) Hashes() []cipher.SHA256 {
	hashes := make([]cipher.SHA256, len(self))
	for i, _ := range self {
		hashes[i] = self[i].Hash()
	}
	return hashes
}

// Returns the sum of contained Transactions' sizes.  It is not the size if
// serialized, since that would have a length prefix.
func (self Transactions) Size() int {
	size := 0
	for i, _ := range self {
		size += self[i].Size()
	}
	return size
}

// Returns the first n transactions whose total size is less than or equal to
// size.
func (self Transactions) TruncateBytesTo(size int) Transactions {
	total := 0
	for i, _ := range self {
		pending := self[i].Size()
		if total+pending > size {
			return self[:i]
		}
		total += pending
	}
	return self
}

// Allows sorting transactions by fee & hash
type SortableTransactions struct {
	Txns   Transactions
	Fees   []uint64
	Hashes []cipher.SHA256
}

// Given a transaction, return its fee or an error if the fee cannot be
// calculated
type FeeCalculator func(*Transaction) (uint64, error)

// Returns transactions sorted by fee per kB, and sorted by lowest hash if
// tied.  Transactions that fail in fee computation are excluded.
func SortTransactions(txns Transactions,
	feeCalc FeeCalculator) Transactions {
	sorted := newSortableTransactions(txns, feeCalc)
	sorted.Sort()
	return sorted.Txns
}

// Returns an array of txns that can be sorted by fee.  On creation, fees are
// calculated, and if any txns have invalid fee, there are removed from
// consideration
func newSortableTransactions(txns Transactions,
	feeCalc FeeCalculator) SortableTransactions {
	newTxns := make(Transactions, len(txns))
	fees := make([]uint64, len(txns))
	hashes := make([]cipher.SHA256, len(txns))
	j := 0
	for i, _ := range txns {
		fee, err := feeCalc(&txns[i])
		if err == nil {
			newTxns[j] = txns[i]
			size := 0
			size, hashes[j] = txns[i].SizeHash()
			// Calculate fee priority based on fee per kb
			fees[j] = (fee * 1024) / uint64(size)
			j++
		}
	}
	return SortableTransactions{
		Txns:   newTxns[:j],
		Fees:   fees[:j],
		Hashes: hashes[:j],
	}
}

// Sorts by tx fee, and then by hash if fee equal
func (self SortableTransactions) Sort() {
	sort.Sort(self)
}

func (self SortableTransactions) IsSorted() bool {
	return sort.IsSorted(self)
}

func (self SortableTransactions) Len() int {
	return len(self.Txns)
}

// Default sorting is fees descending, hash ascending if fees equal
func (self SortableTransactions) Less(i, j int) bool {
	if self.Fees[i] == self.Fees[j] {
		// If fees match, hashes are sorted ascending
		return bytes.Compare(self.Hashes[i][:], self.Hashes[j][:]) < 0
	}
	// Fees are sorted descending
	return self.Fees[i] > self.Fees[j]
}

func (self SortableTransactions) Swap(i, j int) {
	self.Txns[i], self.Txns[j] = self.Txns[j], self.Txns[i]
	self.Fees[i], self.Fees[j] = self.Fees[j], self.Fees[i]
	self.Hashes[i], self.Hashes[j] = self.Hashes[j], self.Hashes[i]
}
