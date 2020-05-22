package bip44wallet

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/stretchr/testify/require"
)

func TestBip44NewWallet(t *testing.T) {
	bip44SkycoinType := bip44.CoinTypeSkycoin
	newBip44Type := bip44.CoinType(1000)

	type expect struct {
		coinType      wallet.CoinType
		bip44CoinType bip44.CoinType
	}

	tt := []struct {
		name           string
		filename       string
		label          string
		seed           string
		opts           []wallet.Option
		seedPassphrase string
		expect         expect
		err            error
	}{
		{
			name:           "skycoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
			},
		},
		{
			name:           "bitcoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeBitcoin),
			},
			expect: expect{
				wallet.CoinTypeBitcoin,
				bip44.CoinTypeBitcoin,
			},
		},
		{
			name:           "skycoin explicit bip44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
				wallet.OptionBip44Coin(&bip44SkycoinType),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
			},
		},
		{
			name:           "skycoin new bip44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
				wallet.OptionBip44Coin(&newBip44Type),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				newBip44Type,
			},
		},
		{
			name:           "no filename",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			err: errors.New("filename not set"),
		},
		{
			name:           "no coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
			},
		},
		{
			name:           "skycoin empty seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           "",
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			err: wallet.ErrMissingSeed,
		},
		{
			name:           "skycoin invalid seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           invalidBip44Seed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			err: errors.New("Mnemonic must have 12, 15, 18, 21 or 24 words"),
		},
		{
			name:           "new coin type, no bi44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType("unknown"),
			},
			err: errors.New("bip44 coin type not set"),
		},
	}

	for _, tc := range tt {
		for _, encrypt := range []bool{false, true} {
			for _, ct := range crypto.TypesInsecure() {
				name := fmt.Sprintf("%s crypto=%v encrypt-%v", tc.name, ct, encrypt)
				opts := tc.opts
				if encrypt {
					opts = append(opts, wallet.OptionEncrypt(true))
					opts = append(opts, wallet.OptionPassword([]byte("pwd")))
					opts = append(opts, wallet.OptionCryptoType(ct))
				}

				t.Run(name, func(t *testing.T) {
					w, err := NewWallet(tc.filename, tc.label, tc.seed, tc.seedPassphrase, opts...)
					require.Equal(t, tc.err, err, fmt.Sprintf("want: %v got: %v", tc.err, err))
					if err != nil {
						return
					}
					require.Equal(t, tc.filename, w.Meta.Filename())
					require.Equal(t, tc.label, w.Meta.Label())
					require.Equal(t, WalletType, w.Meta.Type())
					require.NotEmpty(t, w.Meta.Timestamp())
					require.NotNil(t, w.decoder)
					bip44Coin := w.Bip44Coin()
					require.Equal(t, tc.expect.bip44CoinType, *bip44Coin)
					require.Equal(t, tc.expect.coinType, w.Meta.Coin())

					if encrypt {
						require.Equal(t, ct, w.Meta.CryptoType())
						// confirms that seeds and entry secrets are all empty
						checkNoSensitiveData(t, w)
						return
					}

					require.Equal(t, tc.seed, w.Meta.Seed())
					require.Equal(t, tc.seedPassphrase, w.Meta.SeedPassphrase())
					require.False(t, w.Meta.IsEncrypted())
					require.Empty(t, w.Meta.Secrets())
				})
			}
		}
	}
}

func TestWalletCreateAccount(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)
	require.Equal(t, uint32(0), ai)

	ai, err = w.NewAccount("account2")
	require.Equal(t, uint32(1), ai)

	require.Equal(t, uint32(2), w.accountManager.len())
}

func TestWalletAccountCreateAddresses(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)
	require.Equal(t, uint32(0), ai)

	addrs, err := w.newExternalAddresses(ai, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr := make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinExternalAddresses[:2], addrsStr)

	addrs, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr = make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinChangeAddresses[:2], addrsStr)
}

func TestBip44WalletLock(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.newExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)

	err = w.Lock([]byte("123456"))
	require.NoError(t, err)

	require.Empty(t, w.Seed())
	require.Empty(t, w.SeedPassphrase())
	require.NotEmpty(t, w.Secrets())
	require.True(t, w.IsEncrypted())

	// confirms that no secrets exist in the accounts
	ss := make(wallet.Secrets)
	w.accountManager.packSecrets(ss)
	require.Equal(t, 4, len(ss))
	for k, v := range ss {
		if k == secretBip44AccountPrivateKey {
			require.Empty(t, v)
		} else {
			require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", v)
		}
	}
}

// - Test wallet unlock
func TestBip44WalletUnlock(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin),
		wallet.OptionCryptoType(crypto.CryptoTypeScryptChacha20poly1305Insecure))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.newExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)

	cw := w.Clone().(*Wallet)

	err = cw.Lock([]byte("123456"))
	require.NoError(t, err)

	// generates addresses after locking
	_, err = cw.newExternalAddresses(ai, 2)
	require.NoError(t, err)
	_, err = cw.newChangeAddresses(ai, 3)
	require.NoError(t, err)

	// unlock with wrong password
	_, err = cw.Unlock([]byte("12345"))
	require.Equal(t, errors.New("invalid password"), err)

	// unlock with correct password
	wlt, err := cw.Unlock([]byte("123456"))
	require.NoError(t, err)

	el, err := wlt.EntriesLen(wallet.OptionAccount(ai), wallet.OptionChange(false))
	require.NoError(t, err)
	require.Equal(t, 4, el)

	cl, err := wlt.EntriesLen(wallet.OptionAccount(ai), wallet.OptionChange(true))
	require.NoError(t, err)
	require.Equal(t, 5, cl)

	// confirms that unlocking wallet won't lose data
	require.Empty(t, wlt.Secrets())
	require.False(t, wlt.IsEncrypted())
	require.Equal(t, w.Seed(), wlt.Seed())
	require.Equal(t, w.SeedPassphrase(), wlt.SeedPassphrase())

	// pack the origin wallet's secrets
	originSS := make(wallet.Secrets)
	w.accountManager.packSecrets(originSS)

	// pack the unlocked wallet's secrets
	ss := make(wallet.Secrets)
	wlt.(*Wallet).accountManager.packSecrets(ss)

	// compare these two secrets, they should have the same keys and values
	require.Equal(t, len(originSS)+5, len(ss))
	for k, v := range originSS {
		vv, ok := ss[k]
		require.True(t, ok)
		require.Equal(t, v, vv)
	}
}

func TestBip44WalletNewSerializeDeserialize(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.newExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)

	b, err := w.Serialize()
	require.NoError(t, err)
	t.Log(string(b))

	wlt := Wallet{}
	err = wlt.Deserialize(b)
	require.NoError(t, err)

	// Confirms that serialize/deserialize do not lose meta data
	require.Equal(t, len(w.Meta), len(wlt.Meta))
	for k, v := range wlt.Meta {
		vv, ok := w.Meta[k]
		require.Truef(t, ok, "key:%s", k)
		require.Equal(t, v, vv)
	}

	// confirms that serialize/deserialize do not lose accounts data
	require.Equal(t, w.accountManager.len(), wlt.accountManager.len())
	originSS := make(wallet.Secrets)
	ss := make(wallet.Secrets)
	w.accountManager.packSecrets(originSS)
	wlt.accountManager.packSecrets(ss)

	require.Equal(t, len(originSS), len(ss))
	for k, v := range originSS {
		vv, ok := ss[k]
		require.True(t, ok)
		require.Equal(t, v, vv)
	}
}
