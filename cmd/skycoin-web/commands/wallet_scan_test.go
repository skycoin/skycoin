package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet"
)

// fakeSkyNode answers /api/v1/transactions with one transaction paying paid.
func fakeSkyNode(t *testing.T, paid string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/transactions", r.URL.Path)
		var out []readable.TransactionWithStatus
		for _, a := range strings.Split(r.URL.Query().Get("addrs"), ",") {
			if a == paid {
				out = append(out, readable.TransactionWithStatus{Transaction: readable.Transaction{Out: []readable.TransactionOutput{{Address: paid}}}})
			}
		}
		require.NoError(t, json.NewEncoder(w).Encode(out))
	}))
	t.Cleanup(srv.Close)
	return srv
}

// externalAddrs makes a bip44 wallet of hardhatMnemonic with n external addresses.
func externalAddrs(t *testing.T, n uint64) []string {
	coin := bip44.CoinTypeSkycoin
	svc, err := wallet.NewService(wallet.Config{WalletDir: t.TempDir(), EnableWalletAPI: true, Bip44Coin: &coin})
	require.NoError(t, err)
	w, err := svc.CreateWallet("", wallet.Options{Label: "x", Seed: hardhatMnemonic, Type: wallet.WalletTypeBip44, GenerateN: n})
	require.NoError(t, err)
	addrs, err := w.GetAddresses(wallet.OptionExternal())
	require.NoError(t, err)
	out := make([]string, len(addrs))
	for i, a := range addrs {
		out[i] = a.String()
	}
	return out
}

func newScanWallet(t *testing.T) (*wallet.Service, string) {
	coin := bip44.CoinTypeSkycoin
	svc, err := wallet.NewService(wallet.Config{WalletDir: t.TempDir(), EnableWalletAPI: true, Bip44Coin: &coin})
	require.NoError(t, err)
	w, err := svc.CreateWallet("", wallet.Options{Label: "t", Seed: hardhatMnemonic, Type: wallet.WalletTypeBip44})
	require.NoError(t, err)
	return svc, w.Filename()
}

func walletExternal(t *testing.T, svc *wallet.Service, id string) []string {
	w, err := svc.GetWallet(id)
	require.NoError(t, err)
	addrs, err := w.GetAddresses(wallet.OptionExternal())
	require.NoError(t, err)
	out := make([]string, len(addrs))
	for i, a := range addrs {
		out[i] = a.String()
	}
	return out
}

func TestHandleScanAddresses(t *testing.T) {
	want := externalAddrs(t, 5)
	node := fakeSkyNode(t, want[4])
	svc, id := newScanWallet(t)

	form := url.Values{"id": {id}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet/scan", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.True(t, handleWalletAPI(newCtx(rec, req), "/v1/wallet/scan", svc, node.URL))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var got struct{ Addresses []string }
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, want[1:], got.Addresses)
	require.Equal(t, want, walletExternal(t, svc, id))
}

func TestScanAheadAtStart(t *testing.T) {
	want := externalAddrs(t, 3)
	node := fakeSkyNode(t, want[2])
	svc, id := newScanWallet(t)
	table := &coinTable{
		coins: []*discoveredCoin{{CoinType: "skycoin", remoteNodeURL: node.URL}},
		wlt:   map[int][]*wallet.Service{0: {svc}},
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { scanAhead(ctx, table); close(done) }()
	require.Eventually(t, func() bool { return len(walletExternal(t, svc, id)) == 3 }, 5*time.Second, 10*time.Millisecond)
	cancel()
	<-done
	require.Equal(t, want, walletExternal(t, svc, id))
}
