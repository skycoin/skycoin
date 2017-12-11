package gui

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-querystring/query"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/strand"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/require"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

var (
	genPublic, genSecret         = cipher.GenerateKeyPair()
	genAddress                   = cipher.AddressFromPubKey(genPublic)
	testMaxSize                  = 1024 * 1024
	GenesisPublic, GenesisSecret = cipher.GenerateKeyPair()
	GenesisAddress               = cipher.AddressFromPubKey(GenesisPublic)
)

const (
	TimeIncrement    uint64 = 3600 * 1000
	GenesisTime      uint64 = 1000
	GenesisCoins     uint64 = 1000e6
	GenesisCoinHours uint64 = 1000 * 1000
)

type httpBody struct {
	Id    string `url:"id,omitempty"`
	Dst   string `url:"dst,omitempty"`
	Coins string `url:"coins,omitempty"`
}

// Gateway RPC interface wrapper for daemon state
type FakeGateway struct {
	vrpc visor.RPC

	t *testing.T
	// Backref to Daemon

	// Backref to Visor
	v *visor.Visor
	// Requests are queued on this channel
	requests chan strand.Request
}

// impelemts the wallet.Validator interface
type spendValidator struct {
	uncfm   *visor.UnconfirmedTxnPool
	unspent blockdb.UnspentPool
}

func MakeBlockchain(t *testing.T, db *bolt.DB, seckey cipher.SecKey) *visor.Blockchain {
	pubkey := cipher.PubKeyFromSecKey(seckey)
	b, err := visor.NewBlockchain(db, pubkey)
	require.NoError(t, err)
	gb, err := coin.NewGenesisBlock(GenesisAddress, GenesisCoins, GenesisTime)
	if err != nil {
		panic(fmt.Errorf("create genesis block failed: %v", err))
	}

	sig := cipher.SignHash(gb.HashHeader(), seckey)
	db.Update(func(tx *bolt.Tx) error {
		return b.ExecuteBlockWithTx(tx, &coin.SignedBlock{
			Block: *gb,
			Sig:   sig,
		})
	})
	return b
}

func (gw *FakeGateway) Spend(wltID string, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	var tx *coin.Transaction
	var err error
	//// create spend validator
	unspent := gw.v.Blockchain.Unspent()
	logger.Info("Spend unspent: %s", unspent)
	//sv := newSpendValidator(gw.v.Unconfirmed, unspent)
	//sv := &spendValidator{
	//	uncfm:   gw.v.Unconfirmed,
	//	unspent: unspent,
	//}
	tx = &coin.Transaction{}
	tx.PushOutput(genAddress, 11e6, 255)
	//tx, err = gw.vrpc.CreateAndSignTransaction(wltID, sv, unspent, gw.v.Blockchain.Time(), coins, dest)
	//if err != nil {
	//	err = fmt.Errorf("Create transaction failed: %v", err)
	//	return tx, err
	//}
	//// create and sign transaction
	//tx, err = gw.vrpc.CreateAndSignTransaction(wltID, sv, unspent, gw.v.Blockchain.Time(), coins, dest)
	//if err != nil {
	//	err = fmt.Errorf("Create transaction failed: %v", err)
	//	return
	//}
	//
	//// inject transaction
	//if err = gw.d.Visor.InjectTransaction(*tx, gw.d.Pool); err != nil {
	//	err = fmt.Errorf("Inject transaction failed: %v", err)
	//
	//})

	return tx, err
}

//func () UnspentPool() blockdb.UnspentPool {
//	return nil
//}

func NewVisorConfig() visor.Config {
	c := visor.Config{
		IsMaster: false,

		BlockchainPubkey: cipher.PubKey{},
		BlockchainSeckey: cipher.SecKey{},

		BlockCreationInterval: 10,
		//BlockCreationForceInterval: 120, //create block if no block within this many seconds

		UnconfirmedCheckInterval: time.Hour * 2,
		UnconfirmedMaxAge:        time.Hour * 48,
		UnconfirmedRefreshRate:   time.Minute,
		// UnconfirmedRefreshRate:   time.Minute * 30,
		UnconfirmedResendPeriod: time.Minute,
		MaxBlockSize:            1024 * 32,

		GenesisAddress:    cipher.Address{},
		GenesisSignature:  cipher.Sig{},
		GenesisTimestamp:  0,
		GenesisCoinVolume: 0, //100e12, 100e6 * 10e6
	}

	return c
}

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}

func makeUxOutWithSecret(t *testing.T) (coin.UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret(t)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func makeUxOut(t *testing.T) coin.UxOut {
	ux, _ := makeUxOutWithSecret(t)
	return ux
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *FakeGateway) GetWalletBalance(wltID string) (wallet.BalancePair, error) {
	var balance wallet.BalancePair
	var err error

	var addrs []cipher.Address
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(gw.t)
		uxs = append(uxs, ux)
	}
	addrs = []cipher.Address{uxs[1].Body.Address}
	//addrs, err = gw.vrpc.GetWalletAddresses(wltID)
	//if err != nil {
	//	return balance, err
	//}
	auxs := gw.vrpc.GetUnspent(gw.v).GetUnspentsOfAddrs(addrs)

	var spendUxs coin.AddressUxOuts
	spendUxs, err = gw.vrpc.GetUnconfirmedSpends(gw.v, addrs)
	if err != nil {
		err = fmt.Errorf("get unconfimed spending failed when checking wallet balance: %v", err)
		return balance, err
	}

	var recvUxs coin.AddressUxOuts
	recvUxs, err = gw.vrpc.GetUnconfirmedReceiving(gw.v, addrs)
	if err != nil {
		err = fmt.Errorf("get unconfirmed receiving failed when when checking wallet balance: %v", err)
		return balance, err
	}

	coins1, hours1 := gw.v.AddressBalance(auxs)
	coins2, hours2 := gw.v.AddressBalance(auxs.Sub(spendUxs).Add(recvUxs))
	balance = wallet.BalancePair{
		Confirmed: wallet.Balance{Coins: coins1, Hours: hours1},
		Predicted: wallet.Balance{Coins: coins2, Hours: hours2},
	}
	return balance, err
}

func TestWalletSpendHandler(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	cfg := NewVisorConfig()
	cfg.DBPath = db.Path()
	cfg.IsMaster = false
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress

	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)
	v := &visor.Visor{
		Config:      cfg,
		Unconfirmed: visor.NewUnconfirmedTxnPool(db),
		Blockchain:  bc,
		//db:          db,
	}
	// Setup blockchain

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	tt := []struct {
		name   string
		method string
		url    string
		body   *httpBody
		status int
		err    string
	}{
		{
			"405",
			"GET",
			"/wallet/spend",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
		},
		{
			"400 - no walletId",
			"POST",
			"/wallet/spend",
			&httpBody{},
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
		},
		{
			"400 - no dst",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id: "123",
			},
			http.StatusBadRequest,
			"400 Bad Request - missing destination address \"dst\"",
		},
		{
			"400 - bad dst addr",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:  "123",
				Dst: "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ======",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid destination address: Invalid address length",
		},
		{
			"400 - no coins",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:  "123",
				Dst: "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value",
		},
		{
			"400 - coins is string",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "foo",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value",
		},
		{
			"400 - coins is negative value",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "-123",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value",
		},
		{
			"400 - zero coins",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "0",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value, must > 0",
		},
		{
			"200 - OK",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "12",
			},
			http.StatusOK,
			"400 Bad Request - invalid \"coins\" value, must > 0",
		},
	}
	gateway := &FakeGateway{
		v: v,
		t: t,
	}
	for _, tc := range tt {
		body, _ := query.Values(tc.body)
		req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(body.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WalletSpendHandler(gateway))

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		status := rr.Code;
		if status != tc.status {
			t.Errorf("case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)
		}
		if status != http.StatusOK  {
			if errMsg := rr.Body.String(); strings.TrimSpace(errMsg) != tc.err {
				t.Errorf("case: %s, handler returned wrong error message: got `%v` want `%v`",
					tc.name, errMsg, tc.err)
			}
		}
	}

}
