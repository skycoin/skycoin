// Package plan src/wallet/sweep/plan/build.go c1-coin-wallet
package plan

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

// Output is a confirmed unspent output. Hours are its coin hours at the head block.
type Output struct {
	Hash    cipher.SHA256
	Address cipher.Address
	Coins   uint64
	Hours   uint64
	// Contested says another unconfirmed transaction spends it.
	Contested bool
}

// Sweep is one transaction of a plan: Txn spends Inputs, in order, and pays
// every coin to one address with zero coin hours, so Fee is all the input hours.
type Sweep struct {
	Txn    *coin.Transaction
	Inputs []Output
	Fee    uint64
}

// Group is the outputs that go to one target.
type Group struct {
	Target  int
	Outputs []Output
	Coins   uint64
	Hours   uint64
}

// Pass is what one sweep of a wallet finds.
type Pass struct {
	Groups []Group
	// Exposed counts the addresses r sweeps, Uncovered those it leaves alone
	// because no mapping covers them and level 1 is off.
	Exposed, Uncovered int
	// Coins and Hours wait on exposed addresses, Contested counts the waiting
	// outputs another unconfirmed transaction spends.
	Coins, Hours uint64
	Contested    int
}

// Find sorts the unspent outputs of addrs by target, leaving out those in
// covered, which a sweep already spends.
func (r Rule) Find(addrs []Addr, sent map[cipher.Address]bool, unspent []Output, covered map[cipher.SHA256]bool) Pass {
	routes, exposed, uncovered := r.Route(addrs, sent)
	p := Pass{Exposed: exposed, Uncovered: uncovered}
	idx := make(map[int]int)
	for _, u := range unspent {
		t, ok := routes[u.Address]
		if !ok || covered[u.Hash] {
			continue
		}
		i, ok := idx[t]
		if !ok {
			i = len(p.Groups)
			idx[t] = i
			p.Groups = append(p.Groups, Group{Target: t})
		}
		g := &p.Groups[i]
		g.Outputs = append(g.Outputs, u)
		g.Coins += u.Coins
		g.Hours += u.Hours
		p.Coins += u.Coins
		p.Hours += u.Hours
		if u.Contested {
			p.Contested++
		}
	}
	sort.Slice(p.Groups, func(i, j int) bool { return p.Groups[i].Target < p.Groups[j].Target })
	return p
}

// Build makes the unsigned transactions that sweep outs to to.
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
func Build(outs []Output, to cipher.Address, maxSize uint32) ([]Sweep, error) {
	var withHours, without []Output
	for _, u := range outs {
		if u.Hours == 0 {
			without = append(without, u)
		} else {
			withHours = append(withHours, u)
		}
	}
	if len(withHours) == 0 {
		return nil, ErrNoHours
	}
	sort.Slice(withHours, func(i, j int) bool {
		if withHours[i].Hours != withHours[j].Hours {
			return withHours[i].Hours > withHours[j].Hours
		}
		return bytes.Compare(withHours[i].Hash[:], withHours[j].Hash[:]) < 0
	})

	sizes := sizeTable{to: to}
	var groups [][]Output
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
			sum += withHours[k-1].Hours
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
		in := append([]Output(nil), g...)
		var coins, hours uint64
		for _, u := range in {
			var carry uint64
			coins, carry = bits.Add64(coins, u.Coins, 0)
			if carry != 0 {
				return nil, errors.New("coins overflow")
			}
			hours, carry = bits.Add64(hours, u.Hours, 0)
			if carry != 0 {
				return nil, errors.New("coin hours overflow")
			}
		}
		hashes := make([]cipher.SHA256, len(in))
		for i, u := range in {
			hashes[i] = u.Hash
		}
		txn, err := build(hashes, to, coins)
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
func build(in []cipher.SHA256, to cipher.Address, coins uint64) (*coin.Transaction, error) {
	txn := &coin.Transaction{}
	for _, h := range in {
		if err := txn.PushInput(h); err != nil {
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
		in := make([]cipher.SHA256, len(s.sizes)+1)
		for i := range in {
			in[i] = cipher.SumSHA256([]byte{byte(i), byte(i >> 8), byte(i >> 16)})
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
