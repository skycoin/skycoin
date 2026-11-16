// Package liteclient src/skycoin-lite/liteclient/sweep.go c2-app-lite
package liteclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet/sweep/plan"
)

// DefaultMaxTransactionSize is the node's default largest user transaction
// (src/params UserVerifyTxn), used when a sweep request names none.
const DefaultMaxTransactionSize = 32768

// SweepAddress is an address of the swept wallet and where it sits in a bip44
// wallet. A fresh destination address of the same wallet is listed too, so a
// destination in a shared account is refused.
type SweepAddress struct {
	Address string `json:"address"`
	Account uint32 `json:"account,omitempty"`
	Change  bool   `json:"change,omitempty"`
	Child   uint32 `json:"child,omitempty"`
	Bip44   bool   `json:"bip44,omitempty"`
}

// SweepOutput is a confirmed unspent output. Coins are droplets and Hours are
// the coin hours at the head block, the node's calculated_hours.
type SweepOutput struct {
	Hash      string `json:"hash"`
	Address   string `json:"address"`
	Coins     uint64 `json:"coins"`
	Hours     uint64 `json:"hours"`
	Contested bool   `json:"contested,omitempty"`
}

// SweepRequest is what one sweep pass of a browser wallet knows.
type SweepRequest struct {
	Rule      plan.Rule      `json:"rule"`
	Addresses []SweepAddress `json:"addresses"`
	// Sent lists the addresses that spent in a confirmed block, destinations
	// included. Seen lists those that appear in one.
	Sent    []string      `json:"sent"`
	Seen    []string      `json:"seen"`
	Outputs []SweepOutput `json:"outputs"`
	// Covered are outputs a pending sweep of this page already spends.
	Covered []string `json:"covered"`
	// Destinations maps a target, "-1" for level 1 or a mapping index, to the
	// address it resolved to. A target with waiting coins and no entry is
	// returned in Need.
	Destinations map[string]string `json:"destinations"`
	// Accounts is how many bip44 accounts the wallet has.
	Accounts uint32 `json:"accounts"`
	MaxSize  uint32 `json:"max_size"`
}

// SweepInput is an input of a sweep. The caller adds its secret and hands the
// inputs and outputs to PrepareTransaction, in this order.
type SweepInput struct {
	Hash    string `json:"hash"`
	Address string `json:"address"`
}

// SweepTxOutput is the one output of a sweep, with zero coin hours.
type SweepTxOutput struct {
	Address string `json:"address"`
	Coins   uint64 `json:"coins"`
	Hours   uint64 `json:"hours"`
}

// SweepTransaction is one unsigned sweep.
type SweepTransaction struct {
	Destination string          `json:"destination"`
	Inputs      []SweepInput    `json:"inputs"`
	Outputs     []SweepTxOutput `json:"outputs"`
	Fee         uint64          `json:"fee"`
}

// SweepResult is what a pass found and the transactions it would send.
type SweepResult struct {
	Exposed   int    `json:"exposed"`
	Uncovered int    `json:"uncovered"`
	Coins     uint64 `json:"coins"`
	Hours     uint64 `json:"hours"`
	Contested int    `json:"contested"`
	// Gap is how many external addresses to add to each shared account before
	// sweeping, so a payment past the stored ones is seen.
	Gap map[string]uint64 `json:"gap,omitempty"`
	// Need lists the targets with waiting coins whose destination is not given.
	Need []int `json:"need,omitempty"`
	// NoHours says some waiting coins have no coin hours to burn yet.
	NoHours bool               `json:"no_hours,omitempty"`
	Sweeps  []SweepTransaction `json:"sweeps,omitempty"`
	// Errors are per target, such as an exposed destination.
	Errors []string `json:"errors,omitempty"`
}

func addrSet(list []string) (map[cipher.Address]bool, error) {
	m := make(map[cipher.Address]bool, len(list))
	for _, s := range list {
		a, err := cipher.DecodeBase58Address(s)
		if err != nil {
			return nil, fmt.Errorf("invalid address %q: %w", s, err)
		}
		m[a] = true
	}
	return m, nil
}

func (r SweepRequest) addrs() ([]plan.Addr, error) {
	out := make([]plan.Addr, len(r.Addresses))
	for i, a := range r.Addresses {
		addr, err := cipher.DecodeBase58Address(a.Address)
		if err != nil {
			return nil, fmt.Errorf("invalid address %q: %w", a.Address, err)
		}
		out[i] = plan.Addr{Address: addr, Account: a.Account, Change: a.Change, Child: a.Child, Bip44: a.Bip44}
	}
	return out, nil
}

// SweepCheck validates a rule before it is saved: its shape, that every listed
// address is in the wallet and exposed, and that no fixed destination is.
// The request needs Rule, Addresses, Sent with the fixed destinations, and Accounts.
func SweepCheck(body string) error {
	var r SweepRequest
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return err
	}
	if err := r.Rule.Check(); err != nil {
		return err
	}
	addrs, err := r.addrs()
	if err != nil {
		return err
	}
	sent, err := addrSet(r.Sent)
	if err != nil {
		return err
	}
	if err := r.Rule.CheckListed(addrs, sent, r.Accounts); err != nil {
		return err
	}
	exp := r.Rule.Exposed(addrs, sent)
	dests := []plan.Destination{r.Rule.Destination}
	for _, m := range r.Rule.Mappings {
		dests = append(dests, m.To)
	}
	for _, d := range dests {
		if d.Kind != plan.DestAddress {
			continue
		}
		if err := plan.CheckDestination(cipher.MustDecodeBase58Address(d.Address), exp, sent); err != nil {
			return err
		}
	}
	return nil
}

// SweepPlan runs one sweep pass of a browser wallet with the same code the node
// runs: which addresses are exposed, where each goes, and the transactions that
// move them, burning every coin hour. It signs nothing.
func SweepPlan(body string) (string, error) {
	var r SweepRequest
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return "", err
	}
	if err := r.Rule.Check(); err != nil {
		return "", err
	}
	addrs, err := r.addrs()
	if err != nil {
		return "", err
	}
	sent, err := addrSet(r.Sent)
	if err != nil {
		return "", err
	}
	seen, err := addrSet(r.Seen)
	if err != nil {
		return "", err
	}
	covered := make(map[cipher.SHA256]bool, len(r.Covered))
	for _, h := range r.Covered {
		hash, err := cipher.SHA256FromHex(h)
		if err != nil {
			return "", fmt.Errorf("invalid output hash %q: %w", h, err)
		}
		covered[hash] = true
	}
	outs := make([]plan.Output, len(r.Outputs))
	for i, o := range r.Outputs {
		hash, err := cipher.SHA256FromHex(o.Hash)
		if err != nil {
			return "", fmt.Errorf("invalid output hash %q: %w", o.Hash, err)
		}
		a, err := cipher.DecodeBase58Address(o.Address)
		if err != nil {
			return "", fmt.Errorf("invalid address %q: %w", o.Address, err)
		}
		outs[i] = plan.Output{Hash: hash, Address: a, Coins: o.Coins, Hours: o.Hours, Contested: o.Contested}
	}
	maxSize := r.MaxSize
	if maxSize == 0 {
		maxSize = DefaultMaxTransactionSize
	}

	var res SweepResult
	if needs := r.Rule.GapNeeds(addrs, seen); len(needs) > 0 {
		res.Gap = make(map[string]uint64, len(needs))
		for acct, n := range needs {
			res.Gap[strconv.FormatUint(uint64(acct), 10)] = n
		}
	}
	p := r.Rule.Find(addrs, sent, outs, covered)
	res.Exposed, res.Uncovered = p.Exposed, p.Uncovered
	res.Coins, res.Hours, res.Contested = p.Coins, p.Hours, p.Contested

	exp := r.Rule.Exposed(addrs, sent)
	var order []cipher.Address
	byDest := make(map[cipher.Address][]plan.Output)
	for _, g := range p.Groups {
		s, ok := r.Destinations[strconv.Itoa(g.Target)]
		if !ok || s == "" {
			res.Need = append(res.Need, g.Target)
			continue
		}
		dest, err := cipher.DecodeBase58Address(s)
		if err == nil {
			err = plan.CheckDestination(dest, exp, sent)
		}
		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("%s: %v", targetName(r.Rule, g.Target), err))
			continue
		}
		if _, ok := byDest[dest]; !ok {
			order = append(order, dest)
		}
		byDest[dest] = append(byDest[dest], g.Outputs...)
	}
	if len(res.Need) > 0 {
		// Resolve every destination first, so a pass sends all or waits.
		order = nil
	}
	for _, dest := range order {
		sweeps, err := plan.Build(byDest[dest], dest, maxSize)
		if errors.Is(err, plan.ErrNoHours) {
			res.NoHours = true
			continue
		}
		if err != nil {
			res.Errors = append(res.Errors, err.Error())
			continue
		}
		for _, sw := range sweeps {
			t := SweepTransaction{Destination: dest.String(), Fee: sw.Fee}
			for _, in := range sw.Inputs {
				t.Inputs = append(t.Inputs, SweepInput{Hash: in.Hash.Hex(), Address: in.Address.String()})
			}
			for _, o := range sw.Txn.Out {
				t.Outputs = append(t.Outputs, SweepTxOutput{Address: o.Address.String(), Coins: o.Coins, Hours: o.Hours})
			}
			res.Sweeps = append(res.Sweeps, t)
		}
	}
	b, err := json.Marshal(res)
	return string(b), err
}

func targetName(r plan.Rule, t int) string {
	if t == plan.Level1 {
		return "destination"
	}
	m := r.Mappings[t].From
	if m.Account != nil {
		return fmt.Sprintf("mapping of account %d", *m.Account)
	}
	return "mapping of " + m.Address
}
