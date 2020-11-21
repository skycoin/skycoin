package wallet_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/wallet/bip44wallet"
	"github.com/skycoin/skycoin/src/wallet/collection"
	_ "github.com/skycoin/skycoin/src/wallet/deterministic"
	_ "github.com/skycoin/skycoin/src/wallet/xpubwallet"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/crypto"
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
	for _, ct := range crypto.Types() {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			dir := prepareWltDir()
			s, err := wallet.NewService(wallet.Config{
				WalletDir:       dir,
				CryptoType:      ct,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			// check if the wallet dir is created
			_, err = os.Stat(dir)
			require.NoError(t, err)

			//require.Equal(t, dir, s.wallet.Config.WalletDir)
			wlts, err := s.GetWallets()
			require.NoError(t, err)

			require.Equal(t, 0, len(wlts))

			// test load wallets
			s, err = wallet.NewService(wallet.Config{
				WalletDir:       "./testdata",
				CryptoType:      ct,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			wlts, err = s.GetWallets()
			require.NoError(t, err)
			require.Equal(t, 11, len(wlts))
		})
	}
}

func TestNewServiceDupWallets(t *testing.T) {
	_, err := wallet.NewService(wallet.Config{
		WalletDir:       "./testdata/duplicate_wallets",
		EnableWalletAPI: true,
	})
	require.NotNil(t, err)
	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), "duplicate wallet found with fingerprint deterministic-2M755W9o7933roLASK9PZTmqRsjQUsVen9y in file"), err.Error())
}

func TestNewServiceEmptyWallet(t *testing.T) {
	cases := []struct {
		dir string
		fn  string
	}{
		{
			dir: "./testdata/empty_wallet",
			fn:  "empty.wlt",
		},
		{
			dir: "./testdata/empty_bip44_wallet",
			fn:  "empty.wlt",
		},
	}

	for _, tc := range cases {
		t.Run(filepath.Join(tc.dir, tc.fn), func(t *testing.T) {
			_, err := wallet.NewService(wallet.Config{
				WalletDir:       tc.dir,
				EnableWalletAPI: true,
			})
			testutil.RequireError(t, err, fmt.Sprintf("empty wallet file found: %q", tc.fn))
		})
	}
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
		xpub            string
		err             error
	}{
		{
			name:            "type=xpub encrypt=false",
			encrypt:         false,
			enableWalletAPI: true,
			walletType:      wallet.WalletTypeXPub,
			filename:        "t1.wlt",
			xpub:            "xpub6EFYYRQeAbWLdWQYbtQv8HnemieKNmYUE23RmwphgtMLjz4UaStKADSKNoSSXM5FDcq4gZec2q6n7kdNWfuMdScxK1cXm8tR37kaitHtvuJ",
		},
		{
			name:            "type=xpub encrypt=true password=pwd",
			encrypt:         true,
			enableWalletAPI: true,
			password:        []byte("pwd"),
			walletType:      wallet.WalletTypeXPub,
			filename:        "t1.wlt",
			xpub:            "xpub6EFYYRQeAbWLdWQYbtQv8HnemieKNmYUE23RmwphgtMLjz4UaStKADSKNoSSXM5FDcq4gZec2q6n7kdNWfuMdScxK1cXm8tR37kaitHtvuJ",
			err:             wallet.NewError(fmt.Errorf("xpub wallet does not support encryption")),
		},
		{
			name:            "type=collection encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: true,
			walletType:      wallet.WalletTypeCollection,
			filename:        "t1.wlt",
		},
		{
			name:            "type=bip44 encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: true,
			walletType:      wallet.WalletTypeBip44,
			filename:        "t1.wlt",
			seed:            "voyage say extend find sheriff surge priority merit ignore maple cash argue",
		},
		{
			name:            "encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: true,
			filename:        "t1.wlt",
			seed:            "seed1",
			walletType:      wallet.WalletTypeDeterministic,
		},
		{
			name:            "encrypt=true password=pwd",
			encrypt:         true,
			password:        []byte("pwd"),
			enableWalletAPI: false,
			filename:        "t1.wlt",
			seed:            "seed1",
			err:             wallet.ErrWalletAPIDisabled,
			walletType:      wallet.WalletTypeDeterministic,
		},
		{
			name:            "encrypt=false",
			encrypt:         false,
			enableWalletAPI: true,
			filename:        "t1.wlt",
			seed:            "seed1",
			walletType:      wallet.WalletTypeDeterministic,
		},
		{
			name:            "encrypt=false",
			encrypt:         false,
			enableWalletAPI: false,
			filename:        "t1.wlt",
			seed:            "seed1",
			err:             wallet.ErrWalletAPIDisabled,
			walletType:      wallet.WalletTypeDeterministic,
		},
	}
	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			t.Run(fmt.Sprintf("%v crypto=%v", tc.name, ct), func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: tc.enableWalletAPI,
				})
				require.NoError(t, err)

				w, err := s.CreateWallet(tc.filename, wallet.Options{
					Seed:     tc.seed,
					Encrypt:  tc.encrypt,
					Password: tc.password,
					Type:     tc.walletType,
					XPub:     tc.xpub,
				})

				require.Equal(t, tc.err, err, fmt.Sprintf("expect %v; got: %v", tc.err, err))
				if err != nil {
					return
				}

				require.NoError(t, err)
				require.Equal(t, w.IsEncrypted(), tc.encrypt)
				if tc.encrypt {
					require.NotEmpty(t, w.Secrets())
					checkNoSensitiveData(t, w)

					// checks that the wallet file has been saved to disk
					_, err = os.Stat(filepath.Join(dir, tc.filename))
					require.False(t, os.IsNotExist(err))

					// Confirms that the data saved to the disk is the same as the wallet.Deserialize()
					data, err := ioutil.ReadFile(filepath.Join(dir, tc.filename))
					require.NoError(t, err)

					sd, err := w.Serialize()
					require.NoError(t, err)

					require.Equal(t, sd, data)
				}

				// create wallet with dup wallet name
				var otherSeed string
				var otherXPub string
				//var dupFingerprintErr error
				switch tc.walletType {
				case wallet.WalletTypeDeterministic,
					wallet.WalletTypeBip44:
					otherSeed = bip39.MustNewDefaultMnemonic()
				case wallet.WalletTypeXPub:
					otherXPub = "xpub6Ea7Vm9yPWhgrpmH7oTTc8vFmfp5Hpaf4ZpcjNWWJmpqr68viqmndJGkq6UFZcM6MpSXpqxF93PgvC7PuqByk5Pkx1XmcKMqkZhQbg21JXA"
				}

				_, err = s.CreateWallet(tc.filename, wallet.Options{
					Seed: otherSeed,
					Type: tc.walletType,
					XPub: otherXPub,
				})
				require.Equal(t, err, wallet.ErrWalletNameConflict)

				switch tc.walletType {
				case wallet.WalletTypeDeterministic,
					wallet.WalletTypeBip44, wallet.WalletTypeXPub:
					// create wallet with dup seed or xpub key
					dupWlt := "dup_wallet.wlt"
					_, err = s.CreateWallet(dupWlt, wallet.Options{
						Seed: tc.seed,
						XPub: tc.xpub,
						Type: tc.walletType,
					})
					require.Equal(t, wallet.NewError(fmt.Errorf("fingerprint conflict for %q wallet", tc.walletType)), err)

					// check that the dup wallet is not created
					_, err := s.GetWallet(dupWlt)
					require.Equal(t, err, wallet.ErrWalletNotExist)

					testutil.RequireFileNotExists(t, filepath.Join(dir, dupWlt))

				case wallet.WalletTypeCollection:
				// collection wallets never conflict with each other

				default:
					t.Fatal("unhandled wallet type")
				}
			})
		}
	}
}

type mockTxnsFinder map[cipher.Addresser]bool

func (mb mockTxnsFinder) AddressesActivity(addrs []cipher.Addresser) ([]bool, error) {
	if len(addrs) == 0 {
		return nil, nil
	}
	active := make([]bool, len(addrs))
	for i, addr := range addrs {
		active[i] = mb[addr]
	}
	return active, nil
}

func TestServiceLoadWallet(t *testing.T) {
	// Prepare addresses
	seed := "seed"
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 10)
	addrs := make([]cipher.Address, 10)
	for i, s := range seckeys {
		addrs[i] = cipher.MustAddressFromSecKey(s)
	}

	bip44Seed := "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	bip44AddrStrs := []string{
		"9BSEAEE3XGtQ2X43BCT2XCYgheGLQQigEG",
		"29cnQPHuWHCRF26LEAb2gR83ywnF3F9HduW",
		"2ZUAv9MGSpDKR3dnKMUnrKqLenV22JXAxzP",
		"fwNVThqdzH7JMsStoLrTpkVsemesbdGftm",
		"eyr5KDLTnN6ZZeggeHqDcXnrwmNUi7sGk2",
		"Aee3J9qoFPLoUEJes6YVzdKHdeuvCrMZeJ",
		"29MZS8aiYUdEwcruwCPggVJG9YJLsm92FHa",
		"2Hbm3bwKiEwqNAMAzVJmz5hL1dNTfaA3ju7",
		"WCaSCwSZnVqtkYeiKryeHjR8LbzE3KbkzJ",
		"baRjCy1yHfishGdZi3bVaPaL7VJM7FZCSd",
	}
	bip44Addrs := make([]cipher.Address, len(bip44AddrStrs))
	for i, a := range bip44AddrStrs {
		bip44Addrs[i] = cipher.MustDecodeBase58Address(a)
	}

	bip44SeedPassphrase := "foobar"
	bip44SeedPassphraseAddrStrs := []string{
		"n5SteDkkYdR3VJtMnVYcQ45L16rDDrseG8",
		"mGeG2PDoU4nc9qE1FSSreAjFeKG12zDvur",
		"rhbE3thvA747E81KfaYCujur7GKXjdhvS4",
		"BDEmcU8u4oTf9domk19Nzh65MXoWLLUvJN",
		"cubnvXGENW3gTdcdJADp8XEJaBscpy7gpq",
		"wv37cSiVhjgo6Qrrs994UJ52YU2zWNGJbu",
		"7aEzdSrcm1s2pm5YhshsRmkFy4EuYEnJ49",
		"nQJgxEE2eaggUeGaA73e4DaXq6KAvUiaS4",
		"2G9bhZaJrTKo1LScgtdvVXpQD4P8tKvgkvL",
		"4RqFK3qLz26XbPjgJsiJ3587P7p6DesDHd",
	}
	bip44SeedPassphraseAddrs := make([]cipher.Address, len(bip44SeedPassphraseAddrStrs))
	for i, a := range bip44SeedPassphraseAddrStrs {
		bip44SeedPassphraseAddrs[i] = cipher.MustDecodeBase58Address(a)
	}

	tt := []struct {
		name          string
		opts          wallet.Options
		tf            wallet.TransactionsFinder
		err           error
		expectAddrNum int
		expectAddrs   []cipher.Address
	}{
		{
			name: "raw wallet address=1",
			opts: wallet.Options{
				Type:  wallet.WalletTypeDeterministic,
				Seed:  seed,
				Label: "wallet",
				ScanN: 5,
				TF: mockTxnsFinder{
					addrs[0]: true,
				},
			},
			err:           nil,
			expectAddrNum: 1,
			expectAddrs:   addrs[:1],
		},
		{
			name: "raw wallet address=2",
			opts: wallet.Options{
				Type:  wallet.WalletTypeDeterministic,
				Seed:  seed,
				Label: "wallet",
				ScanN: 5,
				TF: mockTxnsFinder{
					addrs[1]: true,
				},
			},
			err:           nil,
			expectAddrNum: 2,
			expectAddrs:   addrs[:2],
		},
		{
			name: "encrypted wallet address=1",
			opts: wallet.Options{
				Type:     wallet.WalletTypeDeterministic,
				Seed:     seed,
				Label:    "wallet",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				TF: mockTxnsFinder{
					addrs[0]: true,
				},
			},
			err:           nil,
			expectAddrNum: 1,
			expectAddrs:   addrs[:1],
		},
		{
			name: "encrypted wallet address=2",
			opts: wallet.Options{
				Type:     wallet.WalletTypeDeterministic,
				Seed:     seed,
				Label:    "wallet",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				TF: mockTxnsFinder{
					addrs[1]: true,
				},
			},
			err:           nil,
			expectAddrNum: 2,
			expectAddrs:   addrs[:2],
		},
		{
			name: "bip44 raw wallet address=1",
			opts: wallet.Options{
				Type:  wallet.WalletTypeBip44,
				Seed:  bip44Seed,
				Label: "wallet",
				ScanN: 5,
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
				},
			},
			err:           nil,
			expectAddrNum: 1,
			expectAddrs:   bip44Addrs[:1],
		},
		{
			name: "bip44 raw wallet address=2",
			opts: wallet.Options{
				Type:  wallet.WalletTypeBip44,
				Seed:  bip44Seed,
				Label: "wallet",
				ScanN: 5,
				TF: mockTxnsFinder{
					bip44Addrs[1]: true,
				},
			},
			err:           nil,
			expectAddrNum: 2,
			expectAddrs:   bip44Addrs[:2],
		},
		{
			name: "bip44 encrypted wallet address=1",
			opts: wallet.Options{
				Type:     wallet.WalletTypeBip44,
				Seed:     bip44Seed,
				Label:    "wallet",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
				},
			},
			err:           nil,
			expectAddrNum: 1,
			expectAddrs:   bip44Addrs[:1],
		},
		{
			name: "bip44 encrypted wallet address=2",
			opts: wallet.Options{
				Type:     wallet.WalletTypeBip44,
				Seed:     bip44Seed,
				Label:    "wallet",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				TF: mockTxnsFinder{
					bip44Addrs[1]: true,
				},
			},
			err:           nil,
			expectAddrNum: 2,
			expectAddrs:   bip44Addrs[:2],
		},

		{
			name: "bip44 with seed passphrase raw wallet address=1",
			opts: wallet.Options{
				Type:           wallet.WalletTypeBip44,
				Seed:           bip44Seed,
				SeedPassphrase: bip44SeedPassphrase,
				Label:          "wallet",
				ScanN:          5,
				TF: mockTxnsFinder{
					bip44SeedPassphraseAddrs[0]: true,
				},
			},
			err:           nil,
			expectAddrNum: 1,
			expectAddrs:   bip44SeedPassphraseAddrs[:1],
		},
		{
			name: "bip44 with seed passphrase raw wallet address=2",
			opts: wallet.Options{
				Type:           wallet.WalletTypeBip44,
				Seed:           bip44Seed,
				SeedPassphrase: bip44SeedPassphrase,
				Label:          "wallet",
				ScanN:          5,
				TF: mockTxnsFinder{
					bip44SeedPassphraseAddrs[1]: true,
				},
			},
			err:           nil,
			expectAddrNum: 2,
			expectAddrs:   bip44SeedPassphraseAddrs[:2],
		},
		{
			name: "bip44 with seed passphrase encrypted wallet address=1",
			opts: wallet.Options{
				Type:           wallet.WalletTypeBip44,
				Seed:           bip44Seed,
				SeedPassphrase: bip44SeedPassphrase,
				Label:          "wallet",
				Encrypt:        true,
				Password:       []byte("pwd"),
				ScanN:          5,
				TF: mockTxnsFinder{
					bip44SeedPassphraseAddrs[0]: true,
				},
			},
			err:           nil,
			expectAddrNum: 1,
			expectAddrs:   bip44SeedPassphraseAddrs[:1],
		},
		{
			name: "bip44 with seed passphrase encrypted wallet address=2",
			opts: wallet.Options{
				Type:           wallet.WalletTypeBip44,
				Seed:           bip44Seed,
				SeedPassphrase: bip44SeedPassphrase,
				Label:          "wallet",
				Encrypt:        true,
				Password:       []byte("pwd"),
				ScanN:          5,
				TF: mockTxnsFinder{
					bip44SeedPassphraseAddrs[1]: true,
				},
			},
			err:           nil,
			expectAddrNum: 2,
			expectAddrs:   bip44SeedPassphraseAddrs[:2],
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: true,
				})
				require.NoError(t, err)
				wltName := wallet.NewWalletFilename()

				w, err := s.CreateWallet(wltName, tc.opts)
				//w, err := s.loadWallet(wltName, tc.opts, tc.tf)
				require.Equal(t, tc.err, err, fmt.Sprintf("expect: %v, got: %v", tc.err, err))
				if err != nil {
					return
				}

				el, err := w.EntriesLen()
				require.NoError(t, err)

				require.Equal(t, tc.expectAddrNum, el)
				for i, a := range tc.expectAddrs {
					e, err := w.GetEntryAt(i)
					require.NoError(t, err)
					require.Equal(t, a, e.Address.(cipher.Address))
				}

				require.Equal(t, w.IsEncrypted(), tc.opts.Encrypt)
				if w.IsEncrypted() {
					checkNoSensitiveData(t, w)
					// Checks the wallet file doesn't contain sensitive data
					wltPath := filepath.Join(dir, w.Filename())
					lw, err := s.Load(wltPath)
					require.NoError(t, err)
					checkNoSensitiveData(t, lw)
				}
			})
		}
	}
}

func TestServiceNewAddresses(t *testing.T) {
	seed := "seed"
	// Generate adddresses from the seed
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 10)
	addrs := make([]cipher.Address, 10)
	for i, s := range seckeys {
		addrs[i] = cipher.MustAddressFromSecKey(s)
	}

	bip44Seed := "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	bip44AddrStrs := []string{
		"9BSEAEE3XGtQ2X43BCT2XCYgheGLQQigEG",
		"29cnQPHuWHCRF26LEAb2gR83ywnF3F9HduW",
		"2ZUAv9MGSpDKR3dnKMUnrKqLenV22JXAxzP",
		"fwNVThqdzH7JMsStoLrTpkVsemesbdGftm",
		"eyr5KDLTnN6ZZeggeHqDcXnrwmNUi7sGk2",
		"Aee3J9qoFPLoUEJes6YVzdKHdeuvCrMZeJ",
		"29MZS8aiYUdEwcruwCPggVJG9YJLsm92FHa",
		"2Hbm3bwKiEwqNAMAzVJmz5hL1dNTfaA3ju7",
		"WCaSCwSZnVqtkYeiKryeHjR8LbzE3KbkzJ",
		"baRjCy1yHfishGdZi3bVaPaL7VJM7FZCSd",
	}
	bip44Addrs := make([]cipher.Address, len(bip44AddrStrs))
	for i, a := range bip44AddrStrs {
		bip44Addrs[i] = cipher.MustDecodeBase58Address(a)
	}

	tt := []struct {
		name               string
		opts               wallet.Options
		n                  uint64
		pwd                []byte
		walletAPIDisabled  bool
		walletFileModifier func(w string)
		expectAddrNum      int
		expectAddrs        []cipher.Address
		expectErr          error
	}{
		{
			name: "encrypted=false addresses=0",
			opts: wallet.Options{
				Type:  wallet.WalletTypeDeterministic,
				Label: "label",
				Seed:  seed,
			},
			n:             0,
			expectAddrNum: 0,
		},
		{
			name: "encrypted=false addresses=1",
			opts: wallet.Options{
				Label: "label",
				Seed:  seed,
				Type:  wallet.WalletTypeDeterministic,
			},
			n:             2,
			expectAddrNum: 2,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: addrs[1:3],
		},
		{
			name: "encrypted=false addresses=2",
			opts: wallet.Options{
				Label: "label",
				Seed:  seed,
				Type:  wallet.WalletTypeDeterministic,
			},
			n:             2,
			expectAddrNum: 2,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: addrs[1:3],
		},
		{
			name: "encrypted=true addresses=1",
			opts: wallet.Options{
				Label:    "label",
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			n:             1,
			pwd:           []byte("pwd"),
			expectAddrNum: 1,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: addrs[1:2],
		},
		{
			name: "encrypted=true addresses=2",
			opts: wallet.Options{
				Label:    "label",
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			n:             2,
			pwd:           []byte("pwd"),
			expectAddrNum: 2,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: addrs[1:3],
		},
		{
			name: "bip44 encrypted=false addresses=0",
			opts: wallet.Options{
				Type:  wallet.WalletTypeBip44,
				Label: "label",
				Seed:  bip44Seed,
			},
			n:             0,
			expectAddrNum: 0,
		},
		{
			name: "bip44 encrypted=false addresses=1",
			opts: wallet.Options{
				Label: "label",
				Seed:  bip44Seed,
				Type:  wallet.WalletTypeBip44,
			},
			n:             2,
			expectAddrNum: 2,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: bip44Addrs[1:3],
		},
		{
			name: "bip44 encrypted=false addresses=2",
			opts: wallet.Options{
				Label: "label",
				Seed:  bip44Seed,
				Type:  wallet.WalletTypeBip44,
			},
			n:             2,
			expectAddrNum: 2,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: bip44Addrs[1:3],
		},
		{
			name: "bip44 encrypted=true addresses=1",
			opts: wallet.Options{
				Label:    "label",
				Seed:     bip44Seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeBip44,
			},
			n:             1,
			pwd:           []byte("pwd"),
			expectAddrNum: 1,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: bip44Addrs[1:2],
		},
		{
			name: "bip44 encrypted=true addresses=2",
			opts: wallet.Options{
				Label:    "label",
				Seed:     bip44Seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeBip44,
			},
			n:             2,
			pwd:           []byte("pwd"),
			expectAddrNum: 2,
			// CreateWallet will generate a default address, so check from new address
			expectAddrs: bip44Addrs[1:3],
		},
		{
			name: "encrypted=true wrong password",
			opts: wallet.Options{
				Label:    "label",
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			n:             1,
			pwd:           []byte("wrong password"),
			expectAddrNum: 1,
			expectErr:     wallet.ErrInvalidPassword,
		},
		{
			name: "wallet api disabled",
			opts: wallet.Options{
				Seed:  "seed",
				Label: "label",
				Type:  wallet.WalletTypeDeterministic,
			},
			walletAPIDisabled: true,
			expectErr:         wallet.ErrWalletAPIDisabled,
		},
		{
			name: "encrypted=false password provided",
			opts: wallet.Options{
				Label: "label",
				Seed:  seed,
				Type:  wallet.WalletTypeDeterministic,
			},
			n:         1,
			pwd:       []byte("foo"),
			expectErr: wallet.ErrWalletNotEncrypted,
		},
		{
			name: "encrypted=false writable=false",
			opts: wallet.Options{
				Label: "label",
				Seed:  seed,
				Type:  wallet.WalletTypeDeterministic,
			},
			n: 1,
			walletFileModifier: func(fn string) {
				err := os.Chmod(fn, 0555) // no write permission to the wallet file
				require.NoError(t, err)
			},
			expectAddrNum: 1,
			expectErr:     wallet.ErrWalletPermission,
		},
	}

	for _, tc := range tt {
		//for _, ct := range crypto.TypesInsecure() {
		for _, ct := range []crypto.CryptoType{crypto.CryptoTypeSha256Xor} {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.walletAPIDisabled,
				})
				require.NoError(t, err)

				opts := tc.opts
				opts.CryptoType = ct

				wltName := wallet.NewWalletFilename()

				w, err := s.CreateWallet(wltName, opts)
				if err != nil {
					require.Equal(t, tc.expectErr, err, fmt.Sprintf("expect: %v, got: %v", tc.expectErr, err))
					return
				}

				wltPath := filepath.Join(dir, w.Filename())
				if tc.walletFileModifier != nil {
					tc.walletFileModifier(wltPath)
				}

				if w.IsEncrypted() {
					checkNoSensitiveData(t, w)
				}

				naddrs, err := s.NewAddresses(w.Filename(), tc.pwd, tc.n)
				require.Equal(t, tc.expectErr, err)

				// Confirms that no intermediate tmp file exists
				tmpWltPath := filepath.Join(dir, w.Filename()) + ".tmp"
				_, existErr := os.Stat(tmpWltPath)
				require.True(t, os.IsNotExist(existErr))

				if err != nil {
					return
				}

				// Confirms that the wallet addresses number is correct
				require.Len(t, naddrs, tc.expectAddrNum)
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, naddrs[i])
				}

				// Check the wallet again
				w, err = s.GetWallet(wltName)
				require.NoError(t, err)
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, int(tc.n+1), el)

				addrsInWlt, err := w.GetAddresses()
				require.NoError(t, err)

				// Wallet has a default address, so need to start from the second address
				for i, a := range tc.expectAddrs {
					require.Equal(t, a, addrsInWlt[i+1])
				}

				// Load wallet from file and check
				_, err = os.Stat(filepath.Join(dir, w.Filename()))
				require.NoError(t, err)

				lw, err := s.Load(filepath.Join(dir, w.Filename()))
				require.NoError(t, err)

				el, err = lw.EntriesLen()
				require.NoError(t, err)

				require.Equal(t, tc.expectAddrNum+1, el)
				es, err := w.GetEntries()
				require.NoError(t, err)
				for i := range es {
					e, err := lw.GetEntryAt(i)
					require.NoError(t, err)
					require.Equal(t, e, es[i])
				}

				if w.IsEncrypted() {
					checkNoSensitiveData(t, lw)
				}

				// Wallet doesn't exist
				_, err = s.NewAddresses("wallet_not_exist.wlt", tc.pwd, 1)
				require.Equal(t, wallet.ErrWalletNotExist, err)
			})
		}
	}
}

func TestServiceGetAddress(t *testing.T) {
	for _, enableWalletAPI := range []bool{true, false} {
		for _, ct := range crypto.TypesInsecure() {
			t.Run(fmt.Sprintf("enable wallet api=%v crypto=%v", enableWalletAPI, ct), func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
					WalletDir:       "./testdata",
					CryptoType:      ct,
					EnableWalletAPI: enableWalletAPI,
				})
				require.NoError(t, err)

				if !enableWalletAPI {
					dirIsEmpty(t, dir)

					wlts, err := s.GetWallets()
					require.Equal(t, wallet.ErrWalletAPIDisabled, err)
					require.Empty(t, wlts)
					addrs, err := s.GetAddresses("")
					require.Equal(t, wallet.ErrWalletAPIDisabled, err)
					require.Equal(t, 0, len(addrs))
					return
				}

				addrs, err := s.GetAddresses("test1.wlt")
				require.NoError(t, err)
				require.Equal(t, 1, len(addrs))

				// test none exist wallet
				notExistID := "not_exist_id.wlt"
				_, err = s.GetAddresses(notExistID)
				require.Equal(t, wallet.ErrWalletNotExist, err)
			})
		}
	}
}

func TestServiceGetWallet(t *testing.T) {
	for _, walletType := range []string{wallet.WalletTypeDeterministic} {
		for _, enableWalletAPI := range []bool{true, false} {
			for _, ct := range crypto.Types() {
				t.Run(fmt.Sprintf("enable wallet api=%v crypto=%v", enableWalletAPI, ct), func(t *testing.T) {
					dir := prepareWltDir()
					s, err := wallet.NewService(wallet.Config{
						WalletDir:       dir,
						CryptoType:      ct,
						EnableWalletAPI: enableWalletAPI,
					})
					require.NoError(t, err)

					if !enableWalletAPI {
						dirIsEmpty(t, dir)
						wlts, err := s.GetWallets()
						require.Equal(t, wallet.ErrWalletAPIDisabled, err)
						require.Empty(t, wlts)
						w, err := s.GetWallet("")
						require.Equal(t, wallet.ErrWalletAPIDisabled, err)
						require.Nil(t, w)
						return
					}

					opts := wallet.Options{
						Label: "label",
						Type:  walletType,
					}
					switch walletType {
					case wallet.WalletTypeBip44, wallet.WalletTypeDeterministic:
						opts.Seed = bip39.MustNewDefaultMnemonic()
						// TODO: check other wallet types
						//case WalletTypeCollection:
						//case WalletTypeXPub:
						//	opts.XPub = "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v"
						//default:
						//	t.Fatal("unhandled wallet type")
					}

					// Create a wallet
					w, err := s.CreateWallet("t.wlt", opts)
					require.NoError(t, err)

					w1, err := s.GetWallet(w.Filename())
					require.NoError(t, err)

					// Check if change original wallet would change the returned wallet
					w.SetLabel("new_label")

					require.NotEqual(t, "new_label", w1.Label())

					// Get wallet doesn't exist
					wltName := "does_not_exist.wlt"
					_, err = s.GetWallet(wltName)
					require.Equal(t, wallet.ErrWalletNotExist, err)
				})
			}
		}
	}
}

func TestServiceGetWallets(t *testing.T) {
	for _, enableWalletAPI := range []bool{true, false} {
		for _, ct := range crypto.TypesInsecure() {
			t.Run(fmt.Sprintf("enable wallet=%v crypto=%v", enableWalletAPI, ct), func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: enableWalletAPI,
				})
				require.NoError(t, err)

				if !enableWalletAPI {
					dirIsEmpty(t, dir)

					w, err := s.GetWallets()
					require.Equal(t, wallet.ErrWalletAPIDisabled, err)
					var emptyW wallet.Wallets
					require.Equal(t, w, emptyW)
					return
				}

				// Creates a bip44 wallet
				w, err := s.CreateWallet("t.wlt", wallet.Options{
					Label: "label",
					Seed:  bip39.MustNewDefaultMnemonic(),
					Type:  wallet.WalletTypeBip44,
				})
				require.NoError(t, err)

				var wallets []wallet.Wallet
				// Get the default wallet
				wallets = append(wallets, w)

				// Create a new wallet
				wltName := wallet.NewWalletFilename()
				w1, err := s.CreateWallet(wltName, wallet.Options{
					Label: "label1",
					Seed:  bip39.MustNewDefaultMnemonic(),
					Type:  wallet.WalletTypeDeterministic,
				})
				require.NoError(t, err)
				wallets = append(wallets, w1)

				ws, err := s.GetWallets()
				require.NoError(t, err)
				for _, w := range wallets {
					ww, ok := ws[w.Filename()]
					require.True(t, ok)

					wb, err := w.Serialize()
					require.NoError(t, err)

					wwb, err := ww.Serialize()
					require.NoError(t, err)
					require.Equal(t, wb, wwb)
				}
			})
		}
	}
}

func TestServiceUpdateWalletLabel(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             wallet.Options
		updateWltName    string
		label            string
		disableWalletAPI bool
		err              error
	}{
		{
			name:    "ok",
			wltName: "t.wlt",
			opts: wallet.Options{
				Seed:  bip39.MustNewDefaultMnemonic(),
				Label: "label",
			},
			updateWltName: "t.wlt",
			label:         "new-label",
		},
		{
			name:    "wallet doesn't exist",
			wltName: "t.wlt",
			opts: wallet.Options{
				Seed:  bip39.MustNewDefaultMnemonic(),
				Label: "label",
			},
			updateWltName: "t1.wlt",
			label:         "new-label",
			err:           wallet.ErrWalletNotExist,
		},
		{
			name:    "wallet api disabled",
			wltName: "t.wlt",
			opts: wallet.Options{
				Seed:  bip39.MustNewDefaultMnemonic(),
				Label: "label",
			},
			disableWalletAPI: true,
			err:              wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		// TODO: add bip44 test
		for _, walletType := range []string{
			wallet.WalletTypeDeterministic,
			wallet.WalletTypeBip44,
		} {
			for _, ct := range crypto.Types() {
				t.Run(tc.name, func(t *testing.T) {
					// Create the wallet service
					dir := prepareWltDir()
					s, err := wallet.NewService(wallet.Config{
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
					tc.opts.Type = walletType
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
}

func TestServiceEncryptWallet(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             wallet.Options
		encWltName       string
		pwd              []byte
		disableWalletAPI bool
		err              error
	}{
		{
			name:    "ok xpub wallet",
			wltName: "t.wlt",
			opts: wallet.Options{
				XPub: "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v",
				Type: wallet.WalletTypeXPub,
			},
			encWltName: "t.wlt",
			pwd:        []byte("pwd"),
			err:        wallet.NewError(fmt.Errorf("xpub wallet does not support encryption")),
		},
		{
			name:    "ok deterministic wallet",
			wltName: "t.wlt",
			opts: wallet.Options{
				Seed: "seed",
				Type: wallet.WalletTypeDeterministic,
			},
			encWltName: "t.wlt",
			pwd:        []byte("pwd"),
		},
		{
			name:    "ok collection wallet",
			wltName: "t.wlt",
			opts: wallet.Options{
				Type: wallet.WalletTypeCollection,
			},
			encWltName: "t.wlt",
			pwd:        []byte("pwd"),
		},
		{
			name:    "ok bip44 wallet",
			wltName: "t.wlt",
			opts: wallet.Options{
				Type: wallet.WalletTypeBip44,
				Seed: "voyage say extend find sheriff surge priority merit ignore maple cash argue",
			},
			encWltName: "t.wlt",
			pwd:        []byte("pwd"),
		},
		{
			name:    "wallet doesn't exist",
			wltName: "t.wlt",
			opts: wallet.Options{
				Seed: "seed",
				Type: wallet.WalletTypeDeterministic,
			},
			encWltName: "t2.wlt",
			err:        wallet.ErrWalletNotExist,
		},
		{
			name:    "wallet already encrypted",
			wltName: "t.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Type:     wallet.WalletTypeDeterministic,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			encWltName: "t.wlt",
			pwd:        []byte("pwd"),
			err:        wallet.ErrWalletEncrypted,
		},
		{
			name:    "wallet api disabled",
			wltName: "t.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Type:     wallet.WalletTypeDeterministic,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			encWltName:       "t.wlt",
			pwd:              []byte("pwd"),
			disableWalletAPI: true,
			err:              wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				// Create the wallet service
				s, err := wallet.NewService(wallet.Config{
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
				opts := tc.opts
				opts.CryptoType = ct
				w, err := s.CreateWallet(tc.wltName, opts)
				require.NoError(t, err)

				switch w.Type() {
				// Add an entry to a collection wallet, to verify that secrets are hidden
				case wallet.WalletTypeCollection:
					// TODO: open this after collection is supported
					err := s.Update(w.Filename(), func(w wallet.Wallet) error {
						p, s := cipher.GenerateKeyPair()
						return w.(*collection.Wallet).AddEntry(wallet.Entry{
							Public:  p,
							Secret:  s,
							Address: cipher.AddressFromPubKey(p),
						})
					})
					require.NoError(t, err)

				// Add entries to the a bip44 wallet's change chain, to verify that those secrets are hidden
				case wallet.WalletTypeBip44:
					err := s.Update(w.Filename(), func(w wallet.Wallet) error {
						_, err := w.(*bip44wallet.Wallet).PeekChangeAddress(mockTxnsFinder{})
						return err
					})
					require.NoError(t, err)
				}

				// Encrypt the wallet
				encWlt, err := s.EncryptWallet(tc.encWltName, tc.pwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				encWlt1, err := s.GetWallet(tc.encWltName)
				//encWlt1, err := s.getWallet(tc.encWltName)
				require.NoError(t, err)
				ewb, err := encWlt.Serialize()
				require.NoError(t, err)
				ew1b, err := encWlt1.Serialize()
				require.NoError(t, err)
				require.Equal(t, ewb, ew1b)

				// Check the encrypted wallet
				require.True(t, encWlt.IsEncrypted())
				entries, err := encWlt.GetEntries()
				require.NoError(t, err)
				for _, e := range entries {
					require.True(t, e.Secret.Null())
				}
				require.Empty(t, encWlt.Seed())
				require.Empty(t, encWlt.LastSeed())

				// Check the decrypted seeds
				decWlt, err := encWlt.Unlock(tc.pwd)
				require.NoError(t, err)

				entries, err = decWlt.GetEntries()
				require.NoError(t, err)
				for _, e := range entries {
					switch decWlt.Type() {
					case wallet.WalletTypeXPub:
						// xpub wallets never have secret keys
						require.True(t, e.Secret.Null())
					default:
						require.False(t, e.Secret.Null())
					}

				}

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
	type testCase struct {
		name             string
		wltName          string
		opts             wallet.Options
		decryptWltName   string
		password         []byte
		disableWalletAPI bool
		err              error
	}

	tt := []testCase{
		{
			name:    "ok xpub",
			wltName: "test.wlt",
			opts: wallet.Options{
				Type: wallet.WalletTypeXPub,
				XPub: "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v",
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
			// xpub wallet does not support encryption,
			// hence decrypts wallet would only return wallet is not encrypted
			err: wallet.NewError(errors.New("wallet is not encrypted")),
		},
		{
			name:    "ok collection",
			wltName: "test.wlt",
			opts: wallet.Options{
				Type:     wallet.WalletTypeCollection,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
		},
		{
			name:    "ok deterministic",
			wltName: "test.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
		},
		{
			name:    "ok bip44",
			wltName: "test.wlt",
			opts: wallet.Options{
				Seed:     "voyage say extend find sheriff surge priority merit ignore maple cash argue",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeBip44,
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
		},
		{
			name:    "wallet not exist",
			wltName: "test.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			decryptWltName: "t.wlt",
			password:       []byte("pwd"),
			err:            wallet.ErrWalletNotExist,
		},
		{
			name:    "wallet not encrypted",
			wltName: "test.wlt",
			opts: wallet.Options{
				Seed: "seed",
				Type: wallet.WalletTypeDeterministic,
			},
			decryptWltName: "test.wlt",
			password:       []byte("pwd"),
			err:            wallet.ErrWalletNotEncrypted,
		},
		{
			name:    "invalid password",
			wltName: "test.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			decryptWltName: "test.wlt",
			password:       []byte("wrong password"),
			err:            wallet.ErrInvalidPassword,
		},
		{
			name:    "wallet api disabled",
			wltName: "test.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			decryptWltName:   "test.wlt",
			password:         []byte("pwd"),
			disableWalletAPI: true,
			err:              wallet.ErrWalletAPIDisabled,
		},
	}

	verifyDecryptedXPubWlt := func(tc testCase, wlt wallet.Wallet) {
		// XPub wlt doesn't have anything to encrypt or decrypt
		require.Equal(t, tc.opts.XPub, wlt.XPub())
		require.Empty(t, wlt.Secrets())
		require.Empty(t, wlt.Seed())
		require.Empty(t, wlt.LastSeed())
		entries, err := wlt.GetEntries()
		require.NoError(t, err)

		for _, e := range entries {
			require.True(t, e.Secret.Null())
		}
	}

	verifyDecryptedDeterministicWlt := func(tc testCase, wlt wallet.Wallet) {
		// Checks the "encrypted" meta info
		require.False(t, wlt.IsEncrypted())
		// Checks the seed
		require.Equal(t, tc.opts.Seed, wlt.Seed())

		// Checks the last seed
		entryNum, err := wlt.EntriesLen()
		require.NoError(t, err)
		lsd, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(wlt.Seed()), entryNum)
		require.Equal(t, hex.EncodeToString(lsd), wlt.LastSeed())

		// Checks the entries
		entries, err := wlt.GetEntries()
		require.NoError(t, err)

		for i := range seckeys {
			a := cipher.MustAddressFromSecKey(seckeys[i])
			require.Equal(t, a, entries[i].Address)
			require.Equal(t, seckeys[i], entries[i].Secret)
		}

		require.Empty(t, wlt.Secrets())
	}

	verifyDecryptedCollectionWlt := func(_ testCase, wlt wallet.Wallet) {
		// Checks the "encrypted" meta info
		require.False(t, wlt.IsEncrypted())
		require.Empty(t, wlt.Seed())
		require.Empty(t, wlt.LastSeed())

		// Checks the entries
		entries, err := wlt.GetEntries()
		require.NoError(t, err)

		for _, e := range entries {
			require.False(t, e.Secret.Null())
			a := cipher.MustAddressFromSecKey(e.Secret)
			require.Equal(t, a, e.Address)
			p := cipher.MustPubKeyFromSecKey(e.Secret)
			require.Equal(t, p, e.Public)
		}

		require.Empty(t, wlt.Secrets())
	}

	verifyDecryptedBip44Wlt := func(tc testCase, wlt wallet.Wallet) {
		// Checks the "encrypted" meta info
		require.False(t, wlt.IsEncrypted())
		// Checks the seed
		require.Equal(t, tc.opts.Seed, wlt.Seed())
		require.Empty(t, wlt.LastSeed())

		// Checks the entries
		entries, err := wlt.GetEntries()
		require.NoError(t, err)
		for _, e := range entries {
			require.False(t, e.Secret.Null())
			a := cipher.MustAddressFromSecKey(e.Secret)
			require.Equal(t, a, e.Address)
			p := cipher.MustPubKeyFromSecKey(e.Secret)
			require.Equal(t, p, e.Public)
		}

		require.Empty(t, wlt.Secrets())
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
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

				_, err = s.CreateWallet(tc.wltName, tc.opts)
				require.NoError(t, err)

				_, err = s.DecryptWallet(tc.decryptWltName, tc.password)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				wltType := tc.opts.Type
				if wltType == "" {
					wltType = wallet.WalletTypeBip44
				}

				//verify := verifyDecryptedCollectionWlt
				verify := verifyDecryptedDeterministicWlt
				switch wltType {
				case wallet.WalletTypeCollection:
					verify = verifyDecryptedCollectionWlt
				case wallet.WalletTypeBip44:
					verify = verifyDecryptedBip44Wlt
				case wallet.WalletTypeDeterministic:
					verify = verifyDecryptedDeterministicWlt
				case wallet.WalletTypeXPub:
					verify = verifyDecryptedXPubWlt
				default:
					t.Fatal("unhandled wallet type")
				}

				// Checks the decrypted wallet in service
				w, err := s.GetWallet(tc.wltName)
				require.NoError(t, err)
				verify(tc, w)

				// Checks the existence of the wallet file
				fn := filepath.Join(dir, tc.wltName)
				testutil.RequireFileExists(t, fn)

				// Loads wallet from the file and check if it's decrypted
				w1, err := s.Load(fn)
				require.NoError(t, err)
				verify(tc, w1)
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

	bip44Seed := "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	bip44AddrStrs := []string{
		"9BSEAEE3XGtQ2X43BCT2XCYgheGLQQigEG",
		"29cnQPHuWHCRF26LEAb2gR83ywnF3F9HduW",
		"2ZUAv9MGSpDKR3dnKMUnrKqLenV22JXAxzP",
		"fwNVThqdzH7JMsStoLrTpkVsemesbdGftm",
		"eyr5KDLTnN6ZZeggeHqDcXnrwmNUi7sGk2",
		"Aee3J9qoFPLoUEJes6YVzdKHdeuvCrMZeJ",
		"29MZS8aiYUdEwcruwCPggVJG9YJLsm92FHa",
		"2Hbm3bwKiEwqNAMAzVJmz5hL1dNTfaA3ju7",
		"WCaSCwSZnVqtkYeiKryeHjR8LbzE3KbkzJ",
		"baRjCy1yHfishGdZi3bVaPaL7VJM7FZCSd",
		"296oQmJJgx35NDApi7YYzj1AryM8fZcjwf3",
		"cxxxRfy3RRy2YbFTcptRbVTQYcHY1ejRB5",
		"omLGQm1Z2Y9Bga8v6NQ2hgrpRm1nATzGK9",
		"2EpZP1E8gTJy799t5CVrZUcxjyFFHwshr6X",
		"2hgaPG2oNVrkonPxjv4Sx9au6ruw1Y8pjUi",
		"2bHfa8yjhWB5mTip8j1FjNhB1TGbSBkX3Xu",
		"VYu5ePSB7ReKm2pysC5JRdCUiTBgDn5Tkw",
		"2crRqwG3BaurEqNa7eiB5oUNKTQPETfKrFW",
		"2LVTqqNSTBKE51UC7bZ39bZ6wwmR3sibHBX",
		"2A8C3h1gsw92Q4Uhn4b385onKrhzuH8UTwE",
	}
	bip44Addrs := make([]cipher.Address, len(bip44AddrStrs))
	for i, a := range bip44AddrStrs {
		bip44Addrs[i] = cipher.MustDecodeBase58Address(a)
	}

	bip44ChangeAddrStrs := []string{
		"oHvj7oy8maES9HJiQHJTp4GvcUcpz3voDq",
		"2SGMfTFV2zbQzGw7aJm1D5EeEPgych5ixuC",
		"2ymjULRdbiFoUNJKNhWbQ3JqdE8TXnZkyU",
		"muvdio7V8vkbUPPJsumVWEqScKHZho6Xmx",
		"qyQuA2RW4H6NFqhRQmHDt2q28E654PoBsH",
		"24UcBj42Q7GBx1rmeU5AE7JhcRbpDq9utCV",
		"2TwjgBxKhK84Qa34HoNc5KbsYfp1NBZyPXk",
		"fDYKWyDnDaQmpigGAJNK4ZoBZ4WTg5dtRN",
		"rRMXXm8ufMcqqnLqRX5CoJ8wA6rdWPJEiW",
		"2D6HN4fSHZqsGaPzCDUjddF6NS4hBsVGRag",
		"2bXjTvatPW4Z3SaQBr96zJwwYcxp7iemeGC",
		"tbmXTJWtqtuaAdVvFWS1rmL42omdD4Hkd3",
		"27NT6t4xJoqcs6PsM2EgNNhx4MMw3ioFVER",
		"Xym6pAVg2Xjp9sXvCUKrLgq6PSEt1izsRd",
		"49EqkHPgWPS9W1jDiJXSKdKdNur3pJjADw",
		"f52mwYtpVNhgMMmJZHQYVBQVHwHXCbP5yX",
		"2ffeXAKtVUmxucs8mpM2EKMqzYCN3TB7DYx",
		"2NJERX84xLHMidXtYDtT1rGzhf4HJmCkz47",
		"wik4VEcja3pz2Wo6SwwfoDGzVzB7ZoosQB",
		"2gj3qjYQ4MMkv7CdB3PvkbEEPhXikZBD4q6",
	}
	bip44ChangeAddrs := make([]cipher.Address, len(bip44ChangeAddrStrs))
	for i, a := range bip44ChangeAddrStrs {
		bip44ChangeAddrs[i] = cipher.MustDecodeBase58Address(a)
	}

	xpub := "xpub6E5WPk37XdM79dy6oJ7iH6NkCvVzxmrCo4zMFFHSZMc5ymZYhReQFWaDcGNZeYYe1ahY2e3RcRZDHLHC98FfzPRfNRcU6ecURpS4RCQRP2w"
	xpubAddrStrs := []string{
		"2mhaS6SE2TPSmRRbJvngWQSNXCCVuTic5Zg",
		"2bq2itwDKteqigxVS9eYJv4Ww9SEfuyGcib",
		"B7eMXM6nLUqqzkFcosXR3HSVkQ6yUz53n4",
		"niAy17kBb8vB2pFey8eZnE92e6x9bFGLHp",
		"N8JbzcqWEPkn6CF3JdZDmEydzECyZ4NhRv",
		"TcyY3F4xHPCtFFkdDBTC93y684Fmxg2rPd",
		"kcVFbcrVqAVUHrirp7r3HYHUzrtdFuybez",
		"bzA7UeUmkuFWn3waGh3z1eQ5xV3TAZpgX2",
		"2K29ZX6vaqrTRZJbFKX7hzu646wL3pJfF6H",
		"VKpFCpN4yp46uYbffaeCg7XEvQd6pHNSkV",
		"2bGUCmS8BcFiX6VQbq6DVvjDsCdu5fgoSSr",
		"2Yx8dzMgzU5Y2vrAomAtkajWNmUCJ31xNF",
		"wd2NGBkCygq7cCP899gbiBciM7ZFqphRDM",
		"2akitSsnetXoc3ejzY4pA8dGjWM76uzPxo3",
		"WFTbYvNJGAq9wWG644sMvgt4EJ6CvkbvDK",
		"2FReiyjcQuQBCvjKQzyrhc7QMwdbdBwLvCT",
		"2RnSh4sZrxCEK5fUSCLR8JuwzGd6K5mzdrA",
		"ujJHitcLhQZB393qrdJvdfM4AyEPMmDhcV",
		"25jWWmrgU8Z9HYVCsAXYmCjF1jQqe5DDnqr",
		"rsaKzohU5erbR6FX1whWu9Ke4q2jLkLBeJ",
	}
	xpubAddrs := make([]cipher.Address, len(xpubAddrStrs))
	for i, a := range xpubAddrStrs {
		xpubAddrs[i] = cipher.MustDecodeBase58Address(a)
	}

	tf := make(mockTxnsFinder, 20)

	type exp struct {
		err      error
		seed     string
		lastSeed string
		xpub     string
		entryNum int
		addrs    []cipher.Address
	}

	tt := []struct {
		name             string
		opts             wallet.Options
		disableWalletAPI bool
		expect           exp
	}{
		{
			name: "no coins and scan 0, unencrypted",
			opts: wallet.Options{
				Seed: seed,
				Type: wallet.WalletTypeDeterministic,
				TF:   tf,
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "no coins and scan 0, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
				TF:       tf,
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "no coins and scan 1, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 1,
				Type:  wallet.WalletTypeDeterministic,
				TF:    tf,
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "no coins and scan 1, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    1,
				Type:     wallet.WalletTypeDeterministic,
				TF:       tf,
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "no coins and scan 10, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 10,
				Type:  wallet.WalletTypeDeterministic,
				TF:    tf,
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "scan 1 get 1, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 1,
				Type:  wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "scan 5 get 5, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 5,
				Type:  wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[4]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[4],
				entryNum: 5,
				addrs:    addrs,
			},
		},
		{
			name: "scan 5 get 1, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 5,
				Type:  wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "scan 5 get 2, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 5,
				Type:  wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[1],
				entryNum: 2,
				addrs:    addrs,
			},
		},

		{
			name: "scan 5 get 3, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				Type:     wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[2]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[2],
				entryNum: 3,
				addrs:    addrs,
			},
		},
		{
			name: "scan 5 get 4, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 5,
				Type:  wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[2]: true,
					addrs[3]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[3],
				entryNum: 4,
				addrs:    addrs,
			},
		},
		{
			name: "scan 5 get 5, unencrypted",
			opts: wallet.Options{
				Seed:  seed,
				ScanN: 5,
				Type:  wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[3]: true,
					addrs[4]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[4],
				entryNum: 5,
				addrs:    addrs,
			},
		},
		{
			name: "scan 0 get 1, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    0,
				Type:     wallet.WalletTypeDeterministic,
				TF:       mockTxnsFinder{},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "scan 1 get 1, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    1,
				Type:     wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "scan 2 get 1, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    2,
				Type:     wallet.WalletTypeDeterministic,
				TF:       mockTxnsFinder{},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[0],
				entryNum: 1,
				addrs:    addrs,
			},
		},
		{
			name: "scan 2 get 2, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    2,
				Type:     wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[1],
				entryNum: 2,
				addrs:    addrs,
			},
		},
		{
			name: "scan 5 get 5, encrypted",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				Type:     wallet.WalletTypeDeterministic,
				TF: mockTxnsFinder{
					addrs[3]: true,
					addrs[4]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     seed,
				lastSeed: childSeeds[4],
				entryNum: 5,
				addrs:    addrs,
			},
		},
		{
			name: "bip44 no coins and scan 0, unencrypted",
			opts: wallet.Options{
				Seed: bip44Seed,
				Type: wallet.WalletTypeBip44,
				TF:   tf,
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 no coins and scan 0, encrypted",
			opts: wallet.Options{
				Seed:     bip44Seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeBip44,
				TF:       tf,
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 no coins and scan 1, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 1,
				Type:  wallet.WalletTypeBip44,
				TF:    tf,
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 no coins and scan 1, encrypted",
			opts: wallet.Options{
				Seed:     bip44Seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    1,
				Type:     wallet.WalletTypeBip44,
				TF:       tf,
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 no coins and scan 10, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 10,
				Type:  wallet.WalletTypeBip44,
				TF:    tf,
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 0 get 1, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 1,
				Type:  wallet.WalletTypeBip44,
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 0 get 1, encrypted",
			opts: wallet.Options{
				Seed:     bip44Seed,
				ScanN:    1,
				Type:     wallet.WalletTypeBip44,
				Encrypt:  true,
				Password: []byte("pwd"),
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 1 get 1, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 1,
				Type:  wallet.WalletTypeBip44,
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 1 get 1, encrypted",
			opts: wallet.Options{
				Seed:     bip44Seed,
				ScanN:    1,
				Type:     wallet.WalletTypeBip44,
				Encrypt:  true,
				Password: []byte("pwd"),
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 2 get 2, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 2,
				Type:  wallet.WalletTypeBip44,
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
					bip44Addrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 2,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 2 get 2, encrypted",
			opts: wallet.Options{
				Seed:     bip44Seed,
				ScanN:    2,
				Type:     wallet.WalletTypeBip44,
				Encrypt:  true,
				Password: []byte("pwd"),
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
					bip44Addrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 2,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 5 get 1, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 5,
				Type:  wallet.WalletTypeBip44,
				TF: mockTxnsFinder{
					bip44Addrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 1,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 5 get 2, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 5,
				Type:  wallet.WalletTypeBip44,
				TF: mockTxnsFinder{
					bip44Addrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 2,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 5 get 5, unencrypted",
			opts: wallet.Options{
				Seed:  bip44Seed,
				ScanN: 5,
				Type:  wallet.WalletTypeBip44,
				TF: mockTxnsFinder{
					bip44Addrs[4]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 5,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "bip44 scan 5 get 5, encrypted",
			opts: wallet.Options{
				Seed:     bip44Seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				Type:     wallet.WalletTypeBip44,
				TF: mockTxnsFinder{
					bip44Addrs[3]: true,
					bip44Addrs[4]: true,
				},
			},
			expect: exp{
				err:      nil,
				seed:     bip44Seed,
				entryNum: 5,
				addrs:    bip44Addrs,
			},
		},
		{
			name: "xpub no coins and scan 0, unencrypted",
			opts: wallet.Options{
				XPub: xpub,
				Type: wallet.WalletTypeXPub,
				TF:   tf,
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub no coins and scan 0",
			opts: wallet.Options{
				XPub: xpub,
				Type: wallet.WalletTypeXPub,
				TF:   tf,
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub no coins and scan 1, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 1,
				Type:  wallet.WalletTypeXPub,
				TF:    tf,
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub no coins and scan 1",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 1,
				Type:  wallet.WalletTypeXPub,
				TF:    tf,
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub no coins and scan 10, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 10,
				Type:  wallet.WalletTypeXPub,
				TF:    tf,
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 0 get 1, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 0,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 0 get 1",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 0,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 1 get 1, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 1,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 1 get 1",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 1,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 2 get 1, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 2,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 2 get 1",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 2,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 2 get 2, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 2,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 2,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 2 get 2",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 2,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 2,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 5 get 5, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 5,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[4]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 5,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 5 get 4, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 5,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[1]: true,
					xpubAddrs[3]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 4,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 5 get 3, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 5,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[1]: true,
					xpubAddrs[2]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 3,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 5 get 2, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 5,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 2,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 5 get 1, unencrypted",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 5,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[0]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 1,
				addrs:    xpubAddrs,
			},
		},

		{
			name: "xpub scan 5 get 2",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 5,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[1]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 2,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "xpub scan 5 get 5",
			opts: wallet.Options{
				XPub:  xpub,
				ScanN: 5,
				Type:  wallet.WalletTypeXPub,
				TF: mockTxnsFinder{
					xpubAddrs[3]: true,
					xpubAddrs[4]: true,
				},
			},
			expect: exp{
				err:      nil,
				xpub:     xpub,
				entryNum: 5,
				addrs:    xpubAddrs,
			},
		},
		{
			name: "wallet api disabled",
			opts: wallet.Options{
				Seed:     seed,
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    5,
				Type:     wallet.WalletTypeDeterministic,
				TF:       mockTxnsFinder{},
			},
			disableWalletAPI: true,
			expect: exp{
				err: wallet.ErrWalletAPIDisabled,
			},
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
				})
				require.NoError(t, err)

				wltName := wallet.NewWalletFilename()
				w, err := s.CreateWallet(wltName, tc.opts)
				require.Equal(t, tc.expect.err, err, fmt.Sprintf("expected: %v; got : %v", tc.expect.err, err))
				if err != nil {
					return
				}

				//require.NoError(t, w.Validate())

				if !w.IsEncrypted() {
					require.Equal(t, tc.expect.seed, w.Seed())
					require.Equal(t, tc.expect.lastSeed, w.LastSeed())
				}
				require.Equal(t, tc.expect.xpub, w.XPub())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, tc.expect.entryNum, el)
				entries, err := w.GetEntries()
				require.NoError(t, err)

				for i, e := range entries {
					require.Equal(t, tc.expect.addrs[i].String(), e.Address.String())
				}
			})
		}
	}
}

func TestServiceScanAddresses(t *testing.T) {
	seed := "seed"
	// Generate adddresses from the seed
	var addrs []cipher.Address
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 10)
	for _, s := range seckeys {
		addrs = append(addrs, cipher.MustAddressFromSecKey(s))
	}

	bip44Seed := "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	bip44AddrStrs := []string{
		"9BSEAEE3XGtQ2X43BCT2XCYgheGLQQigEG",
		"29cnQPHuWHCRF26LEAb2gR83ywnF3F9HduW",
		"2ZUAv9MGSpDKR3dnKMUnrKqLenV22JXAxzP",
		"fwNVThqdzH7JMsStoLrTpkVsemesbdGftm",
		"eyr5KDLTnN6ZZeggeHqDcXnrwmNUi7sGk2",
		"Aee3J9qoFPLoUEJes6YVzdKHdeuvCrMZeJ",
		"29MZS8aiYUdEwcruwCPggVJG9YJLsm92FHa",
		"2Hbm3bwKiEwqNAMAzVJmz5hL1dNTfaA3ju7",
		"WCaSCwSZnVqtkYeiKryeHjR8LbzE3KbkzJ",
		"baRjCy1yHfishGdZi3bVaPaL7VJM7FZCSd",
	}
	bip44Addrs := make([]cipher.Address, len(bip44AddrStrs))
	for i, a := range bip44AddrStrs {
		bip44Addrs[i] = cipher.MustDecodeBase58Address(a)
	}

	xpub := "xpub6E5WPk37XdM79dy6oJ7iH6NkCvVzxmrCo4zMFFHSZMc5ymZYhReQFWaDcGNZeYYe1ahY2e3RcRZDHLHC98FfzPRfNRcU6ecURpS4RCQRP2w"
	xpubAddrStrs := []string{
		"2mhaS6SE2TPSmRRbJvngWQSNXCCVuTic5Zg",
		"2bq2itwDKteqigxVS9eYJv4Ww9SEfuyGcib",
		"B7eMXM6nLUqqzkFcosXR3HSVkQ6yUz53n4",
		"niAy17kBb8vB2pFey8eZnE92e6x9bFGLHp",
		"N8JbzcqWEPkn6CF3JdZDmEydzECyZ4NhRv",
		"TcyY3F4xHPCtFFkdDBTC93y684Fmxg2rPd",
		"kcVFbcrVqAVUHrirp7r3HYHUzrtdFuybez",
		"bzA7UeUmkuFWn3waGh3z1eQ5xV3TAZpgX2",
		"2K29ZX6vaqrTRZJbFKX7hzu646wL3pJfF6H",
		"VKpFCpN4yp46uYbffaeCg7XEvQd6pHNSkV",
	}

	xpubAddrs := make([]cipher.Address, len(xpubAddrStrs))
	for i, a := range xpubAddrStrs {
		xpubAddrs[i] = cipher.MustDecodeBase58Address(a)
	}

	type testCases []struct {
		name              string
		opts              wallet.Options
		scanN             uint64
		password          []byte
		walletAPIDisabled bool
		expectAddrs       []cipher.Address
		expectErr         error
	}

	generateTestCasesFunc := func(walletType, seed, xpub string, addrs []cipher.Address) testCases {
		password := []byte("pwd")
		var encrypt bool
		var tt testCases
		switch walletType {
		case wallet.WalletTypeXPub, wallet.WalletTypeBip44:
			password = []byte("")
		default:
			encrypt = true
			tt = testCases{
				{
					name: "scan 0 encrypted wrong password",
					opts: wallet.Options{
						Type:     walletType,
						Label:    "label",
						Seed:     seed,
						XPub:     xpub,
						Encrypt:  true,
						Password: []byte("pwd"),
					},
					scanN:       0,
					password:    []byte("incorrect password"),
					expectAddrs: nil,
					expectErr:   wallet.ErrInvalidPassword,
				},
			}
		}

		tt = append(tt, testCases{
			{
				name: "scan 0 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
				},
				scanN:       0,
				expectAddrs: []cipher.Address{},
				expectErr:   nil,
			},
			{
				name: "scan 0 encrypted",
				opts: wallet.Options{
					Type:     walletType,
					Label:    "label",
					Seed:     seed,
					XPub:     xpub,
					Encrypt:  encrypt,
					Password: password,
				},
				scanN:       0,
				password:    password,
				expectAddrs: []cipher.Address{},
				expectErr:   nil,
			},
			{
				name: "scan 1 get 0 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
					TF:    mockTxnsFinder{},
				},
				scanN:       1,
				expectAddrs: []cipher.Address{},
				expectErr:   nil,
			},
			{
				name: "scan 1 get 0 encrypted",
				opts: wallet.Options{
					Type:     walletType,
					Label:    "label",
					Seed:     seed,
					XPub:     xpub,
					Encrypt:  encrypt,
					Password: password,
					TF:       mockTxnsFinder{},
				},
				scanN:       1,
				password:    password,
				expectAddrs: []cipher.Address{},
				expectErr:   nil,
			},
			{
				name: "scan 1 get 1 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
					TF: mockTxnsFinder{
						addrs[1]: true,
					},
				},
				scanN:       1,
				expectAddrs: addrs[1:2],
				expectErr:   nil,
			},
			{
				name: "scan 1 get 1 encrypted",
				opts: wallet.Options{
					Type:     walletType,
					Label:    "label",
					Seed:     seed,
					XPub:     xpub,
					Encrypt:  encrypt,
					Password: password,
					TF: mockTxnsFinder{
						addrs[1]: true,
					},
				},
				scanN:       1,
				password:    password,
				expectAddrs: addrs[1:2],
				expectErr:   nil,
			},
			{
				name: "have 2 scan 1 get 1 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
					TF: mockTxnsFinder{
						addrs[1]: true,
						addrs[2]: true,
					},
				},
				scanN:       1,
				expectAddrs: addrs[1:2],
				expectErr:   nil,
			},
			{
				name: "have 2 scan 1 get 1 encrypted",
				opts: wallet.Options{
					Type:     walletType,
					Label:    "label",
					Seed:     seed,
					XPub:     xpub,
					Encrypt:  encrypt,
					Password: password,
					TF: mockTxnsFinder{
						addrs[1]: true,
						addrs[2]: true,
					},
				},
				scanN:       1,
				password:    password,
				expectAddrs: addrs[1:2],
				expectErr:   nil,
			},
			{
				name: "scan 2 get 1 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
					TF: mockTxnsFinder{
						addrs[1]: true,
					},
				},
				scanN:       2,
				expectAddrs: addrs[1:2],
				expectErr:   nil,
			},
			{
				name: "scan 2 get 1 encrypted",
				opts: wallet.Options{
					Type:     walletType,
					Label:    "label",
					Seed:     seed,
					XPub:     xpub,
					Encrypt:  encrypt,
					Password: password,
					TF: mockTxnsFinder{
						addrs[1]: true,
					},
				},
				scanN:       2,
				password:    password,
				expectAddrs: addrs[1:2],
				expectErr:   nil,
			},
			{
				name: "scan 2 get 2 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
					TF: mockTxnsFinder{
						addrs[1]: true,
						addrs[2]: true,
					},
				},
				scanN:       2,
				expectAddrs: addrs[1:3],
				expectErr:   nil,
			},
			{
				name: "scan 10 get 2 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
					TF: mockTxnsFinder{
						addrs[1]: true,
						addrs[2]: true,
					},
				},
				scanN:       10,
				expectAddrs: addrs[1:3],
				expectErr:   nil,
			},
			{
				name: "scan 10 get 5 unencrypted",
				opts: wallet.Options{
					Type:  walletType,
					Label: "label",
					Seed:  seed,
					XPub:  xpub,
					TF: mockTxnsFinder{
						addrs[1]: true,
						addrs[2]: true,
						addrs[5]: true,
					},
				},
				scanN:       10,
				expectAddrs: addrs[1:6],
				expectErr:   nil,
			},
		}...)
		return tt
	}

	testData := []struct {
		walletType string
		seed       string
		xpub       string
		addrs      []cipher.Address
	}{
		{
			walletType: wallet.WalletTypeDeterministic,
			addrs:      addrs[:],
			seed:       seed,
		},
		{
			walletType: wallet.WalletTypeBip44,
			addrs:      bip44Addrs[:],
			seed:       bip44Seed,
		},
		{
			walletType: wallet.WalletTypeXPub,
			addrs:      xpubAddrs[:],
			xpub:       xpub,
		},
	}

	for _, d := range testData {
		tt := generateTestCasesFunc(d.walletType, d.seed, d.xpub, d.addrs)
		for _, tc := range tt {
			for _, ct := range crypto.TypesInsecure() {
				name := fmt.Sprintf("crypto=%v type=%v %v", ct, d.walletType, tc.name)
				t.Run(name, func(t *testing.T) {
					dir := prepareWltDir()
					s, err := wallet.NewService(wallet.Config{
						WalletDir:       dir,
						CryptoType:      ct,
						EnableWalletAPI: !tc.walletAPIDisabled,
					})
					require.NoError(t, err)

					wltName := wallet.NewWalletFilename()
					w, err := s.CreateWallet(wltName, tc.opts)
					require.NoError(t, err)
					//require.NoError(t, w.Validate())

					addrs, err := s.ScanAddresses(w.Filename(), tc.password, tc.scanN, tc.opts.TF)
					require.Equal(t, tc.expectErr, err, fmt.Sprintf("%v", err))
					require.Equal(t, tc.expectAddrs, addrs)
				})
			}
		}
	}
}

func TestGetWalletSeed(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             wallet.Options
		id               string
		pwd              []byte
		disableWalletAPI bool
		enableSeedAPI    bool
		expectErr        error
	}{
		{
			name:    "wallet is not encrypted",
			wltName: "wallet.wlt",
			opts: wallet.Options{
				Seed:  "seed",
				Label: "label",
				Type:  wallet.WalletTypeDeterministic,
			},
			id:            "wallet.wlt",
			enableSeedAPI: true,
			expectErr:     wallet.ErrWalletNotEncrypted,
		},
		{
			name:    "wallet api disabled",
			wltName: "wallet.wlt",
			opts: wallet.Options{
				Seed:  "seed",
				Label: "label",
				Type:  wallet.WalletTypeDeterministic,
			},
			id:               "wallet.wlt",
			enableSeedAPI:    true,
			disableWalletAPI: true,
			expectErr:        wallet.ErrWalletAPIDisabled,
		},
		{
			name:    "ok",
			wltName: "wallet.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Label:    "label",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			enableSeedAPI: true,
			id:            "wallet.wlt",
			pwd:           []byte("pwd"),
		},
		{
			name:    "ok seed passphrase",
			wltName: "wallet.wlt",
			opts: wallet.Options{
				Seed:           bip39.MustNewDefaultMnemonic(),
				SeedPassphrase: "seed-passphrase",
				Label:          "label",
				Encrypt:        true,
				Password:       []byte("pwd"),
				Type:           wallet.WalletTypeBip44,
			},
			enableSeedAPI: true,
			id:            "wallet.wlt",
			pwd:           []byte("pwd"),
		},
		{
			name:    "wallet does not exist",
			wltName: "wallet.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Label:    "label",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			enableSeedAPI: true,
			pwd:           []byte("pwd"),
			id:            "none-exist.wlt",
			expectErr:     wallet.ErrWalletNotExist,
		},
		{
			name:    "disable seed api",
			wltName: "wallet.wlt",
			opts: wallet.Options{
				Seed:     "seed",
				Label:    "label",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     wallet.WalletTypeDeterministic,
			},
			pwd:           []byte("pwd"),
			id:            "wallet.wlt",
			enableSeedAPI: false,
			expectErr:     wallet.ErrSeedAPIDisabled,
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			t.Run(tc.name, func(t *testing.T) {
				dir := prepareWltDir()
				s, err := wallet.NewService(wallet.Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
					EnableSeedAPI:   tc.enableSeedAPI,
				})
				require.NoError(t, err)

				if tc.disableWalletAPI {
					_, _, err = s.GetWalletSeed("", tc.pwd)
					require.Equal(t, tc.expectErr, err)
					return
				}

				// Create a wallet
				_, err = s.CreateWallet(tc.wltName, tc.opts)
				require.NoError(t, err)

				seed, seedPassphrase, err := s.GetWalletSeed(tc.id, tc.pwd)
				require.Equal(t, tc.expectErr, err)
				if err != nil {
					return
				}

				require.Equal(t, tc.opts.Seed, seed)
				require.Equal(t, tc.opts.SeedPassphrase, seedPassphrase)
			})
		}
	}
}

func TestServiceView(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             wallet.Options
		viewWltName      string
		action           func(*testing.T) func(wallet.Wallet) error
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted collection wallet",
			wltName:     "test-view-collection-encrypted.wlt",
			viewWltName: "test-view-collection-encrypted.wlt",
			opts: wallet.Options{
				Type:     wallet.WalletTypeCollection,
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, unencrypted collection wallet",
			wltName:     "test-view-collection-unencrypted.wlt",
			viewWltName: "test-view-collection-unencrypted.wlt",
			opts: wallet.Options{
				Label: "foowlt",
				Type:  wallet.WalletTypeCollection,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
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
			name:        "ok, unencrypted xpub wallet",
			wltName:     "test-view-xpub-unencrypted.wlt",
			viewWltName: "test-view-xpub-unencrypted.wlt",
			opts: wallet.Options{
				Label: "foowlt",
				Type:  wallet.WalletTypeXPub,
				XPub:  "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v",
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					// xpub wallets don't have seeds
					require.Empty(t, w.Seed())
					require.Empty(t, w.LastSeed())

					require.Equal(t, "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v", w.XPub())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, encrypted bip44 wallet",
			wltName:     "test-view-bip44-encrypted.wlt",
			viewWltName: "test-view-bip44-encrypted.wlt",
			opts: wallet.Options{
				Type:     wallet.WalletTypeBip44,
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Seed:     "voyage say extend find sheriff surge priority merit ignore maple cash argue",
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, unencrypted bip44 wallet",
			wltName:     "test-view-bip44-unencrypted.wlt",
			viewWltName: "test-view-bip44-unencrypted.wlt",
			opts: wallet.Options{
				Label:          "foowlt",
				Type:           wallet.WalletTypeBip44,
				Seed:           "voyage say extend find sheriff surge priority merit ignore maple cash argue",
				SeedPassphrase: "foo",
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					require.Equal(t, "voyage say extend find sheriff surge priority merit ignore maple cash argue", w.Seed())
					require.Equal(t, "foo", w.SeedPassphrase())
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
			opts: wallet.Options{
				Type:     wallet.WalletTypeDeterministic,
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
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
			opts: wallet.Options{
				Type:  wallet.WalletTypeDeterministic,
				Seed:  "fooseed",
				Label: "foowlt",
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
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
			opts: wallet.Options{
				Type:  wallet.WalletTypeDeterministic,
				Seed:  "fooseed",
				Label: "foowlt",
			},
			err: wallet.ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-view-api-disabled.wlt",
			viewWltName: "test-view-api-disabled.wlt",
			opts: wallet.Options{
				Type:  wallet.WalletTypeDeterministic,
				Seed:  "fooseed",
				Label: "foowlt",
			},
			disableWalletAPI: true,
			err:              wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := wallet.NewService(wallet.Config{
				WalletDir:       dir,
				CryptoType:      crypto.CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			w, err := s.CreateWallet(tc.wltName, tc.opts)
			require.NoError(t, err)

			if w.Type() == wallet.WalletTypeCollection {
				// TODO: open this after collection type is back
				err := s.UpdateSecrets(w.Filename(), tc.opts.Password, func(w wallet.Wallet) error {
					p, s := cipher.GenerateKeyPair()
					return w.(*collection.Wallet).AddEntry(wallet.Entry{
						Public:  p,
						Secret:  s,
						Address: cipher.AddressFromPubKey(p),
					})
				})
				require.NoError(t, err)

				w, err = s.GetWallet(tc.wltName)
				require.NoError(t, err)
			}

			s.SetEnableWalletAPI(!tc.disableWalletAPI)

			var action func(wallet.Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.View(tc.viewWltName, action)

			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.SetEnableWalletAPI(true)

			// Check that the wallet is unmodified
			w2, err := s.GetWallet(tc.wltName)
			require.NoError(t, err)
			wd, err := w.Serialize()
			require.NoError(t, err)
			w2d, err := w2.Serialize()
			require.NoError(t, err)
			require.Equal(t, wd, w2d)
		})
	}
}

func TestServiceViewSecrets(t *testing.T) {
	mnemonicSeed := bip39.MustNewDefaultMnemonic()

	tt := []struct {
		name             string
		wltName          string
		opts             wallet.Options
		viewWltName      string
		action           func(*testing.T) func(wallet.Wallet) error
		password         []byte
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted wallet",
			wltName:     "test-view-secrets-encrypted.wlt",
			viewWltName: "test-view-secrets-encrypted.wlt",
			opts: wallet.Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeDeterministic,
			},
			password: []byte("pwd"),
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
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
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
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
			name:        "ok, encrypted wallet bip44 seed passphrase",
			wltName:     "test-view-secrets-encrypted-bip44.wlt",
			viewWltName: "test-view-secrets-encrypted-bip44.wlt",
			opts: wallet.Options{
				Seed:           mnemonicSeed,
				SeedPassphrase: "foobar",
				Encrypt:        true,
				Password:       []byte("pwd"),
				Label:          "foowlt",
				Type:           wallet.WalletTypeBip44,
			},
			password: []byte("pwd"),
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should be able to see sensitive data
					require.Equal(t, mnemonicSeed, w.Seed())
					require.Equal(t, "foobar", w.SeedPassphrase())
					require.Empty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, unencrypted wallet bip44",
			wltName:     "test-view-secrets-unencrypted-bip44.wlt",
			viewWltName: "test-view-secrets-unencrypted-bip44.wlt",
			opts: wallet.Options{
				Seed:  mnemonicSeed,
				Label: "foowlt",
				Type:  wallet.WalletTypeBip44,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, mnemonicSeed, w.Seed())
					require.Empty(t, w.SeedPassphrase())
					require.Empty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
		},

		{
			name:        "ok, unencrypted wallet bip44 seed passphrase",
			wltName:     "test-view-secrets-unencrypted-bip44.wlt",
			viewWltName: "test-view-secrets-unencrypted-bip44.wlt",
			opts: wallet.Options{
				Seed:           mnemonicSeed,
				SeedPassphrase: "foobar",
				Label:          "foowlt",
				Type:           wallet.WalletTypeBip44,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, mnemonicSeed, w.Seed())
					require.Equal(t, "foobar", w.SeedPassphrase())
					require.Empty(t, w.LastSeed())

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
			opts: wallet.Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeDeterministic,
			},
			err: wallet.ErrMissingPassword,
		},

		{
			name:        "encrypted wallet but password invalid",
			wltName:     "test-view-secrets-encrypted-wrong-password.wlt",
			viewWltName: "test-view-secrets-encrypted-wrong-password.wlt",
			opts: wallet.Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeDeterministic,
			},
			password: []byte("pwdpwd"),
			err:      wallet.ErrInvalidPassword,
		},

		{
			name:        "unencrypted wallet but password provided",
			wltName:     "test-view-secrets-unencrypted-password.wlt",
			viewWltName: "test-view-secrets-unencrypted-password.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			password: []byte("pwd"),
			err:      wallet.ErrWalletNotEncrypted,
		},

		{
			name:        "wallet doesn't exist",
			wltName:     "test-view-secrets-not-exist.wlt",
			viewWltName: "foo-test-view-secrets-not-exist.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			err: wallet.ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-view-secrets-api-disabled.wlt",
			viewWltName: "test-view-secrets-api-disabled.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			disableWalletAPI: true,
			err:              wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := wallet.NewService(wallet.Config{
				WalletDir:       dir,
				CryptoType:      crypto.CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			w, err := s.CreateWallet(tc.wltName, tc.opts)
			require.NoError(t, err)

			s.SetEnableWalletAPI(!tc.disableWalletAPI)

			var action func(wallet.Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.ViewSecrets(tc.viewWltName, tc.password, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			//s.wallet.Config.EnableWalletAPI = true
			s.SetEnableWalletAPI(true)

			// Check that the wallet is unmodified
			w2, err := s.GetWallet(tc.wltName)
			require.NoError(t, err)

			wd, err := w.Serialize()
			require.NoError(t, err)
			w2d, err := w2.Serialize()
			require.NoError(t, err)
			require.Equal(t, wd, w2d)
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	tt := []struct {
		name             string
		wltName          string
		opts             wallet.Options
		viewWltName      string
		action           func(*testing.T) func(wallet.Wallet) error
		checkWallet      func(*testing.T, wallet.Wallet)
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted wallet",
			wltName:     "test-update-encrypted.wlt",
			viewWltName: "test-update-encrypted.wlt",
			opts: wallet.Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeDeterministic,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should not be able to see sensitive data
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					// The wallet is encrypted so it cannot generate more addresses
					// TODO: bip44 wallet can generate address without decrypting
					_, err := w.GenerateAddresses(1)
					require.Equal(t, wallet.ErrWalletEncrypted, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 1, el)
				checkNoSensitiveData(t, w)
			},
		},
		{
			name:        "ok, bip44 encrypted wallet",
			wltName:     "test-update-encrypted-bip44.wlt",
			viewWltName: "test-update-encrypted-bip44.wlt",
			opts: wallet.Options{
				Seed:     "voyage say extend find sheriff surge priority merit ignore maple cash argue",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeBip44,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should not be able to see sensitive data
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					// The wallet is encrypted so it cannot generate more addresses
					// TODO: bip44 wallet can generate address without decrypting
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)
					//require.Equal(t, wallet.ErrWalletEncrypted, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 2, el)
				checkNoSensitiveData(t, w)
			},
		},
		{
			name:        "ok, encrypted collection wallet",
			wltName:     "test-update-encrypted-coll.wlt",
			viewWltName: "test-update-encrypted-coll.wlt",
			opts: wallet.Options{
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeCollection,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should not be able to see sensitive data
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					// The wallet is encrypted so it cannot generate more addresses
					_, err := w.GenerateAddresses(1)
					require.Equal(t, wallet.NewError(errors.New("A collection wallet does not implement GenerateAddresses")), err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 0, el)
				checkNoSensitiveData(t, w)
			},
		},
		{
			name:        "ok, unencrypted xpub wallet",
			wltName:     "test-update-unencrypted-xpub.wlt",
			viewWltName: "test-update-unencrypted-xpub.wlt",
			opts: wallet.Options{
				XPub:  "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v",
				Label: "foowlt",
				Type:  wallet.WalletTypeXPub,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should not be able to see sensitive data
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					// The wallet is encrypted so it cannot generate more addresses
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 2, el)
				checkNoSensitiveData(t, w)
			},
		},
		{
			name:        "ok, unencrypted wallet",
			wltName:     "test-update-unencrypted.wlt",
			viewWltName: "test-update-unencrypted.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.Seed())
					require.NotEmpty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 1, el)
				entries, err := w.GetEntries()
				require.NoError(t, err)

				require.NotEmpty(t, entries[0].Secret)
			},
		},

		{
			name:        "wallet doesn't exist",
			wltName:     "test-update-not-exist.wlt",
			viewWltName: "foo-test-update-not-exist.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			err: wallet.ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-update-api-disabled.wlt",
			viewWltName: "test-update-api-disabled.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			disableWalletAPI: true,
			err:              wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := wallet.NewService(wallet.Config{
				WalletDir:       dir,
				CryptoType:      crypto.CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			_, err = s.CreateWallet(tc.wltName, tc.opts)
			require.NoError(t, err)

			s.SetEnableWalletAPI(!tc.disableWalletAPI)

			var action func(wallet.Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.Update(tc.viewWltName, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.SetEnableWalletAPI(true)

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
		opts             wallet.Options
		viewWltName      string
		action           func(*testing.T) func(wallet.Wallet) error
		checkWallet      func(*testing.T, wallet.Wallet)
		password         []byte
		disableWalletAPI bool
		err              error
	}{
		{
			name:        "ok, encrypted wallet",
			wltName:     "test-update-secrets-encrypted.wlt",
			viewWltName: "test-update-secrets-encrypted.wlt",
			opts: wallet.Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeDeterministic,
			},
			password: []byte("pwd"),
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
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
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 2, el)
				checkNoSensitiveData(t, w)
			},
		},

		{
			name:        "ok, unencrypted wallet",
			wltName:     "test-update-secrets-unencrypted.wlt",
			viewWltName: "test-update-secrets-unencrypted.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
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
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 2, el)
				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.NotEmpty(t, entries[1].Secret)
			},
		},

		{
			name:        "ok, encrypted bip44 wallet",
			wltName:     "test-update-secrets-bip44-encrypted.wlt",
			viewWltName: "test-update-secrets-bip44-encrypted.wlt",
			opts: wallet.Options{
				Seed:     "voyage say extend find sheriff surge priority merit ignore maple cash argue",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowltbip44",
				Type:     wallet.WalletTypeBip44,
			},
			password: []byte("pwd"),
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowltbip44", w.Label())

					// Should be able to see sensitive data
					require.Equal(t, "voyage say extend find sheriff surge priority merit ignore maple cash argue", w.Seed())
					require.Empty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltbip44foo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 2, el)
				checkNoSensitiveData(t, w)
			},
		},

		{
			name:        "ok, unencrypted bip44 wallet",
			wltName:     "test-update-secrets-bip44-unencrypted.wlt",
			viewWltName: "test-update-secrets-bip44-unencrypted.wlt",
			opts: wallet.Options{
				Seed:  "voyage say extend find sheriff surge priority merit ignore maple cash argue",
				Label: "foowltbip44",
				Type:  wallet.WalletTypeBip44,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowltbip44", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "voyage say extend find sheriff surge priority merit ignore maple cash argue", w.Seed())
					require.Empty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltbip44foo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 2, el)
				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.NotEmpty(t, entries[1].Secret)
			},
		},
		{
			name:        "ok, xpub wallet",
			wltName:     "test-update-secrets-xpub.wlt",
			viewWltName: "test-update-secrets-xpub.wlt",
			opts: wallet.Options{
				XPub:  "xpub6EFYYRQeAbWLdWQYbtQv8HnemieKNmYUE23RmwphgtMLjz4UaStKADSKNoSSXM5FDcq4gZec2q6n7kdNWfuMdScxK1cXm8tR37kaitHtvuJ",
				Label: "foowltxpub",
				Type:  wallet.WalletTypeXPub,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowltxpub", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "xpub6EFYYRQeAbWLdWQYbtQv8HnemieKNmYUE23RmwphgtMLjz4UaStKADSKNoSSXM5FDcq4gZec2q6n7kdNWfuMdScxK1cXm8tR37kaitHtvuJ", w.XPub())
					require.Empty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltxpubfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 2, el)
				entries, err := w.GetEntries()
				require.NoError(t, err)
				require.Equal(t, entries[1].Secret, cipher.SecKey{})
			},
		},

		{
			name:        "ok, collection wallet",
			wltName:     "test-update-secrets-collection.wlt",
			viewWltName: "test-update-secrets-collection.wlt",
			opts: wallet.Options{
				Label: "foowltcollection",
				Type:  wallet.WalletTypeCollection,
			},
			action: func(t *testing.T) func(wallet.Wallet) error {
				return func(w wallet.Wallet) error {
					require.Equal(t, "foowltcollection", w.Label())

					require.Empty(t, w.Seed())
					require.Empty(t, w.LastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.SetLabel(w.Label() + "foo")

					seed := []byte("seed")
					_, keys, err := cipher.GenerateDeterministicKeyPairsSeed(seed, 5)
					require.NoError(t, err)
					for i := range keys {
						pk, err := cipher.PubKeyFromSecKey(keys[i])
						require.NoError(t, err)
						addr := cipher.AddressFromPubKey(pk)
						err = w.(*collection.Wallet).AddEntry(wallet.Entry{
							Address: addr,
							Public:  pk,
							Secret:  keys[i],
						})
						require.NoError(t, err)
					}

					return nil
				}
			},
			checkWallet: func(t *testing.T, w wallet.Wallet) {
				require.Equal(t, "foowltcollectionfoo", w.Label())
				el, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 5, el)
				entries, err := w.GetEntries()
				require.NoError(t, err)

				addrs := []string{
					"2EVNa4CK9SKosT4j1GEn8SuuUUEAXaHAMbM",
					"68enNSvabNYLf97xhb19vmLrrG3yqXPmkV",
					"CHAJD8BMpnZ14iv34VWs23BzkBbBNcb5sH",
					"2ip6roWzqAxwLjNdy36HUAqjCRkh6scWbLD",
					"jT1aTYGN4XSUC6bLTCmaEbEeidqjPQCLnK",
				}
				for i, e := range entries {
					require.Equal(t, e.Address.String(), addrs[i])
				}
			},
		},

		{
			name:        "encrypted wallet but password not provided",
			wltName:     "test-update-secrets-encrypted-no-password.wlt",
			viewWltName: "test-update-secrets-encrypted-no-password.wlt",
			opts: wallet.Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeDeterministic,
			},
			err: wallet.ErrMissingPassword,
		},

		{
			name:        "encrypted wallet but password invalid",
			wltName:     "test-update-secrets-encrypted-wrong-password.wlt",
			viewWltName: "test-update-secrets-encrypted-wrong-password.wlt",
			opts: wallet.Options{
				Seed:     "fooseed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Label:    "foowlt",
				Type:     wallet.WalletTypeDeterministic,
			},
			password: []byte("pwdpwd"),
			err:      wallet.ErrInvalidPassword,
		},

		{
			name:        "unencrypted wallet but password provided",
			wltName:     "test-update-secrets-unencrypted-password.wlt",
			viewWltName: "test-update-secrets-unencrypted-password.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			password: []byte("pwd"),
			err:      wallet.ErrWalletNotEncrypted,
		},

		{
			name:        "wallet doesn't exist",
			wltName:     "test-update-secrets-not-exist.wlt",
			viewWltName: "foo-test-update-secrets-not-exist.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			err: wallet.ErrWalletNotExist,
		},

		{
			name:        "api disabled",
			wltName:     "test-update-secrets-api-disabled.wlt",
			viewWltName: "test-update-secrets-api-disabled.wlt",
			opts: wallet.Options{
				Seed:  "fooseed",
				Label: "foowlt",
				Type:  wallet.WalletTypeDeterministic,
			},
			disableWalletAPI: true,
			err:              wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			s, err := wallet.NewService(wallet.Config{
				WalletDir:       dir,
				CryptoType:      crypto.CryptoTypeSha256Xor,
				EnableWalletAPI: true,
			})
			require.NoError(t, err)

			_, err = s.CreateWallet(tc.wltName, tc.opts)
			require.NoError(t, err)

			s.SetEnableWalletAPI(!tc.disableWalletAPI)

			var action func(wallet.Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.UpdateSecrets(tc.viewWltName, tc.password, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.SetEnableWalletAPI(true)

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

func checkNoSensitiveData(t *testing.T, w wallet.Wallet) {
	require.Empty(t, w.Seed())
	require.Empty(t, w.LastSeed())
	require.Empty(t, w.SeedPassphrase())
	entries, err := w.GetEntries()
	require.NoError(t, err)
	for _, e := range entries {
		require.True(t, e.Secret.Null())
	}
}
