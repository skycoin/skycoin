package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/btc"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/eth"
	"github.com/skycoin/skycoin/src/wallet"
)

// fakeBTC is a Bitcoin backend holding value satoshis at every address, with
// one incoming transaction, recording what it is asked to broadcast.
type fakeBTC struct {
	btc.Backend
	value     int64
	broadcast []string
}

func (f *fakeBTC) GetBalance(addrs []string) (map[string]btc.AddressBalance, error) {
	out := map[string]btc.AddressBalance{}
	for _, a := range addrs {
		out[a] = btc.AddressBalance{Confirmed: f.value}
	}
	return out, nil
}

func (f *fakeBTC) ListUnspent(addrs []string) ([]btc.UTXO, error) {
	var out []btc.UTXO
	for i, a := range addrs {
		out = append(out, btc.UTXO{TxID: strings.Repeat("ab", 32), Vout: uint32(i), Address: a, Value: f.value, Height: 100}) //nolint:gosec
	}
	return out, nil
}

func (f *fakeBTC) GetHistory(addrs []string) ([]btc.Transaction, error) {
	return []btc.Transaction{{
		TxID: "btc-in", BlockHeight: 100, Confirmations: 6, Timestamp: 1700000000,
		Inputs:  []btc.TxInput{{Address: "1Stranger", Value: 200000}},
		Outputs: []btc.TxOutput{{Address: addrs[0], Value: 150000}, {Address: "1Stranger", Value: 49000}},
	}}, nil
}

func (f *fakeBTC) BroadcastTransaction(raw string) (string, error) {
	f.broadcast = append(f.broadcast, raw)
	return "btc-sent", nil
}

func (f *fakeBTC) EstimateFee(blocks int) (int64, error) { return int64(30 / blocks), nil }

// fakeEthRPC answers an Ethereum node's calls: every address holds 1 ETH and
// 5 of the token, and sent transactions are recorded.
func fakeEthRPC(t *testing.T, sent *[]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		var result interface{}
		switch req.Method {
		case "eth_chainId":
			result = "0x1"
		case "eth_blockNumber":
			result = "0x10"
		case "eth_getBalance":
			result = "0xde0b6b3a7640000"
		case "eth_call":
			result = "0x" + strings.Repeat("0", 63) + "5"
		case "eth_getTransactionCount":
			result = "0x5"
		case "eth_estimateGas":
			result = "0xfde8"
		case "eth_getBlockByNumber":
			result = map[string]string{"baseFeePerGas": "0x3b9aca00"}
		case "eth_maxPriorityFeePerGas":
			result = "0x3b9aca00"
		case "eth_sendRawTransaction":
			var raw string
			require.NoError(t, json.Unmarshal(req.Params[0], &raw))
			*sent = append(*sent, raw)
			result = "0xfeed"
		default:
			t.Fatalf("unexpected method %s", req.Method)
		}
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc": "2.0", "id": 1, "result": result}))
	}))
}

type coinRig struct {
	h      http.Handler
	btc    *fakeBTC
	ethTxs *[]string
}

// newCoinRig serves the real coin routes and /ui2 over a Bitcoin coin with a
// fake backend and Ether and a token over a fake RPC node and indexer.
func newCoinRig(t *testing.T) coinRig {
	dir := t.TempDir()
	old := walletDirs
	walletDirs = []string{dir}
	t.Cleanup(func() { walletDirs = old })

	service := func(coin bip44.CoinType) *wallet.Service {
		s, err := wallet.NewService(wallet.Config{
			WalletDir: dir, CryptoType: crypto.CryptoTypeScryptChacha20poly1305Insecure, EnableWalletAPI: true, Bip44Coin: &coin,
		})
		require.NoError(t, err)
		return s
	}
	btcSvc, ethSvc := service(bip44.CoinTypeBitcoin), service(bip44.CoinTypeEthereum)

	fb := &fakeBTC{value: 100000}
	var sent []string
	node := fakeEthRPC(t, &sent)
	t.Cleanup(node.Close)
	ix := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"message":"OK","result":[]}`)) //nolint:errcheck
	}))
	t.Cleanup(ix.Close)
	rpc, err := eth.NewRPC([]string{node.URL}, nil)
	require.NoError(t, err)
	indexer := eth.NewIndexer(ix.URL, nil)
	endpoints := newEthEndpoints(rpc, indexer, dir)
	token, err := eth.ParseAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")
	require.NoError(t, err)

	coins := []*discoveredCoin{
		{ID: 0, CoinName: "Skycoin", CoinSymbol: "SKY", CoinType: "skycoin", remoteNodeURL: "http://127.0.0.1:1"},
		{ID: 1, CoinName: "Bitcoin", CoinSymbol: "BTC", CoinType: "bitcoin", ServerWallets: true},
		{ID: 2, CoinName: "Ethereum", CoinSymbol: "ETH", CoinType: "ethereum", Decimals: 18, ServerWallets: true},
		{ID: 3, CoinName: "DAI", CoinSymbol: "DAI", CoinType: "erc20", Decimals: 18, Contract: token.String(), ServerWallets: true},
	}
	table := &coinTable{
		coins: coins,
		btc:   map[int]*btcHandler{1: {backend: fb, wltService: btcSvc}},
		eth: map[int]*ethHandler{
			2: {rpc: rpc, indexer: indexer, wltService: ethSvc, endpoints: endpoints},
			3: {rpc: rpc, indexer: indexer, wltService: ethSvc, endpoints: endpoints, token: &ethToken{Symbol: "DAI", Contract: token, Decimals: 18}},
		},
		wlt: map[int][]*wallet.Service{1: {btcSvc}, 2: {ethSvc}, 3: {ethSvc}},
	}
	mux := http.NewServeMux()
	routeCoins(mux, table, newCustomCoins(table), coins)
	mountUI2(mux)
	return coinRig{h: mux, btc: fb, ethTxs: &sent}
}

func (rig coinRig) do(t *testing.T, coin, method, path string, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if form != nil {
		req = httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.AddCookie(&http.Cookie{Name: "ui2-coin", Value: coin, Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	rr := httptest.NewRecorder()
	rig.h.ServeHTTP(rr, req)
	return rr
}

// createWallet makes a wallet of coin through /ui2 and returns its id.
func (rig coinRig) createWallet(t *testing.T, coin, password string) string {
	t.Helper()
	rr := rig.do(t, coin, http.MethodPost, "/ui2/create", url.Values{
		"mode": {"load"}, "label": {"Rig " + coin}, "seed": {hardhatMnemonic},
		"password": {password}, "password2": {password},
	})
	loc := rr.Header().Get("HX-Redirect")
	require.True(t, strings.HasPrefix(loc, "wallet?id="), "%d %s", rr.Code, rr.Body.String())
	id, err := url.QueryUnescape(strings.TrimPrefix(loc, "wallet?id="))
	require.NoError(t, err)
	return id
}

func TestUI2Coins(t *testing.T) {
	rig := newCoinRig(t)

	body := rig.do(t, "1", http.MethodGet, "/ui2/", nil).Body.String()
	require.Contains(t, body, `data-coin`, "the coin picker")
	require.Contains(t, body, "Bitcoin (BTC)")
	require.Contains(t, body, `class="onboarding"`)
	require.NotContains(t, body, `href="tx"`, "a gateway coin has no Skycoin transaction pages")

	require.Contains(t, rig.do(t, "1", http.MethodGet, "/ui2/load", nil).Body.String(), `name="legacy"`)
	require.NotContains(t, rig.do(t, "2", http.MethodGet, "/ui2/load", nil).Body.String(), `name="legacy"`)

	body = rig.do(t, "0", http.MethodGet, "/ui2/coins", nil).Body.String()
	for _, want := range []string{"Bitcoin (BTC)", "Ethereum (ETH)", "DAI (DAI)", `hx-post="coins"`} {
		require.Contains(t, body, want)
	}
	body = rig.do(t, "0", http.MethodGet, "/ui2/coins/endpoints", nil).Body.String()
	require.Contains(t, body, `name="rpc_urls"`)
}

func TestUI2Bitcoin(t *testing.T) {
	rig := newCoinRig(t)
	id := rig.createWallet(t, "1", "pw")

	body := rig.do(t, "1", http.MethodGet, "/ui2/", nil).Body.String()
	require.Contains(t, body, "Rig 1")
	require.Contains(t, body, "0.002", "100000 satoshis at each of 2 addresses")

	body = rig.do(t, "1", http.MethodGet, "/ui2/wallet?id="+url.QueryEscape(id), nil).Body.String()
	require.Contains(t, body, "bc1q", "a segwit address by default")
	require.Contains(t, body, `data-qr-coin="bitcoin"`)

	body = rig.do(t, "1", http.MethodGet, "/ui2/gw/txns?id="+url.QueryEscape(id), nil).Body.String()
	require.Contains(t, body, "btc-in")
	require.Contains(t, body, "0.0015", "the 150000 satoshis it received")

	body = rig.do(t, "1", http.MethodGet, "/ui2/send?id="+url.QueryEscape(id), nil).Body.String()
	require.Contains(t, body, "5 sat/byte", "the normal preset, 6 blocks")
	require.Contains(t, body, "15 sat/byte", "the fast preset, 2 blocks")

	send := url.Values{"id": {id}, "to": {"bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"}, "amount": {"0.0005"}, "preset": {"fast"}}
	body = rig.do(t, "1", http.MethodPost, "/ui2/send/review", send).Body.String()
	require.Contains(t, body, `hx-post="send/do"`)
	require.Contains(t, body, `name="password"`)
	require.Empty(t, rig.btc.broadcast)

	send.Set("password", "bad")
	require.NotContains(t, rig.do(t, "1", http.MethodPost, "/ui2/send/do", send).Body.String(), "btc-sent")
	require.Empty(t, rig.btc.broadcast)
	send.Set("password", "pw")
	body = rig.do(t, "1", http.MethodPost, "/ui2/send/do", send).Body.String()
	require.Contains(t, body, "btc-sent", body)
	require.Len(t, rig.btc.broadcast, 1)

	send.Set("preset", "nope")
	require.Contains(t, rig.do(t, "1", http.MethodPost, "/ui2/send/review", send).Body.String(), "Choose a fee.")
	send.Set("preset", "fast")
	send.Set("amount", "0.000000001")
	require.Contains(t, rig.do(t, "1", http.MethodPost, "/ui2/send/review", send).Body.String(), "amount")
}

func TestUI2Ethereum(t *testing.T) {
	rig := newCoinRig(t)
	id := rig.createWallet(t, "2", "")

	body := rig.do(t, "2", http.MethodGet, "/ui2/wallet?id="+url.QueryEscape(id), nil).Body.String()
	require.Contains(t, body, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "the hardhat seed's first address")
	require.Contains(t, body, `<span class="ov-big">2</span>`, "1 ETH at each of 2 addresses")

	body = rig.do(t, "3", http.MethodGet, "/ui2/wallet?id="+url.QueryEscape(id), nil).Body.String()
	require.Contains(t, body, "0.00000000000000001", "5 base units of the token at each of 2 addresses")

	body = rig.do(t, "2", http.MethodGet, "/ui2/send?id="+url.QueryEscape(id), nil).Body.String()
	require.Contains(t, body, "ETH")
	send := url.Values{"id": {id}, "to": {"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}, "amount": {"0.001"}, "preset": {"normal"}}
	body = rig.do(t, "2", http.MethodPost, "/ui2/send/do", send).Body.String()
	require.Contains(t, body, "0xfeed", body)
	require.Len(t, *rig.ethTxs, 1)

	rr := rig.do(t, "0", http.MethodPost, "/ui2/coins/endpoints", url.Values{"rpc_urls": {"http://127.0.0.1:8545\n"}, "indexer_url": {"http://127.0.0.1:9000/api"}})
	require.Equal(t, "endpoints?done=endpoints", rr.Header().Get("HX-Redirect"), rr.Body.String())
	require.Contains(t, rig.do(t, "0", http.MethodGet, "/ui2/coins/endpoints", nil).Body.String(), "http://127.0.0.1:9000/api")
}

func TestUI2CustomCoins(t *testing.T) {
	rig := newCoinRig(t)
	fiber := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"fiber":{"display_name":"Mycoin","ticker":"MYC","coin_hours_display_name":"Hours"}}`)) //nolint:errcheck
	}))
	defer fiber.Close()

	rr := rig.do(t, "0", http.MethodPost, "/ui2/coins", url.Values{"kind": {"fiber"}, "node_url": {fiber.URL + "/"}})
	require.Equal(t, "coins?done=coin-added", rr.Header().Get("HX-Redirect"), rr.Body.String())
	rr = rig.do(t, "0", http.MethodPost, "/ui2/coins", url.Values{"kind": {"erc20"}, "ticker": {"usdt"}, "contract": {"0xdAC17F958D2ee523a2206206994597C13D831ec7"}, "decimals": {"6"}})
	require.Equal(t, "coins?done=coin-added", rr.Header().Get("HX-Redirect"), rr.Body.String())
	require.Contains(t, rig.do(t, "0", http.MethodPost, "/ui2/coins", url.Values{"kind": {"erc20"}, "decimals": {"x"}}).Body.String(), "whole number")

	body := rig.do(t, "0", http.MethodGet, "/ui2/coins?done=coin-added", nil).Body.String()
	require.Contains(t, body, "The coin has been added.")
	require.Contains(t, body, "Mycoin (MYC)")
	require.Contains(t, body, "USDT (USDT)")
	m := strings.Index(body, `hx-vals="{&#34;remove&#34;:&#34;`)
	require.Positive(t, m, "an added coin can be removed")
	id := body[m+len(`hx-vals="{&#34;remove&#34;:&#34;`):]
	id = id[:strings.Index(id, "&#34;")]

	rr = rig.do(t, "0", http.MethodPost, "/ui2/coins", url.Values{"remove": {id}})
	require.Equal(t, "coins?done=coin-removed", rr.Header().Get("HX-Redirect"), rr.Body.String())
	require.NotContains(t, rig.do(t, "0", http.MethodGet, "/ui2/coins", nil).Body.String(), "Mycoin (MYC)")
}
