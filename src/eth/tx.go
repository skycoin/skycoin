package eth

import (
	"errors"
	"math/big"

	"github.com/skycoin/skycoin/src/cipher"
)

// DynamicFeeTx is an EIP-1559 (type 2) transaction with an empty access list.
// https://eips.ethereum.org/EIPS/eip-1559
type DynamicFeeTx struct {
	ChainID              *big.Int
	Nonce                uint64
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
	Gas                  uint64
	To                   Address
	Value                *big.Int
	Data                 []byte
}

// SignedTx is a signed transaction, as sent with eth_sendRawTransaction.
type SignedTx struct {
	Raw  []byte
	Hash [32]byte
}

const dynamicFeeTxType = 0x02

func (tx DynamicFeeTx) fields() rlpList {
	return rlpList{
		rlpUint(tx.ChainID),
		rlpUint(new(big.Int).SetUint64(tx.Nonce)),
		rlpUint(tx.MaxPriorityFeePerGas),
		rlpUint(tx.MaxFeePerGas),
		rlpUint(new(big.Int).SetUint64(tx.Gas)),
		rlpString(tx.To[:]),
		rlpUint(tx.Value),
		rlpString(tx.Data),
		rlpList{},
	}
}

// SigningHash is keccak256(0x02 || rlp(fields)), the hash the sender signs.
func (tx DynamicFeeTx) SigningHash() [32]byte {
	return Keccak256([]byte{dynamicFeeTxType}, tx.fields().encode())
}

// Sign signs with sec. The signature's recovery id is the y parity.
func (tx DynamicFeeTx) Sign(sec cipher.SecKey) (SignedTx, error) {
	sig, err := tx.signature(sec)
	if err != nil {
		return SignedTx{}, err
	}
	fields := append(tx.fields(),
		rlpUint(big.NewInt(int64(sig[64]))),
		rlpUint(new(big.Int).SetBytes(sig[:32])),
		rlpUint(new(big.Int).SetBytes(sig[32:64])),
	)
	raw := append([]byte{dynamicFeeTxType}, fields.encode()...)
	return SignedTx{Raw: raw, Hash: Keccak256(raw)}, nil
}

func (tx DynamicFeeTx) signature(sec cipher.SecKey) (cipher.Sig, error) {
	h := tx.SigningHash()
	sig, err := cipher.SignHash(cipher.SHA256(h), sec)
	if err != nil {
		return cipher.Sig{}, err
	}
	if sig[64] > 1 {
		return cipher.Sig{}, errors.New("signature recovery id is not a y parity")
	}
	return sig, nil
}
