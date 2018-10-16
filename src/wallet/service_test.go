package wallet

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/shopspring/decimal"
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

			require.Equal(t, 0, len(s.wallets))

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
				require.Equal(t, err, ErrSeedUsed)

				// check if the dup wallet is created
				_, ok := s.wallets[dupWlt]
				require.False(t, ok)

				testutil.RequireFileNotExists(t, filepath.Join(dir, dupWlt))
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
				wltName := NewWalletFilename()

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
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed(seed, 10)
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
		{
			name: "encrypted=false password provided",
			opts: Options{
				Label: "label",
				Seed:  string(seed),
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
					var emptyW *Wallet
					require.Equal(t, w, emptyW)
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

				// Creates a wallet
				w, err := s.CreateWallet("t.wlt", Options{
					Label: "label",
					Seed:  "seed",
				}, nil)
				require.NoError(t, err)

				var wallets []*Wallet
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

func TestServiceCreateAndSignTransaction(t *testing.T) {
	headTime := time.Now().UTC().Unix()
	seed := []byte("seed")

	// Generate first keys
	_, secKeys := cipher.MustGenerateDeterministicKeyPairsSeed(seed, 1)
	secKey := secKeys[0]
	addr := cipher.MustAddressFromSecKey(secKey)

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
		uxout := makeUxOut(t, secKey, 2e6, 0)
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
			coins:    2e6,
			dest:     addrs[0],
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
			coins:    2e6,
			dest:     addrs[0],
		},
		{
			name: "encrypted=false has change=yes",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			coins:    1e6,
			dest:     addrs[0],
		},
		{
			name: "encrypted=false spend zero",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			dest:     addrs[0],
			err:      ErrZeroSpend,
		},
		{
			name: "encrypted=false spend fractional coins",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			coins:    1e3,
			dest:     addrs[0],
		},
		{
			name: "encrypted=false not enough confirmed coins",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxouts[:],
			coins:    100e6,
			dest:     addrs[0],
			err:      ErrInsufficientBalance,
		},
		{
			name: "encrypted=false no coin hours in inputs",
			opts: Options{
				Seed: string(seed),
			},
			unspents: uxoutsNoHours[:],
			coins:    1e6,
			dest:     addrsNoHours[0],
			err:      fee.ErrTxnNoFee,
		},
		{
			name: "disable wallet api=true",
			opts: Options{
				Seed:  string(seed),
				Label: "label",
			},
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},

		{
			name: "encrypted=false password provided",
			opts: Options{
				Seed: string(seed),
			},
			pwd: []byte("foo"),
			err: ErrWalletNotEncrypted,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			t.Run(name, func(t *testing.T) {
				addrUxOuts := coin.AddressUxOuts{
					addr: tc.unspents,
				}

				unspents := make(map[cipher.SHA256]coin.UxOut)

				for _, uxs := range addrUxOuts {
					for _, ux := range uxs {
						unspents[ux.Hash()] = ux
					}
				}

				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: !tc.disableWalletAPI,
				})
				require.NoError(t, err)

				if tc.disableWalletAPI {
					_, err = s.CreateAndSignTransaction("", tc.pwd, addrUxOuts, uint64(headTime), tc.coins, tc.dest)
					require.Equal(t, tc.err, err)
					return
				}

				wltName := NewWalletFilename()

				w, err := s.CreateWallet(wltName, tc.opts, nil)
				require.NoError(t, err)

				tx, err := s.CreateAndSignTransaction(w.Filename(), tc.pwd, addrUxOuts, uint64(headTime), tc.coins, tc.dest)

				if tc.err != nil {
					require.Error(t, err)
					require.Equal(t, tc.err, err, err.Error())
					return
				}

				require.NoError(t, err)

				// check the IN of tx
				for _, inUxid := range tx.In {
					_, ok := unspents[inUxid]
					require.True(t, ok)
				}

				err = tx.Verify()
				require.NoError(t, err)
			})
		}
	}
}

func TestServiceCreateAndSignTransactionAdvanced(t *testing.T) {
	headTime := uint64(time.Now().UTC().Unix())
	seed := []byte("seed")

	// Generate first keys
	_, secKeys := cipher.MustGenerateDeterministicKeyPairsSeed(seed, 11)
	secKey := secKeys[0]
	addr := cipher.MustAddressFromSecKey(secKey)

	var extraWalletAddrs []cipher.Address
	for _, s := range secKeys[1:] {
		extraWalletAddrs = append(extraWalletAddrs, cipher.MustAddressFromSecKey(s))
	}

	// Create unspent outputs
	var uxouts []coin.UxOut
	var originalUxouts []coin.UxOut
	addrs := []cipher.Address{}
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey, 2e6, uint64(100+i))
		uxout.Head.Time = headTime
		uxouts = append(uxouts, uxout)
		originalUxouts = append(originalUxouts, uxout)

		a := testutil.MakeAddress()
		addrs = append(addrs, a)
	}

	// shuffle the uxouts to test that the uxout sorting during spend selection is working
	rand.Shuffle(len(uxouts), func(i, j int) {
		uxouts[i], uxouts[j] = uxouts[j], uxouts[i]
	})

	// Create extra unspent outputs. These have the same value as uxouts, but are spendable by
	// keys held in extraWalletAddrs
	extraUxouts := make([][]coin.UxOut, len(extraWalletAddrs))
	for j := range extraWalletAddrs {
		s := secKeys[j+1]

		var uxouts []coin.UxOut
		for i := 0; i < 10; i++ {
			uxout := makeUxOut(t, s, 2e6, uint64(100+i))
			uxout.Head.Time = headTime
			uxouts = append(uxouts, uxout)
		}

		extraUxouts[j] = uxouts
	}

	// Create unspent outputs with no hours
	var uxoutsNoHours []coin.UxOut
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey, 2e6, 0)
		uxout.Head.Time = headTime
		uxoutsNoHours = append(uxoutsNoHours, uxout)
	}

	// shuffle the uxouts to test that the uxout sorting during spend selection is working
	rand.Shuffle(len(uxoutsNoHours), func(i, j int) {
		uxoutsNoHours[i], uxoutsNoHours[j] = uxoutsNoHours[j], uxoutsNoHours[i]
	})

	changeAddress := testutil.MakeAddress()

	validParams := CreateTransactionParams{
		HoursSelection: HoursSelection{
			Type: HoursSelectionTypeManual,
		},
		ChangeAddress: &changeAddress,
		To: []coin.TransactionOutput{
			{
				Address: addrs[0],
				Hours:   10,
				Coins:   1e6,
			},
		},
	}

	validParamsWithPassword := validParams
	validParamsWithPassword.Wallet.Password = []byte("password")

	newShareFactor := func(a string) *decimal.Decimal {
		d, err := decimal.NewFromString(a)
		require.NoError(t, err)
		return &d
	}

	firstAddress := func(uxa coin.UxArray) cipher.Address {
		require.NotEmpty(t, uxa)

		addresses := make([]cipher.Address, len(uxa))
		for i, a := range uxa {
			addresses[i] = a.Body.Address
		}

		sort.Slice(addresses, func(i, j int) bool {
			x := addresses[i].Bytes()
			y := addresses[j].Bytes()
			return bytes.Compare(x, y) < 0
		})

		return addresses[0]
	}

	cases := []struct {
		name             string
		err              error
		params           CreateTransactionParams
		opts             Options
		unspents         []coin.UxOut
		addressUnspents  coin.AddressUxOuts
		chosenUnspents   []coin.UxOut
		headTime         uint64
		disableWalletAPI bool
		walletNotExist   bool
		changeOutput     *coin.TransactionOutput
		toExpectedHours  []uint64
	}{
		{
			name:             "wallet api disabled",
			disableWalletAPI: true,
			err:              ErrWalletAPIDisabled,
		},

		{
			name:   "params invalid",
			params: CreateTransactionParams{},
			err:    NewError(errors.New("To is required")),
		},

		{
			name:           "wallet doesn't exist",
			params:         validParams,
			walletNotExist: true,
			err:            ErrWalletNotExist,
		},

		{
			name:   "wallet encrypted and password not provided",
			params: validParams,
			opts: Options{
				Encrypt: true,
			},
			err: ErrMissingPassword,
		},

		{
			name:   "wallet not encrypted and password provided",
			params: validParamsWithPassword,
			opts: Options{
				Encrypt: false,
			},
			err: ErrWalletNotEncrypted,
		},

		{
			name: "overflowing coin hours in params",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   math.MaxUint64,
						Coins:   1e6,
					},
					{
						Address: addrs[1],
						Hours:   1,
						Coins:   1e6,
					},
				},
			},
			err: NewError(errors.New("total output hours error: uint64 addition overflow")),
		},

		{
			name: "overflowing coins in params",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   math.MaxUint64,
					},
					{
						Address: addrs[1],
						Hours:   1,
						Coins:   1,
					},
				},
			},
			err: NewError(errors.New("total output coins error: uint64 addition overflow")),
		},

		{
			name: "no unspents",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   1e6,
					},
				},
			},
			err: ErrNoUnspents,
		},

		{
			name: "insufficient coins",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   100e6,
					},
				},
			},
			unspents: uxouts[:1],
			err:      ErrInsufficientBalance,
		},

		{
			name: "insufficient hours",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   100e6,
						Coins:   1e6,
					},
				},
			},
			unspents: uxouts[:1],
			err:      ErrInsufficientHours,
		},

		{
			name: "insufficient coins for specified uxouts",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				Wallet: CreateTransactionWalletParams{
					UxOuts: []cipher.SHA256{
						extraUxouts[0][0].Hash(),
					},
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   1,
						Coins:   3e6,
					},
				},
			},
			addressUnspents: coin.AddressUxOuts{
				extraWalletAddrs[0]: []coin.UxOut{extraUxouts[0][0]},
			},
			err: ErrInsufficientBalance,
		},

		{
			name: "insufficient hours for specified uxouts",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				Wallet: CreateTransactionWalletParams{
					UxOuts: []cipher.SHA256{
						extraUxouts[0][0].Hash(),
					},
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   200,
						Coins:   1e6,
					},
				},
			},
			addressUnspents: coin.AddressUxOuts{
				extraWalletAddrs[0]: []coin.UxOut{extraUxouts[0][0]},
			},
			err: ErrInsufficientHours,
		},

		{
			name: "manual, 1 output, no change",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   2e6,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0]},
		},

		{
			name: "manual, 1 output, no change, unknown address in auxs",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				Wallet: CreateTransactionWalletParams{},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   2e6,
					},
				},
			},
			addressUnspents: coin.AddressUxOuts{
				testutil.MakeAddress(): []coin.UxOut{extraUxouts[0][0]},
			},
			err: ErrUnknownAddress,
		},

		{
			name: "manual, 1 output, change",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   2e6 + 1,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   50,
				Coins:   2e6 - 1,
			},
		},

		{
			name: "manual, 1 output, change, unspecified change address",
			params: CreateTransactionParams{
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   2e6 + 1,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput: &coin.TransactionOutput{
				Address: firstAddress([]coin.UxOut{originalUxouts[0], originalUxouts[1]}),
				Hours:   50,
				Coins:   2e6 - 1,
			},
		},

		{
			// there are leftover coin hours and an additional input is added
			// to force change to save the leftover coin hours
			name: "manual, 1 output, forced change",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   0,
						Coins:   2e6 * 2,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   151,
				Coins:   2e6,
			},
		},

		{
			// there are leftover coin hours and no coins change,
			// but there are no more unspents to use to force a change output
			name: "manual, 1 output, forced change rejected no more unspents",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   80,
						Coins:   2e6 * 2,
					},
				},
			},
			unspents:       originalUxouts[:2],
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput:   nil,
		},

		{
			// there are leftover coin hours and no coins change,
			// but the hours cost of saving them with an additional input is less than is leftover
			name: "manual, 1 output, forced change rejected",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   80,
						Coins:   2e6 * 2,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput:   nil,
		},

		{
			name: "manual, multiple outputs",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6 + 1,
					},
					{
						Address: addrs[1],
						Hours:   70,
						Coins:   2e6,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2], originalUxouts[3]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   33,
				Coins:   4e6 - 1,
			},
		},

		{
			name: "manual, multiple outputs, varied addressUnspents",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				Wallet: CreateTransactionWalletParams{},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6 + 1,
					},
					{
						Address: addrs[1],
						Hours:   70,
						Coins:   2e6,
					},
				},
			},
			addressUnspents: coin.AddressUxOuts{
				extraWalletAddrs[0]: []coin.UxOut{extraUxouts[0][0]},
				extraWalletAddrs[3]: []coin.UxOut{extraUxouts[3][1], extraUxouts[3][2]},
				extraWalletAddrs[5]: []coin.UxOut{extraUxouts[5][6]},
			},
			chosenUnspents: []coin.UxOut{extraUxouts[0][0], extraUxouts[3][1], extraUxouts[3][2], extraUxouts[5][6]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   34,
				Coins:   4e6 - 1,
			},
		},

		{
			name: "manual, multiple uxouts, varied addressUnspents, wallet outputs specified",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				Wallet: CreateTransactionWalletParams{
					UxOuts: []cipher.SHA256{
						extraUxouts[0][0].Hash(),
						extraUxouts[3][1].Hash(),
						extraUxouts[3][2].Hash(),
						extraUxouts[5][6].Hash(),

						// this extra output is not necessary to satisfy the spend,
						// it is included to test that when UxOuts are specified,
						// only a subset is used
						extraUxouts[0][8].Hash(),
					},
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6 + 1,
					},
					{
						Address: addrs[1],
						Hours:   70,
						Coins:   2e6,
					},
				},
			},
			addressUnspents: coin.AddressUxOuts{
				extraWalletAddrs[0]: []coin.UxOut{extraUxouts[0][0], extraUxouts[0][8]},
				extraWalletAddrs[3]: []coin.UxOut{extraUxouts[3][1], extraUxouts[3][2]},
				extraWalletAddrs[5]: []coin.UxOut{extraUxouts[5][6]},
			},
			chosenUnspents: []coin.UxOut{
				extraUxouts[0][0],
				extraUxouts[3][1],
				extraUxouts[3][2],
				extraUxouts[5][6],
			},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   34,
				Coins:   4e6 - 1,
			},
		},

		{
			name: "auto, multiple outputs, share factor 0.5",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("0.5"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   76,
				Coins:   2e6 - (1e6 + 1e3),
			},
			toExpectedHours: []uint64{15, 30, 29, 1},
		},

		{
			name: "auto, multiple outputs, share factor 0.5, switch to 1.0 because no change could be made",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("0.5"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e6 - 1e3,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:        []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			chosenUnspents:  []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			toExpectedHours: []uint64{25, 50, 50, 25, 1},
		},

		{
			name: "encrypted, auto, multiple outputs, share factor 0.5",
			opts: Options{
				Encrypt:  true,
				Password: []byte("password"),
			},
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("0.5"),
				},
				Wallet: CreateTransactionWalletParams{
					Password: []byte("password"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   76,
				Coins:   2e6 - (1e6 + 1e3),
			},
			toExpectedHours: []uint64{15, 30, 29, 1},
		},

		{
			name: "auto, multiple outputs, share factor 0",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("0"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   151,
				Coins:   2e6 - (1e6 + 1e3),
			},
			toExpectedHours: []uint64{0, 0, 0, 0},
		},

		{
			name: "auto, multiple outputs, share factor 1",
			params: CreateTransactionParams{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("1"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   0,
				Coins:   2e6 - (1e6 + 1e3),
			},
			toExpectedHours: []uint64{30, 60, 60, 1},
		},

		{
			name:     "no coin hours in inputs",
			unspents: uxoutsNoHours[:],
			params: CreateTransactionParams{
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				ChangeAddress: &changeAddress,
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   1e6,
					},
				},
			},
			err: fee.ErrTxnNoFee,
		},
	}

	var cryptoTypes []CryptoType
	for ct := range cryptoTable {
		cryptoTypes = append(cryptoTypes, ct)
	}

	for _, tc := range cases {
		cts := cryptoTypes
		if !tc.opts.Encrypt {
			cts = cts[:1]
		}

		for _, ct := range cts {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			fmt.Println(name)
			t.Run(name, func(t *testing.T) {
				if tc.headTime == 0 {
					tc.headTime = headTime
				}

				addrUxOuts := coin.AddressUxOuts{
					addr: tc.unspents,
				}

				if tc.addressUnspents != nil {
					addrUxOuts = tc.addressUnspents
				}

				unspents := make(map[cipher.SHA256]coin.UxOut)
				for _, uxs := range addrUxOuts {
					for _, ux := range uxs {
						unspents[ux.Hash()] = ux
					}
				}

				if tc.opts.Seed == "" {
					tc.opts.Seed = string(seed)
				}

				dir := prepareWltDir()
				s, err := NewService(Config{
					WalletDir:       dir,
					CryptoType:      ct,
					EnableWalletAPI: true,
				})
				require.NoError(t, err)

				if tc.walletNotExist {
					tc.params.Wallet.ID = "foo.wlt"
				} else {
					wltName := NewWalletFilename()
					opts := tc.opts
					if opts.Encrypt && len(opts.Password) == 0 {
						opts.Password = []byte("password")
					}
					w, err := s.CreateWallet(wltName, opts, nil)
					require.NoError(t, err)

					if !w.IsEncrypted() {
						_, err := s.NewAddresses(w.Filename(), nil, 10)
						require.NoError(t, err)

						w, err = s.GetWallet(wltName)
						require.NoError(t, err)

						require.Equal(t, 11, len(w.Entries))
						require.Equal(t, w.Entries[0].Address, addr)
						for i, e := range w.Entries[1:] {
							require.Equal(t, e.Address, extraWalletAddrs[i])
						}
					}

					tc.params.Wallet.ID = wltName
				}

				s.enableWalletAPI = !tc.disableWalletAPI

				txn, inputs, err := s.CreateAndSignTransactionAdvanced(tc.params, addrUxOuts, tc.headTime)
				if tc.err != nil {
					require.Equal(t, tc.err, err)
					return
				}

				require.NoError(t, err)

				err = txn.Verify()
				require.NoError(t, err)

				require.Equal(t, len(inputs), len(txn.In))

				// Checks duplicate inputs in array
				inputsMap := make(map[cipher.SHA256]struct{})
				for _, i := range inputs {
					_, ok := inputsMap[i.Hash]
					require.False(t, ok)
					inputsMap[i.Hash] = struct{}{}
				}

				for i, inUxid := range txn.In {
					_, ok := unspents[inUxid]
					require.True(t, ok)

					require.Equal(t, inUxid, inputs[i].Hash)
				}

				// Compare the transaction inputs
				chosenUnspents := make([]coin.UxOut, len(tc.chosenUnspents))
				chosenUnspentHashes := make([]cipher.SHA256, len(tc.chosenUnspents))
				for i, u := range tc.chosenUnspents {
					chosenUnspents[i] = u
					chosenUnspentHashes[i] = u.Hash()
				}
				sort.Slice(chosenUnspentHashes, func(i, j int) bool {
					return bytes.Compare(chosenUnspentHashes[i][:], chosenUnspentHashes[j][:]) < 0
				})
				sort.Slice(chosenUnspents, func(i, j int) bool {
					h1 := chosenUnspents[i].Hash()
					h2 := chosenUnspents[j].Hash()
					return bytes.Compare(h1[:], h2[:]) < 0
				})

				sortedTxnIn := make([]cipher.SHA256, len(txn.In))
				copy(sortedTxnIn[:], txn.In[:])

				sort.Slice(sortedTxnIn, func(i, j int) bool {
					return bytes.Compare(sortedTxnIn[i][:], sortedTxnIn[j][:]) < 0
				})

				require.Equal(t, chosenUnspentHashes, sortedTxnIn)

				sort.Slice(inputs, func(i, j int) bool {
					h1 := inputs[i].Hash
					h2 := inputs[j].Hash
					return bytes.Compare(h1[:], h2[:]) < 0
				})

				chosenUnspentsUxBalances := make([]UxBalance, len(chosenUnspents))
				for i, o := range chosenUnspents {
					b, err := NewUxBalance(tc.headTime, o)
					require.NoError(t, err)
					chosenUnspentsUxBalances[i] = b
				}

				require.Equal(t, chosenUnspentsUxBalances, inputs)

				// Assign expected hours for comparison
				var to []coin.TransactionOutput
				to = append(to, tc.params.To...)

				if len(tc.toExpectedHours) != 0 {
					require.Equal(t, len(tc.toExpectedHours), len(to))
					for i, h := range tc.toExpectedHours {
						to[i].Hours = h
					}
				}

				// Add the change output if specified
				if tc.changeOutput != nil {
					to = append(to, *tc.changeOutput)
				}

				// Compare transaction outputs
				require.Equal(t, to, txn.Out)
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
				decWlt, err := encWlt.Unlock(tc.pwd)
				require.NoError(t, err)
				require.Equal(t, w.seed(), decWlt.seed())
				require.Equal(t, w.lastSeed(), decWlt.lastSeed())

				// Check if the wallet file does exist
				path := filepath.Join(dir, w.Filename())
				testutil.RequireFileExists(t, path)

				// Check if the backup wallet file, which should not exist
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
					lsd, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(wlt.seed()), entryNum)
					require.NoError(t, err)
					require.Equal(t, hex.EncodeToString(lsd), wlt.lastSeed())

					// Checks the entries
					for i := range seckeys {
						a := cipher.MustAddressFromSecKey(seckeys[i])
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

	bg := make(mockBalanceGetter, 20)
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
			balGetter: bg,
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
			balGetter: bg,
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
			balGetter: bg,
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
			balGetter: bg,
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
			balGetter: mockBalanceGetter{
				addrs[5]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[5]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[8]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[4+1]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[10]:  BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[5]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[6]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[7]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[2]:  BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[7]:  BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[12]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[2]:  BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[7]:  BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[13]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[2]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[8]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
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
			balGetter: mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
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

				wltName := NewWalletFilename()
				w, err := s.CreateWallet(wltName, tc.opts, tc.balGetter)
				require.Equal(t, tc.expect.err, err)
				if err != nil {
					return
				}

				require.NoError(t, w.Validate())
				require.Equal(t, tc.expect.entryNum, len(w.Entries))
				for i := range w.Entries {
					require.Equal(t, addrs[i].String(), w.Entries[i].Address.String())
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
		action           func(*testing.T) func(*Wallet) error
		disableWalletAPI bool
		err              error
	}{
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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.setLabel(w.Label() + "foo")

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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())
					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.seed())
					require.NotEmpty(t, w.lastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.setLabel(w.Label() + "foo")

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

			s.enableWalletAPI = !tc.disableWalletAPI

			var action func(*Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.View(tc.viewWltName, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.enableWalletAPI = true

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
		action           func(*testing.T) func(*Wallet) error
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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should be able to see sensitive data
					require.Equal(t, "fooseed", w.seed())
					require.NotEmpty(t, w.lastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.setLabel(w.Label() + "foo")

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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.seed())
					require.NotEmpty(t, w.lastSeed())

					// Modify the wallet pointer in order to check that this references a clone and not the original
					w.setLabel(w.Label() + "foo")

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

			s.enableWalletAPI = !tc.disableWalletAPI

			var action func(*Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.ViewSecrets(tc.viewWltName, tc.password, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.enableWalletAPI = true

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
		action           func(*testing.T) func(*Wallet) error
		checkWallet      func(*testing.T, *Wallet)
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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should not be able to see sensitive data
					checkNoSensitiveData(t, w)

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.setLabel(w.Label() + "foo")

					// The wallet is encrypted so it cannot generate more addresses
					_, err := w.GenerateAddresses(1)
					require.Equal(t, ErrWalletEncrypted, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w *Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Len(t, w.Entries, 1)
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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.seed())
					require.NotEmpty(t, w.lastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.setLabel(w.Label() + "foo")

					return nil
				}
			},
			checkWallet: func(t *testing.T, w *Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Len(t, w.Entries, 1)
				require.NotEmpty(t, w.Entries[0].Secret)
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

			s.enableWalletAPI = !tc.disableWalletAPI

			var action func(*Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.Update(tc.viewWltName, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.enableWalletAPI = true

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
		action           func(*testing.T) func(*Wallet) error
		checkWallet      func(*testing.T, *Wallet)
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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Should be able to see sensitive data
					require.Equal(t, "fooseed", w.seed())
					require.NotEmpty(t, w.lastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.setLabel(w.Label() + "foo")
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w *Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Len(t, w.Entries, 2)
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
			action: func(t *testing.T) func(*Wallet) error {
				return func(w *Wallet) error {
					require.Equal(t, "foowlt", w.Label())

					// Seed is visible because its not encrypted
					require.Equal(t, "fooseed", w.seed())
					require.NotEmpty(t, w.lastSeed())

					// Modify the wallet pointer in order to check that the wallet gets saved
					w.setLabel(w.Label() + "foo")
					_, err := w.GenerateAddresses(1)
					require.NoError(t, err)

					return nil
				}
			},
			checkWallet: func(t *testing.T, w *Wallet) {
				require.Equal(t, "foowltfoo", w.Label())
				require.Len(t, w.Entries, 2)
				require.NotEmpty(t, w.Entries[1].Secret)
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

			s.enableWalletAPI = !tc.disableWalletAPI

			var action func(*Wallet) error
			if tc.action != nil {
				action = tc.action(t)
			}

			err = s.UpdateSecrets(tc.viewWltName, tc.password, action)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			s.enableWalletAPI = true

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

func makeUxOut(t *testing.T, s cipher.SecKey, coins, hours uint64) coin.UxOut { // nolint: unparam
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
	p := cipher.MustPubKeyFromSecKey(s)
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
