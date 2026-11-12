// Package ui2 src/ui2/settings.go c3-ui-wallet
package ui2

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/api"
)

// nodeWallets are the wallets the node keeps, with their addresses.
func nodeWallets(n Node) ([]HistWallet, error) {
	ws, err := n.Wallets()
	if err != nil {
		return nil, err
	}
	out := make([]HistWallet, 0, len(ws))
	for _, w := range ws {
		h := HistWallet{ID: w.Meta.Filename, Label: w.Meta.Label}
		for _, e := range w.Entries {
			h.Addrs = append(h.Addrs, e.Address)
		}
		out = append(out, h)
	}
	return out, nil
}

func sortWallets(ws []HistWallet) {
	sort.SliceStable(ws, func(i, j int) bool { return strings.ToLower(ws[i].Label) < strings.ToLower(ws[j].Label) })
}

// Output is one unspent output.
type Output struct {
	Hash  string
	Coins string
	Hours string
}

// OutputAddr is an address and its unspent outputs.
type OutputAddr struct {
	Address string
	Outputs []Output
}

// OutputWallet is a wallet's addresses with their unspent outputs.
type OutputWallet struct {
	Label string
	Addrs []OutputAddr
}

// Outputs is the unspent outputs page.
type Outputs struct {
	Wallets  []OutputWallet
	Disabled bool
	Err      string
}

// LoadOutputs reads the unspent outputs of every wallet in ws.
func LoadOutputs(n Node, ws []HistWallet) Outputs {
	var out Outputs
	sortWallets(ws)
	var all []string
	for _, w := range ws {
		all = append(all, w.Addrs...)
	}
	by := map[string][]Output{}
	if len(all) > 0 {
		s, err := n.OutputsForAddresses(all)
		if err != nil {
			out.Err = err.Error()
			return out
		}
		for _, o := range s.HeadOutputs {
			by[o.Address] = append(by[o.Address], Output{Hash: o.Hash, Coins: o.Coins, Hours: utoa(o.CalculatedHours)})
		}
	}
	for _, w := range ws {
		ow := OutputWallet{Label: w.Label}
		if ow.Label == "" {
			ow.Label = w.ID
		}
		for _, a := range w.Addrs {
			ow.Addrs = append(ow.Addrs, OutputAddr{Address: a, Outputs: by[a]})
		}
		out.Wallets = append(out.Wallets, ow)
	}
	return out
}

// LoadNodeOutputs is LoadOutputs for the node's wallets.
func LoadNodeOutputs(n Node) Outputs {
	ws, err := nodeWallets(n)
	if err != nil {
		return Outputs{Disabled: walletAPIDisabled(err), Err: err.Error()}
	}
	return LoadOutputs(n, ws)
}

// PendingRow is one transaction waiting in the node's pool.
type PendingRow struct {
	Txid  string
	Coins string
	Hours string
	Time  string
}

// Pending is the pending transactions page. All shows the whole pool,
// otherwise only the transactions of the wallets.
type Pending struct {
	All  bool
	Rows []PendingRow
	Err  string
}

// LoadPending reads the node's pool, keeping with all unset only the
// transactions that touch the addresses of ws. A row sums every output, as
// the Angular wallet's page does.
func LoadPending(n Node, ws []HistWallet, all bool) Pending {
	out := Pending{All: all}
	own := map[string]bool{}
	for _, w := range ws {
		for _, a := range w.Addrs {
			own[a] = true
		}
	}
	txns, err := n.PendingTransactionsVerbose()
	if err != nil {
		out.Err = err.Error()
		return out
	}
	sort.SliceStable(txns, func(i, j int) bool { return txns[i].Received.After(txns[j].Received) })
	for _, t := range txns {
		mine := false
		var c, h uint64
		for _, i := range t.Transaction.In {
			mine = mine || own[i.Address]
		}
		for _, o := range t.Transaction.Out {
			mine = mine || own[o.Address]
			c += droplets(o.Coins)
			h += o.Hours
		}
		if !all && !mine {
			continue
		}
		out.Rows = append(out.Rows, PendingRow{
			Txid:  t.Transaction.Hash,
			Coins: coins(c),
			Hours: utoa(h),
			Time:  t.Received.UTC().Format("2006-01-02 15:04 UTC"),
		})
	}
	return out
}

// BackupWallet is a node wallet on the backup page.
type BackupWallet struct {
	ID        string
	Label     string
	Encrypted bool
}

// Backup is the backup page.
type Backup struct {
	Folder   string
	Wallets  []BackupWallet
	Disabled bool
	Err      string
}

// LoadBackup reads the node's wallet directory and wallets.
func LoadBackup(n Node) Backup {
	var out Backup
	ws, err := n.Wallets()
	if err != nil {
		return Backup{Disabled: walletAPIDisabled(err), Err: err.Error()}
	}
	if f, err := n.WalletFolderName(); err == nil {
		out.Folder = f.Address
	}
	for _, w := range ws {
		out.Wallets = append(out.Wallets, BackupWallet{ID: w.Meta.Filename, Label: w.Meta.Label, Encrypted: w.Meta.Encrypted})
	}
	sort.SliceStable(out.Wallets, func(i, j int) bool {
		return strings.ToLower(out.Wallets[i].Label) < strings.ToLower(out.Wallets[j].Label)
	})
	return out
}

var errSeedAPI = errors.New("backup.err-seed-api")

// WalletSeed asks the node for the seed of encrypted wallet id. The node
// answers only with the INSECURE_WALLET_SEED API set enabled.
func WalletSeed(n Node, id, password string) (*api.WalletSeedResponse, error) {
	if password == "" {
		return nil, errNoPassword
	}
	s, err := n.WalletSeed(id, password)
	if walletAPIDisabled(err) {
		return nil, errSeedAPI
	}
	return s, err
}

// AddressesCSV is the addresses of node wallet id, one per line, as the
// Angular wallet's "save addresses" writes them.
func AddressesCSV(n Node, id string) (label, csv string, err error) {
	w, err := n.Wallet(id)
	if err != nil {
		return "", "", err
	}
	addrs := make([]string, 0, len(w.Entries))
	for _, e := range w.Entries {
		addrs = append(addrs, e.Address)
	}
	label = w.Meta.Label
	if label == "" {
		label = id
	}
	return label, strings.Join(addrs, "\n"), nil
}

// Supply is the coin supply on the status page.
type Supply struct {
	Current string
	Total   string
}

// ago is how long before now t was, rounded to a readable unit, empty for
// a zero time.
func ago(now time.Time, unix int64) string {
	if unix <= 0 {
		return ""
	}
	d := now.Sub(time.Unix(unix, 0))
	switch {
	case d < time.Minute:
		return itoa(max(int(d.Seconds()), 0)) + "s"
	case d < time.Hour:
		return itoa(int(d.Minutes())) + "m"
	case d < 48*time.Hour:
		return itoa(int(d.Hours())) + "h"
	}
	return itoa(int(d.Hours()/24)) + "d"
}

func pendingVals(txid string, all bool) string {
	b, _ := json.Marshal(map[string]string{"txid": txid, "all": strconv.FormatBool(all)}) //nolint:errcheck
	return string(b)
}
