package wallet

import (
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

type mockBalanceGetter map[cipher.Address]BalancePair

// 10 addresses of seed1
var addrsOfSeed1 = []string{
	"2GBifzJEehbDX7Mkk63Prfa4MQQQyRzBLfe",
	"q2kU13X8XsAg8cS8BuSeSVzjPF9AT9ghAa",
	"2WXvTagXtrc1Qq71yjNXw86TC6SRgfVRH1B",
	"2NUNw748b9mT2FHRxgJL5KjBHasLfdP32Sh",
	"2V1CnVzWoXDaCX6wHU4tLJkWaFmLcQBb2q4",
	"wBkMr936thcr57wxyrH6ffvA99JN2Q1MN1",
	"2f92Wht7VQefAyoJUz3SEnfwT6wTdeAcq3L",
	"27UM5jPFYVuve3ceEHAYGaJSmkynQYmwPcH",
	"xjWbVN7ihReasVFwXJSSYYWF7rgQa22auC",
	"2LyanokLYFeBfBsNkRYHp2qtN8naGFJqeUw",
}

var childSeedsOfSeed1 = []string{
	"22b826c586039f8078433be26618ca1024e883d97de2267313bb78068f634c5a",
	"68efbbdf8aa06368cfc55e252d1e782bbd7651e590ee59e94ab579d2e44c20ad",
	"8894c818732375680284be4509d153272726f42296b85ecac1fb66b9dc7484b9",
	"6603375ee19c1e9fffe369e3f62e9deaa6931c1183d7da7f24ecbbd591061502",
	"91a63f939149d423ea39701d8ed16cfb16a3554c184d214d2289018ddb9e73de",
	"f0f4f008aa3e7cd32ee953507856fb46e37b734fd289dc01449133d7e37a1f07",
	"6b194da58a5ba5660cf2b00076cf6a2962fe8fe0523abca5647c87df3352866a",
	"b47a2678f7e797d3ada86e7e36855f572a18ab78dcbe54ed0613bba69fd76f8d",
	"fe064533108dadbef13be3a95f547ba03423aa6a701c40aaaed775cb783b12b3",
	"d554da211321a437e4d08f2a57e3ef255cffa89dd182e0fd52a4fd5bdfcab1ae",
}

func fromAddrString(t *testing.T, addrStrs []string) []cipher.Address {
	addrs := make([]cipher.Address, 0, len(addrStrs))
	for _, addr := range addrStrs {
		a, err := cipher.DecodeBase58Address(addr)
		require.NoError(t, err)
		addrs = append(addrs, a)
	}
	return addrs
}

func (mb mockBalanceGetter) GetBalanceOfAddrs(addrs []cipher.Address) ([]BalancePair, error) {
	var bals []BalancePair
	for _, addr := range addrs {
		bal := mb[addr]
		bals = append(bals, bal)
	}
	return bals, nil
}

func TestNewServiceDisabledWalletAPI(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir, true)
	require.NoError(t, err)
	dirIsEmpty(t, dir)

	// check if the wallet dir is created
	_, err = os.Stat(dir)
	require.NoError(t, err)

	require.Equal(t, "", s.WalletDirectory)

	// check if the default wallet is created
	require.Equal(t, 0, len(s.wallets))
}

func TestNewService(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir, false)
	require.NoError(t, err)

	// check if the wallet dir is created
	_, err = os.Stat(dir)
	require.NoError(t, err)

	require.Equal(t, dir, s.WalletDirectory)

	// check if the default wallet is created
	require.Equal(t, 1, len(s.wallets))

	// check if the default wallet file is created
	for name := range s.wallets {
		wltFile := filepath.Join(dir, name)
		_, err := os.Stat(wltFile)
		require.NoError(t, err)
		break
	}

	// test load wallets
	s, err = NewService("./testdata", false)
	require.NoError(t, err)

	// check if the dup wallet is loaded
	_, ok1 := s.wallets["test3.1.wlt"]
	_, ok2 := s.wallets["test3.wlt"]
	if ok1 && ok2 {
		t.Fatal("load dup wallet")
	}

	require.Equal(t, 3, len(s.wallets))
}

func TestServiceCreateWalletDisabledWalletAPI(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir, true)
	require.NoError(t, err)
	dirIsEmpty(t, dir)

	wltName := "t1.wlt"
	seed := "seed1"
	_, err = s.CreateWallet(wltName, Options{
		Seed: seed,
	})
	dirIsEmpty(t, dir)
	require.Equal(t, ErrWalletApiDisabled, err)
}

func TestServiceCreateWallet(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir, false)
	require.NoError(t, err)

	wltName := "t1.wlt"
	seed := "seed1"
	w, err := s.CreateWallet(wltName, Options{
		Seed: seed,
	})
	require.NoError(t, err)
	require.Equal(t, seed, w.Meta["seed"])
	require.NoError(t, w.Validate())

	// create wallet with dup wallet name
	_, err = s.CreateWallet(wltName, Options{Seed: "seed2"})
	require.Equal(t, err, ErrWalletNameConflict)

	// create wallet with dup seed
	dupWlt := "dup_wallet.wlt"
	_, err = s.CreateWallet(dupWlt, Options{
		Seed: seed,
	})
	require.EqualError(t, err, fmt.Sprintf("duplicate wallet with %v", wltName))

	// check if the dup wallet is created
	_, ok := s.wallets[dupWlt]
	require.False(t, ok)

	_, err = os.Stat(filepath.Join(dir, dupWlt))
	require.True(t, os.IsNotExist(err))
}

func TestServiceScanWalletDisabledWalletAPI(t *testing.T) {
	bg := make(mockBalanceGetter, len(addrsOfSeed1))
	dir := prepareWltDir()
	s, err := NewService(dir, true)
	require.NoError(t, err)
	dirIsEmpty(t, dir)

	wallet, err := s.ScanAheadWalletAddresses("wltName", 1, bg)
	require.Equal(t, err, errors.New("wallet doesn't exist"))
	require.Equal(t, wallet, Wallet{})
}

func TestServiceCreateAndScanWallet(t *testing.T) {
	bg := make(mockBalanceGetter, len(addrsOfSeed1))
	addrs := fromAddrString(t, addrsOfSeed1)
	for _, a := range addrs {
		bg[a] = BalancePair{}
	}

	type exp struct {
		seed             string
		lastSeed         string
		entryNum         int
		confirmedBalance uint64
		predictedBalance uint64
	}

	tt := []struct {
		name      string
		scanN     uint64
		balGetter BalanceGetter
		expect    exp
	}{
		{
			"no coins and scan 0",
			0,
			bg,
			exp{
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"no coins and scan 1",
			1,
			bg,
			exp{
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"no coins and scan 10",
			10,
			bg,
			exp{
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[0],
				entryNum:         1,
				confirmedBalance: 0,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 5",
			5,
			mockBalanceGetter{
				addrs[5]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[5],
				entryNum:         5 + 1,
				confirmedBalance: 10,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 4",
			5,
			mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
		{
			"scan 5 get 4 have 6",
			5,
			mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[6]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},

		{
			"confirmed and predicted",
			5,
			mockBalanceGetter{
				addrs[3]: BalancePair{Confirmed: Balance{Coins: 10, Hours: 100}},
				addrs[4]: BalancePair{Predicted: Balance{Coins: 10, Hours: 100}},
			},
			exp{
				seed:             "seed1",
				lastSeed:         childSeedsOfSeed1[4],
				entryNum:         4 + 1,
				confirmedBalance: 20,
				predictedBalance: 0,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()

			s, err := NewService(dir, false)
			require.NoError(t, err)

			wltName := "t1.wlt"
			seed := "seed1"
			w, err := s.CreateWallet(wltName, Options{
				Seed:  seed,
				Label: "foo",
			})
			require.NoError(t, err)

			require.Equal(t, seed, w.Meta["seed"])
			require.NoError(t, w.Validate())

			w, err = s.ScanAheadWalletAddresses(wltName, tc.scanN, tc.balGetter)
			require.NoError(t, err)

			require.Len(t, w.Entries, tc.expect.entryNum)
			require.Equal(t, tc.expect.lastSeed, w.getLastSeed())
		})
	}
}

func TestServiceNewAddressDisabledWalletAPI(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir, true)
	require.NoError(t, err)
	dirIsEmpty(t, dir)

	require.Empty(t, s.wallets)
	addrs, err := s.NewAddresses("", 1)
	require.Equal(t, ErrWalletNotExist, err)
	require.Equal(t, 0, len(addrs))
}

func TestServiceNewAddress(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir, false)
	require.NoError(t, err)

	// get the default wallet id
	var id string
	for id = range s.wallets {
		break
	}
	addrs, err := s.NewAddresses(id, 1)
	require.NoError(t, err)
	require.Equal(t, 1, len(addrs))

	// check if the wallet file is created
	_, err = os.Stat(filepath.Join(dir, id))
	require.NoError(t, err)

	// wallet doesn't exist
	_, err = s.NewAddresses("not_exist_id.wlt", 1)
	require.Equal(t, ErrWalletNotExist, err)
}

func TestServiceGetAddressDisabledWalletAPI(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir, true)
	require.NoError(t, err)
	dirIsEmpty(t, dir)

	require.Empty(t, s.wallets)
	addrs, err := s.GetAddresses("")
	require.Equal(t, ErrWalletNotExist, err)
	require.Equal(t, 0, len(addrs))
}

func TestServiceGetAddress(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir, false)
	require.NoError(t, err)

	var id string
	for id = range s.wallets {
		break
	}

	addrs, err := s.GetAddresses(id)
	require.NoError(t, err)
	require.Equal(t, 1, len(addrs))

	// test none exist wallet
	notExistID := "not_exist_id.wlt"
	_, err = s.GetAddresses(notExistID)
	require.Equal(t, ErrWalletNotExist, err)
}

func TestServiceGetWalletDisabledWalletAPI(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir, true)
	require.NoError(t, err)
	dirIsEmpty(t, dir)

	require.Empty(t, s.wallets)
	w, err := s.GetWallet("")
	require.Equal(t, ErrWalletNotExist, err)
	require.Equal(t, w, Wallet{})
}

func TestServiceGetWallet(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir, false)
	require.NoError(t, err)

	var id string
	for id = range s.wallets {
		break
	}

	w, err := s.GetWallet(id)
	require.NoError(t, err)

	// modify the returned wallet won't affect the wallet in service
	w.SetLabel("new_label")

	w1, err := s.GetWallet(id)
	require.NoError(t, err)

	require.NotEqual(t, "new_label", w1.GetLabel())
}

func TestServiceReloadWalletsDisabledWalletAPI(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir, true)
	require.NoError(t, err)
	dirIsEmpty(t, dir)

	err = s.ReloadWallets()
	require.Equal(t, ErrWalletApiDisabled, err)
}

func TestServiceReloadWallets(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir, false)
	require.NoError(t, err)

	var defaultWltID string
	for defaultWltID = range s.wallets {
		break
	}

	var defaultAddr string
	for defaultAddr = range s.firstAddrIDMap {
		break
	}

	wltName := "t1.wlt"
	w, err := s.CreateWallet(wltName, Options{Seed: "seed1"})
	require.NoError(t, err)

	err = s.ReloadWallets()
	require.NoError(t, err)

	// check if create dup wallet will return error
	_, ok := s.wallets[defaultWltID]
	require.True(t, ok)

	_, ok = s.wallets[wltName]
	require.True(t, ok)

	// check if the first address of each wallet is reloaded
	_, ok = s.firstAddrIDMap[defaultAddr]
	require.True(t, ok)

	_, ok = s.firstAddrIDMap[w.Entries[0].Address.String()]
	require.True(t, ok)
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

func TestServiceCreateAndSignTx(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir, false)
	require.NoError(t, err)
	var id string
	for id = range s.wallets {
		break
	}

	headTime := time.Now().UTC().Unix()

	wlt, err := s.GetWallet(id)
	require.NoError(t, err)

	secKey := wlt.Entries[0].Secret
	addr := wlt.Entries[0].Address

	var uxouts []coin.UxOut
	addrs := []cipher.Address{}
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey)
		uxouts = append(uxouts, uxout)

		p, _ := cipher.GenerateKeyPair()
		a := cipher.AddressFromPubKey(p)
		addrs = append(addrs, a)
	}

	var uxoutsNoHours []coin.UxOut
	addrsNoHours := []cipher.Address{}
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey)
		uxout.Body.Hours = 0
		uxout.Head.Time = uint64(headTime)
		uxoutsNoHours = append(uxoutsNoHours, uxout)

		p, _ := cipher.GenerateKeyPair()
		a := cipher.AddressFromPubKey(p)
		addrsNoHours = append(addrsNoHours, a)
	}

	tt := []struct {
		name       string
		unspents   []coin.UxOut
		addrUxouts coin.AddressUxOuts
		vld        Validator
		coins      uint64
		dest       cipher.Address
		err        error
	}{
		{
			"ok with no change",
			uxouts[:],
			coin.AddressUxOuts{
				addr: uxouts,
			},
			&dummyValidator{
				ok: false,
			},
			2e6,
			addrs[0],
			nil,
		},
		{
			"ok with change",
			uxouts[:],
			coin.AddressUxOuts{
				addr: uxouts,
			},
			&dummyValidator{
				ok: false,
			},
			1e6,
			addrs[0],
			nil,
		},
		{
			"has unconfirmed spending transaction",
			uxouts[:],
			coin.AddressUxOuts{
				addr: uxouts,
			},
			&dummyValidator{
				ok: true,
			},
			2e6,
			addrs[0],
			errors.New("please spend after your pending transaction is confirmed"),
		},
		{
			"check unconfirmed spend failed",
			uxouts[:],
			coin.AddressUxOuts{
				addr: uxouts,
			},
			&dummyValidator{
				ok:  false,
				err: errors.New("fail intentionally"),
			},
			2e6,
			addrs[0],
			errors.New("checking unconfirmed spending failed: fail intentionally"),
		},
		{
			"spend zero",
			uxouts[:],
			coin.AddressUxOuts{
				addr: uxouts,
			},
			&dummyValidator{
				ok: false,
			},
			0,
			addrs[0],
			errors.New("zero spend amount"),
		},
		{
			"spend fractional coins",
			uxouts[:],
			coin.AddressUxOuts{
				addr: uxouts,
			},
			&dummyValidator{
				ok: false,
			},
			1e3,
			addrs[0],
			nil,
		},
		{
			"not enough confirmed coins",
			uxouts[:],
			coin.AddressUxOuts{
				addr: uxouts,
			},
			&dummyValidator{
				ok: false,
			},
			100e6,
			addrs[0],
			ErrInsufficientBalance,
		},
		{
			"no coin hours in inputs",
			uxoutsNoHours[:],
			coin.AddressUxOuts{
				addr: uxoutsNoHours,
			},
			&dummyValidator{
				ok: false,
			},
			1e6,
			addrsNoHours[0],
			fee.ErrTxnNoFee,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			unspents := &dummyUnspentGetter{
				addrUnspents: tc.addrUxouts,
				unspents:     map[cipher.SHA256]coin.UxOut{},
			}

			for _, ux := range tc.unspents {
				unspents.unspents[ux.Hash()] = ux
			}

			tx, err := s.CreateAndSignTransaction(id, tc.vld, unspents, uint64(headTime), tc.coins, tc.dest)
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

func makeUxBody(t *testing.T, s cipher.SecKey) coin.UxBody {
	p := cipher.PubKeyFromSecKey(s)
	return coin.UxBody{
		SrcTransaction: cipher.SumSHA256(testutil.RandBytes(t, 128)),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          2e6,
		Hours:          100,
	}
}

func makeUxOut(t *testing.T, s cipher.SecKey) coin.UxOut {
	body := makeUxBody(t, s)
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
