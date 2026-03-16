package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
)

// TxDestination is a destination for a Bitcoin transaction
type TxDestination struct {
	Address string
	Value   int64 // satoshis
}

// EstimatedTxSize estimates the size of a P2PKH transaction in bytes.
// P2PKH input: ~148 bytes, P2PKH output: ~34 bytes, overhead: ~10 bytes
func EstimatedTxSize(numInputs, numOutputs int) int {
	return 10 + numInputs*148 + numOutputs*34
}

// SelectUTXOs selects UTXOs using a largest-first algorithm.
// Returns selected UTXOs and the total value of selected UTXOs.
func SelectUTXOs(utxos []UTXO, targetAmount int64) ([]UTXO, int64, error) {
	if len(utxos) == 0 {
		return nil, 0, ErrInsufficientFunds
	}

	// Sort UTXOs by value, largest first
	sorted := make([]UTXO, len(utxos))
	copy(sorted, utxos)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	var selected []UTXO
	var totalValue int64
	for _, utxo := range sorted {
		selected = append(selected, utxo)
		totalValue += utxo.Value
		if totalValue >= targetAmount {
			return selected, totalValue, nil
		}
	}

	return nil, 0, ErrInsufficientFunds
}

// BuildTransaction creates a raw Bitcoin P2PKH transaction.
// keys maps address string to cipher.SecKey for signing inputs.
func BuildTransaction(inputs []UTXO, destinations []TxDestination, changeAddr string, feePerByte int64, keys map[string]cipher.SecKey) (string, error) {
	if len(inputs) == 0 {
		return "", fmt.Errorf("no inputs provided")
	}
	if len(destinations) == 0 {
		return "", fmt.Errorf("no destinations provided")
	}

	// Calculate total input value
	var totalIn int64
	for _, input := range inputs {
		totalIn += input.Value
	}

	// Calculate total output value (before change)
	var totalOut int64
	for _, dest := range destinations {
		totalOut += dest.Value
	}

	// Calculate output count (destinations + possible change)
	numOutputs := len(destinations) + 1 // +1 for possible change
	estimatedSize := EstimatedTxSize(len(inputs), numOutputs)
	fee := int64(estimatedSize) * feePerByte

	change := totalIn - totalOut - fee
	if change < 0 {
		return "", fmt.Errorf("%w: need %d satoshis (including %d fee), have %d",
			ErrInsufficientFunds, totalOut+fee, fee, totalIn)
	}

	// Build outputs
	var outputs []txOut

	for _, dest := range destinations {
		script, err := p2pkhScript(dest.Address)
		if err != nil {
			return "", fmt.Errorf("destination address %s: %w", dest.Address, err)
		}
		outputs = append(outputs, txOut{value: dest.Value, script: script})
	}

	// Add change output if dust threshold is met (546 satoshis)
	if change >= 546 {
		changeScript, err := p2pkhScript(changeAddr)
		if err != nil {
			return "", fmt.Errorf("change address %s: %w", changeAddr, err)
		}
		outputs = append(outputs, txOut{value: change, script: changeScript})
	}

	// Build the unsigned transaction for signing
	// Bitcoin transaction format:
	// version (4 bytes, little-endian)
	// input count (varint)
	// inputs
	// output count (varint)
	// outputs
	// locktime (4 bytes)

	// First, serialize the unsigned tx template (for signing each input)
	// For P2PKH signing, each input is signed individually with SIGHASH_ALL
	var signedInputs []signedInput
	for i, input := range inputs {
		key, ok := keys[input.Address]
		if !ok {
			return "", fmt.Errorf("no key found for input address %s", input.Address)
		}

		// Create the script for the input being signed (the previous output's P2PKH script)
		prevScript, err := p2pkhScript(input.Address)
		if err != nil {
			return "", fmt.Errorf("input address %s: %w", input.Address, err)
		}

		// Build the transaction with the current input's scriptPubKey for signing
		sigHash, err := computeSigHash(inputs, outputs, i, prevScript)
		if err != nil {
			return "", fmt.Errorf("compute sighash for input %d: %w", i, err)
		}

		// Sign the hash
		sig, err := cipher.SignHash(sigHash, key)
		if err != nil {
			return "", fmt.Errorf("sign input %d: %w", i, err)
		}

		// Convert the 65-byte recoverable signature to DER format
		derSig := sigToDER(sig)
		// Append SIGHASH_ALL byte
		derSig = append(derSig, 0x01)

		// Get the compressed public key
		pubKey := cipher.MustPubKeyFromSecKey(key)

		signedInputs = append(signedInputs, signedInput{
			txid:      input.TxID,
			vout:      input.Vout,
			scriptSig: buildP2PKHScriptSig(derSig, pubKey[:]),
		})
	}

	// Serialize the final signed transaction
	var buf bytes.Buffer

	// Version
	writeUint32LE(&buf, 1)

	// Input count
	writeVarInt(&buf, uint64(len(signedInputs)))

	// Inputs
	for _, in := range signedInputs {
		txidBytes, _ := hex.DecodeString(in.txid)
		// Reverse txid (Bitcoin uses internal byte order)
		reverseBytes(txidBytes)
		buf.Write(txidBytes)
		writeUint32LE(&buf, in.vout)
		writeVarInt(&buf, uint64(len(in.scriptSig)))
		buf.Write(in.scriptSig)
		writeUint32LE(&buf, 0xffffffff) // sequence
	}

	// Output count
	writeVarInt(&buf, uint64(len(outputs)))

	// Outputs
	for _, out := range outputs {
		writeUint64LE(&buf, uint64(out.value))
		writeVarInt(&buf, uint64(len(out.script)))
		buf.Write(out.script)
	}

	// Locktime
	writeUint32LE(&buf, 0)

	return hex.EncodeToString(buf.Bytes()), nil
}

type txOut struct {
	value  int64
	script []byte
}

type signedInput struct {
	txid      string
	vout      uint32
	scriptSig []byte
}

// p2pkhScript creates a P2PKH output script for a Bitcoin address
func p2pkhScript(address string) ([]byte, error) {
	addr, err := cipher.DecodeBase58BitcoinAddress(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	// OP_DUP OP_HASH160 <20-byte hash> OP_EQUALVERIFY OP_CHECKSIG
	script := make([]byte, 25)
	script[0] = 0x76 // OP_DUP
	script[1] = 0xa9 // OP_HASH160
	script[2] = 0x14 // Push 20 bytes
	copy(script[3:23], addr.Key[:])
	script[23] = 0x88 // OP_EQUALVERIFY
	script[24] = 0xac // OP_CHECKSIG
	return script, nil
}

// computeSigHash computes the signature hash for a P2PKH input (SIGHASH_ALL)
func computeSigHash(inputs []UTXO, outputs []txOut, sigIndex int, prevScript []byte) (cipher.SHA256, error) {
	var buf bytes.Buffer

	// Version
	writeUint32LE(&buf, 1)

	// Input count
	writeVarInt(&buf, uint64(len(inputs)))

	// Inputs
	for i, input := range inputs {
		txidBytes, err := hex.DecodeString(input.TxID)
		if err != nil {
			return cipher.SHA256{}, fmt.Errorf("decode txid: %w", err)
		}
		reverseBytes(txidBytes)
		buf.Write(txidBytes)
		writeUint32LE(&buf, input.Vout)

		if i == sigIndex {
			// Include the previous output script for the input being signed
			writeVarInt(&buf, uint64(len(prevScript)))
			buf.Write(prevScript)
		} else {
			// Empty script for other inputs
			writeVarInt(&buf, 0)
		}
		writeUint32LE(&buf, 0xffffffff) // sequence
	}

	// Output count
	writeVarInt(&buf, uint64(len(outputs)))

	// Outputs
	for _, out := range outputs {
		writeUint64LE(&buf, uint64(out.value))
		writeVarInt(&buf, uint64(len(out.script)))
		buf.Write(out.script)
	}

	// Locktime
	writeUint32LE(&buf, 0)

	// SIGHASH_ALL
	writeUint32LE(&buf, 1)

	// Double SHA256
	first := sha256.Sum256(buf.Bytes())
	second := sha256.Sum256(first[:])
	return cipher.SHA256(second), nil
}

// sigToDER converts a 65-byte recoverable signature (R[32] || S[32] || V[1]) to DER format
func sigToDER(sig cipher.Sig) []byte {
	r := new(big.Int).SetBytes(sig[0:32])
	s := new(big.Int).SetBytes(sig[32:64])

	// Enforce low-S (BIP 62) — if S > N/2, replace with N - S
	// secp256k1 order N
	n, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	halfN := new(big.Int).Rsh(n, 1)
	if s.Cmp(halfN) > 0 {
		s.Sub(n, s)
	}

	rBytes := r.Bytes()
	sBytes := s.Bytes()

	// Prepend 0x00 if high bit is set (DER integer encoding)
	if len(rBytes) > 0 && rBytes[0]&0x80 != 0 {
		rBytes = append([]byte{0x00}, rBytes...)
	}
	if len(sBytes) > 0 && sBytes[0]&0x80 != 0 {
		sBytes = append([]byte{0x00}, sBytes...)
	}

	// DER encoding: 0x30 <total len> 0x02 <r len> <r> 0x02 <s len> <s>
	totalLen := 2 + len(rBytes) + 2 + len(sBytes)
	der := make([]byte, 0, totalLen+2)
	der = append(der, 0x30, byte(totalLen))
	der = append(der, 0x02, byte(len(rBytes)))
	der = append(der, rBytes...)
	der = append(der, 0x02, byte(len(sBytes)))
	der = append(der, sBytes...)

	return der
}

// buildP2PKHScriptSig creates a P2PKH scriptSig: <sig> <pubkey>
func buildP2PKHScriptSig(sig, pubKey []byte) []byte {
	script := make([]byte, 0, 1+len(sig)+1+len(pubKey))
	script = append(script, byte(len(sig)))
	script = append(script, sig...)
	script = append(script, byte(len(pubKey)))
	script = append(script, pubKey...)
	return script
}

func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func writeUint32LE(buf *bytes.Buffer, v uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	buf.Write(b)
}

func writeUint64LE(buf *bytes.Buffer, v uint64) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	buf.Write(b)
}

func writeVarInt(buf *bytes.Buffer, v uint64) {
	switch {
	case v < 0xfd:
		buf.WriteByte(byte(v))
	case v <= 0xffff:
		buf.WriteByte(0xfd)
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v))
		buf.Write(b)
	case v <= 0xffffffff:
		buf.WriteByte(0xfe)
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v))
		buf.Write(b)
	default:
		buf.WriteByte(0xff)
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)
		buf.Write(b)
	}
}
