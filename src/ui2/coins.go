// Package ui2 src/ui2/coins.go c3-ui-wallet
package ui2

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/api"
)

// Coin is one coin a multi-coin host serves, as skycoin web lists them at
// /api/v1/coins (cmd/skycoin-web/commands/root.go, discoveredCoin).
type Coin struct {
	ID            int    `json:"id"`
	CoinName      string `json:"coinName"`
	CoinSymbol    string `json:"coinSymbol"`
	HoursName     string `json:"hoursName"`
	CoinType      string `json:"coinType"`
	ServerWallets bool   `json:"serverWallets"`
	Decimals      int    `json:"decimals,omitempty"`
	Contract      string `json:"contract,omitempty"`
	CustomID      string `json:"customId,omitempty"`
	IconURL       string `json:"iconUrl,omitempty"`
}

// Gateway is "btc" or "eth" for a coin skycoin web serves from its own
// /v1/btc or /v1/eth routes, empty for a Skycoin fiber coin with a node.
func (c Coin) Gateway() string {
	switch c.CoinType {
	case "bitcoin":
		return "btc"
	case "ethereum", "erc20":
		return "eth"
	}
	return ""
}

// Places is how many decimal places an amount of the coin has, as the
// Angular wallet's BaseCoin.coinDecimals.
func (c Coin) Places() int {
	switch c.Gateway() {
	case "btc":
		return 8
	case "eth":
		if c.Decimals > 0 {
			return c.Decimals
		}
		return 18
	}
	return 6
}

// Label is how the coin is named in the picker.
func (c Coin) Label() string {
	if c.CoinSymbol == "" {
		return c.CoinName
	}
	return c.CoinName + " (" + c.CoinSymbol + ")"
}

// CoinHost reaches a host that serves several coins, as skycoin web does:
// the coin list at /api/v1/coins, each coin's API at /coin/{id}/api, and the
// added coins at /api/v1/coins/custom.
type CoinHost struct {
	host   string
	client *http.Client
	h      http.Handler
}

// NewCoinHost calls h in process with requests for host.
func NewCoinHost(h http.Handler, host string) *CoinHost {
	c := api.HandlerClient(h, host)
	return &CoinHost{host: host, client: c.HTTPClient, h: h}
}

// Node is the API client of coin id.
func (c *CoinHost) Node(id int) *api.Client {
	n := api.HandlerClient(c.h, c.host)
	n.Addr = "http://" + c.host + "/coin/" + strconv.Itoa(id) + "/"
	return n
}

// errorBody is the {"error": "..."} body skycoin web answers with.
type errorBody struct {
	Error string `json:"error"`
}

// call sends a request to path, which starts with "/", and reads the JSON
// answer into out. body, when not nil, is sent as JSON.
func (c *CoinHost) call(method, path string, body, out interface{}) error {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, "http://"+c.host+path, r)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var e errorBody
		if json.Unmarshal(b, &e) == nil && e.Error != "" {
			return errors.New(e.Error)
		}
		return errors.New(resp.Status + " " + strings.TrimSpace(string(b)))
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(b, out)
}

// Coins are the coins the host offers.
func (c *CoinHost) Coins() ([]Coin, error) {
	var out []Coin
	err := c.call(http.MethodGet, "/api/v1/coins", nil, &out)
	return out, err
}

func gwPath(coin Coin, route string) string {
	return "/coin/" + strconv.Itoa(coin.ID) + "/api/v1/" + coin.Gateway() + "/" + route
}

// Units is an amount in a coin's smallest unit, which skycoin web writes as
// a JSON number for Bitcoin and as a decimal string for Ethereum coins.
type Units struct{ big.Int }

// UnmarshalJSON reads a number or a decimal string.
func (u *Units) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		u.SetInt64(0)
		return nil
	}
	if _, ok := u.SetString(s, 10); !ok {
		return errors.New("not an amount: " + s)
	}
	return nil
}

// Format writes u in whole coins with places decimal places, trailing
// zeros dropped.
func (u *Units) Format(places int) string { return formatUnits(&u.Int, places) }

func formatUnits(v *big.Int, places int) string {
	neg := v.Sign() < 0
	s := new(big.Int).Abs(v).String()
	if places > 0 {
		if len(s) <= places {
			s = strings.Repeat("0", places-len(s)+1) + s
		}
		s = strings.TrimRight(s[:len(s)-places]+"."+s[len(s)-places:], "0")
		s = strings.TrimSuffix(s, ".")
	}
	if neg {
		s = "-" + s
	}
	return s
}

var errAmount = errors.New("send.err-coins")

// parseUnits reads a decimal amount of whole coins into the smallest unit.
func parseUnits(s string, places int) (*big.Int, error) {
	s = strings.TrimSpace(s)
	whole, frac, _ := strings.Cut(s, ".")
	if whole == "" && frac == "" || len(frac) > places || strings.ContainsAny(whole+frac, "+-") {
		return nil, errAmount
	}
	v, ok := new(big.Int).SetString(whole+frac+strings.Repeat("0", places-len(frac)), 10)
	if !ok || v.Sign() <= 0 {
		return nil, errAmount
	}
	return v, nil
}

// GwBalance is a gateway coin's balance, in total and per address.
type GwBalance struct {
	Confirmed struct {
		Coins Units `json:"coins"`
	} `json:"confirmed"`
	Addresses map[string]struct {
		Confirmed struct {
			Coins Units `json:"coins"`
		} `json:"confirmed"`
	} `json:"addresses"`
}

// Balance reads the balance of addrs.
func (c *CoinHost) Balance(coin Coin, addrs []string) (*GwBalance, error) {
	var out GwBalance
	if len(addrs) == 0 {
		return &out, nil
	}
	err := c.call(http.MethodGet, gwPath(coin, "balance")+"?addrs="+url.QueryEscape(strings.Join(addrs, ",")), nil, &out)
	return &out, err
}

// GwTx is a transaction in a gateway coin's history.
type GwTx struct {
	TxID          string `json:"txid"`
	Timestamp     int64  `json:"timestamp"`
	BlockHeight   int64  `json:"block_height"`
	Confirmations int64  `json:"confirmations"`
	Fee           Units  `json:"fee"`
	Failed        bool   `json:"failed"`
	Inputs        []struct {
		Address string `json:"address"`
		Value   Units  `json:"value"`
	} `json:"inputs"`
	Outputs []struct {
		Address string `json:"address"`
		Value   Units  `json:"value"`
	} `json:"outputs"`
}

// History reads the transactions of addrs, newest first.
func (c *CoinHost) History(coin Coin, addrs []string) ([]GwTx, error) {
	var out []GwTx
	if len(addrs) == 0 {
		return out, nil
	}
	err := c.call(http.MethodGet, gwPath(coin, "history")+"?addrs="+url.QueryEscape(strings.Join(addrs, ",")), nil, &out)
	return out, err
}

// FeePreset is one of the slow, normal and fast fees the host suggests:
// SatPerByte for Bitcoin, the EIP-1559 caps and MaxFee in wei for Ethereum.
type FeePreset struct {
	SatPerByte           int64  `json:"sat_per_byte"`
	Blocks               int    `json:"blocks"`
	MaxFeePerGas         string `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"`
	MaxFee               Units  `json:"max_fee"`
}

// Fees are the fee presets by name.
type Fees struct {
	Presets map[string]FeePreset `json:"presets"`
}

// Fee reads the fee presets.
func (c *CoinHost) Fee(coin Coin) (*Fees, error) {
	var out Fees
	err := c.call(http.MethodGet, gwPath(coin, "fee"), nil, &out)
	return &out, err
}

// GwSend is a send from a server wallet of a gateway coin, in the body
// each gateway takes.
type GwSend struct {
	WalletID string
	Password string
	To       string
	Amount   *big.Int
	Fee      FeePreset
}

// Send has the host sign and broadcast s, and returns the transaction id.
func (c *CoinHost) Send(coin Coin, s GwSend) (string, error) {
	var body interface{}
	if coin.Gateway() == "btc" {
		if !s.Amount.IsInt64() {
			return "", errAmount
		}
		body = map[string]interface{}{
			"wallet_id":    s.WalletID,
			"password":     s.Password,
			"destinations": []map[string]interface{}{{"address": s.To, "coins": s.Amount.Int64()}},
			"fee_rate":     s.Fee.SatPerByte,
		}
	} else {
		body = map[string]string{
			"wallet_id":                s.WalletID,
			"password":                 s.Password,
			"to":                       s.To,
			"amount":                   s.Amount.String(),
			"max_fee_per_gas":          s.Fee.MaxFeePerGas,
			"max_priority_fee_per_gas": s.Fee.MaxPriorityFeePerGas,
		}
	}
	var out struct {
		TxID string `json:"txid"`
	}
	err := c.call(http.MethodPost, gwPath(coin, "send"), body, &out)
	return out.TxID, err
}

// CreateGwWallet makes a bip44 wallet of a gateway coin. A Bitcoin wallet
// has native segwit addresses unless legacy is set.
func (c *CoinHost) CreateGwWallet(coin Coin, f CreateForm, legacy bool) (string, error) {
	v := url.Values{}
	v.Set("type", "bip44")
	v.Set("seed", f.Seed)
	v.Set("label", f.Label)
	v.Set("scan", "20")
	if f.Passphrase != "" {
		v.Set("seed-passphrase", f.Passphrase)
	}
	if f.Password != "" {
		v.Set("encrypt", "true")
		v.Set("password", f.Password)
	}
	if legacy && coin.Gateway() == "btc" {
		v.Set("segwit", "false")
	}
	req, err := http.NewRequest(http.MethodPost, "http://"+c.host+"/coin/"+strconv.Itoa(coin.ID)+"/api/v1/wallet/create", strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //nolint:errcheck
	var w api.WalletResponse
	b, _ := io.ReadAll(resp.Body) //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(strings.TrimSpace(string(b)))
	}
	if err := json.Unmarshal(b, &w); err != nil {
		return "", err
	}
	return w.Meta.Filename, nil
}

// CustomCoin is a coin added in the wallet, a fiber coin with its node or
// an ERC-20 token (cmd/skycoin-web/commands/custom_coins.go).
type CustomCoin struct {
	ID       string `json:"id,omitempty"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Ticker   string `json:"ticker"`
	NodeURL  string `json:"node_url,omitempty"`
	Contract string `json:"contract,omitempty"`
	Decimals int    `json:"decimals,omitempty"`
}

// CustomCoins lists the added coins.
func (c *CoinHost) CustomCoins() ([]CustomCoin, error) {
	var out []CustomCoin
	err := c.call(http.MethodGet, "/api/v1/coins/custom", nil, &out)
	return out, err
}

// AddCoin adds a fiber coin or an ERC-20 token.
func (c *CoinHost) AddCoin(cc CustomCoin) error {
	return c.call(http.MethodPost, "/api/v1/coins/custom", cc, nil)
}

// RemoveCoin removes an added coin.
func (c *CoinHost) RemoveCoin(id string) error {
	return c.call(http.MethodDelete, "/api/v1/coins/custom?id="+url.QueryEscape(id), nil, nil)
}

// Endpoints are the Ethereum RPC servers and indexer every Ethereum coin
// shares.
type Endpoints struct {
	RPCURLs           []string `json:"rpc_urls"`
	IndexerURL        string   `json:"indexer_url"`
	DefaultRPCURLs    []string `json:"default_rpc_urls,omitempty"`
	DefaultIndexerURL string   `json:"default_indexer_url,omitempty"`
}

// Endpoints reads the Ethereum endpoints through coin.
func (c *CoinHost) Endpoints(coin Coin) (*Endpoints, error) {
	var out Endpoints
	err := c.call(http.MethodGet, gwPath(coin, "endpoints"), nil, &out)
	return &out, err
}

// SetEndpoints replaces them.
func (c *CoinHost) SetEndpoints(coin Coin, e Endpoints) error {
	return c.call(http.MethodPut, gwPath(coin, "endpoints"), map[string]interface{}{"rpc_urls": e.RPCURLs, "indexer_url": e.IndexerURL}, nil)
}
