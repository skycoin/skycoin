package wallet

import (
	"errors"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/stretchr/testify/require"
)

const (
	invalidBip44Seed = "invalid bip44 seed"
)

var (
	testSkycoinExternalAddresses = []string{
		"2JBfeo6y6FQn2rCiuhdQ8F1E6bj6rpnHo5U",
		"28Wn9scn3wb5nkScHiTHgNmLjSUS3F2SqAj",
		"qHVbkuuzzxGE6p6CnLY1JxY9ifK1RxjoNS",
		"2WNKEdCvoR8Mv5a7J5bLeE9syq7vHSzACmk",
		"2Z1ZcRWwsyiRqTYLm6VJF914FAE8uhfgmkX",
	}
	testSkycoinChangeAddresses = []string{
		"WFonrBarSSMPwFzcE9CS8vDbqmLjLZaJbT",
		"hiCAv4i9xxtMYXz6Dpwgi1d5Tu1uGzk3Xd",
		"LvV4KEy2pyAmtWXuqYNB2yqAFzN7m6FPme",
		"2Y8fDayjHFSTkVQkCBAAHFuw7gLqZTNmdwr",
		"28JNdPAxvc7yif4gnMhSVLjWh94WTxq9X8y",
	}
	testBitcoinExternalAddresses = []string{
		"162oiGaSrGy8D7iKVC16ga94jDMbqiHmpJ",
		"18fi8QanVYP2aTmiGEPnFozQSvCES3FQun",
		"1DPpoHtKjbCD5aMpisW2i8UqTfbAp7Kfua",
		"1LeT1QohCj6kzeaHbVJZJX3majqKiBct6j",
		"15rLY7XJkeikn7zt7ovYwAN8wzSREAkyU1",
	}
	testBitcoinChangeAddresses = []string{
		"13yByDTkS4eByiCj8k1S817zpa1kPsqMrp",
		"1J434AfF1X5KqvxKmB8H46mabKQZ9dn3pw",
		"1EjD2sUFxFXmqvpw8phbzoPpdBF4QNXLAH",
		"1MPG25DmFoqmVPEbzP6Cd3geg6N99QezSe",
		"1FWfeXmdJ72pKLjWnNSbCe3Lwm91d1b5qr",
	}
)

func TestNewBip44Account(t *testing.T) {
	tt := []struct {
		name           string
		accountName    string
		index          int
		seed           string
		seedPassphrase string
		coinType       CoinType
		err            error
	}{
		{
			name:           "index 0, skycoin, with seed passphrase",
			accountName:    "test",
			index:          0,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
		},
		{
			name:        "index 0, skycoin, without seed passphrase",
			accountName: "test",
			index:       0,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:           "index 1, skycoin, with seed passphrase",
			accountName:    "test",
			index:          1,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
		},
		{
			name:        "index 1, skycoin, without seed passphrase",
			accountName: "test",
			index:       1,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:           "index 2, bitcoin, with seed passphrase",
			accountName:    "test",
			index:          2,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeBitcoin,
		},
		{
			name:        "index 2, bitcoin, without seed passphrase",
			accountName: "test",
			index:       2,
			seed:        testSeed,
			coinType:    CoinTypeBitcoin,
		},
		{
			name:        "index 0x80000000 -1, skycoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000 - 1,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:        "index 0x80000000 -1, bitcoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000 - 1,
			seed:        testSeed,
			coinType:    CoinTypeBitcoin,
		},
		{
			name:        "index 0x80000000, skycoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
			err:         errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:        "index 0x80000000, bitcoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000,
			seed:        testSeed,
			coinType:    CoinTypeBitcoin,
			err:         errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:        "index uint32(-1), skycoin, without seed passphrase",
			accountName: "test",
			index:       -1,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
			err:         errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:        "skycoin, invalid bip44 seed",
			accountName: "test",
			index:       0,
			seed:        "abc def ghe a",
			coinType:    CoinTypeSkycoin,
			err:         errors.New("Mnemonic must have 12, 15, 18, 21 or 24 words"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ba, err := newBip44Account(bip44AccountCreateOptions{
				name:           tc.accountName,
				index:          uint32(tc.index),
				seed:           tc.seed,
				seedPassphrase: testSeedPassphrase,
				coinType:       tc.coinType,
			})
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.accountName, ba.Name)
			require.Equal(t, uint32(tc.index), ba.Index)
			require.Equal(t, tc.coinType, ba.CoinType)
			require.Equal(t, 2, len(ba.Chains))

			externalKey, changeKey, err := makeChainPubKeys(&ba.Account)
			require.NoError(t, err)
			require.Equal(t, *externalKey, ba.Chains[0].PubKey)
			require.Equal(t, *changeKey, ba.Chains[1].PubKey)
			require.Equal(t, 0, len(ba.Chains[0].Entries))
			require.Equal(t, 0, len(ba.Chains[1].Entries))
		})
	}
}

func TestBip44AccountsNewAddresses(t *testing.T) {
	tt := []struct {
		name           string
		coinType       CoinType
		seed           string
		seedPassphrase string
		num            uint32
		chain          uint32
		expectErr      error
		expectAddrs    []string
	}{
		{
			name:           "skycoin, external chain, 1 address",
			coinType:       CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(1),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testSkycoinExternalAddresses[:1],
		},
		{
			name:           "skycoin, change chain, 1 address",
			coinType:       CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(1),
			chain:          bip44.ChangeChainIndex,
			expectAddrs:    testSkycoinChangeAddresses[:1],
		},
		{
			name:           "skycoin, external chain, 2 addresses",
			coinType:       CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testSkycoinExternalAddresses[:2],
		},
		{
			name:           "skycoin, change chain, 2 addresses",
			coinType:       CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ChangeChainIndex,
			expectAddrs:    testSkycoinChangeAddresses[:2],
		},
		{
			name:           "skycoin, external chain, 2 addresses",
			coinType:       CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testSkycoinExternalAddresses[:2],
		},
		{
			name:           "Bitcoin, change chain, 2 addresses",
			coinType:       CoinTypeBitcoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ChangeChainIndex,
			expectAddrs:    testBitcoinChangeAddresses[:2],
		},
		{
			name:           "Bitcoin, external chain, 2 addresses",
			coinType:       CoinTypeBitcoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testBitcoinExternalAddresses[:2],
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			accounts := bip44Accounts{}
			accountIndex, err := accounts.New(bip44AccountCreateOptions{
				name:           "Test",
				coinType:       tc.coinType,
				seed:           tc.seed,
				seedPassphrase: tc.seedPassphrase,
			})
			require.NoError(t, err)

			require.Equal(t, uint32(0), accountIndex)

			addrs, err := accounts.NewAddresses(accountIndex, tc.chain, tc.num)
			require.NoError(t, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.num, uint32(len(addrs)))
			for i, addr := range addrs {
				switch tc.coinType {
				case CoinTypeSkycoin:
					_, err := cipher.DecodeBase58Address(addr.String())
					require.NoError(t, err)
				case CoinTypeBitcoin:
					_, err := cipher.DecodeBase58BitcoinAddress(addr.String())
					require.NoError(t, err)
				}
				require.Equal(t, tc.expectAddrs[i], addr.String())
			}
		})
	}
}
