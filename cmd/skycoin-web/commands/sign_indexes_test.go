package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestSignTransactionSignIndexes(t *testing.T) {
	coinType := bip44.CoinTypeSkycoin
	svc, err := wallet.NewService(wallet.Config{WalletDir: t.TempDir(), EnableWalletAPI: true, Bip44Coin: &coinType})
	require.NoError(t, err)
	w, err := svc.CreateWallet("", wallet.Options{Label: "t", Seed: hardhatMnemonic, Type: wallet.WalletTypeBip44})
	require.NoError(t, err)
	entries, err := w.GetEntries()
	require.NoError(t, err)
	mine := entries[0].SkycoinAddress()
	other := testutil.MakeAddress()

	var txn coin.Transaction
	require.NoError(t, txn.PushInput(testutil.RandSHA256(t)))
	require.NoError(t, txn.PushInput(testutil.RandSHA256(t)))
	require.NoError(t, txn.PushOutput(other, 1000, 10))
	txn.Sigs = make([]cipher.Sig, 2)
	require.NoError(t, txn.UpdateHeader())
	raw, err := txn.SerializeHexWithPQ()
	require.NoError(t, err)

	sign := func(indexes []int, addrs ...cipher.Address) (int, string) {
		in := make([]string, len(addrs))
		for i, a := range addrs {
			in[i] = a.String()
		}
		body, err := json.Marshal(map[string]interface{}{
			"wallet_id": w.Filename(), "encoded_transaction": raw, "input_addresses": in, "sign_indexes": indexes,
		})
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v2/wallet/transaction/sign", bytes.NewReader(body))
		handleSignTransaction(newCtx(rec, req), svc)
		var got struct {
			Data struct {
				EncodedTransaction string `json:"encoded_transaction"`
			} `json:"data"`
		}
		_ = json.Unmarshal(rec.Body.Bytes(), &got) //nolint:errcheck
		return rec.Code, got.Data.EncodedTransaction
	}

	code, signed := sign([]int{1}, other, mine)
	require.Equal(t, http.StatusOK, code)
	stxn, err := coin.DeserializeTransactionHexWithPQ(signed)
	require.NoError(t, err)
	require.False(t, stxn.InputSigned(0))
	require.True(t, stxn.InputSigned(1))
	require.NoError(t, cipher.VerifyAddressSignedHash(mine, stxn.Sigs[1], cipher.AddSHA256(stxn.InnerHash, stxn.In[1])))

	// Without sign_indexes every input is signed, so a foreign one fails
	code, _ = sign(nil, other, mine)
	require.Equal(t, http.StatusBadRequest, code)

	code, _ = sign([]int{2}, other, mine)
	require.Equal(t, http.StatusBadRequest, code)
	code, _ = sign([]int{1, 1}, other, mine)
	require.Equal(t, http.StatusBadRequest, code)
}
