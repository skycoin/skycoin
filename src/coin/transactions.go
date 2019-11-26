package coin

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
)

var (
	// DebugLevel1 checks for extremely unlikely conditions (10e-40)
	DebugLevel1 = true
	// DebugLevel2 enable checks for impossible conditions
	DebugLevel2 = true
)

//go:generate skyencoder -struct Transaction -unexported
//go:generate skyencoder -struct transactionInputs
//go:generate skyencoder -struct transactionOutputs

type transactionInputs struct {
	In []cipher.SHA256 `enc:",maxlen=65535"`
}

type transactionOutputs struct {
	Out []TransactionOutput `enc:",maxlen=65535"`
}

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
	Length    uint32        // length prefix
	Type      uint8         // transaction type
	InnerHash cipher.SHA256 // inner hash SHA256 of In[],Out[]

	Sigs []cipher.Sig        `enc:",maxlen=65535"` // list of signatures, 64+1 bytes each
	In   []cipher.SHA256     `enc:",maxlen=65535"` // ouputs being spent
	Out  []TransactionOutput `enc:",maxlen=65535"` // ouputs being created
}

// TransactionOutput hash output/name is function of Hash
type TransactionOutput struct {
	Address cipher.Address // address to send to
	Coins   uint64         // amount to be sent in coins
	Hours   uint64         // amount to be sent in coin hours
}

// Verify attempts to determine if the transaction is well formed.
// Verify cannot check transaction signatures, it needs the address from unspents
// Verify cannot check if outputs being spent exist
// Verify cannot check if the transaction would create or destroy coins
// or if the inputs have the required coin base
func (txn *Transaction) Verify() error {
	return txn.verify(true)
}

// VerifyUnsigned attempts to determine if the transaction is well formed,
// but requires the transaction to have at least one null signature.
// Verify cannot check transaction signatures, it needs the address from unspents
// Verify cannot check if outputs being spent exist
// Verify cannot check if the transaction would create or destroy coins
// or if the inputs have the required coin base
func (txn *Transaction) VerifyUnsigned() error {
	return txn.verify(false)
}

func (txn *Transaction) verify(signed bool) error {
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
	if len(txn.Sigs) > math.MaxUint16 {
		return errors.New("Too many signatures and inputs")
	}

	if len(txn.Out) > math.MaxUint16 {
		return errors.New("Too many ouptuts")
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

	// Prevent zero coin outputs
	// Artificial restriction to prevent spam
	for _, txo := range txn.Out {
		if txo.Coins == 0 {
			return errors.New("Zero coin output")
		}
	}

	// Check output coin integer overflow
	coins := uint64(0)
	for _, to := range txn.Out {
		var err error
		coins, err = mathutil.AddUint64(coins, to.Coins)
		if err != nil {
			return errors.New("Output coins overflow")
		}
	}

	// Check that Size and Hash can be computed
	txnSize, txnHash, err := txn.SizeHash()
	if err != nil {
		return err
	}

	// Check txn Size set correctly
	if txn.Length != txnSize {
		return errors.New("Incorrect transaction length")
	}

	// Check for duplicate potential outputs
	outputs := make(map[cipher.SHA256]struct{}, len(txn.Out))
	uxb := UxBody{
		SrcTransaction: txnHash,
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

	// Check inner hash
	innerHash, err := txn.hashInner()
	if err != nil {
		return fmt.Errorf("HashInner failed: %v", err)
	}

	if innerHash != txn.InnerHash {
		return errors.New("InnerHash does not match computed hash")
	}

	// Validate signatures
	for i, sig := range txn.Sigs {
		if sig.Null() {
			// Check that signed transactions do not have any null signatures
			if signed {
				return errors.New("Unsigned input in transaction")
			}
			// Ignore null signatures if the transaction is unsigned
			continue
		}

		hash := cipher.AddSHA256(txn.InnerHash, txn.In[i])
		if err := cipher.VerifySignatureRecoverPubKey(sig, hash); err != nil {
			return err
		}
	}

	// Check that unsigned transactions have at least one non-null signature
	if !signed {
		if !txn.hasNullSignature() {
			return errors.New("Unsigned transaction must contain a null signature")
		}
	}

	return nil
}

func (txn Transaction) verifyInputSignaturesPrelude(uxIn UxArray) error {
	if len(txn.In) != len(uxIn) {
		return errors.New("txn.In != uxIn")
	}
	if len(txn.In) != len(txn.Sigs) {
		return errors.New("txn.In != txn.Sigs")
	}
	if txn.InnerHash != txn.HashInner() {
		return errors.New("Invalid Tx Inner Hash")
	}
	for i := range txn.In {
		if txn.In[i] != uxIn[i].Hash() {
			return errors.New("Ux hash mismatch")
		}
	}
	return nil
}

// VerifyInputSignatures verifies the inputs and signatures
func (txn Transaction) VerifyInputSignatures(uxIn UxArray) error {
	if err := txn.verifyInputSignaturesPrelude(uxIn); err != nil {
		if DebugLevel2 {
			log.Panic(err)
		}
		return err
	}

	// Check signatures against unspent address
	for i := range txn.In {
		if txn.Sigs[i].Null() {
			return errors.New("Unsigned input in transaction")
		}

		hash := cipher.AddSHA256(txn.InnerHash, txn.In[i]) // use inner hash, not outer hash
		err := cipher.VerifyAddressSignedHash(uxIn[i].Body.Address, txn.Sigs[i], hash)
		if err != nil {
			return errors.New("Signature not valid for output being spent")
		}
	}

	return nil
}

// VerifyPartialInputSignatures verifies the inputs and signatures for signatures that are not null
func (txn Transaction) VerifyPartialInputSignatures(uxIn UxArray) error {
	if err := txn.verifyInputSignaturesPrelude(uxIn); err != nil {
		if DebugLevel2 {
			log.Panic(err)
		}
		return err
	}

	// Check signatures against unspent address for signatures that are not null
	for i := range txn.In {
		if txn.Sigs[i].Null() {
			continue
		}
		hash := cipher.AddSHA256(txn.InnerHash, txn.In[i]) // use inner hash, not outer hash
		err := cipher.VerifyAddressSignedHash(uxIn[i].Body.Address, txn.Sigs[i], hash)
		if err != nil {
			return errors.New("Signature not valid for output being spent")
		}
	}

	return nil
}

// PushInput adds a unspent output hash to the inputs of a Transaction.
func (txn *Transaction) PushInput(uxOut cipher.SHA256) error {
	if len(txn.In) >= math.MaxUint16 {
		return errors.New("Max transaction inputs reached")
	}
	txn.In = append(txn.In, uxOut)
	return nil
}

// UxID compute transaction output id
func (txOut TransactionOutput) UxID(txID cipher.SHA256) cipher.SHA256 {
	var x UxBody
	x.Coins = txOut.Coins
	x.Hours = txOut.Hours
	x.Address = txOut.Address
	x.SrcTransaction = txID
	return x.Hash()
}

// PushOutput Adds a TransactionOutput, sending coins & hours to an Address
func (txn *Transaction) PushOutput(dst cipher.Address, coins, hours uint64) error {
	if len(txn.Out) >= math.MaxUint16 {
		return errors.New("Max transaction outputs reached")
	}
	txn.Out = append(txn.Out, TransactionOutput{
		Address: dst,
		Coins:   coins,
		Hours:   hours,
	})
	return nil
}

// SignInput signs a specific input in the transaction.
// InnerHash should already be set to a valid value.
// Returns an error if the input is already signed
func (txn *Transaction) SignInput(key cipher.SecKey, index int) error {
	if index < 0 || index >= len(txn.In) {
		return errors.New("Signature index out of range")
	}

	if len(txn.Sigs) == 0 {
		txn.Sigs = make([]cipher.Sig, len(txn.In))
	}
	if len(txn.In) != len(txn.Sigs) {
		return errors.New("Number of signatures does not match number of inputs")
	}

	if !txn.Sigs[index].Null() {
		return errors.New("Input already signed")
	}

	h := cipher.AddSHA256(txn.InnerHash, txn.In[index])
	txn.Sigs[index] = cipher.MustSignHash(h, key)

	return nil
}

// SignInputs signs all inputs in the transaction
func (txn *Transaction) SignInputs(keys []cipher.SecKey) {
	if len(keys) != len(txn.In) {
		log.Panic("Invalid number of keys")
	}
	if len(keys) > math.MaxUint16 {
		log.Panic("Too many keys")
	}
	if len(keys) == 0 {
		log.Panic("No keys")
	}
	if len(txn.Sigs) > 0 && txn.hasNonNullSignature() {
		log.Panic("Transaction has been signed")
	}

	txn.InnerHash = txn.HashInner() // update hash

	sigs := make([]cipher.Sig, len(txn.In))
	for i, k := range keys {
		h := cipher.AddSHA256(txn.InnerHash, txn.In[i]) // hash to sign
		sigs[i] = cipher.MustSignHash(h, k)
	}
	txn.Sigs = sigs
}

// Size returns the encoded byte size of the transaction
func (txn *Transaction) Size() (uint32, error) {
	buf, err := txn.Serialize()
	if err != nil {
		return 0, err
	}
	return mathutil.IntToUint32(len(buf))
}

// IsFullyUnsigned returns true if the transaction is not signed for any input.
// Unsigned transactions have a full signature array, but the signatures are null.
// Returns true if the signatures array is empty.
func (txn *Transaction) IsFullyUnsigned() bool {
	for _, s := range txn.Sigs {
		if !s.Null() {
			return false
		}
	}

	return true
}

// IsFullySigned returns true if the transaction is fully signed.
// Returns true if the signatures array is empty.
func (txn *Transaction) IsFullySigned() bool {
	if len(txn.Sigs) == 0 {
		return false
	}

	for _, s := range txn.Sigs {
		if s.Null() {
			return false
		}
	}

	return true
}

// hasNonNullSignature returns true if the transaction has at least one non-null signature
func (txn *Transaction) hasNonNullSignature() bool {
	for _, s := range txn.Sigs {
		if !s.Null() {
			return true
		}
	}

	return false
}

// hasNullSignature returns true if the transaction has at least one null signature
func (txn *Transaction) hasNullSignature() bool {
	for _, s := range txn.Sigs {
		if s.Null() {
			return true
		}
	}

	return false
}

// Hash an entire Transaction struct, including the TransactionHeader
func (txn *Transaction) Hash() cipher.SHA256 {
	b, err := txn.Serialize()
	if err != nil {
		log.Panicf("Hash: txn.Serialize failed: %v", err)
	}
	return cipher.SumSHA256(b)
}

// SizeHash returns the encoded size and the hash of it (avoids duplicate encoding)
func (txn *Transaction) SizeHash() (uint32, cipher.SHA256, error) {
	b, err := txn.Serialize()
	if err != nil {
		return 0, cipher.SHA256{}, err
	}
	s, err := mathutil.IntToUint32(len(b))
	if err != nil {
		return 0, cipher.SHA256{}, err
	}
	return s, cipher.SumSHA256(b), nil
}

// UpdateHeader saves the txn body hash to TransactionHeader.Hash
func (txn *Transaction) UpdateHeader() error {
	s, err := txn.Size()
	if err != nil {
		return err
	}
	txn.Length = s
	txn.Type = byte(0x00)
	txn.InnerHash = txn.HashInner()
	return nil
}

// HashInner hashes only the Transaction Inputs & Outputs
// This is what is signed
// Client hashes the inner hash with hash of output being spent and signs it with private key
func (txn *Transaction) HashInner() cipher.SHA256 {
	h, err := txn.hashInner()
	if err != nil {
		log.Panicf("hashInner failed: %v", err)
	}
	return h
}

func (txn *Transaction) hashInner() (cipher.SHA256, error) {
	txnInputs := &transactionInputs{
		In: txn.In,
	}
	txnOutputs := &transactionOutputs{
		Out: txn.Out,
	}
	n1 := encodeSizeTransactionInputs(txnInputs)
	n2 := encodeSizeTransactionOutputs(txnOutputs)
	buf := make([]byte, n1+n2)

	if err := encodeTransactionInputsToBuffer(buf[:n1], txnInputs); err != nil {
		return cipher.SHA256{}, fmt.Errorf("encodeTransactionInputsToBuffer failed: %v", err)
	}

	if err := encodeTransactionOutputsToBuffer(buf[n1:], txnOutputs); err != nil {
		return cipher.SHA256{}, fmt.Errorf("encodeTransactionOutputsToBuffer failed: %v", err)
	}

	return cipher.SumSHA256(buf), nil
}

// MustSerialize serializes the transaction to bytes, panics on error.
// Serialization can fail if the transaction has too many elements in its arrays
func (txn *Transaction) MustSerialize() []byte {
	b, err := encodeTransaction(txn)
	if err != nil {
		log.Panicf("encodeTransaction failed: %v", err)
	}
	return b
}

// Serialize serializes the transaction to bytes.
// Serialization can fail if the transaction has too many elements in its arrays
func (txn *Transaction) Serialize() ([]byte, error) {
	return encodeTransaction(txn)
}

// MustSerializeHex serializes the transaction to a hex string, panics on error.
// Serialization can fail if the transaction has too many elements in its arrays
func (txn *Transaction) MustSerializeHex() string {
	return hex.EncodeToString(txn.MustSerialize())
}

// SerializeHex serializes the transaction to a hex string.
// Serialization can fail if the transaction has too many elements in its arrays
func (txn *Transaction) SerializeHex() (string, error) {
	b, err := txn.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// MustDeserializeTransaction deserializes a transaction, panics on error
func MustDeserializeTransaction(b []byte) Transaction {
	txn, err := DeserializeTransaction(b)
	if err != nil {
		log.Panicf("Failed to deserialize transaction: %v", err)
	}
	return txn
}

// DeserializeTransaction deserializes a transaction
func DeserializeTransaction(b []byte) (Transaction, error) {
	txn := Transaction{}
	if err := decodeTransactionExact(b, &txn); err != nil {
		return Transaction{}, fmt.Errorf("Invalid transaction: %v", err)
	}
	return txn, nil
}

// MustDeserializeTransactionHex deserializes a transaction hex string, panics on error
func MustDeserializeTransactionHex(s string) Transaction {
	txn, err := DeserializeTransactionHex(s)
	if err != nil {
		log.Panicf("Failed to deserialize transaction: %v", err)
	}
	return txn
}

// DeserializeTransactionHex deserializes a transaction hex string
func DeserializeTransactionHex(s string) (Transaction, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Transaction{}, err
	}
	return DeserializeTransaction(b)
}

// OutputHours returns the coin hours sent as outputs. This does not include the fee.
func (txn *Transaction) OutputHours() (uint64, error) {
	hours := uint64(0)
	for i := range txn.Out {
		var err error
		hours, err = mathutil.AddUint64(hours, txn.Out[i].Hours)
		if err != nil {
			return 0, errors.New("Transaction output hours overflow")
		}
	}
	return hours, nil
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

		total, err = mathutil.AddUint64(total, fee)
		if err != nil {
			return 0, errors.New("Transactions fee totals overflow")
		}
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
func (txns Transactions) Size() (uint32, error) {
	var size uint32
	for i := range txns {
		s, err := txns[i].Size()
		if err != nil {
			return 0, err
		}

		size, err = mathutil.AddUint32(size, s)
		if err != nil {
			return 0, err
		}
	}

	return size, nil
}

// TruncateBytesTo returns the first n transactions whose total size is less than or equal to size
func (txns Transactions) TruncateBytesTo(size uint32) (Transactions, error) {
	var total uint32
	for i := range txns {
		pending, err := txns[i].Size()
		if err != nil {
			return nil, err
		}

		pendingTotal, err := mathutil.AddUint32(total, pending)
		if err != nil {
			return txns[:i], nil
		}

		if pendingTotal > size {
			return txns[:i], nil
		}

		total = pendingTotal
	}

	return txns, nil
}

// SortableTransactions allows sorting transactions by fee & hash
type SortableTransactions struct {
	Transactions Transactions
	Fees         []uint64
	Hashes       []cipher.SHA256
}

// FeeCalculator given a transaction, return its fee or an error if the fee cannot be calculated
type FeeCalculator func(*Transaction) (uint64, error)

// SortTransactions returns transactions sorted by fee per kB, and sorted by lowest hash if tied.
// Transactions that fail in fee computation are excluded
func SortTransactions(txns Transactions, feeCalc FeeCalculator) (Transactions, error) {
	sorted, err := NewSortableTransactions(txns, feeCalc)
	if err != nil {
		return nil, err
	}
	sorted.Sort()
	return sorted.Transactions, nil
}

// NewSortableTransactions returns an array of txns that can be sorted by fee.
// On creation, fees are calculated, and if any txns have invalid fee, there are removed from consideration
func NewSortableTransactions(txns Transactions, feeCalc FeeCalculator) (*SortableTransactions, error) {
	newTxns := make(Transactions, len(txns))
	fees := make([]uint64, len(txns))
	hashes := make([]cipher.SHA256, len(txns))
	j := 0
	for i := range txns {
		fee, err := feeCalc(&txns[i])
		if err != nil {
			continue
		}

		size, hash, err := txns[i].SizeHash()
		if err != nil {
			return nil, err
		}

		// Calculate fee priority based on fee per kb
		feeKB, err := mathutil.MultUint64(fee, 1024)

		// If the fee * 1024 would exceed math.MaxUint64, set it to math.MaxUint64 so that
		// this transaction can still be processed
		if err != nil {
			feeKB = math.MaxUint64
		}

		newTxns[j] = txns[i]
		hashes[j] = hash
		fees[j] = feeKB / uint64(size)
		j++
	}

	return &SortableTransactions{
		Transactions: newTxns[:j],
		Fees:         fees[:j],
		Hashes:       hashes[:j],
	}, nil
}

// Sort sorts by tx fee, and then by hash if fee equal
func (txns SortableTransactions) Sort() {
	sort.Sort(txns)
}

// Len returns length of transactions
func (txns SortableTransactions) Len() int {
	return len(txns.Transactions)
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
	txns.Transactions[i], txns.Transactions[j] = txns.Transactions[j], txns.Transactions[i]
	txns.Fees[i], txns.Fees[j] = txns.Fees[j], txns.Fees[i]
	txns.Hashes[i], txns.Hashes[j] = txns.Hashes[j], txns.Hashes[i]
}

// VerifyTransactionCoinsSpending checks that coins are not destroyed or created by the transaction
func VerifyTransactionCoinsSpending(uxIn UxArray, uxOut UxArray) error {
	coinsIn := uint64(0)
	for i := range uxIn {
		var err error
		coinsIn, err = mathutil.AddUint64(coinsIn, uxIn[i].Body.Coins)
		if err != nil {
			return errors.New("Transaction input coins overflow")
		}
	}

	coinsOut := uint64(0)
	for i := range uxOut {
		var err error
		coinsOut, err = mathutil.AddUint64(coinsOut, uxOut[i].Body.Coins)
		if err != nil {
			return errors.New("Transaction output coins overflow")
		}
	}

	if coinsIn < coinsOut {
		return errors.New("Insufficient coins")
	}
	if coinsIn > coinsOut {
		return errors.New("Transactions may not destroy coins")
	}

	return nil
}

// VerifyTransactionHoursSpending checks that hours are not created by the transaction
func VerifyTransactionHoursSpending(headTime uint64, uxIn UxArray, uxOut UxArray) error {
	hoursIn := uint64(0)
	for i := range uxIn {
		uxHours, err := uxIn[i].CoinHours(headTime)
		if err != nil {
			// If the error was specifically an overflow when adding the
			// earned coin hours to the base coin hours, treat the uxHours as 0.
			// Block 13277 spends an input which overflows in this way,
			// so the block will not sync if an error is returned.
			if err == ErrAddEarnedCoinHoursAdditionOverflow {
				uxHours = 0
			} else {
				return err
			}
		}

		hoursIn, err = mathutil.AddUint64(hoursIn, uxHours)
		if err != nil {
			return errors.New("Transaction input hours overflow")
		}
	}

	hoursOut := uint64(0)
	for i := range uxOut {
		// NOTE: addition of hours is not checked for overflow here because
		// this would invalidate existing blocks which had overflowed hours.
		// Hours overflow checks are handled as a "soft" constraint in the network
		// until those blocks are repaired.
		hoursOut += uxOut[i].Body.Hours
	}

	if hoursIn < hoursOut {
		return errors.New("Insufficient coin hours")
	}
	return nil
}
