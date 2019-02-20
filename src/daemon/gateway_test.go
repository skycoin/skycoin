package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/visor"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

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
		result          []visor.UnconfirmedTransaction
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
			res, err := gw.GetWalletUnconfirmedTransactions(tc.walletID)
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
			walletID:        "foo.wlt",
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
			res, _, err := gw.GetWalletBalance(tc.walletID)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.Equal(t, tc.result, res)
		})
	}
}

func TestGateway_WalletCreateTransaction(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		err             error
		wltID           string
		params          transaction.Params
		visorParams     visor.CreateTransactionParams
		txn             *coin.Transaction
		inputs          []visor.TransactionInput
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

			txn, inputs, err := gw.WalletCreateTransaction(tc.wltID, tc.params, tc.visorParams)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.txn, txn)
				require.Equal(t, tc.inputs, inputs)
			}
		})
	}
}

func TestGateway_WalletCreateTransactionSigned(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		err             error
		wltID           string
		password        []byte
		params          transaction.Params
		visorParams     visor.CreateTransactionParams
		txn             *coin.Transaction
		inputs          []visor.TransactionInput
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

			txn, inputs, err := gw.WalletCreateTransactionSigned(tc.wltID, tc.password, tc.params, tc.visorParams)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.txn, txn)
				require.Equal(t, tc.inputs, inputs)
			}
		})
	}
}

func TestGateway_WalletSignTransaction(t *testing.T) {
	tests := []struct {
		name            string
		enableWalletAPI bool
		txn             *coin.Transaction
		wltID           string
		password        []byte
		signIndexes     []int
		err             error
		outTxn          *coin.Transaction
		inputs          []visor.TransactionInput
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

			txn, inputs, err := gw.WalletSignTransaction(tc.wltID, tc.password, tc.txn, tc.signIndexes)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.inputs, inputs)
				require.Equal(t, tc.outTxn, txn)
			}
		})
	}
}
