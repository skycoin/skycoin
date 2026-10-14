package coin

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/util/mathutil"
)

// Transaction types. A TxnTypePQ transaction spends at least one input owned by a
// post-quantum (version 1) address, and carries a PQ witness for those inputs.
// Nodes from before the post-quantum fork reject any type but TxnTypeStandard.
const (
	TxnTypeStandard uint8 = 0
	TxnTypePQ       uint8 = 1
)

var (
	// ErrPQWitnessMalformed is returned for a PQ witness that cannot be decoded
	ErrPQWitnessMalformed = errors.New("malformed post-quantum witness")
)

// pqSigned reports whether input i carries a post-quantum signature.
func (txn *Transaction) pqSigned(i int) bool {
	return i < len(txn.PQ) && len(txn.PQ[i].Sig) != 0
}

// InputSigned reports whether input i carries a signature of either kind.
func (txn *Transaction) InputSigned(i int) bool {
	return !txn.Sigs[i].Null() || txn.pqSigned(i)
}

// verifyPQ checks the witness of a TxnTypePQ transaction against its own keys.
// Whether each key owns its input is checked in VerifyInputSignatures.
func (txn *Transaction) verifyPQ(signed bool) error {
	if len(txn.PQ) != len(txn.In) {
		return errors.New("Invalid number of post-quantum inputs")
	}
	found := false
	for i, p := range txn.PQ {
		if len(p.PubKey) == 0 && len(p.Sig) == 0 {
			continue
		}
		if !txn.Sigs[i].Null() {
			return errors.New("Input has both a secp256k1 and a post-quantum signature")
		}
		if len(p.PubKey) != pqsig.PubKeySize {
			return errors.New("Invalid post-quantum public key size")
		}
		if len(p.Sig) == 0 {
			if signed {
				return errors.New("Unsigned input in transaction")
			}
			found = true
			continue
		}
		if len(p.Sig) != pqsig.SigSize {
			return errors.New("Invalid post-quantum signature size")
		}
		h := cipher.AddSHA256(txn.InnerHash, txn.In[i])
		if err := pqsig.VerifyInput(p.PubKey, h[:], p.Sig); err != nil {
			return err
		}
		found = true
	}
	if !found {
		return errors.New("Post-quantum transaction has no post-quantum input")
	}
	return nil
}

// verifyPQInput checks that p owns an input at addr and signs hash.
func verifyPQInput(addr cipher.Address, p PQInputSig, hash cipher.SHA256) error {
	if cipher.AddressFromPQPubKey(p.PubKey) != addr {
		return cipher.ErrAddressInvalidPubKey
	}
	return pqsig.VerifyInput(p.PubKey, hash[:], p.Sig)
}

// SignPQInput signs input index with a post-quantum key, which makes the transaction
// TxnTypePQ. InnerHash should already be set. Call UpdateHeader afterwards.
func (txn *Transaction) SignPQInput(sk pqsig.SecKey, index int) error {
	if index < 0 || index >= len(txn.In) {
		return errors.New("Signature index out of range")
	}
	if len(txn.Sigs) == 0 {
		txn.Sigs = make([]cipher.Sig, len(txn.In))
	}
	if len(txn.PQ) == 0 {
		txn.PQ = make([]PQInputSig, len(txn.In))
	}
	if len(txn.Sigs) != len(txn.In) || len(txn.PQ) != len(txn.In) {
		return errors.New("Number of signatures does not match number of inputs")
	}
	if txn.InputSigned(index) {
		return errors.New("Input already signed")
	}

	pk, err := sk.PubKey()
	if err != nil {
		return err
	}
	h := cipher.AddSHA256(txn.InnerHash, txn.In[index])
	sig, err := pqsig.SignInput(sk, h[:])
	if err != nil {
		return err
	}
	txn.PQ[index] = PQInputSig{PubKey: pk, Sig: sig}
	txn.Type = TxnTypePQ
	return nil
}

// PQWitnessSize is the encoded size of the witness of a transaction with inputs
// inputs, pqInputs of them signed with ML-DSA-44. It is 0 when pqInputs is 0.
func PQWitnessSize(inputs, pqInputs int) int {
	if pqInputs == 0 {
		return 0
	}
	return 4 + 8*inputs + pqInputs*(pqsig.PubKeySize+pqsig.SigSize)
}

// EncodePQWitness encodes a PQ witness: a count, then each input's public key and
// signature, each length prefixed.
func EncodePQWitness(pq []PQInputSig) []byte {
	n := 4
	for _, p := range pq {
		n += 8 + len(p.PubKey) + len(p.Sig)
	}
	buf := make([]byte, n)
	e := &encoder.Encoder{Buffer: buf}
	e.Uint32(uint32(len(pq))) //nolint:gosec // bounded by the transaction's input count
	for _, p := range pq {
		e.Uint32(uint32(len(p.PubKey))) //nolint:gosec // at most pqsig.PubKeySize
		e.CopyBytes(p.PubKey)
		e.Uint32(uint32(len(p.Sig))) //nolint:gosec // at most pqsig.SigSize
		e.CopyBytes(p.Sig)
	}
	return buf
}

// DecodePQWitness decodes a PQ witness. It returns how many bytes it read.
func DecodePQWitness(b []byte) ([]PQInputSig, int, error) {
	d := &encoder.Decoder{Buffer: b}
	n, err := d.Uint32()
	if err != nil || n > math.MaxUint16 {
		return nil, 0, ErrPQWitnessMalformed
	}
	pq := make([]PQInputSig, n)
	for i := range pq {
		if pq[i].PubKey, err = decodeBytes(d, pqsig.PubKeySize); err != nil {
			return nil, 0, err
		}
		if pq[i].Sig, err = decodeBytes(d, pqsig.SigSize); err != nil {
			return nil, 0, err
		}
	}
	return pq, len(b) - len(d.Buffer), nil
}

func decodeBytes(d *encoder.Decoder, maxLen int) ([]byte, error) {
	l, err := d.Uint32()
	if err != nil || int(l) > maxLen || int(l) > len(d.Buffer) {
		return nil, ErrPQWitnessMalformed
	}
	if l == 0 {
		return nil, nil
	}
	out := append([]byte(nil), d.Buffer[:l]...)
	d.Buffer = d.Buffer[l:]
	return out, nil
}

// SerializeWithPQ is Serialize with the PQ witness appended for a TxnTypePQ
// transaction. Raw transactions are exchanged in this form.
func (txn *Transaction) SerializeWithPQ() ([]byte, error) {
	b, err := txn.Serialize()
	if err != nil || txn.Type != TxnTypePQ {
		return b, err
	}
	return append(b, EncodePQWitness(txn.PQ)...), nil
}

// SerializeHexWithPQ is SerializeWithPQ in hex.
func (txn *Transaction) SerializeHexWithPQ() (string, error) {
	b, err := txn.SerializeWithPQ()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// DeserializeTransactionWithPQ reads a transaction written by SerializeWithPQ.
func DeserializeTransactionWithPQ(b []byte) (Transaction, error) {
	txn, err := deserializeTransactionWithPQ(b)
	if err != nil {
		return Transaction{}, fmt.Errorf("Invalid transaction: %v", err)
	}
	return txn, nil
}

func deserializeTransactionWithPQ(b []byte) (Transaction, error) {
	var txn Transaction
	n, err := decodeTransaction(b, &txn)
	if err != nil {
		return Transaction{}, err
	}
	rest := b[n:]
	if txn.Type != TxnTypePQ {
		if len(rest) != 0 {
			return Transaction{}, encoder.ErrRemainingBytes
		}
		return txn, nil
	}
	pq, m, err := DecodePQWitness(rest)
	if err != nil {
		return Transaction{}, err
	}
	if m != len(rest) {
		return Transaction{}, encoder.ErrRemainingBytes
	}
	txn.PQ = pq
	return txn, nil
}

// DeserializeTransactionHexWithPQ is DeserializeTransactionWithPQ from hex.
func DeserializeTransactionHexWithPQ(s string) (Transaction, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Transaction{}, err
	}
	return DeserializeTransactionWithPQ(b)
}

// verifyInputSignature checks the signature of input i against the address that owns it.
func (txn *Transaction) verifyInputSignature(addr cipher.Address, i int) error {
	hash := cipher.AddSHA256(txn.InnerHash, txn.In[i]) // use inner hash, not outer hash
	if addr.Version == cipher.AddressVersionPQ {
		if !txn.pqSigned(i) {
			return errors.New("Post-quantum address needs a post-quantum signature")
		}
		return verifyPQInput(addr, txn.PQ[i], hash)
	}
	if txn.Sigs[i].Null() {
		return errors.New("Address needs a secp256k1 signature")
	}
	return cipher.VerifyAddressSignedHash(addr, txn.Sigs[i], hash)
}

// SizeWithPQ is the size of the transaction with its PQ witness, which is what a
// node stores and sends. Transaction size limits apply to it.
func (txn *Transaction) SizeWithPQ() (uint32, error) {
	s, err := txn.Size()
	if err != nil || txn.Type != TxnTypePQ {
		return s, err
	}
	return mathutil.AddUint32(s, uint32(len(EncodePQWitness(txn.PQ)))) //nolint:gosec // bounded by the input count and key sizes
}
