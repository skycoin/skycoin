package pqwallet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/wallet"
)

const testSeed = "pq wallet test seed"

func newTestWallet(t *testing.T, n uint64) *Wallet {
	w, err := NewWallet("test.wlt", "test", testSeed, wallet.OptionGenerateN(n), wallet.OptionCryptoType(crypto.CryptoTypeScryptChacha20poly1305Insecure))
	require.NoError(t, err)
	return w
}

func TestNewWalletDerivesPQAddresses(t *testing.T) {
	w := newTestWallet(t, 3)
	entries, err := w.GetEntries()
	require.NoError(t, err)
	require.Len(t, entries, 3)

	for i, e := range entries {
		sk := pqsig.DeriveSecKey([]byte(testSeed), uint32(i)) //nolint:gosec
		pk, err := sk.PubKey()
		require.NoError(t, err)
		require.Equal(t, cipher.AddressFromPQPubKey(pk), e.Address)
		require.Equal(t, cipher.AddressVersionPQ, e.SkycoinAddress().Version)
		require.Equal(t, []byte(sk), e.Secret[:])
		require.True(t, e.Public.Null())
		require.True(t, e.IsPQ())
		require.NoError(t, e.Verify())
	}

	// The same seed gives the same addresses, one at a time or all at once
	w2 := newTestWallet(t, 1)
	_, err = w2.GenerateAddresses(wallet.OptionGenerateN(2))
	require.NoError(t, err)
	a1, err := w.GetAddresses()
	require.NoError(t, err)
	a2, err := w2.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, a1, a2)

	require.Equal(t, "pq-"+entries[0].Address.String(), w.Fingerprint())
}

func TestNewWalletRejectsOtherCoins(t *testing.T) {
	_, err := NewWallet("test.wlt", "test", testSeed, wallet.OptionCoinType(wallet.CoinTypeBitcoin))
	require.Equal(t, ErrCoinType, err)
}

func TestSerializeRoundTrip(t *testing.T) {
	w := newTestWallet(t, 2)
	b, err := w.Serialize()
	require.NoError(t, err)
	require.NotContains(t, string(b), "public_key")

	w2 := &Wallet{}
	require.NoError(t, w2.Deserialize(b))
	e1, err := w.GetEntries()
	require.NoError(t, err)
	e2, err := w2.GetEntries()
	require.NoError(t, err)
	require.Equal(t, e1, e2)

	// A secret that does not own its address is refused
	other := pqsig.DeriveSecKey([]byte("other"), 0)
	bad := strings.Replace(string(b), e1[0].Secret.Hex(), other.Hex(), 1)
	require.Error(t, (&Wallet{}).Deserialize([]byte(bad)))
}

func TestLockUnlock(t *testing.T) {
	w := newTestWallet(t, 2)
	want, err := w.GetEntries()
	require.NoError(t, err)

	password := []byte("pwd")
	require.NoError(t, w.Lock(password))
	require.True(t, w.IsEncrypted())
	require.Empty(t, w.Seed())
	locked, err := w.GetEntries()
	require.NoError(t, err)
	for _, e := range locked {
		require.True(t, e.Secret.Null())
	}

	// A locked wallet still loads, with its addresses and no secrets
	b, err := w.Serialize()
	require.NoError(t, err)
	require.NotContains(t, string(b), want[0].Secret.Hex())
	require.NoError(t, (&Wallet{}).Deserialize(b))

	_, err = w.Unlock([]byte("wrong"))
	require.Equal(t, wallet.ErrInvalidPassword, err)

	uw, err := w.Unlock(password)
	require.NoError(t, err)
	require.Equal(t, testSeed, uw.Seed())
	got, err := uw.GetEntries()
	require.NoError(t, err)
	require.Equal(t, want, got)
}

type activity map[cipher.Address]bool

func (a activity) AddressesActivity(addrs []cipher.Addresser) ([]bool, error) {
	res := make([]bool, len(addrs))
	for i, x := range addrs {
		res[i] = a[x.(cipher.Address)]
	}
	return res, nil
}

func TestScanAddresses(t *testing.T) {
	ref := newTestWallet(t, 5)
	addrs, err := ref.GetAddresses()
	require.NoError(t, err)

	w := newTestWallet(t, 1)
	found, err := w.ScanAddresses(5, activity{addrs[3].(cipher.Address): true})
	require.NoError(t, err)
	require.Equal(t, addrs[1:4], found)
	n, err := w.EntriesLen()
	require.NoError(t, err)
	require.Equal(t, 4, n)
}

func pqUxOut(t *testing.T, addr cipher.Address, coins uint64) coin.UxOut {
	return coin.UxOut{
		Head: coin.UxHead{Time: 100, BkSeq: 2},
		Body: coin.UxBody{
			SrcTransaction: testutil.RandSHA256(t),
			Address:        addr,
			Coins:          coins,
			Hours:          1000,
		},
	}
}

// inputOrder returns uxs in the order txn spends them
func inputOrder(txn *coin.Transaction, uxs coin.UxArray) coin.UxArray {
	byHash := make(map[cipher.SHA256]coin.UxOut, len(uxs))
	for _, ux := range uxs {
		byHash[ux.Hash()] = ux
	}
	out := make(coin.UxArray, len(txn.In))
	for i, h := range txn.In {
		out[i] = byHash[h]
	}
	return out
}

func TestCreateTransactionSigned(t *testing.T) {
	w := newTestWallet(t, 2)
	entries, err := w.GetEntries()
	require.NoError(t, err)
	a0, a1 := entries[0].SkycoinAddress(), entries[1].SkycoinAddress()

	uxs := coin.UxArray{pqUxOut(t, a0, 2e6), pqUxOut(t, a1, 2e6)}
	to := testutil.MakeAddress()
	txn, _, err := wallet.CreateTransactionSigned(w, transaction.Params{
		HoursSelection: transaction.HoursSelection{Type: transaction.HoursSelectionTypeManual},
		To:             []coin.TransactionOutput{{Address: to, Coins: 3e6, Hours: 10}},
		ChangeAddress:  &a1,
	}, coin.NewAddressUxOuts(uxs), 200)
	require.NoError(t, err)

	require.Equal(t, coin.TxnTypePQ, txn.Type)
	require.Len(t, txn.PQ, 2)
	for i := range txn.Sigs {
		require.True(t, txn.Sigs[i].Null())
	}
	uxs = inputOrder(txn, uxs)
	require.NoError(t, txn.Verify())
	require.NoError(t, txn.VerifyInputSignatures(uxs))

	// The witness survives the raw encoding
	b, err := txn.SerializeWithPQ()
	require.NoError(t, err)
	txn2, err := coin.DeserializeTransactionWithPQ(b)
	require.NoError(t, err)
	require.NoError(t, txn2.VerifyInputSignatures(uxs))
}

func TestSignTransaction(t *testing.T) {
	w := newTestWallet(t, 1)
	e, err := w.GetEntryAt(0)
	require.NoError(t, err)
	addr := e.SkycoinAddress()

	uxs := coin.UxArray{pqUxOut(t, addr, 2e6), pqUxOut(t, addr, 1e6)}
	txn, _, err := wallet.CreateTransaction(w, transaction.Params{
		HoursSelection: transaction.HoursSelection{Type: transaction.HoursSelectionTypeManual},
		To:             []coin.TransactionOutput{{Address: testutil.MakeAddress(), Coins: 3e6, Hours: 10}},
		ChangeAddress:  &addr,
	}, coin.NewAddressUxOuts(uxs), 200)
	require.NoError(t, err)
	require.False(t, txn.IsFullySigned())
	uxs = inputOrder(txn, uxs)

	// Signing one input leaves the original untouched and the other input open
	half, err := wallet.SignTransaction(w, txn, []int{0}, uxs)
	require.NoError(t, err)
	require.Empty(t, txn.PQ)
	require.True(t, half.InputSigned(0))
	require.False(t, half.InputSigned(1))

	_, err = wallet.SignTransaction(w, half, []int{0}, uxs)
	require.Error(t, err)

	full, err := wallet.SignTransaction(w, half, nil, uxs)
	require.NoError(t, err)
	require.True(t, full.IsFullySigned())
	require.Equal(t, coin.TxnTypePQ, full.Type)
	require.NoError(t, full.Verify())
	require.NoError(t, full.VerifyInputSignatures(uxs))
}

func TestCreateTransactionTooLarge(t *testing.T) {
	w := newTestWallet(t, 1)
	e, err := w.GetEntryAt(0)
	require.NoError(t, err)
	addr := e.SkycoinAddress()

	var uxs coin.UxArray
	for i := 0; i < 12; i++ {
		uxs = append(uxs, pqUxOut(t, addr, 1e6))
	}
	spend := func(coins uint64) transaction.Params {
		return transaction.Params{
			HoursSelection: transaction.HoursSelection{Type: transaction.HoursSelectionTypeManual},
			To:             []coin.TransactionOutput{{Address: testutil.MakeAddress(), Coins: coins, Hours: 1}},
			ChangeAddress:  &addr,
		}
	}

	// Twelve 3.7 KB witnesses do not fit the default 32 KiB limit
	_, _, err = wallet.CreateTransactionSigned(w, spend(11e6), coin.NewAddressUxOuts(uxs), 200)
	require.Error(t, err)
	require.Contains(t, err.Error(), "post-quantum inputs")

	txn, _, err := wallet.CreateTransactionSigned(w, spend(3e6), coin.NewAddressUxOuts(uxs), 200)
	require.NoError(t, err)
	size, err := txn.SizeWithPQ()
	require.NoError(t, err)
	require.LessOrEqual(t, size, params.UserVerifyTxn.MaxTransactionSize)
}
