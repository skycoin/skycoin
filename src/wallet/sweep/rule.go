// Package sweep src/wallet/sweep/rule.go c1-coin-wallet
package sweep

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/skycoin/skycoin/src/wallet/sweep/plan"
)

// RulesFile is the file in the wallet directory that holds the sweep rules. It
// holds no secrets.
const RulesFile = "sweep.json"

// The rule types live in package plan, which the browser cipher shares.
type (
	// Rule is the sweep rule of one wallet.
	Rule = plan.Rule
	// Destination says where swept coins go.
	Destination = plan.Destination
	// Mapping sends what its source holds to its own destination.
	Mapping = plan.Mapping
	// Source is what a mapping sweeps from.
	Source = plan.Source
)

// Destination kinds
const (
	DestAccount = plan.DestAccount
	DestWallet  = plan.DestWallet
	DestAddress = plan.DestAddress
)

// DefaultGap is plan.DefaultGap.
const DefaultGap = plan.DefaultGap

// ErrDestinationExposed is returned for a destination that is itself exposed.
var ErrDestinationExposed = plan.ErrDestinationExposed

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
