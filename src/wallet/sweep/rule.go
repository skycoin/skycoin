// Package sweep src/wallet/sweep/rule.go c1-coin-wallet
package sweep

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
)

// RulesFile is the file in the wallet directory that holds the sweep rules. It
// holds no secrets.
const RulesFile = "sweep.json"

// DefaultGap is how many unused addresses are kept generated past the last used
// one on the external chain of a shared account, so a payment to an address the
// wallet has not generated yet is still seen.
const DefaultGap = 20

// Destination kinds
const (
	// DestAccount is the next unused address of another account of the same bip44
	// wallet. It is the default for a bip44 wallet.
	DestAccount = "account"
	// DestWallet is the next unused address of a wallet, the same one for a
	// deterministic wallet. It is the default for other wallets.
	DestWallet = "wallet"
	// DestAddress is a fixed address.
	DestAddress = "address"
)

// Destination says where swept coins go.
type Destination struct {
	Kind string `json:"kind"`
	// Account is the destination account for DestAccount. Nil picks the first
	// account after the shared ones and creates it when needed.
	Account *uint32 `json:"account,omitempty"`
	// Wallet is the destination wallet for DestWallet. Empty means the wallet itself.
	Wallet string `json:"wallet,omitempty"`
	// Address is the destination for DestAddress.
	Address string `json:"address,omitempty"`
}

// Rule is the sweep rule of one wallet.
type Rule struct {
	Wallet string `json:"wallet"`
	// SharedAccounts are bip44 accounts whose xpub was shared. Every address of
	// such an account is exposed, including those that never sent.
	SharedAccounts []uint32 `json:"shared_accounts,omitempty"`
	// Addresses are added by hand. Each must have sent, or belong to a shared
	// account, because sweeping a single address that never sent would make it
	// send for the first time.
	Addresses   []string    `json:"addresses,omitempty"`
	Destination Destination `json:"destination"`
	Enabled     bool        `json:"enabled"`
	// Gap overrides DefaultGap when not zero.
	Gap uint32 `json:"gap,omitempty"`
}

func (r Rule) gap() uint32 {
	if r.Gap == 0 {
		return DefaultGap
	}
	return r.Gap
}

func (r Rule) shared() map[uint32]bool {
	m := make(map[uint32]bool, len(r.SharedAccounts))
	for _, a := range r.SharedAccounts {
		m[a] = true
	}
	return m
}

func (r Rule) maxShared() (uint32, bool) {
	var maxIdx uint32
	for _, a := range r.SharedAccounts {
		maxIdx = max(maxIdx, a)
	}
	return maxIdx, len(r.SharedAccounts) > 0
}

// check validates what can be checked without the wallet or the blockchain.
func (r Rule) check() error {
	if r.Wallet == "" {
		return errors.New("wallet is required")
	}
	seen := make(map[uint32]bool)
	for _, a := range r.SharedAccounts {
		if seen[a] {
			return fmt.Errorf("account %d is listed twice", a)
		}
		seen[a] = true
	}
	for _, s := range r.Addresses {
		if _, err := cipher.DecodeBase58Address(s); err != nil {
			return fmt.Errorf("invalid address %q: %w", s, err)
		}
	}
	d := r.Destination
	switch d.Kind {
	case DestAccount:
		if d.Wallet != "" || d.Address != "" {
			return errors.New("an account destination takes only an account")
		}
		if d.Account != nil {
			if seen[*d.Account] {
				return fmt.Errorf("destination account %d is shared", *d.Account)
			}
			if m, ok := r.maxShared(); ok && *d.Account <= m {
				return fmt.Errorf("destination account %d must come after shared account %d", *d.Account, m)
			}
		}
	case DestWallet:
		if d.Account != nil || d.Address != "" {
			return errors.New("a wallet destination takes only a wallet")
		}
	case DestAddress:
		if d.Account != nil || d.Wallet != "" {
			return errors.New("an address destination takes only an address")
		}
		if _, err := cipher.DecodeBase58Address(d.Address); err != nil {
			return fmt.Errorf("invalid destination address %q: %w", d.Address, err)
		}
	default:
		return fmt.Errorf("unknown destination kind %q", d.Kind)
	}
	return nil
}

func loadRules(dir string) (map[string]Rule, error) {
	b, err := os.ReadFile(filepath.Join(dir, RulesFile)) //nolint:gosec // the configured wallet directory
	if errors.Is(err, os.ErrNotExist) {
		return map[string]Rule{}, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Rule
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, fmt.Errorf("%s: %w", RulesFile, err)
	}
	rules := make(map[string]Rule, len(list))
	for _, r := range list {
		rules[r.Wallet] = r
	}
	return rules, nil
}

func saveRules(dir string, rules map[string]Rule) error {
	list := make([]Rule, 0, len(rules))
	for _, r := range rules {
		list = append(list, r)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Wallet < list[j].Wallet })
	b, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, RulesFile+".*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name()) //nolint:errcheck
	if _, err := f.Write(append(b, '\n')); err != nil {
		f.Close() //nolint:errcheck,gosec
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), filepath.Join(dir, RulesFile))
}
