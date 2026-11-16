// Package plan src/wallet/sweep/plan/exposure.go c1-coin-wallet
package plan

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

// Addr is an address of a wallet and where it sits in a bip44 wallet.
type Addr struct {
	Address cipher.Address
	Account uint32
	Change  bool
	Child   uint32
	Bip44   bool
}

// listed reports whether r names a: in a shared account, added by hand, or the
// source of a mapping.
func (r Rule) listed(a Addr, shared, sources map[uint32]bool, addrs map[cipher.Address]bool) bool {
	return (a.Bip44 && (shared[a.Account] || sources[a.Account])) || addrs[a.Address]
}

func (r Rule) namedAddrs() map[cipher.Address]bool {
	m := make(map[cipher.Address]bool)
	for _, s := range r.Addresses {
		if a, err := cipher.DecodeBase58Address(s); err == nil {
			m[a] = true
		}
	}
	for _, mp := range r.Mappings {
		if a, err := cipher.DecodeBase58Address(mp.From.Address); err == nil {
			m[a] = true
		}
	}
	return m
}

// Exposed returns the addresses of addrs that r sweeps from. An address is
// exposed when it belongs to a shared account or has sent. It is swept when r
// lists it, or when r sweeps every address that has sent. A single address that
// never sent is never included, and neither is a post-quantum address, which a
// public key does not expose.
func (r Rule) Exposed(addrs []Addr, sent map[cipher.Address]bool) map[cipher.Address]bool {
	shared, sources, named := r.Shared(), r.sourceAccounts(), r.namedAddrs()
	all := r.SweepsSpent()
	out := make(map[cipher.Address]bool)
	for _, a := range addrs {
		if a.Address.Version != cipher.AddressVersionSecp256k1 {
			continue
		}
		inShared := a.Bip44 && shared[a.Account]
		if !inShared && !sent[a.Address] {
			continue
		}
		if inShared || all || r.listed(a, shared, sources, named) {
			out[a.Address] = true
		}
	}
	return out
}

// Route gives each address r sweeps the target that sends it: the mapping of
// the address, else the mapping of its account, else Level1 when r is Enabled.
// Exposed addresses with no target are left alone and counted in uncovered.
func (r Rule) Route(addrs []Addr, sent map[cipher.Address]bool) (routes map[cipher.Address]int, exposed, uncovered int) {
	byAddr := make(map[cipher.Address]int)
	byAcct := make(map[uint32]int)
	for i, m := range r.Mappings {
		if m.From.Account != nil {
			byAcct[*m.From.Account] = i
		} else if a, err := cipher.DecodeBase58Address(m.From.Address); err == nil {
			byAddr[a] = i
		}
	}
	exp := r.Exposed(addrs, sent)
	routes = make(map[cipher.Address]int, len(exp))
	for _, a := range addrs {
		if !exp[a.Address] {
			continue
		}
		if _, done := routes[a.Address]; done {
			continue
		}
		if t, ok := byAddr[a.Address]; ok {
			routes[a.Address] = t
		} else if t, ok := byAcct[a.Account]; ok && a.Bip44 {
			routes[a.Address] = t
		} else if r.Enabled {
			routes[a.Address] = Level1
		} else {
			uncovered++
		}
	}
	return routes, len(exp), uncovered
}

// CheckListed refuses an address listed by hand or as a mapping source that is
// not in the wallet, or that would be swept for the first time without ever
// having sent, and a source account the wallet does not have.
func (r Rule) CheckListed(addrs []Addr, sent map[cipher.Address]bool, accounts uint32) error {
	shared := r.Shared()
	byAddr := make(map[cipher.Address]Addr, len(addrs))
	for _, a := range addrs {
		byAddr[a.Address] = a
	}
	check := func(s string) error {
		a, err := cipher.DecodeBase58Address(s)
		if err != nil {
			return err
		}
		wa, ok := byAddr[a]
		if !ok {
			return fmt.Errorf("address %s is not in wallet %s", s, r.Wallet)
		}
		if !sent[a] && !(wa.Bip44 && shared[wa.Account]) {
			return fmt.Errorf("address %s has never sent, so its public key is not exposed. Sweeping it would send from it for the first time", s)
		}
		return nil
	}
	for _, s := range r.Addresses {
		if err := check(s); err != nil {
			return err
		}
	}
	for _, m := range r.Mappings {
		if m.From.Address != "" {
			if err := check(m.From.Address); err != nil {
				return err
			}
		}
		if m.From.Account != nil && *m.From.Account >= accounts {
			return fmt.Errorf("wallet %s has no account %d", r.Wallet, *m.From.Account)
		}
	}
	for _, a := range r.SharedAccounts {
		if a >= accounts {
			return fmt.Errorf("wallet %s has no account %d", r.Wallet, a)
		}
	}
	return nil
}

// GapNeeds returns how many addresses to generate on the external chain of each
// shared account, so gap unused ones follow the last that appears on chain.
func (r Rule) GapNeeds(addrs []Addr, seen map[cipher.Address]bool) map[uint32]uint64 {
	type tip struct {
		n        uint32
		lastUsed int64
	}
	tips := make(map[uint32]*tip)
	for _, acct := range r.SharedAccounts {
		tips[acct] = &tip{lastUsed: -1}
	}
	for _, a := range addrs {
		t, ok := tips[a.Account]
		if !ok || !a.Bip44 || a.Change {
			continue
		}
		t.n = max(t.n, a.Child+1)
		if seen[a.Address] {
			t.lastUsed = max(t.lastUsed, int64(a.Child))
		}
	}
	needs := make(map[uint32]uint64)
	for acct, t := range tips {
		want := t.lastUsed + 1 + int64(r.GapLen())
		if want > int64(t.n) {
			needs[acct] = uint64(want - int64(t.n)) //nolint:gosec // positive
		}
	}
	return needs
}

// ErrDestinationExposed is returned for a destination that is itself exposed.
var ErrDestinationExposed = errors.New("the destination is exposed")

// CheckDestination refuses a destination that has sent or that is swept. The
// exposed set holds every generated address of a shared account, both chains.
func CheckDestination(dest cipher.Address, exp, sent map[cipher.Address]bool) error {
	if dest.Version == cipher.AddressVersionPQ {
		return nil
	}
	if dest.Null() {
		return errors.New("the destination is the null address")
	}
	if exp[dest] {
		return fmt.Errorf("%w: %s is swept by this rule", ErrDestinationExposed, dest)
	}
	if sent[dest] {
		return fmt.Errorf("%w: %s has sent, so its public key is on the blockchain", ErrDestinationExposed, dest)
	}
	return nil
}
