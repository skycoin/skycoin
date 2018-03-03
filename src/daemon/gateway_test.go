package daemon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/daemon/strand"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/visor"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

func newGateway(disabledWalletApi bool) (*Gateway, *visor.RPCIfaceMock, *visor.VisorerMock, *visor.BlockchainerMock,
	*visor.UnconfirmedTxnPoolerMock, *VisorerMock) {
	vrpc := visor.NewRPCIfaceMock()
	v := visor.NewVisorerMock()
	bc := visor.NewBlockchainerMock()
	ucf := visor.NewUnconfirmedTxnPoolerMock()
	daemonVisor := NewVisorerMock()

	gw := &Gateway{
		Config: GatewayConfig{
			DisableWalletAPI: disabledWalletApi,
		},
		d: &Daemon{
			Visor: daemonVisor,
		},
		drpc:     RPC{},
		vrpc:     vrpc,
		v:        v,
		requests: make(chan strand.Request, 32),
		quit:     make(chan struct{}),
	}
	go func() {
		select {
		case req := <-gw.requests:
			req.Func()
		}
	}()
	return gw, vrpc, v, bc, ucf, daemonVisor
}

func TestFbyAddresses(t *testing.T) {
	uxs := make(coin.UxArray, 5)
	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxs[i] = coin.UxOut{
			Body: coin.UxBody{
				Address: addrs[i],
			},
		}
	}

	tests := []struct {
		name    string
		addrs   []string
		outputs []coin.UxOut
		want    []coin.UxOut
	}{
		// TODO: Add test cases.
		{
			"filter with one address",
			[]string{addrs[0].String()},
			uxs[:2],
			uxs[:1],
		},
		{
			"filter with multiple addresses",
			[]string{addrs[0].String(), addrs[1].String()},
			uxs[:3],
			uxs[:2],
		},
	}
	for _, tt := range tests {
		// fmt.Printf("want:%+v\n", tt.want)
		outs := FbyAddresses(tt.addrs)(tt.outputs)
		require.Equal(t, outs, coin.UxArray(tt.want))
	}
}

func TestFbyHashes(t *testing.T) {
	uxs := make(coin.UxArray, 5)
	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxs[i] = coin.UxOut{
			Body: coin.UxBody{
				Address: addrs[i],
			},
		}
	}

	type args struct {
		hashes []string
	}
	tests := []struct {
		name    string
		hashes  []string
		outputs coin.UxArray
		want    coin.UxArray
	}{
		// TODO: Add test cases.
		{
			"filter with one hash",
			[]string{uxs[0].Hash().Hex()},
			uxs[:2],
			uxs[:1],
		},
		{
			"filter with multiple hash",
			[]string{uxs[0].Hash().Hex(), uxs[1].Hash().Hex()},
			uxs[:3],
			uxs[:2],
		},
	}
	for _, tt := range tests {
		outs := FbyHashes(tt.hashes)(tt.outputs)
		require.Equal(t, outs, coin.UxArray(tt.want))
	}
}

func TestGateway_GetWallet(t *testing.T) {
	tests := []struct {
		name             string
		disableWalletAPI bool
		walletId         string
		wallet           wallet.Wallet
		getWalletError   error
		err              error
	}{
		{
			name:             "wallet api disabled",
			disableWalletAPI: true,
			err:              wallet.ErrWalletApiDisabled,
		},
		{
			name:             "getWalletError",
			disableWalletAPI: false,
			walletId:         "walletId",
			getWalletError:   errors.New("getWalletError"),
			err:              errors.New("getWalletError"),
		},
		{
			name:             "OK",
			disableWalletAPI: false,
			walletId:         "walletId",
			wallet:           wallet.Wallet{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw, vrpc, _, _, _, _ := newGateway(tc.disableWalletAPI)
			vrpc.On("GetWallet", tc.walletId).Return(tc.wallet, tc.getWalletError)
			w, err := gw.GetWallet(tc.walletId)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.wallet, w)
		})
	}
}

func TestGateway_GetWallets(t *testing.T) {
	tests := []struct {
		name             string
		disableWalletAPI bool
		wallets          wallet.Wallets
		getWalletError   error
		err              error
	}{
		{
			name:             "wallet api disabled",
			disableWalletAPI: true,
			err:              wallet.ErrWalletApiDisabled,
		},
		{
			name:             "OK",
			disableWalletAPI: false,
			wallets:          wallet.Wallets{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw, vrpc, _, _, _, _ := newGateway(tc.disableWalletAPI)
			vrpc.On("GetWallets").Return(tc.wallets)
			w, err := gw.GetWallets()
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.wallets, w)
		})
	}
}

func TestGateway_GetWalletUnconfirmedTxns(t *testing.T) {
	tests := []struct {
		name                     string
		disableWalletAPI         bool
		walletId                 string
		getWalletAddressesResult []cipher.Address
		getWalletError           error
		getUnconfirmedTxnsResult []visor.UnconfirmedTxn
		err                      error
	}{
		{
			name:             "wallet api disabled",
			disableWalletAPI: true,
			err:              wallet.ErrWalletApiDisabled,
		},
		{
			name:             "getWalletError",
			disableWalletAPI: false,
			walletId:         "walletId",
			getWalletError:   errors.New("getWalletError"),
		},
		{
			name:                     "OK",
			disableWalletAPI:         false,
			walletId:                 "walletId",
			getWalletError:           errors.New("getWalletError"),
			getUnconfirmedTxnsResult: []visor.UnconfirmedTxn{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw, vrpc, v, _, _, _ := newGateway(tc.disableWalletAPI)

			vrpc.On("GetWalletAddresses", tc.walletId).Return(tc.getWalletAddressesResult, tc.err)
			v.On("GetUnconfirmedTxns", mock.Anything).Return(tc.getUnconfirmedTxnsResult)
			res, err := gw.GetWalletUnconfirmedTxns(tc.walletId)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.getUnconfirmedTxnsResult, res)
		})
	}
}

func TestGateway_ReloadWallets(t *testing.T) {
	tests := []struct {
		name             string
		disableWalletAPI bool
		err              error
	}{
		{
			name:             "wallet api disabled",
			disableWalletAPI: true,
			err:              wallet.ErrWalletApiDisabled,
		},
		{
			name:             "reloadWalletError",
			disableWalletAPI: false,
			err:              errors.New("reloadWalletError"),
		},
		{
			name:             "OK",
			disableWalletAPI: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw, vrpc, _, _, _, _ := newGateway(tc.disableWalletAPI)

			vrpc.On("ReloadWallets").Return(tc.err)
			err := gw.ReloadWallets()
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
		})
	}
}

func TestGateway_Spend(t *testing.T) {
	tests := []struct {
		name                            string
		disableWalletAPI                bool
		walletId                        string
		coins                           uint64
		dest                            cipher.Address
		createAndSignTransactionResult  coin.Transaction
		createAndSignTransactionError   error
		injectBroadcastTransactionError error
		result                          *coin.Transaction
		err                             error
	}{
		{
			name:             "wallet api disabled",
			disableWalletAPI: true,
			err:              wallet.ErrWalletApiDisabled,
		},
		{
			name:                          "CreateAndSignTransaction error",
			disableWalletAPI:              false,
			createAndSignTransactionError: errors.New("createAndSignTransactionError"),
			err: errors.New("createAndSignTransactionError"),
		},
		{
			name:                            "OK",
			disableWalletAPI:                false,
			createAndSignTransactionResult:  coin.Transaction{},
			injectBroadcastTransactionError: errors.New("injectBroadcastTransactionError"),
			err: errors.New("injectBroadcastTransactionError"),
		},
		{
			name:                           "OK",
			disableWalletAPI:               false,
			createAndSignTransactionResult: coin.Transaction{},
			result: &coin.Transaction{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw, vrpc, v, bc, ucf, daemonVisor := newGateway(tc.disableWalletAPI)
			v.On("GetBlockchain").Return(bc)
			v.On("GetUnconfirmed").Return(ucf)
			bc.On("Unspent").Return(&blockdb.Unspents{})
			bc.On("Time").Return(uint64(123))
			vrpc.On("ReloadWallets").Return(tc.err)
			vrpc.On("CreateAndSignTransaction", tc.walletId, newSpendValidator(ucf, &blockdb.Unspents{}),
				bc.Unspent(), gw.v.GetBlockchain().Time(), tc.coins, tc.dest).
				Return(&tc.createAndSignTransactionResult, tc.createAndSignTransactionError)
			daemonVisor.On("InjectBroadcastTransaction", tc.createAndSignTransactionResult, gw.d.Pool).Return(tc.err)
			res, err := gw.Spend(tc.walletId, tc.coins, tc.dest)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}
