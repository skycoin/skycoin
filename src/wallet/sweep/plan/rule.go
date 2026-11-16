// Package plan src/wallet/sweep/plan/rule.go c1-coin-wallet
//
// Package plan decides what a sweep moves and where, and builds the unsigned
// transactions that move it. It holds no wallet and reads no blockchain, so the
// node's sweep service and the browser cipher run the same code.
//
// An address is exposed when a quantum attacker could compute its private key:
// it has spent, so its public key is on the blockchain, or it belongs to a bip44
// account whose xpub was shared, so every child key follows from that one key.
package plan

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

// DefaultGap is how many unused addresses are kept generated past the last used
// one on the external chain of a shared account, so a payment to an address the
// wallet has not generated yet is still seen.
const DefaultGap = 20

// Destination kinds
const (
	// DestAccount is the next unused address of an account of the same bip44
	// wallet. It is the default for a bip44 wallet.
	DestAccount = "account"
	// DestWallet is the next unused address of a wallet, the same one when no
	// wallet is named. It is the default for other wallets.
	DestWallet = "wallet"
	// DestAddress is a fixed address.
	DestAddress = "address"
)

// Destination says where swept coins go.
type Destination struct {
	Kind string `json:"kind"`
	// Account is the destination account for DestAccount. Nil picks the first
	// account after the shared and source accounts, created when needed.
	Account *uint32 `json:"account,omitempty"`
	// Wallet is the destination wallet for DestWallet. Empty means the wallet itself.
	Wallet string `json:"wallet,omitempty"`
	// Address is the destination for DestAddress.
	Address string `json:"address,omitempty"`
}

// Source is what a mapping sweeps from: one address, or the exposed addresses
// of one bip44 account.
type Source struct {
	Address string  `json:"address,omitempty"`
	Account *uint32 `json:"account,omitempty"`
}

// Mapping sends what its source holds to its own destination.
type Mapping struct {
	From Source      `json:"from"`
	To   Destination `json:"to"`
}

// Rule is the sweep rule of one wallet. It has two levels. With Enabled set,
// every exposed address goes to Destination. Mappings send their sources to
// their own destinations, whether Enabled is set or not, and take precedence.
type Rule struct {
	Wallet  string `json:"wallet"`
	Enabled bool   `json:"enabled"`
	// Destination is where Enabled sends what no mapping covers.
	Destination Destination `json:"destination"`
	// SharedAccounts are bip44 accounts whose xpub was shared. Every address of
	// such an account is exposed, including those that never sent.
	SharedAccounts []uint32 `json:"shared_accounts,omitempty"`
	// Addresses are added by hand. Each must have sent, or belong to a shared
	// account, because sweeping a single address that never sent would make it
	// send for the first time.
	Addresses []string `json:"addresses,omitempty"`
	// IncludeSpent, when not false, also sweeps every address that has sent
	// without being listed. False sweeps only shared accounts, Addresses and
	// mapping sources.
	IncludeSpent *bool     `json:"include_spent,omitempty"`
	Mappings     []Mapping `json:"mappings,omitempty"`
	// Gap overrides DefaultGap when not zero.
	Gap uint32 `json:"gap,omitempty"`
}

// Level1 is the target Route gives an address that no mapping covers.
const Level1 = -1

// Active reports whether r sweeps anything.
func (r Rule) Active() bool {
	return r.Enabled || len(r.Mappings) > 0
}

// SweepsSpent reports whether r sweeps addresses that have sent without being listed.
func (r Rule) SweepsSpent() bool {
	return r.IncludeSpent == nil || *r.IncludeSpent
}

// GapLen is how many unused addresses r keeps past the last used one.
func (r Rule) GapLen() uint32 {
	if r.Gap == 0 {
		return DefaultGap
	}
	return r.Gap
}

// Shared reports the shared accounts as a set.
func (r Rule) Shared() map[uint32]bool {
	m := make(map[uint32]bool, len(r.SharedAccounts))
	for _, a := range r.SharedAccounts {
		m[a] = true
	}
	return m
}

func (r Rule) sourceAccounts() map[uint32]bool {
	m := make(map[uint32]bool)
	for _, mp := range r.Mappings {
		if mp.From.Account != nil {
			m[*mp.From.Account] = true
		}
	}
	return m
}

// Target returns the destination of a target Route gave.
func (r Rule) Target(t int) Destination {
	if t == Level1 {
		return r.Destination
	}
	return r.Mappings[t].To
}

// NeedsBip44 reports whether r names accounts, which only a bip44 wallet has.
func (r Rule) NeedsBip44() bool {
	if len(r.SharedAccounts) > 0 || len(r.sourceAccounts()) > 0 || (r.Enabled && r.Destination.Kind == DestAccount) {
		return true
	}
	for _, m := range r.Mappings {
		if m.To.Kind == DestAccount {
			return true
		}
	}
	return false
}

// DestinationAccount is the account d sends to: the one it names, or the first
// after every shared and source account.
func (r Rule) DestinationAccount(d Destination) uint32 {
	if d.Account != nil {
		return *d.Account
	}
	var next uint32
	for a := range r.Shared() {
		next = max(next, a+1)
	}
	for a := range r.sourceAccounts() {
		next = max(next, a+1)
	}
	return next
}

// Check validates what can be checked without the wallet or the blockchain.
func (r Rule) Check() error {
	if r.Wallet == "" {
		return errors.New("wallet is required")
	}
	shared := make(map[uint32]bool)
	var maxShared uint32
	for _, a := range r.SharedAccounts {
		if shared[a] {
			return fmt.Errorf("account %d is listed twice", a)
		}
		shared[a] = true
		maxShared = max(maxShared, a)
	}
	for _, s := range r.Addresses {
		if _, err := cipher.DecodeBase58Address(s); err != nil {
			return fmt.Errorf("invalid address %q: %w", s, err)
		}
	}
	sources := r.sourceAccounts()
	checkDest := func(what string, d Destination) error {
		switch d.Kind {
		case DestAccount:
			if d.Wallet != "" || d.Address != "" {
				return fmt.Errorf("%s: an account destination takes only an account", what)
			}
			if d.Account == nil {
				return nil
			}
			a := *d.Account
			if shared[a] {
				return fmt.Errorf("%s: destination account %d is shared", what, a)
			}
			if len(shared) > 0 && a <= maxShared {
				return fmt.Errorf("%s: destination account %d must come after shared account %d", what, a, maxShared)
			}
			if sources[a] {
				return fmt.Errorf("%s: destination account %d is a source", what, a)
			}
		case DestWallet:
			if d.Account != nil || d.Address != "" {
				return fmt.Errorf("%s: a wallet destination takes only a wallet", what)
			}
		case DestAddress:
			if d.Account != nil || d.Wallet != "" {
				return fmt.Errorf("%s: an address destination takes only an address", what)
			}
			if _, err := cipher.DecodeBase58Address(d.Address); err != nil {
				return fmt.Errorf("%s: invalid destination address %q: %w", what, d.Address, err)
			}
		default:
			return fmt.Errorf("%s: unknown destination kind %q", what, d.Kind)
		}
		return nil
	}
	if err := checkDest("destination", r.Destination); err != nil {
		return err
	}
	seenAddr := make(map[string]bool)
	seenAcct := make(map[uint32]bool)
	for i, m := range r.Mappings {
		what := fmt.Sprintf("mapping %d", i+1)
		switch {
		case m.From.Address != "" && m.From.Account != nil:
			return fmt.Errorf("%s: a source is an address or an account, not both", what)
		case m.From.Address != "":
			if _, err := cipher.DecodeBase58Address(m.From.Address); err != nil {
				return fmt.Errorf("%s: invalid source address %q: %w", what, m.From.Address, err)
			}
			if seenAddr[m.From.Address] {
				return fmt.Errorf("%s: address %s is mapped twice", what, m.From.Address)
			}
			seenAddr[m.From.Address] = true
		case m.From.Account != nil:
			if seenAcct[*m.From.Account] {
				return fmt.Errorf("%s: account %d is mapped twice", what, *m.From.Account)
			}
			seenAcct[*m.From.Account] = true
		default:
			return fmt.Errorf("%s: the source is empty", what)
		}
		if err := checkDest(what, m.To); err != nil {
			return err
		}
	}
	return nil
}
