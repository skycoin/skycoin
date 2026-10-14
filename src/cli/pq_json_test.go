package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
)

// The JSON form of a type 1 transaction carries its witness, so a transaction
// read back from it still verifies.
func TestPQTransactionJSONRoundTrip(t *testing.T) {
	sk, err := pqsig.NewSecKey()
	require.NoError(t, err)
	pk, err := sk.PubKey()
	require.NoError(t, err)
	ux := coin.UxOut{Body: coin.UxBody{SrcTransaction: testutil.RandSHA256(t), Address: cipher.AddressFromPQPubKey(pk), Coins: 5e6, Hours: 10}}

	var txn coin.Transaction
	require.NoError(t, txn.PushInput(ux.Hash()))
	require.NoError(t, txn.PushOutput(testutil.MakeAddress(), 5e6, 5))
	txn.InnerHash = txn.HashInner()
	require.NoError(t, txn.SignPQInput(sk, 0))
	require.NoError(t, txn.UpdateHeader())

	rTxn, err := readable.NewTransaction(txn, false)
	require.NoError(t, err)
	b, err := json.Marshal(rTxn)
	require.NoError(t, err)
	require.Contains(t, string(b), `"pq":[{"public_key":"`)

	var back readable.Transaction
	require.NoError(t, json.Unmarshal(b, &back))
	got, err := readableToCreatedTransaction(&back).ToTransaction()
	require.NoError(t, err)
	require.Equal(t, txn.PQ, got.PQ)
	require.NoError(t, got.VerifyInputSignatures(coin.UxArray{ux}))
}
