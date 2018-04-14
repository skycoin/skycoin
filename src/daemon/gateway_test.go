package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/visor"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

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

func TestGateway_GetWalletDir(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		result          string
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			res, err := gw.GetWalletDir()
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}

func TestGateway_NewAddresses(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		walletID        string
		n               uint64
		result          []cipher.Address
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			res, err := gw.NewAddresses(tc.walletID, nil, tc.n)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}

func TestGateway_UpdateWalletLabel(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		walletID        string
		label           string
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			err := gw.UpdateWalletLabel(tc.walletID, tc.label)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
		})
	}
}

func TestGateway_GetWallet(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		walletID        string
		result          wallet.Wallet
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			res, err := gw.GetWallet(tc.walletID)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}

func TestGateway_GetWallets(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		wallets         wallet.Wallets
		getWalletError  error
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			w, err := gw.GetWallets()
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.wallets, w)
		})
	}
}

func TestGateway_GetWalletUnconfirmedTxns(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		walletID        string
		result          []visor.UnconfirmedTxn
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			res, err := gw.GetWalletUnconfirmedTxns(tc.walletID)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}

func TestGateway_ReloadWallets(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			err := gw.ReloadWallets()
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
		})
	}
}

func TestGateway_Spend(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		walletID        string
		coins           uint64
		dest            cipher.Address
		result          *coin.Transaction
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			res, err := gw.Spend(tc.walletID, nil, tc.coins, tc.dest)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}

func TestGateway_CreateWallet(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		wltName         string
		options         wallet.Options
		result          wallet.Wallet
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			res, err := gw.CreateWallet(tc.wltName, tc.options)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}

func TestGateway_GetWalletBalance(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		walletID        string
		result          wallet.BalancePair
		err             error
	}{
		{
			name:            "wallet api disabled",
			enableWalletAPI: false,
			err:             wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gw := &Gateway{
				Config: GatewayConfig{
					EnableWalletAPI: tc.enableWalletAPI,
				},
			}
			res, err := gw.GetWalletBalance(tc.walletID)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}
