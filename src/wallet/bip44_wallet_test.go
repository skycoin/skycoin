package wallet

import (
	"fmt"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/stretchr/testify/require"
)

var (
	testSeed           = "enact seek among recall one save armed parrot license ask giant fog"
	testSeedPassPhrase = "12345"
	changeAddrs        = []string{
		"2g8WtbURh3f4sATvg5W7ryswRSbWzzKFEkb",
		"Jh45qg41xW7PKJCWUJKaKZQnsCs3zGVyq1",
		"uPqY1rh6jY8Zoq3XMjyM8ZD7WwJf5A23DF",
		"ACk9wc1p6uhfzQrQMJwsWtbz6HPYEg2oj7",
		"dB4GuLyay1jQdafN3JyrFUxfNjHB3kALdS",
	}
	externalAddrs = []string{
		"2b5EwW3UAwieMRUogEEbXV4BzGP1AqRFRB6",
		"2aSysdajziiE2uexbduqV67VhEd2GdLRB8h",
		"2FVWhKpYt5uAavLZmT2PZxV5mhD2pRdqCwd",
		"RHZtm7cf85NDq7SdNR5K5kxwDVgFUPQbRV",
		"o4d7dY58BV7bPMqmgzAEqDKGCRHwQmhB13",
	}
)

func getChangeAddrs(t *testing.T) []cipher.Address {
	var addrs []cipher.Address
	for _, addr := range changeAddrs {
		a, err := cipher.DecodeBase58Address(addr)
		require.NoError(t, err)
		addrs = append(addrs, a)
	}
	return addrs
}

func getExternalAddrs(t *testing.T) []cipher.Address {
	var addrs []cipher.Address
	for _, addr := range externalAddrs {
		a, err := cipher.DecodeBase58Address(addr)
		require.NoError(t, err)
		addrs = append(addrs, a)
	}
	return addrs
}

func TestBip44WalletAssign(t *testing.T) {
	w, err := NewBip44Wallet("test.wlt", Options{
		Seed:           testSeed,
		Coin:           CoinTypeSkycoin,
		SeedPassphrase: testSeedPassPhrase,
	}, nil)

	require.NoError(t, err)
	_, err = w.NewExternalAddresses(defaultAccount, 5)
	require.NoError(t, err)

	// 5 added external address + 1 default external + 1 default change address
	require.Equal(t, 7, w.EntriesLen())

	_, err = w.NewChangeAddresses(defaultAccount, 2)
	require.NoError(t, err)

	require.Equal(t, 9, w.EntriesLen())

	w1, err := NewBip44Wallet("test1.wlt", Options{
		Seed:           "keep analyst jeans trip erosion race fantasy point spray dinner finger palm",
		Coin:           CoinTypeSkycoin,
		SeedPassphrase: "54321",
	}, nil)

	require.NoError(t, err)

	// Confirms there are two default addresses, one for external and one for change.
	require.Equal(t, 2, w1.EntriesLen())

	// Do assignment
	*w1 = *w

	// Confirms the entries length is correct
	require.Equal(t, 9, w1.EntriesLen())

	es, err := w1.ExternalEntries(defaultAccount)
	require.NoError(t, err)
	require.Equal(t, 6, len(es))

	// Confirms that the seed is the same
	require.Equal(t, testSeed, w1.Seed())
	// Confirms  that the seed passphrase is the same
	require.Equal(t, testSeedPassPhrase, w1.SeedPassphrase())
}

func TestPeekChangeAddress(t *testing.T) {
	w, err := NewBip44Wallet("test1.wlt", Options{
		Coin:           CoinTypeSkycoin,
		Seed:           testSeed,
		SeedPassphrase: testSeedPassPhrase,
	}, nil)
	require.NoError(t, err)

	cAddrs := getChangeAddrs(t)
	addr, err := w.PeekChangeAddress(mockTxnsFinder{})
	require.NoError(t, err)
	require.Equal(t, addr, cAddrs[0])

	addr, err = w.PeekChangeAddress(mockTxnsFinder{cAddrs[0]: true})
	require.NoError(t, err)
	require.Equal(t, addr, cAddrs[1])

	addr, err = w.PeekChangeAddress(mockTxnsFinder{cAddrs[1]: true})
	require.NoError(t, err)
	require.Equal(t, addr, cAddrs[2])
}

func TestWalletScanAddresses(t *testing.T) {
	eAddrs := getExternalAddrs(t)
	cAddrs := getChangeAddrs(t)

	tt := []struct {
		name        string
		scanN       uint32
		txnFinder   TransactionsFinder
		expectAddrs []cipher.Address
		err         error
	}{
		{
			name:      "no txns",
			scanN:     10,
			txnFinder: mockTxnsFinder{},
		},
		{
			name:        "external addr with txn",
			scanN:       10,
			txnFinder:   mockTxnsFinder{eAddrs[1]: true},
			expectAddrs: eAddrs[1:2],
		},
		{
			name:      "change addr with txn",
			scanN:     10,
			txnFinder: mockTxnsFinder{cAddrs[1]: true},
			// The default change address already exist, thus no more new change addresses will be created
			expectAddrs: cAddrs[1:2],
		},
		{
			name:  "external and change addrs with txns",
			scanN: 10,
			txnFinder: mockTxnsFinder{
				eAddrs[1]: true,
				cAddrs[1]: true,
			},
			expectAddrs: []cipher.Address{eAddrs[1], cAddrs[1]},
		},
		{
			name:  "external and change addrs with txns 2",
			scanN: 10,
			txnFinder: mockTxnsFinder{
				eAddrs[2]: true,
				cAddrs[1]: true,
			},
			expectAddrs: []cipher.Address{eAddrs[1], eAddrs[2], cAddrs[1]},
		},
		{
			name:  "external and change addrs with txns 3",
			scanN: 10,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
			expectAddrs: append(getExternalAddrs(t)[1:], getChangeAddrs(t)[1:]...),
		},
		{
			name:  "not enough addresses scanned",
			scanN: 3,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
		},
		{
			name:  "just enough addresses scanned",
			scanN: 4,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
			expectAddrs: append(getExternalAddrs(t)[1:], getChangeAddrs(t)[1:]...),
		},
		{
			name:  "more addresses scanned",
			scanN: 6,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
			expectAddrs: append(getExternalAddrs(t)[1:], getChangeAddrs(t)[1:]...),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewBip44Wallet("test.wlt", Options{
				Coin:           CoinTypeSkycoin,
				Seed:           testSeed,
				SeedPassphrase: testSeedPassPhrase,
			}, tc.txnFinder)
			require.NoError(t, err)

			addrs, err := w.ScanAddresses(uint64(tc.scanN), tc.txnFinder)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expectAddrs, addrs)
		})
	}
}

func TestBip44WalletUnlock(t *testing.T) {
	tt := []struct {
		name                  string
		options               Options
		password              []byte
		changeWalletFunc      func(w Wallet) error
		expectedMeta          Meta
		expectedExternalAddrN int
		expectedChangeAddrN   int
	}{
		{
			name: "no change",
			options: Options{
				Coin:           CoinTypeSkycoin,
				Seed:           testSeed,
				SeedPassphrase: testSeedPassPhrase,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(w Wallet) error {
				return nil
			},
			expectedExternalAddrN: 1,
			expectedChangeAddrN:   1,
		},
		{
			name: "change label",
			options: Options{
				Coin:           CoinTypeSkycoin,
				Seed:           testSeed,
				SeedPassphrase: testSeedPassPhrase,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(w Wallet) error {
				w.SetLabel("change_label")
				return nil
			},
			expectedMeta:          Meta{MetaLabel: "change_label"},
			expectedExternalAddrN: 1,
			expectedChangeAddrN:   1,
		},
		{
			name: "change filename, no commit",
			options: Options{
				Coin:           CoinTypeSkycoin,
				Seed:           testSeed,
				SeedPassphrase: testSeedPassPhrase,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(w Wallet) error {
				w.(*Bip44Wallet).Meta[MetaFilename] = "filename_changed"
				return nil
			},
			expectedExternalAddrN: 1,
			expectedChangeAddrN:   1,
		},
		{
			name: "new external addresses",
			options: Options{
				Coin:           CoinTypeSkycoin,
				Seed:           testSeed,
				SeedPassphrase: testSeedPassPhrase,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(w Wallet) error {
				_, err := w.(*Bip44Wallet).NewExternalAddresses(defaultAccount, 2)
				return err
			},
			expectedExternalAddrN: 1 + 2,
			expectedChangeAddrN:   1,
		},
		{
			name: "new change addresses",
			options: Options{
				Coin:           CoinTypeSkycoin,
				Seed:           testSeed,
				SeedPassphrase: testSeedPassPhrase,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(w Wallet) error {
				_, err := w.(*Bip44Wallet).NewChangeAddresses(defaultAccount, 2)
				return err
			},
			expectedExternalAddrN: 1,
			expectedChangeAddrN:   1 + 2,
		},
		{
			name: "new external and change addresses",
			options: Options{
				Coin:           CoinTypeSkycoin,
				Seed:           testSeed,
				SeedPassphrase: testSeedPassPhrase,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(w Wallet) error {
				_, err := w.(*Bip44Wallet).NewExternalAddresses(defaultAccount, 2)
				if err != nil {
					return err
				}

				_, err = w.(*Bip44Wallet).NewChangeAddresses(defaultAccount, 2)
				return err
			},
			expectedExternalAddrN: 1 + 2,
			expectedChangeAddrN:   1 + 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewBip44Wallet("test.wlt", tc.options, nil)
			require.NoError(t, err)

			err = w.Lock(tc.password)
			require.NoError(t, err)

			wlt, err := w.Unlock(tc.password)
			require.NoError(t, err)
			require.NoError(t, tc.changeWalletFunc(wlt))

			bw := wlt.(*Bip44Wallet)

			for k, v := range tc.expectedMeta {
				fmt.Println("key:", k, "v:", v)
				require.Equal(t, v, bw.Meta[k])
			}
			el, err := bw.ExternalEntriesLen(defaultAccount)
			require.NoError(t, err)
			cl, err := bw.ChangeEntriesLen(defaultAccount)
			require.NoError(t, err)

			require.Equal(t, tc.expectedExternalAddrN, int(el))
			require.Equal(t, tc.expectedChangeAddrN, int(cl))
		})
	}
}
