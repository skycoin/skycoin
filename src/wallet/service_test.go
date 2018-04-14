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

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
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
				WalletDir:       dir,
				CryptoType:      ct,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			// check if the wallet dir is created
			_, err = os.Stat(dir)
			require.NoError(t, err)

			require.Equal(t, dir, s.walletDirectory)

			require.Equal(t, 1, len(s.wallets))

			// test load wallets
			s, err = NewService(Config{
				WalletDir:       "./testdata",
				CryptoType:      ct,
				EnableWalletAPI: true,
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

func TestServiceCreateWallet(t *testing.T) {
	tt := []struct {
		name            string
		encrypt         bool
		password        []byte
		enableWalletAPI bool
		err             error
	}{
		{
			name:            "encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: true,
		},
		{
			name:            "encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: false,
			err:             ErrWalletAPIDisabled,
		},
		{
			name:            "encrypt=false",
			encrypt:         false,
			enableWalletAPI: true,
		},
		{
			name:            "encrypt=false",
			encrypt:         false,
			enableWalletAPI: false,
			err:             ErrWalletAPIDisabled,
		},
	}
	for _, tc := range tt {
		for ct := range cryptoTable {
			t.Run(fmt.Sprintf("%v crypto=%v", tc.name, ct), func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: tc.enableWalletAPI,
				})
				require.NoError(t, err)

				wltName := "t1.wlt"
				seed := "seed1"
				w, err := s.CreateWallet(wltName, Options{
					Seed:     seed,
					Encrypt:  tc.encrypt,
					Password: tc.password,
				}, nil)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.NoError(t, err)
				require.Equal(t, w.IsEncrypted(), tc.encrypt)
				if tc.encrypt {
					require.NotEmpty(t, w.secrets())
					checkNoSensitiveData(t, w)

					// Checks the wallet file doesn't contain sensitive data
					lw, err := Load(filepath.Join(dir, w.Filename()))
					require.NoError(t, err)
					checkNoSensitiveData(t, lw)
				} else {
					require.NoError(t, w.Validate())
				}

				// create wallet with dup wallet name
				_, err = s.CreateWallet(wltName, Options{Seed: "seed2"}, nil)
				require.Equal(t, err, ErrWalletNameConflict)

				// create wallet with dup seed
				dupWlt := "dup_wallet.wlt"
				_, err = s.CreateWallet(dupWlt, Options{
					Seed: seed,
				}, nil)
				require.EqualError(t, err, fmt.Sprintf("wallet %s would be duplicate with %v, same seed", dupWlt, wltName))

				// check if the dup wallet is created
				_, ok := s.wallets[dupWlt]
				require.False(t, ok)

				_, err = os.Stat(filepath.Join(dir, dupWlt))
				require.True(t, os.IsNotExist(err))
			})
		}
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
				ScanN: 5,
			},
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
				ScanN: 5,
			},
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
				ScanN:    5,
			},
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
				ScanN:    5,
			},
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
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: true,
				})
				require.NoError(t, err)
				wltName := newWalletFilename()

				w, err := s.loadWallet(wltName, tc.opts, tc.bg)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.Len(t, w.Entries, tc.expectAddrNum)
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, w.Entries[i].Address)
				}

				require.Equal(t, w.IsEncrypted(), tc.opts.Encrypt)
				if w.IsEncrypted() {
					checkNoSensitiveData(t, w)
					// Checks the wallet file doesn't contain sensitive data
					wltPath := filepath.Join(dir, w.Filename())
					lw, err := Load(wltPath)
					require.NoError(t, err)
					checkNoSensitiveData(t, lw)
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
		name              string
		opts              Options
		n                 uint64
		pwd               []byte
		walletAPIDisabled bool
		expectAddrNum     int
		expectAddrs       []cipher.Address
		expectErr         error
	}{
		{
			name: "encrypted=false addresses=0",
			opts: Options{
				Label: "label",
				Seed:  string(seed),
			},
			n:             0,
			expectAddrNum: 0,
		},
		{
			name: "encrypted=false addresses=1",
			opts: Options{
				Label: "label",
				Seed:  string(seed),
			},
			n:             2,
			expectAddrNum: 2,
			expectAddrs:   addrs[1:3], // CreateWallet will generate a default address, so check from new address
		},
		{
			name: "encrypted=false addresses=2",
			opts: Options{
				Label: "label",
				Seed:  string(seed),
			},
			n:             2,
			expectAddrNum: 2,
			expectAddrs:   addrs[1:3], // CreateWallet will generate a default address, so check from new address
		},
		{
			name: "encrypted=true addresses=1",
			opts: Options{
				Label:    "label",
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			n:             1,
			pwd:           []byte("pwd"),
			expectAddrNum: 1,
			expectAddrs:   addrs[1:2], // CreateWallet will generate a default address, so check from new address
		},
		{
			name: "encrypted=true addresses=2",
			opts: Options{
				Label:    "label",
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			n:             2,
			pwd:           []byte("pwd"),
			expectAddrNum: 2,
			expectAddrs:   addrs[1:3], // CreateWallet will generate a default address, so check from new address
		},
		{
			name: "encrypted=true wrong password",
			opts: Options{
				Label:    "label",
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			n:             1,
			pwd:           []byte("wrong password"),
			expectAddrNum: 1,
			expectErr:     ErrInvalidPassword,
		},
		{
			name: "wallet api disabled",
			opts: Options{
				Seed:  "seed",
				Label: "label",
			},
			walletAPIDisabled: true,
			expectErr:         ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.walletAPIDisabled,
				})
				require.NoError(t, err)

				wltName := newWalletFilename()

				w, err := s.CreateWallet(wltName, tc.opts, nil)
				if err != nil {
					require.Equal(t, tc.expectErr, err)
					return
				}

				if w.IsEncrypted() {
					checkNoSensitiveData(t, w)
				}

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
				if w.IsEncrypted() {
					checkNoSensitiveData(t, lw)
				}

				// Wallet doesn't exist
				_, err = s.NewAddresses("wallet_not_exist.wlt", tc.pwd, 1)
				require.Equal(t, ErrWalletNotExist, err)
			})
		}
	}
}

func TestServiceGetAddress(t *testing.T) {
	for _, enableWalletAPI := range []bool{true, false} {
		for ct := range cryptoTable {
			t.Run(fmt.Sprintf("enable wallet api=%v crypto=%v", enableWalletAPI, ct), func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: enableWalletAPI,
				})
				require.NoError(t, err)

				if !enableWalletAPI {
					dirIsEmpty(t, dir)

					require.Empty(t, s.wallets)
					addrs, err := s.GetAddresses("")
					require.Equal(t, ErrWalletAPIDisabled, err)
					require.Equal(t, 0, len(addrs))
					return
				}

				// Get the default wallet
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
}

func TestServiceGetWallet(t *testing.T) {
	for _, enableWalletAPI := range []bool{true, false} {
		for ct := range cryptoTable {
			t.Run(fmt.Sprintf("enable wallet api=%v crypto=%v", enableWalletAPI, ct), func(t *testing.T) {
				dir := prepareWltDir()

				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: enableWalletAPI,
				})
				require.NoError(t, err)

				if !enableWalletAPI {
					dirIsEmpty(t, dir)

					require.Empty(t, s.wallets)
					w, err := s.GetWallet("")
					require.Equal(t, ErrWalletAPIDisabled, err)
					var emptyW *Wallet
					require.Equal(t, w, emptyW)
					return
				}

				// Get the default wallet
				var w *Wallet
				for _, w = range s.wallets {
				}

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
}

func TestServiceGetWallets(t *testing.T) {
	for _, enableWalletAPI := range []bool{true, false} {
		for ct := range cryptoTable {
			t.Run(fmt.Sprintf("enable wallet=%v crypto=%v", enableWalletAPI, ct), func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: enableWalletAPI,
				})
				require.NoError(t, err)

				if !enableWalletAPI {
					dirIsEmpty(t, dir)

					require.Empty(t, s.wallets)
					w, err := s.GetWallets()
					require.Equal(t, ErrWalletAPIDisabled, err)
					var emptyW Wallets
					require.Equal(t, w, emptyW)
					return
				}

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
				}, nil)
				require.NoError(t, err)
				wallets = append(wallets, w2)

				ws, err := s.GetWallets()
				for _, w := range wallets {
					ww, ok := ws[w.Filename()]
					require.True(t, ok)
					require.Equal(t, w, ww)
				}
			})
		}
	}
}

func TestServiceReloadWalletsDisabledWalletAPI(t *testing.T) {
	for ct := range cryptoTable {
		name := fmt.Sprintf("crypto=%v", ct)
		t.Run(name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:       dir,
				CryptoType:      ct,
				EnableWalletAPI: false,
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
				WalletDir:       dir,
				CryptoType:      ct,
				EnableWalletAPI: true,
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
			w1, err := s.CreateWallet(wltName, Options{Seed: "seed1"}, nil)
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
		name             string
		opts             Options
		pwd              []byte
		unspents         []coin.UxOut
		vld              Validator
		coins            uint64
		dest             cipher.Address
		disableWalletAPI bool
		err              error
	}{
		{
			name: "encrypted=false has change=no",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok: false,
			},
			coins: 2e6,
			dest:  addrs[0],
		},
		{
			name: "encrypted=true has change=no",
			opts: Options{
				Seed:     string(seed),
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			pwd:      []byte("pwd"),
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok: false,
			},
			coins: 2e6,
			dest:  addrs[0],
		},
		{
			name: "encrypted=false has change=yes",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok: false,
			},
			coins: 1e6,
			dest:  addrs[0],
		},
		{
			name: "encrypted=false has unconfirmed spending transaction",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok: true,
			},
			coins: 2e6,
			dest:  addrs[0],
			err:   errors.New("please spend after your pending transaction is confirmed"),
		},
		{
			name: "encrypted=false unconfirmed spend failed",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok:  false,
				err: errors.New("fail intentionally"),
			},
			coins: 2e6,
			dest:  addrs[0],
			err:   errors.New("checking unconfirmed spending failed: fail intentionally"),
		},
		{
			name: "encrypted=false spend zero",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok: false,
			},
			dest: addrs[0],
			err:  errors.New("zero spend amount"),
		},
		{
			name: "encrypted=false spend fractional coins",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok: false,
			},
			coins: 1e3,
			dest:  addrs[0],
		},
		{
			name: "encrypted=false not enough confirmed coins",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			vld: &dummyValidator{
				ok: false,
			},
			coins: 100e6,
			dest:  addrs[0],
			err:   ErrInsufficientBalance,
		},
		{
			name: "encrypted=false no coin hours in inputs",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxoutsNoHours[:],
			vld: &dummyValidator{
				ok: false,
			},
			coins: 1e6,
			dest:  addrsNoHours[0],
			err:   fee.ErrTxnNoFee,
		},
		{
			name: "disable wallet api=true",
			opts: Options{
				Seed:  string(seed),
				Label: "label",
			},
			vld:              &dummyValidator{},
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
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
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
				})
				require.NoError(t, err)

				if tc.disableWalletAPI {
					_, err = s.CreateAndSignTransaction("", tc.pwd, tc.vld, unspents, uint64(headTime), tc.coins, tc.dest)
					require.Equal(t, tc.err, err)
					return
				}

				wltName := newWalletFilename()

				w, err := s.CreateWallet(wltName, tc.opts, nil)
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
		name             string
		wltName          string
		opts             Options
		updateWltName    string
		label            string
		disableWalletAPI bool
		err              error
	}{
		{
			name:    "ok",
			wltName: "t.wlt",
			opts: Options{
				Seed:  "seed",
				Label: "label",
			},
			updateWltName: "t.wlt",
			label:         "new-label",
		},
		{
			name:    "wallet doesn't exist",
			wltName: "t.wlt",
			opts: Options{
				Seed:  "seed",
				Label: "label",
			},
			updateWltName: "t1.wlt",
			label:         "new-label",
			err:           ErrWalletNotExist,
		},
		{
			name:    "wallet api disabled",
			wltName: "t.wlt",
			opts: Options{
				Seed:  "seed",
				Label: "label",
			},
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			t.Run(tc.name, func(t *testing.T) {
				// Create the wallet service
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
				})
				require.NoError(t, err)

				if tc.disableWalletAPI {
					err = s.UpdateWalletLabel("", "new label")
					require.Equal(t, tc.err, err)
					return
				}

				// Create a new wallet
				w, err := s.CreateWallet(tc.wltName, tc.opts, nil)
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
		name             string
		wltName          string
		opts             Options
		encWltName       string
		pwd              []byte
		disableWalletAPI bool
		err              error
	}{
		{
			name:    "ok",
			wltName: "t.wlt",
			opts: Options{
				Seed: "seed",
			},
			encWltName: "t.wlt",
			pwd:        []byte("pwd"),
		},
		{
			name:    "wallet doesn't exist",
			wltName: "t.wlt",
			opts: Options{
				Seed: "seed",
			},
			encWltName: "t2.wlt",
			err:        ErrWalletNotExist,
		},
		{
			name:    "wallet already encrypted",
			wltName: "t.wlt",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			encWltName: "t.wlt",
			pwd:        []byte("pwd"),
			err:        ErrWalletEncrypted,
		},
		{
			name:    "wallet api disabled",
			wltName: "t.wlt",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			encWltName:       "t.wlt",
			pwd:              []byte("pwd"),
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				// Create the wallet service
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
				})
				require.NoError(t, err)

				if tc.disableWalletAPI {
					_, err = s.EncryptWallet("", tc.pwd)
					require.Equal(t, tc.err, err)
					return
				}

				// Create a new wallet
				w, err := s.CreateWallet(tc.wltName, tc.opts, nil)
				require.NoError(t, err)

				// Encrypt the wallet
				encWlt, err := s.EncryptWallet(tc.encWltName, tc.pwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				encWlt1, err := s.getWallet(tc.encWltName)
				require.NoError(t, err)
				require.Equal(t, encWlt, encWlt1)

				// Check the encrypted wallet
				require.True(t, encWlt.IsEncrypted())
				require.Equal(t, cipher.SecKey{}, encWlt.Entries[0].Secret)
				require.Empty(t, encWlt.seed())
				require.Empty(t, encWlt.lastSeed())

				// Check the decrypted seeds
				decWlt, err := encWlt.unlock(tc.pwd)
				require.NoError(t, err)
				require.Equal(t, w.seed(), decWlt.seed())
				require.Equal(t, w.lastSeed(), decWlt.lastSeed())

				// Check if the wallet file does exist
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
		name             string
		wltName          string
		opts             Options
		decryptWltName   string
		password         []byte
		disableWalletAPI bool
		err              error
	}{
		{
			name:    "ok",
			wltName: "test.wlt",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
		},
		{
			name:    "wallet not exist",
			wltName: "test.wlt",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			decryptWltName: "t.wlt",
			password:       []byte("pwd"),
			err:            ErrWalletNotExist,
		},
		{
			name:    "wallet not encrypted",
			wltName: "test.wlt",
			opts: Options{
				Seed: "seed",
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
			err:            ErrWalletNotEncrypted,
		},
		{
			name:    "invalid password",
			wltName: "test.wlt",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			decryptWltName: "test.wlt",
			password:       []byte("wrong password"),
			err:            ErrInvalidPassword,
		},
		{
			name:    "wallet api disabled",
			wltName: "test.wlt",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			decryptWltName:   "test.wlt",
			password:         []byte("pwd"),
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
				})
				require.NoError(t, err)

				if tc.disableWalletAPI {
					_, err = s.DecryptWallet(tc.decryptWltName, tc.password)
					require.Equal(t, tc.err, err)
					return
				}

				_, err = s.CreateWallet(tc.wltName, tc.opts, nil)
				require.NoError(t, err)

				_, err = s.DecryptWallet(tc.decryptWltName, tc.password)
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

					require.Empty(t, wlt.secrets())
					require.Empty(t, wlt.cryptoType())
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

func TestServiceCreateWalletWithScan(t *testing.T) {
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
		name             string
		opts             Options
		balGetter        BalanceGetter
		disableWalletAPI bool
		expect           exp
	}{
		{
			name: "no coins and scan 0, unencrypted",
			opts: Options{
				Seed: "seed1",
			},
			balGetter: bg,
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 0, encrypted",
			opts: Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			balGetter: bg,
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 1, unencrypted",
			opts: Options{
				Seed:  "seed1",
				ScanN: 1,
			},
			balGetter: bg,
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 1, encrypted",
			opts: Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    1,
			},
			balGetter: bg,
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 10, unencrypted",
			opts: Options{
				Seed:  "seed1",
				ScanN: 10,
			},
			balGetter: bg,
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 5, unencrypted",
			opts: Options{
				Seed:  "seed1",
				ScanN: 5,
			},
			balGetter: mockBalanceGetter{
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[5],
				entryNum:         5,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 5, encrypted",
			opts: Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter: mockBalanceGetter{
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[5],
				entryNum:         5,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 4, unencrypted",
			opts: Options{
				Seed:  "seed1",
				ScanN: 5,
			},
			balGetter: mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 4, encrypted",
			opts: Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter: mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 4 have 6, unencrypted",
			opts: Options{
				Seed:  "seed1",
				ScanN: 5,
			},
			balGetter: mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[6]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "confirmed and predicted, unencrypted",
			opts: Options{
				Seed:  "seed1",
				ScanN: 5,
			},
			balGetter: mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
			},
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "confirmed and predicted, encrypted",
			opts: Options{
				Seed:     "seed1",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter: mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[3]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
			},
			expect: exp{
				err:              nil,
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "wallet api disabled",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter:        mockBalanceGetter{},
			disableWalletAPI: true,
			expect: exp{
				err: ErrWalletAPIDisabled,
			},
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
				})
				require.NoError(t, err)

				wltName := newWalletFilename()
				w, err := s.CreateWallet(wltName, tc.opts, tc.balGetter)
				require.Equal(t, tc.expect.err, err)
				if err != nil {
					return
				}

				require.NoError(t, w.Validate())
				require.Len(t, w.Entries, tc.expect.entryNum)
				for i := range w.Entries {
					require.Equal(t, addrsOfSeed1[i], w.Entries[i].Address.String())
				}
			})
		}
	}
}

func TestGetWalletSeed(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             Options
		id               string
		pwd              []byte
		disableWalletAPI bool
		enableSeedAPI    bool
		expectErr        error
	}{
		{
			name:    "wallet is not encrypted",
			wltName: "wallet.wlt",
			opts: Options{
				Seed:  "seed",
				Label: "label",
			},
			id:            "wallet.wlt",
			enableSeedAPI: true,
			expectErr:     ErrWalletNotEncrypted,
		},
		{
			name:    "wallet api disabled",
			wltName: "wallet.wlt",
			opts: Options{
				Seed:  "seed",
				Label: "label",
			},
			id:               "wallet.wlt",
			enableSeedAPI:    true,
			disableWalletAPI: true,
			expectErr:        ErrWalletAPIDisabled,
		},
		{
			name:    "ok",
			wltName: "wallet.wlt",
			opts: Options{
				Seed:     "seed",
				Label:    "label",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			enableSeedAPI: true,
			id:            "wallet.wlt",
			pwd:           []byte("pwd"),
		},
		{
			name:    "wallet does not exist",
			wltName: "wallet.wlt",
			opts: Options{
				Seed:     "seed",
				Label:    "label",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			enableSeedAPI: true,
			pwd:           []byte("pwd"),
			id:            "none-exist.wlt",
			expectErr:     ErrWalletNotExist,
		},
		{
			name:    "disable seed api",
			wltName: "wallet.wlt",
			opts: Options{
				Seed:     "seed",
				Label:    "label",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			pwd:           []byte("pwd"),
			id:            "wallet.wlt",
			enableSeedAPI: false,
			expectErr:     ErrSeedAPIDisabled,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			t.Run(tc.name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
					EnableSeedAPI:   tc.enableSeedAPI,
				})
				require.NoError(t, err)

				if tc.disableWalletAPI {
					_, err = s.GetWalletSeed("", tc.pwd)
					require.Equal(t, tc.expectErr, err)
					return
				}

				// Create a wallet
				_, err = s.CreateWallet(tc.wltName, tc.opts, nil)
				require.NoError(t, err)

				seed, err := s.GetWalletSeed(tc.id, tc.pwd)
				require.Equal(t, tc.expectErr, err)
				if err != nil {
					return
				}

				require.Equal(t, tc.opts.Seed, seed)
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

func checkNoSensitiveData(t *testing.T, w *Wallet) {
	require.Empty(t, w.seed())
	require.Empty(t, w.lastSeed())
	var empty cipher.SecKey
	for _, e := range w.Entries {
		require.Equal(t, empty, e.Secret)
	}
}
