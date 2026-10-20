package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/eth"
	"github.com/skycoin/skycoin/src/wallet"
)

func serveCustom(cc *customCoins, method, target, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if !cc.handle(newCtx(rec, req), strings.TrimPrefix(req.URL.Path, "/api")) {
		rec.WriteHeader(http.StatusTeapot)
	}
	return rec
}

func TestCustomCoins(t *testing.T) {
	fiber := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/health", r.URL.Path)
		_, _ = w.Write([]byte(`{"fiber":{"display_name":"Mycoin","ticker":"MYC","coin_hours_display_name":"Hours"}}`)) //nolint:errcheck
	}))
	defer fiber.Close()

	old := walletDirs
	walletDirs = []string{t.TempDir()}
	defer func() { walletDirs = old }()

	rpc, err := eth.NewRPC([]string{"http://127.0.0.1:1"}, nil)
	require.NoError(t, err)
	ether := &ethHandler{rpc: rpc, indexer: eth.NewIndexer("https://idx.example", nil)}
	newTable := func() *coinTable {
		return &coinTable{
			coins: []*discoveredCoin{{ID: 0, CoinType: "skycoin", CoinSymbol: "SKY"}, {ID: 1, CoinType: "ethereum", CoinSymbol: "ETH"}},
			eth:   map[int]*ethHandler{1: ether},
			wlt:   map[int][]*wallet.Service{},
		}
	}
	table := newTable()
	cc := newCustomCoins(table)

	rec := serveCustom(cc, http.MethodPost, "/api/v1/coins/custom", `{"kind":"fiber","node_url":"`+fiber.URL+`/"}`)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var fc discoveredCoin
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &fc))
	require.Equal(t, 2, fc.ID)
	require.Equal(t, "MYC", fc.CoinSymbol)
	require.Equal(t, "/coin/2", fc.NodeURL)

	rec = serveCustom(cc, http.MethodPost, "/api/v1/coins/custom", `{"kind":"erc20","ticker":"dai","contract":"0x6B175474E89094C44Da98b954EedeAC495271d0F","decimals":18,"icon":"data:image/png;base64,AAAA"}`)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	coin, h, _, _ := table.get(3)
	require.Equal(t, "DAI", coin.CoinSymbol)
	require.Equal(t, "data:image/png;base64,AAAA", coin.IconURL)
	require.Equal(t, 18, h.token.Decimals)
	require.Same(t, rpc, h.rpc)

	for _, bad := range []string{
		`{"kind":"erc20","ticker":"DAI","contract":"0x6B175474E89094C44Da98b954EedeAC495271d0F","decimals":18}`,
		`{"kind":"erc20","ticker":"X","contract":"0x12","decimals":18}`,
		`{"kind":"erc20","ticker":"X","contract":"0x6B175474E89094C44Da98b954EedeAC495271d0F","icon":"data:image/svg+xml;base64,AAAA"}`,
		`{"kind":"fiber","node_url":"http://127.0.0.1:1"}`,
		`{"kind":"btc"}`,
	} {
		rec = serveCustom(cc, http.MethodPost, "/api/v1/coins/custom", bad)
		require.Equal(t, http.StatusBadRequest, rec.Code, bad)
	}

	// A restart adds the saved coins again, at the same indices.
	table2 := newTable()
	require.NoError(t, newCustomCoins(table2).load())
	require.Len(t, table2.coins, 4)
	require.Equal(t, "MYC", table2.coins[2].CoinSymbol)
	require.Equal(t, "DAI", table2.coins[3].CoinSymbol)

	// Removing hides the coin and keeps the others where they are.
	var list []customCoin
	require.NoError(t, json.Unmarshal(serveCustom(cc, http.MethodGet, "/api/v1/coins/custom", "").Body.Bytes(), &list))
	require.Len(t, list, 2)
	rec = serveCustom(cc, http.MethodDelete, "/api/v1/coins/custom?id="+list[0].ID, "")
	require.Equal(t, http.StatusOK, rec.Code)
	gone, _, _, _ := table.get(2)
	require.Nil(t, gone)
	coinsMu.RLock()
	visible := table.visible()
	coinsMu.RUnlock()
	require.Len(t, visible, 3)
	require.Equal(t, 3, visible[2].ID)

	table3 := newTable()
	require.NoError(t, newCustomCoins(table3).load())
	require.Len(t, table3.coins, 3)
	require.Equal(t, "DAI", table3.coins[2].CoinSymbol)
}
