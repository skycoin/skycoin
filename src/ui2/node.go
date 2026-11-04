package ui2

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// Node is the part of the node's REST API the pages read. *api.Client has
// every method, for the local node or a remote one.
type Node interface {
	Health() (*api.HealthResponse, error)
	BlockchainProgress() (*readable.BlockchainProgress, error)
	NetworkConnections(*api.NetworkConnectionsFilter) (*api.Connections, error)
	Wallets() ([]api.WalletResponse, error)
	Wallet(id string) (*api.WalletResponse, error)
	WalletBalance(id string) (*api.BalanceResponse, error)
	TransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error)
}

// walletAPIDisabled reports whether err is the node refusing the WALLET API
// set, as it does when the wallets are kept in the browser.
func walletAPIDisabled(err error) bool {
	var ce api.ClientError
	return errors.As(err, &ce) && ce.StatusCode == http.StatusForbidden
}

func coins(n uint64) string {
	s, err := droplet.ToString(n)
	if err != nil {
		return "?"
	}
	return s
}

func utoa(n uint64) string { return strconv.FormatUint(n, 10) }

func itoa(n int) string { return strconv.Itoa(n) }

func unixTime(sec uint64) string {
	if sec == 0 {
		return ""
	}
	return time.Unix(int64(sec), 0).UTC().Format("2006-01-02 15:04 UTC") //nolint:gosec
}

// Amount is a balance in coins and coin hours, confirmed and as it will be
// once the pending transactions confirm.
type Amount struct {
	Coins, Hours               string
	PendingCoins, PendingHours string
	Pending                    bool
}

func amountOf(b readable.BalancePair) Amount {
	return Amount{
		Coins:        coins(b.Confirmed.Coins),
		Hours:        utoa(b.Confirmed.Hours),
		PendingCoins: coins(b.Predicted.Coins),
		PendingHours: utoa(b.Predicted.Hours),
		Pending:      b.Confirmed != b.Predicted,
	}
}

// WalletRow is one wallet in the wallet list.
type WalletRow struct {
	ID        string
	Label     string
	Type      string
	Encrypted bool
	Addresses int
	Balance   Amount
	Err       string
}

// Wallets is the wallet list page.
type Wallets struct {
	Rows     []WalletRow
	Total    Amount
	Disabled bool
	Err      string
}

// LoadWallets reads every wallet and its balance.
func LoadWallets(n Node) Wallets {
	ws, err := n.Wallets()
	if err != nil {
		return Wallets{Disabled: walletAPIDisabled(err), Err: err.Error()}
	}
	var out Wallets
	var total readable.BalancePair
	for _, w := range ws {
		r := WalletRow{
			ID:        w.Meta.Filename,
			Label:     w.Meta.Label,
			Type:      w.Meta.Type,
			Encrypted: w.Meta.Encrypted,
			Addresses: len(w.Entries),
		}
		if b, err := n.WalletBalance(r.ID); err != nil {
			r.Err = err.Error()
		} else {
			r.Balance = amountOf(b.BalancePair)
			total.Confirmed.Coins += b.Confirmed.Coins
			total.Confirmed.Hours += b.Confirmed.Hours
			total.Predicted.Coins += b.Predicted.Coins
			total.Predicted.Hours += b.Predicted.Hours
		}
		out.Rows = append(out.Rows, r)
	}
	sort.SliceStable(out.Rows, func(i, j int) bool {
		return strings.ToLower(out.Rows[i].Label) < strings.ToLower(out.Rows[j].Label)
	})
	out.Total = amountOf(total)
	return out
}

// AddressRow is one address of a wallet.
type AddressRow struct {
	Address string
	Chain   string
	Index   string
	Coins   string
	Hours   string
	Empty   bool
}

// Wallet is the page of one wallet.
type Wallet struct {
	ID       string
	Meta     readable.WalletMeta
	Balance  Amount
	Addrs    []AddressRow
	Pager    Pager
	Disabled bool
	Err      string
}

// LoadWallet reads wallet id, its balance and the page of addresses q asks for.
func LoadWallet(n Node, id string, q ListQuery) Wallet {
	out := Wallet{ID: id}
	w, err := n.Wallet(id)
	if err != nil {
		out.Disabled, out.Err = walletAPIDisabled(err), err.Error()
		return out
	}
	out.Meta = w.Meta
	b, err := n.WalletBalance(id)
	if err != nil {
		out.Err = err.Error()
		b = &api.BalanceResponse{}
	}
	out.Balance = amountOf(b.BalancePair)

	rows := make([]AddressRow, 0, len(w.Entries))
	for _, e := range w.Entries {
		r := AddressRow{Address: e.Address}
		if e.Change != nil {
			r.Chain = "external"
			if *e.Change == 1 {
				r.Chain = "change"
			}
		}
		if e.ChildNumber != nil {
			r.Index = strconv.FormatUint(uint64(*e.ChildNumber), 10)
		}
		ab := b.Addresses[e.Address]
		r.Coins, r.Hours = coins(ab.Confirmed.Coins), utoa(ab.Confirmed.Hours)
		r.Empty = ab.Confirmed.Coins == 0 && ab.Predicted.Coins == 0
		if q.Match(r.Address, r.Chain) && !(q.Hide && r.Empty) {
			rows = append(rows, r)
		}
	}
	out.Addrs, out.Pager = Paginate(rows, q, ListPageSize)
	return out
}

// TxnRow is one transaction of a wallet, from the wallet's side.
type TxnRow struct {
	Txid      string
	Time      string
	Confirmed bool
	Height    uint64
	// Kind is "in", "out" or "self", a transfer between the wallet's own
	// addresses.
	Kind  string
	Coins string
	Hours string
}

// Txns is a page of a wallet's transaction history.
type Txns struct {
	ID    string
	Rows  []TxnRow
	Pager Pager
	Err   string
}

// LoadTxns reads the history of wallet id, newest first.
func LoadTxns(n Node, id string, q ListQuery) Txns {
	out := Txns{ID: id}
	w, err := n.Wallet(id)
	if err != nil {
		out.Err = err.Error()
		return out
	}
	own := make(map[string]bool, len(w.Entries))
	addrs := make([]string, 0, len(w.Entries))
	for _, e := range w.Entries {
		own[e.Address] = true
		addrs = append(addrs, e.Address)
	}
	var txns []readable.TransactionWithStatusVerbose
	if len(addrs) > 0 {
		if txns, err = n.TransactionsVerbose(addrs); err != nil {
			out.Err = err.Error()
			return out
		}
	}
	sort.SliceStable(txns, func(i, j int) bool {
		a, b := txns[i], txns[j]
		if a.Status.Confirmed != b.Status.Confirmed {
			return !a.Status.Confirmed
		}
		return a.Time > b.Time
	})
	rows := make([]TxnRow, 0, len(txns))
	for _, t := range txns {
		r := txnRow(t, own)
		if q.Match(r.Txid) {
			rows = append(rows, r)
		}
	}
	out.Rows, out.Pager = Paginate(rows, q, TxnPageSize)
	return out
}

// TxnPageSize is how many transactions a history page shows.
const TxnPageSize = 20

func txnRow(t readable.TransactionWithStatusVerbose, own map[string]bool) TxnRow {
	var in, out, inH, outH int64
	foreign := false
	for _, i := range t.Transaction.In {
		if own[i.Address] {
			c, _ := droplet.FromString(i.Coins) //nolint:errcheck
			in += int64(c)                      //nolint:gosec
			inH += int64(i.CalculatedHours)     //nolint:gosec
		}
	}
	for _, o := range t.Transaction.Out {
		if own[o.Address] {
			c, _ := droplet.FromString(o.Coins) //nolint:errcheck
			out += int64(c)                     //nolint:gosec
			outH += int64(o.Hours)              //nolint:gosec
		} else {
			foreign = true
		}
	}
	r := TxnRow{
		Txid:      t.Transaction.Hash,
		Time:      unixTime(t.Time),
		Confirmed: t.Status.Confirmed,
		Height:    t.Status.Height,
	}
	net, netH := out-in, outH-inH
	switch {
	case !foreign && in > 0:
		r.Kind, net, netH = "self", out, outH
	case net < 0:
		r.Kind, net, netH = "out", -net, -netH
	default:
		r.Kind = "in"
	}
	if netH < 0 {
		netH = 0
	}
	r.Coins, r.Hours = coins(uint64(net)), strconv.FormatInt(netH, 10) //nolint:gosec
	return r
}

// Peer is one connection of the node.
type Peer struct {
	Addr     string
	Outgoing bool
	Height   uint64
	Agent    string
	Trusted  bool
}

// Status is the node and blockchain status page.
type Status struct {
	Coin        string
	Version     string
	Height      uint64
	BlockHash   string
	BlockTime   string
	Current     uint64
	Highest     uint64
	Synced      bool
	Percent     string
	Unconfirmed uint64
	Open        int
	Incoming    int
	Outgoing    int
	Uptime      string
	WalletAPI   bool
	Peers       []Peer
	Err         string
}

// LoadStatus reads the node's health, sync progress and connections.
func LoadStatus(n Node) Status {
	var s Status
	h, err := n.Health()
	if err != nil {
		s.Err = err.Error()
		return s
	}
	s.Coin = h.CoinName
	s.Version = h.Version.Version
	s.Height = h.BlockchainMetadata.Head.BkSeq
	s.BlockHash = h.BlockchainMetadata.Head.Hash
	s.BlockTime = unixTime(h.BlockchainMetadata.Head.Time)
	s.Unconfirmed = h.BlockchainMetadata.Unconfirmed
	s.Open, s.Incoming, s.Outgoing = h.OpenConnections, h.IncomingConnections, h.OutgoingConnections
	s.Uptime = h.Uptime.Truncate(time.Second).String()
	s.WalletAPI = h.WalletAPIEnabled

	if p, err := n.BlockchainProgress(); err != nil {
		s.Err = err.Error()
	} else {
		s.Current, s.Highest = p.Current, p.Highest
		s.Synced = p.Highest == 0 || p.Current >= p.Highest
		pct := 100.0
		if p.Highest > 0 && !s.Synced {
			pct = float64(p.Current) * 100 / float64(p.Highest)
		}
		s.Percent = strconv.FormatFloat(pct, 'f', 1, 64)
	}

	if c, err := n.NetworkConnections(nil); err != nil {
		s.Err = err.Error()
	} else {
		for _, x := range c.Connections {
			s.Peers = append(s.Peers, Peer{
				Addr:     x.Addr,
				Outgoing: x.Outgoing,
				Height:   x.Height,
				Agent:    agent(x.UserAgent.Coin, x.UserAgent.Version),
				Trusted:  x.IsTrustedPeer,
			})
		}
		sort.Slice(s.Peers, func(i, j int) bool { return s.Peers[i].Addr < s.Peers[j].Addr })
	}
	return s
}

func agent(coin, version string) string {
	if coin == "" {
		return ""
	}
	return coin + ":" + version
}
