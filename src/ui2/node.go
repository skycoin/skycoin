// Package ui2 src/ui2/node.go c3-ui-wallet
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
	"github.com/skycoin/skycoin/src/wallet"
)

// Node is the part of the node's REST API the pages use. *api.Client has
// every method, for the local node or a remote one.
type Node interface {
	Health() (*api.HealthResponse, error)
	BlockchainProgress() (*readable.BlockchainProgress, error)
	NetworkConnections(*api.NetworkConnectionsFilter) (*api.Connections, error)
	Wallets() ([]api.WalletResponse, error)
	Wallet(id string) (*api.WalletResponse, error)
	WalletBalance(id string) (*api.BalanceResponse, error)
	TransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error)
	WalletXPub(id, path string) (*api.WalletXPub, error)
	NewSeed(entropy int) (string, error)
	VerifySeed(seed string) (bool, error)
	CreateWallet(api.CreateWalletOptions) (*api.WalletResponse, error)
	CreateWalletTemp(api.CreateWalletOptions) (*api.WalletResponse, error)
	RecoverWallet(api.WalletRecoverRequest) (*api.WalletResponse, error)
	NewWalletAddress(id, password string, options ...wallet.Option) ([]string, error)
	NewWalletAccount(id, password, name string) (*api.WalletResponse, error)
	RemovePendingTransaction(txid string) error
	ScanWalletAddresses(id string, n int, password string) ([]string, error)
	UpdateWallet(id, label string) error
	EncryptWallet(id, password string) (*api.WalletResponse, error)
	DecryptWallet(id, password string) (*api.WalletResponse, error)
	WalletCreateTransaction(api.WalletCreateTransactionRequest) (*api.CreateTransactionResponse, error)
	InjectEncodedTransaction(rawtx string) (string, error)
	WalletFolderName() (*api.WalletFolder, error)
	WalletSeed(id, password string) (*api.WalletSeedResponse, error)
	OutputsForAddresses(addrs []string) (*readable.UnspentOutputsSummary, error)
	PendingTransactionsVerbose() ([]readable.UnconfirmedTransactionVerbose, error)
	CoinSupply() (*api.CoinSupply, error)
	Balance(addrs []string) (*api.BalanceResponse, error)
	CreateTransaction(api.CreateTransactionRequest) (*api.CreateTransactionResponse, error)
	VerifyTransaction(api.VerifyTransactionRequest) (*api.VerifyTransactionResponse, error)
	WalletSignTransaction(api.WalletSignTransactionRequest) (*api.CreateTransactionResponse, error)
	OutputsForHashes(hashes []string) (*readable.UnspentOutputsSummary, error)
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
	Temp      bool
	Addresses int
	Balance   Amount
	Err       string
	Href      string
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
			Temp:      w.Meta.Temp,
			Addresses: len(w.Entries),
			Href:      walletURL(w.Meta.Filename),
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
	// Exposed is an address that has spent, so its public key is on the
	// blockchain, and still holds coins. MoveHref moves its coins away.
	Exposed  bool
	MoveHref string
}

// Wallet is the page of one wallet.
type Wallet struct {
	ID      string
	Meta    readable.WalletMeta
	Balance Amount
	Addrs   []AddressRow
	Pager   Pager
	// Exposed counts the addresses that hold coins and have spent.
	Exposed int
	Notice  Notice
	// Accounts are a bip44 wallet's accounts after the first, which the
	// address list does not show.
	Accounts []AccountView
	// AccountIDs are every account index, for the new address form.
	AccountIDs []uint32
	// BrowserAddrs are every address of a browser wallet, in its order.
	BrowserAddrs []string
	Disabled     bool
	Err          string
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
	rows, exposed, err := addressRows(n, w.Entries, b, q, func(a string) string { return sendURL(id, a) })
	out.Exposed = exposed
	if err != nil {
		out.Err = err.Error()
	}
	out.Addrs, out.Pager = Paginate(rows, q, ListPageSize)
	out.Accounts, out.AccountIDs = loadAccounts(n, w.Accounts)
	return out
}

// AccountView is one bip44 account with its addresses and balance.
type AccountView struct {
	Index   uint32
	Name    string
	Balance Amount
	Addrs   []AddressRow
	Err     string
}

func loadAccounts(n Node, accts []readable.WalletAccount) ([]AccountView, []uint32) {
	var out []AccountView
	var ids []uint32
	for _, a := range accts {
		ids = append(ids, a.Index)
		if a.Index == 0 {
			continue
		}
		v := AccountView{Index: a.Index, Name: a.Name}
		var addrs []string
		for _, e := range append(append([]readable.WalletEntry{}, a.ExternalEntries...), a.ChangeEntries...) {
			addrs = append(addrs, e.Address)
		}
		b := &api.BalanceResponse{}
		if len(addrs) > 0 {
			var err error
			if b, err = n.Balance(addrs); err != nil {
				v.Err, b = err.Error(), &api.BalanceResponse{}
			}
		}
		v.Balance = amountOf(b.BalancePair)
		for i, list := range [][]readable.WalletEntry{a.ExternalEntries, a.ChangeEntries} {
			chain := "external"
			if i == 1 {
				chain = "change"
			}
			for j, e := range list {
				ab := b.Addresses[e.Address]
				v.Addrs = append(v.Addrs, AddressRow{Address: e.Address, Chain: chain, Index: strconv.Itoa(j), Coins: coins(ab.Confirmed.Coins), Hours: utoa(ab.Confirmed.Hours), Empty: ab.Confirmed.Coins == 0})
			}
		}
		out = append(out, v)
	}
	return out, ids
}

// addressRows are the rows of entries that q keeps, with their balances in
// b, and how many of all entries are exposed. move is the page that moves
// the coins off an exposed address.
func addressRows(n Node, entries []readable.WalletEntry, b *api.BalanceResponse, q ListQuery, move func(string) string) ([]AddressRow, int, error) {
	exposed := 0
	all := make([]string, 0, len(entries))
	for _, e := range entries {
		all = append(all, e.Address)
	}
	use, err := addressUse(n, all)

	rows := make([]AddressRow, 0, len(entries))
	for _, e := range entries {
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
		r.Exposed = use[e.Address].Spent && !r.Empty
		if r.Exposed {
			r.MoveHref = move(e.Address)
			exposed++
		}
		if q.Match(r.Address, r.Chain) && !(q.Hide && r.Empty) {
			rows = append(rows, r)
		}
	}
	return rows, exposed, err
}

// Peer is one connection of the node.
type Peer struct {
	Addr     string
	Outgoing bool
	Height   uint64
	Agent    string
	Trusted  bool
	Sent     string
	Received string
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
	Supply      Supply
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

	if cs, err := n.CoinSupply(); err == nil {
		s.Supply = Supply{Current: cs.CurrentSupply, Total: cs.TotalSupply}
	}

	now := time.Now()
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
				Sent:     ago(now, x.LastSent),
				Received: ago(now, x.LastReceived),
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
