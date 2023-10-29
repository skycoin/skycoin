package bip44wallet

import (
	"errors"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/require"
)

const (
	testSeed           = "attitude coach wet rely typical habit alien security deny imitate spike slab"
	testSeedPassphrase = "pwd"
	invalidBip44Seed   = "invalid bip44 seed"
)

var (
	testSkycoinExternalXPubKey   = "xpub6EMRsT95ntbCFRR2Z6WppnGss1SijAkarfKoRM8tft66tuJh2nt4aJi13S21hUCLZL4cbFBXgHuxipmsS7dj1DW1s4NRup3hzxWfqUdGYv7"
	testSkycoinInternalXPubKey   = "xpub6EMRsT95ntbCGrt4gKqcJTx8rFbBLSvPzxFGfq9DVqFyA6UmDYXAoeTNFs3nmuycUhJG1hC1R5rSbEMK1EiSHotne9hYG55pyPLj8kLuutb"
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
		index          int64
		seed           string
		seedPassphrase string
		coinType       wallet.CoinType
		bip44CoinType  bip44.CoinType
		err            error
	}{
		{
			name:           "index 0, skycoin, with seed passphrase",
			accountName:    "test",
			index:          0,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       wallet.CoinTypeSkycoin,
			bip44CoinType:  bip44.CoinTypeSkycoin,
		},
		{
			name:          "index 0, skycoin, without seed passphrase",
			accountName:   "test",
			index:         0,
			seed:          testSeed,
			coinType:      wallet.CoinTypeSkycoin,
			bip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:           "index 1, skycoin, with seed passphrase",
			accountName:    "test",
			index:          1,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       wallet.CoinTypeSkycoin,
			bip44CoinType:  bip44.CoinTypeSkycoin,
		},
		{
			name:          "index 1, skycoin, without seed passphrase",
			accountName:   "test",
			index:         1,
			seed:          testSeed,
			coinType:      wallet.CoinTypeSkycoin,
			bip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:           "index 2, bitcoin, with seed passphrase",
			accountName:    "test",
			index:          2,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       wallet.CoinTypeBitcoin,
			bip44CoinType:  bip44.CoinTypeBitcoin,
		},
		{
			name:          "index 2, bitcoin, without seed passphrase",
			accountName:   "test",
			index:         2,
			seed:          testSeed,
			coinType:      wallet.CoinTypeBitcoin,
			bip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:          "index 0x80000000 -1, skycoin, without seed passphrase",
			accountName:   "test",
			index:         0x80000000 - 1,
			seed:          testSeed,
			coinType:      wallet.CoinTypeSkycoin,
			bip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:          "index 0x80000000 -1, bitcoin, without seed passphrase",
			accountName:   "test",
			index:         0x80000000 - 1,
			seed:          testSeed,
			coinType:      wallet.CoinTypeBitcoin,
			bip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:          "index 0x80000000, skycoin, without seed passphrase",
			accountName:   "test",
			index:         0x80000000,
			seed:          testSeed,
			coinType:      wallet.CoinTypeSkycoin,
			bip44CoinType: bip44.CoinTypeSkycoin,
			err:           errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:          "index 0x80000000, bitcoin, without seed passphrase",
			accountName:   "test",
			index:         0x80000000,
			seed:          testSeed,
			coinType:      wallet.CoinTypeBitcoin,
			bip44CoinType: bip44.CoinTypeBitcoin,
			err:           errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:          "index uint32(-1), skycoin, without seed passphrase",
			accountName:   "test",
			index:         -1,
			seed:          testSeed,
			coinType:      wallet.CoinTypeSkycoin,
			bip44CoinType: bip44.CoinTypeSkycoin,
			err:           errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:          "skycoin, invalid bip44 seed",
			accountName:   "test",
			index:         0,
			seed:          "abc def ghe a",
			coinType:      wallet.CoinTypeSkycoin,
			bip44CoinType: bip44.CoinTypeSkycoin,
			err:           errors.New("Mnemonic must have 12, 15, 18, 21 or 24 words"),
		},
		{
			name:           "index 0, skycoin, with seed passphrase, no bip44 coin type",
			accountName:    "test",
			index:          0,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       wallet.CoinTypeSkycoin,
			err:            errors.New("newBip44Account missing bip44 coin type"),
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
				bip44CoinType:  &tc.bip44CoinType,
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
		coinType       wallet.CoinType
		seed           string
		seedPassphrase string
		bip44CoinType  bip44.CoinType
		num            uint32
		chain          uint32
		expectErr      error
		expectAddrs    []string
	}{
		{
			name:           "skycoin, external chain, 1 address",
			coinType:       wallet.CoinTypeSkycoin,
			bip44CoinType:  bip44.CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(1),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testSkycoinExternalAddresses[:1],
		},
		{
			name:           "skycoin, change chain, 1 address",
			coinType:       wallet.CoinTypeSkycoin,
			bip44CoinType:  bip44.CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(1),
			chain:          bip44.ChangeChainIndex,
			expectAddrs:    testSkycoinChangeAddresses[:1],
		},
		{
			name:           "skycoin, external chain, 2 addresses",
			coinType:       wallet.CoinTypeSkycoin,
			bip44CoinType:  bip44.CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testSkycoinExternalAddresses[:2],
		},
		{
			name:           "skycoin, change chain, 2 addresses",
			coinType:       wallet.CoinTypeSkycoin,
			bip44CoinType:  bip44.CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ChangeChainIndex,
			expectAddrs:    testSkycoinChangeAddresses[:2],
		},
		{
			name:           "skycoin, external chain, 2 addresses",
			coinType:       wallet.CoinTypeSkycoin,
			bip44CoinType:  bip44.CoinTypeSkycoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testSkycoinExternalAddresses[:2],
		},
		{
			name:           "Bitcoin, change chain, 2 addresses",
			coinType:       wallet.CoinTypeBitcoin,
			bip44CoinType:  bip44.CoinTypeBitcoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ChangeChainIndex,
			expectAddrs:    testBitcoinChangeAddresses[:2],
		},
		{
			name:           "Bitcoin, external chain, 2 addresses",
			coinType:       wallet.CoinTypeBitcoin,
			bip44CoinType:  bip44.CoinTypeBitcoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          bip44.ExternalChainIndex,
			expectAddrs:    testBitcoinExternalAddresses[:2],
		},
		{
			name:           "Bitcoin, invalid chain",
			coinType:       wallet.CoinTypeBitcoin,
			bip44CoinType:  bip44.CoinTypeBitcoin,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			num:            uint32(2),
			chain:          2,
			expectErr:      errors.New("invalid chain index: 2"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			accounts := bip44Accounts{}
			accountIndex, err := accounts.new(bip44AccountCreateOptions{
				name:           "Test",
				coinType:       tc.coinType,
				bip44CoinType:  &tc.bip44CoinType,
				seed:           tc.seed,
				seedPassphrase: tc.seedPassphrase,
			})
			require.NoError(t, err)

			require.Equal(t, uint32(0), accountIndex)

			addrs, err := accounts.newAddresses(accountIndex, tc.chain, tc.num)
			require.Equal(t, tc.expectErr, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.num, uint32(len(addrs)))

			act, err := accounts.account(accountIndex)
			require.NoError(t, err)
			entries := act.Chains[tc.chain].Entries
			require.Equal(t, len(entries), len(addrs))

			for i, addr := range addrs {
				// Confirms that the the secrets key matches the addresses
				secKey := entries[i].Secret

				switch tc.coinType {
				case wallet.CoinTypeSkycoin:
					_, err := cipher.DecodeBase58Address(addr.String())
					require.NoError(t, err)

					addrFromSecKey, err := cipher.AddressFromSecKey(secKey)
					require.NoError(t, err)
					require.Equal(t, addrFromSecKey.String(), addr.String())
				case wallet.CoinTypeBitcoin:
					_, err := cipher.DecodeBase58BitcoinAddress(addr.String())
					require.NoError(t, err)

					addrFromSecKey, err := cipher.BitcoinAddressFromSecKey(secKey)
					require.NoError(t, err)
					require.Equal(t, addrFromSecKey.String(), addr.String())
				}
				require.Equal(t, tc.expectAddrs[i], addr.String())
			}
		})
	}
}

func requireBip44AccountEqual(t *testing.T, a, b *bip44Account) {
	require.Equal(t, a.Account, b.Account)
	require.Equal(t, a.Name, b.Name)
	require.Equal(t, a.Index, b.Index)
	require.Equal(t, a.CoinType, b.CoinType)

	require.Equal(t, len(a.Chains), len(b.Chains))
	aDecoder := wallet.ResolveAddressDecoder(a.CoinType)
	bDecoder := wallet.ResolveAddressDecoder(b.CoinType)
	for j, c := range a.Chains {
		cc := b.Chains[j]
		require.Equal(t, c.PubKey, cc.PubKey)
		require.Equal(t, c.ChainIndex, cc.ChainIndex)
		// verify that the cloned addressFromPubKey func performs the same operation.
		require.Equal(t, aDecoder.AddressFromPubKey(cipher.MustNewPubKey(c.PubKey.Key)), bDecoder.AddressFromPubKey(cipher.MustNewPubKey(c.PubKey.Key)))
		require.Equal(t, len(c.Entries), len(cc.Entries))
		for i, e := range c.Entries {
			ce := cc.Entries[i]
			require.Equal(t, e.Address.String(), ce.Address.String())
			require.Equal(t, e.Public[:], ce.Public[:])
			require.Equal(t, e.Secret[:], ce.Secret[:])
			require.Equal(t, e.ChildNumber, ce.ChildNumber)
			require.Equal(t, e.Change, ce.Change)
		}
	}
}

func TestBip44AccountsClone(t *testing.T) {
	accounts := bip44Accounts{}
	var bip44CoinType = bip44.CoinTypeSkycoin
	accountIndex, err := accounts.new(bip44AccountCreateOptions{
		name:           "Test",
		coinType:       wallet.CoinTypeSkycoin,
		bip44CoinType:  &bip44CoinType,
		seed:           testSeed,
		seedPassphrase: testSeedPassphrase,
	})
	require.NoError(t, err)

	require.Equal(t, uint32(0), accountIndex)

	_, err = accounts.newAddresses(accountIndex, 0, 1)
	require.NoError(t, err)

	cloneAccounts := accounts.clone().(*bip44Accounts)

	require.Equal(t, len(accounts.accounts), len(cloneAccounts.accounts))
	for i, a := range accounts.accounts {
		b := cloneAccounts.accounts[i]
		requireBip44AccountEqual(t, a, b)
	}

	addrs, err := accounts.newAddresses(0, 0, 1)
	require.NoError(t, err)
	caddrs, err := cloneAccounts.newAddresses(0, 0, 1)
	require.NoError(t, err)
	require.Equal(t, addrs, caddrs)
}

func TestBip44AccountErase(t *testing.T) {
	// create a bip44 account
	bip44CoinType := bip44.CoinTypeSkycoin
	a, err := newBip44Account(bip44AccountCreateOptions{
		name:           "Test",
		coinType:       wallet.CoinTypeSkycoin,
		bip44CoinType:  &bip44CoinType,
		seed:           testSeed,
		seedPassphrase: testSeedPassphrase,
	})
	require.NoError(t, err)

	_, err = a.newAddresses(bip44.ExternalChainIndex, 5)
	require.NoError(t, err)

	_, err = a.newAddresses(bip44.ChangeChainIndex, 5)
	require.NoError(t, err)

	// Confirms that the account private key is not empty
	require.NotEmpty(t, a.Account)
	// Confirms that the chains secrets are not empty
	for _, c := range a.Chains {
		for _, e := range c.Entries {
			require.NotEmpty(t, e.Secret)
		}
	}

	// wipes sensitive data
	a.erase()

	// Confirms that the account privatee key is empty
	require.Empty(t, a.Account)
	// Confirms that the secrets in chains are empty
	for _, c := range a.Chains {
		for _, e := range c.Entries {
			require.Equal(t, cipher.SecKey{}, e.Secret)
		}
	}
	// erase multiple times should have no side effect
	for i := 0; i < 10; i++ {
		a.erase()
	}
}

func TestBip44AccountPackSecrets(t *testing.T) {
	// create a bip44 account
	bip44CoinType := bip44.CoinTypeSkycoin
	a, err := newBip44Account(bip44AccountCreateOptions{
		name:           "Test",
		coinType:       wallet.CoinTypeSkycoin,
		bip44CoinType:  &bip44CoinType,
		seed:           testSeed,
		seedPassphrase: testSeedPassphrase,
	})
	require.NoError(t, err)

	_, err = a.newAddresses(bip44.ExternalChainIndex, 5)
	require.NoError(t, err)
	_, err = a.newAddresses(bip44.ChangeChainIndex, 5)
	require.NoError(t, err)

	ss := make(wallet.Secrets)
	a.packSecrets(ss)

	// 11 = 1 account private key + 5 external chain secrets + 5 change chain secrets
	require.Equal(t, 11, len(ss))
	sk, ok := ss.Get(a.accountKeyName())
	require.True(t, ok)

	// Confirms that the Secrets contains account private key
	require.Equal(t, a.Account.String(), sk)
	for _, c := range a.Chains {
		for _, e := range c.Entries {
			s, ok := ss.Get(e.Address.String())
			require.True(t, ok)
			require.Equal(t, e.Secret.Hex(), s)
		}
	}
}

func TestBip44AccountUnpackSecrets(t *testing.T) {
	// create a bip44 account
	bip44CoinType := bip44.CoinTypeSkycoin
	a, err := newBip44Account(bip44AccountCreateOptions{
		name:           "Test",
		coinType:       wallet.CoinTypeSkycoin,
		bip44CoinType:  &bip44CoinType,
		seed:           testSeed,
		seedPassphrase: testSeedPassphrase,
	})
	require.NoError(t, err)

	_, err = a.newAddresses(bip44.ExternalChainIndex, 5)
	require.NoError(t, err)
	_, err = a.newAddresses(bip44.ChangeChainIndex, 5)
	require.NoError(t, err)
	ss := make(wallet.Secrets)
	a.packSecrets(ss)

	ca := a.Clone()

	// erase sensitive data from cloned account
	ca.erase()
	// unpack from Secrets
	err = ca.unpackSecrets(ss)
	require.NoError(t, err)

	// compare the account and cloned account
	requireBip44AccountEqual(t, a, &ca)
}

func TestAccountSyncSecrets(t *testing.T) {
	// create a bip44 account
	bip44CoinType := bip44.CoinTypeSkycoin
	a, err := newBip44Account(bip44AccountCreateOptions{
		name:           "Test",
		coinType:       wallet.CoinTypeSkycoin,
		bip44CoinType:  &bip44CoinType,
		seed:           testSeed,
		seedPassphrase: testSeedPassphrase,
	})
	require.NoError(t, err)

	eAddrs, err := a.newAddresses(bip44.ExternalChainIndex, 5)
	require.NoError(t, err)
	cAddrs, err := a.newAddresses(bip44.ChangeChainIndex, 5)
	require.NoError(t, err)
	ss := make(wallet.Secrets)
	a.packSecrets(ss)
	require.Equal(t, 11, len(ss))

	// wipes secrets
	a.erase()
	nEAddrs, err := a.newAddresses(bip44.ExternalChainIndex, 2)
	nCAddrs, err := a.newAddresses(bip44.ChangeChainIndex, 2)
	require.NoError(t, err)
	require.NoError(t, a.syncSecrets(ss))
	require.Equal(t, 15, len(ss))

	// confirms that all addresses has secrets now
	for _, a := range append(eAddrs, append(cAddrs, append(nEAddrs, nCAddrs...)...)...) {
		_, ok := ss.Get(a.String())
		require.True(t, ok)
	}
}
