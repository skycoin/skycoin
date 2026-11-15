// Package transaction src/transaction/split.go c1-coin-chain
package transaction

import (
	"errors"
	"fmt"
	"math/bits"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/mathutil"
)

// TooLargeError is returned when a transaction would be over the size limit once signed.
type TooLargeError struct {
	Inputs int
	Size   uint64
	Limit  uint32
}

func (e TooLargeError) Error() string {
	return fmt.Sprintf("the transaction spends %d outputs and would be %d bytes signed, more than the %d bytes allowed", e.Inputs, e.Size, e.Limit)
}

// CheckSize returns a TooLargeError when txn, spending inputs, would be over the
// size limit once signed. Each post-quantum input adds its witness when signed.
func CheckSize(txn *coin.Transaction, inputs []UxBalance) error {
	size, err := txn.Size()
	if err != nil {
		return err
	}
	var nPQ int
	for _, u := range inputs {
		if u.Address.Version == cipher.AddressVersionPQ {
			nPQ++
		}
	}
	total := uint64(size) + uint64(coin.PQWitnessSize(len(txn.In), nPQ)) //nolint:gosec // sizes are non-negative
	limit := params.UserVerifyTxn.MaxTransactionSize
	if total > uint64(limit) {
		return NewError(TooLargeError{Inputs: len(txn.In), Size: total, Limit: limit})
	}
	return nil
}

// BuildFunc builds one transaction that pays to from some of uxb. It returns the
// transaction and the outputs it spends, in input order.
type BuildFunc func(to []coin.TransactionOutput, uxb []UxBalance) (*coin.Transaction, []UxBalance, error)

// Split pays to with as many transactions as it takes to keep each one under the
// size limit. No two parts spend the same output and no part spends the change of
// another, so every part can be injected at once, in any order.
//
// A recipient whose amount does not fit in one part is paid across two. In manual
// hours mode its hours are divided in proportion to the coins of each piece.
func Split(to []coin.TransactionOutput, uxb []UxBalance, build BuildFunc) ([]*coin.Transaction, [][]UxBalance, error) {
	avail := append([]UxBalance(nil), uxb...)
	sortSpendsCoinsHighToLow(avail)
	remaining := append([]coin.TransactionOutput(nil), to...)

	var txns []*coin.Transaction
	var inputs [][]UxBalance
	n := len(avail)
	for len(remaining) > 0 {
		if len(avail) == 0 {
			return nil, nil, ErrInsufficientBalance
		}
		n = min(n, len(avail))
		chunk := avail[:n]

		var chunkCoins uint64
		for _, u := range chunk {
			var err error
			if chunkCoins, err = mathutil.AddUint64(chunkCoins, u.Coins); err != nil {
				return nil, nil, err
			}
		}
		var availCoins uint64
		for _, u := range avail {
			var err error
			if availCoins, err = mathutil.AddUint64(availCoins, u.Coins); err != nil {
				return nil, nil, err
			}
		}
		// Leave the smallest output of a full chunk unpaid, so the part has
		// change to keep its unsent hours instead of giving them all away.
		// A send of the whole balance has no coins to spare for that.
		payable := chunkCoins
		want := totalCoins(remaining)
		if n > 1 && want > chunkCoins && availCoins-chunk[n-1].Coins >= want {
			payable -= chunk[n-1].Coins
		}
		part, rest := takeCoins(remaining, payable)

		txn, used, err := build(part, chunk)
		var tl TooLargeError
		if errors.As(err, &tl) && n > 1 {
			m := int(uint64(tl.Inputs) * uint64(tl.Limit) / tl.Size) //nolint:gosec // bounded by n
			n = max(1, min(n-1, m))
			continue
		}
		if err != nil {
			return nil, nil, err
		}

		txns = append(txns, txn)
		inputs = append(inputs, used)
		avail = uxBalancesSub(avail, used)
		remaining = rest
	}

	return txns, inputs, nil
}

func totalCoins(to []coin.TransactionOutput) uint64 {
	var t uint64
	for _, o := range to {
		t += o.Coins
	}
	return t
}

// takeCoins splits to into the outputs paying at most limit coins and the rest.
func takeCoins(to []coin.TransactionOutput, limit uint64) (part, rest []coin.TransactionOutput) {
	for i, o := range to {
		if o.Coins <= limit {
			part = append(part, o)
			limit -= o.Coins
			continue
		}
		if limit > 0 {
			// o.Hours*limit/o.Coins cannot overflow, since limit < o.Coins.
			hi, lo := bits.Mul64(o.Hours, limit)
			hours, _ := bits.Div64(hi, lo, o.Coins)
			part = append(part, coin.TransactionOutput{Address: o.Address, Coins: limit, Hours: hours})
			o.Coins -= limit
			o.Hours -= hours
		}
		rest = append(rest, o)
		rest = append(rest, to[i+1:]...)
		return part, rest
	}
	return part, nil
}

// CreateSplit is Create, except that when one transaction would be too large it
// pays the outputs with several transactions made by Split.
func CreateSplit(p Params, auxs coin.AddressUxOuts, headTime uint64) ([]*coin.Transaction, [][]UxBalance, error) {
	txn, inputs, err := Create(p, auxs, headTime)
	if err == nil {
		return []*coin.Transaction{txn}, [][]UxBalance{inputs}, nil
	}
	if !errors.As(err, new(TooLargeError)) {
		return nil, nil, err
	}

	flat := auxs.Flatten()
	uxb, err := NewUxBalances(flat, headTime)
	if err != nil {
		return nil, nil, err
	}
	byHash := make(map[cipher.SHA256]coin.UxOut, len(flat))
	for _, u := range flat {
		byHash[u.Hash()] = u
	}
	return Split(p.To, uxb, func(to []coin.TransactionOutput, chunk []UxBalance) (*coin.Transaction, []UxBalance, error) {
		pp := p
		pp.To = to
		ux := make(coin.UxArray, len(chunk))
		for i, u := range chunk {
			ux[i] = byHash[u.Hash]
		}
		return Create(pp, coin.NewAddressUxOuts(ux), headTime)
	})
}
