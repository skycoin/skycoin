package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestNewAddressesChain(t *testing.T) {
	dir, err := os.MkdirTemp("", "webwallets")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) }) //nolint:errcheck
	coin := bip44.CoinTypeSkycoin
	svc, err := wallet.NewService(wallet.Config{WalletDir: dir, EnableWalletAPI: true, Bip44Coin: &coin})
	require.NoError(t, err)
	w, err := svc.CreateWallet("", wallet.Options{Label: "t", Seed: hardhatMnemonic, Type: wallet.WalletTypeBip44})
	require.NoError(t, err)

	newAddr := func(chain string) (int, []string) {
		form := url.Values{"id": {w.Filename()}, "num": {"1"}}
		if chain != "" {
			form.Set("chain", chain)
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet/newAddress", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handleNewAddresses(newCtx(rec, req), svc)
		var got struct{ Addresses []string }
		_ = json.Unmarshal(rec.Body.Bytes(), &got) //nolint:errcheck
		return rec.Code, got.Addresses
	}

	seed, err := bip39.NewSeed(hardhatMnemonic, "")
	require.NoError(t, err)
	c, err := bip44.NewCoin(seed, coin)
	require.NoError(t, err)
	acct, err := c.Account(0)
	require.NoError(t, err)
	change, err := acct.Change()
	require.NoError(t, err)
	k, err := change.NewPublicChildKey(1) // a new wallet already holds change address 0
	require.NoError(t, err)
	wantChange := cipher.AddressFromPubKey(cipher.MustNewPubKey(k.Key)).String()

	code, addrs := newAddr("change")
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, []string{wantChange}, addrs)

	code, addrs = newAddr("")
	require.Equal(t, http.StatusOK, code)
	require.NotEqual(t, wantChange, addrs[0])

	code, _ = newAddr("sideways")
	require.Equal(t, http.StatusBadRequest, code)
}
