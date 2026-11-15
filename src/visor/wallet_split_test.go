package visor

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/collection"
)

// newSplitTestVisor returns a visor on a real blockchain and unconfirmed pool, with
// a collection wallet whose addresses own nOuts confirmed outputs of one coin each.
func newSplitTestVisor(t *testing.T, nOuts int) (*Visor, string) {
	db, shutdown := prepareDB(t)
	t.Cleanup(shutdown)
	bc := MakeBlockchain(t, db, GenesisSecret)

	entries, addrs := makeEntries(3)
	ws, err := wallet.NewService(wallet.Config{
		EnableWalletAPI: true,
		CryptoType:      crypto.CryptoTypeScryptChacha20poly1305Insecure,
		WalletDir:       t.TempDir(),
	})
	require.NoError(t, err)
	const wltID = "split.wlt"
	_, err = ws.CreateWallet(wltID, wallet.Options{
		Label: "split",
		Coin:  wallet.CoinTypeSkycoin,
		Type:  wallet.WalletTypeCollection,
	})
	require.NoError(t, err)
	require.NoError(t, ws.Update(wltID, func(w wallet.Wallet) error {
		for _, e := range entries {
			if err := w.(*collection.Wallet).AddEntry(e); err != nil {
				return err
			}
		}
		return nil
	}))

	// Fund the wallet with nOuts outputs in one block. A transaction cannot hold two equal outputs, so the hours differ.
	var txn coin.Transaction
	require.NoError(t, db.View("", func(tx *dbutil.Tx) error {
		uxs, err := bc.Unspent().GetAll(tx)
		require.NoError(t, err)
		require.Len(t, uxs, 1)
		tim, err := bc.Time(tx)
		require.NoError(t, err)
		inHours, err := uxs[0].CoinHours(tim)
		require.NoError(t, err)

		require.NoError(t, txn.PushInput(uxs[0].Hash()))
		var outHours uint64
		for i := range nOuts {
			h := uint64(100 + i) //nolint:gosec // test value
			require.NoError(t, txn.PushOutput(addrs[i%len(addrs)], 1e6, h))
			outHours += h
		}
		require.NoError(t, txn.PushOutput(GenesisAddress, GenesisCoins-uint64(nOuts)*1e6, (inHours-outHours)/2)) //nolint:gosec // test value
		require.NoError(t, txn.SignInputs([]cipher.SecKey{GenesisSecret}))
		return txn.UpdateHeader()
	}))
	ExecuteGenesisSpendTransaction(t, db, bc, txn)

	pool, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)
	v := &Visor{
		Config:      Config{Distribution: params.MainNetDistribution},
		db:          db,
		blockchain:  bc,
		unconfirmed: pool,
		wallets:     ws,
		tf:          mockTxnsFinder{},
	}
	return v, wltID
}

func splitTestParams(to cipher.Address, coins uint64) transaction.Params {
	half := decimal.New(5, -1)
	return transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type:        transaction.HoursSelectionTypeAuto,
			Mode:        transaction.HoursSelectionModeShare,
			ShareFactor: &half,
		},
		To: []coin.TransactionOutput{{Address: to, Coins: coins}},
	}
}

// TestWalletCreateTransactionsSplitInject sends 650 coins from 700 outputs of one
// coin, which needs more inputs than fit in one transaction. Every part is
// injected before any confirms, then each is executed in its own block.
func TestWalletCreateTransactionsSplitInject(t *testing.T) {
	v, wltID := newSplitTestVisor(t, 700)
	to := testutil.MakeAddress()
	p := splitTestParams(to, 650e6)

	_, _, err := v.WalletCreateTransactionSigned(wltID, nil, p, CreateTransactionParams{})
	require.ErrorAs(t, err, new(transaction.TooLargeError))

	txns, inputs, err := v.WalletCreateTransactionsSigned(wltID, nil, p, CreateTransactionParams{})
	require.NoError(t, err)
	require.Greater(t, len(txns), 1)
	require.Len(t, inputs, len(txns))

	seen := make(map[cipher.SHA256]struct{})
	var nIn int
	for _, txn := range txns {
		require.True(t, txn.FullySigned())
		size, err := txn.Size()
		require.NoError(t, err)
		require.LessOrEqual(t, size, params.UserVerifyTxn.MaxTransactionSize)
		for _, h := range txn.In {
			_, dup := seen[h]
			require.False(t, dup, "an input is spent by two parts")
			seen[h] = struct{}{}
		}
		nIn += len(txn.In)
		t.Logf("part: %d inputs, %d outputs, %d bytes", len(txn.In), len(txn.Out), size)
	}
	require.Greater(t, nIn, 650)

	for _, txn := range txns {
		known, _, _, err := v.InjectUserTransaction(*txn)
		require.NoError(t, err)
		require.False(t, known)
	}
	unconfirmed, err := v.GetAllUnconfirmedTransactions()
	require.NoError(t, err)
	require.Len(t, unconfirmed, len(txns))

	// Execute the parts one block at a time. A part that reused an input of an
	// earlier one would fail here.
	for i, txn := range txns {
		var block *coin.Block
		require.NoError(t, v.db.View("", func(tx *dbutil.Tx) error {
			var err error
			block, err = v.blockchain.NewBlock(tx, coin.Transactions{*txn}, GenesisTime+TimeIncrement*uint64(i+2)) //nolint:gosec // test value
			return err
		}))
		sb := coin.SignedBlock{Block: *block, Sig: cipher.MustSignHash(block.HashHeader(), GenesisSecret)}
		require.NoError(t, v.db.Update("", func(tx *dbutil.Tx) error {
			return v.blockchain.ExecuteBlock(tx, &sb)
		}))
	}

	var received uint64
	require.NoError(t, v.db.View("", func(tx *dbutil.Tx) error {
		hashes, err := v.blockchain.Unspent().GetUnspentHashesOfAddrs(tx, []cipher.Address{to})
		require.NoError(t, err)
		uxs, err := v.blockchain.Unspent().GetArray(tx, hashes.Flatten())
		require.NoError(t, err)
		for _, u := range uxs {
			received += u.Body.Coins
		}
		return nil
	}))
	require.Equal(t, uint64(650e6), received)
}

// TestWalletCreateTransactionsSplitFits checks that a send that fits stays one transaction.
func TestWalletCreateTransactionsSplitFits(t *testing.T) {
	v, wltID := newSplitTestVisor(t, 50)
	p := splitTestParams(testutil.MakeAddress(), 20e6)

	txns, _, err := v.WalletCreateTransactions(wltID, p, CreateTransactionParams{})
	require.NoError(t, err)
	require.Len(t, txns, 1)

	txn, _, err := v.WalletCreateTransaction(wltID, p, CreateTransactionParams{})
	require.NoError(t, err)
	require.Equal(t, txn, txns[0])
}

// TestWalletCreateTransactionBeforeInjectReusesInputs shows why sends had to be
// broadcast one at a time. Outputs come from the confirmed unspent pool, less those
// spent by unconfirmed transactions, so a second transaction made before the first
// is injected picks the same outputs.
func TestWalletCreateTransactionBeforeInjectReusesInputs(t *testing.T) {
	v, wltID := newSplitTestVisor(t, 50)
	p := splitTestParams(testutil.MakeAddress(), 20e6)
	wp := CreateTransactionParams{IgnoreUnconfirmed: true}

	a, _, err := v.WalletCreateTransactionSigned(wltID, nil, p, wp)
	require.NoError(t, err)
	b, _, err := v.WalletCreateTransactionSigned(wltID, nil, p, wp)
	require.NoError(t, err)
	require.Equal(t, a.In, b.In)

	_, _, _, err = v.InjectUserTransaction(*a)
	require.NoError(t, err)
	c, _, err := v.WalletCreateTransactionSigned(wltID, nil, p, wp)
	require.NoError(t, err)
	for _, h := range c.In {
		require.NotContains(t, a.In, h)
	}
}

// TestCreateTransactionsSplit checks the address based create, which a wallet
// that keeps its keys elsewhere uses, splits the same way.
func TestCreateTransactionsSplit(t *testing.T) {
	v, wltID := newSplitTestVisor(t, 700)
	w, err := v.wallets.GetWallet(wltID)
	require.NoError(t, err)
	addrsers, err := w.GetAddresses()
	require.NoError(t, err)
	addrs := make([]cipher.Address, len(addrsers))
	for i, a := range addrsers {
		addrs[i] = a.(cipher.Address)
	}
	p := splitTestParams(testutil.MakeAddress(), 650e6)
	p.ChangeAddress = &addrs[0]
	wp := CreateTransactionParams{Addresses: addrs}

	_, _, err = v.CreateTransaction(p, wp)
	require.ErrorAs(t, err, new(transaction.TooLargeError))

	txns, inputs, err := v.CreateTransactions(p, wp)
	require.NoError(t, err)
	require.Greater(t, len(txns), 1)
	require.Len(t, inputs, len(txns))
	seen := make(map[cipher.SHA256]struct{})
	for i, txn := range txns {
		require.Len(t, inputs[i], len(txn.In))
		for _, h := range txn.In {
			_, dup := seen[h]
			require.False(t, dup, "an input is spent by two parts")
			seen[h] = struct{}{}
		}
	}

	small := splitTestParams(testutil.MakeAddress(), 20e6)
	small.ChangeAddress = &addrs[0]
	txns, _, err = v.CreateTransactions(small, wp)
	require.NoError(t, err)
	require.Len(t, txns, 1)
}
