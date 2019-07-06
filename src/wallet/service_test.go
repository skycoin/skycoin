package wallet

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/testutil"
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

			require.Equal(t, dir, s.config.WalletDir)

			require.Equal(t, 0, len(s.wallets))

			// test load wallets
			s, err = NewService(Config{
				WalletDir:       "./testdata",
				CryptoType:      ct,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			require.Equal(t, 7, len(s.wallets))

		})
	}
}

func TestNewServiceDupWallets(t *testing.T) {
	_, err := NewService(Config{
		WalletDir:       "./testdata/duplicate_wallets",
		EnableWalletAPI: true,
	})
	require.NotNil(t, err)
	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), "duplicate wallet found with fingerprint deterministic-2M755W9o7933roLASK9PZTmqRsjQUsVen9y in file"), err.Error())
}

func TestNewServiceEmptyWallet(t *testing.T) {
	_, err := NewService(Config{
		WalletDir:       "./testdata/empty_wallet",
		EnableWalletAPI: true,
	})
	testutil.RequireError(t, err, "empty wallet file found: \"empty.wlt\"")
}

func TestServiceCreateWallet(t *testing.T) {
	tt := []struct {
		name            string
		encrypt         bool
		password        []byte
		enableWalletAPI bool
		walletType      string
		filename        string
		seed            string
		err             error
	}{
		{
			name:            "type=collection encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: true,
			walletType:      WalletTypeCollection,
			filename:        "t1.wlt",
		},
		{
			name:            "encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: true,
			filename:        "t1.wlt",
			seed:            "seed1",
		},
		{
			name:            "encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: false,
			filename:        "t1.wlt",
			seed:            "seed1",
			err:             ErrWalletAPIDisabled,
		},
		{
			name:            "encrypt=false",
			encrypt:         false,
			enableWalletAPI: true,
			filename:        "t1.wlt",
			seed:            "seed1",
		},
		{
			name:            "encrypt=false",
			encrypt:         false,
			enableWalletAPI: false,
			filename:        "t1.wlt",
			seed:            "seed1",
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

				w, err := s.CreateWallet(tc.filename, Options{
					Seed:     tc.seed,
					Encrypt:  tc.encrypt,
					Password: tc.password,
					Type:     tc.walletType,
				}, nil)

				if tc.err == nil {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					require.Equal(t, tc.err, err, "%s != %s", tc.err.Error(), err.Error())
					return
				}

				require.NoError(t, err)
				require.Equal(t, w.IsEncrypted(), tc.encrypt)
				if tc.encrypt {
					require.NotEmpty(t, w.Secrets())
					checkNoSensitiveData(t, w)

					// Checks the wallet file doesn't contain sensitive data
					lw, err := Load(filepath.Join(dir, w.Filename()))
					require.NoError(t, err)
					checkNoSensitiveData(t, lw)
				} else {
					require.NoError(t, w.Validate())
				}

				// create wallet with dup wallet name
				_, err = s.CreateWallet(tc.filename, Options{
					Seed: tc.seed + "2",
				}, nil)
				require.Equal(t, err, ErrWalletNameConflict)

				if tc.walletType != WalletTypeCollection {
					// create wallet with dup seed
					dupWlt := "dup_wallet.wlt"
					_, err = s.CreateWallet(dupWlt, Options{
						Seed: tc.seed,
					}, nil)
					require.Equal(t, err, ErrSeedUsed)

					// check if the dup wallet is created
					_, ok := s.wallets[dupWlt]
					require.False(t, ok)

					testutil.RequireFileNotExists(t, filepath.Join(dir, dupWlt))
				}
			})
		}
	}
}

func TestServiceLoadWallet(t *testing.T) {
	// Prepare addresss
	seed := "seed"
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 10)
	var addrs []cipher.Address
	for _, s := range seckeys {
		addrs = append(addrs, cipher.MustAddressFromSecKey(s))
	}

	tt := []struct {
		name          string
		opts          Options
		tf            TransactionsFinder
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
			mockTxnsFinder{
				addrs[0]: true,
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
			mockTxnsFinder{
				addrs[1]: true,
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
			mockTxnsFinder{
				addrs[0]: true,
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
			mockTxnsFinder{
				addrs[1]: true,
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
				wltName := NewWalletFilename()

				w, err := s.loadWallet(wltName, tc.opts, tc.tf)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.Equal(t, w.EntriesLen(), tc.expectAddrNum)
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, w.GetEntryAt(i).Address)
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
	seed := "seed"
	// Generate adddresses from the seed
	var addrs []cipher.Address
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 10)
	for _, s := range seckeys {
		addrs = append(addrs, cipher.MustAddressFromSecKey(s))
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
				Seed:  seed,
			},
			n:             0,
			expectAddrNum: 0,
		},
		{
			name: "encrypted=false addresses=1",
			opts: Options{
				Label: "label",
				Seed:  seed,
			},
			n:             2,
			expectAddrNum: 2,
			expectAddrs:   addrs[1:3], // CreateWallet will generate a default address, so check from new address
		},
		{
			name: "encrypted=false addresses=2",
			opts: Options{
				Label: "label",
				Seed:  seed,
			},
			n:             2,
			expectAddrNum: 2,
			expectAddrs:   addrs[1:3], // CreateWallet will generate a default address, so check from new address
		},
		{
			name: "encrypted=true addresses=1",
			opts: Options{
				Label:    "label",
				Seed:     seed,
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
				Seed:     seed,
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
				Seed:     seed,
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
		{
			name: "encrypted=false password provided",
			opts: Options{
				Label: "label",
				Seed:  seed,
			},
			n:         1,
			pwd:       []byte("foo"),
			expectErr: ErrWalletNotEncrypted,
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

				wltName := NewWalletFilename()

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
				require.Equal(t, w.EntriesLen(), int(tc.n+1))

				// Wallet has a default address, so need to start from the second address
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, w.GetEntryAt(i+1).Address)
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
					WalletDir:       "./testdata",
					CryptoType:      ct,
					EnableWalletAPI: enableWalletAPI,
				})
				require.NoError(t, err)

				if !enableWalletAPI {
					dirIsEmpty(t, dir)

					require.Empty(t, s.wallets)
					addrs, err := s.GetSkycoinAddresses("")
					require.Equal(t, ErrWalletAPIDisabled, err)
					require.Equal(t, 0, len(addrs))
					return
				}

				addrs, err := s.GetSkycoinAddresses("test1.wlt")
				require.NoError(t, err)
				require.Equal(t, 1, len(addrs))

				// test none exist wallet
				notExistID := "not_exist_id.wlt"
				_, err = s.GetSkycoinAddresses(notExistID)
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
					require.Nil(t, w)
					return
				}

				// Create a wallet
				w, err := s.CreateWallet("t.wlt", Options{
					Label: "label",
					Seed:  "seed",
				}, nil)
				require.NoError(t, err)

				w1, err := s.GetWallet(w.Filename())
				require.NoError(t, err)

				// Check if change original wallet would change the returned wallet
				w.SetLabel("new_label")

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

				// Creates a wallet
				w, err := s.CreateWallet("t.wlt", Options{
					Label: "label",
					Seed:  "seed",
				}, nil)
				require.NoError(t, err)

				var wallets []Wallet
				// Get the default wallet
				wallets = append(wallets, w)

				// Create a new wallet
				wltName := NewWalletFilename()
				w1, err := s.CreateWallet(wltName, Options{
					Label: "label1",
					Seed:  "seed1",
				}, nil)
				require.NoError(t, err)
				wallets = append(wallets, w1)

				ws, err := s.GetWallets()
				require.NoError(t, err)
				for _, w := range wallets {
					ww, ok := ws[w.Filename()]
					require.True(t, ok)
					require.Equal(t, w, ww)
				}
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
			name:    "ok collection wallet",
			wltName: "t.wlt",
			opts: Options{
				Type: WalletTypeCollection,
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

				// Add an entry to a collection wallet, to verify that secrets are hidden
				if w.Type() == WalletTypeCollection {
					err := s.Update(w.Filename(), func(w Wallet) error {
						p, s := cipher.GenerateKeyPair()
						return w.(*CollectionWallet).AddEntry(Entry{
							Public:  p,
							Secret:  s,
							Address: cipher.AddressFromPubKey(p),
						})
					})
					require.NoError(t, err)
				}

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
				require.Equal(t, cipher.SecKey{}, encWlt.GetEntryAt(0).Secret)
				require.Empty(t, encWlt.Seed())
				require.Empty(t, encWlt.LastSeed())

				// Check the decrypted seeds
				decWlt, err := Unlock(encWlt, tc.pwd)
				require.NoError(t, err)
				require.Equal(t, w.Seed(), decWlt.Seed())
				require.Equal(t, w.LastSeed(), decWlt.LastSeed())

				// Check if the wallet file does exist
				path := filepath.Join(dir, w.Filename())
				testutil.RequireFileExists(t, path)

				// Check that the temporary backup wallet file does not exist
				bakPath := path + ".bak"
				testutil.RequireFileNotExists(t, bakPath)
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
			name:    "ok collection",
			wltName: "test.wlt",
			opts: Options{
				Type:     WalletTypeCollection,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
		},
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

				verifyDecryptedWlt := func(wlt Wallet) {
					// Checks the "encrypted" meta info
					require.False(t, wlt.IsEncrypted())
					// Checks the seed
					require.Equal(t, tc.opts.Seed, wlt.Seed())
					// Checks the last seed
					entryNum := wlt.EntriesLen()
					lsd, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(wlt.Seed()), entryNum)
					require.NoError(t, err)
					require.Equal(t, hex.EncodeToString(lsd), wlt.LastSeed())

					// Checks the entries
					for i := range seckeys {
						a := cipher.MustAddressFromSecKey(seckeys[i])
						require.Equal(t, a, wlt.GetEntryAt(i).Address)
						require.Equal(t, seckeys[i], wlt.GetEntryAt(i).Secret)
					}

					require.Empty(t, wlt.Secrets())
					require.Empty(t, wlt.CryptoType())
				}

				// Checks the decrypted wallet in service
				w, err := s.getWallet(tc.wltName)
				require.NoError(t, err)
				verifyDecryptedWlt(w)

				// Checks the existence of the wallet file
				fn := filepath.Join(dir, tc.wltName)
				testutil.RequireFileExists(t, fn)

				// Loads wallet from the file and check if it's decrypted
				w1, err := Load(fn)
				require.NoError(t, err)
				verifyDecryptedWlt(w1)
			})
		}
	}
}

func TestServiceCreateWalletWithScan(t *testing.T) {
	seed := "seed1"
	addrs := make([]cipher.Address, 20)
	childSeeds := make([]string, 20)
	lastSeed := []byte(seed)
	for i := range addrs {
		s, pk, _, err := cipher.DeterministicKeyPairIterator(lastSeed)
		require.NoError(t, err)
		addrs[i] = cipher.AddressFromPubKey(pk)
		childSeeds[i] = hex.EncodeToString(s)
		lastSeed = s
	}

	tf := make(mockTxnsFinder, 20)

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
		balGetter        TransactionsFinder
		disableWalletAPI bool
		expect           exp
	}{
		{
			name: "no coins and scan 0, unencrypted",
			opts: Options{
				Seed: "seed1",
			},
			balGetter: tf,
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 0, encrypted",
			opts: Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			balGetter: tf,
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 1, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 1,
			},
			balGetter: tf,
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 1, encrypted",
			opts: Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    1,
			},
			balGetter: tf,
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "no coins and scan 10, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 10,
			},
			balGetter: tf,
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 5, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[5]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[5],
				entryNum:         5 + 1,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 8, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[5]: true,
				addrs[8]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[7],
				entryNum:         8 + 1,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 10, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[4+1]: true,
				addrs[10]:  true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[9],
				entryNum:         10 + 1,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},

		{
			name: "scan 5 get 5, encrypted",
			opts: Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter: mockTxnsFinder{
				addrs[5]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[5],
				entryNum:         5 + 1,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 4, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[3]: true,
				addrs[4]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[4],
				entryNum:         4 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 4, encrypted",
			opts: Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter: mockTxnsFinder{
				addrs[3]: true,
				addrs[4]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[4],
				entryNum:         4 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 4 have 6, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[3]: true,
				addrs[4]: true,
				addrs[6]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[6],
				entryNum:         6 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 2 have 7, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[2]: true,
				addrs[7]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[7],
				entryNum:         7 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 2 get 7 have 12, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[2]:  true,
				addrs[7]:  true,
				addrs[12]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[12],
				entryNum:         12 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 2 get 7 have 13, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[2]:  true,
				addrs[7]:  true,
				addrs[13]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[7],
				entryNum:         7 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "scan 5 get 2 have 8, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[2]: true,
				addrs[8]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[2],
				entryNum:         2 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "confirmed and predicted, unencrypted",
			opts: Options{
				Seed:  seed,
				ScanN: 5,
			},
			balGetter: mockTxnsFinder{
				addrs[3]: true,
				addrs[4]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[4],
				entryNum:         4 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "confirmed and predicted, encrypted",
			opts: Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter: mockTxnsFinder{
				addrs[3]: true,
				addrs[4]: true,
			},
			expect: exp{
				err:              nil,
				seed:             seed,
				lastSeed:         childSeeds[4],
				entryNum:         4 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			name: "wallet api disabled",
			opts: Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
			},
			balGetter:        mockTxnsFinder{},
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

				wltName := NewWalletFilename()
				w, err := s.CreateWallet(wltName, tc.opts, tc.balGetter)
				require.Equal(t, tc.expect.err, err)
				if err != nil {
					return
				}

				require.NoError(t, w.Validate())
				require.Equal(t, tc.expect.entryNum, w.EntriesLen())
				for i, e := range w.GetEntries() {
					require.Equal(t, addrs[i].String(), e.Address.String())
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

func TestServiceView(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             Options
		viewWltName      string
		action           func(*testing.T) func(Wallet) error
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted collection wallet",
			wltName:     "test-view-collection-encrypted.wlt",
			viewWltName: "test-view-collection-encrypted.wlt",
			opts: Options{
				Type:     WalletTypeCollection,
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					fmt.Println("checking label")
					require.Equal(t, "foowlt", w.Label())
					fmt.Println("checking sensitive")
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that this references a clone and not the original
					fmt.Println("modifying label")
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, unencrypted collection wallet",
			wltName:     "test-view-collection-unencrypted.wlt",
			viewWltName: "test-view-collection-unencrypted.wlt",
			opts: Options{
				Label: "foowlt",
				Type:  WalletTypeCollection,
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					// Collection wallets don't have seeds
					require.Empty(t, w.Seed())
					require.Empty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, encrypted wallet",
			wltName:     "test-view-encrypted.wlt",
			viewWltName: "test-view-encrypted.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, unencrypted wallet",
			wltName:     "test-view-unencrypted.wlt",
			viewWltName: "test-view-unencrypted.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.Seed())
					require.NotEmpty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "wallet doesn't exist",
			wltName:     "test-view-not-exist.wlt",
			viewWltName: "foo-test-view-not-exist.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			err: ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-view-api-disabled.wlt",
			viewWltName: "test-view-api-disabled.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:       dir,
				CryptoType:      CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			w, err := s.CreateWallet(tc.wltName, tc.opts, nil)
			require.NoError(t, err)

			if w.Type() == WalletTypeCollection {
				err := s.UpdateSecrets(w.Filename(), tc.opts.Password, func(w Wallet) error {
					p, s := cipher.GenerateKeyPair()
					return w.(*CollectionWallet).AddEntry(Entry{
						Public:  p,
						Secret:  s,
						Address: cipher.AddressFromPubKey(p),
					})
				})
				require.NoError(t, err)

				w, err = s.GetWallet(tc.wltName)
				require.NoError(t, err)
			}

			s.config.EnableWalletAPI = !tc.disableWalletAPI

			var action func(Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.View(tc.viewWltName, action)

			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.config.EnableWalletAPI = true

			// Check that the wallet is unmodified
			w2, err := s.GetWallet(tc.wltName)
			require.NoError(t, err)
			require.Equal(t, w, w2)
		})
	}
}

func TestServiceViewSecrets(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             Options
		viewWltName      string
		action           func(*testing.T) func(Wallet) error
		password         []byte
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted wallet",
			wltName:     "test-view-secrets-encrypted.wlt",
			viewWltName: "test-view-secrets-encrypted.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			password: []byte("pwd"),
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should be able to see sensitive data
					require.Equal(t, "fooseed", w.Seed())
					require.NotEmpty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, unencrypted wallet",
			wltName:     "test-view-secrets-unencrypted.wlt",
			viewWltName: "test-view-secrets-unencrypted.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.Seed())
					require.NotEmpty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "encrypted wallet but password not provided",
			wltName:     "test-view-secrets-encrypted-no-password.wlt",
			viewWltName: "test-view-secrets-encrypted-no-password.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			err: ErrMissingPassword,
		},

		{
			name:        "encrypted wallet but password invalid",
			wltName:     "test-view-secrets-encrypted-wrong-password.wlt",
			viewWltName: "test-view-secrets-encrypted-wrong-password.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			password: []byte("pwdpwd"),
			err:      ErrInvalidPassword,
		},

		{
			name:        "unencrypted wallet but password provided",
			wltName:     "test-view-secrets-unencrypted-password.wlt",
			viewWltName: "test-view-secrets-unencrypted-password.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			password: []byte("pwd"),
			err:      ErrWalletNotEncrypted,
		},

		{
			name:        "wallet doesn't exist",
			wltName:     "test-view-secrets-not-exist.wlt",
			viewWltName: "foo-test-view-secrets-not-exist.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			err: ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-view-secrets-api-disabled.wlt",
			viewWltName: "test-view-secrets-api-disabled.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:       dir,
				CryptoType:      CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			w, err := s.CreateWallet(tc.wltName, tc.opts, nil)
			require.NoError(t, err)

			s.config.EnableWalletAPI = !tc.disableWalletAPI

			var action func(Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.ViewSecrets(tc.viewWltName, tc.password, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.config.EnableWalletAPI = true

			// Check that the wallet is unmodified
			w2, err := s.GetWallet(tc.wltName)
			require.NoError(t, err)
			require.Equal(t, w, w2)
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             Options
		viewWltName      string
		action           func(*testing.T) func(Wallet) error
		checkWallet      func(*testing.T, Wallet)
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted wallet",
			wltName:     "test-update-encrypted.wlt",
			viewWltName: "test-update-encrypted.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should not be able to see sensitive data
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					// The wallet is encrypted so it cannot generate more addresses
					_, err := w.GenerateAddresses(1)
					require.Equal(t, ErrWalletEncrypted, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Equal(t, 1, w.EntriesLen())
				checkNoSensitiveData(t, w)
			},
		},

		{
			name:        "ok, unencrypted wallet",
			wltName:     "test-update-unencrypted.wlt",
			viewWltName: "test-update-unencrypted.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.Seed())
					require.NotEmpty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
			checkWallet: func(t *testing.T, w Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Equal(t, 1, w.EntriesLen())
				require.NotEmpty(t, w.GetEntryAt(0).Secret)
			},
		},

		{
			name:        "wallet doesn't exist",
			wltName:     "test-update-not-exist.wlt",
			viewWltName: "foo-test-update-not-exist.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			err: ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-update-api-disabled.wlt",
			viewWltName: "test-update-api-disabled.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:       dir,
				CryptoType:      CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			_, err = s.CreateWallet(tc.wltName, tc.opts, nil)
			require.NoError(t, err)

			s.config.EnableWalletAPI = !tc.disableWalletAPI

			var action func(Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.Update(tc.viewWltName, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.config.EnableWalletAPI = true

			// Check that the wallet was modified as expected
			w, err := s.GetWallet(tc.wltName)
			require.NoError(t, err)
			tc.checkWallet(t, w)

			// Even if secrets were modified, wallet should still be encrypted
			require.Equal(t, tc.opts.Encrypt, w.IsEncrypted())
			if w.IsEncrypted() {
				checkNoSensitiveData(t, w)
			}
		})
	}
}

func TestServiceUpdateSecrets(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             Options
		viewWltName      string
		action           func(*testing.T) func(Wallet) error
		checkWallet      func(*testing.T, Wallet)
		password         []byte
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted wallet",
			wltName:     "test-update-secrets-encrypted.wlt",
			viewWltName: "test-update-secrets-encrypted.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			password: []byte("pwd"),
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should be able to see sensitive data
					require.Equal(t, "fooseed", w.Seed())
					require.NotEmpty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Equal(t, 2, w.EntriesLen())
				checkNoSensitiveData(t, w)
			},
		},

		{
			name:        "ok, unencrypted wallet",
			wltName:     "test-update-secrets-unencrypted.wlt",
			viewWltName: "test-update-secrets-unencrypted.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			action: func(t *testing.T) func(Wallet) error {
				return func(w Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.Seed())
					require.NotEmpty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Equal(t, 2, w.EntriesLen())
				require.NotEmpty(t, w.GetEntryAt(1).Secret)
			},
		},

		{
			name:        "encrypted wallet but password not provided",
			wltName:     "test-update-secrets-encrypted-no-password.wlt",
			viewWltName: "test-update-secrets-encrypted-no-password.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			err: ErrMissingPassword,
		},

		{
			name:        "encrypted wallet but password invalid",
			wltName:     "test-update-secrets-encrypted-wrong-password.wlt",
			viewWltName: "test-update-secrets-encrypted-wrong-password.wlt",
			opts: Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			password: []byte("pwdpwd"),
			err:      ErrInvalidPassword,
		},

		{
			name:        "unencrypted wallet but password provided",
			wltName:     "test-update-secrets-unencrypted-password.wlt",
			viewWltName: "test-update-secrets-unencrypted-password.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			password: []byte("pwd"),
			err:      ErrWalletNotEncrypted,
		},

		{
			name:        "wallet doesn't exist",
			wltName:     "test-update-secrets-not-exist.wlt",
			viewWltName: "foo-test-update-secrets-not-exist.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			err: ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-update-secrets-api-disabled.wlt",
			viewWltName: "test-update-secrets-api-disabled.wlt",
			opts: Options{
				Seed:  "fooseed",
				Label: "foowlt",
			},
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := NewService(Config{
				WalletDir:       dir,
				CryptoType:      CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			_, err = s.CreateWallet(tc.wltName, tc.opts, nil)
			require.NoError(t, err)

			s.config.EnableWalletAPI = !tc.disableWalletAPI

			var action func(Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.UpdateSecrets(tc.viewWltName, tc.password, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.config.EnableWalletAPI = true

			// Check that the wallet was modified as expected
			w, err := s.GetWallet(tc.wltName)
			require.NoError(t, err)
			tc.checkWallet(t, w)

			// Even if secrets were modified, wallet should still be encrypted
			require.Equal(t, tc.opts.Encrypt, w.IsEncrypted())
			if w.IsEncrypted() {
				checkNoSensitiveData(t, w)
			}
		})
	}
}

func checkNoSensitiveData(t *testing.T, w Wallet) {
	fmt.Println("check seed empty")
	require.Empty(t, w.Seed())
	fmt.Println("check lastseed empty")
	require.Empty(t, w.LastSeed())
	fmt.Println("check secret entries empty")
	for _, e := range w.GetEntries() {
		require.True(t, e.Secret.Null())
	}
}
