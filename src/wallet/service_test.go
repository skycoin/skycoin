// build ignore

package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/stretchr/testify/require"
)

func prepareWltDir() string {
	dir, err := ioutil.TempDir("", "wallets")
	if err != nil {
		panic(err)
	}

	return dir
}

func dirIsEmpty(t *testing.T, dir string) {
	f, err := os.Open(dir)
	require.NoError(t, err)
	names, err := f.Readdirnames(1)
	require.Equal(t, io.EOF, err)
	require.Empty(t, names)
}

func TestNewService(t *testing.T) {
	for ct := range cryptoTable {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: false,
			})
			require.NoError(t, err)

			// check if the wallet dir is created
			_, err = os.Stat(dir)
			require.NoError(t, err)

			require.Equal(t, dir, s.walletDirectory)

			require.Equal(t, 1, len(s.wallets))

			// test load wallets
			s, err = NewService(Config{
				WalletDir:        "./testdata",
				CryptoType:       ct,
				DisableWalletAPI: false,
			})
			require.NoError(t, err)

			// check if the dup wallet is loaded
			_, ok1 := s.wallets["test3.1.wlt"]
			_, ok2 := s.wallets["test3.wlt"]
			if ok1 && ok2 {
				t.Fatal("load dup wallet")
			}

			require.Equal(t, 4, len(s.wallets))

		})
	}
}

func TestServiceCreateWalletDisabledWalletAPI(t *testing.T) {
	for ct := range cryptoTable {
		name := fmt.Sprintf("crypto=%v", ct)
		t.Run(name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: true,
			})
			require.NoError(t, err)
			dirIsEmpty(t, dir)

			wltName := "t1.wlt"
			seed := "seed1"
			_, err = s.CreateWallet(wltName, Options{
				Seed: seed,
			})
			dirIsEmpty(t, dir)
			require.Equal(t, ErrWalletAPIDisabled, err)
		})
	}
}

func TestServiceCreateWallet(t *testing.T) {
	for ct := range cryptoTable {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: false,
			})
			require.NoError(t, err)

			wltName := "t1.wlt"
			seed := "seed1"
			w, err := s.CreateWallet(wltName, Options{
				Seed: seed,
			})
			require.NoError(t, err)
			require.NoError(t, w.Validate())

			// create wallet with dup wallet name
			_, err = s.CreateWallet(wltName, Options{Seed: "seed2"})
			require.Equal(t, err, ErrWalletNameConflict)

			// create wallet with dup seed
			dupWlt := "dup_wallet.wlt"
			_, err = s.CreateWallet(dupWlt, Options{
				Seed: seed,
			})
			require.EqualError(t, err, fmt.Sprintf("wallet %s would be duplicate with %v, same seed", dupWlt, wltName))

			// check if the dup wallet is created
			_, ok := s.wallets[dupWlt]
			require.False(t, ok)

			_, err = os.Stat(filepath.Join(dir, dupWlt))
			require.True(t, os.IsNotExist(err))
		})
	}
}

func TestServiceLoadWallet(t *testing.T) {
	// Prepare addresss
	seed := "seed"
	_, seckeys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), 10)
	var addrs []cipher.Address
	for _, s := range seckeys {
		addrs = append(addrs, cipher.AddressFromSecKey(s))
	}

	tt := []struct {
		name          string
		opts          Options
		scanN         uint64
		bg            BalanceGetter
		err           error
		expectAddrNum int
		expectAddrs   []cipher.Address
	}{
		{
			"raw wallet address=1",
			Options{
				Seed:  "seed",
				Label: "wallet",
			},
			5,
			mockBalanceGetter{
				addrs[0]: BalancePair{Confirmed: Balance{Coins: 1e6, Hours: 100}},
			},
			nil,
			1,
			addrs[:1],
		},
		{
			"raw wallet address=2",
			Options{
				Seed:  "seed",
				Label: "wallet",
			},
			5,
			mockBalanceGetter{
				addrs[1]: BalancePair{Confirmed: Balance{Coins: 1e6, Hours: 100}},
			},
			nil,
			2,
			addrs[:2],
		},
		{
			"encrypted wallet address=1",
			Options{
				Seed:     "seed",
				Label:    "wallet",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			5,
			mockBalanceGetter{
				addrs[0]: BalancePair{Confirmed: Balance{Coins: 1e6, Hours: 100}},
			},
			nil,
			1,
			addrs[:1],
		},
		{
			"encrypted wallet address=2",
			Options{
				Seed:     "seed",
				Label:    "wallet",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			5,
			mockBalanceGetter{
				addrs[1]: BalancePair{Confirmed: Balance{Coins: 1e6, Hours: 100}},
			},
			nil,
			2,
			addrs[:2],
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:        dir,
					CryptoType:       ct,
					DisableWalletAPI: false,
				})
				require.NoError(t, err)
				wltName := newWalletFilename()

				w, err := s.loadWallet(wltName, tc.opts, tc.scanN, tc.bg)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.Len(t, w.Entries, tc.expectAddrNum)
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, w.Entries[i].Address)
				}
			})
		}
	}

}

func TestServiceNewAddress(t *testing.T) {
	seed := []byte("seed")
	// Generate adddresses from the seed
	var addrs []cipher.Address
	_, seckeys := cipher.GenerateDeterministicKeyPairsSeed(seed, 10)
	for _, s := range seckeys {
		addrs = append(addrs, cipher.AddressFromSecKey(s))
	}

	tt := []struct {
		name          string
		opts          Options
		n             uint64
		pwd           []byte
		expectAddrNum int
		expectAddrs   []cipher.Address
		expectErr     error
	}{
		{
			"encrypted=false addresses=0",
			Options{
				Label: "label",
				Seed:  string(seed),
			},
			0,
			nil,
			0,
			nil, // CreateWallet will generate a default address, so check from new address
			nil,
		},
		{
			"encrypted=false addresses=1",
			Options{
				Label: "label",
				Seed:  string(seed),
			},
			2,
			nil,
			2,
			addrs[1:3], // CreateWallet will generate a default address, so check from new address
			nil,
		},
		{
			"encrypted=false addresses=2",
			Options{
				Label: "label",
				Seed:  string(seed),
			},
			2,
			nil,
			2,
			addrs[1:3], // CreateWallet will generate a default address, so check from new address
			nil,
		},
		{
			"encrypted=true addresses=1",
			Options{
				Label:    "label",
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			1,
			[]byte("pwd"),
			1,
			addrs[1:2], // CreateWallet will generate a default address, so check from new address
			nil,
		},
		{
			"encrypted=true addresses=2",
			Options{
				Label:    "label",
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			2,
			[]byte("pwd"),
			2,
			addrs[1:3], // CreateWallet will generate a default address, so check from new address
			nil,
		},
		{
			"encrypted=true wrong password",
			Options{
				Label:    "label",
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			1,
			[]byte("wrong password"),
			1,
			nil,
			ErrInvalidPassword,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:        dir,
					CryptoType:       ct,
					DisableWalletAPI: false,
				})
				require.NoError(t, err)

				wltName := newWalletFilename()

				w, err := s.CreateWallet(wltName, tc.opts)
				require.NoError(t, err)

				naddrs, err := s.NewAddresses(w.Filename(), tc.pwd, tc.n)
				require.Equal(t, tc.expectErr, err)
				if err != nil {
					return
				}

				require.Len(t, naddrs, tc.expectAddrNum)
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, naddrs[i])
				}

				// Check the wallet again
				w, ok := s.wallets[wltName]
				require.True(t, ok)
				require.Len(t, w.Entries, int(tc.n+1))

				// Wallet has a default address, so need to start from the second address
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, w.Entries[i+1].Address)
				}

				// Load wallet from file and check
				_, err = os.Stat(filepath.Join(dir, w.Filename()))
				require.NoError(t, err)

				lw, err := Load(filepath.Join(dir, w.Filename()))
				require.NoError(t, err)
				require.Equal(t, lw, w)

				// Wallet doesn't exist
				_, err = s.NewAddresses("wallet_not_exist.wlt", tc.pwd, 1)
				require.Equal(t, ErrWalletNotExist, err)
			})
		}
	}
}

func TestServiceNewAddressDisabledWalletAPI(t *testing.T) {
	for ct := range cryptoTable {
		name := fmt.Sprintf("crypto=%v", ct)
		t.Run(name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: true,
			})
			require.NoError(t, err)
			dirIsEmpty(t, dir)

			require.Empty(t, s.wallets)
			addrs, err := s.NewAddresses("", nil, 1)
			require.Equal(t, ErrWalletNotExist, err)
			require.Equal(t, 0, len(addrs))
		})
	}
}

func TestServiceGetAddressDisabledWalletAPI(t *testing.T) {
	for ct := range cryptoTable {
		name := fmt.Sprintf("crypto=%v", ct)
		t.Run(name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: true,
			})
			require.NoError(t, err)
			dirIsEmpty(t, dir)

			require.Empty(t, s.wallets)
			addrs, err := s.GetAddresses("")
			require.Equal(t, ErrWalletNotExist, err)
			require.Equal(t, 0, len(addrs))
		})
	}
}

func TestServiceGetAddress(t *testing.T) {
	for ct := range cryptoTable {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: false,
			})
			require.NoError(t, err)

			// get the defaut wallet
			var w *Wallet
			for _, w = range s.wallets {
			}

			addrs, err := s.GetAddresses(w.Filename())
			require.NoError(t, err)
			require.Equal(t, 1, len(addrs))

			// test none exist wallet
			notExistID := "not_exist_id.wlt"
			_, err = s.GetAddresses(notExistID)
			require.Equal(t, ErrWalletNotExist, err)
		})
	}
}

func TestServiceGetWalletDisabledWalletAPI(t *testing.T) {
	for ct := range cryptoTable {
		name := fmt.Sprintf("crypto=%v", ct)
		t.Run(name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: true,
			})
			require.NoError(t, err)
			dirIsEmpty(t, dir)

			require.Empty(t, s.wallets)
			w, err := s.GetWallet("")
			require.Equal(t, ErrWalletNotExist, err)
			var emptyW *Wallet
			require.Equal(t, w, emptyW)
		})
	}
}

func TestServiceGetWallet(t *testing.T) {
	for ct := range cryptoTable {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			dir := prepareWltDir()

			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: false,
			})
			require.NoError(t, err)

			// Get the defaut wallet
			var w *Wallet
			for _, w = range s.wallets {
			}
			require.NoError(t, err)

			w1, err := s.GetWallet(w.Filename())
			require.NoError(t, err)

			// Check if change original wallet would change the returned wallet
			w.setLabel("new_label")

			require.NotEqual(t, "new_label", w1.Label())

			// Get wallet doesn't exist
			wltName := "does_not_exist.wlt"
			_, err = s.GetWallet(wltName)
			require.Equal(t, ErrWalletNotExist, err)
		})
	}
}

func TestServiceGetWallets(t *testing.T) {
	for ct := range cryptoTable {
		dir := prepareWltDir()
		s, err := NewService(Config{
			WalletDir:        dir,
			CryptoType:       ct,
			DisableWalletAPI: false,
		})
		require.NoError(t, err)

		var wallets []*Wallet
		// Get the default wallet
		var w1 *Wallet
		for _, w1 = range s.wallets {
		}
		wallets = append(wallets, w1)

		// Create a new wallet
		wltName := newWalletFilename()
		w2, err := s.CreateWallet(wltName, Options{
			Seed:  "seed",
			Label: "label",
		})
		require.NoError(t, err)
		wallets = append(wallets, w2)

		ws := s.GetWallets()
		for _, w := range wallets {
			ww, ok := ws[w.Filename()]
			require.True(t, ok)
			require.Equal(t, w, ww)
		}
	}
}

func TestServiceReloadWalletsDisabledWalletAPI(t *testing.T) {
	for ct := range cryptoTable {
		name := fmt.Sprintf("crypto=%v", ct)
		t.Run(name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: true,
			})
			require.NoError(t, err)
			dirIsEmpty(t, dir)

			err = s.ReloadWallets()
			require.Equal(t, ErrWalletAPIDisabled, err)
		})
	}
}

func TestServiceReloadWallets(t *testing.T) {
	for ct := range cryptoTable {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			dir := prepareWltDir()

			s, err := NewService(Config{
				WalletDir:        dir,
				CryptoType:       ct,
				DisableWalletAPI: false,
			})
			require.NoError(t, err)

			var w *Wallet
			for _, w = range s.wallets {
			}

			defaultWltID := w.Filename()

			var defaultAddr string
			for defaultAddr = range s.firstAddrIDMap {
				break
			}

			wltName := "t1.wlt"
			w1, err := s.CreateWallet(wltName, Options{Seed: "seed1"})
			require.NoError(t, err)

			fmt.Println(dir, w1.Filename())

			err = s.ReloadWallets()
			require.NoError(t, err)

			// check if create dup wallet will return error
			_, ok := s.wallets[defaultWltID]
			require.True(t, ok)

			_, ok = s.wallets["t1.wlt"]
			require.True(t, ok)

			// check if the first address of each wallet is reloaded
			_, ok = s.firstAddrIDMap[defaultAddr]
			require.True(t, ok)

			_, ok = s.firstAddrIDMap[w1.Entries[0].Address.String()]
			require.True(t, ok)

		})
	}
}

func TestServiceCreateAndSignTx(t *testing.T) {
	headTime := time.Now().UTC().Unix()
	seed := []byte("seed")

	// Generate first keys
	_, secKeys := cipher.GenerateDeterministicKeyPairsSeed(seed, 1)
	secKey := secKeys[0]
	addr := cipher.AddressFromSecKey(secKey)

	// Create unspent outptus
	var uxouts []coin.UxOut
	addrs := []cipher.Address{}
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey, 2e6, 100)
		uxouts = append(uxouts, uxout)

		p, _ := cipher.GenerateKeyPair()
		a := cipher.AddressFromPubKey(p)
		addrs = append(addrs, a)
	}

	// Create unspent outputs with no hours
	var uxoutsNoHours []coin.UxOut
	addrsNoHours := []cipher.Address{}
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey, 2e6, 100)
		uxout.Body.Hours = 0
		uxout.Head.Time = uint64(headTime)
		uxoutsNoHours = append(uxoutsNoHours, uxout)

		p, _ := cipher.GenerateKeyPair()
		a := cipher.AddressFromPubKey(p)
		addrsNoHours = append(addrsNoHours, a)
	}

	tt := []struct {
		name     string
		opts     Options
		pwd      []byte
		unspents []coin.UxOut
		vld      Validator
		coins    uint64
		dest     cipher.Address
		err      error
	}{
		{
			"encrypted=false has change=no",
			Options{
				Seed: string(seed),
			},
			nil,
			uxouts[:],
			&dummyValidator{
				ok: false,
			},
			2e6,
			addrs[0],
			nil,
		},
		{
			"encrypted=true has change=no",
			Options{
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			uxouts[:],
			&dummyValidator{
				ok: false,
			},
			2e6,
			addrs[0],
			nil,
		},
		{
			"encrypted=false has change=yes",
			Options{
				Seed: string(seed),
			},
			nil,
			uxouts[:],
			&dummyValidator{
				ok: false,
			},
			1e6,
			addrs[0],
			nil,
		},
		{
			"encrypted=false has unconfirmed spending transaction",
			Options{
				Seed: string(seed),
			},
			nil,
			uxouts[:],
			&dummyValidator{
				ok: true,
			},
			2e6,
			addrs[0],
			errors.New("please spend after your pending transaction is confirmed"),
		},
		{
			"encrypted=false unconfirmed spend failed",
			Options{
				Seed: string(seed),
			},
			nil,
			uxouts[:],
			&dummyValidator{
				ok:  false,
				err: errors.New("fail intentionally"),
			},
			2e6,
			addrs[0],
			errors.New("checking unconfirmed spending failed: fail intentionally"),
		},
		{
			"encrypted=false spend zero",
			Options{
				Seed: string(seed),
			},
			nil,
			uxouts[:],
			&dummyValidator{
				ok: false,
			},
			0,
			addrs[0],
			errors.New("zero spend amount"),
		},
		{
			"encrypted=false spend fractional coins",
			Options{
				Seed: string(seed),
			},
			nil,
			uxouts[:],
			&dummyValidator{
				ok: false,
			},
			1e3,
			addrs[0],
			nil,
		},
		{
			"encrypted=false not enough confirmed coins",
			Options{
				Seed: string(seed),
			},
			nil,
			uxouts[:],
			&dummyValidator{
				ok: false,
			},
			100e6,
			addrs[0],
			ErrInsufficientBalance,
		},
		{
			"encrypted=false no coin hours in inputs",
			Options{
				Seed: string(seed),
			},
			nil,
			uxoutsNoHours[:],
			&dummyValidator{
				ok: false,
			},
			1e6,
			addrsNoHours[0],
			fee.ErrTxnNoFee,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				unspents := &dummyUnspentGetter{
					addrUnspents: coin.AddressUxOuts{
						addr: tc.unspents,
					},
					unspents: map[cipher.SHA256]coin.UxOut{},
				}

				for _, ux := range tc.unspents {
					unspents.unspents[ux.Hash()] = ux
				}

				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:        dir,
					CryptoType:       ct,
					DisableWalletAPI: false,
				})
				require.NoError(t, err)

				wltName := newWalletFilename()

				w, err := s.CreateWallet(wltName, tc.opts)
				require.NoError(t, err)

				tx, err := s.CreateAndSignTransaction(w.Filename(), tc.pwd, tc.vld, unspents, uint64(headTime), tc.coins, tc.dest)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				// check the IN of tx
				for _, inUxid := range tx.In {
					_, ok := unspents.unspents[inUxid]
					require.True(t, ok)
				}

				err = tx.Verify()
				require.NoError(t, err)
			})
		}
	}
}

func TestServiceUpdateWalletLabel(t *testing.T) {
	tt := []struct {
		name          string
		wltName       string
		opts          Options
		updateWltName string
		label         string
		err           error
	}{
		{
			"ok",
			"t.wlt",
			Options{
				Seed:  "seed",
				Label: "label",
			},
			"t.wlt",
			"new-label",
			nil,
		},
		{
			"wallet doesn't exist",
			"t.wlt",
			Options{
				Seed:  "seed",
				Label: "label",
			},
			"t1.wlt",
			"new-label",
			ErrWalletNotExist,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			t.Run(tc.name, func(t *testing.T) {
				// Create the wallet service
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:        dir,
					CryptoType:       ct,
					DisableWalletAPI: false,
				})
				require.NoError(t, err)

				// Create a new wallet
				w, err := s.CreateWallet(tc.wltName, tc.opts)
				require.NoError(t, err)

				err = s.UpdateWalletLabel(tc.updateWltName, tc.label)
				require.Equal(t, tc.err, err)

				if err != nil {
					return
				}

				nw, err := s.GetWallet(w.Filename())
				require.NoError(t, err)
				require.NotEqual(t, w.Label(), nw.Label())

				require.Equal(t, tc.label, nw.Label())
			})
		}
	}
}

func TestServiceEncryptWallet(t *testing.T) {
	tt := []struct {
		name       string
		wltName    string
		opts       Options
		encWltName string
		pwd        []byte
		err        error
	}{
		{
			"ok",
			"t.wlt",
			Options{
				Seed: "seed",
			},
			"t.wlt",
			[]byte("pwd"),
			nil,
		},
		{
			"wallet doesn't exist",
			"t.wlt",
			Options{
				Seed: "seed",
			},
			"t2.wlt",
			nil,
			ErrWalletNotExist,
		},
		{
			"wallet already encrypted",
			"t.wlt",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			"t.wlt",
			[]byte("pwd"),
			ErrWalletEncrypted,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				// Create the wallet service
				s, err := NewService(Config{
					WalletDir:        dir,
					CryptoType:       ct,
					DisableWalletAPI: false,
				})
				require.NoError(t, err)

				// Create a new wallet
				w, err := s.CreateWallet(tc.wltName, tc.opts)
				require.NoError(t, err)

				// Encrypt the wallet
				err = s.EncryptWallet(tc.encWltName, tc.pwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				encWlt, err := s.getWallet(tc.encWltName)
				require.NoError(t, err)

				// Check the encrypted wallet
				require.True(t, encWlt.IsEncrypted())
				require.Equal(t, cipher.SecKey{}, encWlt.Entries[0].Secret)
				require.NotEqual(t, w.seed(), encWlt.seed())
				require.NotEqual(t, w.lastSeed(), encWlt.lastSeed())

				// Check the decrypted seeds
				decWlt, err := encWlt.unlock(tc.pwd)
				require.NoError(t, err)
				require.Equal(t, w.seed(), decWlt.seed())
				require.Equal(t, w.lastSeed(), decWlt.lastSeed())

				// Check if the wallet file exist
				path := filepath.Join(dir, w.Filename())
				_, err = os.Stat(path)
				require.True(t, !os.IsNotExist(err))

				// Check if the backup wallet file, which should not exist
				bakPath := path + ".bak"
				_, err = os.Stat(bakPath)
				require.True(t, os.IsNotExist(err))
			})
		}
	}
}

func TestServiceDecryptWallet(t *testing.T) {
	tt := []struct {
		name           string
		wltName        string
		opts           Options
		decryptWltName string
		password       []byte
		err            error
	}{
		{
			"ok",
			"test.wlt",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			"test.wlt",
			[]byte("pwd"),
			nil,
		},
		{
			"wallet not exist",
			"test.wlt",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			"t.wlt",
			[]byte("pwd"),
			ErrWalletNotExist,
		},
		{
			"wallet not encrypted",
			"test.wlt",
			Options{
				Seed: "seed",
			},
			"test.wlt",
			[]byte("pwd"),
			ErrWalletNotEncrypted,
		},
		{
			"invalid password",
			"test.wlt",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			"test.wlt",
			[]byte("wrong password"),
			ErrInvalidPassword,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {

			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:        dir,
					CryptoType:       ct,
					DisableWalletAPI: false,
				})
				require.NoError(t, err)

				_, err = s.CreateWallet(tc.wltName, tc.opts)
				require.NoError(t, err)

				err = s.DecryptWallet(tc.decryptWltName, tc.password)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				verifyDecryptedWlt := func(wlt *Wallet) {
					// Checks the "encrypted" meta info
					require.False(t, wlt.IsEncrypted())
					// Checks the seed
					require.Equal(t, tc.opts.Seed, wlt.seed())
					// Checks the last seed
					entryNum := len(wlt.Entries)
					lsd, seckeys := cipher.GenerateDeterministicKeyPairsSeed([]byte(wlt.seed()), entryNum)
					require.NoError(t, err)
					require.Equal(t, hex.EncodeToString(lsd), wlt.lastSeed())

					// Checks the entries
					for i := range seckeys {
						a := cipher.AddressFromSecKey(seckeys[i])
						require.Equal(t, a, wlt.Entries[i].Address)
						require.Equal(t, seckeys[i], wlt.Entries[i].Secret)
					}
				}

				// Checks the decrypted wallet in service
				w, err := s.getWallet(tc.wltName)
				require.NoError(t, err)
				verifyDecryptedWlt(w)

				// Checks the existence of the wallet file
				fn := filepath.Join(dir, tc.wltName)
				_, err = os.Stat(fn)
				require.True(t, !os.IsNotExist(err))

				// Loads wallet from the file and check if it's decrypted
				w1, err := Load(fn)
				require.NoError(t, err)
				verifyDecryptedWlt(w1)
			})

		}
	}
}

func TestServiceScanAheadWalletAddresses(t *testing.T) {
	bg := make(mockBalanceGetter, len(addrsOfSeed1))
	addrs := fromAddrString(t, addrsOfSeed1)
	for _, a := range addrs {
		bg[a] = BalancePair{}
	}

	type exp struct {
		err              error
		seed             string
		lastSeed         string
		entryNum         int
		confirmedBalance uint64
		predictedBalance uint64
	}

	tt := []struct {
		name      string
		opts      Options
		pwd       []byte
		scanN     uint64
		balGetter BalanceGetter
		expect    exp
	}{
		{
			"no coins and scan 0, unencrypted",
			Options{
				Seed: "seed1",
			},
			nil,
			0,
			bg,
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"no coins and scan 0, encrypted",
			Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			0,
			bg,
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"no coins and scan 1, unencrypted",
			Options{
				Seed: "seed1",
			},
			nil,
			1,
			bg,
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"no coins and scan 1, encrypted",
			Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			1,
			bg,
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"no coins and scan 10, unencrypted",
			Options{
				Seed: "seed1",
			},
			nil,
			10,
			bg,
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 5, unencrypted",
			Options{
				Seed: "seed1",
			},
			nil,
			5,
			mockBalanceGetter{
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[5],
				entryNum:         5,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 5, encrypted",
			Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			5,
			mockBalanceGetter{
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[5],
				entryNum:         5,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 4, unencrypted",
			Options{
				Seed: "seed1",
			},
			nil,
			5,
			mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 4, encrypted",
			Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			5,
			mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 4 have 6, unencrypted",
			Options{
				Seed: "seed1",
			},
			nil,
			5,
			mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[6]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			"confirmed and predicted, unencrypted",
			Options{
				Seed: "seed1",
			},
			nil,
			5,
			mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			"confirmed and predicted, encrypted",
			Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			5,
			mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			"confirmed and predicted, wrong password",
			Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("wrong password"),
			5,
			mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				err: ErrInvalidPassword,
			},
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:        dir,
					CryptoType:       ct,
					DisableWalletAPI: false,
				})
				require.NoError(t, err)

				wltName := newWalletFilename()
				w, err := s.CreateWallet(wltName, tc.opts)
				require.NoError(t, err)

				require.NoError(t, w.Validate())

				w1, err := s.ScanAheadWalletAddresses(wltName, tc.pwd, tc.scanN, tc.balGetter)
				require.Equal(t, tc.expect.err, err)
				if err != nil {
					return
				}

				require.Len(t, w1.Entries, tc.expect.entryNum)
				for i := range w1.Entries {
					require.Equal(t, addrsOfSeed1[i], w1.Entries[i].Address.String())
				}
			})
		}
	}
}

type dummyValidator struct {
	ok  bool
	err error
}

func (dvld dummyValidator) HasUnconfirmedSpendTx(addr []cipher.Address) (bool, error) {
	return dvld.ok, dvld.err
}

type dummyUnspentGetter struct {
	addrUnspents coin.AddressUxOuts
	unspents     map[cipher.SHA256]coin.UxOut
}

func (dug dummyUnspentGetter) GetUnspentsOfAddrs(addrs []cipher.Address) coin.AddressUxOuts {
	return dug.addrUnspents
}

func (dug dummyUnspentGetter) Get(uxid cipher.SHA256) (coin.UxOut, bool) {
	uxout, ok := dug.unspents[uxid]
	return uxout, ok
}

func makeUxOut(t *testing.T, s cipher.SecKey, coins, hours uint64) coin.UxOut {
	body := makeUxBody(t, s, coins, hours)
	tm := rand.Int31n(1000)
	seq := rand.Int31n(100)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  uint64(tm),
			BkSeq: uint64(seq),
		},
		Body: body,
	}
}

func makeUxBody(t *testing.T, s cipher.SecKey, coins, hours uint64) coin.UxBody {
	p := cipher.PubKeyFromSecKey(s)
	return coin.UxBody{
		SrcTransaction: cipher.SumSHA256(testutil.RandBytes(t, 128)),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          coins,
		Hours:          hours,
	}
}
