package wallet

import (
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
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
		"2aSysdajziiE2uexbduqV67VhEd2GdLRB8h",
		"2FVWhKpYt5uAavLZmT2PZxV5mhD2pRdqCwd",
		"RHZtm7cf85NDq7SdNR5K5kxwDVgFUPQbRV",
		"o4d7dY58BV7bPMqmgzAEqDKGCRHwQmhB13",
		"EZmQLkLfwhiZ9dAwgVf2bLTvP5wtwwnqVf",
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
		Coin:           meta.CoinTypeSkycoin,
		SeedPassphrase: testSeedPassPhrase,
	}, nil)

	require.NoError(t, err)
	_, err = w.NewExternalAddresses(defaultAccount, 5)
	require.NoError(t, err)

	require.Equal(t, 5, w.EntriesLen())

	_, err = w.NewChangeAddresses(defaultAccount, 2)
	require.NoError(t, err)

	require.Equal(t, 7, w.EntriesLen())

	w1, err := NewBip44Wallet("test1.wlt", Options{
		Seed:           "keep analyst jeans trip erosion race fantasy point spray dinner finger palm",
		Coin:           meta.CoinTypeSkycoin,
		SeedPassphrase: "54321",
	}, nil)

	require.NoError(t, err)

	// Confirms there is one default address
	require.Equal(t, 1, w1.EntriesLen())

	// Do assignment
	*w1 = *w

	// Confirms the entries length is correct
	require.Equal(t, 7, w1.EntriesLen())

	es, err := w1.ExternalEntries(defaultAccount)
	require.NoError(t, err)
	require.Equal(t, 5, len(es))

	// Confirms that the seed is the same
	require.Equal(t, testSeed, w1.Seed())
	// Confirms  that the seed passphrase is the same
	require.Equal(t, testSeedPassPhrase, w1.SeedPassphrase())
}

type mockTxnsFinder struct {
	v map[cipher.Address]bool
}

func (mtf mockTxnsFinder) AddressesActivity(addrs []cipher.Address) ([]bool, error) {
	ret := make([]bool, len(addrs))
	for i, a := range addrs {
		_, ok := mtf.v[a]
		ret[i] = ok
	}
	return ret, nil
}

func TestPeekChangeAddress(t *testing.T) {
	w, err := NewBip44Wallet("test1.wlt", Options{
		Coin:           meta.CoinTypeSkycoin,
		Seed:           testSeed,
		SeedPassphrase: testSeedPassPhrase,
	}, nil)
	require.NoError(t, err)

	cAddrs := getChangeAddrs(t)
	addr, err := w.PeekChangeAddress(mockTxnsFinder{map[cipher.Address]bool{}})
	require.NoError(t, err)
	require.Equal(t, addr, cAddrs[0])

	addr, err = w.PeekChangeAddress(mockTxnsFinder{map[cipher.Address]bool{cAddrs[0]: true}})
	require.NoError(t, err)
	require.Equal(t, addr, cAddrs[1])

	addr, err = w.PeekChangeAddress(mockTxnsFinder{map[cipher.Address]bool{cAddrs[1]: true}})
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
			txnFinder: mockTxnsFinder{map[cipher.Address]bool{}},
		},
		{
			name:        "external addr with txn",
			scanN:       10,
			txnFinder:   mockTxnsFinder{map[cipher.Address]bool{eAddrs[0]: true}},
			expectAddrs: eAddrs[:1],
		},
		{
			name:        "change addr with txn",
			scanN:       10,
			txnFinder:   mockTxnsFinder{map[cipher.Address]bool{cAddrs[0]: true}},
			expectAddrs: cAddrs[:1],
		},
		{
			name:  "external and change addrs with txns",
			scanN: 10,
			txnFinder: mockTxnsFinder{map[cipher.Address]bool{
				eAddrs[0]: true,
				cAddrs[0]: true,
			}},
			expectAddrs: []cipher.Address{eAddrs[0], cAddrs[0]},
		},
		{
			name:  "external and change addrs with txns 2",
			scanN: 10,
			txnFinder: mockTxnsFinder{map[cipher.Address]bool{
				eAddrs[1]: true,
				cAddrs[0]: true,
			}},
			expectAddrs: []cipher.Address{eAddrs[0], eAddrs[1], cAddrs[0]},
		},
		{
			name:  "external and change addrs with txns 3",
			scanN: 10,
			txnFinder: mockTxnsFinder{map[cipher.Address]bool{
				eAddrs[4]: true,
				cAddrs[4]: true,
			}},
			expectAddrs: append(getExternalAddrs(t), getChangeAddrs(t)...),
		},
		{
			name:  "not enough addresses scanned",
			scanN: 3,
			txnFinder: mockTxnsFinder{map[cipher.Address]bool{
				eAddrs[4]: true,
				cAddrs[4]: true,
			}},
		},
		{
			name:  "just enough addresses scanned",
			scanN: 5,
			txnFinder: mockTxnsFinder{map[cipher.Address]bool{
				eAddrs[4]: true,
				cAddrs[4]: true,
			}},
			expectAddrs: append(getExternalAddrs(t), getChangeAddrs(t)...),
		},
		{
			name:  "more addresses scanned",
			scanN: 6,
			txnFinder: mockTxnsFinder{map[cipher.Address]bool{
				eAddrs[4]: true,
				cAddrs[4]: true,
			}},
			expectAddrs: append(getExternalAddrs(t), getChangeAddrs(t)...),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewBip44Wallet("test.wlt", Options{
				Coin:           meta.CoinTypeSkycoin,
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
