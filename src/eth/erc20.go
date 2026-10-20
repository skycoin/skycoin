package eth

import (
	"errors"
	"math/big"
)

// ERC-20 function selectors, the first 4 bytes of the keccak of the signature.
// https://eips.ethereum.org/EIPS/eip-20
var (
	transferSelector  = []byte{0xa9, 0x05, 0x9c, 0xbb}
	balanceOfSelector = []byte{0x70, 0xa0, 0x82, 0x31}
)

// ERC20Transfer is the call data of transfer(to, amount).
func ERC20Transfer(to Address, amount *big.Int) ([]byte, error) {
	word, err := abiWord(amount)
	if err != nil {
		return nil, err
	}
	data := append([]byte{}, transferSelector...)
	data = append(data, leftPad32(to[:])...)
	return append(data, word...), nil
}

// ERC20BalanceOf is the call data of balanceOf(owner).
func ERC20BalanceOf(owner Address) []byte {
	return append(append([]byte{}, balanceOfSelector...), leftPad32(owner[:])...)
}

func abiWord(v *big.Int) ([]byte, error) {
	if v.Sign() < 0 || v.BitLen() > 256 {
		return nil, errors.New("amount does not fit a uint256")
	}
	return leftPad32(v.Bytes()), nil
}

func leftPad32(b []byte) []byte {
	return append(make([]byte, 32-len(b)), b...)
}
