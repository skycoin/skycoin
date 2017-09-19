package coin

import (
	"bytes"
	"errors"
	"math"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var (
	// DebugLevel1 checks for extremely unlikely conditions (10e-40)
	DebugLevel1 = true
	// DebugLevel2 enable checks for impossible conditions
	DebugLevel2 = true
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

The inner hash is SHA256 hash of the serialization of Input and Output array
The outer hash is the hash of the whole transaction serialization
*/

// Transaction transaction struct
type Transaction struct {
	Length    uint32        //length prefix
	Type      uint8         //transaction type
	InnerHash cipher.SHA256 //inner hash SHA256 of In[],Out[]

	Sigs []cipher.Sig        //list of signatures, 64+1 bytes each
	In   []cipher.SHA256     //ouputs being spent
	Out  []TransactionOutput //ouputs being created
}

// TransactionOutput hash output/name is function of Hash
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
func (txn *Transaction) Verify() error {

	h := txn.HashInner()
	if h != txn.InnerHash {
		return errors.New("Invalid header hash")
	}

	if len(txn.In) == 0 {
		return errors.New("No inputs")
	}
	if len(txn.Out) == 0 {
		return errors.New("No outputs")
	}

	// Check signature index fields
	if len(txn.Sigs) != len(txn.In) {
		return errors.New("Invalid number of signatures")
	}
	if len(txn.Sigs) >= math.MaxUint16 {
		return errors.New("Too many signatures and inputs")
	}

	// Check duplicate inputs
	uxOuts := make(map[cipher.SHA256]struct{}, len(txn.In))
	for i := range txn.In {
		uxOuts[txn.In[i]] = struct{}{}
	}
	if len(uxOuts) != len(txn.In) {
		return errors.New("Duplicate spend")
	}

	if txn.Type != 0 {
		return errors.New("transaction type invalid")
	}
	if txn.Length != uint32(txn.Size()) {
		return errors.New("transaction size prefix invalid")
	}

	// Check for duplicate potential outputs
	outputs := make(map[cipher.SHA256]struct{}, len(txn.Out))
	uxb := UxBody{
		SrcTransaction: txn.Hash(),
	}
	for _, to := range txn.Out {
		uxb.Coins = to.Coins
		uxb.Hours = to.Hours
		uxb.Address = to.Address
		outputs[uxb.Hash()] = struct{}{}
	}
	if len(outputs) != len(txn.Out) {
		return errors.New("Duplicate output in transaction")
	}

	// Validate signature
	for i, sig := range txn.Sigs {
		hash := cipher.AddSHA256(txn.InnerHash, txn.In[i])
		if err := cipher.VerifySignedHash(sig, hash); err != nil {
			return err
		}
	}

	// Artificial restriction to prevent spam
	for _, txo := range txn.Out {
		if txo.Coins == 0 {
			return errors.New("Zero coin output")
		}
	}

	return nil
}

// VerifyInput verifies the input
func (txn Transaction) VerifyInput(uxIn UxArray) error {
	if DebugLevel2 {
		if len(txn.In) != len(txn.Sigs) || len(txn.In) != len(uxIn) {
			logger.Panic("tx.In != tx.Sigs != uxIn")
		}
		if txn.InnerHash != txn.HashInner() {
			logger.Panic("Invalid Tx Header Hash")
		}
	}

	// Check signatures against unspent address
	for i := range txn.In {
		hash := cipher.AddSHA256(txn.InnerHash, txn.In[i]) //use inner hash, not outer hash
		err := cipher.ChkSig(uxIn[i].Body.Address, hash, txn.Sigs[i])
		if err != nil {
			return errors.New("Signature not valid for output being spent")
		}
	}
	if DebugLevel2 {
		// Check that hashes match.
		// This would imply a bug with UnspentPool.GetMultiple
		if len(txn.In) != len(uxIn) {
			logger.Panic("tx.In does not match uxIn")
		}
		for i := range txn.In {
			if txn.In[i] != uxIn[i].Hash() {
				logger.Panic("impossible error: Ux hash mismatch")
			}
		}
	}
	return nil
}

// PushInput adds a UxArray to the Transaction given the hash of a UxOut.
// Returns the signature index for later signing
func (txn *Transaction) PushInput(uxOut cipher.SHA256) uint16 {
	if len(txn.In) >= math.MaxUint16 {
		logger.Panic("Max transaction inputs reached")
	}
	txn.In = append(txn.In, uxOut)
	return uint16(len(txn.In) - 1)
}

// UxID compute transaction output id
func (txOut TransactionOutput) UxID(TxID cipher.SHA256) cipher.SHA256 {
	var x UxBody
	x.Coins = txOut.Coins
	x.Hours = txOut.Hours
	x.Address = txOut.Address
	x.SrcTransaction = TxID
	return x.Hash()
}

// PushOutput Adds a TransactionOutput, sending coins & hours to an Address
func (txn *Transaction) PushOutput(dst cipher.Address, coins, hours uint64) {
	to := TransactionOutput{
		Address: dst,
		Coins:   coins,
		Hours:   hours,
	}
	txn.Out = append(txn.Out, to)
}

// SignInputs signs all inputs in the transaction
func (txn *Transaction) SignInputs(keys []cipher.SecKey) {
	txn.InnerHash = txn.HashInner() //update hash

	if len(txn.Sigs) != 0 {
		logger.Panic("Transaction has been signed")
	}
	if len(keys) != len(txn.In) {
		logger.Panic("Invalid number of keys")
	}
	if len(keys) > math.MaxUint16 {
		logger.Panic("Too many key")
	}
	if len(keys) == 0 {
		logger.Panic("No keys")
	}
	sigs := make([]cipher.Sig, len(txn.In))
	innerHash := txn.HashInner()
	for i, k := range keys {
		h := cipher.AddSHA256(innerHash, txn.In[i]) // hash to sign
		sigs[i] = cipher.SignHash(h, k)
	}
	txn.Sigs = sigs
}

// Size returns the encoded byte size of the transaction
func (txn *Transaction) Size() int {
	return len(txn.Serialize())
}

// Hash an entire Transaction struct, including the TransactionHeader
func (txn *Transaction) Hash() cipher.SHA256 {
	b := txn.Serialize()
	return cipher.SumSHA256(b)
}

// SizeHash returns the encoded size and the hash of it (avoids duplicate encoding)
func (txn *Transaction) SizeHash() (int, cipher.SHA256) {
	b := txn.Serialize()
	return len(b), cipher.SumSHA256(b)
}

// TxID returns transaction ID as byte string
func (txn *Transaction) TxID() []byte {
	hash := txn.Hash()
	return hash[0:32]
}

// TxIDHex returns transaction ID as hex
func (txn *Transaction) TxIDHex() string {
	return txn.Hash().Hex()
}

// UpdateHeader saves the txn body hash to TransactionHeader.Hash
func (txn *Transaction) UpdateHeader() {
	txn.Length = uint32(txn.Size())
	txn.Type = byte(0x00)
	txn.InnerHash = txn.HashInner()
}

// HashInner hashes only the Transaction Inputs & Outputs
// This is what is signed
// Client hashes the inner hash with hash of output being spent and signs it with private key
func (txn *Transaction) HashInner() cipher.SHA256 {
	b1 := encoder.Serialize(txn.In)
	b2 := encoder.Serialize(txn.Out)
	b3 := append(b1, b2...)
	return cipher.SumSHA256(b3)
}

// Serialize serialize the transaction
func (txn *Transaction) Serialize() []byte {
	return encoder.Serialize(*txn)
}

// TransactionDeserialize deserialize transaction
func TransactionDeserialize(b []byte) Transaction {
	t := Transaction{}
	if err := encoder.DeserializeRaw(b, &t); err != nil {
		logger.Panic("Failed to deserialize transaction")
	}
	return t
}

// OutputHours returns the coin hours sent as outputs. This does not include the fee.
func (txn *Transaction) OutputHours() uint64 {
	hours := uint64(0)
	for i := range txn.Out {
		hours += txn.Out[i].Hours
	}
	return hours
}

// Transactions transaction slice
type Transactions []Transaction

// Fees calculates all the fees in Transactions
func (txns Transactions) Fees(calc FeeCalculator) (uint64, error) {
	total := uint64(0)
	for i := range txns {
		fee, err := calc(&txns[i])
		if err != nil {
			return 0, err
		}
		total += fee
	}
	return total, nil
}

// Hashes caculate transactions hashes
func (txns Transactions) Hashes() []cipher.SHA256 {
	hashes := make([]cipher.SHA256, len(txns))
	for i := range txns {
		hashes[i] = txns[i].Hash()
	}
	return hashes
}

// Size returns the sum of contained Transactions' sizes.  It is not the size if
// serialized, since that would have a length prefix.
func (txns Transactions) Size() int {
	size := 0
	for i := range txns {
		size += txns[i].Size()
	}
	return size
}

// TruncateBytesTo returns the first n transactions whose total size is less than or equal to
// size.
func (txns Transactions) TruncateBytesTo(size int) Transactions {
	total := 0
	for i := range txns {
		pending := txns[i].Size()
		if total+pending > size {
			return txns[:i]
		}
		total += pending
	}
	return txns
}

// SortableTransactions allows sorting transactions by fee & hash
type SortableTransactions struct {
	Txns   Transactions
	Fees   []uint64
	Hashes []cipher.SHA256
}

// FeeCalculator given a transaction, return its fee or an error if the fee cannot be
// calculated
type FeeCalculator func(*Transaction) (uint64, error)

// SortTransactions returns transactions sorted by fee per kB, and sorted by lowest hash if
// tied.  Transactions that fail in fee computation are excluded.
func SortTransactions(txns Transactions,
	feeCalc FeeCalculator) Transactions {
	sorted := NewSortableTransactions(txns, feeCalc)
	sorted.Sort()
	return sorted.Txns
}

// NewSortableTransactions returns an array of txns that can be sorted by fee.  On creation, fees are
// calculated, and if any txns have invalid fee, there are removed from
// consideration
func NewSortableTransactions(txns Transactions, feeCalc FeeCalculator) SortableTransactions {
	newTxns := make(Transactions, len(txns))
	fees := make([]uint64, len(txns))
	hashes := make([]cipher.SHA256, len(txns))
	j := 0
	for i := range txns {
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

// Sort sorts by tx fee, and then by hash if fee equal
func (txns SortableTransactions) Sort() {
	sort.Sort(txns)
}

// IsSorted checks if transactions are sorted
func (txns SortableTransactions) IsSorted() bool {
	return sort.IsSorted(txns)
}

// Len returns length of transactions
func (txns SortableTransactions) Len() int {
	return len(txns.Txns)
}

// Less default sorting is fees descending, hash ascending if fees equal
func (txns SortableTransactions) Less(i, j int) bool {
	if txns.Fees[i] == txns.Fees[j] {
		// If fees match, hashes are sorted ascending
		return bytes.Compare(txns.Hashes[i][:], txns.Hashes[j][:]) < 0
	}
	// Fees are sorted descending
	return txns.Fees[i] > txns.Fees[j]
}

// Swap swaps txns
func (txns SortableTransactions) Swap(i, j int) {
	txns.Txns[i], txns.Txns[j] = txns.Txns[j], txns.Txns[i]
	txns.Fees[i], txns.Fees[j] = txns.Fees[j], txns.Fees[i]
	txns.Hashes[i], txns.Hashes[j] = txns.Hashes[j], txns.Hashes[i]
}

// VerifyTransactionSpending checks that coins will not be destroyed and that enough coins are hours
// are being spent for the outputs
func VerifyTransactionSpending(headTime uint64, uxIn UxArray, uxOut UxArray) error {
	coinsIn := uint64(0)
	hoursIn := uint64(0)
	for i := range uxIn {
		coinsIn += uxIn[i].Body.Coins
		hoursIn += uxIn[i].CoinHours(headTime)
	}
	coinsOut := uint64(0)
	hoursOut := uint64(0)
	for i := range uxOut {
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
