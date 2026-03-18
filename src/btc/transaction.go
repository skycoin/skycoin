package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

// TxDestination is a destination for a Bitcoin transaction
type TxDestination struct {
	Address string
	Value   int64 // satoshis
}

// EstimatedTxSize estimates the virtual size of a transaction in vbytes.
// For P2PKH: input ~148 bytes, output ~34 bytes, overhead ~10 bytes.
// For P2WPKH: input ~68 vbytes (witness discounted), output ~31 bytes, overhead ~10.75 bytes.
func EstimatedTxSize(numInputs, numOutputs int) int {
	return 10 + numInputs*148 + numOutputs*34
}

// EstimatedSegwitTxVSize estimates the virtual size of a segwit P2WPKH transaction.
func EstimatedSegwitTxVSize(numInputs, numOutputs int) int {
	// Non-witness: 10 + 41*inputs + 31*outputs
	// Witness: 1 + 107*inputs (approx)
	// vsize = (non_witness_weight*4 + witness_weight) / 4
	nonWitness := 10 + 41*numInputs + 31*numOutputs
	witness := 1 + 107*numInputs
	weight := nonWitness*4 + witness
	vsize := (weight + 3) / 4 // round up
	return vsize
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

	// Determine if any inputs are segwit
	hasSegwit := false
	for _, input := range inputs {
		if isSegwitAddress(input.Address) {
			hasSegwit = true
			break
		}
	}

	// Build outputs (supports both P2PKH and P2WPKH destinations)
	var outputs []txOut

	for _, dest := range destinations {
		script, err := outputScript(dest.Address)
		if err != nil {
			return "", fmt.Errorf("destination address %s: %w", dest.Address, err)
		}
		outputs = append(outputs, txOut{value: dest.Value, script: script})
	}

	// Add change output if dust threshold is met (546 satoshis for P2PKH, 294 for P2WPKH)
	dustThreshold := int64(546)
	if isSegwitAddress(changeAddr) {
		dustThreshold = 294
	}
	if change >= dustThreshold {
		changeScript, err := outputScript(changeAddr)
		if err != nil {
			return "", fmt.Errorf("change address %s: %w", changeAddr, err)
		}
		outputs = append(outputs, txOut{value: change, script: changeScript})
	}

	// Sign each input
	type inputSigning struct {
		txid      string
		vout      uint32
		scriptSig []byte
		witness   [][]byte // nil for non-segwit inputs
	}

	var signedInputs []inputSigning
	for i, input := range inputs {
		key, ok := keys[input.Address]
		if !ok {
			return "", fmt.Errorf("no key found for input address %s", input.Address)
		}

		pubKey := cipher.MustPubKeyFromSecKey(key)

		if isSegwitAddress(input.Address) {
			// BIP143 segwit signing for P2WPKH
			sigHash, err := computeSegwitSigHash(inputs, outputs, i, input.Value, pubKey)
			if err != nil {
				return "", fmt.Errorf("compute segwit sighash for input %d: %w", i, err)
			}

			sig, err := cipher.SignHash(sigHash, key)
			if err != nil {
				return "", fmt.Errorf("sign input %d: %w", i, err)
			}

			derSig := sigToDER(sig)
			derSig = append(derSig, 0x01) // SIGHASH_ALL

			signedInputs = append(signedInputs, inputSigning{
				txid:      input.TxID,
				vout:      input.Vout,
				scriptSig: nil, // empty scriptSig for native segwit
				witness:   [][]byte{derSig, pubKey[:]},
			})
		} else {
			// Legacy P2PKH signing
			prevScript, err := p2pkhScript(input.Address)
			if err != nil {
				return "", fmt.Errorf("input address %s: %w", input.Address, err)
			}

			sigHash, err := computeSigHash(inputs, outputs, i, prevScript)
			if err != nil {
				return "", fmt.Errorf("compute sighash for input %d: %w", i, err)
			}

			sig, err := cipher.SignHash(sigHash, key)
			if err != nil {
				return "", fmt.Errorf("sign input %d: %w", i, err)
			}

			derSig := sigToDER(sig)
			derSig = append(derSig, 0x01) // SIGHASH_ALL

			signedInputs = append(signedInputs, inputSigning{
				txid:      input.TxID,
				vout:      input.Vout,
				scriptSig: buildP2PKHScriptSig(derSig, pubKey[:]),
			})
		}
	}

	// Serialize the final signed transaction
	var buf bytes.Buffer

	// Version
	writeUint32LE(&buf, 1)

	// Segwit marker and flag (if any segwit inputs)
	if hasSegwit {
		buf.WriteByte(0x00) // marker
		buf.WriteByte(0x01) // flag
	}

	// Input count
	writeVarInt(&buf, uint64(len(signedInputs)))

	// Inputs
	for _, in := range signedInputs {
		txidBytes, err := hex.DecodeString(in.txid)
		if err != nil {
			return "", fmt.Errorf("invalid txid hex %q: %w", in.txid, err)
		}
		reverseBytes(txidBytes)
		buf.Write(txidBytes)
		writeUint32LE(&buf, in.vout)
		if in.scriptSig != nil {
			writeVarInt(&buf, uint64(len(in.scriptSig)))
			buf.Write(in.scriptSig)
		} else {
			writeVarInt(&buf, 0) // empty scriptSig for segwit
		}
		writeUint32LE(&buf, 0xffffffff) // sequence
	}

	// Output count
	writeVarInt(&buf, uint64(len(outputs)))

	// Outputs
	for _, out := range outputs {
		writeUint64LE(&buf, uint64(out.value)) //nolint:gosec // satoshi amounts are non-negative
		writeVarInt(&buf, uint64(len(out.script)))
		buf.Write(out.script)
	}

	// Witness data (if segwit)
	if hasSegwit {
		for _, in := range signedInputs {
			if in.witness != nil {
				writeVarInt(&buf, uint64(len(in.witness)))
				for _, item := range in.witness {
					writeVarInt(&buf, uint64(len(item)))
					buf.Write(item)
				}
			} else {
				writeVarInt(&buf, 0) // empty witness for non-segwit inputs
			}
		}
	}

	// Locktime
	writeUint32LE(&buf, 0)

	return hex.EncodeToString(buf.Bytes()), nil
}

type txOut struct {
	value  int64
	script []byte
}

// outputScript creates the appropriate output script for a Bitcoin address.
// Supports both P2PKH (1...) and P2WPKH (bc1q...) addresses.
func outputScript(address string) ([]byte, error) {
	if strings.HasPrefix(strings.ToLower(address), "bc1") || strings.HasPrefix(strings.ToLower(address), "tb1") {
		return p2wpkhScript(address)
	}
	return p2pkhScript(address)
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

// p2wpkhScript creates a P2WPKH output script for a bech32 segwit address
func p2wpkhScript(address string) ([]byte, error) {
	segAddr, err := cipher.DecodeBech32BitcoinAddress(address)
	if err != nil {
		return nil, fmt.Errorf("invalid segwit address: %w", err)
	}

	// OP_0 <20-byte witness program>
	script := make([]byte, 22)
	script[0] = 0x00 // OP_0 (witness version 0)
	script[1] = 0x14 // Push 20 bytes
	copy(script[2:22], segAddr.Key[:])
	return script, nil
}

// isSegwitAddress returns true if the address is a bech32 segwit address
func isSegwitAddress(address string) bool {
	lower := strings.ToLower(address)
	return strings.HasPrefix(lower, "bc1") || strings.HasPrefix(lower, "tb1")
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
		writeUint64LE(&buf, uint64(out.value)) //nolint:gosec // satoshi amounts are non-negative
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

// computeSegwitSigHash computes the BIP143 signature hash for a P2WPKH input.
func computeSegwitSigHash(inputs []UTXO, outputs []txOut, sigIndex int, inputValue int64, pubKey cipher.PubKey) (cipher.SHA256, error) {
	// BIP143 defines a new digest algorithm for segwit:
	// Double SHA256 of the serialization of:
	//  1. nVersion (4 bytes LE)
	//  2. hashPrevouts (32 bytes)
	//  3. hashSequence (32 bytes)
	//  4. outpoint (32+4 bytes)
	//  5. scriptCode (varint + script)
	//  6. value (8 bytes LE)
	//  7. nSequence (4 bytes LE)
	//  8. hashOutputs (32 bytes)
	//  9. nLockTime (4 bytes LE)
	// 10. nHashType (4 bytes LE)

	// hashPrevouts = SHA256(SHA256(all input outpoints))
	var prevoutsData bytes.Buffer
	for _, input := range inputs {
		txidBytes, err := hex.DecodeString(input.TxID)
		if err != nil {
			return cipher.SHA256{}, fmt.Errorf("decode txid: %w", err)
		}
		reverseBytes(txidBytes)
		prevoutsData.Write(txidBytes)
		writeUint32LE(&prevoutsData, input.Vout)
	}
	hashPrevouts := doubleSHA256(prevoutsData.Bytes())

	// hashSequence = SHA256(SHA256(all input sequences))
	var seqData bytes.Buffer
	for range inputs {
		writeUint32LE(&seqData, 0xffffffff)
	}
	hashSequence := doubleSHA256(seqData.Bytes())

	// hashOutputs = SHA256(SHA256(all outputs))
	var outputsData bytes.Buffer
	for _, out := range outputs {
		writeUint64LE(&outputsData, uint64(out.value)) //nolint:gosec // satoshi amounts are non-negative
		writeVarInt(&outputsData, uint64(len(out.script)))
		outputsData.Write(out.script)
	}
	hashOutputs := doubleSHA256(outputsData.Bytes())

	// scriptCode for P2WPKH is OP_DUP OP_HASH160 <20-byte-key-hash> OP_EQUALVERIFY OP_CHECKSIG
	keyHash := cipher.BitcoinPubKeyRipemd160(pubKey)
	scriptCode := make([]byte, 25)
	scriptCode[0] = 0x76
	scriptCode[1] = 0xa9
	scriptCode[2] = 0x14
	copy(scriptCode[3:23], keyHash[:])
	scriptCode[23] = 0x88
	scriptCode[24] = 0xac

	// Build the preimage
	var buf bytes.Buffer

	// 1. nVersion
	writeUint32LE(&buf, 1)

	// 2. hashPrevouts
	buf.Write(hashPrevouts[:])

	// 3. hashSequence
	buf.Write(hashSequence[:])

	// 4. outpoint (the input being signed)
	txidBytes, err := hex.DecodeString(inputs[sigIndex].TxID)
	if err != nil {
		return cipher.SHA256{}, fmt.Errorf("decode txid: %w", err)
	}
	reverseBytes(txidBytes)
	buf.Write(txidBytes)
	writeUint32LE(&buf, inputs[sigIndex].Vout)

	// 5. scriptCode
	writeVarInt(&buf, uint64(len(scriptCode)))
	buf.Write(scriptCode)

	// 6. value of the input
	writeUint64LE(&buf, uint64(inputValue)) //nolint:gosec // inputValue validated non-negative before call

	// 7. nSequence
	writeUint32LE(&buf, 0xffffffff)

	// 8. hashOutputs
	buf.Write(hashOutputs[:])

	// 9. nLockTime
	writeUint32LE(&buf, 0)

	// 10. nHashType (SIGHASH_ALL)
	writeUint32LE(&buf, 1)

	result := doubleSHA256(buf.Bytes())
	return cipher.SHA256(result), nil
}

func doubleSHA256(data []byte) [32]byte {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
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
