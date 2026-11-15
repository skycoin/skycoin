// Package sweep src/wallet/sweep/plan.go c1-coin-wallet
//
// Package sweep moves coins off exposed addresses as soon as they confirm.
//
// An address is exposed when a quantum attacker could compute its private key:
// it has spent, so its public key is on the blockchain, or it belongs to a bip44
// account whose xpub was shared, so every child key follows from that one key.
package sweep

import (
	"bytes"
	"errors"
	"math/bits"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

// ErrNoHours is returned when no output to sweep has coin hours. A transaction
// must burn at least one coin hour, so those outputs wait until they accrue some.
var ErrNoHours = errors.New("the outputs to sweep have no coin hours to burn yet")

// Sweep is one transaction of a plan: Txn spends Inputs, in order, and pays
// every coin to one address with zero coin hours, so Fee is all the input hours.
type Sweep struct {
	Txn    *coin.Transaction
	Inputs coin.UxArray
	Fee    uint64
}

// Plan makes the unsigned transactions that sweep uxs to to.
//
// Each transaction burns every coin hour of its inputs, the most the fee rules
// allow (src/util/fee VerifyTransactionFeeForHours only sets a minimum). The
// block publisher orders transactions by fee per kB, then by lowest hash
// (coin.SortTransactions), so a competing spend of the same outputs that is
// paid only from their own hours can at best tie a sweep.
//
// Inputs are grouped to keep that true. Sorted by hours, the first transaction
// takes the prefix with the highest fee per byte, so no subset of its inputs
// pays more per byte than the whole. The rest are planned the same way. Outputs
// without hours raise no fee, so they join the last transaction while it fits.
func Plan(uxs coin.UxArray, headTime uint64, to cipher.Address, maxSize uint32) ([]Sweep, error) {
	type uxh struct {
		ux    coin.UxOut
		hash  cipher.SHA256
		hours uint64
	}
	var withHours, without []uxh
	for _, ux := range uxs {
		h, err := ux.CoinHours(headTime)
		if err != nil {
			return nil, err
		}
		u := uxh{ux: ux, hash: ux.Hash(), hours: h}
		if h == 0 {
			without = append(without, u)
		} else {
			withHours = append(withHours, u)
		}
	}
	if len(withHours) == 0 {
		return nil, ErrNoHours
	}
	sort.Slice(withHours, func(i, j int) bool {
		if withHours[i].hours != withHours[j].hours {
			return withHours[i].hours > withHours[j].hours
		}
		return bytes.Compare(withHours[i].hash[:], withHours[j].hash[:]) < 0
	})

	sizes := sizeTable{to: to}
	var groups [][]uxh
	for len(withHours) > 0 {
		var best, sum, bestHours uint64
		var bestSize uint32
		for k := 1; k <= len(withHours); k++ {
			size, err := sizes.size(k)
			if err != nil {
				return nil, err
			}
			if size > maxSize {
				break
			}
			sum += withHours[k-1].hours
			// sum/size >= bestHours/bestSize, compared as 128 bit products
			if best == 0 || !ratioLess(sum, size, bestHours, bestSize) {
				best, bestHours, bestSize = uint64(k), sum, size
			}
		}
		if best == 0 {
			return nil, errors.New("one input does not fit in a transaction")
		}
		groups = append(groups, withHours[:best])
		withHours = withHours[best:]
	}

	last := groups[len(groups)-1]
	for _, u := range without {
		size, err := sizes.size(len(last) + 1)
		if err != nil {
			return nil, err
		}
		if size > maxSize {
			break
		}
		last = append(last, u)
	}
	groups[len(groups)-1] = last

	plan := make([]Sweep, 0, len(groups))
	for _, g := range groups {
		in := make(coin.UxArray, len(g))
		var coins, hours uint64
		for i, u := range g {
			in[i] = u.ux
			var carry uint64
			coins, carry = bits.Add64(coins, u.ux.Body.Coins, 0)
			if carry != 0 {
				return nil, errors.New("coins overflow")
			}
			hours, carry = bits.Add64(hours, u.hours, 0)
			if carry != 0 {
				return nil, errors.New("coin hours overflow")
			}
		}
		txn, err := build(in, to, coins)
		if err != nil {
			return nil, err
		}
		plan = append(plan, Sweep{Txn: txn, Inputs: in, Fee: hours})
	}
	return plan, nil
}

// ratioLess reports whether a/as < b/bs.
func ratioLess(a uint64, as uint32, b uint64, bs uint32) bool {
	ah, al := bits.Mul64(a, uint64(bs))
	bh, bl := bits.Mul64(b, uint64(as))
	return ah < bh || (ah == bh && al < bl)
}

// build makes the unsigned transaction spending in to one output of coins with
// zero coin hours.
func build(in coin.UxArray, to cipher.Address, coins uint64) (*coin.Transaction, error) {
	txn := &coin.Transaction{}
	for _, ux := range in {
		if err := txn.PushInput(ux.Hash()); err != nil {
			return nil, err
		}
	}
	if err := txn.PushOutput(to, coins, 0); err != nil {
		return nil, err
	}
	txn.Sigs = make([]cipher.Sig, len(txn.In))
	if err := txn.UpdateHeader(); err != nil {
		return nil, err
	}
	return txn, nil
}

// sizeTable gives the size of a signed sweep with n inputs. Signatures have a
// fixed width, so it depends on n only.
type sizeTable struct {
	to    cipher.Address
	sizes []uint32
}

func (s *sizeTable) size(n int) (uint32, error) {
	for len(s.sizes) < n {
		in := make(coin.UxArray, len(s.sizes)+1)
		for i := range in {
			in[i].Body.SrcTransaction = cipher.SumSHA256([]byte{byte(i), byte(i >> 8), byte(i >> 16)})
		}
		txn, err := build(in, s.to, 1)
		if err != nil {
			return 0, err
		}
		size, err := txn.Size()
		if err != nil {
			return 0, err
		}
		s.sizes = append(s.sizes, size)
	}
	return s.sizes[n-1], nil
}
