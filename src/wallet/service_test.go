package wallet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareWltDir() string {
	dir, err := ioutil.TempDir("", "wallets")
	if err != nil {
		panic(err)
	}

	return dir
}

func TestNewService(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir)
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
	s, err = NewService("./testdata")
	require.NoError(t, err)

	// check if the dup wallet is loaded
	_, ok1 := s.wallets["test3.1.wlt"]
	_, ok2 := s.wallets["test3.wlt"]
	if ok1 && ok2 {
		t.Fatal("load dup wallet")
	}

	require.Equal(t, 3, len(s.wallets))
}

func TestServiceCreateWallet(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir)
	require.NoError(t, err)

	wltName := "t1.wlt"
	seed := "seed1"
	w, err := s.CreateWallet(wltName, OptSeed(seed))
	require.NoError(t, err)
	require.Equal(t, seed, w.Meta["seed"])
	require.NoError(t, w.Validate())

	// create walelt with dup wallet name
	_, err = s.CreateWallet(wltName)
	require.EqualError(t, err, "wallet name would conflict with existing wallet, renaming")

	// create wallet with dup seed
	dupWlt := "dup_wallet.wlt"
	_, err = s.CreateWallet(dupWlt, OptSeed(seed))
	require.EqualError(t, err, fmt.Sprintf("duplicate wallet with %v", wltName))

	// check if the dup wallet is created
	_, ok := s.wallets[dupWlt]
	require.False(t, ok)

	_, err = os.Stat(filepath.Join(dir, dupWlt))
	require.True(t, os.IsNotExist(err))
}

func TestServiceNewAddress(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir)
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
	require.Equal(t, errWalletNotExist("not_exist_id.wlt"), err)
}

func TestServiceGetAddress(t *testing.T) {
	dir := prepareWltDir()
	s, err := NewService(dir)
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
	require.Equal(t, errWalletNotExist(notExistID), err)
}

func TestServiceGetWallet(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir)
	require.NoError(t, err)

	var id string
	for id = range s.wallets {
		break
	}

	w, ok := s.GetWallet(id)
	require.True(t, ok)

	// modify the returned wallet won't affect the wallet in service
	w.SetLabel("new_label")

	w1, ok := s.GetWallet(id)
	require.True(t, ok)

	require.NotEqual(t, "new_label", w1.GetLabel())
}

func TestServiceReloadWallets(t *testing.T) {
	dir := prepareWltDir()

	s, err := NewService(dir)
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
	w, err := s.CreateWallet(wltName)
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

	s, err := NewService(dir)
	require.NoError(t, err)
	var id string
	for id = range s.wallets {
		break
	}

	headTime := time.Now().UTC().Unix()

	wlt, ok := s.GetWallet(id)
	require.True(t, ok)
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

	tt := []struct {
		name       string
		unspents   []coin.UxOut
		addrUxouts coin.AddressUxOuts
		vld        Validator
		amt        Balance
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
			Balance{Coins: 2e6},
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
			Balance{Coins: 1e6},
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
			Balance{Coins: 2e6},
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
			Balance{Coins: 2e6},
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
			Balance{},
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
			Balance{Coins: 1000},
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
			Balance{Coins: 100e6},
			addrs[0],
			errors.New("not enough confirmed coins"),
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

			tx, err := s.CreateAndSignTransaction(id, tc.vld, unspents, uint64(headTime), tc.amt, tc.dest)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			// check the IN of tx
			for _, inUxid := range tx.In {
				_, ok := unspents.unspents[inUxid]
				require.True(t, ok)
			}

			require.NoError(t, tx.Verify())
		})
	}
}

func makeUxBody(t *testing.T, s cipher.SecKey) coin.UxBody {
	body, _ := makeUxBodyWithSecret(t, s)
	return body
}

func makeUxOut(t *testing.T, s cipher.SecKey) coin.UxOut {
	ux, _ := makeUxOutWithSecret(t, s)
	return ux
}

func makeUxBodyWithSecret(t *testing.T, s cipher.SecKey) (coin.UxBody, cipher.SecKey) {
	p := cipher.PubKeyFromSecKey(s)
	return coin.UxBody{
		SrcTransaction: cipher.SumSHA256(randBytes(t, 128)),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          2e6,
		Hours:          100,
	}, s
}

func makeUxOutWithSecret(t *testing.T, s cipher.SecKey) (coin.UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret(t, s)
	tm := rand.Int31n(1000)
	seq := rand.Int31n(100)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  uint64(tm),
			BkSeq: uint64(seq),
		},
		Body: body,
	}, sec
}

func randBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	x, err := rand.Read(b)
	assert.Equal(t, n, x) //end unit testing.
	assert.Nil(t, err)
	return b
}
